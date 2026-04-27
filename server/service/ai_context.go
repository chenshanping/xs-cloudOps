package service

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"server/global"
	"server/model"
)

// 构建上下文消息
func (s *AIService) buildContextMessages(conversationID uint, newMessage string, fileIDs []uint) ([]ChatMessage, error) {
	var messages []ChatMessage

	var history []model.AIMessage
	if conversationID > 0 {
		var conversation model.AIConversation
		if err := global.DB.First(&conversation, conversationID).Error; err != nil {
			return nil, err
		}

		query := global.DB.Where("conversation_id = ?", conversationID)
		if conversation.ContextClearedAt != nil {
			query = query.Where("created_at > ?", conversation.ContextClearedAt)
		}

		if err := query.Order("created_at ASC").Limit(20).Find(&history).Error; err != nil {
			return nil, err
		}
	}

	for _, msg := range history {
		messages = append(messages, ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	if len(fileIDs) > 0 {
		content, err := s.buildMessageWithFiles(newMessage, fileIDs)
		if err != nil {
			messages = append(messages, ChatMessage{Role: "user", Content: newMessage})
		} else {
			messages = append(messages, ChatMessage{Role: "user", Content: content})
		}
	} else {
		messages = append(messages, ChatMessage{Role: "user", Content: newMessage})
	}

	return messages, nil
}

// 构建带文件的消息内容
func (s *AIService) buildMessageWithFiles(message string, fileIDs []uint) (interface{}, error) {
	var files []model.SysFile
	if err := global.DB.Preload("Storage").Where("id IN ? AND status = 1", fileIDs).Find(&files).Error; err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return message, nil
	}

	var textParts []string
	var imageParts []ContentPart
	hasImages := false

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name))

		if imageFileExts[ext] {
			imageURL := file.URL
			if file.Storage != nil && file.Storage.Type == model.StorageTypeLocal {
				base64URL, err := s.localFileToBase64(file)
				if err == nil {
					imageURL = base64URL
				}
			}
			imageParts = append(imageParts, ContentPart{
				Type:     "image_url",
				ImageURL: &ImageURL{URL: imageURL},
			})
			hasImages = true
		} else if textFileExts[ext] {
			content, err := s.readFileContent(file)
			if err != nil {
				textParts = append(textParts, fmt.Sprintf("[文件: %s] (读取失败: %v)", file.Name, err))
			} else {
				textParts = append(textParts, fmt.Sprintf("[文件: %s]\n```\n%s\n```", file.Name, content))
			}
		} else {
			textParts = append(textParts, fmt.Sprintf("[文件: %s] (不支持的文件类型)", file.Name))
		}
	}

	fullText := message
	if len(textParts) > 0 {
		fullText = strings.Join(textParts, "\n\n") + "\n\n请基于以上文件内容回答：" + message
	}

	if hasImages {
		parts := []ContentPart{{Type: "text", Text: fullText}}
		parts = append(parts, imageParts...)
		return parts, nil
	}

	return fullText, nil
}

// 编码FileIDs为JSON字符串
func (s *AIService) encodeFileIDs(ids []uint) string {
	if len(ids) == 0 {
		return ""
	}
	data, _ := json.Marshal(ids)
	return string(data)
}
