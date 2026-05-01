<template>
  <div class="uncategorized-card">
    <div class="uncategorized-header">
      <span>{{ title }}</span>
      <a-checkbox
        :checked="checked"
        :indeterminate="indeterminate"
        @change="emit('toggle-all', $event.target.checked)"
      >
        全选当前分组直授 API
      </a-checkbox>
    </div>
    <div class="uncategorized-body">
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
              <a-tag v-if="checkedApiIds.includes(api.id)" color="blue">直接授权</a-tag>
              <a-tag v-if="inheritedApiIds.includes(api.id)" color="green">菜单继承</a-tag>
              <span
                v-if="inheritedApiIds.includes(api.id) && inheritedApiSourceMap[api.id]?.length"
                class="api-source"
              >
                来源：{{ inheritedApiSourceMap[api.id].join('、') }}
              </span>
            </div>
          </a-checkbox>
        </div>
      </template>
      <a-empty v-else description="暂无匹配的未归类 API" />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Api } from '@/types'

interface Props {
  checked: boolean
  checkedApiIds: number[]
  inheritedApiIds: number[]
  inheritedApiSourceMap: Record<number, string[]>
  indeterminate: boolean
  title?: string
  visibleApis: Api[]
}

defineProps<Props>()

const emit = defineEmits<{
  'toggle-all': [checked: boolean]
  'toggle-api': [apiId: number, checked: boolean]
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
.uncategorized-card {
  border: 1px solid var(--permission-border);
  border-radius: 8px;
  overflow: hidden;
  background: var(--permission-surface);
}

.uncategorized-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 16px;
  font-weight: 500;
  background: var(--permission-surface-soft);
  border-bottom: 1px solid var(--permission-border);
}

.uncategorized-body {
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
  background: var(--permission-hover);
}

.api-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.api-group {
  color: var(--permission-text-default);
  font-size: 12px;
  background: var(--permission-code-bg);
  padding: 2px 6px;
  border-radius: 10px;
}

.api-path {
  font-size: 13px;
  color: var(--permission-text-secondary);
  font-family: monospace;
}

.api-desc {
  font-size: 13px;
  color: var(--permission-text-strong);
}

.api-source {
  font-size: 12px;
  color: var(--permission-text-muted);
}

.uncategorized-body :deep(.ant-tag) {
  margin-right: 0;
  font-size: 10px;
  padding: 0 4px;
  line-height: 16px;
  min-width: 50px;
  text-align: center;
}
</style>
