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
          <a-switch
            :checked="record.status === 1"
            :loading="isRoleStatusLoading(record.id)"
            :disabled="isRoleRowUpdating(record.id) && !isRoleStatusLoading(record.id)"
            @change="(checked: boolean) => handleStatusToggle(record, checked)"
          />
        </template>
        <template v-if="column.key === 'is_super_admin'">
          <a-switch
            :checked="record.is_super_admin"
            :loading="isRoleSuperAdminLoading(record.id)"
            :disabled="isRoleRowUpdating(record.id) && !isRoleSuperAdminLoading(record.id)"
            @change="(checked: boolean) => handleSuperAdminToggle(record, checked)"
          />
        </template>
        <template v-if="column.key === 'users'">
          <div class="role-users-cell">
            <span class="role-users-count">{{ getRoleUserCount(record) }} 人</span>
            <a-button
              type="link"
              size="small"
              class="role-users-view"
              :disabled="getRoleUserCount(record) === 0"
              @click="handleViewRoleUsers(record)"
            >
              查看
            </a-button>
          </div>
        </template>
        <template v-if="column.key === 'data_scope'">
          {{ formatDataScope(record.data_scope) }}
        </template>
        <template v-if="column.key === 'action'">
          <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'system:role:edit'">编辑</a-button>
          <a-button type="link" size="small" @click="handleAssignPermissions(record)" v-permission="'system:role:assign'">分配权限</a-button>
          <a-button
            v-if="getRoleUserCount(record) > 0"
            type="link"
            size="small"
            danger
            @click="handleDelete(record)"
            v-permission="'system:role:delete'"
          >
            删除
          </a-button>
          <a-popconfirm v-else title="确定删除吗？" @confirm="handleDelete(record)">
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
      :dept-options="deptSelectOptions"
    />

    <RoleUsersDrawer
      v-model:open="roleUsersDrawerVisible"
      :role-name="roleUsersRoleName"
      :loading="roleUsersLoading"
      :users="roleUsers"
      :pagination="roleUsersPagination"
      @update:open="handleRoleUsersDrawerOpenChange"
      @page-change="handleRoleUsersPageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { PlusOutlined } from '@ant-design/icons-vue'
import ProTable from '@/components/ProTable.vue'
import RolePermissionDrawer from './components/RolePermissionDrawer.vue'
import RoleFormDrawer from './components/RoleFormDrawer.vue'
import RoleUsersDrawer from './components/RoleUsersDrawer.vue'
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
  handleRoleUsersDrawerOpenChange,
  handleRoleUsersPageChange,
  handleStatusToggle,
  handleSuperAdminToggle,
  handleViewRoleUsers,
  getRoleUserCount,
  isEdit,
  isRoleRowUpdating,
  isRoleStatusLoading,
  isRoleSuperAdminLoading,
  loading,
  permissionDrawerVisible,
  roleUsers,
  roleUsersDrawerVisible,
  roleUsersLoading,
  roleUsersPagination,
  roleUsersRoleName,
  tableData
} = useRolePage()
</script>

<style scoped>
.role-users-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  line-height: 1;
}

.role-users-count {
  color: #262626;
  font-weight: 500;
}

.role-users-view {
  padding-inline: 0;
}
</style>
