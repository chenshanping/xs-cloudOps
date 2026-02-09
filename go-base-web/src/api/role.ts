import request from '@/utils/request'
import type { ApiResponse, Role } from '@/types'

// 获取角色列表
export function getRoleList() {
  return request.get<any, ApiResponse<Role[]>>('/roles')
}

// 获取角色详情
export function getRole(id: number) {
  return request.get<any, ApiResponse<Role>>(`/roles/${id}`)
}

// 创建角色
export function createRole(data: any) {
  return request.post<any, ApiResponse>('/roles', data)
}

// 更新角色
export function updateRole(id: number, data: any) {
  return request.put<any, ApiResponse>(`/roles/${id}`, data)
}

// 删除角色
export function deleteRole(id: number) {
  return request.delete<any, ApiResponse>(`/roles/${id}`)
}

// 分配菜单
export function assignMenus(id: number, menuIds: number[]) {
  return request.put<any, ApiResponse>(`/roles/${id}/menus`, { menu_ids: menuIds })
}

// 分配API
export function assignApis(id: number, apiIds: number[]) {
  return request.put<any, ApiResponse>(`/roles/${id}/apis`, { api_ids: apiIds })
}
