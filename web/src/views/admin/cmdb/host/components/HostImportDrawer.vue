<template>
  <a-drawer
    :open="open"
    title="导入主机"
    width="760"
    placement="right"
    destroy-on-close
    @close="handleClose"
  >
    <div class="import-layout">
      <div class="import-section">
        <div class="import-section__header">
          <div class="import-section__title">导入文件</div>
          <a-button v-permission="'cmdb:host:import'" :loading="templateLoading" @click="emit('download-template')">
            下载模板
          </a-button>
        </div>
        <a-alert
          type="info"
          show-icon
          style="margin-bottom: 16px"
          message="请先下载模板再填写。模板仅保留导入必需字段；IP、负责人、备注不在模板中，创建后可由系统回填或在详情中补充。"
        />

        <a-upload-dragger
          :before-upload="handleBeforeUpload"
          :file-list="fileList"
          :multiple="false"
          accept=".xlsx,.xls"
          @remove="handleRemove"
        >
          <p class="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p class="ant-upload-text">点击或拖拽文件到此区域</p>
          <p class="ant-upload-hint">仅支持 `.xlsx` / `.xls`，每次导入一个文件。</p>
        </a-upload-dragger>
      </div>

      <div v-if="result" class="import-section">
        <div class="import-section__title">导入结果</div>
        <div class="result-summary">
          <div class="result-item">
            <span>总行数</span>
            <strong>{{ result.total }}</strong>
          </div>
          <div class="result-item success">
            <span>成功</span>
            <strong>{{ result.success_count }}</strong>
          </div>
          <div class="result-item danger">
            <span>失败</span>
            <strong>{{ result.failure_count }}</strong>
          </div>
        </div>

        <a-table
          size="small"
          :columns="resultColumns"
          :data-source="result.rows"
          :pagination="false"
          row-key="row"
          :scroll="{ x: 900, y: 360 }"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'created'">
              <a-tag :color="record.created ? 'success' : 'error'">{{ record.created ? '已创建' : '失败' }}</a-tag>
            </template>
            <template v-else-if="column.key === 'verify_success'">
              <a-tag :color="record.verify_success ? 'success' : 'default'">{{ record.verify_success ? '成功' : '未成功' }}</a-tag>
            </template>
            <template v-else-if="column.key === 'error_message'">
              <span class="cell-text danger-text">{{ record.error_message || '-' }}</span>
            </template>
            <template v-else-if="column.key === 'verify_message'">
              <span class="cell-text">{{ record.verify_message || '-' }}</span>
            </template>
          </template>
        </a-table>
      </div>
    </div>

    <template #footer>
      <a-space>
        <a-button @click="handleClose">关闭</a-button>
        <a-button type="primary" :loading="importing" :disabled="!selectedFile" @click="handleImport">
          开始导入
        </a-button>
      </a-space>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { InboxOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import type { UploadFile } from 'ant-design-vue'
import type { CmdbHostImportResult } from '@/api/cmdb'

const props = defineProps<{
  open: boolean
  importing?: boolean
  templateLoading?: boolean
  result?: CmdbHostImportResult | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', file: File): void
  (e: 'download-template'): void
}>()

const selectedFile = ref<File>()
const fileList = ref<UploadFile[]>([])

const resultColumns = [
  { title: '行号', dataIndex: 'row', key: 'row', width: 70, align: 'center' },
  { title: '主机名称', dataIndex: 'name', key: 'name', width: 150 },
  { title: '创建结果', key: 'created', width: 100, align: 'center' },
  { title: '校验结果', key: 'verify_success', width: 100, align: 'center' },
  { title: '错误信息', key: 'error_message', width: 220 },
  { title: '校验说明', key: 'verify_message' },
]

watch(
  () => props.open,
  (open) => {
    if (!open) {
      selectedFile.value = undefined
      fileList.value = []
    }
  }
)

const handleBeforeUpload = (file: File) => {
  const isExcel = file.name.endsWith('.xlsx') || file.name.endsWith('.xls')
  if (!isExcel) {
    message.warning('仅支持 Excel 文件')
    return false
  }
  selectedFile.value = file
  fileList.value = [{
    uid: String(Date.now()),
    name: file.name,
    status: 'done',
  }]
  return false
}

const handleRemove = () => {
  selectedFile.value = undefined
  fileList.value = []
}

const handleImport = () => {
  if (!selectedFile.value) {
    message.warning('请先选择导入文件')
    return
  }
  emit('submit', selectedFile.value)
}

const handleClose = () => {
  emit('update:open', false)
}
</script>

<style scoped>
.import-layout {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.import-section {
  padding: 16px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.import-section__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.import-section__header .import-section__title {
  margin-bottom: 0;
}

.import-section__title {
  margin-bottom: 14px;
  color: var(--app-text-strong);
  font-size: 14px;
  font-weight: 600;
}

.result-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.result-item {
  padding: 14px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.result-item span {
  color: var(--app-text-muted);
  font-size: 12px;
}

.result-item strong {
  display: block;
  margin-top: 8px;
  color: var(--app-text-strong);
  font-size: 18px;
}

.result-item.success strong {
  color: #389e0d;
}

.result-item.danger strong {
  color: #cf1322;
}

.cell-text {
  color: var(--app-text-color);
  font-size: 12px;
  line-height: 1.6;
  word-break: break-all;
}

.danger-text {
  color: #cf1322;
}

@media (max-width: 900px) {
  .result-summary {
    grid-template-columns: 1fr;
  }
}
</style>
