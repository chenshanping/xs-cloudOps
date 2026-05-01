import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types'

// AI模型信息
export interface ModelInfo {
  id: string
  name: string
  description: string
}

// AI对话
export interface AIConversation {
  id: number
  user_id: number
  title: string
  model: string
  context_cleared_at?: string
  created_at: string
  updated_at: string
}

// AI消息
export interface AIMessage {
  id: number
  conversation_id: number
  role: 'user' | 'assistant' | 'system'
  content: string
  reasoning_content: string
  file_ids?: string  // JSON数组字符串
  created_at: string
}

export interface BatchDeleteConversationResult {
  success_count: number
  failed_count: number
  failed_msgs: string[]
}

export interface AIModelConfig {
  id: string
  name: string
  description: string
  is_thinking?: boolean
  support_vision?: boolean
  support_tools?: boolean
  search_strategy?: 'none' | 'builtin' | 'tool'
  support_embedding?: boolean
  support_rerank?: boolean
  is_free?: boolean
  tags?: string[]
  temperature?: number | null
  context_window?: number | null
}

export interface AIProviderConfig {
  name: string
  api_key: string
  base_url: string
  models: AIModelConfig[]
}

export interface AIAdminConfig {
  default_provider: string
  providers: AIProviderConfig[]
}

export interface AITestRequest {
  api_key: string
  base_url: string
  model: string
}

export interface AIProviderRemoteModel {
  id: string
  name?: string
  description?: string
  object?: string
  created?: number
  owned_by?: string
  is_thinking?: boolean
  support_vision?: boolean
  support_tools?: boolean
  search_strategy?: 'none' | 'builtin' | 'tool'
  support_embedding?: boolean
  support_rerank?: boolean
  is_free?: boolean
  temperature?: number | null
  context_window?: number | null
  tags?: string[]
}

export interface AIProviderModelsFetchRequest {
  api_key: string
  base_url: string
  provider_name?: string
}

export interface AIProviderModelsFetchResponse {
  models: AIProviderRemoteModel[]
}

// 对话请求
export interface ChatRequest {
  conversation_id?: number
  model?: string
  message: string
  file_ids?: number[]
  enable_search?: boolean
  enable_thinking?: boolean
  save_conversation?: boolean  // 是否保存对话记录，默认true
}

// 获取模型列表
export function getModels() {
  return request.get<any, ApiResponse<ModelInfo[]>>('/ai/models')
}

// 获取对话列表
export function getConversations(params: { page?: number; page_size?: number }) {
  return request.get<any, ApiResponse<PageResponse<AIConversation>>>('/ai/conversations', { params })
}

// 创建对话
export function createConversation(data: { title?: string; model?: string }) {
  return request.post<any, ApiResponse<AIConversation>>('/ai/conversations', data)
}

// 获取对话详情
export function getConversation(id: number) {
  return request.get<any, ApiResponse<{ conversation: AIConversation; messages: AIMessage[] }>>(`/ai/conversations/${id}`)
}

// 更新对话标题
export function updateConversationTitle(id: number, title: string) {
  return request.put<any, ApiResponse>(`/ai/conversations/${id}`, { title })
}

// 删除对话
export function deleteConversation(id: number) {
  return request.delete<any, ApiResponse>(`/ai/conversations/${id}`)
}

export function deleteConversations(ids: number[]) {
  return request.delete<any, ApiResponse<BatchDeleteConversationResult>>('/ai/conversations/batch', { data: { ids } })
}

export function getAIConfig() {
  return request.get<any, ApiResponse<AIAdminConfig>>('/ai/config')
}

export function updateAIConfig(data: AIAdminConfig) {
  return request.put<any, ApiResponse>('/ai/config', data)
}

// 获取对话消息
export function getMessages(conversationId: number) {
  return request.get<any, ApiResponse<AIMessage[]>>(`/ai/conversations/${conversationId}/messages`)
}

// 清空对话消息
export function clearMessages(conversationId: number) {
  return request.delete<any, ApiResponse>(`/ai/conversations/${conversationId}/messages`)
}

// 清空上下文（保留聊天记录）
export function clearContext(conversationId: number) {
  return request.post<any, ApiResponse>(`/ai/conversations/${conversationId}/clear-context`)
}

// 删除单条消息
export function deleteMessage(messageId: number) {
  return request.delete<any, ApiResponse>(`/ai/messages/${messageId}`)
}

export function aiTest(config: AITestRequest) {
  return request.post<any, ApiResponse>('/ai/test', config, { silent: true })
}

