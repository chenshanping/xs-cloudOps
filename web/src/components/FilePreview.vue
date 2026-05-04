<template>
  <a-modal
    v-model:open="visible"
    :title="title"
    :footer="null"
    :width="modalWidth"
    wrapClassName="file-preview-modal"
    :bodyStyle="modalBodyStyle"
    :destroy-on-close="true"
    @cancel="handleClose"
  >
    <div class="file-preview">
      <div class="preview-toolbar">
        <a-space wrap>
          <a-button @click="handleDownload">
            <DownloadOutlined /> 下载文件
          </a-button>
          <span v-if="previewNotice" class="preview-notice">{{ previewNotice }}</span>
        </a-space>
      </div>
      <div class="preview-shell">
        <div class="preview-stage">
          <template v-if="previewDescriptor.kind === 'unsupported'">
            <div class="unsupported-preview">
              <FileOutlined class="file-icon" />
              <p class="file-name">{{ name }}</p>
              <p class="file-info">{{ formatFileSize(size) }} | {{ ext.toUpperCase() }}</p>
              <p class="unsupported-message">{{ unsupportedMessage }}</p>
            </div>
          </template>

          <template v-else-if="previewDescriptor.kind === 'pdf'">
            <div v-if="pdfPreviewUrl" class="pdf-preview-wrap">
              <iframe :src="pdfPreviewUrl" class="pdf-iframe" />
            </div>
            <div v-else class="unsupported-preview">
              <a-spin v-if="pdfLoading" tip="加载 PDF 中..." />
              <p v-else>{{ previewNotice || 'PDF 加载失败，请下载查看' }}</p>
            </div>
          </template>

          <template v-else-if="previewSource">
            <VueFilesPreview
              :url="previewSource"
              :name="name"
              width="100%"
              height="100%"
              @rendered="handleRendered"
              @error="handlePreviewError"
            />
          </template>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, nextTick, onUnmounted, ref, watch, type CSSProperties } from 'vue'
import { DownloadOutlined, FileOutlined } from '@ant-design/icons-vue'
import { VueFilesPreview } from 'vue-files-preview'
import 'vue-files-preview/lib/style.css'

import request from '@/utils/request'
import { formatFileSize } from '@/utils/upload'
import {
  getFilePreviewDescriptor,
  getUnsupportedPreviewMessage,
  type FilePreviewDescriptor,
} from './file-preview-utils'

interface Props {
  open: boolean
  url: string
  name: string
  ext: string
  size?: number
  mimeType?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 0,
  mimeType: '',
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const visible = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val),
})

const title = computed(() => props.name || '文件预览')
const previewSource = computed(() => getAbsoluteUrl(props.url))

const previewNotice = ref('')
const blobUrl = ref('')
const pdfBlobUrl = ref('')
const pdfLoading = ref(false)

const pdfPreviewUrl = computed(() => {
  if (!pdfBlobUrl.value) return ''
  return `${pdfBlobUrl.value}#toolbar=1&navpanes=0&scrollbar=1&view=FitH`
})

const previewDescriptor = computed<FilePreviewDescriptor>(() => getFilePreviewDescriptor({
  ext: props.ext,
  mimeType: props.mimeType,
  size: props.size,
  name: props.name,
}))

const unsupportedMessage = computed(() => getUnsupportedPreviewMessage(previewDescriptor.value))

const modalBodyStyle: CSSProperties = {
  padding: '16px 24px 24px',
  display: 'flex',
  flexDirection: 'column',
  flex: '1 1 auto',
  overflow: 'hidden',
  minHeight: 0,
}

const modalWidth = computed(() => {
  const kind = previewDescriptor.value.kind
  if (kind === 'video') return 800
  if (kind === 'audio' || kind === 'unsupported') return 600
  return '85vw'
})

