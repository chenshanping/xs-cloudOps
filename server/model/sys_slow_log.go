package model

import "time"

type SysSlowLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	SQL       string    `json:"sql" gorm:"type:text;comment:SQL语句"`
	Rows      int64     `json:"rows" gorm:"comment:影响行数"`
	Latency   float64   `json:"latency" gorm:"comment:耗时(ms)"`
	Source    string    `json:"source" gorm:"size:500;comment:调用来源"`
	CreatedAt time.Time `json:"created_at"`
}

func (SysSlowLog) TableName() string {
	return "sys_slow_log"
}
