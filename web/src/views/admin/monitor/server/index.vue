<template>
  <PageWrapper :class="['server-monitor-page', { 'server-monitor-page--dark': uiStore.isDark }]">
    <div class="server-monitor-page__content">
      <ServerDependencyBadge
        class="monitor-dependency"
        :health="dependency"
        :loading="dependencyLoading"
        :error="dependencyError"
        @refresh="loadDependency"
      />

      <a-card class="monitor-tabs-card" :bordered="false">
        <template #extra>
          <a-button :loading="activeLoading" @click="refreshActiveTab">
            <ReloadOutlined /> 刷新当前
          </a-button>
        </template>
        <a-tabs v-model:activeKey="activeTab" @change="handleTabChange">
          <a-tab-pane key="server" tab="系统">
            <ServerOverviewPanel :data="serverInfo" :loading="serverLoading" :error="serverError" />
          </a-tab-pane>
          <a-tab-pane key="runtime" tab="Go Runtime">
            <ServerRuntimePanel :data="runtimeInfo" :loading="runtimeLoading" :error="runtimeError" />
          </a-tab-pane>
          <a-tab-pane key="db" tab="数据库">
            <ServerDbPanel :data="dbStats" :loading="dbLoading" :error="dbError" />
          </a-tab-pane>
          <a-tab-pane key="redis" tab="Redis 缓存">
            <ServerCachePanel
              ref="cachePanelRef"
              :data="redisInfo"
              :loading="redisLoading"
              :clear-loading="clearLoading"
              :error="redisError"
              @clear="handleClearCache"
            />
          </a-tab-pane>
          <a-tab-pane key="oss" tab="OSS">
            <ServerOssPanel :data="ossHealth" :loading="ossLoading" :error="ossError" />
          </a-tab-pane>
        </a-tabs>
      </a-card>
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import { useUiStore } from '@/store/ui'
import ServerCachePanel from './components/ServerCachePanel.vue'
import ServerDbPanel from './components/ServerDbPanel.vue'
import ServerDependencyBadge from './components/ServerDependencyBadge.vue'
import ServerOssPanel from './components/ServerOssPanel.vue'
import ServerOverviewPanel from './components/ServerOverviewPanel.vue'
import ServerRuntimePanel from './components/ServerRuntimePanel.vue'
import { useServerMonitor, type MonitorTabKey } from './useServerMonitor'

const uiStore = useUiStore()
const activeTab = ref<MonitorTabKey>('server')
const cachePanelRef = ref<InstanceType<typeof ServerCachePanel> | null>(null)

const {
  dependency,
  serverInfo,
  runtimeInfo,
  dbStats,
  redisInfo,
  ossHealth,
  dependencyLoading,
  serverLoading,
  runtimeLoading,
  dbLoading,
  redisLoading,
  ossLoading,
  clearLoading,
  dependencyError,
  serverError,
  runtimeError,
  dbError,
  redisError,
  ossError,
  loadDependency,
  loadTab,
  loadTabIfEmpty,
  loadServerInfo,
  clearRedisCache,
} = useServerMonitor()

const activeLoading = computed(() => {
  if (activeTab.value === 'server') return serverLoading.value
  if (activeTab.value === 'runtime') return runtimeLoading.value
  if (activeTab.value === 'db') return dbLoading.value
  if (activeTab.value === 'redis') return redisLoading.value || clearLoading.value
  return ossLoading.value
})

function handleTabChange(key: string | number) {
  const tab = String(key) as MonitorTabKey
  activeTab.value = tab
  loadTabIfEmpty(tab)
}

function refreshActiveTab() {
  loadTab(activeTab.value)
}

async function handleClearCache(prefix: string) {
  const result = await clearRedisCache(prefix)
  cachePanelRef.value?.notifyClearResult(result)
}

onMounted(() => {
  loadDependency()
  loadServerInfo()
})
</script>

