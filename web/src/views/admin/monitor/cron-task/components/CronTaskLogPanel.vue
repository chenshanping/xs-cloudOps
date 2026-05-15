<template>
  <div class="cron-task-log-panel">
    <div class="cron-task-log-panel__toolbar">
      <a-space wrap :size="[12, 12]">
        <div class="filter-item">
          <span class="filter-item__label">状态</span>
          <a-select v-model:value="searchForm.status" allow-clear placeholder="全部状态" style="width: 132px">
            <a-select-option value="running">执行中</a-select-option>
            <a-select-option value="success">成功</a-select-option>
            <a-select-option value="failure">失败</a-select-option>
            <a-select-option value="skipped">已跳过</a-select-option>
          </a-select>
        </div>
        <div class="filter-item">
          <span class="filter-item__label">触发方式</span>
          <a-select v-model:value="searchForm.triggered_by" allow-clear placeholder="全部方式" style="width: 132px">
            <a-select-option value="schedule">调度</a-select-option>
            <a-select-option value="manual">手动</a-select-option>
          </a-select>
        </div>
        <div class="filter-item">
          <span class="filter-item__label">执行时间</span>
          <a-range-picker
            v-model:value="searchForm.time_range"
            show-time
            format="YYYY-MM-DD HH:mm:ss"
            style="width: 360px"
          />
        </div>
        <a-space :size="8">
          <a-button type="primary" @click="handleSearch">筛选</a-button>
          <a-button @click="handleReset">重置</a-button>
          <a-button @click="fetchData">刷新</a-button>
        </a-space>
      </a-space>
    </div>

    <div class="cron-task-log-panel__surface">
      <div class="cron-task-log-panel__meta">
        <div>
          <div class="meta-label">执行日志</div>
          <div class="meta-title">查看 {{ task.name }} 的最近执行记录</div>
        </div>
        <a-tag color="blue">{{ pagination.total }} 条</a-tag>
      </div>

      <a-table
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        :pagination="paginationConfig"
        :custom-row="customRow"
        :row-class-name="rowClassName"
        row-key="id"
        :scroll="{ x: 1040 }"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
          </template>
          <template v-else-if="column.key === 'triggered_by'">
            <a-tag :color="record.triggered_by === 'manual' ? 'blue' : 'purple'">{{ record.triggered_by === 'manual' ? '手动' : '调度' }}</a-tag>
          </template>
          <template v-else-if="column.key === 'started_at'">{{ formatTime(record.started_at) }}</template>
          <template v-else-if="column.key === 'finished_at'">{{ formatTime(record.finished_at) }}</template>
          <template v-else-if="column.key === 'duration_ms'">{{ formatDuration(record.duration_ms) }}</template>
          <template v-else-if="column.key === 'summary'">
            <span class="summary-text">{{ record.summary || record.error_message || '-' }}</span>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-button type="link" size="small" @click.stop="handleDetail(record)" v-permission="'monitor:cron:logs:view'">详情</a-button>
          </template>
        </template>
        <template #emptyText>
          <a-empty description="当前任务暂无执行日志" />
        </template>
      </a-table>
    </div>

    <div ref="detailAnchorRef" class="cron-task-log-panel__detail">
      <div class="cron-task-log-panel__detail-header">
        <div>
          <div class="meta-label">日志详情</div>
          <div class="meta-title">{{ currentLog ? `执行记录 #${currentLog.id}` : '点击上方记录查看执行详情' }}</div>
        </div>
        <a-button v-if="currentLog" type="link" @click="clearCurrentLog">收起</a-button>
      </div>

      <a-spin :spinning="detailLoading">
        <template v-if="currentLog">
          <div class="detail-hero">
            <div class="detail-hero__route">
              <a-tag :color="getStatusColor(currentLog.status)">{{ getStatusText(currentLog.status) }}</a-tag>
              <strong>{{ currentLog.summary || currentLog.error_message || `任务 ${task.name} 执行详情` }}</strong>
            </div>
            <div class="detail-hero__sub">
              <span>触发方式：{{ currentLog.triggered_by === 'manual' ? '手动' : '调度' }}</span>
              <span>开始时间：{{ formatTime(currentLog.started_at) }}</span>
            </div>
          </div>

          <a-descriptions :column="2" bordered size="small">
            <a-descriptions-item label="任务名称">{{ task.name }}</a-descriptions-item>
            <a-descriptions-item label="注册任务">{{ currentLog.task_code }}</a-descriptions-item>
            <a-descriptions-item label="执行状态">
              <a-tag :color="getStatusColor(currentLog.status)">{{ getStatusText(currentLog.status) }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="触发方式">{{ currentLog.triggered_by === 'manual' ? '手动' : '调度' }}</a-descriptions-item>
            <a-descriptions-item label="开始时间">{{ formatTime(currentLog.started_at) }}</a-descriptions-item>
            <a-descriptions-item label="结束时间">{{ formatTime(currentLog.finished_at) }}</a-descriptions-item>
            <a-descriptions-item label="耗时">{{ formatDuration(currentLog.duration_ms) }}</a-descriptions-item>
            <a-descriptions-item label="触发用户">{{ currentLog.actor_user_id || '-' }}</a-descriptions-item>
          </a-descriptions>

          <div class="detail-grid">
            <div class="detail-card">
              <div class="detail-card__title">执行摘要</div>
              <pre>{{ currentLog.summary || '-' }}</pre>
            </div>
            <div class="detail-card detail-card--danger">
              <div class="detail-card__title">错误信息</div>
              <pre :class="{ 'detail-card__danger-text': currentLog.error_message }">{{ currentLog.error_message || '-' }}</pre>
            </div>
          </div>
        </template>
        <a-empty v-else description="上方点击某条执行记录后，这里会展示完整详情" />
      </a-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, reactive, ref, watch } from 'vue'
