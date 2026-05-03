<template>
  <a-layout-header
    :class="['header', { 'header-dark': headerDark }]"
    :style="headerStyle"
  >
    <div class="header-left">
      <a-button
        v-if="showCollapseTrigger"
        class="header-icon-button"
        type="text"
        @click="toggleSidebar"
      >
        <MenuUnfoldOutlined v-if="uiStore.layout.sidebarCollapsed" />
        <MenuFoldOutlined v-else />
      </a-button>

      <div v-if="showBrand" class="brand" @click="router.push('/dashboard')">
        <img :src="configStore.get('sys_logo')" alt="logo" class="brand-logo" />
        <span class="brand-title">{{ configStore.get('sys_name') }}</span>
      </div>

      <a-menu
        v-if="showTopNavigation"
        :items="topMenuItems"
        :selected-keys="selectedTopKeys"
        mode="horizontal"
        class="top-menu"
        @click="handleTopMenuClick"
      />

      <a-breadcrumb v-else class="breadcrumb">
        <a-breadcrumb-item v-for="item in breadcrumbs" :key="item.path || item.title">
          <router-link v-if="item.path" :to="item.path">{{ item.title }}</router-link>
          <span v-else>{{ item.title }}</span>
        </a-breadcrumb-item>
      </a-breadcrumb>
    </div>

    <div class="header-right">
      <a-tooltip title="刷新当前页">
        <a-button
          class="header-icon-button"
          type="text"
          @click="handleRefresh"
        >
          <ReloadOutlined />
        </a-button>
      </a-tooltip>
      <a-tooltip :title="uiStore.isDark ? '切换浅色模式' : '切换深色模式'">
        <a-button
          class="header-icon-button"
          type="text"
          @click="uiStore.updateTheme({ mode: uiStore.isDark ? 'light' : 'dark' })"
        >
          <BulbOutlined v-if="!uiStore.isDark" />
          <BgColorsOutlined v-else />
        </a-button>
      </a-tooltip>
      <a-button class="header-icon-button" type="text" @click="uiStore.toggleSettings(true)">
        <SettingOutlined />
      </a-button>

      <a-dropdown overlayClassName="header-user-dropdown">
        <span class="user-info">
          <a-avatar :size="30" :src="userStore.user?.avatar_file_url">
            {{ userStore.user?.nickname?.charAt(0) || 'U' }}
          </a-avatar>
          <span class="username">{{ userStore.user?.nickname || userStore.user?.username }}</span>
        </span>
        <template #overlay>
          <a-menu>
            <a-menu-item @click="router.push('/profile')">
              <UserOutlined />
              个人中心
            </a-menu-item>
            <a-menu-item @click="router.push('/ai')">
              <SvgIcon name="svg:aiChat" />
              AI助手
            </a-menu-item>
            <a-menu-divider />
            <a-menu-item @click="handleLogout">
              <LogoutOutlined />
              退出登录
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </a-layout-header>
</template>

