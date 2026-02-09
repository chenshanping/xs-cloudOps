<template>
  <a-config-provider :locale="zhCN">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { useConfigStore } from '@/store/config'
import { watch, onMounted } from 'vue'

const configStore = useConfigStore()

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

// 监听 logo 配置变化
watch(
  () => configStore.get('sys_logo'),
  (newLogo) => updateFavicon(newLogo),
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
}
</style>
