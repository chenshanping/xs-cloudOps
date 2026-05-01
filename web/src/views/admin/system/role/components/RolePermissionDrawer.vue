<template>
  <a-drawer
    v-model:open="visible"
    :title="`分配权限 - ${roleName}`"
    width="1200"
    placement="right"
    :class="['permission-drawer', { 'permission-drawer--dark': uiStore.isDark }]"
  >
    <div :class="['permission-shell', { 'permission-shell--dark': uiStore.isDark }]">
      <div class="permission-header">
        <span class="selected-count">已选菜单 {{ assignableSelectedMenuKeys.length }} 个</span>
        <span class="selected-count">直接 API {{ checkedApiIds.length }} 个</span>
        <span class="selected-count">继承 API {{ inheritedApiIds.length }} 个</span>
        <span class="selected-count">生效 API {{ effectiveApiIds.length }} 个</span>
        <a-input-search
          v-if="activeTab !== 'dataScopes'"
          v-model:value="searchText"
          :placeholder="searchPlaceholder"
          style="width: 280px; margin-left: auto"
          allow-clear
        />
      </div>

      <div
        class="permission-layout"
        :class="{ 'permission-layout--full': activeTab === 'dataScopes' }"
      >
        <div v-if="activeTab !== 'dataScopes'" class="permission-sidebar">
          <RolePermissionDrawerSidebar
            :top-menus="topMenus"
            :selected-top-menu-id="selectedTopMenuId"
            :checked-menu-keys="assignableSelectedMenuKeys"
            @select="selectedTopMenuId = $event"
          />
        </div>

        <div class="permission-content">
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
                    :inherited-api-ids="inheritedApiIds"
                    :inherited-api-source-map="inheritedApiSourceMap"
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
                    :inherited-api-ids="inheritedApiIds"
                    :inherited-api-source-map="inheritedApiSourceMap"
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
                      :inherited-api-ids="inheritedApiIds"
                      :inherited-api-source-map="inheritedApiSourceMap"
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

            <a-tab-pane key="dataScopes" tab="数据权限">
              <RolePermissionDataScopePanel
                v-model:model-value="featureDataScopes"
                :dept-options="deptOptions"
                :default-data-scope-label="formatDataScopeLabel(defaultDataScope)"
              />
            </a-tab-pane>
          </a-tabs>
        </div>
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
import RolePermissionDataScopePanel from './RolePermissionDataScopePanel.vue'
import RolePermissionPageSection from './RolePermissionPageSection.vue'
import RolePermissionUncategorizedApis from './RolePermissionUncategorizedApis.vue'
import {
  groupUncategorizedApis,
  type UncategorizedApiGroup,
  type FilteredPermissionSection
} from './permissionDrawer'
import { useRolePermissionDrawer } from './useRolePermissionDrawer'
import { useUiStore } from '@/store/ui'
import type { Api } from '@/types'
import type { RolePermissionDeptOption } from './dataScopeResources'

interface Props {
  roleId: number
  roleName: string
  deptOptions: RolePermissionDeptOption[]
}

const props = defineProps<Props>()
const uiStore = useUiStore()
const visible = defineModel<boolean>('open', { default: false })
const activeTab = ref<'menus' | 'apis' | 'dataScopes'>('menus')
const activeMenuSectionId = ref('')
const activeApiSectionId = ref('')

