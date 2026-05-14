<template>
  <PageWrapper class="ai-config">
    <a-spin :spinning="loading">
      <AdminSplitLayout class="ai-config-layout" :aside-width="320" :content-min-width="960">
        <template #aside>
          <aside class="provider-sidebar">
            <div class="sidebar-header">
              <div class="sidebar-title">AI 平台</div>
              <a-button type="primary" @click="openCreateProviderDrawer" v-permission="'ai:config:createProvider'">
                <template #icon><PlusOutlined /></template>
                新增平台
              </a-button>
            </div>

            <div v-if="formData.providers.length > 0" class="provider-list">
              <article
                v-for="(provider, index) in formData.providers"
                :key="`${provider.name || 'provider'}-${index}`"
                class="provider-item"
                :class="{ 'provider-item--active': index === selectedProviderIndex }"
              >
                <button
                  type="button"
                  class="provider-item__main"
                  @click="selectProvider(index)"
                >
                  <div class="provider-item__summary">
                    <div class="provider-item__name">{{ provider.name || '未命名平台' }}</div>
                    <span class="provider-item__count">{{ provider.models.length }} 个模型</span>
                  </div>
                </button>

                <div class="provider-item__actions">
                  <a-button
                    type="link"
                    size="small"
                    class="provider-action-button"
                    :loading="providerSaving"
                    :disabled="providerSaving"
                    @click.stop="openEditProviderDrawer(index)"
                    v-permission="'ai:config:editProvider'"
                  >
                    编辑
                  </a-button>
                  <a-button
                    v-if="!isDefaultProvider(index)"
                    type="link"
                    size="small"
                    class="provider-action-button"
                    :loading="providerSaving"
                    :disabled="providerSaving"
                    @click.stop="handleSetDefaultProvider(index)"
                    v-permission="'ai:config:save'"
                  >
                    设为默认
                  </a-button>
                </div>
              </article>
            </div>

            <a-empty v-else description="暂未配置 AI 平台">
              <a-button type="primary" @click="openCreateProviderDrawer" v-permission="'ai:config:createProvider'">新增第一个平台</a-button>
            </a-empty>
          </aside>
        </template>

        <section class="workspace-panel">
          <template v-if="selectedProvider">
            <div class="workspace-hero">
              <div>
                <div class="workspace-hero__title">{{ selectedProvider.name || '未命名平台' }}</div>
                <div class="workspace-hero__subtitle">
                  {{ selectedProvider.base_url || '请先配置 Base URL' }}
                </div>
              </div>
              <a-space wrap>
                <a-button @click="openEditProviderDrawer" v-permission="'ai:config:editProvider'">
                  <template #icon><EditOutlined /></template>
                  编辑平台
                </a-button>
                <a-button @click="openCreateModelDrawer" v-permission="'ai:config:createModel'">
                  <template #icon><PlusOutlined /></template>
                  新增模型
                </a-button>
                <a-button @click="openImportDrawer" v-permission="'ai:config:importModel'">
                  <template #icon><CloudDownloadOutlined /></template>
                  自动识别导入
                </a-button>
                <a-button type="primary" :loading="saving" @click="handleSave" v-permission="'ai:config:save'">
                  <template #icon><SaveOutlined /></template>
                  保存配置
                </a-button>
              </a-space>
            </div>

            <div class="workspace-summary">
              <span class="summary-chip">已导入 {{ selectedProvider.models.length }} 个模型</span>
              <span class="summary-chip">当前视图 {{ filteredProviderModels.length }} 个</span>
              <span class="summary-chip">已选批量 {{ selectedBatchIndices.length }} 个</span>
              <span class="summary-chip">默认平台 {{ isDefaultProvider(selectedProviderIndex) ? '是' : '否' }}</span>
            </div>

            <div class="workspace-toolbar">
              <div class="workspace-toolbar__main">
                <div class="toolbar-title">模型工作区</div>
                <div class="toolbar-meta">
                  <span v-if="activeModel">当前模型：{{ activeModel.name || activeModel.id || '未命名模型' }}</span>
                  <span v-else>请选择一个模型继续编辑、测试或排序</span>
                </div>
              </div>
              <a-space wrap>
                <a-button :disabled="!activeModel" @click="openEditModelDrawer()" v-permission="'ai:config:editModel'">
                  编辑模型
                </a-button>
                <a-button :disabled="!activeModel" :loading="testingModel" @click="handleTestActiveModel" v-permission="'ai:config:test'">
                  测试模型
                </a-button>
                <a-button :disabled="activeModelIndex <= 0" @click="moveActiveModel(-1)" v-permission="'ai:config:editModel'">上移</a-button>
                <a-button
                  :disabled="activeModelIndex < 0 || activeModelIndex >= selectedProvider.models.length - 1"
                  @click="moveActiveModel(1)"
                  v-permission="'ai:config:editModel'"
                >
                  下移
                </a-button>
                <a-popconfirm title="确定删除当前模型吗？" @confirm="handleRemoveActiveModel">
                  <a-button danger :disabled="!activeModel" v-permission="'ai:config:deleteModel'">删除当前</a-button>
                </a-popconfirm>
                <a-popconfirm title="确定删除已选中的模型吗？" @confirm="handleBatchRemoveModels">
                  <a-button danger :disabled="selectedBatchIndices.length === 0" v-permission="'ai:config:deleteModel'">批量删除</a-button>
                </a-popconfirm>
              </a-space>
            </div>

            <div class="workspace-filters">
              <a-tabs v-model:activeKey="activeCapabilityTab" class="workspace-tabs">
                <a-tab-pane
                  v-for="item in capabilityTabOptions"
                  :key="item.value"
                  :tab="`${item.label} (${getCapabilityCount(item.value)})`"
                />
              </a-tabs>
              <a-input-search
                v-model:value="modelKeyword"
                allow-clear
                placeholder="搜索模型 ID / 名称 / 描述"
                class="workspace-search"
              />
            </div>

            <div v-if="filteredProviderModels.length > 0" class="model-grid">
              <article
                v-for="entry in filteredProviderModels"
                :key="`${entry.index}-${entry.model.id || entry.model.name || 'model'}`"
                class="model-card"
                :class="{ 'model-card--active': entry.index === activeModelIndex }"
                @click="selectActiveModel(entry.index)"
              >
                <div class="model-card__head">
                  <div class="model-card__title-wrap">
                    <div class="model-card__name">{{ entry.model.name || entry.model.id || '未命名模型' }}</div>
                    <div class="model-card__id">{{ entry.model.id || '待填写模型 ID' }}</div>
                  </div>
                  <a-checkbox
                    :checked="selectedBatchSet.has(entry.index)"
                    @click.stop
                    @change="toggleBatchModel(entry.index)"
                  />
                </div>

                <div class="model-card__tags">
                  <a-tag
                    v-for="tag in getModelCapabilityTags(entry.model)"
                    :key="tag"
                    :color="capabilityTagMetaMap[tag].color"
                  >
                    {{ capabilityTagMetaMap[tag].label }}
                  </a-tag>
                  <span v-if="getModelCapabilityTags(entry.model).length === 0" class="model-card__placeholder">未识别能力</span>
                </div>

                <div class="model-card__meta">
                  <span>联网 {{ formatSearchStrategyLabel(entry.model.search_strategy) }}</span>
                  <span>温度 {{ formatTemperature(entry.model.temperature) }}</span>
                  <span>上下文 {{ formatContextWindow(entry.model.context_window) }}</span>
                </div>

                <div v-if="entry.model.description" class="model-card__desc">{{ entry.model.description }}</div>

                <div class="model-card__footer">
                  <span>排序 {{ entry.index + 1 }}</span>
                  <a-button type="link" size="small" @click.stop="openEditModelDrawer(entry.index)" v-permission="'ai:config:editModel'">编辑</a-button>
                </div>
              </article>
            </div>

            <a-empty v-else description="当前筛选条件下没有模型">
              <a-button type="primary" @click="openCreateModelDrawer" v-permission="'ai:config:createModel'">新增模型</a-button>
            </a-empty>
          </template>

          <a-empty v-else description="请先新增 AI 平台">
            <a-button type="primary" @click="openCreateProviderDrawer" v-permission="'ai:config:createProvider'">新增平台</a-button>
          </a-empty>
        </section>
      </AdminSplitLayout>
    </a-spin>

    <ProviderEditorDrawer
      v-model:open="providerDrawerOpen"
      :is-edit="providerDrawerMode === 'edit'"
      :initial-value="editingProvider"
      :is-default="editingProviderIsDefault"
      :existing-names="providerNames"
      :submitting="providerSaving"
      :can-delete="hasPermission('ai:config:deleteProvider')"
      :can-submit="hasPermission(providerDrawerMode === 'edit' ? 'ai:config:editProvider' : 'ai:config:createProvider')"
      @submit="handleProviderSubmit"
      @remove="handleProviderRemove"
    />

    <ProviderModelEditorDrawer
      v-model:open="modelEditorOpen"
      :is-edit="modelEditorMode === 'edit'"
      :model="modelEditorInitialValue"
      :provider-name="selectedProvider?.name ?? ''"
      :can-delete="hasPermission('ai:config:deleteModel')"
      :can-submit="hasPermission(modelEditorMode === 'edit' ? 'ai:config:editModel' : 'ai:config:createModel')"
      :can-test="hasPermission('ai:config:test')"
      @submit="handleModelSubmit"
      @remove="handleModelRemoveFromDrawer"
      @test="handleTestModel"
    />

    <ProviderRemoteModelImportDrawer
      v-model:open="importDrawerOpen"
      :provider-name="selectedProvider?.name ?? ''"
      :api-key="selectedProvider?.api_key ?? ''"
      :provider-base-url="selectedProvider?.base_url ?? ''"
      :existing-models="selectedProvider?.models ?? []"
      :can-import="hasPermission('ai:config:importModel')"
      @import="handleImportModels"
    />
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { Modal, message } from 'ant-design-vue'
import {
  CloudDownloadOutlined,
  EditOutlined,
  PlusOutlined,
  SaveOutlined,
} from '@ant-design/icons-vue'
import AdminSplitLayout from '@/components/AdminSplitLayout.vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import { aiTest, getAIConfig, updateAIConfig } from '@/api/ai'
import { cloneFromSnapshot, createSnapshot, isSnapshotDirty } from '../config-tab-guard'
import ProviderEditorDrawer from './ai-config/ProviderEditorDrawer.vue'
import ProviderModelEditorDrawer from './ai-config/ProviderModelEditorDrawer.vue'
import ProviderRemoteModelImportDrawer from './ai-config/ProviderRemoteModelImportDrawer.vue'
import { usePermission } from '@/utils/permission'
import {
  capabilityTabOptions,
  capabilityTagMetaMap,
  createEmptyModel,
  createEmptyProvider,
  filterModelsByCapabilityAndKeyword,
  formatSearchStrategyLabel,
  getModelCapabilityTags,
  matchesModelCapability,
  mergeImportedModels,
  normalizeAIConfig,
  normalizeModel,
  serializeModel,
  type AIConfigState,
  type AIModel,
  type AIModelCapabilityKey,
} from './ai-config/state'

