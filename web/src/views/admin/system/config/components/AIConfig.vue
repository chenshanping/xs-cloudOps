<template>
  <div class="ai-config">
    <a-spin :spinning="loading">
      <div class="ai-config-layout">
        <section class="provider-panel">
          <div class="panel-header">
            <div>
              <div class="panel-title">AI 平台</div>
              <div class="panel-subtitle">选择平台并维护本地导入模型</div>
            </div>
            <a-button type="primary" @click="openCreateProviderDrawer">
              <template #icon><PlusOutlined /></template>
              新增平台
            </a-button>
          </div>

          <div v-if="formData.providers.length > 0" class="provider-list">
            <button
              v-for="(provider, index) in formData.providers"
              :key="`${provider.name || 'provider'}-${index}`"
              type="button"
              class="provider-item"
              :class="{ 'provider-item--active': index === selectedProviderIndex }"
              @click="selectProvider(index)"
            >
              <div class="provider-item__header">
                <span class="provider-item__name">{{ provider.name || '未命名平台' }}</span>
                <a-tag v-if="isDefaultProvider(index)" color="blue">默认</a-tag>
              </div>
              <div class="provider-item__meta">{{ provider.models.length }} 个已导入模型</div>
            </button>
          </div>

          <a-empty v-else description="暂未配置 AI 平台">
            <a-button type="primary" @click="openCreateProviderDrawer">新增第一个平台</a-button>
          </a-empty>
        </section>

        <section class="models-panel">
          <template v-if="selectedProvider">
            <div class="panel-header">
              <div>
                <div class="panel-title">{{ selectedProvider.name || '未命名平台' }}</div>
                <div class="panel-subtitle">
                  {{ selectedProvider.base_url || '请先配置 Base URL' }}
                </div>
              </div>
              <a-space wrap>
                <a-button @click="openEditProviderDrawer">
                  <template #icon><EditOutlined /></template>
                  编辑平台
                </a-button>
                <a-button type="primary" @click="openModelManagerDrawer">
                  <template #icon><CloudDownloadOutlined /></template>
                  管理平台模型
                </a-button>
              </a-space>
            </div>

            <a-alert
              class="page-alert"
              type="info"
              show-icon
              message="模型导入只更新当前编辑态；仍需点击底部“保存配置”后才会真正落库。"
            />

            <div class="models-header">
              <div>
                <div class="models-title">已导入模型</div>
                <div class="models-subtitle">支持本地维护显示名称、描述、顺序与测试。</div>
              </div>
              <a-button @click="addModel">
                <template #icon><PlusOutlined /></template>
                手动添加模型
              </a-button>
            </div>

            <a-table
              :data-source="selectedProvider.models"
              :columns="modelColumns"
              :pagination="false"
              :row-key="(_record, index) => `model-${index}`"
              size="small"
              class="models-table"
              :locale="{ emptyText: '当前平台还没有已导入模型' }"
            >
              <template #bodyCell="{ column, index }">
                <template v-if="column.key === 'id'">
                  <a-input
                    v-model:value="selectedProvider.models[index].id"
                    size="small"
                    placeholder="模型 ID"
                  />
                </template>
                <template v-else-if="column.key === 'name'">
                  <a-input
                    v-model:value="selectedProvider.models[index].name"
                    size="small"
                    placeholder="显示名称"
                  />
                </template>
                <template v-else-if="column.key === 'description'">
                  <a-input
                    v-model:value="selectedProvider.models[index].description"
                    size="small"
                    placeholder="模型描述"
                  />
                </template>
                <template v-else-if="column.key === 'actions'">
                  <div class="model-actions">
                    <a-tooltip title="上移">
                      <a-button
                        type="text"
                        size="small"
                        :disabled="index === 0"
                        @click="moveModel(index, -1)"
                      >
                        <template #icon><UpOutlined /></template>
                      </a-button>
                    </a-tooltip>
                    <a-tooltip title="下移">
                      <a-button
                        type="text"
                        size="small"
                        :disabled="index === selectedProvider.models.length - 1"
                        @click="moveModel(index, 1)"
                      >
                        <template #icon><DownOutlined /></template>
                      </a-button>
                    </a-tooltip>
                    <a-tooltip title="测试">
                      <a-button
                        type="text"
                        size="small"
                        :loading="testingModel === `${selectedProviderIndex}-${index}`"
                        @click="testModel(index)"
                      >
                        <template #icon><ThunderboltOutlined /></template>
                      </a-button>
                    </a-tooltip>
                    <a-tooltip title="删除">
                      <a-button type="text" danger size="small" @click="removeModel(index)">
                        <template #icon><DeleteOutlined /></template>
                      </a-button>
                    </a-tooltip>
                  </div>
                </template>
              </template>
            </a-table>
          </template>

          <a-empty v-else description="请先新增 AI 平台">
            <a-button type="primary" @click="openCreateProviderDrawer">新增平台</a-button>
          </a-empty>
        </section>
      </div>

      <div class="page-actions">
        <a-button type="primary" :loading="saving" @click="handleSave">保存配置</a-button>
      </div>
    </a-spin>

    <ProviderEditorDrawer
      v-model:open="providerDrawerOpen"
      :is-edit="providerDrawerMode === 'edit'"
      :initial-value="editingProvider"
      :is-default="editingProviderIsDefault"
      :existing-names="providerNames"
      @submit="handleProviderSubmit"
      @remove="handleProviderRemove"
    />

    <ProviderModelManagerDrawer
      v-model:open="modelManagerDrawerOpen"
      :provider-name="selectedProvider?.name ?? ''"
      :api-key="selectedProvider?.api_key ?? ''"
      :provider-base-url="selectedProvider?.base_url ?? ''"
      :local-models="selectedProvider?.models ?? []"
      @import="handleImportModels"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import {
  CloudDownloadOutlined,
  DeleteOutlined,
  DownOutlined,
  EditOutlined,
  PlusOutlined,
  ThunderboltOutlined,
  UpOutlined,
} from '@ant-design/icons-vue'
import { aiTest, getAIConfig, updateAIConfig } from '@/api/ai'
import { cloneFromSnapshot, createSnapshot, isSnapshotDirty } from '../config-tab-guard'
import ProviderEditorDrawer from './ai-config/ProviderEditorDrawer.vue'
import ProviderModelManagerDrawer from './ai-config/ProviderModelManagerDrawer.vue'
import {
  createEmptyModel,
  createEmptyProvider,
  mergeImportedModels,
  normalizeAIConfig,
  type AIConfigState,
  type RemoteProviderModel,
} from './ai-config/state'

