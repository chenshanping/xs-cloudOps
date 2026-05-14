<template>
  <a-spin :spinning="loading">
    <a-alert v-if="error" class="panel-alert" type="error" :message="error" show-icon />
    <a-empty v-else-if="!data" description="暂无 Redis 指标" />
    <div v-else class="cache-panel">
      <a-alert
        v-if="!data.reachable"
        class="panel-alert"
        type="error"
        :message="data.error || 'Redis 不可达'"
        show-icon
      />
      <div class="metric-grid metric-grid--four">
        <div class="metric-card" :class="data.reachable ? 'metric-card--success' : 'metric-card--danger'">
          <div class="metric-card__label">Redis 状态</div>
          <div class="metric-card__value metric-card__value--sm">{{ data.reachable ? '正常' : '异常' }}</div>
          <div class="metric-card__hint">Ping {{ data.ping_latency_ms }} ms</div>
        </div>
        <div class="metric-card metric-card--primary">
          <div class="metric-card__label">Keys</div>
          <div class="metric-card__value">{{ data.db_size }}</div>
          <div class="metric-card__hint">当前 DB</div>
        </div>
        <div class="metric-card">
          <div class="metric-card__label">命令总数</div>
          <div class="metric-card__value metric-card__value--sm">{{ data.total_commands_processed }}</div>
          <div class="metric-card__hint">版本 {{ data.version || '-' }}</div>
        </div>
        <div class="metric-card">
          <div class="metric-card__label">命中率</div>
          <div class="metric-card__value">{{ data.hit_rate.toFixed(2) }}%</div>
          <div class="metric-card__hint">{{ data.keyspace_hits }} / {{ data.keyspace_misses }}</div>
        </div>
      </div>

      <div class="monitor-section-grid">
        <a-card class="monitor-card" title="内存与连接" size="small">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="已用内存">{{ data.used_memory_human || formatBytes(data.used_memory) }}</a-descriptions-item>
            <a-descriptions-item label="峰值内存">{{ formatBytes(data.used_memory_peak) }}</a-descriptions-item>
            <a-descriptions-item label="客户端连接">{{ data.connected_clients }}</a-descriptions-item>
            <a-descriptions-item label="运行时间">{{ formatDuration(data.uptime_seconds) }}</a-descriptions-item>
            <a-descriptions-item label="采集时间" :span="2">{{ formatDateTime(data.collected_at) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>

        <a-card class="monitor-card" title="缓存清理白名单" size="small">
          <div class="prefix-list">
            <div v-for="item in data.prefix_counts" :key="item.prefix" class="prefix-row">
              <div class="prefix-row__main">
                <code>{{ item.prefix }}</code>
                <span>{{ item.count }} keys</span>
                <a-tag v-if="item.truncated" color="orange">已截断</a-tag>
                <a-tooltip v-if="item.error" :title="item.error">
                  <a-tag color="red">计数失败</a-tag>
                </a-tooltip>
              </div>
              <a-button
                danger
                size="small"
                :loading="clearLoading && pendingPrefix === item.prefix"
                @click="openClearConfirm(item.prefix)"
                v-permission="'monitor:cache:clear'"
              >清理</a-button>
            </div>
          </div>
        </a-card>
      </div>
    </div>
  </a-spin>

  <a-modal
    v-model:open="confirmOpen"
    title="确认清理 Redis 缓存"
    ok-text="确认清理"
    cancel-text="取消"
    :confirm-loading="clearLoading"
    :ok-button-props="{ danger: true, disabled: confirmInput !== pendingPrefix }"
    @ok="handleClear"
  >
    <a-alert type="warning" show-icon message="该操作将删除白名单前缀下的全部 Redis key，无法撤销。" />
    <div class="confirm-prefix">
      <div>请输入完整前缀确认：</div>
      <code>{{ pendingPrefix }}</code>
      <a-input v-model:value="confirmInput" placeholder="输入完整 prefix" />
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import type { ClearCacheResult, RedisInfo } from '@/api/monitor'
import { formatBytes, formatDateTime, formatDuration } from '../format'

defineProps<{
  data: RedisInfo | null
  loading: boolean
  clearLoading: boolean
  error: string
}>()

const emit = defineEmits<{
  clear: [prefix: string]
}>()

const confirmOpen = ref(false)
const pendingPrefix = ref('')
const confirmInput = ref('')

function openClearConfirm(prefix: string) {
  pendingPrefix.value = prefix
  confirmInput.value = ''
  confirmOpen.value = true
}

function handleClear() {
  if (confirmInput.value !== pendingPrefix.value) {
    message.warning('请完整输入待清理前缀')
    return
  }
  emit('clear', pendingPrefix.value)
  confirmOpen.value = false
}

defineExpose({
  notifyClearResult(result: ClearCacheResult | null) {
    if (result?.truncated) {
      message.warning('缓存清理达到扫描上限，可能仍有剩余 key')
    }
  },
})
</script>
