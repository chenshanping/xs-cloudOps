package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
)

func setupStorageServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysConfig{}, &model.SysFile{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})

	return db
}

func TestStorageServiceGetStorageByTypeUsesTypeSpecificConfig(t *testing.T) {
	db := setupStorageServiceTestDB(t)

	configs := []model.SysConfig{
		{Key: StorageTypeConfigKey, Value: string(model.StorageTypeMinio)},
		{Key: StorageConfigKey(model.StorageTypeLocal), Value: `{"base_path":"local-dir","base_url":"/api/v1/upload"}`},
		{Key: StorageConfigKey(model.StorageTypeMinio), Value: `{"endpoint":"127.0.0.1:9000","access_key_id":"minio","secret_access_key":"secret","bucket_name":"files","use_ssl":false}`},
	}
	if err := db.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}

	localStorage, err := Storage.GetStorageByType(model.StorageTypeLocal)
	if err != nil {
		t.Fatalf("GetStorageByType(local) error: %v", err)
	}
	if localStorage.Type != model.StorageTypeLocal {
		t.Fatalf("local storage type = %s, want %s", localStorage.Type, model.StorageTypeLocal)
	}
	if localStorage.Config != `{"base_path":"local-dir","base_url":"/api/v1/upload"}` {
		t.Fatalf("local storage config = %s", localStorage.Config)
	}

	minioStorage, err := Storage.GetStorageByType(model.StorageTypeMinio)
	if err != nil {
		t.Fatalf("GetStorageByType(minio) error: %v", err)
	}
	if minioStorage.Type != model.StorageTypeMinio {
		t.Fatalf("minio storage type = %s, want %s", minioStorage.Type, model.StorageTypeMinio)
	}
	if minioStorage.Config != `{"endpoint":"127.0.0.1:9000","access_key_id":"minio","secret_access_key":"secret","bucket_name":"files","use_ssl":false}` {
		t.Fatalf("minio storage config = %s", minioStorage.Config)
	}
}

func TestAIServiceReadFileContentUsesTypeSpecificSystemConfig(t *testing.T) {
	db := setupStorageServiceTestDB(t)

	dir := t.TempDir()
	filePath := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(filePath, []byte("hello type config"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	configs := []model.SysConfig{
		{Key: StorageTypeConfigKey, Value: string(model.StorageTypeMinio)},
		{Key: StorageConfigKey(model.StorageTypeLocal), Value: `{"base_path":"` + filepath.ToSlash(dir) + `","base_url":"/api/v1/upload"}`},
		{Key: StorageConfigKey(model.StorageTypeMinio), Value: `{"endpoint":"127.0.0.1:9000","access_key_id":"minio","secret_access_key":"secret","bucket_name":"files","use_ssl":false}`},
	}
	if err := db.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}

	content, err := AI.readFileContent(model.SysFile{
		Name:        "notes.txt",
		Path:        "notes.txt",
		URL:         "/api/v1/upload/notes.txt",
		StorageType: string(model.StorageTypeLocal),
	})
	if err != nil {
		t.Fatalf("readFileContent returned error: %v", err)
	}
	if content != "hello type config" {
		t.Fatalf("unexpected content: %q", content)
	}
}
