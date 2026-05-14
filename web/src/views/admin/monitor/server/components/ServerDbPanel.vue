<template>
  <a-spin :spinning="loading">
    <a-alert v-if="error" class="panel-alert" type="error" :message="error" show-icon />
    <a-empty v-else-if="!data" description="暂无数据库指标" />
    <div v-else class="db-panel">
      <a-alert
        v-if="!data.reachable"
        class="panel-alert"
        type="error"
        :message="data.error || '数据库不可达'"
        show-icon
      />
      <div class="metric-grid metric-grid--four">
        <div class="metric-card" :class="data.reachable ? 'metric-card--success' : 'metric-card--danger'">
          <div class="metric-card__label">连接状态</div>
          <div class="metric-card__value metric-card__value--sm">{{ data.reachable ? '正常' : '异常' }}</div>
          <div class="metric-card__hint">Ping {{ data.ping_latency_ms }} ms</div>
        </div>
        <div class="metric-card metric-card--primary">
          <div class="metric-card__label">Open</div>
          <div class="metric-card__value">{{ data.open_connections }}</div>
          <div class="metric-card__hint">最大 {{ maxOpenText }}</div>
        </div>
        <div class="metric-card">
          <div class="metric-card__label">In Use</div>
          <div class="metric-card__value">{{ data.in_use }}</div>
          <div class="metric-card__hint">Idle {{ data.idle }}</div>
        </div>
        <div class="metric-card">
          <div class="metric-card__label">等待次数</div>
          <div class="metric-card__value">{{ data.wait_count }}</div>
          <div class="metric-card__hint">累计 {{ data.wait_duration_ms }} ms</div>
        </div>
      </div>

      <div class="monitor-section-grid">
        <a-card class="monitor-card" title="连接池占用" size="small">
          <div class="usage-row">
            <div>
              <strong>Open / Max Open</strong>
              <span>{{ data.open_connections }} / {{ maxOpenText }}</span>
            </div>
            <a-progress :percent="openPercent" :show-info="true" />
          </div>
          <div class="usage-row">
            <div>
              <strong>In Use / Open</strong>
              <span>{{ data.in_use }} / {{ data.open_connections || 0 }}</span>
            </div>
            <a-progress :percent="inUsePercent" :show-info="true" />
          </div>
        </a-card>
        <a-card class="monitor-card" title="关闭统计" size="small">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="Max Idle Closed">{{ data.max_idle_closed }}</a-descriptions-item>
            <a-descriptions-item label="Max Idle Time Closed">{{ data.max_idle_time_closed }}</a-descriptions-item>
            <a-descriptions-item label="Max Lifetime Closed">{{ data.max_lifetime_closed }}</a-descriptions-item>
            <a-descriptions-item label="采集时间">{{ formatDateTime(data.collected_at) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </div>
    </div>
  </a-spin>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { DBStats } from '@/api/monitor'
import { formatDateTime } from '../format'

const props = defineProps<{
  data: DBStats | null
  loading: boolean
  error: string
}>()

const maxOpenText = computed(() => props.data?.max_open_connections ? String(props.data.max_open_connections) : '不限')
const openPercent = computed(() => {
  if (!props.data?.max_open_connections) return 0
  return Math.min(100, Number(((props.data.open_connections / props.data.max_open_connections) * 100).toFixed(2)))
})
const inUsePercent = computed(() => {
  if (!props.data?.open_connections) return 0
  return Math.min(100, Number(((props.data.in_use / props.data.open_connections) * 100).toFixed(2)))
})
</script>
