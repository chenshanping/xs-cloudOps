package v1

import (
	"errors"

	"github.com/gin-gonic/gin"

	"server/model/response"
	"server/service"
	"server/service/monitorsvc"
)

type MonitorApi struct{}

var Monitor = new(MonitorApi)

type ClearCacheRequest struct {
	Prefix string `json:"prefix" binding:"required"`
}

func (a *MonitorApi) Server(c *gin.Context) {
	data, err := service.Monitor.CollectServerInfo(c.Request.Context())
	if err != nil {
		response.Fail(c, "获取服务器信息失败")
		return
	}
	response.OkWithData(c, data)
}

func (a *MonitorApi) Runtime(c *gin.Context) {
	data, err := service.Monitor.CollectRuntime(c.Request.Context())
	if err != nil {
		response.Fail(c, "获取运行时信息失败")
		return
	}
	response.OkWithData(c, data)
}

func (a *MonitorApi) DB(c *gin.Context) {
	data, err := service.Monitor.CollectDB(c.Request.Context())
	if err != nil {
		response.Fail(c, "获取数据库状态失败")
		return
	}
	response.OkWithData(c, data)
}

func (a *MonitorApi) Redis(c *gin.Context) {
	data, err := service.Monitor.CollectRedis(c.Request.Context())
	if err != nil {
		response.Fail(c, "获取 Redis 状态失败")
		return
	}
	response.OkWithData(c, data)
}

func (a *MonitorApi) ClearCache(c *gin.Context) {
	var req ClearCacheRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	data, err := service.Monitor.ClearByPrefix(c.Request.Context(), req.Prefix)
	if errors.Is(err, monitorsvc.ErrCachePrefixNotAllowed) {
		response.FailWithCode(c, response.BAD_REQUEST, "cache_prefix_not_allowed")
		return
	}
	if err != nil {
		response.Fail(c, "清理缓存失败")
		return
	}
	response.OkWithData(c, data)
}

func (a *MonitorApi) Oss(c *gin.Context) {
	data, err := service.Monitor.CollectOss(c.Request.Context())
	if err != nil {
		response.Fail(c, "获取存储健康状态失败")
		return
	}
	response.OkWithData(c, data)
}

func (a *MonitorApi) Dependency(c *gin.Context) {
	data, err := service.Monitor.CollectDependency(c.Request.Context())
	if err != nil {
		response.Fail(c, "获取依赖健康状态失败")
		return
	}
	response.OkWithData(c, data)
}
