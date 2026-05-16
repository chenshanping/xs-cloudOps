<template>
  <a-drawer
    :open="open"
    :title="isEdit ? '编辑模型' : '新增模型'"
    width="720"
    placement="right"
    :mask-closable="false"
    @close="handleClose"
  >
    <div class="editor-shell">
      <div class="editor-summary">
        <div>
          <div class="editor-summary__title">{{ providerName || '未命名平台' }}</div>
          <div class="editor-summary__subtitle">编辑当前平台下的模型分组、能力、参数和描述</div>
        </div>
        <div class="editor-summary__tags">
          <a-tag
            v-for="tag in capabilityTags"
            :key="tag"
            :color="capabilityTagMetaMap[tag].color"
          >
            {{ capabilityTagMetaMap[tag].label }}
          </a-tag>
          <span v-if="capabilityTags.length === 0" class="editor-summary__placeholder">尚未标记能力</span>
        </div>
      </div>

      <a-form ref="formRef" :model="formState" layout="vertical" class="editor-form">
        <div class="editor-grid editor-grid--basic">
          <a-form-item
            label="模型 ID"
            name="id"
            :rules="[{ required: true, message: '请输入模型 ID' }]"
          >
            <a-input v-model:value="formState.id" placeholder="例如：qwen3.5-plus-2026-04-20" />
          </a-form-item>
          <a-form-item label="显示名称" name="name">
            <a-input v-model:value="formState.name" placeholder="默认回退为模型 ID" />
          </a-form-item>
          <a-form-item label="模型分组" name="group">
            <a-input
              v-model:value="formState.group"
              placeholder="例如：deepseek-v4、mimo-v2.5"
              @update:value="handleGroupInput"
            />
          </a-form-item>
        </div>

        <div class="editor-grid editor-grid--switches">
          <a-form-item label="推理模型">
            <a-switch v-model:checked="formState.is_thinking" />
          </a-form-item>
          <a-form-item label="视觉能力">
            <a-switch v-model:checked="formState.support_vision" />
          </a-form-item>
          <a-form-item label="工具调用">
            <a-switch v-model:checked="formState.support_tools" />
          </a-form-item>
          <a-form-item label="嵌入模型">
            <a-switch v-model:checked="formState.support_embedding" />
          </a-form-item>
          <a-form-item label="重排模型">
            <a-switch v-model:checked="formState.support_rerank" />
          </a-form-item>
          <a-form-item label="免费模型">
            <a-switch v-model:checked="formState.is_free" />
          </a-form-item>
        </div>

        <div class="editor-grid editor-grid--params">
          <a-form-item label="联网策略">
            <a-select
              v-model:value="formState.search_strategy"
              :options="searchStrategyOptions"
              placeholder="选择联网策略"
            />
          </a-form-item>
          <a-form-item label="温度">
            <a-input-number
              v-model:value="formState.temperature"
              :min="0"
              :max="2"
              :step="0.1"
              style="width: 100%"
              placeholder="例如 0.7"
            />
          </a-form-item>
          <a-form-item label="上下文数">
            <a-input-number
              v-model:value="formState.context_window"
              :min="0"
              :step="1024"
              style="width: 100%"
              placeholder="例如 128000"
            />
          </a-form-item>
        </div>

        <a-form-item label="描述">
          <a-textarea
            v-model:value="formState.description"
            :auto-size="{ minRows: 6, maxRows: 12 }"
            placeholder="补充模型用途、能力边界或推荐场景"
          />
        </a-form-item>
      </a-form>
    </div>

    <template #footer>
      <div class="editor-footer">
        <a-button v-if="isEdit && canDelete" danger @click="emit('remove')">删除模型</a-button>
        <a-space>
          <a-button @click="handleClose">取消</a-button>
          <a-button :disabled="!canTest" @click="handleTest">测试模型</a-button>
          <a-button type="primary" :disabled="!canSubmit" @click="handleSubmit">确认</a-button>
        </a-space>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import type { FormInstance } from 'ant-design-vue'
