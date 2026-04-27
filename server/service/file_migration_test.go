package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"server/global"
	"server/model"
	modelrequest "server/model/request"
	serviceoss "server/service/oss"
)

type fakeObjectStorageClient struct {
	files            map[string][]byte
	uploadStarted    chan struct{}
	uploadRelease    chan struct{}
	uploadStartedOnce sync.Once
}

func newFakeObjectStorageClient() *fakeObjectStorageClient {
	return &fakeObjectStorageClient{files: make(map[string][]byte)}
}

func (c *fakeObjectStorageClient) Upload(ctx context.Context, key string, reader io.Reader, size int64) error {
	if c.uploadStarted != nil {
		c.uploadStartedOnce.Do(func() {
			close(c.uploadStarted)
		})
	}
	if c.uploadRelease != nil {
		select {
		case <-c.uploadRelease:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	c.files[key] = data
	return nil
}

func (c *fakeObjectStorageClient) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	data, ok := c.files[key]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (c *fakeObjectStorageClient) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := c.files[key]
	return ok, nil
}

func (c *fakeObjectStorageClient) Delete(ctx context.Context, key string) error {
	delete(c.files, key)
	return nil
}

func (c *fakeObjectStorageClient) GetURL(key string) string {
	return "https://fake-minio/" + key
}

func (c *fakeObjectStorageClient) GetSignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	return c.GetURL(key), nil
}

func (c *fakeObjectStorageClient) GetUploadCredential(ctx context.Context, key string, expires time.Duration) (*serviceoss.UploadCredential, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *fakeObjectStorageClient) InitMultipartUpload(ctx context.Context, key string) (*serviceoss.MultipartUpload, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *fakeObjectStorageClient) GetMultipartUploadURL(ctx context.Context, uploadID, key string, partNumber int, expires time.Duration) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *fakeObjectStorageClient) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []serviceoss.Part) error {
	return fmt.Errorf("not implemented")
}

func (c *fakeObjectStorageClient) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	return fmt.Errorf("not implemented")
}

func (c *fakeObjectStorageClient) ListParts(ctx context.Context, key, uploadID string) ([]serviceoss.Part, error) {
	return nil, fmt.Errorf("not implemented")
}

func seedFileMigrationConfigs(t *testing.T, dbDir string) {
	t.Helper()

	configJSON, err := json.Marshal(model.LocalConfig{BasePath: dbDir, BaseURL: "/api/v1/upload"})
	if err != nil {
		t.Fatalf("marshal local config: %v", err)
	}

	configs := []model.SysConfig{
		{Key: StorageTypeConfigKey, Value: string(model.StorageTypeLocal)},
		{Key: StorageConfigKey(model.StorageTypeLocal), Value: string(configJSON)},
		{Key: StorageConfigKey(model.StorageTypeMinio), Value: "{}"},
	}
	if err := global.DB.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}
}

