package model

type SysRole struct {
	BaseModel
	Name         string    `json:"name" gorm:"size:50;comment:角色名称"`
	Code         string    `json:"code" gorm:"size:50;uniqueIndex;comment:角色编码"`
	Sort         int       `json:"sort" gorm:"default:0;comment:排序"`
	Status       int       `json:"status" gorm:"default:1;comment:状态 1启用 0禁用"`
	IsSuperAdmin bool      `json:"is_super_admin" gorm:"default:false;comment:是否显式超管"`
	DataScope    int       `json:"data_scope" gorm:"default:1;comment:数据范围 1全部 2自定义 3本部门 4本部门及下级 5仅本人"`
	Remark       string    `json:"remark" gorm:"size:255;comment:备注"`
	Menus        []SysMenu `json:"menus" gorm:"many2many:sys_role_menu;"`
	Apis         []SysApi  `json:"apis" gorm:"many2many:sys_role_api;"`
	Depts        []SysDept `json:"depts,omitempty" gorm:"many2many:sys_role_dept;"`
	Users        []SysUser `json:"users,omitempty" gorm:"many2many:sys_user_role;"`
	UserCount    int64     `json:"user_count" gorm:"column:user_count;->"`
}

const (
	DataScopeAll             = 1
	DataScopeCustom          = 2
	DataScopeDept            = 3
	DataScopeDeptAndChildren = 4
	DataScopeSelf            = 5
)

func (SysRole) TableName() string {
	return "sys_role"
}
