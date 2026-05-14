<template>
  <PageWrapper class="role-page">
    <div class="role-page__content">
      <div class="role-summary">
        <div class="summary-item summary-item--primary">
          <div class="summary-item__icon">
            <TeamOutlined />
          </div>
          <div>
            <div class="summary-item__label">角色总数</div>
            <div class="summary-item__value">{{ roleCount }}</div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-item__icon summary-item__icon--success">
            <CheckCircleOutlined />
          </div>
          <div>
            <div class="summary-item__label">启用角色</div>
            <div class="summary-item__value">{{ enabledRoleCount }}</div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-item__icon summary-item__icon--warning">
            <CrownOutlined />
          </div>
          <div>
            <div class="summary-item__label">超管角色</div>
            <div class="summary-item__value">{{ superAdminRoleCount }}</div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-item__icon summary-item__icon--muted">
            <UsergroupAddOutlined />
          </div>
          <div>
            <div class="summary-item__label">关联用户</div>
            <div class="summary-item__value">{{ boundUserCount }}</div>
          </div>
        </div>
      </div>

      <div class="role-table-card">
        <div class="role-table-card__header">
          <div>
            <h2 class="role-table-card__title">角色列表</h2>
            <div class="role-table-card__subtitle">维护角色信息、权限分配和数据权限</div>
          </div>
          <a-space>
            <a-input-search
              v-model:value="searchKeyword"
              class="role-table-card__search"
              placeholder="搜索角色名称 / 编码"
              allow-clear
            />
            <a-button type="primary" @click="handleAdd" v-permission="'system:role:add'">
              <PlusOutlined /> 新增
            </a-button>
          </a-space>
        </div>

        <ProTable
          :columns="columns"
          :data-source="filteredTableData"
          :loading="loading"
          :scroll="{ x: 1160 }"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'name'">
              <div class="role-name-cell">
                <div class="role-name-cell__main">
                  <span class="role-name-cell__dot"></span>
                  <span>{{ record.name }}</span>
                  <a-tag v-if="record.is_super_admin" color="orange">超管</a-tag>
                </div>
                <div class="role-name-cell__code">{{ record.code }}</div>
              </div>
            </template>
            <template v-if="column.key === 'code'">
              <a-typography-text code>{{ record.code }}</a-typography-text>
            </template>
            <template v-if="column.key === 'status'">
              <a-switch
                :checked="record.status === 1"
                :loading="isRoleStatusLoading(record.id)"
                :disabled="isRoleRowUpdating(record.id) && !isRoleStatusLoading(record.id)"
                @change="(checked: boolean) => handleStatusToggle(record, checked)"
                v-permission="'system:role:edit'"
              />
            </template>
            <template v-if="column.key === 'is_super_admin'">
              <a-switch
                :checked="record.is_super_admin"
                :loading="isRoleSuperAdminLoading(record.id)"
                :disabled="isRoleRowUpdating(record.id) && !isRoleSuperAdminLoading(record.id)"
                @change="(checked: boolean) => handleSuperAdminToggle(record, checked)"
                v-permission="'system:role:edit'"
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
              <a-space :size="0" class="role-action-group">
                <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'system:role:edit'">编辑</a-button>
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
                <template v-if="hasAssignPermission || canUseAdvancedApiPermissions || hasDataScopePermission">
                  <a-dropdown :trigger="['click']">
                    <a-button type="link" size="small" class="role-action-more" @click.stop>
                      更多
                      <DownOutlined />
                    </a-button>
                    <template #overlay>
                      <a-menu @click="({ key }) => handleMoreMenuClick(record, key as string)">
                        <a-menu-item
                          v-if="hasAssignPermission"
                          key="menu"
                          :disabled="menuPermissionDrawerVisible && currentId === record.id"
                        >
                          <AppstoreOutlined />
                          <span>菜单权限</span>
                        </a-menu-item>
                        <a-menu-item
                          v-if="canUseAdvancedApiPermissions"
                          key="api"
                          :disabled="apiPermissionDrawerVisible && currentId === record.id"
                        >
                          <ApiOutlined />
                          <span>API 权限</span>
                        </a-menu-item>
                        <a-menu-item
                          v-if="hasDataScopePermission"
                          key="data"
                          :disabled="dataScopeDrawerVisible && currentId === record.id"
                        >
                          <DatabaseOutlined />
                          <span>数据权限</span>
                        </a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                </template>
              </a-space>
            </template>
          </template>
        </ProTable>
      </div>

      <RoleFormDrawer
        v-model:open="drawerVisible"
        :title="drawerTitle"
        :is-edit="isEdit"
        :initial-value="drawerInitialValue"
        @submit="handleDrawerSubmit"
      />

      <RolePermissionMenuDrawer
        v-model:open="menuPermissionDrawerVisible"
        :role-id="currentId"
        :role-name="currentRoleName"
      />

      <RolePermissionApiDrawer
        v-if="canUseAdvancedApiPermissions"
        v-model:open="apiPermissionDrawerVisible"
        :role-id="currentId"
        :role-name="currentRoleName"
      />

      <RoleDataScopeDrawer
        v-model:open="dataScopeDrawerVisible"
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
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  ApiOutlined,
  AppstoreOutlined,
  CheckCircleOutlined,
  CrownOutlined,
  DatabaseOutlined,
  DownOutlined,
  PlusOutlined,
  TeamOutlined,
  UsergroupAddOutlined
} from '@ant-design/icons-vue'
import ProTable from '@/components/ProTable.vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import RolePermissionMenuDrawer from './components/RolePermissionMenuDrawer.vue'
import RolePermissionApiDrawer from './components/RolePermissionApiDrawer.vue'
import RoleDataScopeDrawer from './components/RoleDataScopeDrawer.vue'
import RoleFormDrawer from './components/RoleFormDrawer.vue'
import RoleUsersDrawer from './components/RoleUsersDrawer.vue'
import { useRolePage } from './hooks/useRolePage'
import { useUserStore } from '@/store/user'
import { usePermission } from '@/utils/permission'
import type { Role } from '@/types'

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
  handleAssignApiPermissions,
  handleAssignDataScope,
  handleAssignMenuPermissions,
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
  apiPermissionDrawerVisible,
  dataScopeDrawerVisible,
  filteredTableData,
  menuPermissionDrawerVisible,
  roleUsers,
  roleUsersDrawerVisible,
  roleUsersLoading,
  roleUsersPagination,
  roleUsersRoleName,
  searchKeyword,
  tableData
} = useRolePage()

