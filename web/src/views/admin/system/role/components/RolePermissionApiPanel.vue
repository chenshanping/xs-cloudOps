<template>
  <div class="panel">
    <div class="panel-header">
      <span>相关 API 权限</span>
      <a-checkbox
        v-if="section.apis.length"
        :checked="checked"
        :indeterminate="indeterminate"
        @change="emit('toggle-all', section, $event.target.checked)"
      >
        全选当前页面 API
      </a-checkbox>
    </div>
    <div class="panel-body">
      <template v-if="visibleApis.length">
        <div
          v-for="api in visibleApis"
          :key="api.id"
          class="api-row"
        >
          <a-checkbox
            :checked="checkedApiIds.includes(api.id)"
            @change="emit('toggle-api', api.id, $event.target.checked)"
          >
            <div class="api-content">
              <a-tag :color="getMethodColor(api.method)" size="small">{{ api.method }}</a-tag>
              <span class="api-group">{{ api.group || '未分组' }}</span>
              <span class="api-path">{{ api.path }}</span>
              <span class="api-desc">{{ api.description }}</span>
            </div>
          </a-checkbox>
        </div>
      </template>
      <a-empty v-else description="暂无匹配 API" />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Api } from '@/types'
import type { PermissionPageSection } from './permissionDrawer'

interface Props {
  section: PermissionPageSection
  visibleApis: Api[]
  checkedApiIds: number[]
  checked: boolean
  indeterminate: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  'toggle-api': [apiId: number, checked: boolean]
  'toggle-all': [section: PermissionPageSection, checked: boolean]
}>()

const getMethodColor = (method: string) => {
  const colors: Record<string, string> = {
    GET: 'green',
    POST: 'blue',
    PUT: 'orange',
    DELETE: 'red',
    PATCH: 'purple'
  }
  return colors[method.toUpperCase()] || 'default'
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
  border-bottom: 1px solid #f5f5f5;
}

.panel-body {
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.api-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px;
  border-radius: 6px;
  transition: background 0.2s;
}

.api-row:hover {
  background: #f5f5f5;
}

.api-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.api-group {
  color: #262626;
  font-size: 12px;
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 10px;
}

.api-path {
  font-size: 13px;
  color: #666;
  font-family: monospace;
}

.api-desc {
  font-size: 13px;
  color: #333;
}

.panel-body :deep(.ant-tag) {
  margin-right: 0;
  font-size: 10px;
  padding: 0 4px;
  line-height: 16px;
  min-width: 50px;
  text-align: center;
}
</style>
