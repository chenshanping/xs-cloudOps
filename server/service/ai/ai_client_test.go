package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"server/global"
)

func TestAIServiceTestConfigBuildsAssistantMessageForXiaomiMimoTTS(t *testing.T) {
	previousLog := global.Log
	global.Log = zap.NewNop().Sugar()
	t.Cleanup(func() {
		global.Log = previousLog
	})

	var payload ChatCompletionRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("request path = %s, want /chat/completions", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"test","object":"chat.completion","created":1,"model":"mimo-v2-tts","choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`))
	}))
	defer server.Close()

	if err := Default.TestConfig("sk-test", server.URL, "mimo-v2-tts"); err != nil {
		t.Fatalf("TestConfig error: %v", err)
	}

	if len(payload.Messages) < 2 {
		t.Fatalf("messages len = %d, want at least 2", len(payload.Messages))
	}
	last := payload.Messages[len(payload.Messages)-1]
	if last.Role != "assistant" {
		t.Fatalf("last role = %s, want assistant", last.Role)
	}
	if last.Content != "你好，请回复OK" {
		t.Fatalf("last content = %#v, want test content", last.Content)
	}
}

func TestAIServiceTestConfigKeepsSingleUserMessageForNormalModel(t *testing.T) {
	previousLog := global.Log
	global.Log = zap.NewNop().Sugar()
	t.Cleanup(func() {
		global.Log = previousLog
	})

	var payload ChatCompletionRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"test","object":"chat.completion","created":1,"model":"gpt-4o-mini","choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`))
	}))
	defer server.Close()

	if err := Default.TestConfig("sk-test", server.URL, "gpt-4o-mini"); err != nil {
		t.Fatalf("TestConfig error: %v", err)
	}

	if len(payload.Messages) != 1 {
		t.Fatalf("messages len = %d, want 1", len(payload.Messages))
	}
	if payload.Messages[0].Role != "user" {
		t.Fatalf("first role = %s, want user", payload.Messages[0].Role)
	}
}
