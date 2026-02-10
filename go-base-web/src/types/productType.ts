// 产品类型

// 产品类型响应体
export interface ProductType {
  id: number
  name?: string
  icon?: string
  created_at?: string
  updated_at?: string
}

// 创建产品类型请求体
export interface CreateProductTypeRequest {
  name?: string
  icon?: string
}

// 更新产品类型请求体
export type UpdateProductTypeRequest = Partial<CreateProductTypeRequest>

// 产品类型查询参数
export interface ProductTypeQuery {
  page?: number
  page_size?: number
  sort_field?: string
  sort_order?: 'ascend' | 'descend'
}

// 选项项（用于下拉选择）
export interface OptionItem {
  id: number
  name: string
  count?: number
}