const getAbsoluteUrl = (url?: string) => {
  if (!url) return ''
  if (/^https?:\/\//.test(url)) return url
  if (url.startsWith('/')) return `${window.location.origin}${url}`
  return `${window.location.origin}/${url}`
}

const loadPdf = async () => {
  if (!props.url) return
  pdfLoading.value = true
  try {
    const fileBlob = await request.get<Blob>(props.url, {
      responseType: 'blob',
      silent: true,
      timeout: 30000,
    })
    const blob = new Blob([fileBlob as unknown as BlobPart], { type: 'application/pdf' })
    pdfBlobUrl.value = URL.createObjectURL(blob)
    previewNotice.value = 'PDF 预览已加载'
  } catch {
    previewNotice.value = 'PDF 加载失败，请下载查看'
  } finally {
    pdfLoading.value = false
  }
}

const handleRendered = () => {
  previewNotice.value = '预览已加载'
}

const handlePreviewError = () => {
  previewNotice.value = '预览加载失败，请尝试下载查看'
}

const triggerNativeDownload = (href: string) => {
  const a = document.createElement('a')
  a.href = href
  a.download = props.name
  a.target = '_blank'
  a.rel = 'noopener'
  a.click()
}

const handleDownload = async () => {
  if (!props.url) return

  if (blobUrl.value) {
    triggerNativeDownload(blobUrl.value)
    return
  }

  try {
    const fileBlob = await request.get<Blob>(props.url, {
      responseType: 'blob',
      silent: true,
      timeout: 30000,
    })
    blobUrl.value = URL.createObjectURL(fileBlob as unknown as Blob)
    triggerNativeDownload(blobUrl.value)
  } catch {
    triggerNativeDownload(props.url)
  } finally {
    previewNotice.value = '已触发文件下载'
  }
}

const resetState = () => {
  previewNotice.value = ''
  pdfLoading.value = false
  if (blobUrl.value) {
    URL.revokeObjectURL(blobUrl.value)
    blobUrl.value = ''
  }
  if (pdfBlobUrl.value) {
    URL.revokeObjectURL(pdfBlobUrl.value)
    pdfBlobUrl.value = ''
  }
}

const handleClose = () => {
  visible.value = false
  resetState()
}

watch(() => props.open, async (open) => {
  if (open) {
    await nextTick()
    if (previewDescriptor.value.kind === 'pdf') {
      loadPdf()
    }
  } else {
    resetState()
  }
})

onUnmounted(() => {
  resetState()
})
</script>

<style scoped>
.file-preview {
  display: flex;
  flex: 1;
  min-height: 0;
  height: 100%;
  width: 100%;
  flex-direction: column;
}

.preview-toolbar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 12px;
  flex-shrink: 0;
}

.preview-notice {
  color: var(--app-text-secondary);
}

.preview-shell {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  height: 100%;
  overflow: hidden;
}

.preview-stage {
  display: flex;
  flex-direction: column;
  flex: 1;
  height: 100%;
  width: 100%;
  min-height: 0;
  min-width: 0;
  overflow: hidden;
}

.pdf-preview-wrap {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  height: 100%;
  width: 100%;
  overflow: hidden;
  background: var(--app-surface-soft);
}

.pdf-iframe {
  flex: 1;
  width: 100%;
  height: 100%;
  min-height: 100%;
  border: none;
  display: block;
  background: var(--app-surface-color);
}

.unsupported-preview {
  display: flex;
  height: 100%;
  flex: 1;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
}

.unsupported-preview .file-icon {
  font-size: 64px;
  color: var(--app-text-muted);
  margin-bottom: 16px;
}

.unsupported-preview .file-name {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 8px;
}

.unsupported-preview .file-info {
  color: var(--app-text-muted);
  margin-bottom: 24px;
}

.unsupported-message {
  color: var(--app-text-secondary);
  margin-bottom: 0;
}

:global(.file-preview-modal .ant-modal) {
  top: 16px;
  padding-bottom: 16px;
}

:global(.file-preview-modal .ant-modal-content) {
  display: flex;
  flex-direction: column;
  height: calc(100dvh - 32px);
  max-height: calc(100dvh - 32px);
  background: var(--app-surface-color);
  color: var(--app-text-color);
}

:global(.file-preview-modal .ant-modal-header) {
  flex: none;
  background: var(--app-surface-color);
  border-bottom-color: var(--app-border-color);
}

:global(.file-preview-modal .ant-modal-body) {
  display: flex;
  flex: 1;
  min-height: 0;
  overflow: hidden;
  background: var(--app-surface-color);
}

@media (max-width: 768px) {
  :global(.file-preview-modal .ant-modal) {
    top: 8px;
    width: calc(100vw - 16px) !important;
    padding-bottom: 8px;
  }

  :global(.file-preview-modal .ant-modal-content) {
    height: calc(100dvh - 16px);
    max-height: calc(100dvh - 16px);
  }

  :global(.file-preview-modal .ant-modal-body) {
    padding: 12px !important;
  }
}
</style>
