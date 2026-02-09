<template>
  <AuthLayout>
    <div class="form-header">
      <h2>{{ configStore.get('login_title') }}</h2>
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

      <!-- 验证码 -->
      <a-form-item v-if="captchaConfig?.login_captcha_enabled" name="captcha">
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
import { getCaptchaConfig } from '@/api/captcha'

interface CaptchaConfig {
  login_captcha_enabled: boolean
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
  register_captcha_enabled: false,
  register_email_verify: false,
})
const loading = ref(false)
const rememberMe = ref(false)
const captchaRef = ref<InstanceType<typeof Captcha> | null>(null)

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
  
  if (captchaConfig.value.login_captcha_enabled) {
    baseRules.captcha = [{ required: true, message: '请输入验证码' }]
  }
  
  return baseRules
})

const handleLogin = async () => {
  loading.value = true
  try {
    const loginData: any = {
      username: formState.username,
      password: formState.password,
    }
    
    if (captchaConfig?.value?.login_captcha_enabled && captchaRef.value) {
      loginData.captcha_id = captchaRef.value.getCaptchaId()
      loginData.captcha = formState.captcha
    }

    await userStore.loginAction(loginData)
    await userStore.getUserInfoAction()
    message.success('登录成功')

    const roles = userStore.user?.roles || []
    console.log(roles)
    if (roles.length > 0 && roles[0].code === 'user') {
      router.push('/front')
    } else {
      router.push('/')
    }
  } catch (error) {
    captchaRef.value?.refresh()
  } finally {
    loading.value = false
  }
}

onMounted(async() => {
  // 加载系统配置
  await configStore.loadConfigs()
  
  try {
    const captchaRes = await getCaptchaConfig()
    captchaConfig.value = {
      login_captcha_enabled: captchaRes.data?.login_captcha_enabled,
      register_captcha_enabled: captchaRes.data?.register_captcha_enabled,
      register_email_verify: captchaRes.data?.register_email_verify,
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

  h2 {
    font-size: 26px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0 0 8px 0;
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
</style>
