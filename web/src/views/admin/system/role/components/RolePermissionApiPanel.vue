<template>
  <div class="panel">
    <div class="panel-header">
      <span>高级 API 直授</span>
      <a-checkbox
        v-if="section.apis.length"
        :checked="checked"
        :indeterminate="indeterminate"
        @change="emit('toggle-all', section, $event.target.checked)"
      >
        全选当前分组直授 API
      </a-checkbox>
    </div>
    <div class="panel-hint">
      已显示「继承」的 API 会随菜单/按钮自动生效，只有未绑定或临时兜底接口需要在这里直接授权。
    </div>
    <div class="panel-body">
      <template v-if="apiGroups.length">
        <div
          v-for="group in apiGroups"
          :key="group.id"
          :class="['api-group', { 'api-group--heuristic': group.type === 'heuristic' }]"
        >
          <div class="api-group__header">
            <a-checkbox
              :checked="isGroupChecked(group)"
              :indeterminate="isGroupIndeterminate(group)"
              @change="onGroupToggle(group, $event.target.checked)"
            >
              <span class="api-group__title">
                <a-tag
                  :color="group.type === 'heuristic' ? 'warning' : 'cyan'"
                  class="api-group__type"
                >
                  {{ group.type === 'heuristic' ? '推测' : '菜单' }}
                </a-tag>
                <span class="api-group__label">{{ group.label }}</span>
              </span>
            </a-checkbox>
            <span class="api-group__count">{{ groupCheckedCount(group) }}/{{ group.apis.length }}</span>
          </div>
          <div class="api-group__body">
            <div
              v-for="api in group.apis"
              :key="api.id"
              class="api-row"
            >
              <a-checkbox
                :checked="checkedApiIds.includes(api.id)"
                @change="emit('toggle-api', api.id, $event.target.checked)"
              >
                <div class="api-content">
                  <a-tag :color="getMethodColor(api.method)" size="small">{{ api.method }}</a-tag>
                  <span class="api-path">{{ api.path }}</span>
                  <span class="api-desc">{{ api.description }}</span>
                  <a-tag v-if="checkedApiIds.includes(api.id)" color="blue">直接</a-tag>
                  <a-tag v-if="inheritedApiIds.includes(api.id)" color="green">继承</a-tag>
                  <a-tooltip
                    v-if="inheritedApiIds.includes(api.id) && inheritedApiSourceMap[api.id]?.length"
                    :title="`继承自菜单：${inheritedApiSourceMap[api.id].join('、')}`"
                  >
                    <InfoCircleOutlined class="api-source-icon" />
                  </a-tooltip>
                </div>
              </a-checkbox>
            </div>
          </div>
        </div>
      </template>
      <a-empty v-else description="暂无匹配 API" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { InfoCircleOutlined } from '@ant-design/icons-vue'
import type { Api } from '@/types'
import type { PermissionPageSection } from './permissionDrawer'

interface Props {
  section: PermissionPageSection
  visibleApis: Api[]
  checkedApiIds: number[]
  inheritedApiIds: number[]
  inheritedApiSourceMap: Record<number, string[]>
  checked: boolean
  indeterminate: boolean
}

interface ApiGroup {
  id: string
  label: string
  type: 'menu' | 'heuristic'
  apis: Api[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'toggle-api': [apiId: number, checked: boolean]
  'toggle-all': [section: PermissionPageSection, checked: boolean]
}>()

const apiGroups = computed<ApiGroup[]>(() => {
  const menuOrder: string[] = []
  const menuBuckets = new Map<string, Api[]>()
  const heuristicApis: Api[] = []

  // 以 section.menuItems 的层级顺序预创建分组，让 API 展示顺序与菜单一致
  props.section.menuItems.forEach(item => {
    const name = item.menu.name
    if (!menuBuckets.has(name)) {
      menuBuckets.set(name, [])
      menuOrder.push(name)
    }
  })

  props.visibleApis.forEach(api => {
    const source = props.section.apiSourceMap[api.id]
    if (source === 'heuristic') {
      heuristicApis.push(api)
      return
    }
    const menuNames = props.section.apiMenuSourceMap[api.id] || []
    // 优先选择 menuOrder 里出现过的菜单名作为归属分组
    const primary = menuNames.find(n => menuBuckets.has(n))
      || menuNames[0]
      || props.section.pageMenu.name
    if (!menuBuckets.has(primary)) {
      menuBuckets.set(primary, [])
      menuOrder.push(primary)
    }
    menuBuckets.get(primary)!.push(api)
  })

  const result: ApiGroup[] = []
  menuOrder.forEach(name => {
    const apis = menuBuckets.get(name) || []
    if (apis.length === 0) return
    result.push({
      id: `menu-${name}`,
      label: name,
      type: 'menu',
      apis
    })
  })
  if (heuristicApis.length) {
    result.push({
      id: 'heuristic',
      label: '按名称/路径推测归类',
      type: 'heuristic',
      apis: heuristicApis
    })
  }
  return result
})

const groupCheckedCount = (group: ApiGroup) =>
  group.apis.filter(api => props.checkedApiIds.includes(api.id)).length

const isGroupChecked = (group: ApiGroup) =>
  group.apis.length > 0 && groupCheckedCount(group) === group.apis.length

const isGroupIndeterminate = (group: ApiGroup) => {
  const checked = groupCheckedCount(group)
  return checked > 0 && checked < group.apis.length
}

const onGroupToggle = (group: ApiGroup, checked: boolean) => {
  group.apis.forEach(api => {
    const already = props.checkedApiIds.includes(api.id)
    if (checked && !already) {
      emit('toggle-api', api.id, true)
    } else if (!checked && already) {
      emit('toggle-api', api.id, false)
    }
  })
}

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
  border-bottom: 1px solid var(--permission-border-soft);
}

.panel-hint {
  padding: 8px 16px;
  color: var(--permission-text-muted);
  font-size: 12px;
  background: var(--permission-surface-soft);
  border-bottom: 1px solid var(--permission-border-soft);
}

.panel-body {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.api-group {
  border: 1px solid var(--permission-border-soft);
  border-radius: 8px;
  background: var(--permission-surface);
  overflow: hidden;
}

.api-group--heuristic {
  border-style: dashed;
  border-color: var(--app-warning-color, #faad14);
}

.api-group__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  background: var(--permission-surface-soft);
  border-bottom: 1px solid var(--permission-border-soft);
}

.api-group__title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.api-group__type {
  margin: 0;
  font-size: 11px;
  line-height: 16px;
  padding: 0 6px;
  min-width: 40px;
  text-align: center;
}

.api-group__label {
  font-weight: 500;
  color: var(--permission-text-strong);
}

.api-group__count {
  font-size: 12px;
  color: var(--permission-text-muted);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.api-group__body {
  padding: 6px 12px 8px 28px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.api-row {
  display: flex;
  align-items: center;
  padding: 4px 6px;
  border-radius: 4px;
  transition: background 0.15s;
}

.api-row:hover {
  background: var(--permission-hover);
}

.api-content {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
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

.api-source-icon {
  font-size: 13px;
  color: var(--permission-text-muted);
  cursor: help;
}

.panel-body :deep(.api-row .ant-tag) {
  margin-right: 0;
  font-size: 10px;
  padding: 0 4px;
  line-height: 16px;
  min-width: 40px;
  text-align: center;
}
</style>
