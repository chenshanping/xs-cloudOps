package oss

import (
	"context"
	"io"
	"time"
)

// Client 对象存储客户端接口
type Client interface {
	// Upload 上传文件
	Upload(ctx context.Context, key string, reader io.Reader, size int64) error

	// Open 打开文件读取流
	Open(ctx context.Context, key string) (io.ReadCloser, error)

	// Exists 判断文件是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// Delete 删除文件
	Delete(ctx context.Context, key string) error

	// GetURL 获取文件访问URL
	GetURL(key string) string

	// GetSignedURL 获取签名URL（用于私有文件访问）
	GetSignedURL(ctx context.Context, key string, expires time.Duration) (string, error)

	// GetUploadCredential 获取前端直传凭证
	GetUploadCredential(ctx context.Context, key string, expires time.Duration) (*UploadCredential, error)

	// InitMultipartUpload 初始化分片上传
	InitMultipartUpload(ctx context.Context, key string) (*MultipartUpload, error)

	// GetMultipartUploadURL 获取分片上传URL
	GetMultipartUploadURL(ctx context.Context, uploadID, key string, partNumber int, expires time.Duration) (string, error)

	// CompleteMultipartUpload 完成分片上传
	CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []Part) error

	// AbortMultipartUpload 取消分片上传
	AbortMultipartUpload(ctx context.Context, key, uploadID string) error

	// ListParts 列出已上传的分片
	ListParts(ctx context.Context, key, uploadID string) ([]Part, error)
}

// UploadCredential 上传凭证
type UploadCredential struct {
	// 公共字段
	Provider   string            `json:"provider"`    // 提供商: local/aliyun/tencent/minio
	UploadURL  string            `json:"upload_url"`  // 上传地址
	Key        string            `json:"key"`         // 文件key
	Expires    int64             `json:"expires"`     // 过期时间戳
	Headers    map[string]string `json:"headers"`     // 需要携带的请求头
	FormData   map[string]string `json:"form_data"`   // 表单数据（用于POST方式）
	Method     string            `json:"method"`      // 上传方法 PUT/POST
	PreviewURL string            `json:"preview_url"` // 预览URL

	// STS临时凭证（用于SDK直传）
	AccessKeyID     string `json:"access_key_id,omitempty"`
	AccessKeySecret string `json:"access_key_secret,omitempty"`
	SecurityToken   string `json:"security_token,omitempty"`
	Bucket          string `json:"bucket,omitempty"`
	Region          string `json:"region,omitempty"`
	Endpoint        string `json:"endpoint,omitempty"`
}

// MultipartUpload 分片上传信息
type MultipartUpload struct {
	UploadID  string `json:"upload_id"`
	Key       string `json:"key"`
	Bucket    string `json:"bucket,omitempty"`
	ChunkSize int64  `json:"chunk_size"` // 推荐分片大小
}

// Part 分片信息
type Part struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
	Size       int64  `json:"size,omitempty"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	Type   string // 存储类型
	Config string // JSON配置
}
