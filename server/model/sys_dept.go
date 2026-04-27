package model

type SysDept struct {
	BaseModel
	ParentID        uint      `json:"parent_id" gorm:"default:0;index;comment:父部门ID"`
	Ancestors       string    `json:"ancestors" gorm:"size:500;comment:祖级列表"`
	Name            string    `json:"name" gorm:"size:100;comment:部门名称"`
	Sort            int       `json:"sort" gorm:"default:0;comment:排序"`
	Status          int       `json:"status" gorm:"default:1;comment:状态 1启用 0禁用"`
	Remark          string    `json:"remark" gorm:"size:255;comment:备注"`
	DirectUserCount int64     `json:"direct_user_count" gorm:"-"`
	TotalUserCount  int64     `json:"total_user_count" gorm:"-"`
	HasChildren     bool      `json:"has_children" gorm:"-"`
	Bindable        bool      `json:"bindable" gorm:"-"`
	Manageable      bool      `json:"manageable" gorm:"-"`
	Children        []SysDept `json:"children" gorm:"-"`
}

type ManageableDeptTreeResponse struct {
	Tree                []SysDept `json:"tree"`
	UnassignedUserCount int64     `json:"unassigned_user_count"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}
