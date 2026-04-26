package model

import "time"

type SysOperationLog struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	UserID       uint      `json:"user_id" gorm:"comment:用户ID"`
	Username     string    `json:"username" gorm:"size:50;comment:用户名"`
	IP           string    `json:"ip" gorm:"size:50;comment:IP地址"`
	Method       string    `json:"method" gorm:"size:10;comment:请求方法"`
	Path         string    `json:"path" gorm:"size:200;comment:请求路径"`
	Group        string    `json:"group" gorm:"size:50;comment:路由分组"`
	Summary      string    `json:"summary" gorm:"size:100;comment:路由描述"`
	Request      string    `json:"request" gorm:"type:text;comment:请求参数"`
	Response     string    `json:"response" gorm:"type:text;comment:响应结果"`
	Status       int       `json:"status" gorm:"comment:HTTP状态码"`
	BusinessCode int       `json:"business_code" gorm:"comment:业务状态码"`
	Latency      int64     `json:"latency" gorm:"comment:耗时(ms)"`
	UserAgent    string    `json:"user_agent" gorm:"size:500;comment:User-Agent"`
	CreatedAt    time.Time `json:"created_at"`
}

func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}
