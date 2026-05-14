<template>
  <PageWrapper class="config-page">
    <div class="config-page__content">
      <a-card :loading="loading">
      <a-alert
        v-if="hasUnsavedChanges"
        class="unsaved-alert"
        type="warning"
        show-icon
        message="当前配置有未保存修改"
      >
        <template #action>
          <a-button
            type="primary"
            size="small"
            :disabled="!currentTabDirty || quickSaving"
            :loading="quickSaving"
            @click="handleQuickSave"
            v-permission="'system:config:edit'"
          >
            保存当前页签
          </a-button>
        </template>
      </a-alert>

      <a-tabs :activeKey="activeTab" @change="handleTabChange">
        <a-tab-pane key="basic">
          <template #tab>
            <span class="config-tab-title">
              基础配置
              <span v-if="dirtyState.basic" class="dirty-mark">*</span>
            </span>
          </template>
          <SystemConfig
            v-if="!loading"
            ref="basicTabRef"
            @dirty-change="setTabDirty('basic', $event)"
          />
        </a-tab-pane>
        <a-tab-pane key="file">
          <template #tab>
            <span class="config-tab-title">
              文件设置
              <span v-if="dirtyState.file" class="dirty-mark">*</span>
            </span>
          </template>
          <FileSettings
            v-if="!loading"
            ref="fileTabRef"
            @dirty-change="setTabDirty('file', $event)"
          />
        </a-tab-pane>
        <a-tab-pane key="login">
          <template #tab>
            <span class="config-tab-title">
              登录与注册
              <span v-if="dirtyState.login" class="dirty-mark">*</span>
            </span>
          </template>
          <LoginRegisterConfig
            v-if="!loading"
            ref="loginTabRef"
            @dirty-change="setTabDirty('login', $event)"
          />
        </a-tab-pane>
        <a-tab-pane key="email">
          <template #tab>
            <span class="config-tab-title">
              邮箱与安全
              <span v-if="dirtyState.email" class="dirty-mark">*</span>
            </span>
          </template>
          <EmailConfig
            v-if="!loading"
            ref="emailTabRef"
            @dirty-change="setTabDirty('email', $event)"
          />
        </a-tab-pane>
      </a-tabs>
      </a-card>

      <a-modal
        :open="leavePromptOpen"
        title="当前配置有未保存修改"
        :maskClosable="false"
        :closable="false"
        @cancel="resolveLeavePrompt('cancel')"
      >
        <p class="leave-prompt-text">继续操作前请先确认是否保存当前修改。</p>

        <template #footer>
          <a-space>
            <a-button @click="resolveLeavePrompt('cancel')">取消</a-button>
            <a-button danger @click="resolveLeavePrompt('discard')">直接离开</a-button>
            <a-button type="primary" @click="resolveLeavePrompt('save')" v-permission="'system:config:edit'">保存并继续</a-button>
          </a-space>
        </template>
      </a-modal>
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import PageWrapper from '@/components/page/PageWrapper.vue'
import { useConfigStore } from '@/store/config'
import SystemConfig from './components/SystemConfig.vue'
import FileSettings from './components/FileSettings.vue'
import LoginRegisterConfig from './components/LoginRegisterConfig.vue'
import EmailConfig from './components/EmailConfig.vue'
import type { ConfigTabGuardHandle, ConfigTabKey } from './config-tab-guard'
import { hasDirtyConfigTabs } from './config-tab-guard'

type LeavePromptAction = 'save' | 'discard' | 'cancel'

const CONFIG_TAB_KEYS: ConfigTabKey[] = ['basic', 'file', 'login', 'email']

const activeTab = ref<ConfigTabKey>('basic')
const configStore = useConfigStore()
const loading = ref(true)
const quickSaving = ref(false)
const leavePromptOpen = ref(false)
const leavePromptResolver = ref<((action: LeavePromptAction) => void) | null>(null)

const dirtyState = reactive<Record<ConfigTabKey, boolean>>({
  basic: false,
  file: false,
  login: false,
  email: false,
  ai: false,
})

