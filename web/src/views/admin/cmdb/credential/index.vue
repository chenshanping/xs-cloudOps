<template>
  <PageWrapper class="cmdb-credential-page">
    <div class="cmdb-credential-page__content">
      <div class="summary-grid">
        <div class="summary-item">
          <KeyOutlined />
          <span>凭据总数</span>
          <strong>{{ pagination.total }}</strong>
        </div>
        <div class="summary-item success">
          <LockOutlined />
          <span>密码认证</span>
          <strong>{{ passwordCount }}</strong>
        </div>
        <div class="summary-item warning">
          <SafetyCertificateOutlined />
          <span>私钥认证</span>
          <strong>{{ privateKeyCount }}</strong>
        </div>
      </div>

      <ProTable
        title="SSH 凭据"
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        :scroll="{ x: 1100 }"
        @search="handleSearch"
        @reset="handleReset"
        @change="handleTableChange"
      >
        <template #search>
          <a-form-item label="凭据名称">
            <a-input v-model:value="searchForm.name" placeholder="请输入凭据名称" allow-clear />
          </a-form-item>
          <a-form-item label="认证方式">
            <a-select v-model:value="searchForm.auth_type" placeholder="请选择认证方式" allow-clear style="width: 160px">
              <a-select-option v-for="item in cmdbAuthTypeOptions" :key="item.value" :value="item.value">
                {{ item.label }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </template>

        <template #toolbar>
          <a-button type="primary" v-permission="'cmdb:credential:create'" @click="handleAdd">
            <PlusOutlined /> 新增凭据
          </a-button>
        </template>

        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a-button
              v-if="canViewCredential"
              type="link"
              class="name-button"
              v-permission="'cmdb:credential:view'"
              @click="handleView(record)"
            >
              <div class="name-cell">
                <strong>{{ record.name }}</strong>
                <span>{{ record.username }}</span>
              </div>
            </a-button>
            <div v-else class="name-cell">
              <strong>{{ record.name }}</strong>
              <span>{{ record.username }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'auth_type'">
            <a-tag :color="record.auth_type === 'private_key' ? 'purple' : 'blue'">
              {{ getCmdbAuthTypeLabel(record.auth_type) }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'bind_count'">
            <a-badge :count="record.bind_count" :number-style="{ backgroundColor: '#1677ff' }" />
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space :size="0">
              <a-button type="link" size="small" v-permission="'cmdb:credential:view'" @click="handleView(record)">查看</a-button>
              <a-button type="link" size="small" v-permission="'cmdb:credential:update'" @click="handleEdit(record)">编辑</a-button>
              <a-button type="link" size="small" danger v-permission="'cmdb:credential:delete'" @click="handleDelete(record)">删除</a-button>
            </a-space>
          </template>
        </template>
      </ProTable>

      <CredentialFormDrawer
        v-model:open="drawerVisible"
        :title="drawerTitle"
        :mode="drawerMode"
        :initial-value="drawerInitialValue"
        :detail="detailData"
        @submit="handleSubmit"
      />
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import {
  ExclamationCircleOutlined,
  KeyOutlined,
  LockOutlined,
  PlusOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import ProTable from '@/components/ProTable.vue'
import {
  createCmdbCredential,
  deleteCmdbCredential,
  getCmdbCredential,
  getCmdbCredentials,
  updateCmdbCredential,
  type CmdbCredentialPayload,
  type CmdbCredentialSummary,
} from '@/api/cmdb'
import { usePermission, useTableColumns } from '@/utils/permission'
import { cmdbAuthTypeOptions, getCmdbAuthTypeLabel } from '../shared'
import CredentialFormDrawer from './components/CredentialFormDrawer.vue'

type DrawerMode = 'create' | 'edit' | 'view'

interface CredentialDrawerValue {
  name: string
  auth_type: 'password' | 'private_key'
  username: string
  remark: string
}

const { hasPermission } = usePermission()
const loading = ref(false)
const tableData = ref<CmdbCredentialSummary[]>([])
const drawerVisible = ref(false)
const drawerTitle = ref('新增凭据')
const drawerMode = ref<DrawerMode>('create')
const editingId = ref<number>()
const drawerInitialValue = ref<Partial<CredentialDrawerValue>>({})
const detailData = ref<CmdbCredentialSummary | null>(null)

const searchForm = reactive({
  name: '',
  auth_type: undefined as 'password' | 'private_key' | undefined,
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
})

const columns = useTableColumns(
  [
    { title: '凭据名称', key: 'name', width: 240, fixed: 'left' },
    { title: '认证方式', key: 'auth_type', width: 120, align: 'center' },
    { title: '登录用户', dataIndex: 'username', key: 'username', width: 180 },
    { title: '绑定主机', key: 'bind_count', width: 110, align: 'center' },
    { title: '备注', dataIndex: 'remark', key: 'remark', ellipsis: true },
  ],
  { title: '操作', key: 'action', width: 180, fixed: 'right', align: 'center' },
  ['cmdb:credential:view', 'cmdb:credential:update', 'cmdb:credential:delete']
)

const canViewCredential = computed(() => hasPermission('cmdb:credential:view'))
const passwordCount = computed(() => tableData.value.filter(item => item.auth_type === 'password').length)
const privateKeyCount = computed(() => tableData.value.filter(item => item.auth_type === 'private_key').length)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getCmdbCredentials({
      page: pagination.current,
      page_size: pagination.pageSize,
      name: searchForm.name || undefined,
      auth_type: searchForm.auth_type,
    })
    tableData.value = res.data.list || []
    pagination.total = res.data.total || 0
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.auth_type = undefined
  pagination.current = 1
  fetchData()
}

const handleTableChange = (pager: { current?: number; pageSize?: number }) => {
  pagination.current = pager.current || 1
  pagination.pageSize = pager.pageSize || 10
  fetchData()
}

const handleAdd = () => {
  drawerMode.value = 'create'
  drawerTitle.value = '新增凭据'
  editingId.value = undefined
  detailData.value = null
  drawerInitialValue.value = {
    name: '',
    auth_type: 'password',
    username: '',
    remark: '',
  }
  drawerVisible.value = true
}

const handleEdit = (record: CmdbCredentialSummary) => {
  drawerMode.value = 'edit'
  drawerTitle.value = '编辑凭据'
  editingId.value = record.id
  detailData.value = null
  drawerInitialValue.value = {
    name: record.name,
    auth_type: record.auth_type,
    username: record.username,
    remark: record.remark,
  }
  drawerVisible.value = true
}

const handleView = async (record: CmdbCredentialSummary) => {
  drawerMode.value = 'view'
  drawerTitle.value = '凭据详情'
  editingId.value = record.id
  drawerInitialValue.value = {}
  detailData.value = record
  drawerVisible.value = true
  try {
    const res = await getCmdbCredential(record.id)
    detailData.value = res.data
  } catch {
    detailData.value = record
  }
}

const handleSubmit = async (values: CmdbCredentialPayload) => {
  try {
    if (drawerMode.value === 'edit' && editingId.value) {
      await updateCmdbCredential(editingId.value, values)
      message.success('凭据更新成功')
    } else {
      await createCmdbCredential(values)
      message.success('凭据创建成功')
    }
    drawerVisible.value = false
    await fetchData()
  } catch {
    // handled by interceptor
  }
}

const handleDelete = (record: CmdbCredentialSummary) => {
  Modal.confirm({
    title: '确认删除凭据',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除凭据「${record.name}」吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deleteCmdbCredential(record.id)
      message.success('凭据删除成功')
      await fetchData()
    },
  })
}

fetchData()
</script>

<style scoped>
.cmdb-credential-page__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.summary-item {
  display: grid;
  grid-template-columns: 34px 1fr auto;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  color: var(--app-text-strong);
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.summary-item :deep(.anticon) {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  color: #1677ff;
  background: rgba(22, 119, 255, 0.1);
  border-radius: 8px;
}

.summary-item.success :deep(.anticon) {
  color: #389e0d;
  background: rgba(82, 196, 26, 0.12);
}

.summary-item.warning :deep(.anticon) {
  color: #d46b08;
  background: rgba(250, 140, 22, 0.12);
}

.summary-item span {
  color: var(--app-text-muted);
  font-size: 13px;
}

.summary-item strong {
  font-size: 22px;
  font-weight: 650;
}

.name-button {
  height: auto;
  padding: 0;
  text-align: left;
}

.name-cell {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.name-cell strong {
  color: var(--app-text-strong);
}

.name-cell span {
  color: var(--app-text-muted);
  font-size: 12px;
}

@media (max-width: 960px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }
}
</style>
