<template>
  <div class="login-register-config">
    <div class="config-layout">
      <!-- 左侧表单 -->
      <div class="config-form">
        <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
          <!-- 注册配置 -->
          <a-divider orientation="left">注册配置</a-divider>
          
          <a-form-item label="开启注册功能">
            <a-switch v-model:checked="enableRegister" />
            <div class="form-tip">开启后，登录页面将显示注册和忘记密码入口</div>
          </a-form-item>
          <a-form-item label="注册用户默认头像">
            <ImageUpload 
              v-model="formData.register_logo" 
              :width="120" 
              :height="60" 
              placeholder="注册用户默认头像"
            />
          </a-form-item>

          <!-- 登录页配置 -->
          <a-divider orientation="left">登录页配置</a-divider>
          
          <a-form-item label="登录页标题">
            <a-input v-model:value="formData.login_title" placeholder="欢迎回来" />
          </a-form-item>
          
          <a-form-item label="标语">
            <a-input v-model:value="formData.login_slogan" placeholder="智能化企业管理平台" />
          </a-form-item>
          
          <a-form-item label="背景图片">
            <ImageUpload 
              v-model="formData.login_bg_image" 
              :width="160" 
              :height="90"
              :max-size="10 * 1024 * 1024"
              placeholder="上传背景图"
            />
          </a-form-item>

          <!-- 操作按钮 -->
          <a-form-item :wrapper-col="{ offset: 8, span: 16 }" style="margin-top: 24px">
            <a-space>
              <a-button type="primary" :loading="saving" @click="handleSave">保存配置</a-button>
              <a-button @click="advancedDrawerVisible = true">
                <SettingOutlined /> 高级设置
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </div>

      <!-- 右侧预览 -->
      <div class="config-preview">
        <div class="preview-title">登录页预览</div>
        <div class="login-preview-new">
          <!-- 左侧背景区 -->
          <div 
            class="login-preview-left" 
            :style="formData.login_bg_image ? { backgroundImage: `url(${formData.login_bg_image})` } : {}"
          >
            <div class="preview-overlay">
              <div class="preview-content">
                <img v-if="configStore.get('sys_logo')" :src="configStore.get('sys_logo')" alt="Logo" class="preview-logo" />
                <div class="preview-sys-name">{{ configStore.get('sys_name') || 'Go RBAC Admin' }}</div>
                <div class="preview-slogan">{{ formData.login_slogan || '智能化企业管理平台' }}</div>
                <div class="preview-desc" v-if="formData.login_desc">{{ formData.login_desc }}</div>
                <div class="preview-features" v-if="featureList.length > 0">
                  <div class="preview-feature" v-for="(f, i) in featureList.slice(0, 4)" :key="i">
                    <component :is="getIconComponent(f.icon)" class="preview-feature-icon" />
                    <span>{{ f.title }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <!-- 右侧表单区 -->
          <div class="login-preview-right">
            <div class="preview-form-title">{{ formData.login_title || '欢迎回来' }}</div>
            <div class="preview-form-subtitle">登录您的账户</div>
            <div class="preview-form-input"><UserOutlined /> 用户名</div>
            <div class="preview-form-input"><LockOutlined /> 密码</div>
            <div class="preview-form-btn">登 录</div>
            <div class="preview-form-links" v-if="enableRegister">
              <span class="preview-link">注册账号</span>
              <span class="preview-link">忘记密码</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 高级设置抽屉 -->
    <a-drawer
      v-model:open="advancedDrawerVisible"
      title="登录页高级设置"
      placement="right"
      :width="500"
    >
      <a-form :label-col="{ span: 6 }" :wrapper-col="{ span: 18 }">
        <!-- 登录页副标题 -->
        <a-form-item label="副标题">
          <a-input v-model:value="formData.login_subtitle" placeholder="企业级后台管理系统" />
        </a-form-item>
        
        <a-form-item label="描述">
          <a-textarea 
            v-model:value="formData.login_desc" 
            placeholder="平台描述文字" 
            :rows="2"
          />
        </a-form-item>
        
        <a-form-item label="背景渐变色">
          <a-input v-model:value="formData.login_bg_color" placeholder="linear-gradient(135deg, #667eea 0%, #764ba2 100%)" />
          <div class="form-tip">上传背景图后将覆盖渐变色</div>
        </a-form-item>

        <!-- 特性标签 -->
        <a-divider orientation="left">特性标签</a-divider>
        
        <a-form-item label="最大数量">
          <a-input-number v-model:value="formData.login_features_max" :min="1" :max="10" />
        </a-form-item>
        
        <div class="feature-editor-drawer">
          <div 
            v-for="(feature, index) in featureList" 
            :key="index" 
            class="feature-card"
          >
            <div class="feature-card-header">
              <a-select 
                v-model:value="feature.icon" 
                style="width: 100px"
                size="small"
                placeholder="图标"
              >
                <a-select-option v-for="icon in iconOptions" :key="icon.value" :value="icon.value">
                  <component :is="getIconComponent(icon.value)" /> {{ icon.label }}
                </a-select-option>
              </a-select>
              <a-button type="text" danger size="small" @click="removeFeature(index)">
                <DeleteOutlined />
              </a-button>
            </div>
            <a-input v-model:value="feature.title" placeholder="标题" size="small" />
            <a-input v-model:value="feature.desc" placeholder="描述" size="small" />
          </div>
          <a-button 
            type="dashed" 
            block
            @click="addFeature" 
            :disabled="featureList.length >= formData.login_features_max"
          >
            <PlusOutlined /> 添加特性 ({{ featureList.length }}/{{ formData.login_features_max }})
          </a-button>
        </div>

        <!-- 展示图片 -->
        <a-divider orientation="left">展示图片</a-divider>
        
        <a-form-item label="最大数量">
          <a-input-number v-model:value="formData.login_images_max" :min="1" :max="10" />
        </a-form-item>
        
        <div class="image-editor-drawer">
          <div 
            v-for="(img, index) in imageList" 
            :key="index" 
            class="image-card"
          >
            <ImageUpload 
              v-model="img.url" 
              :width="60" 
              :height="40" 
              placeholder="上传"
            />
            <div class="image-card-info">
              <a-input v-model:value="img.title" placeholder="图片标题" size="small" />
              <a-button type="text" danger size="small" @click="removeImage(index)">
                <DeleteOutlined />
              </a-button>
            </div>
          </div>
          <a-button 
            type="dashed" 
            block
            @click="addImage" 
            :disabled="imageList.length >= formData.login_images_max"
          >
            <PlusOutlined /> 添加图片 ({{ imageList.length }}/{{ formData.login_images_max }})
          </a-button>
        </div>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useConfigStore } from '@/store/config'
