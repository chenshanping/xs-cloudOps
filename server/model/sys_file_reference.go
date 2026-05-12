package model

import "time"

// SysFileReference 记录业务数据与文件的活动引用关系。
// 该表只保存“当前仍在使用”的绑定关系，不使用软删除。
type SysFileReference struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FileID    uint      `json:"file_id" gorm:"index;uniqueIndex:uk_sys_file_reference,priority:4;comment:文件ID"`
	RefTable  string    `json:"ref_table" gorm:"size:100;index:idx_ref_table_id_field,priority:1;uniqueIndex:uk_sys_file_reference,priority:2;comment:引用表名"`
	RefID     uint      `json:"ref_id" gorm:"index:idx_ref_table_id_field,priority:2;uniqueIndex:uk_sys_file_reference,priority:3;comment:引用记录ID"`
	RefField  string    `json:"ref_field" gorm:"size:100;index:idx_ref_table_id_field,priority:3;uniqueIndex:uk_sys_file_reference,priority:1;comment:引用字段"`
}

func (SysFileReference) TableName() string {
	return "sys_file_reference"
}
