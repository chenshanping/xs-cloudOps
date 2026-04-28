package ai

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExaMCPClientSearchUsesLifecycleAndParsesSources(t *testing.T) {
	t.Helper()

	var methods []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		defer r.Body.Close()

		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal payload: %v", err)
		}

		method, _ := payload["method"].(string)
		methods = append(methods, method)

		switch method {
		case "initialize":
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = io.WriteString(w, "event: message\n")
			_, _ = io.WriteString(w, "data: {\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"protocolVersion\":\"2025-03-26\",\"capabilities\":{\"tools\":{\"listChanged\":true}},\"serverInfo\":{\"name\":\"exa-test\",\"version\":\"1.0.0\"}}}\n\n")
		case "notifications/initialized":
			w.WriteHeader(http.StatusAccepted)
		case "tools/call":
			params := payload["params"].(map[string]interface{})
			if params["name"] != "web_search_exa" {
				t.Fatalf("tool name = %v, want web_search_exa", params["name"])
			}
			args := params["arguments"].(map[string]interface{})
			if args["query"] != "2026 伦敦世乒赛 男团名单" {
				t.Fatalf("query = %v", args["query"])
			}

			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = io.WriteString(w, "event: message\n")
			_, _ = io.WriteString(w, "data: {\"jsonrpc\":\"2.0\",\"id\":2,\"result\":{\"content\":[{\"type\":\"text\",\"text\":\"Title: London 2026: Squad lists for every country revealed!\\nURL: https://www.tabletennisengland.co.uk/news/2026/london-2026-squad-lists-for-every-country-revealed/\\nPublished: 2026-04-12T17:25:28.000Z\\nHighlights:\\nChina: LIANG Jingkun, LIN Shidong, WANG Chuqin\\n\\n---\\n\\nTitle: Here are the draws for ITTF World Team Table Tennis Championships Finals London 2026\\nURL: https://news.tabletennis.tv/news/draws-ittf-world-team-championships-london\\nPublished: N/A\\nHighlights:\\nThe ceremony decided the groups for both the Men’s and Women’s events.\"}]}}\n\n")
		default:
			t.Fatalf("unexpected method %q", method)
		}
	}))
	defer server.Close()

	client := &ExaMCPClient{
		Endpoint:   server.URL,
		HTTPClient: server.Client(),
	}

	grounding, err := client.Search(context.Background(), "2026 伦敦世乒赛 男团名单")
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(methods) != 3 {
		t.Fatalf("methods = %v, want 3 lifecycle calls", methods)
	}
	if methods[0] != "initialize" || methods[1] != "notifications/initialized" || methods[2] != "tools/call" {
		t.Fatalf("methods = %v, want initialize -> notifications/initialized -> tools/call", methods)
	}
	if grounding.Query != "2026 伦敦世乒赛 男团名单" {
		t.Fatalf("grounding.Query = %q", grounding.Query)
	}
	if len(grounding.Sources) != 2 {
		t.Fatalf("len(grounding.Sources) = %d, want 2", len(grounding.Sources))
	}
	if grounding.Sources[0].Title != "London 2026: Squad lists for every country revealed!" {
		t.Fatalf("first source title = %q", grounding.Sources[0].Title)
	}
	if grounding.Sources[0].URL != "https://www.tabletennisengland.co.uk/news/2026/london-2026-squad-lists-for-every-country-revealed/" {
		t.Fatalf("first source url = %q", grounding.Sources[0].URL)
	}
	if !strings.Contains(grounding.Evidence, "China: LIANG Jingkun") {
		t.Fatalf("grounding.Evidence = %q, want parsed highlights", grounding.Evidence)
	}
}

func TestPrepareSearchGroundingCallsSearcherAndPrependsEvidence(t *testing.T) {
	t.Helper()

	searcher := &stubAISearcher{
		grounding: &AISearchGrounding{
			Query:    "2026 伦敦世乒赛 男团名单",
			Evidence: "来源 1: Table Tennis England\n- URL: https://example.com/1\n- 摘要: 中国队名单已公布。",
			Sources: []AISearchSource{
				{
					Title: "Table Tennis England",
					URL:   "https://example.com/1",
				},
			},
		},
	}

	svc := &AIService{searcher: searcher}
	original := []ChatMessage{{Role: "user", Content: "2026 伦敦世乒赛男团名单？"}}

	messages, footer, err := svc.prepareSearchGrounding(context.Background(), original, "2026 伦敦世乒赛 男团名单", true)
	if err != nil {
		t.Fatalf("prepareSearchGrounding() error = %v", err)
	}

	if searcher.called != 1 {
		t.Fatalf("searcher called = %d, want 1", searcher.called)
	}
	if len(messages) != 2 {
		t.Fatalf("len(messages) = %d, want 2", len(messages))
	}
	if messages[0].Role != "system" {
		t.Fatalf("messages[0].Role = %s, want system", messages[0].Role)
	}

	systemContent, ok := messages[0].Content.(string)
	if !ok {
		t.Fatalf("messages[0].Content type = %T, want string", messages[0].Content)
	}
	if !strings.Contains(systemContent, "Exa 联网搜索") {
		t.Fatalf("systemContent = %q, want Exa grounding guidance", systemContent)
	}
	if !strings.Contains(systemContent, "https://example.com/1") {
		t.Fatalf("systemContent = %q, want source URL", systemContent)
	}
	if !strings.Contains(footer, "[Table Tennis England](https://example.com/1)") {
		t.Fatalf("footer = %q, want markdown source link", footer)
	}
}

func TestBuildChatCompletionRequestDoesNotForwardVendorSearchFlags(t *testing.T) {
	req := buildChatCompletionRequest(
		"qwen3-max",
		[]ChatMessage{{Role: "user", Content: "今天上海天气怎么样？"}},
		true,
		true,
		true,
	)

	if req.EnableSearch {
		t.Fatal("expected provider native enable_search to stay false when using Exa grounding")
	}
	if req.SearchOptions != nil {
		t.Fatalf("search options = %#v, want nil", req.SearchOptions)
	}
	if req.ExtraBody != nil {
		t.Fatalf("extra body = %#v, want nil", req.ExtraBody)
	}
	if len(req.Messages) != 2 {
		t.Fatalf("messages len = %d, want 2", len(req.Messages))
	}
}

func TestPrepareSearchGroundingReturnsErrorWhenSearcherFails(t *testing.T) {
	svc := &AIService{
		searcher: &stubAISearcher{
			err: errors.New("upstream timeout"),
		},
	}

	_, _, err := svc.prepareSearchGrounding(context.Background(), []ChatMessage{{Role: "user", Content: "test"}}, "2026 伦敦世乒赛 男团名单", true)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "联网搜索失败") {
		t.Fatalf("error = %v, want search failure prefix", err)
	}
}

type stubAISearcher struct {
	called    int
	grounding *AISearchGrounding
	err       error
}

func (s *stubAISearcher) Search(_ context.Context, _ string) (*AISearchGrounding, error) {
	s.called++
	return s.grounding, s.err
}
