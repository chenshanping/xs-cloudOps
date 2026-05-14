<template>
  <PageWrapper class="cron-log-page">
    <div class="cron-log-page__content">
      <ProTable
        title="任务执行日志"
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        :scroll="{ x: 1120 }"
        @search="handleSearch"
        @reset="handleReset"
        @change="handleTableChange"
      >
        <template #search>
          <a-form-item label="任务编码"><a-input v-model:value="searchForm.task_code" placeholder="注册任务编码" allow-clear style="width: 180px" /></a-form-item>
          <a-form-item label="状态">
            <a-select v-model:value="searchForm.status" placeholder="请选择" allow-clear style="width: 130px">
              <a-select-option value="running">执行中</a-select-option>
              <a-select-option value="success">成功</a-select-option>
              <a-select-option value="failure">失败</a-select-option>
              <a-select-option value="skipped">已跳过</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="触发方式">
            <a-select v-model:value="searchForm.triggered_by" placeholder="请选择" allow-clear style="width: 130px">
              <a-select-option value="schedule">调度</a-select-option>
              <a-select-option value="manual">手动</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="执行时间">
            <a-range-picker
              v-model:value="searchForm.time_range"
              show-time
              format="YYYY-MM-DD HH:mm:ss"
              style="width: 360px"
            />
          </a-form-item>
        </template>

        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'"><a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag></template>
          <template v-else-if="column.key === 'triggered_by'"><a-tag :color="record.triggered_by === 'manual' ? 'blue' : 'purple'">{{ record.triggered_by === 'manual' ? '手动' : '调度' }}</a-tag></template>
          <template v-else-if="column.key === 'started_at'">{{ formatTime(record.started_at) }}</template>
          <template v-else-if="column.key === 'finished_at'">{{ formatTime(record.finished_at) }}</template>
          <template v-else-if="column.key === 'duration_ms'">{{ formatDuration(record.duration_ms) }}</template>
          <template v-else-if="column.key === 'summary'"><span class="ellipsis-text">{{ record.summary || record.error_message || '-' }}</span></template>
          <template v-else-if="column.key === 'action'"><a-button type="link" size="small" @click="handleDetail(record)" v-permission="'monitor:cron:logs:view'">详情</a-button></template>
        </template>
      </ProTable>

      <a-drawer :open="detailVisible" width="620" title="执行日志详情" @close="detailVisible = false">
        <a-descriptions v-if="currentLog" :column="1" bordered size="small">
          <a-descriptions-item label="注册任务">{{ currentLog.task_code }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="getStatusColor(currentLog.status)">{{ getStatusText(currentLog.status) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="触发方式">{{ currentLog.triggered_by === 'manual' ? '手动' : '调度' }}</a-descriptions-item>
          <a-descriptions-item label="开始时间">{{ formatTime(currentLog.started_at) }}</a-descriptions-item>
          <a-descriptions-item label="结束时间">{{ formatTime(currentLog.finished_at) }}</a-descriptions-item>
          <a-descriptions-item label="耗时">{{ formatDuration(currentLog.duration_ms) }}</a-descriptions-item>
          <a-descriptions-item label="摘要"><pre>{{ currentLog.summary || '-' }}</pre></a-descriptions-item>
          <a-descriptions-item label="错误信息"><pre :class="{ danger: currentLog.error_message }">{{ currentLog.error_message || '-' }}</pre></a-descriptions-item>
        </a-descriptions>
      </a-drawer>
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import type { Dayjs } from 'dayjs'
import ProTable from '@/components/ProTable.vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import { getCronLogDetail, getCronLogList, type CronLog } from '@/api/cron'
import { formatTime } from '@/utils/format'

const route = useRoute()
const loading = ref(false)
const tableData = ref<CronLog[]>([])
const currentLog = ref<CronLog | null>(null)
const detailVisible = ref(false)
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const searchForm = reactive({
  task_id: undefined as number | undefined,
  task_code: '',
  status: undefined as string | undefined,
  triggered_by: undefined as string | undefined,
  time_range: null as [Dayjs, Dayjs] | null,
})
const columns = [
  { title: '注册任务', dataIndex: 'task_code', key: 'task_code', width: 210, ellipsis: true, align: 'center' },
  { title: '状态', key: 'status', width: 100, align: 'center' },
  { title: '触发方式', key: 'triggered_by', width: 100, align: 'center' },
  { title: '开始时间', key: 'started_at', width: 180, align: 'center' },
  { title: '结束时间', key: 'finished_at', width: 180, align: 'center' },
  { title: '耗时', key: 'duration_ms', width: 100, align: 'center' },
  { title: '摘要', key: 'summary', width: 220, ellipsis: true },
  { title: '操作', key: 'action', width: 90, fixed: 'right', align: 'center' },
]

const getStatusColor = (status?: string) => ({ success: 'green', failure: 'red', running: 'processing', skipped: 'orange' }[status || ''] || 'default')
const getStatusText = (status?: string) => ({ success: '成功', failure: '失败', running: '执行中', skipped: '已跳过' }[status || ''] || '-')
const formatDuration = (duration?: number) => duration ? `${duration}ms` : '-'

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getCronLogList({
      page: pagination.current,
      page_size: pagination.pageSize,
      task_id: searchForm.task_id,
      task_code: searchForm.task_code || undefined,
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

const handleSearch = () => { pagination.current = 1; fetchData() }
const handleReset = () => {
  searchForm.task_id = undefined
  searchForm.task_code = ''
  searchForm.status = undefined
  searchForm.triggered_by = undefined
  searchForm.time_range = null
  handleSearch()
}
const handleTableChange = (pag: any) => { pagination.current = pag.current; pagination.pageSize = pag.pageSize; fetchData() }
const handleDetail = (record: CronLog) => { currentLog.value = record; detailVisible.value = true }
const openLogFromQuery = async () => {
  const logID = Number(route.query.log_id || 0)
  if (!logID) return
  const res = await getCronLogDetail(logID)
  currentLog.value = res.data
  detailVisible.value = true
}

onMounted(async () => {
  const taskID = Number(route.query.task_id || 0)
  if (taskID) searchForm.task_id = taskID
  await fetchData()
  await openLogFromQuery()
})
</script>

<style scoped>
.cron-log-page__content { display: flex; flex-direction: column; gap: 16px; }
.ellipsis-text { display: inline-block; max-width: 220px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; vertical-align: bottom; }
pre { margin: 0; white-space: pre-wrap; word-break: break-word; font-family: Consolas, monospace; }
pre.danger { color: #cf1322; }
</style>
