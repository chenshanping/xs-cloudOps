<template>
  <div class="data-scope-panel">
    <a-alert
      type="info"
      show-icon
      :message="`角色默认数据范围：${defaultDataScopeLabel}`"
      description="这里只配置首批业务功能的覆盖范围。未单独配置时，将自动回退到角色默认数据范围。"
    />

    <div class="data-scope-cards">
      <a-card
        v-for="item in modelValue"
        :key="item.resource_code"
        size="small"
        class="data-scope-card"
      >
        <template #title>
          <div class="card-title">
            <span>{{ getResourceDefinition(item.resource_code)?.label || item.resource_code }}</span>
            <a-tag v-if="item.data_scope === 0">继承默认</a-tag>
          </div>
        </template>

        <p class="card-description">{{ getResourceDefinition(item.resource_code)?.description || '' }}</p>

        <a-form layout="vertical">
          <a-form-item label="数据范围">
            <a-select
              :value="item.data_scope"
              :options="getScopeOptions(item.resource_code)"
              @change="value => handleScopeChange(item.resource_code, value)"
            />
          </a-form-item>

          <a-form-item
            v-if="item.data_scope === 2"
            label="自定义部门"
            :validate-status="item.dept_ids.length === 0 ? 'error' : undefined"
            :help="item.dept_ids.length === 0 ? getDeptValidationMessage(item.resource_code) : undefined"
          >
            <a-tree-select
              :value="item.dept_ids"
              :tree-data="deptOptions"
              placeholder="请选择部门"
              tree-default-expand-all
              tree-checkable
              multiple
              @update:value="value => handleDeptChange(item.resource_code, value)"
            />
          </a-form-item>
        </a-form>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  findDataScopeResourceDefinition,
  getSupportedFeatureScopeOptions,
  type DataScopeResourceDefinition,
  type RoleFeatureDataScopeFormItem,
  type RolePermissionDeptOption
} from './dataScopeResources'

const props = defineProps<{
  modelValue: RoleFeatureDataScopeFormItem[]
  resourceDefinitions: DataScopeResourceDefinition[]
  deptOptions: RolePermissionDeptOption[]
  defaultDataScopeLabel: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: RoleFeatureDataScopeFormItem[]): void
}>()

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
.data-scope-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: var(--permission-text-default);
}

.data-scope-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 16px;
}

.data-scope-card {
  border-radius: 10px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-description {
  margin-bottom: 16px;
  color: var(--permission-text-secondary);
  font-size: 13px;
}
</style>
