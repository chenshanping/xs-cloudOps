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
        <div class="drawer-section__title">基础信息</div>
        <a-row :gutter="16">
          <a-col :span="16">
            <a-form-item label="分组名称" name="name">
              <a-input v-model:value="formState.name" placeholder="请输入分组名称" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="排序" name="sort">
              <a-input-number v-model:value="formState.sort" style="width: 100%" :min="0" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="备注" name="remark">
          <a-textarea v-model:value="formState.remark" :rows="4" :maxlength="500" show-count placeholder="请输入备注" />
        </a-form-item>
      </div>

      <div class="drawer-section">
        <div class="drawer-section__title">状态设置</div>
        <div class="status-panel">
          <div>
            <div class="status-panel__label">启用状态</div>
            <div class="status-panel__hint">停用后不影响已有主机关联，仅表示该分组不再推荐使用。</div>
          </div>
          <a-switch v-model:checked="formState.statusChecked" checked-children="启用" un-checked-children="停用" />
        </div>
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
import type { CmdbHostGroupPayload } from '@/api/cmdb'

interface GroupFormValue {
  name: string
  sort: number
  remark: string
  statusChecked: boolean
}

const props = defineProps<{
  open: boolean
  title: string
  initialValue?: Partial<GroupFormValue>
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: CmdbHostGroupPayload): void
}>()

const formRef = ref<FormInstance>()
const formState = reactive<GroupFormValue>({
  name: '',
  sort: 0,
  remark: '',
  statusChecked: true,
})

const formRules: Record<string, Rule[]> = {
  name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }],
}

watch(
  () => [props.open, props.initialValue] as const,
  ([open]) => {
    if (!open) {
      return
    }
    Object.assign(formState, {
      name: props.initialValue?.name ?? '',
      sort: props.initialValue?.sort ?? 0,
      remark: props.initialValue?.remark ?? '',
      statusChecked: props.initialValue?.statusChecked ?? true,
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
    sort: formState.sort || 0,
    remark: formState.remark.trim(),
    status: formState.statusChecked ? 1 : 0,
  })
}
</script>

<style scoped>
.drawer-section + .drawer-section {
  margin-top: 16px;
}

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

.status-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.status-panel__label {
  color: var(--app-text-strong);
  font-size: 13px;
  font-weight: 600;
}

.status-panel__hint {
  margin-top: 4px;
  color: var(--app-text-muted);
  font-size: 12px;
  line-height: 1.6;
}
</style>
