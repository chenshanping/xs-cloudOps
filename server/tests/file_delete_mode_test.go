package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"server/global"
	"server/model"
	. "server/service"
	"server/testutil"
)

func seedLocalFileDeleteConfig(t *testing.T, dbDir string, deleteMode string) {
	t.Helper()

	configJSON, err := json.Marshal(model.LocalConfig{BasePath: dbDir, BaseURL: "/api/v1/upload"})
	if err != nil {
		t.Fatalf("marshal local config: %v", err)
	}

	configs := []model.SysConfig{
		{Key: FileDeleteModeConfigKey, Value: deleteMode},
		{Key: StorageTypeConfigKey, Value: string(model.StorageTypeLocal)},
		{Key: StorageConfigKey(model.StorageTypeLocal), Value: string(configJSON)},
	}
	if err := global.DB.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}
}

func TestDeleteFileLogicalModeKeepsPhysicalFile(t *testing.T) {
	testutil.SetupFileServiceTestDB(t)
	dir := t.TempDir()
	seedLocalFileDeleteConfig(t, dir, FileDeleteModeLogical)

	fullPath := filepath.Join(dir, "demo.txt")
	if err := os.WriteFile(fullPath, []byte("demo"), 0o600); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	file := model.SysFile{
		Name:        "demo.txt",
		Path:        "demo.txt",
		URL:         "/api/v1/upload/demo.txt",
		StorageType: string(model.StorageTypeLocal),
		Status:      1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	if err := File.DeleteFile(file.ID); err != nil {
		t.Fatalf("DeleteFile logical error: %v", err)
	}

	var updated model.SysFile
	if err := global.DB.First(&updated, file.ID).Error; err != nil {
		t.Fatalf("query updated file: %v", err)
	}
	if updated.Status != 0 {
		t.Fatalf("logical delete status = %d, want 0", updated.Status)
	}
	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("logical delete removed physical file: %v", err)
	}
}

func TestDeleteFilePhysicalModeRemovesRecordAndPhysicalFile(t *testing.T) {
	testutil.SetupFileServiceTestDB(t)
	dir := t.TempDir()
	seedLocalFileDeleteConfig(t, dir, FileDeleteModePhysical)

	fullPath := filepath.Join(dir, "demo.txt")
	if err := os.WriteFile(fullPath, []byte("demo"), 0o600); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	file := model.SysFile{
		Name:        "demo.txt",
		Path:        "demo.txt",
		URL:         "/api/v1/upload/demo.txt",
		StorageType: string(model.StorageTypeLocal),
		Status:      1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	if err := File.DeleteFile(file.ID); err != nil {
		t.Fatalf("DeleteFile physical error: %v", err)
	}

	var count int64
	if err := global.DB.Unscoped().Model(&model.SysFile{}).Where("id = ?", file.ID).Count(&count).Error; err != nil {
		t.Fatalf("count deleted file: %v", err)
	}
	if count != 0 {
		t.Fatalf("physical delete retained db record count = %d", count)
	}
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		t.Fatalf("physical delete did not remove physical file, stat err = %v", err)
	}
}

