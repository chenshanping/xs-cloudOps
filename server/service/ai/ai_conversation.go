package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/service/core"
	"server/service/filesvc"
)

func (s *AIService) GetModels() ([]ModelInfo, error) {
	var models []ModelInfo
	aiConfig, err := s.GetAdminConfig()
	if err != nil {
		return nil, err
	}
	defaultProvider := aiConfig.DefaultProvider
	for _, provider := range aiConfig.Providers {
		for _, m := range provider.Models {
			if defaultProvider == provider.Name {
				models = append(models, ModelInfo{
					ID:          m.ID,
					Name:        m.Name,
					Description: m.Description,
				})
			}
		}
	}

	return models, nil
}

func (s *AIService) GetConversations(userID uint, input ConversationListInput) ([]model.AIConversation, int64, error) {
	var conversations []model.AIConversation
	var total int64

	paging := input.Normalize()
	db := global.DB.Model(&model.AIConversation{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(paging.Offset()).Limit(paging.PageSize).Order("updated_at DESC").Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

func (s *AIService) GetMessages(conversationID uint, userID uint) ([]model.AIMessage, error) {
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return nil, errors.New("对话不存在")
	}

	var messages []model.AIMessage
	if err := global.DB.Where("conversation_id = ?", conversationID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *AIService) CreateConversation(userID uint, input CreateConversationInput) (*model.AIConversation, error) {
	modelName := input.Model
	if modelName == "" {
		defaultModelName, err := s.defaultModelID()
		if err != nil {
			return nil, err
		}
		modelName = defaultModelName
	}

	title := input.Title
	if title == "" {
		title = "新对话"
	}

	conversation := &model.AIConversation{
		UserID: userID,
		Title:  title,
		Model:  modelName,
	}

	if err := global.DB.Create(conversation).Error; err != nil {
		return nil, err
	}

	return conversation, nil
}

func (s *AIService) DeleteConversation(conversationID uint, userID uint) error {
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("对话不存在")
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.clearConversationMessageFileRefs(tx, conversationID); err != nil {
			return err
		}
		if err := tx.Where("conversation_id = ?", conversationID).Delete(&model.AIMessage{}).Error; err != nil {
			return err
		}
		return tx.Delete(&conversation).Error
	})
}

func (s *AIService) BatchDeleteConversations(userID uint, ids []uint) (int, []string) {
	normalized := core.NormalizeIDs(ids)
	if len(normalized) == 0 {
		return 0, []string{"请选择要删除的对话"}
	}

	var successCount int
	failedMsgs := make([]string, 0)

	for _, id := range normalized {
		if err := s.DeleteConversation(id, userID); err != nil {
			failedMsgs = append(failedMsgs, fmt.Sprintf("ID %d: %s", id, err.Error()))
			continue
		}
		successCount++
	}

	return successCount, failedMsgs
}

func (s *AIService) Chat(userID uint, input AIChatInput) (*model.AIMessage, error) {
	var conversation *model.AIConversation
	var err error

	if input.ConversationID > 0 {
		conversation = &model.AIConversation{}
		if err := global.DB.Where("id = ? AND user_id = ?", input.ConversationID, userID).First(conversation).Error; err != nil {
			return nil, errors.New("对话不存在")
		}
	} else {
		conversation, err = s.CreateConversation(userID, CreateConversationInput{Model: input.Model})
		if err != nil {
			return nil, err
		}
	}

	messages, err := s.buildContextMessages(conversation.ID, input.Message, input.FileIDs)
	if err != nil {
		return nil, err
	}

	userMessage := &model.AIMessage{
		ConversationID: conversation.ID,
		Role:           "user",
		Content:        input.Message,
		FileIDs:        s.encodeFileIDs(input.FileIDs),
	}
	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(userMessage).Error; err != nil {
			return err
		}
		return s.registerMessageFileRefs(tx, userMessage.ID, input.FileIDs)
	}); err != nil {
		return nil, err
	}

	if conversation.Title == "新对话" {
		title := input.Message
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		global.DB.Model(conversation).Update("title", title)
	}

	modelName := conversation.Model
	if input.Model != "" {
		modelName = input.Model
	}

	messages, sourcesMarkdown, err := s.prepareSearchGrounding(context.Background(), messages, input.Message, input.EnableSearch)
	if err != nil {
		return nil, err
	}
	response, err := s.callAPI(modelName, messages, false, input.EnableSearch, input.EnableThinking)
	if err != nil {
		return nil, err
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(response, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	if len(chatResp.Choices) == 0 {
		return nil, errors.New("API返回空响应")
	}
	assistantContent := appendSourcesToContent(chatResp.Choices[0].Message.Content, sourcesMarkdown)

	assistantMessage := &model.AIMessage{
		ConversationID:   conversation.ID,
		Role:             "assistant",
		Content:          assistantContent,
		ReasoningContent: chatResp.Choices[0].Message.ReasoningContent,
	}
	if err := global.DB.Create(assistantMessage).Error; err != nil {
		return nil, err
	}

	global.DB.Model(conversation).Update("updated_at", assistantMessage.CreatedAt)
	return assistantMessage, nil
}

func (s *AIService) ChatStream(userID uint, input AIChatInput) (
	conversationID uint,
	reader io.ReadCloser,
	saveConversation bool,
	sourcesMarkdown string,
	err error,
) {
	saveConversation = input.SaveConversation == nil || *input.SaveConversation

	var messages []ChatMessage
	var modelName string

	if saveConversation {
		var conversation *model.AIConversation

		if input.ConversationID > 0 {
			conversation = &model.AIConversation{}
			if err := global.DB.Where("id = ? AND user_id = ?", input.ConversationID, userID).First(conversation).Error; err != nil {
				return 0, nil, false, "", errors.New("对话不存在")
			}
		} else {
			conversation, err = s.CreateConversation(userID, CreateConversationInput{Model: input.Model})
			if err != nil {
				return 0, nil, false, "", err
			}
		}

		messages, err = s.buildContextMessages(conversation.ID, input.Message, input.FileIDs)
		if err != nil {
			return 0, nil, false, "", err
		}

		userMessage := &model.AIMessage{
			ConversationID: conversation.ID,
			Role:           "user",
			Content:        input.Message,
			FileIDs:        s.encodeFileIDs(input.FileIDs),
		}
		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(userMessage).Error; err != nil {
				return err
			}
			return s.registerMessageFileRefs(tx, userMessage.ID, input.FileIDs)
		}); err != nil {
			return 0, nil, false, "", err
		}

		if conversation.Title == "新对话" {
			title := input.Message
			if len(title) > 50 {
				title = title[:50] + "..."
			}
			global.DB.Model(conversation).Update("title", title)
		}

		conversationID = conversation.ID
		modelName = conversation.Model
		if input.Model != "" {
			modelName = input.Model
		}
	} else {
		messages = []ChatMessage{{Role: "user", Content: input.Message}}
		modelName = input.Model
		if modelName == "" {
			modelName, err = s.defaultModelID()
			if err != nil {
				return 0, nil, saveConversation, "", err
			}
		}
	}

	messages, sourcesMarkdown, err = s.prepareSearchGrounding(context.Background(), messages, input.Message, input.EnableSearch)
	if err != nil {
		return 0, nil, saveConversation, "", err
	}

	reader, err = s.callStreamAPI(modelName, messages, input.EnableSearch, input.EnableThinking)
	if err != nil {
		return 0, nil, saveConversation, "", err
	}

	return conversationID, reader, saveConversation, sourcesMarkdown, nil
}

