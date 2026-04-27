package service

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
