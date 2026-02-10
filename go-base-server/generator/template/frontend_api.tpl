import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types'
import type { {{.ModelName}}, Create{{.ModelName}}Request, Update{{.ModelName}}Request, {{.ModelName}}Query, OptionItem{{if .LinkToUser}}, SaveMy{{.ModelName}}Request{{end}} } from '@/types/{{.ModuleName}}'

// 获取{{.Description}}列表
export function get{{.ModelName}}List(params: {{.ModelName}}Query) {
  return request.get<any, ApiResponse<PageResponse<{{.ModelName}}>>>('/{{.RoutePath}}', { params })
}

// 获取{{.Description}}详情
export function get{{.ModelName}}(id: number) {
  return request.get<any, ApiResponse<{{.ModelName}}>>(`/{{.RoutePath}}/${id}`)
}

// 创建{{.Description}}
export function create{{.ModelName}}(data: Create{{.ModelName}}Request) {
  return request.post<any, ApiResponse<{{.ModelName}}>>('/{{.RoutePath}}', data)
}

// 更新{{.Description}}
export function update{{.ModelName}}(id: number, data: Update{{.ModelName}}Request) {
  return request.put<any, ApiResponse<{{.ModelName}}>>(`/{{.RoutePath}}/${id}`, data)
}

// 删除{{.Description}}
export function delete{{.ModelName}}(id: number) {
  return request.delete<any, ApiResponse>(`/{{.RoutePath}}/${id}`)
}

// 批量删除{{.Description}}
export function batchDelete{{.ModelName}}(ids: number[]) {
  return request.delete<any, ApiResponse>('/{{.RoutePath}}/batch', { data: { ids } })
}

// 获取{{.Description}}选项列表
export function get{{.ModelName}}Options(params?: { display_field?: string; count_table?: string; count_field?: string; exclude_deleted?: boolean; count_created_by?: number }) {
  return request.get<any, ApiResponse<OptionItem[]>>('/{{.RoutePath}}/options', { params })
}
{{- if .HasCreatedBy}}

// 获取{{.Description}}创建人选项列表
export function get{{.ModelName}}CreatorOptions() {
  return request.get<any, ApiResponse<OptionItem[]>>('/{{.RoutePath}}/creator/options')
}
{{- end}}
{{- if .HasAudit}}

// 审批{{.Description}}
export function audit{{.ModelName}}(id: number, data: { audit_status: number; audit_remark: string }) {
  return request.post<any, ApiResponse>(`/{{.RoutePath}}/${id}/audit`, data)
}
{{- end}}
{{- if .LinkToUser}}

// 获取我的{{.Description}}信息
export function getMy{{.ModelName}}() {
  return request.get<any, ApiResponse<{{.ModelName}} | null>>('/{{.RoutePath}}/my')
}

// 保存我的{{.Description}}信息
export function saveMy{{.ModelName}}(data: SaveMy{{.ModelName}}Request) {
  return request.post<any, ApiResponse>('/{{.RoutePath}}/my', data)
}
{{- end}}
{{- if .HasStats}}
{{- range $chart := .StatsCharts}}

// 获取{{$.Description}}按{{$chart.Title}}分组统计
export function get{{$.ModelName}}Stats{{$chart.Field}}() {
  return request.get<any, ApiResponse<{ group_key: any; name?: string; value: number }[]>>('/{{$.RoutePath}}/stats/{{$chart.Column}}')
}
{{- end}}
{{- if .HasStatsTrend}}

// 获取{{.Description}}趋势统计
export function get{{.ModelName}}TrendStats(days?: number) {
  return request.get<any, ApiResponse<{ date: string; value: number }[]>>('/{{.RoutePath}}/stats/trend', { params: { days } })
}
{{- end}}
{{- end}}
