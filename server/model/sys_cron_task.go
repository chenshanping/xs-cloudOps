package model

import (
	"time"

	"gorm.io/datatypes"
)

type CronTaskStatus int

const (
	CronTaskStatusDisabled CronTaskStatus = 0
	CronTaskStatusEnabled  CronTaskStatus = 1
)

type SysCronTask struct {
	BaseModel
	Code           string         `json:"code" gorm:"size:64;not null;index:idx_sys_cron_task_code;comment:任务实例编码"`
	TaskCode       string         `json:"task_code" gorm:"size:64;not null;index:idx_sys_cron_task_task_code;comment:注册任务编码"`
	Name           string         `json:"name" gorm:"size:128;not null;comment:任务名称"`
	CronExpr       string         `json:"cron_expr" gorm:"size:64;not null;comment:Cron表达式"`
	Params         datatypes.JSON `json:"params" gorm:"type:json;comment:任务参数"`
	Status         CronTaskStatus `json:"status" gorm:"not null;default:0;index:idx_sys_cron_task_status;comment:状态 0禁用 1启用"`
	LastRunAt      *time.Time     `json:"last_run_at" gorm:"comment:上次执行时间"`
	LastStatus     string         `json:"last_status" gorm:"size:16;comment:上次执行状态"`
	LastDurationMS int64          `json:"last_duration_ms" gorm:"comment:上次执行耗时毫秒"`
	NextRunAt      *time.Time     `json:"next_run_at" gorm:"comment:下次执行时间"`
	Remark         string         `json:"remark" gorm:"size:255;comment:备注"`
	Sort           int            `json:"sort" gorm:"not null;default:0;comment:排序"`
	CreatedBy      uint           `json:"created_by" gorm:"comment:创建人ID"`
}

func (SysCronTask) TableName() string {
	return "sys_cron_task"
}
