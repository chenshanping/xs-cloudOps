<template>
  <div class="ai-chat-container" :class="{ 'front-style': frontStyle }">
    <!-- 左侧对话列表 -->
    <div class="conversation-sidebar">
      <div class="sidebar-header">
        <a-button type="primary" block @click="handleNewChat">
          <template #icon><PlusOutlined /></template>
          新对话
        </a-button>
        <a-input
          v-model:value="searchKeyword"
          placeholder="搜索对话"
          allowClear
          size="small"
          style="margin-top: 8px"
        >
          <template #prefix><SearchOutlined /></template>
        </a-input>
      </div>
      
      <div class="conversation-list">
        <div
          v-for="conv in filteredConversations"
          :key="conv.id"
          :class="['conversation-item', { active: aiStore.currentConversation?.id === conv.id }]"
          @click="handleSelectConversation(conv.id)"
        >
          <MessageOutlined v-if="frontStyle" class="conv-icon" />
          <!-- 编辑模式 -->
          <a-input
            v-if="editingConvId === conv.id"
            v-model:value="editingTitle"
            size="small"
            class="conv-title-input"
            @click.stop
            @pressEnter="saveConvTitle(conv.id)"
            @blur="saveConvTitle(conv.id)"
            ref="titleInputRef"
          />
          <span v-else class="conv-title" @dblclick.stop="startEditTitle(conv)">
            {{ conv.title }}
          </span>
          <div class="conv-actions">
            <EditOutlined class="edit-icon" @click.stop="startEditTitle(conv)" />
            <a-popconfirm
              title="确定删除这个对话吗？"
              @confirm="handleDeleteConversation(conv.id)"
              placement="right"
            >
              <DeleteOutlined class="delete-icon" @click.stop />
            </a-popconfirm>
          </div>
        </div>
        
        <a-empty v-if="filteredConversations.length === 0" :description="searchKeyword ? '无匹配对话' : '暂无对话'" :image="frontStyle ? Empty.PRESENTED_IMAGE_SIMPLE : undefined" />
      </div>
    </div>
    
    <!-- 右侧聊天区域 -->
    <div class="chat-main">
      <!-- 顶部工具栏 -->
      <div class="chat-header">
        <div class="header-left">
          <a-select
            v-model:value="aiStore.currentModel"
            style="width: 300px"
            @change="handleModelChange"
          >
            <a-select-option v-for="model in aiStore.models" :key="model.id" :value="model.id">
              {{ model.name }}
            </a-select-option>
          </a-select>
        </div>
        <div class="header-right">
          <a-space>
            <a-tooltip title="联网搜索">
              <a-switch
                v-model:checked="aiStore.enableSearch"
                checked-children="联网"
                un-checked-children="关闭"
                size="small"
              />
            </a-tooltip>
            <a-tooltip title="思考模式">
              <a-switch
                v-model:checked="aiStore.enableThinking"
                checked-children="思考"
                un-checked-children="关闭"
                size="small"
              />
            </a-tooltip>
          </a-space>
        </div>
      </div>
      
      
      <!-- 消息列表 -->
      <div class="message-list" ref="messageListRef">
        <template v-if="aiStore.messages.length > 0 || aiStore.streaming">
          <div
            v-for="(msg, idx) in aiStore.messages"
            :key="msg.id"
            :class="['message-item', msg.role]"
          >
            <div class="message-avatar">
              <a-avatar 
                v-if="msg.role === 'user'" 
                :src="frontStyle ? userStore.user?.avatar_file_url : undefined"
                :style="{ backgroundColor: '#1890ff' }"
              >
                <template #icon><UserOutlined /></template>
              </a-avatar>
              <a-avatar v-else :style="{ backgroundColor: '#52c41a' }">
                <template #icon><RobotOutlined /></template>
              </a-avatar>
            </div>
            
            <div class="message-content">
              <!-- 思考过程 -->
              <div v-if="msg.reasoning_content" class="reasoning-panel">
                <div class="reasoning-header" @click="toggleReasoning(msg.id)">
                  <ThunderboltOutlined />
                  <span>思考过程</span>
                  <DownOutlined :class="{ rotated: !expandedReasoning[msg.id] }" />
                </div>
                <div v-show="expandedReasoning[msg.id]" class="reasoning-content">
                  {{ msg.reasoning_content }}
                </div>
              </div>
              
              <!-- 消息正文 -->
              <div class="message-text" v-html="renderMarkdown(msg.content)"></div>
              
              <!-- AI消息操作按钮 -->
              <div v-if="msg.role === 'assistant' && msg.content" class="message-actions">
                <a-tooltip title="复制">
                  <a-button type="text" size="small" @click="handleCopyMessage(msg.content)">
                    <template #icon><CopyOutlined /></template>
                  </a-button>
                </a-tooltip>
                <a-tooltip v-if="idx === aiStore.messages.length - 1 && !aiStore.streaming" title="重新生成">
                  <a-button type="text" size="small" @click="handleRegenerate">
                    <template #icon><ReloadOutlined /></template>
                  </a-button>
                </a-tooltip>
              </div>
            </div>
          </div>
          
          <!-- 流式输出占位消息 -->
          <div v-if="aiStore.streaming" class="message-item assistant">
            <div class="message-avatar">
              <a-avatar :style="{ backgroundColor: '#52c41a' }">
                <template #icon><RobotOutlined /></template>
              </a-avatar>
            </div>
            <div class="message-content">
              <div v-if="aiStore.streamingReasoning" class="reasoning-panel">
                <div class="reasoning-header">
                  <ThunderboltOutlined />
                  <span>思考中...</span>
                </div>
                <div class="reasoning-content">{{ aiStore.streamingReasoning }}</div>
              </div>
              <div class="message-text" v-if="aiStore.streamingContent" v-html="renderMarkdown(aiStore.streamingContent)"></div>
              <div v-else class="typing-indicator">
                <span></span><span></span><span></span>
              </div>
            </div>
          </div>
        </template>
        
        <!-- 空状态 -->
        <div v-else class="empty-chat">
          <RobotOutlined class="empty-icon" />
          <h3 v-if="frontStyle">开始您的 AI 对话</h3>
          <p>{{ frontStyle ? '输入任何问题，AI 助手将为您解答' : '开始一个新的对话吧' }}</p>
        </div>
      </div>
      
      <!-- 输入区域 -->
      <div class="chat-input-area">
        <!-- 附件列表 -->
        <div v-if="attachments.length > 0" class="attachment-list">
          <div v-for="(att, i) in attachments" :key="i" class="attachment-item">
            <LoadingOutlined v-if="att.uploading" class="att-icon loading" />
            <PaperClipOutlined v-else class="att-icon" />
            <span class="att-name" :title="att.file_name">{{ att.file_name }}</span>
            <span class="att-size">{{ formatSize(att.file_size) }}</span>
            <a-progress v-if="att.uploading" :percent="att.progress" size="small" :show-info="false" style="width: 60px" />
            <CloseOutlined class="att-remove" @click="removeAttachment(i)" />
          </div>
        </div>
        <div class="input-wrapper">
          <a-tooltip title="上传文件">
            <a-button type="text" class="attach-btn" @click="triggerFileInput" :disabled="aiStore.streaming || attachments.length >= 5">
              <template #icon><PaperClipOutlined /></template>
            </a-button>
          </a-tooltip>
          <input
            ref="fileInputRef"
            type="file"
            multiple
            :accept="acceptFileTypes"
            style="display: none"
            @change="handleFileSelect"
          />
          <a-textarea
            v-model:value="inputMessage"
            placeholder="输入消息，按 Enter 发送，Shift + Enter 换行"
            :auto-size="{ minRows: 1, maxRows: 6 }"
            @keydown="handleKeydown"
            :disabled="aiStore.streaming"
          />
          <div class="input-actions">
            <a-popconfirm
              v-if="aiStore.currentConversation && aiStore.messages.length > 0"
              title="清空上下文后，AI将不记得之前的对话内容"
              @confirm="handleClearContext"
              placement="top"
            >
              <a-button :disabled="aiStore.streaming">
                <template #icon><ClearOutlined /></template>
                <span v-if="!frontStyle">清空上下文</span>
              </a-button>
            </a-popconfirm>

            <a-button
              v-if="aiStore.streaming"
              type="primary"
              danger
              @click="aiStore.stopStreaming"
            >
              <template #icon><PauseCircleOutlined /></template>
              停止
            </a-button>
            <a-button
              v-else
              type="primary"
              @click="handleSend"
              :loading="aiStore.streaming"
              :disabled="!inputMessage.trim()"
            >
              <template #icon><SendOutlined /></template>
              发送
            </a-button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Mermaid 预览抽屉 -->
    <a-drawer
      v-model:open="mermaidPreviewVisible"
      title="📊 Mermaid 图表预览"
      placement="bottom"
      :height="'90vh'"
    >
      <a-spin :spinning="mermaidLoading" tip="正在渲染图表...">
        <div class="mermaid-preview-content">
          <div class="mermaid" v-html="mermaidSvg"></div>
        </div>
      </a-spin>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, reactive, ref, computed, watch } from 'vue'
