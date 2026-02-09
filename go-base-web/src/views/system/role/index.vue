<template>
  <div class="role-page">
    <ProTable
      :columns="columns"
      :data-source="tableData"
      :loading="loading"
    >
      <!-- 工具栏 -->
      <template #toolbar>
        <a-button type="primary" @click="handleAdd" v-permission="'system:role:add'">
          <PlusOutlined /> 新增
        </a-button>
      </template>

      <!-- 表格单元格 -->
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template v-if="column.key === 'action'">
          <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'system:role:edit'">编辑</a-button>
          <a-button type="link" size="small" @click="handleAssignPermissions(record)" v-permission="'system:role:assign'">分配权限</a-button>
          <a-popconfirm title="确定删除吗？" @confirm="handleDelete(record)">
            <a-button type="link" size="small" danger v-permission="'system:role:delete'">删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </ProTable>

    <!-- 新增/编辑弹窗 -->
    <a-modal v-model:open="modalVisible" :title="modalTitle" @ok="handleModalOk">
      <a-form :model="formState" :label-col="{ span: 5 }">
        <a-form-item label="角色名称" required>
          <a-input v-model:value="formState.name" />
        </a-form-item>
        <a-form-item label="角色编码" required>
          <a-input v-model:value="formState.code" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="formState.sort" :min="0" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model:checked="formState.statusChecked" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model:value="formState.remark" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 分配权限组件 -->
    <AssignPermission
      v-model:open="permissionDrawerVisible"
      :role-id="currentId"
      :role-name="currentRoleName"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import ProTable from '@/components/ProTable.vue'
import AssignPermission from './AssignPermission.vue'
import { getRoleList, createRole, updateRole, deleteRole } from '@/api/role'
import { useTableColumns } from '@/utils/permission'
import type { Role } from '@/types'

const loading = ref(false)
const tableData = ref<Role[]>([])
const modalVisible = ref(false)
const modalTitle = ref('新增角色')
const isEdit = ref(false)
const currentId = ref(0)
const currentRoleName = ref('')
const permissionDrawerVisible = ref(false)

const formState = reactive({
  name: '',
  code: '',
  sort: 0,
  statusChecked: true,
  remark: ''
})

// 使用工具函数动态生成列配置
const columns = useTableColumns(
  [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '角色名称', dataIndex: 'name', key: 'name' },
    { title: '角色编码', dataIndex: 'code', key: 'code' },
    { title: '排序', dataIndex: 'sort', key: 'sort', width: 80 },
    { title: '状态', key: 'status', width: 80 },
    { title: '备注', dataIndex: 'remark', key: 'remark' },
  ],
  { title: '操作', key: 'action', width: 200 },
  ['system:role:edit', 'system:role:delete', 'system:role:assign']
)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getRoleList()
    tableData.value = res.data
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  modalTitle.value = '新增角色'
  Object.assign(formState, { name: '', code: '', sort: 0, statusChecked: true, remark: '' })
  modalVisible.value = true
}

const handleEdit = (record: Role) => {
  isEdit.value = true
  modalTitle.value = '编辑角色'
  currentId.value = record.id
  Object.assign(formState, {
    name: record.name,
    code: record.code,
    sort: record.sort,
    statusChecked: record.status === 1,
    remark: record.remark
  })
  modalVisible.value = true
}

const handleModalOk = async () => {
  const data = {
    name: formState.name,
    code: formState.code,
    sort: formState.sort,
    status: formState.statusChecked ? 1 : 0,
    remark: formState.remark
  }
  if (isEdit.value) {
    await updateRole(currentId.value, data)
    message.success('更新成功')
  } else {
    await createRole(data)
    message.success('创建成功')
  }
  modalVisible.value = false
  fetchData()
}

const handleDelete = async (record: Role) => {
  await deleteRole(record.id)
  message.success('删除成功')
  fetchData()
}

const handleAssignPermissions = (record: Role) => {
  currentId.value = record.id
  currentRoleName.value = record.name
  permissionDrawerVisible.value = true
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* 角色页面样式 */
</style>
