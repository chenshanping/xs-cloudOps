package v1

import "testing"

func TestBatchConversationDeleteFailureMessage(t *testing.T) {
	t.Run("returns fallback when list is empty", func(t *testing.T) {
		if got := batchConversationDeleteFailureMessage(nil); got != "删除失败" {
			t.Fatalf("expected fallback message, got %q", got)
		}
	})

	t.Run("returns single failure message directly", func(t *testing.T) {
		want := "ID 3: 对话不存在"
		if got := batchConversationDeleteFailureMessage([]string{want}); got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})

	t.Run("joins multiple failure messages", func(t *testing.T) {
		want := "ID 3: 对话不存在；ID 5: 无权限删除该对话"
		got := batchConversationDeleteFailureMessage([]string{
			"ID 3: 对话不存在",
			"ID 5: 无权限删除该对话",
		})
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})
}
