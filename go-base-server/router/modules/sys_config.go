package modules

import (
	"go-base-server/model/request"

	"github.com/gin-gonic/gin"

	v1 "go-base-server/api/v1"
	"go-base-server/router/registry"
)

func init() {
	RegisterModule(&SysConfigModule{})
}

type SysConfigModule struct{}

// TestEmailRequest 发送测试邮件请求
type TestEmailRequest struct {
	Email string `json:"email" binding:"required,email" comment:"接收测试邮件的邮箱地址"`
}

func (m *SysConfigModule) Name() string {
	return "配置管理"
}

func (m *SysConfigModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 公开路由，无需认证
	R(rg, "POST", "/configs/keys", m.Name(), "批量获取配置", v1.Config.GetConfigsByKeys)
}

func (m *SysConfigModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/configs", m.Name(), "配置列表", v1.Config.GetConfigList, registry.WithAuth())
	R(rg, "GET", "/configs/key/:key", m.Name(), "根据key获取配置", v1.Config.GetConfigByKey, registry.WithAuth())
	R(rg, "POST", "/configs", m.Name(), "创建配置", v1.Config.CreateConfig, registry.WithAuth())
	R(rg, "PUT", "/configs/:id", m.Name(), "更新配置", v1.Config.UpdateConfig, registry.WithAuth())
	R(rg, "PUT", "/configs/batch", m.Name(), "批量更新配置", v1.Config.BatchUpdateConfigs, registry.WithAuth())
	R(rg, "DELETE", "/configs/:id", m.Name(), "删除配置", v1.Config.DeleteConfig, registry.WithAuth())
	R(rg, "POST", "/configs/test-email", m.Name(), "发送测试邮件", v1.Config.SendTestEmail, registry.WithAuth(), registry.WithRequest(request.TestEmailRequest{}))
}
