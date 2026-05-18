<template>
  <PageWrapper class="cmdb-host-page">
    <div class="cmdb-host-page__content">
      <div class="summary-grid">
        <div class="summary-item">
          <DesktopOutlined />
          <span>主机总数</span>
          <strong>{{ pagination.total }}</strong>
        </div>
        <div class="summary-item success">
          <CheckCircleOutlined />
          <span>校验成功</span>
          <strong>{{ successCount }}</strong>
        </div>
        <div class="summary-item danger">
          <CloseCircleOutlined />
          <span>校验失败</span>
          <strong>{{ failedCount }}</strong>
        </div>
        <div class="summary-item warning">
          <ClockCircleOutlined />
          <span>待校验</span>
          <strong>{{ pendingCount }}</strong>
        </div>
      </div>

      <ProTable
        title="主机台账"
        :columns="columns"
        :data-source="tableData"
        :loading="loading || referenceLoading"
        :pagination="pagination"
        row-key="id"
        :scroll="{ x: 1680 }"
        @search="handleSearch"
        @reset="handleReset"
        @change="handleTableChange"
      >
        <template #search>
          <a-form-item label="关键字">
            <a-input v-model:value="searchForm.keyword" placeholder="主机名 / IP / SSH 地址 / 系统主机名" allow-clear />
          </a-form-item>
          <a-form-item label="分组">
            <a-select v-model:value="searchForm.group_id" placeholder="请选择分组" allow-clear style="width: 170px">
              <a-select-option v-for="item in groupOptions" :key="item.id" :value="item.id">
                {{ item.name }}
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="标签">
            <a-select v-model:value="searchForm.tag_id" placeholder="请选择标签" allow-clear style="width: 170px">
              <a-select-option v-for="item in tagOptions" :key="item.id" :value="item.id">
                {{ item.name }}
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="校验状态">
            <a-select v-model:value="searchForm.verify_status" placeholder="请选择状态" allow-clear style="width: 160px">
              <a-select-option v-for="item in cmdbVerifyStatusOptions" :key="item.value" :value="item.value">
                {{ item.label }}
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="环境">
            <a-input v-model:value="searchForm.environment" placeholder="如 prod / test" allow-clear />
          </a-form-item>
        </template>

        <template #toolbar>
          <a-space>
            <a-button type="primary" v-permission="'cmdb:host:create'" @click="handleAdd">
              <PlusOutlined /> 新增主机
            </a-button>
            <a-button v-permission="'cmdb:host:import'" @click="handleOpenImport">
              <UploadOutlined /> 导入主机
            </a-button>
          </a-space>
        </template>

        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a-button
              v-if="canViewHost"
              type="link"
              class="name-button"
              v-permission="'cmdb:host:view'"
              @click="handleView(record)"
            >
              <div class="name-cell">
                <strong>{{ record.name }}</strong>
                <span>{{ record.hostname || record.ssh_host }}</span>
              </div>
            </a-button>
            <div v-else class="name-cell">
              <strong>{{ record.name }}</strong>
              <span>{{ record.hostname || record.ssh_host }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'group'">
            <a-tag color="blue">{{ record.group.name }}</a-tag>
          </template>
          <template v-else-if="column.key === 'tags'">
            <a-space v-if="record.tags?.length" :size="[4, 4]" wrap>
              <a-tag v-for="item in record.tags" :key="item.id" :color="item.color || '#1677ff'">{{ item.name }}</a-tag>
            </a-space>
            <span v-else class="muted-text">-</span>
          </template>
          <template v-else-if="column.key === 'network'">
            <div class="ip-stack">
              <span>内网：{{ record.private_ip || '-' }}</span>
              <span>公网：{{ record.public_ip || '-' }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'ssh'">
            <div class="ip-stack">
              <span>{{ record.ssh_host }}:{{ record.ssh_port }}</span>
              <span class="muted-text">{{ record.credential_summary.name }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'verify_status'">
            <div class="verify-cell">
              <a-tag :color="getCmdbVerifyStatusMeta(record.verify_status).badge">
                {{ getCmdbVerifyStatusMeta(record.verify_status).text }}
              </a-tag>
              <a-tooltip :title="record.verify_message || '暂无校验消息'">
                <span class="verify-cell__message">{{ record.verify_message || '暂无校验消息' }}</span>
              </a-tooltip>
            </div>
          </template>
          <template v-else-if="column.key === 'system'">
            <div class="system-stack">
              <span class="system-stack__primary">{{ record.platform_version || record.platform || record.os || '-' }}</span>
              <span class="system-stack__secondary">
                内核：{{ record.kernel_version || '-' }} · 规格：{{ record.cpu_cores || 0 }}C / {{ formatCmdbMemory(record.memory_mb) }}
              </span>
            </div>
          </template>
          <template v-else-if="column.key === 'updated_at'">
            {{ formatTime(record.updated_at) }}
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space :size="0">
              <a-button type="link" size="small" v-permission="'cmdb:host:view'" @click="handleView(record)">查看</a-button>
              <a-button type="link" size="small" v-permission="'cmdb:host:update'" @click="handleEdit(record)">编辑</a-button>
              <a-button type="link" size="small" v-permission="'cmdb:host:create'" @click="handleCopy(record)">复制</a-button>
              <a-button type="link" size="small" v-permission="'cmdb:host:verify'" @click="handleVerify(record)">校验</a-button>
              <a-button type="link" size="small" danger v-permission="'cmdb:host:delete'" @click="handleDelete(record)">删除</a-button>
            </a-space>
          </template>
        </template>
      </ProTable>

      <HostFormDrawer
        v-model:open="drawerVisible"
        :title="drawerTitle"
        :group-options="groupOptions"
        :tag-options="tagOptions"
        :credential-options="credentialOptions"
        :initial-value="drawerInitialValue"
        :submitting="submitLoading"
        @submit="handleSubmit"
      />

      <HostImportDrawer
        v-model:open="importVisible"
        :importing="importLoading"
        :template-loading="templateLoading"
        :result="importResult"
        @submit="handleImport"
        @download-template="handleDownloadTemplate"
      />

      <HostDetailDrawer
        v-model:open="detailVisible"
        :host="detailData"
        :verifying="detailVerifyLoading"
        @verify="handleVerify"
      />
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import {
  CheckCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  DesktopOutlined,
  ExclamationCircleOutlined,
  PlusOutlined,
  UploadOutlined,
} from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import ProTable from '@/components/ProTable.vue'
import {
  createCmdbHost,
  deleteCmdbHost,
  downloadCmdbHostImportTemplate,
  getCmdbHost,
  getCmdbHosts,
  importCmdbHosts,
  updateCmdbHost,
  verifyCmdbHost,
  type CmdbHostItem,
  type CmdbHostPayload,
} from '@/api/cmdb'
import { formatTime } from '@/utils/format'
import { usePermission, useTableColumns } from '@/utils/permission'
import {
  cmdbVerifyStatusOptions,
  formatCmdbMemory,
  getCmdbVerifyStatusMeta,
  useCmdbReferenceOptions,
} from '../shared'
import HostDetailDrawer from './components/HostDetailDrawer.vue'
import HostFormDrawer from './components/HostFormDrawer.vue'
import HostImportDrawer from './components/HostImportDrawer.vue'

interface HostDrawerValue {
  name: string
  group_id?: number
  tag_ids: number[]
  environment: string
  owner: string
  private_ip: string
  public_ip: string
  ssh_host: string
  ssh_port: number
  credential_id?: number
  remark: string
}

const { hasPermission } = usePermission()
const {
  loading: referenceLoading,
  groupOptions,
  tagOptions,
  credentialOptions,
  loadReferenceOptions,
} = useCmdbReferenceOptions()

const loading = ref(false)
const tableData = ref<CmdbHostItem[]>([])
const drawerVisible = ref(false)
const drawerTitle = ref('新增主机')
const editingId = ref<number>()
const submitLoading = ref(false)
const drawerInitialValue = ref<Partial<HostDrawerValue>>({})
const detailVisible = ref(false)
const detailData = ref<CmdbHostItem | null>(null)
const detailVerifyLoading = ref(false)
const importVisible = ref(false)
const importLoading = ref(false)
const templateLoading = ref(false)
const importResult = ref<any>(null)

const searchForm = reactive({
  keyword: '',
  group_id: undefined as number | undefined,
  tag_id: undefined as number | undefined,
  verify_status: undefined as string | undefined,
  environment: '',
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
})

const columns = useTableColumns(
  [
    { title: '主机', key: 'name', width: 240, fixed: 'left' },
    { title: '分组', key: 'group', width: 130, align: 'center' },
    { title: '标签', key: 'tags', width: 220 },
    { title: '环境', dataIndex: 'environment', key: 'environment', width: 100, align: 'center' },
    { title: '网络信息', key: 'network', width: 220 },
    { title: 'SSH 连接', key: 'ssh', width: 220 },
    { title: '校验状态', key: 'verify_status', width: 220 },
    { title: '系统信息', key: 'system', width: 280 },
    { title: '更新时间', key: 'updated_at', width: 170, align: 'center' },
  ],
  { title: '操作', key: 'action', width: 260, fixed: 'right', align: 'center' },
  ['cmdb:host:view', 'cmdb:host:update', 'cmdb:host:create', 'cmdb:host:verify', 'cmdb:host:delete']
)

const canViewHost = computed(() => hasPermission('cmdb:host:view'))
const successCount = computed(() => tableData.value.filter(item => item.verify_status === 'success').length)
const failedCount = computed(() => tableData.value.filter(item => item.verify_status === 'failed').length)
const pendingCount = computed(() => tableData.value.filter(item => item.verify_status === 'pending').length)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getCmdbHosts({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword || undefined,
      group_id: searchForm.group_id,
      tag_id: searchForm.tag_id,
      verify_status: searchForm.verify_status,
      environment: searchForm.environment || undefined,
    })
    tableData.value = res.data.list || []
    pagination.total = res.data.total || 0
  } finally {
    loading.value = false
  }
}

const refreshAll = async () => {
  await Promise.all([loadReferenceOptions(), fetchData()])
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.group_id = undefined
  searchForm.tag_id = undefined
  searchForm.verify_status = undefined
  searchForm.environment = ''
  pagination.current = 1
  fetchData()
}

const handleTableChange = (pager: { current?: number; pageSize?: number }) => {
  pagination.current = pager.current || 1
  pagination.pageSize = pager.pageSize || 10
  fetchData()
}

const handleAdd = () => {
  editingId.value = undefined
  drawerTitle.value = '新增主机'
  drawerInitialValue.value = {
    name: '',
    group_id: groupOptions.value[0]?.id,
    tag_ids: [],
    environment: '',
    owner: '',
    private_ip: '',
    public_ip: '',
    ssh_host: '',
    ssh_port: 22,
    credential_id: credentialOptions.value[0]?.id,
    remark: '',
  }
  drawerVisible.value = true
}

const handleEdit = (record: CmdbHostItem) => {
  editingId.value = record.id
  drawerTitle.value = '编辑主机'
  drawerInitialValue.value = {
    name: record.name,
    group_id: record.group.id,
    tag_ids: record.tags?.map(item => item.id) || [],
    environment: record.environment,
    owner: record.owner,
    private_ip: record.private_ip,
    public_ip: record.public_ip,
    ssh_host: record.ssh_host,
    ssh_port: record.ssh_port,
    credential_id: record.credential_summary.id,
    remark: record.remark,
  }
  drawerVisible.value = true
}

const handleCopy = (record: CmdbHostItem) => {
  const newName = `${record.name}-copy`
  Modal.confirm({
    title: '复制主机',
    icon: createVNode(ExclamationCircleOutlined),
    content: `将以「${newName}」为名创建一台新主机，分组、标签、SSH 凭据等配置会从「${record.name}」复制过来，IP 留空待校验回填。是否继续？`,
    okText: '继续',
    cancelText: '取消',
    onOk() {
      editingId.value = undefined
      drawerTitle.value = '复制主机'
      drawerInitialValue.value = {
        name: newName,
        group_id: record.group.id,
        tag_ids: record.tags?.map(item => item.id) || [],
        environment: record.environment,
        owner: record.owner,
        // 复制场景下不带源主机的 IP：新主机是另一台机器，留空以便 SSH 校验时自动回填。
        private_ip: '',
        public_ip: '',
        ssh_host: record.ssh_host,
        ssh_port: record.ssh_port,
        credential_id: record.credential_summary.id,
        remark: record.remark,
      }
      drawerVisible.value = true
    },
  })
}

const handleSubmit = async (values: CmdbHostPayload) => {
  if (submitLoading.value) {
    return
  }
  submitLoading.value = true
  try {
    if (editingId.value) {
      const res = await updateCmdbHost(editingId.value, values)
      message.success('主机更新成功')
      if (detailData.value?.id === editingId.value) {
        detailData.value = res.data
      }
    } else {
      await createCmdbHost(values)
      message.success('主机创建成功')
    }
    drawerVisible.value = false
    await fetchData()
  } catch {
    // handled by interceptor
  } finally {
    submitLoading.value = false
  }
}

const handleView = async (record: CmdbHostItem) => {
  detailData.value = record
  detailVisible.value = true
  try {
    const res = await getCmdbHost(record.id)
    detailData.value = res.data
  } catch {
    detailData.value = record
  }
}

const handleVerify = async (record: CmdbHostItem) => {
  const isDetailVerify = detailVisible.value && detailData.value?.id === record.id
  if (isDetailVerify) {
    detailVerifyLoading.value = true
  }
  try {
    const res = await verifyCmdbHost(record.id)
    message.success('主机校验已完成')
    if (detailData.value?.id === record.id) {
      detailData.value = res.data
    }
    await fetchData()
  } catch {
    // handled by interceptor
  } finally {
    if (isDetailVerify) {
      detailVerifyLoading.value = false
    }
  }
}

const handleDelete = (record: CmdbHostItem) => {
  Modal.confirm({
    title: '确认删除主机',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除主机「${record.name}」吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deleteCmdbHost(record.id)
      message.success('主机删除成功')
      if (detailData.value?.id === record.id) {
        detailVisible.value = false
        detailData.value = null
      }
      await fetchData()
    },
  })
}

