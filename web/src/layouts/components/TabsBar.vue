<template>
  <div :class="['tabs-wrapper', { 'tabs-wrapper-dark': tabsDark }]">
    <a-tabs
      v-model:activeKey="activeTab"
      hide-add
      size="small"
      type="editable-card"
      @change="onTabChange"
      @edit="onTabEdit"
    >
      <a-tab-pane
        v-for="tab in tabs"
        :key="tab.path"
        :closable="tab.path !== '/dashboard'"
        :tab="tab.title"
      />
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUiStore } from '@/store/ui'
import { useUserStore } from '@/store/user'

interface Tab {
  path: string
  title: string
  name?: string
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const uiStore = useUiStore()
const tabsDark = computed(() => uiStore.isDark || uiStore.theme.headerDark)

const tabs = ref<Tab[]>([{ path: '/dashboard', title: '首页', name: 'Dashboard' }])
const activeTab = ref('/dashboard')

const cachedViews = computed(() => tabs.value.map((tab) => tab.name).filter(Boolean) as string[])
defineExpose({ cachedViews })

const getPageTitle = (path: string): string => {
  if (path === '/dashboard') return '首页'
  if (path === '/profile') return '个人中心'
  if (path === '/ai') return 'AI助手'

  const menus = userStore.menus || []
  for (const menu of menus) {
    if (menu.path === path) return menu.name
    if (menu.children) {
      for (const child of menu.children) {
        if (child.path === path) return child.name
      }
    }
  }
  return '未命名页面'
}

watch(
  () => route.path,
  (path) => {
    activeTab.value = path
    if (!tabs.value.find((tab) => tab.path === path)) {
      tabs.value.push({
        path,
        title: getPageTitle(path),
        name: route.name as string,
      })
    }
  },
  { immediate: true },
)

const onTabChange = (key: string) => {
  router.push(key)
}

const onTabEdit = (targetKey: string | MouseEvent, action: string) => {
  if (action !== 'remove' || typeof targetKey !== 'string') {
    return
  }

  const index = tabs.value.findIndex((tab) => tab.path === targetKey)
  if (index === -1) {
    return
  }

  tabs.value.splice(index, 1)
  if (activeTab.value === targetKey) {
    const nextTab = tabs.value[index] || tabs.value[index - 1]
    if (nextTab) {
      router.push(nextTab.path)
    }
  }
}
</script>

<style scoped>
.tabs-wrapper {
  padding: 8px 20px 0;
  background: #ffffff;
  border-bottom: 1px solid #f0f0f0;
}

.tabs-wrapper-dark {
  background: #0f172a;
  border-bottom-color: rgba(255, 255, 255, 0.08);
}

.tabs-wrapper :deep(.ant-tabs) {
  margin-bottom: 0;
}

.tabs-wrapper :deep(.ant-tabs-nav) {
  margin-bottom: 0;
}

.tabs-wrapper :deep(.ant-tabs-tab) {
  padding: 6px 12px;
  font-size: 13px;
  border-radius: 10px 10px 0 0;
}

.tabs-wrapper-dark :deep(.ant-tabs-tab) {
  color: rgba(255, 255, 255, 0.8);
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.08);
}

.tabs-wrapper-dark :deep(.ant-tabs-tab:hover) {
  color: #ffffff;
}

.tabs-wrapper-dark :deep(.ant-tabs-tab-active) {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.12);
}

.tabs-wrapper-dark :deep(.ant-tabs-tab-active .ant-tabs-tab-btn),
.tabs-wrapper-dark :deep(.ant-tabs-tab .ant-tabs-tab-btn),
.tabs-wrapper-dark :deep(.ant-tabs-tab-remove) {
  color: inherit;
}

.tabs-wrapper-dark :deep(.ant-tabs-tab-remove:hover) {
  color: #ffffff;
}

.tabs-wrapper-dark :deep(.ant-tabs-nav::before) {
  border-bottom-color: rgba(255, 255, 255, 0.08);
}
</style>
