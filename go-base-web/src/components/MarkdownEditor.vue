<template>
  <div class="markdown-editor">
    <div class="editor-toolbar">
      <a-space>
        <a-tooltip title="粗体">
          <a-button size="small" @click="insertFormat('**', '粗体文字')">
            <template #icon><BoldOutlined /></template>
          </a-button>
        </a-tooltip>
        <a-tooltip title="斜体">
          <a-button size="small" @click="insertFormat('*', '斜体文字')">
            <template #icon><ItalicOutlined /></template>
          </a-button>
        </a-tooltip>
        <a-divider type="vertical" />
        <a-tooltip title="标题">
          <a-dropdown>
            <a-button size="small">
              H <DownOutlined />
            </a-button>
            <template #overlay>
              <a-menu @click="handleHeadingClick">
                <a-menu-item key="1">H1</a-menu-item>
                <a-menu-item key="2">H2</a-menu-item>
                <a-menu-item key="3">H3</a-menu-item>
                <a-menu-item key="4">H4</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </a-tooltip>
        <a-divider type="vertical" />
        <a-tooltip title="无序列表">
          <a-button size="small" @click="insertList('- ')">
            <template #icon><UnorderedListOutlined /></template>
          </a-button>
        </a-tooltip>
        <a-tooltip title="有序列表">
          <a-button size="small" @click="insertList('1. ')">
            <template #icon><OrderedListOutlined /></template>
          </a-button>
        </a-tooltip>
        <a-divider type="vertical" />
        <a-tooltip title="引用">
          <a-button size="small" @click="insertQuote">
            <template #icon><MessageOutlined /></template>
          </a-button>
        </a-tooltip>
        <a-tooltip title="代码">
          <a-button size="small" @click="insertCode">
            <template #icon><CodeOutlined /></template>
          </a-button>
        </a-tooltip>
        <a-divider type="vertical" />
        <a-tooltip title="预览">
          <a-button size="small" :type="showPreview ? 'primary' : 'default'" @click="togglePreview">
            <template #icon><EyeOutlined /></template>
          </a-button>
        </a-tooltip>
      </a-space>
    </div>
    
    <div class="editor-content" :class="{ 'split-view': showPreview }">
      <div class="editor-pane">
        <a-textarea
          ref="textareaRef"
          v-model:value="localContent"
          :rows="rows"
          :placeholder="placeholder"
          @change="handleChange"
          class="editor-textarea"
        />
      </div>
      
      <div v-if="showPreview" class="preview-pane">
        <div class="markdown-body" v-html="renderedContent"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { 
  BoldOutlined, 
  ItalicOutlined, 
  UnorderedListOutlined, 
  OrderedListOutlined,
  MessageOutlined,
  CodeOutlined,
  EyeOutlined,
  DownOutlined
} from '@ant-design/icons-vue'
import { marked } from 'marked'
import type { TextAreaRef } from 'ant-design-vue/es/input/TextArea'

// 配置 marked
marked.setOptions({
  breaks: true,
  gfm: true
})

