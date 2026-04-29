<template>
  <a-drawer
    v-model:open="visible"
    :title="`分配权限 - ${roleName}`"
    width="1200"
    placement="right"
    class="permission-drawer"
  >
    <div class="permission-header">
      <span class="selected-count">已选菜单 {{ assignableSelectedMenuKeys.length }} 个</span>
      <span class="selected-count">已选 API {{ checkedApiIds.length }} 个</span>
      <a-input-search
        v-model:value="searchText"
        :placeholder="searchPlaceholder"
        style="width: 280px; margin-left: auto"
        allow-clear
      />
    </div>

    <div class="permission-layout">
      <div class="permission-sidebar">
        <RolePermissionDrawerSidebar
          :top-menus="topMenus"
          :selected-top-menu-id="selectedTopMenuId"
          :checked-menu-keys="assignableSelectedMenuKeys"
          @select="selectedTopMenuId = $event"
        />
      </div>

      <div class="permission-content">
        <template v-if="topMenus.length">
          <a-tabs v-model:activeKey="activeTab" class="permission-tabs">
            <a-tab-pane key="menus" tab="菜单权限">
              <a-tabs
                v-if="menuTabSections.length"
                v-model:activeKey="activeMenuSectionId"
                size="small"
                class="section-tabs"
              >
                <a-tab-pane
                  v-for="section in menuTabSections"
                  :key="section.id"
                  :tab="section.raw.pageMenu.name"
                >
                  <RolePermissionPageSection
                    mode="menu"
                    :section="section.raw"
                    :visible-menu-items="section.visibleMenuItems"
                    :visible-apis="section.visibleApis"
                    :exact-checked-menu-keys="assignableSelectedMenuKeys"
                    :checked-menu-keys="checkedMenuKeys"
                    :checked-api-ids="checkedApiIds"
                    @toggle-menu="handleMenuToggle"
                    @toggle-section-menus="handleSectionMenusToggle"
                    @keep-page-only="handleSectionKeepPageOnly"
                    @select-child-permissions="handleSectionSelectChildPermissions"
                    @clear-child-permissions="handleSectionClearChildPermissions"
                    @toggle-api="handleApiToggle"
                    @toggle-section-apis="handleSectionApisToggle"
                  />
                </a-tab-pane>
              </a-tabs>
              <a-empty
                v-else
                class="content-empty"
                description="当前一级菜单下暂无匹配的菜单权限项"
              />
            </a-tab-pane>

            <a-tab-pane key="apis" tab="API权限">
              <a-tabs
                v-if="apiTabItems.length"
                v-model:activeKey="activeApiSectionId"
                size="small"
                class="section-tabs"
              >
                <a-tab-pane
                  v-for="item in apiTabItems"
                  :key="item.id"
                  :tab="item.label"
                >
                  <RolePermissionPageSection
                    v-if="item.kind === 'section' && item.section"
                    mode="api"
                    :section="item.section.raw"
                    :visible-menu-items="item.section.visibleMenuItems"
                    :visible-apis="item.section.visibleApis"
                    :exact-checked-menu-keys="assignableSelectedMenuKeys"
                    :checked-menu-keys="checkedMenuKeys"
                    :checked-api-ids="checkedApiIds"
                    @toggle-menu="handleMenuToggle"
                    @toggle-section-menus="handleSectionMenusToggle"
                    @keep-page-only="handleSectionKeepPageOnly"
                    @select-child-permissions="handleSectionSelectChildPermissions"
                    @clear-child-permissions="handleSectionClearChildPermissions"
                    @toggle-api="handleApiToggle"
                    @toggle-section-apis="handleSectionApisToggle"
                  />

                  <div v-else class="system-api-groups">
                    <RolePermissionUncategorizedApis
                      v-for="group in item.groups"
                      :key="group.id"
                      :title="group.label"
                      :checked="getApiGroupChecked(group.apis)"
                      :indeterminate="getApiGroupIndeterminate(group.apis)"
                      :checked-api-ids="checkedApiIds"
                      :visible-apis="group.apis"
                      @toggle-all="handleApiGroupToggle(group.apis, $event)"
                      @toggle-api="handleApiToggle"
                    />
                  </div>
                </a-tab-pane>
              </a-tabs>
              <a-empty
                v-else
                class="content-empty"
                description="当前一级菜单下暂无匹配的 API 权限项"
              />
            </a-tab-pane>
          </a-tabs>
        </template>
        <a-empty v-else class="content-empty" description="暂无可分配权限数据" />
      </div>
    </div>

    <template #footer>
      <div style="display: flex; justify-content: flex-end; gap: 8px">
        <a-button @click="visible = false">取消</a-button>
        <a-button type="primary" :loading="saveLoading" @click="handleSavePermissions">保存</a-button>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import RolePermissionDrawerSidebar from './RolePermissionDrawerSidebar.vue'
