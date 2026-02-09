package modules

import (
	"github.com/gin-gonic/gin"

	v1 "go-base-server/api/v1"
	"go-base-server/model/request"
	"go-base-server/router/registry"
)

func init() {
	RegisterModule(&MenuModule{})
}

type MenuModule struct{}

func (m *MenuModule) Name() string {
	return "菜单管理"
}

func (m *MenuModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *MenuModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/menus", m.Name(), "菜单列表", v1.Menu.GetMenuTree, registry.WithAuth())
	R(rg, "GET", "/menus/:id", m.Name(), "菜单详情", v1.Menu.GetMenu, registry.WithAuth())
	R(rg, "GET", "/user/menus", m.Name(), "获取用户菜单", v1.Menu.GetUserMenus, registry.WithAuth())
	R(rg, "POST", "/menus", m.Name(), "创建菜单", v1.Menu.CreateMenu,
		registry.WithAuth(), registry.WithRequest(request.CreateMenuRequest{}))
	R(rg, "PUT", "/menus/:id", m.Name(), "更新菜单", v1.Menu.UpdateMenu,
		registry.WithAuth(), registry.WithRequest(request.UpdateMenuRequest{}))
	R(rg, "DELETE", "/menus/:id", m.Name(), "删除菜单", v1.Menu.DeleteMenu, registry.WithAuth())
}