interface Props {
  modelValue: string
  rows?: number
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  rows: 12,
  placeholder: '请输入内容，支持 Markdown 格式'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const localContent = ref(props.modelValue)
const showPreview = ref(false)
const textareaRef = ref<TextAreaRef>()

watch(() => props.modelValue, (val) => {
  localContent.value = val
})

const handleChange = () => {
  emit('update:modelValue', localContent.value)
}

const renderedContent = computed(() => {
  if (!localContent.value) return '<div style="color: #999; text-align: center; padding: 40px;">暂无预览内容</div>'
  try {
    return marked.parse(localContent.value) as string
  } catch {
    return localContent.value
  }
})

// 获取光标位置
const getCursorPosition = (): { start: number; end: number } => {
  const textarea = textareaRef.value?.$el?.querySelector('textarea')
  if (!textarea) return { start: 0, end: 0 }
  return {
    start: textarea.selectionStart || 0,
    end: textarea.selectionEnd || 0
  }
}

// 设置光标位置
const setCursorPosition = (start: number, end?: number) => {
  const textarea = textareaRef.value?.$el?.querySelector('textarea')
  if (!textarea) return
  textarea.focus()
  textarea.setSelectionRange(start, end || start)
}

// 插入格式
const insertFormat = (format: string, placeholder: string) => {
  const { start, end } = getCursorPosition()
  const selectedText = localContent.value.substring(start, end) || placeholder
  const beforeText = localContent.value.substring(0, start)
  const afterText = localContent.value.substring(end)
  
  localContent.value = `${beforeText}${format}${selectedText}${format}${afterText}`
  handleChange()
  
  // 设置新的光标位置
  setTimeout(() => {
    const newPosition = start + format.length
    setCursorPosition(newPosition, newPosition + selectedText.length)
  }, 0)
}

// 插入标题
const handleHeadingClick = ({ key }: { key: string }) => {
  const { start } = getCursorPosition()
  const beforeText = localContent.value.substring(0, start)
  const afterText = localContent.value.substring(start)
  const headingPrefix = '#'.repeat(parseInt(key)) + ' '
  
  localContent.value = `${beforeText}${headingPrefix}标题文字${afterText}`
  handleChange()
  
  setTimeout(() => {
    setCursorPosition(start + headingPrefix.length, start + headingPrefix.length + 4)
  }, 0)
}

// 插入列表
const insertList = (prefix: string) => {
  const { start } = getCursorPosition()
  const beforeText = localContent.value.substring(0, start)
  const afterText = localContent.value.substring(start)
  
  localContent.value = `${beforeText}${prefix}列表项${afterText}`
  handleChange()
  
  setTimeout(() => {
    setCursorPosition(start + prefix.length, start + prefix.length + 3)
  }, 0)
}

// 插入引用
const insertQuote = () => {
  const { start, end } = getCursorPosition()
  const selectedText = localContent.value.substring(start, end) || '引用内容'
  const beforeText = localContent.value.substring(0, start)
  const afterText = localContent.value.substring(end)
  
  localContent.value = `${beforeText}> ${selectedText}${afterText}`
  handleChange()
  
  setTimeout(() => {
    setCursorPosition(start + 2, start + 2 + selectedText.length)
  }, 0)
}

// 插入代码块
const insertCode = () => {
  const { start, end } = getCursorPosition()
  const selectedText = localContent.value.substring(start, end) || '代码'
  const beforeText = localContent.value.substring(0, start)
  const afterText = localContent.value.substring(end)
  
  if (selectedText.includes('\n')) {
    // 多行代码块
    localContent.value = `${beforeText}\`\`\`\n${selectedText}\n\`\`\`${afterText}`
    handleChange()
    setTimeout(() => {
      setCursorPosition(start + 4, start + 4 + selectedText.length)
    }, 0)
  } else {
    // 单行代码
    localContent.value = `${beforeText}\`${selectedText}\`${afterText}`
    handleChange()
    setTimeout(() => {
      setCursorPosition(start + 1, start + 1 + selectedText.length)
    }, 0)
  }
}

// 切换预览
const togglePreview = () => {
  showPreview.value = !showPreview.value
}
</script>

<style scoped lang="less">
.markdown-editor {
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  overflow: hidden;
  
  .editor-toolbar {
    padding: 8px 12px;
    background: #fafafa;
    border-bottom: 1px solid #e8e8e8;
  }
  
  .editor-content {
    display: flex;
    
    &.split-view {
      .editor-pane {
        width: 50%;
        border-right: 1px solid #e8e8e8;
      }
      
      .preview-pane {
        width: 50%;
      }
    }
    
    .editor-pane {
      width: 100%;
      
      .editor-textarea {
        border: none;
        border-radius: 0;
        resize: none;
        font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
        font-size: 14px;
        
        &:focus {
          box-shadow: none;
        }
      }
      
      :deep(.ant-input) {
        border: none;
        border-radius: 0;
        
        &:focus {
          box-shadow: none;
        }
      }
    }
    
    .preview-pane {
      padding: 12px;
      overflow-y: auto;
      background: #fff;
      max-height: 400px;
      
      .markdown-body {
        line-height: 1.8;
        
        :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
          margin-top: 20px;
          margin-bottom: 12px;
          font-weight: 600;
          line-height: 1.4;
          color: #333;
          
          &:first-child {
            margin-top: 0;
          }
        }
        
        :deep(h1) {
          font-size: 24px;
          padding-bottom: 10px;
          border-bottom: 2px solid #e8e8e8;
        }
        
        :deep(h2) {
          font-size: 20px;
          padding-bottom: 8px;
          border-bottom: 1px solid #e8e8e8;
        }
        
        :deep(h3) { font-size: 18px; }
        :deep(h4) { font-size: 16px; }
        :deep(h5) { font-size: 14px; }
        :deep(h6) { font-size: 13px; }
        
        :deep(p) {
          margin-bottom: 12px;
          color: #555;
        }
        
        :deep(ul), :deep(ol) {
          padding-left: 24px;
          margin-bottom: 12px;
          
          li {
            margin-bottom: 6px;
            line-height: 1.8;
            color: #555;
          }
        }
        
        :deep(ul) {
          list-style-type: disc;
        }
        
        :deep(ol) {
          list-style-type: decimal;
        }
        
        :deep(blockquote) {
          border-left: 4px solid #1890ff;
          padding: 10px 12px;
          margin: 12px 0;
          background: #f0f8ff;
          color: #666;
          
          p {
            margin-bottom: 0;
          }
        }
        
        :deep(code) {
          background: #f5f5f5;
          padding: 2px 6px;
          border-radius: 3px;
          font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
          font-size: 13px;
          color: #e83e8c;
        }
        
        :deep(pre) {
          background: #282c34;
          color: #abb2bf;
          padding: 12px;
          border-radius: 4px;
          overflow-x: auto;
          margin: 12px 0;
          
          code {
            background: none;
            padding: 0;
            color: inherit;
          }
        }
        
        :deep(strong) {
          font-weight: 600;
          color: #333;
        }
        
        :deep(em) {
          font-style: italic;
        }
      }
    }
  }
}
</style>