interface ProviderEditorSubmitValue {
  name: string
  api_key: string
  base_url: string
  isDefault: boolean
}

const emit = defineEmits<{
  (e: 'dirty-change', value: boolean): void
}>()

const loading = ref(true)
const saving = ref(false)
const initialized = ref(false)
const { hasPermission } = usePermission()
const selectedProviderIndex = ref(0)
const activeModelIndex = ref(-1)
const activeCapabilityTab = ref<AIModelCapabilityKey>('all')
const modelKeyword = ref('')
const selectedBatchIndices = ref<number[]>([])
const testingModel = ref(false)

const providerDrawerOpen = ref(false)
const providerDrawerMode = ref<'create' | 'edit'>('create')
const editingProviderIndex = ref<number | null>(null)

const modelEditorOpen = ref(false)
const modelEditorMode = ref<'create' | 'edit'>('create')
const editingModelIndex = ref<number | null>(null)
const modelEditorInitialValue = ref<AIModel>(createEmptyModel())

const importDrawerOpen = ref(false)
const providerSaving = ref(false)

const formData = reactive<AIConfigState>({
  default_provider: '',
  providers: [],
})

const selectedProvider = computed(() => formData.providers[selectedProviderIndex.value] ?? null)
const providerNames = computed(() => formData.providers.map(provider => provider.name))
const editingProvider = computed(() => (
  editingProviderIndex.value === null ? undefined : formData.providers[editingProviderIndex.value]
))
const editingProviderIsDefault = computed(() => (
  editingProviderIndex.value === null
    ? formData.providers.length === 0
    : isDefaultProvider(editingProviderIndex.value)
))
const activeModel = computed(() => {
  const provider = selectedProvider.value
  if (!provider || activeModelIndex.value < 0 || activeModelIndex.value >= provider.models.length) {
    return null
  }
  return provider.models[activeModelIndex.value]
})
const selectedBatchSet = computed(() => new Set(selectedBatchIndices.value))
const filteredProviderModels = computed(() => {
  const provider = selectedProvider.value
  if (!provider) {
    return []
  }
  const filtered = filterModelsByCapabilityAndKeyword(
    provider.models,
    activeCapabilityTab.value,
    modelKeyword.value,
  )
  const visibleSet = new Set(filtered)
  return provider.models
    .map((model, index) => ({ model, index }))
    .filter(entry => visibleSet.has(entry.model))
})