import { Empty, message } from 'ant-design-vue'
import {
  ClearOutlined,
  CloseOutlined,
  CopyOutlined,
  DeleteOutlined,
  DownOutlined,
  EditOutlined,
  LoadingOutlined,
  MessageOutlined,
  PaperClipOutlined,
  PauseCircleOutlined,
  PlusOutlined,
  ReloadOutlined,
  RobotOutlined,
  SearchOutlined,
  SendOutlined,
  ThunderboltOutlined,
  UserOutlined
} from '@ant-design/icons-vue'
import { useAIStore } from '@/store/ai'
import { useUserStore } from '@/store/user'
import { updateConversationTitle } from '@/api/ai'
import { calculateMD5, multipartUpload, formatFileSize } from '@/utils/upload'
import type { FileInfo } from '@/types/file'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'
import DOMPurify from 'dompurify'
import mermaid from 'mermaid'

// Props
interface Props {
  frontStyle?: boolean  // 是否使用前台样式
}

const props = withDefaults(defineProps<Props>(), {
  frontStyle: false
})

// 配置 Mermaid
mermaid.initialize({
  startOnLoad: false,
  theme: 'default',
  securityLevel: 'loose',
  flowchart: {
    useMaxWidth: false,
    htmlLabels: true,
    curve: 'basis'
  },
  sequence: {
    useMaxWidth: false
  }
})

