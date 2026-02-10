import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types'
import type { ProductType, CreateProductTypeRequest, UpdateProductTypeRequest, ProductTypeQuery, OptionItem } from '@/types/productType'

// 获取产品类型列表
export function getProductTypeList(params: ProductTypeQuery) {
  return request.get<any, ApiResponse<PageResponse<ProductType>>>('/productType', { params })
}

// 获取产品类型详情
export function getProductType(id: number) {
  return request.get<any, ApiResponse<ProductType>>(`/productType/${id}`)
}

// 创建产品类型
export function createProductType(data: CreateProductTypeRequest) {
  return request.post<any, ApiResponse<ProductType>>('/productType', data)
}

// 更新产品类型
export function updateProductType(id: number, data: UpdateProductTypeRequest) {
  return request.put<any, ApiResponse<ProductType>>(`/productType/${id}`, data)
}

// 删除产品类型
export function deleteProductType(id: number) {
  return request.delete<any, ApiResponse>(`/productType/${id}`)
}

// 批量删除产品类型
export function batchDeleteProductType(ids: number[]) {
  return request.delete<any, ApiResponse>('/productType/batch', { data: { ids } })
}

// 获取产品类型选项列表
export function getProductTypeOptions(params?: { display_field?: string; count_table?: string; count_field?: string; exclude_deleted?: boolean; count_created_by?: number }) {
  return request.get<any, ApiResponse<OptionItem[]>>('/productType/options', { params })
}
