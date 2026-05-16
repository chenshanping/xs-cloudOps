<template>
  <a-drawer
    :open="open"
    :title="`自动识别导入 - ${providerName || '未命名平台'}`"
    width="1180"
    placement="right"
    :body-style="{ padding: 0, overflow: 'hidden' }"
    :mask-closable="false"
    @close="handleClose"
  >
    <div class="import-drawer">
      <aside class="import-sidebar">
        <div class="import-sidebar__group">
          <div class="import-sidebar__title">模型类型</div>
          <div class="import-sidebar__menu">
            <button
              v-for="item in primaryFilterOptions"
              :key="item.value"
              type="button"
              :class="[
                'import-sidebar__item',
                primaryFilter === item.value ? 'import-sidebar__item--active' : '',
              ]"
              @click="primaryFilter = item.value"
            >
              <span>{{ item.label }}</span>
              <span class="import-sidebar__count">{{ getPrimaryFilterCount(item.value) }}</span>
            </button>
          </div>
        </div>

        <div class="import-sidebar__divider"></div>

        <div class="import-sidebar__group">
          <div class="import-sidebar__title">能力特征</div>
          <div class="import-sidebar__menu">
            <button
              v-for="item in featureFilterOptions"
              :key="item.value"
              type="button"
              :class="[
                'import-sidebar__item',
                activeFeatureFilters.includes(item.value) ? 'import-sidebar__item--active' : '',
                getFeatureFilterCount(item.value) === 0 && !activeFeatureFilters.includes(item.value)
                  ? 'import-sidebar__item--disabled'
                  : '',
              ]"
              :disabled="getFeatureFilterCount(item.value) === 0 && !activeFeatureFilters.includes(item.value)"
              @click="toggleFeatureFilter(item.value)"
            >
              <span>{{ item.label }}</span>
              <span class="import-sidebar__count">{{ getFeatureFilterCount(item.value) }}</span>
            </button>
          </div>
        </div>
      </aside>

      <section class="import-content">
        <div class="import-panel">
          <div class="import-panel__summary">
            <div>
              <div class="import-panel__title">远端模型列表</div>
              <div class="import-panel__subtitle">
                默认按模型分组收起；远端能力由后端识别，这里只负责筛选、选择并回写到当前页面
              </div>
            </div>
            <div class="import-panel__stats">
              <span>总数 {{ remoteModels.length }}</span>
              <span>当前 {{ filteredRemoteModels.length }}</span>
              <span>已导入 {{ existingModelIDs.size }}</span>
              <span>已选 {{ targetKeys.length }}</span>
            </div>
          </div>

          <div class="import-panel__toolbar">
            <a-input-search
              v-model:value="keyword"
              allow-clear
              class="import-panel__search"
              placeholder="搜索模型名称 / 标识 / 分组"
            />
            <a-space wrap>
              <a-button :loading="loading" :disabled="!canImport" @click="handleFetch">重新获取</a-button>
              <a-button @click="handleResetFilters">重置筛选</a-button>
              <a-button :disabled="filteredImportableModels.length === 0" @click="handleSelectFiltered">
                全选当前筛选
              </a-button>
              <a-button :disabled="targetKeys.length === 0" @click="handleClearSelection">
                清空选择
              </a-button>
            </a-space>
          </div>

          <div v-if="activeFilterTags.length > 0" class="import-panel__filters">
            <a-tag
              v-for="tag in activeFilterTags"
              :key="tag.key"
              closable
              color="blue"
              @close.prevent="removeFilterTag(tag.key)"
            >
              {{ tag.label }}
            </a-tag>
          </div>

          <div class="import-panel__body">
            <a-spin :spinning="loading">
              <template v-if="groupedRemoteModels.length > 0">
                <a-collapse v-model:activeKey="expandedGroupKeys" ghost class="import-group-collapse">
                  <a-collapse-panel
                    v-for="group in groupedRemoteModels"
                    :key="group.name"
                    :header="group.name"
                    class="import-group-panel"
                  >
                    <template #extra>
                      <div class="import-group__extra" @click.stop>
                        <span class="import-group__count">{{ group.items.length }}</span>
                        <span v-if="group.importableCount > 0" class="import-group__meta">
                          待导入 {{ group.importableCount }}
                        </span>
                        <span v-if="group.importedCount > 0" class="import-group__meta">
                          已导入 {{ group.importedCount }}
                        </span>
                        <span v-if="group.selectedCount > 0" class="import-group__meta import-group__meta--selected">
                          已选 {{ group.selectedCount }}
                        </span>
                        <a-button
                          v-if="group.importableCount > 0"
                          type="link"
                          size="small"
                          class="import-group__action"
                          @click.stop="handleSelectGroup(group.items)"
                        >
                          全选本组
                        </a-button>
                        <a-button
                          v-if="group.selectedCount > 0"
                          type="link"
                          size="small"
                          class="import-group__action"
                          @click.stop="handleClearGroup(group.items)"
                        >
                          清空选择
                        </a-button>
                        <a-button
                          v-if="group.importableCount > 0"
                          type="link"
                          size="small"
                          class="import-group__action"
                          @click.stop="handleImportGroup(group.items)"
                        >
                          导入本组
                        </a-button>
                      </div>
                    </template>

                    <div class="import-list">
                      <article
                        v-for="model in group.items"
                        :key="model.id"
                        :class="[
                          'import-model-row',
                          targetKeySet.has(model.id) ? 'import-model-row--selected' : '',
                          existingModelIDs.has(model.id) ? 'import-model-row--imported' : '',
                        ]"
                        @click="handleToggleModel(model.id)"
                      >
                        <div class="import-model-row__head">
                          <a-checkbox
                            :checked="targetKeySet.has(model.id)"
                            :disabled="existingModelIDs.has(model.id)"
                            @click.stop
                            @change="handleToggleModel(model.id)"
                          />
                          <div class="import-model-row__title-wrap">
                            <div class="import-model-row__title">{{ model.name || model.id }}</div>
                            <div class="import-model-row__key">{{ model.id }}</div>
                          </div>
                          <a-tag :color="existingModelIDs.has(model.id) ? 'default' : 'green'">
                            {{ existingModelIDs.has(model.id) ? '已导入' : '待导入' }}
                          </a-tag>
                        </div>

                        <div class="capability-tags capability-tags--row">
                          <a-tag
                            v-for="tag in getModelCapabilityTags(model)"
                            :key="`${model.id}-${tag}`"
                            :color="capabilityTagMetaMap[tag].color"
                          >
                            {{ capabilityTagMetaMap[tag].label }}
                          </a-tag>
                          <span v-if="getModelCapabilityTags(model).length === 0" class="capability-tags__empty">
                            未识别
                          </span>
                        </div>

                        <div class="import-model-row__meta">
                          <span>联网 {{ formatSearchStrategyLabel(model.search_strategy) }}</span>
                          <span>温度 {{ formatTemperature(model.temperature) }}</span>
                          <span>上下文 {{ formatContextWindow(model.context_window) }}</span>
                        </div>

                        <div v-if="model.description" class="import-model-row__desc">
                          {{ model.description }}
                        </div>
                      </article>
                    </div>
                  </a-collapse-panel>
                </a-collapse>
              </template>

              <a-empty
                v-else
                :description="remoteModels.length === 0 ? '暂无远端模型，请先获取模型列表' : '当前筛选条件下没有模型'"
              />
            </a-spin>
          </div>
        </div>
      </section>
    </div>

    <template #footer>
      <div class="import-drawer__footer">
        <div class="import-drawer__footer-summary">已选择 {{ targetKeys.length }} 个模型</div>
        <a-space>
          <a-button @click="handleClose">关闭</a-button>
          <a-button type="primary" :disabled="targetKeys.length === 0 || !canImport" @click="handleConfirm">
            导入已选模型
          </a-button>
        </a-space>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { fetchAIProviderModels } from '@/api/ai'
