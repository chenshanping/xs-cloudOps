<template>
  <div class="captcha-wrapper">
    <a-input
      v-model:value="captchaCode"
      :placeholder="placeholder"
      :size="size"
      @change="handleChange"
    >
      <template #prefix>
        <SafetyCertificateOutlined />
      </template>
    </a-input>
    <div class="captcha-image" @click="refreshCaptcha">
      <img v-if="captchaImage" :src="captchaImage" alt="验证码" />
      <div v-else class="captcha-loading">
        <LoadingOutlined v-if="loading" />
        <span v-else>点击获取</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, defineExpose } from 'vue'
import { SafetyCertificateOutlined, LoadingOutlined } from '@ant-design/icons-vue'
import { getCaptcha } from '@/api/captcha'

const props = withDefaults(defineProps<{
  modelValue?: string
  captchaId?: string
  placeholder?: string
  size?: 'large' | 'default' | 'small'
  autoLoad?: boolean
}>(), {
  modelValue: '',
  captchaId: '',
  placeholder: '请输入验证码',
  size: 'large',
  autoLoad: true
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'update:captchaId', value: string): void
}>()

const captchaCode = ref(props.modelValue)
const captchaImage = ref('')
const loading = ref(false)
const currentCaptchaId = ref(props.captchaId)

// 监听外部传入的值
watch(() => props.modelValue, (val:any) => {
  captchaCode.value = val
})

// 刷新验证码
const refreshCaptcha = async () => {
  loading.value = true
  try {
    const res = await getCaptcha()
    captchaImage.value = res.data.captcha_image
    currentCaptchaId.value = res.data.captcha_id
    emit('update:captchaId', res.data.captcha_id)
    // 清空输入
    captchaCode.value = ''
    emit('update:modelValue', '')
  } catch (e) {
    console.error('获取验证码失败', e)
  } finally {
    loading.value = false
  }
}

// 输入变化
const handleChange = () => {
  emit('update:modelValue', captchaCode.value)
}

// 获取验证码ID
const getCaptchaId = () => {
  return currentCaptchaId.value
}

// 暴露方法
defineExpose({
  refresh: refreshCaptcha,
  getCaptchaId
})

onMounted(() => {
  if (props.autoLoad) {
    refreshCaptcha()
  }
})
</script>

<style scoped>
.captcha-wrapper {
  display: flex;
  gap: 12px;
}

.captcha-wrapper :deep(.ant-input-affix-wrapper) {
  flex: 1;
}

.captcha-image {
  width: 120px;
  height: 40px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
  flex-shrink: 0;
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.captcha-loading {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  color: #999;
  font-size: 12px;
}
</style>