const getConfigState = (): AIConfigState => ({
  default_provider: formData.default_provider,
  providers: formData.providers.map(provider => ({
    name: provider.name,
    api_key: provider.api_key,
    base_url: provider.base_url,
    models: provider.models.map(model => serializeModel(model)),
  })),
})

const getProviderOnlyState = (): AIConfigState => ({
  default_provider: formData.default_provider,
  providers: formData.providers.map(provider => ({
    name: provider.name,
    api_key: provider.api_key,
    base_url: provider.base_url,
    models: [],
  })),
})

const getModelOnlyState = () => formData.providers.map(provider => ({
  name: provider.name,
  models: provider.models.map(model => serializeModel(model)),
}))

const baselineSnapshot = ref(createSnapshot(getConfigState()))
const hasUnsavedChanges = computed(() => initialized.value && isSnapshotDirty(baselineSnapshot.value, getConfigState()))
const providerBaselineSnapshot = ref(createSnapshot(getProviderOnlyState()))
const modelBaselineSnapshot = ref(createSnapshot(getModelOnlyState()))
const hasModelUnsavedChanges = computed(() => initialized.value && isSnapshotDirty(modelBaselineSnapshot.value, getModelOnlyState()))

watch(
  hasUnsavedChanges,
  value => emit('dirty-change', value),
  { immediate: true },
)

