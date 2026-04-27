package model

import "fmt"

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal   StorageType = "local"
	StorageTypeAliyun  StorageType = "aliyun"
	StorageTypeTencent StorageType = "tencent"
	StorageTypeMinio   StorageType = "minio"
)

// StorageProfile 运行时使用的存储配置
type StorageProfile struct {
	Name   string      `json:"name,omitempty"`
	Type   StorageType `json:"type"`
	Config string      `json:"-"`
}

func (p StorageProfile) CacheKey() string {
	return fmt.Sprintf("%s::%s", p.Type, p.Config)
}

// LegacyStorageRecord 仅用于兼容历史 sys_storage 表数据迁移
type LegacyStorageRecord struct {
	BaseModel
	Name      string      `gorm:"size:100;comment:存储名称"`
	Type      StorageType `gorm:"size:20;comment:存储类型 local/aliyun/tencent/minio"`
	Config    string      `gorm:"type:text;comment:存储配置(JSON)"`
	IsDefault int         `gorm:"default:0;comment:是否默认 0否 1是"`
	Status    int         `gorm:"default:1;comment:状态 0禁用 1启用"`
	Remark    string      `gorm:"size:255;comment:备注"`
}

func (LegacyStorageRecord) TableName() string {
	return "sys_storage"
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	BasePath string `json:"base_path"`
	BaseURL  string `json:"base_url"`
}

// AliyunOSSConfig 阿里云OSS配置
type AliyunOSSConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	Region          string `json:"region"`
	RoleArn         string `json:"role_arn"`
}

// TencentCOSConfig 腾讯云COS配置
type TencentCOSConfig struct {
	Region    string `json:"region"`
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	AppID     string `json:"app_id"`
}

// MinioConfig MinIO配置
type MinioConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	BucketName      string `json:"bucket_name"`
	UseSSL          bool   `json:"use_ssl"`
}
