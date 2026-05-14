<template>
  <PageWrapper class="login-log-page">
    <div class="login-log-page__content">
      <div class="log-summary">
        <div class="summary-item summary-item--primary">
          <div class="summary-item__icon">
            <LoginOutlined />
          </div>
          <div>
            <div class="summary-item__label">日志总数</div>
            <div class="summary-item__value">{{ pagination.total }}</div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-item__icon summary-item__icon--success">
            <CheckCircleOutlined />
          </div>
          <div>
            <div class="summary-item__label">当前页成功</div>
            <div class="summary-item__value">{{ currentPageSuccessCount }}</div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-item__icon summary-item__icon--danger">
            <CloseCircleOutlined />
          </div>
          <div>
            <div class="summary-item__label">当前页失败</div>
            <div class="summary-item__value">{{ currentPageFailureCount }}</div>
          </div>
        </div>
      </div>
      <ProTable
        title="登录日志列表"
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        :pagination="pagination"
        :scroll="{ x: 1180, y: 520 }"
        @search="handleSearch"
        @reset="handleReset"
        @change="handleTableChange"
      >
        <template #search>
          <a-form-item label="用户名">
            <a-input v-model:value="searchForm.username" placeholder="请输入用户名" allowClear style="width: 180px" />
          </a-form-item>
          <a-form-item label="状态">
            <a-select v-model:value="searchForm.status" placeholder="请选择状态" allowClear style="width: 128px">
              <a-select-option :value="1">成功</a-select-option>
              <a-select-option :value="0">失败</a-select-option>
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
          <template v-if="column.key === 'location'">
            <span class="location-text">{{ record.location || '-' }}</span>
          </template>
          <template v-if="column.key === 'browser'">
            <a-tag class="client-tag" color="blue">{{ record.browser || '未知浏览器' }}</a-tag>
          </template>
          <template v-if="column.key === 'os'">
            <a-tag class="client-tag" color="geekblue">{{ record.os || '未知系统' }}</a-tag>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag class="status-tag" :color="record.status === 1 ? 'green' : 'red'">{{ record.status === 1 ? '成功' : '失败' }}</a-tag>
          </template>
          <template v-if="column.key === 'msg'">
            <span :class="['message-text', record.status === 1 ? 'success' : 'danger']" :title="record.msg">{{ record.msg || '-' }}</span>
          </template>
          <template v-if="column.key === 'created_at'">
            <span class="time-text">{{ formatTime(record.created_at) }}</span>
          </template>
        </template>
      </ProTable>
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, ref, reactive, onMounted } from 'vue'
import { CheckCircleOutlined, CloseCircleOutlined, LoginOutlined } from '@ant-design/icons-vue'
import { getLoginLogList } from '@/api/log'
import { formatTime } from '@/utils/format'
import ProTable from '@/components/ProTable.vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import type { LoginLog } from '@/types'

const loading = ref(false)
const tableData = ref<LoginLog[]>([])
const searchForm = reactive({ username: '', status: undefined as number | undefined })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const columns = [
  { title: '用户名', dataIndex: 'username', key: 'username', width: 130, align: 'center' },
  { title: 'IP地址', dataIndex: 'ip', key: 'ip', width: 150, align: 'center' },
  { title: '登录地点', dataIndex: 'location', key: 'location', width: 160, align: 'center', ellipsis: true },
  { title: '浏览器', dataIndex: 'browser', key: 'browser', width: 140, align: 'center', ellipsis: true },
  { title: '操作系统', dataIndex: 'os', key: 'os', width: 140, align: 'center', ellipsis: true },
  { title: '状态', key: 'status', width: 90, align: 'center' },
  { title: '消息', dataIndex: 'msg', key: 'msg', width: 220, ellipsis: true },
  { title: '登录时间', key: 'created_at', width: 180, align: 'center' }
]
const currentPageSuccessCount = computed(() => tableData.value.filter(item => item.status === 1).length)
const currentPageFailureCount = computed(() => tableData.value.filter(item => item.status === 0).length)
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
.login-log-page__content {
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
  padding: 12px 14px;
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
.summary-item__icon--success {
  color: #389e0d;
  background: rgba(82, 196, 26, 0.12);
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
.login-log-page :deep(.ant-card) {
  border-radius: 8px;
  box-shadow: 0 6px 18px rgba(15, 35, 66, 0.04);
}
.login-log-page :deep(.ant-table-thead > tr > th) {
  color: var(--app-text-strong);
  font-size: 12px;
  font-weight: 600;
}
.login-log-page :deep(.ant-table-tbody > tr > td) {
  font-size: 12px;
}
.user-name {
  color: var(--app-text-strong);
  font-weight: 500;
}
.ip-text {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Courier New', monospace;
}
.location-text,
.time-text {
  color: var(--app-text-color);
}
.client-tag {
  max-width: 116px;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
}
.status-tag {
  min-width: 52px;
  text-align: center;
}
.message-text {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  white-space: nowrap;
}
.message-text.success {
  color: #389e0d;
}
.message-text.danger {
  color: #cf1322;
}
@media (max-width: 768px) {
  .log-summary {
    grid-template-columns: 1fr;
  }
}
</style>
