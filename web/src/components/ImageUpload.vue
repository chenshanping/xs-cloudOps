<template>
  <div class="image-upload">
    <!-- 已上传的图片预览 -->
    <div class="image-preview" v-if="currentUrl">
      <img :src="currentUrl" alt="preview" />
      <div class="image-actions">
        <EyeOutlined @click="handlePreview" />
        <DeleteOutlined @click="handleRemove" />
      </div>
      <a-progress
        v-if="uploading"
        :percent="progress"
        size="small"
        :show-info="false"
        class="upload-progress"
      />
    </div>

    <!-- 上传按钮 -->
    <a-upload
      v-else
      :show-upload-list="false"
      :before-upload="handleBeforeUpload"
      :custom-request="handleUpload"
      accept="image/*"
    >
      <div class="upload-btn" :style="{ width: `${width}px`, height: `${height}px` }">
        <LoadingOutlined v-if="uploading" />
        <PlusOutlined v-else />
        <div class="upload-text">{{ uploading ? '上传中...' : placeholder }}</div>
      </div>
    </a-upload>

    <!-- 图片预览弹窗 -->
    <a-modal
      v-model:open="previewVisible"
      :footer="null"
      :title="null"
      centered
    >
      <img :src="currentUrl" style="width: 100%" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, LoadingOutlined, EyeOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import type { UploadProps } from 'ant-design-vue'
import { calculateMD5, multipartUpload } from '@/utils/upload'

interface Props {
  /** 图片URL (v-model) */
  modelValue?: string
  /** 文件ID (v-model:fileId) */
  fileId?: number
  /** 图片URL (v-model:url) - 备用 */
  url?: string
  /** 上传按钮宽度 */
  width?: number
  /** 上传按钮高度 */
  height?: number
  /** 最大文件大小(字节)，默认2MB */
  maxSize?: number
  /** 占位文字 */
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  fileId: 0,
  url: '',
  width: 120,
  height: 120,
  maxSize: 2 * 1024 * 1024,
  placeholder: '上传图片',
})

const emit = defineEmits<{
  (e: 'update:modelValue', url: string): void
  (e: 'update:fileId', id: number): void
  (e: 'update:url', url: string): void
  (e: 'success', data: { id: number; url: string }): void
  (e: 'error', error: Error): void
}>()

const uploading = ref(false)
const progress = ref(0)
const localUrl = ref('')
const previewVisible = ref(false)

const currentUrl = computed(() => localUrl.value || props.modelValue || props.url)

// 上传前校验
const handleBeforeUpload: UploadProps['beforeUpload'] = (file) => {
  if (!file.type.startsWith('image/')) {
    message.error('只能上传图片文件')
    return false
  }
  if (file.size > props.maxSize) {
    const maxSizeMB = (props.maxSize / 1024 / 1024).toFixed(1)
    message.error(`图片大小不能超过 ${maxSizeMB}MB`)
    return false
  }
  return true
}

// 执行上传
const handleUpload: UploadProps['customRequest'] = async (options) => {
  const file = options.file as File
  uploading.value = true
  progress.value = 0

  try {
    // 本地预览
    const reader = new FileReader()
    reader.onload = (e) => {
      localUrl.value = e.target?.result as string
    }
    reader.readAsDataURL(file)

    // 计算MD5
    const md5 = await calculateMD5(file, (p) => {
      progress.value = Math.round(p * 0.1)
    })

    // 上传文件
    const result = await multipartUpload(file, md5, (p) => {
      progress.value = 10 + Math.round(p * 0.9)
    })

    // 上传成功
    localUrl.value = result.url
    emit('update:modelValue', result.url)
    emit('update:fileId', result.id)
    emit('update:url', result.url)
    emit('success', { id: result.id, url: result.url })
    message.success('上传成功')
    options.onSuccess?.(result)
  } catch (error) {
    localUrl.value = ''
    emit('error', error as Error)
    message.error((error as Error).message || '上传失败')
    options.onError?.(error as Error)
  } finally {
    uploading.value = false
    progress.value = 0
  }
}

// 预览
const handlePreview = () => {
  previewVisible.value = true
}

// 移除
const handleRemove = () => {
  localUrl.value = ''
  emit('update:modelValue', '')
  emit('update:fileId', 0)
  emit('update:url', '')
}

// 暴露方法
defineExpose({
  clear: handleRemove,
})
</script>

<style scoped>
.image-upload {
  display: inline-block;
}

.image-preview {
  position: relative;
  display: inline-block;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  overflow: hidden;
}

.image-preview img {
  max-width: 200px;
  max-height: 120px;
  display: block;
  object-fit: contain;
}

.image-actions {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  background: rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.3s;
}

.image-preview:hover .image-actions {
  opacity: 1;
}

.image-actions :deep(.anticon) {
  color: #fff;
  font-size: 18px;
  cursor: pointer;
}

.image-actions :deep(.anticon:hover) {
  color: #1890ff;
}

.upload-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 1px dashed #d9d9d9;
  border-radius: 4px;
  background: #fafafa;
  cursor: pointer;
  transition: all 0.3s;
  color: #999;
}

.upload-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.upload-text {
  font-size: 12px;
}

.upload-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
}
</style>
