import request from '@/utils/request'
import type { Storage } from '@/types/storage'

// 获取存储配置列表
export function getStorageList() {
  return request.get<Storage[]>('/storages')
}

// 获取存储配置详情
export function getStorage(id: number) {
  return request.get<Storage>(`/storages/${id}`)
}

// 创建存储配置
export function createStorage(data: Partial<Storage>) {
  return request.post<Storage>('/storages', data)
}

// 更新存储配置
export function updateStorage(id: number, data: Partial<Storage>) {
  return request.put(`/storages/${id}`, data)
}

// 删除存储配置
export function deleteStorage(id: number) {
  return request.delete(`/storages/${id}`)
}

// 设置默认存储
export function setDefaultStorage(id: number) {
  return request.put(`/storages/${id}/default`)
}

// 测试存储配置
export function testStorage(data: Partial<Storage>) {
  return request.post('/storages/test', data)
}
