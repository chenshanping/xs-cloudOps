package v1

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"server/middleware"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/service/cronsvc"
)

type CronApi struct{}

var Cron = new(CronApi)

func (a *CronApi) ListTasks(c *gin.Context) {
	var req request.CronTaskListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	tasks, total, err := service.Cron.ListTasks(&req)
	if err != nil {
		response.Fail(c, "获取定时任务失败")
		return
	}
	response.OkWithPage(c, tasks, total, req.Page, req.PageSize)
}

func (a *CronApi) CreateTask(c *gin.Context) {
	var req request.CronTaskSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.Cron.CreateTask(&req, middleware.GetUserID(c)); err != nil {
		writeCronError(c, err)
		return
	}
	response.OkWithMessage(c, "创建成功")
}

func (a *CronApi) UpdateTask(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var req request.CronTaskSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.Cron.UpdateTask(id, &req); err != nil {
		writeCronError(c, err)
		return
	}
	response.OkWithMessage(c, "更新成功")
}

func (a *CronApi) DeleteTask(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := service.Cron.DeleteTask(id); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "删除成功")
}

func (a *CronApi) EnableTask(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := service.Cron.EnableTask(id); err != nil {
		writeCronError(c, err)
		return
	}
	response.OkWithMessage(c, "启用成功")
}

func (a *CronApi) DisableTask(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := service.Cron.DisableTask(id); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "停用成功")
}

func (a *CronApi) RunNow(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	logID, err := service.Cron.RunNow(id, middleware.GetUserID(c))
	if err != nil {
		writeCronError(c, err)
		return
	}
	response.OkWithData(c, gin.H{"log_id": logID})
}

func (a *CronApi) Registry(c *gin.Context) {
	response.OkWithData(c, service.Cron.Registry())
}

func (a *CronApi) ListLogs(c *gin.Context) {
	var req request.CronLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	logs, total, err := service.Cron.ListLogs(&req)
	if err != nil {
		response.Fail(c, "获取执行日志失败")
		return
	}
	response.OkWithPage(c, logs, total, req.Page, req.PageSize)
}

func (a *CronApi) GetLog(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	log, err := service.Cron.GetLog(id)
	if err != nil {
		response.Fail(c, "获取执行日志详情失败")
		return
	}
	response.OkWithData(c, log)
}

func parseUintParam(c *gin.Context, name string) (uint, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return 0, false
	}
	return uint(id), true
}

func writeCronError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, cronsvc.ErrTaskAlreadyRunning):
		response.FailWithCode(c, 409, "cron_task_already_running")
	case err.Error() == "cron_task_code_not_registered":
		response.BadRequest(c, "cron_task_code_not_registered")
	case err.Error() == "cron_expr_invalid":
		response.BadRequest(c, "cron_expr_invalid")
	default:
		response.Fail(c, err.Error())
	}
}
