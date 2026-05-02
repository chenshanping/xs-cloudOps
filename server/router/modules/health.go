package modules

import (
	v1 "server/api/v1"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterModule(&HealthModule{})
}

type HealthModule struct{}

func (m *HealthModule) Name() string {
	return "健康检查"
}

func (m *HealthModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/health", m.Name(), "服务健康检查", v1.Health.GetHealth)
}

func (m *HealthModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
}
