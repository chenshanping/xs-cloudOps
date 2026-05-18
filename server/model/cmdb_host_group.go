package model

import "time"

type CmdbHostGroup struct {
	ID        uint      `json:"id" gorm:"primaryKey;comment:主键ID"`
	Name      string    `json:"name" gorm:"size:100;uniqueIndex;not null;comment:分组名称"`
	Sort      int       `json:"sort" gorm:"default:0;not null;comment:排序值"`
	Remark    string    `json:"remark" gorm:"size:500;default:'';comment:备注"`
	Status    int       `json:"status" gorm:"default:1;not null;comment:状态:1启用,0禁用"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (CmdbHostGroup) TableName() string {
	return "cmdb_host_group"
}
