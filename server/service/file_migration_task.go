package service

import (
	"fmt"
	"sync"
	"time"

	modelrequest "server/model/request"
	modelresponse "server/model/response"
)

const (
	fileMigrationTaskStatusScanning = "SCANNING"
	fileMigrationTaskStatusRunning  = "RUNNING"
	fileMigrationTaskStatusSuccess  = "SUCCESS"
	fileMigrationTaskStatusFailed   = "FAILED"
	fileMigrationTaskItemLimit      = 200
)

type fileMigrationTask struct {
	mu     sync.RWMutex
	status modelresponse.FileMigrationTaskStatus
}

func newFileMigrationTask(req modelrequest.FileMigrationRequest) *fileMigrationTask {
	now := time.Now()
	return &fileMigrationTask{
		status: modelresponse.FileMigrationTaskStatus{
			TaskID:            fmt.Sprintf("file-migration-%d", now.UnixNano()),
			Status:            fileMigrationTaskStatusScanning,
			Message:           "正在预检查迁移文件",
			Scope:             req.Scope,
			SourceStorageType: req.SourceStorageType,
			TargetStorageType: req.TargetStorageType,
			StartedAt:         now.Format(time.RFC3339),
			Items:             make([]modelresponse.FileMigrationItem, 0),
		},
	}
}

func (t *fileMigrationTask) snapshot() *modelresponse.FileMigrationTaskStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()

	status := t.status
	status.Items = append([]modelresponse.FileMigrationItem(nil), t.status.Items...)
	return &status
}

func (t *fileMigrationTask) setPrecheck(result *modelresponse.FileMigrationResult) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.status.TotalCount = result.TotalCount
	t.status.TotalSize = result.TotalSize
	t.status.PendingCount = result.PendingCount
	t.status.PendingSize = result.PendingSize
	t.status.SkippedCount = result.SkippedCount
	t.status.SkippedSize = result.SkippedSize
	t.status.ConflictCount = result.ConflictCount
	t.status.ConflictSize = result.ConflictSize
	t.status.MissingSourceCount = result.MissingSourceCount
	t.status.MissingSourceSize = result.MissingSourceSize
	t.status.Message = "预检查完成，开始执行迁移"
}

func (t *fileMigrationTask) markRunning() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.status.Status = fileMigrationTaskStatusRunning
	t.status.Message = "正在迁移文件"
}

func (t *fileMigrationTask) setCurrentFile(fileID uint, fileName string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.status.CurrentFileID = fileID
	t.status.CurrentFileName = fileName
}

func (t *fileMigrationTask) recordHandled(item modelresponse.FileMigrationItem, size int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.status.ProcessedCount++
	t.status.ProcessedSize += size
	switch item.Action {
	case "MIGRATED":
		t.status.MigratedCount++
	case "WARNING":
		t.status.MigratedCount++
		t.status.WarningCount++
	case "FAILED":
		t.status.FailedCount++
	}

	t.status.Items = append(t.status.Items, item)
	if len(t.status.Items) > fileMigrationTaskItemLimit {
		t.status.Items = append([]modelresponse.FileMigrationItem(nil), t.status.Items[len(t.status.Items)-fileMigrationTaskItemLimit:]...)
	}
}

func (t *fileMigrationTask) finishSuccess(message string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.status.Status = fileMigrationTaskStatusSuccess
	t.status.Message = message
	t.status.CurrentFileID = 0
	t.status.CurrentFileName = ""
	t.status.FinishedAt = time.Now().Format(time.RFC3339)
}

func (t *fileMigrationTask) finishFailed(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.status.Status = fileMigrationTaskStatusFailed
	if err != nil {
		t.status.Message = err.Error()
	} else {
		t.status.Message = "迁移任务失败"
	}
	t.status.CurrentFileID = 0
	t.status.CurrentFileName = ""
	t.status.FinishedAt = time.Now().Format(time.RFC3339)
}

func (t *fileMigrationTask) isActive() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.status.Status == fileMigrationTaskStatusScanning || t.status.Status == fileMigrationTaskStatusRunning
}

type fileMigrationTaskManager struct {
	mu      sync.RWMutex
	current *fileMigrationTask
}

func newFileMigrationTaskManager() *fileMigrationTaskManager {
	return &fileMigrationTaskManager{}
}

func (m *fileMigrationTaskManager) start(req modelrequest.FileMigrationRequest, runner func(task *fileMigrationTask)) (*modelresponse.FileMigrationTaskStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.current != nil && m.current.isActive() {
		return nil, fmt.Errorf("已有文件迁移任务正在执行，请稍后再试")
	}

	task := newFileMigrationTask(req)
	m.current = task
	go runner(task)
	return task.snapshot(), nil
}

func (m *fileMigrationTaskManager) snapshot() *modelresponse.FileMigrationTaskStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current == nil {
		return nil
	}
	return m.current.snapshot()
}

var fileMigrationTasks = newFileMigrationTaskManager()
