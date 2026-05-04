// 文件信息
export interface FileInfo {
  id: number
  name: string
  path: string
  url: string
  size: number
  ext: string
  mime_type: string
  md5: string
  storage_type: string
  uploader_id: number
  status: number
  reference_count: number
  created_at: string
  updated_at: string
}

// 上传凭证
export interface UploadCredential {
  provider: string
  upload_url: string
  key: string
  expires: number
  headers: Record<string, string>
  form_data: Record<string, string>
  method: string
  preview_url: string
  access_key_id?: string
  access_key_secret?: string
  security_token?: string
  bucket?: string
  region?: string
  endpoint?: string
}

// 分片上传信息
export interface MultipartUpload {
  upload_id: string
  key: string
  bucket?: string
  chunk_size: number
}

// 分片信息
export interface Part {
  part_number: number
  etag: string
  size?: number
}

// 初始化分片上传响应
export interface InitMultipartUploadResponse {
  instant_upload: boolean
  file?: FileInfo
  upload_id?: string
  key?: string
  chunk_size?: number
  total_parts?: number
  upload_urls?: string[]
}

// 上传进度
export interface UploadProgress {
  file: File
  status: 'pending' | 'calculating' | 'uploading' | 'success' | 'error'
  progress: number
  md5?: string
  error?: string
  result?: FileInfo
}

export interface FileMigrationItem {
  file_id: number
  file_name: string
  source_storage_type: string
  target_storage_type: string
  old_url: string
  new_url: string
  action: 'PENDING' | 'SKIP' | 'CONFLICT' | 'MISSING_SOURCE' | 'MIGRATED' | 'WARNING' | 'FAILED'
  message: string
}

export interface FileMigrationResult {
  target_storage_type: string
  total_count: number
  total_size: number
  pending_count: number
  pending_size: number
  skipped_count: number
  skipped_size: number
  conflict_count: number
  conflict_size: number
  missing_source_count: number
  missing_source_size: number
  migrated_count: number
  failed_count: number
  warning_count: number
  items: FileMigrationItem[]
}

export type FileMigrationTaskState = 'SCANNING' | 'RUNNING' | 'SUCCESS' | 'FAILED'

export interface FileMigrationTaskStatus {
  task_id: string
  status: FileMigrationTaskState
  message: string
  scope: FileMigrationScope
  source_storage_type: string
  target_storage_type: string
  total_count: number
  total_size: number
  pending_count: number
  pending_size: number
  skipped_count: number
  skipped_size: number
  conflict_count: number
  conflict_size: number
  missing_source_count: number
  missing_source_size: number
  processed_count: number
  processed_size: number
  migrated_count: number
  failed_count: number
  warning_count: number
  current_file_id: number
  current_file_name: string
  started_at: string
  finished_at: string
  items: FileMigrationItem[]
}

export type FileMigrationScope = 'all' | 'filter' | 'selected'

export interface FileMigrationFilters {
  name?: string
  ext?: string
  referenced?: boolean
}

export interface FileMigrationRequest {
  scope: FileMigrationScope
  ids?: number[]
  source_storage_type: string
  target_storage_type: string
  source_config?: string
  filters?: FileMigrationFilters
}

// 允许的文件类型
export const imageTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/bmp', 'image/svg+xml']
export const documentTypes = [
  'application/pdf',
  'application/msword',
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  'application/vnd.ms-excel',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  'application/vnd.ms-powerpoint',
  'application/vnd.openxmlformats-officedocument.presentationml.presentation',
]
export const videoTypes = ['video/mp4', 'video/avi', 'video/quicktime', 'video/x-ms-wmv', 'video/x-flv']
export const audioTypes = ['audio/mpeg', 'audio/wav', 'audio/ogg']
export const archiveTypes = ['application/zip', 'application/x-rar-compressed', 'application/x-7z-compressed', 'application/x-tar', 'application/gzip']

// 文件扩展名映射
export const extMimeMap: Record<string, string> = {
  jpg: 'image/jpeg',
  jpeg: 'image/jpeg',
  png: 'image/png',
  gif: 'image/gif',
  webp: 'image/webp',
  bmp: 'image/bmp',
  svg: 'image/svg+xml',
  pdf: 'application/pdf',
  doc: 'application/msword',
  docx: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  xls: 'application/vnd.ms-excel',
  xlsx: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  ppt: 'application/vnd.ms-powerpoint',
  pptx: 'application/vnd.openxmlformats-officedocument.presentationml.presentation',
  mp4: 'video/mp4',
  mp3: 'audio/mpeg',
  zip: 'application/zip',
  rar: 'application/x-rar-compressed',
  txt: 'text/plain',
  json: 'application/json',
}
