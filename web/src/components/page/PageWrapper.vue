<template>
  <section class="page-wrapper">
    <div
      v-if="hasHeader"
      :class="['page-wrapper__header', { 'page-wrapper__header--bordered': headerBordered }]"
    >
      <div class="page-wrapper__header-main">
        <slot name="header">
          <div v-if="title" class="page-wrapper__title">{{ title }}</div>
          <div v-if="description" class="page-wrapper__description">{{ description }}</div>
        </slot>
      </div>
      <div v-if="$slots.actions" class="page-wrapper__actions">
        <slot name="actions" />
      </div>
    </div>

    <div :class="['page-wrapper__content', contentClass]" :style="contentStyle">
      <slot />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'

const props = withDefaults(defineProps<{
  title?: string
  description?: string
  contentClass?: string
  gap?: number
  headerBordered?: boolean
}>(), {
  title: '',
  description: '',
  contentClass: '',
  gap: 16,
  headerBordered: false,
})

const slots = useSlots()

const hasHeader = computed(() => Boolean(props.title || props.description || slots.header || slots.actions))
const contentStyle = computed(() => ({
  '--page-wrapper-gap': `${props.gap}px`,
}))
</script>

<style scoped>
.page-wrapper {
  min-width: 0;
}

.page-wrapper__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.page-wrapper__header--bordered {
  padding-bottom: 12px;
  border-bottom: 1px solid var(--ant-color-border-secondary, #f0f0f0);
}

.page-wrapper__header-main {
  min-width: 0;
}

.page-wrapper__title {
  color: var(--ant-color-text-heading, #262626);
  font-size: 18px;
  font-weight: 600;
  line-height: 1.4;
}

.page-wrapper__description {
  margin-top: 4px;
  color: var(--ant-color-text-description, #8c8c8c);
  font-size: 13px;
  line-height: 1.6;
}

.page-wrapper__actions {
  flex-shrink: 0;
}

.page-wrapper__content {
  display: flex;
  flex-direction: column;
  gap: var(--page-wrapper-gap);
  min-width: 0;
}

@media (max-width: 768px) {
  .page-wrapper__header {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
