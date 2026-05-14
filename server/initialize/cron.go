package initialize

import (
	"context"
	"time"

	"server/global"
	"server/service"
)

func InitCronScheduler(ctx context.Context) {
	service.Cron.Scheduler.Start(ctx)
	if err := service.Cron.Scheduler.LoadFromDB(); err != nil {
		global.Log.Errorf("加载定时任务失败: %v", err)
		return
	}
	global.Log.Info("定时任务调度器启动成功")
}

func ShutdownCronScheduler() {
	service.Cron.Scheduler.Stop(30 * time.Second)
	global.Log.Info("定时任务调度器已停止")
}
