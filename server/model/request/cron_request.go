package request

import "encoding/json"

type CronTaskListRequest struct {
	PageRequest
	Code      string `json:"code" form:"code"`
	TaskCode  string `json:"task_code" form:"task_code"`
	Name      string `json:"name" form:"name"`
	Status    *int   `json:"status" form:"status"`
	SortField string `json:"sort_field" form:"sort_field"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}

type CronTaskSaveRequest struct {
	Code     string          `json:"code" binding:"required"`
	TaskCode string          `json:"task_code" binding:"required"`
	Name     string          `json:"name" binding:"required"`
	CronExpr string          `json:"cron_expr" binding:"required"`
	Params   json.RawMessage `json:"params"`
	Remark   string          `json:"remark"`
	Sort     int             `json:"sort"`
}

type CronLogListRequest struct {
	PageRequest
	TaskID      uint   `json:"task_id" form:"task_id"`
	TaskCode    string `json:"task_code" form:"task_code"`
	Status      string `json:"status" form:"status"`
	TriggeredBy string `json:"triggered_by" form:"triggered_by"`
	StartTime   string `json:"start_time" form:"start_time"`
	EndTime     string `json:"end_time" form:"end_time"`
}
