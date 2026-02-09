package service

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-base-server/global"
	"go-base-server/model"
	"go-base-server/model/request"
)

type AIService struct{}

var AI = new(AIService)

// 模型信息（包含平台信息）
type ModelInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ProviderName string `json:"provider_name"` // 所属平台名称
}

// OpenAI兼容的消息格式
type ChatMessage struct {
	Role             string      `json:"role"`
	Content          interface{} `json:"content"`
	ReasoningContent string      `json:"reasoning_content,omitempty"`
}

// 多模态内容块
type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

// 文本类文件扩展名
var textFileExts = map[string]bool{
	".txt": true, ".md": true, ".csv": true, ".json": true, ".log": true,
	".go": true, ".py": true, ".js": true, ".ts": true, ".java": true,
	".html": true, ".css": true, ".sql": true, ".xml": true, ".yaml": true, ".yml": true,
	".sh": true, ".bat": true, ".c": true, ".cpp": true, ".h": true, ".rs": true,
	".rb": true, ".php": true, ".swift": true, ".kt": true, ".vue": true, ".jsx": true, ".tsx": true,
	".ini": true, ".toml": true, ".env": true, ".conf": true,
}

// 图片类文件扩展名
var imageFileExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true,
}

// OpenAI兼容的请求格式
type ChatCompletionRequest struct {
	Model          string                 `json:"model"`
	Messages       []ChatMessage          `json:"messages"`
	Stream         bool                   `json:"stream"`
	EnableSearch   bool                   `json:"enable_search,omitempty"`
	EnableThinking bool                   `json:"enable_thinking,omitempty"`
	ExtraBody      map[string]interface{} `json:"extra_body,omitempty"`
}

