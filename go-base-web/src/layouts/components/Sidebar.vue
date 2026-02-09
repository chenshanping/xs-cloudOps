<template>
  <a-layout-sider 
    v-model:collapsed="collapsed" 
    collapsible 
    class="sider"
    :style="{ background: configStore.get('menu_bg_color') }"
  >
    <!-- Logo 区域 -->
    <div class="logo">
      <img :src="configStore.get('sys_logo')" alt="logo" class="logo-img" />
      <span v-if="!collapsed" class="logo-title">{{ configStore.get('sys_name') }}后台管理系统</span>
    </div>
    <!-- 菜单区域 -->
    <div class="menu-wrapper">
      <a-menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        mode="inline"
        :style="menuStyle"
      >
        <a-menu-item key="/dashboard" @click="$router.push('/dashboard')">
          <template #icon><DashboardOutlined /></template>
          <span>首页</span>
        </a-menu-item>
        <template v-for="menu in menuList" :key="menu.id">
          <a-sub-menu v-if="menu.type === 1" :key="`/menu-${menu.id}`">
            <template v-if="menu.icon" #icon>
              <component v-if="!menu.icon.startsWith('custom-')" :is="getIconComponent(menu.icon)" />
              <img v-else :src="getCustomIconUrl(menu.icon)" style="width: 1em; height: 1em" alt="" />
            </template>
            <template #title>{{ menu.name }}</template>
            <a-menu-item
              v-for="child in menu.children"
              :key="child.path"
              @click="$router.push(child.path)"
            >
              <template v-if="child.icon" #icon>
                <component v-if="!child.icon.startsWith('custom-')" :is="getIconComponent(child.icon)" />
                <img v-else :src="getCustomIconUrl(child.icon)" style="width: 1em; height: 1em" alt="" />
              </template>
              {{ child.name }}
            </a-menu-item>
          </a-sub-menu>
          <a-menu-item v-else-if="menu.type === 2" :key="menu.path" @click="$router.push(menu.path)">
            <template v-if="menu.icon" #icon>
              <component v-if="!menu.icon.startsWith('custom-')" :is="getIconComponent(menu.icon)" />
              <img v-else :src="getCustomIconUrl(menu.icon)" style="width: 1em; height: 1em" alt="" />
            </template>
            <span>{{ menu.name }}</span>
          </a-menu-item>
        </template>
      </a-menu>
    </div>
  </a-layout-sider>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useConfigStore } from '@/store/config'
import * as AntIcons from '@ant-design/icons-vue'
import { DashboardOutlined } from '@ant-design/icons-vue'
import type { Menu } from '@/types'

const route = useRoute()
const userStore = useUserStore()
const configStore = useConfigStore()

const collapsed = ref(false)
const selectedKeys = ref<string[]>([route.path])
const openKeys = ref<string[]>([])

const menuList = computed(() => {
  const menus = userStore.menus || []
  return menus.filter(m => m.status === 1)
})

// 菜单样式
const menuStyle = computed(() => ({
  '--menu-bg-color': configStore.get('menu_bg_color'),
  '--menu-text-color': configStore.get('menu_text_color'),
  '--menu-active-bg-color': configStore.get('menu_active_bg_color'),
  '--menu-active-text-color': configStore.get('menu_active_text_color')
}))

const findParentMenu = (menus: Menu[], path: string): Menu | null => {
  for (const menu of menus) {
    if (menu.path === path) return menu
    if (menu.children) {
      for (const child of menu.children) {
        if (child.path === path) return menu
      }
    }
  }
  return null
}

// 预加载所有自定义图标
const iconModules = import.meta.glob('@/assets/icons/*.svg', { eager: true, query: '?url', import: 'default' })

const getIconComponent = (iconName?: string) => {
  if (!iconName) return 'DatabaseOutlined'
  let name = iconName
  if (name.startsWith('official-')) {
    name = name.replace('official-', '')
  } else if (name.startsWith('custom-')) {
    return null
  }
  return (AntIcons as any)[name] || 'DatabaseOutlined'
}

const getCustomIconUrl = (iconName?: string) => {
  if (!iconName) return ''
  const name = iconName.replace('custom-', '')
  const key = `/src/assets/icons/${name}.svg`
  return (iconModules[key] as string) || ''
}

watch(() => route.path, (path) => {
  selectedKeys.value = [path]
  const menus = userStore.menus || []
  const menu = findParentMenu(menus, path)
  if (menu) {
    openKeys.value = [`/menu-${menu.id}`]
  }
}, { immediate: true })
</script>

<style scoped>
.sider {
  display: flex;
  flex-direction: column;
}

.sider :deep(.ant-layout-sider-trigger) {
  background: rgba(0, 0, 0, 0.2);
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  background: rgba(0, 0, 0, 0.2);
}

.logo-img {
  width: 32px;
  height: 32px;
}

.logo-title {
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  margin-left: 12px;
  white-space: nowrap;
}

.menu-wrapper {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

.menu-wrapper::-webkit-scrollbar {
  width: 6px;
}

.menu-wrapper::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

/* 菜单样式覆盖 */
.menu-wrapper :deep(.ant-menu) {
  background: var(--menu-bg-color) !important;
  color: var(--menu-text-color);
}

.menu-wrapper :deep(.ant-menu-item),
.menu-wrapper :deep(.ant-menu-submenu-title) {
  color: var(--menu-text-color) !important;
}

.menu-wrapper :deep(.ant-menu-item-selected) {
  background-color: var(--menu-active-bg-color) !important;
  color: var(--menu-active-text-color) !important;
}

.menu-wrapper :deep(.ant-menu-item:hover),
.menu-wrapper :deep(.ant-menu-submenu-title:hover) {
  color: var(--menu-active-text-color) !important;
}

.menu-wrapper :deep(.ant-menu-sub) {
  background: rgba(0, 0, 0, 0.2) !important;
}
</style>
