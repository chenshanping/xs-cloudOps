package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/router/registry"
)

func init() {
	RegisterModule(&EchartModule{})
}

type EchartModule struct{}

func (m *EchartModule) Name() string {
	return "图表统计"
}

func (m *EchartModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *EchartModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	// 图表统计接口
	R(rg, "GET", "/echart/user-role-stats", m.Name(), "用户角色占比", v1.Echart.GetUserRoleStats, registry.WithAuth())
	R(rg, "GET", "/echart/user-status-stats", m.Name(), "用户状态统计", v1.Echart.GetUserStatusStats, registry.WithAuth())
	R(rg, "GET", "/echart/role-status-stats", m.Name(), "角色状态统计", v1.Echart.GetRoleStatusStats, registry.WithAuth())
	R(rg, "GET", "/echart/user-register-trend", m.Name(), "用户注册趋势", v1.Echart.GetUserRegisterTrend, registry.WithAuth())
}