watch(selectedProviderIndex, () => {
  activeCapabilityTab.value = 'all'
  modelKeyword.value = ''
  syncModelSelection(0)
})

const syncSelectedProviderIndex = (preferredIndex = selectedProviderIndex.value) => {
  if (formData.providers.length === 0) {
    selectedProviderIndex.value = 0
    return
  }
  selectedProviderIndex.value = Math.min(Math.max(preferredIndex, 0), formData.providers.length - 1)
}

const syncModelSelection = (preferredIndex = activeModelIndex.value) => {
  const provider = selectedProvider.value
  if (!provider || provider.models.length === 0) {
    activeModelIndex.value = -1
    selectedBatchIndices.value = []
    return
  }
  activeModelIndex.value = Math.min(Math.max(preferredIndex, 0), provider.models.length - 1)
  selectedBatchIndices.value = selectedBatchIndices.value.filter(index => index >= 0 && index < provider.models.length)
}

const applyConfigState = (state?: Partial<AIConfigState> | null) => {
  const normalized = normalizeAIConfig(state)
  formData.default_provider = normalized.default_provider
  formData.providers = normalized.providers
  if (!formData.default_provider && formData.providers.length > 0) {
    formData.default_provider = formData.providers[0].name
  }
  syncSelectedProviderIndex()
  syncModelSelection(0)
}

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
    providerBaselineSnapshot.value = createSnapshot(getProviderOnlyState())
    modelBaselineSnapshot.value = createSnapshot(getModelOnlyState())
    initialized.value = true
    loading.value = false
  }
}

const selectProvider = (index: number) => {
  selectedProviderIndex.value = index
}

const selectActiveModel = (index: number) => {
  activeModelIndex.value = index
}

const toggleBatchModel = (index: number) => {
  if (selectedBatchSet.value.has(index)) {
    selectedBatchIndices.value = selectedBatchIndices.value.filter(item => item !== index)
    return
  }
  selectedBatchIndices.value = [...selectedBatchIndices.value, index].sort((a, b) => a - b)
}

