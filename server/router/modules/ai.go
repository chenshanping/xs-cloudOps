package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/model/request"
	"server/router/registry"
)

func init() {
	RegisterModule(&AIModule{})
}

type AIModule struct{}

func (m *AIModule) Name() string {
	return "AI对话"
}

func (m *AIModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 模型列表可以公开访问
	R(rg, "GET", "/ai/models", m.Name(), "获取模型列表", v1.AI.GetModels)
}

func (m *AIModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	// 对话管理
	R(rg, "GET", "/ai/conversations", m.Name(), "获取对话列表", v1.AI.GetConversations,
		registry.WithAuth(), registry.WithRequest(request.ConversationListRequest{}))
	R(rg, "POST", "/ai/conversations", m.Name(), "创建对话", v1.AI.CreateConversation,
		registry.WithAuth(), registry.WithRequest(request.CreateConversationRequest{}))
	R(rg, "DELETE", "/ai/conversations/batch", m.Name(), "批量删除对话", v1.AI.BatchDeleteConversations,
		registry.WithAuth(), registry.WithRequest(request.BatchConversationDeleteRequest{}))
	R(rg, "GET", "/ai/conversations/:id", m.Name(), "获取对话详情", v1.AI.GetConversation, registry.WithAuth())
	R(rg, "PUT", "/ai/conversations/:id", m.Name(), "更新对话标题", v1.AI.UpdateConversationTitle, registry.WithAuth())
	R(rg, "DELETE", "/ai/conversations/:id", m.Name(), "删除对话", v1.AI.DeleteConversation, registry.WithAuth())
	R(rg, "GET", "/ai/conversations/:id/messages", m.Name(), "获取对话消息", v1.AI.GetMessages, registry.WithAuth())
	R(rg, "DELETE", "/ai/conversations/:id/messages", m.Name(), "清空对话消息", v1.AI.ClearMessages, registry.WithAuth())
	R(rg, "POST", "/ai/conversations/:id/clear-context", m.Name(), "清空上下文", v1.AI.ClearContext, registry.WithAuth())
	R(rg, "DELETE", "/ai/messages/:id", m.Name(), "删除单条消息", v1.AI.DeleteMessage, registry.WithAuth())

	// 对话功能
	R(rg, "POST", "/ai/chat", m.Name(), "AI对话", v1.AI.Chat,
		registry.WithAuth(), registry.WithRequest(request.AIChatRequest{}))
	R(rg, "POST", "/ai/chat/stream", m.Name(), "AI流式对话", v1.AI.ChatStream,
		registry.WithAuth(), registry.WithRequest(request.AIChatRequest{}))

	// 配置测试
	R(rg, "POST", "/ai/test", m.Name(), "测试AI配置", v1.AI.TestConfig, registry.WithAuth(), registry.WithRequest(request.AITestRequest{}))
	R(rg, "POST", "/ai/providers/models/fetch", m.Name(), "拉取平台模型列表", v1.AI.FetchProviderModels,
		registry.WithAuth(), registry.WithRequest(request.AIProviderModelsFetchRequest{}))
}
