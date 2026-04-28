import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import { useUserStore } from '@/store/user';
import { message } from 'ant-design-vue';
import router from '@/router';

type HandledRequestError = Error & {
  handledByMessage?: boolean
  errorSource?: 'request'
}

// 扩展 AxiosRequestConfig 支持自定义配置
declare module 'axios' {
  interface AxiosRequestConfig {
    silent?: boolean // 静默模式，不自动弹出错误提示
    _retry?: boolean // 标记是否已重试
  }
}

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

function extractErrorMessage(payload: any, fallback: string) {
  if (!payload || typeof payload !== 'object') {
    return fallback
  }

  const messageText = payload.message || payload.msg || payload.error
  if (typeof messageText === 'string' && messageText.trim()) {
    return messageText.trim()
  }

  return fallback
}

function createHandledRequestError(errorMessage: string): HandledRequestError {
  const error = new Error(errorMessage) as HandledRequestError
  error.handledByMessage = true
  error.errorSource = 'request'
  return error
}

// Token 刷新状态
let isRefreshing = false
let refreshSubscribers: Array<(token: string) => void> = []

// 添加等待刷新的请求
function subscribeTokenRefresh(callback: (token: string) => void) {
  refreshSubscribers.push(callback)
}

// 刷新完成后通知所有等待的请求
function onTokenRefreshed(newToken: string) {
  refreshSubscribers.forEach(callback => callback(newToken))
  refreshSubscribers = []
}

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  async (response) => {
    // 如果是 blob 类型响应（文件下载）
    if (response.config.responseType === 'blob') {
      // 检查是否是错误响应（服务器返回JSON错误而不是文件）
      const contentType = response.headers['content-type']
      if (contentType && contentType.includes('application/json')) {
        // 将 blob 转换为 JSON 读取错误信息
        const text = await response.data.text()
        const errorData = JSON.parse(text)
        const errorMessage = extractErrorMessage(errorData, '操作失败')
        message.error(errorMessage)
        return Promise.reject(createHandledRequestError(errorMessage))
      }
      // 正常的文件响应，直接返回
      return response.data
    }
    
    const res = response.data
    const silent = response.config.silent
    
    if (res.code !== 200) {
      // 401: Token失效，尝试刷新
      if (res.code === 401 && !silent) {
        return handle401Error(response)
      }
      
      // 其他错误
      const errorMessage = extractErrorMessage(res, '请求失败')
      if (!silent) {
        message.error(errorMessage)
      }
      
      return Promise.reject(createHandledRequestError(errorMessage))
    }
    return res
  },
  (error) => {
    const silent = error.config?.silent
    
    // 获取更友好的错误提示
    let errorMessage = '网络错误'
    if (error.response) {
      // 服务器返回了错误状态码
      const status = error.response.status
      const responseData = error.response.data

      if (responseData instanceof Blob) {
        errorMessage = `请求失败 (${status})`
      } else {
        const backendMessage = extractErrorMessage(responseData, '')
        if (backendMessage) {
          errorMessage = backendMessage
        } else {
          switch (status) {
            case 500:
              errorMessage = '服务器内部错误'
              break
            case 502:
              errorMessage = '网关错误'
              break
            case 503:
              errorMessage = '服务不可用'
              break
            case 504:
              errorMessage = '网关超时'
              break
            default:
              errorMessage = `请求失败 (${status})`
          }
        }
      }
    } else if (error.code === 'ECONNABORTED') {
      errorMessage = '请求超时'
    } else if (!navigator.onLine) {
      errorMessage = '网络已断开'
    } else {
      errorMessage = '无法连接服务器'
    }
    
    if (!silent) {
      message.error(errorMessage)
    }
    
    // 返回一个带有友好消息的错误
    const friendlyError = createHandledRequestError(errorMessage)
    return Promise.reject(friendlyError)
  }
)

// 处理 401 错误，尝试刷新 Token
async function handle401Error(response: AxiosResponse) {
  const originalConfig = response.config
  const userStore = useUserStore()
  
  // 如果是刷新接口本身失败，或已经重试过，直接登出
  if (originalConfig.url === '/auth/refresh' || originalConfig._retry) {
    await userStore.logoutAction()
    router.push('/login')
    return Promise.reject(createHandledRequestError('Token已失效'))
  }
  
  // 标记已重试
  originalConfig._retry = true
  
  // 如果正在刷新，等待刷新完成后重试
  if (isRefreshing) {
    return new Promise((resolve) => {
      subscribeTokenRefresh((newToken: string) => {
        originalConfig.headers.Authorization = `Bearer ${newToken}`
        resolve(request(originalConfig))
      })
    })
  }
  
  isRefreshing = true
  
  try {
    // 调用刷新接口
    const res = await request.post('/auth/refresh', {}, { silent: true })
    const newToken = res.data.token
    
    // 更新本地 Token
    userStore.token = newToken
    localStorage.setItem('token', newToken)
    
    // 通知所有等待的请求
    onTokenRefreshed(newToken)
    
    // 重试原请求
    originalConfig.headers.Authorization = `Bearer ${newToken}`
    return request(originalConfig)
  } catch (error) {
    // 刷新失败，登出
    await userStore.logoutAction()
    router.push('/login')
    return Promise.reject(createHandledRequestError('Token刷新失败'))
  } finally {
    isRefreshing = false
  }
}

export default request
