<template>
  <div class="tabs-wrapper">
    <a-tabs
      v-model:activeKey="activeTab"
      type="editable-card"
      hide-add
      size="small"
      @edit="onTabEdit"
      @change="onTabChange"
    >
      <a-tab-pane
        v-for="tab in tabs"
        :key="tab.path"
        :tab="tab.title"
        :closable="tab.path !== '/dashboard'"
      />
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'

interface Tab {
  path: string
  title: string
  name?: string
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const tabs = ref<Tab[]>([{ path: '/dashboard', title: '首页', name: 'Dashboard' }])
const activeTab = ref('/dashboard')

// 暴露缓存视图给父组件
const cachedViews = computed(() => tabs.value.map(t => t.name).filter(Boolean) as string[])
defineExpose({ cachedViews })

// 获取页面标题
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

// 监听路由变化
watch(() => route.path, (path) => {
  activeTab.value = path
  if (!tabs.value.find(t => t.path === path)) {
    tabs.value.push({
      path,
      title: getPageTitle(path),
      name: route.name as string
    })
  }
}, { immediate: true })

// Tab 切换
const onTabChange = (key: string) => {
  router.push(key)
}

// Tab 关闭
const onTabEdit = (targetKey: string | MouseEvent, action: string) => {
  if (action === 'remove' && typeof targetKey === 'string') {
    const index = tabs.value.findIndex(t => t.path === targetKey)
    if (index > -1) {
      tabs.value.splice(index, 1)
      if (activeTab.value === targetKey) {
        const newTab = tabs.value[index] || tabs.value[index - 1]
        if (newTab) {
          router.push(newTab.path)
        }
      }
    }
  }
}
</script>

<style scoped>
.tabs-wrapper {
  background: #fff;
  padding: 6px 16px 0;
  border-bottom: 1px solid #f0f0f0;
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
}
</style>
