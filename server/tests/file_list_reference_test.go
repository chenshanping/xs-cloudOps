package tests

import (
	"encoding/json"
	"strconv"
	"testing"

	"server/global"
	"server/model"
	. "server/service"
	"server/testutil"
)

func TestGetFileListIncludesReferenceCount(t *testing.T) {
	testutil.SetupFileServiceTestDB(t)

	files := []model.SysFile{
		{Name: "unused.txt", Path: "unused.txt", URL: "/unused.txt", Status: 1},
		{Name: "message.txt", Path: "message.txt", URL: "/message.txt", Status: 1},
		{Name: "avatar.txt", Path: "avatar.txt", URL: "/avatar.txt", Status: 1},
	}
	if err := global.DB.Create(&files).Error; err != nil {
		t.Fatalf("create files: %v", err)
	}

	user := model.SysUser{
		Username:     "avatar-user",
		Password:     "pwd",
		Nickname:     "Avatar User",
		Status:       1,
		AvatarFileID: files[2].ID,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	messageFileIDs, err := json.Marshal([]uint{files[1].ID, files[2].ID})
	if err != nil {
		t.Fatalf("marshal file ids: %v", err)
	}
	aiMessage := model.AIMessage{
		ConversationID: 1,
		Role:           "user",
		Content:        "hello",
		FileIDs:        string(messageFileIDs),
	}
	if err := global.DB.Create(&aiMessage).Error; err != nil {
		t.Fatalf("create ai message: %v", err)
	}
	if err := FileReference.BackfillFileReferences(); err != nil {
		t.Fatalf("backfill file references: %v", err)
	}

	list, total, err := File.GetFileList(1, 10, "", "", nil)
	if err != nil {
		t.Fatalf("GetFileList error: %v", err)
	}
	if total != 3 {
		t.Fatalf("GetFileList total = %d, want 3", total)
	}

	countByID := make(map[uint]int64, len(list))
	for _, item := range list {
		countByID[item.ID] = item.ReferenceCount
	}

	if got := countByID[files[0].ID]; got != 0 {
		t.Fatalf("unused file reference_count = %d, want 0", got)
	}
	if got := countByID[files[1].ID]; got != 1 {
		t.Fatalf("message file reference_count = %d, want 1", got)
	}
	if got := countByID[files[2].ID]; got != 2 {
		t.Fatalf("avatar file reference_count = %d, want 2", got)
	}
}

func TestGetFileListFiltersByReferencedStatus(t *testing.T) {
	testutil.SetupFileServiceTestDB(t)

	files := []model.SysFile{
		{Name: "unused.txt", Path: "unused.txt", URL: "/unused.txt", Status: 1},
		{Name: "message.txt", Path: "message.txt", URL: "/message.txt", Status: 1},
		{Name: "avatar.txt", Path: "avatar.txt", URL: "/avatar.txt", Status: 1},
	}
	if err := global.DB.Create(&files).Error; err != nil {
		t.Fatalf("create files: %v", err)
	}

	user := model.SysUser{
		Username:     "avatar-user",
		Password:     "pwd",
		Nickname:     "Avatar User",
		Status:       1,
		AvatarFileID: files[2].ID,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	messageFileIDs, err := json.Marshal([]uint{files[1].ID})
	if err != nil {
		t.Fatalf("marshal file ids: %v", err)
	}
	aiMessage := model.AIMessage{
		ConversationID: 1,
		Role:           "user",
		Content:        "hello",
		FileIDs:        string(messageFileIDs),
	}
	if err := global.DB.Create(&aiMessage).Error; err != nil {
		t.Fatalf("create ai message: %v", err)
	}
	if err := FileReference.BackfillFileReferences(); err != nil {
		t.Fatalf("backfill file references: %v", err)
	}

	referenced := true
	referencedList, referencedTotal, err := File.GetFileList(1, 10, "", "", &referenced)
	if err != nil {
		t.Fatalf("GetFileList referenced=true error: %v", err)
	}
	if referencedTotal != 2 {
		t.Fatalf("GetFileList referenced=true total = %d, want 2", referencedTotal)
	}
	for _, item := range referencedList {
		if item.ReferenceCount <= 0 {
			t.Fatalf("referenced file %d has reference_count = %d, want > 0", item.ID, item.ReferenceCount)
		}
	}

	referenced = false
	unusedList, unusedTotal, err := File.GetFileList(1, 10, "", "", &referenced)
	if err != nil {
		t.Fatalf("GetFileList referenced=false error: %v", err)
	}
	if unusedTotal != 1 {
		t.Fatalf("GetFileList referenced=false total = %d, want 1", unusedTotal)
	}
	if len(unusedList) != 1 {
		t.Fatalf("GetFileList referenced=false len = %d, want 1", len(unusedList))
	}
	if unusedList[0].ID != files[0].ID {
		t.Fatalf("GetFileList referenced=false file id = %d, want %d", unusedList[0].ID, files[0].ID)
	}
	if unusedList[0].ReferenceCount != 0 {
		t.Fatalf("unused file reference_count = %d, want 0", unusedList[0].ReferenceCount)
	}
}

func TestGetFileListIncludesConfigFileReferences(t *testing.T) {
	testutil.SetupFileServiceTestDB(t)

	file := model.SysFile{
		Name:   "logo.png",
		Path:   "logo.png",
		URL:    "/logo.png",
		Status: 1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	configs := []model.SysConfig{
		{Name: "系统Logo文件ID", Key: "sys_logo_file_id", Value: strconv.FormatUint(uint64(file.ID), 10), ValueType: "string"},
		{Name: "系统Logo", Key: "sys_logo", Value: file.URL, ValueType: "string"},
	}
	if err := global.DB.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}
	if err := FileReference.BackfillFileReferences(); err != nil {
		t.Fatalf("backfill file references: %v", err)
	}

	list, total, err := File.GetFileList(1, 10, "", "", nil)
	if err != nil {
		t.Fatalf("GetFileList error: %v", err)
	}
	if total != 1 {
		t.Fatalf("GetFileList total = %d, want 1", total)
	}
	if len(list) != 1 {
		t.Fatalf("GetFileList len = %d, want 1", len(list))
	}
	if got := list[0].ReferenceCount; got != 1 {
		t.Fatalf("config referenced file reference_count = %d, want 1", got)
	}
}