const basicTabRef = ref<ConfigTabGuardHandle | null>(null)
const fileTabRef = ref<ConfigTabGuardHandle | null>(null)
const loginTabRef = ref<ConfigTabGuardHandle | null>(null)
const emailTabRef = ref<ConfigTabGuardHandle | null>(null)

const hasUnsavedChanges = computed(() => hasDirtyConfigTabs(dirtyState))
const currentTabDirty = computed(() => dirtyState[activeTab.value])

const getTabHandle = (tab: ConfigTabKey): ConfigTabGuardHandle | null => {
  switch (tab) {
    case 'basic':
      return basicTabRef.value
    case 'file':
      return fileTabRef.value
    case 'login':
      return loginTabRef.value
    case 'email':
      return emailTabRef.value
    default:
      return null
  }
}

const setTabDirty = (tab: ConfigTabKey, value: boolean) => {
  dirtyState[tab] = value
}

const promptUnsavedChanges = () => {
  leavePromptOpen.value = true
  return new Promise<LeavePromptAction>((resolve) => {
    leavePromptResolver.value = resolve
  })
}

const resolveLeavePrompt = (action: LeavePromptAction) => {
  leavePromptOpen.value = false
  const resolve = leavePromptResolver.value
  leavePromptResolver.value = null
  resolve?.(action)
}

const closeCurrentTransientUi = () => {
  getTabHandle(activeTab.value)?.closeTransientUi()
}

const saveCurrentTab = async () => {
  const currentTab = getTabHandle(activeTab.value)
  if (!currentTab) {
    return true
  }
  return currentTab.save()
}

const discardCurrentTabChanges = () => {
  getTabHandle(activeTab.value)?.discardChanges()
}

const handleQuickSave = async () => {
  if (!currentTabDirty.value) {
    return
  }
  quickSaving.value = true
  try {
    await saveCurrentTab()
  } finally {
    quickSaving.value = false
  }
}

const confirmLeaveCurrentTab = async () => {
  const currentTab = getTabHandle(activeTab.value)
  if (!currentTab?.isDirty()) {
    return true
  }

  const action = await promptUnsavedChanges()
  if (action === 'cancel') {
    return false
  }
  if (action === 'discard') {
    discardCurrentTabChanges()
    return true
  }

  const saved = await saveCurrentTab()
  if (!saved) {
    return false
  }
  closeCurrentTransientUi()
  return true
}

const handleTabChange = async (nextKey: string) => {
  if (!CONFIG_TAB_KEYS.includes(nextKey as ConfigTabKey)) {
    return
  }

  const nextTab = nextKey as ConfigTabKey
  if (nextTab === activeTab.value) {
    return
  }

  const canLeave = await confirmLeaveCurrentTab()
  if (!canLeave) {
    return
  }

  closeCurrentTransientUi()
  activeTab.value = nextTab
}

const handleBeforeUnload = (event: BeforeUnloadEvent) => {
  if (!hasUnsavedChanges.value) {
    return
  }
  event.preventDefault()
  event.returnValue = ''
}

onBeforeRouteLeave(async () => {
  if (!hasUnsavedChanges.value) {
    return true
  }

  const currentTab = getTabHandle(activeTab.value)
  const action = await promptUnsavedChanges()
  if (action === 'cancel') {
    return false
  }
  if (action === 'discard') {
    currentTab?.discardChanges()
    return true
  }

  const saved = await saveCurrentTab()
  if (!saved) {
    return false
  }
  closeCurrentTransientUi()
  return true
})

onMounted(async () => {
  window.addEventListener('beforeunload', handleBeforeUnload)
  await configStore.loadConfigs(true, 'all')
  if (window.location.search.includes('tab=storage') || window.location.search.includes('tab=file')) {
    activeTab.value = 'file'
  }
  loading.value = false
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
})
</script>

<style scoped>
.config-page__content {
  min-width: 0;
}

.unsaved-alert {
  margin-bottom: 16px;
}

.config-tab-title {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.dirty-mark {
  color: #ff4d4f;
  font-weight: 700;
}

.leave-prompt-text {
  margin: 0;
  color: #595959;
}
</style>
