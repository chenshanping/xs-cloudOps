<template>
  <a-drawer
    v-model:open="visible"
    :width="drawerWidth"
    placement="right"
    :mask-closable="false"
    destroyOnClose
    :class="['permission-drawer', { 'permission-drawer--dark': uiStore.isDark }]"
  >
    <template #title>
      <RoleDrawerContextHeader
        title="API 权限"
        :role-name="roleName"
      />
    </template>

    <div :class="['permission-shell', { 'permission-shell--dark': uiStore.isDark }]">
      <div class="permission-header">
        <span class="selected-count">直授 {{ checkedApiIds.length }} 个</span>
        <span class="selected-count">继承 {{ inheritedApiIds.length }} 个</span>
        <span class="selected-count">生效 {{ effectiveApiIds.length }} 个</span>
        <a-input-search
          v-model:value="searchText"
          placeholder="搜索分组、路径、方法或描述"
          :disabled="permissionLoading"
          style="width: 280px; margin-left: auto"
          allow-clear
        />
      </div>

      <a-spin :spinning="permissionLoading">
        <div class="permission-layout">
          <div class="permission-sidebar">
            <RolePermissionApiGroupSidebar
              :groups="groupItems"
              :selected-group-id="selectedGroupId"
              @select="selectedGroupId = $event"
            />
          </div>

          <div class="permission-content">
            <RolePermissionUncategorizedApis
              v-if="activeGroup"
              :title="activeGroup.label"
              :checked="getApiGroupChecked(activeGroup.apis)"
              :indeterminate="getApiGroupIndeterminate(activeGroup.apis)"
              :checked-api-ids="checkedApiIds"
              :inherited-api-ids="inheritedApiIds"
              :inherited-api-source-map="inheritedApiSourceMap"
              :visible-apis="activeGroup.apis"
              @toggle-all="handleApiGroupToggle(activeGroup.apis, $event)"
              @toggle-api="handleApiToggle"
            />
            <a-empty v-else class="content-empty" description="暂无可分配 API" />
          </div>
        </div>
      </a-spin>
    </div>

    <template #footer>
      <div style="display: flex; justify-content: flex-end; gap: 8px">
        <a-button @click="visible = false">取消</a-button>
        <a-button type="primary" :loading="saveLoading" :disabled="permissionLoading" @click="handleSavePermissions">
          保存
        </a-button>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import RoleDrawerContextHeader from './RoleDrawerContextHeader.vue'
import RolePermissionApiGroupSidebar from './RolePermissionApiGroupSidebar.vue'
import RolePermissionUncategorizedApis from './RolePermissionUncategorizedApis.vue'
import { useResponsiveDrawerWidth } from './useResponsiveDrawerWidth'
import { useRolePermissionApiDrawer } from './useRolePermissionApiDrawer'
import { useUiStore } from '@/store/ui'
import type { Api } from '@/types'

interface Props {
  roleId: number
  roleName: string
}

const props = defineProps<Props>()
const uiStore = useUiStore()
const visible = defineModel<boolean>('open', { default: false })
const { drawerWidth } = useResponsiveDrawerWidth(1100, 920)

const {
  activeGroup,
  apiGroups,
  checkedApiIds,
  effectiveApiIds,
  handleApiToggle,
  handleSavePermissions,
  inheritedApiIds,
  inheritedApiSourceMap,
  permissionLoading,
  saveLoading,
  searchText,
  selectedGroupId
} = useRolePermissionApiDrawer(props, visible)

const groupItems = computed(() =>
  apiGroups.value.map(group => ({
    id: group.id,
    label: group.label,
    count: group.apis.length
  }))
)

const getApiGroupChecked = (apis: Api[]) =>
  apis.length > 0 && apis.every(api => checkedApiIds.value.includes(api.id))

const getApiGroupIndeterminate = (apis: Api[]) => {
  const checkedCount = apis.filter(api => checkedApiIds.value.includes(api.id)).length
  return checkedCount > 0 && checkedCount < apis.length
}

const handleApiGroupToggle = (apis: Api[], checked: boolean) => {
  const ids = apis.map(api => api.id)
  checkedApiIds.value = checked
    ? Array.from(new Set([...checkedApiIds.value, ...ids]))
    : checkedApiIds.value.filter(id => !ids.includes(id))
}
</script>

<style scoped>
@import './rolePermissionDrawerShared.css';
</style>