const isDefaultProvider = (index: number) => formData.providers[index]?.name === formData.default_provider

const ensureDefaultProvider = () => {
  if (formData.providers.length === 0) {
    formData.default_provider = ''
    return
  }
  if (!formData.providers.some(provider => provider.name === formData.default_provider)) {
    formData.default_provider = formData.providers[0].name
  }
}

const openCreateProviderDrawer = () => {
  if (!hasPermission('ai:config:createProvider')) {
    message.warning('无权新增 AI 平台')
    return
  }
  providerDrawerMode.value = 'create'
  editingProviderIndex.value = null
  providerDrawerOpen.value = true
}

const openEditProviderDrawer = (index = selectedProviderIndex.value) => {
  if (!hasPermission('ai:config:editProvider')) {
    message.warning('无权编辑 AI 平台')
    return
  }
  if (index < 0 || index >= formData.providers.length) {
    message.warning('请先选择要编辑的平台')
    return
  }
  providerDrawerMode.value = 'edit'
  editingProviderIndex.value = index
  providerDrawerOpen.value = true
}

const syncSnapshotsAfterProviderSave = () => {
  baselineSnapshot.value = createSnapshot(getConfigState())
  providerBaselineSnapshot.value = createSnapshot(getProviderOnlyState())
  modelBaselineSnapshot.value = createSnapshot(getModelOnlyState())
}

const persistProviderChanges = async (successMessage: string) => {
  providerSaving.value = true
  try {
    await updateAIConfig(getConfigState())
    syncSnapshotsAfterProviderSave()
    message.success(successMessage)
    return true
  } catch {
    message.error('保存平台配置失败')
    return false
  } finally {
    providerSaving.value = false
  }
}

const confirmDiscardModelDraftsForProviderAction = async () => {
  if (!hasModelUnsavedChanges.value) {
    return true
  }
  return new Promise<boolean>((resolve) => {
    Modal.confirm({
      title: '当前模型配置有未保存修改',
      content: '平台操作需要立即保存。继续后将放弃当前模型草稿修改，是否继续？',
      okText: '继续',
      cancelText: '取消',
      onOk: () => {
        discardChanges()
        resolve(true)
      },
      onCancel: () => resolve(false),
    })
  })
}

const handleProviderSubmit = async (value: ProviderEditorSubmitValue) => {
  if (!hasPermission(providerDrawerMode.value === 'edit' ? 'ai:config:editProvider' : 'ai:config:createProvider')) {
    message.warning('无权修改 AI 平台')
    return
  }
  const canContinue = await confirmDiscardModelDraftsForProviderAction()
  if (!canContinue) {
    return
  }
  const previousSnapshot = baselineSnapshot.value
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
    syncModelSelection(-1)
  } else if (editingProviderIndex.value !== null) {
    const provider = formData.providers[editingProviderIndex.value]
    const previousName = provider.name
    provider.name = value.name
    provider.api_key = value.api_key
    provider.base_url = value.base_url
    if (value.isDefault || formData.default_provider === previousName) {
      formData.default_provider = value.name
    }
    ensureDefaultProvider()
    selectedProviderIndex.value = editingProviderIndex.value
  }
  const saved = await persistProviderChanges(providerDrawerMode.value === 'create' ? '平台新增成功' : '平台保存成功')
  if (saved) {
    providerDrawerOpen.value = false
    return
  }
  applyConfigState(cloneFromSnapshot<AIConfigState>(previousSnapshot))
}

const handleProviderRemove = async () => {
  if (!hasPermission('ai:config:deleteProvider')) {
    message.warning('无权删除 AI 平台')
    return
  }
  const canContinue = await confirmDiscardModelDraftsForProviderAction()
  if (!canContinue) {
    return
  }
  const previousSnapshot = baselineSnapshot.value
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
  syncModelSelection(0)
  const saved = await persistProviderChanges(`已删除平台 ${removed.name || '未命名平台'}`)
  if (!saved) {
    applyConfigState(cloneFromSnapshot<AIConfigState>(previousSnapshot))
    return
  }
  providerDrawerOpen.value = false
}

