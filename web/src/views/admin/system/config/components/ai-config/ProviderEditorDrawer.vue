<template>
  <a-drawer
    :open="open"
    :title="isEdit ? '编辑 AI 平台' : '新增 AI 平台'"
    width="50%"
    placement="right"
    @close="handleClose"
  >
    <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 5 }">
      <a-form-item label="平台名称" name="name">
        <a-input v-model:value="formState.name" placeholder="如：阿里云百炼" />
      </a-form-item>
      <a-form-item label="API Key" name="api_key">
        <a-input-password v-model:value="formState.api_key" placeholder="请输入平台 API Key" />
      </a-form-item>
      <a-form-item label="Base URL" name="base_url">
        <a-input v-model:value="formState.base_url" placeholder="如：https://dashscope.aliyuncs.com/compatible-mode/v1" />
      </a-form-item>
      <a-form-item label="默认平台">
        <a-switch v-model:checked="formState.isDefault" />
      </a-form-item>
    </a-form>

    <template #footer>
      <div class="drawer-footer">
        <a-button v-if="isEdit && canDelete" danger :loading="submitting" :disabled="submitting" @click="emit('remove')">删除平台</a-button>
        <a-space>
          <a-button :disabled="submitting" @click="handleClose">取消</a-button>
          <a-button type="primary" :loading="submitting" :disabled="!canSubmit || submitting" @click="handleSubmit">确认</a-button>
        </a-space>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import type { AIProvider } from './state'

interface ProviderEditorSubmitValue {
  name: string
  api_key: string
  base_url: string
  isDefault: boolean
}

const props = defineProps<{
  open: boolean
  isEdit: boolean
  initialValue?: Partial<AIProvider>
  isDefault: boolean
  existingNames: string[]
  submitting: boolean
  canDelete: boolean
  canSubmit: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: ProviderEditorSubmitValue): void
  (e: 'remove'): void
}>()

const formRef = ref<FormInstance>()

const formState = reactive<ProviderEditorSubmitValue>({
  name: '',
  api_key: '',
  base_url: '',
  isDefault: false,
})

const formRules = computed<Record<string, Rule[]>>(() => ({
  name: [
    { required: true, message: '请输入平台名称', trigger: 'blur' },
    {
      trigger: 'blur',
      validator: async (_rule, value?: string) => {
        const trimmed = value?.trim() ?? ''
        if (!trimmed) {
          return
        }
        const duplicate = props.existingNames.some(name =>
          name === trimmed && name !== (props.initialValue?.name ?? ''),
        )
        if (duplicate) {
          throw new Error('平台名称不能重复')
        }
      },
    },
  ],
  base_url: [{ required: true, message: '请输入 Base URL', trigger: 'blur' }],
}))

watch(
  () => [props.open, props.initialValue, props.isDefault],
  () => {
    if (!props.open) {
      return
    }
    Object.assign(formState, {
      name: props.initialValue?.name ?? '',
      api_key: props.initialValue?.api_key ?? '',
      base_url: props.initialValue?.base_url ?? '',
      isDefault: props.isDefault,
    })
  },
  { immediate: true, deep: true },
)

const handleClose = () => {
  emit('update:open', false)
}

const handleSubmit = async () => {
  if (!props.canSubmit) {
    return
  }
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  emit('submit', {
    name: formState.name.trim(),
    api_key: formState.api_key.trim(),
    base_url: formState.base_url.trim(),
    isDefault: formState.isDefault,
  })
}
</script>

<style scoped>
.drawer-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
