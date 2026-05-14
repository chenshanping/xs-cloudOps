<template>
  <PageWrapper class="dept-page">
    <div class="dept-page__content">
      <div class="dept-page__summary">
      <div class="summary-item summary-item--primary">
        <div class="summary-item__icon">
          <ApartmentOutlined />
        </div>
        <div>
          <div class="summary-item__label">可管理部门</div>
          <div class="summary-item__value">{{ deptCount }}</div>
        </div>
      </div>
      <div class="summary-item">
        <div class="summary-item__icon summary-item__icon--success">
          <TeamOutlined />
        </div>
        <div>
          <div class="summary-item__label">已绑定用户</div>
          <div class="summary-item__value">{{ boundUserCount }}</div>
        </div>
      </div>
      <div class="summary-item">
        <div class="summary-item__icon summary-item__icon--warning">
          <UsergroupAddOutlined />
        </div>
        <div>
          <div class="summary-item__label">未绑定用户</div>
          <div class="summary-item__value">{{ unassignedUserCount }}</div>
        </div>
      </div>
      <div class="summary-item">
        <div class="summary-item__icon summary-item__icon--muted">
          <CheckCircleOutlined />
        </div>
        <div>
          <div class="summary-item__label">启用部门</div>
          <div class="summary-item__value">{{ enabledDeptCount }}</div>
        </div>
      </div>
    </div>

      <a-card :bordered="false" class="dept-table-card">
        <div class="dept-table-card__header">
        <div>
          <h2 class="dept-table-card__title">部门结构</h2>
          <div class="dept-table-card__subtitle">人数统计区分直属用户与包含下级部门的总人数</div>
        </div>
        <a-button type="primary" @click="handleAdd()" v-permission="'system:dept:add'">
          <PlusOutlined /> 新增部门
        </a-button>
      </div>

        <a-table
          :columns="columns"
          :data-source="tableData"
          :loading="loading"
          row-key="id"
          :pagination="false"
          default-expand-all-rows
          :row-class-name="getRowClassName"
          :scroll="{ x: 1040 }"
        >
          <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div class="dept-name-cell">
              <div class="dept-name-cell__main">
                <span class="dept-name-cell__mark"></span>
                <span class="dept-name-cell__text">{{ record.name }}</span>
                <a-tag v-if="record.manageable === false" color="default">只读</a-tag>
              </div>
              <div class="dept-name-cell__meta">
                <span v-if="getChildDeptCount(record) > 0">{{ getChildDeptCount(record) }} 个下级</span>
                <span v-else>末级部门</span>
              </div>
            </div>
          </template>
          <template v-if="column.key === 'direct_user_count'">
            <a-tooltip title="仅统计直接归属于当前部门的用户">
              <span :class="['user-count-pill', { 'user-count-pill--empty': getDirectUserCount(record) === 0 }]">
                {{ getDirectUserCount(record) }}
              </span>
            </a-tooltip>
          </template>
          <template v-if="column.key === 'total_user_count'">
            <a-tooltip title="统计当前部门及所有下级部门绑定的用户">
              <span :class="['user-count-pill', 'user-count-pill--total', { 'user-count-pill--empty': getTotalUserCount(record) === 0 }]">
                {{ getTotalUserCount(record) }}
              </span>
            </a-tooltip>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 1 ? 'green' : 'red'">
              {{ record.status === 1 ? '启用' : '禁用' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'created_at'">
            {{ formatTime(record.created_at) }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space v-if="record.manageable !== false" :size="0">
              <a-button type="link" size="small" @click="handleAdd(record.id)" v-permission="'system:dept:add'">新增下级</a-button>
              <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'system:dept:edit'">编辑</a-button>
              <a-popconfirm :title="getDeleteConfirmTitle(record)" @confirm="handleDelete(record)">
                <a-button type="link" size="small" danger v-permission="'system:dept:delete'">删除</a-button>
              </a-popconfirm>
            </a-space>
            <span v-else>-</span>
          </template>
          </template>
        </a-table>
      </a-card>

    <DeptFormDrawer
      v-model:open="drawerVisible"
      :title="drawerTitle"
      :tree-options="treeOptions"
      :initial-value="drawerInitialValue"
      @submit="handleSubmit"
    />
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import { ApartmentOutlined, CheckCircleOutlined, PlusOutlined, TeamOutlined, UsergroupAddOutlined } from '@ant-design/icons-vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import DeptFormDrawer from './components/DeptFormDrawer.vue'
import { createDept, deleteDept, getManageableDeptTree, updateDept } from '@/api/dept'
import type { Dept } from '@/types'
import { formatTime } from '@/utils/format'

const loading = ref(false)
const tableData = ref<Dept[]>([])
const drawerVisible = ref(false)
const drawerTitle = ref('新增部门')
const currentId = ref(0)
const drawerInitialValue = ref<Record<string, any>>({})
const unassignedUserCount = ref(0)

const columns = [
  { title: '部门名称', dataIndex: 'name', key: 'name', width: 300 },
  { title: '直属用户', key: 'direct_user_count', width: 120, align: 'center' },
  { title: '含下级用户', key: 'total_user_count', width: 140, align: 'center' },
  { title: '排序', dataIndex: 'sort', key: 'sort', width: 80 },
  { title: '状态', key: 'status', width: 100 },
  { title: '备注', dataIndex: 'remark', key: 'remark' },
  { title: '创建时间', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 220, fixed: 'right' }
]

const treeOptions = computed(() => [
  {
    id: 0,
    name: '顶级部门',
    children: tableData.value
  }
])

const flattenDepts = (depts: Dept[]): Dept[] =>
  depts.flatMap((dept) => [dept, ...(dept.children ? flattenDepts(dept.children) : [])])

const allDepts = computed(() => flattenDepts(tableData.value))
const deptCount = computed(() => allDepts.value.length)
const enabledDeptCount = computed(() => allDepts.value.filter((dept) => dept.status === 1).length)
const boundUserCount = computed(() => {
  const rootTotal = tableData.value.reduce((sum, dept) => sum + getTotalUserCount(dept), 0)
  if (rootTotal > 0) {
    return rootTotal
  }
  return allDepts.value.reduce((sum, dept) => sum + getDirectUserCount(dept), 0)
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getManageableDeptTree('system:dept-management')
    tableData.value = Array.isArray(res.data?.tree) ? res.data.tree : []
    unassignedUserCount.value = res.data?.unassigned_user_count || 0
  } finally {
    loading.value = false
  }
}

const getDirectUserCount = (record: Dept) => record.direct_user_count || 0
const getTotalUserCount = (record: Dept) => record.total_user_count || getDirectUserCount(record)
const getChildDeptCount = (record: Dept) => record.children?.length || 0
const getRowClassName = (record: Dept) => record.manageable === false ? 'dept-row--readonly' : ''

const getDeleteConfirmTitle = (record: Dept) => {
  const total = getTotalUserCount(record)
  if (total > 0) {
    return `该部门及下级当前绑定 ${total} 名用户，确定继续删除吗？`
  }
  return '确定删除该部门吗？'
}

const handleAdd = (parentId = 0) => {
  currentId.value = 0
  drawerTitle.value = '新增部门'
  drawerInitialValue.value = {
    parent_id: parentId,
    name: '',
    sort: 0,
    statusChecked: true,
    remark: ''
  }
  drawerVisible.value = true
}

const handleEdit = (record: Dept) => {
  currentId.value = record.id
  drawerTitle.value = '编辑部门'
  drawerInitialValue.value = {
    parent_id: record.parent_id || 0,
    name: record.name,
    sort: record.sort,
    statusChecked: record.status === 1,
    remark: record.remark || ''
  }
  drawerVisible.value = true
}

const handleSubmit = async (values: { parent_id: number; name: string; sort: number; status: number; remark: string }) => {
  if (currentId.value > 0) {
    await updateDept(currentId.value, values)
    message.success('更新成功')
  } else {
    await createDept(values)
    message.success('创建成功')
  }
  drawerVisible.value = false
  fetchData()
}

const handleDelete = async (record: Dept) => {
  await deleteDept(record.id)
  message.success('删除成功')
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dept-page__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.dept-page__summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.summary-item {
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 88px;
  padding: 16px;
  background: #fff;
  border: 1px solid #edf0f5;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgb(15 23 42 / 4%);
}

.summary-item--primary {
  border-color: #d6e4ff;
}

.summary-item__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  color: #1677ff;
  font-size: 20px;
  background: #eef5ff;
  border-radius: 8px;
}

.summary-item__icon--success {
  color: #389e0d;
  background: #f0f8e8;
}

.summary-item__icon--warning {
  color: #d46b08;
  background: #fff4e6;
}

.summary-item__icon--muted {
  color: #45556c;
  background: #f3f5f8;
}

.summary-item__label {
  color: #8c8c8c;
  font-size: 12px;
}

.summary-item__value {
  margin-top: 4px;
  color: #262626;
  font-size: 24px;
  font-weight: 650;
  line-height: 1;
}

.dept-table-card {
  border-radius: 8px;
  box-shadow: 0 4px 16px rgb(15 23 42 / 4%);
}

.dept-table-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.dept-table-card__title {
  margin: 0;
  color: #262626;
  font-size: 18px;
  font-weight: 650;
  line-height: 1.4;
}

.dept-table-card__subtitle {
  margin-top: 4px;
  color: #8c8c8c;
  font-size: 12px;
}

.dept-name-cell {
  min-width: 0;
}

.dept-name-cell__main {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.dept-name-cell__mark {
  width: 8px;
  height: 8px;
  background: #1677ff;
  border-radius: 50%;
  box-shadow: 0 0 0 4px #e6f4ff;
}

.dept-name-cell__text {
  overflow: hidden;
  color: #262626;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dept-name-cell__meta {
  display: flex;
  gap: 10px;
  margin-top: 4px;
  padding-left: 16px;
  color: #8c8c8c;
  font-size: 12px;
}

.user-count-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 38px;
  height: 26px;
  padding: 0 10px;
  color: #0958d9;
  font-weight: 650;
  background: #e6f4ff;
  border: 1px solid #bae0ff;
  border-radius: 999px;
}

.user-count-pill--total {
  color: #237804;
  background: #f6ffed;
  border-color: #b7eb8f;
}

.user-count-pill--empty {
  color: #8c8c8c;
  background: #fafafa;
  border-color: #f0f0f0;
}

:deep(.dept-row--readonly) {
  color: #8c8c8c;
}

:deep(.ant-table-thead > tr > th) {
  color: #595959;
  font-size: 12px;
  font-weight: 650;
  background: #fafafa;
}

@media (max-width: 1180px) {
  .dept-page__summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .dept-page__summary {
    grid-template-columns: 1fr;
  }

  .dept-table-card__header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
