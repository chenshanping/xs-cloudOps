<template>
  <AuthLayout>
    <!-- 重置成功 -->
    <template v-if="step === 3">
      <a-result
        status="success"
        title="密码重置成功"
        sub-title="您现在可以使用新密码登录了"
      >
        <template #extra>
          <a-button type="primary" @click="goToLogin">
            立即登录
          </a-button>
        </template>
      </a-result>
    </template>

    <!-- 表单 -->
    <template v-else>
      <div class="form-header">
        <h2>忘记密码</h2>
        <template v-if="emailVerifyEnabled">
          <p v-if="step === 1">请输入您的注册邮箱，我们将发送验证码</p>
          <p v-else>请输入邮箱验证码和新密码</p>
        </template>
        <template v-else>
          <p>请输入您的账号和新密码</p>
        </template>
      </div>

      <!-- 邮箱验证模式的步骤条 -->
      <a-steps v-if="emailVerifyEnabled" :current="step - 1" size="small" class="reset-steps">
        <a-step title="验证邮箱" />
        <a-step title="重置密码" />
        <a-step title="完成" />
      </a-steps>

      <a-form
        ref="formRef"
        :model="formState"
        :rules="rules"
        @finish="handleSubmit"
        layout="vertical"
      >
        <!-- ========== 邮箱验证模式 ========== -->
        <template v-if="emailVerifyEnabled">
          <!-- 步骤1：输入邮箱 -->
          <template v-if="step === 1">
            <a-form-item name="email">
              <a-input
                v-model:value="formState.email"
                placeholder="请输入邮箱"
                size="large"
              >
                <template #prefix><MailOutlined /></template>
              </a-input>
            </a-form-item>

            <a-form-item>
              <a-button
                type="primary"
                html-type="submit"
                size="large"
                block
                :loading="loading"
              >
                发送验证码
              </a-button>
            </a-form-item>
          </template>

          <!-- 步骤2：输入验证码和新密码 -->
          <template v-else-if="step === 2">
            <a-form-item>
              <a-alert
                :message="`验证码已发送至 ${formState.email}`"
                type="info"
                show-icon
                closable
              />
            </a-form-item>

            <a-form-item name="email_code">
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
                  :disabled="countdown > 0"
                  :loading="sendingCode"
                  @click="handleResendCode"
                >
                  {{ countdown > 0 ? `${countdown}s` : '重发' }}
                </a-button>
              </div>
            </a-form-item>

            <a-form-item name="password">
              <a-input-password
                v-model:value="formState.password"
                placeholder="请输入新密码"
                size="large"
              >
                <template #prefix><LockOutlined /></template>
              </a-input-password>
            </a-form-item>

            <a-form-item name="confirm_password">
              <a-input-password
                v-model:value="formState.confirm_password"
                placeholder="请确认新密码"
                size="large"
              >
                <template #prefix><LockOutlined /></template>
              </a-input-password>
            </a-form-item>

            <a-form-item>
              <a-space style="width: 100%" direction="vertical" :size="12">
                <a-button
                  type="primary"
                  html-type="submit"
                  size="large"
                  block
                  :loading="loading"
                >
                  重置密码
                </a-button>
                <a-button
                  size="large"
                  block
                  @click="step = 1"
                >
                  返回上一步
                </a-button>
              </a-space>
            </a-form-item>
          </template>
        </template>

        <!-- ========== 账号验证模式（无邮箱） ========== -->
        <template v-else>
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
              placeholder="请输入新密码"
              size="large"
            >
              <template #prefix><LockOutlined /></template>
            </a-input-password>
          </a-form-item>

          <a-form-item name="confirm_password">
            <a-input-password
              v-model:value="formState.confirm_password"
              placeholder="请确认新密码"
              size="large"
            >
              <template #prefix><LockOutlined /></template>
            </a-input-password>
          </a-form-item>

          <!-- 图形验证码 -->
          <a-form-item name="captcha">
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
              重置密码
            </a-button>
          </a-form-item>
        </template>
      </a-form>

      <div class="form-footer">
        <router-link to="/login">
          <ArrowLeftOutlined /> 返回登录
        </router-link>
      </div>
    </template>
  </AuthLayout>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  MailOutlined,
  ArrowLeftOutlined,
  SafetyOutlined,
  LockOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import type { Rule, FormInstance } from 'ant-design-vue/es/form'
import { sendEmailCode, resetPasswordByEmail, resetPasswordByUsername } from '@/api/auth'
import { getCaptchaConfig } from '@/api/captcha'
import AuthLayout from '@/layouts/AuthLayout.vue'
import Captcha from '@/components/Captcha.vue'

interface CaptchaConfig {
  login_captcha_enabled: boolean
  register_captcha_enabled: boolean
  register_email_verify: boolean
}

