<template>
  <a-layout :class="['layout', `layout-${uiStore.layout.mode}`, { 'layout-dark': uiStore.isDark }]">
    <Sidebar />

    <a-layout class="main-layout">
      <Header v-if="uiStore.effectiveShowHeader" />
      <TabsBar v-if="uiStore.effectiveShowTabs" ref="tabsBarRef" />

      <a-layout-content :style="contentStyle" class="content">
        <ErrorBoundary>
          <router-view v-slot="{ Component }">
            <Suspense>
              <template #default>
                <keep-alive :include="cachedViews">
                  <component :is="Component" :key="$route.fullPath" />
                </keep-alive>
              </template>
              <template #fallback>
                <div class="loading-container">
                  <a-spin size="large" tip="加载中..." />
                </div>
              </template>
            </Suspense>
          </router-view>
        </ErrorBoundary>
      </a-layout-content>
    </a-layout>

    <LayoutSettingsDrawer />

    <a-button
      v-if="!uiStore.effectiveShowHeader"
      class="floating-settings"
      shape="circle"
      type="primary"
      @click="uiStore.toggleSettings(true)"
    >
      <template #icon>
        <SettingOutlined />
      </template>
    </a-button>
  </a-layout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { SettingOutlined } from '@ant-design/icons-vue'
import ErrorBoundary from '@/components/ErrorBoundary.vue'
import { useUiStore } from '@/store/ui'
import { Header, LayoutSettingsDrawer, Sidebar, TabsBar } from './components'

const uiStore = useUiStore()

const tabsBarRef = ref<InstanceType<typeof TabsBar> | null>(null)
const cachedViews = computed(() => tabsBarRef.value?.cachedViews || [])

const contentStyle = computed(() => {
  const padding = `${uiStore.contentPadding}px`
  return {
    margin: padding,
    padding,
    borderRadius: '16px',
    minHeight: '280px',
    overflow: 'auto',
    background: uiStore.isDark ? '#111827' : '#ffffff',
    boxShadow: uiStore.isDark ? 'none' : '0 12px 32px rgba(15, 23, 42, 0.06)',
  }
})
</script>

<style scoped>
.layout {
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(22, 119, 255, 0.14), transparent 28%),
    #f5f7fb;
}

.layout-dark {
  background:
    radial-gradient(circle at top left, rgba(56, 189, 248, 0.12), transparent 24%),
    #020617;
}

.main-layout {
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: transparent;
}

.content {
  flex: 1;
}

.loading-container {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.floating-settings {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 44px;
  height: 44px;
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.18);
}

@media (max-width: 960px) {
  .content {
    margin: 12px !important;
    padding: 12px !important;
    border-radius: 12px !important;
  }
}
</style>
