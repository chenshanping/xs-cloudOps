// 产品信息
import type { ProductType } from './productType'

// 产品信息响应体
export interface Product {
  id: number
  type_id: number
  name: string
  num?: number
  price?: number
  status: string
  product_type?: ProductType
  created_at?: string
  updated_at?: string
}

// 创建产品信息请求体
export interface CreateProductRequest {
  name: string
  num?: number
  price?: number
  status: string
  type_id?: number
}

// 更新产品信息请求体
export type UpdateProductRequest = Partial<CreateProductRequest>

// 产品信息查询参数
export interface ProductQuery {
  page?: number
  page_size?: number
  name?: string
  num?: number
  type_id?: number
  sort_field?: string
  sort_order?: 'ascend' | 'descend'
}

// 选项项（用于下拉选择）
export interface OptionItem {
  id: number
  name: string
  count?: number
}
