package v1

import (
	"bufio"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	appconfig "server/config"
	"server/model/request"
	"server/model/response"
	"server/service"
)

type AIApi struct{}

var AI = new(AIApi)

func batchConversationDeleteFailureMessage(failedMsgs []string) string {
	if len(failedMsgs) == 0 {
		return "删除失败"
	}

	normalized := make([]string, 0, len(failedMsgs))
	for _, failedMsg := range failedMsgs {
		message := strings.TrimSpace(failedMsg)
		if message != "" {
			normalized = append(normalized, message)
		}
	}

	if len(normalized) == 0 {
		return "删除失败"
	}
	if len(normalized) == 1 {
		return normalized[0]
	}
	return strings.Join(normalized, "；")
}

// 获取模型列表
func (a *AIApi) GetModels(c *gin.Context) {
	models, err := service.AI.GetModels()
	if err != nil {
		response.Fail(c, "获取模型列表失败")
		return
	}
	response.OkWithData(c, models)
}

// 获取对话列表
func (a *AIApi) GetConversations(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.ConversationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := service.ConversationListInput{
		Page:     req.Page,
		PageSize: req.PageSize,
	}.Normalize()

	conversations, total, err := service.AI.GetConversations(userID, input)
	if err != nil {
		response.Fail(c, "获取对话列表失败")
		return
	}

	response.OkWithPage(c, conversations, total, input.Page, input.PageSize)
}

// 获取对话消息
func (a *AIApi) GetMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	messages, err := service.AI.GetMessages(uint(conversationID), userID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, messages)
}

// 创建对话
func (a *AIApi) CreateConversation(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	conversation, err := service.AI.CreateConversation(userID, service.CreateConversationInput{
		Title: req.Title,
		Model: req.Model,
	})
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, conversation)
}

// 删除对话
func (a *AIApi) DeleteConversation(c *gin.Context) {
	userID := c.GetUint("user_id")
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.AI.DeleteConversation(uint(conversationID), userID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

func (a *AIApi) BatchDeleteConversations(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.BatchConversationDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if len(req.Ids) == 0 {
		response.BadRequest(c, "请选择要删除的对话")
		return
	}

	successCount, failedMsgs := service.AI.BatchDeleteConversations(userID, req.Ids)
	if len(failedMsgs) == 0 {
		response.OkWithMessage(c, "batch_delete_success")
		return
	}
	if successCount > 0 {
		response.OkWithData(c, gin.H{
			"success_count": successCount,
			"failed_count":  len(failedMsgs),
			"failed_msgs":   failedMsgs,
		})
		return
	}

	response.Fail(c, batchConversationDeleteFailureMessage(failedMsgs))
}

// 普通对话（非流式）
func (a *AIApi) Chat(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	message, err := service.AI.Chat(userID, service.AIChatInput{
		ConversationID:   req.ConversationID,
		Model:            req.Model,
		Message:          req.Message,
		FileIDs:          req.FileIDs,
		EnableSearch:     req.EnableSearch,
		EnableThinking:   req.EnableThinking,
		SaveConversation: req.SaveConversation,
	})
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, message)
}

// 流式对话（SSE）
func (a *AIApi) ChatStream(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	conversationID, reader, saveConversation, sourcesMarkdown, err := service.AI.ChatStream(userID, service.AIChatInput{
		ConversationID:   req.ConversationID,
		Model:            req.Model,
		Message:          req.Message,
		FileIDs:          req.FileIDs,
		EnableSearch:     req.EnableSearch,
		EnableThinking:   req.EnableThinking,
		SaveConversation: req.SaveConversation,
	})
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	defer reader.Close()

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// 发送对话ID（只有保存对话时才发送）
	if saveConversation {
		c.SSEvent("conversation_id", conversationID)
		c.Writer.Flush()
	}

	// 收集完整响应用于保存
	accumulator := service.NewAIStreamAccumulator()

	// 解析并转发流式响应
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		event, err := accumulator.HandleLine(scanner.Text())
		if err != nil || event == nil {
			continue
		}
		c.SSEvent(event.Name, event.Data)
		c.Writer.Flush()
		if event.Done {
			break
		}
	}

	if sourcesMarkdown != "" {
		footerData := map[string]string{
			"content":           "\n\n---\n" + sourcesMarkdown,
			"reasoning_content": "",
			"finish_reason":     "",
		}
		if payload, marshalErr := json.Marshal(footerData); marshalErr == nil {
			c.SSEvent("message", string(payload))
			c.Writer.Flush()
		}
	}

	// 保存助手消息（只有启用保存时才保存）
	if saveConversation && (accumulator.FullContent() != "" || accumulator.FullReasoning() != "") {
		content := accumulator.FullContent()
		if sourcesMarkdown != "" {
			content = content + "\n\n---\n" + sourcesMarkdown
		}
		_ = service.AI.SaveAssistantMessage(conversationID, content, accumulator.FullReasoning())
	}
}

