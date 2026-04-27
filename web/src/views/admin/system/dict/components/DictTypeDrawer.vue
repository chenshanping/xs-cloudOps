<template>
  <a-drawer
    :open="open"
    :title="title"
    width="520"
    placement="right"
    @close="handleClose"
  >
    <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 5 }">
      <a-form-item label="字典名称" name="name">
        <a-input v-model:value="formState.name" placeholder="请输入字典名称" />
      </a-form-item>
      <a-form-item label="字典类型" name="type">
        <a-input
          v-model:value="formState.type"
          :disabled="isEdit"
          placeholder="请输入字典类型（英文与下划线）"
        />
      </a-form-item>
      <a-form-item label="状态" name="statusChecked">
        <a-switch
          v-model:checked="formState.statusChecked"
          checked-children="正常"
          un-checked-children="停用"
        />
      </a-form-item>
      <a-form-item label="备注" name="remark">
        <a-textarea v-model:value="formState.remark" :rows="4" placeholder="请输入备注" />
      </a-form-item>
    </a-form>

    <template #footer>
      <a-space>
        <a-button @click="handleClose">取消</a-button>
        <a-button type="primary" :loading="submitting" @click="handleSubmit">保存</a-button>
      </a-space>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'

interface DictTypeFormState {
  name: string
  type: string
  statusChecked: boolean
  remark: string
}

const props = defineProps<{
  open: boolean
  title: string
  isEdit: boolean
  submitting: boolean
  initialValue?: Partial<{
    name: string
    type: string
    status: number
    remark: string
  }>
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: { name: string; type: string; status: number; remark: string }): void
}>()

const formRef = ref<FormInstance>()

const formState = reactive<DictTypeFormState>({
  name: '',
  type: '',
  statusChecked: true,
  remark: '',
})

const formRules: Record<string, Rule[]> = {
  name: [{ required: true, message: '请输入字典名称', trigger: 'blur' }],
  type: [
    { required: true, message: '请输入字典类型', trigger: 'blur' },
    {
      pattern: /^[A-Za-z][A-Za-z0-9_]*$/,
      message: '字典类型需以字母开头，只能包含字母、数字和下划线',
      trigger: 'blur',
    },
  ],
}

watch(
  () => [props.open, props.initialValue],
  () => {
    if (!props.open) {
      return
    }

    Object.assign(formState, {
      name: props.initialValue?.name ?? '',
      type: props.initialValue?.type ?? '',
      statusChecked: (props.initialValue?.status ?? 1) === 1,
      remark: props.initialValue?.remark ?? '',
    })
  },
  { immediate: true, deep: true },
)

const handleClose = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  emit('submit', {
    name: formState.name.trim(),
    type: formState.type.trim(),
    status: formState.statusChecked ? 1 : 0,
    remark: formState.remark.trim(),
  })
}
</script>
