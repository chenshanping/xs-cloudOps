package service

import "testing"

func TestAIStreamAccumulatorIgnoresNonDataLine(t *testing.T) {
	acc := NewAIStreamAccumulator()

	event, err := acc.HandleLine("event: message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event != nil {
		t.Fatalf("expected nil event for non-data line")
	}
}

func TestAIStreamAccumulatorBuildsMessagePayloadAndAggregatesContent(t *testing.T) {
	acc := NewAIStreamAccumulator()

	event, err := acc.HandleLine(`data: {"choices":[{"delta":{"content":"hello","reasoning_content":"think"},"finish_reason":""}]}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event == nil {
		t.Fatalf("expected stream event")
	}
	if event.Name != "message" {
		t.Fatalf("unexpected event name: %s", event.Name)
	}
	if event.Data != `{"content":"hello","finish_reason":"","reasoning_content":"think"}` {
		t.Fatalf("unexpected event data: %s", event.Data)
	}
	if acc.FullContent() != "hello" {
		t.Fatalf("unexpected full content: %s", acc.FullContent())
	}
	if acc.FullReasoning() != "think" {
		t.Fatalf("unexpected full reasoning: %s", acc.FullReasoning())
	}
}

func TestAIStreamAccumulatorHandlesDoneMarker(t *testing.T) {
	acc := NewAIStreamAccumulator()

	event, err := acc.HandleLine("data: [DONE]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event == nil {
		t.Fatalf("expected done event")
	}
	if !event.Done {
		t.Fatalf("expected done event flag")
	}
	if event.Data != "[DONE]" {
		t.Fatalf("unexpected done payload: %s", event.Data)
	}
}
