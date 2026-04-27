import request from '@/utils/request'
import type {
  FileInfo,
  UploadCredential,
  InitMultipartUploadResponse,
  Part,
  FileMigrationRequest,
  FileMigrationResult,
  FileMigrationTaskStatus,
} from '@/types/file'
import type { PageResponse } from '@/types'

// 获取文件列表
export function getFileList(params: { page: number; page_size: number; name?: string; ext?: string }) {
  return request.get<PageResponse<FileInfo>>('/files', { params })
}

// 获取文件详情
export function getFile(id: number) {
  return request.get<FileInfo>(`/files/${id}`)
}

// 删除文件
export function deleteFile(id: number) {
  return request.delete(`/files/${id}`)
}

// 批量删除文件
export function batchDeleteFiles(ids: number[]) {
  return request.delete('/files/batch', { data: { ids } })
}

export function previewFileMigration(data: FileMigrationRequest) {
  return request.post<FileMigrationResult>('/files/migrate/preview', data)
}

export function executeFileMigration(data: FileMigrationRequest) {
  return request.post<FileMigrationTaskStatus>('/files/migrate/execute', data)
}

export function getCurrentFileMigrationTask() {
  return request.get<FileMigrationTaskStatus | null>('/files/migrate/task/current')
}

// 获取上传凭证
export function getUploadCredential(filename: string) {
  return request.post<UploadCredential>('/files/credential', { filename })
}

// 检查MD5（秒传）
export function checkFileMD5(md5: string) {
  return request.post<{ exists: boolean; file?: FileInfo }>('/files/check-md5', { md5 })
}

// 初始化分片上传
export function initMultipartUpload(filename: string, fileSize: number, md5: string) {
  return request.post<InitMultipartUploadResponse>('/files/multipart/init', {
    filename,
    file_size: fileSize,
    md5,
  })
}

// 获取已上传的分片列表
export function getUploadedParts(uploadId: string, key: string) {
  return request.get<Part[]>('/files/multipart/parts', {
    params: { upload_id: uploadId, key },
  })
}

// 完成分片上传
export function completeMultipartUpload(data: {
  upload_id: string
  key: string
  filename: string
  file_size: number
  md5: string
  parts: Part[]
}) {
  return request.post<FileInfo>('/files/multipart/complete', data)
}

// 取消分片上传
export function abortMultipartUpload(uploadId: string, key: string) {
  return request.post('/files/multipart/abort', {
    upload_id: uploadId,
    key,
  })
}

// 保存已上传的文件记录
export function saveUploadedFile(data: {
  filename: string
  key: string
  url: string
  file_size: number
  md5: string
}) {
  return request.post<FileInfo>('/files/save', data)
}
