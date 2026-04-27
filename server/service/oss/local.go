package oss

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"server/model"
)

// LocalClient 本地存储客户端
type LocalClient struct {
	config model.LocalConfig
	secret string // 用于签名的密钥
}

// NewLocalClient 创建本地存储客户端
func NewLocalClient(configJSON string) (*LocalClient, error) {
	var config model.LocalConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析本地存储配置失败: %v", err)
	}

	// 本地存储的访问前缀后端固定为 /api/v1/upload，方便前端通过 /api/v1 代理访问图片
	config.BaseURL = "/api/v1/upload"

	// 确保目录存在
	if err := os.MkdirAll(config.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %v", err)
	}

	return &LocalClient{
		config: config,
		secret: "local-storage-secret-key",
	}, nil
}

// Upload 上传文件
func (c *LocalClient) Upload(ctx context.Context, key string, reader io.Reader, size int64) error {
	fullPath := filepath.Join(c.config.BasePath, key)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// Open 打开文件读取流
func (c *LocalClient) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(c.config.BasePath, key)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Exists 判断文件是否存在
func (c *LocalClient) Exists(ctx context.Context, key string) (bool, error) {
	fullPath := filepath.Join(c.config.BasePath, key)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Delete 删除文件
func (c *LocalClient) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(c.config.BasePath, key)
	err := os.Remove(fullPath)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}

// GetURL 获取文件访问URL
func (c *LocalClient) GetURL(key string) string {
	return strings.TrimSuffix(c.config.BaseURL, "/") + "/" + key
}

// GetSignedURL 获取签名URL
func (c *LocalClient) GetSignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	expireTime := time.Now().Add(expires).Unix()
	signature := c.generateSignature(key, expireTime)
	return fmt.Sprintf("%s?expires=%d&signature=%s", c.GetURL(key), expireTime, signature), nil
}

// GetUploadCredential 获取上传凭证
func (c *LocalClient) GetUploadCredential(ctx context.Context, key string, expires time.Duration) (*UploadCredential, error) {
	expireTime := time.Now().Add(expires).Unix()
	signature := c.generateSignature(key, expireTime)

	// 上传接口统一走后端 API 前缀 /api/v1，不再复用静态访问前缀 /upload
	uploadURL := "/api/v1/files/upload/local"

	return &UploadCredential{
		Provider:  "local",
		UploadURL: uploadURL,
		Key:       key,
		Expires:   expireTime,
		Method:    "POST",
		Headers: map[string]string{
			"X-Upload-Signature": signature,
			"X-Upload-Expires":   strconv.FormatInt(expireTime, 10),
			"X-Upload-Key":       key,
		},
		PreviewURL: c.GetURL(key),
	}, nil
}

// InitMultipartUpload 初始化分片上传
func (c *LocalClient) InitMultipartUpload(ctx context.Context, key string) (*MultipartUpload, error) {
	uploadID := fmt.Sprintf("%s-%d", hex.EncodeToString([]byte(key))[:16], time.Now().UnixNano())

	// 创建临时分片目录
	chunkDir := filepath.Join(c.config.BasePath, ".chunks", uploadID)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		return nil, fmt.Errorf("创建分片目录失败: %v", err)
	}

	return &MultipartUpload{
		UploadID:  uploadID,
		Key:       key,
		ChunkSize: 5 * 1024 * 1024, // 5MB
	}, nil
}

// GetMultipartUploadURL 获取分片上传URL
func (c *LocalClient) GetMultipartUploadURL(ctx context.Context, uploadID, key string, partNumber int, expires time.Duration) (string, error) {
	expireTime := time.Now().Add(expires).Unix()
	data := fmt.Sprintf("%s:%s:%d", uploadID, key, partNumber)
	signature := c.generateSignature(data, expireTime)

	// 分片上传接口也统一走后端 API 前缀 /api/v1
	uploadURL := fmt.Sprintf("/api/v1/files/upload/chunk?upload_id=%s&key=%s&part_number=%d&expires=%d&signature=%s",
		uploadID, key, partNumber, expireTime, signature)

	return uploadURL, nil
}

// CompleteMultipartUpload 完成分片上传
func (c *LocalClient) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []Part) error {
	chunkDir := filepath.Join(c.config.BasePath, ".chunks", uploadID)

	// 目标文件路径
	fullPath := filepath.Join(c.config.BasePath, key)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建目标文件
	destFile, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer destFile.Close()

	// 按分片序号排序
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].PartNumber < parts[j].PartNumber
	})

	// 合并分片
	for _, part := range parts {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", part.PartNumber))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("打开分片文件失败: %v", err)
		}

		if _, err := io.Copy(destFile, chunkFile); err != nil {
			chunkFile.Close()
			return fmt.Errorf("合并分片失败: %v", err)
		}
		chunkFile.Close()
	}

	// 删除分片目录
	os.RemoveAll(chunkDir)

	return nil
}

// AbortMultipartUpload 取消分片上传
func (c *LocalClient) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	chunkDir := filepath.Join(c.config.BasePath, ".chunks", uploadID)
	return os.RemoveAll(chunkDir)
}

// ListParts 列出已上传的分片
func (c *LocalClient) ListParts(ctx context.Context, key, uploadID string) ([]Part, error) {
	chunkDir := filepath.Join(c.config.BasePath, ".chunks", uploadID)

	entries, err := os.ReadDir(chunkDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Part{}, nil
		}
		return nil, err
	}

	var parts []Part
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		partNumber, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		info, _ := entry.Info()
		parts = append(parts, Part{
			PartNumber: partNumber,
			Size:       info.Size(),
		})
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].PartNumber < parts[j].PartNumber
	})

	return parts, nil
}

// generateSignature 生成签名
func (c *LocalClient) generateSignature(data string, expires int64) string {
	message := fmt.Sprintf("%s:%d", data, expires)
	h := hmac.New(sha256.New, []byte(c.secret))
	h.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// VerifySignature 验证签名
func (c *LocalClient) VerifySignature(data string, expires int64, signature string) bool {
	if time.Now().Unix() > expires {
		return false
	}
	expected := c.generateSignature(data, expires)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// UploadChunk 上传分片（本地存储专用）
func (c *LocalClient) UploadChunk(uploadID string, partNumber int, reader io.Reader) error {
	chunkDir := filepath.Join(c.config.BasePath, ".chunks", uploadID)
	chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", partNumber))

	file, err := os.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("创建分片文件失败: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("写入分片失败: %v", err)
	}

	return nil
}