interface ProviderEditorSubmitValue {
  name: string
  api_key: string
  base_url: string
  isDefault: boolean
}

const modelColumns = [
  { title: '模型 ID', key: 'id', width: 220 },
  { title: '显示名称', key: 'name', width: 220 },
  { title: '模型描述', key: 'description' },
  { title: '操作', key: 'actions', width: 180, fixed: 'right' as const },
]

const loading = ref(true)
const saving = ref(false)
const selectedProviderIndex = ref(0)
const providerDrawerOpen = ref(false)
const providerDrawerMode = ref<'create' | 'edit'>('create')
const editingProviderIndex = ref<number | null>(null)
const modelManagerDrawerOpen = ref(false)
const testingModel = ref<string | null>(null)
const initialized = ref(false)

const emit = defineEmits<{
  (e: 'dirty-change', value: boolean): void
}>()

const formData = reactive<AIConfigState>({
  default_provider: '',
  providers: [],
})

const selectedProvider = computed(() => formData.providers[selectedProviderIndex.value] ?? null)
const providerNames = computed(() => formData.providers.map(provider => provider.name))
const editingProvider = computed(() => {
  if (editingProviderIndex.value === null) {
    return undefined
  }
  return formData.providers[editingProviderIndex.value]
})
const editingProviderIsDefault = computed(() => {
  if (editingProviderIndex.value === null) {
    return formData.providers.length === 0
  }
  return isDefaultProvider(editingProviderIndex.value)
})