import RolePermissionPageSection from './RolePermissionPageSection.vue'
import RolePermissionUncategorizedApis from './RolePermissionUncategorizedApis.vue'
import {
  groupUncategorizedApis,
  type UncategorizedApiGroup,
  type FilteredPermissionSection
} from './permissionDrawer'
import { useRolePermissionDrawer } from './useRolePermissionDrawer'
import type { Api } from '@/types'

interface Props {
  roleId: number
  roleName: string
}

const props = defineProps<Props>()
const visible = defineModel<boolean>('open', { default: false })
const activeTab = ref<'menus' | 'apis'>('menus')
const activeMenuSectionId = ref('')
const activeApiSectionId = ref('')

const {
  assignableSelectedMenuKeys,
  checkedApiIds,
  checkedMenuKeys,
  filteredSections,
  filteredUncategorizedApis,
  handleApiToggle,
  handleMenuToggle,
  handleSavePermissions,
  handleSectionClearChildPermissions,
  handleSectionKeepPageOnly,
  handleSectionApisToggle,
  handleSectionMenusToggle,
  handleSectionSelectChildPermissions,
  saveLoading,
  searchText,
  selectedTopMenuId,
  topMenus
} = useRolePermissionDrawer(props, visible)

const searchPlaceholder = computed(() => (
  activeTab.value === 'menus'
    ? '搜索菜单名称、权限码或路由'
    : '搜索接口路径、分组、方法或描述'
))

const menuTabSections = computed(() => filteredSections.value)
const systemApiTopMenuId = computed(() =>
  topMenus.value.find(menu => menu.name === '系统管理' || menu.path === '/system')?.id ?? null
)

type ApiTabItem =
  | { id: string; label: string; kind: 'section'; section: FilteredPermissionSection }
  | { id: string; label: string; kind: 'system'; groups: UncategorizedApiGroup[]; section?: undefined }

const apiTabItems = computed<ApiTabItem[]>(() => {
  const items: ApiTabItem[] = filteredSections.value.map(section => ({
    id: section.id,
    label: section.raw.pageMenu.name,
    kind: 'section',
    section
  }))

  const uncategorizedGroups = groupUncategorizedApis(filteredUncategorizedApis.value)
  if (
    uncategorizedGroups.length &&
    selectedTopMenuId.value != null &&
    selectedTopMenuId.value === systemApiTopMenuId.value
  ) {
    items.push({
      id: 'system-apis',
      label: '系统接口',
      kind: 'system',
      groups: uncategorizedGroups
    })
  }
  return items
})

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

watch(menuTabSections, sections => {
  if (!sections.length) {
    activeMenuSectionId.value = ''
    return
  }
  if (!sections.some(section => section.id === activeMenuSectionId.value)) {
    activeMenuSectionId.value = sections[0].id
  }
}, { immediate: true })

watch(apiTabItems, items => {
  if (!items.length) {
    activeApiSectionId.value = ''
    return
  }
  if (!items.some(item => item.id === activeApiSectionId.value)) {
    activeApiSectionId.value = items[0].id
  }
}, { immediate: true })

watch(selectedTopMenuId, () => {
  activeMenuSectionId.value = ''
  activeApiSectionId.value = ''
})
</script>

<style scoped>
.permission-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 16px;
}

.selected-count {
  color: #666;
  font-size: 13px;
}

.permission-layout {
  display: flex;
  height: calc(100vh - 250px);
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
}

.permission-sidebar {
  width: 240px;
  border-right: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
}

.permission-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: #fafafa;
}

.permission-tabs :deep(.ant-tabs-content-holder) {
  min-height: 100%;
}

.section-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 12px;
}

.section-tabs :deep(.ant-tabs-tab) {
  padding-top: 8px;
  padding-bottom: 8px;
}

.section-tabs :deep(.ant-tabs-content-holder) {
  min-height: 100%;
}

.content-empty {
  margin: auto 0;
}

.system-api-groups {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

@media (max-width: 1280px) {
  .permission-layout {
    height: auto;
    min-height: 640px;
    flex-direction: column;
  }

  .permission-sidebar {
    width: 100%;
    border-right: 0;
    border-bottom: 1px solid #f0f0f0;
  }
}
</style>
