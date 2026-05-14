<template>
  <div class="params-form">
    <a-empty v-if="schemaEntries.length === 0" description="该任务暂无参数" :image="simpleImage" />
    <a-row v-else :gutter="16">
      <a-col v-for="[key, field] in schemaEntries" :key="key" :span="12">
        <a-form-item :label="field.description || key" :required="field.required">
          <a-input-number
            v-if="field.type === 'int'"
            :value="localValue[key]"
            :min="field.min"
            :max="field.max"
            style="width: 100%"
            @change="(value) => updateField(key, value)"
          />
          <a-switch
            v-else-if="field.type === 'bool'"
            :checked="Boolean(localValue[key])"
            @change="(checked) => updateField(key, checked)"
          />
          <a-input
            v-else
            :value="localValue[key]"
            allow-clear
            @update:value="(value) => updateField(key, value)"
          />
          <div class="param-meta">
            <span>{{ key }}</span>
            <span v-if="field.type === 'int' && (field.min !== undefined || field.max !== undefined)">
              {{ field.min ?? '-' }} ~ {{ field.max ?? '-' }}
            </span>
          </div>
        </a-form-item>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { Empty } from 'ant-design-vue'
import type { CronParamDefinition } from '@/api/cron'

const props = defineProps<{
  modelValue: Record<string, any>
  schema: Record<string, CronParamDefinition>
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Record<string, any>): void
}>()

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const localValue = reactive<Record<string, any>>({})
const schemaEntries = computed(() => Object.entries(props.schema || {}))

const applyValue = () => {
  Object.keys(localValue).forEach((key) => delete localValue[key])
  schemaEntries.value.forEach(([key, field]) => {
    localValue[key] = props.modelValue?.[key] ?? field.default ?? defaultValue(field.type)
  })
}

const defaultValue = (type: string) => {
  if (type === 'int') return 0
  if (type === 'bool') return false
  return ''
}

const updateField = (key: string, value: any) => {
  localValue[key] = value
  emit('update:modelValue', { ...localValue })
}

watch(() => [props.modelValue, props.schema], applyValue, { immediate: true, deep: true })
</script>

<style scoped>
.params-form {
  padding: 12px;
  background: var(--app-bg-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}
.param-meta {
  display: flex;
  justify-content: space-between;
  margin-top: 4px;
  color: var(--app-text-muted);
  font-size: 12px;
  line-height: 1.4;
}
</style>
