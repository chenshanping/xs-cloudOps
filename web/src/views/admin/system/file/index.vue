<template>
  <div class="file-page">
    <a-card>
      <a-tabs v-model:activeKey="activeTab">
        <a-tab-pane key="list" tab="文件列表">
          <ProTable
            :title="'文件列表'"
            :columns="columns"
            :data-source="fileList"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            :row-selection="{ selectedRowKeys, onChange: onSelectChange }"
            @search="handleSearch"
            @reset="handleReset"
            @change="handleTableChange"
          >
            <template #search>
              <a-form-item>
                <a-input v-model:value="searchForm.name" placeholder="文件名" allowClear style="width: 200px" />
              </a-form-item>
              <a-form-item>
                <a-select v-model:value="searchForm.ext" placeholder="文件类型" allowClear style="width: 120px">
                  <a-select-option value="">全部</a-select-option>
                  <a-select-option v-for="item in FILE_TYPES" :key="item.value" :value="item.value">
                    {{ item.label }}
                  </a-select-option>
                </a-select>
              </a-form-item>
            </template>

            <template #toolbar>
              <a-button type="primary" danger :disabled="selectedRowKeys.length === 0" @click="handleBatchDelete">
                <DeleteOutlined /> 批量删除
                <span v-if="selectedRowKeys.length > 0">({{ selectedRowKeys.length }})</span>
              </a-button>
            </template>

            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'name'">
                <div class="file-name-cell">
                  <img v-if="isExtImg(record.ext)" :src="record.url" style="width: 50px;height: 50px;" />
                  <component v-else :is="getFileIconComponent(record.ext)" class="file-icon" />
                  <a-tooltip :title="record.name">
                    <span class="file-name">{{ record.name }}</span>
                  </a-tooltip>
                </div>
              </template>
              <template v-if="column.key === 'size'">
                {{ formatFileSize(record.size) }}
              </template>
              <template v-if="column.key === 'ext'">
                <a-tag :color="getFileTypeInfo(record.ext)?.color">{{ getFileTypeInfo(record.ext)?.label || record.ext?.toUpperCase() }}</a-tag>
              </template>
              <template v-if="column.key === 'storage'">
                <a-tag color="blue">
                  {{ getStorageLabel(record) }}
                </a-tag>
              </template>
              <template v-if="column.key === 'created_at'">
                {{ formatTime(record.created_at) }}
              </template>
              <template v-if="column.key === 'action'">
                <a-button type="link" size="small" @click="handlePreview(record)">预览</a-button>
                <a-button type="link" size="small" @click="handleCopyUrl(record)">复制链接</a-button>
                <a-popconfirm title="确定删除吗？" @confirm="handleDelete(record)">
                  <a-button type="link" size="small" danger>删除</a-button>
                </a-popconfirm>
              </template>
            </template>
          </ProTable>
        </a-tab-pane>

        <a-tab-pane key="upload" tab="上传文件">
          <FileUpload
            ref="fileUploadRef"
            :multiple="true"
            @success="handleUploadSuccess"
          />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <FilePreview
      v-model:open="previewVisible"
      :url="previewFile?.url || ''"
      :name="previewFile?.name || ''"
      :ext="previewFile?.ext || ''"
      :size="previewFile?.size"
      :mime-type="previewFile?.mime_type"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import ProTable from '@/components/ProTable.vue'
