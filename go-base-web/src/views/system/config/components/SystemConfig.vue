<template>
  <div class="system-config">
    <div class="config-layout">
      <!-- 左侧表单 -->
      <div class="config-form">
        <a-form :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
          <!-- 基本设置 -->
          <a-divider orientation="left">基本设置</a-divider>
          
          <a-form-item label="系统名称">
            <a-input v-model:value="formData.sys_name" placeholder="请输入系统名称" />
          </a-form-item>
          
          <a-form-item label="系统Logo">
            <ImageUpload 
              v-model="formData.sys_logo" 
              :width="120" 
              :height="60" 
              :max-size="5*1024*1024"
              placeholder="上传Logo"
            />
          </a-form-item>
          
          <a-form-item label="前台模式">
            <a-radio-group v-model:value="formData.front_mode">
              <a-radio value="full">完整前台</a-radio>
              <a-radio value="profile">仅个人中心</a-radio>
            </a-radio-group>
            <div class="form-tip">
              完整前台: 显示全部前台页面；仅个人中心: 只显示个人资料页面（用于身份认证）
            </div>
          </a-form-item>

          

          <!-- 菜单配置 -->
          <a-divider orientation="left">菜单配置</a-divider>
          
          <a-form-item label="菜单背景色">
            <div class="color-input">
              <a-input v-model:value="formData.menu_bg_color" placeholder="#001529 或 linear-gradient(...)" />
              <input 
                type="color" 
                :value="extractColor(formData.menu_bg_color)" 
                @input="updateMenuBgColor" 
                class="color-picker" 
                title="点击选择纯色"
              />
            </div>
            <div class="form-tip">支持渐变色，如: linear-gradient(135deg, #667eea, #764ba2)</div>
          </a-form-item>
          
          <a-form-item label="菜单文字颜色">
            <div class="color-input">
              <a-input v-model:value="formData.menu_text_color" />
              <input type="color" :value="rgbaToHex(formData.menu_text_color)" @input="formData.menu_text_color = $event.target.value" class="color-picker" />
            </div>
          </a-form-item>
          
          <a-form-item label="菜单激活背景色">
            <div class="color-input">
              <a-input v-model:value="formData.menu_active_bg_color" />
              <input type="color" v-model="formData.menu_active_bg_color" class="color-picker" />
            </div>
          </a-form-item>
          
          <a-form-item label="菜单激活文字颜色">
            <div class="color-input">
              <a-input v-model:value="formData.menu_active_text_color" />
              <input type="color" v-model="formData.menu_active_text_color" class="color-picker" />
            </div>
          </a-form-item>

          <!-- 头部配置 -->
          <a-divider orientation="left">头部配置</a-divider>
          
          <a-form-item label="头部背景色">
            <div class="color-input">
              <a-input v-model:value="formData.header_bg_color" placeholder="#ffffff 或 linear-gradient(...)" />
              <input 
                type="color" 
                :value="extractColor(formData.header_bg_color)" 
                @input="updateHeaderBgColor" 
                class="color-picker" 
                title="点击选择纯色"
              />
            </div>
            <div class="form-tip">支持渐变色</div>
          </a-form-item>
          
          <a-form-item label="头部文字颜色">
            <div class="color-input">
              <a-input v-model:value="formData.header_text_color" />
              <input type="color" v-model="formData.header_text_color" class="color-picker" />
            </div>
          </a-form-item>

          <!-- 操作按钮 -->
          <a-form-item :wrapper-col="{ offset: 8, span: 16 }" style="margin-top: 24px">
            <a-space>
              <a-button type="primary" :loading="saving" @click="handleSave">保存配置</a-button>
              <!-- <a-button @click="handleReset">恢复默认</a-button> -->
            </a-space>
          </a-form-item>
        </a-form>
      </div>

      <!-- 右侧预览 -->
      <div class="config-preview">
        <div class="preview-title">布局预览</div>
        <!-- 布局预览 -->
        <div class="preview-container">
          <div class="preview-sidebar" :style="{ background: formData.menu_bg_color }">
            <div class="preview-logo">
              <img v-if="formData.sys_logo" :src="formData.sys_logo" alt="logo" />
              <span :style="{ color: formData.menu_text_color }">{{ formData.sys_name }}</span>
            </div>
            <div class="preview-menu">
              <div class="preview-menu-item" :style="{ color: formData.menu_text_color }">
                普通菜单项
              </div>
              <div 
                class="preview-menu-item active" 
                :style="{ 
                  background: formData.menu_active_bg_color, 
                  color: formData.menu_active_text_color 
                }"
              >
                激活菜单项
              </div>
              <div class="preview-menu-item" :style="{ color: formData.menu_text_color }">
                普通菜单项
              </div>
            </div>
          </div>
          <div class="preview-main">
            <div 
              class="preview-header" 
              :style="{ 
                background: formData.header_bg_color, 
                color: formData.header_text_color 
              }"
            >
              头部区域
            </div>
            <div class="preview-content">内容区域</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { useConfigStore } from '@/store/config'
