import type { ApiResponse, PageResponse, User } from '@/types';
import request from '@/utils/request';


// 获取用户列表
export function getUserList(params: any) {
  return request.get<any, ApiResponse<PageResponse<User>>>('/users', { params })
}

// 获取用户选项（轻量级，用于下拉选择）
export function getUserOptions() {
  return request.get<any, ApiResponse<{ id: number; username: string; nickname: string }[]>>('/users/options')
}

// 获取用户详情
export function getUser(id: number) {
  return request.get<any, ApiResponse<User>>(`/users/${id}`)
}

// 创建用户
export function createUser(data: any) {
  return request.post<any, ApiResponse>('/users', data)
}

// 获取用户默认密码（用于新增用户表单预填）
export function getUserDefaultPassword() {
  return request.get<any, ApiResponse<{ password: string }>>('/users/default-password')
}

// 更新用户
export function updateUser(id: number, data: any) {
  return request.put<any, ApiResponse>(`/users/${id}`, data)
}

// 删除用户
export function deleteUser(id: number) {
  return request.delete<any, ApiResponse>(`/users/${id}`)
}

// 批量删除用户
export function batchDeleteUsers(ids: number[]) {
  return request.delete<any, ApiResponse>('/users/batch', { data: { ids } })
}

// 修改用户状态
export function updateUserStatus(id: number, status: number) {
  return request.put<any, ApiResponse>(`/users/${id}/status`, { status })
}

// 批量修改用户状态
export function batchUpdateUserStatus(ids: number[], status: number) {
  return request.put<any, ApiResponse>('/users/batch-status', { ids, status })
}

// 重置密码
export function resetPassword(id: number) {
  return request.put<any, ApiResponse>(`/users/${id}/password`, {})
}

// 批量重置密码
export function batchResetPassword(ids: number[]) {
  return request.put<any, ApiResponse>('/users/batch-password', { ids })
}

// 修改密码
export function changePassword(data: { old_password: string; new_password: string }) {
  return request.put<any, ApiResponse>('/user/password', data)
}

// 获取个人资料
export function getProfile() {
  return request.get<any, ApiResponse<User>>('/user/profile')
}

// 更新个人资料
export function updateProfile(data: { nickname?: string; email?: string; phone?: string }) {
  return request.put<any, ApiResponse>('/user/profile', data)
}

// 更新头像
export function updateAvatar(fileId: number) {
  return request.put<any, ApiResponse>('/user/avatar', { file_id: fileId })
}

export function forceUserOffline(id: number) {
  return request.post(`/users/${id}/offline`)
}

// 下载用户导入模板
export function downloadUserImportTemplate(deptId?: number) {
  return request.get('/users/import-template', { params: { dept_id: deptId }, responseType: 'blob' })
}

// 导入用户
export interface ImportError {
  row: number
  column: string
  value: string
  message: string
}

export interface ImportResult {
  total_count: number
  success_count: number
  failed_count: number
  errors: ImportError[]
}

export function importUsers(file: File, deptId: number) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('dept_id', String(deptId))
  return request.post<any, ApiResponse<ImportResult>>('/users/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 导出用户
export function exportUsers(deptId: number, ids?: number[]) {
  const params: Record<string, any> = { dept_id: deptId }
  if (ids && ids.length > 0) {
    params.ids = ids.join(',')
  }
  return request.get('/users/export', { params, responseType: 'blob' })
}

// 字段配置
export interface FieldConfig {
  key: string      // 字段名
  label: string    // 中文标签
  required: boolean // 是否必填
  type: string     // 字段类型: text, image, file, images, files, select
}

// 用户身份信息
export interface UserProfile {
  key: string           // 身份标识，如 "doctor", "merchant"
  name: string          // 显示名称，如 "医生", "商家"
  description: string   // 描述
  icon: string          // 图标（ant-design图标名）
  data: any             // 身份数据
  has_profile: boolean  // 当前用户是否有此身份
  is_complete: boolean  // 是否已完善（所有必填字段都有值）
  fields: FieldConfig[] // 字段配置
}

// 获取当前用户所有身份信息
export function getUserProfiles() {
  return request.get<any, ApiResponse<UserProfile[]>>('/user/profiles')
}

// 获取系统已注册的所有身份类型
export function getRegisteredProfiles() {
  return request.get<any, ApiResponse<UserProfile[]>>('/user/profiles/types')
}

// 获取指定用户的身份信息（管理员使用）
export function getUserProfilesById(id: number) {
  return request.get<any, ApiResponse<UserProfile[]>>(`/users/${id}/profiles`)
}
