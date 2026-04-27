package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"server/global"
	"server/model"
)

// 获取可用模型列表（优先从数据库读取）
func (s *AIService) GetModels() []ModelInfo {
	var models []ModelInfo
	aiConfig := Config.GetAIConfig()
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

	return models
}

// 获取对话列表
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

// 获取对话消息
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

// 创建新对话
func (s *AIService) CreateConversation(userID uint, input CreateConversationInput) (*model.AIConversation, error) {
	modelName := input.Model
	if modelName == "" {
		modelName = s.defaultModelID()
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

// 删除对话
func (s *AIService) DeleteConversation(conversationID uint, userID uint) error {
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("对话不存在")
	}

	if err := global.DB.Where("conversation_id = ?", conversationID).Delete(&model.AIMessage{}).Error; err != nil {
		return err
	}

	return global.DB.Delete(&conversation).Error
}

// 普通对话（非流式）
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
	if err := global.DB.Create(userMessage).Error; err != nil {
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

	assistantMessage := &model.AIMessage{
		ConversationID:   conversation.ID,
		Role:             "assistant",
		Content:          chatResp.Choices[0].Message.Content,
		ReasoningContent: chatResp.Choices[0].Message.ReasoningContent,
	}
	if err := global.DB.Create(assistantMessage).Error; err != nil {
		return nil, err
	}

	global.DB.Model(conversation).Update("updated_at", assistantMessage.CreatedAt)
	return assistantMessage, nil
}

// 流式对话 - 返回reader供SSE使用
func (s *AIService) ChatStream(userID uint, input AIChatInput) (
	conversationID uint,
	reader io.ReadCloser,
	saveConversation bool,
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
				return 0, nil, false, errors.New("对话不存在")
			}
		} else {
			conversation, err = s.CreateConversation(userID, CreateConversationInput{Model: input.Model})
			if err != nil {
				return 0, nil, false, err
			}
		}

		messages, err = s.buildContextMessages(conversation.ID, input.Message, input.FileIDs)
		if err != nil {
			return 0, nil, false, err
		}

		userMessage := &model.AIMessage{
			ConversationID: conversation.ID,
			Role:           "user",
			Content:        input.Message,
			FileIDs:        s.encodeFileIDs(input.FileIDs),
		}
		if err := global.DB.Create(userMessage).Error; err != nil {
			return 0, nil, false, err
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
			modelName = s.defaultModelID()
		}
	}

	reader, err = s.callStreamAPI(modelName, messages, input.EnableSearch, input.EnableThinking)
	if err != nil {
		return 0, nil, saveConversation, err
	}

	return conversationID, reader, saveConversation, nil
}

// 保存流式对话的助手消息
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

// 更新对话标题
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

// 清空对话消息
func (s *AIService) ClearMessages(conversationID uint, userID uint) error {
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("对话不存在")
	}

	return global.DB.Where("conversation_id = ?", conversationID).Delete(&model.AIMessage{}).Error
}

// 清空上下文（保留聊天记录，但新消息不带历史上下文）
func (s *AIService) ClearContext(conversationID uint, userID uint) error {
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("对话不存在")
	}

	now := time.Now()
	return global.DB.Model(&conversation).Update("context_cleared_at", now).Error
}

// 删除单条消息
func (s *AIService) DeleteMessage(messageID uint, userID uint) error {
	var message model.AIMessage
	if err := global.DB.First(&message, messageID).Error; err != nil {
		return errors.New("消息不存在")
	}

	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", message.ConversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("无权限删除此消息")
	}

	return global.DB.Delete(&message).Error
}

// 获取单个对话
func (s *AIService) GetConversation(conversationID uint, userID uint) (*model.AIConversation, error) {
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return nil, errors.New("对话不存在")
	}
	return &conversation, nil
}
