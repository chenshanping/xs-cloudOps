package cronsvc

import (
	"server/global"
	"server/model"
)

func (s *Scheduler) LoadFromDB() error {
	markRunningLogsFailed()
	var tasks []model.SysCronTask
	if err := global.DB.Where("status = ?", model.CronTaskStatusEnabled).Find(&tasks).Error; err != nil {
		return err
	}
	for _, task := range tasks {
		if err := s.AddTask(task); err != nil {
			global.Log.Errorf("加载定时任务失败(task=%d, code=%s): %v", task.ID, task.TaskCode, err)
		}
	}
	return nil
}
