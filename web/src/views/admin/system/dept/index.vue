<template>
  <div class="dept-page">
    <a-card>
      <template #extra>
        <a-button type="primary" @click="handleAdd()" v-permission="'system:dept:add'">
          <PlusOutlined /> 新增部门
        </a-button>
      </template>

      <a-table
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        row-key="id"
        :pagination="false"
        default-expand-all-rows
      >
        <template #bodyCell="{ column, record }">
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
              <a-popconfirm title="确定删除该部门吗？" @confirm="handleDelete(record)">
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
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
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

const columns = [
  { title: '部门名称', dataIndex: 'name', key: 'name', width: 220 },
  { title: '排序', dataIndex: 'sort', key: 'sort', width: 80 },
  { title: '状态', key: 'status', width: 100 },
  { title: '备注', dataIndex: 'remark', key: 'remark' },
  { title: '创建时间', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 220 }
]

const treeOptions = computed(() => [
  {
    id: 0,
    name: '顶级部门',
    children: tableData.value
  }
])

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getManageableDeptTree('system:dept-management')
    tableData.value = Array.isArray(res.data?.tree) ? res.data.tree : []
  } finally {
    loading.value = false
  }
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
