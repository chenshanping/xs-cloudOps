import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import { useUserStore } from '@/store/user';
import { message } from 'ant-design-vue';
import router from '@/router';

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
  (response) => {
    const res = response.data
    const silent = response.config.silent
    
    if (res.code !== 200) {
      // 401: Token失效，尝试刷新
      if (res.code === 401 && !silent) {
        return handle401Error(response)
      }
      
      // 其他错误
      if (!silent) {
        message.error(res.message || '请求失败')
      }
      
      return Promise.reject(new Error(res.message || '请求失败'))
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
    const friendlyError = new Error(errorMessage)
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
    return Promise.reject(new Error('Token已失效'))
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
    return Promise.reject(new Error('Token刷新失败'))
  } finally {
    isRefreshing = false
  }
}

export default request
