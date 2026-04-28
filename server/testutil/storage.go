package testutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	serviceoss "server/service/oss"
)

// FakeObjectStorageClient implements serviceoss.Client with in-memory storage
// for use in integration tests. It supports synchronization hooks via
// UploadStarted and UploadRelease channels.
type FakeObjectStorageClient struct {
	Files              map[string][]byte
	UploadStarted      chan struct{}
	UploadRelease      chan struct{}
	uploadStartedOnce  sync.Once
}

// NewFakeObjectStorageClient creates a new FakeObjectStorageClient.
func NewFakeObjectStorageClient() *FakeObjectStorageClient {
	return &FakeObjectStorageClient{Files: make(map[string][]byte)}
}

func (c *FakeObjectStorageClient) Upload(ctx context.Context, key string, reader io.Reader, size int64) error {
	if c.UploadStarted != nil {
		c.uploadStartedOnce.Do(func() {
			close(c.UploadStarted)
		})
	}
	if c.UploadRelease != nil {
		select {
		case <-c.UploadRelease:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	c.Files[key] = data
	return nil
}

func (c *FakeObjectStorageClient) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	data, ok := c.Files[key]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (c *FakeObjectStorageClient) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := c.Files[key]
	return ok, nil
}

func (c *FakeObjectStorageClient) Delete(ctx context.Context, key string) error {
	delete(c.Files, key)
	return nil
}

func (c *FakeObjectStorageClient) GetURL(key string) string {
	return "https://fake-minio/" + key
}

func (c *FakeObjectStorageClient) GetSignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	return c.GetURL(key), nil
}

func (c *FakeObjectStorageClient) GetUploadCredential(ctx context.Context, key string, expires time.Duration) (*serviceoss.UploadCredential, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *FakeObjectStorageClient) InitMultipartUpload(ctx context.Context, key string) (*serviceoss.MultipartUpload, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *FakeObjectStorageClient) GetMultipartUploadURL(ctx context.Context, uploadID, key string, partNumber int, expires time.Duration) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *FakeObjectStorageClient) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []serviceoss.Part) error {
	return fmt.Errorf("not implemented")
}

func (c *FakeObjectStorageClient) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	return fmt.Errorf("not implemented")
}

func (c *FakeObjectStorageClient) ListParts(ctx context.Context, key, uploadID string) ([]serviceoss.Part, error) {
	return nil, fmt.Errorf("not implemented")
}
