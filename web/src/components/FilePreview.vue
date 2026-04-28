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
          <a-button v-if="canRetryLoad" @click="loadPreviewContent">
            <ReloadOutlined /> 重试预览
          </a-button>
          <span v-if="previewNotice" class="preview-notice">{{ previewNotice }}</span>
        </a-space>
      </div>
      <div class="preview-shell">
        <a-spin :spinning="loading" tip="加载预览中..." class="preview-spin">
          <div class="preview-stage">
            <template v-if="previewDescriptor.kind === 'image'">
              <div class="image-preview">
                <div class="image-stage">
                  <img
                    :src="url"
                    :alt="name"
                    :style="imageStyle"
                    @load="handleImageLoad"
                  />
                </div>
                <div class="image-toolbar">
                  <a-button-group>
                    <a-button :disabled="scale >= 3" @click="zoomIn">
                      <ZoomInOutlined />
                    </a-button>
                    <a-button :disabled="scale <= 0.5" @click="zoomOut">
                      <ZoomOutOutlined />
                    </a-button>
                    <a-button @click="resetZoom">
                      <ExpandOutlined />
                    </a-button>
                    <a-button @click="rotateLeft">
                      <RotateLeftOutlined />
                    </a-button>
                    <a-button @click="rotateRight">
                      <RotateRightOutlined />
                    </a-button>
                  </a-button-group>
                </div>
              </div>
            </template>

            <template v-else-if="previewDescriptor.kind === 'video'">
              <div class="media-preview">
                <video :src="url" controls class="video-preview">
                  您的浏览器不支持视频播放
                </video>
              </div>
            </template>

            <template v-else-if="previewDescriptor.kind === 'audio'">
              <div class="audio-preview">
                <div class="audio-icon">
                  <SoundOutlined />
                </div>
                <audio :src="url" controls class="audio-player" />
              </div>
            </template>

            <template v-else-if="previewDescriptor.kind === 'pdf'">
              <div v-if="pdfPreviewUrl" class="preview-office-wrap preview-office-wrap--pdf">
                <iframe
                  :src="pdfPreviewUrl"
                  class="pdf-preview"
                />
              </div>
              <div v-else class="status-preview">
                <a-empty :description="errorMessage || '正在加载 PDF 预览'" />
              </div>
            </template>

            <template
              v-else-if="
                previewDescriptor.kind === 'docx' ||
                previewDescriptor.kind === 'excel' ||
                previewDescriptor.kind === 'pptx'
              "
            >
              <div
                v-if="officePreviewComponent && previewSource"
                :class="['preview-office-wrap', `preview-office-wrap--${previewDescriptor.kind}`]"
              >
                <component
                  :is="officePreviewComponent"
                  :src="previewSource"
                  :options="officePreviewOptions"
                  class="office-viewer"
                  @rendered="handleOfficeRendered"
                  @error="handleOfficeError"
                />
              </div>
              <div v-else class="status-preview">
                <a-empty :description="errorMessage || '正在加载文档预览'" />
              </div>
            </template>

            <template v-else-if="previewDescriptor.kind === 'text'">
              <div v-if="textContent" class="text-preview">
                <div class="text-scroll">
                  <pre>{{ textContent }}</pre>
                </div>
              </div>
              <div v-else class="status-preview">
                <a-empty :description="errorMessage || '正在加载文本预览'" />
              </div>
            </template>

            <template v-else>
              <div class="unsupported-preview">
                <FileOutlined class="file-icon" />
                <p class="file-name">{{ name }}</p>
                <p class="file-info">{{ formatFileSize(size) }} | {{ ext.toUpperCase() }}</p>
                <p class="unsupported-message">{{ unsupportedMessage }}</p>
              </div>
            </template>
          </div>
        </a-spin>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, nextTick, onUnmounted, ref, watch, withDefaults, defineProps, type CSSProperties } from 'vue'
import {
  DownloadOutlined,
  ExpandOutlined,
  FileOutlined,
  ReloadOutlined,
  RotateLeftOutlined,
  RotateRightOutlined,
  SoundOutlined,
  ZoomInOutlined,
  ZoomOutOutlined,
} from '@ant-design/icons-vue'

import request from '@/utils/request'
import { formatFileSize } from '@/utils/upload'
import {
  getFilePreviewDescriptor,
  getPreferredPreviewMimeType,
  getUnsupportedPreviewMessage,
  type FilePreviewDescriptor,
} from './file-preview-utils'

import '@vue-office/docx/lib/index.css'
import '@vue-office/excel/lib/index.css'

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

const loading = ref(false)
const textContent = ref('')
const scale = ref(1)
const rotate = ref(0)
const errorMessage = ref('')
const previewNotice = ref('')
const pdfUrl = ref('')
const blobUrl = ref('')
const loadSequence = ref(0)

const VueOfficeDocx = defineAsyncComponent(() => import('@vue-office/docx'))
const VueOfficeExcel = defineAsyncComponent(() => import('@vue-office/excel'))
const VueOfficePptx = defineAsyncComponent(() => import('@vue-office/pptx'))

