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
      <a-form-item label="数据范围" name="data_scope">
        <a-select v-model:value="formState.data_scope" placeholder="请选择数据范围">
          <a-select-option :value="1">全部数据</a-select-option>
          <a-select-option :value="2">自定义部门</a-select-option>
          <a-select-option :value="3">本部门</a-select-option>
          <a-select-option :value="4">本部门及下级</a-select-option>
          <a-select-option :value="5">仅本人</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item v-if="formState.data_scope === 2" label="自定义部门" name="dept_ids">
        <a-tree-select
          v-model:value="formState.dept_ids"
          :tree-data="deptOptions"
          placeholder="请选择部门"
          tree-default-expand-all
          tree-checkable
          multiple
        />
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

interface TreeSelectOption {
  key: string | number
  title: string
  value: number
  disabled?: boolean
  children?: TreeSelectOption[]
}

interface RoleFormValue {
  name: string
  code: string
  sort: number
  statusChecked: boolean
  data_scope: number
  dept_ids: number[]
  remark: string
}

const props = defineProps<{
  open: boolean
  title: string
  isEdit: boolean
  deptOptions: TreeSelectOption[]
  initialValue?: Partial<RoleFormValue>
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: {
    name: string
    code: string
    sort: number
    status: number
    data_scope: number
    dept_ids: number[]
    remark: string
  }): void
}>()

const formRef = ref<FormInstance>()

const formState = reactive<RoleFormValue>({
  name: '',
  code: '',
  sort: 0,
  statusChecked: true,
  data_scope: 1,
  dept_ids: [],
  remark: ''
})

const formRules = computed<Record<string, Rule[]>>(() => ({
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }],
  data_scope: [{ required: true, message: '请选择数据范围', trigger: 'change' }],
  dept_ids: formState.data_scope === 2 ? [{ required: true, message: '请选择自定义部门', trigger: 'change' }] : []
}))

watch(
  () => formState.data_scope,
  value => {
    if (value !== 2) {
      formState.dept_ids = []
    }
  }
)

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
      data_scope: props.initialValue?.data_scope ?? 1,
      dept_ids: props.initialValue?.dept_ids ? [...props.initialValue.dept_ids] : [],
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
    data_scope: formState.data_scope,
    dept_ids: formState.data_scope === 2 ? [...formState.dept_ids] : [],
    remark: formState.remark
  })
}
</script>
