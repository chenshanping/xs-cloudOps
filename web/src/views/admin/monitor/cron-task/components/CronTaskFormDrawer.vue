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
          <a-col :span="16">
            <a-form-item label="执行计划" name="cron_expr">
              <div class="schedule-builder">
                <a-radio-group v-model:value="scheduleMode" button-style="solid" size="small" @change="handleScheduleModeChange">
                  <a-radio-button value="simple">简单设置</a-radio-button>
                  <a-radio-button value="advanced">高级模式</a-radio-button>
                </a-radio-group>

                <template v-if="scheduleMode === 'simple'">
                  <div class="schedule-builder__row">
                    <a-select v-model:value="simpleSchedule.type" style="width: 110px">
                      <a-select-option value="daily">每天</a-select-option>
                      <a-select-option value="weekly">每周</a-select-option>
                      <a-select-option value="monthly">每月</a-select-option>
                    </a-select>
                    <a-time-picker v-model:value="timeValue" format="HH:mm" value-format="HH:mm" style="width: 132px" />
                    <a-select v-if="simpleSchedule.type === 'weekly'" v-model:value="simpleSchedule.weekday" style="width: 120px">
                      <a-select-option v-for="item in weekdayOptions" :key="item.value" :value="item.value">{{ item.label }}</a-select-option>
                    </a-select>
                    <a-select v-if="simpleSchedule.type === 'monthly'" v-model:value="simpleSchedule.monthDay" style="width: 120px">
                      <a-select-option v-for="day in monthDayOptions" :key="day" :value="day">{{ day }} 日</a-select-option>
                    </a-select>
                  </div>
                  <div class="schedule-builder__hint">{{ scheduleDescription }}</div>
                  <div class="schedule-builder__cron">生成表达式：<code>{{ form.cron_expr }}</code></div>
                </template>

                <template v-else>
                  <a-input v-model:value="form.cron_expr" placeholder="如 0 2 * * 1" allow-clear />
                  <div class="schedule-builder__hint" :class="{ 'schedule-builder__hint--warning': !canDescribeAdvancedCron }">
                    {{ advancedScheduleHint }}
                  </div>
                </template>
              </div>
            </a-form-item>
          </a-col>
          <a-col :span="8">
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
import {
  buildCronExprFromSimple,
  defaultSimpleSchedule,
  describeCronExpr,
  describeSimpleSchedule,
  parseSimpleCronExpr,
  type ScheduleMode,
  type SimpleScheduleState,
} from './cronSchedule'

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
const scheduleMode = ref<ScheduleMode>('simple')
const form = reactive<CronTaskPayload>({
  code: '',
  task_code: '',
  name: '',
  cron_expr: '0 2 * * *',
  params: {},
  remark: '',
  sort: 0,
})
const simpleSchedule = reactive<SimpleScheduleState>(defaultSimpleSchedule())
const weekdayOptions = [
  { label: '周一', value: 1 },
  { label: '周二', value: 2 },
  { label: '周三', value: 3 },
  { label: '周四', value: 4 },
  { label: '周五', value: 5 },
  { label: '周六', value: 6 },
  { label: '周日', value: 0 },
]
const monthDayOptions = Array.from({ length: 31 }, (_, index) => index + 1)

const rules = {
  code: [{ required: true, message: '请输入实例编码', trigger: 'blur' }],
  task_code: [{ required: true, message: '请选择注册任务', trigger: 'change' }],
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  cron_expr: [{ required: true, message: '请输入 Cron 表达式', trigger: 'blur' }],
}

const selectedRegisteredTask = computed(() => props.registry.find((item) => item.code === form.task_code))
const scheduleDescription = computed(() => describeSimpleSchedule(simpleSchedule))
const canDescribeAdvancedCron = computed(() => !!parseSimpleCronExpr(form.cron_expr))
const advancedScheduleHint = computed(() => describeCronExpr(form.cron_expr))
const timeValue = computed({
  get: () => simpleSchedule.time,
  set: (value?: string | null) => {
    simpleSchedule.time = value || '02:00'
  },
})

const defaultParamsFromSchema = (task?: RegisteredCronTask) => {
  const params: Record<string, any> = {}
  Object.entries(task?.param_schema || {}).forEach(([key, field]) => {
    params[key] = field.default ?? (field.type === 'int' ? 0 : field.type === 'bool' ? false : '')
  })
  return params
}

const applySimpleSchedule = (value: SimpleScheduleState) => {
  simpleSchedule.type = value.type
  simpleSchedule.time = value.time
  simpleSchedule.weekday = value.weekday
  simpleSchedule.monthDay = value.monthDay
}

const resetForm = () => {
  const firstTask = props.registry[0]
  form.code = ''
  form.task_code = firstTask?.code || ''
  form.name = firstTask?.name || ''
  form.cron_expr = '0 2 * * *'
  scheduleMode.value = 'simple'
  applySimpleSchedule(parseSimpleCronExpr(form.cron_expr) || defaultSimpleSchedule())
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
  const parsedSchedule = parseSimpleCronExpr(record.cron_expr)
  scheduleMode.value = parsedSchedule ? 'simple' : 'advanced'
  applySimpleSchedule(parsedSchedule || defaultSimpleSchedule())
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

const handleScheduleModeChange = () => {
  if (scheduleMode.value !== 'simple') {
    return
  }
  applySimpleSchedule(parseSimpleCronExpr(form.cron_expr) || simpleSchedule)
  form.cron_expr = buildCronExprFromSimple(simpleSchedule)
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

watch(
  simpleSchedule,
  () => {
    if (scheduleMode.value === 'simple') {
      form.cron_expr = buildCronExprFromSimple(simpleSchedule)
    }
  },
  { deep: true },
)
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
.schedule-builder {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
  background: var(--app-bg-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}
.schedule-builder__row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.schedule-builder__hint {
  color: #1677ff;
  font-size: 13px;
  font-weight: 500;
}
.schedule-builder__hint--warning {
  color: #d46b08;
}
.schedule-builder__cron {
  color: var(--app-text-muted);
  font-size: 12px;
}
.schedule-builder__cron code {
  padding: 2px 6px;
  color: #0958d9;
  background: rgba(22, 119, 255, 0.08);
  border-radius: 4px;
}
</style>
