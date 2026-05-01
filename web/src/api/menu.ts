import request from '@/utils/request'
import type { Api, ApiResponse, Menu } from '@/types'

// 获取菜单列表(树形)
export function getMenuTree() {
  return request.get<any, ApiResponse<Menu[]>>('/menus')
}

// 获取菜单详情
export function getMenu(id: number) {
  return request.get<any, ApiResponse<Menu>>(`/menus/${id}`)
}

// 创建菜单
export function createMenu(data: any) {
  return request.post<any, ApiResponse<Menu>>('/menus', data)
}

// 更新菜单
export function updateMenu(id: number, data: any) {
  return request.put<any, ApiResponse>(`/menus/${id}`, data)
}

export function getMenuApis(id: number) {
  return request.get<any, ApiResponse<Api[]>>(`/menus/${id}/apis`)
}

export function updateMenuApis(id: number, apiIds: number[]) {
  return request.put<any, ApiResponse>(`/menus/${id}/apis`, { api_ids: apiIds })
}

// 删除菜单
export function deleteMenu(id: number) {
  return request.delete<any, ApiResponse>(`/menus/${id}`)
}

// 获取用户菜单(根据角色筛选)
export function getUserMenus() {
  return request.get<any, ApiResponse>('/user/menus')
}