<script setup lang="ts">
import { computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { ItemType, MenuInfo } from 'ant-design-vue/es/menu/src/interface'
import {
  BgColorsOutlined,
  BulbOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ReloadOutlined,
  SettingOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import SvgIcon from '@/components/SvgIcon.vue'
import { useConfigStore } from '@/store/config'
import { useUiStore } from '@/store/ui'
import { useUserStore } from '@/store/user'
import {
  filterEnabledMenus,
  filterVisibleMenus,
  findMenuTrail,
  firstNavigablePath,
  getBreadcrumbs,
  normalizePath,
} from './layout-menu'
import MenuIcon from './MenuIcon.vue'

interface BreadcrumbItem {
  path?: string
  title: string
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const configStore = useConfigStore()
const uiStore = useUiStore()

const normalizedMenus = computed(() => filterEnabledMenus(userStore.menus || []))
const visibleMenus = computed(() => filterVisibleMenus(userStore.menus || []))

const breadcrumbs = computed<BreadcrumbItem[]>(() => {
  const items: BreadcrumbItem[] = [{ path: '/dashboard', title: '首页' }]

  if (route.path === '/dashboard') {
    return items
  }

  const trail = getBreadcrumbs(normalizedMenus.value, route.path)
  if (trail.length) {
    return [...items, ...trail]
  }

  if (route.path === '/profile') {
    return [...items, { title: '个人中心' }]
  }

  return items
})

const activeTrail = computed(() => findMenuTrail(normalizedMenus.value, route.path))

const selectedTopKeys = computed(() => {
  const top = activeTrail.value[0]
  return top ? [String(top.id)] : []
})

const buildTopMenuChildren = (menus: typeof normalizedMenus.value[number]['children'] = []): ItemType[] => {
  return menus.map((menu) => {
    const icon = menu.icon ? () => h(MenuIcon, { icon: menu.icon }) : undefined
    if (menu.type === 1 && menu.children?.length) {
      return {
        key: String(menu.id),
        label: menu.name,
        icon,
        children: buildTopMenuChildren(menu.children),
      }
    }
    return {
      key: normalizePath(menu.path),
      label: menu.name,
      icon,
    }
  })
}

const topMenuItems = computed<ItemType[]>(() => {
  return visibleMenus.value.map((menu) => {
    const icon = menu.icon ? () => h(MenuIcon, { icon: menu.icon }) : undefined

    if (menu.type === 1 && menu.children?.length) {
      return {
        key: String(menu.id),
        label: menu.name,
        icon,
        children: buildTopMenuChildren(menu.children),
      }
    }

    return {
      key: normalizePath(menu.path),
      label: menu.name,
      icon,
    }
  })
})

const showTopNavigation = computed(() => uiStore.layout.mode === 'top' || uiStore.layout.mode === 'mixed')
const showCollapseTrigger = computed(() => uiStore.effectiveShowSidebar && uiStore.layout.mode !== 'top')
const showBrand = computed(() => uiStore.layout.mode !== 'sidebar' || !uiStore.effectiveShowSidebar)

const headerDark = computed(() => uiStore.isDark || uiStore.theme.headerDark)

const headerStyle = computed(() => {
  const background = headerDark.value ? 'var(--app-elevated-bg)' : 'var(--app-surface-color)'
  const color = headerDark.value ? 'var(--app-text-strong)' : 'var(--app-text-strong)'

  return {
    background,
    color,
    borderBottom: `1px solid ${headerDark.value ? 'var(--app-border-color)' : 'var(--app-border-color)'}`,
    boxShadow: headerDark.value ? '0 1px 0 rgba(148, 163, 184, 0.08)' : '0 1px 4px rgba(15, 23, 42, 0.06)',
  }
})

const toggleSidebar = () => {
  uiStore.updateLayout({ sidebarCollapsed: !uiStore.layout.sidebarCollapsed })
}

const handleTopMenuClick = ({ key }: MenuInfo) => {
  if (typeof key === 'string' && key.startsWith('/')) {
    router.push(key)
    return
  }

  const target = visibleMenus.value.find((menu) => String(menu.id) === String(key))
  if (!target) {
    return
  }

  const destination = firstNavigablePath(target)
  if (destination) {
    router.push(destination)
  }
}

const handleLogout = () => {
  userStore.logoutAction()
  router.push('/login')
}

const handleRefresh = () => {
  window.location.reload()
}
</script>

<style scoped>
.header {
  height: 64px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.header-dark {
  color: var(--app-text-strong);
}

.header-dark .brand-title,
.header-dark .username {
  color: var(--app-text-strong);
}

.header-left,
.header-right {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left {
  flex: 1;
}

.header-icon-button {
  width: 36px;
  height: 36px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
}

.header-dark .header-icon-button {
  color: var(--app-text-strong);
}

.header-dark .header-icon-button:hover,
.header-dark .header-icon-button:focus {
  color: var(--app-text-strong);
  background: var(--app-hover-bg);
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  cursor: pointer;
}

.brand-logo {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  object-fit: cover;
}

.brand-title {
  font-size: 16px;
  font-weight: 600;
  white-space: nowrap;
}

.top-menu {
  flex: 1;
  min-width: 0;
  border-bottom: none;
  background: transparent;
}

.top-menu :deep(.ant-menu-item) {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border-radius: 10px;
}

.header-dark .top-menu :deep(.ant-menu-overflow-item),
.header-dark .top-menu :deep(.ant-menu-item),
.header-dark .top-menu :deep(.ant-menu-submenu-title),
.header-dark .top-menu :deep(.ant-menu-title-content),
.header-dark .top-menu :deep(.ant-menu-item .ant-menu-item-icon),
.header-dark .top-menu :deep(.ant-menu-submenu-title .ant-menu-item-icon) {
  color: var(--app-text-secondary) !important;
}

.header-dark .top-menu :deep(.ant-menu-item:hover),
.header-dark .top-menu :deep(.ant-menu-submenu-title:hover) {
  color: var(--app-text-strong) !important;
  background: var(--app-hover-bg);
}

.header-dark .top-menu :deep(.ant-menu-item-selected),
.header-dark .top-menu :deep(.ant-menu-submenu-selected > .ant-menu-submenu-title) {
  color: var(--app-text-strong) !important;
  background: var(--app-hover-bg);
}

.header-dark .top-menu :deep(.ant-menu-horizontal) {
  background: transparent;
  border-bottom: none;
}

.header-dark .top-menu :deep(.ant-menu-item::after),
.header-dark .top-menu :deep(.ant-menu-submenu::after) {
  border-bottom-color: transparent !important;
}

.breadcrumb {
  min-width: 0;
}

.header-dark .breadcrumb :deep(.ant-breadcrumb-link),
.header-dark .breadcrumb :deep(.ant-breadcrumb-link a) {
  color: var(--app-text-secondary);
}

.header-dark .breadcrumb :deep(.ant-breadcrumb-link a:hover) {
  color: var(--app-text-strong);
}

.header-dark .breadcrumb :deep(.ant-breadcrumb-separator) {
  color: var(--app-text-muted);
}

.header-dark .breadcrumb :deep(.ant-breadcrumb > span:last-child .ant-breadcrumb-link),
.header-dark .breadcrumb :deep(.ant-breadcrumb > span:last-child .ant-breadcrumb-link a) {
  color: var(--app-text-strong);
}

.user-info {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 4px 10px;
  border-radius: 999px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.user-info:hover {
  background: rgba(148, 163, 184, 0.14);
}

.header-dark .user-info:hover {
  background: var(--app-hover-bg);
}

.username {
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>

<style>
.header-user-dropdown .ant-dropdown-menu {
  min-width: 140px;
}
</style>
