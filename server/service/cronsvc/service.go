package cronsvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
)

type CronTaskService struct {
	Scheduler *Scheduler
}

var Default = &CronTaskService{Scheduler: NewScheduler()}

var cronTaskSortColumns = map[string]string{
	"id":               "id",
	"code":             "code",
	"task_code":        "task_code",
	"name":             "name",
	"status":           "status",
	"sort":             "sort",
	"last_run_at":      "last_run_at",
	"last_duration_ms": "last_duration_ms",
	"next_run_at":      "next_run_at",
	"created_at":       "created_at",
	"updated_at":       "updated_at",
}

func (s *CronTaskService) ListTasks(req *request.CronTaskListRequest) ([]model.SysCronTask, int64, error) {
	var tasks []model.SysCronTask
	var total int64
	db := global.DB.Model(&model.SysCronTask{})
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.TaskCode != "" {
		db = db.Where("task_code = ?", req.TaskCode)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	orderStr := "sort ASC, id DESC"
	if sortField, ok := cronTaskSortColumns[req.SortField]; ok {
		order := "ASC"
		if req.SortOrder == "descend" {
			order = "DESC"
		}
		orderStr = sortField + " " + order
	}
	if err := db.Order(orderStr).Offset(req.GetOffset()).Limit(req.PageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (s *CronTaskService) GetTask(id uint) (*model.SysCronTask, error) {
	var task model.SysCronTask
	if err := global.DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *CronTaskService) CreateTask(req *request.CronTaskSaveRequest, actorID uint) error {
	task, err := s.buildTaskFromRequest(req, actorID)
	if err != nil {
		return err
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.ensureCodeAvailable(tx, task.Code, 0); err != nil {
			return err
		}
		return tx.Create(task).Error
	})
}

func (s *CronTaskService) UpdateTask(id uint, req *request.CronTaskSaveRequest) error {
	var existing model.SysCronTask
	if err := global.DB.First(&existing, id).Error; err != nil {
		return err
	}
	task, err := s.buildTaskFromRequest(req, existing.CreatedBy)
	if err != nil {
		return err
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.ensureCodeAvailable(tx, task.Code, id); err != nil {
			return err
		}
		updates := map[string]interface{}{
			"code":        task.Code,
			"task_code":   task.TaskCode,
			"name":        task.Name,
			"cron_expr":   task.CronExpr,
			"params":      task.Params,
			"remark":      task.Remark,
			"sort":        task.Sort,
			"next_run_at": task.NextRunAt,
		}
		if err := tx.Model(&existing).Updates(updates).Error; err != nil {
			return err
		}
		if existing.Status == model.CronTaskStatusEnabled {
			var updated model.SysCronTask
			if err := tx.First(&updated, id).Error; err != nil {
				return err
			}
			return s.Scheduler.UpdateTask(updated)
		}
		return nil
	})
}

func (s *CronTaskService) DeleteTask(id uint) error {
	var task model.SysCronTask
	if err := global.DB.First(&task, id).Error; err != nil {
		return err
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&task).Error; err != nil {
			return err
		}
		s.Scheduler.RemoveTask(id)
		return nil
	})
}

func (s *CronTaskService) EnableTask(id uint) error {
	var task model.SysCronTask
	if err := global.DB.First(&task, id).Error; err != nil {
		return err
	}
	if _, ok := Get(task.TaskCode); !ok {
		return errors.New("cron_task_code_not_registered")
	}
	next, err := ValidateCronExpr(task.CronExpr)
	if err != nil {
		return err
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&task).Updates(map[string]interface{}{
			"status":      model.CronTaskStatusEnabled,
			"next_run_at": next,
		}).Error; err != nil {
			return err
		}
		task.Status = model.CronTaskStatusEnabled
		task.NextRunAt = next
		return s.Scheduler.UpdateTask(task)
	})
}

func (s *CronTaskService) DisableTask(id uint) error {
	var task model.SysCronTask
	if err := global.DB.First(&task, id).Error; err != nil {
		return err
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&task).Updates(map[string]interface{}{
			"status":      model.CronTaskStatusDisabled,
			"next_run_at": nil,
		}).Error; err != nil {
			return err
		}
		s.Scheduler.RemoveTask(id)
		return nil
	})
}

func (s *CronTaskService) RunNow(id uint, actorID uint) (uint, error) {
	var task model.SysCronTask
	if err := global.DB.First(&task, id).Error; err != nil {
		return 0, err
	}
	return s.Scheduler.RunNow(task, actorID)
}

func (s *CronTaskService) ListLogs(req *request.CronLogListRequest) ([]model.SysCronLog, int64, error) {
	var logs []model.SysCronLog
	var total int64
	db := global.DB.Model(&model.SysCronLog{})
	if req.TaskID > 0 {
		db = db.Where("task_id = ?", req.TaskID)
	}
	if req.TaskCode != "" {
		db = db.Where("task_code = ?", req.TaskCode)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.TriggeredBy != "" {
		db = db.Where("triggered_by = ?", req.TriggeredBy)
	}
	if req.StartTime != "" {
		db = db.Where("started_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		db = db.Where("started_at <= ?", req.EndTime)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset(req.GetOffset()).Limit(req.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (s *CronTaskService) GetLog(id uint) (*model.SysCronLog, error) {
	var log model.SysCronLog
	if err := global.DB.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (s *CronTaskService) Registry() []RegisteredTask {
	return ListRegisteredTasks()
}

func (s *CronTaskService) buildTaskFromRequest(req *request.CronTaskSaveRequest, actorID uint) (*model.SysCronTask, error) {
	code := strings.TrimSpace(req.Code)
	taskCode := strings.TrimSpace(req.TaskCode)
	name := strings.TrimSpace(req.Name)
	cronExpr := strings.TrimSpace(req.CronExpr)
	if code == "" || taskCode == "" || name == "" || cronExpr == "" {
		return nil, errors.New("任务编码、任务类型、任务名称和Cron表达式不能为空")
	}
	handler, ok := Get(taskCode)
	if !ok {
		return nil, errors.New("cron_task_code_not_registered")
	}
	next, err := ValidateCronExpr(cronExpr)
	if err != nil {
		return nil, err
	}
	_, cleanedParams, err := ValidateParams(req.Params, handler.ParamSchema)
	if err != nil {
		return nil, err
	}
	return &model.SysCronTask{
		Code:      code,
		TaskCode:  taskCode,
		Name:      name,
		CronExpr:  cronExpr,
		Params:    cleanedParams,
		Status:    model.CronTaskStatusDisabled,
		NextRunAt: next,
		Remark:    strings.TrimSpace(req.Remark),
		Sort:      req.Sort,
		CreatedBy: actorID,
	}, nil
}

func (s *CronTaskService) ensureCodeAvailable(tx *gorm.DB, code string, excludeID uint) error {
	var count int64
	db := tx.Model(&model.SysCronTask{}).Where("code = ?", code)
	if excludeID > 0 {
		db = db.Where("id <> ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("任务实例编码已存在")
	}
	return nil
}

func SeedTaskParams(retainDays, batchLimit int) []byte {
	encoded, _ := json.Marshal(map[string]interface{}{
		"retain_days": retainDays,
		"batch_limit": batchLimit,
	})
	return encoded
}