const handleImport = async (file: File) => {
  importLoading.value = true
  try {
    const res = await importCmdbHosts(file)
    importResult.value = res.data
    if (res.data.success_count > 0) {
      await fetchData()
    }
    if (res.data.failure_count > 0) {
      message.warning(`导入完成：成功 ${res.data.success_count} 条，失败 ${res.data.failure_count} 条`)
    } else {
      message.success(`导入成功，共 ${res.data.success_count} 条`)
    }
  } catch {
    // handled by interceptor
  } finally {
    importLoading.value = false
  }
}

const handleOpenImport = () => {
  importResult.value = null
  importVisible.value = true
}

const handleDownloadTemplate = async () => {
  templateLoading.value = true
  try {
    const res = await downloadCmdbHostImportTemplate()
    const blob = new Blob([res as any], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = '主机导入模板.xlsx'
    link.click()
    window.URL.revokeObjectURL(url)
    message.success('模板下载成功')
  } catch {
    message.error('模板下载失败')
  } finally {
    templateLoading.value = false
  }
}

refreshAll()
</script>

<style scoped>
.cmdb-host-page__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
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

.summary-item.danger :deep(.anticon) {
  color: #cf1322;
  background: rgba(255, 77, 79, 0.12);
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

.name-cell span,
.muted-text {
  color: var(--app-text-muted);
  font-size: 12px;
}

.ip-stack {
  display: flex;
  flex-direction: column;
  gap: 3px;
  line-height: 1.6;
}

.system-stack {
  display: flex;
  flex-direction: column;
  gap: 3px;
  line-height: 1.6;
}

.system-stack__primary {
  color: var(--app-text-strong);
  font-weight: 500;
  word-break: break-word;
}

.system-stack__secondary {
  color: var(--app-text-muted);
  font-size: 12px;
  word-break: break-word;
}

.verify-cell {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.verify-cell__message {
  color: var(--app-text-muted);
  font-size: 12px;
  line-height: 1.6;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 1180px) {
  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }
}
</style>
