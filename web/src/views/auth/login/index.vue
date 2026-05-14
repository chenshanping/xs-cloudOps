<template>
  <AuthLayout>
    <div class="form-header">
      <div class="form-header__title-row">
        <img
          v-if="configStore.get('sys_logo')"
          :src="configStore.get('sys_logo')"
          alt="Logo"
          class="form-header__logo"
        />
        <h2>{{ configStore.get('login_title') }}</h2>
      </div>
      <p>登录您的账户</p>
    </div>

    <a-form
      :model="formState"
      :rules="rules"
      @finish="handleLogin"
      layout="vertical"
    >
      <a-form-item name="username">
        <a-input
          v-model:value="formState.username"
          placeholder="请输入用户名"
          size="large"
        >
          <template #prefix><UserOutlined /></template>
        </a-input>
      </a-form-item>

      <a-form-item name="password">
        <a-input-password
          v-model:value="formState.password"
          placeholder="请输入密码"
          size="large"
        >
          <template #prefix><LockOutlined /></template>
        </a-input-password>
      </a-form-item>

      <!-- 图形验证码 -->
      <a-form-item v-if="captchaConfig?.login_captcha_enabled && captchaConfig?.login_captcha_type !== 'slider'" name="captcha">
        <Captcha
          ref="captchaRef"
          v-model="formState.captcha"
        />
      </a-form-item>
      
      <div class="form-options">
        <a-checkbox v-model:checked="rememberMe">记住我</a-checkbox>
        <router-link v-if="enableRegister" to="/forgot-password" class="forgot-link">忘记密码？</router-link>
      </div>

      <a-form-item>
        <a-button
          type="primary"
          html-type="submit"
          size="large"
          block
          :loading="loading"
        >
          登 录
        </a-button>
      </a-form-item>
    </a-form>

    <div class="form-footer" v-if="enableRegister">
      <span>还没有账号？</span>
      <router-link to="/register">立即注册</router-link>
    </div>
    
    <!-- 滑动验证码弹窗 -->
    <a-modal
      v-model:open="sliderModalVisible"
      title="安全验证"
      :footer="null"
      :width="340"
      centered
      :maskClosable="false"
    >
      <div class="slider-modal-content">
        <SliderCaptcha
          ref="sliderCaptchaRef"
          :background-image="sliderBgImage"
          @success="handleSliderSuccess"
          @fail="handleSliderFail"
        />
      </div>
    </a-modal>
  </AuthLayout>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useConfigStore } from '@/store/config'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import Captcha from '@/components/Captcha.vue'
import SliderCaptcha from '@/components/SliderCaptcha.vue'
import { getCaptchaConfig } from '@/api/captcha'

interface CaptchaConfig {
  login_captcha_enabled: boolean
  login_captcha_type: string
  register_captcha_enabled: boolean
  register_email_verify: boolean
}

const router = useRouter()
const userStore = useUserStore()
const configStore = useConfigStore()

// 是否开启注册功能
const enableRegister = computed(() => configStore.get('enable_register') === 'true')

const captchaConfig = ref<CaptchaConfig>({
  login_captcha_enabled: false,
  login_captcha_type: 'digit',
  register_captcha_enabled: false,
  register_email_verify: false,
})
const loading = ref(false)
const rememberMe = ref(false)
const captchaRef = ref<InstanceType<typeof Captcha> | null>(null)
const sliderCaptchaRef = ref<InstanceType<typeof SliderCaptcha> | null>(null)
const sliderVerified = ref(false)
const sliderCaptchaId = ref('')
const sliderModalVisible = ref(false)
const sliderBgImage = ref('')
const pendingLoginData = ref<any>(null)

const formState = reactive({
  username: '',
  password: '',
  captcha: '',
})

// 动态验证规则：只在验证码启用时才验证 captcha 字段
const rules = computed(() => {
  const baseRules: any = {
    username: [{ required: true, message: '请输入用户名' }],
    password: [{ required: true, message: '请输入密码' }],
  }
  
  // 滑动验证码不需要表单验证，其他类型需要
  if (captchaConfig.value.login_captcha_enabled && captchaConfig.value.login_captcha_type !== 'slider') {
    baseRules.captcha = [{ required: true, message: '请输入验证码' }]
  }
  
  return baseRules
})