const handleSetDefaultProvider = async (index: number) => {
  if (!hasPermission('ai:config:save')) {
    message.warning('无权保存 AI 配置')
    return
  }
  if (index < 0 || index >= formData.providers.length || isDefaultProvider(index)) {
    return
  }
  const canContinue = await confirmDiscardModelDraftsForProviderAction()
  if (!canContinue) {
    return
  }
  const previousSnapshot = baselineSnapshot.value
  formData.default_provider = formData.providers[index].name
  ensureDefaultProvider()
  selectedProviderIndex.value = index
  const saved = await persistProviderChanges(`默认平台已切换为 ${formData.default_provider}`)
  if (!saved) {
    applyConfigState(cloneFromSnapshot<AIConfigState>(previousSnapshot))
  }
}

const hasModelIDConflict = (value: AIModel, ignoreIndex = -1) => {
  const provider = selectedProvider.value
  if (!provider || !value.id.trim()) {
    return false
  }
  return provider.models.some((model, index) => index !== ignoreIndex && model.id.trim() === value.id.trim())
}

const openCreateModelDrawer = () => {
  if (!hasPermission('ai:config:createModel')) {
    message.warning('无权新增模型')
    return
  }
  if (!selectedProvider.value) {
    message.warning('请先选择平台')
    return
  }
  modelEditorMode.value = 'create'
  editingModelIndex.value = null
  modelEditorInitialValue.value = createEmptyModel()
  modelEditorOpen.value = true
}

const openEditModelDrawer = (index = activeModelIndex.value) => {
  if (!hasPermission('ai:config:editModel')) {
    message.warning('无权编辑模型')
    return
  }
  const provider = selectedProvider.value
  if (!provider || index < 0 || index >= provider.models.length) {
    message.warning('请先选择要编辑的模型')
    return
  }
  modelEditorMode.value = 'edit'
  editingModelIndex.value = index
  modelEditorInitialValue.value = normalizeModel(provider.models[index])
  modelEditorOpen.value = true
}

const handleModelSubmit = (value: AIModel) => {
  if (!hasPermission(modelEditorMode.value === 'edit' ? 'ai:config:editModel' : 'ai:config:createModel')) {
    message.warning('无权修改模型')
    return
  }
  const provider = selectedProvider.value
  if (!provider) {
    return
  }
  const normalized = normalizeModel(value)
  const ignoreIndex = modelEditorMode.value === 'edit' ? (editingModelIndex.value ?? -1) : -1
  if (hasModelIDConflict(normalized, ignoreIndex)) {
    message.warning(`当前平台下已存在模型 ID：${normalized.id}`)
    return
  }

  if (modelEditorMode.value === 'create') {
    provider.models.push(normalized)
    activeModelIndex.value = provider.models.length - 1
  } else if (editingModelIndex.value !== null) {
    provider.models[editingModelIndex.value] = normalized
    activeModelIndex.value = editingModelIndex.value
  }

  syncModelSelection(activeModelIndex.value)
  modelEditorOpen.value = false
}

const handleModelRemoveFromDrawer = () => {
  handleRemoveModelByIndex(editingModelIndex.value ?? activeModelIndex.value)
  modelEditorOpen.value = false
}

const handleRemoveModelByIndex = (index: number) => {
  if (!hasPermission('ai:config:deleteModel')) {
    message.warning('无权删除模型')
    return
  }
  const provider = selectedProvider.value
  if (!provider || index < 0 || index >= provider.models.length) {
    return
  }
  provider.models.splice(index, 1)
  selectedBatchIndices.value = selectedBatchIndices.value
    .filter(item => item !== index)
    .map(item => (item > index ? item - 1 : item))
  syncModelSelection(Math.min(index, provider.models.length - 1))
  message.success('模型已删除')
}

const handleRemoveActiveModel = () => {
  handleRemoveModelByIndex(activeModelIndex.value)
}

const handleBatchRemoveModels = () => {
  if (!hasPermission('ai:config:deleteModel')) {
    message.warning('无权删除模型')
    return
  }
  const provider = selectedProvider.value
  if (!provider || selectedBatchIndices.value.length === 0) {
    return
  }
  const indices = [...selectedBatchIndices.value].sort((a, b) => b - a)
  for (const index of indices) {
    provider.models.splice(index, 1)
  }
  selectedBatchIndices.value = []
  syncModelSelection(Math.min(activeModelIndex.value, provider.models.length - 1))
  message.success(`已删除 ${indices.length} 个模型`)
}

