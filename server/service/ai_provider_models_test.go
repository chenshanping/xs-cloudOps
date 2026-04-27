package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestFetchProviderModelsUsesDirectModelsWhenBaseURLAlreadyContainsV1(t *testing.T) {
	var requestedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		if got := r.Header.Get("Authorization"); got != "Bearer sk-test" {
			t.Fatalf("authorization header = %s, want %s", got, "Bearer sk-test")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[{"id":"gpt-4o","object":"model","created":1712345678,"owned_by":"openai"},{"id":""}]}`))
	}))
	defer server.Close()

	models, err := AI.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/compatible-mode/v1")
	if err != nil {
		t.Fatalf("fetch provider models: %v", err)
	}
	if requestedPath != "/compatible-mode/v1/models" {
		t.Fatalf("requested path = %s, want %s", requestedPath, "/compatible-mode/v1/models")
	}
	if len(models) != 1 {
		t.Fatalf("models len = %d, want 1", len(models))
	}
	if models[0].ID != "gpt-4o" {
		t.Fatalf("model id = %s, want %s", models[0].ID, "gpt-4o")
	}
	if models[0].Object != "model" {
		t.Fatalf("model object = %s, want %s", models[0].Object, "model")
	}
	if models[0].Created != 1712345678 {
		t.Fatalf("model created = %d, want %d", models[0].Created, int64(1712345678))
	}
	if models[0].OwnedBy != "openai" {
		t.Fatalf("model owned_by = %s, want %s", models[0].OwnedBy, "openai")
	}
}

func TestFetchProviderModelsFallsBackWhenV1ModelsEndpointMissing(t *testing.T) {
	var v1Attempts int32
	var fallbackAttempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/proxy/v1/models":
			atomic.AddInt32(&v1Attempts, 1)
			http.Error(w, `{"error":{"message":"not found"}}`, http.StatusNotFound)
		case "/proxy/models":
			atomic.AddInt32(&fallbackAttempts, 1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[{"id":"deepseek-v3"}]}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	models, err := AI.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/proxy")
	if err != nil {
		t.Fatalf("fetch provider models with fallback: %v", err)
	}
	if atomic.LoadInt32(&v1Attempts) != 1 {
		t.Fatalf("v1 attempts = %d, want 1", v1Attempts)
	}
	if atomic.LoadInt32(&fallbackAttempts) != 1 {
		t.Fatalf("fallback attempts = %d, want 1", fallbackAttempts)
	}
	if len(models) != 1 || models[0].ID != "deepseek-v3" {
		t.Fatalf("models = %#v, want one deepseek-v3 item", models)
	}
}

func TestFetchProviderModelsSanitizesUpstreamErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"message":"invalid api key sk-secret-123456"}}`, http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := AI.fetchProviderModelsWithClient(server.Client(), "sk-secret-123456", server.URL+"/v1")
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "sk-secret-123456") {
		t.Fatalf("error leaked api key: %v", err)
	}
	if !strings.Contains(err.Error(), "API Key") {
		t.Fatalf("error = %v, want sanitized API Key guidance", err)
	}
}
