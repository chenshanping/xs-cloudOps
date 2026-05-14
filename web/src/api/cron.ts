import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types'

export type CronTaskStatus = 0 | 1
export type CronLogStatus = 'running' | 'success' | 'failure' | 'skipped'
export type CronTriggeredBy = 'schedule' | 'manual'
export type CronParamType = 'int' | 'string' | 'bool'

export interface CronParamDefinition {
  type: CronParamType
  required: boolean
  default?: any
  description: string
  min?: number
  max?: number
}

export interface RegisteredCronTask {
  code: string
  name: string
  description: string
  param_schema: Record<string, CronParamDefinition>
}

export interface CronTask {
  id: number
  code: string
  task_code: string
  name: string
  cron_expr: string
  params: Record<string, any>
  status: CronTaskStatus
  last_run_at?: string
  last_status?: CronLogStatus | string
  last_duration_ms?: number
  next_run_at?: string
  remark?: string
  sort: number
  created_by?: number
  created_at: string
  updated_at: string
}

export interface CronLog {
  id: number
  task_id: number
  task_code: string
  started_at: string
  finished_at?: string
  duration_ms: number
  status: CronLogStatus
  summary?: string
  error_message?: string
  triggered_by: CronTriggeredBy | string
  actor_user_id?: number
  created_at: string
}

export interface CronTaskPayload {
  code: string
  task_code: string
  name: string
  cron_expr: string
  params: Record<string, any>
  remark?: string
  sort?: number
}

export interface RunNowResult {
  log_id: number
}

export function getCronTaskList(params: Record<string, any>) {
  return request.get<any, ApiResponse<PageResponse<CronTask>>>('/monitor/cron-task', { params })
}

export function createCronTask(data: CronTaskPayload) {
  return request.post<any, ApiResponse>('/monitor/cron-task', data)
}

export function updateCronTask(id: number, data: CronTaskPayload) {
  return request.put<any, ApiResponse>(`/monitor/cron-task/${id}`, data)
}

export function deleteCronTask(id: number) {
  return request.delete<any, ApiResponse>(`/monitor/cron-task/${id}`)
}

export function enableCronTask(id: number) {
  return request.post<any, ApiResponse>(`/monitor/cron-task/${id}/enable`)
}

export function disableCronTask(id: number) {
  return request.post<any, ApiResponse>(`/monitor/cron-task/${id}/disable`)
}

export function runCronTaskNow(id: number) {
  return request.post<any, ApiResponse<RunNowResult>>(`/monitor/cron-task/${id}/run`)
}

export function getCronRegistry() {
  return request.get<any, ApiResponse<RegisteredCronTask[]>>('/monitor/cron-task/registry')
}

export function getCronLogList(params: Record<string, any>) {
  return request.get<any, ApiResponse<PageResponse<CronLog>>>('/monitor/cron-log', { params })
}

export function getCronLogDetail(id: number) {
  return request.get<any, ApiResponse<CronLog>>(`/monitor/cron-log/${id}`)
}
