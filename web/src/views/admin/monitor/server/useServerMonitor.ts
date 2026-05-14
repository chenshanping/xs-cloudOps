import { ref } from 'vue'
import { message } from 'ant-design-vue'
import {
  clearCacheByPrefix,
  getDbStats,
  getDependencyHealth,
  getOssHealth,
  getRedisInfo,
  getRuntimeInfo,
  getServerInfo,
  type ClearCacheResult,
  type DBStats,
  type DependencyHealth,
  type OssHealth,
  type RedisInfo,
  type RuntimeInfo,
  type ServerInfo,
} from '@/api/monitor'

export type MonitorTabKey = 'server' | 'runtime' | 'db' | 'redis' | 'oss'

function getErrorMessage(error: unknown) {
  if (error instanceof Error && error.message) {
    return error.message
  }
  return '加载失败，请稍后重试'
}

export function useServerMonitor() {
  const dependency = ref<DependencyHealth | null>(null)
  const serverInfo = ref<ServerInfo | null>(null)
  const runtimeInfo = ref<RuntimeInfo | null>(null)
  const dbStats = ref<DBStats | null>(null)
  const redisInfo = ref<RedisInfo | null>(null)
  const ossHealth = ref<OssHealth | null>(null)

  const dependencyLoading = ref(false)
  const serverLoading = ref(false)
  const runtimeLoading = ref(false)
  const dbLoading = ref(false)
  const redisLoading = ref(false)
  const ossLoading = ref(false)
  const clearLoading = ref(false)

  const dependencyError = ref('')
  const serverError = ref('')
  const runtimeError = ref('')
  const dbError = ref('')
  const redisError = ref('')
  const ossError = ref('')

  async function loadDependency() {
    dependencyLoading.value = true
    dependencyError.value = ''
    try {
      const res = await getDependencyHealth()
      dependency.value = res.data
    } catch (error) {
      dependencyError.value = getErrorMessage(error)
    } finally {
      dependencyLoading.value = false
    }
  }

  async function loadServerInfo() {
    serverLoading.value = true
    serverError.value = ''
    try {
      const res = await getServerInfo()
      serverInfo.value = res.data
    } catch (error) {
      serverError.value = getErrorMessage(error)
    } finally {
      serverLoading.value = false
    }
  }

  async function loadRuntimeInfo() {
    runtimeLoading.value = true
    runtimeError.value = ''
    try {
      const res = await getRuntimeInfo()
      runtimeInfo.value = res.data
    } catch (error) {
      runtimeError.value = getErrorMessage(error)
    } finally {
      runtimeLoading.value = false
    }
  }

  async function loadDbStats() {
    dbLoading.value = true
    dbError.value = ''
    try {
      const res = await getDbStats()
      dbStats.value = res.data
    } catch (error) {
      dbError.value = getErrorMessage(error)
    } finally {
      dbLoading.value = false
    }
  }

  async function loadRedisInfo() {
    redisLoading.value = true
    redisError.value = ''
    try {
      const res = await getRedisInfo()
      redisInfo.value = res.data
    } catch (error) {
      redisError.value = getErrorMessage(error)
    } finally {
      redisLoading.value = false
    }
  }

  async function loadOssHealth() {
    ossLoading.value = true
    ossError.value = ''
    try {
      const res = await getOssHealth()
      ossHealth.value = res.data
    } catch (error) {
      ossError.value = getErrorMessage(error)
    } finally {
      ossLoading.value = false
    }
  }

  async function clearRedisCache(prefix: string): Promise<ClearCacheResult | null> {
    clearLoading.value = true
    try {
      const res = await clearCacheByPrefix(prefix)
      message.success(`已清理 ${res.data.deleted} 个缓存键`)
      await loadRedisInfo()
      await loadDependency()
      return res.data
    } catch {
      return null
    } finally {
      clearLoading.value = false
    }
  }

  async function loadTab(key: MonitorTabKey) {
    if (key === 'server') return loadServerInfo()
    if (key === 'runtime') return loadRuntimeInfo()
    if (key === 'db') return loadDbStats()
    if (key === 'redis') return loadRedisInfo()
    return loadOssHealth()
  }

  async function loadTabIfEmpty(key: MonitorTabKey) {
    if (key === 'server' && !serverInfo.value) return loadServerInfo()
    if (key === 'runtime' && !runtimeInfo.value) return loadRuntimeInfo()
    if (key === 'db' && !dbStats.value) return loadDbStats()
    if (key === 'redis' && !redisInfo.value) return loadRedisInfo()
    if (key === 'oss' && !ossHealth.value) return loadOssHealth()
  }

  return {
    dependency,
    serverInfo,
    runtimeInfo,
    dbStats,
    redisInfo,
    ossHealth,
    dependencyLoading,
    serverLoading,
    runtimeLoading,
    dbLoading,
    redisLoading,
    ossLoading,
    clearLoading,
    dependencyError,
    serverError,
    runtimeError,
    dbError,
    redisError,
    ossError,
    loadDependency,
    loadServerInfo,
    loadRuntimeInfo,
    loadDbStats,
    loadRedisInfo,
    loadOssHealth,
    clearRedisCache,
    loadTab,
    loadTabIfEmpty,
  }
}