// OpenAI兼容的响应格式
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role             string `json:"role"`
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// 流式响应的chunk格式
type ChatCompletionChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role             string `json:"role,omitempty"`
			Content          string `json:"content,omitempty"`
			ReasoningContent string `json:"reasoning_content,omitempty"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason,omitempty"`
	} `json:"choices"`
}

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
func (s *AIService) GetConversations(userID uint, req *request.ConversationListRequest) ([]model.AIConversation, int64, error) {
	var conversations []model.AIConversation
	var total int64

	db := global.DB.Model(&model.AIConversation{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	if err := db.Offset(offset).Limit(req.PageSize).Order("updated_at DESC").Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// 获取对话消息
func (s *AIService) GetMessages(conversationID uint, userID uint) ([]model.AIMessage, error) {
	// 验证对话属于当前用户
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
func (s *AIService) CreateConversation(userID uint, req *request.CreateConversationRequest) (*model.AIConversation, error) {
	modelName := req.Model
	if modelName == "" {
		// 使用默认平台的第一个模型
		aiConfig := Config.GetAIConfig()
		if provider := aiConfig.GetDefaultProvider(); provider != nil && len(provider.Models) > 0 {
			modelName = provider.Models[0].ID
		}
	}

	title := req.Title
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
	// 验证对话属于当前用户
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("对话不存在")
	}

	// 删除消息
	if err := global.DB.Where("conversation_id = ?", conversationID).Delete(&model.AIMessage{}).Error; err != nil {
		return err
	}

	// 删除对话
	return global.DB.Delete(&conversation).Error
}

// 普通对话（非流式）
func (s *AIService) Chat(userID uint, req *request.AIChatRequest) (*model.AIMessage, error) {
	// 获取或创建对话
	var conversation *model.AIConversation
	var err error

	if req.ConversationID > 0 {
		conversation = &model.AIConversation{}
		if err := global.DB.Where("id = ? AND user_id = ?", req.ConversationID, userID).First(conversation).Error; err != nil {
			return nil, errors.New("对话不存在")
		}
	} else {
		conversation, err = s.CreateConversation(userID, &request.CreateConversationRequest{
			Model: req.Model,
		})
		if err != nil {
			return nil, err
		}
	}

	// 获取历史消息构建上下文
	messages, err := s.buildContextMessages(conversation.ID, req.Message, req.FileIDs)
	if err != nil {
		return nil, err
	}

	// 保存用户消息
	userMessage := &model.AIMessage{
		ConversationID: conversation.ID,
		Role:           "user",
		Content:        req.Message,
		FileIDs:        s.encodeFileIDs(req.FileIDs),
	}
	if err := global.DB.Create(userMessage).Error; err != nil {
		return nil, err
	}

	// 更新对话标题（如果是第一条消息）
	if conversation.Title == "新对话" {
		title := req.Message
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		global.DB.Model(conversation).Update("title", title)
	}

	// 调用API
	modelName := conversation.Model
	if req.Model != "" {
		modelName = req.Model
	}

	response, err := s.callAPI(modelName, messages, false, req.EnableSearch, req.EnableThinking)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(response, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, errors.New("API返回空响应")
	}

	// 保存助手消息
	assistantMessage := &model.AIMessage{
		ConversationID:   conversation.ID,
		Role:             "assistant",
		Content:          chatResp.Choices[0].Message.Content,
		ReasoningContent: chatResp.Choices[0].Message.ReasoningContent,
	}
	if err := global.DB.Create(assistantMessage).Error; err != nil {
		return nil, err
	}

	// 更新对话时间
	global.DB.Model(conversation).Update("updated_at", assistantMessage.CreatedAt)

	return assistantMessage, nil
}

// 流式对话 - 返回reader供SSE使用
func (s *AIService) ChatStream(userID uint, req *request.AIChatRequest) (
	conversationID uint,
	reader io.ReadCloser,
	saveConversation bool,
	err error,
) {
	// 判断是否保存对话（默认为true）
	saveConversation = req.SaveConversation == nil || *req.SaveConversation

	var messages []ChatMessage
	var modelName string

	if saveConversation {
		// 需要保存对话 - 获取或创建对话
		var conversation *model.AIConversation

		if req.ConversationID > 0 {
			conversation = &model.AIConversation{}
			if err := global.DB.Where("id = ? AND user_id = ?", req.ConversationID, userID).First(conversation).Error; err != nil {
				return 0, nil, false, errors.New("对话不存在")
			}
		} else {
			conversation, err = s.CreateConversation(userID, &request.CreateConversationRequest{
				Model: req.Model,
			})
			if err != nil {
				return 0, nil, false, err
			}
		}

		// 获取历史消息构建上下文
		messages, err = s.buildContextMessages(conversation.ID, req.Message, req.FileIDs)
		if err != nil {
			return 0, nil, false, err
		}

		// 保存用户消息
		userMessage := &model.AIMessage{
			ConversationID: conversation.ID,
			Role:           "user",
			Content:        req.Message,
			FileIDs:        s.encodeFileIDs(req.FileIDs),
		}
		if err := global.DB.Create(userMessage).Error; err != nil {
			return 0, nil, false, err
		}

		// 更新对话标题（如果是第一条消息）
		if conversation.Title == "新对话" {
			title := req.Message
			if len(title) > 50 {
				title = title[:50] + "..."
			}
			global.DB.Model(conversation).Update("title", title)
		}

		conversationID = conversation.ID
		modelName = conversation.Model
		if req.Model != "" {
			modelName = req.Model
		}
	} else {
		// 不保存对话 - 直接构建消息
		messages = []ChatMessage{
			{Role: "user", Content: req.Message},
		}
		modelName = req.Model
		if modelName == "" {
			aiConfig := Config.GetAIConfig()
			provider := aiConfig.GetDefaultProvider()
			modelName = provider.Models[0].ID
		}
		conversationID = 0
	}

	// 调用流式API
	reader, err = s.callStreamAPI(modelName, messages, req.EnableSearch, req.EnableThinking)
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

	// 更新对话时间
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
	// 验证对话属于当前用户
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("对话不存在")
	}

	// 删除消息
	return global.DB.Where("conversation_id = ?", conversationID).Delete(&model.AIMessage{}).Error
}

// 清空上下文（保留聊天记录，但新消息不带历史上下文）
func (s *AIService) ClearContext(conversationID uint, userID uint) error {
	// 验证对话属于当前用户
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("对话不存在")
	}

	// 更新上下文清空时间
	now := time.Now()
	return global.DB.Model(&conversation).Update("context_cleared_at", now).Error
}

// 删除单条消息
func (s *AIService) DeleteMessage(messageID uint, userID uint) error {
	// 获取消息
	var message model.AIMessage
	if err := global.DB.First(&message, messageID).Error; err != nil {
		return errors.New("消息不存在")
	}

	// 验证对话属于当前用户
	var conversation model.AIConversation
	if err := global.DB.Where("id = ? AND user_id = ?", message.ConversationID, userID).First(&conversation).Error; err != nil {
		return errors.New("无权限删除此消息")
	}

	// 删除消息
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

// 构建上下文消息
func (s *AIService) buildContextMessages(conversationID uint, newMessage string, fileIDs []uint) ([]ChatMessage, error) {
	var messages []ChatMessage

	// 获取历史消息
	var history []model.AIMessage
	if conversationID > 0 {
		// 先获取对话信息，检查上下文清空时间
		var conversation model.AIConversation
		if err := global.DB.First(&conversation, conversationID).Error; err != nil {
			return nil, err
		}

		query := global.DB.Where("conversation_id = ?", conversationID)

		// 如果有上下文清空时间，只获取该时间之后的消息
		if conversation.ContextClearedAt != nil {
			query = query.Where("created_at > ?", conversation.ContextClearedAt)
		}

		if err := query.Order("created_at ASC").
			Limit(20). // 限制上下文长度
			Find(&history).Error; err != nil {
			return nil, err
		}
	}

	// 构建消息列表
	for _, msg := range history {
		messages = append(messages, ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 处理附件文件
	if len(fileIDs) > 0 {
		content, err := s.buildMessageWithFiles(newMessage, fileIDs)
		if err != nil {
			// 文件读取失败不阻断发送，回退到纯文本
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
			// 图片文件 - 使用多模态
			imageURL := file.URL
			// 本地存储的相对路径需要转换为完整URL
			if file.Storage != nil && file.Storage.Type == model.StorageTypeLocal {
				// 将本地文件转base64
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
			// 文本文件 - 读取内容
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

	// 拼接文本部分
	fullText := message
	if len(textParts) > 0 {
		fullText = strings.Join(textParts, "\n\n") + "\n\n请基于以上文件内容回答：" + message
	}

	if hasImages {
		// 多模态格式
		parts := []ContentPart{{Type: "text", Text: fullText}}
		parts = append(parts, imageParts...)
		return parts, nil
	}

	return fullText, nil
}

// 读取文件内容
func (s *AIService) readFileContent(file model.SysFile) (string, error) {
	const maxSize = 100 * 1024 // 100KB

	var data []byte
	var err error

	if file.Storage != nil && file.Storage.Type == model.StorageTypeLocal {
		// 本地存储 - 直接读取磁盘文件
		var config model.LocalConfig
		if jsonErr := json.Unmarshal([]byte(file.Storage.Config), &config); jsonErr != nil {
			return "", fmt.Errorf("解析存储配置失败: %v", jsonErr)
		}
		fullPath := filepath.Join(config.BasePath, file.Path)
		data, err = os.ReadFile(fullPath)
	} else {
		// 远程存储 - HTTP GET
		data, err = s.httpGetFileContent(file.URL)
	}

	if err != nil {
		return "", err
	}

	if len(data) > maxSize {
		return string(data[:maxSize]) + "\n...(文件内容已截断，超过100KB限制)", nil
	}
	return string(data), nil
}

// 本地图片转base64
func (s *AIService) localFileToBase64(file model.SysFile) (string, error) {
	const maxImageSize = 5 * 1024 * 1024 // 5MB

	var config model.LocalConfig
	if err := json.Unmarshal([]byte(file.Storage.Config), &config); err != nil {
		return "", err
	}
	fullPath := filepath.Join(config.BasePath, file.Path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	if len(data) > maxImageSize {
		return "", fmt.Errorf("图片超过5MB限制")
	}

	// 根据扩展名确定MIME类型
	mimeType := file.MimeType
	if mimeType == "" {
		ext := strings.ToLower(filepath.Ext(file.Name))
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		default:
			mimeType = "image/png"
		}
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, b64), nil
}

// 编码FileIDs为JSON字符串
func (s *AIService) encodeFileIDs(ids []uint) string {
	if len(ids) == 0 {
		return ""
	}
	data, _ := json.Marshal(ids)
	return string(data)
}

// HTTP获取文件内容
func (s *AIService) httpGetFileContent(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载文件失败: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %v", err)
	}
	return data, nil
}

// 调用API（非流式）
func (s *AIService) callAPI(modelName string, messages []ChatMessage, stream bool, enableSearch bool, enableThinking bool) ([]byte, error) {
	// 根据模型找到对应的平台
	aiConfig := Config.GetAIConfig()
	provider := aiConfig.GetProviderByModel(modelName)
	if provider == nil {
		// 找不到则使用默认平台
		provider = aiConfig.GetDefaultProvider()
	}
	if provider == nil {
		return nil, errors.New("未配置AI平台")
	}
	if provider.APIKey == "" {
		return nil, errors.New("未配置API Key")
	}

	reqBody := ChatCompletionRequest{
		Model:          modelName,
		Messages:       messages,
		Stream:         stream,
		EnableSearch:   enableSearch,
		EnableThinking: enableThinking,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", provider.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+provider.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败: %s", string(body))
	}

	return body, nil
}

// 调用流式API
func (s *AIService) callStreamAPI(modelName string, messages []ChatMessage, enableSearch bool, enableThinking bool) (io.ReadCloser, error) {
	// 根据模型找到对应的平台
	aiConfig := Config.GetAIConfig()
	provider := aiConfig.GetProviderByModel(modelName)
	if provider == nil {
		// 找不到则使用默认平台
		provider = aiConfig.GetDefaultProvider()
	}
	if provider == nil {
		return nil, errors.New("未配置AI平台")
	}
	if provider.APIKey == "" {
		return nil, errors.New("未配置API Key")
	}

	reqBody := ChatCompletionRequest{
		Model:          modelName,
		Messages:       messages,
		Stream:         true,
		EnableSearch:   enableSearch,
		EnableThinking: enableThinking,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", provider.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+provider.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API请求失败: %s", string(body))
	}

	return resp.Body, nil
}

// 解析SSE流数据
func ParseSSEStream(reader io.Reader, onChunk func(chunk *ChatCompletionChunk) error) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk ChatCompletionChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // 跳过解析失败的行
		}

		if err := onChunk(&chunk); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// 测试AI配置
func (s *AIService) TestConfig(apiKey, baseURL, model string) error {
	if apiKey == "" {
		return errors.New("未配置API Key")
	}
	if baseURL == "" {
		return errors.New("未配置Base URL")
	}
	if model == "" {
		return errors.New("未选择模型")
	}

	reqBody := ChatCompletionRequest{
		Model: model,
		Messages: []ChatMessage{
			{Role: "user", Content: "你好，请回复OK"},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	global.Log.Info("AI 请求日志:" + baseURL + "/chat/completions" + "\n参数为:" + string(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API请求失败(%d): %s", resp.StatusCode, string(body))
	}

	// 验证响应格式
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return errors.New("API返回空响应")
	}

	return nil
}