const getConfigState = (): AIConfigState => ({
  default_provider: formData.default_provider,
  providers: formData.providers.map(provider => ({
    name: provider.name,
    api_key: provider.api_key,
    base_url: provider.base_url,
    models: provider.models.map(model => ({
      id: model.id,
      name: model.name,
      description: model.description,
    })),
  })),
})

const syncSelectedProviderIndex = (preferredIndex = selectedProviderIndex.value) => {
  if (formData.providers.length === 0) {
    selectedProviderIndex.value = 0
    return
  }
  selectedProviderIndex.value = Math.min(Math.max(preferredIndex, 0), formData.providers.length - 1)
}

const applyConfigState = (state?: Partial<AIConfigState> | null) => {
  const normalized = normalizeAIConfig(state)
  formData.default_provider = normalized.default_provider
  formData.providers = normalized.providers
  if (!formData.default_provider && formData.providers.length > 0) {
    formData.default_provider = formData.providers[0].name
  }
  syncSelectedProviderIndex()
}

const baselineSnapshot = ref(createSnapshot(getConfigState()))
const hasUnsavedChanges = computed(() => initialized.value && isSnapshotDirty(baselineSnapshot.value, getConfigState()))

watch(
  hasUnsavedChanges,
  value => {
    emit('dirty-change', value)
  },
  { immediate: true },
)

const loadConfig = async () => {
  loading.value = true
  try {
    const res = await getAIConfig()
    applyConfigState(res.data)
  } catch (error: any) {
    message.error(error.message || '加载 AI 配置失败')
    applyConfigState()
  } finally {
    baselineSnapshot.value = createSnapshot(getConfigState())
    initialized.value = true
    loading.value = false
  }
}

const openCreateProviderDrawer = () => {
  providerDrawerMode.value = 'create'
  editingProviderIndex.value = null
  providerDrawerOpen.value = true
}

const openEditProviderDrawer = () => {
  if (!selectedProvider.value) {
    message.warning('请先选择要编辑的平台')
    return
  }
  providerDrawerMode.value = 'edit'
  editingProviderIndex.value = selectedProviderIndex.value
  providerDrawerOpen.value = true
}

const selectProvider = (index: number) => {
  selectedProviderIndex.value = index
}

const isDefaultProvider = (index: number) => formData.providers[index]?.name === formData.default_provider

const ensureDefaultProvider = () => {
  if (formData.providers.length === 0) {
    formData.default_provider = ''
    return
  }
  const hasDefaultProvider = formData.providers.some(provider => provider.name === formData.default_provider)
  if (!hasDefaultProvider) {
    formData.default_provider = formData.providers[0].name
  }
}

const handleProviderSubmit = (value: ProviderEditorSubmitValue) => {
  if (providerDrawerMode.value === 'create') {
    formData.providers.push({
      ...createEmptyProvider(),
      name: value.name,
      api_key: value.api_key,
      base_url: value.base_url,
    })
    selectedProviderIndex.value = formData.providers.length - 1
    if (value.isDefault || formData.providers.length === 1 || !formData.default_provider) {
      formData.default_provider = value.name
    } else {
      ensureDefaultProvider()
    }
  } else if (editingProviderIndex.value !== null) {
    const provider = formData.providers[editingProviderIndex.value]
    const previousName = provider.name
    provider.name = value.name
    provider.api_key = value.api_key
    provider.base_url = value.base_url

    if (value.isDefault) {
      formData.default_provider = value.name
    } else if (formData.default_provider === previousName) {
      formData.default_provider = value.name
    }
    ensureDefaultProvider()
    selectedProviderIndex.value = editingProviderIndex.value
  }

  providerDrawerOpen.value = false
}

