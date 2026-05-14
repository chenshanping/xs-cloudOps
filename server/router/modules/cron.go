package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/model/request"
	"server/router/registry"
)

func init() {
	RegisterModule(&CronModule{})
}

type CronModule struct{}

func (m *CronModule) Name() string {
	return "定时任务"
}

func (m *CronModule) RegisterPublicRoutes(rg *gin.RouterGroup) {}

func (m *CronModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/monitor/cron-task", m.Name(), "定时任务列表", v1.Cron.ListTasks,
		registry.WithAuth(), registry.WithRequest(request.CronTaskListRequest{}))
	R(rg, "POST", "/monitor/cron-task", m.Name(), "创建定时任务", v1.Cron.CreateTask,
		registry.WithAuth(), registry.WithRequest(request.CronTaskSaveRequest{}))
	R(rg, "GET", "/monitor/cron-task/registry", m.Name(), "定时任务注册列表", v1.Cron.Registry, registry.WithAuth())
	R(rg, "PUT", "/monitor/cron-task/:id", m.Name(), "更新定时任务", v1.Cron.UpdateTask,
		registry.WithAuth(), registry.WithRequest(request.CronTaskSaveRequest{}))
	R(rg, "DELETE", "/monitor/cron-task/:id", m.Name(), "删除定时任务", v1.Cron.DeleteTask, registry.WithAuth())
	R(rg, "POST", "/monitor/cron-task/:id/enable", m.Name(), "启用定时任务", v1.Cron.EnableTask, registry.WithAuth())
	R(rg, "POST", "/monitor/cron-task/:id/disable", m.Name(), "停用定时任务", v1.Cron.DisableTask, registry.WithAuth())
	R(rg, "POST", "/monitor/cron-task/:id/run", m.Name(), "立即执行定时任务", v1.Cron.RunNow, registry.WithAuth())
	R(rg, "GET", "/monitor/cron-log", m.Name(), "定时任务执行日志", v1.Cron.ListLogs,
		registry.WithAuth(), registry.WithRequest(request.CronLogListRequest{}))
	R(rg, "GET", "/monitor/cron-log/:id", m.Name(), "定时任务执行日志详情", v1.Cron.GetLog, registry.WithAuth())
}
