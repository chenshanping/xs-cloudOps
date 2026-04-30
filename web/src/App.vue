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

const applyThemeVariables = () => {
  const root = document.documentElement
  const primaryColor = uiStore.theme.primaryColor
  const isDark = uiStore.isDark

  const palette = isDark
    ? {
      '--app-layout-bg': '#0b1220',
      '--app-layout-radial': 'rgba(56, 189, 248, 0.12)',
      '--app-surface-color': '#111b2d',
      '--app-surface-soft': '#162235',
      '--app-elevated-bg': '#1a2638',
      '--app-border-color': 'rgba(148, 163, 184, 0.18)',
      '--app-border-strong': 'rgba(148, 163, 184, 0.28)',
      '--app-hover-bg': 'rgba(255, 255, 255, 0.08)',
      '--app-code-bg': 'rgba(148, 163, 184, 0.16)',
      '--app-text-color': 'rgba(255, 255, 255, 0.88)',
      '--app-text-secondary': 'rgba(255, 255, 255, 0.72)',
      '--app-text-muted': 'rgba(255, 255, 255, 0.56)',
      '--app-text-strong': '#ffffff',
      '--app-card-shadow': '0 0 0 1px rgba(148, 163, 184, 0.08)',
    }
    : {
      '--app-layout-bg': '#f5f7fb',
      '--app-layout-radial': 'rgba(22, 119, 255, 0.14)',
      '--app-surface-color': '#ffffff',
      '--app-surface-soft': '#f8fafc',
      '--app-elevated-bg': '#ffffff',
      '--app-border-color': '#e5e7eb',
      '--app-border-strong': '#d1d5db',
      '--app-hover-bg': 'rgba(15, 23, 42, 0.06)',
      '--app-code-bg': '#f3f4f6',
      '--app-text-color': '#111827',
      '--app-text-secondary': '#4b5563',
      '--app-text-muted': '#6b7280',
      '--app-text-strong': '#0f172a',
      '--app-card-shadow': '0 12px 32px rgba(15, 23, 42, 0.06)',
    }

  root.dataset.theme = isDark ? 'dark' : 'light'
  root.style.colorScheme = isDark ? 'dark' : 'light'
  root.style.setProperty('--app-primary-color', primaryColor)
  root.style.setProperty('--app-primary-color-soft', hexToRgba(primaryColor, isDark ? 0.22 : 0.12))
  Object.entries(palette).forEach(([name, value]) => {
    root.style.setProperty(name, value)
  })

  document.body.style.background = palette['--app-layout-bg']
  document.body.style.color = palette['--app-text-color']
}

const antdTheme = computed(() => ({
  cssVar: true,
  algorithm: uiStore.isDark
    ? antdThemeFactory.darkAlgorithm
    : antdThemeFactory.defaultAlgorithm,
  token: {
    colorPrimary: uiStore.theme.primaryColor,
    colorInfo: uiStore.theme.primaryColor,
    colorLink: uiStore.theme.primaryColor,
    borderRadius: 10,
    colorBgLayout: uiStore.isDark ? '#0b1220' : '#f5f7fb',
    colorBgContainer: uiStore.isDark ? '#111b2d' : '#ffffff',
    colorBgElevated: uiStore.isDark ? '#1a2638' : '#ffffff',
    colorTextBase: uiStore.isDark ? 'rgba(255, 255, 255, 0.88)' : '#111827',
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
  () => applyThemeVariables(),
  { immediate: true }
)

watch(
  () => uiStore.isDark,
  () => applyThemeVariables(),
  { immediate: true }
)

// 确保配置加载后更新 favicon
onMounted(async () => {
  await configStore.loadConfigs()
  updateFavicon(configStore.get('sys_logo'))
})
</script>

<style>
html,
body {
  margin: 0;
  padding: 0;
  background: var(--app-layout-bg, #f5f7fb);
  color: var(--app-text-color, #111827);
}

#app {
  min-height: 100vh;
}
</style>
