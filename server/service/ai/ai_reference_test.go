package ai

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
)

func setupAIConversationTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.AIConversation{}, &model.AIMessage{}, &model.SysFileReference{}); err != nil {
		t.Fatalf("auto migrate ai models: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})
}

func TestDeleteConversationClearsMessageFileRefs(t *testing.T) {
	setupAIConversationTestDB(t)

	conversation := model.AIConversation{UserID: 1, Title: "测试对话", Model: "demo"}
	if err := global.DB.Create(&conversation).Error; err != nil {
		t.Fatalf("create conversation: %v", err)
	}

	message := model.AIMessage{ConversationID: conversation.ID, Role: "user", Content: "hello", FileIDs: "[1,2]"}
	if err := global.DB.Create(&message).Error; err != nil {
		t.Fatalf("create message: %v", err)
	}
	if err := Default.registerMessageFileRefs(global.DB, message.ID, []uint{1, 2}); err != nil {
		t.Fatalf("register message file refs: %v", err)
	}

	if err := Default.DeleteConversation(conversation.ID, conversation.UserID); err != nil {
		t.Fatalf("DeleteConversation: %v", err)
	}

	var refCount int64
	if err := global.DB.Model(&model.SysFileReference{}).
		Where("ref_table = ? AND ref_id = ?", "ai_message", message.ID).
		Count(&refCount).Error; err != nil {
		t.Fatalf("count ai refs: %v", err)
	}
	if refCount != 0 {
		t.Fatalf("expected ai refs to be cleared, got %d", refCount)
	}
}

func TestClearMessagesClearsMessageFileRefs(t *testing.T) {
	setupAIConversationTestDB(t)

	conversation := model.AIConversation{UserID: 2, Title: "测试清空", Model: "demo"}
	if err := global.DB.Create(&conversation).Error; err != nil {
		t.Fatalf("create conversation: %v", err)
	}

	message := model.AIMessage{ConversationID: conversation.ID, Role: "user", Content: "hello", FileIDs: "[3]"}
	if err := global.DB.Create(&message).Error; err != nil {
		t.Fatalf("create message: %v", err)
	}
	if err := Default.registerMessageFileRefs(global.DB, message.ID, []uint{3}); err != nil {
		t.Fatalf("register message file refs: %v", err)
	}

	if err := Default.ClearMessages(conversation.ID, conversation.UserID); err != nil {
		t.Fatalf("ClearMessages: %v", err)
	}

	var refCount int64
	if err := global.DB.Model(&model.SysFileReference{}).
		Where("ref_table = ? AND ref_id = ?", "ai_message", message.ID).
		Count(&refCount).Error; err != nil {
		t.Fatalf("count ai refs: %v", err)
	}
	if refCount != 0 {
		t.Fatalf("expected ai refs to be cleared, got %d", refCount)
	}
}
