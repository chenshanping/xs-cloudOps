package request

// ProductTypeQueryRequest 产品类型查询请求（用于导出）
type ProductTypeQueryRequest struct {
}

// ProductTypeListRequest 产品类型列表请求
type ProductTypeListRequest struct {
	PageRequest
	SortField string `json:"sort_field" form:"sort_field" comment:"排序字段"`
	SortOrder string `json:"sort_order" form:"sort_order" comment:"排序方式 asc/desc"`
}

// CreateProductTypeRequest 创建产品类型请求
type CreateProductTypeRequest struct {
	Name string `json:"name" comment:"产品类型名称"`
	Icon string `json:"icon" comment:"类型图标"`
	Status string `json:"status"`
}

// UpdateProductTypeRequest 更新产品类型请求
type UpdateProductTypeRequest struct {
	Name string `json:"name" comment:"产品类型名称"`
	Icon string `json:"icon" comment:"类型图标"`
	Status string `json:"status"`
}

// BatchDeleteProductTypeRequest 批量删除产品类型请求
type BatchDeleteProductTypeRequest struct {
	Ids []uint `json:"ids" binding:"required" comment:"ID列表"`
}
