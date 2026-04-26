<template>
  <component v-if="antIcon" :is="antIcon" />
  <img v-else-if="customIconUrl" :src="customIconUrl" class="menu-icon" alt="" />
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'
import * as AntIcons from '@ant-design/icons-vue'

const props = defineProps<{
  icon?: string
}>()

const iconModules = import.meta.glob('@/assets/icons/*.svg', {
  eager: true,
  query: '?url',
  import: 'default',
})

const customIconUrl = computed(() => {
  if (!props.icon?.startsWith('custom-')) {
    return ''
  }
  const name = props.icon.replace('custom-', '')
  const key = `/src/assets/icons/${name}.svg`
  return (iconModules[key] as string) || ''
})

const antIcon = computed<Component | null>(() => {
  if (!props.icon || props.icon.startsWith('custom-')) {
    return null
  }

  const iconName = props.icon.startsWith('official-')
    ? props.icon.replace('official-', '')
    : props.icon

  return ((AntIcons as Record<string, Component>)[iconName] || null) as Component | null
})
</script>

<style scoped>
.menu-icon {
  width: 1em;
  height: 1em;
  object-fit: contain;
}
</style>
