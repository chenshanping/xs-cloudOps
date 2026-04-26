package model

type SysRole struct {
	BaseModel
	Name   string    `json:"name" gorm:"size:50;comment:角色名称"`
	Code   string    `json:"code" gorm:"size:50;uniqueIndex;comment:角色编码"`
	Sort   int       `json:"sort" gorm:"default:0;comment:排序"`
	Status int       `json:"status" gorm:"default:1;comment:状态 1启用 0禁用"`
	Remark string    `json:"remark" gorm:"size:255;comment:备注"`
	Menus  []SysMenu `json:"menus" gorm:"many2many:sys_role_menu;"`
	Apis   []SysApi  `json:"apis" gorm:"many2many:sys_role_api;"`
}

func (SysRole) TableName() string {
	return "sys_role"
}
