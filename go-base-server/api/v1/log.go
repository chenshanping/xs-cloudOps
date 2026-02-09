package v1

import (
	"github.com/gin-gonic/gin"

	"go-base-server/model/request"
	"go-base-server/model/response"
	"go-base-server/service"
)

type LogApi struct{}

var Log = new(LogApi)

// 获取操作日志列表
func (a *LogApi) GetOperationLogList(c *gin.Context) {
	var req request.LogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	logs, total, err := service.Log.GetOperationLogList(&req)
	if err != nil {
		response.Fail(c, "获取操作日志失败")
		return
	}

	response.OkWithPage(c, logs, total, req.Page, req.PageSize)
}

// 获取登录日志列表
func (a *LogApi) GetLoginLogList(c *gin.Context) {
	var req request.LogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	logs, total, err := service.Log.GetLoginLogList(&req)
	if err != nil {
		response.Fail(c, "获取登录日志失败")
		return
	}

	response.OkWithPage(c, logs, total, req.Page, req.PageSize)
}

// 获取路由分组列表
func (a *LogApi) GetRouteGroups(c *gin.Context) {
	groups, err := service.Log.GetRouteGroups()
	if err != nil {
		response.Fail(c, "获取路由分组失败")
		return
	}
	response.OkWithData(c, groups)
}

// 获取慢查询日志列表
func (a *LogApi) GetSlowLogList(c *gin.Context) {
	var req request.SlowLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	logs, total, err := service.Log.GetSlowLogList(&req)
	if err != nil {
		response.Fail(c, "获取慢查询日志失败")
		return
	}

	response.OkWithPage(c, logs, total, req.Page, req.PageSize)
}
