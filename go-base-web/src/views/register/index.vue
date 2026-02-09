<template>
  <AuthLayout>
    <div class="form-header">
      <h2>创建账户</h2>
      <p>填写以下信息完成注册</p>
    </div>

    <a-form
      ref="formRef"
      :model="formState"
      :rules="rules"
      @finish="handleRegister"
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

      <!-- 邮箱（仅在启用邮箱验证时显示） -->
      <a-form-item v-if="captchaConfig?.register_email_verify" name="email">
        <a-input
          v-model:value="formState.email"
          placeholder="请输入邮箱"
          size="large"
        >
          <template #prefix><MailOutlined /></template>
        </a-input>
      </a-form-item>

      <!-- 邮箱验证码 -->
      <a-form-item v-if="captchaConfig?.register_email_verify" name="email_code">
        <div class="email-code-row">
          <a-input
            v-model:value="formState.email_code"
            placeholder="请输入邮箱验证码"
            size="large"
          >
            <template #prefix><SafetyOutlined /></template>
          </a-input>
          <a-button
            size="large"
            :disabled="emailCodeCountdown > 0"
            :loading="sendingEmailCode"
            @click="handleSendEmailCode"
          >
            {{ emailCodeCountdown > 0 ? `${emailCodeCountdown}s` : '发送验证码' }}
          </a-button>
        </div>
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

      <a-form-item name="confirm_password">
        <a-input-password
          v-model:value="formState.confirm_password"
          placeholder="请确认密码"
          size="large"
        >
          <template #prefix><LockOutlined /></template>
        </a-input-password>
      </a-form-item>

      <!-- 图形验证码 -->
      <a-form-item v-if="captchaConfig?.register_captcha_enabled" name="captcha">
        <Captcha ref="captchaRef" v-model="formState.captcha" />
      </a-form-item>

      <a-form-item>
        <a-button
          type="primary"
          html-type="submit"
          size="large"
          block
          :loading="loading"
        >
          注 册
        </a-button>
      </a-form-item>
    </a-form>

    <div class="form-footer">
      <span>已有账户？</span>
      <router-link to="/login">立即登录</router-link>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { reactive, ref, inject, computed, type Ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  UserOutlined,
  LockOutlined,
  MailOutlined,
  SafetyOutlined,
} from '@ant-design/icons-vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import { register, sendEmailCode, type RegisterParams } from '@/api/auth'
import AuthLayout from '@/layouts/AuthLayout.vue'
import Captcha from '@/components/Captcha.vue'
import { getCaptchaConfig } from '@/api/captcha'

interface CaptchaConfig {
  login_captcha_enabled: boolean
  register_captcha_enabled: boolean
  register_email_verify: boolean
}

const router = useRouter()
const formRef = ref<FormInstance>()
const captchaRef = ref<InstanceType<typeof Captcha>>()
const captchaConfig = ref<CaptchaConfig>({
  login_captcha_enabled: false,
  register_captcha_enabled: false,
  register_email_verify: false,
})

const loading = ref(false)
const sendingEmailCode = ref(false)
const emailCodeCountdown = ref(0)

const formState = reactive({
  username: '',
  email: '',
  email_code: '',
  password: '',
  confirm_password: '',
  captcha: '',
})

const validateConfirmPassword = async (_rule: Rule, value: string) => {
  if (value !== formState.password) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

// 动态验证规则：根据配置动态生成
const rules = computed<Record<string, Rule[]>>(() => {
  const baseRules: Record<string, Rule[]> = {
    username: [
      { required: true, message: '请输入用户名' },
      { min: 3, max: 20, message: '用户名长度为3-20个字符' },
    ],
    password: [
      { required: true, message: '请输入密码' },
      { min: 6, max: 20, message: '密码长度为6-20个字符' },
    ],
    confirm_password: [
      { required: true, message: '请确认密码' },
      { validator: validateConfirmPassword },
    ],
  }

  // 根据配置动态添加验证规则
  if (captchaConfig?.value?.register_email_verify) {
    baseRules.email = [
      { required: true, message: '请输入邮箱' },
      { type: 'email', message: '请输入有效的邮箱地址' },
    ]
    baseRules.email_code = [{ required: true, message: '请输入邮箱验证码' }]
  }

  if (captchaConfig?.value?.register_captcha_enabled) {
    baseRules.captcha = [{ required: true, message: '请输入验证码' }]
  }

  return baseRules
})

const handleSendEmailCode = async () => {
  try {
    await formRef.value?.validateFields(['email'])
  } catch {
    return
  }

  sendingEmailCode.value = true
  try {
    await sendEmailCode({ email: formState.email, type: 'register' })
    message.success('验证码已发送，请查收邮件')

    emailCodeCountdown.value = 60
    const timer = setInterval(() => {
      emailCodeCountdown.value--
      if (emailCodeCountdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error: any) {
    message.error(error.message || '发送验证码失败')
  } finally {
    sendingEmailCode.value = false
  }
}

const handleRegister = async () => {
  loading.value = true
  try {
    const registerData: RegisterParams = {
      username: formState.username,
      password: formState.password,
    }

    if (captchaConfig?.value?.register_email_verify) {
      registerData.email = formState.email
      registerData.email_code = formState.email_code
    }

    if (captchaConfig?.value?.register_captcha_enabled && captchaRef.value) {
      registerData.captcha_id = captchaRef.value.getCaptchaId()
      registerData.captcha_code = formState.captcha
    }

    await register(registerData)
    message.success('注册成功，请登录')
    router.push({
      path: '/login',
      query: { username: formState.username }
    })
  } catch (error: any) {
    captchaRef.value?.refresh()
  } finally {
    loading.value = false
  }
}
onMounted(async()=>{
   const captchaRes=await getCaptchaConfig()
   captchaConfig.value = {
      login_captcha_enabled: captchaRes.data?.login_captcha_enabled,
      register_captcha_enabled: captchaRes.data?.register_captcha_enabled,
      register_email_verify: captchaRes.data?.register_email_verify,
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

.email-code-row {
  display: flex;
  gap: 12px;

  .ant-input-affix-wrapper {
    flex: 1;
  }

  .ant-btn {
    width: 120px;
    flex-shrink: 0;
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
