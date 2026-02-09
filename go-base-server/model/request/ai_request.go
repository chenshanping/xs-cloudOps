package request

// AI对话请求
type AIChatRequest struct {
	ConversationID   uint   `json:"conversation_id" comment:"对话ID(为0则创建新对话)"`
	Model            string `json:"model" comment:"模型名称"`
	Message          string `json:"message" binding:"required" comment:"用户消息"`
	FileIDs          []uint `json:"file_ids" comment:"附件文件ID列表"`
	EnableSearch     bool   `json:"enable_search" comment:"是否启用联网搜索"`
	EnableThinking   bool   `json:"enable_thinking" comment:"是否启用思考模式"`
	SaveConversation *bool  `json:"save_conversation" comment:"是否保存对话记录(默认true)"`
}

// 创建对话请求
type CreateConversationRequest struct {
	Title string `json:"title" comment:"对话标题"`
	Model string `json:"model" comment:"模型名称"`
}

// 对话列表请求
type ConversationListRequest struct {
	PageRequest
}

// 对话消息列表请求
type MessageListRequest struct {
	ConversationID uint `json:"conversation_id" form:"conversation_id" binding:"required" comment:"对话ID"`
}

type AITestRequest struct {
	APIKey  string `json:"api_key" binding:"required"`
	BaseURL string `json:"base_url" binding:"required"`
	Model   string `json:"model" binding:"required"`
}
