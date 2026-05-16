package response

type AIProviderModelItem struct {
	ID               string   `json:"id"`
	Name             string   `json:"name,omitempty"`
	Group            string   `json:"group,omitempty"`
	Description      string   `json:"description,omitempty"`
	Object           string   `json:"object,omitempty"`
	Created          int64    `json:"created,omitempty"`
	OwnedBy          string   `json:"owned_by,omitempty"`
	IsThinking       bool     `json:"is_thinking"`
	SupportVision    bool     `json:"support_vision"`
	SupportTools     bool     `json:"support_tools"`
	SearchStrategy   string   `json:"search_strategy,omitempty"`
	SupportEmbedding bool     `json:"support_embedding"`
	SupportRerank    bool     `json:"support_rerank"`
	IsFree           bool     `json:"is_free"`
	Temperature      *float64 `json:"temperature,omitempty"`
	ContextWindow    *int     `json:"context_window,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

type AIProviderModelsFetchResponse struct {
	Models []AIProviderModelItem `json:"models"`
}
