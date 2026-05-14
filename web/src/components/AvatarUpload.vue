<template>
  <div
    ref="uploadRootRef"
    class="avatar-upload"
    tabindex="0"
    @click="focusUploadRoot"
    @paste="handlePaste"
  >
    <a-upload
      :show-upload-list="false"
      :before-upload="handleBeforeUpload"
      :custom-request="handleUpload"
      accept="image/*"
    >
      <div class="avatar-container" :class="{ uploading: uploading }">
        <a-avatar :size="size" :src="currentUrl || undefined">
          <template #icon>
            <UserOutlined />
          </template>
        </a-avatar>
        <div class="avatar-overlay">
          <LoadingOutlined v-if="uploading" class="upload-icon" />
          <CameraOutlined v-else class="upload-icon" />
        </div>
        <a-progress
          v-if="uploading"
          type="circle"
          :percent="progress"
          :width="size"
          :show-info="false"
          class="upload-progress"
        />
      </div>
    </a-upload>
    <div class="avatar-tip" v-if="tip">{{ tip }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import { UserOutlined, CameraOutlined, LoadingOutlined } from '@ant-design/icons-vue'
import type { UploadProps } from 'ant-design-vue'
import type { FileInfo } from '@/types/file'
import { calculateMD5, multipartUpload } from '@/utils/upload'
import { getClipboardFiles, isEditablePasteTarget } from '@/utils/upload-paste'

interface Props {
  fileId?: number
  url?: string
  size?: number
  maxSize?: number // 最大文件大小（字节），默认 2MB
  tip?: string
}

const props = withDefaults(defineProps<Props>(), {
  fileId: undefined,
  url: '',
  size: 100,
  maxSize: 2 * 1024 * 1024,
  tip: '点击上传头像'
})

const emit = defineEmits<{
  (e: 'update:fileId', v: number | undefined): void
  (e: 'update:url', v: string): void
  (e: 'success', file: FileInfo): void
  (e: 'error', error: Error): void
}>()

const uploading = ref(false)
const progress = ref(0)
const localUrl = ref('')
const uploadRootRef = ref<HTMLElement>()

const currentUrl = computed(() => localUrl.value || props.url)

// 当父组件清空 fileId 时，同步清空 localUrl
watch(() => props.fileId, (newFileId) => {
  if (!newFileId) {
    localUrl.value = ''
  }
})

// 上传前校验
const handleBeforeUpload: UploadProps['beforeUpload'] = (file) => {
  // 校验文件类型
  if (!file.type.startsWith('image/')) {
    message.error('只能上传图片文件')
    return false
  }

  // 校验文件大小
  if (file.size > props.maxSize) {
    const maxSizeMB = (props.maxSize / 1024 / 1024).toFixed(1)
    message.error(`图片大小不能超过 ${maxSizeMB}MB`)
    return false
  }

  return true
}

const uploadAvatar = async (
  file: File,
  callbacks?: {
    onSuccess?: (result: FileInfo) => void
    onError?: (error: Error) => void
  }
) => {
  uploading.value = true
  progress.value = 0

  try {
    const reader = new FileReader()
    reader.onload = (e) => {
      localUrl.value = e.target?.result as string
    }
    reader.readAsDataURL(file)

    const md5 = await calculateMD5(file, (p) => {
      progress.value = Math.round(p * 0.1)
    })

    const result = await multipartUpload(file, md5, undefined, (p) => {
      progress.value = 10 + Math.round(p * 0.9)
    })

    emit('update:fileId', result.id)
    emit('update:url', result.url)
    localUrl.value = result.url
    emit('success', result)
    message.success('头像上传成功')
    callbacks?.onSuccess?.(result)
    return result
  } catch (error) {
    localUrl.value = ''
    emit('error', error as Error)
    message.error((error as Error).message || '上传失败')
    callbacks?.onError?.(error as Error)
    throw error
  } finally {
    uploading.value = false
    progress.value = 0
  }
}

// 自定义上传
const handleUpload: UploadProps['customRequest'] = async (options) => {
  const file = options.file as File
  await uploadAvatar(file, {
    onSuccess: (result) => options.onSuccess?.(result),
    onError: (error) => options.onError?.(error),
  })
}

const focusUploadRoot = () => {
  uploadRootRef.value?.focus()
}

const handlePaste = async (event: ClipboardEvent) => {
  if (isEditablePasteTarget(event.target)) {
    return
  }

  const file = getClipboardFiles(event, {
    multiple: false,
  }).find(item => item.type.startsWith('image/'))

  if (!file) {
    return
  }

  if (!handleBeforeUpload(file)) {
    return
  }

  event.preventDefault()
  await uploadAvatar(file)
}
</script>

<style scoped>
.avatar-upload {
  display: inline-block;
}

.avatar-container {
  position: relative;
  cursor: pointer;
  border-radius: 50%;
  overflow: hidden;
  transition: all 0.3s;
}

.avatar-container:hover .avatar-overlay {
  opacity: 1;
}

.avatar-container.uploading .avatar-overlay {
  opacity: 1;
  background: rgba(0, 0, 0, 0.6);
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.4);
  opacity: 0;
  transition: opacity 0.3s;
  border-radius: 50%;
}

.upload-icon {
  color: #fff;
  font-size: 24px;
}

.upload-progress {
  position: absolute;
  top: 0;
  left: 0;
}

.upload-progress :deep(.ant-progress-circle-trail) {
  stroke: rgba(255, 255, 255, 0.3);
}

.upload-progress :deep(.ant-progress-circle-path) {
  stroke: #1890ff;
}

.avatar-tip {
  text-align: center;
  margin-top: 8px;
  color: #999;
  font-size: 12px;
}
</style>
