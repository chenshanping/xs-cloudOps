<template>
  <a-card :bordered="false" class="search-card-shell" :body-style="bodyStyle">
    <div class="search-card-shell__body">
      <slot />
    </div>
  </a-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  bodyPadding?: string
}>(), {
  bodyPadding: '16px 20px 0',
})

const bodyStyle = computed(() => ({
  padding: props.bodyPadding,
}))
</script>

<style scoped>
.search-card-shell {
  border: 1px solid var(--ant-color-border-secondary, #f0f0f0);
  border-radius: var(--app-radius-md, 10px);
  background: var(--ant-color-bg-container, #ffffff);
  box-shadow: var(--app-card-shadow, 0 4px 16px rgb(15 23 42 / 4%));
}

.search-card-shell__body {
  min-width: 0;
}
</style>
