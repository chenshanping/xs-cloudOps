<template>
  <div class="group-list-header">API 分组</div>
  <div class="group-list-content">
    <div
      v-for="group in groups"
      :key="group.id"
      :class="['group-item', { active: selectedGroupId === group.id }]"
      @click="$emit('select', group.id)"
    >
      <span class="group-item-name">{{ group.label }}</span>
      <a-badge :count="group.count" :overflow-count="999" class="group-item-badge" />
    </div>
  </div>
</template>

<script setup lang="ts">
interface ApiGroupSidebarItem {
  id: number
  label: string
  count: number
}

defineProps<{
  groups: ApiGroupSidebarItem[]
  selectedGroupId: number | null
}>()

defineEmits<{
  select: [id: number]
}>()
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
}

.group-item-badge {
  margin-inline-start: auto;
}

.group-item-badge :deep(.ant-badge-count) {
  background: var(--permission-code-bg);
  color: var(--permission-text-secondary);
  box-shadow: none;
}

.group-item.active .group-item-badge :deep(.ant-badge-count) {
  background: var(--app-primary-color, #1677ff);
  color: #fff;
}
</style>
