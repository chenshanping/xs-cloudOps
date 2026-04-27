package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"server/model"

	"net/http"
	"net/http/httptest"
)

func TestAIServiceReadFileContentReadsLocalTextFile(t *testing.T) {
	db := setupStorageServiceTestDB(t)

	dir := t.TempDir()
	filePath := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(filePath, []byte("hello local file"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	configJSON, err := json.Marshal(model.LocalConfig{BasePath: dir})
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}

	configs := []model.SysConfig{
		{Key: StorageTypeConfigKey, Value: string(model.StorageTypeLocal)},
		{Key: StorageConfigKey(model.StorageTypeLocal), Value: string(configJSON)},
	}
	if err := db.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}

	content, err := AI.readFileContent(model.SysFile{
		Name:        "notes.txt",
		Path:        "notes.txt",
		StorageType: string(model.StorageTypeLocal),
	})
	if err != nil {
		t.Fatalf("readFileContent returned error: %v", err)
	}
	if content != "hello local file" {
		t.Fatalf("unexpected content: %q", content)
	}
}

func TestAIServiceHttpGetFileContentReadsRemoteContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("remote-content"))
	}))
	defer server.Close()

	content, err := AI.httpGetFileContent(server.URL)
	if err != nil {
		t.Fatalf("httpGetFileContent returned error: %v", err)
	}
	if string(content) != "remote-content" {
		t.Fatalf("unexpected content: %q", string(content))
	}
}

func TestAIServiceLocalFileToBase64ConvertsImage(t *testing.T) {
	db := setupStorageServiceTestDB(t)

	dir := t.TempDir()
	filePath := filepath.Join(dir, "avatar.png")
	if err := os.WriteFile(filePath, []byte("png-bytes"), 0o600); err != nil {
		t.Fatalf("write temp image: %v", err)
	}

	configJSON, err := json.Marshal(model.LocalConfig{BasePath: dir})
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}

	configs := []model.SysConfig{
		{Key: StorageTypeConfigKey, Value: string(model.StorageTypeLocal)},
		{Key: StorageConfigKey(model.StorageTypeLocal), Value: string(configJSON)},
	}
	if err := db.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}

	dataURL, err := AI.localFileToBase64(model.SysFile{
		Name:        "avatar.png",
		Path:        "avatar.png",
		MimeType:    "image/png",
		StorageType: string(model.StorageTypeLocal),
	})
	if err != nil {
		t.Fatalf("localFileToBase64 returned error: %v", err)
	}
	if !strings.HasPrefix(dataURL, "data:image/png;base64,") {
		t.Fatalf("unexpected data url: %s", dataURL)
	}
}