const {
  assignableSelectedMenuKeys,
  checkedApiIds,
  checkedMenuKeys,
  defaultDataScope,
  effectiveApiIds,
  featureDataScopes,
  filteredSections,
  filteredUncategorizedApis,
  formatDataScopeLabel,
  handleApiToggle,
  handleMenuToggle,
  handleSavePermissions,
  handleSectionClearChildPermissions,
  handleSectionKeepPageOnly,
  handleSectionApisToggle,
  handleSectionMenusToggle,
  handleSectionSelectChildPermissions,
  inheritedApiIds,
  inheritedApiSourceMap,
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
.permission-shell {
  --permission-surface: #ffffff;
  --permission-surface-soft: #fafafa;
  --permission-border: #f0f0f0;
  --permission-border-soft: #f5f5f5;
  --permission-hover: #f5f5f5;
  --permission-text-secondary: #666666;
  --permission-text-muted: #999999;
  --permission-text-strong: #333333;
  --permission-text-default: #262626;
  --permission-code-bg: #f5f5f5;
  color: var(--permission-text-default);
}

.permission-shell--dark {
  --permission-surface: #161d29;
  --permission-surface-soft: #1c2433;
  --permission-border: rgba(148, 163, 184, 0.2);
  --permission-border-soft: rgba(148, 163, 184, 0.16);
  --permission-hover: rgba(148, 163, 184, 0.12);
  --permission-text-secondary: rgba(255, 255, 255, 0.8);
  --permission-text-muted: rgba(255, 255, 255, 0.62);
  --permission-text-strong: rgba(255, 255, 255, 0.93);
  --permission-text-default: rgba(255, 255, 255, 0.87);
  --permission-code-bg: rgba(148, 163, 184, 0.14);
}

:deep(.permission-drawer .ant-drawer-header) {
  background: #ffffff;
  border-bottom-color: #f0f0f0;
}

:deep(.permission-drawer.permission-drawer--dark .ant-drawer-header) {
  background: #161d29;
  border-bottom-color: rgba(148, 163, 184, 0.2);
}

:deep(.permission-drawer .ant-drawer-title) {
  color: #333333;
}

:deep(.permission-drawer.permission-drawer--dark .ant-drawer-title) {
  color: rgba(255, 255, 255, 0.93);
}

:deep(.permission-drawer .ant-drawer-close) {
  color: #666666;
}

:deep(.permission-drawer.permission-drawer--dark .ant-drawer-close) {
  color: rgba(255, 255, 255, 0.8);
}

:deep(.permission-drawer .ant-drawer-close:hover) {
  color: #333333;
}

:deep(.permission-drawer.permission-drawer--dark .ant-drawer-close:hover) {
  color: rgba(255, 255, 255, 0.93);
}

:deep(.permission-drawer .ant-drawer-footer) {
  background: #ffffff;
  border-top-color: #f0f0f0;
}

:deep(.permission-drawer.permission-drawer--dark .ant-drawer-body) {
  background: #161d29;
  color: rgba(255, 255, 255, 0.87);
}

:deep(.permission-drawer.permission-drawer--dark .ant-drawer-footer) {
  background: #161d29;
  border-top-color: rgba(148, 163, 184, 0.2);
}

.permission-shell :deep(.ant-tabs-nav) {
  margin-bottom: 16px;
}

.permission-shell :deep(.ant-tabs-tab) {
  color: var(--permission-text-strong) !important;
}

.permission-shell :deep(.ant-tabs-tab:hover) {
  color: var(--permission-text-strong) !important;
}

.permission-shell :deep(.ant-tabs-tab-btn) {
  color: inherit !important;
}

.permission-shell :deep(.ant-tabs-tab-active .ant-tabs-tab-btn) {
  color: var(--app-primary-color, #006be6) !important;
}

.permission-shell :deep(.ant-tabs-ink-bar) {
  background: var(--app-primary-color, #006be6);
}

.permission-shell :deep(.ant-input-affix-wrapper),
.permission-shell :deep(.ant-input-search .ant-input),
.permission-shell :deep(.ant-input-group-addon) {
  background: var(--permission-surface);
  border-color: var(--permission-border);
  color: var(--permission-text-default);
}

.permission-shell :deep(.ant-input),
.permission-shell :deep(.ant-input-affix-wrapper input) {
  color: var(--permission-text-default);
}

.permission-shell :deep(.ant-input::placeholder),
.permission-shell :deep(.ant-input-affix-wrapper input::placeholder) {
  color: var(--permission-text-muted) !important;
}

.permission-shell :deep(.ant-input-prefix),
.permission-shell :deep(.ant-input-suffix),
.permission-shell :deep(.ant-input-search-button),
.permission-shell :deep(.ant-input-search-button .anticon) {
  color: var(--permission-text-secondary);
}

.permission-shell :deep(.ant-checkbox-wrapper),
.permission-shell :deep(.ant-checkbox-wrapper span),
.permission-shell :deep(.ant-form-item-label > label),
.permission-shell :deep(.ant-card-head-title),
.permission-shell :deep(.ant-alert-message),
.permission-shell :deep(.ant-alert-description) {
  color: var(--permission-text-default);
}

.permission-shell :deep(.ant-card) {
  background: var(--permission-surface);
  border-color: var(--permission-border);
}

.permission-shell :deep(.ant-card-head) {
  background: var(--permission-surface-soft);
  border-bottom-color: var(--permission-border);
}

.permission-shell :deep(.ant-card-body) {
  color: var(--permission-text-default);
}

.permission-shell :deep(.ant-empty-description) {
  color: var(--permission-text-secondary);
}

.permission-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid var(--permission-border);
  margin-bottom: 16px;
}

.selected-count {
  color: var(--permission-text-strong);
  font-size: 13px;
}

.permission-layout {
  display: flex;
  height: calc(100vh - 250px);
  border: 1px solid var(--permission-border);
  border-radius: 8px;
  overflow: hidden;
  background: var(--permission-surface);
}

.permission-sidebar {
  width: 240px;
  border-right: 1px solid var(--permission-border);
  display: flex;
  flex-direction: column;
  color: var(--permission-text-default);
}

.permission-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: var(--permission-surface-soft);
  color: var(--permission-text-default);
}

.permission-layout--full .permission-content {
  padding-top: 20px;
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
    border-bottom: 1px solid var(--permission-border);
  }
}
</style>
