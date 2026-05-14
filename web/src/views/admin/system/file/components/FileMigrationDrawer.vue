<template>
  <a-drawer
    :open="open"
    title="文件迁移"
    width="80%"
    destroy-on-close
    @close="handleClose"
  >
    <a-space direction="vertical" style="width: 100%" :size="16">
      <a-alert
        type="info"
        show-icon
        message="文件迁移只处理已上传的历史文件，不会修改系统当前的默认上传存储配置。"
        :description="migrationNoticeDescription"
      />

      <a-form layout="vertical">
        <a-form-item label="迁移范围">
          <a-radio-group v-model:value="form.scope" :options="scopeOptions" option-type="button" button-style="solid" />
        </a-form-item>
        <a-form-item
          label="源存储（历史文件当前所在位置）"
          extra="这里用于识别历史文件现在所在的存储位置，不会修改后续新上传文件的默认去向。"
        >
          <a-select
            v-model:value="form.sourceStorageType"
            placeholder="请选择源存储"
            :options="storageOptionsWithBucket"
          />
        </a-form-item>
        <a-form-item
          label="目标存储（迁移后写入的位置）"
          extra="这里表示本次迁移后历史文件要搬去的位置；选择这里不会自动修改系统默认上传存储。"
        >
          <a-select
            v-model:value="form.targetStorageType"
            placeholder="请选择目标存储"
            :options="storageOptionsWithBucket"
          />
        </a-form-item>
      </a-form>

      <a-card v-if="sameStorageSelected && form.sourceStorageType !== 'local'" size="small" title="源存储配置（旧桶）" style="margin-bottom: 12px">
        <a-form layout="vertical">
          <a-form-item label="配置来源">
            <a-radio-group v-model:value="sourceConfigMode" option-type="button" button-style="solid">
              <a-radio-button value="system">使用系统配置</a-radio-button>
              <a-radio-button value="custom">手动填写</a-radio-button>
            </a-radio-group>
            <div v-if="sourceConfigMode === 'system'" style="margin-top: 6px; font-size: 12px; color: #64748b">
              已加载系统配置的连接信息（AK/SK/Endpoint 等），桶名已清空，请填写旧桶名称。
              <template v-if="targetBucketLabel">当前目标桶：<strong>{{ targetBucketLabel }}</strong></template>
            </div>
          </a-form-item>
          <template v-if="form.sourceStorageType === 'aliyun'">
            <a-form-item label="Endpoint">
              <a-input v-model:value="sourceConfig.endpoint" placeholder="如：oss-cn-hangzhou.aliyuncs.com" />
            </a-form-item>
            <a-form-item label="AccessKey ID">
              <a-input v-model:value="sourceConfig.access_key_id" />
            </a-form-item>
            <a-form-item label="AccessKey Secret">
              <a-input-password v-model:value="sourceConfig.access_key_secret" />
            </a-form-item>
            <a-form-item label="Bucket（旧桶）" :extra="targetBucketLabel ? `目标桶：${targetBucketLabel}，请填写不同的旧桶名` : ''">
              <a-input v-model:value="sourceConfig.bucket_name" placeholder="请输入旧桶名称" />
            </a-form-item>
            <a-form-item label="Region">
              <a-input v-model:value="sourceConfig.region" placeholder="如：cn-hangzhou" />
            </a-form-item>
          </template>
          <template v-else-if="form.sourceStorageType === 'tencent'">
            <a-form-item label="Region">
              <a-input v-model:value="sourceConfig.region" placeholder="如：ap-guangzhou" />
            </a-form-item>
            <a-form-item label="SecretId">
              <a-input v-model:value="sourceConfig.secret_id" />
            </a-form-item>
            <a-form-item label="SecretKey">
              <a-input-password v-model:value="sourceConfig.secret_key" />
            </a-form-item>
            <a-form-item label="Bucket（旧桶）" :extra="targetBucketLabel ? `目标桶：${targetBucketLabel}，请填写不同的旧桶名` : ''">
              <a-input v-model:value="sourceConfig.bucket" placeholder="请输入旧桶名称" />
            </a-form-item>
            <a-form-item label="AppID">
              <a-input v-model:value="sourceConfig.app_id" />
            </a-form-item>
          </template>
          <template v-else-if="form.sourceStorageType === 'minio'">
            <a-form-item label="Endpoint">
              <a-input v-model:value="sourceConfig.endpoint" placeholder="如：127.0.0.1:9000" />
            </a-form-item>
            <a-form-item label="AccessKey ID">
              <a-input v-model:value="sourceConfig.access_key_id" />
            </a-form-item>
            <a-form-item label="SecretAccessKey">
              <a-input-password v-model:value="sourceConfig.secret_access_key" />
            </a-form-item>
            <a-form-item label="Bucket（旧桶）" :extra="targetBucketLabel ? `目标桶：${targetBucketLabel}，请填写不同的旧桶名` : ''">
              <a-input v-model:value="sourceConfig.bucket_name" placeholder="请输入旧桶名称" />
            </a-form-item>
            <a-form-item label="使用SSL">
              <a-switch v-model:checked="sourceConfig.use_ssl" />
            </a-form-item>
          </template>
        </a-form>
      </a-card>

      <a-alert
        v-if="selectionAlert"
        :type="selectionAlert.type"
        show-icon
        :message="selectionAlert.message"
        :description="selectionAlert.description"
      />

      <div class="migration-actions">
        <a-space>
          <a-button :disabled="!canPreview" :loading="previewLoading" @click="handlePreview" v-permission="'system:file:migrate'">
            预检查统计
          </a-button>
          <a-button type="primary" :disabled="!canExecute" :loading="executeLoading" @click="handleExecute" v-permission="'system:file:migrate'">
            开始迁移
          </a-button>
        </a-space>
      </div>

      <template v-if="taskStatus">
        <a-alert
          :type="taskAlertType"
          show-icon
          :message="taskStatus.message || taskStatusLabel"
        />

        <div class="task-panel">
          <div class="task-panel__header">
            <a-tag :color="taskTagColor">{{ taskStatusLabel }}</a-tag>
            <span class="task-panel__meta">
              {{ getSourceStorageLabel(taskStatus.source_storage_type) }} → {{ getStorageLabel(taskStatus.target_storage_type) }}
            </span>
          </div>
          <div v-if="taskStatus.current_file_name" class="task-panel__current">
            当前文件：{{ taskStatus.current_file_name }}
          </div>
          <div class="task-progress">
            <div class="task-progress__item">
              <div class="task-progress__label">文件进度</div>
              <a-progress :percent="countProgressPercent" :status="taskProgressStatus">
                <template #format>
                  {{ taskStatus.processed_count }}/{{ taskStatus.total_count }}
                </template>
              </a-progress>
            </div>
            <div class="task-progress__item">
              <div class="task-progress__label">容量进度</div>
              <a-progress :percent="sizeProgressPercent" :status="taskProgressStatus">
                <template #format>
                  {{ formatFileSize(taskStatus.processed_size) }}/{{ formatFileSize(taskStatus.total_size) }}
                </template>
              </a-progress>
            </div>
          </div>
        </div>
      </template>

      <template v-if="summarySource">
        <a-alert
          :type="summaryAlertType"
          show-icon
          :message="summaryText"
        />

        <a-row :gutter="[12, 12]" class="summary-grid">
          <a-col :span="12">
            <div class="summary-card">
              <div class="summary-card__label">总文件</div>
              <div class="summary-card__value">{{ summarySource.total_count }}</div>
              <div class="summary-card__meta">{{ formatFileSize(summarySource.total_size) }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="summary-card">
              <div class="summary-card__label">可迁移</div>
              <div class="summary-card__value">{{ summarySource.pending_count }}</div>
              <div class="summary-card__meta">{{ formatFileSize(summarySource.pending_size) }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="summary-card">
              <div class="summary-card__label">冲突</div>
              <div class="summary-card__value">{{ summarySource.conflict_count }}</div>
              <div class="summary-card__meta">{{ formatFileSize(summarySource.conflict_size) }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="summary-card">
              <div class="summary-card__label">源文件缺失</div>
              <div class="summary-card__value">{{ summarySource.missing_source_count }}</div>
              <div class="summary-card__meta">{{ formatFileSize(summarySource.missing_source_size) }}</div>
            </div>
          </a-col>
        </a-row>

        <a-table
          :columns="columns"
          :data-source="summarySource.items"
          :pagination="false"
          size="small"
          :row-key="getRowKey"
          :scroll="{ y: 360 }"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'storage'">
              <a-space>
                <a-tag>{{ getSourceStorageLabel(record.source_storage_type) }}</a-tag>
                <span>→</span>
                <a-tag color="blue">{{ getStorageLabel(record.target_storage_type) }}</a-tag>
              </a-space>
            </template>
            <template v-else-if="column.key === 'action'">
              <a-tag :color="getActionColor(record.action)">
                {{ getActionLabel(record.action) }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'message'">
              <a-typography-text :type="record.action === 'FAILED' ? 'danger' : undefined">
                {{ record.message }}
              </a-typography-text>
            </template>
          </template>
        </a-table>
      </template>
    </a-space>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { executeFileMigration, getCurrentFileMigrationTask, previewFileMigration } from '@/api/file'
import type {
  FileMigrationItem,
  FileMigrationRequest,
  FileMigrationResult,
  FileMigrationScope,
  FileMigrationTaskStatus,
} from '@/types/file'
import type { StorageType } from '@/types/storage'
import { storageTypeOptions } from '@/types/storage'
import { formatFileSize } from '@/utils/upload'
import { useConfigStore } from '@/store/config'
import { getDefaultStorageConfig, getStorageBucketLabel, parseStorageConfig, STORAGE_CONFIG_KEY_MAP } from '@/views/admin/system/config/storage-config-state'
import { shouldRestoreMigrationTaskOnOpen } from './file-migration-task-state.js'

const props = defineProps<{
  open: boolean
  selectedIds: number[]
  currentDefaultStorageType: string
  currentFilters: {
    name?: string
    ext?: string
    referenced?: boolean
  }
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  success: []
}>()

const configStore = useConfigStore()

const previewLoading = ref(false)
const executeLoading = ref(false)
const result = ref<FileMigrationResult | null>(null)
const taskStatus = ref<FileMigrationTaskStatus | null>(null)
const pollingTimer = ref<number | null>(null)
const watchTaskId = ref('')
const terminalNoticeTaskKey = ref('')

const form = reactive({
  scope: 'all' as FileMigrationScope,
  sourceStorageType: '',
  targetStorageType: '',
})

const sourceConfigMode = ref<'system' | 'custom'>('system')
const sourceConfig = reactive<Record<string, any>>({})

const getBucketFieldKey = (type: string): string => {
  if (type === 'tencent') return 'bucket'
  if (type === 'local') return 'base_path'
  return 'bucket_name'
}

const targetBucketLabel = computed(() => {
  if (!form.targetStorageType) return ''
  const configKey = STORAGE_CONFIG_KEY_MAP[form.targetStorageType as StorageType]
  const configJSON = configKey ? configStore.get(configKey) : ''
  return getStorageBucketLabel(form.targetStorageType as StorageType, configJSON)
})

const fillSourceConfigFromSystem = async (type: string) => {
  Object.keys(sourceConfig).forEach((key) => delete sourceConfig[key])
  if (!type) return
  await configStore.loadConfigs(false, 'all')
  const configKey = STORAGE_CONFIG_KEY_MAP[type as StorageType]
  const configJSON = configKey ? configStore.get(configKey) : ''
  const parsed = parseStorageConfig(type as StorageType, configJSON)
  Object.assign(sourceConfig, parsed)
  // 清空桶名字段，迫使用户填写旧桶名
  const bucketKey = getBucketFieldKey(type)
  sourceConfig[bucketKey] = ''
}

const fillSourceConfigEmpty = (type: string) => {
  Object.keys(sourceConfig).forEach((key) => delete sourceConfig[key])
  if (!type) return
  Object.assign(sourceConfig, getDefaultStorageConfig(type as StorageType))
}

watch(
  () => form.sourceStorageType,
  (type) => {
    sourceConfigMode.value = 'system'
    fillSourceConfigFromSystem(type)
  },
)

watch(sourceConfigMode, (mode) => {
  if (mode === 'system') {
    fillSourceConfigFromSystem(form.sourceStorageType)
  } else {
    fillSourceConfigEmpty(form.sourceStorageType)
  }
})

const hasActiveFilters = computed(() =>
  !!props.currentFilters.name || !!props.currentFilters.ext || props.currentFilters.referenced === true
)

const storageOptionsWithBucket = computed(() =>
  storageTypeOptions.map((opt) => {
    const configKey = STORAGE_CONFIG_KEY_MAP[opt.value as StorageType]
    const configJSON = configKey ? configStore.get(configKey) : ''
    const bucket = getStorageBucketLabel(opt.value as StorageType, configJSON)
    return {
      ...opt,
      label: bucket ? `${opt.label} (${bucket})` : opt.label,
    }
  })
)

const scopeOptions = computed(() => [
  { label: '全部文件', value: 'all' },
  { label: '当前筛选结果', value: 'filter', disabled: !hasActiveFilters.value },
  { label: '手动选中项', value: 'selected', disabled: props.selectedIds.length === 0 },
])

const columns = [
  { title: '文件ID', dataIndex: 'file_id', key: 'file_id', width: 88 },
  { title: '文件名', dataIndex: 'file_name', key: 'file_name', ellipsis: true },
  { title: '迁移方向', key: 'storage', width: 220 },
  { title: '结果', key: 'action', width: 110 },
  { title: '说明', key: 'message', ellipsis: true },
]

const isTaskActive = computed(() => taskStatus.value?.status === 'SCANNING' || taskStatus.value?.status === 'RUNNING')

const canPreview = computed(() => {
  if (isTaskActive.value) {
    return false
  }
  if (!form.sourceStorageType || !form.targetStorageType) {
    return false
  }
  if (form.scope === 'selected') {
    return props.selectedIds.length > 0
  }
  if (form.scope === 'filter') {
    return hasActiveFilters.value
  }
  return true
})

const canExecute = computed(() => canPreview.value && !isTaskActive.value)
const sameStorageSelected = computed(() => {
  return !!form.sourceStorageType && form.sourceStorageType === form.targetStorageType
})

const scopeDescription = computed(() => {
  switch (form.scope) {
    case 'selected':
      return `将迁移手动选中的 ${props.selectedIds.length} 个文件`
    case 'filter':
      return hasActiveFilters.value ? '将迁移当前筛选结果中的文件' : '当前没有有效筛选条件'
    default:
      return '将迁移指定源存储下的全部文件'
  }
})
const currentDefaultStorageLabel = computed(() => getStorageLabel(props.currentDefaultStorageType))
const migrationNoticeDescription = computed(() => {
  return `当前默认上传存储：${currentDefaultStorageLabel.value}。${scopeDescription.value}`
})
const selectionAlert = computed(() => {
  if (!form.sourceStorageType || !form.targetStorageType) {
    return null
  }

  if (sameStorageSelected.value) {
    return {
      type: 'info',
      message: '同类型跃桶迁移',
      description: '源存储与目标存储类型相同，请在下方填写旧桶的连接信息，系统会从旧桶读取文件并写入当前配置的桶。',
    }
  }

  if (props.currentDefaultStorageType && form.targetStorageType === props.currentDefaultStorageType) {
    return {
      type: 'success',
      message: '目标存储已是当前默认上传位置',
      description: '迁移完成后，历史文件会与后续新上传文件存放在同一存储。',
    }
  }

  if (props.currentDefaultStorageType && form.sourceStorageType === props.currentDefaultStorageType) {
    return {
      type: 'info',
      message: '当前默认上传存储仍指向源存储',
      description: '本次迁移不会自动切换默认上传存储；如需修改后续上传位置，请前往 系统设置 > 文件设置。',
    }
  }

  return null
})

const summarySource = computed(() => taskStatus.value || result.value)

const summaryText = computed(() => {
  const current = summarySource.value
  if (!current) {
    return ''
  }
  if (taskStatus.value) {
    return `共 ${current.total_count} 个文件，已处理 ${taskStatus.value.processed_count}，成功 ${taskStatus.value.migrated_count}，告警 ${taskStatus.value.warning_count}，失败 ${taskStatus.value.failed_count}`
  }
  return `共匹配 ${current.total_count} 个文件，待迁移 ${current.pending_count}，冲突 ${current.conflict_count}，源文件缺失 ${current.missing_source_count}，其他跳过 ${current.skipped_count}`
})

const summaryAlertType = computed(() => {
  if (taskStatus.value) {
    if (taskStatus.value.status === 'FAILED') {
      return 'error'
    }
    if (taskStatus.value.failed_count > 0 || taskStatus.value.warning_count > 0) {
      return 'warning'
    }
    if (taskStatus.value.status === 'SUCCESS') {
      return 'success'
    }
    return 'info'
  }
  if (!result.value) {
    return 'info'
  }
  if (result.value.conflict_count > 0 || result.value.missing_source_count > 0 || result.value.failed_count > 0) {
    return 'warning'
  }
  return 'success'
})

const taskStatusLabel = computed(() => {
  switch (taskStatus.value?.status) {
    case 'SCANNING':
      return '预检查中'
    case 'RUNNING':
      return '迁移中'
    case 'SUCCESS':
      return '已完成'
    case 'FAILED':
      return '执行失败'
    default:
      return ''
  }
})

const taskAlertType = computed(() => {
  switch (taskStatus.value?.status) {
    case 'FAILED':
      return 'error'
    case 'SUCCESS':
      return taskStatus.value.failed_count > 0 || taskStatus.value.warning_count > 0 ? 'warning' : 'success'
    default:
      return 'info'
  }
})

const taskTagColor = computed(() => {
  switch (taskStatus.value?.status) {
    case 'SCANNING':
      return 'processing'
    case 'RUNNING':
      return 'blue'
    case 'SUCCESS':
      return taskStatus.value.failed_count > 0 || taskStatus.value.warning_count > 0 ? 'warning' : 'success'
    case 'FAILED':
      return 'error'
    default:
      return 'default'
  }
})

const taskProgressStatus = computed(() => {
  if (taskStatus.value?.status === 'FAILED') {
    return 'exception'
  }
  if (taskStatus.value?.status === 'SUCCESS') {
    return taskStatus.value.failed_count > 0 ? 'exception' : 'success'
  }
  return 'active'
})

const countProgressPercent = computed(() => {
  if (!taskStatus.value?.total_count) {
    return taskStatus.value?.status === 'SUCCESS' ? 100 : 0
  }
  return Math.min(100, Math.round((taskStatus.value.processed_count / taskStatus.value.total_count) * 100))
})

const sizeProgressPercent = computed(() => {
  if (!taskStatus.value?.total_size) {
    return taskStatus.value?.status === 'SUCCESS' ? 100 : 0
  }
  return Math.min(100, Math.round((taskStatus.value.processed_size / taskStatus.value.total_size) * 100))
})

watch(
  () => props.open,
  async (open) => {
    if (!open) {
      resetState()
      return
    }
    await configStore.loadConfigs(false, 'all')
    await fetchCurrentTask()
  }
)

watch(
  () => [...props.selectedIds],
  () => {
    result.value = null
    if (form.scope === 'selected' && props.selectedIds.length === 0) {
      form.scope = hasActiveFilters.value ? 'filter' : 'all'
    }
  }
)

watch(
  () => hasActiveFilters.value,
  (value) => {
    if (!value && form.scope === 'filter') {
      form.scope = props.selectedIds.length > 0 ? 'selected' : 'all'
    }
  }
)

onBeforeUnmount(() => {
  stopPolling()
})

const getStorageLabel = (value: string) => {
  const typeLabel = storageTypeOptions.find((item) => item.value === value)?.label || value || '-'
  const configKey = STORAGE_CONFIG_KEY_MAP[value as StorageType]
  const configJSON = configKey ? configStore.get(configKey) : ''
  const bucket = getStorageBucketLabel(value as StorageType, configJSON)
  return bucket ? `${typeLabel} (${bucket})` : typeLabel
}

const sourceBucketFromConfig = computed(() => {
  if (!sameStorageSelected.value || !form.sourceStorageType) return ''
  const bucketKey = getBucketFieldKey(form.sourceStorageType)
  return sourceConfig[bucketKey] || ''
})

const getSourceStorageLabel = (value: string) => {
  if (sameStorageSelected.value && sourceBucketFromConfig.value) {
    const typeLabel = storageTypeOptions.find((item) => item.value === value)?.label || value || '-'
    return `${typeLabel} (${sourceBucketFromConfig.value})`
  }
  return getStorageLabel(value)
}

const getActionLabel = (action: string) => {
  switch (action) {
    case 'PENDING':
      return '待迁移'
    case 'SKIP':
      return '跳过'
    case 'CONFLICT':
      return '冲突'
    case 'MISSING_SOURCE':
      return '源文件缺失'
    case 'MIGRATED':
      return '成功'
    case 'WARNING':
      return '告警'
    case 'FAILED':
      return '失败'
    default:
      return action
  }
}

const getActionColor = (action: string) => {
  switch (action) {
    case 'PENDING':
      return 'processing'
    case 'SKIP':
      return 'default'
    case 'CONFLICT':
      return 'warning'
    case 'MISSING_SOURCE':
      return 'volcano'
    case 'MIGRATED':
      return 'success'
    case 'WARNING':
      return 'warning'
    case 'FAILED':
      return 'error'
    default:
      return 'default'
  }
}

const getRowKey = (record: FileMigrationItem) => `${record.file_id}-${record.action}-${record.message}`

const buildPayload = (): FileMigrationRequest => {
  const payload: FileMigrationRequest = {
    scope: form.scope,
    source_storage_type: form.sourceStorageType,
    target_storage_type: form.targetStorageType,
  }

  if (sameStorageSelected.value) {
    payload.source_config = JSON.stringify(sourceConfig)
  }

  if (form.scope === 'selected') {
    payload.ids = props.selectedIds
  }

  if (form.scope === 'filter') {
    payload.filters = {
      name: props.currentFilters.name || '',
      ext: props.currentFilters.ext || '',
      referenced: props.currentFilters.referenced,
    }
  }

  return payload
}

const getInvalidMigrationWarning = () => {
  return '请先完善迁移范围和源/目标存储'
}

const handlePreview = async () => {
  if (!canPreview.value) {
    message.warning(getInvalidMigrationWarning())
    return
  }
  previewLoading.value = true
  try {
    const res = await previewFileMigration(buildPayload())
    result.value = res.data
    taskStatus.value = null
    stopPolling()
  } catch {
    // 请求层已统一提示错误，这里阻止错误冒泡到页面级 ErrorBoundary
  } finally {
    previewLoading.value = false
  }
}

const handleExecute = async () => {
  if (!canExecute.value) {
    message.warning(getInvalidMigrationWarning())
    return
  }
  executeLoading.value = true
  try {
    const res = await executeFileMigration(buildPayload())
    taskStatus.value = res.data
    result.value = null
    watchTaskId.value = res.data.task_id
    message.success('迁移任务已启动')
    startPolling()
  } catch {
    // 请求层已统一提示错误，这里阻止错误冒泡到页面级 ErrorBoundary
  } finally {
    executeLoading.value = false
  }
}

const handleClose = () => {
  emit('update:open', false)
}

const startPolling = () => {
  if (pollingTimer.value !== null) {
    return
  }
  pollingTimer.value = window.setInterval(() => {
    fetchCurrentTask(true)
  }, 2000)
}

const stopPolling = () => {
  if (pollingTimer.value !== null) {
    window.clearInterval(pollingTimer.value)
    pollingTimer.value = null
  }
}

const maybeNotifyTaskFinished = (status: FileMigrationTaskStatus) => {
  if (status.status !== 'SUCCESS' && status.status !== 'FAILED') {
    return
  }

  const taskKey = `${status.task_id}-${status.finished_at}`
  if (terminalNoticeTaskKey.value === taskKey || watchTaskId.value !== status.task_id) {
    return
  }

  terminalNoticeTaskKey.value = taskKey
  emit('success')
  if (status.status === 'FAILED' || status.failed_count > 0) {
    message.warning(status.message || '迁移完成，但存在失败项')
    return
  }
  if (status.warning_count > 0) {
    message.warning('迁移完成，请检查告警详情')
    return
  }
  message.success('文件迁移完成')
}

const fetchCurrentTask = async (notifyOnFinish = false) => {
  try {
    const res = await getCurrentFileMigrationTask()
    if (!notifyOnFinish && !shouldRestoreMigrationTaskOnOpen(res.data)) {
      taskStatus.value = null
      stopPolling()
      return
    }

    taskStatus.value = res.data

    if (!res.data) {
      stopPolling()
      return
    }

    if (res.data.status === 'SCANNING' || res.data.status === 'RUNNING') {
      if (!watchTaskId.value) {
        watchTaskId.value = res.data.task_id
      }
      startPolling()
      return
    }

    stopPolling()
    if (notifyOnFinish) {
      maybeNotifyTaskFinished(res.data)
    }
  } catch {
    stopPolling()
  }
}

const resetState = () => {
  stopPolling()
  form.scope = 'all'
  form.sourceStorageType = ''
  form.targetStorageType = ''
  sourceConfigMode.value = 'system'
  Object.keys(sourceConfig).forEach((key) => delete sourceConfig[key])
  result.value = null
  taskStatus.value = null
  watchTaskId.value = ''
}
</script>

<style scoped>
.migration-actions {
  display: flex;
  justify-content: flex-end;
}

.task-panel {
  padding: 16px;
  border: 1px solid #eef2f7;
  border-radius: 12px;
  background: #fafcff;
}

.task-panel__header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.task-panel__meta,
.task-panel__current {
  color: #64748b;
  font-size: 13px;
}

.task-progress {
  display: grid;
  gap: 12px;
  margin-top: 12px;
}

.task-progress__label {
  margin-bottom: 6px;
  color: #0f172a;
  font-size: 13px;
  font-weight: 500;
}

.summary-grid {
  margin: 0;
}

.summary-card {
  height: 100%;
  padding: 14px 16px;
  border: 1px solid #eef2f7;
  border-radius: 12px;
  background: #fff;
}

.summary-card__label {
  color: #64748b;
  font-size: 13px;
}

.summary-card__value {
  margin-top: 8px;
  color: #0f172a;
  font-size: 24px;
  font-weight: 600;
  line-height: 1.2;
}

.summary-card__meta {
  margin-top: 4px;
  color: #94a3b8;
  font-size: 12px;
}
</style>
