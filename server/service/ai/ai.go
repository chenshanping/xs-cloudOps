package ai

type AIService struct {
	searcher AISearcher
}

var Default = &AIService{}
