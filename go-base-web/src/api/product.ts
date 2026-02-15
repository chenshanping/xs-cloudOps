import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types'
import type { Product, CreateProductRequest, UpdateProductRequest, ProductQuery, OptionItem } from '@/types/product'

// 获取产品信息列表
export function getProductList(params: ProductQuery) {
  return request.get<any, ApiResponse<PageResponse<Product>>>('/product', { params })
}

// 获取产品信息详情
export function getProduct(id: number) {
  return request.get<any, ApiResponse<Product>>(`/product/${id}`)
}

// 创建产品信息
export function createProduct(data: CreateProductRequest) {
  return request.post<any, ApiResponse<Product>>('/product', data)
}

// 更新产品信息
export function updateProduct(id: number, data: UpdateProductRequest) {
  return request.put<any, ApiResponse<Product>>(`/product/${id}`, data)
}

// 删除产品信息
export function deleteProduct(id: number) {
  return request.delete<any, ApiResponse>(`/product/${id}`)
}

// 批量删除产品信息
export function batchDeleteProduct(ids: number[]) {
  return request.delete<any, ApiResponse>('/product/batch', { data: { ids } })
}

// 获取产品信息选项列表
export function getProductOptions(params?: { display_field?: string; count_table?: string; count_field?: string; exclude_deleted?: boolean; count_created_by?: number }) {
  return request.get<any, ApiResponse<OptionItem[]>>('/product/options', { params })
}

// 获取产品信息按产品类型分组统计
export function getProductStatsTypeId() {
  return request.get<any, ApiResponse<{ group_key: any; name?: string; value: number }[]>>('/product/stats/type_id')
}

// 获取产品信息按产品状态分组统计
export function getProductStatsStatus() {
  return request.get<any, ApiResponse<{ group_key: any; name?: string; value: number }[]>>('/product/stats/status')
}

// 获取产品信息趋势统计
export function getProductTrendStats(days?: number) {
  return request.get<any, ApiResponse<{ date: string; value: number }[]>>('/product/stats/trend', { params: { days } })
}

// 导出产品信息
export function exportProduct(params?: ProductQuery) {
  return request.get('/product/export', { 
    params,
    responseType: 'blob'
  })
}

// 导入产品信息
export function importProduct(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<any, ApiResponse<{
    success_count: number
    fail_count: number
    total: number
    errors: string[]
  }>>('/product/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 下载导入模板
export function downloadTemplateProduct() {
  return request.get('/product/template', {
    responseType: 'blob'
  })
}
