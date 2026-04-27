package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/model/request"
	"server/router/registry"
)

func init() {
	RegisterModule(&SysApiModule{})
}

type SysApiModule struct{}

func (m *SysApiModule) Name() string {
	return "API管理"
}

func (m *SysApiModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *SysApiModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/apis", m.Name(), "API列表", v1.Api.GetApiList,
		registry.WithAuth(), registry.WithRequest(request.ApiListRequest{}))
	R(rg, "GET", "/apis/all", m.Name(), "全部API", v1.Api.GetAllApis, registry.WithAuth())
	R(rg, "GET", "/apis/groups", m.Name(), "API分组", v1.Api.GetApiGroups, registry.WithAuth())
	R(rg, "GET", "/apis/:id", m.Name(), "API详情", v1.Api.GetApi, registry.WithAuth())
	R(rg, "POST", "/apis", m.Name(), "创建API", v1.Api.CreateApi,
		registry.WithAuth(), registry.WithRequest(request.CreateApiRequest{}))
	R(rg, "PUT", "/apis/:id", m.Name(), "更新API", v1.Api.UpdateApi,
		registry.WithAuth(), registry.WithRequest(request.UpdateApiRequest{}))
	R(rg, "DELETE", "/apis/:id", m.Name(), "删除API", v1.Api.DeleteApi, registry.WithAuth())
	R(rg, "POST", "/apis/sync", m.Name(), "同步API", v1.Api.SyncApis, registry.WithAuth())
}
