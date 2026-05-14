package cronsvc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
)

func setupCronServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if err := db.AutoMigrate(&model.SysCronTask{}, &model.SysCronLog{}, &model.SysLoginLog{}, &model.SysOperationLog{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	previousDB := global.DB
	previousLog := global.Log
	global.DB = db
	global.Log = zap.NewNop().Sugar()
	t.Cleanup(func() {
		global.DB = previousDB
		global.Log = previousLog
	})
	return db
}

func TestValidateParamsStripsUnknownAndRejectsWrongType(t *testing.T) {
	handler, ok := Get("cleanup_login_logs")
	if !ok {
		t.Fatal("cleanup_login_logs not registered")
	}
	cleaned, _, err := ValidateParams([]byte(`{"retain_days":30,"batch_limit":1000,"evil":"x"}`), handler.ParamSchema)
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := cleaned["evil"]; exists {
		t.Fatal("unknown param was not stripped")
	}
	if cleaned["retain_days"] != 30 || cleaned["batch_limit"] != 1000 {
		t.Fatalf("unexpected cleaned params: %#v", cleaned)
	}
	if _, _, err := ValidateParams([]byte(`{"retain_days":"bad","batch_limit":1000}`), handler.ParamSchema); err == nil {
		t.Fatal("expected wrong param type to fail")
	}
}

func TestValidateCronExpr(t *testing.T) {
	if next, err := ValidateCronExpr("0 2 * * *"); err != nil || next == nil {
		t.Fatalf("valid expr rejected: next=%v err=%v", next, err)
	}
	if _, err := ValidateCronExpr("@@@@"); err == nil || err.Error() != "cron_expr_invalid" {
		t.Fatalf("invalid expr error = %v", err)
	}
}

func TestCreateTaskRejectsUnregisteredTaskCode(t *testing.T) {
	setupCronServiceTestDB(t)
	svc := &CronTaskService{Scheduler: NewScheduler()}
	req := requestLikeTask("bad", "rm_rf_root")
	err := svc.CreateTask(&req, 1)
	if err == nil || err.Error() != "cron_task_code_not_registered" {
		t.Fatalf("expected unregistered task code error, got %v", err)
	}
}

func TestListTasksIgnoresUnsupportedSortField(t *testing.T) {
	setupCronServiceTestDB(t)
	svc := &CronTaskService{Scheduler: NewScheduler()}
	task := model.SysCronTask{Code: "safe_sort", TaskCode: "cleanup_login_logs", Name: "排序测试", CronExpr: "0 2 * * *", Params: datatypes.JSON([]byte(`{}`))}
	if err := global.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}
	req := request.CronTaskListRequest{
		PageRequest: request.PageRequest{Page: 1, PageSize: 10},
		SortField:   "id;DROP TABLE sys_cron_task",
		SortOrder:   "descend",
	}
	tasks, total, err := svc.ListTasks(&req)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(tasks) != 1 {
		t.Fatalf("unexpected list result: total=%d len=%d", total, len(tasks))
	}
}