const userStore = useUserStore()
const { hasPermission } = usePermission()
const canUseAdvancedApiPermissions = computed(() => userStore.permissions.includes('*'))
const hasAssignPermission = computed(() => hasPermission('system:role:assign'))
const hasDataScopePermission = computed(() => hasPermission('system:role:dataScope'))

const handleMoreMenuClick = (record: Role, key: string) => {
  if (key === 'menu') {
    handleAssignMenuPermissions(record)
  } else if (key === 'api') {
    handleAssignApiPermissions(record)
  } else if (key === 'data') {
    handleAssignDataScope(record)
  }
}

const roleCount = computed(() => tableData.value.length)
const enabledRoleCount = computed(() => tableData.value.filter(role => role.status === 1).length)
const superAdminRoleCount = computed(() => tableData.value.filter(role => role.is_super_admin).length)
const boundUserCount = computed(() => tableData.value.reduce((sum, role) => sum + getRoleUserCount(role), 0))
</script>

<style scoped>
.role-page__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.role-summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 16px;
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

.role-table-card {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.role-table-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.role-table-card__title {
  margin: 0;
  color: #262626;
  font-size: 18px;
  font-weight: 650;
  line-height: 1.4;
}

.role-table-card__subtitle {
  margin-top: 4px;
  color: #8c8c8c;
  font-size: 12px;
}

.role-table-card__search {
  width: 280px;
}

.role-name-cell__main {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #262626;
  font-weight: 600;
}

.role-name-cell__dot {
  width: 8px;
  height: 8px;
  background: #1677ff;
  border-radius: 50%;
  box-shadow: 0 0 0 4px #e6f4ff;
}

.role-name-cell__code {
  margin-top: 4px;
  padding-left: 16px;
  color: #8c8c8c;
  font-size: 12px;
}

.role-users-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  line-height: 1;
}

.role-users-count {
  color: var(--app-text-color, #262626);
  font-weight: 500;
}

.role-users-view {
  padding-inline: 0;
}

.role-action-group {
  flex-wrap: nowrap;
  gap: 2px;
}

.role-action-group :deep(.ant-space-item) {
  display: inline-flex;
  align-items: center;
}

.role-action-more {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.role-table-card :deep(.ant-table-thead > tr > th) {
  color: #595959;
  font-size: 12px;
  font-weight: 650;
  background: #fafafa;
}

@media (max-width: 1180px) {
  .role-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .role-summary {
    grid-template-columns: 1fr;
  }

  .role-table-card__header {
    align-items: flex-start;
    flex-direction: column;
  }

  .role-table-card__search {
    width: 100%;
  }
}
</style>
