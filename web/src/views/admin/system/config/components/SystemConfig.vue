<template>
  <div class="system-config">
    <div class="config-form">
      <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="系统名称">
          <a-input v-model:value="basicForm.sys_name" placeholder="请输入系统名称" />
        </a-form-item>

        <a-form-item label="系统Logo">
          <ImageUpload
            v-model="basicForm.sys_logo"
            :width="120"
            :height="60"
            :max-size="5 * 1024 * 1024"
            placeholder="上传Logo"
          />
        </a-form-item>

        <a-form-item label="前台模式">
          <a-radio-group v-model:value="basicForm.front_mode">
            <a-radio value="full">完整前台</a-radio>
            <a-radio value="profile">仅个人中心</a-radio>
          </a-radio-group>
          <div class="form-tip">
            完整前台: 显示全部前台页面；仅个人中心: 只显示个人资料页面
          </div>
        </a-form-item>

        <a-form-item label="用户身份按钮">
          <a-switch
            :checked="basicForm.user_profile_button_visible === 'true'"
            @change="(checked: boolean) => basicForm.user_profile_button_visible = checked ? 'true' : 'false'"
          />
          <div class="form-tip">
            控制后台用户管理列表中的“身份”按钮是否显示，默认隐藏
          </div>
        </a-form-item>

        <a-form-item label="用户默认密码">
          <a-space-compact style="width: 100%">
            <a-input-password
              v-model:value="basicForm.user_default_password"
              placeholder="请输入用户默认密码"
            />
            <a-button @click="handleGeneratePassword">随机生成</a-button>
          </a-space-compact>
          <div class="form-tip">
            后台用户管理单条/批量重置密码默认使用该值，建议至少 6 位
          </div>
        </a-form-item>

        <a-form-item :wrapper-col="{ offset: 8, span: 16 }" style="margin-top: 24px">
          <a-button type="primary" :loading="basicSaving" @click="handleBasicSave">保存基础配置</a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useConfigStore } from '@/store/config'
import ImageUpload from '@/components/ImageUpload.vue'
import { cloneFromSnapshot, createSnapshot, isSnapshotDirty } from '../config-tab-guard'

const configStore = useConfigStore()
const emit = defineEmits<{
  (e: 'dirty-change', value: boolean): void
}>()

const basicSaving = ref(false)

const basicForm = reactive({
  sys_name: configStore.get('sys_name'),
  sys_logo: configStore.get('sys_logo'),
  front_mode: configStore.get('front_mode') || 'full',
  user_profile_button_visible: configStore.get('user_profile_button_visible') || 'false',
  user_default_password: configStore.get('user_default_password') || '123456',
})

const getBasicState = () => ({
  sys_name: basicForm.sys_name,
  sys_logo: basicForm.sys_logo,
  front_mode: basicForm.front_mode,
  user_profile_button_visible: basicForm.user_profile_button_visible,
  user_default_password: basicForm.user_default_password,
})

const applyBasicState = (state: ReturnType<typeof getBasicState>) => {
  basicForm.sys_name = state.sys_name
  basicForm.sys_logo = state.sys_logo
  basicForm.front_mode = state.front_mode
  basicForm.user_profile_button_visible = state.user_profile_button_visible
  basicForm.user_default_password = state.user_default_password
}

const baselineSnapshot = ref(createSnapshot(getBasicState()))
const hasUnsavedChanges = computed(() => isSnapshotDirty(baselineSnapshot.value, getBasicState()))

watch(hasUnsavedChanges, (value) => {
  emit('dirty-change', value)
}, { immediate: true })

const updateTitle = () => {
  document.title = configStore.get('sys_name') || 'Go RBAC Admin'
}

const save = async () => {
  basicForm.user_default_password = basicForm.user_default_password.trim()
  if (basicForm.user_default_password.length < 6) {
    message.warning('用户默认密码至少 6 位')
    return false
  }

  basicSaving.value = true
  try {
    const configs: Record<string, string> = {
      sys_name: basicForm.sys_name,
      sys_logo: basicForm.sys_logo,
      front_mode: basicForm.front_mode,
      user_profile_button_visible: basicForm.user_profile_button_visible,
      user_default_password: basicForm.user_default_password,
    }
    await configStore.updateConfigs(configs)
    baselineSnapshot.value = createSnapshot(getBasicState())
    updateTitle()
    message.success('基础配置保存成功')
    return true
  } catch {
    message.error('保存失败')
    return false
  } finally {
    basicSaving.value = false
  }
}

const discardChanges = () => {
  const restored = cloneFromSnapshot<ReturnType<typeof getBasicState>>(baselineSnapshot.value)
  applyBasicState(restored)
}

const closeTransientUi = () => {}

const handleBasicSave = async () => {
  await save()
}

const handleGeneratePassword = () => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789'
  let password = ''
  for (let i = 0; i < 12; i += 1) {
    const index = Math.floor(Math.random() * chars.length)
    password += chars[index]
  }
  basicForm.user_default_password = password
}

defineExpose({
  isDirty: () => hasUnsavedChanges.value,
  save,
  discardChanges,
  closeTransientUi,
})
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