import {
  capabilityTagMetaMap,
  createEmptyModel,
  filterModelsByCapabilityAndKeyword,
  formatSearchStrategyLabel,
  getModelCapabilityTags,
  groupRemoteModelsByDisplayGroup,
  matchesModelCapability,
  normalizeModel,
  normalizeRemoteProviderModel,
  type AIModel,
  type AIModelCapabilityKey,
  type RemoteProviderModel,
} from './state'

type PrimaryFilterKey = 'all' | 'reasoning' | 'vision' | 'embedding' | 'rerank'
type FeatureFilterKey = 'search' | 'free' | 'tool'

const primaryFilterOptions: Array<{ label: string; value: PrimaryFilterKey }> = [
  { label: '全部模型', value: 'all' },
  { label: '推理模型', value: 'reasoning' },
  { label: '视觉模型', value: 'vision' },
  { label: '嵌入模型', value: 'embedding' },
  { label: '重排模型', value: 'rerank' },
]

const featureFilterOptions: Array<{ label: string; value: FeatureFilterKey }> = [
  { label: '联网', value: 'search' },
  { label: '免费', value: 'free' },
  { label: '工具', value: 'tool' },
]

const props = defineProps<{
  open: boolean
  providerName: string
  apiKey: string
  providerBaseUrl: string
  existingModels: AIModel[]
  canImport: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'import', value: AIModel[]): void
}>()

