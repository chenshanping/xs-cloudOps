import request from '@/utils/request'
import type { ApiResponse } from '@/types'

// 图表数据项类型
export interface ChartItem {
  name: string
  value: number
}

// 趋势数据项类型
export interface TrendItem {
  date: string
  count: number
}

// 获取用户角色占比统计
export function getUserRoleStats() {
  return request.get<any, ApiResponse<ChartItem[]>>('/echart/user-role-stats', { silent: true })
}

// 获取用户状态统计
export function getUserStatusStats() {
  return request.get<any, ApiResponse<ChartItem[]>>('/echart/user-status-stats', { silent: true })
}

// 获取角色状态统计
export function getRoleStatusStats() {
  return request.get<any, ApiResponse<ChartItem[]>>('/echart/role-status-stats', { silent: true })
}

// 获取用户注册趋势（近30天）
export function getUserRegisterTrend() {
  return request.get<any, ApiResponse<TrendItem[]>>('/echart/user-register-trend', { silent: true })
}
