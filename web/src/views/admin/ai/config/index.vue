<template>
  <div class="ai-config-page">
    <a-card>
      <a-alert
        v-if="hasUnsavedChanges"
        class="unsaved-alert"
        type="warning"
        show-icon
        message="当前配置有未保存修改"
      >
        <template #action>
          <a-button
            v-permission="'ai:config:save'"
            type="primary"
            size="small"
            :loading="quickSaving"
            @click="handleQuickSave"
          >
            保存当前页面
          </a-button>
        </template>
      </a-alert>

      <AIConfig
        ref="aiConfigRef"
        @dirty-change="setDirtyState"
      />
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
          <a-button type="primary" @click="resolveLeavePrompt('save')" v-permission="'ai:config:save'">保存并继续</a-button>
        </a-space>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import AIConfig from '@/views/admin/system/config/components/AIConfig.vue'
import type { ConfigTabGuardHandle } from '@/views/admin/system/config/config-tab-guard'

type LeavePromptAction = 'save' | 'discard' | 'cancel'

const aiConfigRef = ref<ConfigTabGuardHandle | null>(null)
const dirty = ref(false)
const quickSaving = ref(false)
const leavePromptOpen = ref(false)
const leavePromptResolver = ref<((action: LeavePromptAction) => void) | null>(null)

const hasUnsavedChanges = computed(() => dirty.value)

const setDirtyState = (value: boolean) => {
  dirty.value = value
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

const saveCurrentPage = async () => {
  if (!aiConfigRef.value) {
    return true
  }
  return aiConfigRef.value.save()
}

const discardCurrentPageChanges = () => {
  aiConfigRef.value?.discardChanges()
}

const closeTransientUi = () => {
  aiConfigRef.value?.closeTransientUi()
}

const handleQuickSave = async () => {
  if (!hasUnsavedChanges.value) {
    return
  }
  quickSaving.value = true
  try {
    await saveCurrentPage()
  } finally {
    quickSaving.value = false
  }
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

  const action = await promptUnsavedChanges()
  if (action === 'cancel') {
    return false
  }
  if (action === 'discard') {
    discardCurrentPageChanges()
    return true
  }

  const saved = await saveCurrentPage()
  if (!saved) {
    return false
  }
  closeTransientUi()
  return true
})

onMounted(() => {
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
})
</script>

<style scoped>
.ai-config-page {
  height: 100%;
}

.unsaved-alert {
  margin-bottom: 16px;
}

.leave-prompt-text {
  margin: 0;
  color: #595959;
}
</style>
