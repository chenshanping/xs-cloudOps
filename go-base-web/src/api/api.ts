import request from '@/utils/request'
import type { ApiResponse, PageResponse, Api } from '@/types'

// 获取API列表
export function getApiList(params: any) {
  return request.get<any, ApiResponse<PageResponse<Api>>>('/apis', { params })
}

// 获取全部API
export function getAllApis() {
  return request.get<any, ApiResponse<Api[]>>('/apis/all')
}

// API分组统计结果
export interface ApiGroupStats {
  group: string
  api_count: number
}

// 获取API分组(含数量)
export function getApiGroups() {
  return request.get<any, ApiResponse<ApiGroupStats[]>>('/apis/groups')
}

// 获取API详情
export function getApi(id: number) {
  return request.get<any, ApiResponse<Api>>(`/apis/${id}`)
}

// 创建API
export function createApi(data: any) {
  return request.post<any, ApiResponse>('/apis', data)
}

// 更新API
export function updateApi(id: number, data: any) {
  return request.put<any, ApiResponse>(`/apis/${id}`, data)
}

// 删除API
export function deleteApi(id: number) {
  return request.delete(`/apis/${id}`)
}

// 同步API路由到数据库
export function syncApis() {
  return request.post<{ added: number; updated: number; total: number }>('/apis/sync')
}
