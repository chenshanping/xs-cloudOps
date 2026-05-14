<template>
  <section class="admin-split-layout" :style="layoutStyle">
    <aside class="admin-split-layout__aside">
      <slot name="aside" />
    </aside>
    <div class="admin-split-layout__main">
      <slot />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  asideWidth?: number
  contentMinWidth?: number
  gap?: number
}>(), {
  asideWidth: 260,
  contentMinWidth: 900,
  gap: 16,
})

const layoutStyle = computed(() => ({
  '--admin-split-aside-width': `${props.asideWidth}px`,
  '--admin-split-content-min-width': `${props.contentMinWidth}px`,
  '--admin-split-gap': `${props.gap}px`,
}))
</script>

<style scoped>
.admin-split-layout {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: var(--admin-split-gap);
  width: 100%;
  min-width: 0;
}

.admin-split-layout__aside {
  flex: 0 0 auto;
  width: min(100%, var(--admin-split-aside-width));
}

.admin-split-layout__main {
  flex: 1 1 var(--admin-split-content-min-width);
  min-width: min(100%, var(--admin-split-content-min-width));
}
</style>
