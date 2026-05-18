<template>
  <PageWrapper class="cmdb-group-page">
    <div class="cmdb-group-page__content">
      <div class="summary-grid">
        <div class="summary-item">
          <ApartmentOutlined />
          <span>分组总数</span>
          <strong>{{ tableData.length }}</strong>
        </div>
        <div class="summary-item success">
          <CheckCircleOutlined />
          <span>已启用</span>
          <strong>{{ enabledCount }}</strong>
        </div>
        <div class="summary-item warning">
          <StopOutlined />
          <span>已停用</span>
          <strong>{{ disabledCount }}</strong>
        </div>
      </div>

      <ProTable
        title="主机分组"
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        row-key="id"
        :scroll="{ x: 900 }"
        @search="handleSearch"
        @reset="handleReset"
      >
        <template #search>
          <a-form-item label="分组名称">
            <a-input v-model:value="searchForm.name" placeholder="请输入分组名称" allow-clear />
          </a-form-item>
          <a-form-item label="状态">
            <a-select v-model:value="searchForm.status" placeholder="请选择状态" allow-clear style="width: 140px">
              <a-select-option :value="1">启用</a-select-option>
              <a-select-option :value="0">停用</a-select-option>
            </a-select>
          </a-form-item>
        </template>

        <template #toolbar>
          <a-button type="primary" v-permission="'cmdb:group:create'" @click="handleAdd">
            <PlusOutlined /> 新增分组
          </a-button>
        </template>

        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div class="group-name-cell">
              <span class="group-name-cell__dot"></span>
              <span>{{ record.name }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="record.status === 1 ? 'success' : 'default'">{{ record.status === 1 ? '启用' : '停用' }}</a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space :size="0">
              <a-button type="link" size="small" v-permission="'cmdb:group:update'" @click="handleEdit(record)">编辑</a-button>
              <a-button type="link" size="small" danger v-permission="'cmdb:group:delete'" @click="handleDelete(record)">删除</a-button>
            </a-space>
          </template>
        </template>
      </ProTable>

      <GroupFormDrawer
        v-model:open="drawerVisible"
        :title="drawerTitle"
        :initial-value="drawerInitialValue"
        @submit="handleSubmit"
      />
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { ApartmentOutlined, CheckCircleOutlined, PlusOutlined, StopOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import ProTable from '@/components/ProTable.vue'
import { createCmdbHostGroup, deleteCmdbHostGroup, getCmdbHostGroups, updateCmdbHostGroup, type CmdbHostGroup } from '@/api/cmdb'
import { useTableColumns } from '@/utils/permission'
import GroupFormDrawer from './components/GroupFormDrawer.vue'

interface GroupDrawerValue {
  name: string
  sort: number
  remark: string
  statusChecked: boolean
}

const loading = ref(false)
const tableData = ref<CmdbHostGroup[]>([])
const drawerVisible = ref(false)
const drawerTitle = ref('新增分组')
const editingId = ref<number>()
const drawerInitialValue = ref<Partial<GroupDrawerValue>>({})

const searchForm = reactive({
  name: '',
  status: undefined as number | undefined,
})

const columns = useTableColumns(
  [
    { title: '分组名称', key: 'name', width: 220 },
    { title: '排序', dataIndex: 'sort', key: 'sort', width: 90, align: 'center' },
    { title: '状态', key: 'status', width: 100, align: 'center' },
    { title: '备注', dataIndex: 'remark', key: 'remark', ellipsis: true },
  ],
  { title: '操作', key: 'action', width: 150, fixed: 'right', align: 'center' },
  ['cmdb:group:update', 'cmdb:group:delete']
)

const enabledCount = computed(() => tableData.value.filter(item => item.status === 1).length)
const disabledCount = computed(() => tableData.value.filter(item => item.status !== 1).length)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getCmdbHostGroups({
      name: searchForm.name || undefined,
      status: searchForm.status,
    })
    tableData.value = res.data || []
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  fetchData()
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.status = undefined
  fetchData()
}

const handleAdd = () => {
  editingId.value = undefined
  drawerTitle.value = '新增分组'
  drawerInitialValue.value = { name: '', sort: 0, remark: '', statusChecked: true }
  drawerVisible.value = true
}

const handleEdit = (record: CmdbHostGroup) => {
  editingId.value = record.id
  drawerTitle.value = '编辑分组'
  drawerInitialValue.value = {
    name: record.name,
    sort: record.sort,
    remark: record.remark,
    statusChecked: record.status === 1,
  }
  drawerVisible.value = true
}

const handleSubmit = async (values: { name: string; sort: number; status: number; remark: string }) => {
  try {
    if (editingId.value) {
      await updateCmdbHostGroup(editingId.value, values)
      message.success('分组更新成功')
    } else {
      await createCmdbHostGroup(values)
      message.success('分组创建成功')
    }
    drawerVisible.value = false
    await fetchData()
  } catch {
    // handled by interceptor
  }
}

const handleDelete = (record: CmdbHostGroup) => {
  Modal.confirm({
    title: '确认删除分组',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除分组「${record.name}」吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deleteCmdbHostGroup(record.id)
      message.success('分组删除成功')
      await fetchData()
    },
  })
}

fetchData()
</script>

<style scoped>
.cmdb-group-page__content {
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

.group-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--app-text-strong);
  font-weight: 600;
}

.group-name-cell__dot {
  width: 8px;
  height: 8px;
  background: #1677ff;
  border-radius: 50%;
  box-shadow: 0 0 0 4px #e6f4ff;
}

@media (max-width: 960px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }
}
</style>
