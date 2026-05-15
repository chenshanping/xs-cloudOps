<template>
  <PageWrapper class="cron-task-page">
    <div class="cron-task-page__content">
      <div class="cron-summary">
        <div class="summary-item"><ScheduleOutlined /><span>任务总数</span><strong>{{ pagination.total }}</strong></div>
        <div class="summary-item success"><CheckCircleOutlined /><span>当前页启用</span><strong>{{ enabledCount }}</strong></div>
        <div class="summary-item warning"><SyncOutlined /><span>执行中</span><strong>{{ runningCount }}</strong></div>
        <div class="summary-item danger"><CloseCircleOutlined /><span>当前页失败</span><strong>{{ failedCount }}</strong></div>
      </div>

    <ProTable
      title="定时任务列表"
      :columns="columns"
      :data-source="tableData"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      :scroll="{ x: 1420 }"
      @search="handleSearch"
      @reset="handleReset"
      @change="handleTableChange"
    >
      <template #search>
        <a-form-item label="实例编码"><a-input v-model:value="searchForm.code" placeholder="请输入实例编码" allow-clear style="width: 180px" /></a-form-item>
        <a-form-item label="任务名称"><a-input v-model:value="searchForm.name" placeholder="请输入任务名称" allow-clear style="width: 180px" /></a-form-item>
        <a-form-item label="注册任务">
          <a-select v-model:value="searchForm.task_code" placeholder="请选择" allow-clear style="width: 220px">
            <a-select-option v-for="item in registry" :key="item.code" :value="item.code">{{ item.name }}（{{ item.code }}）</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchForm.status" placeholder="请选择" allow-clear style="width: 120px">
            <a-select-option :value="1">启用</a-select-option>
            <a-select-option :value="0">停用</a-select-option>
          </a-select>
        </a-form-item>
      </template>
      <template #toolbar>
        <a-button type="primary" @click="handleAdd" v-permission="'monitor:cron:create'"><PlusOutlined /> 新增任务</a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'name'">
          <a-button v-if="hasViewPermission" type="link" class="task-name-button" @click="handleViewTask(record)" v-permission="'monitor:cron:view'">
            <div class="task-name">
              <strong>{{ record.name }}</strong>
              <span>{{ record.code }}</span>
            </div>
          </a-button>
          <div v-else class="task-name"><strong>{{ record.name }}</strong><span>{{ record.code }}</span></div>
        </template>
        <template v-else-if="column.key === 'task_code'"><a-tag color="blue">{{ record.task_code }}</a-tag></template>
        <template v-else-if="column.key === 'cron_expr'"><code>{{ record.cron_expr }}</code></template>
        <template v-else-if="column.key === 'status'"><a-tag :color="record.status === 1 ? 'green' : 'default'">{{ record.status === 1 ? '启用' : '停用' }}</a-tag></template>
        <template v-else-if="column.key === 'last_status'"><a-tag :color="getRunStatusColor(record.last_status)">{{ getRunStatusText(record.last_status) }}</a-tag></template>
        <template v-else-if="column.key === 'last_duration_ms'">{{ formatDuration(record.last_duration_ms) }}</template>
        <template v-else-if="column.key === 'last_run_at'">{{ formatTime(record.last_run_at) }}</template>
        <template v-else-if="column.key === 'next_run_at'">{{ record.status === 1 ? formatTime(record.next_run_at) : '-' }}</template>
        <template v-else-if="column.key === 'action'">
          <a-space :size="0">
            <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'monitor:cron:update'">编辑</a-button>
            <a-button v-if="record.status !== 1" type="link" size="small" @click="handleEnable(record)" v-permission="'monitor:cron:enable'">启用</a-button>
            <a-button v-else type="link" size="small" danger @click="handleDisable(record)" v-permission="'monitor:cron:disable'">停用</a-button>
            <a-button type="link" size="small" @click="handleRunNow(record)" v-permission="'monitor:cron:runNow'">立即执行</a-button>
            <a-button type="link" size="small" @click="handleViewLogs(record)" v-permission="'monitor:cron:logs:view'">日志</a-button>
            <a-button type="link" size="small" danger @click="handleDelete(record)" v-permission="'monitor:cron:delete'">删除</a-button>
          </a-space>
        </template>
      </template>
    </ProTable>

      <CronTaskFormDrawer v-model:open="drawerVisible" :title="isEdit ? '编辑定时任务' : '新增定时任务'" :is-edit="isEdit" :record="currentRecord" :registry="registry" :submitting="submitting" @submit="handleSubmit" />
      <CronTaskDetailDrawer
        v-model:open="detailVisible"
        :task="detailTask"
        :registry="registry"
        :initial-tab="detailInitialTab"
        :initial-log-id="detailLogId"
        @run-now="handleRunNow"
      />
      <CronRunNowConfirm ref="runNowConfirmRef" @success="handleRunNowSuccess" />
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { CheckCircleOutlined, CloseCircleOutlined, PlusOutlined, ScheduleOutlined, SyncOutlined } from '@ant-design/icons-vue'
import ProTable from '@/components/ProTable.vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import { formatTime } from '@/utils/format'
import type { CronTask } from '@/api/cron'
import { useUserStore } from '@/store/user'
import CronRunNowConfirm from './components/CronRunNowConfirm.vue'
import CronTaskDetailDrawer from './components/CronTaskDetailDrawer.vue'
import CronTaskFormDrawer from './components/CronTaskFormDrawer.vue'
import { useCronTaskPage } from './useCronTaskPage'