const loading = ref(false)
const keyword = ref('')
const primaryFilter = ref<PrimaryFilterKey>('all')
const activeFeatureFilters = ref<FeatureFilterKey[]>([])
const remoteModels = ref<RemoteProviderModel[]>([])
const targetKeys = ref<string[]>([])
const expandedGroupKeys = ref<string[]>([])

const existingModelIDs = computed(() => new Set(props.existingModels.map(model => model.id).filter(Boolean)))
const targetKeySet = computed(() => new Set(targetKeys.value))

const filteredRemoteModels = computed(() => {
  const primaryFiltered = filterModelsByCapabilityAndKeyword(
    remoteModels.value,
    primaryFilter.value as AIModelCapabilityKey,
    keyword.value,
  )
  return primaryFiltered.filter(model => (
    activeFeatureFilters.value.every(filter => matchesModelCapability(model, filter))
  ))
})

const filteredImportableModels = computed(() => (
  filteredRemoteModels.value.filter(model => !existingModelIDs.value.has(model.id))
))

const groupedRemoteModels = computed(() => (
  groupRemoteModelsByDisplayGroup(filteredRemoteModels.value).map(group => ({
    ...group,
    importableCount: group.items.filter(model => !existingModelIDs.value.has(model.id)).length,
    importedCount: group.items.filter(model => existingModelIDs.value.has(model.id)).length,
    selectedCount: group.items.filter(model => targetKeySet.value.has(model.id)).length,
  }))
))

const selectedModels = computed(() => (
  targetKeys.value
    .map(id => remoteModels.value.find(model => model.id === id))
    .filter((model): model is RemoteProviderModel => Boolean(model))
))

const activeFilterTags = computed(() => {
  const tags: Array<{ key: string; label: string }> = []
  if (primaryFilter.value !== 'all') {
    tags.push({
      key: `primary:${primaryFilter.value}`,
      label: primaryFilterOptions.find(item => item.value === primaryFilter.value)?.label ?? primaryFilter.value,
    })
  }
  for (const item of activeFeatureFilters.value) {
    tags.push({
      key: `feature:${item}`,
      label: featureFilterOptions.find(option => option.value === item)?.label ?? item,
    })
  }
  if (keyword.value.trim()) {
    tags.push({
      key: 'keyword',
      label: `搜索：${keyword.value.trim()}`,
    })
  }
  return tags
})

watch(
  () => props.open,
  open => {
    if (!open) {
      handleResetFilters()
      targetKeys.value = []
      expandedGroupKeys.value = []
      return
    }
    handleResetFilters()
    targetKeys.value = []
    expandedGroupKeys.value = []
    if (props.apiKey.trim() && props.providerBaseUrl.trim()) {
      void handleFetch()
    }
  },
)

watch(
  groupedRemoteModels,
  groups => {
    const validKeys = new Set(groups.map(group => group.name))
    expandedGroupKeys.value = expandedGroupKeys.value.filter(key => validKeys.has(key))
  },
  { deep: true },
)

watch(keyword, value => {
  if (!props.open) {
    return
  }
  expandedGroupKeys.value = value.trim()
    ? groupedRemoteModels.value.map(group => group.name)
    : []
})

const getPrimaryFilterCount = (value: PrimaryFilterKey) => (
  filterModelsByCapabilityAndKeyword(remoteModels.value, value as AIModelCapabilityKey).length
)

const getFeatureFilterCount = (value: FeatureFilterKey) => (
  filteredByPrimary(remoteModels.value).filter(model => matchesModelCapability(model, value)).length
)

const filteredByPrimary = (models: RemoteProviderModel[]) => (
  filterModelsByCapabilityAndKeyword(models, primaryFilter.value as AIModelCapabilityKey)
)

const toggleFeatureFilter = (value: FeatureFilterKey) => {
  if (activeFeatureFilters.value.includes(value)) {
    activeFeatureFilters.value = activeFeatureFilters.value.filter(item => item !== value)
    return
  }
  activeFeatureFilters.value = [...activeFeatureFilters.value, value]
}

const handleResetFilters = () => {
  keyword.value = ''
  primaryFilter.value = 'all'
  activeFeatureFilters.value = []
  expandedGroupKeys.value = []
}

