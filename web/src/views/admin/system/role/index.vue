<template>
  <div class="role-page">
    <ProTable
      :columns="columns"
      :data-source="tableData"
      :loading="loading"
    >
      <template #toolbar>
        <a-button type="primary" @click="handleAdd" v-permission="'system:role:add'">
          <PlusOutlined /> 新增
        </a-button>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template v-if="column.key === 'data_scope'">
          {{ formatDataScope(record.data_scope) }}
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

    <RoleFormDrawer
      v-model:open="drawerVisible"
      :title="drawerTitle"
      :is-edit="isEdit"
      :dept-options="deptSelectOptions"
      :initial-value="drawerInitialValue"
      @submit="handleDrawerSubmit"
    />

    <RolePermissionDrawer
      v-model:open="permissionDrawerVisible"
      :role-id="currentId"
      :role-name="currentRoleName"
    />
  </div>
</template>

<script setup lang="ts">
import { PlusOutlined } from '@ant-design/icons-vue'
import ProTable from '@/components/ProTable.vue'
import RolePermissionDrawer from './components/RolePermissionDrawer.vue'
import RoleFormDrawer from './components/RoleFormDrawer.vue'
import { useRolePage } from './hooks/useRolePage'

const {
  columns,
  currentId,
  currentRoleName,
  deptSelectOptions,
  drawerInitialValue,
  drawerTitle,
  drawerVisible,
  formatDataScope,
  handleAdd,
  handleAssignPermissions,
  handleDelete,
  handleDrawerSubmit,
  handleEdit,
  isEdit,
  loading,
  permissionDrawerVisible,
  tableData
} = useRolePage()
</script>
