import type { AxiosRequestConfig } from 'axios'
import request from '@/utils/request'
import type { ApiResponse, Role } from '@/types'

export interface RoleUpsertPayload {
  name: string
  code: string
  sort: number
  status: number
  is_super_admin: boolean
  data_scope: number
  dept_ids: number[]
  remark: string
}

// 获取角色列表
export function getRoleList() {
  return request.get<any, ApiResponse<Role[]>>('/roles')
}

// 获取角色详情
export function getRole(id: number) {
  return request.get<any, ApiResponse<Role>>(`/roles/${id}`)
}

// 创建角色
export function createRole(data: RoleUpsertPayload) {
  return request.post<any, ApiResponse>('/roles', data)
}

// 更新角色
export function updateRole(id: number, data: RoleUpsertPayload) {
	return request.put<any, ApiResponse>(`/roles/${id}`, data)
}

// 删除角色
export function deleteRole(id: number) {
  return request.delete<any, ApiResponse>(`/roles/${id}`)
}

// 分配菜单
export function assignMenus(id: number, menuIds: number[], config?: AxiosRequestConfig) {
  return request.put<any, ApiResponse>(`/roles/${id}/menus`, { menu_ids: menuIds }, config)
}

// 分配API
export function assignApis(id: number, apiIds: number[], config?: AxiosRequestConfig) {
  return request.put<any, ApiResponse>(`/roles/${id}/apis`, { api_ids: apiIds }, config)
}
