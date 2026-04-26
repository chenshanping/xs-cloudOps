<template>
  <AuthLayout>
    <!-- 链接无效 -->
    <a-result
      v-if="invalidToken"
      status="error"
      title="链接无效或已过期"
      sub-title="请重新申请密码重置"
    >
      <template #extra>
        <a-button type="primary" @click="$router.push('/forgot-password')">
          重新申请
        </a-button>
        <a-button @click="$router.push('/login')">
          返回登录
        </a-button>
      </template>
    </a-result>

    <!-- 重置成功 -->
    <a-result
      v-else-if="resetSuccess"
      status="success"
      title="密码重置成功"
      sub-title="您现在可以使用新密码登录"
    >
      <template #extra>
        <a-button type="primary" @click="$router.push('/login')">
          立即登录
        </a-button>
      </template>
    </a-result>

    <!-- 重置表单 -->
    <template v-else>
      <div class="form-header">
        <h2>重置密码</h2>
        <p>请输入您的新密码</p>
      </div>

      <a-form
        :model="formState"
        :rules="rules"
        @finish="handleSubmit"
        layout="vertical"
      >
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
import { reactive, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { LockOutlined, ArrowLeftOutlined } from '@ant-design/icons-vue'
import type { Rule } from 'ant-design-vue/es/form'
import { resetPassword } from '@/api/auth'
import AuthLayout from '@/layouts/AuthLayout.vue'

const route = useRoute()
const loading = ref(false)
const invalidToken = ref(false)
const resetSuccess = ref(false)

const formState = reactive({
  password: '',
  confirm_password: '',
})

const validateConfirmPassword = async (_rule: Rule, value: string) => {
  if (value !== formState.password) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

const rules: Record<string, Rule[]> = {
  password: [
    { required: true, message: '请输入新密码' },
    { min: 6, max: 20, message: '密码长度为6-20个字符' },
  ],
  confirm_password: [
    { required: true, message: '请确认新密码' },
    { validator: validateConfirmPassword },
  ],
}

const handleSubmit = async () => {
  const token = route.query.token as string
  if (!token) {
    invalidToken.value = true
    return
  }

  loading.value = true
  try {
    await resetPassword({
      token,
      new_password: formState.password,
    })
    resetSuccess.value = true
  } catch (error: any) {
    if (error.message?.includes('无效') || error.message?.includes('过期')) {
      invalidToken.value = true
    } else {
      message.error(error.message || '重置密码失败')
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  const token = route.query.token as string
  if (!token) {
    invalidToken.value = true
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
</style>
