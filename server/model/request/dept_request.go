package request

// 创建部门请求
type CreateDeptRequest struct {
	ParentID uint   `json:"parent_id" comment:"父部门ID"`
	Name     string `json:"name" binding:"required" comment:"部门名称"`
	Sort     int    `json:"sort" comment:"排序"`
	Status   int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	Remark   string `json:"remark" comment:"备注"`
}

// 更新部门请求
type UpdateDeptRequest struct {
	ParentID uint   `json:"parent_id" comment:"父部门ID"`
	Name     string `json:"name" binding:"required" comment:"部门名称"`
	Sort     int    `json:"sort" comment:"排序"`
	Status   int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	Remark   string `json:"remark" comment:"备注"`
}
