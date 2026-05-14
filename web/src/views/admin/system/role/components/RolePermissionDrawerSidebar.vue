<template>
  <div class="group-list-header">{{ headerTitle }}</div>
  <div class="group-list-content">
    <div
      v-for="menu in topMenus"
      :key="menu.id"
      :class="['group-item', { active: selectedTopMenuId === menu.id }]"
      @click="$emit('select', menu.id)"
    >
      <span class="group-item-name">{{ menu.name }}</span>
      <a-badge
        :count="getSelectedCount(menu)"
        :number-style="{ backgroundColor: '#52c41a' }"
        :show-zero="false"
      />
      <span class="group-item-total">({{ getTotalCount(menu) }})</span>
    </div>

    <!-- 虚拟入口：未分类 API（所有顶菜单之外的全局入口，只出现一次） -->
    <div
      v-if="uncategorizedCount > 0"
      :class="['group-item', 'group-item--system', { active: selectedTopMenuId === uncategorizedMenuId }]"
      @click="$emit('select', uncategorizedMenuId)"
    >
      <span class="group-item-name">未分类 API</span>
      <a-tag color="default" class="group-item-tag">{{ uncategorizedCount }}</a-tag>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Menu } from '@/types'
import { collectAssignableMenuIds } from './permissionDrawer'

interface Props {
  topMenus: Menu[]
  selectedTopMenuId: number | null
  checkedMenuKeys: number[]
  uncategorizedCount: number
  uncategorizedMenuId: number
  headerTitle?: string
}

const props = defineProps<Props>()

defineEmits<{
  select: [id: number]
}>()

const getSelectedCount = (menu: Menu) => {
  const ids = collectAssignableMenuIds(menu)
  return ids.filter(id => props.checkedMenuKeys.includes(id)).length
}

const getTotalCount = (menu: Menu) => collectAssignableMenuIds(menu).length
</script>

<style scoped>
.group-list-header {
  padding: 12px 16px;
  font-weight: 500;
  background: var(--permission-surface-soft);
  border-bottom: 1px solid var(--permission-border);
  color: var(--permission-text-strong);
}

.group-list-content {
  flex: 1;
  overflow-y: auto;
}

.group-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  border-bottom: 1px solid var(--permission-border-soft);
  transition: background 0.2s;
  color: var(--permission-text-strong);
}

.group-item:hover {
  background: var(--permission-hover);
  color: var(--permission-text-default);
}

.group-item.active {
  background: var(--app-primary-color-soft, rgba(0, 107, 230, 0.12));
  border-left: 3px solid var(--app-primary-color, #006be6);
  color: var(--permission-text-default);
}

.group-item-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: inherit;
}

.group-item-total {
  color: var(--permission-text-muted);
  font-size: 12px;
}

.group-item--system {
  margin-top: 8px;
  border-top: 1px dashed var(--permission-border);
  background: var(--permission-surface-soft);
}

.group-item-tag {
  margin: 0;
  font-size: 11px;
  line-height: 16px;
  padding: 0 6px;
}
</style>
