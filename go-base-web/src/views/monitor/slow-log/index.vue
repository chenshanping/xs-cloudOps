<template>
  <div class="slow-log-page">
    <ProTable
      :columns="columns"
      :data-source="tableData"
      :loading="loading"
      :pagination="pagination"
      @search="handleSearch"
      @reset="handleReset"
      @change="handleTableChange"
    >
      <template #search>
        <a-form-item label="SQL关键字"><a-input v-model:value="searchForm.sql" placeholder="请输入SQL关键字" allowClear /></a-form-item>
        <a-form-item label="最小耗时(ms)"><a-input-number v-model:value="searchForm.min_latency" :min="0" placeholder="请输入" /></a-form-item>
      </template>
      <template #bodyCell="{ column, record }">
        <!-- <template v-if="column.key === 'sql'">
          <a-tooltip :title="record.sql">
            <span class="sql-text">{{ record.sql }}</span>
          </a-tooltip>
        </template> -->
        <template v-if="column.key === 'latency'">
          <a-tag :color="getLatencyColor(record.latency)">{{ record.latency.toFixed(2) }}ms</a-tag>
        </template>
        <template v-if="column.key === 'created_at'">{{ formatTime(record.created_at) }}</template>
        <template v-if="column.key === 'action'">
          <a-button type="link" size="small" @click="handleViewDetail(record)">详情</a-button>
        </template>
      </template>
    </ProTable>
    <a-modal v-model:open="detailVisible" title="慢查询详情" :footer="null" width="800px">
      <a-descriptions :column="1" bordered size="small">
        <a-descriptions-item label="耗时"><a-tag :color="getLatencyColor(currentLog?.latency || 0)">{{ currentLog?.latency?.toFixed(2) }}ms</a-tag></a-descriptions-item>
        <a-descriptions-item label="影响行数">{{ currentLog?.rows }}</a-descriptions-item>
        <a-descriptions-item label="调用来源"><code>{{ currentLog?.source || '-' }}</code></a-descriptions-item>
        <a-descriptions-item label="执行时间">{{ formatTime(currentLog?.created_at) }}</a-descriptions-item>
        <a-descriptions-item label="SQL语句">
          <pre class="sql-detail"><code v-html="highlightedSql"></code></pre>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { getSlowLogList } from '@/api/log'
import { formatTime } from '@/utils/format'
import ProTable from '@/components/ProTable.vue'
import type { SlowLog } from '@/types'
import hljs from 'highlight.js/lib/core'
import sql from 'highlight.js/lib/languages/sql'
import 'highlight.js/styles/github.css'

hljs.registerLanguage('sql', sql)

const loading = ref(false)
const tableData = ref<SlowLog[]>([])
const detailVisible = ref(false)
const currentLog = ref<SlowLog | null>(null)
const searchForm = reactive({ sql: '', min_latency: undefined as number | undefined })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 ,align: 'center'},
  { title: 'SQL语句', dataIndex: 'sql',key: 'sql', ellipsis: true,align: 'center' },
  { title: '影响行数', dataIndex: 'rows', key: 'rows', width: 100 ,align: 'center'},
  { title: '耗时', key: 'latency', width: 120 ,align: 'center'},
  { title: '调用来源', dataIndex: 'source', key: 'source', width: 200, ellipsis: true,align: 'center' },
  { title: '执行时间', key: 'created_at', width: 180,align: 'center' },
  { title: '操作', key: 'action', width: 80 ,align: 'center'}
]

const getLatencyColor = (latency: number) => {
  if (latency >= 1000) return 'red'
  if (latency >= 500) return 'orange'
  return 'gold'
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getSlowLogList({ page: pagination.current, page_size: pagination.pageSize, sql: searchForm.sql, min_latency: searchForm.min_latency })
    tableData.value = res.data.list
    pagination.total = res.data.total
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; fetchData() }
const handleReset = () => { searchForm.sql = ''; searchForm.min_latency = undefined; handleSearch() }
const handleTableChange = (pag: any) => { pagination.current = pag.current; pagination.pageSize = pag.pageSize; fetchData() }
const handleViewDetail = (record: SlowLog) => { currentLog.value = record; detailVisible.value = true }

const highlightedSql = computed(() => {
  if (!currentLog.value?.sql) return ''
  return hljs.highlight(currentLog.value.sql, { language: 'sql' }).value
})

onMounted(() => fetchData())
</script>

<style scoped>
.sql-text {
  display: inline-block;
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.sql-detail {
  max-height: 300px;
  overflow: auto;
  margin: 0;
  padding: 12px;
  background: #f5f5f5;
  border-radius: 4px;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
