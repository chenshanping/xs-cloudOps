<template>
  <div class="dependency-badge">
    <div class="dependency-badge__items">
      <div class="dependency-pill" :class="{ 'dependency-pill--ok': health?.db?.reachable, 'dependency-pill--bad': health && !health.db?.reachable }">
        <DatabaseOutlined />
        <span>DB</span>
        <strong>{{ health?.db?.reachable ? `${health.db.ping_latency_ms}ms` : '异常' }}</strong>
      </div>
      <div class="dependency-pill" :class="{ 'dependency-pill--ok': health?.redis?.reachable, 'dependency-pill--bad': health && !health.redis?.reachable }">
        <CloudServerOutlined />
        <span>Redis</span>
        <strong>{{ health?.redis?.reachable ? `${health.redis.ping_latency_ms}ms` : '异常' }}</strong>
      </div>
      <div class="dependency-pill" :class="{ 'dependency-pill--ok': ossOk, 'dependency-pill--idle': health?.oss && !health.oss.enabled, 'dependency-pill--bad': health?.oss?.enabled && !health.oss.reachable }">
        <CloudOutlined />
        <span>OSS</span>
        <strong>{{ ossText }}</strong>
      </div>
    </div>
    <a-space>
      <a-alert v-if="error" class="dependency-badge__error" type="warning" :message="error" show-icon />
      <a-button size="small" :loading="loading" @click="$emit('refresh')" v-permission="'monitor:dependency:view'">
        <ReloadOutlined /> 重新检测
      </a-button>
    </a-space>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CloudOutlined, CloudServerOutlined, DatabaseOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import type { DependencyHealth } from '@/api/monitor'

const props = defineProps<{
  health: DependencyHealth | null
  loading: boolean
  error: string
}>()

defineEmits<{
  refresh: []
}>()

const ossOk = computed(() => Boolean(props.health?.oss?.enabled && props.health.oss.reachable))
const ossText = computed(() => {
  const oss = props.health?.oss
  if (!oss) return '待检测'
  if (!oss.enabled) return '未启用'
  return oss.reachable ? `${oss.latency_ms}ms` : '异常'
})
</script>

<style scoped>
.dependency-badge {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 14px;
  border: 1px solid var(--monitor-border);
  border-radius: 14px;
  background: var(--monitor-surface);
  box-shadow: var(--monitor-shadow-soft);
}

.dependency-badge__items {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.dependency-badge__error {
  max-width: 320px;
}

.dependency-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-width: 104px;
  padding: 7px 10px;
  border: 1px solid var(--monitor-border);
  border-radius: 999px;
  color: var(--monitor-text-muted);
  background: var(--monitor-surface-soft);
}

.dependency-pill strong {
  margin-left: auto;
  color: var(--monitor-text);
  font-size: 12px;
  font-weight: 600;
}

.dependency-pill--ok {
  border-color: rgba(82, 196, 26, 0.28);
  color: #389e0d;
  background: rgba(82, 196, 26, 0.08);
}

.dependency-pill--bad {
  border-color: rgba(255, 77, 79, 0.28);
  color: #cf1322;
  background: rgba(255, 77, 79, 0.08);
}

.dependency-pill--idle {
  border-color: rgba(250, 173, 20, 0.28);
  color: #d48806;
  background: rgba(250, 173, 20, 0.08);
}

@media (max-width: 920px) {
  .dependency-badge {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