const moveActiveModel = (direction: number) => {
  if (!hasPermission('ai:config:editModel')) {
    message.warning('无权调整模型排序')
    return
  }
  const provider = selectedProvider.value
  if (!provider || activeModelIndex.value < 0) {
    return
  }
  const currentIndex = activeModelIndex.value
  const nextIndex = currentIndex + direction
  if (nextIndex < 0 || nextIndex >= provider.models.length) {
    return
  }
  const current = provider.models[currentIndex]
  provider.models[currentIndex] = provider.models[nextIndex]
  provider.models[nextIndex] = current
  selectedBatchIndices.value = selectedBatchIndices.value.map((item) => {
    if (item === currentIndex) {
      return nextIndex
    }
    if (item === nextIndex) {
      return currentIndex
    }
    return item
  }).sort((a, b) => a - b)
  activeModelIndex.value = nextIndex
}

const runModelTest = async (model: AIModel) => {
  if (!hasPermission('ai:config:test')) {
    message.warning('无权测试模型')
    return
  }
  const provider = selectedProvider.value
  if (!provider) {
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

  testingModel.value = true
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
    testingModel.value = false
  }
}

const handleTestModel = async (model: AIModel) => {
  await runModelTest(model)
}

const handleTestActiveModel = async () => {
  if (!activeModel.value) {
    message.warning('请先选择模型')
    return
  }
  await runModelTest(activeModel.value)
}

const openImportDrawer = () => {
  if (!hasPermission('ai:config:importModel')) {
    message.warning('无权导入模型')
    return
  }
  if (!selectedProvider.value) {
    message.warning('请先选择平台')
    return
  }
  importDrawerOpen.value = true
}

const handleImportModels = (models: AIModel[]) => {
  if (!hasPermission('ai:config:importModel')) {
    message.warning('无权导入模型')
    return
  }
  const provider = selectedProvider.value
  if (!provider) {
    return
  }
  const merged = mergeImportedModels(provider.models, models)
  provider.models = merged.models
  if (merged.importedCount > 0) {
    const firstImportedId = models.find(model => model.id)?.id ?? ''
    const nextIndex = provider.models.findIndex(model => model.id === firstImportedId)
    syncModelSelection(nextIndex >= 0 ? nextIndex : provider.models.length - 1)
    message.success(`已导入 ${merged.importedCount} 个模型${merged.skippedCount > 0 ? `，跳过 ${merged.skippedCount} 个重复模型` : ''}`)
    return
  }
  message.info('所选模型都已存在，未重复导入')
}

const getCapabilityCount = (capability: AIModelCapabilityKey) => {
  const provider = selectedProvider.value
  if (!provider) {
    return 0
  }
  return provider.models.filter(model => matchesModelCapability(model, capability)).length
}

const formatTemperature = (value: number | null) => (typeof value === 'number' ? value.toFixed(1) : '-')
const formatContextWindow = (value: number | null) => (typeof value === 'number' && value > 0 ? `${value}` : '-')

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
    for (let index = 0; index < provider.models.length; index += 1) {
      const normalized = serializeModel(provider.models[index])
      provider.models[index] = normalized
      if (!normalized.id) {
        message.warning(`请填写 ${provider.name} 平台下的模型 ID`)
        return false
      }
      if (modelIDSet.has(normalized.id)) {
        message.warning(`平台 ${provider.name} 下存在重复模型 ID：${normalized.id}`)
        return false
      }
      modelIDSet.add(normalized.id)
    }
  }
  ensureDefaultProvider()
  return true
}

