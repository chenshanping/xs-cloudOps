package request

// ProductQueryRequest 产品信息查询请求（用于导出）
type ProductQueryRequest struct {
	Name *string `json:"name" form:"name" comment:"产品名称"`
	Num *int `json:"num" form:"num" comment:"产品数量"`
	TypeId *uint `json:"type_id" form:"type_id" comment:"产品类型ID"`
}

// ProductListRequest 产品信息列表请求
type ProductListRequest struct {
	PageRequest
	Name *string `json:"name" form:"name" comment:"产品名称"`
	Num *int `json:"num" form:"num" comment:"产品数量"`
	TypeId *uint `json:"type_id" form:"type_id" comment:"产品类型ID"`
	SortField string `json:"sort_field" form:"sort_field" comment:"排序字段"`
	SortOrder string `json:"sort_order" form:"sort_order" comment:"排序方式 asc/desc"`
}

// CreateProductRequest 创建产品信息请求
type CreateProductRequest struct {
	Name string `json:"name" binding:"required" comment:"产品名称"`
	Num int `json:"num" comment:"产品数量"`
	Price float64 `json:"price" comment:"产品单价"`
	Status string `json:"status" binding:"required" comment:"状态"`
	TypeId uint `json:"type_id" comment:"产品类型ID"`
}

// UpdateProductRequest 更新产品信息请求
type UpdateProductRequest struct {
	Name string `json:"name" comment:"产品名称"`
	Num int `json:"num" comment:"产品数量"`
	Price float64 `json:"price" comment:"产品单价"`
	Status string `json:"status" comment:"状态"`
	TypeId uint `json:"type_id" comment:"产品类型ID"`
}

// BatchDeleteProductRequest 批量删除产品信息请求
type BatchDeleteProductRequest struct {
	Ids []uint `json:"ids" binding:"required" comment:"ID列表"`
}
