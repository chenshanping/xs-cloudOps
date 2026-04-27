package response

type AIProviderModelItem struct {
	ID      string `json:"id"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
}

type AIProviderModelsFetchResponse struct {
	Models []AIProviderModelItem `json:"models"`
}
