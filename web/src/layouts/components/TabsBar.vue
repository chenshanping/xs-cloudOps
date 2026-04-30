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
      <template #rightExtra>
        <a-dropdown placement="bottomRight">
          <a-button size="small" type="text">标签操作</a-button>
          <template #overlay>
            <a-menu @click="handleBatchAction">
              <a-menu-item key="close-left">关闭左侧</a-menu-item>
              <a-menu-item key="close-right">关闭右侧</a-menu-item>
              <a-menu-item key="close-other">关闭其他</a-menu-item>
              <a-menu-item key="close-all">关闭全部</a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </template>

      <a-tab-pane
        v-for="tab in tabs"
        :key="tab.key"
        :closable="!tab.affix"
        :tab="tab.title"
      />
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { MenuInfo } from 'ant-design-vue/es/menu/src/interface'
import { useTabsStore } from '@/store/tabs'
import { useUiStore } from '@/store/ui'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const uiStore = useUiStore()
const tabsStore = useTabsStore()
const tabsDark = computed(() => uiStore.isDark || uiStore.theme.headerDark)

const tabs = computed(() => tabsStore.tabs)
const activeTab = computed({
  get: () => tabsStore.activeKey,
  set: (value: string) => tabsStore.setActiveKey(value)
})

const findMenuTitle = (path: string, menus = userStore.menus || []): string | null => {
  for (const menu of menus) {
    if (menu.path === path) {
      return menu.name
    }
    if (menu.children?.length) {
      const title = findMenuTitle(path, menu.children)
      if (title) {
        return title
      }
    }
  }
  return null
}

const getPageTitle = () => {
  if (typeof route.meta.title === 'string' && route.meta.title) {
    return route.meta.title
  }

  if (route.path === '/dashboard') return '首页'
  if (route.path === '/profile') return '个人中心'
  if (route.path === '/ai') return 'AI助手'

  return findMenuTitle(route.path) || '未命名页面'
}

watch(
  () => [route.fullPath, route.path, route.name, route.meta.title, userStore.menus],
  () => {
    tabsStore.openTab({
      key: route.fullPath,
      path: route.path,
      fullPath: route.fullPath,
      title: getPageTitle(),
      name: typeof route.name === 'string' ? route.name : undefined,
      affix: route.path === '/dashboard'
    })
  },
  { immediate: true, deep: true },
)

const onTabChange = (key: string) => {
  router.push(key)
}

const onTabEdit = (targetKey: string | MouseEvent, action: string) => {
  if (action !== 'remove' || typeof targetKey !== 'string') {
    return
  }

  const nextKey = tabsStore.removeTab(targetKey)
  if (nextKey) {
    router.push(nextKey)
  }
}

const handleBatchAction = ({ key }: MenuInfo) => {
  const currentKey = activeTab.value

  if (key === 'close-left') {
    tabsStore.removeLeftTabs(currentKey)
    return
  }

  if (key === 'close-right') {
    tabsStore.removeRightTabs(currentKey)
    return
  }

  if (key === 'close-other') {
    tabsStore.removeOtherTabs(currentKey)
    return
  }

  if (key === 'close-all') {
    tabsStore.removeAllTabs()
    if (route.fullPath !== '/dashboard') {
      router.push('/dashboard')
    }
  }
}
</script>

<style scoped>
.tabs-wrapper {
  padding: 8px 20px 0;
  background: var(--app-surface-color);
  border-bottom: 1px solid var(--app-border-color);
}

.tabs-wrapper-dark {
  background: var(--app-elevated-bg);
  border-bottom-color: var(--app-border-color);
}

.tabs-wrapper :deep(.ant-tabs) {
  margin-bottom: 0;
}

.tabs-wrapper :deep(.ant-tabs-nav) {
  margin-bottom: 0;
}

.tabs-wrapper :deep(.ant-tabs-extra-content) {
  display: flex;
  align-items: center;
}

.tabs-wrapper :deep(.ant-tabs-tab) {
  padding: 6px 12px;
  font-size: 13px;
  border-radius: 10px 10px 0 0;
}

.tabs-wrapper-dark :deep(.ant-tabs-tab) {
  color: var(--app-text-secondary);
  background: rgba(255, 255, 255, 0.05);
  border-color: var(--app-border-color);
}

.tabs-wrapper-dark :deep(.ant-tabs-tab:hover) {
  color: var(--app-text-strong);
}

.tabs-wrapper-dark :deep(.ant-tabs-tab-active) {
  background: rgba(255, 255, 255, 0.1);
  border-color: var(--app-border-strong);
}

.tabs-wrapper-dark :deep(.ant-tabs-tab-active .ant-tabs-tab-btn),
.tabs-wrapper-dark :deep(.ant-tabs-tab .ant-tabs-tab-btn),
.tabs-wrapper-dark :deep(.ant-tabs-tab-remove) {
  color: inherit;
}

.tabs-wrapper-dark :deep(.ant-tabs-tab-remove:hover) {
  color: var(--app-text-strong);
}

.tabs-wrapper-dark :deep(.ant-tabs-nav::before) {
  border-bottom-color: var(--app-border-color);
}
</style>
