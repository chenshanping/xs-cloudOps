<template>
  <div class="no-permission-page">
    <div class="card">
      <div class="icon-wrapper">
        <LockOutlined class="lock-icon" />
      </div>
      <h1>{{ sysName }}</h1>
      <a-result
        status="403"
        title="暂无访问权限"
        sub-title="您的账号尚未分配任何功能权限，请联系管理员开通。"
      >
        <template #extra>
          <a-space>
            <a-button @click="handleRefresh">
              <template #icon><ReloadOutlined /></template>
              刷新权限
            </a-button>
            <a-button type="primary" danger @click="handleLogout">
              <template #icon><LogoutOutlined /></template>
              退出登录
            </a-button>
          </a-space>
        </template>
      </a-result>
      <div class="tips">
        <p>当前登录账号：<strong>{{ userStore.user?.username }}</strong></p>
        <p class="hint">管理员分配权限后，点击"刷新权限"即可进入系统</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { LockOutlined, ReloadOutlined, LogoutOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '@/store/user'
import { useConfigStore } from '@/store/config'

const router = useRouter()
const userStore = useUserStore()
const configStore = useConfigStore()

const sysName = configStore.get('sys_name')

const handleRefresh = async () => {
  try {
    await userStore.getUserInfoAction()
    if (userStore.menus && userStore.menus.length > 0) {
      message.success('权限已更新，正在跳转...')
      router.replace('/dashboard')
    } else {
      message.info('暂无新权限，请联系管理员')
    }
  } catch {
    message.error('刷新失败，请重试')
  }
}

const handleLogout = async () => {
  await userStore.logoutAction()
  router.replace('/login')
}
</script>

<style scoped lang="less">
.no-permission-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 24px;

  .card {
    background: #fff;
    border-radius: 16px;
    padding: 48px 40px 40px;
    max-width: 520px;
    width: 100%;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
    text-align: center;

    .icon-wrapper {
      margin-bottom: 16px;

      .lock-icon {
        font-size: 48px;
        color: #faad14;
      }
    }

    h1 {
      font-size: 22px;
      font-weight: 600;
      color: #1a1a1a;
      margin-bottom: 8px;
    }

    .tips {
      margin-top: 8px;
      padding-top: 20px;
      border-top: 1px solid #f0f0f0;

      p {
        margin: 0;
        font-size: 14px;
        color: #666;
        line-height: 2;
      }

      .hint {
        color: #999;
        font-size: 13px;
      }
    }
  }
}

@media (max-width: 480px) {
  .no-permission-page .card {
    padding: 32px 20px 28px;
  }
}
</style>
