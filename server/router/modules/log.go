package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/model/request"
	"server/router/registry"
)

func init() {
	RegisterModule(&LogModule{})
}

type LogModule struct{}

func (m *LogModule) Name() string {
	return "日志管理"
}

func (m *LogModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *LogModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/logs/operation", m.Name(), "操作日志列表", v1.Log.GetOperationLogList,
		registry.WithAuth(), registry.WithRequest(request.LogListRequest{}))
	R(rg, "GET", "/logs/login", m.Name(), "登录日志列表", v1.Log.GetLoginLogList,
		registry.WithAuth(), registry.WithRequest(request.LogListRequest{}))
	R(rg, "GET", "/logs/slow", m.Name(), "慢查询日志列表", v1.Log.GetSlowLogList,
		registry.WithAuth(), registry.WithRequest(request.SlowLogListRequest{}))
	R(rg, "GET", "/logs/route-groups", m.Name(), "获取路由分组列表", v1.Log.GetRouteGroups,
		registry.WithAuth())
}
