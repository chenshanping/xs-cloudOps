import request from '@/utils/request'
import type { ApiResponse } from '@/types'

// 字典类型
export interface DictType {
  id: number
  name: string
  type: string
  status: number
  remark: string
  created_at: string
  updated_at: string
}

// 字典数据
export interface DictData {
  id: number
  dict_type: string
  label: string
  value: string
  sort: number
  status: number
  tag_type: string
  is_default: number
  remark: string
  created_at: string
  updated_at: string
}

// 分页响应
export interface PageResponse<T> {
  list: T[]
  total: number
}

// ==================== 字典类型 ====================

// 获取字典类型列表（分页）
export function getDictTypeList(params?: { page?: number; page_size?: number; name?: string; type?: string; status?: number }) {
  return request.get<any, ApiResponse<PageResponse<DictType>>>('/dict/types', { params })
}

// 获取所有字典类型（不分页）
export function getAllDictTypes() {
  return request.get<any, ApiResponse<DictType[]>>('/dict/types/all')
}

// 获取字典类型详情
export function getDictType(id: number) {
  return request.get<any, ApiResponse<DictType>>(`/dict/types/${id}`)
}

// 创建字典类型
export function createDictType(data: { name: string; type: string; status?: number; remark?: string }) {
  return request.post<any, ApiResponse>('/dict/types', data)
}

// 更新字典类型
export function updateDictType(id: number, data: { name?: string; type?: string; status?: number; remark?: string }) {
  return request.put<any, ApiResponse>(`/dict/types/${id}`, data)
}

// 删除字典类型
export function deleteDictType(id: number) {
  return request.delete<any, ApiResponse>(`/dict/types/${id}`)
}

// ==================== 字典数据 ====================

// 获取字典数据列表（分页）
export function getDictDataList(params: { dict_type: string; page?: number; page_size?: number; label?: string; status?: number }) {
  return request.get<any, ApiResponse<PageResponse<DictData>>>('/dict/data', { params })
}

// 根据字典类型获取字典数据（不分页，用于下拉框）
export function getDictDataByType(type: string) {
  return request.get<any, ApiResponse<DictData[]>>(`/dict/type/${type}`)
}

// 获取字典数据详情
export function getDictData(id: number) {
  return request.get<any, ApiResponse<DictData>>(`/dict/data/${id}`)
}

// 创建字典数据
export function createDictData(data: {
  dict_type: string
  label: string
  value: string
  sort?: number
  status?: number
  tag_type?: string
  is_default?: number
  remark?: string
}) {
  return request.post<any, ApiResponse>('/dict/data', data)
}

// 更新字典数据
export function updateDictData(id: number, data: {
  label?: string
  value?: string
  sort?: number
  status?: number
  tag_type?: string
  is_default?: number
  remark?: string
}) {
  return request.put<any, ApiResponse>(`/dict/data/${id}`, data)
}

// 删除字典数据
export function deleteDictData(id: number) {
  return request.delete<any, ApiResponse>(`/dict/data/${id}`)
}
