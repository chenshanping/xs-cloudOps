<template>
  <a-drawer
    :open="open"
    :title="title"
    width="520"
    placement="right"
    @close="handleClose"
  >
    <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 6 }">
      <a-form-item label="上级部门" name="parent_id">
        <a-tree-select
          v-model:value="formState.parent_id"
          :tree-data="treeOptions"
          :field-names="{ label: 'name', value: 'id', children: 'children' }"
          placeholder="请选择上级部门"
          tree-default-expand-all
        />
      </a-form-item>
      <a-form-item label="部门名称" name="name">
        <a-input v-model:value="formState.name" placeholder="请输入部门名称" />
      </a-form-item>
      <a-form-item label="排序" name="sort">
        <a-input-number v-model:value="formState.sort" :min="0" style="width: 100%" />
      </a-form-item>
      <a-form-item label="状态" name="statusChecked">
        <a-switch v-model:checked="formState.statusChecked" />
      </a-form-item>
      <a-form-item label="备注" name="remark">
        <a-textarea v-model:value="formState.remark" :rows="4" placeholder="请输入备注" />
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
import { reactive, ref, watch } from 'vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import type { Dept } from '@/types'

interface DeptFormValue {
  parent_id: number
  name: string
  sort: number
  statusChecked: boolean
  remark: string
}

const props = defineProps<{
  open: boolean
  title: string
  treeOptions: Dept[]
  initialValue?: Partial<DeptFormValue>
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: { parent_id: number; name: string; sort: number; status: number; remark: string }): void
}>()

const formRef = ref<FormInstance>()

const formState = reactive<DeptFormValue>({
  parent_id: 0,
  name: '',
  sort: 0,
  statusChecked: true,
  remark: ''
})

const formRules: Record<string, Rule[]> = {
  name: [{ required: true, message: '请输入部门名称', trigger: 'blur' }]
}

watch(
  () => [props.open, props.initialValue],
  () => {
    if (!props.open) {
      return
    }
    Object.assign(formState, {
      parent_id: props.initialValue?.parent_id ?? 0,
      name: props.initialValue?.name ?? '',
      sort: props.initialValue?.sort ?? 0,
      statusChecked: props.initialValue?.statusChecked ?? true,
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
    parent_id: formState.parent_id || 0,
    name: formState.name,
    sort: formState.sort || 0,
    status: formState.statusChecked ? 1 : 0,
    remark: formState.remark
  })
}
</script>
