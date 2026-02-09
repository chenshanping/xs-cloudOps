<template>
  <!-- 自定义 SVG（svg:xxx 格式） -->
  <span v-if="isSvg" class="svg-icon" v-html="svgContent"></span>
  <!-- Ant Design 图标 -->
  <component v-else-if="antIcon" :is="antIcon" />
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import * as AntIcons from '@ant-design/icons-vue'

const props = defineProps<{
  name?: string
}>()

// 判断是否为自定义 SVG（svg:xxx 格式）
const isSvg = computed(() => props.name?.startsWith('svg:'))
const svgName = computed(() => props.name?.replace('svg:', '') || '')

// 获取 Ant Design 图标组件
const antIcon = computed(() => {
  if (!props.name || isSvg.value) return null
  return (AntIcons as Record<string, Component>)[props.name] || null
})

// 预加载所有 SVG
const svgModules = import.meta.glob<string>('@/assets/icons/*.svg', { 
  eager: true,
  query: '?raw',
  import: 'default'
})

const svgContent = ref('')

const loadSvg = () => {
  if (!isSvg.value) {
    svgContent.value = ''
    return
  }
  
  let content = ''
  for (const [path, svg] of Object.entries(svgModules)) {
    if (path.endsWith(`/${svgName.value}.svg`)) {
      content = svg as string
      break
    }
  }
  
  if (content) {
    svgContent.value = content
      .replace(/width="[^"]*"/g, '')
      .replace(/height="[^"]*"/g, '')
      .replace(/<svg/, '<svg width="1em" height="1em"')
  } else {
    console.warn(`SVG not found: ${svgName.value}`)
    svgContent.value = ''
  }
}

watch(() => props.name, loadSvg, { immediate: true })
</script>

<style scoped>
.svg-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  vertical-align: -0.15em;
  line-height: 1;
}

.svg-icon :deep(svg) {
  width: 1em;
  height: 1em;
  fill: currentColor;
}
</style>
