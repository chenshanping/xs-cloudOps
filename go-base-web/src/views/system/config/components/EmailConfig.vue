<template>
  <div class="email-config">
    <a-form :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
      <!-- 安全配置 -->
      <a-divider orientation="left">安全配置</a-divider>
      
      <a-form-item label="登录验证码">
        <a-switch 
          v-model:checked="loginCaptchaEnabled" 
          checked-children="开" 
          un-checked-children="关"
        />
        <span class="switch-tip">开启后登录需要输入图形验证码</span>
      </a-form-item>
      
      <a-form-item label="注册验证码">
        <a-switch 
          v-model:checked="registerCaptchaEnabled" 
          checked-children="开" 
          un-checked-children="关"
        />
        <span class="switch-tip">开启后注册需要输入图形验证码</span>
      </a-form-item>
      
      <a-form-item label="邮箱验证">
        <a-switch 
          v-model:checked="registerEmailVerify" 
          checked-children="开" 
          un-checked-children="关"
        />
        <span class="switch-tip">开启后注册和忘记密码需要验证邮箱</span>
      </a-form-item>

      <!-- 保存按钮 -->
      <a-form-item :wrapper-col="{ offset: 6, span: 16 }" style="margin-top: 24px">
        <a-space>
          <a-button type="primary" :loading="saving" @click="handleSave">保存配置</a-button>
          <a-button @click="smtpDrawerVisible = true">
            <SettingOutlined /> SMTP配置
          </a-button>
        </a-space>
      </a-form-item>
    </a-form>

    <!-- SMTP配置抽屉 -->
    <a-drawer
      v-model:open="smtpDrawerVisible"
      title="SMTP 邮箱配置"
      placement="right"
      :width="450"
    >
      <a-form :label-col="{ span: 6 }" :wrapper-col="{ span: 18 }">
        <a-form-item label="SMTP服务器">
          <a-input v-model:value="formData.email_smtp_host" placeholder="如: smtp.qq.com" />
        </a-form-item>
        
        <a-form-item label="SMTP端口">
          <a-input-number 
            v-model:value="formData.email_smtp_port" 
            :min="1" 
            :max="65535" 
            style="width: 100%"
            placeholder="如: 587 或 465"
          />
        </a-form-item>
        
        <a-form-item label="发件人邮箱">
          <a-input v-model:value="formData.email_username" placeholder="your-email@example.com" />
        </a-form-item>
        
        <a-form-item label="授权码">
          <a-input-password v-model:value="formData.email_password" placeholder="SMTP授权码" />
          <div class="form-tip">QQ邮箱等需要使用授权码而非登录密码</div>
        </a-form-item>
        
        <a-form-item label="发件人名称">
          <a-input v-model:value="formData.email_from_name" placeholder="系统邮件" />
        </a-form-item>

        <a-form-item :wrapper-col="{ offset: 6 }">
          <a-button type="primary" ghost @click="testEmail" :loading="testing">
            <MailOutlined /> 发送测试邮件
          </a-button>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, createVNode } from 'vue'
import { message, Modal, Input } from 'ant-design-vue'
import { MailOutlined, SettingOutlined } from '@ant-design/icons-vue'
import { useConfigStore } from '@/store/config'
import { sendTestEmail } from '@/api/config'

const configStore = useConfigStore()
const saving = ref(false)
const testing = ref(false)
const smtpDrawerVisible = ref(false)

// 直接从 store 初始化数据
const formData = reactive({
  email_smtp_host: configStore.get('email_smtp_host'),
  email_smtp_port: parseInt(configStore.get('email_smtp_port')) || 587,
  email_username: configStore.get('email_username'),
  email_password: configStore.get('email_password'),
  email_from_name: configStore.get('email_from_name'),
  login_captcha_enabled: configStore.get('login_captcha_enabled') || '0',
  register_captcha_enabled: configStore.get('register_captcha_enabled') || '0',
  register_email_verify: configStore.get('register_email_verify') || '0',
  frontend_url: configStore.get('frontend_url') || 'http://localhost:5173'
})

// 开关状态转换
const loginCaptchaEnabled = computed({
  get: () => formData.login_captcha_enabled === '1',
  set: (val) => { formData.login_captcha_enabled = val ? '1' : '0' }
})

const registerCaptchaEnabled = computed({
  get: () => formData.register_captcha_enabled === '1',
  set: (val) => { formData.register_captcha_enabled = val ? '1' : '0' }
})

const registerEmailVerify = computed({
  get: () => formData.register_email_verify === '1',
  set: (val) => { formData.register_email_verify = val ? '1' : '0' }
})

// 保存配置
const handleSave = async () => {
  saving.value = true
  try {
    const configs: Record<string, string> = {}
    for (const key in formData) {
      const val = (formData as any)[key]
      configs[key] = typeof val === 'number' ? String(val) : val
    }
    await configStore.updateConfigs(configs)
    message.success('配置保存成功')
  } catch (e) {
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 恢复默认
const handleReset = () => {
  formData.email_smtp_host = ''
  formData.email_smtp_port = 587
  formData.email_username = ''
  formData.email_password = ''
  formData.email_from_name = '系统邮件'
  formData.login_captcha_enabled = '0'
  formData.register_captcha_enabled = '0'
  formData.register_email_verify = '0'
  formData.frontend_url = 'http://localhost:5173'
}

// 测试邮件
const testEmailAddress = ref('')

const testEmail = () => {
  if (!formData.email_smtp_host || !formData.email_username || !formData.email_password) {
    message.warning('请先填写完整的SMTP配置并保存')
    return
  }
  
  // 默认使用发件人邮箱
  testEmailAddress.value = formData.email_username
  
  Modal.confirm({
    title: '发送测试邮件',
    content: createVNode('div', {}, [
      createVNode('p', { style: 'margin-bottom: 12px; color: #666;' }, '请输入接收测试邮件的邮箱地址：'),
      createVNode(Input, {
        value: testEmailAddress.value,
        'onUpdate:value': (val: string) => { testEmailAddress.value = val },
        placeholder: '请输入邮箱地址',
        style: 'width: 100%'
      })
    ]),
    okText: '发送',
    cancelText: '取消',
    onOk: async () => {
      if (!testEmailAddress.value) {
        message.warning('请输入邮箱地址')
        return Promise.reject()
      }
      testing.value = true
      try {
        await sendTestEmail(testEmailAddress.value)
        message.success('测试邮件发送成功，请检查收件箱')
      } catch (e: any) {
        message.error('发送失败，请检查SMTP配置')
        return Promise.reject()
      } finally {
        testing.value = false
      }
    }
  })
}

</script>

<style scoped>
.email-config {
  max-width: 600px;
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: #999;
}

.switch-tip {
  margin-left: 12px;
  font-size: 12px;
  color: #999;
}
</style>