const previewDescriptor = computed<FilePreviewDescriptor>(() => getFilePreviewDescriptor({
  ext: props.ext,
  mimeType: props.mimeType,
  size: props.size,
  name: props.name,
}))

const unsupportedMessage = computed(() => getUnsupportedPreviewMessage(previewDescriptor.value))
const canRetryLoad = computed(() => (
  previewDescriptor.value.kind === 'text' ||
  previewDescriptor.value.kind === 'pdf' ||
  previewDescriptor.value.kind === 'docx' ||
  previewDescriptor.value.kind === 'excel' ||
  previewDescriptor.value.kind === 'pptx'
))

const officePreviewComponent = computed(() => {
  if (previewDescriptor.value.kind === 'docx') {
    return VueOfficeDocx
  }
  if (previewDescriptor.value.kind === 'excel') {
    return VueOfficeExcel
  }
  if (previewDescriptor.value.kind === 'pptx') {
    return VueOfficePptx
  }
  return null
})

const modalBodyStyle: CSSProperties = {
  padding: '16px 24px 24px',
  display: 'flex',
  flexDirection: 'column',
  flex: '1 1 auto',
  overflow: 'hidden',
  minHeight: 0,
}

const officePreviewOptions = computed(() => ({
  inWrapper: true,
}))

const pdfPreviewUrl = computed(() => {
  if (!pdfUrl.value) {
    return ''
  }
  return `${pdfUrl.value}#toolbar=1&navpanes=0&scrollbar=1&view=FitH`
})

const modalWidth = computed(() => {
  if (previewDescriptor.value.kind === 'image' || previewDescriptor.value.kind === 'pdf') return '82vw'
  if (
    previewDescriptor.value.kind === 'docx' ||
    previewDescriptor.value.kind === 'excel' ||
    previewDescriptor.value.kind === 'pptx'
  ) return '88vw'
  if (previewDescriptor.value.kind === 'video') return 800
  return 600
})

const imageStyle = computed(() => ({
  transform: `scale(${scale.value}) rotate(${rotate.value}deg)`,
  transition: 'transform 0.3s',
  maxWidth: '100%',
  maxHeight: '100%',
}))

const isOfficePreviewKind = computed(() => (
  previewDescriptor.value.kind === 'docx' ||
  previewDescriptor.value.kind === 'excel' ||
  previewDescriptor.value.kind === 'pptx'
))

const zoomIn = () => {
  scale.value = Math.min(scale.value + 0.25, 3)
}

const zoomOut = () => {
  scale.value = Math.max(scale.value - 0.25, 0.5)
}

const resetZoom = () => {
  scale.value = 1
  rotate.value = 0
}

const rotateLeft = () => {
  rotate.value -= 90
}

const rotateRight = () => {
  rotate.value += 90
}

const handleImageLoad = () => {
  // 图片加载完成
}

const getAbsoluteUrl = (url?: string) => {
  if (!url) {
    return ''
  }
  if (/^https?:\/\//.test(url)) {
    return url
  }
  if (url.startsWith('/')) {
    return `${window.location.origin}${url}`
  }
  return `${window.location.origin}/${url}`
}

const resetPreviewState = () => {
  loading.value = false
  textContent.value = ''
  errorMessage.value = ''
  previewNotice.value = ''
  pdfUrl.value = ''
  if (blobUrl.value) {
    URL.revokeObjectURL(blobUrl.value)
    blobUrl.value = ''
  }
}

const setBlobUrl = (blob: Blob) => {
  if (blobUrl.value) {
    URL.revokeObjectURL(blobUrl.value)
  }
  blobUrl.value = URL.createObjectURL(blob)
  return blobUrl.value
}

const normalizePreviewBlob = (blob: Blob) => {
  const preferredMimeType = getPreferredPreviewMimeType(previewDescriptor.value, props.mimeType || blob.type)
  if (!preferredMimeType || blob.type === preferredMimeType) {
    return blob
  }
  return new Blob([blob], { type: preferredMimeType })
}

const loadBinaryBlob = async () => {
  return request.get<Blob>(props.url, {
    responseType: 'blob',
    silent: true,
    timeout: 30000,
  })
}

const loadPreviewContent = async () => {
  const currentSequence = ++loadSequence.value
  resetPreviewState()

  if (!props.open || !props.url) {
    return
  }

  if (previewDescriptor.value.kind === 'unsupported') {
    previewNotice.value = unsupportedMessage.value
    return
  }

  if (
    previewDescriptor.value.kind === 'image' ||
    previewDescriptor.value.kind === 'video' ||
    previewDescriptor.value.kind === 'audio'
  ) {
    return
  }

  if (isOfficePreviewKind.value) {
    if (!previewSource.value) {
      errorMessage.value = '当前文档缺少可预览地址'
      return
    }
    loading.value = true
    return
  }

  loading.value = true
  try {
    const fileBlob = await loadBinaryBlob()
    if (currentSequence !== loadSequence.value) {
      return
    }

    const normalizedBlob = normalizePreviewBlob(fileBlob)

    if (previewDescriptor.value.kind === 'text') {
      textContent.value = await normalizedBlob.text()
      previewNotice.value = '文本内容已加载'
      return
    }

    if (previewDescriptor.value.kind === 'pdf') {
      pdfUrl.value = setBlobUrl(normalizedBlob)
      previewNotice.value = 'PDF 预览已加载'
      return
    }

  } catch (error) {
    errorMessage.value = '预览加载失败，请尝试下载查看'
  } finally {
    if (currentSequence === loadSequence.value) {
      loading.value = false
    }
  }
}

