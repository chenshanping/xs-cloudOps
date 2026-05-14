package model

import "time"

type CronLogStatus string

const (
	CronLogStatusRunning CronLogStatus = "running"
	CronLogStatusSuccess CronLogStatus = "success"
	CronLogStatusFailure CronLogStatus = "failure"
	CronLogStatusSkipped CronLogStatus = "skipped"
)

const (
	CronTriggeredBySchedule = "schedule"
	CronTriggeredByManual   = "manual"
)

type SysCronLog struct {
	ID           uint          `json:"id" gorm:"primarykey"`
	TaskID       uint          `json:"task_id" gorm:"not null;index:idx_sys_cron_log_task_id;comment:任务ID"`
	TaskCode     string        `json:"task_code" gorm:"size:64;not null;comment:注册任务编码"`
	StartedAt    time.Time     `json:"started_at" gorm:"not null;index:idx_sys_cron_log_started_at;comment:开始时间"`
	FinishedAt   *time.Time    `json:"finished_at" gorm:"comment:结束时间"`
	DurationMS   int64         `json:"duration_ms" gorm:"comment:耗时毫秒"`
	Status       CronLogStatus `json:"status" gorm:"size:16;not null;index:idx_sys_cron_log_status;comment:执行状态"`
	Summary      string        `json:"summary" gorm:"type:text;comment:执行摘要"`
	ErrorMessage string        `json:"error_message" gorm:"type:text;comment:错误信息"`
	TriggeredBy  string        `json:"triggered_by" gorm:"size:32;not null;comment:触发方式"`
	ActorUserID  uint          `json:"actor_user_id" gorm:"comment:手动触发用户ID"`
	CreatedAt    time.Time     `json:"created_at"`
}

func (SysCronLog) TableName() string {
	return "sys_cron_log"
}