func (s *AIService) SaveAssistantMessage(conversationID uint, content, reasoningContent string) error {
	message := &model.AIMessage{
		ConversationID:   conversationID,
		Role:             "assistant",
		Content:          content,
		ReasoningContent: reasoningContent,
	}
	if err := global.DB.Create(message).Error; err != nil {
		return err
	}

	global.DB.Model(&model.AIConversation{}).Where("id = ?", conversationID).Update("updated_at", message.CreatedAt)
	return nil
}

func (s *AIService) UpdateConversationTitle(conversationID uint, userID uint, title string) error {
	result := global.DB.Model(&model.AIConversation{}).Where("id = ? AND user_id = ?", conversationID, userID).Update("title", title)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("对话不存在")
	}
	return nil
}

func (s *AIService) ClearMessages(conversationID uint, userID uint) error {
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("对话不存在")
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.clearConversationMessageFileRefs(tx, conversationID); err != nil {
			return err
		}
		return tx.Where("conversation_id = ?", conversationID).Delete(&model.AIMessage{}).Error
	})
}

func (s *AIService) ClearContext(conversationID uint, userID uint) error {
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("对话不存在")
	}

	now := time.Now()
	return global.DB.Model(&conversation).Update("context_cleared_at", now).Error
}

func (s *AIService) DeleteMessage(messageID uint, userID uint) error {
	var message model.AIMessage
	if err := global.DB.First(&message, messageID).Error; err != nil {
		return errors.New("消息不存在")
	}

	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", message.ConversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("无权限删除此消息")
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := filesvc.Reference.ClearRefs(tx, "ai_message", message.ID); err != nil {
			return err
		}
		return tx.Delete(&message).Error
	})
}

func (s *AIService) GetConversation(conversationID uint, userID uint) (*model.AIConversation, error) {
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return nil, errors.New("对话不存在")
	}
	return &conversation, nil
}

func (s *AIService) registerMessageFileRefs(tx *gorm.DB, messageID uint, fileIDs []uint) error {
	refs := make([]filesvc.FileRef, 0, len(fileIDs))
	for _, fileID := range fileIDs {
		refs = append(refs, filesvc.FileRef{
			FileID: fileID,
			Field:  "attachment",
		})
	}
	return filesvc.Reference.ReplaceRefs(tx, "ai_message", messageID, refs)
}

func (s *AIService) clearConversationMessageFileRefs(tx *gorm.DB, conversationID uint) error {
	var messageIDs []uint
	if err := tx.Model(&model.AIMessage{}).
		Where("conversation_id = ?", conversationID).
		Pluck("id", &messageIDs).Error; err != nil {
		return err
	}

	for _, messageID := range messageIDs {
		if err := filesvc.Reference.ClearRefs(tx, "ai_message", messageID); err != nil {
			return err
		}
	}
	return nil
}
