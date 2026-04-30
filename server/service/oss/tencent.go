package oss

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"server/model"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// TencentClient 腾讯云COS客户端
type TencentClient struct {
	config model.TencentCOSConfig
	client *cos.Client
}

// NewTencentClient 创建腾讯云COS客户端
func NewTencentClient(configJSON string) (*TencentClient, error) {
	var config model.TencentCOSConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析腾讯云COS配置失败: %v", err)
	}

	bucketURL, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", config.Bucket, config.Region))
	serviceURL, _ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", config.Region))

	client := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL, ServiceURL: serviceURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	})

	return &TencentClient{
		config: config,
		client: client,
	}, nil
}

// Upload 上传文件
func (c *TencentClient) Upload(ctx context.Context, key string, reader io.Reader, size int64) error {
	_, err := c.client.Object.Put(ctx, key, reader, nil)
	return err
}

// Open 打开文件读取流
func (c *TencentClient) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := c.client.Object.Get(ctx, key, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Exists 判断文件是否存在
func (c *TencentClient) Exists(ctx context.Context, key string) (bool, error) {
	resp, err := c.client.Object.Head(ctx, key, nil)
	if err == nil {
		return true, nil
	}
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, err
}

// Delete 删除文件
func (c *TencentClient) Delete(ctx context.Context, key string) error {
	_, err := c.client.Object.Delete(ctx, key)
	return err
}

// GetURL 获取文件访问URL
func (c *TencentClient) GetURL(key string) string {
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", c.config.Bucket, c.config.Region, key)
}

// GetSignedURL 获取签名URL
func (c *TencentClient) GetSignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	presignedURL, err := c.client.Object.GetPresignedURL(ctx, http.MethodGet, key, c.config.SecretID, c.config.SecretKey, expires, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// GetUploadCredential 获取上传凭证
func (c *TencentClient) GetUploadCredential(ctx context.Context, key string, expires time.Duration) (*UploadCredential, error) {
	presignedURL, err := c.client.Object.GetPresignedURL(ctx, http.MethodPut, key, c.config.SecretID, c.config.SecretKey, expires, nil)
	if err != nil {
		return nil, fmt.Errorf("生成签名URL失败: %v", err)
	}

	return &UploadCredential{
		Provider:   "tencent",
		UploadURL:  presignedURL.String(),
		Key:        key,
		Expires:    time.Now().Add(expires).Unix(),
		Method:     "PUT",
		Headers:    map[string]string{},
		PreviewURL: c.GetURL(key),
		Bucket:     c.config.Bucket,
		Region:     c.config.Region,
	}, nil
}

// InitMultipartUpload 初始化分片上传
func (c *TencentClient) InitMultipartUpload(ctx context.Context, key string) (*MultipartUpload, error) {
	result, _, err := c.client.Object.InitiateMultipartUpload(ctx, key, nil)
	if err != nil {
		return nil, fmt.Errorf("初始化分片上传失败: %v", err)
	}

	return &MultipartUpload{
		UploadID:  result.UploadID,
		Key:       key,
		Bucket:    c.config.Bucket,
		ChunkSize: 5 * 1024 * 1024, // 5MB
	}, nil
}

// GetMultipartUploadURL 获取分片上传URL
func (c *TencentClient) GetMultipartUploadURL(ctx context.Context, uploadID, key string, partNumber int, expires time.Duration) (string, error) {
	opt := &cos.PresignedURLOptions{
		Query: &url.Values{},
	}
	opt.Query.Add("uploadId", uploadID)
	opt.Query.Add("partNumber", fmt.Sprintf("%d", partNumber))

	presignedURL, err := c.client.Object.GetPresignedURL(ctx, http.MethodPut, key, c.config.SecretID, c.config.SecretKey, expires, opt)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// CompleteMultipartUpload 完成分片上传
func (c *TencentClient) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []Part) error {
	opt := &cos.CompleteMultipartUploadOptions{}
	for _, p := range parts {
		opt.Parts = append(opt.Parts, cos.Object{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		})
	}
	_, _, err := c.client.Object.CompleteMultipartUpload(ctx, key, uploadID, opt)
	return err
}

// AbortMultipartUpload 取消分片上传
func (c *TencentClient) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	_, err := c.client.Object.AbortMultipartUpload(ctx, key, uploadID)
	return err
}

// ListParts 列出已上传的分片
func (c *TencentClient) ListParts(ctx context.Context, key, uploadID string) ([]Part, error) {
	result, _, err := c.client.Object.ListParts(ctx, key, uploadID, nil)
	if err != nil {
		return nil, err
	}

	var parts []Part
	for _, p := range result.Parts {
		parts = append(parts, Part{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
			Size:       int64(p.Size),
		})
	}
	return parts, nil
}
