package ai

import (
	"os"
	"path/filepath"
	"testing"

	"server/model"
	"server/service/configsvc"
	"server/service/storagesvc"
	"server/testutil"
)

func TestAIServiceReadFileContentUsesTypeSpecificSystemConfig(t *testing.T) {
	db := testutil.SetupStorageServiceTestDB(t)

	dir := t.TempDir()
	filePath := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(filePath, []byte("hello type config"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	configs := []model.SysConfig{
		{Key: configsvc.StorageTypeConfigKey, Value: string(model.StorageTypeMinio)},
		{Key: storagesvc.StorageConfigKey(model.StorageTypeLocal), Value: `{"base_path":"` + filepath.ToSlash(dir) + `","base_url":"/api/v1/upload"}`},
		{Key: storagesvc.StorageConfigKey(model.StorageTypeMinio), Value: `{"endpoint":"127.0.0.1:9000","access_key_id":"minio","secret_access_key":"secret","bucket_name":"files","use_ssl":false}`},
	}
	if err := db.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}

	content, err := Default.readFileContent(model.SysFile{
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
