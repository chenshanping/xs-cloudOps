package ai

import (
	"strings"
	"testing"
)

func TestBuildChatCompletionRequestAddsGroundingGuidanceWhenSearchEnabled(t *testing.T) {
	req := buildChatCompletionRequest(
		"qwen3-max",
		[]ChatMessage{{Role: "user", Content: "今天上海天气怎么样？"}},
		true,
		true,
		true,
	)

	if req.EnableSearch {
		t.Fatal("expected provider-native enable_search to stay false")
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
	if req.Messages[0].Role != "system" {
		t.Fatalf("first message role = %s, want system", req.Messages[0].Role)
	}

	systemContent, ok := req.Messages[0].Content.(string)
	if !ok {
		t.Fatalf("system content type = %T, want string", req.Messages[0].Content)
	}
	if !strings.Contains(systemContent, "联网搜索") {
		t.Fatalf("system content = %q, want search guidance", systemContent)
	}
	if !strings.Contains(systemContent, "来源摘要") {
		t.Fatalf("system content = %q, want grounded-source guidance", systemContent)
	}
	if !strings.Contains(systemContent, "不要提及“知识截止日期”") {
		t.Fatalf("system content = %q, want cutoff suppression guidance", systemContent)
	}
	if !strings.Contains(systemContent, "名单") {
		t.Fatalf("system content = %q, want roster guidance", systemContent)
	}
	if !strings.Contains(systemContent, "无法从权威来源确认") {
		t.Fatalf("system content = %q, want no-guess guidance", systemContent)
	}
}

func TestBuildChatCompletionRequestLeavesMessagesUntouchedWhenSearchDisabled(t *testing.T) {
	original := []ChatMessage{{Role: "user", Content: "解释一下 Go context"}}
	req := buildChatCompletionRequest("qwen-plus", original, false, false, false)

	if req.EnableSearch {
		t.Fatal("expected enable_search to stay false")
	}
	if req.SearchOptions != nil {
		t.Fatalf("search options = %#v, want nil", req.SearchOptions)
	}
	if len(req.Messages) != 1 {
		t.Fatalf("messages len = %d, want 1", len(req.Messages))
	}
	if req.Messages[0].Role != "user" {
		t.Fatalf("message role = %s, want user", req.Messages[0].Role)
	}
}