const removeFilterTag = (key: string) => {
  if (key === 'keyword') {
    keyword.value = ''
    return
  }
  if (key.startsWith('primary:')) {
    primaryFilter.value = 'all'
    return
  }
  if (key.startsWith('feature:')) {
    const value = key.replace('feature:', '') as FeatureFilterKey
    activeFeatureFilters.value = activeFeatureFilters.value.filter(item => item !== value)
  }
}

const handleClose = () => {
  emit('update:open', false)
}

const formatTemperature = (value: number | null) => (typeof value === 'number' ? value.toFixed(1) : '-')
const formatContextWindow = (value: number | null) => (typeof value === 'number' && value > 0 ? `${value}` : '-')

const handleFetch = async () => {
  if (!props.canImport) {
    return
  }
  if (!props.apiKey.trim()) {
    message.warning('请先填写当前平台的 API Key')
    return
  }
  if (!props.providerBaseUrl.trim()) {
    message.warning('请先填写当前平台的 Base URL')
    return
  }

  loading.value = true
  try {
    const res = await fetchAIProviderModels({
      api_key: props.apiKey,
      base_url: props.providerBaseUrl,
      provider_name: props.providerName,
    })
    remoteModels.value = (res.data.models ?? []).map(model => normalizeRemoteProviderModel(model))
    targetKeys.value = targetKeys.value.filter(key => (
      remoteModels.value.some(model => model.id === key && !existingModelIDs.value.has(model.id))
    ))
    message.success(`已获取 ${remoteModels.value.length} 个平台模型`)
  } catch (error: any) {
    message.error(error.message || '获取模型列表失败')
  } finally {
    loading.value = false
  }
}

const handleToggleModel = (id: string) => {
  if (existingModelIDs.value.has(id)) {
    return
  }
  if (targetKeySet.value.has(id)) {
    targetKeys.value = targetKeys.value.filter(item => item !== id)
    return
  }
  targetKeys.value = [...targetKeys.value, id]
}

const handleSelectFiltered = () => {
  const merged = new Set([...targetKeys.value, ...filteredImportableModels.value.map(model => model.id)])
  targetKeys.value = Array.from(merged)
}

const handleSelectGroup = (models: RemoteProviderModel[]) => {
  const importableIDs = models
    .filter(model => !existingModelIDs.value.has(model.id))
    .map(model => model.id)
  const merged = new Set([...targetKeys.value, ...importableIDs])
  targetKeys.value = Array.from(merged)
}

const handleClearGroup = (models: RemoteProviderModel[]) => {
  const ids = new Set(models.map(model => model.id))
  targetKeys.value = targetKeys.value.filter(id => !ids.has(id))
}

const handleClearSelection = () => {
  targetKeys.value = []
}

const emitImportedModels = (models: RemoteProviderModel[]) => {
  if (!props.canImport) {
    return
  }
  const modelsToImport = models.filter(model => !existingModelIDs.value.has(model.id))
  if (modelsToImport.length === 0) {
    message.warning('请先选择要导入的模型')
    return
  }

  const importedModels = modelsToImport.map(model => normalizeModel({
    ...createEmptyModel(),
    id: model.id,
    name: model.name || model.id,
    group: model.group || '',
    description: model.description || '',
    is_thinking: model.is_thinking,
    support_vision: model.support_vision,
    support_tools: model.support_tools,
    search_strategy: model.search_strategy,
    support_embedding: model.support_embedding,
    support_rerank: model.support_rerank,
    is_free: model.is_free,
    temperature: model.temperature,
    context_window: model.context_window,
  }))

  emit('import', importedModels)
  handleClose()
}

const handleImportGroup = (models: RemoteProviderModel[]) => {
  emitImportedModels(models)
}

const handleConfirm = () => {
  emitImportedModels(selectedModels.value)
}
</script>

<style scoped>
.import-drawer {
  display: flex;
  height: calc(100vh - 132px);
  min-height: 560px;
}