const userStore = useUserStore()
const runNowConfirmRef = ref<InstanceType<typeof CronRunNowConfirm>>()
const detailVisible = ref(false)
const detailTask = ref<CronTask | null>(null)
const detailInitialTab = ref<'detail' | 'logs'>('detail')
const detailLogId = ref<number>()
const {
  loading, submitting, tableData, registry, drawerVisible, isEdit, currentRecord,
  pagination, searchForm, enabledCount, failedCount, runningCount, fetchRegistry, fetchData,
  handleSearch, handleReset, handleTableChange, handleAdd, handleEdit, handleSubmit,
  handleDelete: baseHandleDelete, handleEnable: baseHandleEnable, handleDisable: baseHandleDisable,
} = useCronTaskPage()

const columns = [
  { title: '任务名称', key: 'name', width: 240, fixed: 'left' },
  { title: '注册任务', key: 'task_code', width: 210, align: 'center' },
  { title: 'Cron 表达式', key: 'cron_expr', width: 140, align: 'center' },
  { title: '状态', key: 'status', width: 90, align: 'center' },
  { title: '上次状态', key: 'last_status', width: 110, align: 'center' },
  { title: '上次耗时', key: 'last_duration_ms', width: 110, align: 'center' },
  { title: '上次执行时间', key: 'last_run_at', width: 180, align: 'center' },
  { title: '下次执行时间', key: 'next_run_at', width: 180, align: 'center' },
  { title: '操作', key: 'action', width: 260, fixed: 'right', align: 'center' },
]

const getRunStatusColor = (status?: string) => ({ success: 'green', failure: 'red', running: 'processing', skipped: 'orange' }[status || ''] || 'default')
const getRunStatusText = (status?: string) => ({ success: '成功', failure: '失败', running: '执行中', skipped: '已跳过' }[status || ''] || '-')
const formatDuration = (duration?: number) => duration ? `${duration}ms` : '-'
const hasViewPermission = computed(() => userStore.hasPermission('monitor:cron:view'))
const syncDetailTask = (taskId: number) => {
  const matched = tableData.value.find((item) => item.id === taskId) || null
  if (detailVisible.value && detailTask.value?.id === taskId) {
    detailTask.value = matched
    if (!matched) {
      detailVisible.value = false
    }
  }
  return matched
}
const openTaskDrawer = (task: CronTask, tab: 'detail' | 'logs' = 'detail', logId?: number) => {
  detailTask.value = task
  detailInitialTab.value = tab
  detailLogId.value = logId
  detailVisible.value = true
}
const handleViewTask = (record: CronTask) => openTaskDrawer(record, 'detail')
const handleRunNow = (record: CronTask) => runNowConfirmRef.value?.confirmRun(record)
const handleRunNowSuccess = async (logId: number, task: CronTask) => {
  await fetchData()
  openTaskDrawer(syncDetailTask(task.id) || task, 'logs', logId)
}
const handleViewLogs = (record: CronTask) => openTaskDrawer(record, 'logs')
const handleEnable = async (record: CronTask) => {
  await baseHandleEnable(record)
  syncDetailTask(record.id)
}
const handleDisable = async (record: CronTask) => {
  await baseHandleDisable(record)
  syncDetailTask(record.id)
}
const handleDelete = (record: CronTask) => {
  baseHandleDelete(record, () => {
    if (detailTask.value?.id === record.id) {
      detailTask.value = null
      detailVisible.value = false
    }
  })
}

onMounted(async () => {
  await fetchRegistry()
  await fetchData()
})
</script>

<style scoped>
.cron-task-page__content { display: flex; flex-direction: column; gap: 16px; }
.cron-summary { display: grid; grid-template-columns: repeat(4, minmax(140px, 1fr)); gap: 12px; }
.summary-item { display: grid; grid-template-columns: 34px 1fr auto; align-items: center; gap: 10px; padding: 14px 16px; color: var(--app-text-strong); background: var(--app-surface-color); border: 1px solid var(--app-border-color); border-radius: 8px; }
.summary-item :deep(.anticon) { display: flex; align-items: center; justify-content: center; width: 34px; height: 34px; color: #1677ff; background: rgba(22, 119, 255, 0.10); border-radius: 8px; }
.summary-item.success :deep(.anticon) { color: #389e0d; background: rgba(82, 196, 26, 0.12); }
.summary-item.warning :deep(.anticon) { color: #d46b08; background: rgba(250, 140, 22, 0.12); }
.summary-item.danger :deep(.anticon) { color: #cf1322; background: rgba(255, 77, 79, 0.12); }
.summary-item span { color: var(--app-text-muted); font-size: 13px; }
.summary-item strong { font-size: 22px; font-weight: 650; }
.task-name-button {
  height: auto;
  padding: 0;
  text-align: left;
}
.task-name { display: flex; flex-direction: column; gap: 3px; }
.task-name strong { color: var(--app-text-strong); }
.task-name span { color: var(--app-text-muted); font-size: 12px; }
code { padding: 2px 6px; color: #0958d9; background: rgba(22, 119, 255, 0.08); border-radius: 4px; }
@media (max-width: 1200px) { .cron-summary { grid-template-columns: repeat(2, minmax(140px, 1fr)); } }
</style>
