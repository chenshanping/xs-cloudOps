<template>
  <a-config-provider :locale="zhCN" :theme="antdTheme">
    <a-app>
      <router-view />
    </a-app>
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { theme as antdThemeFactory } from 'ant-design-vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { useConfigStore } from '@/store/config'
import { useUiStore } from '@/store/ui'

const configStore = useConfigStore()
const uiStore = useUiStore()

// 更新 favicon
const updateFavicon = (url: string) => {
  if (!url) return
  let link = document.querySelector("link[rel*='icon']") as HTMLLinkElement
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }
  link.href = url
}

const hexToRgba = (hex: string, alpha: number) => {
  const normalized = hex.replace('#', '')
  if (!/^[0-9a-fA-F]{6}$/.test(normalized)) {
    return `rgba(0, 107, 230, ${alpha})`
  }

  const r = Number.parseInt(normalized.slice(0, 2), 16)
  const g = Number.parseInt(normalized.slice(2, 4), 16)
  const b = Number.parseInt(normalized.slice(4, 6), 16)

  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

const updateThemeVariables = (color: string) => {
  document.documentElement.style.setProperty('--app-primary-color', color)
  document.documentElement.style.setProperty('--app-primary-color-soft', hexToRgba(color, 0.12))
}

const antdTheme = computed(() => ({
  algorithm: uiStore.isDark
    ? antdThemeFactory.darkAlgorithm
    : antdThemeFactory.defaultAlgorithm,
  token: {
    colorPrimary: uiStore.theme.primaryColor,
    colorInfo: uiStore.theme.primaryColor,
    colorLink: uiStore.theme.primaryColor,
    borderRadius: 10,
    colorBgLayout: uiStore.isDark ? '#020617' : '#f5f7fb',
    colorBgContainer: uiStore.isDark ? '#111827' : '#ffffff',
  },
}))

// 监听 logo 配置变化
watch(
  () => configStore.get('sys_logo'),
  (newLogo) => updateFavicon(newLogo),
  { immediate: true }
)

watch(
  () => uiStore.theme.primaryColor,
  (color) => updateThemeVariables(color),
  { immediate: true }
)

// 确保配置加载后更新 favicon
onMounted(async () => {
  await configStore.loadConfigs()
  updateFavicon(configStore.get('sys_logo'))
})
</script>

<style>
body {
  margin: 0;
  padding: 0;
  background: #f5f7fb;
}

#app {
  min-height: 100vh;
}
</style>
