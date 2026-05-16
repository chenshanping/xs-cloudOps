package ai

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
		_, _ = w.Write([]byte(`{"object":"list","data":[{"id":"gpt-4o","name":"GPT-4o","description":"旗舰模型","object":"model","created":1712345678,"owned_by":"openai","tags":["reasoning","vision","search","tool","free","reasoning",""],"context_window":128000},{"id":""}]}`))
	}))
	defer server.Close()

	models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/compatible-mode/v1", "OpenAI")
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
	if models[0].Name != "GPT-4o" {
		t.Fatalf("model name = %s, want %s", models[0].Name, "GPT-4o")
	}
	if models[0].Description != "旗舰模型" {
		t.Fatalf("model description = %s, want %s", models[0].Description, "旗舰模型")
	}
	if !models[0].IsThinking || !models[0].SupportVision || !models[0].SupportTools || !models[0].IsFree {
		t.Fatalf("model capability flags = %#v, want true flags from tags", models[0])
	}
	if models[0].SearchStrategy == "" || models[0].SearchStrategy == "none" {
		t.Fatalf("model search_strategy = %q, want searchable model", models[0].SearchStrategy)
	}
	if models[0].ContextWindow == nil || *models[0].ContextWindow != 128000 {
		t.Fatalf("model context_window = %#v, want 128000", models[0].ContextWindow)
	}
	if len(models[0].Tags) < 5 {
		t.Fatalf("model tags = %#v, want normalized tags retained", models[0].Tags)
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

	models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/proxy", "OpenAI Compatible")
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

	_, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-secret-123456", server.URL+"/v1", "OpenAI")
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

func TestFetchProviderModelsInfersCapabilitiesFromModelName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"bge-reranker-v2-m3","name":"BGE Reranker V2 M3"},{"id":"qwen2.5-vl-72b-instruct","name":"Qwen2.5 VL 72B Instruct"},{"id":"deepseek-r1","name":"DeepSeek R1"},{"id":"deepseek-chat","name":"DeepSeek Chat"}]}`))
	}))
	defer server.Close()

	models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/v1", "OpenAI Compatible")
	if err != nil {
		t.Fatalf("fetch provider models: %v", err)
	}
	if len(models) != 4 {
		t.Fatalf("models len = %d, want 4", len(models))
	}
	if !models[0].SupportEmbedding || !models[0].SupportRerank {
		t.Fatalf("reranker model capabilities = %#v, want embedding and rerank", models[0])
	}
	if !models[1].SupportVision {
		t.Fatalf("vision model capabilities = %#v, want support_vision", models[1])
	}
	if !models[2].IsThinking {
		t.Fatalf("reasoning model capabilities = %#v, want is_thinking", models[2])
	}
	if !models[3].SupportTools {
		t.Fatalf("tool capable model capabilities = %#v, want support_tools", models[3])
	}
	if models[3].SearchStrategy != searchStrategyNone {
		t.Fatalf("tool capable model search_strategy = %q, want %q", models[3].SearchStrategy, searchStrategyNone)
	}
}

func TestFetchProviderModelsPreservesExplicitToolSearchStrategy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"custom-agent","name":"Custom Agent","search_strategy":"tool"}]}`))
	}))
	defer server.Close()

	models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/v1", "OpenAI Compatible")
	if err != nil {
		t.Fatalf("fetch provider models: %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("models len = %d, want 1", len(models))
	}
	if models[0].SearchStrategy != searchStrategyTool {
		t.Fatalf("explicit tool search_strategy = %q, want %q", models[0].SearchStrategy, searchStrategyTool)
	}
	if !models[0].SupportTools {
		t.Fatalf("explicit tool search model should imply support_tools: %#v", models[0])
	}
}

func TestFetchProviderModelsInfersProviderAwareBuiltinSearch(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		modelID      string
		wantSearch   string
	}{
		{
			name:         "dashscope qwen max",
			providerName: "阿里云百炼",
			modelID:      "qwen-max",
			wantSearch:   searchStrategyBuiltin,
		},
		{
			name:         "perplexity sonar pro",
			providerName: "Perplexity",
			modelID:      "sonar-pro",
			wantSearch:   searchStrategyBuiltin,
		},
		{
			name:         "openai gpt 4o mini",
			providerName: "OpenAI",
			modelID:      "gpt-4o-mini",
			wantSearch:   searchStrategyBuiltin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[{"id":"` + tt.modelID + `","name":"` + tt.modelID + `"}]}`))
			}))
			defer server.Close()

			models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/v1", tt.providerName)
			if err != nil {
				t.Fatalf("fetch provider models: %v", err)
			}
			if len(models) != 1 {
				t.Fatalf("models len = %d, want 1", len(models))
			}
			if models[0].SearchStrategy != tt.wantSearch {
				t.Fatalf("search_strategy = %q, want %q for provider %s model %s", models[0].SearchStrategy, tt.wantSearch, tt.providerName, tt.modelID)
			}
		})
	}
}

func TestFetchProviderModelsDoesNotMarkDashScopeKimiK26AsVision(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"kimi-k2.6","name":"Kimi K2.6"}]}`))
	}))
	defer server.Close()

	models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/v1", "阿里云百炼")
	if err != nil {
		t.Fatalf("fetch provider models: %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("models len = %d, want 1", len(models))
	}
	if models[0].SupportVision {
		t.Fatalf("dashscope kimi-k2.6 support_vision = %v, want false", models[0].SupportVision)
	}
	if !models[0].IsThinking {
		t.Fatalf("dashscope kimi-k2.6 is_thinking = %v, want true", models[0].IsThinking)
	}
	if !models[0].SupportTools {
		t.Fatalf("dashscope kimi-k2.6 support_tools = %v, want true", models[0].SupportTools)
	}
}

