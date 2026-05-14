<template>
  <PageWrapper class="operation-log-page">
    <div class="operation-log-page__content">
      <div class="log-summary">
        <div class="summary-item summary-item--primary">
          <div class="summary-item__icon">
            <FileSearchOutlined />
          </div>
          <div>
            <div class="summary-item__label">日志总数</div>
            <div class="summary-item__value">{{ pagination.total }}</div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-item__icon summary-item__icon--warning">
            <WarningOutlined />
          </div>
          <div>
            <div class="summary-item__label">当前页异常</div>
            <div class="summary-item__value">{{ currentPageErrorCount }}</div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-item__icon summary-item__icon--danger">
            <ThunderboltOutlined />
          </div>
          <div>
            <div class="summary-item__label">慢请求</div>
            <div class="summary-item__value">{{ slowRequestCount }}</div>
          </div>
        </div>
      </div>
      <AdminSplitLayout class="operation-log-page__layout" :aside-width="252" :content-min-width="1024">
        <template #aside>
          <div class="left-tree">
            <div class="tree-header">
              <div>
                <div class="tree-title">路由分组</div>
                <div class="tree-subtitle">当前：{{ selectedGroupTitle }}</div>
              </div>
              <a-tag color="blue">{{ routeGroupCount }}</a-tag>
            </div>
            <a-spin :spinning="groupLoading">
              <a-tree
                v-if="groupTreeData.length"
                class="group-tree"
                :tree-data="groupTreeData"
                :selected-keys="[selectedGroup]"
                @select="handleGroupSelect"
              />
              <a-empty v-else description="暂无分组" :image="simpleImage" />
            </a-spin>
          </div>
        </template>
        <div class="right-content">
          <ProTable
            title="操作日志列表"
            :columns="columns"
            :data-source="tableData"
            :loading="loading"
            :pagination="pagination"
            :scroll="{ x: 1420, y: 520 }"
            @search="handleSearch"
            @reset="handleReset"
            @change="handleTableChange"
          >
            <template #search>
              <a-form-item label="用户名">
                <a-input v-model:value="searchForm.username" placeholder="请输入用户名" allowClear style="width: 180px" />
              </a-form-item>
              <a-form-item label="请求方法">
                <a-select v-model:value="searchForm.method" placeholder="请选择方法" allowClear style="width: 128px">
                  <a-select-option value="GET">GET</a-select-option>
                  <a-select-option value="POST">POST</a-select-option>
                  <a-select-option value="PUT">PUT</a-select-option>
                  <a-select-option value="DELETE">DELETE</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="路由描述">
                <a-input v-model:value="searchForm.summary" placeholder="请输入路由描述" allowClear style="width: 220px" />
              </a-form-item>
              <a-form-item label="业务状态码">
                <a-select v-model:value="searchForm.business_code" placeholder="请选择状态码" allowClear style="width: 148px">
                  <a-select-option :value="200">200</a-select-option>
                  <a-select-option :value="400">400</a-select-option>
                  <a-select-option :value="403">403</a-select-option>
                  <a-select-option :value="500">500</a-select-option>
                </a-select>
              </a-form-item>
            </template>
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'username'">
                <span class="user-name">{{ record.username || '未知用户' }}</span>
              </template>
              <template v-if="column.key === 'ip'">
                <span class="ip-text">{{ record.ip || '-' }}</span>
              </template>
              <template v-if="column.key === 'method'">
                <a-tag class="method-tag" :color="getMethodColor(record.method)">{{ record.method || '-' }}</a-tag>
              </template>
              <template v-if="column.key === 'group'">
                <span class="group-text">{{ record.group || '-' }}</span>
              </template>
              <template v-if="column.key === 'summary'">
                <span class="summary-text" :title="record.summary">{{ record.summary || '-' }}</span>
              </template>
              <template v-if="column.key === 'path'">
                <span class="path-text" :title="record.path">{{ record.path || '-' }}</span>
              </template>
              <template v-if="column.key === 'status'">
                <a-tag :color="getStatusColor(record.status)">{{ record.status || '-' }}</a-tag>
              </template>
              <template v-if="column.key === 'business_code'">
                <a-tag :color="getStatusColor(record.business_code)">{{ record.business_code || '-' }}</a-tag>
              </template>
              <template v-if="column.key === 'latency'">
                <span class="latency-cell">
                  <span :class="['latency-dot', getLatencyLevel(record.latency)]"></span>
                  <span :class="['latency-value', getLatencyLevel(record.latency)]">{{ formatLatency(record.latency) }}</span>
                </span>
              </template>
              <template v-if="column.key === 'created_at'">
                <span class="time-text">{{ formatTime(record.created_at) }}</span>
              </template>
              <template v-if="column.key === 'action'">
                <a-button class="detail-button" type="link" size="small" @click="handleViewDetail(record)">详情</a-button>
              </template>
            </template>
          </ProTable>
        </div>
      </AdminSplitLayout>
    <!-- 抽屉详情 -->
      <a-drawer v-model:open="detailVisible" title="操作日志详情" width="760" placement="right">
        <div class="detail-hero">
          <div class="detail-route">
            <a-tag class="method-tag" :color="getMethodColor(currentLog?.method || '')">{{ currentLog?.method || '-' }}</a-tag>
            <strong>{{ currentLogSummary }}</strong>
          </div>
          <div class="detail-path">{{ currentLog?.path || '-' }}</div>
        </div>
        <a-descriptions class="detail-descriptions" :column="2" bordered size="small">
          <a-descriptions-item label="用户名">{{ currentLog?.username || '-' }}</a-descriptions-item>
          <a-descriptions-item label="用户ID">{{ currentLog?.user_id }}</a-descriptions-item>
          <a-descriptions-item label="IP地址">{{ currentLog?.ip || '-' }}</a-descriptions-item>
          <a-descriptions-item label="请求方法">
            <a-tag :color="getMethodColor(currentLog?.method || '')">{{ currentLog?.method || '-' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="路由分组">{{ currentLog?.group || '-' }}</a-descriptions-item>
          <a-descriptions-item label="路由描述">{{ currentLog?.summary || '-' }}</a-descriptions-item>
          <a-descriptions-item label="请求路径" :span="2">{{ currentLog?.path || '-' }}</a-descriptions-item>
          <a-descriptions-item label="HTTP状态码">
            <a-tag :color="getStatusColor(currentLog?.status)">{{ currentLog?.status || '-' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="业务状态码">
            <a-tag :color="getStatusColor(currentLog?.business_code)">{{ currentLog?.business_code || '-' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="耗时">
            <span class="latency-cell">
              <span :class="['latency-dot', getLatencyLevel(currentLog?.latency)]"></span>
              <span :class="['latency-value', getLatencyLevel(currentLog?.latency)]">{{ formatLatency(currentLog?.latency) }}</span>
            </span>
          </a-descriptions-item>
          <a-descriptions-item label="请求时间">{{ formatTime(currentLog?.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="User-Agent" :span="2">{{ currentLog?.user_agent || '-' }}</a-descriptions-item>
        </a-descriptions>
        <div class="payload-grid">
          <div class="payload-card">
            <div class="payload-title">请求参数</div>
            <div class="json-viewer-wrap">
              <pre v-if="currentLog?.request" class="json-pre">{{ formatJsonPretty(currentLog.request) }}</pre>
              <a-empty v-else description="无请求参数" />
            </div>
          </div>
          <div class="payload-card">
            <div class="payload-title">响应结果</div>
            <div class="json-viewer-wrap">
              <pre v-if="currentLog?.response" class="json-pre">{{ formatJsonPretty(currentLog.response) }}</pre>
              <a-empty v-else description="无响应结果" />
            </div>
          </div>
        </div>
      </a-drawer>
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, ref, reactive, onMounted } from 'vue'
import { Empty } from 'ant-design-vue'
import { FileSearchOutlined, ThunderboltOutlined, WarningOutlined } from '@ant-design/icons-vue'
import { getOperationLogList, getRouteGroups } from '@/api/log'
import { formatTime } from '@/utils/format'
import AdminSplitLayout from '@/components/AdminSplitLayout.vue'
import ProTable from '@/components/ProTable.vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import type { OperationLog } from '@/types'

type OperationLogRecord = OperationLog & {
  group?: string
  summary?: string
  business_code?: number | string
}

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const loading = ref(false)
const groupLoading = ref(false)
const tableData = ref<OperationLogRecord[]>([])
const detailVisible = ref(false)
const currentLog = ref<OperationLogRecord | null>(null)
const searchForm = reactive({ 
  username: '', 
  method: undefined as string | undefined, 
  group: '',
  summary: '',
  business_code: undefined as number | undefined
})
const sortInfo = reactive({ field: '', order: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

// 分组树相关
const groupTreeData = ref<{ key: string; title: string }[]>([])
const selectedGroup = ref<string>('')
const columns = [
  { title: '用户名', dataIndex: 'username', key: 'username', width: 110, align: 'center' },
  { title: 'IP', dataIndex: 'ip', key: 'ip', width: 130, align: 'center' },
  { title: '方法', key: 'method', width: 90, align: 'center' },
  { title: '分组', dataIndex: 'group', key: 'group', align: 'center', width: 150 },
  { title: '描述', dataIndex: 'summary', key: 'summary', width: 150, align: 'center', ellipsis: true },
  { title: '请求路径', dataIndex: 'path', key: 'path', width: 220, ellipsis: true, align: 'center' },
  { title: 'HTTP状态', key: 'status', width: 100, align: 'center' },
  { title: '业务码', key: 'business_code', width: 90, align: 'center' },
  { title: '耗时', key: 'latency', dataIndex: 'latency', width: 100, sorter: true, align: 'center' },
  { title: '请求时间', key: 'created_at', width: 180, align: 'center' },
  { title: '操作', key: 'action', width: 90, fixed: 'right', align: 'center' }
]

const formatJsonPretty = (str?: string) => {
  if (!str) return ''
  try { return JSON.stringify(JSON.parse(str), null, 2) } catch { return str }
}

const getMethodColor = (m: string) => ({ GET: 'green', POST: 'blue', PUT: 'orange', DELETE: 'red' }[m] || 'default')
const isErrorStatus = (value?: number | string) => {
  if (value === undefined || value === null || value === '') return false
  const code = Number(value)
  return !Number.isNaN(code) && (code < 200 || code >= 400)
}
const getStatusColor = (value?: number | string) => {
  if (value === undefined || value === null || value === '') return 'default'
  const code = Number(value)
  if (Number.isNaN(code)) return 'default'
  if (code >= 200 && code < 400) return 'green'
  if (code === 401 || code === 403) return 'orange'
  if (code >= 500) return 'red'
  return 'gold'
}
const getLatencyLevel = (latency?: number) => {
  const value = Number(latency || 0)
  if (value >= 1000) return 'danger'
  if (value >= 300) return 'warning'
  return 'normal'
}
const formatLatency = (latency?: number) => {
  if (latency === undefined || latency === null) return '-'
  return `${latency}ms`
}
const routeGroupCount = computed(() => Math.max(groupTreeData.value.length - 1, 0))
const selectedGroupTitle = computed(() => {
  const matched = groupTreeData.value.find(item => item.key === selectedGroup.value)
  return matched?.title || '全部'
})
const currentPageErrorCount = computed(() => tableData.value.filter(item => isErrorStatus(item.status) || isErrorStatus(item.business_code)).length)
const slowRequestCount = computed(() => tableData.value.filter(item => Number(item.latency || 0) >= 1000).length)
const currentLogSummary = computed(() => currentLog.value?.summary || currentLog.value?.path || '接口调用')
// 获取分组数据
const fetchGroups = async () => {
  groupLoading.value = true
  try {
    const res = await getRouteGroups()
    groupTreeData.value = [
      { key: '', title: `全部 (${res.data.reduce((sum, item) => sum + item.count, 0)})` },
      ...res.data.map(item => ({ key: item.group, title: `${item.group} (${item.count})` }))
    ]
  } finally { groupLoading.value = false }
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getOperationLogList({
      page: pagination.current,
      page_size: pagination.pageSize,
      username: searchForm.username || undefined,
      method: searchForm.method || undefined,
      group: searchForm.group || undefined,
      summary: searchForm.summary || undefined,
      business_code: searchForm.business_code || undefined,
      sort_field: sortInfo.field || undefined,
      sort_order: sortInfo.order || undefined
    })
    tableData.value = res.data.list
    pagination.total = res.data.total
  } finally { loading.value = false }
}

const handleGroupSelect = (keys: string[]) => {
  selectedGroup.value = keys[0] ?? ''
  searchForm.group = selectedGroup.value
  pagination.current = 1
  fetchData()
}

const handleSearch = () => { pagination.current = 1; fetchData() }
const handleReset = () => {
  searchForm.username = ''
  searchForm.method = undefined
  searchForm.summary = ''
  searchForm.business_code = undefined
  sortInfo.field = ''
  sortInfo.order = ''
  // 不重置分组选择
  handleSearch()
}
const handleTableChange = (pag: any, _filters: any, sorter: any) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  sortInfo.field = sorter.field || ''
  sortInfo.order = sorter.order || ''
  fetchData()
}
const handleViewDetail = (record: OperationLogRecord) => { currentLog.value = record; detailVisible.value = true }
onMounted(() => {
  fetchGroups()
  fetchData()
})
</script>

<style scoped>
.operation-log-page__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.log-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(160px, 1fr));
  gap: 16px;
}
.summary-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
  box-shadow: 0 6px 18px rgba(15, 35, 66, 0.04);
}
.summary-item--primary {
  border-color: rgba(22, 119, 255, 0.18);
}
.summary-item__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  color: #1677ff;
  font-size: 18px;
  background: rgba(22, 119, 255, 0.10);
  border-radius: 10px;
}
.summary-item__icon--warning {
  color: #d46b08;
  background: rgba(250, 140, 22, 0.12);
}
.summary-item__icon--danger {
  color: #cf1322;
  background: rgba(255, 77, 79, 0.12);
}
.summary-item__label {
  margin-bottom: 4px;
  color: var(--app-text-muted);
  font-size: 12px;
}
.summary-item__value {
  color: var(--app-text-strong);
  font-size: 22px;
  font-weight: 650;
  line-height: 1;
}
.left-tree {
  width: 100%;
  background: var(--app-surface-color);
  border-radius: 8px;
  padding: 14px;
  border: 1px solid var(--app-border-color);
  box-shadow: 0 6px 18px rgba(15, 35, 66, 0.04);
}
.tree-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  font-weight: 500;
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--app-border-color);
  color: var(--app-text-strong);
}
.tree-title {
  font-size: 15px;
  line-height: 1.3;
}
.tree-subtitle {
  max-width: 160px;
  margin-top: 4px;
  color: var(--app-text-muted);
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.left-tree :deep(.ant-tree) {
  max-height: 560px;
  overflow-y: auto;
  background: transparent;
  color: var(--app-text-color);
}
.left-tree :deep(.ant-tree-node-content-wrapper) {
  border-radius: 6px;
}
.left-tree :deep(.ant-tree-node-selected) {
  font-weight: 500;
}
.right-content {
  min-width: 0;
}
.right-content :deep(.ant-card) {
  border-radius: 8px;
  box-shadow: 0 6px 18px rgba(15, 35, 66, 0.04);
}
.right-content :deep(.ant-table-thead > tr > th) {
  color: var(--app-text-strong);
  font-size: 12px;
  font-weight: 600;
}
.right-content :deep(.ant-table-tbody > tr > td) {
  font-size: 12px;
}
.right-content :deep(.ant-table-body) {
  max-height: 520px;
  overflow-y: auto !important;
}
.user-name {
  color: var(--app-text-strong);
  font-weight: 500;
}
.ip-text,
.group-text,
.summary-text,
.path-text,
.time-text {
  color: var(--app-text-color);
}
.summary-text,
.path-text {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  white-space: nowrap;
}
.method-tag {
  min-width: 56px;
  text-align: center;
}
.latency-cell {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.latency-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
}
.latency-dot.normal {
  background: #52c41a;
}
.latency-dot.warning {
  background: #fa8c16;
}
.latency-dot.danger {
  background: #ff4d4f;
}
.latency-value.normal {
  color: #389e0d;
}
.latency-value.warning {
  color: #d46b08;
}
.latency-value.danger {
  color: #cf1322;
  font-weight: 600;
}
.detail-button {
  padding: 0;
}
.detail-hero {
  margin-bottom: 16px;
  padding: 14px 16px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}
.detail-route {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--app-text-strong);
  font-size: 15px;
}
.detail-path {
  margin-top: 8px;
  color: var(--app-text-muted);
  font-size: 12px;
  word-break: break-all;
}
.detail-descriptions {
  margin-bottom: 16px;
}
.payload-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 14px;
}
.payload-card {
  padding: 12px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}
.payload-title {
  margin-bottom: 10px;
  color: var(--app-text-strong);
  font-size: 13px;
  font-weight: 600;
}
.json-viewer-wrap {
  max-height: 300px;
  overflow: auto;
  background: var(--app-surface-soft);
  border-radius: 6px;
  border: 1px solid var(--app-border-color);
}
.json-pre {
  margin: 0;
  padding: 12px;
  background: var(--app-surface-soft);
  font-size: 12px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Courier New', monospace;
  white-space: pre-wrap;
  word-break: break-all;
  color: var(--app-text-color);
  line-height: 1.5;
}
@media (max-width: 1200px) {
  .left-tree {
    min-width: 0;
  }
}
@media (max-width: 768px) {
  .log-summary {
    grid-template-columns: 1fr;
  }
}
</style>
