<template>
  <div class="drawer-context-header">
    <div v-if="eyebrow" class="drawer-context-header__eyebrow">{{ eyebrow }}</div>
    <div class="drawer-context-header__title-row">
      <span class="drawer-context-header__title">{{ title }}</span>
      <a-tag v-if="roleName" color="blue" class="drawer-context-header__tag">{{ roleName }}</a-tag>
    </div>
    <div v-if="subtitle" class="drawer-context-header__subtitle">{{ subtitle }}</div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  eyebrow?: string
  title: string
  subtitle?: string
  roleName?: string
}

defineProps<Props>()
</script>

<style scoped>
.drawer-context-header__eyebrow {
  font-size: 12px;
  font-weight: 500;
  color: inherit;
  opacity: 0.68;
}

.drawer-context-header__title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
}

.drawer-context-header__title {
  color: inherit;
  font-size: 16px;
  font-weight: 650;
  line-height: 1.4;
}

.drawer-context-header__tag {
  margin-inline-end: 0;
  font-size: 12px;
}

.drawer-context-header__subtitle {
  margin-top: 6px;
  color: inherit;
  opacity: 0.6;
  font-size: 12px;
  line-height: 1.5;
}
</style>
