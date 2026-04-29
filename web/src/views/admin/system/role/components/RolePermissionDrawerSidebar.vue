<template>
  <div class="group-list-header">一级菜单</div>
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
  </div>
</template>

<script setup lang="ts">
import type { Menu } from '@/types'
import { collectAssignableMenuIds } from './permissionDrawer'

interface Props {
  topMenus: Menu[]
  selectedTopMenuId: number | null
  checkedMenuKeys: number[]
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
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
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
  border-bottom: 1px solid #f5f5f5;
  transition: background 0.2s;
}

.group-item:hover {
  background: #f5f5f5;
}

.group-item.active {
  background: var(--app-primary-color-soft, rgba(0, 107, 230, 0.12));
  border-left: 3px solid var(--app-primary-color, #006be6);
}

.group-item-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.group-item-total {
  color: #999;
  font-size: 12px;
}
</style>