export function fetchAIProviderModels(data: AIProviderModelsFetchRequest) {
  return request.post<any, ApiResponse<AIProviderModelsFetchResponse>>('/ai/providers/models/fetch', data, { silent: true })
}

// 普通对话（非流式）
export function chat(data: ChatRequest) {
  return request.post<any, ApiResponse<AIMessage>>('/ai/chat', data)
}

// 流式对话（SSE）
export function chatStream(
  data: ChatRequest,
  onMessage: (data: { content: string; reasoning_content: string; finish_reason: string }) => void,
  onConversationId: (id: number) => void,
  onError?: (error: Error) => void,
  onDone?: () => void
): () => void {
  const token = localStorage.getItem('token')
  
  // 使用fetch进行SSE请求（POST方法）
  const controller = new AbortController()
  
  // 设置超时
  const timeoutId = setTimeout(() => {
    controller.abort()
    onError?.(new Error('请求超时，请重试'))
  }, 180000) // 180秒超时
  
  fetch('/api/v1/ai/chat/stream', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      'Accept': 'text/event-stream'
    },
    body: JSON.stringify(data),
    signal: controller.signal
  })
    .then(async response => {
      clearTimeout(timeoutId)
      
      const contentType = response.headers.get('Content-Type') || ''
      
      // 如果返回的是 JSON（而不是 SSE），说明是错误响应
      if (contentType.includes('application/json')) {
        const errorData = await response.json()
        // 解析错误信息
        let errorMessage = errorData.message || '请求失败'
        // 如果 message 内部还包含 JSON，尝试提取具体错误
        if (errorMessage.includes('{"error"')) {
          try {
            const innerError = JSON.parse(errorMessage.replace(/^[^{]*/, ''))
            errorMessage = innerError.error?.message || errorMessage
          } catch {
            // 保持原始消息
          }
        }
        throw new Error(errorMessage)
      }
      
      // HTTP 状态码错误
      if (!response.ok) {
        let errorMessage = `请求失败 (${response.status})`
        try {
          const text = await response.text()
          if (text) errorMessage = text
        } catch {
          // 忽略
        }
        throw new Error(errorMessage)
      }
      
      const reader = response.body?.getReader()
      if (!reader) {
        throw new Error('无法获取响应流')
      }
      
      const decoder = new TextDecoder()
      let buffer = ''
      let hasReceivedData = false
      
      const processLine = (line: string) => {
        // 跳过空行和 event: 行
        if (!line.trim() || line.startsWith('event:')) {
          return
        }
        
        // 处理 data: 行（兼容有无空格两种格式）
        if (line.startsWith('data:')) {
          const data = line.startsWith('data: ') ? line.slice(6) : line.slice(5)
          
          // 处理conversation_id事件
          if (!isNaN(Number(data)) && data.trim() !== '') {
            onConversationId(Number(data))
            hasReceivedData = true
            return
          }
          
          if (data === '[DONE]') {
            return
          }
          
          try {
            const parsed = JSON.parse(data)
            // 检查是否是错误响应
            if (parsed.error) {
              throw new Error(parsed.error.message || parsed.error)
            }
            onMessage(parsed)
            hasReceivedData = true
          } catch (e) {
            // 如果是JSON解析错误且不是空字符串，记录一下
            if (data.trim() && !(e instanceof SyntaxError)) {
              console.warn('SSE数据处理警告:', data, e)
            }
          }
        }
      }
      
      try {
        while (true) {
          const { done, value } = await reader.read()
          
          if (value) {
            buffer += decoder.decode(value, { stream: true })
            const lines = buffer.split('\n')
            buffer = lines.pop() || ''
            
            for (const line of lines) {
              processLine(line)
            }
          }
          
          if (done) {
            // 处理剩余的buffer
            if (buffer.trim()) {
              processLine(buffer)
            }
            break
          }
        }
        
        // 检查是否收到了数据
        if (!hasReceivedData) {
          throw new Error('未收到有效响应数据')
        }
        
        // 流结束时确保调用onDone
        onDone?.()
      } catch (readError) {
        throw readError
      } finally {
        reader.releaseLock()
      }
    })
    .catch(error => {
      clearTimeout(timeoutId)
      if (error.name === 'AbortError') {
        // 用户取消或超时，不重复调用 onError
        return
      }
      console.error('流式请求错误:', error)
      onError?.(error)
    })
  
  // 返回取消函数
  return () => {
    clearTimeout(timeoutId)
    controller.abort()
  }
}