func TestPreviewFileMigrationMarksPendingAndSkippedItems(t *testing.T) {
	db := setupFileServiceTestDB(t)
	dir := t.TempDir()
	seedFileMigrationConfigs(t, dir)

	fakeMinio := newFakeObjectStorageClient()
	resetBuilder := serviceoss.RegisterClientBuilderForTest(model.StorageTypeMinio, func(storage *model.StorageProfile) (serviceoss.Client, error) {
		return fakeMinio, nil
	})
	t.Cleanup(resetBuilder)

	if err := os.WriteFile(filepath.Join(dir, "pending.txt"), []byte("pending"), 0o600); err != nil {
		t.Fatalf("write pending source file: %v", err)
	}

	files := []model.SysFile{
		{
			Name:        "pending.txt",
			Path:        "pending.txt",
			URL:         "/api/v1/upload/pending.txt",
			StorageType: string(model.StorageTypeLocal),
			Status:      1,
		},
		{
			Name:        "already-minio.txt",
			Path:        "already-minio.txt",
			URL:         "https://fake-minio/already-minio.txt",
			StorageType: string(model.StorageTypeMinio),
			Status:      1,
		},
		{
			Name:        "conflict.txt",
			Path:        "conflict.txt",
			URL:         "/api/v1/upload/conflict.txt",
			StorageType: string(model.StorageTypeLocal),
			Status:      1,
		},
		{
			Name:        "conflict-target.txt",
			Path:        "conflict.txt",
			URL:         "https://fake-minio/conflict.txt",
			StorageType: string(model.StorageTypeMinio),
			Status:      1,
		},
	}
	if err := db.Create(&files).Error; err != nil {
		t.Fatalf("create files: %v", err)
	}

	result, err := File.PreviewFileMigration(modelrequest.FileMigrationRequest{
		Scope:             "all",
		SourceStorageType: string(model.StorageTypeLocal),
		TargetStorageType: string(model.StorageTypeMinio),
	})
	if err != nil {
		t.Fatalf("PreviewFileMigration error: %v", err)
	}

	if result.TotalCount != 2 {
		t.Fatalf("total count = %d, want 2", result.TotalCount)
	}
	if result.PendingCount != 1 {
		t.Fatalf("pending count = %d, want 1", result.PendingCount)
	}
	if result.ConflictCount != 1 {
		t.Fatalf("conflict count = %d, want 1", result.ConflictCount)
	}
}

func TestPreviewFileMigrationIncludesPrecheckStats(t *testing.T) {
	db := setupFileServiceTestDB(t)
	dir := t.TempDir()
	seedFileMigrationConfigs(t, dir)

	fakeMinio := newFakeObjectStorageClient()
	fakeMinio.files["conflict.txt"] = []byte("target-conflict")
	resetBuilder := serviceoss.RegisterClientBuilderForTest(model.StorageTypeMinio, func(storage *model.StorageProfile) (serviceoss.Client, error) {
		return fakeMinio, nil
	})
	t.Cleanup(resetBuilder)

	if err := os.WriteFile(filepath.Join(dir, "pending.txt"), []byte("pending"), 0o600); err != nil {
		t.Fatalf("write pending source file: %v", err)
	}
	conflictContent := []byte("source-conflict")
	conflictSum := md5.Sum(conflictContent)
	if err := os.WriteFile(filepath.Join(dir, "conflict.txt"), conflictContent, 0o600); err != nil {
		t.Fatalf("write conflict source file: %v", err)
	}

	files := []model.SysFile{
		{
			Name:        "pending.txt",
			Path:        "pending.txt",
			URL:         "/api/v1/upload/pending.txt",
			Size:        int64(len("pending")),
			StorageType: string(model.StorageTypeLocal),
			Status:      1,
		},
		{
			Name:        "missing.txt",
			Path:        "missing.txt",
			URL:         "/api/v1/upload/missing.txt",
			Size:        int64(len("missing")),
			StorageType: string(model.StorageTypeLocal),
			Status:      1,
		},
		{
			Name:        "conflict.txt",
			Path:        "conflict.txt",
			URL:         "/api/v1/upload/conflict.txt",
			MD5:         fmt.Sprintf("%x", conflictSum),
			Size:        int64(len(conflictContent)),
			StorageType: string(model.StorageTypeLocal),
			Status:      1,
		},
	}
	if err := db.Create(&files).Error; err != nil {
		t.Fatalf("create files: %v", err)
	}

	result, err := File.PreviewFileMigration(modelrequest.FileMigrationRequest{
		Scope:             "all",
		SourceStorageType: string(model.StorageTypeLocal),
		TargetStorageType: string(model.StorageTypeMinio),
	})
	if err != nil {
		t.Fatalf("PreviewFileMigration error: %v", err)
	}

	if result.TotalCount != 3 {
		t.Fatalf("total count = %d, want 3", result.TotalCount)
	}
	if result.TotalSize != int64(len("pending")+len("missing"))+int64(len(conflictContent)) {
		t.Fatalf("total size = %d", result.TotalSize)
	}
	if result.PendingCount != 1 {
		t.Fatalf("pending count = %d, want 1", result.PendingCount)
	}
	if result.ConflictCount != 1 {
		t.Fatalf("conflict count = %d, want 1", result.ConflictCount)
	}
	if result.MissingSourceCount != 1 {
		t.Fatalf("missing source count = %d, want 1", result.MissingSourceCount)
	}
}

