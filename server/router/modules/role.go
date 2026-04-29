package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/model/request"
	"server/router/registry"
)

func init() {
	RegisterModule(&RoleModule{})
}

type RoleModule struct{}

func (m *RoleModule) Name() string {
	return "角色管理"
}

func (m *RoleModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *RoleModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/roles", m.Name(), "角色列表", v1.Role.GetRoleList, registry.WithAuth())
	R(rg, "GET", "/roles/:id", m.Name(), "角色详情", v1.Role.GetRole, registry.WithAuth())
	R(rg, "POST", "/roles", m.Name(), "创建角色", v1.Role.CreateRole,
		registry.WithAuth(), registry.WithRequest(request.CreateRoleRequest{}))
	R(rg, "PUT", "/roles/:id", m.Name(), "更新角色", v1.Role.UpdateRole,
		registry.WithAuth(), registry.WithRequest(request.UpdateRoleRequest{}))
	R(rg, "DELETE", "/roles/:id", m.Name(), "删除角色", v1.Role.DeleteRole, registry.WithAuth())
	R(rg, "PUT", "/roles/:id/menus", m.Name(), "分配菜单", v1.Role.AssignMenus,
		registry.WithAuth(), registry.WithRequest(request.AssignMenusRequest{}))
	R(rg, "PUT", "/roles/:id/apis", m.Name(), "分配API", v1.Role.AssignApis,
		registry.WithAuth(), registry.WithRequest(request.AssignApisRequest{}))
	R(rg, "PUT", "/roles/:id/data-scopes", m.Name(), "分配数据权限", v1.Role.AssignDataScopes,
		registry.WithAuth(), registry.WithRequest(request.AssignRoleDataScopesRequest{}))
}
