import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types'
import type { Doctor, CreateDoctorRequest, UpdateDoctorRequest, DoctorQuery, OptionItem, SaveMyDoctorRequest } from '@/types/doctor'

// 获取医生列表
export function getDoctorList(params: DoctorQuery) {
  return request.get<any, ApiResponse<PageResponse<Doctor>>>('/doctor', { params })
}

// 获取医生详情
export function getDoctor(id: number) {
  return request.get<any, ApiResponse<Doctor>>(`/doctor/${id}`)
}

// 创建医生
export function createDoctor(data: CreateDoctorRequest) {
  return request.post<any, ApiResponse<Doctor>>('/doctor', data)
}

// 更新医生
export function updateDoctor(id: number, data: UpdateDoctorRequest) {
  return request.put<any, ApiResponse<Doctor>>(`/doctor/${id}`, data)
}

// 删除医生
export function deleteDoctor(id: number) {
  return request.delete<any, ApiResponse>(`/doctor/${id}`)
}

// 批量删除医生
export function batchDeleteDoctor(ids: number[]) {
  return request.delete<any, ApiResponse>('/doctor/batch', { data: { ids } })
}

// 获取医生选项列表
export function getDoctorOptions(params?: { display_field?: string; count_table?: string; count_field?: string; exclude_deleted?: boolean; count_created_by?: number }) {
  return request.get<any, ApiResponse<OptionItem[]>>('/doctor/options', { params })
}

// 获取医生创建人选项列表
export function getDoctorCreatorOptions() {
  return request.get<any, ApiResponse<OptionItem[]>>('/doctor/creator/options')
}

// 审批医生
export function auditDoctor(id: number, data: { audit_status: number; audit_remark: string }) {
  return request.post<any, ApiResponse>(`/doctor/${id}/audit`, data)
}

// 获取我的医生信息
export function getMyDoctor() {
  return request.get<any, ApiResponse<Doctor | null>>('/doctor/my')
}

// 保存我的医生信息
export function saveMyDoctor(data: SaveMyDoctorRequest) {
  return request.post<any, ApiResponse>('/doctor/my', data)
}