import type { Dayjs } from 'dayjs'
import { getCronLogDetail, getCronLogList, type CronLog, type CronTask } from '@/api/cron'
import { formatTime } from '@/utils/format'

const props = defineProps<{
  task: CronTask
  active: boolean
  autoOpenLogId?: number
}>()

const loading = ref(false)
const detailLoading = ref(false)
const tableData = ref<CronLog[]>([])
const currentLog = ref<CronLog | null>(null)
const detailAnchorRef = ref<HTMLElement>()
const initializedTaskId = ref<number>()
const handledAutoOpenKey = ref('')
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const searchForm = reactive({
  status: undefined as string | undefined,
  triggered_by: undefined as string | undefined,
  time_range: null as [Dayjs, Dayjs] | null,
})

const columns = [
  { title: '状态', key: 'status', width: 100, align: 'center' },
  { title: '触发方式', key: 'triggered_by', width: 100, align: 'center' },
  { title: '开始时间', key: 'started_at', width: 180, align: 'center' },
  { title: '结束时间', key: 'finished_at', width: 180, align: 'center' },
  { title: '耗时', key: 'duration_ms', width: 100, align: 'center' },
  { title: '摘要', key: 'summary', ellipsis: true },
  { title: '操作', key: 'action', width: 88, align: 'center', fixed: 'right' },
]

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
}))

const getStatusColor = (status?: string) => ({ success: 'green', failure: 'red', running: 'processing', skipped: 'orange' }[status || ''] || 'default')
const getStatusText = (status?: string) => ({ success: '成功', failure: '失败', running: '执行中', skipped: '已跳过' }[status || ''] || '-')
const formatDuration = (duration?: number) => duration || duration === 0 ? `${duration}ms` : '-'

const resetSearchForm = () => {
  searchForm.status = undefined
  searchForm.triggered_by = undefined
  searchForm.time_range = null
}

const clearCurrentLog = () => {
  currentLog.value = null
}

const fetchData = async () => {
  if (!props.task?.id) {
    return
  }
  loading.value = true
  try {
    const res = await getCronLogList({
      page: pagination.current,
      page_size: pagination.pageSize,
      task_id: props.task.id,
      status: searchForm.status,
      triggered_by: searchForm.triggered_by,
      start_time: searchForm.time_range?.[0]?.format('YYYY-MM-DD HH:mm:ss'),
      end_time: searchForm.time_range?.[1]?.format('YYYY-MM-DD HH:mm:ss'),
    })
    tableData.value = res.data.list || []
    pagination.total = res.data.total
  } finally {
    loading.value = false
  }
}

