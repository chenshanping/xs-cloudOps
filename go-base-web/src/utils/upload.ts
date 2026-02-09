import SparkMD5 from 'spark-md5'
import type { FileInfo, Part, UploadCredential } from '@/types/file'
import { initMultipartUpload, completeMultipartUpload, saveUploadedFile, checkFileMD5 } from '@/api/file'

// 默认分片大小 5MB
const DEFAULT_CHUNK_SIZE = 5 * 1024 * 1024

// 断点续传存储key前缀
const UPLOAD_CACHE_PREFIX = 'upload_progress_'

// 上传进度缓存
interface UploadCache {
  uploadId: string
  key: string
  storageId: number
  filename: string
  fileSize: number
  md5: string
  chunkSize: number
  uploadedParts: Part[]
  uploadUrls: string[]
  timestamp: number
}

/**
 * 计算文件MD5
 */
export function calculateMD5(file: File, onProgress?: (progress: number) => void): Promise<string> {
  return new Promise((resolve, reject) => {
    const chunkSize = 2 * 1024 * 1024 // 2MB chunks for MD5 calculation
    const chunks = Math.ceil(file.size / chunkSize)
    let currentChunk = 0
    const spark = new SparkMD5.ArrayBuffer()
    const reader = new FileReader()

    reader.onload = (e) => {
      spark.append(e.target?.result as ArrayBuffer)
      currentChunk++

      if (onProgress) {
        onProgress(Math.round((currentChunk / chunks) * 100))
      }

      if (currentChunk < chunks) {
        loadNext()
      } else {
        resolve(spark.end())
      }
    }

    reader.onerror = () => {
      reject(new Error('文件读取失败'))
    }

    function loadNext() {
      const start = currentChunk * chunkSize
      const end = Math.min(start + chunkSize, file.size)
      reader.readAsArrayBuffer(file.slice(start, end))
    }

    loadNext()
  })
}

/**
 * 文件分片
 */
export function sliceFile(file: File, chunkSize: number = DEFAULT_CHUNK_SIZE): Blob[] {
  const chunks: Blob[] = []
  let start = 0

  while (start < file.size) {
    const end = Math.min(start + chunkSize, file.size)
    chunks.push(file.slice(start, end))
    start = end
  }

  return chunks
}

/**
 * 直传到对象存储
 */
export async function uploadToOSS(
  file: File,
  credential: UploadCredential,
  onProgress?: (progress: number) => void
): Promise<void> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    }

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        resolve()
      } else {
        reject(new Error(`上传失败: ${xhr.status}`))
      }
    }

    xhr.onerror = () => reject(new Error('网络错误'))

    xhr.open(credential.method, credential.upload_url, true)

    // 设置请求头
    for (const [key, value] of Object.entries(credential.headers || {})) {
      xhr.setRequestHeader(key, value)
    }

    // 根据上传方式发送
    if (credential.method === 'PUT') {
      xhr.send(file)
    } else {
      // POST方式，使用FormData
      const formData = new FormData()
      for (const [key, value] of Object.entries(credential.form_data || {})) {
        formData.append(key, value)
      }
      formData.append('file', file)
      xhr.send(formData)
    }
  })
}

/**
 * 上传单个分片
 */
