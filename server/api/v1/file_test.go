package v1

import (
	"errors"
	"testing"
)

func TestFileDeleteFailureMessage(t *testing.T) {
	t.Run("returns fallback when error is nil", func(t *testing.T) {
		if got := fileDeleteFailureMessage(nil); got != "删除失败" {
			t.Fatalf("expected fallback message, got %q", got)
		}
	})

	t.Run("returns concrete error message", func(t *testing.T) {
		want := "文件正在被引用：用户头像正在使用，无法删除"
		if got := fileDeleteFailureMessage(errors.New(want)); got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})
}

func TestBatchFileDeleteFailureMessage(t *testing.T) {
	t.Run("returns fallback when list is empty", func(t *testing.T) {
		if got := batchFileDeleteFailureMessage(nil); got != "删除失败" {
			t.Fatalf("expected fallback message, got %q", got)
		}
	})

	t.Run("returns first concrete message when only one exists", func(t *testing.T) {
		want := "ID 7: 文件正在被引用：用户头像正在使用，无法删除"
		if got := batchFileDeleteFailureMessage([]string{want}); got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})

	t.Run("joins multiple failure messages", func(t *testing.T) {
		want := "ID 7: 文件正在被引用：用户头像正在使用，无法删除；ID 9: 文件正在被引用：AI对话附件正在使用，无法删除"
		got := batchFileDeleteFailureMessage([]string{
			"ID 7: 文件正在被引用：用户头像正在使用，无法删除",
			"ID 9: 文件正在被引用：AI对话附件正在使用，无法删除",
		})
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})
}
