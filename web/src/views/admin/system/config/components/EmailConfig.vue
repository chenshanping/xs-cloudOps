<template>
  <div class="email-config">
    <a-form :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
      <!-- 登录安全配置 -->
      <a-divider orientation="left">登录安全</a-divider>
      
      <a-form-item label="登录验证码">
        <a-switch 
          v-model:checked="loginCaptchaEnabled" 
          checked-children="开" 
          un-checked-children="关"
        />
        <span class="switch-tip">开启后登录需要输入验证码</span>
      </a-form-item>

      <a-form-item label="验证码类型">
        <a-select 
          v-model:value="formData.login_captcha_type" 
          :disabled="!loginCaptchaEnabled"
          style="width: 200px"
        >
          <a-select-option value="digit">数字验证码</a-select-option>
          <a-select-option value="math">算术验证码</a-select-option>
          <a-select-option value="string">字符串验证码</a-select-option>
          <a-select-option value="slider">滑动验证码</a-select-option>
        </a-select>
        <div class="form-tip">仅开启验证码时可选</div>
      </a-form-item>

      <a-form-item label="滑动验证码背景" v-if="formData.login_captcha_type === 'slider'">
        <ImageUpload 
          v-model="formData.slider_captcha_bg" 
          v-model:fileId="sliderCaptchaBgFileId"
          :width="280" 
          :height="160" 
          :max-size="2*1024*1024"
          placeholder="上传背景图"
        />
        <div class="form-tip">建议尺寸 280x160 像素，留空使用默认渐变背景</div>
      </a-form-item>

      <a-form-item label="最大重试次数">
        <a-input-number 
          v-model:value="formData.login_max_retry" 
          :min="1" 
          :max="20" 
          style="width: 200px"
        />
        <div class="form-tip">达到次数后账户将被临时锁定，0表示不限制</div>
      </a-form-item>

      <a-form-item label="锁定时间(分钟)">
        <a-input-number 
          v-model:value="formData.login_lock_time" 
          :min="1" 
          :max="1440" 
          style="width: 200px"
        />
        <div class="form-tip">账户锁定后多久自动解锁</div>
      </a-form-item>

      <!-- 注册安全配置 -->
      <a-divider orientation="left">注册安全</a-divider>
      
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
import { computed, createVNode, reactive, ref, watch } from 'vue'
import { message, Modal, Input } from 'ant-design-vue'
import { MailOutlined, SettingOutlined } from '@ant-design/icons-vue'
import { useConfigStore } from '@/store/config'
import { sendTestEmail } from '@/api/config'
import ImageUpload from '@/components/ImageUpload.vue'
import { cloneFromSnapshot, createSnapshot, isSnapshotDirty } from '../config-tab-guard'

const configStore = useConfigStore()
const emit = defineEmits<{
  (e: 'dirty-change', value: boolean): void
}>()
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
  login_captcha_type: configStore.get('login_captcha_type') || 'digit',
  login_max_retry: parseInt(configStore.get('login_max_retry')) || 5,
  login_lock_time: parseInt(configStore.get('login_lock_time')) || 15,
  register_email_verify: configStore.get('register_email_verify') || '0',
  frontend_url: configStore.get('frontend_url') || 'http://localhost:5173',
  slider_captcha_bg: configStore.get('slider_captcha_bg') || '',
  slider_captcha_bg_file_id: configStore.get('slider_captcha_bg_file_id') || '0'
})

const sliderCaptchaBgFileId = computed({
  get: () => Number(formData.slider_captcha_bg_file_id) || 0,
  set: (val) => { formData.slider_captcha_bg_file_id = String(val || 0) }
})

// 开关状态转换
const loginCaptchaEnabled = computed({
  get: () => formData.login_captcha_enabled === '1',
  set: (val) => { formData.login_captcha_enabled = val ? '1' : '0' }
})

const registerEmailVerify = computed({
  get: () => formData.register_email_verify === '1',
  set: (val) => { formData.register_email_verify = val ? '1' : '0' }
})

const getConfigState = () => ({
  email_smtp_host: formData.email_smtp_host,
  email_smtp_port: formData.email_smtp_port,
  email_username: formData.email_username,
  email_password: formData.email_password,
  email_from_name: formData.email_from_name,
  login_captcha_enabled: formData.login_captcha_enabled,
  login_captcha_type: formData.login_captcha_type,
  login_max_retry: formData.login_max_retry,
  login_lock_time: formData.login_lock_time,
  register_email_verify: formData.register_email_verify,
  frontend_url: formData.frontend_url,
  slider_captcha_bg: formData.slider_captcha_bg,
  slider_captcha_bg_file_id: formData.slider_captcha_bg_file_id,
})

const applyConfigState = (state: ReturnType<typeof getConfigState>) => {
  formData.email_smtp_host = state.email_smtp_host
  formData.email_smtp_port = Number(state.email_smtp_port) || 587
  formData.email_username = state.email_username
  formData.email_password = state.email_password
  formData.email_from_name = state.email_from_name
  formData.login_captcha_enabled = state.login_captcha_enabled
  formData.login_captcha_type = state.login_captcha_type
  formData.login_max_retry = Number(state.login_max_retry) || 5
  formData.login_lock_time = Number(state.login_lock_time) || 15
  formData.register_email_verify = state.register_email_verify
  formData.frontend_url = state.frontend_url
  formData.slider_captcha_bg = state.slider_captcha_bg
  formData.slider_captcha_bg_file_id = state.slider_captcha_bg_file_id
}

const baselineSnapshot = ref(createSnapshot(getConfigState()))
const hasUnsavedChanges = computed(() => isSnapshotDirty(baselineSnapshot.value, getConfigState()))

watch(hasUnsavedChanges, (value) => {
  emit('dirty-change', value)
}, { immediate: true })

// 保存配置
const save = async () => {
  saving.value = true
  try {
    const configs: Record<string, string> = {}
    for (const [key, val] of Object.entries(getConfigState())) {
      configs[key] = typeof val === 'number' ? String(val) : val
    }
    await configStore.updateConfigs(configs)
    baselineSnapshot.value = createSnapshot(getConfigState())
    message.success('配置保存成功')
    return true
  } catch (e) {
    message.error('保存失败')
    return false
  } finally {
    saving.value = false
  }
}

const discardChanges = () => {
  const restored = cloneFromSnapshot<ReturnType<typeof getConfigState>>(baselineSnapshot.value)
  applyConfigState(restored)
  smtpDrawerVisible.value = false
}

const closeTransientUi = () => {
  smtpDrawerVisible.value = false
}

const handleSave = async () => {
  await save()
}

// 恢复默认
const handleReset = () => {
  formData.email_smtp_host = ''
  formData.email_smtp_port = 587
  formData.email_username = ''
  formData.email_password = ''
  formData.email_from_name = '系统邮件'
  formData.login_captcha_enabled = '0'
  formData.login_captcha_type = 'digit'
  formData.login_max_retry = 5
  formData.login_lock_time = 15
  formData.register_email_verify = '0'
  formData.frontend_url = 'http://localhost:5173'
  formData.slider_captcha_bg = ''
  formData.slider_captcha_bg_file_id = '0'
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

defineExpose({
  isDirty: () => hasUnsavedChanges.value,
  save,
  discardChanges,
  closeTransientUi,
})

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
