package cronsvc

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"server/global"
	"server/model"
)

var ErrTaskAlreadyRunning = errors.New("cron_task_already_running")

type scheduledTask struct {
	task    model.SysCronTask
	entryID cron.EntryID
	mu      sync.Mutex
}

type Scheduler struct {
	mu      sync.RWMutex
	cron    *cron.Cron
	ctx     context.Context
	cancel  context.CancelFunc
	tasks   map[uint]*scheduledTask
	started bool
	wg      sync.WaitGroup
}

func NewScheduler() *Scheduler {
	return &Scheduler{tasks: make(map[uint]*scheduledTask)}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return
	}
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.cron = cron.New(cron.WithParser(cronParser), cron.WithChain(cron.Recover(cron.DefaultLogger)))
	s.cron.Start()
	s.started = true
}

func (s *Scheduler) Stop(timeout time.Duration) {
	s.mu.Lock()
	if !s.started || s.cron == nil {
		s.mu.Unlock()
		return
	}
	if s.cancel != nil {
		s.cancel()
	}
	ctx := s.cron.Stop()
	s.started = false
	s.mu.Unlock()

	stopCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	case <-ctx.Done():
	case <-stopCtx.Done():
		return
	}
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-stopCtx.Done():
	}
}

func (s *Scheduler) AddTask(task model.SysCronTask) error {
	if task.Status != model.CronTaskStatusEnabled {
		return nil
	}
	handler, ok := Get(task.TaskCode)
	if !ok {
		return errors.New("cron_task_code_not_registered")
	}
	if _, err := ValidateCronExpr(task.CronExpr); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.started || s.cron == nil {
		return errors.New("cron_scheduler_not_started")
	}
	if current, ok := s.tasks[task.ID]; ok {
		s.cron.Remove(current.entryID)
	}
	item := &scheduledTask{task: task}
	entryID, err := s.cron.AddFunc(task.CronExpr, func() {
		s.execute(item, handler, model.CronTriggeredBySchedule, 0)
	})
	if err != nil {
		return err
	}
	item.entryID = entryID
	s.tasks[task.ID] = item
	return nil
}

func (s *Scheduler) RemoveTask(taskID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if current, ok := s.tasks[taskID]; ok {
		if s.cron != nil {
			s.cron.Remove(current.entryID)
		}
		delete(s.tasks, taskID)
	}
}

func (s *Scheduler) UpdateTask(task model.SysCronTask) error {
	if task.Status != model.CronTaskStatusEnabled {
		s.RemoveTask(task.ID)
		return nil
	}
	handler, ok := Get(task.TaskCode)
	if !ok {
		return errors.New("cron_task_code_not_registered")
	}
	if _, err := ValidateCronExpr(task.CronExpr); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.started || s.cron == nil {
		return errors.New("cron_scheduler_not_started")
	}
	item, ok := s.tasks[task.ID]
	if !ok {
		item = &scheduledTask{}
		s.tasks[task.ID] = item
	} else if item.entryID != 0 {
		s.cron.Remove(item.entryID)
	}
	item.task = task
	entryID, err := s.cron.AddFunc(task.CronExpr, func() {
		s.execute(item, handler, model.CronTriggeredBySchedule, 0)
	})
	if err != nil {
		return err
	}
	item.entryID = entryID
	return nil
}

func (s *Scheduler) RunNow(task model.SysCronTask, actorID uint) (uint, error) {
	handler, ok := Get(task.TaskCode)
	if !ok {
		return 0, errors.New("cron_task_code_not_registered")
	}
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return 0, errors.New("cron_scheduler_not_started")
	}
	item, ok := s.tasks[task.ID]
	if ok {
		item.task = task
	} else {
		item = &scheduledTask{task: task}
		s.tasks[task.ID] = item
	}
	s.wg.Add(1)
	s.mu.Unlock()
	if !item.mu.TryLock() {
		s.wg.Done()
		return 0, ErrTaskAlreadyRunning
	}
	logID, err := StartLog(task.ID, task.TaskCode, model.CronTriggeredByManual, actorID)
	if err != nil {
		item.mu.Unlock()
		s.wg.Done()
		return 0, err
	}
	go func() {
		defer s.wg.Done()
		s.executeLocked(item, handler, logID)
	}()
	return logID, nil
}

func (s *Scheduler) getOrCreateTaskItem(task model.SysCronTask) *scheduledTask {
	s.mu.Lock()
	defer s.mu.Unlock()
	if item, ok := s.tasks[task.ID]; ok {
		item.task = task
		return item
	}
	item := &scheduledTask{task: task}
	s.tasks[task.ID] = item
	return item
}

func (s *Scheduler) execute(item *scheduledTask, handler TaskHandler, triggeredBy string, actorID uint) {
	if !item.mu.TryLock() {
		logID, err := StartLog(item.task.ID, item.task.TaskCode, triggeredBy, actorID)
		if err == nil {
			_ = FinishLog(logID, model.CronLogStatusSkipped, "上一次执行尚未结束，本次触发已跳过", nil)
		}
		return
	}
	logID, err := StartLog(item.task.ID, item.task.TaskCode, triggeredBy, actorID)
	if err != nil {
		item.mu.Unlock()
		global.Log.Errorf("创建定时任务执行日志失败(task=%d): %v", item.task.ID, err)
		return
	}
	s.executeLocked(item, handler, logID)
}

func (s *Scheduler) executeLocked(item *scheduledTask, handler TaskHandler, logID uint) {
	defer item.mu.Unlock()
	status := model.CronLogStatusSuccess
	summary := ""
	var runErr error
	defer func() {
		if recovered := recover(); recovered != nil {
			status = model.CronLogStatusFailure
			runErr = panicToError(recovered)
		}
		if err := FinishLog(logID, status, summary, runErr); err != nil {
			global.Log.Errorf("更新定时任务执行日志失败(log=%d): %v", logID, err)
		}
		updateTaskRunState(item.task.ID, logID, status)
	}()

	params, err := ParamsToMap(item.task.Params)
	if err != nil {
		status = model.CronLogStatusFailure
		runErr = err
		return
	}
	ctx := s.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	summary, runErr = handler.Run(ctx, params)
	if runErr != nil {
		status = model.CronLogStatusFailure
	}
}

func updateTaskRunState(taskID uint, logID uint, status model.CronLogStatus) {
	now := time.Now()
	updates := map[string]interface{}{
		"last_run_at":      now,
		"last_status":      string(status),
		"next_run_at":      nil,
		"last_duration_ms": 0,
	}
	var currentLog model.SysCronLog
	if err := global.DB.First(&currentLog, logID).Error; err == nil {
		updates["last_duration_ms"] = currentLog.DurationMS
	}
	var task model.SysCronTask
	if err := global.DB.First(&task, taskID).Error; err == nil && task.Status == model.CronTaskStatusEnabled {
		updates["next_run_at"] = nextRunAt(task.CronExpr)
	}
	if err := global.DB.Model(&model.SysCronTask{}).Where("id = ?", taskID).Updates(updates).Error; err != nil {
		global.Log.Errorf("更新定时任务状态失败(task=%d): %v", taskID, err)
	}
}
