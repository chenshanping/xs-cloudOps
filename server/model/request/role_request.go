package request

// 创建角色请求
type CreateRoleRequest struct {
	Name         string `json:"name" binding:"required" comment:"角色名称"`
	Code         string `json:"code" binding:"required" comment:"角色编码"`
	Sort         int    `json:"sort" comment:"排序"`
	Status       int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	IsSuperAdmin bool   `json:"is_super_admin" comment:"是否显式超管"`
	DataScope    int    `json:"data_scope" comment:"数据范围"`
	Remark       string `json:"remark" comment:"备注"`
	DeptIds      []uint `json:"dept_ids" comment:"自定义数据范围部门ID列表"`
}

// 更新角色请求
type UpdateRoleRequest struct {
	Name         string `json:"name" comment:"角色名称"`
	Code         string `json:"code" comment:"角色编码"`
	Sort         int    `json:"sort" comment:"排序"`
	Status       int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	IsSuperAdmin bool   `json:"is_super_admin" comment:"是否显式超管"`
	DataScope    int    `json:"data_scope" comment:"数据范围"`
	Remark       string `json:"remark" comment:"备注"`
	DeptIds      []uint `json:"dept_ids" comment:"自定义数据范围部门ID列表"`
}

type RoleFeatureDataScopeAssignment struct {
	ResourceCode string `json:"resource_code" comment:"业务功能资源码"`
	DataScope    int    `json:"data_scope" comment:"数据范围"`
	DeptIds      []uint `json:"dept_ids" comment:"自定义数据范围部门ID列表"`
}

type AssignRoleDataScopesRequest struct {
	Scopes []RoleFeatureDataScopeAssignment `json:"scopes" comment:"角色业务功能数据范围列表"`
}

type SaveRolePermissionsRequest struct {
	MenuIds      []uint                           `json:"menu_ids" comment:"菜单ID列表"`
	DirectApiIds []uint                           `json:"direct_api_ids" comment:"直接API ID列表"`
	Scopes       []RoleFeatureDataScopeAssignment `json:"scopes" comment:"角色业务功能数据范围列表"`
}
