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
