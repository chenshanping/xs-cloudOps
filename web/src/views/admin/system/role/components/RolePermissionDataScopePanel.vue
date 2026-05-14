<template>
  <div :class="['permission-shell', { 'permission-shell--dark': uiStore.isDark }]">
    <div class="permission-header">
      <span class="selected-count">资源 {{ modelValue.length }} 个</span>
      <span class="selected-count">已自定义 {{ customizedCount }} 个</span>
      <span class="selected-count default-scope">
        默认范围
        <a-tag color="blue" class="default-scope-tag">{{ defaultDataScopeLabel }}</a-tag>
      </span>
    </div>

    <div class="permission-layout">
      <div class="permission-sidebar">
        <div class="group-list-header">业务资源</div>
        <div class="group-list-content">
          <div
            v-for="item in modelValue"
            :key="item.resource_code"
            :class="['group-item', { active: selectedResourceCode === item.resource_code }]"
            @click="selectedResourceCode = item.resource_code"
          >
            <span class="group-item-name">
              {{ getResourceDefinition(item.resource_code)?.label || item.resource_code }}
            </span>
            <a-tag
              :color="getScopeTagColor(item.data_scope)"
              class="group-item-tag"
            >
              {{ formatScopeLabel(item.data_scope) }}
            </a-tag>
          </div>
        </div>
      </div>

      <div class="permission-content">
        <template v-if="activeItem">
          <div class="resource-meta">
            <div class="resource-title">
              {{ activeResource?.label || activeItem.resource_code }}
            </div>
            <div v-if="activeResource?.description" class="resource-description">
              {{ activeResource.description }}
            </div>
            <div class="resource-meta-tags">
              <a-tag>当前：{{ formatScopeLabel(activeItem.data_scope) }}</a-tag>
              <a-tag v-if="activeItem.data_scope === 0" color="default">
                回退到默认 {{ defaultDataScopeLabel }}
              </a-tag>
            </div>
          </div>

          <a-form layout="vertical" class="resource-form">
            <a-form-item label="数据范围">
              <a-select
                :value="activeItem.data_scope"
                :options="getScopeOptions(activeItem.resource_code)"
                @change="(value: unknown) => handleScopeChange(activeItem!.resource_code, Number(value))"
              />
            </a-form-item>

            <a-form-item
              v-if="activeItem.data_scope === 2"
              label="自定义部门"
              :validate-status="activeItem.dept_ids.length === 0 ? 'error' : undefined"
              :help="activeItem.dept_ids.length === 0 ? getDeptValidationMessage(activeItem.resource_code) : undefined"
            >
              <a-tree-select
                :value="activeItem.dept_ids"
                :tree-data="deptOptions"
                placeholder="请选择部门"
                tree-default-expand-all
                tree-checkable
                multiple
                style="width: 100%"
                @update:value="(value: number[]) => handleDeptChange(activeItem!.resource_code, value)"
              />
            </a-form-item>
          </a-form>
        </template>
        <a-empty v-else class="content-empty" description="请选择左侧业务资源" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
  findDataScopeResourceDefinition,
  formatDataScopeLabel,
  getSupportedFeatureScopeOptions,
  type DataScopeResourceDefinition,
  type RoleFeatureDataScopeFormItem,
  type RolePermissionDeptOption
} from './dataScopeResources'
import { useUiStore } from '@/store/ui'

const props = defineProps<{
  modelValue: RoleFeatureDataScopeFormItem[]
  resourceDefinitions: DataScopeResourceDefinition[]
  deptOptions: RolePermissionDeptOption[]
  defaultDataScopeLabel: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: RoleFeatureDataScopeFormItem[]): void
}>()

const uiStore = useUiStore()
const selectedResourceCode = ref<string>('')

watch(
  () => props.modelValue,
  items => {
    if (!items.length) {
      selectedResourceCode.value = ''
      return
    }
    if (!items.some(item => item.resource_code === selectedResourceCode.value)) {
      selectedResourceCode.value = items[0].resource_code
    }
  },
  { immediate: true, deep: false }
)

const activeItem = computed(() =>
  props.modelValue.find(item => item.resource_code === selectedResourceCode.value) || null
)

const activeResource = computed(() =>
  activeItem.value ? findDataScopeResourceDefinition(props.resourceDefinitions, activeItem.value.resource_code) : null
)

const customizedCount = computed(() =>
  props.modelValue.filter(item => item.data_scope > 0).length
)

const formatScopeLabel = (value: number) =>
  value === 0 ? '继承默认' : formatDataScopeLabel(value)

const getScopeTagColor = (value: number) => {
  switch (value) {
    case 0:
      return 'default'
    case 1:
      return 'blue'
    case 2:
      return 'geekblue'
    case 3:
    case 4:
      return 'cyan'
    case 5:
      return 'orange'
    default:
      return 'default'
  }
}

const getResourceDefinition = (resourceCode: string) =>
  findDataScopeResourceDefinition(props.resourceDefinitions, resourceCode)

const getScopeOptions = (resourceCode: string) =>
  getSupportedFeatureScopeOptions(getResourceDefinition(resourceCode))

const getDeptValidationMessage = (resourceCode: string) => {
  const resourceLabel = getResourceDefinition(resourceCode)?.label || resourceCode
  return `请选择「${resourceLabel}」的自定义部门`
}

const updateItems = (
  resourceCode: string,
  updater: (item: RoleFeatureDataScopeFormItem) => RoleFeatureDataScopeFormItem
) => {
  emit('update:modelValue', props.modelValue.map(item => (
    item.resource_code === resourceCode ? updater(item) : item
  )))
}

const handleScopeChange = (resourceCode: string, value: number) => {
  updateItems(resourceCode, item => ({
    ...item,
    data_scope: value,
    dept_ids: value === 2 ? item.dept_ids : []
  }))
}

const handleDeptChange = (resourceCode: string, value: number[]) => {
  updateItems(resourceCode, item => ({
    ...item,
    dept_ids: value || []
  }))
}
</script>

<style scoped>
@import './rolePermissionDrawerShared.css';

.default-scope {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.default-scope-tag {
  margin-inline-end: 0;
}

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

.group-item-tag {
  margin-inline-end: 0;
  font-size: 12px;
}

.resource-meta {
  padding: 4px 0 4px;
  border-bottom: 1px solid var(--permission-border-soft);
}

.resource-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--permission-text-default);
  margin-bottom: 6px;
}

.resource-description {
  color: var(--permission-text-secondary);
  font-size: 13px;
  line-height: 1.6;
  margin-bottom: 10px;
}

.resource-meta-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding-bottom: 12px;
}

.resource-form {
  padding-top: 4px;
  max-width: 560px;
}
</style>
