package model

type SysRoleDataScope struct {
	BaseModel
	RoleID       uint      `json:"role_id" gorm:"index:idx_role_resource,unique;comment:角色ID"`
	ResourceCode string    `json:"resource_code" gorm:"size:100;index:idx_role_resource,unique;comment:业务功能资源码"`
	DataScope    int       `json:"data_scope" gorm:"default:1;comment:数据范围 1全部 2自定义 3本部门 4本部门及下级 5仅本人"`
	Depts        []SysDept `json:"depts,omitempty" gorm:"many2many:sys_role_data_scope_dept;"`
}

func (SysRoleDataScope) TableName() string {
	return "sys_role_data_scope"
}
