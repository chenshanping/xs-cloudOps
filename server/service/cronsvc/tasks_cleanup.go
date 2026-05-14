package cronsvc

import (
	"context"
	"fmt"
	"time"

	"server/global"
	"server/model"
)

const (
	defaultRetainDays = 30
	defaultBatchLimit = 1000
	maxBatchLimit     = 10000
)

func cleanupLogParamSchema() map[string]ParamDefinition {
	minRetain := 1
	maxRetain := 3650
	minBatch := 1
	maxBatch := maxBatchLimit
	return map[string]ParamDefinition{
		"retain_days": {
			Type:        ParamTypeInt,
			Required:    true,
			Default:     defaultRetainDays,
			Description: "保留最近多少天的日志",
			Min:         &minRetain,
			Max:         &maxRetain,
		},
		"batch_limit": {
			Type:        ParamTypeInt,
			Required:    true,
			Default:     defaultBatchLimit,
			Description: "单批最多删除记录数",
			Min:         &minBatch,
			Max:         &maxBatch,
		},
	}
}

func cleanupLoginLogs(ctx context.Context, params map[string]interface{}) (string, error) {
	return cleanupLogs(ctx, "登录日志", &model.SysLoginLog{}, params)
}

func cleanupOperationLogs(ctx context.Context, params map[string]interface{}) (string, error) {
	return cleanupLogs(ctx, "操作日志", &model.SysOperationLog{}, params)
}

func cleanupLogs(ctx context.Context, label string, modelValue interface{}, params map[string]interface{}) (string, error) {
	retainDays := getIntParam(params, "retain_days", defaultRetainDays)
	batchLimit := getIntParam(params, "batch_limit", defaultBatchLimit)
	if batchLimit <= 0 {
		batchLimit = defaultBatchLimit
	}
	if batchLimit > maxBatchLimit {
		batchLimit = maxBatchLimit
	}

	cutoff := time.Now().AddDate(0, 0, -retainDays)
	var total int64
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Sprintf("%s清理中止，已删除%d条，截止时间%s", label, total, cutoff.Format(time.RFC3339)), err
		}

		var ids []uint
		if err := global.DB.WithContext(ctx).
			Model(modelValue).
			Where("created_at < ?", cutoff).
			Order("id ASC").
			Limit(batchLimit).
			Pluck("id", &ids).Error; err != nil {
			return "", err
		}
		if len(ids) == 0 {
			break
		}

		result := global.DB.WithContext(ctx).Where("id IN ?", ids).Delete(modelValue)
		if result.Error != nil {
			return "", result.Error
		}
		total += result.RowsAffected
		if len(ids) < batchLimit {
			break
		}
	}

	return fmt.Sprintf("%s清理完成，已删除%d条，保留%d天，截止时间%s", label, total, retainDays, cutoff.Format(time.RFC3339)), nil
}
