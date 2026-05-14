<template>
  <div class="test-page">
    <a-card :bordered="false" class="test-page__panel">
      <div class="test-page__header">
        <div>
          <h2 class="test-page__title">测试页面</h2>
          <div class="test-page__subtitle">公共组件联调入口</div>
        </div>
        <a-button type="primary" @click="openRichTextDrawer">
          <template #icon>
            <EditOutlined />
          </template>
          打开富文本 Demo
        </a-button>
      </div>

      <a-divider />

      <a-empty v-if="!savedContent.content" description="暂无富文本内容" />
      <article v-else class="rich-preview">
        <h3 class="rich-preview__title">{{ savedContent.title }}</h3>
        <div class="rich-preview__body" v-html="safePreviewHtml" />
      </article>
    </a-card>

    <RichTextDemoDrawer
      v-model:open="richTextDrawerOpen"
      :initial-value="savedContent"
      @submit="handleRichTextSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import DOMPurify from 'dompurify'
import { EditOutlined } from '@ant-design/icons-vue'
import RichTextDemoDrawer from './components/RichTextDemoDrawer.vue'

const richTextDrawerOpen = ref(false)

const savedContent = reactive({
  title: '富文本示例',
  content: '',
})

const safePreviewHtml = computed(() => DOMPurify.sanitize(savedContent.content))

const openRichTextDrawer = () => {
  richTextDrawerOpen.value = true
}

const handleRichTextSubmit = (value: { title: string; content: string }) => {
  savedContent.title = value.title
  savedContent.content = value.content
}
</script>

<style scoped>
.test-page {
  padding: 24px;
}

.test-page__panel {
  border-radius: 8px;
}

.test-page__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.test-page__title {
  margin: 0;
  color: #262626;
  font-size: 18px;
  font-weight: 600;
  line-height: 1.4;
}

.test-page__subtitle {
  margin-top: 4px;
  color: #8c8c8c;
  font-size: 13px;
}

.rich-preview {
  min-height: 180px;
  padding: 16px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
}

.rich-preview__title {
  margin: 0 0 12px;
  color: #262626;
  font-size: 16px;
  font-weight: 600;
}

.rich-preview__body {
  color: #262626;
  font-size: 14px;
  line-height: 1.7;
}

.rich-preview__body :deep(p) {
  margin: 0 0 8px;
}

.rich-preview__body :deep(img),
.rich-preview__body :deep(video) {
  max-width: 100%;
  border-radius: 4px;
}

@media (max-width: 768px) {
  .test-page {
    padding: 12px;
  }

  .test-page__header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
