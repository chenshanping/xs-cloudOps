<template>
  <a-drawer
    :open="open"
    :title="task ? `任务详情 · ${task.name}` : '任务详情'"
    width="980"
    placement="right"
    destroy-on-close
    @close="handleClose"
  >
    <template #extra>
      <a-space :size="8">
        <a-button
          v-if="task && canRunNow"
          type="primary"
          @click="handleRunNow"
          v-permission="'monitor:cron:runNow'"
        >
          立即执行
        </a-button>
        <a-button @click="handleClose">关闭</a-button>
      </a-space>
    </template>

    <template v-if="task">
      <div class="detail-overview">
        <div>
          <div class="detail-overview__eyebrow">实例编码</div>
          <div class="detail-overview__title">{{ task.code }}</div>
          <div class="detail-overview__desc">{{ registryInfo?.description || task.remark || '当前任务未填写额外说明' }}</div>
        </div>
        <a-space wrap :size="[8, 8]">
          <a-tag :color="task.status === 1 ? 'green' : 'default'">{{ task.status === 1 ? '启用中' : '已停用' }}</a-tag>
          <a-tag color="blue">{{ task.task_code }}</a-tag>
          <a-tag v-if="task.last_status" :color="runStatusColor(task.last_status)">{{ runStatusText(task.last_status) }}</a-tag>
        </a-space>
      </div>

      <a-tabs v-model:activeKey="activeTab">
        <a-tab-pane v-if="canViewDetail" key="detail" tab="任务信息">
          <div class="detail-layout">
            <div class="detail-section">
              <div class="section-title">调度信息</div>
              <a-descriptions :column="2" bordered size="small">
                <a-descriptions-item label="任务名称">{{ task.name }}</a-descriptions-item>
                <a-descriptions-item label="注册任务">{{ task.task_code }}</a-descriptions-item>
                <a-descriptions-item label="Cron 表达式"><code>{{ task.cron_expr }}</code></a-descriptions-item>
                <a-descriptions-item label="排序">{{ task.sort }}</a-descriptions-item>
                <a-descriptions-item label="当前状态">
                  <a-tag :color="task.status === 1 ? 'green' : 'default'">{{ task.status === 1 ? '启用' : '停用' }}</a-tag>
                </a-descriptions-item>
                <a-descriptions-item label="上次状态">
                  <a-tag :color="runStatusColor(task.last_status)">{{ runStatusText(task.last_status) }}</a-tag>
                </a-descriptions-item>
                <a-descriptions-item label="上次执行">{{ formatTime(task.last_run_at) }}</a-descriptions-item>
                <a-descriptions-item label="下次执行">{{ task.status === 1 ? formatTime(task.next_run_at) : '-' }}</a-descriptions-item>
                <a-descriptions-item label="上次耗时">{{ formatDuration(task.last_duration_ms) }}</a-descriptions-item>
                <a-descriptions-item label="更新时间">{{ formatTime(task.updated_at) }}</a-descriptions-item>
                <a-descriptions-item label="备注" :span="2">{{ task.remark || '-' }}</a-descriptions-item>
              </a-descriptions>
            </div>

            <div class="detail-grid">
              <div class="detail-section">
                <div class="section-title">任务说明</div>
                <div class="registry-card">
                  <div class="registry-card__label">注册任务</div>
                  <div class="registry-card__title">{{ registryInfo?.name || task.task_code }}</div>
                  <div class="registry-card__content">{{ registryInfo?.description || '当前注册任务没有提供额外说明' }}</div>
                </div>
              </div>

              <div class="detail-section">
                <div class="section-title">执行参数</div>
                <pre class="params-view">{{ paramsText }}</pre>
              </div>
            </div>

            <div class="detail-section">
              <div class="section-title">参数说明</div>
              <div v-if="parameterEntries.length" class="param-list">
                <div v-for="[name, definition] in parameterEntries" :key="name" class="param-item">
                  <div class="param-item__header">
                    <strong>{{ name }}</strong>
                    <a-space :size="6">
                      <a-tag color="blue">{{ definition.type }}</a-tag>
                      <a-tag v-if="definition.required" color="red">必填</a-tag>
                      <a-tag v-else>可选</a-tag>
                    </a-space>
                  </div>
                  <div class="param-item__content">{{ definition.description || '暂无说明' }}</div>
                  <div class="param-item__meta">
                    <span v-if="definition.default !== undefined">默认值：{{ formatDefinitionValue(definition.default) }}</span>
                    <span v-if="definition.min !== undefined || definition.max !== undefined">
                      范围：{{ definition.min ?? '-' }} ~ {{ definition.max ?? '-' }}
                    </span>
                  </div>
                </div>
              </div>
              <a-empty v-else description="该任务没有定义额外参数说明" />
            </div>
          </div>
        </a-tab-pane>

        <a-tab-pane v-if="canViewLogs" key="logs" tab="执行日志">
          <CronTaskLogPanel
            v-if="task"
            :task="task"
            :active="open && activeTab === 'logs'"
            :auto-open-log-id="initialLogId"
          />
        </a-tab-pane>
      </a-tabs>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { CronParamDefinition, CronTask, RegisteredCronTask } from '@/api/cron'
