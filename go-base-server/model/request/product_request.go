package request

// ProductListRequest 产品信息列表请求
type ProductListRequest struct {
	PageRequest
	TypeId *uint `json:"type_id" form:"type_id" comment:"产品分类ID"`
	SortField string `json:"sort_field" form:"sort_field" comment:"排序字段"`
	SortOrder string `json:"sort_order" form:"sort_order" comment:"排序方式 asc/desc"`
}

// CreateProductRequest 创建产品信息请求
type CreateProductRequest struct {
	Name string `json:"name" comment:"产品名称"`
	TypeId uint `json:"type_id" comment:"产品分类ID"`
}

// UpdateProductRequest 更新产品信息请求
type UpdateProductRequest struct {
	Name string `json:"name" comment:"产品名称"`
	TypeId uint `json:"type_id" comment:"产品分类ID"`
}

// BatchDeleteProductRequest 批量删除产品信息请求
type BatchDeleteProductRequest struct {
	Ids []uint `json:"ids" binding:"required" comment:"ID列表"`
}
