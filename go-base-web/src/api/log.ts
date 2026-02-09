import request from '@/utils/request'
import type { ApiResponse, PageResponse, OperationLog, LoginLog, SlowLog } from '@/types'

// 获取操作日志列表
export function getOperationLogList(params: any) {
  return request.get<any, ApiResponse<PageResponse<OperationLog>>>('/logs/operation', { params })
}

// 获取登录日志列表
export function getLoginLogList(params: any) {
  return request.get<any, ApiResponse<PageResponse<LoginLog>>>('/logs/login', { params })
}

// 获取慢查询日志列表
export function getSlowLogList(params: any) {
  return request.get<any, ApiResponse<PageResponse<SlowLog>>>('/logs/slow', { params })
}

// 获取路由分组列表
export function getRouteGroups() {
  return request.get<any, ApiResponse<{ group: string; count: number }[]>>('/logs/route-groups')
}
