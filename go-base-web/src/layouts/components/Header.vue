<template>
  <a-layout-header 
    class="header" 
    :style="{ 
      background: configStore.get('header_bg_color'),
      color: configStore.get('header_text_color')
    }"
  >
    <div class="header-left">
      <!-- 面包屑 -->
      <a-breadcrumb class="breadcrumb">
        <a-breadcrumb-item v-for="item in breadcrumbs" :key="item.path || item.title">
          <router-link v-if="item.path" :to="item.path">{{ item.title }}</router-link>
          <span v-else>{{ item.title }}</span>
        </a-breadcrumb-item>
      </a-breadcrumb>
    </div>
    <div class="header-right">
      <a-dropdown overlayClassName="header-user-dropdown">
        <span class="user-info">
          <a-avatar :size="28" :src="userStore.user?.avatar_file_url">{{ userStore.user?.nickname?.charAt(0) || 'U' }}</a-avatar>

          <span class="username">{{ userStore.user?.nickname || userStore.user?.username }}</span>
        </span>
        <template #overlay>
          <a-menu>
            <a-menu-item @click="$router.push('/profile')">
              <UserOutlined />
              个人中心
            </a-menu-item>
            <a-menu-item @click="$router.push('/ai')">
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
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useConfigStore } from '@/store/config'
import { UserOutlined, LogoutOutlined } from '@ant-design/icons-vue'
import aiChatIcon from '@/assets/icons/aiChat.svg?url'
import SvgIcon from '@/components/SvgIcon.vue'
interface Breadcrumb {
  path?: string
  title: string
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const configStore = useConfigStore()

const breadcrumbs = computed<Breadcrumb[]>(() => {
  const items: Breadcrumb[] = [{ path: '/dashboard', title: '首页' }]
  
  if (route.path === '/dashboard') return items
  
  const menus = userStore.menus || []
  for (const menu of menus) {
    if (menu.path === route.path) {
      items.push({ title: menu.name })
      break
    }
    if (menu.children) {
      for (const child of menu.children) {
        if (child.path === route.path) {
          items.push({ title: menu.name })
          items.push({ title: child.name })
          break
        }
      }
    }
  }
  
  if (route.path === '/profile') {
    items.push({ title: '个人中心' })
  }
  
  return items
})

const handleLogout = () => {
  userStore.logoutAction()
  router.push('/login')
}
</script>

<style scoped>
.header {
  height: 48px;
  line-height: 48px;
  padding: 0 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
}

.breadcrumb {
  font-size: 14px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 12px;
  border-radius: 4px;
  transition: background 0.3s;
}

.user-info:hover {
  background: #f5f5f5;
}

.username {
  margin-left: 8px;
  font-size: 14px;
}
</style>

<style>
/* 下拉菜单全局样式，防止收缩 */
.header-user-dropdown .ant-dropdown-menu {
  min-width: 120px;
}
</style>