// Mermaid 预览抽屉相关
const mermaidPreviewVisible = ref(false)
const mermaidPreviewContent = ref('')
const mermaidSvg = ref('')
const mermaidLoading = ref(false)
let mermaidCounter = 0

// 配置 marked 和 highlight.js
const renderer = new marked.Renderer()

// 自定义代码块渲染
renderer.code = function(token: { text?: string; lang?: string } | string) {
  const code = typeof token === 'string' ? token : (token.text || '')
  const lang = (typeof token === 'string' ? '' : (token.lang || '')).trim().toLowerCase()
  
  // Mermaid 图表：生成预览按钮
  if (lang === 'mermaid') {
    const btnId = `mermaid-btn-${mermaidCounter++}`
    return `<div class="mermaid-preview-container">
      <button class="mermaid-preview-btn" id="${btnId}" data-mermaid-code="${encodeURIComponent(code)}">
        📊 预览 Mermaid 图表
      </button>
    </div>`
  }
  
  const language = lang && hljs.getLanguage(lang) ? lang : 'plaintext'
  const highlighted = hljs.highlight(code, { language }).value
  const langLabel = lang || 'code'
  
  return `<div class="code-block">
    <div class="code-header">
      <span class="code-lang">${langLabel}</span>
      <button class="copy-btn" onclick="navigator.clipboard.writeText(decodeURIComponent('${encodeURIComponent(code)}')).then(() => { this.textContent = '✓ 已复制'; setTimeout(() => this.textContent = '复制', 2000) })">复制</button>
    </div>
    <pre><code class="hljs language-${language}">${highlighted}</code></pre>
  </div>`
}

marked.setOptions({
  renderer,
  breaks: true,
  gfm: true
})

// 打开 Mermaid 预览抽屉
let mermaidRenderCounter = 0

