<template>
  <a-drawer
    :open="open"
    :title="title"
    width="520"
    placement="right"
    destroy-on-close
    @close="handleClose"
  >
    <a-form ref="formRef" :model="formState" :rules="formRules" layout="vertical">
      <div class="drawer-section">
        <div class="drawer-section__title">标签信息</div>
        <a-form-item label="标签名称" name="name">
          <a-input v-model:value="formState.name" placeholder="请输入标签名称" allow-clear />
        </a-form-item>
        <a-form-item label="标签颜色" name="color">
          <div class="color-row">
            <input v-model="formState.color" class="color-picker" type="color" />
            <a-input v-model:value="formState.color" placeholder="如 #1677ff" allow-clear />
            <a-tag :color="formState.color || '#1677ff'">预览</a-tag>
          </div>
        </a-form-item>
        <a-form-item label="备注" name="remark">
          <a-textarea v-model:value="formState.remark" :rows="4" :maxlength="500" show-count placeholder="请输入备注" />
        </a-form-item>
      </div>
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
import type { CmdbHostTagPayload } from '@/api/cmdb'

interface TagFormValue {
  name: string
  color: string
  remark: string
}

const props = defineProps<{
  open: boolean
  title: string
  initialValue?: Partial<TagFormValue>
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: CmdbHostTagPayload): void
}>()

const formRef = ref<FormInstance>()
const formState = reactive<TagFormValue>({
  name: '',
  color: '#1677ff',
  remark: '',
})

const formRules: Record<string, Rule[]> = {
  name: [{ required: true, message: '请输入标签名称', trigger: 'blur' }],
}

watch(
  () => [props.open, props.initialValue] as const,
  ([open]) => {
    if (!open) {
      return
    }
    Object.assign(formState, {
      name: props.initialValue?.name ?? '',
      color: props.initialValue?.color || '#1677ff',
      remark: props.initialValue?.remark ?? '',
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
    name: formState.name.trim(),
    color: formState.color.trim() || '#1677ff',
    remark: formState.remark.trim(),
  })
}
</script>

<style scoped>
.drawer-section {
  padding: 16px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.drawer-section__title {
  margin-bottom: 14px;
  color: var(--app-text-strong);
  font-size: 14px;
  font-weight: 600;
}

.color-row {
  display: grid;
  grid-template-columns: 52px 1fr auto;
  gap: 10px;
  align-items: center;
}

.color-picker {
  width: 52px;
  height: 34px;
  padding: 0;
  cursor: pointer;
  background: transparent;
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}
</style>
