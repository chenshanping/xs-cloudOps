<template>
  <a-spin :spinning="loading">
    <a-alert v-if="error" class="panel-alert" type="error" :message="error" show-icon />
    <a-empty v-else-if="!data" description="暂无系统指标" />
    <div v-else class="overview-panel">
      <div class="metric-grid metric-grid--four">
        <div class="metric-card metric-card--primary">
          <div class="metric-card__label">CPU 使用率</div>
          <div class="metric-card__value">{{ data.cpu.usage_percent.toFixed(2) }}%</div>
          <a-progress :percent="formatPercent(data.cpu.usage_percent)" :status="progressStatus(data.cpu.usage_percent)" :show-info="false" />
        </div>
        <div class="metric-card">
          <div class="metric-card__label">内存使用率</div>
          <div class="metric-card__value">{{ data.memory.usage_percent.toFixed(2) }}%</div>
          <a-progress :percent="formatPercent(data.memory.usage_percent)" :status="progressStatus(data.memory.usage_percent)" :show-info="false" />
        </div>
        <div class="metric-card">
          <div class="metric-card__label">逻辑核心</div>
          <div class="metric-card__value">{{ data.cpu.logical_core }}</div>
          <div class="metric-card__hint">物理核心 {{ data.cpu.physical_core || '-' }}</div>
        </div>
        <div class="metric-card">
          <div class="metric-card__label">主机运行</div>
          <div class="metric-card__value metric-card__value--sm">{{ formatDuration(data.host.uptime_seconds) }}</div>
          <div class="metric-card__hint">{{ data.data_source }}</div>
        </div>
      </div>

      <div class="monitor-section-grid">
        <a-card class="monitor-card" title="主机信息" size="small">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="主机名">{{ data.host.hostname || '-' }}</a-descriptions-item>
            <a-descriptions-item label="平台">{{ data.host.platform || data.host.os || '-' }}</a-descriptions-item>
            <a-descriptions-item label="版本">{{ data.host.platform_version || '-' }}</a-descriptions-item>
            <a-descriptions-item label="内核">{{ data.host.kernel_version || '-' }}</a-descriptions-item>
            <a-descriptions-item label="架构">{{ data.host.architecture || '-' }}</a-descriptions-item>
            <a-descriptions-item label="启动时间">{{ formatDateTime(data.host.boot_time) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>

        <a-card class="monitor-card" title="进程信息" size="small">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="PID">{{ data.process.pid }}</a-descriptions-item>
            <a-descriptions-item label="进程名">{{ data.process.binary_name || '-' }}</a-descriptions-item>
            <a-descriptions-item label="Go 版本">{{ data.process.go_version }}</a-descriptions-item>
            <a-descriptions-item label="GOMAXPROCS">{{ data.process.go_max_procs }}</a-descriptions-item>
            <a-descriptions-item label="启动时间">{{ formatDateTime(data.process.started_at) }}</a-descriptions-item>
            <a-descriptions-item label="运行时长">{{ formatDuration(data.process.uptime_seconds) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </div>

      <div class="monitor-section-grid monitor-section-grid--wide">
        <a-card class="monitor-card" title="内存 / Swap" size="small">
          <div class="usage-list">
            <div class="usage-row">
              <div>
                <strong>物理内存</strong>
                <span>{{ formatBytes(data.memory.used) }} / {{ formatBytes(data.memory.total) }}</span>
              </div>
              <a-progress :percent="formatPercent(data.memory.usage_percent)" :status="progressStatus(data.memory.usage_percent)" />
            </div>
            <div class="usage-row">
              <div>
                <strong>Swap</strong>
                <span>{{ formatBytes(data.swap.used) }} / {{ formatBytes(data.swap.total) }}</span>
              </div>
              <a-progress :percent="formatPercent(data.swap.usage_percent)" :status="progressStatus(data.swap.usage_percent)" />
            </div>
          </div>
        </a-card>

        <a-card class="monitor-card" title="系统负载" size="small">
          <div class="load-strip">
            <a-statistic title="1 min" :value="data.load.load_1" :precision="2" />
            <a-statistic title="5 min" :value="data.load.load_5" :precision="2" />
            <a-statistic title="15 min" :value="data.load.load_15" :precision="2" />
          </div>
        </a-card>
      </div>

      <a-card class="monitor-card" title="磁盘分区" size="small">
        <div v-if="data.disks.length" class="disk-grid">
          <div v-for="disk in data.disks" :key="disk.mountpoint" class="disk-card">
            <div class="disk-card__top">
              <strong>{{ disk.mountpoint }}</strong>
              <a-tag>{{ disk.fs_type || '-' }}</a-tag>
            </div>
            <a-progress :percent="formatPercent(disk.usage_percent)" :status="progressStatus(disk.usage_percent)" />
            <div class="disk-card__meta">
              <span>已用 {{ formatBytes(disk.used) }}</span>
              <span>剩余 {{ formatBytes(disk.free) }}</span>
              <span>总计 {{ formatBytes(disk.total) }}</span>
            </div>
          </div>
        </div>
        <a-empty v-else description="暂无磁盘数据" />
      </a-card>
    </div>
  </a-spin>
</template>

<script setup lang="ts">
import type { ServerInfo } from '@/api/monitor'
import { formatBytes, formatDateTime, formatDuration, formatPercent, progressStatus } from '../format'

defineProps<{
  data: ServerInfo | null
  loading: boolean
  error: string
}>()
</script>
