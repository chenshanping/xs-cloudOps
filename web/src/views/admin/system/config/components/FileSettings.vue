<template>
  <div class="file-settings">
    <div class="config-form">
      <a-alert
        class="storage-notice"
        type="info"
        show-icon
        message="这里配置的是默认上传存储位置，仅影响后续新上传文件。"
        description="切换后不会自动迁移已有历史文件；如需处理存量文件，请前往 文件管理 > 文件迁移。"
      />

      <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="文件删除方式">
          <a-radio-group v-model:value="fileDeleteMode">
            <a-radio value="logical">逻辑删除</a-radio>
            <a-radio value="physical">物理删除</a-radio>
          </a-radio-group>
          <div class="form-tip">
            逻辑删除仅隐藏文件记录；物理删除会同时删除存储中的真实文件、分片临时文件和相关数据库记录。
          </div>
        </a-form-item>

        <a-form-item label="默认上传存储" :extra="storageTypeExtra">
          <a-space direction="vertical" size="small" style="width: 100%">
            <a-select v-model:value="storageType" @change="handleStorageTypeChange">
              <a-select-option v-for="item in storageTypeOptions" :key="item.value" :value="item.value">
                {{ item.label }}
              </a-select-option>
            </a-select>
            <a-space>
              <a-tag color="blue">{{ currentStorageLabel }}</a-tag>
              <a-tag color="processing">当前默认上传位置</a-tag>
              <a-button @click="configVisible = true">查看当前类型配置</a-button>
            </a-space>
          </a-space>
        </a-form-item>

        <a-form-item :wrapper-col="{ offset: 8, span: 16 }" style="margin-top: 24px">
          <a-button type="primary" :loading="saving" @click="handleSave">保存文件设置</a-button>
        </a-form-item>
      </a-form>
    </div>

    <a-modal
      v-model:open="configVisible"
      :title="`${currentStorageLabel}配置`"
      width="640px"
      @ok="handleStorageSave"
    >
      <a-form :label-col="{ span: 7 }" :wrapper-col="{ span: 16 }">
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
      </a-form>

      <template #footer>
        <a-space>
          <a-button @click="configVisible = false">关闭</a-button>
          <a-button
            v-if="storageType !== 'local'"
            :loading="storageTesting"
            :class="{ 'validated-btn': storageValidated }"
            @click="handleTestStorage"
          >
            <template v-if="storageValidated">
              <CheckCircleOutlined /> 已验证
            </template>
            <template v-else>测试连接</template>
          </a-button>
          <a-button type="primary" :loading="storageSaving" @click="handleStorageSave">保存配置</a-button>
        </a-space>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { CheckCircleOutlined } from '@ant-design/icons-vue'
import { useConfigStore } from '@/store/config'
import { testStorageConfig } from '@/api/config'
import type { StorageType } from '@/types/storage'
import { storageTypeOptions } from '@/types/storage'
import {
  changeActiveStorageType,
  createStorageDraftsState,
  getActiveStorageConfig,
} from '../storage-config-state'
import { cloneFromSnapshot, createSnapshot, isSnapshotDirty } from '../config-tab-guard'

const configStore = useConfigStore()
const emit = defineEmits<{
  (e: 'dirty-change', value: boolean): void
}>()

const fileDeleteMode = ref(configStore.get('file_delete_mode') || 'logical')
const configVisible = ref(false)
const saving = ref(false)
const storageSaving = ref(false)
const storageTesting = ref(false)
const storageValidated = ref(false)

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
const currentStorageLabel = computed(() => {
  return storageTypeOptions.find((item) => item.value === storageType.value)?.label || storageType.value
})
const localPreviewUrl = computed(() => '/api/v1/upload')
const storageTypeExtra = '这里仅影响后续新上传文件的默认存储位置；切换后，已上传文件仍保留在原存储位置。'

