import type { AxiosRequestConfig } from 'axios'
import type { ApiResponse, Dept, ManageableDeptTreeData } from '@/types'
import request from '@/utils/request'

export function getDeptTree() {
  return request.get<any, ApiResponse<Dept[]>>('/depts/tree')
}

export function getManageableDeptTree(config?: AxiosRequestConfig) {
  return request.get<any, ApiResponse<ManageableDeptTreeData>>('/depts/manageable-tree', config)
}

export function getDept(id: number) {
  return request.get<any, ApiResponse<Dept>>(`/depts/${id}`)
}

export function createDept(data: any) {
  return request.post<any, ApiResponse>('/depts', data)
}

export function updateDept(id: number, data: any) {
  return request.put<any, ApiResponse>(`/depts/${id}`, data)
}

export function deleteDept(id: number) {
  return request.delete<any, ApiResponse>(`/depts/${id}`)
}
