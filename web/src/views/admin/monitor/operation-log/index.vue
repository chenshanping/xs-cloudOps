<template>
  <div class="operation-log-page">
    <div class="page-layout">
      <!-- 左侧分组树 -->
      <div class="left-tree">
        <div class="tree-header">路由分组</div>
        <a-spin :spinning="groupLoading">
          <a-tree
            v-if="groupTreeData.length"
            :tree-data="groupTreeData"
            :selected-keys="[selectedGroup]"
            @select="handleGroupSelect"
          />
          <a-empty v-else description="暂无分组" :image="simpleImage" />
        </a-spin>
      </div>
      <!-- 右侧表格 -->
      <div class="right-content">
        <ProTable
          :columns="columns"
          :data-source="tableData"
          :loading="loading"
          :pagination="pagination"
          :scroll="{ x: 1400,y: 400 }"
          @search="handleSearch"
          @reset="handleReset"
          @change="handleTableChange"
        >
          <template #search>
            <a-form-item label="用户名"><a-input v-model:value="searchForm.username" placeholder="请输入用户名" allowClear style="width: 150px" /></a-form-item>
            <a-form-item label="请求方法">
              <a-select v-model:value="searchForm.method" placeholder="请选择" allowClear style="width: 100px">
                <a-select-option value="GET">GET</a-select-option>
                <a-select-option value="POST">POST</a-select-option>
                <a-select-option value="PUT">PUT</a-select-option>
                <a-select-option value="DELETE">DELETE</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="路由描述"><a-input v-model:value="searchForm.summary" placeholder="请输入路由描述" allowClear style="width: 150px" /></a-form-item>
            <a-form-item label="业务状态码">
              <!-- <a-input v-model:value="searchForm.business_code" placeholder="请输入业务状态码" allowClear /> -->
            
              <a-select v-model:value="searchForm.business_code" placeholder="请选择" allowClear style="width: 120px">
                <a-select-option :value="200">200</a-select-option>
                 <a-select-option :value="400">400</a-select-option>
                <a-select-option :value="403">403</a-select-option>
                <a-select-option :value="500">500</a-select-option>
              </a-select>
            </a-form-item>
          </template>
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'method'">
              <a-tag :color="getMethodColor(record.method)">{{ record.method }}</a-tag>
            </template>
            <template v-if="column.key === 'status'">
              <a-tag :color="record.status === 200 ? 'green' : 'red'">{{ record.status }}</a-tag>
            </template>
            <template v-if="column.key === 'business_code'">
              <a-tag :color="record.business_code === 200 ? 'green' : 'red'">{{ record.business_code }}</a-tag>
            </template>
            <template v-if="column.key === 'latency'">
              <span :style="{ color: record.latency > 200 ? 'red' : 'inherit' }">{{ record.latency }}ms</span>
            </template>
            <template v-if="column.key === 'created_at'">
              {{ formatTime(record.created_at) }}
            </template>
            <template v-if="column.key === 'action'">
              <a-button type="link" size="small" @click="handleViewDetail(record)">详情</a-button>
            </template>
          </template>
        </ProTable>
      </div>
    </div>
    <!-- 抽屉详情 -->
    <a-drawer v-model:open="detailVisible" title="日志详情" width="640" placement="right">
      <a-descriptions :column="2" bordered size="small">
        <a-descriptions-item label="用户名">{{ currentLog?.username }}</a-descriptions-item>
        <a-descriptions-item label="用户ID">{{ currentLog?.user_id }}</a-descriptions-item>
        <a-descriptions-item label="IP地址">{{ currentLog?.ip }}</a-descriptions-item>
        <a-descriptions-item label="请求方法">
          <a-tag :color="getMethodColor(currentLog?.method || '')">{{ currentLog?.method }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="路由分组">{{ currentLog?.group }}</a-descriptions-item>
        <a-descriptions-item label="路由描述">{{ currentLog?.summary }}</a-descriptions-item>
        <a-descriptions-item label="请求路径" :span="2">{{ currentLog?.path }}</a-descriptions-item>
        <a-descriptions-item label="HTTP状态码">
          <a-tag :color="currentLog?.status === 200 ? 'green' : 'red'">{{ currentLog?.status }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="业务状态码">
          <a-tag :color="currentLog?.business_code === 200 ? 'green' : 'red'">{{ currentLog?.business_code }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="耗时">
          <span :style="{ color: (currentLog?.latency || 0) > 200 ? 'red' : 'inherit' }">{{ currentLog?.latency }}ms</span>
        </a-descriptions-item>
        <a-descriptions-item label="请求时间">{{ formatTime(currentLog?.created_at) }}</a-descriptions-item>
        <a-descriptions-item label="User-Agent" :span="2">{{ currentLog?.user_agent }}</a-descriptions-item>
      </a-descriptions>
      <a-divider>请求参数</a-divider>
      <div class="json-viewer-wrap">
        <pre v-if="currentLog?.request" class="json-pre">{{ formatJsonPretty(currentLog.request) }}</pre>
        <a-empty v-else description="无请求参数" />
      </div>
      <a-divider>响应结果</a-divider>
      <div class="json-viewer-wrap">
        <pre v-if="currentLog?.response" class="json-pre">{{ formatJsonPretty(currentLog.response) }}</pre>
        <a-empty v-else description="无响应结果" />
      </div>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Empty } from 'ant-design-vue'
import { getOperationLogList, getRouteGroups } from '@/api/log'
import { formatTime } from '@/utils/format'
import ProTable from '@/components/ProTable.vue'
import type { OperationLog } from '@/types'

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const loading = ref(false)
const groupLoading = ref(false)
const tableData = ref<OperationLog[]>([])
const detailVisible = ref(false)
const currentLog = ref<OperationLog | null>(null)
const searchForm = reactive({ 
  username: '', 
  method: undefined as string | undefined, 
  group: '',
  summary: '',
  business_code: ''
})
const sortInfo = reactive({ field: '', order: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

// 分组树相关
const groupTreeData = ref<{ key: string; title: string }[]>([])
const selectedGroup = ref<string>('')
const columns = [
  { title: '用户名', dataIndex: 'username', key: 'username', width: 100,align: 'center' },
  { title: 'IP', dataIndex: 'ip', key: 'ip', width: 120 ,align: 'center'},
  { title: '方法', key: 'method', width: 80,align: 'center' },
  { title: '分组', dataIndex: 'group', key: 'group', align: 'center',width: 150 },
  { title: '描述', dataIndex: 'summary', key: 'summary', width: 150, align: 'center', ellipsis: true },
  { title: '请求路径', dataIndex: 'path', key: 'path', width: 200, ellipsis: true, align: 'center' },
  { title: 'HTTP状态', key: 'status', width: 100, align: 'center' },
  { title: '业务码', key: 'business_code', width: 80, align: 'center' },
  { title: '耗时', key: 'latency', dataIndex: 'latency', width: 80, sorter: true, align: 'center' },
  { title: '请求时间', key: 'created_at', width: 180, align: 'center' },
  { title: '操作', key: 'action', width: 80, fixed: 'right', align: 'center' }
]

const formatJsonPretty = (str?: string) => {
  if (!str) return ''
  try { return JSON.stringify(JSON.parse(str), null, 2) } catch { return str }
}

const getMethodColor = (m: string) => ({ GET: 'green', POST: 'blue', PUT: 'orange', DELETE: 'red' }[m] || 'default')
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
  searchForm.business_code = ''
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
const handleViewDetail = (record: OperationLog) => { currentLog.value = record; detailVisible.value = true }
onMounted(() => {
  fetchGroups()
  fetchData()
})
</script>

<style scoped>
.page-layout {
  display: flex;
  gap: 16px;
}
.left-tree {
  width: 220px;
  flex-shrink: 0;
  background: var(--app-surface-color);
  border-radius: 4px;
  padding: 12px;
  border: 1px solid var(--app-border-color);
}
.tree-header {
  font-weight: 500;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--app-border-color);
  color: var(--app-text-strong);
}
.left-tree :deep(.ant-tree) {
  max-height: 500px;
  overflow-y: auto;
  background: transparent;
  color: var(--app-text-color);
}
.right-content {
  flex: 1;
  min-width: 0;
}
.right-content :deep(.ant-table-body) {
  max-height: 500px;
  overflow-y: auto !important;
}
.json-viewer-wrap {
  max-height: 300px;
  overflow: auto;
  background: var(--app-surface-soft);
  border-radius: 4px;
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
</style>