// 滑动验证码成功
const handleSliderSuccess = async (captchaId: string) => {
  sliderVerified.value = true
  sliderCaptchaId.value = captchaId
  sliderModalVisible.value = false
  
  // 验证成功后自动提交登录
  if (pendingLoginData.value) {
    await doLogin({
      ...pendingLoginData.value,
      captcha_id: captchaId,
      captcha: 'slider_verified'
    })
  }
}

// 滑动验证码失败
const handleSliderFail = () => {
  sliderVerified.value = false
  sliderCaptchaId.value = ''
}

// 实际登录请求
const doLogin = async (loginData: any) => {
  loading.value = true
  try {
    await userStore.loginAction(loginData)
    await userStore.refreshAccessAction()
    message.success('登录成功')

    const roles = userStore.user?.roles || []
    if (roles.length > 0 && roles[0].code === 'user') {
      router.push('/front')
    } else {
      router.push('/')
    }
  } catch (error) {
    // 刷新验证码
    if (captchaConfig.value.login_captcha_type === 'slider') {
      sliderCaptchaRef.value?.refresh()
      sliderVerified.value = false
      sliderCaptchaId.value = ''
    } else {
      captchaRef.value?.refresh()
    }
  } finally {
    loading.value = false
    pendingLoginData.value = null
  }
}

const handleLogin = async () => {
  const loginData: any = {
    username: formState.username,
    password: formState.password,
  }
  
  // 滑动验证码 - 弹出Modal
  if (captchaConfig.value.login_captcha_enabled && captchaConfig.value.login_captcha_type === 'slider') {
    pendingLoginData.value = loginData
    sliderModalVisible.value = true
    // 刷新滑动验证码
    setTimeout(() => {
      sliderCaptchaRef.value?.refresh()
    }, 100)
    return
  }
  
  // 普通图形验证码
  if (captchaConfig?.value?.login_captcha_enabled && captchaRef.value) {
    loginData.captcha_id = captchaRef.value.getCaptchaId()
    loginData.captcha = formState.captcha
  }

  await doLogin(loginData)
}

onMounted(async() => {
  // 加载系统配置
  await configStore.loadConfigs(false, 'public')
  
  try {
    const captchaRes = await getCaptchaConfig()
    console.log('[Login] captcha config:', captchaRes.data)
    captchaConfig.value = {
      login_captcha_enabled: captchaRes.data?.login_captcha_enabled,
      login_captcha_type: captchaRes.data?.login_captcha_type || 'digit',
      register_captcha_enabled: captchaRes.data?.register_captcha_enabled,
      register_email_verify: captchaRes.data?.register_email_verify,
    }
    // 滑动验证码背景图
    if (captchaRes.data?.slider_captcha_bg) {
      sliderBgImage.value = captchaRes.data.slider_captcha_bg
      console.log('[Login] slider bg image:', sliderBgImage.value)
    }
  } catch {
    // 后端不可用时使用默认配置
  }
  
  // 从路由参数中获取用户名（注册后跳转过来的）
  const username = router.currentRoute.value.query.username as string
  if (username) {
    formState.username = username
  }
})
</script>
<style scoped lang="scss">
.form-header {
  text-align: center;
  margin-bottom: 32px;

  .form-header__title-row {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    margin-bottom: 8px;
  }

  .form-header__logo {
    width: 32px;
    height: 32px;
    border-radius: 8px;
    object-fit: cover;
    box-shadow: 0 6px 16px rgba(102, 126, 234, 0.18);
  }

  h2 {
    font-size: 26px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0;
  }

  p {
    color: #666;
    margin: 0;
  }
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  :deep(.ant-checkbox-wrapper) {
    color: #333;
  }

  .forgot-link {
    color: #667eea;
    font-size: 14px;

    &:hover {
      color: #764ba2;
    }
  }
}

.form-footer {
  text-align: center;
  margin-top: 24px;
  color: #666;

  a {
    color: #667eea;
    margin-left: 8px;

    &:hover {
      text-decoration: underline;
    }
  }
}

:deep(.ant-input-affix-wrapper) {
  border-radius: 8px;
}

:deep(.ant-btn) {
  height: 44px;
  border-radius: 8px;
  font-size: 15px;
}

.slider-modal-content {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}
</style>
