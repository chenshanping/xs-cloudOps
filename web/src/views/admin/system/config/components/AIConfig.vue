<template>
  <div class="ai-config">
    <a-spin :spinning="loading">
      <a-form :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
        <!-- 平台配置 -->
        <a-collapse v-model:activeKey="activeProviders" accordion>
          <a-collapse-panel v-for="(provider, pIndex) in formData.providers" :key="pIndex">
            <template #header>
              <span>{{ provider.name || '新平台' }}</span>
              <a-tag v-if="isDefaultProvider(pIndex)" color="blue" style="margin-left: 8px">默认</a-tag>
            </template>
            <template #extra>
              <a-space @click.stop>
                <a-button 
                  v-if="!isDefaultProvider(pIndex) && provider.name" 
                  type="link" 
                  size="small" 
                  @click="setDefaultProvider(pIndex)"
                >
                  设为默认
                </a-button>
                <a-button type="link" danger size="small" @click="removeProvider(pIndex)">
                  <template #icon><DeleteOutlined /></template>
                </a-button>
              </a-space>
            </template>
            
            <a-form-item label="平台名称" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
              <a-input v-model:value="provider.name" placeholder="如：阿里云百炼" />
            </a-form-item>
            
            <a-form-item label="API Key" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
              <a-input-password v-model:value="provider.api_key" placeholder="请输入API Key" />
            </a-form-item>
            
            <a-form-item label="Base URL" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
              <a-input v-model:value="provider.base_url" placeholder="如：https://dashscope.aliyuncs.com/compatible-mode/v1" />
            </a-form-item>

            <!-- 模型列表 -->
            <a-form-item label="模型列表" :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
              <div class="models-list">
                <div class="models-header">
                  <span class="col-id">模型ID</span>
                  <span class="col-name">显示名称</span>
                  <span class="col-desc">模型描述</span>
                  <span class="col-actions">操作</span>
                </div>
                <div v-for="(model, mIndex) in provider.models" :key="mIndex" class="model-item">
                  <a-input v-model:value="model.id" placeholder="模型ID" size="small" class="col-id" />
                  <a-input v-model:value="model.name" placeholder="显示名称" size="small" class="col-name" />
                  <a-input v-model:value="model.description" placeholder="模型描述" size="small" class="col-desc" />
                  <div class="col-actions">
                    <a-tooltip title="上移">
                      <a-button 
                        type="text" 
                        size="small" 
                        :disabled="mIndex === 0"
                        @click="moveModel(pIndex, mIndex, -1)"
                      >
                        <template #icon><UpOutlined /></template>
                      </a-button>
                    </a-tooltip>
                    <a-tooltip title="下移">
                      <a-button 
                        type="text" 
                        size="small" 
                        :disabled="mIndex === provider.models.length - 1"
                        @click="moveModel(pIndex, mIndex, 1)"
                      >
                        <template #icon><DownOutlined /></template>
                      </a-button>
                    </a-tooltip>
                    <a-tooltip title="测试">
                      <a-button 
                        type="text" 
                        size="small"
                        :loading="testingModel === `${pIndex}-${mIndex}`"
                        @click="testModel(pIndex, mIndex)"
                      >
                        <template #icon><ThunderboltOutlined /></template>
                      </a-button>
                    </a-tooltip>
                    <a-tooltip title="删除">
                      <a-button type="text" danger size="small" @click="removeModel(pIndex, mIndex)">
                        <template #icon><DeleteOutlined /></template>
                      </a-button>
                    </a-tooltip>
                  </div>
                </div>
                <a-button type="dashed" size="small" block @click="addModel(pIndex)" style="margin-top: 8px">
                  <template #icon><PlusOutlined /></template>
                  添加模型
                </a-button>
              </div>
            </a-form-item>
          </a-collapse-panel>
        </a-collapse>

        <div style="margin-top: 16px">
          <a-button type="dashed" block @click="addProvider">
            <template #icon><PlusOutlined /></template>
            添加平台
          </a-button>
        </div>

        <!-- 操作按钮 -->
        <a-form-item :wrapper-col="{ offset: 4, span: 18 }" style="margin-top: 24px">
          <a-button type="primary" :loading="saving" @click="handleSave">保存配置</a-button>
        </a-form-item>
      </a-form>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { DeleteOutlined, PlusOutlined, UpOutlined, DownOutlined, ThunderboltOutlined } from '@ant-design/icons-vue'
import { getConfigByKey, updateConfig, createConfig, aiTest } from '@/api/config'
import { cloneFromSnapshot, createSnapshot, isSnapshotDirty } from '../config-tab-guard'

interface AIModel {
  id: string
  name: string
  description: string
}

interface AIProvider {
  name: string
  api_key: string
  base_url: string
  models: AIModel[]
}

interface AIConfig {
  default_provider: string
  providers: AIProvider[]
}

const loading = ref(true)
const saving = ref(false)
const configId = ref<number | null>(null)
const activeProviders = ref<number[]>([0])
const testingModel = ref<string | null>(null)
const initialized = ref(false)
const emit = defineEmits<{
  (e: 'dirty-change', value: boolean): void
}>()

const formData = reactive<AIConfig>({
  default_provider: '',
  providers: []
})

const getConfigState = (): AIConfig => ({
  default_provider: formData.default_provider,
  providers: formData.providers,
})

const applyConfigState = (state: AIConfig) => {
  formData.default_provider = state.default_provider || ''
  formData.providers = state.providers || []
  activeProviders.value = formData.providers.length > 0 ? [0] : []
}

