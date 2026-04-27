<template>
  <a-drawer
    :open="open"
    :title="title"
    width="560"
    placement="right"
    @close="handleClose"
  >
    <a-alert
      type="info"
      show-icon
      style="margin-bottom: 16px"
      :message="`当前字典类型：${currentTypeName}`"
      :description="currentTypeCode"
    />

    <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 5 }">
      <a-form-item label="字典标签" name="label">
        <a-input v-model:value="formState.label" placeholder="请输入字典标签（显示名称）" />
      </a-form-item>
      <a-form-item label="字典键值" name="value">
        <a-input v-model:value="formState.value" placeholder="请输入字典键值" />
      </a-form-item>
      <a-form-item label="排序" name="sort">
        <a-input-number v-model:value="formState.sort" :min="0" style="width: 100%" />
      </a-form-item>
      <a-form-item label="标签颜色" name="tag_type">
        <a-select v-model:value="formState.tag_type" allow-clear placeholder="请选择标签颜色">
          <a-select-option v-for="option in colorOptions" :key="option.value" :value="option.value">
            <a-tag :color="option.value || 'default'">{{ option.label }}</a-tag>
          </a-select-option>
        </a-select>
        <div class="dict-data-drawer__preview">
          <span class="dict-data-drawer__preview-label">预览</span>
          <a-tag :color="formState.tag_type || 'default'">{{ formState.label || '字典标签' }}</a-tag>
        </div>
      </a-form-item>
      <a-form-item label="状态" name="statusChecked">
        <a-switch
          v-model:checked="formState.statusChecked"
          checked-children="正常"
          un-checked-children="停用"
        />
      </a-form-item>
      <a-form-item label="是否默认" name="isDefaultChecked">
        <a-switch
          v-model:checked="formState.isDefaultChecked"
          checked-children="是"
          un-checked-children="否"
        />
      </a-form-item>
      <a-form-item label="备注" name="remark">
        <a-textarea v-model:value="formState.remark" :rows="3" placeholder="请输入备注" />
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

interface DictDataFormState {
  label: string
  value: string
  sort: number
  tag_type: string
  statusChecked: boolean
  isDefaultChecked: boolean
  remark: string
}

const colorOptions = [
  { label: '默认', value: '' },
  { label: '成功/绿色', value: 'success' },
  { label: '处理中/蓝色', value: 'processing' },
  { label: '警告/橙色', value: 'warning' },
  { label: '错误/红色', value: 'error' },
  { label: '灰色', value: 'default' },
  { label: '粉色', value: 'pink' },
  { label: '紫色', value: 'purple' },
  { label: '青色', value: 'cyan' },
]

const props = defineProps<{
  open: boolean
  title: string
  submitting: boolean
  currentTypeName: string
  currentTypeCode: string
  initialValue?: Partial<{
    label: string
    value: string
    sort: number
    tag_type: string
    status: number
    is_default: number
    remark: string
  }>
}>()

const emit = defineEmits<{
  (
    e: 'submit',
    value: {
      label: string
      value: string
      sort: number
      tag_type: string
      status: number
      is_default: number
      remark: string
    },
  ): void
  (e: 'update:open', value: boolean): void
}>()

const formRef = ref<FormInstance>()

const formState = reactive<DictDataFormState>({
  label: '',
  value: '',
  sort: 0,
  tag_type: '',
  statusChecked: true,
  isDefaultChecked: false,
  remark: '',
})

const formRules: Record<string, Rule[]> = {
  label: [{ required: true, message: '请输入字典标签', trigger: 'blur' }],
  value: [{ required: true, message: '请输入字典键值', trigger: 'blur' }],
}

watch(
  () => [props.open, props.initialValue],
  () => {
    if (!props.open) {
      return
    }

    Object.assign(formState, {
      label: props.initialValue?.label ?? '',
      value: props.initialValue?.value ?? '',
      sort: props.initialValue?.sort ?? 0,
      tag_type: props.initialValue?.tag_type ?? '',
      statusChecked: (props.initialValue?.status ?? 1) === 1,
      isDefaultChecked: (props.initialValue?.is_default ?? 0) === 1,
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
    label: formState.label.trim(),
    value: formState.value.trim(),
    sort: formState.sort || 0,
    tag_type: formState.tag_type,
    status: formState.statusChecked ? 1 : 0,
    is_default: formState.isDefaultChecked ? 1 : 0,
    remark: formState.remark.trim(),
  })
}
</script>

<style scoped>
.dict-data-drawer__preview {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.dict-data-drawer__preview-label {
  color: #8c8c8c;
  font-size: 12px;
}
</style>