// 更新对话标题
func (a *AIApi) UpdateConversationTitle(c *gin.Context) {
	userID := c.GetUint("user_id")
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 验证对话属于当前用户并更新
	result := service.AI.UpdateConversationTitle(uint(conversationID), userID, req.Title)
	if result != nil {
		response.Fail(c, result.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// 清空对话消息
func (a *AIApi) ClearMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.AI.ClearMessages(uint(conversationID), userID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "清空成功")
}

// 清空上下文（保留聊天记录）
func (a *AIApi) ClearContext(c *gin.Context) {
	userID := c.GetUint("user_id")
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.AI.ClearContext(uint(conversationID), userID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "上下文已清空")
}

// 删除单条消息
func (a *AIApi) DeleteMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	messageID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.AI.DeleteMessage(uint(messageID), userID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

// 获取单个对话详情
func (a *AIApi) GetConversation(c *gin.Context) {
	userID := c.GetUint("user_id")
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	conversation, err := service.AI.GetConversation(uint(conversationID), userID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	// 同时获取消息
	messages, _ := service.AI.GetMessages(uint(conversationID), userID)

	response.OkWithData(c, gin.H{
		"conversation": conversation,
		"messages":     messages,
	})
}

// 测试AI配置
func (a *AIApi) TestConfig(c *gin.Context) {
	var req request.AITestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	err := service.AI.TestConfig(req.APIKey, req.BaseURL, req.Model)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "测试成功")
}

func (a *AIApi) FetchProviderModels(c *gin.Context) {
	var req request.AIProviderModelsFetchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	models, err := service.AI.FetchProviderModels(req.APIKey, req.BaseURL, req.ProviderName)
	if err != nil {
		if fetchErr, ok := err.(*service.AIProviderModelFetchError); ok {
			response.FailWithCode(c, fetchErr.Code, fetchErr.Message)
			return
		}
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, response.AIProviderModelsFetchResponse{Models: models})
}

// AdminListAIUsers 管理员视角分页查询有 AI 对话记录的用户列表（page 分页）
func (a *AIApi) AdminListAIUsers(c *gin.Context) {
	var req request.AdminAIUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := service.AdminAIUserListInput{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}.Normalize()

	items, total, err := service.AI.AdminListAIUsers(input)
	if err != nil {
		response.Fail(c, "获取用户列表失败")
		return
	}

	response.OkWithPage(c, items, total, input.Page, input.PageSize)
}

// AdminListConversations 管理员视角分页查询所有对话（cursor 分页）
func (a *AIApi) AdminListConversations(c *gin.Context) {
	var req request.AdminConversationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	items, nextCursor, hasMore, err := service.AI.AdminListConversations(service.AdminConversationListInput{
		CursorInput: service.CursorInput{
			Cursor: req.Cursor,
			Limit:  req.Limit,
		},
		UserID:  req.UserID,
		Keyword: req.Keyword,
	})
	if err != nil {
		response.Fail(c, "获取对话列表失败")
		return
	}

	response.OkWithData(c, gin.H{
		"list":        items,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}

// AdminListMessages 管理员视角分页查询某对话的消息（cursor 分页）
func (a *AIApi) AdminListMessages(c *gin.Context) {
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.AdminMessageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	messages, nextCursor, hasMore, err := service.AI.AdminListMessages(uint(conversationID), service.CursorInput{
		Cursor: req.Cursor,
		Limit:  req.Limit,
	})
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, gin.H{
		"list":        messages,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}

// AdminDeleteConversation 管理员删除任意对话
func (a *AIApi) AdminDeleteConversation(c *gin.Context) {
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.AI.AdminDeleteConversation(uint(conversationID)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

func (a *AIApi) GetAdminConfig(c *gin.Context) {
	cfg, err := service.AI.GetAdminConfig()
	if err != nil {
		response.Fail(c, "获取 AI 配置失败")
		return
	}

	response.OkWithData(c, cfg)
}

func (a *AIApi) UpdateAdminConfig(c *gin.Context) {
	var req appconfig.AI
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.AI.SaveAdminConfig(&req); err != nil {
		response.Fail(c, "保存 AI 配置失败")
		return
	}

	response.OkWithMessage(c, "更新成功")
}