// 预处理 Mermaid 代码，修复常见语法问题
function preprocessMermaidCode(code: string): string {
  let result = code
  
  // 处理 ([text]) stadium 形状 - 给文本添加引号
  result = result.replace(/\(\[([^\]"]+)\]\)/g, '(["$1"])')
  
  // 处理 ((text)) 圆形
  result = result.replace(/\(\(([^)"]+)\)\)/g, '(("$1"))')
  
  // 处理 {text} 菱形 - 包含特殊字符时添加引号
  result = result.replace(/\{([^}"]+)\}/g, (match, text) => {
    if (/[?/()\u4e00-\u9fa5]/.test(text)) {
      return `{"${text}"}`
    }
    return match
  })
  
  // 处理 [text] 矩形 - 包含特殊字符时添加引号
  result = result.replace(/\[([^\]"]+)\]/g, (match, text) => {
    if (/[?/()]/.test(text)) {
      return `["${text}"]`
    }
    return match
  })
  
  // 处理边标签 |text| - 包含特殊字符时添加引号
  result = result.replace(/\|([^|"]+)\|/g, (match, text) => {
    if (/[\u4e00-\u9fa5?/()]/.test(text)) {
      return `|"${text}"|`
    }
    return match
  })
  
  return result
}

async function openMermaidPreview(code: string) {
  mermaidPreviewContent.value = code
  mermaidLoading.value = true
  mermaidSvg.value = ''
  mermaidPreviewVisible.value = true
  
  await nextTick()
  await new Promise(resolve => setTimeout(resolve, 200))
  
  const id = `mermaid-svg-${mermaidRenderCounter++}`
  const processedCode = preprocessMermaidCode(code)
  
  try {
    const { svg } = await mermaid.render(id, processedCode)
    mermaidSvg.value = svg
        .replace(/width="[^"]*"/, 'width="100%"')
        .replace(/height="[^"]*"/, 'height="1200px"')
  } catch (error: any) {
    console.error('Mermaid 渲染失败:', error)
    mermaidSvg.value = `<div style="color: #ff4d4f; padding: 20px; text-align: left;">
      <p><strong>图表渲染失败</strong></p>
      <p style="font-size: 12px; color: #666;">${error.message || '请检查 Mermaid 语法'}</p>
      <details style="margin-top: 12px;">
        <summary style="cursor: pointer; color: #1890ff;">查看原始代码</summary>
        <pre style="background: #f5f5f5; padding: 12px; border-radius: 4px; overflow: auto; font-size: 12px; margin-top: 8px; white-space: pre-wrap;">${code}</pre>
      </details>
    </div>`
  } finally {
    mermaidLoading.value = false
    const tempElement = document.getElementById('d' + id)
    if (tempElement) {
      tempElement.remove()
    }
  }
}

// 使用事件委托处理 Mermaid 按钮点击
function handleMermaidClick(e: MouseEvent) {
  const target = e.target as HTMLElement
  const btn = target.closest('.mermaid-preview-btn') as HTMLElement
  if (btn) {
    const code = decodeURIComponent(btn.dataset.mermaidCode || '')
    if (code) {
      openMermaidPreview(code)
    }
  }
}

const aiStore = useAIStore()
const userStore = useUserStore()
const inputMessage = ref('')
const messageListRef = ref<HTMLElement | null>(null)
const expandedReasoning = reactive<Record<number, boolean>>({})

// 侧栏搜索
const searchKeyword = ref('')
const filteredConversations = computed(() => {
  if (!searchKeyword.value.trim()) return aiStore.conversations
  const kw = searchKeyword.value.toLowerCase()
  return aiStore.conversations.filter(c => c.title.toLowerCase().includes(kw))
})

// 对话标题编辑
const editingConvId = ref<number | null>(null)
const editingTitle = ref('')
const titleInputRef = ref<any>(null)

function startEditTitle(conv: any) {
  editingConvId.value = conv.id
  editingTitle.value = conv.title
  nextTick(() => {
    // focus input
    const inputs = document.querySelectorAll('.conv-title-input input')
    if (inputs.length > 0) (inputs[0] as HTMLInputElement).focus()
  })
}

async function saveConvTitle(id: number) {
  if (!editingTitle.value.trim()) {
    editingConvId.value = null
    return
  }
  try {
    await updateConversationTitle(id, editingTitle.value.trim())
    const conv = aiStore.conversations.find(c => c.id === id)
    if (conv) conv.title = editingTitle.value.trim()
    if (aiStore.currentConversation?.id === id) {
      aiStore.currentConversation.title = editingTitle.value.trim()
    }
  } catch {
    message.error('修改标题失败')
  }
  editingConvId.value = null
}

// === 文件上传 ===
interface Attachment {
  file_id: number
  file_url: string
  file_name: string
  file_ext: string
  file_size: number
  uploading: boolean
  progress: number
}

const attachments = ref<Attachment[]>([])
const fileInputRef = ref<HTMLInputElement | null>(null)
const acceptFileTypes = '.txt,.md,.csv,.json,.log,.go,.py,.js,.ts,.java,.html,.css,.sql,.xml,.yaml,.yml,.sh,.c,.cpp,.h,.rs,.rb,.php,.vue,.jsx,.tsx,.ini,.toml,.conf,.jpg,.jpeg,.png,.gif,.webp'

function formatSize(bytes: number) {
  return formatFileSize(bytes)
}

function triggerFileInput() {
  fileInputRef.value?.click()
}

async function handleFileSelect(e: Event) {
  const input = e.target as HTMLInputElement
  const files = Array.from(input.files || [])
  input.value = '' // reset
  
  const remaining = 5 - attachments.value.length
  if (files.length > remaining) {
    message.warning(`最多上传 5 个文件，还可添加 ${remaining} 个`)
  }
  
  for (const file of files.slice(0, remaining)) {
    const att: Attachment = {
      file_id: 0,
      file_url: '',
      file_name: file.name,
      file_ext: file.name.split('.').pop() || '',
      file_size: file.size,
      uploading: true,
      progress: 0
    }
    attachments.value.push(att)
    const idx = attachments.value.length - 1
    
    try {
      const md5 = await calculateMD5(file, (p) => {
        attachments.value[idx].progress = Math.round(p * 0.1)
      })
      const result: FileInfo = await multipartUpload(file, md5, undefined, (p) => {
        attachments.value[idx].progress = 10 + Math.round(p * 0.9)
      })
      console.log('[AIChat] upload result:', JSON.stringify(result))
      attachments.value[idx].file_id = result.id
      attachments.value[idx].file_url = result.url
      attachments.value[idx].uploading = false
      attachments.value[idx].progress = 100
      console.log('[AIChat] attachment after upload:', JSON.stringify(attachments.value[idx]))
    } catch (err: any) {
      message.error(`上传失败: ${file.name}`)
      attachments.value.splice(idx, 1)
    }
  }
}

function removeAttachment(index: number) {
  attachments.value.splice(index, 1)
}

// 复制消息
function handleCopyMessage(content: string) {
  navigator.clipboard.writeText(content).then(() => {
    message.success('已复制')
  }).catch(() => {
    message.error('复制失败')
  })
}

// 重新生成
function handleRegenerate() {
  aiStore.regenerateLastMessage()
}

// 初始化
onMounted(async () => {
  await aiStore.fetchModels()
  await aiStore.fetchConversations()
  
  if (messageListRef.value) {
    messageListRef.value.addEventListener('click', handleMermaidClick)
  }
})

// 组件卸载清理
onUnmounted(() => {
  if (messageListRef.value) {
    messageListRef.value.removeEventListener('click', handleMermaidClick)
  }
  aiStore.stopStreaming()
})

// 监听消息变化，自动滚动到底部
watch(
  () => aiStore.messages.length,
  () => {
    nextTick(() => scrollToBottom())
  }
)

// 监听流式内容变化
watch(
  () => aiStore.streamingContent,
  () => {
    nextTick(() => scrollToBottom())
  }
)

// 监听 Mermaid 抽屉关闭，重置状态
watch(mermaidPreviewVisible, (visible) => {
  if (!visible) {
    mermaidSvg.value = ''
    mermaidPreviewContent.value = ''
    mermaidLoading.value = false
    document.querySelectorAll('[id^="dmermaid-svg-"]').forEach(el => el.remove())
  }
})

// 滚动到底部
function scrollToBottom() {
  if (messageListRef.value) {
    messageListRef.value.scrollTop = messageListRef.value.scrollHeight
  }
}

// 新建对话
function handleNewChat() {
  aiStore.currentConversation = null
  aiStore.messages = []
}

// 选择对话
async function handleSelectConversation(id: number) {
  try {
    await aiStore.selectConversation(id)
  } catch (error) {
    message.error('加载对话失败')
  }
}

// 删除对话
async function handleDeleteConversation(id: number) {
  try {
    await aiStore.removeConversation(id)
    message.success('删除成功')
  } catch (error) {
    message.error('删除失败')
  }
}

// 模型变更
function handleModelChange(modelId: string) {
  aiStore.setModel(modelId)
}

// 发送消息
function handleSend() {
  if (!inputMessage.value.trim()) return
  // 检查是否有文件正在上传
  if (attachments.value.some(a => a.uploading)) {
    message.warning('请等待文件上传完成')
    return
  }
  const msg = inputMessage.value
  console.log('[AIChat] handleSend attachments:', JSON.stringify(attachments.value))
  const fileIds = attachments.value.filter(a => a.file_id > 0).map(a => a.file_id)
  console.log('[AIChat] handleSend fileIds:', fileIds)
  inputMessage.value = ''
  attachments.value = []
  nextTick(() => {
    aiStore.sendMessage(msg, fileIds.length > 0 ? fileIds : undefined)
  })
}

// 键盘事件
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

// 切换思考过程展开/收起
function toggleReasoning(msgId: number) {
  expandedReasoning[msgId] = !expandedReasoning[msgId]
}

// 清空上下文（保留聊天记录）
async function handleClearContext() {
  try {
    await aiStore.clearCurrentContext()
    message.success('上下文已清空，AI将不记得之前的对话')
  } catch (error) {
    message.error('清空失败')
  }
}

// Markdown渲染（带XSS防护）
function renderMarkdown(text: string): string {
  if (!text) return ''
  
  try {
    const html = marked.parse(text) as string
    return DOMPurify.sanitize(html, {
      ADD_TAGS: ['button'],
      ADD_ATTR: ['onclick', 'data-mermaid-code']
    })
  } catch (error) {
    console.error('Markdown 解析失败:', error)
    return DOMPurify.sanitize(text.replace(/\n/g, '<br>'))
  }
}
</script>

<style scoped lang="less">
.ai-chat-container {
  display: flex;
  height: calc(100vh - 120px);
  background: #f5f5f5;
  border-radius: 8px;
  overflow: hidden;
  
  // 前台样式变体
  &.front-style {
    height: calc(100vh - 160px);
    min-height: 500px;
    border-radius: 12px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    
    .conversation-sidebar {
      width: 260px;
      background: #fafafa;
    }
    
    .message-list {
      background: #f9f9f9;
      
      .message-item.user .message-content .message-text {
        background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
      }
      
      .message-item.assistant .message-content .message-text {
        background: #fff;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
      }
      
      .empty-chat {
        .empty-icon {
          font-size: 80px;
          color: #d9d9d9;
          margin-bottom: 24px;
        }
        
        h3 {
          font-size: 20px;
          color: #666;
          margin-bottom: 8px;
        }
      }
    }
  }
}

.conversation-sidebar {
  width: 280px;
  background: #fff;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  
  .sidebar-header {
    padding: 16px;
    border-bottom: 1px solid #e8e8e8;
  }
  
  .conversation-list {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
    
    .conversation-item {
      display: flex;
      align-items: center;
      padding: 12px;
      border-radius: 8px;
      cursor: pointer;
      margin-bottom: 4px;
      transition: all 0.2s;
      
      &:hover {
        background: #f5f5f5;
        
        .conv-actions {
          opacity: 1;
        }
      }
      
      &.active {
        background: #e6f7ff;
        border: 1px solid #91d5ff;
      }
      
      .conv-icon {
        color: #999;
        margin-right: 8px;
      }
      
      .conv-title {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-size: 14px;
      }
      
      .conv-title-input {
        flex: 1;
        min-width: 0;
      }
      
      .conv-actions {
        opacity: 0;
        transition: opacity 0.2s;
        display: flex;
        gap: 4px;
        
        .edit-icon, .delete-icon {
          color: #999;
          font-size: 12px;
        }
        
        .edit-icon:hover {
          color: #1890ff;
        }
        
        .delete-icon:hover {
          color: #ff4d4f;
        }
      }
    }
  }
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
  
  .chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #e8e8e8;
  }
  
  .message-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    
    .message-item {
      display: flex;
      margin-bottom: 24px;
      
      &.user {
        flex-direction: row-reverse;
        
        .message-content {
          margin-right: 12px;
          margin-left: 48px;
          
          .message-text {
            background: #1890ff;
            color: #fff;
          }
        }
      }
      
      &.assistant {
        .message-content {
          margin-left: 12px;
          margin-right: 0;
          width: calc(100% - 44px);
          
          .message-text {
            background: #f5f5f5;
          }
        }
      }
      
      .message-avatar {
        flex-shrink: 0;
      }
      
      .message-content {
        max-width: 80%;
        
        .reasoning-panel {
          margin-bottom: 8px;
          border: 1px solid #d9d9d9;
          border-radius: 8px;
          overflow: hidden;
          background: #fff;
          
          .reasoning-header {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 8px 12px;
            background: #fafafa;
            cursor: pointer;
            font-size: 13px;
            color: #666;
            
            .anticon-down {
              margin-left: auto;
              transition: transform 0.2s;
              
              &.rotated {
                transform: rotate(-90deg);
              }
            }
          }
          
          .reasoning-content {
            padding: 12px;
            font-size: 13px;
            color: #666;
            white-space: pre-wrap;
            max-height: 300px;
            overflow-y: auto;
            border-top: 1px solid #e8e8e8;
          }
        }
        
        .message-text {
          padding: 12px 16px;
          border-radius: 12px;
          line-height: 1.6;
          word-break: break-word;
          
          // 代码块容器
          :deep(.code-block) {
            margin: 12px 0;
            border-radius: 8px;
            overflow: hidden;
            background: #1e1e1e;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
            
            .code-header {
              display: flex;
              justify-content: space-between;
              align-items: center;
              padding: 8px 12px;
              background: #2d2d2d;
              border-bottom: 1px solid #3d3d3d;
              
              .code-lang {
                font-size: 12px;
                color: #9cdcfe;
                font-weight: 500;
                text-transform: uppercase;
                letter-spacing: 0.5px;
              }
              
              .copy-btn {
                background: transparent;
                border: 1px solid #555;
                color: #ccc;
                padding: 4px 10px;
                border-radius: 4px;
                font-size: 12px;
                cursor: pointer;
                transition: all 0.2s;
                
                &:hover {
                  background: #3d3d3d;
                  border-color: #888;
                  color: #fff;
                }
              }
            }
            
            pre {
              margin: 0;
              padding: 16px;
              background: #0d1117 !important;
              overflow-x: auto;
              
              &::-webkit-scrollbar {
                height: 6px;
              }
              
              &::-webkit-scrollbar-thumb {
                background: #444;
                border-radius: 3px;
              }
              
              code.hljs {
                background: transparent !important;
                padding: 0 !important;
                font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
                font-size: 13px;
                line-height: 1.6;
                white-space: pre;
              }
            }
          }
          
          // 普通 pre 样式
          :deep(pre:not(.code-block pre)) {
            background: #1e1e1e;
            padding: 16px;
            border-radius: 8px;
            overflow-x: auto;
            margin: 12px 0;
            
            code {
              color: #d4d4d4;
              font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
              font-size: 13px;
            }
          }
          
          // 行内代码
          :deep(code:not(pre code)) {
            background: rgba(0, 0, 0, 0.08);
            color: #c7254e;
            padding: 2px 6px;
            border-radius: 4px;
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.9em;
          }
          
          // 标题样式
          :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
            margin: 16px 0 8px 0;
            font-weight: 600;
            line-height: 1.4;
            color: #1a1a1a;
            
            &:first-child {
              margin-top: 0;
            }
          }
          
          :deep(h1) {
            font-size: 1.5em;
            padding-bottom: 8px;
            border-bottom: 2px solid #1890ff;
          }
          
          :deep(h2) {
            font-size: 1.3em;
            padding-bottom: 6px;
            border-bottom: 1px solid #e8e8e8;
          }
          
          :deep(h3) { font-size: 1.15em; color: #1890ff; }
          :deep(h4) { font-size: 1.05em; color: #52c41a; }
          :deep(h5) { font-size: 1em; color: #666; }
          :deep(h6) { font-size: 0.9em; color: #999; }
          
          // 列表样式
          :deep(ul), :deep(ol) {
            margin: 8px 0;
            padding-left: 24px;
            
            li {
              margin: 4px 0;
              line-height: 1.6;
            }
          }
          
          :deep(ul) { list-style-type: disc; }
          :deep(ol) { list-style-type: decimal; }
          
          // 链接样式
          :deep(a) {
            color: #1890ff;
            text-decoration: none;
            
            &:hover {
              text-decoration: underline;
            }
          }
          
          // 强调样式
          :deep(strong) { font-weight: 600; color: #1a1a1a; }
          :deep(em) { font-style: italic; color: #666; }
          
          // 表格样式
          :deep(table) {
            width: 100%;
            border-collapse: collapse;
            margin: 12px 0;
            font-size: 14px;
            
            th, td {
              border: 1px solid #d9d9d9;
              padding: 8px 12px;
              text-align: left;
            }
            
            th {
              background: #fafafa;
              font-weight: 600;
              color: #1a1a1a;
            }
            
            tr:nth-child(even) td {
              background: #fafafa;
            }
            
            tr:hover td {
              background: #e6f7ff;
            }
          }
        }
        
        // 消息操作按钮
        .message-actions {
          display: flex;
          gap: 4px;
          margin-top: 4px;
          opacity: 0;
          transition: opacity 0.2s;
          
          :deep(.ant-btn) {
            color: #999;
            &:hover {
              color: #1890ff;
            }
          }
        }
      }
      
      &:hover .message-actions {
        opacity: 1;
      }
    }
    
    // 打字指示器
    .typing-indicator {
      display: flex;
      gap: 4px;
      padding: 12px 16px;
      
      span {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #999;
        animation: typing 1.4s infinite;
        
        &:nth-child(2) { animation-delay: 0.2s; }
        &:nth-child(3) { animation-delay: 0.4s; }
      }
    }
    
    @keyframes typing {
      0%, 60%, 100% { opacity: 0.3; transform: scale(0.8); }
      30% { opacity: 1; transform: scale(1); }
    }
    
    .empty-chat {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: #999;
      
      .empty-icon {
        font-size: 64px;
        color: #ccc;
      }
      
      p {
        margin-top: 16px;
        font-size: 16px;
      }
    }
  }
  
  .chat-input-area {
    padding: 16px;
    border-top: 1px solid #e8e8e8;
    
    .attachment-list {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-bottom: 8px;
      
      .attachment-item {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 4px 8px;
        background: #f5f5f5;
        border: 1px solid #e8e8e8;
        border-radius: 6px;
        font-size: 12px;
        max-width: 250px;
        
        .att-icon {
          color: #1890ff;
          font-size: 14px;
          flex-shrink: 0;
          
          &.loading {
            animation: spin 1s linear infinite;
          }
        }
        
        .att-name {
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          color: #333;
        }
        
        .att-size {
          color: #999;
          flex-shrink: 0;
        }
        
        .att-remove {
          color: #999;
          cursor: pointer;
          flex-shrink: 0;
          
          &:hover {
            color: #ff4d4f;
          }
        }
      }
    }
    
    @keyframes spin {
      from { transform: rotate(0deg); }
      to { transform: rotate(360deg); }
    }
    
    .input-wrapper {
      display: flex;
      gap: 8px;
      align-items: flex-end;
      
      .attach-btn {
        flex-shrink: 0;
        color: #999;
        
        &:hover {
          color: #1890ff;
        }
      }
      
      :deep(.ant-input) {
        flex: 1;
        resize: none;
      }
      
      .input-actions {
        display: flex;
        gap: 8px;
        flex-shrink: 0;
      }
    }
  }
}

// Mermaid 预览按钮样式
:deep(.mermaid-preview-container) {
  margin: 12px 0;
  text-align: center;
  
  .mermaid-preview-btn {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
    border: none;
    padding: 12px 24px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.35);
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(102, 126, 234, 0.45);
    }
    
    &:active {
      transform: translateY(0);
    }
  }
}

// Mermaid 预览抽屉内容
.mermaid-preview-content {
  width: 100%;
  height: calc(90vh - 100px);
  overflow: auto;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 20px;
  
  .mermaid {
    width: 100%;
    min-height: 100px;
    
    svg {
      width: 80% !important;
      max-width: none !important;
      height: auto !important;
      min-height: 100px;
    }
  }
}

@media (max-width: 768px) {
  .conversation-sidebar {
    display: none;
  }
}
</style>
