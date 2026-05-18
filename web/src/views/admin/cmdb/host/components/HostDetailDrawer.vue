<template>
  <a-drawer
    :open="open"
    :title="host ? `主机详情 · ${host.name}` : '主机详情'"
    width="920"
    placement="right"
    destroy-on-close
    @close="handleClose"
  >
    <template #extra>
      <a-space>
        <a-button
          v-if="host"
          type="primary"
          :loading="verifying"
          v-permission="'cmdb:host:verify'"
          @click="emit('verify', host)"
        >
          重新校验
        </a-button>
        <a-button @click="handleClose">关闭</a-button>
      </a-space>
    </template>

    <template v-if="host">
      <div class="detail-hero">
        <div>
          <div class="detail-hero__eyebrow">CMDB 主机台账</div>
          <div class="detail-hero__title">{{ host.name }}</div>
          <div class="detail-hero__desc">
            {{ host.hostname || host.ssh_host }} · {{ host.environment || '未设置环境' }} · {{ host.group.name }}
          </div>
        </div>
        <a-space wrap :size="[8, 8]">
          <a-tag :color="verifyMeta.badge">{{ verifyMeta.text }}</a-tag>
          <a-tag color="blue">{{ host.group.name }}</a-tag>
          <a-tag v-if="host.environment">{{ host.environment }}</a-tag>
        </a-space>
      </div>

      <div class="detail-layout">
        <div class="detail-section">
          <div class="detail-section__title">基础信息</div>
          <a-descriptions :column="2" bordered size="small">
            <a-descriptions-item label="主机名称">{{ host.name }}</a-descriptions-item>
            <a-descriptions-item label="负责人">{{ host.owner || '-' }}</a-descriptions-item>
            <a-descriptions-item label="内网 IP">{{ host.private_ip || '-' }}</a-descriptions-item>
            <a-descriptions-item label="公网 IP">{{ host.public_ip || '-' }}</a-descriptions-item>
            <a-descriptions-item label="主机分组">{{ host.group.name }}</a-descriptions-item>
            <a-descriptions-item label="环境">{{ host.environment || '-' }}</a-descriptions-item>
            <a-descriptions-item label="主机标签" :span="2">
              <a-space v-if="host.tags?.length" :size="[6, 6]" wrap>
                <a-tag v-for="item in host.tags" :key="item.id" :color="item.color || '#1677ff'">{{ item.name }}</a-tag>
              </a-space>
              <span v-else>-</span>
            </a-descriptions-item>
            <a-descriptions-item label="备注" :span="2">{{ host.remark || '-' }}</a-descriptions-item>
          </a-descriptions>
        </div>

        <div class="detail-grid">
          <div class="detail-section">
            <div class="detail-section__title">SSH 连接</div>
            <a-descriptions :column="1" bordered size="small">
              <a-descriptions-item label="SSH 地址">{{ host.ssh_host }}</a-descriptions-item>
              <a-descriptions-item label="SSH 端口">{{ host.ssh_port }}</a-descriptions-item>
              <a-descriptions-item label="凭据名称">{{ host.credential_summary.name }}</a-descriptions-item>
              <a-descriptions-item label="认证方式">{{ getCmdbAuthTypeLabel(host.credential_summary.auth_type) }}</a-descriptions-item>
              <a-descriptions-item label="登录用户">{{ host.credential_summary.username }}</a-descriptions-item>
            </a-descriptions>
          </div>

          <div class="detail-section">
            <div class="detail-section__title">校验结果</div>
            <div class="verify-card">
              <div class="verify-card__status">
                <a-badge :status="verifyMeta.color" />
                <strong>{{ verifyMeta.text }}</strong>
              </div>
              <div class="verify-card__message">{{ host.verify_message || '暂无校验消息' }}</div>
              <div class="verify-card__time">最后校验：{{ formatTime(host.last_verified_at) }}</div>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <div class="detail-section__title">系统回填信息</div>
          <div class="metrics-grid">
            <div class="metric-item">
              <span>主机名</span>
              <strong>{{ host.hostname || '-' }}</strong>
              <p>SSH 地址：{{ host.ssh_host }}:{{ host.ssh_port }}</p>
            </div>
            <div class="metric-item">
              <span>系统信息</span>
              <strong>{{ host.platform_version || host.platform || host.os || '-' }}</strong>
              <p>操作系统：{{ host.os || '-' }}</p>
            </div>
            <div class="metric-item">
              <span>内核 / 架构</span>
              <strong>{{ host.kernel_version || '-' }}</strong>
              <p>架构：{{ host.architecture || '-' }}</p>
            </div>
            <div class="metric-item">
              <span>资源规格</span>
              <strong>{{ host.cpu_cores || 0 }}C / {{ formatCmdbMemory(host.memory_mb) }}</strong>
              <p>CPU / 内存</p>
            </div>
          </div>
        </div>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatTime } from '@/utils/format'
import { type CmdbHostItem } from '@/api/cmdb'
import { formatCmdbMemory, getCmdbAuthTypeLabel, getCmdbVerifyStatusMeta } from '../../shared'

const props = defineProps<{
  open: boolean
  host: CmdbHostItem | null
  verifying?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'verify', host: CmdbHostItem): void
}>()

const verifyMeta = computed(() => getCmdbVerifyStatusMeta(props.host?.verify_status))

const handleClose = () => {
  emit('update:open', false)
}
</script>

<style scoped>
.detail-hero {
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

.detail-hero__eyebrow {
  color: #1677ff;
  font-size: 12px;
  font-weight: 600;
}

.detail-hero__title {
  margin-top: 6px;
  color: var(--app-text-strong);
  font-size: 18px;
  font-weight: 700;
}

.detail-hero__desc {
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

.detail-section__title {
  margin-bottom: 14px;
  color: var(--app-text-strong);
  font-size: 14px;
  font-weight: 600;
}

.verify-card {
  padding: 14px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.verify-card__status {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--app-text-strong);
}

.verify-card__message {
  margin-top: 10px;
  color: var(--app-text-color);
  font-size: 13px;
  line-height: 1.7;
  word-break: break-all;
}

.verify-card__time {
  margin-top: 10px;
  color: var(--app-text-muted);
  font-size: 12px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.metric-item {
  padding: 14px;
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.metric-item span {
  color: var(--app-text-muted);
  font-size: 12px;
}

.metric-item strong {
  display: block;
  margin-top: 8px;
  color: var(--app-text-strong);
  font-size: 14px;
  word-break: break-word;
}

.metric-item p {
  margin: 8px 0 0;
  color: var(--app-text-muted);
  font-size: 12px;
  line-height: 1.6;
  word-break: break-word;
}

@media (max-width: 960px) {
  .detail-hero,
  .detail-grid {
    grid-template-columns: 1fr;
    flex-direction: column;
  }

  .metrics-grid {
    grid-template-columns: 1fr;
  }
}
</style>
