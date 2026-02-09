package model

import "time"

// AI对话会话
type AIConversation struct {
	BaseModel
	UserID           uint       `json:"user_id" gorm:"index;comment:用户ID"`
	Title            string     `json:"title" gorm:"size:255;comment:对话标题"`
	Model            string     `json:"model" gorm:"size:100;comment:使用的模型"`
	ContextClearedAt *time.Time `json:"context_cleared_at" gorm:"comment:上下文清空时间"`
}

func (AIConversation) TableName() string {
	return "ai_conversations"
}

// AI对话消息
type AIMessage struct {
	BaseModel
	ConversationID   uint   `json:"conversation_id" gorm:"index;comment:对话ID"`
	Role             string `json:"role" gorm:"size:20;comment:角色(user/assistant/system)"`
	Content          string `json:"content" gorm:"type:longtext;comment:消息内容"`
	ReasoningContent string `json:"reasoning_content" gorm:"type:longtext;comment:思考过程"`
	FileIDs          string `json:"file_ids" gorm:"size:500;comment:附件文件ID(JSON数组)"`
}

func (AIMessage) TableName() string {
	return "ai_messages"
}