async function uploadChunk(
  chunk: Blob,
  url: string,
  partNumber: number,
  onProgress?: (loaded: number, total: number) => void
): Promise<string> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()

    // 判断是走后端代理还是直传 OSS：
    // 1. 以 /api/ 开头 => 后端接口，需要带 JWT，使用 POST，并从 JSON 响应里拿 ETag。
    // 2. 其他情况 => 预签名直传 URL（阿里云/Tencent 等），使用 PUT，从响应头里拿 ETag。
    const isBackendAPI = url.startsWith('/api/')

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable && onProgress) {
        onProgress(e.loaded, e.total)
      }
    }

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        if (isBackendAPI) {
          // 后端 /api 接口返回的是 JSON：{ code, message, data: { part_number, etag } }
          try {
            const resp = xhr.responseText ? JSON.parse(xhr.responseText) : {}
            const etag = resp?.data?.etag || ''
            resolve(etag)
          } catch (e) {
            // MinIO / 本地存储在某些情况下可能不需要 ETag（本地存储完全不用 ETag），解析失败时返回空字符串即可。
            resolve('')
          }
        } else {
          // 直连 OSS/COS：ETag 在响应头里
          const rawETag = xhr.getResponseHeader('ETag') || xhr.getResponseHeader('etag') || ''
          resolve(rawETag.replace(/"/g, ''))
        }
      } else {
        reject(new Error(`分片上传失败: ${xhr.status}`))
      }
    }

    xhr.onerror = () => reject(new Error('网络错误'))

    // 后端接口走 POST，直传 OSS/COS 等预签名 URL 必须使用 PUT（签名时就是按 PUT 算的）
    xhr.open(isBackendAPI ? 'POST' : 'PUT', url, true)

    if (isBackendAPI) {
      // 携带登录 Token，保证通过后端 JWT 鉴权
      const token = localStorage.getItem('token')
      if (token) {
        xhr.setRequestHeader('Authorization', `Bearer ${token}`)
      }
    }

    xhr.send(chunk)
  })

}

/**
 * 获取缓存的上传进度
 */
function getUploadCache(md5: string): UploadCache | null {
  const key = UPLOAD_CACHE_PREFIX + md5
  const cached = localStorage.getItem(key)
  if (cached) {
    const data = JSON.parse(cached) as UploadCache
    // 检查是否过期（24小时）
    if (Date.now() - data.timestamp < 24 * 60 * 60 * 1000) {
      return data
    }
    localStorage.removeItem(key)
  }
  return null
}

/**
 * 保存上传进度到缓存
 */
function saveUploadCache(md5: string, data: Omit<UploadCache, 'timestamp'>) {
  const key = UPLOAD_CACHE_PREFIX + md5
  localStorage.setItem(key, JSON.stringify({ ...data, timestamp: Date.now() }))
}

/**
 * 清除上传进度缓存
 */
function clearUploadCache(md5: string) {
  const key = UPLOAD_CACHE_PREFIX + md5
  localStorage.removeItem(key)
}

/**
 * 分片上传（支持断点续传）
 */
export async function multipartUpload(
  file: File,
  md5: string,
  storageId?: number,
  onProgress?: (progress: number, stage: string) => void
): Promise<FileInfo> {
  // 1. 检查是否可以秒传
  const checkResult = await checkFileMD5(md5)
  if (checkResult.data.exists && checkResult.data.file) {
    onProgress?.(100, '秒传完成')
    return checkResult.data.file
  }

  // 2. 检查是否有缓存的上传进度
  let uploadId: string
  let key: string
  let uploadStorageId: number
  let chunkSize: number
  let uploadedParts: Part[] = []
  let uploadUrls: string[]

  const cache = getUploadCache(md5)
  if (cache && cache.fileSize === file.size && cache.filename === file.name) {
    // 恢复上传
    uploadId = cache.uploadId
    key = cache.key
    uploadStorageId = cache.storageId
    chunkSize = cache.chunkSize
    uploadedParts = cache.uploadedParts
    uploadUrls = cache.uploadUrls
    onProgress?.(0, '恢复上传')
  } else {
    // 初始化新上传
    const initResult = await initMultipartUpload(file.name, file.size, md5, storageId)
    const data = initResult.data

    // 如果服务端返回秒传结果
    if (data.instant_upload && data.file) {
      onProgress?.(100, '秒传完成')
      return data.file
    }

    uploadId = data.upload_id!
    key = data.key!
    uploadStorageId = data.storage_id!
    chunkSize = data.chunk_size!
    uploadUrls = data.upload_urls!
  }

  // 3. 分片文件
  const chunks = sliceFile(file, chunkSize)
  const totalChunks = chunks.length
  const uploadedSet = new Set(uploadedParts.map((p) => p.part_number))

  // 4. 上传分片
  let completedChunks = uploadedParts.length
  const allParts: Part[] = [...uploadedParts]

  for (let i = 0; i < totalChunks; i++) {
    const partNumber = i + 1

    // 跳过已上传的分片
    if (uploadedSet.has(partNumber)) {
      continue
    }

    try {
      const etag = await uploadChunk(chunks[i], uploadUrls[i], partNumber, (loaded, total) => {
        const chunkProgress = loaded / total
        const overallProgress = ((completedChunks + chunkProgress) / totalChunks) * 100
        onProgress?.(Math.round(overallProgress), `上传中 ${completedChunks + 1}/${totalChunks}`)
      })

      const part: Part = { part_number: partNumber, etag }
      allParts.push(part)
      completedChunks++

      // 保存进度
      saveUploadCache(md5, {
        uploadId,
        key,
        storageId: uploadStorageId,
        filename: file.name,
        fileSize: file.size,
        md5,
        chunkSize,
        uploadedParts: allParts,
        uploadUrls,
      })
    } catch (error) {
      throw new Error(`分片 ${partNumber} 上传失败: ${(error as Error).message}`)
    }
  }

  // 5. 完成上传
  onProgress?.(99, '合并文件')
  const result = await completeMultipartUpload({
    upload_id: uploadId,
    key,
    filename: file.name,
    file_size: file.size,
    md5,
    storage_id: uploadStorageId,
    parts: allParts.sort((a, b) => a.part_number - b.part_number),
  })

  // 清除缓存
  clearUploadCache(md5)
  onProgress?.(100, '上传完成')

  return result.data
}

