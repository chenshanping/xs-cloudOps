<template>
  <a-layout-sider
    v-if="uiStore.effectiveShowSidebar"
    :collapsed="uiStore.layout.sidebarCollapsed"
    :collapsed-width="80"
    :theme="sidebarTheme"
    :trigger="null"
    :width="uiStore.layout.sidebarWidth"
    :class="['sider', { 'sider-dark': sidebarTheme === 'dark' }]"
  >
    <div class="logo" @click="router.push(backendHomePath)">
      <img :src="configStore.get('sys_logo')" alt="logo" class="logo-img" />
      <span v-if="!uiStore.layout.sidebarCollapsed" class="logo-title">
        {{ sidebarTitle }}
      </span>
    </div>

    <div class="menu-wrapper">
      <a-menu
        :inline-collapsed="uiStore.layout.sidebarCollapsed"
        :items="menuItems"
        :mode="'inline'"
        :open-keys="openKeys"
        :selected-keys="selectedKeys"
        :theme="sidebarTheme"
        class="sidebar-menu"
        @click="handleMenuClick"
        @openChange="handleOpenChange"
      />
    </div>
  </a-layout-sider>
</template>

<script setup lang="ts">
import { computed, h, ref, watch } from 'vue'
import type { ItemType, MenuInfo } from 'ant-design-vue/es/menu/src/interface'
import { useRoute, useRouter } from 'vue-router'
import { useConfigStore } from '@/store/config'
import { useUiStore } from '@/store/ui'
import { useUserStore } from '@/store/user'
import type { Menu } from '@/types'
import {
  filterVisibleMenus,
  findMenuTrail,
  firstNavigableMenuPath,
  firstNavigablePath,
  getMixedSidebarMenus,
  normalizePath,
} from './layout-menu'
import MenuIcon from './MenuIcon.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const configStore = useConfigStore()
const uiStore = useUiStore()

const openKeys = ref<string[]>([])

const sidebarTheme = computed(() => (uiStore.isDark || uiStore.theme.sidebarDark ? 'dark' : 'light'))

const sidebarTitle = computed(() => {
  const base = configStore.get('sys_name') || '后台管理'
  return uiStore.layout.mode === 'mixed' ? `${base}` : `${base}后台`
})

const normalizedMenus = computed(() => filterVisibleMenus(userStore.menus || []))
const backendHomePath = computed(() => firstNavigableMenuPath(normalizedMenus.value) || '/no-permission')

const currentMenus = computed(() => {
  if (uiStore.layout.mode === 'mixed') {
    const mixedMenus = getMixedSidebarMenus(normalizedMenus.value, route.path)
    return mixedMenus.length ? mixedMenus : normalizedMenus.value
  }
  return normalizedMenus.value
})

const selectedKeys = computed(() => [normalizePath(route.path)])

const buildMenuItems = (menus: Menu[]): ItemType[] => {
  return menus.map((menu) => {
    const icon = menu.icon ? () => h(MenuIcon, { icon: menu.icon }) : undefined

    if (menu.type === 1 && menu.children?.length) {
      return {
        key: `menu-${menu.id}`,
        icon,
        label: menu.name,
        children: buildMenuItems(menu.children),
      }
    }

    return {
      key: normalizePath(menu.path),
      icon,
      label: menu.name,
    }
  })
}

const menuItems = computed<ItemType[]>(() => buildMenuItems(currentMenus.value))

const getSubmenuKey = (menu: Menu) => `menu-${menu.id}`

const findTopSubmenuKey = (menus: Menu[], targetKey: string, topKey?: string): string | null => {
  for (const menu of menus) {
    const currentTopKey = topKey || getSubmenuKey(menu)
    if (getSubmenuKey(menu) === targetKey) {
      return currentTopKey
    }
    if (menu.children?.length) {
      const matchedKey = findTopSubmenuKey(menu.children, targetKey, currentTopKey)
      if (matchedKey) {
        return matchedKey
      }
    }
  }
  return null
}

watch(
  () => [route.path, currentMenus.value, uiStore.layout.mode],
  () => {
    const trail = findMenuTrail(currentMenus.value, route.path)
    openKeys.value = trail
      .filter((menu) => menu.type === 1)
      .map((menu) => `menu-${menu.id}`)
  },
  { immediate: true, deep: true },
)

