<template>
  <a-spin :spinning="loading">
    <a-alert v-if="error" class="panel-alert" type="error" :message="error" show-icon />
    <a-empty v-else-if="!data" description="暂无 OSS 健康数据" />
    <div v-else class="oss-panel">
      <div class="metric-grid metric-grid--three">
        <div class="metric-card" :class="statusClass">
          <div class="metric-card__label">OSS 状态</div>
          <div class="metric-card__value metric-card__value--sm">{{ statusText }}</div>
          <div class="metric-card__hint">{{ data.provider || 'local' }}</div>
        </div>
        <div class="metric-card metric-card--primary">
          <div class="metric-card__label">探测耗时</div>
          <div class="metric-card__value">{{ data.latency_ms }}</div>
          <div class="metric-card__hint">ms</div>
        </div>
        <div class="metric-card">
          <div class="metric-card__label">采集时间</div>
          <div class="metric-card__value metric-card__value--xs">{{ formatDateTime(data.collected_at) }}</div>
          <div class="metric-card__hint">实时探测</div>
        </div>
      </div>

      <a-card class="monitor-card" title="存储健康详情" size="small">
        <a-alert
          v-if="!data.enabled"
          type="info"
          show-icon
          message="当前使用本地存储或未启用远端 OSS。"
          description="服务监控不会展示任何 AccessKey、SecretKey 或签名 URL。"
        />
        <a-alert
          v-else-if="!data.reachable"
          type="error"
          show-icon
          :message="data.error || 'OSS 探测失败'"
          description="请检查当前激活存储配置、网络连通性和服务端访问策略。"
        />
        <a-result v-else status="success" title="对象存储可达" sub-title="已完成轻量 Exists 探测，未暴露任何凭证字段。" />
      </a-card>
    </div>
  </a-spin>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { OssHealth } from '@/api/monitor'
import { formatDateTime } from '../format'

const props = defineProps<{
  data: OssHealth | null
  loading: boolean
  error: string
}>()

const statusText = computed(() => {
  if (!props.data) return '-'
  if (!props.data.enabled) return '未启用'
  return props.data.reachable ? '正常' : '异常'
})

const statusClass = computed(() => {
  if (!props.data?.enabled) return 'metric-card--warning'
  return props.data.reachable ? 'metric-card--success' : 'metric-card--danger'
})
</script>
