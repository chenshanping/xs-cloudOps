package modules

import (
	"server/router/registry"

	"github.com/gin-gonic/gin"
)

// R 注册路由并记录元信息
func R(rg *gin.RouterGroup, method, path, group, summary string, handler gin.HandlerFunc, opts ...registry.RouteOption) {
	fullPath := rg.BasePath() + path
	registry.Register(method, fullPath, group, summary, handler, opts...)
	switch method {
	case "GET":
		rg.GET(path, handler)
	case "POST":
		rg.POST(path, handler)
	case "PUT":
		rg.PUT(path, handler)
	case "DELETE":
		rg.DELETE(path, handler)
	}
}