import {
  FileImageOutlined,
  FilePdfOutlined,
  FileWordOutlined,
  FileExcelOutlined,
  FilePptOutlined,
  FileZipOutlined,
  FileTextOutlined,
  VideoCameraOutlined,
  AudioOutlined,
  FileOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import FileUpload from '@/components/FileUpload.vue'
import FilePreview from '@/components/FilePreview.vue'
import type { FileInfo } from '@/types/file'
import { getFileList, deleteFile, batchDeleteFiles } from '@/api/file'
import { useConfigStore } from '@/store/config'
import { formatFileSize } from '@/utils/upload'
import { formatTime } from '@/utils/format'
import { storageTypeOptions } from '@/types/storage'

const configStore = useConfigStore()

const activeTab = ref('list')
const loading = ref(false)
const fileList = ref<FileInfo[]>([])
const fileUploadRef = ref()
const selectedRowKeys = ref<number[]>([])
const previewVisible = ref(false)
const previewFile = ref<FileInfo | null>(null)

const searchForm = reactive({
  name: '',
  ext: '',
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
})

const currentStorageLabel = computed(() => {
  const currentType = configStore.get('storage_type') || 'local'
  return storageTypeOptions.find((item) => item.value === currentType)?.label || currentType
})

const FILE_TYPES = [
  { label: '图片', value: 'jpg,jpeg,png,gif,webp,bmp,svg', color: 'green' },
  { label: 'PDF', value: 'pdf', color: 'red' },
  { label: '文档', value: 'doc,docx,txt', color: 'blue' },
  { label: '视频', value: 'mp4,avi,mov,wmv,flv,mkv', color: 'purple' },
  { label: '压缩包', value: 'zip,rar,7z,tar,gz', color: 'orange' },
  { label: '音频', value: 'mp3,wav,flac,aac', color: 'cyan' },
]

const columns = [
  { title: '文件名', key: 'name', ellipsis: true },
  { title: '大小', key: 'size', width: 100 },
  { title: '类型', key: 'ext', width: 80 },
  { title: '存储', key: 'storage', width: 160 },
  { title: '上传时间', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 200 },
]

const getFileTypeInfo = (ext: string) => {
  if (!ext) return null
  const lowerExt = ext.toLowerCase()
  return FILE_TYPES.find(type => type.value.split(',').includes(lowerExt))
}

const isExtImg = (ext: string) => {
  const imgType = FILE_TYPES.find(t => t.label === '图片')
  return imgType ? imgType.value.split(',').includes(ext?.toLowerCase()) : false
}

const getFileIconComponent = (ext: string) => {
  const iconMap: Record<string, any> = {
    jpg: FileImageOutlined,
    jpeg: FileImageOutlined,
    png: FileImageOutlined,
    gif: FileImageOutlined,
    webp: FileImageOutlined,
    bmp: FileImageOutlined,
    svg: FileImageOutlined,
    pdf: FilePdfOutlined,
    doc: FileWordOutlined,
    docx: FileWordOutlined,
    xls: FileExcelOutlined,
    xlsx: FileExcelOutlined,
    ppt: FilePptOutlined,
    pptx: FilePptOutlined,
    zip: FileZipOutlined,
    rar: FileZipOutlined,
    '7z': FileZipOutlined,
    txt: FileTextOutlined,
    mp4: VideoCameraOutlined,
    avi: VideoCameraOutlined,
    mov: VideoCameraOutlined,
    mp3: AudioOutlined,
    wav: AudioOutlined,
  }
  return iconMap[ext?.toLowerCase()] || FileOutlined
}

const getStorageLabel = (record: FileInfo) => {
  if (record.storage_type) {
    return storageTypeOptions.find((item) => item.value === record.storage_type)?.label || record.storage_type
  }
  return `系统配置(${currentStorageLabel.value})`
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getFileList({
      page: pagination.current,
      page_size: pagination.pageSize,
      name: searchForm.name,
      ext: searchForm.ext,
    })
    fileList.value = res.data.list
    pagination.total = res.data.total
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchList()
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.ext = ''
  pagination.current = 1
  fetchList()
}

const handleTableChange = (pag: any) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

const handlePreview = (record: FileInfo) => {
  previewFile.value = record
  previewVisible.value = true
}

const handleCopyUrl = async (record: FileInfo) => {
  try {
    await navigator.clipboard.writeText(record.url)
    message.success('链接已复制到剪贴板')
  } catch {
    message.error('复制失败')
  }
}

const handleDelete = async (record: FileInfo) => {
  await deleteFile(record.id)
  message.success('删除成功')
  fetchList()
}

const onSelectChange = (keys: number[]) => {
  selectedRowKeys.value = keys
}

const handleBatchDelete = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的文件')
    return
  }
  Modal.confirm({
    title: '确认批量删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个文件吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        const res = await batchDeleteFiles(selectedRowKeys.value)
        if (res.data?.failed_count > 0) {
          message.warning(`成功删除 ${res.data.success_count} 个，失败 ${res.data.failed_count} 个`)
        } else {
          message.success('批量删除成功')
        }
        selectedRowKeys.value = []
        fetchList()
      } catch {
        // 错误已由 request 拦截器处理
      }
    }
  })
}

const handleUploadSuccess = (file: FileInfo) => {
  message.success(`${file.name} 上传成功`)
  activeTab.value = 'list'
  fetchList()
}

fetchList()
</script>

<style scoped>
.file-page {
  height: 100%;
}

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 18px;
  color: #1890ff;
}

.file-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
