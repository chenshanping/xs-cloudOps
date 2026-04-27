package service

import (
	"encoding/json"
	"strings"
)

type AIStreamEvent struct {
	Name string
	Data string
	Done bool
}

type AIStreamAccumulator struct {
	fullContent   strings.Builder
	fullReasoning strings.Builder
}

func NewAIStreamAccumulator() *AIStreamAccumulator {
	return &AIStreamAccumulator{}
}

func (a *AIStreamAccumulator) HandleLine(line string) (*AIStreamEvent, error) {
	if !strings.HasPrefix(line, "data: ") {
		return nil, nil
	}

	data := strings.TrimPrefix(line, "data: ")
	if data == "[DONE]" {
		return &AIStreamEvent{
			Name: "message",
			Data: "[DONE]",
			Done: true,
		}, nil
	}

	var chunk ChatCompletionChunk
	if err := json.Unmarshal([]byte(data), &chunk); err != nil {
		return nil, err
	}
	if len(chunk.Choices) == 0 {
		return nil, nil
	}

	delta := chunk.Choices[0].Delta
	if delta.Content != "" {
		a.fullContent.WriteString(delta.Content)
	}
	if delta.ReasoningContent != "" {
		a.fullReasoning.WriteString(delta.ReasoningContent)
	}

	eventData := map[string]interface{}{
		"content":           delta.Content,
		"reasoning_content": delta.ReasoningContent,
		"finish_reason":     chunk.Choices[0].FinishReason,
	}
	jsonData, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	return &AIStreamEvent{
		Name: "message",
		Data: string(jsonData),
	}, nil
}

func (a *AIStreamAccumulator) FullContent() string {
	return a.fullContent.String()
}

func (a *AIStreamAccumulator) FullReasoning() string {
	return a.fullReasoning.String()
}
