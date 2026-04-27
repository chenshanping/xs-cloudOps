<template>
  <div class="system-config">
    <a-tabs v-model:activeKey="activeTab">
      <a-tab-pane key="basic" tab="基础设置">
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

            <a-form-item :wrapper-col="{ offset: 8, span: 16 }" style="margin-top: 24px">
              <a-button type="primary" :loading="basicSaving" @click="handleBasicSave">保存基础配置</a-button>
            </a-form-item>
          </a-form>
        </div>
      </a-tab-pane>

      <a-tab-pane key="storage" tab="存储设置">
        <div class="config-form">
          <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
            <a-form-item label="存储类型">
              <a-select v-model:value="storageType" @change="handleStorageTypeChange">
                <a-select-option v-for="item in storageTypeOptions" :key="item.value" :value="item.value">
                  {{ item.label }}
                </a-select-option>
              </a-select>
            </a-form-item>

            <template v-if="storageType === 'local'">
              <a-form-item label="存储路径">
                <a-input v-model:value="storageConfig.base_path" placeholder="如：uploads" />
              </a-form-item>
              <a-form-item label="访问地址">
                <a-input :value="localPreviewUrl" disabled />
                <div class="form-tip">本地存储访问前缀固定由后端提供，不再单独配置。</div>
              </a-form-item>
            </template>

            <template v-else-if="storageType === 'aliyun'">
              <a-form-item label="Endpoint">
                <a-input v-model:value="storageConfig.endpoint" placeholder="如：oss-cn-hangzhou.aliyuncs.com" />
              </a-form-item>
              <a-form-item label="AccessKey ID">
                <a-input v-model:value="storageConfig.access_key_id" />
              </a-form-item>
              <a-form-item label="AccessKey Secret">
                <a-input-password v-model:value="storageConfig.access_key_secret" />
              </a-form-item>
              <a-form-item label="Bucket">
                <a-input v-model:value="storageConfig.bucket_name" />
              </a-form-item>
              <a-form-item label="Region">
                <a-input v-model:value="storageConfig.region" placeholder="如：cn-hangzhou" />
              </a-form-item>
            </template>

            <template v-else-if="storageType === 'tencent'">
              <a-form-item label="Region">
                <a-input v-model:value="storageConfig.region" placeholder="如：ap-guangzhou" />
              </a-form-item>
              <a-form-item label="SecretId">
                <a-input v-model:value="storageConfig.secret_id" />
              </a-form-item>
              <a-form-item label="SecretKey">
                <a-input-password v-model:value="storageConfig.secret_key" />
              </a-form-item>
              <a-form-item label="Bucket">
                <a-input v-model:value="storageConfig.bucket" placeholder="如：mybucket-1234567890" />
              </a-form-item>
              <a-form-item label="AppID">
                <a-input v-model:value="storageConfig.app_id" />
              </a-form-item>
            </template>

            <template v-else-if="storageType === 'minio'">
              <a-form-item label="Endpoint">
                <a-input v-model:value="storageConfig.endpoint" placeholder="如：127.0.0.1:9000" />
              </a-form-item>
              <a-form-item label="AccessKey ID">
                <a-input v-model:value="storageConfig.access_key_id" />
              </a-form-item>
              <a-form-item label="SecretAccessKey">
                <a-input-password v-model:value="storageConfig.secret_access_key" />
              </a-form-item>
              <a-form-item label="Bucket">
                <a-input v-model:value="storageConfig.bucket_name" />
              </a-form-item>
              <a-form-item label="使用SSL">
                <a-switch v-model:checked="storageConfig.use_ssl" />
              </a-form-item>
            </template>

            <a-form-item :wrapper-col="{ offset: 8, span: 16 }" style="margin-top: 24px">
              <a-space>
                <a-button :loading="storageTesting" @click="handleTestStorage">测试连接</a-button>
                <a-button type="primary" :loading="storageSaving" @click="handleStorageSave">保存存储配置</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </div>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { useConfigStore } from '@/store/config'
import { testStorageConfig } from '@/api/config'
import ImageUpload from '@/components/ImageUpload.vue'
import type { StorageType } from '@/types/storage'
import { storageTypeOptions } from '@/types/storage'
import {
  buildStorageSavePayload,
  changeActiveStorageType,
  createStorageDraftsState,
  getActiveStorageConfig,
} from '../storage-config-state'

const configStore = useConfigStore()
const route = useRoute()

const basicSaving = ref(false)
const storageSaving = ref(false)
const storageTesting = ref(false)
const activeTab = ref(route.query.tab === 'storage' ? 'storage' : 'basic')

const BASIC_CONFIG_KEYS = [
  'sys_name',
  'sys_logo',
  'front_mode',
  'user_profile_button_visible',
] as const

const basicForm = reactive({
  sys_name: configStore.get('sys_name'),
  sys_logo: configStore.get('sys_logo'),
  front_mode: configStore.get('front_mode') || 'full',
  user_profile_button_visible: configStore.get('user_profile_button_visible') || 'false',
})

const storageState = reactive(createStorageDraftsState({
  storage_type: configStore.get('storage_type'),
  storage_local_config: configStore.get('storage_local_config'),
  storage_aliyun_config: configStore.get('storage_aliyun_config'),
  storage_tencent_config: configStore.get('storage_tencent_config'),
  storage_minio_config: configStore.get('storage_minio_config'),
}))
const storageType = computed<StorageType>({
  get: () => storageState.activeType,
  set: (value) => {
    storageState.activeType = changeActiveStorageType(storageState, value).activeType
  },
})
const storageConfig = computed(() => getActiveStorageConfig(storageState))

const localPreviewUrl = computed(() => '/api/v1/upload')

const updateTitle = () => {
  const sysName = configStore.get('sys_name') || 'Go RBAC Admin'
  const pageTitle = route.meta?.title as string
  document.title = pageTitle ? `${pageTitle} - ${sysName}` : sysName
}

const handleStorageTypeChange = (value: StorageType) => {
  storageType.value = value
}

const handleBasicSave = async () => {
  basicSaving.value = true
  try {
    const configs: Record<string, string> = {}
    for (const key of BASIC_CONFIG_KEYS) {
      configs[key] = basicForm[key]
    }
    await configStore.updateConfigs(configs)
    updateTitle()
    message.success('基础配置保存成功')
  } catch {
    message.error('保存失败')
  } finally {
    basicSaving.value = false
  }
}

const handleTestStorage = async () => {
  storageTesting.value = true
  try {
    await testStorageConfig({
      type: storageType.value,
      config: JSON.stringify(storageConfig.value),
    })
    message.success('连接测试成功')
  } catch {
    message.error('连接测试失败')
  } finally {
    storageTesting.value = false
  }
}

const handleStorageSave = async () => {
  storageSaving.value = true
  try {
    await configStore.updateConfigs(buildStorageSavePayload(storageState))
    message.success('存储配置保存成功')
  } catch {
    message.error('保存失败')
  } finally {
    storageSaving.value = false
  }
}

watch(
  () => route.query.tab,
  (tab) => {
    activeTab.value = tab === 'storage' ? 'storage' : 'basic'
  }
)
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