import { formatTime } from '@/utils/format'
import { useUserStore } from '@/store/user'
import CronTaskLogPanel from './CronTaskLogPanel.vue'

const props = defineProps<{
  open: boolean
  task: CronTask | null
  registry: RegisteredCronTask[]
  initialTab?: 'detail' | 'logs'
  initialLogId?: number
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'run-now', task: CronTask): void
}>()

const userStore = useUserStore()
const activeTab = ref<'detail' | 'logs'>('detail')

const canViewDetail = computed(() => userStore.hasPermission('monitor:cron:view'))
const canViewLogs = computed(() => userStore.hasPermission('monitor:cron:logs:view'))
const canRunNow = computed(() => userStore.hasPermission('monitor:cron:runNow'))
const registryInfo = computed(() => props.registry.find((item) => item.code === props.task?.task_code))
const parameterEntries = computed(() => Object.entries(registryInfo.value?.param_schema || {}) as Array<[string, CronParamDefinition]>)
const paramsText = computed(() => {
  if (!props.task?.params || Object.keys(props.task.params).length === 0) {
    return '{}'
  }
  return JSON.stringify(props.task.params, null, 2)
})

const runStatusColor = (status?: string) => ({ success: 'green', failure: 'red', running: 'processing', skipped: 'orange' }[status || ''] || 'default')
const runStatusText = (status?: string) => ({ success: '成功', failure: '失败', running: '执行中', skipped: '已跳过' }[status || ''] || '-')
const formatDuration = (duration?: number) => duration || duration === 0 ? `${duration}ms` : '-'
const formatDefinitionValue = (value: unknown) => typeof value === 'string' ? value : JSON.stringify(value)

const resolveDefaultTab = (requested?: 'detail' | 'logs'): 'detail' | 'logs' => {
  if (requested === 'logs' && canViewLogs.value) {
    return 'logs'
  }
  if (canViewDetail.value) {
    return 'detail'
  }
  if (canViewLogs.value) {
    return 'logs'
  }
  return 'detail'
}

const handleClose = () => {
  emit('update:open', false)
}

const handleRunNow = () => {
  if (!props.task) {
    return
  }
  emit('run-now', props.task)
}

watch(
  () => [props.open, props.task?.id, props.initialTab] as const,
  ([open]) => {
    if (!open) {
      return
    }
    activeTab.value = resolveDefaultTab(props.initialTab)
  },
  { immediate: true },
)
</script>

<style scoped>
.detail-overview {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(22, 119, 255, 0.08), rgba(22, 119, 255, 0.02));
  border: 1px solid rgba(22, 119, 255, 0.16);
  border-radius: 10px;
}

.detail-overview__eyebrow {
  color: #1677ff;
  font-size: 12px;
  font-weight: 600;
}

.detail-overview__title {
  margin-top: 6px;
  color: var(--app-text-strong);
  font-size: 18px;
  font-weight: 700;
}

.detail-overview__desc {
  margin-top: 8px;
  color: var(--app-text-muted);
  font-size: 13px;
  line-height: 1.6;
}

.detail-layout {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.detail-section {
  padding: 16px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.section-title {
  margin-bottom: 14px;
  color: var(--app-text-strong);
  font-size: 14px;
  font-weight: 600;
}

.registry-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 100%;
  padding: 14px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.registry-card__label {
  color: var(--app-text-muted);
  font-size: 12px;
}

.registry-card__title {
  color: var(--app-text-strong);
  font-size: 15px;
  font-weight: 600;
}

.registry-card__content {
  color: var(--app-text-color);
  font-size: 13px;
  line-height: 1.7;
}

.params-view {
  min-height: 126px;
  margin: 0;
  padding: 14px;
  overflow: auto;
  color: var(--app-text-color);
  background: #0f172a;
  border-radius: 8px;
  font-size: 12px;
  font-family: Consolas, monospace;
  line-height: 1.6;
}

.param-list {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
}

.param-item {
  padding: 14px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.param-item__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.param-item__content {
  margin-top: 8px;
  color: var(--app-text-color);
  font-size: 13px;
  line-height: 1.7;
}

.param-item__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 8px;
  color: var(--app-text-muted);
  font-size: 12px;
}

code {
  padding: 2px 6px;
  color: #0958d9;
  background: rgba(22, 119, 255, 0.08);
  border-radius: 4px;
}

@media (max-width: 900px) {
  .detail-overview {
    flex-direction: column;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
