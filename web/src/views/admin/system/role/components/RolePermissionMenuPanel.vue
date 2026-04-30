<template>
  <div class="panel">
    <div class="panel-header">
      <span>菜单 / 按钮权限</span>
      <div class="panel-actions">
        <a-space size="small" wrap>
          <a-button size="small" @click="emit('keep-page-only', section)">仅页面</a-button>
          <a-button size="small" @click="emit('select-child-permissions', section)">全选子权限</a-button>
          <a-button size="small" @click="emit('clear-child-permissions', section)">清空子权限</a-button>
        </a-space>
        <a-checkbox
          v-if="section.menuItems.length"
          :checked="checked"
          :indeterminate="indeterminate"
          @change="emit('toggle-all', section, $event.target.checked)"
        >
          全选当前页面菜单
        </a-checkbox>
      </div>
    </div>
    <div class="panel-body">
      <template v-if="visibleMenuItems.length">
        <div
          v-for="item in visibleMenuItems"
          :key="item.menu.id"
          class="menu-row"
          :style="{ paddingLeft: `${item.level * 20 + 8}px` }"
        >
          <a-checkbox
            :checked="isMenuChecked(item.menu)"
            :indeterminate="isMenuIndeterminate(item.menu)"
            @change="emit('toggle-menu', item.menu, $event.target.checked)"
          >
            <span class="menu-label">
              <a-tag :color="getMenuTagColor(item.menu.type)" size="small">
                {{ getMenuTypeText(item.menu.type) }}
              </a-tag>
              <span>{{ item.menu.name }}</span>
            </span>
          </a-checkbox>
          <span v-if="item.menu.permission" class="permission-code">
            {{ item.menu.permission }}
          </span>
        </div>
      </template>
      <a-empty v-else description="暂无匹配菜单权限" />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Menu } from '@/types'
import type { PermissionMenuItem, PermissionPageSection } from './permissionDrawer'
import { collectAssignableMenuIds } from './permissionDrawer'

interface Props {
  section: PermissionPageSection
  visibleMenuItems: PermissionMenuItem[]
  exactCheckedMenuKeys: number[]
  checkedMenuKeys: number[]
  checked: boolean
  indeterminate: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'toggle-menu': [menu: Menu, checked: boolean]
  'toggle-all': [section: PermissionPageSection, checked: boolean]
  'keep-page-only': [section: PermissionPageSection]
  'select-child-permissions': [section: PermissionPageSection]
  'clear-child-permissions': [section: PermissionPageSection]
}>()

const isMenuChecked = (menu: Menu) => {
  if (menu.type === 1) {
    const ids = collectAssignableMenuIds(menu)
    return ids.length > 0 && ids.every(id => props.exactCheckedMenuKeys.includes(id))
  }
  return props.exactCheckedMenuKeys.includes(menu.id)
}

const isMenuIndeterminate = (menu: Menu) => {
  const childIds = collectAssignableMenuIds(menu).filter(id => id !== menu.id)
  if (!childIds.length) {
    return false
  }
  const checkedCount = childIds.filter(id => props.exactCheckedMenuKeys.includes(id)).length
  if (menu.type === 1) {
    return checkedCount > 0 && checkedCount < childIds.length
  }
  if (props.exactCheckedMenuKeys.includes(menu.id)) {
    return false
  }
  return checkedCount > 0
}

const getMenuTypeText = (type: number) => {
  switch (type) {
    case 1:
      return '目录'
    case 2:
      return '菜单'
    case 3:
      return '按钮'
    default:
      return '权限'
  }
}

const getMenuTagColor = (type: number) => {
  switch (type) {
    case 1:
      return 'purple'
    case 2:
      return 'blue'
    case 3:
      return 'orange'
    default:
      return 'default'
  }
}
</script>

<style scoped>
.panel {
  min-width: 0;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  font-weight: 500;
  border-bottom: 1px solid var(--permission-border-soft);
}

.panel-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  flex-wrap: wrap;
}

.panel-body {
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.menu-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px;
  border-radius: 6px;
  transition: background 0.2s;
}

.menu-row:hover {
  background: var(--permission-hover);
}

.menu-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.permission-code {
  font-size: 12px;
  color: var(--permission-text-muted);
  font-family: monospace;
  background: var(--permission-code-bg);
  padding: 2px 6px;
  border-radius: 3px;
}

.panel-body :deep(.ant-tag) {
  margin-right: 0;
  font-size: 10px;
  padding: 0 4px;
  line-height: 16px;
}
</style>