import { 
  UserOutlined, 
  LockOutlined,
  DeleteOutlined,
  PlusOutlined,
  CheckCircleOutlined,
  SafetyOutlined,
  LineChartOutlined,
  ThunderboltOutlined,
  RocketOutlined,
  SettingOutlined,
  CloudOutlined,
  TeamOutlined
} from '@ant-design/icons-vue'
import ImageUpload from '@/components/ImageUpload.vue'
import { cloneFromSnapshot, createSnapshot, isSnapshotDirty } from '../config-tab-guard'

const configStore = useConfigStore()
const emit = defineEmits<{
  (e: 'dirty-change', value: boolean): void
}>()
const saving = ref(false)
const advancedDrawerVisible = ref(false)

// 是否开启注册
const enableRegister = ref(configStore.get('enable_register') === 'true')
// const enableRegister = ref(false)

// 图标选项
const iconOptions = [
  { value: 'CheckCircleOutlined', label: '成功' },
  { value: 'SafetyOutlined', label: '安全' },
  { value: 'LineChartOutlined', label: '图表' },
  { value: 'ThunderboltOutlined', label: '闪电' },
  { value: 'RocketOutlined', label: '火箭' },
  { value: 'SettingOutlined', label: '设置' },
  { value: 'CloudOutlined', label: '云' },
  { value: 'TeamOutlined', label: '团队' }
]

// 图标映射
const iconMap: Record<string, any> = {
  CheckCircleOutlined,
  SafetyOutlined,
  LineChartOutlined,
  ThunderboltOutlined,
  RocketOutlined,
  SettingOutlined,
  CloudOutlined,
  TeamOutlined
}

