<template>
  <div class="front-layout">
    <!-- 顶部导航 -->
    <header class="front-header">
      <div class="header-content">
        <div class="logo" @click="handleLogoClick">
<!--          <RobotOutlined class="logo-icon" />-->
          <img :src="configStore.get('sys_logo')" alt="logo" class="logo-img" />
          <span class="logo-text">{{ sysName }}</span>
        </div>
        
        <!-- 完整前台模式才显示导航菜单 -->
        <nav v-if="isFullMode" class="nav-menu">
          <router-link 
            v-for="item in navMenus" 
            :key="item.path"
            :to="item.path" 
            class="nav-item" 
            active-class="active"
          >
            <SvgIcon :name="item.meta?.icon" />
            {{ item.meta?.title }}
          </router-link>
        </nav>
        
        <!-- profile 模式显示提示文字 -->
        <div v-else class="profile-mode-hint">
          <span>身份认证中心</span>
        </div>
        
        <div class="header-right">
          <template v-if="userStore.token">
            <a-dropdown>
              <div class="user-info">
                <a-avatar :src="userStore.user?.avatar_file_url" :size="32">
                  <template #icon><UserOutlined /></template>
                </a-avatar>
                <span class="username">{{ userStore.user?.nickname || userStore.user?.username }}</span>
                <DownOutlined />
              </div>
              <template #overlay>
                <a-menu>
                  <a-menu-item key="profile" @click="router.push('/front/profile')">
                    <UserOutlined /> 个人中心
                  </a-menu-item>
                  <a-menu-divider />
                  <a-menu-item key="logout" @click="handleLogout">
                    <LogoutOutlined /> 退出登录
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </template>
          <template v-else>
            <a-button type="primary" @click="router.push('/login')">登录</a-button>
          </template>
        </div>
      </div>
    </header>
    
    <!-- 主内容区 -->
    <main class="front-main">
      <router-view />
    </main>
    
    <!-- 底部 -->
    <footer class="front-footer">
      <p>© {{ new Date().getFullYear() }} {{ sysName }} All Rights Reserved</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { UserOutlined, DownOutlined, LogoutOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '@/store/user'
import { useConfigStore } from '@/store/config'
import SvgIcon from '@/components/SvgIcon.vue'

const router = useRouter()
const userStore = useUserStore()
const configStore = useConfigStore()

// 前台模式: 'full' = 完整前台, 'profile' = 仅个人中心
const frontMode = computed(() => configStore.get('front_mode') || 'full')
const isFullMode = computed(() => frontMode.value === 'full')

// 从路由获取导航菜单
const navMenus = computed(() => {
  const frontRoute = router.getRoutes().find(r => r.name === 'FrontLayout')
  if (!frontRoute?.children) return []
  
  return frontRoute.children
    .filter(child => child.meta?.showInNav)
    .map(child => ({
      path: `/front/${child.path}`,
      meta: child.meta
    }))
})

const sysName = computed(() => configStore.get('sys_name') || 'AI Assistant')

// Logo 点击处理
const handleLogoClick = () => {
  if (isFullMode.value) {
    router.push('/front')
  } else {
    router.push('/front/profile')
  }
}

const handleLogout = () => {
  userStore.logoutAction()
  router.push('/login')
}
</script>

<style scoped lang="less">
.front-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}
.logo-img {
  width: 32px;
  height: 32px;
}
.front-header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 100;
  
  .header-content {
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 24px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    transition: opacity 0.2s;
    
    &:hover {
      opacity: 0.8;
    }
    
    .logo-icon {
      font-size: 28px;
      color: #1890ff;
    }
    
    .logo-text {
      font-size: 20px;
      font-weight: 600;
      color: #1a1a1a;
    }
  }
  
  .profile-mode-hint {
    font-size: 16px;
    color: #666;
    font-weight: 500;
  }
  
  .nav-menu {
    display: flex;
    gap: 8px;
    
    .nav-item {
      padding: 8px 16px;
      border-radius: 6px;
      color: #666;
      text-decoration: none;
      display: flex;
      align-items: center;
      gap: 6px;
      transition: all 0.2s;
      
      &:hover {
        color: #1890ff;
        background: #e6f7ff;
      }
      
      &.active {
        color: #1890ff;
        background: #e6f7ff;
        font-weight: 500;
      }
    }
  }
  
  .header-right {
    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      padding: 4px 8px;
      border-radius: 6px;
      transition: background 0.2s;
      
      &:hover {
        background: #f5f5f5;
      }
      
      .username {
        max-width: 100px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

.front-main {
  flex: 1;
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
  padding: 24px;
}

.front-footer {
  background: #fff;
  padding: 24px;
  text-align: center;
  color: #999;
  border-top: 1px solid #e8e8e8;
  
  p {
    margin: 0;
  }
}
</style>
