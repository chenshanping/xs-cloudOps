package v1

import (
	"bufio"
	"strconv"

	"github.com/gin-gonic/gin"

	"server/model/request"
	"server/model/response"
	"server/service"
)

type AIApi struct{}

var AI = new(AIApi)

// 获取模型列表
func (a *AIApi) GetModels(c *gin.Context) {
	models := service.AI.GetModels()
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

	conversationID, reader, saveConversation, err := service.AI.ChatStream(userID, service.AIChatInput{
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

	// 保存助手消息（只有启用保存时才保存）
	if saveConversation && (accumulator.FullContent() != "" || accumulator.FullReasoning() != "") {
		_ = service.AI.SaveAssistantMessage(conversationID, accumulator.FullContent(), accumulator.FullReasoning())
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