const openLogDetail = async (logId: number, fallback?: CronLog) => {
  detailLoading.value = true
  if (fallback) {
    currentLog.value = fallback
  }
  try {
    const res = await getCronLogDetail(logId)
    currentLog.value = res.data
    await nextTick()
    detailAnchorRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  } finally {
    detailLoading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleReset = () => {
  resetSearchForm()
  pagination.current = 1
  fetchData()
}

const handleTableChange = (pag: any) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchData()
}

const handleDetail = (record: CronLog) => {
  void openLogDetail(record.id, record)
}

const customRow = (record: CronLog) => ({
  onClick: () => handleDetail(record),
})

const rowClassName = (record: CronLog) => record.id === currentLog.value?.id ? 'cron-task-log-panel__row--active' : ''

watch(
  () => props.task.id,
  () => {
    initializedTaskId.value = undefined
    handledAutoOpenKey.value = ''
    tableData.value = []
    currentLog.value = null
    pagination.current = 1
    pagination.total = 0
    resetSearchForm()
  },
)

watch(
  () => [props.active, props.task.id] as const,
  async ([active, taskId]) => {
    if (!active || !taskId) {
      return
    }
    if (initializedTaskId.value !== taskId) {
      await fetchData()
      initializedTaskId.value = taskId
    }
  },
  { immediate: true },
)

watch(
  () => [props.active, props.task.id, props.autoOpenLogId] as const,
  async ([active, taskId, logId]) => {
    if (!active || !taskId || !logId) {
      return
    }
    const key = `${taskId}:${logId}`
    if (handledAutoOpenKey.value === key) {
      return
    }
    if (initializedTaskId.value !== taskId) {
      await fetchData()
      initializedTaskId.value = taskId
    }
    await openLogDetail(logId)
    handledAutoOpenKey.value = key
  },
  { immediate: true },
)
</script>

<style scoped>
.cron-task-log-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.cron-task-log-panel__toolbar,
.cron-task-log-panel__surface,
.cron-task-log-panel__detail {
  padding: 16px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.cron-task-log-panel__meta,
.cron-task-log-panel__detail-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.meta-label,
.filter-item__label {
  color: var(--app-text-muted);
  font-size: 12px;
}

.meta-title {
  margin-top: 4px;
  color: var(--app-text-strong);
  font-size: 15px;
  font-weight: 600;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.summary-text {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  color: var(--app-text-color);
  text-overflow: ellipsis;
  vertical-align: middle;
  white-space: nowrap;
}

.cron-task-log-panel__surface :deep(.ant-table-row) {
  cursor: pointer;
}

.cron-task-log-panel__surface :deep(.cron-task-log-panel__row--active > td) {
  background: rgba(22, 119, 255, 0.08);
}

.detail-hero {
  margin-bottom: 16px;
  padding: 14px 16px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.detail-hero__route {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--app-text-strong);
  font-size: 15px;
}

.detail-hero__sub {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 8px;
  color: var(--app-text-muted);
  font-size: 12px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 14px;
  margin-top: 16px;
}

.detail-card {
  padding: 12px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.detail-card--danger {
  border-color: rgba(255, 77, 79, 0.24);
}

.detail-card__title {
  margin-bottom: 10px;
  color: var(--app-text-strong);
  font-size: 13px;
  font-weight: 600;
}

.detail-card pre {
  margin: 0;
  color: var(--app-text-color);
  white-space: pre-wrap;
  word-break: break-word;
  font-family: Consolas, monospace;
}

.detail-card__danger-text {
  color: #cf1322;
}

@media (max-width: 900px) {
  .cron-task-log-panel__toolbar,
  .cron-task-log-panel__surface,
  .cron-task-log-panel__detail {
    padding: 14px;
  }
}
</style>
