package model

import "time"

type CmdbHostTagRel struct {
	ID        uint      `json:"id" gorm:"primaryKey;comment:主键ID"`
	HostID    uint      `json:"host_id" gorm:"not null;uniqueIndex:uk_cmdb_host_tag_rel;comment:主机ID"`
	TagID     uint      `json:"tag_id" gorm:"not null;uniqueIndex:uk_cmdb_host_tag_rel;index:idx_cmdb_host_tag_rel_tag_id;comment:标签ID"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
}

func (CmdbHostTagRel) TableName() string {
	return "cmdb_host_tag_rel"
}