const getIconComponent = (name: string) => iconMap[name] || CheckCircleOutlined

// 登录页相关配置键
const LOGIN_CONFIG_KEYS = [
  'login_bg_image',
  'login_title',
  'login_subtitle',
  'login_bg_color',
  'login_slogan',
  'login_desc',
  'login_features',
  'login_features_max',
  'login_images',
  'login_images_max',
  'enable_register'
] as const

// 直接从 store 初始化数据
const formData = reactive({
  login_bg_image: configStore.get('login_bg_image'),
  login_title: configStore.get('login_title'),
  login_subtitle: configStore.get('login_subtitle'),
  login_bg_color: configStore.get('login_bg_color'),
  login_slogan: configStore.get('login_slogan'),
  login_desc: configStore.get('login_desc'),
  login_features: configStore.get('login_features'),
  login_features_max: parseInt(configStore.get('login_features_max')) || 4,
  login_images: configStore.get('login_images'),
  login_images_max: parseInt(configStore.get('login_images_max')) || 4,
  register_logo: configStore.get('register_logo')
})

// 特性标签列表
interface FeatureItem {
  icon: string
  title: string
  desc: string
}
const featureList = ref<FeatureItem[]>([])

// 图片列表
interface ImageItem {
  url: string
  title: string
}
const imageList = ref<ImageItem[]>([])

// 添加/删除特性
const addFeature = () => {
  if (featureList.value.length < formData.login_features_max) {
    featureList.value.push({ icon: 'CheckCircleOutlined', title: '', desc: '' })
  }
}
const removeFeature = (index: number) => {
  featureList.value.splice(index, 1)
}

// 添加/删除图片
const addImage = () => {
  if (imageList.value.length < formData.login_images_max) {
    imageList.value.push({ url: '', title: '' })
  }
}
const removeImage = (index: number) => {
  imageList.value.splice(index, 1)
}

// 同步特性列表到formData
watch(featureList, (val) => {
  formData.login_features = JSON.stringify(val)
}, { deep: true })

// 同步图片列表到formData
watch(imageList, (val) => {
  formData.login_images = JSON.stringify(val)
}, { deep: true })

// 初始化解析特性列表
try {
  featureList.value = formData.login_features ? JSON.parse(formData.login_features) : []
} catch {
  featureList.value = []
}
// 初始化解析图片列表
try {
  imageList.value = formData.login_images ? JSON.parse(formData.login_images) : []
} catch {
  imageList.value = []
}

const getConfigState = () => ({
  register_logo: formData.register_logo || '',
  login_bg_image: formData.login_bg_image,
  login_title: formData.login_title,
  login_subtitle: formData.login_subtitle,
  login_bg_color: formData.login_bg_color,
  login_slogan: formData.login_slogan,
  login_desc: formData.login_desc,
  login_features: JSON.stringify(featureList.value),
  login_features_max: formData.login_features_max,
  login_images: JSON.stringify(imageList.value),
  login_images_max: formData.login_images_max,
  enable_register: enableRegister.value ? 'true' : 'false',
})

const applyConfigState = (state: ReturnType<typeof getConfigState>) => {
  enableRegister.value = state.enable_register === 'true'
  formData.register_logo = state.register_logo
  formData.login_bg_image = state.login_bg_image
  formData.login_title = state.login_title
  formData.login_subtitle = state.login_subtitle
  formData.login_bg_color = state.login_bg_color
  formData.login_slogan = state.login_slogan
  formData.login_desc = state.login_desc
  formData.login_features_max = Number(state.login_features_max) || 4
  formData.login_images_max = Number(state.login_images_max) || 4
  try {
    featureList.value = state.login_features ? JSON.parse(state.login_features) : []
  } catch {
    featureList.value = []
  }
  try {
    imageList.value = state.login_images ? JSON.parse(state.login_images) : []
  } catch {
    imageList.value = []
  }
}

const baselineSnapshot = ref(createSnapshot(getConfigState()))
const hasUnsavedChanges = computed(() => isSnapshotDirty(baselineSnapshot.value, getConfigState()))

watch(hasUnsavedChanges, (value) => {
  emit('dirty-change', value)
}, { immediate: true })