func TestDeleteFilePhysicalModeAllowsMissingPhysicalFile(t *testing.T) {
	testutil.SetupFileServiceTestDB(t)
	dir := t.TempDir()
	seedLocalFileDeleteConfig(t, dir, FileDeleteModePhysical)

	file := model.SysFile{
		Name:        "missing.txt",
		Path:        "missing.txt",
		URL:         "/api/v1/upload/missing.txt",
		StorageType: string(model.StorageTypeLocal),
		Status:      1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	if err := File.DeleteFile(file.ID); err != nil {
		t.Fatalf("DeleteFile physical with missing file error: %v", err)
	}

	var count int64
	if err := global.DB.Unscoped().Model(&model.SysFile{}).Where("id = ?", file.ID).Count(&count).Error; err != nil {
		t.Fatalf("count deleted file: %v", err)
	}
	if count != 0 {
		t.Fatalf("physical delete retained db record count = %d", count)
	}
}

func TestDeleteFileRejectsReferencedAvatarFile(t *testing.T) {
	testutil.SetupFileServiceTestDB(t)
	dir := t.TempDir()
	seedLocalFileDeleteConfig(t, dir, FileDeleteModePhysical)

	file := model.SysFile{
		Name:        "avatar.png",
		Path:        "avatar.png",
		URL:         "/api/v1/upload/avatar.png",
		StorageType: string(model.StorageTypeLocal),
		Status:      1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	user := model.SysUser{
		Username:     "avatar-user",
		Password:     "pwd",
		Nickname:     "Avatar User",
		Status:       1,
		AvatarFileID: file.ID,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	err := File.DeleteFile(file.ID)
	if err == nil {
		t.Fatalf("expected referenced avatar file delete to fail")
	}
	if !strings.Contains(err.Error(), "文件正在被引用") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteFileRejectsReferencedAIMessageFile(t *testing.T) {
	testutil.SetupFileServiceTestDB(t)
	dir := t.TempDir()
	seedLocalFileDeleteConfig(t, dir, FileDeleteModePhysical)

	file := model.SysFile{
		Name:        "doc.txt",
		Path:        "doc.txt",
		URL:         "/api/v1/upload/doc.txt",
		StorageType: string(model.StorageTypeLocal),
		Status:      1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	messageFileIDs, err := json.Marshal([]uint{file.ID, file.ID + 1})
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

	err = File.DeleteFile(file.ID)
	if err == nil {
		t.Fatalf("expected referenced ai file delete to fail")
	}
	if !strings.Contains(err.Error(), "AI对话附件") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteFilePhysicalModeRemovesChunkRecordsAndTempDirectory(t *testing.T) {
	testutil.SetupFileServiceTestDB(t)
	dir := t.TempDir()
	seedLocalFileDeleteConfig(t, dir, FileDeleteModePhysical)

	fullPath := filepath.Join(dir, "demo.txt")
	if err := os.WriteFile(fullPath, []byte("demo"), 0o600); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	chunkDir := filepath.Join(dir, ".chunks", "upload-1")
	if err := os.MkdirAll(chunkDir, 0o755); err != nil {
		t.Fatalf("mkdir chunk dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(chunkDir, "1"), []byte("chunk"), 0o600); err != nil {
		t.Fatalf("write chunk file: %v", err)
	}

	file := model.SysFile{
		Name:        "demo.txt",
		Path:        "demo.txt",
		URL:         "/api/v1/upload/demo.txt",
		MD5:         "demo-md5",
		StorageType: string(model.StorageTypeLocal),
		Status:      1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	chunks := []model.SysFileChunk{
		{
			UploadID:    "upload-1",
			FileHash:    "demo-md5",
			ChunkIndex:  1,
			ChunkHash:   "chunk-1",
			StorageType: string(model.StorageTypeLocal),
			StoragePath: file.Path,
			Status:      0,
		},
		{
			UploadID:    "upload-1",
			FileHash:    "demo-md5",
			ChunkIndex:  2,
			ChunkHash:   "chunk-2",
			StorageType: string(model.StorageTypeLocal),
			StoragePath: file.Path,
			Status:      0,
		},
	}
	if err := global.DB.Create(&chunks).Error; err != nil {
		t.Fatalf("create chunks: %v", err)
	}

	if err := File.DeleteFile(file.ID); err != nil {
		t.Fatalf("DeleteFile physical error: %v", err)
	}

	var chunkCount int64
	if err := global.DB.Model(&model.SysFileChunk{}).Where("file_hash = ?", file.MD5).Count(&chunkCount).Error; err != nil {
		t.Fatalf("count chunk rows: %v", err)
	}
	if chunkCount != 0 {
		t.Fatalf("chunk rows still exist: %d", chunkCount)
	}
	if _, err := os.Stat(chunkDir); !os.IsNotExist(err) {
		t.Fatalf("chunk directory still exists, stat err = %v", err)
	}
}