const handleOfficeRendered = () => {
  loading.value = false
  previewNotice.value = '文档预览已加载'
}

const handleOfficeError = () => {
  errorMessage.value = '文档解析失败，请下载后查看'
  loading.value = false
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
  if (!props.url) {
    return
  }

  if (blobUrl.value) {
    triggerNativeDownload(blobUrl.value)
    return
  }

  try {
    const fileBlob = await loadBinaryBlob()
    triggerNativeDownload(setBlobUrl(fileBlob))
  } catch {
    triggerNativeDownload(props.url)
  } finally {
    previewNotice.value = '已触发文件下载'
  }
}

const handleClose = () => {
  visible.value = false
  scale.value = 1
  rotate.value = 0
  resetPreviewState()
}

watch(() => [props.open, props.url, props.ext, props.size, props.mimeType], async ([open]) => {
  if (open) {
    await nextTick()
    loadPreviewContent()
    return
  }
  resetPreviewState()
})

onUnmounted(() => {
  resetPreviewState()
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
  color: #666;
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

:deep(.preview-spin.ant-spin-nested-loading) {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  min-width: 0;
  height: 100%;
  width: 100%;
}

:deep(.preview-spin .ant-spin-container) {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  min-width: 0;
  height: 100%;
  width: 100%;
}

.image-preview {
  display: flex;
  height: 100%;
  flex: 1;
  min-height: 0;
  flex-direction: column;
  overflow: hidden;
  background: #f5f5f5;
}

.image-stage {
  flex: 1;
  min-height: 0;
  overflow: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.image-preview img {
  cursor: move;
  transform-origin: center center;
}

.image-toolbar {
  display: flex;
  justify-content: center;
  padding: 12px 16px;
  border-top: 1px solid #f0f0f0;
  background: #fff;
}

.media-preview {
  display: flex;
  height: 100%;
  align-items: center;
  justify-content: center;
  background: #000;
  overflow: auto;
}

.video-preview {
  width: 100%;
  height: 100%;
  background: #000;
}

.audio-preview {
  display: flex;
  height: 100%;
  flex: 1;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
}

.audio-icon {
  font-size: 64px;
  color: #1890ff;
  margin-bottom: 24px;
}

.audio-player {
  width: 100%;
}

.pdf-preview {
  flex: 1;
  width: 100%;
  height: 100%;
  min-height: 100%;
  border: none;
  display: block;
  background: #fff;
}

.preview-office-wrap {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  min-width: 0;
  height: 100%;
  width: 100%;
  overflow: auto;
  background: #f5f5f5;
  box-sizing: border-box;
}

.office-viewer {
  display: block;
  width: 100%;
  background: #fff;
  box-sizing: border-box;
}

:deep(.preview-office-wrap--docx .office-viewer.vue-office-docx) {
  height: auto !important;
  min-height: 100% !important;
}

:deep(.preview-office-wrap--docx .office-viewer .vue-office-docx-main),
:deep(.preview-office-wrap--docx .office-viewer > div) {
  min-height: 100%;
}

:deep(.preview-office-wrap--excel .office-viewer.vue-office-excel),
:deep(.preview-office-wrap--pptx .office-viewer.vue-office-pptx) {
  flex: 1 0 auto;
  height: 100% !important;
  min-height: 100% !important;
}

:deep(.preview-office-wrap--excel .office-viewer .vue-office-excel-main),
:deep(.preview-office-wrap--excel .office-viewer > div),
:deep(.preview-office-wrap--pptx .office-viewer .vue-office-pptx-main),
:deep(.preview-office-wrap--pptx .office-viewer > div) {
  height: 100%;
  min-height: 100%;
}

.text-preview {
  height: 100%;
  flex: 1;
  min-height: 0;
  background: #f5f5f5;
  overflow: hidden;
}

.text-scroll {
  height: 100%;
  overflow: auto;
  padding: 16px;
  box-sizing: border-box;
}

.text-preview pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
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
  color: #999;
  margin-bottom: 16px;
}

.unsupported-preview .file-name {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 8px;
}

.unsupported-preview .file-info {
  color: #999;
  margin-bottom: 24px;
}

.unsupported-message {
  color: #666;
  margin-bottom: 0;
}

.status-preview {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  height: 100%;
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
}

:global(.file-preview-modal .ant-modal-header) {
  flex: none;
}

:global(.file-preview-modal .ant-modal-body) {
  display: flex;
  flex: 1;
  min-height: 0;
  overflow: hidden;
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
