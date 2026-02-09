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

const router = useRouter()
const error = ref<string>('')

// API 错误关键词，这些错误已经由 request.ts 处理并显示 message，不需要显示错误页面
const apiErrorKeywords = [
  '服务器内部错误',
  '网关错误',
  '服务不可用',
  '网关超时',
  '请求失败',
  '请求超时',
  '网络已断开',
  '无法连接服务器',
  '网络错误',
  'Token',
]

onErrorCaptured((err: any) => {
  const errMsg = err.message || '未知错误'
  
  // 检查是否是 API 错误，如果是则不显示错误页面
  const isApiError = apiErrorKeywords.some(keyword => errMsg.includes(keyword))
  if (isApiError) {
    console.warn('API错误已由 request.ts 处理:', errMsg)
    return false // 阻止错误继续向上传播，但不显示错误页面
  }
  
  console.error('捕获到组件错误:', err)
  error.value = errMsg
  return false // 阻止错误继续向上传播
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