func TestPreviewFileMigrationFilterScopeUsesSearchConditions(t *testing.T) {
	db := setupFileServiceTestDB(t)
	dir := t.TempDir()
	seedFileMigrationConfigs(t, dir)

	fakeMinio := newFakeObjectStorageClient()
	resetBuilder := serviceoss.RegisterClientBuilderForTest(model.StorageTypeMinio, func(storage *model.StorageProfile) (serviceoss.Client, error) {
		return fakeMinio, nil
	})
	t.Cleanup(resetBuilder)

	if err := os.WriteFile(filepath.Join(dir, "report-a.pdf"), []byte("report-a"), 0o600); err != nil {
		t.Fatalf("write filtered source file: %v", err)
	}

	files := []model.SysFile{
		{Name: "report-a.pdf", Path: "report-a.pdf", URL: "/api/v1/upload/report-a.pdf", Ext: "pdf", StorageType: string(model.StorageTypeLocal), Status: 1},
		{Name: "report-b.doc", Path: "report-b.doc", URL: "/api/v1/upload/report-b.doc", Ext: "doc", StorageType: string(model.StorageTypeLocal), Status: 1},
		{Name: "image-a.pdf", Path: "image-a.pdf", URL: "/api/v1/upload/image-a.pdf", Ext: "pdf", StorageType: string(model.StorageTypeLocal), Status: 1},
	}
	if err := db.Create(&files).Error; err != nil {
		t.Fatalf("create files: %v", err)
	}

	result, err := File.PreviewFileMigration(modelrequest.FileMigrationRequest{
		Scope:             "filter",
		SourceStorageType: string(model.StorageTypeLocal),
		TargetStorageType: string(model.StorageTypeMinio),
		Filters: modelrequest.FileMigrationFilters{
			Name: "report",
			Ext:  "pdf",
		},
	})
	if err != nil {
		t.Fatalf("PreviewFileMigration filter error: %v", err)
	}

	if result.TotalCount != 1 {
		t.Fatalf("total count = %d, want 1", result.TotalCount)
	}
	if result.PendingCount != 1 {
		t.Fatalf("pending count = %d, want 1", result.PendingCount)
	}
}