import ImageUpload from '@/components/ImageUpload.vue'

const configStore = useConfigStore()
const route = useRoute()
const saving = ref(false)

// 系统配置相关键
const SYSTEM_CONFIG_KEYS = [
  'sys_name',
  'sys_logo',
  'register_logo',
  'menu_bg_color',
  'menu_text_color',
  'menu_active_text_color',
  'menu_active_bg_color',
  'header_bg_color',
  'header_text_color',
  'front_mode'
] as const

// 直接从 store 初始化数据
const formData = reactive({
  sys_name: configStore.get('sys_name'),
  sys_logo: configStore.get('sys_logo'),
  register_logo: configStore.get('register_logo'),
  menu_bg_color: configStore.get('menu_bg_color'),
  menu_text_color: configStore.get('menu_text_color'),
  menu_active_text_color: configStore.get('menu_active_text_color'),
  menu_active_bg_color: configStore.get('menu_active_bg_color'),
  header_bg_color: configStore.get('header_bg_color'),
  header_text_color: configStore.get('header_text_color'),
  front_mode: configStore.get('front_mode') || 'full'
})
const updateHeaderBgColor = (e: Event) => {
  formData.header_bg_color = (e.target as HTMLInputElement).value
}
const updateMenuBgColor = (e: Event) => {
  formData.menu_bg_color = (e.target as HTMLInputElement).value
  
}
// 从颜色值中提取纯色（用于颜色选择器显示）
const extractColor = (color: string): string => {
  if (!color) return '#ffffff'
  // 如果是纯 hex 色
  if (color.startsWith('#')) {
    return color.length === 4 
      ? `#${color[1]}${color[1]}${color[2]}${color[2]}${color[3]}${color[3]}` 
      : color
  }
  // 从渐变色中提取第一个颜色
  const hexMatch = color.match(/#[0-9a-fA-F]{3,6}/)
  if (hexMatch) return hexMatch[0].length === 4 
    ? `#${hexMatch[0][1]}${hexMatch[0][1]}${hexMatch[0][2]}${hexMatch[0][2]}${hexMatch[0][3]}${hexMatch[0][3]}`
    : hexMatch[0]
  // rgba 转 hex
  const rgbaMatch = color.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)/)
  if (rgbaMatch) {
    const r = parseInt(rgbaMatch[1]).toString(16).padStart(2, '0')
    const g = parseInt(rgbaMatch[2]).toString(16).padStart(2, '0')
    const b = parseInt(rgbaMatch[3]).toString(16).padStart(2, '0')
    return `#${r}${g}${b}`
  }
  return '#ffffff'
}

// rgba 转 hex（简单处理）
const rgbaToHex = (rgba: string): string => {
  return extractColor(rgba)
}

// 更新页面标题
const updateTitle = () => {
  const sysName = configStore.get('sys_name') || 'Go RBAC Admin'
  const pageTitle = route.meta?.title as string
  document.title = pageTitle ? `${pageTitle} - ${sysName}` : sysName
}

// 保存配置
const handleSave = async () => {
  saving.value = true
  try {
    const configs: Record<string, string> = {}
    for (const key of SYSTEM_CONFIG_KEYS) {
      configs[key] = (formData as any)[key]
    }
    await configStore.updateConfigs(configs)
    updateTitle()
    message.success('配置保存成功')
  } catch (e) {
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 恢复默认
const handleReset = () => {
  formData.sys_name = 'Go RBAC Admin'
  formData.sys_logo = '/src/assets/logo.svg'
  formData.register_logo = ''
  formData.menu_bg_color = '#001529'
  formData.menu_text_color = 'rgba(255, 255, 255, 0.65)'
  formData.menu_active_text_color = '#ffffff'
  formData.menu_active_bg_color = '#1890ff'
  formData.header_bg_color = '#ffffff'
  formData.header_text_color = '#333333'
}

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
  border-bottom: 1px solid #e8e8e8;
}

.color-input {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-picker {
  width: 32px;
  height: 32px;
  padding: 0;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  cursor: pointer;
  flex-shrink: 0;
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: #999;
}

.preview-container {
  display: flex;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  overflow: hidden;
  height: 300px;
}

.preview-sidebar {
  width: 140px;
  display: flex;
  flex-direction: column;
}

.preview-logo {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  background: rgba(0, 0, 0, 0.2);
}

.preview-logo img {
  width: 20px;
  height: 20px;
}

.preview-logo span {
  font-size: 12px;
  font-weight: 500;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-menu {
  flex: 1;
  padding: 8px 0;
}

.preview-menu-item {
  padding: 10px 12px;
  font-size: 12px;
  cursor: pointer;
}

.preview-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #f0f2f5;
}

.preview-header {
  height: 36px;
  display: flex;
  align-items: center;
  padding: 0 12px;
  font-size: 12px;
}

.preview-content {
  flex: 1;
  margin: 8px;
  padding: 12px;
  background: #fff;
  border-radius: 4px;
  font-size: 12px;
  color: #666;
}

.preview-title {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e8e8e8;
}
</style>