const handleMenuClick = ({ key }: MenuInfo) => {
  if (typeof key === 'string' && key.startsWith('/')) {
    router.push(key)
  }
}

const handleOpenChange = (keys: string[]) => {
  const latestOpenKey = keys.find(key => !openKeys.value.includes(key))
  if (!latestOpenKey) {
    openKeys.value = keys
    return
  }

  const latestTopKey = findTopSubmenuKey(currentMenus.value, latestOpenKey)
  if (!latestTopKey) {
    openKeys.value = keys
    return
  }

  openKeys.value = keys.filter(key => findTopSubmenuKey(currentMenus.value, key) === latestTopKey)
}
</script>

<style scoped>
.sider {
  height: 100vh;
  position: sticky;
  top: 0;
  inset-inline-start: 0;
  z-index: 30;
  display: flex;
  flex-direction: column;
  background: var(--app-surface-color);
  border-inline-end: 1px solid rgba(148, 163, 184, 0.12);
}

.sider-dark {
  background: var(--app-elevated-bg);
  border-inline-end-color: var(--app-border-color);
}

.sider :deep(.ant-layout-sider-children) {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: inherit;
}

.logo {
  height: var(--app-header-height, 64px);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 var(--app-page-gap, 20px);
  cursor: pointer;
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);
}

.sider-dark .logo {
  border-bottom-color: var(--app-border-color);
}

.logo-img {
  width: 34px;
  height: 34px;
  border-radius: var(--app-radius-sm, 10px);
  object-fit: cover;
  flex-shrink: 0;
}

.logo-title {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 15px;
  font-weight: 600;
}

.sider-dark .logo-title {
  color: rgba(255, 255, 255, 0.92);
}

.menu-wrapper {
  flex: 1;
  overflow-y: auto;
  padding: 8px 8px 12px;
}

.sidebar-menu {
  height: 100%;
  background: transparent;
  border-inline-end: none;
}

.sidebar-menu :deep(.ant-menu) {
  background: transparent !important;
  border-inline-end: none !important;
}

.sidebar-menu :deep(.ant-menu-item),
.sidebar-menu :deep(.ant-menu-submenu-title) {
  border-radius: var(--app-radius-sm, 10px);
}

.sidebar-menu :deep(.ant-menu-item),
.sidebar-menu :deep(.ant-menu-submenu-title) {
  margin-inline: 0;
  width: 100%;
}

.sidebar-menu :deep(.ant-menu-sub.ant-menu-inline) {
  background: transparent !important;
}

.sider-dark .sidebar-menu :deep(.ant-menu-item),
.sider-dark .sidebar-menu :deep(.ant-menu-submenu-title),
.sider-dark .sidebar-menu :deep(.ant-menu-title-content),
.sider-dark .sidebar-menu :deep(.ant-menu-item .ant-menu-item-icon),
.sider-dark .sidebar-menu :deep(.ant-menu-submenu-title .ant-menu-item-icon),
.sider-dark .sidebar-menu :deep(.ant-menu-submenu-arrow) {
  color: var(--app-text-secondary) !important;
}

.sider-dark .sidebar-menu :deep(.ant-menu-item:hover),
.sider-dark .sidebar-menu :deep(.ant-menu-submenu-title:hover) {
  color: var(--app-text-strong) !important;
  background: var(--app-hover-bg) !important;
}

.sider-dark .sidebar-menu :deep(.ant-menu-item-selected),
.sider-dark .sidebar-menu :deep(.ant-menu-submenu-selected > .ant-menu-submenu-title) {
  color: var(--app-text-strong) !important;
}

.sider-dark .sidebar-menu :deep(.ant-menu-item-selected .ant-menu-item-icon),
.sider-dark .sidebar-menu :deep(.ant-menu-submenu-selected > .ant-menu-submenu-title .ant-menu-item-icon),
.sider-dark .sidebar-menu :deep(.ant-menu-item-selected .ant-menu-title-content),
.sider-dark .sidebar-menu :deep(.ant-menu-submenu-selected > .ant-menu-submenu-title .ant-menu-title-content) {
  color: var(--app-text-strong) !important;
}
</style>