const save = async () => {
  if (!hasPermission('ai:config:save')) {
    message.warning('无权保存 AI 配置')
    return false
  }
  if (!validateBeforeSave()) {
    return false
  }
  saving.value = true
  try {
    await updateAIConfig(getConfigState())
    baselineSnapshot.value = createSnapshot(getConfigState())
    providerBaselineSnapshot.value = createSnapshot(getProviderOnlyState())
    modelBaselineSnapshot.value = createSnapshot(getModelOnlyState())
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
  modelEditorOpen.value = false
  importDrawerOpen.value = false
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
.ai-config-layout {
  min-height: 680px;
  color: var(--app-text-color);
}

.provider-sidebar,
.workspace-panel {
  border: 1px solid var(--app-border-color);
  border-radius: 18px;
  background: var(--app-surface-color);
}

.provider-sidebar {
  display: flex;
  flex-direction: column;
  padding: 18px 16px;
  background: linear-gradient(180deg, var(--app-surface-color) 0%, var(--app-surface-soft) 100%);
}

.workspace-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px;
  min-width: 0;
}

.sidebar-header,
.workspace-hero,
.workspace-toolbar,
.workspace-filters,
.model-card__head,
.model-card__footer {
  display: flex;
  gap: 16px;
  align-items: flex-start;
  justify-content: space-between;
}

.sidebar-title,
.workspace-hero__title,
.toolbar-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--app-text-strong);
}

.sidebar-subtitle,
.workspace-hero__subtitle,
.toolbar-meta,
.summary-chip,
.model-card__id,
.model-card__meta,
.model-card__desc,
.model-card__placeholder,
.model-card__footer {
  color: var(--app-text-secondary);
}

.sidebar-subtitle,
.workspace-hero__subtitle,
.toolbar-meta {
  margin-top: 4px;
}

.provider-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 18px;
  min-height: 0;
  overflow-y: auto;
  padding-right: 4px;
}

.provider-item {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border: 1px solid var(--app-border-color);
  border-radius: 16px;
  background: var(--app-surface-color);
  padding: 12px 14px;
  transition: all 0.2s ease;
}

.provider-item:hover {
  border-color: var(--app-primary-color);
  background: var(--app-hover-bg);
}

.provider-item__main {
  min-width: 0;
  flex: 1;
  padding: 0;
  border: 0;
  background: transparent;
  color: var(--app-text-color);
  cursor: pointer;
  text-align: left;
}

.provider-item__summary {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.provider-item__actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.provider-item--active {
  border-color: var(--app-primary-color);
  background: var(--app-primary-color-soft);
  box-shadow: 0 12px 24px rgba(24, 144, 255, 0.12);
}

.provider-item__name,
.model-card__name {
  font-weight: 600;
  color: var(--app-text-strong);
  min-width: 0;
}

.provider-item__count {
  flex-shrink: 0;
  font-size: 12px;
  color: var(--app-text-secondary);
}

.provider-action-button {
  padding-inline: 2px;
  height: 22px;
  font-size: 12px;
}

.workspace-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.summary-chip {
  padding: 6px 12px;
  border-radius: 999px;
  border: 1px solid var(--app-border-color);
  background: var(--app-surface-soft);
  font-size: 13px;
}

.workspace-toolbar {
  padding: 14px 16px;
  border: 1px solid var(--app-border-color);
  border-radius: 16px;
  background: var(--app-surface-soft);
}

.workspace-toolbar__main {
  min-width: 0;
}

.workspace-tabs {
  min-width: 0;
  flex: 1;
}

.workspace-search {
  width: 320px;
  flex-shrink: 0;
}

.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 14px;
}

.model-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border: 1px solid var(--app-border-color);
  border-radius: 18px;
  background: linear-gradient(180deg, var(--app-surface-color) 0%, var(--app-surface-soft) 100%);
  transition: all 0.2s ease;
  cursor: pointer;
}

.model-card:hover,
.model-card--active {
  border-color: var(--app-primary-color);
  box-shadow: 0 16px 28px rgba(15, 23, 42, 0.1);
  transform: translateY(-2px);
}

.model-card__title-wrap {
  min-width: 0;
  flex: 1;
}

.model-card__id,
.model-card__meta,
.model-card__desc {
  font-size: 12px;
  word-break: break-all;
}

.model-card__tags,
.model-card__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.ai-config :deep(.ant-empty-description) {
  color: var(--app-text-secondary);
}

.ai-config :deep(.ant-tabs-nav) {
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .sidebar-header,
  .workspace-hero,
  .workspace-toolbar,
  .workspace-filters,
  .model-card__head,
  .model-card__footer {
    flex-direction: column;
  }

  .workspace-search {
    width: 100%;
  }

  .model-grid {
    grid-template-columns: 1fr;
  }
}
</style>
