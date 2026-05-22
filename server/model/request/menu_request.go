package request

// 分配菜单请求
type AssignMenusRequest struct {
	MenuIds []uint `json:"menu_ids" comment:"菜单ID列表"`
}

// 分配API请求
type AssignApisRequest struct {
	ApiIds []uint `json:"api_ids" comment:"API ID列表"`
}

type UpdateMenuApisRequest struct {
	ApiIds []uint `json:"api_ids" comment:"菜单关联API ID列表"`
}

// 创建菜单请求
type CreateMenuRequest struct {
	ParentID     uint   `json:"parent_id" comment:"父菜单ID"`
	Name         string `json:"name" binding:"required" comment:"菜单名称"`
	Path         string `json:"path" comment:"路由路径"`
	Component    string `json:"component" comment:"组件路径"`
	Icon         string `json:"icon" comment:"图标"`
	Sort         int    `json:"sort" comment:"排序"`
	Type         int    `json:"type" binding:"required,oneof=1 2 3" comment:"类型(1:目录,2:菜单,3:按钮)"`
	Permission   string `json:"permission" comment:"权限标识"`
	Status       int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	Hidden       int    `json:"hidden" comment:"是否隐藏(0:显示,1:隐藏)"`
	IsStandalone int    `json:"is_standalone" binding:"oneof=0 1" comment:"是否独立页(0:否,1:是)"`
}

// 更新菜单请求
type UpdateMenuRequest struct {
	ParentID     uint   `json:"parent_id" comment:"父菜单ID"`
	Name         string `json:"name" comment:"菜单名称"`
	Path         string `json:"path" comment:"路由路径"`
	Component    string `json:"component" comment:"组件路径"`
	Icon         string `json:"icon" comment:"图标"`
	Sort         int    `json:"sort" comment:"排序"`
	Type         int    `json:"type" comment:"类型(1:目录,2:菜单,3:按钮)"`
	Permission   string `json:"permission" comment:"权限标识"`
	Status       int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	Hidden       int    `json:"hidden" comment:"是否隐藏(0:显示,1:隐藏)"`
	IsStandalone int    `json:"is_standalone" binding:"oneof=0 1" comment:"是否独立页(0:否,1:是)"`
}