func TestExecuteFileMigrationMovesLocalFileToTargetStorage(t *testing.T) {
	db := setupFileServiceTestDB(t)
	dir := t.TempDir()
	seedFileMigrationConfigs(t, dir)

	fakeMinio := newFakeObjectStorageClient()
	resetBuilder := serviceoss.RegisterClientBuilderForTest(model.StorageTypeMinio, func(storage *model.StorageProfile) (serviceoss.Client, error) {
		return fakeMinio, nil
	})
	t.Cleanup(resetBuilder)

	content := []byte("migration-demo")
	sum := md5.Sum(content)
	filePath := filepath.Join(dir, "demo.txt")
	if err := os.WriteFile(filePath, content, 0o600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	file := model.SysFile{
		Name:        "demo.txt",
		Path:        "demo.txt",
		URL:         "/api/v1/upload/demo.txt",
		MD5:         fmt.Sprintf("%x", sum),
		Size:        int64(len(content)),
		StorageType: string(model.StorageTypeLocal),
		Status:      1,
	}
	if err := db.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	result, err := File.ExecuteFileMigration(modelrequest.FileMigrationRequest{
		Scope:             "all",
		SourceStorageType: string(model.StorageTypeLocal),
		TargetStorageType: string(model.StorageTypeMinio),
	})
	if err != nil {
		t.Fatalf("ExecuteFileMigration error: %v", err)
	}

	if result.MigratedCount != 1 {
		t.Fatalf("migrated count = %d, want 1", result.MigratedCount)
	}
	if result.FailedCount != 0 {
		t.Fatalf("failed count = %d, want 0", result.FailedCount)
	}
	if result.WarningCount != 0 {
		t.Fatalf("warning count = %d, want 0", result.WarningCount)
	}

	var updated model.SysFile
	if err := db.First(&updated, file.ID).Error; err != nil {
		t.Fatalf("query updated file: %v", err)
	}
	if updated.StorageType != string(model.StorageTypeMinio) {
		t.Fatalf("storage type = %s, want %s", updated.StorageType, model.StorageTypeMinio)
	}
	if updated.URL != "https://fake-minio/demo.txt" {
		t.Fatalf("url = %s", updated.URL)
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("source file still exists, stat err = %v", err)
	}

	migratedContent, ok := fakeMinio.files["demo.txt"]
	if !ok {
		t.Fatalf("target object missing")
	}
	if string(migratedContent) != string(content) {
		t.Fatalf("target object content mismatch")
	}
}

func TestStartFileMigrationTaskTracksProgressAndBlocksConcurrentStart(t *testing.T) {
	db := setupFileServiceTestDB(t)
	dir := t.TempDir()
	seedFileMigrationConfigs(t, dir)

	fakeMinio := newFakeObjectStorageClient()
	fakeMinio.uploadStarted = make(chan struct{})
	fakeMinio.uploadRelease = make(chan struct{})
	resetBuilder := serviceoss.RegisterClientBuilderForTest(model.StorageTypeMinio, func(storage *model.StorageProfile) (serviceoss.Client, error) {
		return fakeMinio, nil
	})
	t.Cleanup(resetBuilder)

	content := []byte("migration-task-demo")
	sum := md5.Sum(content)
	filePath := filepath.Join(dir, "task-demo.txt")
	if err := os.WriteFile(filePath, content, 0o600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	file := model.SysFile{
		Name:        "task-demo.txt",
		Path:        "task-demo.txt",
		URL:         "/api/v1/upload/task-demo.txt",
		MD5:         fmt.Sprintf("%x", sum),
		Size:        int64(len(content)),
		StorageType: string(model.StorageTypeLocal),
		Status:      1,
	}
	if err := db.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	task, err := File.StartFileMigrationTask(modelrequest.FileMigrationRequest{
		Scope:             "all",
		SourceStorageType: string(model.StorageTypeLocal),
		TargetStorageType: string(model.StorageTypeMinio),
	})
	if err != nil {
		t.Fatalf("StartFileMigrationTask error: %v", err)
	}
	if task == nil || task.TaskID == "" {
		t.Fatalf("expected task id")
	}

	select {
	case <-fakeMinio.uploadStarted:
	case <-time.After(2 * time.Second):
		t.Fatalf("upload did not start in time")
	}

	if _, err := File.StartFileMigrationTask(modelrequest.FileMigrationRequest{
		Scope:             "all",
		SourceStorageType: string(model.StorageTypeLocal),
		TargetStorageType: string(model.StorageTypeMinio),
	}); err == nil {
		t.Fatalf("expected concurrent task start to fail")
	}

	close(fakeMinio.uploadRelease)

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		status := File.GetCurrentFileMigrationTask()
		if status != nil && status.Status == "SUCCESS" {
			if status.MigratedCount != 1 {
				t.Fatalf("migrated count = %d, want 1", status.MigratedCount)
			}
			if status.ProcessedCount != 1 {
				t.Fatalf("processed count = %d, want 1", status.ProcessedCount)
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("task did not finish in time")
}