const handleStorageTypeChange = (value: StorageType) => {
  storageType.value = value
  storageValidated.value = false
}

watch(storageConfig, () => {
  storageValidated.value = false
}, { deep: true })

const getFileSettingsState = () => ({
  file_delete_mode: fileDeleteMode.value,
  storage_type: storageState.activeType,
  drafts: {
    local: { ...storageState.drafts.local },
    aliyun: { ...storageState.drafts.aliyun },
    tencent: { ...storageState.drafts.tencent },
    minio: { ...storageState.drafts.minio },
  },
})

const applyFileSettingsState = (state: ReturnType<typeof getFileSettingsState>) => {
  fileDeleteMode.value = state.file_delete_mode
  storageState.activeType = state.storage_type
  storageState.drafts.local = { ...state.drafts.local }
  storageState.drafts.aliyun = { ...state.drafts.aliyun }
  storageState.drafts.tencent = { ...state.drafts.tencent }
  storageState.drafts.minio = { ...state.drafts.minio }
}

const buildFullSavePayload = (): Record<string, string> => ({
  file_delete_mode: fileDeleteMode.value,
  storage_type: storageState.activeType,
  storage_local_config: JSON.stringify(storageState.drafts.local),
  storage_aliyun_config: JSON.stringify(storageState.drafts.aliyun),
  storage_tencent_config: JSON.stringify(storageState.drafts.tencent),
  storage_minio_config: JSON.stringify(storageState.drafts.minio),
})

const baselineSnapshot = ref(createSnapshot(getFileSettingsState()))
const hasUnsavedChanges = computed(() => isSnapshotDirty(baselineSnapshot.value, getFileSettingsState()))

watch(hasUnsavedChanges, (value) => {
  emit('dirty-change', value)
}, { immediate: true })

const save = async () => {
  saving.value = true
  try {
    await configStore.updateConfigs(buildFullSavePayload())
    baselineSnapshot.value = createSnapshot(getFileSettingsState())
    message.success('文件设置保存成功')
    return true
  } catch {
    message.error('保存失败')
    return false
  } finally {
    saving.value = false
  }
}

const discardChanges = () => {
  const restored = cloneFromSnapshot<ReturnType<typeof getFileSettingsState>>(baselineSnapshot.value)
  applyFileSettingsState(restored)
  configVisible.value = false
}

const closeTransientUi = () => {
  configVisible.value = false
}

const handleSave = async () => {
  await save()
}

const handleTestStorage = async () => {
  storageTesting.value = true
  try {
    await testStorageConfig({
      type: storageType.value,
      config: JSON.stringify(storageConfig.value),
    })
    storageValidated.value = true
    message.success('连接测试成功')
  } catch (error: any) {
    storageValidated.value = false
    message.error(error?.message || '连接测试失败')
  } finally {
    storageTesting.value = false
  }
}

const handleStorageSave = async () => {
  if (storageType.value === 'local') {
    await doStorageSave()
    return
  }

  if (!storageValidated.value) {
    storageTesting.value = true
    try {
      await testStorageConfig({
        type: storageType.value,
        config: JSON.stringify(storageConfig.value),
      })
      storageValidated.value = true
    } catch (error: any) {
      storageValidated.value = false
      message.error(error?.message || '连接验证失败，请检查配置后重试')
      return
    } finally {
      storageTesting.value = false
    }
  }

  await doStorageSave()
}

const doStorageSave = async () => {
  storageSaving.value = true
  try {
    await configStore.updateConfigs(buildFullSavePayload())
    baselineSnapshot.value = createSnapshot(getFileSettingsState())
    message.success('存储配置保存成功')
    configVisible.value = false
  } catch {
    message.error('保存失败')
  } finally {
    storageSaving.value = false
  }
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

.storage-notice {
  margin-bottom: 16px;
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: #999;
}

.validated-btn {
  color: #52c41a;
  border-color: #b7eb8f;
}
</style>
