<template>
  <a-modal
    v-model:open="visible"
    :title="title"
    :footer="null"
    :width="modalWidth"
    :centered="true"
    :destroy-on-close="true"
    @cancel="handleClose"
  >
    <div class="file-preview">
      <!-- 图片预览 -->
      <template v-if="isImage">
        <div class="image-preview">
          <img
            :src="url"
            :alt="name"
            :style="imageStyle"
            @load="handleImageLoad"
          />
          <div class="image-toolbar">
            <a-button-group>
              <a-button @click="zoomIn" :disabled="scale >= 3">
                <ZoomInOutlined />
              </a-button>
              <a-button @click="zoomOut" :disabled="scale <= 0.5">
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

      <!-- 视频预览 -->
      <template v-else-if="isVideo">
        <video
          :src="url"
          controls
          class="video-preview"
        >
          您的浏览器不支持视频播放
        </video>
      </template>

      <!-- 音频预览 -->
      <template v-else-if="isAudio">
        <div class="audio-preview">
          <div class="audio-icon">
            <SoundOutlined />
          </div>
          <audio :src="url" controls class="audio-player" />
        </div>
      </template>

      <!-- PDF预览 -->
      <template v-else-if="isPdf">
        <iframe :src="url" class="pdf-preview" />
      </template>

      <!-- 文本预览 -->
      <template v-else-if="isText">
        <div class="text-preview">
          <a-spin v-if="loading" />
          <pre v-else>{{ textContent }}</pre>
        </div>
      </template>

      <!-- 不支持预览 -->
      <template v-else>
        <div class="unsupported-preview">
          <FileOutlined class="file-icon" />
          <p class="file-name">{{ name }}</p>
          <p class="file-info">{{ formatFileSize(size) }} | {{ ext.toUpperCase() }}</p>
          <a-button type="primary" @click="handleDownload">
            <DownloadOutlined /> 下载文件
          </a-button>
        </div>
      </template>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch,withDefaults,defineProps } from 'vue'
import { ZoomInOutlined,ZoomOutOutlined,ExpandOutlined,RotateLeftOutlined,RotateRightOutlined } from '@ant-design/icons-vue';

import { imageTypes, videoTypes, audioTypes } from '@/types/file'
import { formatFileSize } from '@/utils/upload'

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

const loading = ref(false)
const textContent = ref('')
const scale = ref(1)
const rotate = ref(0)

// 判断文件类型
const isImage = computed(() => {
  const ext = props.ext.toLowerCase()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(ext) ||
    imageTypes.includes(props.mimeType)
})

const isVideo = computed(() => {
  const ext = props.ext.toLowerCase()
  return ['mp4', 'webm', 'ogg', 'avi', 'mov', 'wmv', 'flv'].includes(ext) ||
    videoTypes.includes(props.mimeType)
})

const isAudio = computed(() => {
  const ext = props.ext.toLowerCase()
  return ['mp3', 'wav', 'ogg', 'aac', 'flac'].includes(ext) ||
    audioTypes.includes(props.mimeType)
})

const isPdf = computed(() => {
  return props.ext.toLowerCase() === 'pdf' || props.mimeType === 'application/pdf'
})

const isText = computed(() => {
  const ext = props.ext.toLowerCase()
  return ['txt', 'md', 'json', 'xml', 'html', 'css', 'js', 'ts', 'vue', 'go', 'py'].includes(ext)
})

// 弹窗宽度
const modalWidth = computed(() => {
  if (isImage.value || isPdf.value) return '80%'
  if (isVideo.value) return 800
  return 600
})

// 图片样式
const imageStyle = computed(() => ({
  transform: `scale(${scale.value}) rotate(${rotate.value}deg)`,
  transition: 'transform 0.3s',
  maxWidth: '100%',
  maxHeight: '70vh',
}))

// 图片操作
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

// 加载文本内容
const loadTextContent = async () => {
  if (!isText.value || !props.url) return
  loading.value = true
  try {
    const response = await fetch(props.url)
    textContent.value = await response.text()
  } catch {
    textContent.value = '加载失败'
  } finally {
    loading.value = false
  }
}

// 下载文件
const handleDownload = () => {
  const a = document.createElement('a')
  a.href = props.url
  a.download = props.name
  a.click()
}

// 关闭
const handleClose = () => {
  visible.value = false
  scale.value = 1
  rotate.value = 0
  textContent.value = ''
}

// 监听打开状态
watch(() => props.open, (val) => {
  if (val && isText.value) {
    loadTextContent()
  }
})
</script>

<style scoped>
.file-preview {
  min-height: 200px;
}

.image-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.image-preview img {
  cursor: move;
}

.image-toolbar {
  padding: 8px 0;
}

.video-preview {
  width: 100%;
  max-height: 70vh;
}

.audio-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
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
  width: 100%;
  height: 70vh;
  border: none;
}

.text-preview {
  max-height: 70vh;
  overflow: auto;
  background: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
}

.text-preview pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.unsupported-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
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
</style>
