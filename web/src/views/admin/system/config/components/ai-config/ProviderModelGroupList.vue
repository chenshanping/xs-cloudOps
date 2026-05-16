<template>
  <a-collapse v-model:activeKey="activeGroupKeys" ghost class="saved-group-collapse">
    <a-collapse-panel
      v-for="group in groupedEntries"
      :key="group.name"
      :header="group.name"
      class="saved-group-panel"
    >
      <template #extra>
        <div class="saved-group-extra" @click.stop>
          <span class="saved-group-count">{{ group.items.length }}</span>
          <span v-if="getGroupSelectedCount(group.items) > 0" class="saved-group-meta">
            已选 {{ getGroupSelectedCount(group.items) }}
          </span>
          <a-button
            v-if="getGroupSelectableCount(group.items) > 0"
            type="link"
            size="small"
            class="saved-group-action"
            @click.stop="emit('select-group', group.items.map(item => item.index))"
          >
            全选本组
          </a-button>
          <a-button
            v-if="getGroupSelectedCount(group.items) > 0"
            type="link"
            size="small"
            class="saved-group-action"
            @click.stop="emit('clear-group', group.items.map(item => item.index))"
          >
            清空选择
          </a-button>
        </div>
      </template>

      <div class="saved-model-list">
        <article
          v-for="entry in group.items"
          :key="`${entry.index}-${entry.model.id || entry.model.name || 'model'}`"
          :class="[
            'saved-model-row',
            entry.index === activeModelIndex ? 'saved-model-row--active' : '',
          ]"
          @click="emit('select-model', entry.index)"
        >
          <div class="saved-model-row__head">
            <div class="saved-model-row__title-wrap">
              <div class="saved-model-row__title">
                {{ entry.model.name || entry.model.id || '未命名模型' }}
              </div>
              <div class="saved-model-row__key">{{ entry.model.id || '待填写模型 ID' }}</div>
            </div>
            <div class="saved-model-row__controls">
              <a-checkbox
                :checked="selectedBatchSet.has(entry.index)"
                @click.stop
                @change="emit('toggle-batch', entry.index)"
              />
              <a-button type="link" size="small" @click.stop="emit('edit-model', entry.index)">
                编辑
              </a-button>
            </div>
          </div>

          <div class="saved-model-row__tags">
            <a-tag
              v-for="tag in getModelCapabilityTags(entry.model)"
              :key="`${entry.index}-${tag}`"
              :color="capabilityTagMetaMap[tag].color"
            >
              {{ capabilityTagMetaMap[tag].label }}
            </a-tag>
            <span v-if="getModelCapabilityTags(entry.model).length === 0" class="saved-model-row__placeholder">
              未识别能力
            </span>
          </div>

          <div class="saved-model-row__meta">
            <span>联网 {{ formatSearchStrategyLabel(entry.model.search_strategy) }}</span>
            <span>温度 {{ formatTemperature(entry.model.temperature) }}</span>
            <span>上下文 {{ formatContextWindow(entry.model.context_window) }}</span>
            <span>排序 {{ entry.index + 1 }}</span>
          </div>

          <div v-if="entry.model.description" class="saved-model-row__desc">
            {{ entry.model.description }}
          </div>
        </article>
      </div>
    </a-collapse-panel>
  </a-collapse>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
  capabilityTagMetaMap,
  formatSearchStrategyLabel,
  getModelCapabilityTags,
  groupIndexedModelsByDisplayGroup,
  type IndexedAIModelEntry,
} from './state'

const props = defineProps<{
  entries: IndexedAIModelEntry[]
  activeModelIndex: number
  selectedBatchIndices: number[]
  keyword: string
}>()

const emit = defineEmits<{
  (e: 'select-model', index: number): void
  (e: 'edit-model', index: number): void
  (e: 'toggle-batch', index: number): void
  (e: 'select-group', indices: number[]): void
  (e: 'clear-group', indices: number[]): void
}>()

