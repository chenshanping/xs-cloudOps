<template>
  <a-card :bordered="false" class="table-card-shell" :body-style="bodyStyle">
    <div v-if="hasHeader" class="table-card-shell__header">
      <div class="table-card-shell__main">
        <slot name="header">
          <div v-if="title" class="table-card-shell__title">{{ title }}</div>
          <div v-if="description" class="table-card-shell__description">{{ description }}</div>
        </slot>
      </div>
      <div v-if="$slots.toolbar" class="table-card-shell__toolbar">
        <slot name="toolbar" />
      </div>
    </div>

    <div class="table-card-shell__content">
      <slot />
    </div>
  </a-card>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'

const props = withDefaults(defineProps<{
  title?: string
  description?: string
  bodyPadding?: string
}>(), {
  title: '',
  description: '',
  bodyPadding: '16px 20px',
})

const slots = useSlots()

const hasHeader = computed(() => Boolean(props.title || props.description || slots.header || slots.toolbar))
const bodyStyle = computed(() => ({
  padding: props.bodyPadding,
}))
</script>

<style scoped>
.table-card-shell {
  border: 1px solid var(--ant-color-border-secondary, #f0f0f0);
  border-radius: var(--app-radius-md, 10px);
  background: var(--ant-color-bg-container, #ffffff);
  box-shadow: var(--app-card-shadow, 0 4px 16px rgb(15 23 42 / 4%));
}

.table-card-shell__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.table-card-shell__main {
  min-width: 0;
}

.table-card-shell__title {
  color: var(--ant-color-text-heading, #262626);
  font-size: 16px;
  font-weight: 600;
  line-height: 1.4;
}

.table-card-shell__description {
  margin-top: 4px;
  color: var(--ant-color-text-description, #8c8c8c);
  font-size: 13px;
  line-height: 1.5;
}

.table-card-shell__toolbar {
  flex-shrink: 0;
}

.table-card-shell__content {
  min-width: 0;
}

@media (max-width: 768px) {
  .table-card-shell__header {
    flex-direction: column;
    align-items: stretch;
  }

  .table-card-shell__toolbar {
    width: 100%;
  }
}
</style>
