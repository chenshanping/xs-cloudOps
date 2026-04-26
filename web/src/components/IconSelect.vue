<template>
  <div class="icon-select">
    <a-tabs v-model:activeKey="activeTab" size="small">
      <!-- Ant Design Icons 页签 -->
      <a-tab-pane key="official" tab="Ant Design Icons">
        <div class="official-icons-wrapper">
          <div style="margin-bottom: 12px; display: flex; gap: 8px;">
            <a-input
              v-model:value="searchText"
              placeholder="搜索图标名称"
              allow-clear
              style="flex: 1;"
            />
            <a-select v-model:value="iconTypeFilter" style="width: 120px;">
              <a-select-option value="all">所有类型</a-select-option>
              <a-select-option value="Outlined">Outlined</a-select-option>
              <a-select-option value="Filled">Filled</a-select-option>
              <a-select-option value="TwoTone">TwoTone</a-select-option>
            </a-select>
          </div>
          <a-empty v-if="filteredOfficialIcons.length === 0" description="没有找到图标" style="margin-top: 20px" />
          <div v-else class="icon-grid">
            <div
              v-for="icon in filteredOfficialIcons"
              :key="`official-${icon.full}`"
              class="icon-item"
              :class="{ active: modelValue === `official-${icon.full}` }"
              @click="$emit('update:modelValue', `official-${icon.full}`)"
              :title="`${icon.name} (${icon.type})`"
            >
              <component :is="getIconComponent(icon.full)" />
              <span class="icon-name">{{ icon.name }}</span>
              <span class="icon-type">{{ icon.type }}</span>
            </div>
          </div>
        </div>
      </a-tab-pane>
      
      <!-- 自定义SVG 页签 -->
      <a-tab-pane key="custom" tab="自定义SVG">
        <a-input
          v-model:value="searchText"
          placeholder="搜索图标"
          style="margin-bottom: 12px"
          allow-clear
        />
        <a-empty v-if="filteredCustomIcons.length === 0" description="没有找到自定义图标" style="margin-top: 20px" />
        <div v-else class="icon-grid">
          <div
            v-for="icon in filteredCustomIcons"
            :key="`custom-${icon}`"
            class="icon-item"
            :class="{ active: modelValue === `custom-${icon}` }"
            @click="$emit('update:modelValue', `custom-${icon}`)"
            :title="icon"
          >
            <img :src="getCustomIconUrl(icon)" alt="" style="width: 24px; height: 24px" />
            <span class="icon-name">{{ icon }}</span>
          </div>
        </div>
      </a-tab-pane>
    </a-tabs>
    
    <!-- 图标预览 -->
    <div v-if="modelValue" class="icon-preview">
      <span style="margin-right: 8px">已选择：</span>
      <component v-if="modelValue.startsWith('official-') && previewIconComponent" :is="previewIconComponent" style="font-size: 24px" />
      <img v-else-if="modelValue.startsWith('custom-')" :src="getCustomIconUrl(currentIconName)" style="width: 24px; height: 24px" alt="" />
      <span style="margin-left: 8px">{{ displayName }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import * as AntIcons from '@ant-design/icons-vue'
import { allAntIcons } from '@/utils/ant-icons'

const props = defineProps<{
  modelValue: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

// 动态生成的官方图标列表（包含所有类型）
const officialIcons = ref(allAntIcons)

// 自定义SVG图标列表
const customIcons = ref<string[]>([])
const activeTab = ref('official')
const searchText = ref('')
const iconTypeFilter = ref('all')

// 预加载所有自定义图标
const iconModules = import.meta.glob('@/assets/icons/*.svg', { eager: true, query: '?url', import: 'default' })

// 动态扫描自定义SVG图标
const loadCustomIcons = () => {
  try {
    const icons = Object.keys(iconModules)
      .map(path => path.replace('/src/assets/icons/', '').replace('.svg', ''))
      .filter(name => name.length > 0)
    customIcons.value = [...new Set(icons)]
  } catch (error) {
    console.error('加载自定义SVG图标失败:', error)
    customIcons.value = []
  }
}

loadCustomIcons()

const isOfficialIcon = computed(() => {
  return props.modelValue?.startsWith('official-') || false
})

const currentIconName = computed(() => {
  const value = props.modelValue || ''
  if (value.startsWith('official-')) {
    return value.replace('official-', '')
  }
  if (value.startsWith('custom-')) {
    return value.replace('custom-', '')
  }
  return ''
})

const displayName = computed(() => {
  const value = props.modelValue || ''
  if (value.startsWith('official-')) {
    const fullName = value.replace('official-', '')
    const icon = officialIcons.value.find(i => i.full === fullName)
    return icon ? `${icon.name} (${icon.type})` : fullName
  }
  return currentIconName.value
})

const previewIconComponent = computed(() => {
  const value = props.modelValue || ''
  if (!value.startsWith('official-')) return null
  const iconName = value.replace('official-', '')
  return (AntIcons as any)[iconName] || null
})

const filteredOfficialIcons = computed(() => {
  let filtered = officialIcons.value

  // 类型过滤
  if (iconTypeFilter.value !== 'all') {
    filtered = filtered.filter(icon => icon.type === iconTypeFilter.value)
  }

  // 搜索过滤
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    filtered = filtered.filter(icon =>
      icon.name.toLowerCase().includes(search) ||
      icon.full.toLowerCase().includes(search)
    )
  }

  return filtered
})

const filteredCustomIcons = computed(() => {
  if (!searchText.value) return customIcons.value
  return customIcons.value.filter(icon =>
    icon.toLowerCase().includes(searchText.value.toLowerCase())
  )
})

// 切换tab时清空搜索
watch(activeTab, () => {
  searchText.value = ''
})

const getCustomIconUrl = (iconName: string) => {
  const key = `/src/assets/icons/${iconName}.svg`
  return (iconModules[key] as string) || ''
}

const getIconComponent = (iconName: string) => {
  return (AntIcons as any)[iconName] || null
}
</script>

<style scoped>
.icon-select {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.official-icons-wrapper {
  max-height: 400px;
  overflow-y: auto;
}

.icon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 8px;
  max-height: 350px;
  overflow-y: auto;
  padding: 8px;
}

.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 12px 8px;
  border: 2px solid #f0f0f0;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 24px;
  min-height: 90px;
}

.icon-item:hover {
  border-color: #1890ff;
  background-color: #f5f7ff;
}

.icon-item.active {
  border-color: #1890ff;
  background-color: #e6f7ff;
}

.icon-name {
  font-size: 11px;
  color: #666;
  margin-top: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  text-align: center;
  font-weight: 500;
}

.icon-type {
  font-size: 9px;
  color: #999;
  margin-top: 2px;
}

.icon-preview {
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f5f5f5;
  border-radius: 4px;
  border-left: 4px solid #1890ff;
}
</style>