func TestFetchProviderModelsMarksDashScopeQwen35PlusAsReasoning(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"qwen3.5-plus-2026-04-20","name":"Qwen3.5 Plus 2026-04-20"}]}`))
	}))
	defer server.Close()

	models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/v1", "阿里云百炼")
	if err != nil {
		t.Fatalf("fetch provider models: %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("models len = %d, want 1", len(models))
	}
	if !models[0].IsThinking {
		t.Fatalf("dashscope qwen3.5-plus-2026-04-20 is_thinking = %v, want true", models[0].IsThinking)
	}
}

func TestFetchProviderModelsDoesNotMarkLegacyQwen3MaxSnapshotAsReasoning(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"qwen3-max-2025-09-23","name":"Qwen3 Max 2025-09-23"}]}`))
	}))
	defer server.Close()

	models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/v1", "阿里云百炼")
	if err != nil {
		t.Fatalf("fetch provider models: %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("models len = %d, want 1", len(models))
	}
	if models[0].IsThinking {
		t.Fatalf("dashscope qwen3-max-2025-09-23 is_thinking = %v, want false", models[0].IsThinking)
	}
}

func TestFetchProviderModelsAppliesXiaomiMimoPricingPageCapabilities(t *testing.T) {
	tests := []struct {
		name         string
		modelID      string
		wantThinking bool
		wantVision   bool
		wantTools    bool
		wantSearch   string
		wantFree     bool
	}{
		{
			name:         "mimo v2.5 pro reasoning",
			modelID:      "mimo-v2.5-pro",
			wantThinking: true,
			wantVision:   false,
			wantTools:    true,
			wantSearch:   searchStrategyNone,
			wantFree:     false,
		},
		{
			name:         "mimo v2.5 multimodal",
			modelID:      "mimo-v2.5",
			wantThinking: false,
			wantVision:   true,
			wantTools:    true,
			wantSearch:   searchStrategyNone,
			wantFree:     false,
		},
		{
			name:         "mimo v2.5 tts free and not tool",
			modelID:      "mimo-v2.5-tts",
			wantThinking: false,
			wantVision:   false,
			wantTools:    false,
			wantSearch:   searchStrategyNone,
			wantFree:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[{"id":"` + tt.modelID + `","name":"` + tt.modelID + `"}]}`))
			}))
			defer server.Close()

			models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/v1", "Xiaomi MiMo")
			if err != nil {
				t.Fatalf("fetch provider models: %v", err)
			}
			if len(models) != 1 {
				t.Fatalf("models len = %d, want 1", len(models))
			}

			got := models[0]
			if got.IsThinking != tt.wantThinking {
				t.Fatalf("is_thinking = %v, want %v for %s", got.IsThinking, tt.wantThinking, tt.modelID)
			}
			if got.SupportVision != tt.wantVision {
				t.Fatalf("support_vision = %v, want %v for %s", got.SupportVision, tt.wantVision, tt.modelID)
			}
			if got.SupportTools != tt.wantTools {
				t.Fatalf("support_tools = %v, want %v for %s", got.SupportTools, tt.wantTools, tt.modelID)
			}
			if got.SearchStrategy != tt.wantSearch {
				t.Fatalf("search_strategy = %q, want %q for %s", got.SearchStrategy, tt.wantSearch, tt.modelID)
			}
			if got.IsFree != tt.wantFree {
				t.Fatalf("is_free = %v, want %v for %s", got.IsFree, tt.wantFree, tt.modelID)
			}
		})
	}
}

func TestFetchProviderModelsDerivesStableModelGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"mimo-v2.5-pro","name":"MiMo V2.5 Pro","owned_by":"xiaomi"},{"id":"deepseek-v4-flash","name":"DeepSeek V4 Flash","owned_by":"system"},{"id":"deepseek-ai/DeepSeek-R1-0528","name":"DeepSeek R1 0528","owned_by":"deepseek"}]}`))
	}))
	defer server.Close()

	models, err := Default.fetchProviderModelsWithClient(server.Client(), "sk-test", server.URL+"/v1", "OpenAI Compatible")
	if err != nil {
		t.Fatalf("fetch provider models: %v", err)
	}
	if len(models) != 3 {
		t.Fatalf("models len = %d, want 3", len(models))
	}

	groupByID := make(map[string]string, len(models))
	for _, item := range models {
		groupByID[item.ID] = item.Group
	}

	if groupByID["mimo-v2.5-pro"] != "mimo-v2.5" {
		t.Fatalf("mimo-v2.5-pro group = %q, want %q", groupByID["mimo-v2.5-pro"], "mimo-v2.5")
	}
	if groupByID["deepseek-v4-flash"] != "deepseek-v4" {
		t.Fatalf("deepseek-v4-flash group = %q, want %q", groupByID["deepseek-v4-flash"], "deepseek-v4")
	}
	if groupByID["deepseek-ai/DeepSeek-R1-0528"] != "deepseek-ai" {
		t.Fatalf("deepseek-ai/DeepSeek-R1-0528 group = %q, want %q", groupByID["deepseek-ai/DeepSeek-R1-0528"], "deepseek-ai")
	}
}
