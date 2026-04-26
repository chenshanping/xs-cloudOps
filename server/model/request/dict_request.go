package request

// 字典类型列表请求
type DictTypeListRequest struct {
	PageRequest
	Name   string `json:"name" form:"name" comment:"字典名称"`
	Type   string `json:"type" form:"type" comment:"字典类型"`
	Status *int   `json:"status" form:"status" comment:"状态"`
}

// 创建字典类型请求
type CreateDictTypeRequest struct {
	Name   string `json:"name" binding:"required" comment:"字典名称"`
	Type   string `json:"type" binding:"required" comment:"字典类型"`
	Status int    `json:"status" comment:"状态(1:正常,0:停用)"`
	Remark string `json:"remark" comment:"备注"`
}

// 更新字典类型请求
type UpdateDictTypeRequest struct {
	Name   string `json:"name" comment:"字典名称"`
	Type   string `json:"type" comment:"字典类型"`
	Status int    `json:"status" comment:"状态(1:正常,0:停用)"`
	Remark string `json:"remark" comment:"备注"`
}

// 字典数据列表请求
type DictDataListRequest struct {
	PageRequest
	DictType string `json:"dict_type" form:"dict_type" binding:"required" comment:"字典类型"`
	Label    string `json:"label" form:"label" comment:"字典标签"`
	Status   *int   `json:"status" form:"status" comment:"状态"`
}

// 创建字典数据请求
type CreateDictDataRequest struct {
	DictType  string `json:"dict_type" binding:"required" comment:"字典类型"`
	Label     string `json:"label" binding:"required" comment:"字典标签"`
	Value     string `json:"value" binding:"required" comment:"字典键值"`
	Sort      int    `json:"sort" comment:"排序"`
	Status    int    `json:"status" comment:"状态(1:正常,0:停用)"`
	TagType   string `json:"tag_type" comment:"标签类型"`
	IsDefault int    `json:"is_default" comment:"是否默认"`
	Remark    string `json:"remark" comment:"备注"`
}

// 更新字典数据请求
type UpdateDictDataRequest struct {
	DictType  string `json:"dict_type" comment:"字典类型"`
	Label     string `json:"label" comment:"字典标签"`
	Value     string `json:"value" comment:"字典键值"`
	Sort      int    `json:"sort" comment:"排序"`
	Status    int    `json:"status" comment:"状态(1:正常,0:停用)"`
	TagType   string `json:"tag_type" comment:"标签类型"`
	IsDefault int    `json:"is_default" comment:"是否默认"`
	Remark    string `json:"remark" comment:"备注"`
}
