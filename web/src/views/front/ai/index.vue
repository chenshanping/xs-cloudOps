<template>
  <AIChat front-style class="ai-fullwidth" />
</template>

<script setup lang="ts">
import AIChat from '@/components/AIChat.vue'
</script>

<style scoped>
.ai-fullwidth {
  width: calc(100vw - 48px);
  margin-left: calc((100% - 100vw + 48px) / 2);
}
</style>
