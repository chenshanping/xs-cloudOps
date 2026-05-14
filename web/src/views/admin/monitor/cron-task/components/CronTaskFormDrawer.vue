<template>
  <a-drawer :open="open" :title="title" width="640" placement="right" @close="handleClose">
    <a-form ref="formRef" :model="form" :rules="rules" layout="vertical">
      <a-card size="small" class="form-section" title="基础信息">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="实例编码" name="code">
              <a-input v-model:value="form.code" placeholder="如 cleanup_login_logs_default" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="任务名称" name="name">
              <a-input v-model:value="form.name" placeholder="请输入任务名称" allow-clear />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="注册任务" name="task_code">
          <a-select v-model:value="form.task_code" placeholder="请选择注册任务" @change="handleTaskCodeChange">
            <a-select-option v-for="item in registry" :key="item.code" :value="item.code">
              {{ item.name }}（{{ item.code }}）
            </a-select-option>
          </a-select>
        </a-form-item>
        <div v-if="selectedRegisteredTask" class="registry-desc">
          {{ selectedRegisteredTask.description }}
        </div>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Cron 表达式" name="cron_expr">
              <a-input v-model:value="form.cron_expr" placeholder="如 0 2 * * *" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="排序" name="sort">
              <a-input-number v-model:value="form.sort" style="width: 100%" :min="0" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="备注" name="remark">
          <a-textarea v-model:value="form.remark" :rows="3" placeholder="请输入备注" allow-clear />
        </a-form-item>
      </a-card>
      <a-card size="small" class="form-section" title="任务参数">
        <CronTaskParamsForm v-model="form.params" :schema="selectedRegisteredTask?.param_schema || {}" />
      </a-card>
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
import { computed, reactive, ref, watch } from 'vue'
import type { FormInstance } from 'ant-design-vue'
import type { CronTask, CronTaskPayload, RegisteredCronTask } from '@/api/cron'
import CronTaskParamsForm from './CronTaskParamsForm.vue'

const props = defineProps<{
  open: boolean
  title: string
  isEdit: boolean
  record: CronTask | null
  registry: RegisteredCronTask[]
  submitting?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: CronTaskPayload): void
}>()

const formRef = ref<FormInstance>()
const form = reactive<CronTaskPayload>({
  code: '',
  task_code: '',
  name: '',
  cron_expr: '0 2 * * *',
  params: {},
  remark: '',
  sort: 0,
})

const rules = {
  code: [{ required: true, message: '请输入实例编码', trigger: 'blur' }],
  task_code: [{ required: true, message: '请选择注册任务', trigger: 'change' }],
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  cron_expr: [{ required: true, message: '请输入 Cron 表达式', trigger: 'blur' }],
}

const selectedRegisteredTask = computed(() => props.registry.find((item) => item.code === form.task_code))

const defaultParamsFromSchema = (task?: RegisteredCronTask) => {
  const params: Record<string, any> = {}
  Object.entries(task?.param_schema || {}).forEach(([key, field]) => {
    params[key] = field.default ?? (field.type === 'int' ? 0 : field.type === 'bool' ? false : '')
  })
  return params
}

const resetForm = () => {
  const firstTask = props.registry[0]
  form.code = ''
  form.task_code = firstTask?.code || ''
  form.name = firstTask?.name || ''
  form.cron_expr = '0 2 * * *'
  form.params = defaultParamsFromSchema(firstTask)
  form.remark = ''
  form.sort = 0
  formRef.value?.clearValidate()
}

const fillForm = (record: CronTask) => {
  form.code = record.code
  form.task_code = record.task_code
  form.name = record.name
  form.cron_expr = record.cron_expr
  form.params = { ...(record.params || {}) }
  form.remark = record.remark || ''
  form.sort = record.sort || 0
  formRef.value?.clearValidate()
}

const handleTaskCodeChange = () => {
  const task = selectedRegisteredTask.value
  if (!props.isEdit && task && !form.name) {
    form.name = task.name
  }
  form.params = defaultParamsFromSchema(task)
}

const handleClose = () => {
  emit('update:open', false)
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  emit('submit', { ...form, params: { ...form.params } })
}

watch(() => props.open, (open) => {
  if (!open) return
  if (props.isEdit && props.record) {
    fillForm(props.record)
  } else {
    resetForm()
  }
})

watch(() => props.registry, () => {
  if (props.open && !props.isEdit) resetForm()
}, { deep: true })
</script>

<style scoped>
.form-section + .form-section {
  margin-top: 16px;
}
.registry-desc {
  margin: -8px 0 16px;
  padding: 10px 12px;
  color: var(--app-text-muted);
  font-size: 13px;
  background: var(--app-bg-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}
</style>
