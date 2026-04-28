export const FILE_PREVIEW_SIZE_LIMIT = 20 * 1024 * 1024

export type FilePreviewKind =
  | 'image'
  | 'video'
  | 'audio'
  | 'pdf'
  | 'docx'
  | 'excel'
  | 'pptx'
  | 'text'
  | 'unsupported'

export type FilePreviewUnsupportedReason =
  | 'empty'
  | 'legacy-office'
  | 'too-large'
  | 'unknown'

export interface FilePreviewDescriptor {
  kind: FilePreviewKind
  reason?: FilePreviewUnsupportedReason
}

interface FilePreviewInput {
  ext?: string
  mimeType?: string
  size?: number
  name?: string
}

const imageExts = new Set(['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'])
const videoExts = new Set(['mp4', 'webm', 'ogg', 'avi', 'mov', 'wmv', 'flv', 'mkv'])
const audioExts = new Set(['mp3', 'wav', 'ogg', 'aac', 'flac'])
const textExts = new Set(['txt', 'md', 'json', 'xml', 'html', 'css', 'js', 'ts', 'vue', 'go', 'py', 'yaml', 'yml'])
const officePreviewExts = new Set(['docx', 'xlsx', 'pptx'])
const legacyOfficeExts = new Set(['doc', 'xls', 'ppt'])

function normalizeExt(input: FilePreviewInput) {
  if (input.ext?.trim()) {
    return input.ext.trim().toLowerCase()
  }
  const name = input.name?.trim() || ''
  const segments = name.split('.')
  if (segments.length < 2) {
    return ''
  }
  return segments[segments.length - 1]?.toLowerCase() || ''
}

export function getFilePreviewDescriptor(input: FilePreviewInput): FilePreviewDescriptor {
  const ext = normalizeExt(input)
  const mimeType = (input.mimeType || '').toLowerCase()
  const size = input.size || 0

  if (!ext && !mimeType) {
    return { kind: 'unsupported', reason: 'empty' }
  }

  if (imageExts.has(ext) || mimeType.startsWith('image/')) {
    return { kind: 'image' }
  }
  if (videoExts.has(ext) || mimeType.startsWith('video/')) {
    return { kind: 'video' }
  }
  if (audioExts.has(ext) || mimeType.startsWith('audio/')) {
    return { kind: 'audio' }
  }
  if (textExts.has(ext) || mimeType.startsWith('text/')) {
    return { kind: 'text' }
  }
  if (ext === 'pdf' || mimeType === 'application/pdf') {
    if (size > FILE_PREVIEW_SIZE_LIMIT) {
      return { kind: 'unsupported', reason: 'too-large' }
    }
    return { kind: 'pdf' }
  }
  if (legacyOfficeExts.has(ext)) {
    return { kind: 'unsupported', reason: 'legacy-office' }
  }
  if (officePreviewExts.has(ext)) {
    if (size > FILE_PREVIEW_SIZE_LIMIT) {
      return { kind: 'unsupported', reason: 'too-large' }
    }
    if (ext === 'docx') {
      return { kind: 'docx' }
    }
    if (ext === 'pptx') {
      return { kind: 'pptx' }
    }
    return { kind: 'excel' }
  }

  return { kind: 'unsupported', reason: 'unknown' }
}

export function getUnsupportedPreviewMessage(descriptor: FilePreviewDescriptor) {
  switch (descriptor.reason) {
    case 'empty':
      return '缺少可预览的文件信息'
    case 'legacy-office':
      return '当前旧版 Office 格式暂不支持在线预览，请下载后查看'
    case 'too-large':
      return '文件过大，已关闭在线预览，请下载后查看'
    default:
      return '当前文件暂不支持在线预览，请下载后查看'
  }
}

export function getPreferredPreviewMimeType(descriptor: FilePreviewDescriptor, mimeType?: string) {
  const normalizedMimeType = (mimeType || '').trim().toLowerCase()
  if (descriptor.kind === 'pdf') {
    return 'application/pdf'
  }
  if (normalizedMimeType) {
    return normalizedMimeType
  }
  if (descriptor.kind === 'text') {
    return 'text/plain'
  }
  return ''
}
