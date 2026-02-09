package oss

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go-base-server/model"
)

// MinioClient MinIO客户端
type MinioClient struct {
	config model.MinioConfig
	client *minio.Client
	core   *minio.Core
}

// NewMinioClient 创建MinIO客户端
func NewMinioClient(configJSON string) (*MinioClient, error) {
	var config model.MinioConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析MinIO配置失败: %v", err)
	}

	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("创建MinIO客户端失败: %v", err)
	}

	// 创建Core客户端用于分片上传
	core, err := minio.NewCore(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("创建MinIO Core客户端失败: %v", err)
	}

	return &MinioClient{
		config: config,
		client: client,
		core:   core,
	}, nil
}

// Upload 上传文件
func (c *MinioClient) Upload(ctx context.Context, key string, reader io.Reader, size int64) error {
	_, err := c.client.PutObject(ctx, c.config.BucketName, key, reader, size, minio.PutObjectOptions{})
	return err
}

// Delete 删除文件
func (c *MinioClient) Delete(ctx context.Context, key string) error {
	return c.client.RemoveObject(ctx, c.config.BucketName, key, minio.RemoveObjectOptions{})
}

// GetURL 获取文件访问URL
func (c *MinioClient) GetURL(key string) string {
	protocol := "http"
	if c.config.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, c.config.Endpoint, c.config.BucketName, key)
}

// GetSignedURL 获取签名URL
func (c *MinioClient) GetSignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := c.client.PresignedGetObject(ctx, c.config.BucketName, key, expires, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// GetUploadCredential 获取上传凭证
func (c *MinioClient) GetUploadCredential(ctx context.Context, key string, expires time.Duration) (*UploadCredential, error) {
	presignedURL, err := c.client.PresignedPutObject(ctx, c.config.BucketName, key, expires)
	if err != nil {
		return nil, fmt.Errorf("生成签名URL失败: %v", err)
	}

	return &UploadCredential{
		Provider:   "minio",
		UploadURL:  presignedURL.String(),
		Key:        key,
		Expires:    time.Now().Add(expires).Unix(),
		Method:     "PUT",
		Headers:    map[string]string{},
		PreviewURL: c.GetURL(key),
		Bucket:     c.config.BucketName,
		Endpoint:   c.config.Endpoint,
	}, nil
}

// InitMultipartUpload 初始化分片上传
func (c *MinioClient) InitMultipartUpload(ctx context.Context, key string) (*MultipartUpload, error) {
	uploadID, err := c.core.NewMultipartUpload(ctx, c.config.BucketName, key, minio.PutObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("初始化分片上传失败: %v", err)
	}

	return &MultipartUpload{
		UploadID:  uploadID,
		Key:       key,
		Bucket:    c.config.BucketName,
		ChunkSize: 5 * 1024 * 1024, // 5MB
	}, nil
}

// GetMultipartUploadURL 获取分片上传URL(MinIO使用SDK上传，不支持presigned分片URL)
func (c *MinioClient) GetMultipartUploadURL(ctx context.Context, uploadID, key string, partNumber int, expires time.Duration) (string, error) {
	// MinIO SDK不直接支持presigned multipart upload URL
	// 需要使用 PutObjectPart 方法直接上传
	// 返回一个自定义的URL，前端需要通过后端代理上传
	return fmt.Sprintf("/api/v1/files/upload/chunk?provider=minio&upload_id=%s&key=%s&part_number=%d",
		uploadID, key, partNumber), nil
}

// CompleteMultipartUpload 完成分片上传
func (c *MinioClient) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []Part) error {
	var completeParts []minio.CompletePart
	for _, p := range parts {
		completeParts = append(completeParts, minio.CompletePart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		})
	}

	_, err := c.core.CompleteMultipartUpload(ctx, c.config.BucketName, key, uploadID, completeParts, minio.PutObjectOptions{})
	return err
}

// AbortMultipartUpload 取消分片上传
func (c *MinioClient) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	return c.core.AbortMultipartUpload(ctx, c.config.BucketName, key, uploadID)
}

// ListParts 列出已上传的分片
func (c *MinioClient) ListParts(ctx context.Context, key, uploadID string) ([]Part, error) {
	result, err := c.core.ListObjectParts(ctx, c.config.BucketName, key, uploadID, 0, 1000)
	if err != nil {
		return nil, err
	}

	var parts []Part
	for _, p := range result.ObjectParts {
		parts = append(parts, Part{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
			Size:       p.Size,
		})
	}
	return parts, nil
}

// UploadPart 上传分片（MinIO专用，用于代理上传）
func (c *MinioClient) UploadPart(ctx context.Context, key, uploadID string, partNumber int, reader io.Reader, size int64) (string, error) {
	part, err := c.core.PutObjectPart(ctx, c.config.BucketName, key, uploadID, partNumber, reader, size, minio.PutObjectPartOptions{})
	if err != nil {
		return "", err
	}
	return part.ETag, nil
}
