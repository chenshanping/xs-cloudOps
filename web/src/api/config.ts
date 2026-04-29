import request from '@/utils/request';
import type { StorageType } from '@/types/storage'


export interface SysConfig {
  id: number
  name: string
  key: string
  value: string
  value_type: 'string' | 'json'
  remark: string
  created_at: string
  updated_at: string
}

// 获取配置列表
export function getConfigList(key?: string) {
  return request.get<SysConfig[]>('/configs', { params: { key } })
}

// 根据 key 获取配置
export function getConfigByKey(key: string) {
  return request.get<SysConfig>(`/configs/key/${key}`)
}

// 批量获取配置
export function getConfigsByKeys(keys: string[]) {
  return request.post<Record<string, SysConfig>>('/configs/keys', { keys })
}

// 创建配置
export function createConfig(data: Partial<SysConfig>) {
  return request.post<SysConfig>('/configs', data)
}

// 更新配置
export function updateConfig(id: number, data: Partial<SysConfig>) {
  return request.put(`/configs/${id}`, data)
}

// 批量更新配置
export function batchUpdateConfigs(configs: Record<string, string>) {
  return request.put('/configs/batch', configs)
}

// 删除配置
export function deleteConfig(id: number) {
  return request.delete(`/configs/${id}`)
}

// 发送测试邮件
export function sendTestEmail(email: string) {
  return request.post('/configs/test-email', { email }, { silent: true })
}

export function testStorageConfig(data: { type: StorageType; config: string }) {
  return request.post('/configs/storage/test', data, { silent: true })
}