// 保存配置
const save = async () => {
  // 验证：开启注册功能时需要填写默认头像
  if (enableRegister.value && !formData.register_logo) {
    message.warning('开启注册功能时，请上传注册用户默认头像')
    return false
  }
  
  saving.value = true
  try {
    const configs: Record<string, string> = {}
    for (const [key, value] of Object.entries(getConfigState())) {
      configs[key] = typeof value === 'number' ? String(value) : value
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
  advancedDrawerVisible.value = false
}

const closeTransientUi = () => {
  advancedDrawerVisible.value = false
}

const handleSave = async () => {
  await save()
}

// 恢复默认
const handleReset = () => {
  enableRegister.value = false
  formData.login_bg_image = ''
  formData.login_title = '欢迎回来'
  formData.login_subtitle = '企业级后台管理系统'
  formData.login_bg_color = 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  formData.login_slogan = '智能化企业管理平台'
  formData.login_desc = '实时数据联动，智能分析策略，多渠道展示整合，让管理更简单，让体验更美好。'
  formData.login_features_max = 4
  formData.login_images_max = 4
  featureList.value = [
    { icon: 'CheckCircleOutlined', title: '智能分析', desc: '动态调整业务策略' },
    { icon: 'SafetyOutlined', title: '安全合规', desc: '多层级权限保障' },
    { icon: 'LineChartOutlined', title: '数据洞察', desc: '全方位数据分析' },
    { icon: 'ThunderboltOutlined', title: '高效管理', desc: '简化日常运营流程' }
  ]
  imageList.value = []
}

defineExpose({
  isDirty: () => hasUnsavedChanges.value,
  save,
  discardChanges,
  closeTransientUi,
})

</script>

<style scoped>
.config-layout {
  display: flex;
  gap: 32px;
}

.config-form {
  flex: 1;
  min-width: 400px;
  max-width: 500px;
}

.config-preview {
  width: 380px;
  flex-shrink: 0;
  position: sticky;
  top: 16px;
  align-self: flex-start;
}

.preview-title {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--app-border-color);
  color: var(--app-text-strong);
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--app-text-muted);
}

/* 登录页预览 */
.login-preview-new {
  height: 300px;
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.login-preview-left {
  flex: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-size: cover;
  background-position: center;
  position: relative;
}

.preview-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-content {
  text-align: center;
  color: #fff;
  padding: 16px;
}

.preview-logo {
  width: 32px;
  height: 32px;
  margin-bottom: 8px;
  border-radius: 6px;
}

.preview-sys-name {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 4px;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
}

.preview-slogan {
  font-size: 10px;
  color: #fff;
  opacity: 0.9;
  margin-bottom: 6px;
}

.preview-desc {
  font-size: 9px;
  color: #fff;
  opacity: 0.8;
  margin-bottom: 12px;
  max-width: 180px;
  line-height: 1.4;
}

.preview-features {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 6px;
  text-align: left;
}

.preview-feature {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 6px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(4px);
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  font-size: 8px;
  color: #fff;
}

.preview-feature-icon {
  color: #fff;
  font-size: 10px;
}

.login-preview-right {
  width: 140px;
  background: var(--app-surface-color);
  padding: 20px 12px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.05);
}

.preview-form-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--app-text-strong);
  margin-bottom: 2px;
  text-align: center;
}

.preview-form-subtitle {
  font-size: 9px;
  color: var(--app-text-secondary);
  margin-bottom: 12px;
  text-align: center;
}

.preview-form-input {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  font-size: 10px;
  color: var(--app-text-muted);
  margin-bottom: 8px;
}

.preview-form-btn {
  padding: 6px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  text-align: center;
  border-radius: 4px;
  font-size: 10px;
}

.preview-form-links {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
}

.preview-link {
  font-size: 9px;
  color: #667eea;
  cursor: pointer;
}

/* 抽屉内特性编辑器 */
.feature-editor-drawer {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.feature-card {
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.feature-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 抽屉内图片编辑器 */
.image-editor-drawer {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.image-card {
  background: var(--app-surface-soft);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
  padding: 12px;
  display: flex;
  gap: 12px;
  align-items: center;
}

.image-card-info {
  flex: 1;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 8px;
}
</style>
