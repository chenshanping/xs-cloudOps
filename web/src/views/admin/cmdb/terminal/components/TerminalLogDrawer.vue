<template>
  <a-drawer
    :open="open"
    :title="session ? `会话日志 · ${session.host_name}` : '会话日志'"
    width="920"
    placement="right"
    destroy-on-close
    @close="handleClose"
  >
    <template #extra>
      <a-space>
        <a-select
          v-model:value="innerStreamType"
          allow-clear
          placeholder="全部日志"
          style="width: 160px"
          @change="handleFilterChange"
        >
          <a-select-option value="input">输入日志</a-select-option>
          <a-select-option value="output">输出日志</a-select-option>
          <a-select-option value="system">系统日志</a-select-option>
        </a-select>
        <a-button :loading="loading" @click="emit('refresh')">刷新</a-button>
        <a-button @click="handleClose">关闭</a-button>
      </a-space>
    </template>

    <div class="log-layout">
      <div v-if="session" class="session-summary">
        <div class="summary-item">
          <span>主机</span>
          <strong>{{ session.host_name }}</strong>
        </div>
        <div class="summary-item">
          <span>发起人</span>
          <strong>{{ session.username_snapshot }}</strong>
        </div>
        <div class="summary-item">
          <span>状态</span>
          <strong>{{ statusLabel }}</strong>
        </div>
        <div class="summary-item">
          <span>断开原因</span>
          <strong>{{ session.disconnect_reason || '-' }}</strong>
        </div>
      </div>

      <div class="log-panel">
        <div class="log-panel__meta">
          <span>共 {{ pagination.total }} 条</span>
          <span>按会话原始顺序展示</span>
        </div>

        <div v-if="logs.length" class="log-list">
          <div
            v-for="item in logs"
            :key="item.id"
            :class="['log-item', `log-item--${item.stream_type}`]"
          >
            <div class="log-item__header">
              <a-tag :color="getStreamTagColor(item.stream_type)">{{ getStreamLabel(item.stream_type) }}</a-tag>
              <span>#{{ item.seq }}</span>
              <span>{{ formatTime(item.created_at) }}</span>
            </div>
            <pre class="log-item__content">{{ item.content || '-' }}</pre>
          </div>
        </div>
        <a-empty v-else description="暂无会话日志" />
      </div>

      <div class="log-pagination">
        <a-pagination
          :current="pagination.current"
          :page-size="pagination.pageSize"
          :total="pagination.total"
          :show-size-changer="true"
          :page-size-options="['20', '50', '100']"
          :show-total="(total: number) => `共 ${total} 条`"
          @change="handlePageChange"
        />
      </div>
    </div>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { CmdbTerminalLogItem, CmdbTerminalSessionItem } from '@/api/cmdb'
import { formatTime } from '@/utils/format'

const props = defineProps<{
  open: boolean
  loading?: boolean
  session?: CmdbTerminalSessionItem | null
  logs: CmdbTerminalLogItem[]
  pagination: {
    current: number
    pageSize: number
    total: number
  }
  streamType?: 'input' | 'output' | 'system'
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'change', payload: { page: number; pageSize: number; streamType?: 'input' | 'output' | 'system' }): void
  (e: 'refresh'): void
}>()

const innerStreamType = ref<'input' | 'output' | 'system' | undefined>(props.streamType)

watch(
  () => props.streamType,
  value => {
    innerStreamType.value = value
  }
)

const statusLabel = computed(() => {
  switch (props.session?.status) {
    case 'prepared':
      return '准备中'
    case 'active':
      return '在线中'
    case 'failed':
      return '连接失败'
    case 'closed':
      return '已关闭'
    default:
      return '-'
  }
})

const getStreamLabel = (value: string) => {
  switch (value) {
    case 'input':
      return '输入'
    case 'output':
      return '输出'
    case 'system':
      return '系统'
    default:
      return value || '-'
  }
}

const getStreamTagColor = (value: string) => {
  switch (value) {
    case 'input':
      return 'blue'
    case 'output':
      return 'green'
    case 'system':
      return 'orange'
    default:
      return 'default'
  }
}

const handleClose = () => {
  emit('update:open', false)
}

const handleFilterChange = () => {
  emit('change', {
    page: 1,
    pageSize: props.pagination.pageSize,
    streamType: innerStreamType.value,
  })
}

const handlePageChange = (page: number, pageSize: number) => {
  emit('change', {
    page,
    pageSize,
    streamType: innerStreamType.value,
  })
}
</script>

<style scoped>
.log-layout {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 100%;
}

.session-summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.summary-item {
  padding: 14px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.summary-item span {
  color: var(--app-text-muted);
  font-size: 12px;
}

.summary-item strong {
  display: block;
  margin-top: 8px;
  color: var(--app-text-strong);
  font-size: 14px;
  word-break: break-word;
}

.log-panel {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
  padding: 16px;
  background: #0f172a;
  border: 1px solid rgba(148, 163, 184, 0.16);
  border-radius: 10px;
}

.log-panel__meta {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  color: rgba(226, 232, 240, 0.72);
  font-size: 12px;
}

.log-list {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 12px;
  min-height: 0;
  overflow: auto;
}

.log-item {
  padding: 12px;
  background: rgba(15, 23, 42, 0.9);
  border: 1px solid rgba(148, 163, 184, 0.16);
  border-left-width: 3px;
  border-radius: 8px;
}

.log-item--input {
  border-left-color: #1677ff;
}

.log-item--output {
  border-left-color: #52c41a;
}

.log-item--system {
  border-left-color: #fa8c16;
}

.log-item__header {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  margin-bottom: 10px;
  color: rgba(226, 232, 240, 0.72);
  font-size: 12px;
}

.log-item__content {
  margin: 0;
  color: #e2e8f0;
  font-size: 12px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
}

.log-pagination {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 960px) {
  .session-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .session-summary {
    grid-template-columns: 1fr;
  }
}
</style>
