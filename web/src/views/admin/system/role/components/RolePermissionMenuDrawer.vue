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
        title="菜单权限"
        :role-name="roleName"
      />
    </template>

    <div :class="['permission-shell', { 'permission-shell--dark': uiStore.isDark }]">
      <div class="permission-header">
        <span class="selected-count">已选菜单/按钮 {{ assignableSelectedMenuKeys.length }} 个</span>
        <a-input-search
          v-model:value="searchText"
          placeholder="搜索菜单名称、权限码或路由"
          :disabled="permissionLoading"
          style="width: 280px; margin-left: auto"
          allow-clear
        />
      </div>

      <a-spin :spinning="permissionLoading">
        <div class="permission-layout">
          <div class="permission-sidebar">
            <RolePermissionDrawerSidebar
              header-title="一级菜单"
              :top-menus="topMenus"
              :selected-top-menu-id="selectedTopMenuId"
              :checked-menu-keys="assignableSelectedMenuKeys"
              :uncategorized-count="0"
              :uncategorized-menu-id="-1"
              @select="selectedTopMenuId = $event"
            />
          </div>

          <div class="permission-content">
            <a-empty
              v-if="!filteredSections.length"
              class="content-empty"
              description="当前一级菜单下暂无可分配的菜单/按钮"
            />
            <template v-else>
              <a-tabs
                v-if="filteredSections.length > 1"
                v-model:activeKey="activeSectionId"
                size="small"
                class="section-tabs"
              >
                <a-tab-pane
                  v-for="section in filteredSections"
                  :key="section.id"
                  :tab="section.raw.pageMenu.name"
                >
                  <RolePermissionPageSection
                    mode="menu"
                    :section="section.raw"
                    :visible-menu-items="section.visibleMenuItems"
                    :visible-apis="[]"
                    :exact-checked-menu-keys="assignableSelectedMenuKeys"
                    :checked-menu-keys="checkedMenuKeys"
                    :checked-api-ids="[]"
                    :inherited-api-ids="[]"
                    :inherited-api-source-map="{}"
                    :show-inherited-api-meta="false"
                    @toggle-menu="handleMenuToggle"
                    @toggle-section-menus="handleSectionMenusToggle"
                    @keep-page-only="handleSectionKeepPageOnly"
                    @select-child-permissions="handleSectionSelectChildPermissions"
                    @clear-child-permissions="handleSectionClearChildPermissions"
                  />
                </a-tab-pane>
              </a-tabs>
              <RolePermissionPageSection
                v-else
                mode="menu"
                :section="filteredSections[0].raw"
                :visible-menu-items="filteredSections[0].visibleMenuItems"
                :visible-apis="[]"
                :exact-checked-menu-keys="assignableSelectedMenuKeys"
                :checked-menu-keys="checkedMenuKeys"
                :checked-api-ids="[]"
                :inherited-api-ids="[]"
                :inherited-api-source-map="{}"
                :show-inherited-api-meta="false"
                @toggle-menu="handleMenuToggle"
                @toggle-section-menus="handleSectionMenusToggle"
                @keep-page-only="handleSectionKeepPageOnly"
                @select-child-permissions="handleSectionSelectChildPermissions"
                @clear-child-permissions="handleSectionClearChildPermissions"
              />
            </template>
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
import { ref, watch } from 'vue'
import RoleDrawerContextHeader from './RoleDrawerContextHeader.vue'
import RolePermissionDrawerSidebar from './RolePermissionDrawerSidebar.vue'
import RolePermissionPageSection from './RolePermissionPageSection.vue'
import { useResponsiveDrawerWidth } from './useResponsiveDrawerWidth'
import { useRolePermissionMenuDrawer } from './useRolePermissionMenuDrawer'
import { useUiStore } from '@/store/ui'

interface Props {
  roleId: number
  roleName: string
}

const props = defineProps<Props>()
const uiStore = useUiStore()
const visible = defineModel<boolean>('open', { default: false })
const activeSectionId = ref('')
const { drawerWidth } = useResponsiveDrawerWidth(1100, 920)

const {
  assignableSelectedMenuKeys,
  checkedMenuKeys,
  filteredSections,
  handleMenuToggle,
  handleSavePermissions,
  handleSectionClearChildPermissions,
  handleSectionKeepPageOnly,
  handleSectionMenusToggle,
  handleSectionSelectChildPermissions,
  permissionLoading,
  saveLoading,
  searchText,
  selectedTopMenuId,
  topMenus
} = useRolePermissionMenuDrawer(props, visible)

watch([filteredSections, activeSectionId], ([sections]) => {
  if (!sections.length) {
    activeSectionId.value = ''
    return
  }
  if (!sections.some(section => section.id === activeSectionId.value)) {
    activeSectionId.value = sections[0].id
  }
}, { immediate: true })

watch(selectedTopMenuId, () => {
  activeSectionId.value = filteredSections.value[0]?.id || ''
})
</script>

<style scoped>
@import './rolePermissionDrawerShared.css';
</style>