const handleProviderRemove = () => {
  const index = editingProviderIndex.value ?? selectedProviderIndex.value
  if (index < 0 || index >= formData.providers.length) {
    return
  }

  const removed = formData.providers[index]
  formData.providers.splice(index, 1)
  if (formData.default_provider === removed.name) {
    formData.default_provider = formData.providers[0]?.name ?? ''
  }
  providerDrawerOpen.value = false
  ensureDefaultProvider()
  syncSelectedProviderIndex(index)
  message.success(`已删除平台 ${removed.name || '未命名平台'}`)
}

const openModelManagerDrawer = () => {
  if (!selectedProvider.value) {
    message.warning('请先选择平台')
    return
  }
  modelManagerDrawerOpen.value = true
}

const addModel = () => {
  if (!selectedProvider.value) {
    return
  }
  selectedProvider.value.models.push(createEmptyModel())
}

const removeModel = (modelIndex: number) => {
  selectedProvider.value?.models.splice(modelIndex, 1)
}

const moveModel = (modelIndex: number, direction: number) => {
  const models = selectedProvider.value?.models
  if (!models) {
    return
  }
  const targetIndex = modelIndex + direction
  if (targetIndex < 0 || targetIndex >= models.length) {
    return
  }
  const current = models[modelIndex]
  models[modelIndex] = models[targetIndex]
  models[targetIndex] = current
}

const testModel = async (modelIndex: number) => {
  const provider = selectedProvider.value
  const model = provider?.models[modelIndex]
  if (!provider || !model) {
    return
  }

  if (!provider.api_key.trim()) {
    message.warning('请先填写当前平台的 API Key')
    return
  }
  if (!provider.base_url.trim()) {
    message.warning('请先填写当前平台的 Base URL')
    return
  }
  if (!model.id.trim()) {
    message.warning('请先填写模型 ID')
    return
  }

  testingModel.value = `${selectedProviderIndex.value}-${modelIndex}`
  try {
    await aiTest({
      api_key: provider.api_key,
      base_url: provider.base_url,
      model: model.id,
    })
    message.success(`模型 ${model.name || model.id} 测试成功`)
  } catch (error: any) {
    message.error(error.message || '测试失败')
  } finally {
    testingModel.value = null
  }
}

const handleImportModels = (models: RemoteProviderModel[]) => {
  if (!selectedProvider.value) {
    return
  }
  const result = mergeImportedModels(selectedProvider.value.models, models)
  selectedProvider.value.models = result.models

  if (result.importedCount > 0) {
    message.success(`已导入 ${result.importedCount} 个模型，重复跳过 ${result.skippedCount} 个`)
    return
  }
  message.info('所选模型都已存在，未重复导入')
}

const validateBeforeSave = () => {
  const providerNameSet = new Set<string>()

  for (const provider of formData.providers) {
    provider.name = provider.name.trim()
    provider.api_key = provider.api_key.trim()
    provider.base_url = provider.base_url.trim()

    if (!provider.name) {
      message.warning('请填写平台名称')
      return false
    }
    if (providerNameSet.has(provider.name)) {
      message.warning(`平台名称“${provider.name}”不能重复`)
      return false
    }
    providerNameSet.add(provider.name)

    if (!provider.base_url) {
      message.warning(`请填写 ${provider.name} 的 Base URL`)
      return false
    }

    const modelIDSet = new Set<string>()
    for (const model of provider.models) {
      model.id = model.id.trim()
      model.name = (model.name || model.id).trim()
      model.description = model.description.trim()

      if (!model.id) {
        message.warning(`请填写 ${provider.name} 平台下的模型 ID`)
        return false
      }
      if (modelIDSet.has(model.id)) {
        message.warning(`平台 ${provider.name} 下存在重复模型 ID：${model.id}`)
        return false
      }
      modelIDSet.add(model.id)
    }
  }

  ensureDefaultProvider()
  return true
}

