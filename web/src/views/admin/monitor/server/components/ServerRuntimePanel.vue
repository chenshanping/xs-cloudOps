<template>
  <a-spin :spinning="loading">
    <a-alert v-if="error" class="panel-alert" type="error" :message="error" show-icon />
    <a-empty v-else-if="!data" description="暂无运行时指标" />
    <div v-else class="runtime-panel">
      <div class="metric-grid metric-grid--four">
        <div class="metric-card metric-card--primary">
          <div class="metric-card__label">Goroutines</div>
          <div class="metric-card__value">{{ data.goroutines }}</div>
          <div class="metric-card__hint">当前协程数</div>
        </div>
        <div class="metric-card">
          <div class="metric-card__label">Heap Alloc</div>
          <div class="metric-card__value metric-card__value--sm">{{ formatBytes(data.heap_alloc) }}</div>
          <div class="metric-card__hint">堆上已分配</div>
        </div>
        <div class="metric-card">
          <div class="metric-card__label">GC 次数</div>
          <div class="metric-card__value">{{ data.num_gc }}</div>
          <div class="metric-card__hint">强制 {{ data.num_forced_gc }}</div>
        </div>
        <div class="metric-card">
          <div class="metric-card__label">进程运行</div>
          <div class="metric-card__value metric-card__value--sm">{{ formatDuration(data.process_uptime_seconds) }}</div>
          <div class="metric-card__hint">{{ data.go_version }}</div>
        </div>
      </div>

      <div class="monitor-section-grid">
        <a-card class="monitor-card" title="内存分配" size="small">
          <div class="runtime-bars">
            <div class="runtime-bar">
              <div><strong>Heap Inuse</strong><span>{{ formatBytes(data.heap_inuse) }}</span></div>
              <a-progress :percent="heapInusePercent" />
            </div>
            <div class="runtime-bar">
              <div><strong>Heap Sys</strong><span>{{ formatBytes(data.heap_sys) }}</span></div>
              <a-progress :percent="100" status="active" />
            </div>
            <div class="runtime-bar">
              <div><strong>Next GC</strong><span>{{ formatBytes(data.next_gc) }}</span></div>
              <a-progress :percent="nextGcPercent" />
            </div>
          </div>
        </a-card>
        <a-card class="monitor-card" title="栈与调度" size="small">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="NumCPU">{{ data.num_cpu }}</a-descriptions-item>
            <a-descriptions-item label="GOMAXPROCS">{{ data.go_max_procs }}</a-descriptions-item>
            <a-descriptions-item label="Stack Inuse">{{ formatBytes(data.stack_inuse) }}</a-descriptions-item>
            <a-descriptions-item label="Stack Sys">{{ formatBytes(data.stack_sys) }}</a-descriptions-item>
            <a-descriptions-item label="Heap Objects">{{ data.heap_objects }}</a-descriptions-item>
            <a-descriptions-item label="采集时间">{{ formatDateTime(data.collected_at) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </div>

      <a-card class="monitor-card" title="GC 细节" size="small">
        <div class="runtime-gc-grid">
          <a-statistic title="Pause Total" :value="pauseMs" suffix="ms" :precision="2" />
          <a-statistic title="GC CPU Fraction" :value="data.gc_cpu_fraction * 100" suffix="%" :precision="3" />
          <a-statistic title="Last GC" :value="formatDateTime(data.last_gc)" />
        </div>
      </a-card>
    </div>
  </a-spin>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { RuntimeInfo } from '@/api/monitor'
import { formatBytes, formatDateTime, formatDuration } from '../format'

const props = defineProps<{
  data: RuntimeInfo | null
  loading: boolean
  error: string
}>()

const heapInusePercent = computed(() => {
  if (!props.data?.heap_sys) return 0
  return Math.min(100, Number(((props.data.heap_inuse / props.data.heap_sys) * 100).toFixed(2)))
})

const nextGcPercent = computed(() => {
  if (!props.data?.next_gc) return 0
  return Math.min(100, Number(((props.data.heap_alloc / props.data.next_gc) * 100).toFixed(2)))
})

const pauseMs = computed(() => Number(((props.data?.pause_total_ns || 0) / 1000000).toFixed(2)))
</script>
