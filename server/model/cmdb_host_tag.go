package model

import "time"

type CmdbHostTag struct {
	ID        uint      `json:"id" gorm:"primaryKey;comment:主键ID"`
	Name      string    `json:"name" gorm:"size:100;uniqueIndex;not null;comment:标签名称"`
	Color     string    `json:"color" gorm:"size:30;default:'';comment:标签颜色"`
	Remark    string    `json:"remark" gorm:"size:500;default:'';comment:备注"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (CmdbHostTag) TableName() string {
	return "cmdb_host_tag"
}
