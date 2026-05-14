package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/router/registry"
)

func init() {
	RegisterModule(&MonitorModule{})
}

type MonitorModule struct{}

func (m *MonitorModule) Name() string {
	return "服务监控"
}

func (m *MonitorModule) RegisterPublicRoutes(rg *gin.RouterGroup) {}

func (m *MonitorModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/monitor/server", m.Name(), "服务器指标", v1.Monitor.Server, registry.WithAuth())
	R(rg, "GET", "/monitor/runtime", m.Name(), "运行时指标", v1.Monitor.Runtime, registry.WithAuth())
	R(rg, "GET", "/monitor/db", m.Name(), "数据库连接池指标", v1.Monitor.DB, registry.WithAuth())
	R(rg, "GET", "/monitor/redis", m.Name(), "Redis 缓存指标", v1.Monitor.Redis, registry.WithAuth())
	R(rg, "POST", "/monitor/redis/clear", m.Name(), "清理 Redis 缓存", v1.Monitor.ClearCache, registry.WithAuth())
	R(rg, "GET", "/monitor/oss", m.Name(), "对象存储健康", v1.Monitor.Oss, registry.WithAuth())
	R(rg, "GET", "/monitor/dependency", m.Name(), "依赖健康概览", v1.Monitor.Dependency, registry.WithAuth())
}