<style scoped>
.server-monitor-page {
  --monitor-surface: #ffffff;
  --monitor-surface-soft: #f7f9fc;
  --monitor-border: #edf0f5;
  --monitor-text: #1f2937;
  --monitor-text-muted: #64748b;
  --monitor-shadow-soft: 0 10px 28px rgba(15, 23, 42, 0.06);
  --monitor-primary: var(--app-primary-color, #0960be);
  min-height: 100%;
  padding: 2px;
  color: var(--monitor-text);
}

.server-monitor-page__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.server-monitor-page--dark {
  --monitor-surface: #161d29;
  --monitor-surface-soft: #111827;
  --monitor-border: rgba(148, 163, 184, 0.18);
  --monitor-text: rgba(255, 255, 255, 0.9);
  --monitor-text-muted: rgba(255, 255, 255, 0.62);
  --monitor-shadow-soft: 0 18px 38px rgba(0, 0, 0, 0.22);
}

.monitor-dependency {
  margin-bottom: 12px;
}

.monitor-tabs-card {
  border-radius: 18px;
  background: var(--monitor-surface);
  box-shadow: var(--monitor-shadow-soft);
}

.monitor-tabs-card :deep(.ant-card-head) {
  border-bottom-color: var(--monitor-border);
}

.monitor-tabs-card :deep(.ant-tabs-nav) {
  margin-bottom: 18px;
}

.panel-alert {
  margin-bottom: 14px;
}

:deep(.metric-grid) {
  display: grid;
  gap: 12px;
  margin-bottom: 12px;
}

:deep(.metric-grid--three) {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

:deep(.metric-grid--four) {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

:deep(.metric-card) {
  min-height: 112px;
  padding: 16px;
  border: 1px solid var(--monitor-border);
  border-radius: 16px;
  background: linear-gradient(180deg, var(--monitor-surface), var(--monitor-surface-soft));
}

:deep(.metric-card--primary) {
  border-color: rgba(9, 96, 190, 0.24);
  background: linear-gradient(180deg, rgba(9, 96, 190, 0.1), var(--monitor-surface));
}

:deep(.metric-card--success) {
  border-color: rgba(82, 196, 26, 0.28);
  background: linear-gradient(180deg, rgba(82, 196, 26, 0.1), var(--monitor-surface));
}

:deep(.metric-card--warning) {
  border-color: rgba(250, 173, 20, 0.3);
  background: linear-gradient(180deg, rgba(250, 173, 20, 0.1), var(--monitor-surface));
}

:deep(.metric-card--danger) {
  border-color: rgba(255, 77, 79, 0.3);
  background: linear-gradient(180deg, rgba(255, 77, 79, 0.1), var(--monitor-surface));
}

:deep(.metric-card__label) {
  color: var(--monitor-text-muted);
  font-size: 12px;
}

:deep(.metric-card__value) {
  margin-top: 8px;
  color: var(--monitor-text);
  font-size: 30px;
  font-weight: 700;
  line-height: 1;
}

:deep(.metric-card__value--sm) {
  font-size: 22px;
}

:deep(.metric-card__value--xs) {
  font-size: 16px;
  line-height: 1.35;
}

:deep(.metric-card__hint) {
  margin-top: 8px;
  color: var(--monitor-text-muted);
  font-size: 12px;
}

:deep(.monitor-section-grid) {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}

:deep(.monitor-section-grid--wide) {
  grid-template-columns: 1.5fr 1fr;
}

:deep(.monitor-card) {
  border: 1px solid var(--monitor-border);
  border-radius: 16px;
  background: var(--monitor-surface);
}

:deep(.monitor-card .ant-card-head) {
  border-bottom-color: var(--monitor-border);
}

:deep(.usage-list),
:deep(.runtime-bars) {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

:deep(.usage-row),
:deep(.runtime-bar) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

:deep(.usage-row > div),
:deep(.runtime-bar > div) {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  color: var(--monitor-text-muted);
}

:deep(.usage-row strong),
:deep(.runtime-bar strong) {
  color: var(--monitor-text);
}

:deep(.load-strip),
:deep(.runtime-gc-grid) {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

:deep(.disk-grid) {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 12px;
}

:deep(.disk-card) {
  padding: 12px;
  border: 1px solid var(--monitor-border);
  border-radius: 14px;
  background: var(--monitor-surface-soft);
}

:deep(.disk-card__top),
:deep(.prefix-row),
:deep(.prefix-row__main) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.disk-card__top),
:deep(.prefix-row) {
  justify-content: space-between;
}

:deep(.disk-card__meta) {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 14px;
  color: var(--monitor-text-muted);
  font-size: 12px;
}

:deep(.prefix-list) {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

:deep(.prefix-row) {
  padding: 10px;
  border: 1px solid var(--monitor-border);
  border-radius: 12px;
  background: var(--monitor-surface-soft);
}

:deep(.prefix-row code),
.confirm-prefix code {
  padding: 2px 6px;
  border-radius: 6px;
  color: var(--monitor-primary);
  background: rgba(9, 96, 190, 0.08);
}

.confirm-prefix {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 14px;
}

@media (max-width: 1180px) {
  :deep(.metric-grid--four) {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 820px) {
  :deep(.metric-grid--three),
  :deep(.metric-grid--four),
  :deep(.monitor-section-grid),
  :deep(.monitor-section-grid--wide),
  :deep(.load-strip),
  :deep(.runtime-gc-grid) {
    grid-template-columns: 1fr;
  }
}
</style>
