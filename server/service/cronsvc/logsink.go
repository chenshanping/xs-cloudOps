package cronsvc

import (
	"fmt"
	"runtime/debug"
	"time"
	"unicode/utf8"

	"server/global"
	"server/model"
)

const (
	maxCronLogSummaryBytes = 4096
	maxCronLogErrorBytes   = 4096
)

func StartLog(taskID uint, taskCode, triggeredBy string, actorID uint) (uint, error) {
	log := model.SysCronLog{
		TaskID:      taskID,
		TaskCode:    taskCode,
		StartedAt:   time.Now(),
		Status:      model.CronLogStatusRunning,
		TriggeredBy: triggeredBy,
		ActorUserID: actorID,
	}
	if err := global.DB.Create(&log).Error; err != nil {
		return 0, err
	}
	return log.ID, nil
}

func FinishLog(logID uint, status model.CronLogStatus, summary string, runErr error) error {
	finishedAt := time.Now()
	var log model.SysCronLog
	if err := global.DB.First(&log, logID).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"finished_at":   finishedAt,
		"duration_ms":   finishedAt.Sub(log.StartedAt).Milliseconds(),
		"status":        status,
		"summary":       truncateString(summary, maxCronLogSummaryBytes),
		"error_message": "",
	}
	if runErr != nil {
		updates["error_message"] = truncateString(runErr.Error(), maxCronLogErrorBytes)
	}
	return global.DB.Model(&log).Updates(updates).Error
}

func markRunningLogsFailed() {
	message := "服务重启前任务未正常结束"
	finishedAt := time.Now()
	if err := global.DB.Model(&model.SysCronLog{}).
		Where("status = ?", model.CronLogStatusRunning).
		Updates(map[string]interface{}{
			"status":        model.CronLogStatusFailure,
			"finished_at":   finishedAt,
			"error_message": message,
		}).Error; err != nil {
		global.Log.Errorf("标记未完成定时任务日志失败: %v", err)
	}
}

func panicToError(recovered interface{}) error {
	return fmt.Errorf("panic: %v\n%s", recovered, string(debug.Stack()))
}

func truncateString(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	cut := maxBytes
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	for cut > 0 && !utf8.ValidString(s[:cut]) {
		cut--
	}
	return s[:cut] + "..."
}