const activeGroupKeys = ref<string[]>([])
const selectedBatchSet = computed(() => new Set(props.selectedBatchIndices))
const groupedEntries = computed(() => groupIndexedModelsByDisplayGroup(props.entries))

watch(
  [groupedEntries, () => props.keyword, () => props.activeModelIndex, () => props.selectedBatchIndices],
  ([groups, keyword]) => {
    if (groups.length === 0) {
      activeGroupKeys.value = []
      return
    }
    if (String(keyword ?? '').trim()) {
      activeGroupKeys.value = groups.map(group => group.name)
      return
    }

    const focusKeys = groups
      .filter(group => group.items.some(item => (
        item.index === props.activeModelIndex || selectedBatchSet.value.has(item.index)
      )))
      .map(group => group.name)

    activeGroupKeys.value = focusKeys.length > 0 ? focusKeys : [groups[0].name]
  },
  { immediate: true, deep: true },
)

const getGroupSelectedCount = (items: IndexedAIModelEntry[]) => (
  items.filter(item => selectedBatchSet.value.has(item.index)).length
)

const getGroupSelectableCount = (items: IndexedAIModelEntry[]) => (
  items.filter(item => !selectedBatchSet.value.has(item.index)).length
)

const formatTemperature = (value: number | null) => (typeof value === 'number' ? value.toFixed(1) : '-')
const formatContextWindow = (value: number | null) => (typeof value === 'number' && value > 0 ? `${value}` : '-')
</script>

<style scoped>
.saved-group-collapse {
  background: transparent;
}

.saved-group-collapse :deep(.ant-collapse-item) {
  margin-bottom: 12px;
  overflow: hidden;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 18px;
}

.saved-group-collapse :deep(.ant-collapse-item:last-child) {
  margin-bottom: 0;
}

.saved-group-collapse :deep(.ant-collapse-header) {
  align-items: center !important;
  padding: 14px 18px !important;
}

.saved-group-collapse :deep(.ant-collapse-header-text) {
  font-size: 14px;
  font-weight: 600;
  color: var(--app-text-strong);
}

.saved-group-collapse :deep(.ant-collapse-content-box) {
  padding: 0 !important;
}

.saved-group-extra {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.saved-group-count {
  min-width: 24px;
  padding: 0 8px;
  font-size: 12px;
  line-height: 22px;
  color: var(--app-primary-color);
  text-align: center;
  background: var(--app-primary-color-soft);
  border-radius: 999px;
}

.saved-group-meta,
.saved-model-row__key,
.saved-model-row__meta,
.saved-model-row__desc,
.saved-model-row__placeholder {
  color: var(--app-text-secondary);
}

.saved-group-action {
  padding-inline: 0;
}

.saved-model-list {
  display: flex;
  flex-direction: column;
}

.saved-model-row {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px 18px;
  background: var(--app-surface-color);
  border-top: 1px solid var(--app-border-color);
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    box-shadow 0.2s ease,
    border-color 0.2s ease;
}

.saved-model-row:first-child {
  border-top: 0;
}

.saved-model-row:hover {
  background: var(--app-hover-bg);
}

.saved-model-row--active {
  background: var(--app-primary-color-soft);
  box-shadow: inset 3px 0 0 var(--app-primary-color);
}

.saved-model-row__head {
  display: flex;
  gap: 16px;
  align-items: flex-start;
  justify-content: space-between;
}

.saved-model-row__title-wrap {
  min-width: 0;
  flex: 1;
}

.saved-model-row__title {
  font-weight: 600;
  color: var(--app-text-strong);
}

.saved-model-row__controls {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-shrink: 0;
}

.saved-model-row__key,
.saved-model-row__meta,
.saved-model-row__desc {
  margin-top: 6px;
  font-size: 12px;
  word-break: break-all;
}

.saved-model-row__tags,
.saved-model-row__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

@media (max-width: 768px) {
  .saved-model-row__head {
    flex-direction: column;
    align-items: stretch;
  }

  .saved-model-row__controls {
    justify-content: space-between;
  }
}
</style>