.import-sidebar {
  width: 220px;
  padding: 22px 14px 18px;
  background: linear-gradient(180deg, #f8fafc 0%, #f1f5f9 100%);
  border-right: 1px solid #e2e8f0;
}

.import-sidebar__title {
  padding: 0 10px 12px;
  font-size: 13px;
  font-weight: 600;
  color: #334155;
}

.import-sidebar__group + .import-sidebar__group {
  margin-top: 14px;
}

.import-sidebar__divider {
  height: 1px;
  margin: 14px 8px;
  background: #dbe5f0;
}

.import-sidebar__menu {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.import-sidebar__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  width: 100%;
  padding: 10px 12px;
  color: #475569;
  cursor: pointer;
  background: rgb(255 255 255 / 92%);
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  transition: all 0.2s ease;
}

.import-sidebar__item:hover {
  border-color: #93c5fd;
  transform: translateY(-1px);
}

.import-sidebar__item--active {
  color: #0f172a;
  background: linear-gradient(135deg, #eff6ff 0%, #f8fbff 100%);
  border-color: #93c5fd;
  box-shadow: 0 8px 18px rgb(37 99 235 / 8%);
}

.import-sidebar__item--disabled {
  color: #94a3b8;
  cursor: not-allowed;
  background: #f8fafc;
  border-color: #e2e8f0;
}

.import-sidebar__item--disabled:hover {
  border-color: #e2e8f0;
  transform: none;
}

.import-sidebar__count,
.import-group__count {
  min-width: 24px;
  padding: 0 8px;
  font-size: 12px;
  line-height: 22px;
  color: #2563eb;
  text-align: center;
  background: rgb(37 99 235 / 10%);
  border-radius: 999px;
}

.import-content {
  flex: 1;
  min-width: 0;
  background: #fff;
}

.import-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 100%;
  padding: 24px;
}

.import-panel__summary,
.import-panel__toolbar,
.import-drawer__footer,
.import-model-row__head {
  display: flex;
  gap: 16px;
  align-items: flex-start;
  justify-content: space-between;
}

.import-panel__title {
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
}

.import-panel__subtitle,
.import-model-row__key,
.import-model-row__meta,
.import-model-row__desc,
.capability-tags__empty,
.import-drawer__footer-summary {
  color: #64748b;
}

.import-panel__subtitle {
  margin-top: 4px;
  font-size: 13px;
}

.import-panel__stats {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: flex-end;
}

.import-panel__stats span {
  padding: 6px 10px;
  background: rgb(255 255 255 / 78%);
  border-radius: 999px;
  border: 1px solid #e2e8f0;
  color: #475569;
}

.import-panel__search {
  max-width: 360px;
}

.import-panel__body {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding-right: 4px;
}

.import-panel__filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.import-group-collapse {
  background: transparent;
}

.import-group-collapse :deep(.ant-collapse-item) {
  margin-bottom: 12px;
  overflow: hidden;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
}

.import-group-collapse :deep(.ant-collapse-item:last-child) {
  margin-bottom: 0;
}

.import-group-collapse :deep(.ant-collapse-header) {
  align-items: center !important;
  padding: 14px 18px !important;
}

.import-group-collapse :deep(.ant-collapse-header-text) {
  font-size: 14px;
  font-weight: 600;
  color: #0f172a;
}

.import-group-collapse :deep(.ant-collapse-content-box) {
  padding: 0 !important;
}

.import-group__extra {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.import-group__meta {
  font-size: 12px;
  color: #64748b;
}

.import-group__meta--selected {
  color: #2563eb;
}

.import-group__action {
  padding-inline: 0;
}

.import-list {
  display: flex;
  flex-direction: column;
}

.capability-tags,
.import-model-row__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.import-model-row {
  padding: 16px 18px;
  background: #fff;
  border-top: 1px solid #edf2f7;
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    box-shadow 0.2s ease;
}

.import-model-row:first-child {
  border-top: 0;
}

.import-model-row:hover {
  background: #f8fbff;
}

.import-model-row--selected {
  background: #eff6ff;
  box-shadow: inset 3px 0 0 #2563eb;
}

.import-model-row--imported {
  background: #f8fafc;
  cursor: default;
}

.import-model-row--imported:hover {
  background: #f8fafc;
}

.import-model-row__title-wrap {
  min-width: 0;
  flex: 1;
}

.import-model-row__title {
  font-weight: 600;
  color: #0f172a;
}

.import-model-row__key,
.import-model-row__meta,
.import-model-row__desc {
  margin-top: 6px;
  font-size: 12px;
  word-break: break-all;
}

.capability-tags--row {
  margin-top: 12px;
}

@media (max-width: 768px) {
  .import-drawer {
    flex-direction: column;
    height: auto;
    min-height: 0;
  }

  .import-sidebar {
    width: 100%;
    border-right: 0;
    border-bottom: 1px solid #e2e8f0;
  }

  .import-sidebar__menu {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .import-panel {
    min-height: 0;
    padding: 18px 16px;
  }

  .import-panel__summary,
  .import-panel__toolbar,
  .import-drawer__footer,
  .import-model-row__head {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
