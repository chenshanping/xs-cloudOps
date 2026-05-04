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
  | 'code'
  | 'markdown'
  | 'epub'
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

const imageExts = new Set(['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg', 'ico'])
const videoExts = new Set(['mp4', 'webm', 'ogg', 'avi', 'mov', 'wmv', 'flv', 'mkv', 'mpeg'])
const audioExts = new Set(['mp3', 'wav', 'ogg', 'aac', 'flac', 'wma'])
const textExts = new Set(['txt'])
const codeExts = new Set([
  'html', 'css', 'less', 'scss', 'js', 'json', 'ts', 'vue',
  'c', 'cpp', 'java', 'py', 'go', 'php', 'lua', 'rb', 'pl',
  'swift', 'vb', 'cs', 'sh', 'rs', 'vim', 'log', 'lock',
  'xml', 'mht', 'mhtml', 'mod', 'yaml', 'yml',
])
const markdownExts = new Set(['md'])
const docxExts = new Set(['docx'])
const excelExts = new Set(['xlsx', 'xls', 'csv', 'fods', 'ods', 'ots', 'xlsm', 'xlt', 'xltm'])
const pptxExts = new Set(['pptx', 'ppt', 'fodp', 'odp', 'otp', 'pot', 'potm', 'potx', 'pps', 'ppsm', 'ppsx', 'pptm'])
const legacyOfficeExts = new Set(['doc', 'docm', 'dot', 'dotm', 'dotx', 'fodt', 'odt', 'ott', 'rtf'])

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
  if (markdownExts.has(ext)) {
    return { kind: 'markdown' }
  }
  if (textExts.has(ext)) {
    return { kind: 'text' }
  }
  if (codeExts.has(ext) || mimeType.startsWith('text/')) {
    return { kind: 'code' }
  }
  if (ext === 'pdf' || mimeType === 'application/pdf') {
    if (size > FILE_PREVIEW_SIZE_LIMIT) {
      return { kind: 'unsupported', reason: 'too-large' }
    }
    return { kind: 'pdf' }
  }
  if (ext === 'epub') {
    return { kind: 'epub' }
  }
  if (legacyOfficeExts.has(ext)) {
    return { kind: 'unsupported', reason: 'legacy-office' }
  }
  if (docxExts.has(ext)) {
    if (size > FILE_PREVIEW_SIZE_LIMIT) {
      return { kind: 'unsupported', reason: 'too-large' }
    }
    return { kind: 'docx' }
  }
  if (excelExts.has(ext)) {
    if (size > FILE_PREVIEW_SIZE_LIMIT) {
      return { kind: 'unsupported', reason: 'too-large' }
    }
    return { kind: 'excel' }
  }
  if (pptxExts.has(ext)) {
    if (size > FILE_PREVIEW_SIZE_LIMIT) {
      return { kind: 'unsupported', reason: 'too-large' }
    }
    return { kind: 'pptx' }
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

