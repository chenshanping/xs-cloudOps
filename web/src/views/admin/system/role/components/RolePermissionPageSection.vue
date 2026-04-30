<template>
  <div class="section-card">
    <div class="section-header">
      <div class="section-title">
        <span>{{ section.pageMenu.name }}</span>
        <span v-if="section.pageMenu.permission" class="permission-code">
          {{ section.pageMenu.permission }}
        </span>
      </div>
      <div class="section-counts">
        <span v-if="showMenus">菜单 {{ visibleMenuItems.length }}</span>
        <span v-if="showApis">API {{ visibleApis.length }}</span>
      </div>
    </div>

    <div class="section-grid" :class="gridClass">
      <RolePermissionMenuPanel
        v-if="showMenus"
        :section="section"
        :visible-menu-items="visibleMenuItems"
        :exact-checked-menu-keys="exactCheckedMenuKeys"
        :checked-menu-keys="checkedMenuKeys"
        :checked="isSectionMenuChecked"
        :indeterminate="isSectionMenuIndeterminate"
        @toggle-menu="(menu, checked) => emit('toggle-menu', menu, checked)"
        @toggle-all="(sectionValue, checked) => emit('toggle-section-menus', sectionValue, checked)"
        @keep-page-only="sectionValue => emit('keep-page-only', sectionValue)"
        @select-child-permissions="sectionValue => emit('select-child-permissions', sectionValue)"
        @clear-child-permissions="sectionValue => emit('clear-child-permissions', sectionValue)"
      />

      <RolePermissionApiPanel
        v-if="showApis"
        :section="section"
        :visible-apis="visibleApis"
        :checked-api-ids="checkedApiIds"
        :checked="isSectionApiChecked"
        :indeterminate="isSectionApiIndeterminate"
        @toggle-api="(apiId, checked) => emit('toggle-api', apiId, checked)"
        @toggle-all="(sectionValue, checked) => emit('toggle-section-apis', sectionValue, checked)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Api, Menu } from '@/types'
import type { PermissionMenuItem, PermissionPageSection } from './permissionDrawer'
import RolePermissionMenuPanel from './RolePermissionMenuPanel.vue'
import RolePermissionApiPanel from './RolePermissionApiPanel.vue'

interface Props {
  section: PermissionPageSection
  visibleMenuItems: PermissionMenuItem[]
  visibleApis: Api[]
  exactCheckedMenuKeys: number[]
  checkedMenuKeys: number[]
  checkedApiIds: number[]
  mode?: 'split' | 'menu' | 'api'
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'toggle-menu': [menu: Menu, checked: boolean]
  'toggle-section-menus': [section: PermissionPageSection, checked: boolean]
  'keep-page-only': [section: PermissionPageSection]
  'select-child-permissions': [section: PermissionPageSection]
  'clear-child-permissions': [section: PermissionPageSection]
  'toggle-api': [apiId: number, checked: boolean]
  'toggle-section-apis': [section: PermissionPageSection, checked: boolean]
}>()

const showMenus = computed(() => props.mode !== 'api')
const showApis = computed(() => props.mode !== 'menu')

const getSectionMenuIds = () =>
  props.section.menuItems
    .filter(item => item.menu.type !== 1)
    .map(item => item.menu.id)
const getSectionApiIds = () => props.section.apis.map(api => api.id)

const gridClass = computed(() => {
  if (showMenus.value && showApis.value) {
    return 'section-grid--split'
  }
  return 'section-grid--single'
})

const isSectionMenuChecked = computed(() => {
  const ids = getSectionMenuIds()
  return ids.length > 0 && ids.every(id => props.exactCheckedMenuKeys.includes(id))
})

const isSectionMenuIndeterminate = computed(() => {
  const ids = getSectionMenuIds()
  const checkedCount = ids.filter(id => props.exactCheckedMenuKeys.includes(id)).length
  return checkedCount > 0 && checkedCount < ids.length
})

const isSectionApiChecked = computed(() => {
  const ids = getSectionApiIds()
  return ids.length > 0 && ids.every(id => props.checkedApiIds.includes(id))
})

const isSectionApiIndeterminate = computed(() => {
  const ids = getSectionApiIds()
  const checkedCount = ids.filter(id => props.checkedApiIds.includes(id)).length
  return checkedCount > 0 && checkedCount < ids.length
})
</script>

<style scoped>
.section-card {
  border: 1px solid var(--permission-border);
  border-radius: 8px;
  overflow: hidden;
  background: var(--permission-surface);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 16px;
  background: var(--permission-surface-soft);
  border-bottom: 1px solid var(--permission-border);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  font-weight: 500;
}

.section-counts {
  display: flex;
  gap: 12px;
  color: var(--permission-text-muted);
  font-size: 12px;
  white-space: nowrap;
}

.section-grid {
  display: grid;
}

.section-grid--split {
  grid-template-columns: 1fr 1fr;
}

.section-grid--single {
  grid-template-columns: 1fr;
}

.section-grid--split > * + * {
  border-left: 1px solid var(--permission-border);
}

.permission-code {
  font-size: 12px;
  color: var(--permission-text-muted);
  font-family: monospace;
  background: var(--permission-code-bg);
  padding: 2px 6px;
  border-radius: 3px;
}

@media (max-width: 1200px) {
  .section-grid--split {
    grid-template-columns: 1fr;
  }

  .section-grid--split > * + * {
    border-left: 0;
    border-top: 1px solid var(--permission-border);
  }
}
</style>
