<template>
  <div class="system-config">
    <div class="config-form">
      <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
        <a-divider orientation="left">基本设置</a-divider>

        <a-form-item label="系统名称">
          <a-input v-model:value="formData.sys_name" placeholder="请输入系统名称" />
        </a-form-item>

        <a-form-item label="系统Logo">
          <ImageUpload
            v-model="formData.sys_logo"
            :width="120"
            :height="60"
            :max-size="5 * 1024 * 1024"
            placeholder="上传Logo"
          />
        </a-form-item>

        <a-form-item label="前台模式">
          <a-radio-group v-model:value="formData.front_mode">
            <a-radio value="full">完整前台</a-radio>
            <a-radio value="profile">仅个人中心</a-radio>
          </a-radio-group>
          <div class="form-tip">
            完整前台: 显示全部前台页面；仅个人中心: 只显示个人资料页面
          </div>
        </a-form-item>

        <a-form-item label="用户身份按钮">
          <a-switch
            :checked="formData.user_profile_button_visible === 'true'"
            @change="(checked: boolean) => formData.user_profile_button_visible = checked ? 'true' : 'false'"
          />
          <div class="form-tip">
            控制后台用户管理列表中的“身份”按钮是否显示，默认隐藏
          </div>
        </a-form-item>

        <a-form-item :wrapper-col="{ offset: 8, span: 16 }" style="margin-top: 24px">
          <a-space>
            <a-button type="primary" :loading="saving" @click="handleSave">保存配置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { useConfigStore } from '@/store/config'
import ImageUpload from '@/components/ImageUpload.vue'

const configStore = useConfigStore()
const route = useRoute()
const saving = ref(false)

const SYSTEM_CONFIG_KEYS = [
  'sys_name',
  'sys_logo',
  'front_mode',
  'user_profile_button_visible',
] as const

const formData = reactive({
  sys_name: configStore.get('sys_name'),
  sys_logo: configStore.get('sys_logo'),
  front_mode: configStore.get('front_mode') || 'full',
  user_profile_button_visible: configStore.get('user_profile_button_visible') || 'false',
})

const updateTitle = () => {
  const sysName = configStore.get('sys_name') || 'Go RBAC Admin'
  const pageTitle = route.meta?.title as string
  document.title = pageTitle ? `${pageTitle} - ${sysName}` : sysName
}

const handleSave = async () => {
  saving.value = true
  try {
    const configs: Record<string, string> = {}
    for (const key of SYSTEM_CONFIG_KEYS) {
      configs[key] = formData[key]
    }
    await configStore.updateConfigs(configs)
    updateTitle()
    message.success('配置保存成功')
  } catch {
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.config-form {
  width: 100%;
  max-width: 560px;
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: #999;
}
</style>
