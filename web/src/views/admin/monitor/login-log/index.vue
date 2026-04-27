<template>
  <div class="login-log-page">
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
        <a-form-item label="用户名"><a-input v-model:value="searchForm.username" placeholder="请输入用户名" allowClear /></a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchForm.status" placeholder="请选择" allowClear style="width: 120px">
            <a-select-option :value="1">成功</a-select-option>
            <a-select-option :value="0">失败</a-select-option>
          </a-select>
        </a-form-item>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">{{ record.status === 1 ? '成功' : '失败' }}</a-tag>
        </template>
        <template v-if="column.key === 'created_at'">{{ formatTime(record.created_at) }}</template>
      </template>
    </ProTable>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getLoginLogList } from '@/api/log'
import { formatTime } from '@/utils/format'
import ProTable from '@/components/ProTable.vue'
import type { LoginLog } from '@/types'

const loading = ref(false)
const tableData = ref<LoginLog[]>([])
const searchForm = reactive({ username: '', status: undefined as number | undefined })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '用户名', dataIndex: 'username', key: 'username', width: 120 },
  { title: 'IP地址', dataIndex: 'ip', key: 'ip', width: 140 },
  { title: '登录地点', dataIndex: 'location', key: 'location', width: 150 },
  { title: '浏览器', dataIndex: 'browser', key: 'browser', width: 120 },
  { title: '操作系统', dataIndex: 'os', key: 'os', width: 120 },
  { title: '状态', key: 'status', width: 80 },
  { title: '消息', dataIndex: 'msg', key: 'msg', ellipsis: true },
  { title: '登录时间', key: 'created_at', width: 180 }
]
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getLoginLogList({ page: pagination.current, page_size: pagination.pageSize, username: searchForm.username, status: searchForm.status })
    tableData.value = res.data.list
    pagination.total = res.data.total
  } finally { loading.value = false }
}
const handleSearch = () => { pagination.current = 1; fetchData() }
const handleReset = () => { searchForm.username = ''; searchForm.status = undefined; handleSearch() }
const handleTableChange = (pag: any) => { pagination.current = pag.current; pagination.pageSize = pag.pageSize; fetchData() }
onMounted(() => fetchData())
</script>

<style scoped>
</style>
