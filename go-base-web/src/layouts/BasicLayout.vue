<template>
  <a-layout class="layout">
    <!-- 左侧菜单 -->
    <Sidebar />

    <!-- 右侧内容区 -->
    <a-layout class="right-layout">
      <!-- 顶部栏 -->
      <Header />

      <!-- Tab标签栏 -->
      <TabsBar ref="tabsBarRef" />

      <!-- 内容区 -->
      <a-layout-content class="content">
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
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Sidebar, Header, TabsBar } from './components'
import ErrorBoundary from '@/components/ErrorBoundary.vue'

const tabsBarRef = ref<InstanceType<typeof TabsBar> | null>(null)
const cachedViews = computed(() => tabsBarRef.value?.cachedViews || [])
</script>

<style scoped>
.layout {
  min-height: 100vh;
}

.right-layout {
  display: flex;
  flex-direction: column;
}

.content {
  margin: 16px;
  padding: 16px;
  background: #fff;
  border-radius: 4px;
  min-height: 280px;
  overflow: auto;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}
</style>
