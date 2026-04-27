<template>
  <div v-if="error" class="error-container">
    <a-result
      status="error"
      title="页面加载失败"
      :sub-title="error"
    >
      <template #extra>
        <a-button type="primary" @click="reload">
          刷新页面
        </a-button>
        <a-button @click="goBack">
          返回上一页
        </a-button>
      </template>
    </a-result>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue'
import { useRouter } from 'vue-router'

type CapturedBoundaryError = Error & {
  handledByMessage?: boolean
  errorSource?: string
}

const router = useRouter()
const error = ref<string>('')

onErrorCaptured((err: CapturedBoundaryError) => {
  const errMsg = err?.message || '未知错误'

  // request.ts 已经弹过消息的接口错误，不再升级成整页错误页
  if (err?.handledByMessage || err?.errorSource === 'request') {
    console.warn('接口错误已由消息提示处理:', errMsg)
    return false
  }

  console.error('捕获到组件错误:', err)
  error.value = errMsg
  return false
})

const reload = () => {
  error.value = ''
  window.location.reload()
}

const goBack = () => {
  error.value = ''
  router.back()
}
</script>

<style scoped>
.error-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}
</style>
