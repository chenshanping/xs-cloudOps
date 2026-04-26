package router

import (
	"server/global"
	"server/middleware"
	"server/router/modules"
	"server/swagger"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 跨域中间件
	r.Use(middleware.Cors())

	// 静态资源：
	r.Static("/api/v1/upload", "./uploads")
	// 公开路由组
	public := r.Group("/api/v1")

	// 需要认证的路由组
	private := r.Group("/api/v1")
	private.Use(middleware.JWTAuth())
	private.Use(middleware.OperationLog())
	private.Use(middleware.CasbinAuth())

	// 加载所有已注册的路由模块
	for _, m := range modules.GetAllModules() {
		m.RegisterPublicRoutes(public)
		m.RegisterPrivateRoutes(private)
	}

	// 注册 Swagger 文档（在所有路由注册完成后）
	swagger.RegisterRoutes(r, swagger.Config{
		Title:       "Go Base Server API",
		Description: "基于 Gin + GORM 的后台管理系统 API",
		Version:     "1.0",
		Host:        global.Config.Server.Host,
		BasePath:    "/api/v1",
	})

	return r
}
