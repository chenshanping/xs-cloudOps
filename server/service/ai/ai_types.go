package ai

import "context"

type ConversationListInput struct {
	Page     int
	PageSize int
}

func (in ConversationListInput) Normalize() ConversationListInput {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	return in
}

func (in ConversationListInput) Offset() int {
	normalized := in.Normalize()
	return (normalized.Page - 1) * normalized.PageSize
}

type CreateConversationInput struct {
	Title string
	Model string
}

// CursorInput cursor 分页通用入参
type CursorInput struct {
	Cursor uint // 上一页最后一条的 id；0 表示首页
	Limit  int  // 每页条数
}

func (in CursorInput) Normalize() CursorInput {
	if in.Limit <= 0 {
		in.Limit = 20
	}
	if in.Limit > 100 {
		in.Limit = 100
	}
	return in
}

// AdminConversationListInput 管理员侧对话列表查询入参
type AdminConversationListInput struct {
	CursorInput
	UserID  uint   // 按用户筛选；0 表示不筛选
	Keyword string // 按标题模糊查询；空表示不筛选
}

// AdminConversationItem 管理员视角的对话条目（带用户信息）
type AdminConversationItem struct {
	ID               uint   `json:"id"`
	UserID           uint   `json:"user_id"`
	Username         string `json:"username"`
	Nickname         string `json:"nickname"`
	Title            string `json:"title"`
	Model            string `json:"model"`
	MessageCount     int64  `json:"message_count"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	ContextClearedAt string `json:"context_cleared_at,omitempty"`
}

// AdminAIUserListInput 管理员侧 AI 活跃用户列表查询入参（page 分页）
type AdminAIUserListInput struct {
	Page     int
	PageSize int
	Keyword  string // 按用户名/昵称模糊查询
}

func (in AdminAIUserListInput) Normalize() AdminAIUserListInput {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 50
	}
	if in.PageSize > 200 {
		in.PageSize = 200
	}
	return in
}

func (in AdminAIUserListInput) Offset() int {
	normalized := in.Normalize()
	return (normalized.Page - 1) * normalized.PageSize
}

// AdminAIUserItem 有 AI 对话记录的用户条目
type AdminAIUserItem struct {
	ID                uint   `json:"id"`
	Username          string `json:"username"`
	Nickname          string `json:"nickname"`
	ConversationCount int64  `json:"conversation_count"`
	LastActiveAt      string `json:"last_active_at"`
}

type AIChatInput struct {
	ConversationID   uint
	Model            string
	Message          string
	FileIDs          []uint
	EnableSearch     bool
	EnableThinking   bool
	SaveConversation *bool
}

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
	SearchOptions  *SearchOptions         `json:"search_options,omitempty"`
	ExtraBody      map[string]interface{} `json:"extra_body,omitempty"`
}

type SearchOptions struct {
	ForcedSearch bool `json:"forced_search,omitempty"`
}

type AISearchSource struct {
	Title     string
	URL       string
	Published string
	Snippet   string
}

type AISearchGrounding struct {
	Query    string
	Evidence string
	Sources  []AISearchSource
}

type AISearcher interface {
	Search(ctx context.Context, query string) (*AISearchGrounding, error)
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
