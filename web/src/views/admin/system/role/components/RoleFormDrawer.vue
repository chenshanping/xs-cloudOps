<template>
  <a-drawer
    :open="open"
    :title="title"
    width="560"
    placement="right"
    @close="handleClose"
  >
    <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 5 }">
      <a-form-item label="角色名称" name="name">
        <a-input v-model:value="formState.name" />
      </a-form-item>
      <a-form-item label="角色编码" name="code">
        <a-input v-model:value="formState.code" :disabled="isEdit" />
      </a-form-item>
      <a-form-item label="排序">
        <a-input-number v-model:value="formState.sort" :min="0" style="width: 100%" />
      </a-form-item>
      <a-form-item label="状态">
        <a-switch v-model:checked="formState.statusChecked" />
      </a-form-item>
      <a-form-item label="超管角色">
        <a-switch v-model:checked="formState.is_super_admin" />
      </a-form-item>
      <a-form-item label="备注">
        <a-textarea v-model:value="formState.remark" />
      </a-form-item>
    </a-form>

    <template #footer>
      <a-space>
        <a-button @click="handleClose">取消</a-button>
        <a-button type="primary" @click="handleSubmit">保存</a-button>
      </a-space>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'

interface RoleFormValue {
  name: string
  code: string
  sort: number
  statusChecked: boolean
  is_super_admin: boolean
  remark: string
}

const props = defineProps<{
  open: boolean
  title: string
  isEdit: boolean
  initialValue?: Partial<RoleFormValue>
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: {
    name: string
    code: string
    sort: number
    status: number
    is_super_admin: boolean
    remark: string
  }): void
}>()

const formRef = ref<FormInstance>()

const formState = reactive<RoleFormValue>({
  name: '',
  code: '',
  sort: 0,
  statusChecked: true,
  is_super_admin: false,
  remark: ''
})

const formRules = computed<Record<string, Rule[]>>(() => ({
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }],
}))

watch(
  () => [props.open, props.initialValue],
  () => {
    if (!props.open) {
      return
    }
    Object.assign(formState, {
      name: props.initialValue?.name ?? '',
      code: props.initialValue?.code ?? '',
      sort: props.initialValue?.sort ?? 0,
      statusChecked: props.initialValue?.statusChecked ?? true,
      is_super_admin: props.initialValue?.is_super_admin ?? false,
      remark: props.initialValue?.remark ?? ''
    })
  },
  { immediate: true, deep: true }
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
    name: formState.name,
    code: formState.code,
    sort: formState.sort || 0,
    status: formState.statusChecked ? 1 : 0,
    is_super_admin: formState.is_super_admin,
    remark: formState.remark
  })
}
</script>
