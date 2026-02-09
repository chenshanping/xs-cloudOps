// 存储类型
export type StorageType = 'local' | 'aliyun' | 'tencent' | 'minio'

// 存储配置
export interface Storage {
  id: number
  name: string
  type: StorageType
  config: string
  is_default: number
  status: number
  remark: string
  created_at: string
  updated_at: string
}

// 本地存储配置
export interface LocalConfig {
  base_path: string
  base_url: string
}

// 阿里云OSS配置
export interface AliyunOSSConfig {
  endpoint: string
  access_key_id: string
  access_key_secret: string
  bucket_name: string
  region: string
  role_arn?: string
}

// 腾讯云COS配置
export interface TencentCOSConfig {
  region: string
  secret_id: string
  secret_key: string
  bucket: string
  app_id: string
}

// MinIO配置
export interface MinioConfig {
  endpoint: string
  access_key_id: string
  secret_access_key: string
  bucket_name: string
  use_ssl: boolean
}

// 存储配置联合类型
export type StorageConfig = LocalConfig | AliyunOSSConfig | TencentCOSConfig | MinioConfig

// 存储类型选项
export const storageTypeOptions = [
  { label: '本地存储', value: 'local' },
  { label: '阿里云 OSS', value: 'aliyun' },
  { label: '腾讯云 COS', value: 'tencent' },
  { label: 'MinIO', value: 'minio' },
]