const save = async () => {
  if (!validateBeforeSave()) {
    return false
  }

  saving.value = true
  try {
    await updateAIConfig(getConfigState())
    baselineSnapshot.value = createSnapshot(getConfigState())
    message.success('保存成功')
    return true
  } catch {
    message.error('保存失败')
    return false
  } finally {
    saving.value = false
  }
}

const discardChanges = () => {
  const restored = cloneFromSnapshot<AIConfigState>(baselineSnapshot.value)
  applyConfigState(restored)
  closeTransientUi()
}

const closeTransientUi = () => {
  providerDrawerOpen.value = false
  modelManagerDrawerOpen.value = false
}

const handleSave = async () => {
  await save()
}

onMounted(() => {
  loadConfig()
})

defineExpose({
  isDirty: () => hasUnsavedChanges.value,
  save,
  discardChanges,
  closeTransientUi,
})
</script>

<style scoped>
.ai-config {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: var(--app-text-color);
}

.ai-config-layout {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  gap: 16px;
  align-items: start;
}

.provider-panel,
.models-panel {
  border: 1px solid var(--app-border-color);
  border-radius: 10px;
  background: var(--app-surface-color);
  padding: 16px;
}

.panel-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.panel-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--app-text-strong);
}

.panel-subtitle {
  margin-top: 4px;
  color: var(--app-text-secondary);
  word-break: break-all;
}

.provider-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.provider-item {
  width: 100%;
  border: 1px solid var(--app-border-color);
  border-radius: 10px;
  padding: 14px 16px;
  background: var(--app-surface-soft);
  cursor: pointer;
  text-align: left;
  transition: all 0.2s ease;
  color: var(--app-text-color);
}

.provider-item:hover,
.provider-item--active {
  border-color: var(--app-primary-color);
  background: var(--app-primary-color-soft);
}

.provider-item__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.provider-item__name {
  min-width: 0;
  font-weight: 600;
  color: var(--app-text-strong);
}

.provider-item__meta {
  margin-top: 8px;
  color: var(--app-text-secondary);
  font-size: 12px;
}

.page-alert {
  margin-bottom: 16px;
}

.models-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.models-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--app-text-strong);
}

.models-subtitle {
  margin-top: 4px;
  color: var(--app-text-secondary);
}

.models-table :deep(.ant-table-cell) {
  vertical-align: middle;
}

.ai-config :deep(.ant-empty-description) {
  color: var(--app-text-secondary);
}

.ai-config :deep(.ant-alert) {
  border-color: var(--app-border-color);
}

.ai-config :deep(.ant-alert-message),
.ai-config :deep(.ant-alert-description) {
  color: var(--app-text-color);
}

.models-table :deep(.ant-table) {
  background: transparent;
  color: var(--app-text-color);
}

.models-table :deep(.ant-table-container) {
  border: 1px solid var(--app-border-color);
  border-radius: 10px;
  overflow: hidden;
}

.models-table :deep(.ant-table-thead > tr > th) {
  background: var(--app-surface-soft);
  color: var(--app-text-strong);
  border-bottom-color: var(--app-border-color);
}

.models-table :deep(.ant-table-tbody > tr > td) {
  background: var(--app-surface-color);
  border-bottom-color: var(--app-border-color);
  color: var(--app-text-color);
}

.models-table :deep(.ant-table-tbody > tr:hover > td) {
  background: var(--app-hover-bg);
}

.models-table :deep(.ant-table-placeholder > td) {
  background: var(--app-surface-color);
}

.models-table :deep(.ant-input) {
  background: var(--app-surface-soft);
  border-color: var(--app-border-color);
  color: var(--app-text-color);
}

.models-table :deep(.ant-input::placeholder) {
  color: var(--app-text-muted);
}

.model-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
}

.page-actions {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 960px) {
  .ai-config-layout {
    grid-template-columns: 1fr;
  }
}
</style>
