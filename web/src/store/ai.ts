import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import {
  getModels,
  getConversations,
  createConversation,
  getConversation,
  deleteConversation,
  deleteConversations,
  clearMessages,
  clearContext,
  deleteMessage,
  chatStream,
  type ModelInfo,
  type AIConversation,
  type AIMessage,
  type ChatRequest,
  type BatchDeleteConversationResult
} from '@/api/ai'
import type { ApiResponse } from '@/types'
import {
  loadAIPreferences,
  persistAIPreferences,
} from './ai-preferences'

export const useAIStore = defineStore('ai', () => {
  const savedPreferences = loadAIPreferences()

  // 状态
  const models = ref<ModelInfo[]>([])
  const conversations = ref<AIConversation[]>([])
  const currentConversation = ref<AIConversation | null>(null)
  const messages = ref<AIMessage[]>([])
  const loading = ref(false)
  const streaming = ref(false)
  const currentModel = ref('')
  const enableSearch = ref(savedPreferences.enableSearch)
  const enableThinking = ref(savedPreferences.enableThinking)
  
  // 临时的流式内容
  const streamingContent = ref('')
  const streamingReasoning = ref('')
  
  // 计算属性
  const hasConversations = computed(() => conversations.value.length > 0)

  function persistPreferences() {
    persistAIPreferences(undefined, {
      enableSearch: enableSearch.value,
      enableThinking: enableThinking.value,
    })
  }

  watch([enableSearch, enableThinking], persistPreferences)
  
  // 获取模型列表
  async function fetchModels() {
    try {
      const res = await getModels()
      models.value = res.data
      currentModel.value = res.data[0]?.id ?? ''
    } catch (error) {
      console.error('获取模型列表失败:', error)
    }
  }
  
  // 获取对话列表
  async function fetchConversations() {
    try {
      const res = await getConversations({ page: 1, page_size: 100 })
      conversations.value = res.data.list || []
    } catch (error) {
      console.error('获取对话列表失败:', error)
    }
  }
  
  // 创建新对话
  async function newConversation() {
    try {
      const res = await createConversation({ model: currentModel.value })
      const conversation = res.data
      conversations.value.unshift(conversation)
      currentConversation.value = conversation
      messages.value = []
      return conversation
    } catch (error) {
      console.error('创建对话失败:', error)
      throw error
    }
  }
  
  // 选择对话
  async function selectConversation(id: number) {
    try {
      loading.value = true
      const res = await getConversation(id)
      currentConversation.value = res.data.conversation
      messages.value = res.data.messages || []
      currentModel.value = res.data.conversation.model
    } catch (error) {
      console.error('获取对话详情失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  // 删除对话
  async function removeConversation(id: number) {
    try {
      await deleteConversation(id)
      conversations.value = conversations.value.filter(c => c.id !== id)
      if (currentConversation.value?.id === id) {
        currentConversation.value = null
        messages.value = []
      }
    } catch (error) {
      console.error('删除对话失败:', error)
      throw error
    }
  }

  async function removeConversationsBatch(ids: number[]) {
    try {
      const res = await deleteConversations(ids)
      await fetchConversations()

      if (currentConversation.value) {
        const stillExists = conversations.value.some(c => c.id === currentConversation.value?.id)
        if (!stillExists) {
          currentConversation.value = null
          messages.value = []
        }
      }

      return res as ApiResponse<BatchDeleteConversationResult>
    } catch (error) {
      console.error('批量删除对话失败:', error)
      throw error
    }
  }
  
  // 清空当前对话消息
  async function clearCurrentMessages() {
    if (!currentConversation.value) return
    try {
      await clearMessages(currentConversation.value.id)
      messages.value = []
    } catch (error) {
      console.error('清空消息失败:', error)
      throw error
    }
  }
  
  // 清空上下文（保留聊天记录）
  async function clearCurrentContext() {
    if (!currentConversation.value) return
    try {
      await clearContext(currentConversation.value.id)
      // 更新本地对话的 context_cleared_at
      currentConversation.value.context_cleared_at = new Date().toISOString()
    } catch (error) {
      console.error('清空上下文失败:', error)
      throw error
    }
  }
  
  // 删除单条消息
  async function removeMessage(messageId: number) {
    try {
      await deleteMessage(messageId)
      messages.value = messages.value.filter(m => m.id !== messageId)
    } catch (error) {
      console.error('删除消息失败:', error)
      throw error
    }
  }
  
  // 发送消息（流式）
  let cancelStream: (() => void) | null = null
  
  async function sendMessage(content: string, fileIds?: number[]) {
    if (!content.trim() || streaming.value) return
    
    // 添加用户消息到列表
    const userMessage: AIMessage = {
      id: Date.now(),
      conversation_id: currentConversation.value?.id || 0,
      role: 'user',
      content: content,
      reasoning_content: '',
      file_ids: fileIds?.length ? JSON.stringify(fileIds) : undefined,
      created_at: new Date().toISOString()
    }
    messages.value.push(userMessage)
    
    // 重置流式内容
    streamingContent.value = ''
    streamingReasoning.value = ''
    streaming.value = true
    lastUserContent.value = content
    
    const request: ChatRequest = {
      conversation_id: currentConversation.value?.id,
      model: currentModel.value,
      message: content,
      file_ids: fileIds?.length ? fileIds : undefined,
      enable_search: enableSearch.value,
      enable_thinking: enableThinking.value
    }
    console.log('[AIStore] sendMessage fileIds:', fileIds, 'request:', JSON.stringify(request))
    
    cancelStream = chatStream(
      request,
      // onMessage: 只更新流式ref，不重建数组
      (data) => {
        if (data.content) {
          streamingContent.value += data.content
        }
        if (data.reasoning_content) {
          streamingReasoning.value += data.reasoning_content
        }
      },
      // onConversationId
      (id) => {
        if (!currentConversation.value) {
          const newConv: AIConversation = {
            id,
            user_id: 0,
            title: content.slice(0, 50) + (content.length > 50 ? '...' : ''),
            model: currentModel.value,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString()
          }
          currentConversation.value = newConv
          conversations.value.unshift(newConv)
        }
        userMessage.conversation_id = id
      },
      // onError
      (error) => {
        console.error('流式对话错误:', error)
        // 流式失败时将错误写入消息列表
        messages.value.push({
          id: Date.now() + 1,
          conversation_id: currentConversation.value?.id || 0,
          role: 'assistant',
          content: '抱歉，发生了错误: ' + error.message,
          reasoning_content: '',
          created_at: new Date().toISOString()
        })
        streaming.value = false
      },
      // onDone: 将流式内容写入messages数组
      () => {
        if (streamingContent.value || streamingReasoning.value) {
          messages.value.push({
            id: Date.now() + 1,
            conversation_id: currentConversation.value?.id || 0,
            role: 'assistant',
            content: streamingContent.value,
            reasoning_content: streamingReasoning.value,
            created_at: new Date().toISOString()
          })
        }
        streaming.value = false
        fetchConversations()
      }
    )
  }
  
  // 重新生成最后一条AI回复
  const lastUserContent = ref('')
  
  async function regenerateLastMessage() {
    if (streaming.value || messages.value.length < 2) return
    // 移除最后一条AI消息
    const lastMsg = messages.value[messages.value.length - 1]
    if (lastMsg.role === 'assistant') {
      messages.value.pop()
    }
    // 移除最后一条用户消息
    const lastUser = messages.value[messages.value.length - 1]
    if (lastUser?.role === 'user') {
      const content = lastUser.content
      messages.value.pop()
      await sendMessage(content)
    }
  }
  
  // 停止流式输出
  function stopStreaming() {
    if (cancelStream) {
      cancelStream()
      cancelStream = null
    }
    streaming.value = false
  }
  
  // 设置当前模型
  function setModel(modelId: string) {
    currentModel.value = modelId
  }
  
  // 切换联网搜索
  function toggleSearch() {
    enableSearch.value = !enableSearch.value
  }
  
  // 切换思考模式
  function toggleThinking() {
    enableThinking.value = !enableThinking.value
  }
  
  return {
    // 状态
    models,
    conversations,
    currentConversation,
    messages,
    loading,
    streaming,
    currentModel,
    enableSearch,
    enableThinking,
    streamingContent,
    streamingReasoning,
    
    // 计算属性
    hasConversations,
    
    // 方法
    fetchModels,
    fetchConversations,
    newConversation,
    selectConversation,
    removeConversation,
    removeConversationsBatch,
    clearCurrentMessages,
    clearCurrentContext,
    sendMessage,
    regenerateLastMessage,
    lastUserContent,
    stopStreaming,
    setModel,
    toggleSearch,
    toggleThinking
  }
})
