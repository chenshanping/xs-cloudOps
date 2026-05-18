import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types'

export type CmdbVerifyStatus = 'pending' | 'success' | 'failed'
export type CmdbCredentialAuthType = 'password' | 'private_key'

export interface CmdbHostGroup {
  id: number
  name: string
  sort: number
  status: number
  remark: string
  created_at?: string
  updated_at?: string
}

export interface CmdbHostTag {
  id: number
  name: string
  color: string
  remark: string
  created_at?: string
  updated_at?: string
}

export interface CmdbCredentialSummary {
  id: number
  name: string
  auth_type: CmdbCredentialAuthType
  username: string
  remark: string
  bind_count: number
}

export interface CmdbHostCredentialSummary {
  id: number
  name: string
  auth_type: CmdbCredentialAuthType
  username: string
}

export interface CmdbHostTagView {
  id: number
  name: string
  color: string
}

export interface CmdbHostItem {
  id: number
  name: string
  group: CmdbHostGroup
  tags: CmdbHostTagView[]
  environment: string
  owner: string
  remark: string
  private_ip: string
  public_ip: string
  ssh_host: string
  ssh_port: number
  credential_summary: CmdbHostCredentialSummary
  verify_status: CmdbVerifyStatus
  verify_message: string
  last_verified_at?: string
  hostname: string
  os: string
  platform: string
  platform_version: string
  kernel_version: string
  architecture: string
  cpu_cores: number
  memory_mb: number
  updated_at: string
}

export interface CmdbHostImportRowResult {
  row: number
  name: string
  created: boolean
  verify_success: boolean
  error_message: string
  verify_message: string
}

export interface CmdbHostImportResult {
  total: number
  success_count: number
  failure_count: number
  rows: CmdbHostImportRowResult[]
}

export interface CmdbHostGroupPayload {
  name: string
  sort: number
  status: number
  remark: string
}

export interface CmdbHostTagPayload {
  name: string
  color: string
  remark: string
}

export interface CmdbCredentialPayload {
  name: string
  auth_type: CmdbCredentialAuthType
  username: string
  password: string
  private_key: string
  passphrase: string
  remark: string
}

export interface CmdbHostPayload {
  name: string
  group_id: number
  tag_ids: number[]
  environment: string
  owner: string
  private_ip: string
  public_ip: string
  ssh_host: string
  ssh_port: number
  credential_id: number
  remark: string
}

export function getCmdbHostGroups(params?: { name?: string; status?: number }) {
  return request.get<any, ApiResponse<CmdbHostGroup[]>>('/cmdb/host-groups', { params })
}

export function createCmdbHostGroup(data: CmdbHostGroupPayload) {
  return request.post<any, ApiResponse>('/cmdb/host-groups', data)
}

export function updateCmdbHostGroup(id: number, data: CmdbHostGroupPayload) {
  return request.put<any, ApiResponse>(`/cmdb/host-groups/${id}`, data)
}

export function deleteCmdbHostGroup(id: number) {
  return request.delete<any, ApiResponse>(`/cmdb/host-groups/${id}`)
}

export function getCmdbHostTags(params?: { name?: string }) {
  return request.get<any, ApiResponse<CmdbHostTag[]>>('/cmdb/host-tags', { params })
}

export function createCmdbHostTag(data: CmdbHostTagPayload) {
  return request.post<any, ApiResponse>('/cmdb/host-tags', data)
}

export function updateCmdbHostTag(id: number, data: CmdbHostTagPayload) {
  return request.put<any, ApiResponse>(`/cmdb/host-tags/${id}`, data)
}

export function deleteCmdbHostTag(id: number) {
  return request.delete<any, ApiResponse>(`/cmdb/host-tags/${id}`)
}

export function getCmdbCredentials(params: { page: number; page_size: number; name?: string; auth_type?: string }) {
  return request.get<any, ApiResponse<PageResponse<CmdbCredentialSummary>>>('/cmdb/ssh-credentials', { params })
}

export function getCmdbCredentialOptions() {
  return request.get<any, ApiResponse<CmdbCredentialSummary[]>>('/cmdb/ssh-credentials/options')
}

export function getCmdbCredential(id: number) {
  return request.get<any, ApiResponse<CmdbCredentialSummary>>(`/cmdb/ssh-credentials/${id}`)
}

export function createCmdbCredential(data: CmdbCredentialPayload) {
  return request.post<any, ApiResponse>('/cmdb/ssh-credentials', data)
}

export function updateCmdbCredential(id: number, data: CmdbCredentialPayload) {
  return request.put<any, ApiResponse>(`/cmdb/ssh-credentials/${id}`, data)
}

export function deleteCmdbCredential(id: number) {
  return request.delete<any, ApiResponse>(`/cmdb/ssh-credentials/${id}`)
}

export function getCmdbHosts(params: {
  page: number
  page_size: number
  keyword?: string
  group_id?: number
  tag_id?: number
  verify_status?: string
  environment?: string
}) {
  return request.get<any, ApiResponse<PageResponse<CmdbHostItem>>>('/cmdb/hosts', { params })
}

export function getCmdbHost(id: number) {
  return request.get<any, ApiResponse<CmdbHostItem>>(`/cmdb/hosts/${id}`)
}

export function createCmdbHost(data: CmdbHostPayload) {
  return request.post<any, ApiResponse<CmdbHostItem>>('/cmdb/hosts', data)
}

export function updateCmdbHost(id: number, data: CmdbHostPayload) {
  return request.put<any, ApiResponse<CmdbHostItem>>(`/cmdb/hosts/${id}`, data)
}

export function deleteCmdbHost(id: number) {
  return request.delete<any, ApiResponse>(`/cmdb/hosts/${id}`)
}

export function verifyCmdbHost(id: number) {
  return request.post<any, ApiResponse<CmdbHostItem>>(`/cmdb/hosts/${id}/verify`)
}

export function importCmdbHosts(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<any, ApiResponse<CmdbHostImportResult>>('/cmdb/hosts/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function downloadCmdbHostImportTemplate() {
  return request.get('/cmdb/hosts/import-template', { responseType: 'blob' })
}
