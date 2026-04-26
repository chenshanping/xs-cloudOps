package model

type SysMenu struct {
	BaseModel
	ParentID   uint      `json:"parent_id" gorm:"default:0;comment:父菜单ID"`
	Name       string    `json:"name" gorm:"size:50;comment:菜单名称"`
	Path       string    `json:"path" gorm:"size:200;comment:路由路径"`
	Component  string    `json:"component" gorm:"size:200;comment:组件路径"`
	Icon       string    `json:"icon" gorm:"size:50;comment:图标"`
	Sort       int       `json:"sort" gorm:"default:0;comment:排序"`
	Type       int       `json:"type" gorm:"default:1;comment:类型 1目录 2菜单 3按钮"`
	Permission string    `json:"permission" gorm:"size:100;comment:权限标识"`
	Status     int       `json:"status" gorm:"default:1;comment:状态 1启用 0禁用"`
	Hidden     int       `json:"hidden" gorm:"default:0;comment:是否隐藏 0显示 1隐藏"`
	Children   []SysMenu `json:"children" gorm:"-"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
