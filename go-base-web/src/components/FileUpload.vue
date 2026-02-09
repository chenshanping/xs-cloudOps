<template>
  <div class="file-upload">
    <a-upload-dragger
      v-model:file-list="fileList"
      :multiple="multiple"
      :accept="accept"
      :before-upload="handleBeforeUpload"
      :custom-request="handleCustomUpload"
      :show-upload-list="false"
      @drop="handleDrop"
    >
      <p class="ant-upload-drag-icon">
        <InboxOutlined />
      </p>
      <p class="ant-upload-text">点击或拖拽文件到此区域上传</p>
      <p class="ant-upload-hint" v-if="hint">{{ hint }}</p>
      <p class="ant-upload-hint" v-else>
        支持单个或批量上传{{ maxSize ? `，单个文件最大 ${formatFileSize(maxSize)}` : '' }}
      </p>
    </a-upload-dragger>

    <!-- 上传列表 -->
    <div class="upload-list" v-if="uploadList.length > 0">
      <div
        v-for="(item, index) in uploadList"
        :key="index"
        class="upload-item"
        :class="item.status"
      >
        <div class="upload-item-info">
          <FileOutlined class="file-icon" />
          <span class="file-name" :title="item.file.name">{{ item.file.name }}</span>
          <span class="file-size">{{ formatFileSize(item.file.size) }}</span>
        </div>
        <div class="upload-item-status">
          <template v-if="item.status === 'calculating'">
            <span class="status-text">计算MD5...</span>
          </template>
          <template v-else-if="item.status === 'uploading'">
            <a-progress :percent="item.progress" size="small" :show-info="false" />
            <span class="progress-text">{{ item.progress }}%</span>
          </template>
          <template v-else-if="item.status === 'success'">
            <CheckCircleFilled class="success-icon" />
          </template>
          <template v-else-if="item.status === 'error'">
            <CloseCircleFilled class="error-icon" />
            <span class="error-text" :title="item.error">{{ item.error }}</span>
          </template>
        </div>
        <div class="upload-item-actions">
          <a-button
            v-if="item.status === 'error'"
            type="link"
            size="small"
            @click="retryUpload(index)"
          >
            重试
          </a-button>
          <a-button
            type="link"
            size="small"
            danger
            @click="removeItem(index)"
          >
            删除
          </a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { InboxOutlined, FileOutlined, CheckCircleFilled, CloseCircleFilled } from '@ant-design/icons-vue'
import type { UploadProps } from 'ant-design-vue'
import type { FileInfo, UploadProgress } from '@/types/file'
import {
  calculateMD5,
  multipartUpload,
  formatFileSize,
  validateFileType,
  validateFileSize,
} from '@/utils/upload'

interface Props {
  multiple?: boolean
  accept?: string
  maxSize?: number // 最大文件大小（字节）
  hint?: string
  storageId?: number
}

const props = withDefaults(defineProps<Props>(), {
  multiple: true,
  accept: '',
  maxSize: 0,
})

const emit = defineEmits<{
  (e: 'success', file: FileInfo): void
  (e: 'error', error: Error, file: File): void
}>()

const fileList = ref<any[]>([])
const uploadList = ref<UploadProgress[]>([])

// 上传前校验
const handleBeforeUpload: UploadProps['beforeUpload'] = (file) => {
  // 校验文件类型
  if (props.accept && !validateFileType(file, props.accept)) {
    message.error(`不支持的文件类型: ${file.name}`)
    return false
  }

  // 校验文件大小
  if (props.maxSize && !validateFileSize(file, props.maxSize)) {
    message.error(`文件过大: ${file.name}，最大支持 ${formatFileSize(props.maxSize)}`)
    return false
  }

  return true
}

// 自定义上传
const handleCustomUpload: UploadProps['customRequest'] = async (options) => {
  const file = options.file as File
  const index = uploadList.value.length

  // 添加到上传列表
  uploadList.value.push({
    file,
    status: 'calculating',
    progress: 0,
  })

  try {
    // 计算MD5
    const md5 = await calculateMD5(file, (progress) => {
      uploadList.value[index].progress = Math.round(progress * 0.1) // MD5计算占10%
    })
    uploadList.value[index].md5 = md5
    uploadList.value[index].status = 'uploading'

    // 执行上传
    const result = await multipartUpload(file, md5, props.storageId, (progress, stage) => {
      uploadList.value[index].progress = 10 + Math.round(progress * 0.9) // 上传占90%
    })

    // 上传成功
    uploadList.value[index].status = 'success'
    uploadList.value[index].progress = 100
    uploadList.value[index].result = result
    emit('success', result)
    options.onSuccess?.(result)
  } catch (error) {
    // 上传失败
    uploadList.value[index].status = 'error'
    uploadList.value[index].error = (error as Error).message
    emit('error', error as Error, file)
    options.onError?.(error as Error)
  }
}

// 重试上传
const retryUpload = async (index: number) => {
  const item = uploadList.value[index]
  item.status = 'uploading'
  item.progress = 0
  item.error = undefined

  try {
    let md5 = item.md5
    if (!md5) {
      item.status = 'calculating'
      md5 = await calculateMD5(item.file)
      item.md5 = md5
      item.status = 'uploading'
    }

    const result = await multipartUpload(item.file, md5, props.storageId, (progress) => {
      item.progress = progress
    })

    item.status = 'success'
    item.progress = 100
    item.result = result
    emit('success', result)
  } catch (error) {
    item.status = 'error'
    item.error = (error as Error).message
    emit('error', error as Error, item.file)
  }
}

// 移除项
const removeItem = (index: number) => {
  uploadList.value.splice(index, 1)
}

// 拖拽处理
const handleDrop = (e: DragEvent) => {
  console.log('drop', e)
}

// 暴露方法
defineExpose({
  uploadList,
  clearList: () => {
    uploadList.value = []
    fileList.value = []
  },
})
</script>

<style scoped>
.file-upload {
  width: 100%;
}

.upload-list {
  margin-top: 16px;
}

.upload-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  margin-bottom: 8px;
  background: #fafafa;
  border-radius: 4px;
  transition: all 0.3s;
}

.upload-item.uploading {
  background: #e6f7ff;
}

.upload-item.success {
  background: #f6ffed;
}

.upload-item.error {
  background: #fff2f0;
}

.upload-item-info {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.file-icon {
  font-size: 16px;
  color: #1890ff;
}

.file-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  color: #999;
  font-size: 12px;
}

.upload-item-status {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 150px;
}

.status-text {
  color: #999;
  font-size: 12px;
}

.progress-text {
  font-size: 12px;
  color: #1890ff;
  width: 40px;
}

.success-icon {
  color: #52c41a;
  font-size: 16px;
}

.error-icon {
  color: #ff4d4f;
  font-size: 16px;
}

.error-text {
  color: #ff4d4f;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.upload-item-actions {
  display: flex;
  gap: 4px;
}
</style>
