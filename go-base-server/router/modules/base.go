package modules

import "github.com/gin-gonic/gin"

// RouterModule 路由模块接口
type RouterModule interface {
	Name() string
	RegisterPublicRoutes(rg *gin.RouterGroup)
	RegisterPrivateRoutes(rg *gin.RouterGroup)
}

// 所有已注册的模块
var registeredModules []RouterModule

// RegisterModule 注册路由模块
func RegisterModule(m RouterModule) {
	registeredModules = append(registeredModules, m)
}

// GetAllModules 获取所有已注册的模块
func GetAllModules() []RouterModule {
	return registeredModules
}