import {
  capabilityTagMetaMap,
  createEmptyModel,
  getRemoteModelGroupName,
  getModelCapabilityTags,
  isMeaningfulExplicitModelGroup,
  normalizeModel,
  searchStrategyOptions,
  type AIModel,
} from './state'

const props = defineProps<{
  open: boolean
  isEdit: boolean
  model?: Partial<AIModel> | null
  providerName: string
  canDelete: boolean
  canSubmit: boolean
  canTest: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: AIModel): void
  (e: 'remove'): void
  (e: 'test', value: AIModel): void
}>()

const formRef = ref<FormInstance>()
const formState = reactive<AIModel>(createEmptyModel())
const capabilityTags = computed(() => getModelCapabilityTags(formState))
const groupManualOverride = ref(false)
const lastAutoDerivedGroup = ref('')

const syncDerivedGroup = (force = false) => {
  const nextGroup = getRemoteModelGroupName({
    id: formState.id,
    name: formState.name,
  })
  const shouldReplaceCurrentGroup = force
    || !isMeaningfulExplicitModelGroup(formState.group)
    || !groupManualOverride.value
    || formState.group.trim() === lastAutoDerivedGroup.value

  if (!shouldReplaceCurrentGroup) {
    return
  }

  formState.group = nextGroup
  lastAutoDerivedGroup.value = nextGroup
}

watch(
  () => [props.open, props.model],
  () => {
    if (!props.open) {
      return
    }
    const normalized = normalizeModel(props.model)
    Object.assign(formState, normalized)
    const originalGroup = String(props.model?.group ?? '').trim()
    groupManualOverride.value = isMeaningfulExplicitModelGroup(originalGroup)
    lastAutoDerivedGroup.value = groupManualOverride.value ? '' : normalized.group
    syncDerivedGroup(!groupManualOverride.value)
  },
  { immediate: true, deep: true },
)

watch(
  () => [formState.id, formState.name],
  () => {
    if (!props.open) {
      return
    }
    syncDerivedGroup()
  },
)

const handleGroupInput = (value: string) => {
  const normalizedValue = value.trim()
  if (isMeaningfulExplicitModelGroup(normalizedValue) && normalizedValue !== lastAutoDerivedGroup.value) {
    groupManualOverride.value = true
    return
  }

  groupManualOverride.value = false
  syncDerivedGroup(true)
}

const handleClose = () => {
  emit('update:open', false)
}

const buildPayload = () => normalizeModel(formState)

const handleSubmit = async () => {
  if (!props.canSubmit) {
    return
  }
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  emit('submit', buildPayload())
}

const handleTest = async () => {
  if (!props.canTest) {
    return
  }
  try {
    await formRef.value?.validateFields(['id'])
  } catch {
    return
  }
  emit('test', buildPayload())
}
</script>

<style scoped>
.editor-shell {
  display: flex;
  flex-direction: column;
  gap: 18px;
  color: var(--app-text-color);
}

.editor-summary {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px 18px;
  border: 1px solid var(--app-border-color);
  border-radius: 18px;
  background: linear-gradient(180deg, var(--app-surface-color) 0%, var(--app-surface-soft) 100%);
}

.editor-summary__title {
  font-size: 16px;
  font-weight: 600;
  color: var(--app-text-strong);
}

.editor-summary__subtitle,
.editor-summary__placeholder {
  margin-top: 4px;
  color: var(--app-text-secondary);
}

.editor-summary__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.editor-form {
  min-width: 0;
}

.editor-grid {
  display: grid;
  gap: 12px;
}

.editor-grid--basic,
.editor-grid--params {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.editor-grid--switches {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.editor-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

@media (max-width: 960px) {
  .editor-grid--basic,
  .editor-grid--params,
  .editor-grid--switches {
    grid-template-columns: 1fr;
  }

  .editor-footer {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