const router = useRouter()
const formRef = ref<FormInstance>()
const captchaRef = ref<InstanceType<typeof Captcha>>()
const loading = ref(false)
const sendingCode = ref(false)
const countdown = ref(0)
const step = ref(1) // 1: 输入邮箱/账号, 2: 输入验证码和密码, 3: 完成
const emailVerifyEnabled = ref(true) // 是否启用邮箱验证

const formState = reactive({
  email: '',
  email_code: '',
  username: '',
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

const rules = computed<Record<string, Rule[]>>(() => {
  // 邮箱验证模式
  if (emailVerifyEnabled.value) {
    if (step.value === 1) {
      return {
        email: [
          { required: true, message: '请输入邮箱' },
          { type: 'email', message: '请输入有效的邮箱地址' },
        ],
      }
    } else {
      return {
        email_code: [{ required: true, message: '请输入邮箱验证码' }],
        password: [
          { required: true, message: '请输入新密码' },
          { min: 6, max: 20, message: '密码长度为6-20个字符' },
        ],
        confirm_password: [
          { required: true, message: '请确认新密码' },
          { validator: validateConfirmPassword },
        ],
      }
    }
  } else {
    // 账号验证模式
    return {
      username: [
        { required: true, message: '请输入用户名' },
      ],
      password: [
        { required: true, message: '请输入新密码' },
        { min: 6, max: 20, message: '密码长度为6-20个字符' },
      ],
      confirm_password: [
        { required: true, message: '请确认新密码' },
        { validator: validateConfirmPassword },
      ],
      captcha: [{ required: true, message: '请输入验证码' }],
    }
  }
})

const startCountdown = () => {
  countdown.value = 60
  const timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
    }
  }, 1000)
}

const handleSendCode = async () => {
  try {
    await sendEmailCode({ email: formState.email, type: 'reset_password' })
    message.success('验证码已发送，请查收邮件')
    startCountdown()
  } catch (error: any) {
    message.error(error.message || '发送验证码失败')
    throw error
  }
}

const handleResendCode = async () => {
  sendingCode.value = true
  try {
    await handleSendCode()
  } finally {
    sendingCode.value = false
  }
}

// 跳转到登录页，并传递用户名
const goToLogin = () => {
  const query: any = {}
  if (formState.username) {
    query.username = formState.username
  }
  router.push({ path: '/login', query })
}

const handleSubmit = async () => {
  // 邮箱验证模式
  if (emailVerifyEnabled.value) {
    if (step.value === 1) {
      // 步骤1：发送验证码
      loading.value = true
      try {
        await handleSendCode()
        step.value = 2
      } catch {
        // 错误已在 handleSendCode 中处理
      } finally {
        loading.value = false
      }
    } else if (step.value === 2) {
      // 步骤2：重置密码
      loading.value = true
      try {
        await resetPasswordByEmail({
          email: formState.email,
          email_code: formState.email_code,
          new_password: formState.password,
        })
        message.success('密码重置成功')
        step.value = 3
      } catch (error: any) {
      } finally {
        loading.value = false
      }
    }
  } else {
    // 账号验证模式
    loading.value = true
    try {
      await resetPasswordByUsername({
        username: formState.username,
        new_password: formState.password,
        captcha_id: captchaRef.value?.getCaptchaId() || '',
        captcha: formState.captcha,
      })
      message.success('密码重置成功')
      step.value = 3
    } catch (error: any) {
      // message.error(error.message || '重置密码失败')
      captchaRef.value?.refresh()
    } finally {
      loading.value = false
    }
  }
}

onMounted(async () => {
  try {
    const res = await getCaptchaConfig()
    emailVerifyEnabled.value = res.data?.register_email_verify ?? true
  } catch (error) {
    console.error('获取配置失败', error)
  }
})
</script>

<style scoped lang="scss">
.form-header {
  text-align: center;
  margin-bottom: 24px;

  h2 {
    font-size: 26px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0 0 8px 0;
  }

  p {
    color: #666;
    margin: 0;
    line-height: 1.5;
    font-size: 14px;
  }
}

.reset-steps {
  margin-bottom: 32px;
}

.email-code-row {
  display: flex;
  gap: 12px;

  .ant-input-affix-wrapper {
    flex: 1;
  }

  .ant-btn {
    width: 80px;
    flex-shrink: 0;
  }
}

.form-footer {
  text-align: center;
  margin-top: 24px;

  a {
    color: #667eea;
    display: inline-flex;
    align-items: center;
    gap: 8px;

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

:deep(.ant-result) {
  padding: 0;

  .ant-result-title {
    color: #1a1a2e;
  }

  .ant-result-subtitle {
    color: #666;
  }
}

:deep(.ant-alert) {
  border-radius: 8px;
}
</style>
