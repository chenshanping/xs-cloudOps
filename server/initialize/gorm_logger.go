package initialize

import (
	"context"
	"fmt"
	"server/global"
	"server/model"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

type SlowLogger struct {
	logger.Interface
	SlowThreshold time.Duration
}

func NewSlowLogger(base logger.Interface, threshold time.Duration) *SlowLogger {
	return &SlowLogger{
		Interface:     base,
		SlowThreshold: threshold,
	}
}

func (l *SlowLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 先调用原始logger
	l.Interface.Trace(ctx, begin, fc, err)

	elapsed := time.Since(begin)
	if elapsed >= l.SlowThreshold {
		sql, rows := fc()
		// 排除对慢日志表本身的操作，避免死循环
		if strings.Contains(sql, "sys_slow_log") {
			return
		}

		// 获取调用来源
		source := getCallerSource()

		// 异步写入慢日志
		go func() {
			slowLog := model.SysSlowLog{
				SQL:     sql,
				Rows:    rows,
				Latency: float64(elapsed.Milliseconds()),
				Source:  source,
			}
			if global.DB != nil {
				global.DB.Create(&slowLog)
			}
		}()
	}
}

func getCallerSource() string {
	for i := 4; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// 排除gorm内部调用
		if !strings.Contains(file, "gorm.io") && !strings.Contains(file, "gorm-adapter") {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	return ""
}
