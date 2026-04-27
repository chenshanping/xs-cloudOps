package oss

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"server/model"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliyunClient 阿里云OSS客户端
type AliyunClient struct {
	config model.AliyunOSSConfig
	client *oss.Client
	bucket *oss.Bucket
}

// NewAliyunClient 创建阿里云OSS客户端
func NewAliyunClient(configJSON string) (*AliyunClient, error) {
	var config model.AliyunOSSConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析阿里云OSS配置失败: %v", err)
	}

	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("获取Bucket失败: %v", err)
	}

	return &AliyunClient{
		config: config,
		client: client,
		bucket: bucket,
	}, nil
}

// Upload 上传文件
func (c *AliyunClient) Upload(ctx context.Context, key string, reader io.Reader, size int64) error {
	return c.bucket.PutObject(key, reader)
}

// Open 打开文件读取流
func (c *AliyunClient) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	return c.bucket.GetObject(key)
}

// Exists 判断文件是否存在
func (c *AliyunClient) Exists(ctx context.Context, key string) (bool, error) {
	return c.bucket.IsObjectExist(key)
}

// Delete 删除文件
func (c *AliyunClient) Delete(ctx context.Context, key string) error {
	return c.bucket.DeleteObject(key)
}

// GetURL 获取文件访问URL
func (c *AliyunClient) GetURL(key string) string {
	return fmt.Sprintf("https://%s.%s/%s", c.config.BucketName, c.config.Endpoint, key)
}

// GetSignedURL 获取签名URL
func (c *AliyunClient) GetSignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	return c.bucket.SignURL(key, oss.HTTPGet, int64(expires.Seconds()))
}

// GetUploadCredential 获取上传凭证（使用签名URL方式）
func (c *AliyunClient) GetUploadCredential(ctx context.Context, key string, expires time.Duration) (*UploadCredential, error) {
	signedURL, err := c.bucket.SignURL(key, oss.HTTPPut, int64(expires.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("生成签名URL失败: %v", err)
	}

	return &UploadCredential{
		Provider:   "aliyun",
		UploadURL:  signedURL,
		Key:        key,
		Expires:    time.Now().Add(expires).Unix(),
		Method:     "PUT",
		Headers:    map[string]string{},
		PreviewURL: c.GetURL(key),
		Bucket:     c.config.BucketName,
		Region:     c.config.Region,
		Endpoint:   c.config.Endpoint,
	}, nil
}

// InitMultipartUpload 初始化分片上传
func (c *AliyunClient) InitMultipartUpload(ctx context.Context, key string) (*MultipartUpload, error) {
	result, err := c.bucket.InitiateMultipartUpload(key)
	if err != nil {
		return nil, fmt.Errorf("初始化分片上传失败: %v", err)
	}

	return &MultipartUpload{
		UploadID:  result.UploadID,
		Key:       key,
		Bucket:    c.config.BucketName,
		ChunkSize: 5 * 1024 * 1024, // 5MB
	}, nil
}

// GetMultipartUploadURL 获取分片上传URL
func (c *AliyunClient) GetMultipartUploadURL(ctx context.Context, uploadID, key string, partNumber int, expires time.Duration) (string, error) {
	options := []oss.Option{
		oss.AddParam("uploadId", uploadID),
		oss.AddParam("partNumber", fmt.Sprintf("%d", partNumber)),
	}
	return c.bucket.SignURL(key, oss.HTTPPut, int64(expires.Seconds()), options...)
}

// CompleteMultipartUpload 完成分片上传
func (c *AliyunClient) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []Part) error {
	var ossParts []oss.UploadPart
	for _, p := range parts {
		ossParts = append(ossParts, oss.UploadPart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		})
	}

	_, err := c.bucket.CompleteMultipartUpload(oss.InitiateMultipartUploadResult{
		UploadID: uploadID,
		Key:      key,
		Bucket:   c.config.BucketName,
	}, ossParts)
	return err
}

// AbortMultipartUpload 取消分片上传
func (c *AliyunClient) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	return c.bucket.AbortMultipartUpload(oss.InitiateMultipartUploadResult{
		UploadID: uploadID,
		Key:      key,
		Bucket:   c.config.BucketName,
	})
}

// ListParts 列出已上传的分片
func (c *AliyunClient) ListParts(ctx context.Context, key, uploadID string) ([]Part, error) {
	result, err := c.bucket.ListUploadedParts(oss.InitiateMultipartUploadResult{
		UploadID: uploadID,
		Key:      key,
		Bucket:   c.config.BucketName,
	})
	if err != nil {
		return nil, err
	}

	var parts []Part
	for _, p := range result.UploadedParts {
		parts = append(parts, Part{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
			Size:       int64(p.Size),
		})
	}
	return parts, nil
}
