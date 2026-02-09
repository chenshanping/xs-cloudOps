<template>
  <div class="config-page">
    <a-card :loading="loading">
      <a-tabs v-model:activeKey="activeTab">
        <a-tab-pane key="system" tab="系统配置">
          <SystemConfig v-if="!loading" />
        </a-tab-pane>
        <a-tab-pane key="login" tab="登录与注册">
          <LoginRegisterConfig v-if="!loading" />
        </a-tab-pane>
        <a-tab-pane key="email" tab="邮箱与安全">
          <EmailConfig v-if="!loading" />
        </a-tab-pane>
        <a-tab-pane key="ai" tab="AI配置">
          <AIConfig v-if="!loading" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useConfigStore } from '@/store/config'
import SystemConfig from './components/SystemConfig.vue'
import LoginRegisterConfig from './components/LoginRegisterConfig.vue'
import EmailConfig from './components/EmailConfig.vue'
import AIConfig from './components/AIConfig.vue'

const activeTab = ref('system')
const configStore = useConfigStore()
const loading = ref(true)

onMounted(async () => {
  // 强制刷新配置，确保获取最新数据
  await configStore.loadConfigs(true)
  loading.value = false
})
</script>

<style scoped>
.config-page {
  height: 100%;
}
</style>
