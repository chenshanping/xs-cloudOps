<template>
  <a-modal
    v-model:open="visible"
    :title="title"
    :width="800"
    :footer="null"
  >
    <div class="markdown-viewer">
      <div v-if="!content" class="empty-content">
        <FileTextOutlined class="empty-icon" />
        <div>暂无内容</div>
      </div>
      <div v-else class="markdown-body" v-html="renderedContent"></div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { FileTextOutlined } from '@ant-design/icons-vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'

// 配置 marked 和 highlight.js
const renderer = new marked.Renderer()

// 自定义代码块渲染
renderer.code = function(token: { text?: string; lang?: string } | string) {
  const code = typeof token === 'string' ? token : (token.text || '')
  const lang = (typeof token === 'string' ? '' : (token.lang || '')).trim().toLowerCase()
  
  const language = lang && hljs.getLanguage(lang) ? lang : 'plaintext'
  const highlighted = hljs.highlight(code, { language }).value
  const langLabel = lang || 'code'
  
  return `<div class="code-block">
    <div class="code-header">
      <span class="code-lang">${langLabel}</span>
    </div>
    <pre><code class="hljs language-${language}">${highlighted}</code></pre>
  </div>`
}

marked.setOptions({
  renderer,
  breaks: true,
  gfm: true
})

interface Props {
  open: boolean
  title?: string
  content: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '查看详情'
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const visible = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val)
})

// 渲染 Markdown
function renderMarkdown(text: string): string {
  if (!text) return ''
  
  try {
    // 解码 HTML 实体
    const textarea = document.createElement('textarea')
    textarea.innerHTML = text
    const decodedContent = textarea.value
    
    // 调试输出
    if (process.env.NODE_ENV === 'development') {
      console.log('=== Markdown Viewer Debug ===')
      console.log('原始内容:', text)
      console.log('解码后:', decodedContent)
    }
    
    // 使用 marked.parse 解析
    const parsed = marked.parse(decodedContent) as string
    
    if (process.env.NODE_ENV === 'development') {
      console.log('解析后:', parsed)
      console.log('==========================')
    }
    
    return parsed
  } catch (error) {
    console.error('Markdown 解析失败:', error)
    return text.replace(/\n/g, '<br>')
  }
}

const renderedContent = computed(() => renderMarkdown(props.content))
</script>

<style scoped>
.markdown-viewer {
  max-height: 600px;
  overflow-y: auto;
  color: var(--app-text-color);
}

.empty-icon {
  font-size: 48px;
  color: var(--app-text-muted);
  margin-bottom: 12px;
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: var(--app-text-muted);
  font-size: 14px;
}

.markdown-body {
  padding: 16px;
  background: var(--app-surface-soft);
  border-radius: 8px;
  border: 1px solid var(--app-border-color);
  line-height: 1.8;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4),
.markdown-body :deep(h5),
.markdown-body :deep(h6) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.4;
  color: var(--app-text-strong);
}

.markdown-body :deep(h1:first-child),
.markdown-body :deep(h2:first-child),
.markdown-body :deep(h3:first-child),
.markdown-body :deep(h4:first-child),
.markdown-body :deep(h5:first-child),
.markdown-body :deep(h6:first-child) {
  margin-top: 0;
}

.markdown-body :deep(h1) {
  font-size: 28px;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--app-border-color);
}

.markdown-body :deep(h2) {
  font-size: 24px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--app-border-color);
}

.markdown-body :deep(h3) {
  font-size: 20px;
}

.markdown-body :deep(h4) {
  font-size: 18px;
}

.markdown-body :deep(h5) {
  font-size: 16px;
}

.markdown-body :deep(h6) {
  font-size: 14px;
}

.markdown-body :deep(p) {
  margin-bottom: 16px;
  color: var(--app-text-secondary);
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 28px;
  margin-bottom: 16px;
}

.markdown-body :deep(ul li),
.markdown-body :deep(ol li) {
  margin-bottom: 8px;
  line-height: 1.8;
  color: var(--app-text-secondary);
}

.markdown-body :deep(ul) {
  list-style-type: disc;
}

.markdown-body :deep(ol) {
  list-style-type: decimal;
}

.markdown-body :deep(blockquote) {
  border-left: 4px solid var(--app-primary-color);
  padding: 12px 16px;
  margin: 16px 0;
  background: var(--app-primary-color-soft);
  color: var(--app-text-secondary);
}

.markdown-body :deep(blockquote p) {
  margin-bottom: 0;
}

.markdown-body :deep(code) {
  background: var(--app-code-bg);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
  font-size: 13px;
  color: var(--app-text-strong);
}

.markdown-body :deep(pre) {
  background: #282c34;
  color: #abb2bf;
  padding: 16px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 16px 0;
}

.markdown-body :deep(pre code) {
  background: none;
  padding: 0;
  color: inherit;
}

.markdown-body :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 16px 0;
}

.markdown-body :deep(table th),
.markdown-body :deep(table td) {
  border: 1px solid var(--app-border-color);
  padding: 10px 12px;
  text-align: left;
}

.markdown-body :deep(table th) {
  background: var(--app-code-bg);
  font-weight: 600;
  color: var(--app-text-strong);
}

.markdown-body :deep(table tr:hover) {
  background: var(--app-hover-bg);
}

.markdown-body :deep(hr) {
  border: none;
  border-top: 1px solid var(--app-border-color);
  margin: 24px 0;
}

.markdown-body :deep(a) {
  color: var(--app-primary-color);
  text-decoration: none;
}

.markdown-body :deep(a:hover) {
  text-decoration: underline;
}

.markdown-body :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 6px;
  margin: 16px 0;
}

.markdown-body :deep(strong) {
  font-weight: 600;
  color: var(--app-text-strong);
}

.markdown-body :deep(em) {
  font-style: italic;
}

/* 代码块样式 */
.markdown-body :deep(.code-block) {
  margin: 16px 0;
  border-radius: 6px;
  overflow: hidden;
  background: #282c34;
}

.markdown-body :deep(.code-header) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #21252b;
  border-bottom: 1px solid #181a1f;
}

.markdown-body :deep(.code-lang) {
  font-size: 12px;
  color: #abb2bf;
  font-weight: 500;
}

.markdown-body :deep(.code-block pre) {
  margin: 0;
  padding: 16px;
  background: #282c34;
  overflow-x: auto;
}

.markdown-body :deep(.code-block code) {
  background: none;
  padding: 0;
  color: #abb2bf;
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
}
</style>