const baselineSnapshot = ref(createSnapshot(getConfigState()))
const hasUnsavedChanges = computed(() => initialized.value && isSnapshotDirty(baselineSnapshot.value, getConfigState()))

watch(hasUnsavedChanges, (value) => {
  emit('dirty-change', value)
}, { immediate: true })

// 加载配置
const loadConfig = async () => {
  loading.value = true
  try {
    const res = await getConfigByKey('ai_config')
    if (res.data) {
      configId.value = res.data.id
      if (res.data.value) {
        const config = JSON.parse(res.data.value) as AIConfig
        applyConfigState(config)
      }
    }
  } catch (e: any) {
    // 配置不存在，使用默认空配置
    if (e.response?.status !== 404) {
      message.error('加载配置失败')
    }
  } finally {
    baselineSnapshot.value = createSnapshot(getConfigState())
    initialized.value = true
    loading.value = false
  }
}

// 添加平台
const addProvider = () => {
  formData.providers.push({
    name: '',
    api_key: '',
    base_url: '',
    models: []
  })
  activeProviders.value = [formData.providers.length - 1]
}

// 删除平台
const removeProvider = (index: number) => {
  const removed = formData.providers[index]
  formData.providers.splice(index, 1)
  // 如果删除的是默认平台，自动设置第一个为默认
  if (removed.name === formData.default_provider) {
    formData.default_provider = formData.providers[0]?.name || ''
  }
}

// 设置默认平台
const setDefaultProvider = (index: number) => {
  formData.default_provider = formData.providers[index].name
  message.info('已设为默认，请记得保存配置')
}

// 判断是否为默认平台
const isDefaultProvider = (index: number) => {
  return formData.providers[index].name === formData.default_provider
}

// 添加模型
const addModel = (providerIndex: number) => {
  formData.providers[providerIndex].models.push({
    id: '',
    name: '',
    description: ''
  })
}

// 删除模型
const removeModel = (providerIndex: number, modelIndex: number) => {
  formData.providers[providerIndex].models.splice(modelIndex, 1)
}

// 移动模型顺序
const moveModel = (providerIndex: number, modelIndex: number, direction: number) => {
  const models = formData.providers[providerIndex].models
  const newIndex = modelIndex + direction
  if (newIndex < 0 || newIndex >= models.length) return
  const temp = models[modelIndex]
  models[modelIndex] = models[newIndex]
  models[newIndex] = temp
}

// 测试模型
const testModel = async (providerIndex: number, modelIndex: number) => {
  const provider = formData.providers[providerIndex]
  const model = provider.models[modelIndex]
  
  if (!provider.api_key) {
    message.warning('请先填写 API Key')
    return
  }
  if (!provider.base_url) {
    message.warning('请先填写 Base URL')
    return
  }
  if (!model.id) {
    message.warning('请先填写模型 ID')
    return
  }

  testingModel.value = `${providerIndex}-${modelIndex}`
  try {
    let data={
        api_key: provider.api_key,
      base_url: provider.base_url,
      model: model.id
    }
    await aiTest(data)
    message.success(`模型 ${model.name || model.id} 测试成功`)
    
  } catch (e: any) {
    // 拦截器抛出的是 new Error(res.message)，所以用 e.message
    message.error(e.message || '测试失败')
  } finally {
    testingModel.value = null
  }
}

// 保存配置
const save = async () => {
  // 验证
  for (const p of formData.providers) {
    if (!p.name) {
      message.warning('请填写平台名称')
      return false
    }
    if (!p.base_url) {
      message.warning(`请填写 ${p.name} 的 Base URL`)
      return false
    }
  }

  saving.value = true
  try {
    const value = JSON.stringify(getConfigState())
    if (configId.value) {
      await updateConfig(configId.value, { value })
    } else {
      const res = await createConfig({
        name: 'AI配置',
        key: 'ai_config',
        value,
        value_type: 'json',
        remark: 'AI平台配置，包含平台名称、API Key、基础URL和模型列表'
      })
      configId.value = res.data.id
    }
    baselineSnapshot.value = createSnapshot(getConfigState())
    message.success('保存成功')
    return true
  } catch (e) {
    message.error('保存失败')
    return false
  } finally {
    saving.value = false
  }
}

const discardChanges = () => {
  const restored = cloneFromSnapshot<AIConfig>(baselineSnapshot.value)
  applyConfigState(restored)
}

const closeTransientUi = () => {}

const handleSave = async () => {
  await save()
}

onMounted(() => {
  loadConfig()
})

defineExpose({
  isDirty: () => hasUnsavedChanges.value,
  save,
  discardChanges,
  closeTransientUi,
})
</script>

<style scoped>
.ai-config {
  max-width: 900px;
}

.models-list {
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  background: #fafafa;
  overflow: hidden;
}

.models-header {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
  background: #f0f0f0;
  font-size: 12px;
  color: #666;
  font-weight: 500;
}

.model-item {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
  border-bottom: 1px solid #f0f0f0;
  align-items: center;
}

.model-item:hover {
  background: #fff;
}

.model-item:last-child {
  border-bottom: none;
}

.col-id {
  width: 160px;
  flex-shrink: 0;
}

.col-name {
  width: 140px;
  flex-shrink: 0;
}

.col-desc {
  flex: 1;
  min-width: 0;
}

.col-actions {
  width: 130px;
  flex-shrink: 0;
  display: flex;
  justify-content: flex-end;
  gap: 2px;
}

:deep(.ant-collapse-header) {
  align-items: center !important;
}

:deep(.ant-collapse-extra) {
  margin-left: auto;
}
</style>