func TestRunNowRejectsConcurrentSameTask(t *testing.T) {
	setupCronServiceTestDB(t)
	code := "test_blocking_task"
	started := make(chan struct{})
	release := make(chan struct{})
	Register(TaskHandler{
		Code:        code,
		Name:        "测试阻塞任务",
		Description: "测试同任务串行",
		ParamSchema: map[string]ParamDefinition{},
		Run: func(ctx context.Context, params map[string]interface{}) (string, error) {
			close(started)
			<-release
			return "ok", nil
		},
	})

	scheduler := NewScheduler()
	scheduler.Start(context.Background())
	defer scheduler.Stop(time.Second)
	svc := &CronTaskService{Scheduler: scheduler}
	task := model.SysCronTask{Code: "blocking", TaskCode: code, Name: "阻塞任务", CronExpr: "0 2 * * *", Params: datatypes.JSON([]byte(`{}`))}
	if err := global.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RunNow(task.ID, 1); err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("task did not start")
	}
	if _, err := svc.RunNow(task.ID, 1); !errors.Is(err, ErrTaskAlreadyRunning) {
		t.Fatalf("expected ErrTaskAlreadyRunning, got %v", err)
	}
	close(release)
	deadline := time.After(time.Second)
	for {
		var count int64
		if err := global.DB.Model(&model.SysCronLog{}).Where("task_id = ? AND status = ?", task.ID, model.CronLogStatusSuccess).Count(&count).Error; err != nil {
			t.Fatal(err)
		}
		if count == 1 {
			break
		}
		select {
		case <-deadline:
			t.Fatal("task did not finish")
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func TestScheduledTriggerWritesSkippedLogWhenTaskAlreadyRunning(t *testing.T) {
	setupCronServiceTestDB(t)
	task := model.SysCronTask{Code: "skip_task", TaskCode: "cleanup_login_logs", Name: "跳过任务", CronExpr: "0 2 * * *", Params: datatypes.JSON([]byte(`{}`))}
	item := &scheduledTask{task: task}
	item.mu.Lock()
	scheduler := NewScheduler()
	scheduler.execute(item, TaskHandler{Run: func(ctx context.Context, params map[string]interface{}) (string, error) {
		t.Fatal("handler should not run while task lock is held")
		return "", nil
	}}, model.CronTriggeredBySchedule, 0)
	item.mu.Unlock()
	var log model.SysCronLog
	if err := global.DB.First(&log).Error; err != nil {
		t.Fatal(err)
	}
	if log.Status != model.CronLogStatusSkipped {
		t.Fatalf("log status = %s, want %s", log.Status, model.CronLogStatusSkipped)
	}
}

func TestSchedulerStopWaitsForManualRunNowToFinish(t *testing.T) {
	setupCronServiceTestDB(t)
	code := "test_stop_wait_task"
	started := make(chan struct{})
	release := make(chan struct{})
	Register(TaskHandler{
		Code:        code,
		Name:        "测试停机等待任务",
		Description: "测试停机等待手动任务收尾",
		ParamSchema: map[string]ParamDefinition{},
		Run: func(ctx context.Context, params map[string]interface{}) (string, error) {
			close(started)
			<-release
			return "done", nil
		},
	})
	scheduler := NewScheduler()
	scheduler.Start(context.Background())
	task := model.SysCronTask{Code: "stop_wait", TaskCode: code, Name: "停机等待任务", CronExpr: "0 2 * * *", Params: datatypes.JSON([]byte(`{}`))}
	if err := global.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := scheduler.RunNow(task, 1); err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("task did not start")
	}
	stopped := make(chan struct{})
	go func() {
		scheduler.Stop(time.Second)
		close(stopped)
	}()
	select {
	case <-stopped:
		t.Fatal("scheduler stopped before manual task finished")
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("scheduler did not stop after manual task finished")
	}
}

func TestUpdateTaskRunStateUsesCurrentExecutionLog(t *testing.T) {
	setupCronServiceTestDB(t)
	task := model.SysCronTask{Code: "duration_task", TaskCode: "cleanup_login_logs", Name: "耗时任务", CronExpr: "0 2 * * *", Params: datatypes.JSON([]byte(`{}`)), Status: model.CronTaskStatusEnabled}
	if err := global.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}
	original := model.SysCronLog{TaskID: task.ID, TaskCode: task.TaskCode, StartedAt: time.Now().Add(-time.Second), DurationMS: 800, Status: model.CronLogStatusSuccess, TriggeredBy: model.CronTriggeredBySchedule}
	if err := global.DB.Create(&original).Error; err != nil {
		t.Fatal(err)
	}
	skipped := model.SysCronLog{TaskID: task.ID, TaskCode: task.TaskCode, StartedAt: time.Now(), DurationMS: 1, Status: model.CronLogStatusSkipped, TriggeredBy: model.CronTriggeredBySchedule}
	if err := global.DB.Create(&skipped).Error; err != nil {
		t.Fatal(err)
	}
	updateTaskRunState(task.ID, original.ID, model.CronLogStatusSuccess)
	var updated model.SysCronTask
	if err := global.DB.First(&updated, task.ID).Error; err != nil {
		t.Fatal(err)
	}
	if updated.LastDurationMS != 800 {
		t.Fatalf("last duration = %d, want 800", updated.LastDurationMS)
	}
}

func TestListLogsFiltersTriggeredBy(t *testing.T) {
	setupCronServiceTestDB(t)
	svc := &CronTaskService{Scheduler: NewScheduler()}
	scheduled := model.SysCronLog{TaskID: 1, TaskCode: "cleanup_login_logs", StartedAt: time.Now(), Status: model.CronLogStatusSuccess, TriggeredBy: model.CronTriggeredBySchedule}
	manual := model.SysCronLog{TaskID: 1, TaskCode: "cleanup_login_logs", StartedAt: time.Now(), Status: model.CronLogStatusSuccess, TriggeredBy: model.CronTriggeredByManual}
	if err := global.DB.Create(&scheduled).Error; err != nil {
		t.Fatal(err)
	}
	if err := global.DB.Create(&manual).Error; err != nil {
		t.Fatal(err)
	}
	req := request.CronLogListRequest{
		PageRequest: request.PageRequest{Page: 1, PageSize: 10},
		TriggeredBy: model.CronTriggeredByManual,
	}
	logs, total, err := svc.ListLogs(&req)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(logs) != 1 || logs[0].TriggeredBy != model.CronTriggeredByManual {
		t.Fatalf("unexpected logs: total=%d logs=%+v", total, logs)
	}
}

func TestCleanupLoginLogsDeletesOnlyExpiredRows(t *testing.T) {
	db := setupCronServiceTestDB(t)
	now := time.Now()
	oldLog := model.SysLoginLog{Username: "old", CreatedAt: now.AddDate(0, 0, -60)}
	recentLog := model.SysLoginLog{Username: "recent", CreatedAt: now.AddDate(0, 0, -1)}
	if err := db.Create(&oldLog).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&recentLog).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := cleanupLoginLogs(context.Background(), map[string]interface{}{"retain_days": 30, "batch_limit": 1000}); err != nil {
		t.Fatal(err)
	}
	var count int64
	if err := db.Model(&model.SysLoginLog{}).Where("username = ?", "old").Count(&count).Error; err != nil || count != 0 {
		t.Fatalf("old log count=%d err=%v", count, err)
	}
	if err := db.Model(&model.SysLoginLog{}).Where("username = ?", "recent").Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("recent log count=%d err=%v", count, err)
	}
}

func TestCleanupLogsRespectsCanceledContext(t *testing.T) {
	setupCronServiceTestDB(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := cleanupOperationLogs(ctx, map[string]interface{}{"retain_days": 30, "batch_limit": 1000})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func requestLikeTask(code, taskCode string) request.CronTaskSaveRequest {
	return request.CronTaskSaveRequest{
		Code:     code,
		TaskCode: taskCode,
		Name:     "测试任务",
		CronExpr: "0 2 * * *",
		Params:   []byte(`{}`),
	}
}