/**
 * 简单上传（小文件）
 */
export async function simpleUpload(
  file: File,
  credential: UploadCredential,
  md5: string,
  onProgress?: (progress: number) => void
): Promise<FileInfo> {
  // 上传文件
  await uploadToOSS(file, credential, onProgress)

  // 保存文件记录
  const result = await saveUploadedFile({
    filename: file.name,
    key: credential.key,
    url: credential.preview_url,
    file_size: file.size,
    md5,
    storage_id: 0, // 使用默认存储
  })

  return result.data
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 获取文件图标
 */
export function getFileIcon(ext: string): string {
  const iconMap: Record<string, string> = {
    // 图片
    jpg: 'file-image',
    jpeg: 'file-image',
    png: 'file-image',
    gif: 'file-image',
    webp: 'file-image',
    bmp: 'file-image',
    svg: 'file-image',
    // 文档
    pdf: 'file-pdf',
    doc: 'file-word',
    docx: 'file-word',
    xls: 'file-excel',
    xlsx: 'file-excel',
    ppt: 'file-ppt',
    pptx: 'file-ppt',
    txt: 'file-text',
    // 视频
    mp4: 'video-camera',
    avi: 'video-camera',
    mov: 'video-camera',
    wmv: 'video-camera',
    flv: 'video-camera',
    // 音频
    mp3: 'audio',
    wav: 'audio',
    // 压缩包
    zip: 'file-zip',
    rar: 'file-zip',
    '7z': 'file-zip',
    tar: 'file-zip',
    gz: 'file-zip',
    // 代码
    js: 'code',
    ts: 'code',
    json: 'code',
    html: 'code',
    css: 'code',
  }
  return iconMap[ext.toLowerCase()] || 'file'
}

/**
 * 验证文件类型
 */
export function validateFileType(file: File, accept: string): boolean {
  if (!accept) return true

  const acceptList = accept.split(',').map((s) => s.trim())
  const fileExt = '.' + file.name.split('.').pop()?.toLowerCase()
  const fileMime = file.type

  return acceptList.some((a) => {
    if (a.startsWith('.')) {
      // 扩展名匹配
      return fileExt === a.toLowerCase()
    } else if (a.endsWith('/*')) {
      // MIME类型通配符
      const prefix = a.replace('/*', '')
      return fileMime.startsWith(prefix)
    } else {
      // 完整MIME类型
      return fileMime === a
    }
  })
}

/**
 * 验证文件大小
 */
export function validateFileSize(file: File, maxSize: number): boolean {
  return file.size <= maxSize
}
