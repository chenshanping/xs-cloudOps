package main

import (
	"fmt"
	"go-base-server/global"
	"go-base-server/initialize"
	"go-base-server/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	initialize.InitConfig()

	// 初始化日志
	initialize.InitLogger()

	// 初始化数据库
	initialize.InitGorm()

	// 初始化Redis
	initialize.InitRedis()

	// 初始化Casbin（依赖数据库）
	initialize.InitCasbin()

	// 初始化数据库表及默认数据（需要Enforcer以便同步策略）
	initialize.InitDBTables()

	// 设置Gin模式
	gin.SetMode(global.Config.Server.Mode)

	// 初始化路由
	r := router.InitRouter()

	// 启动服务
	addr := fmt.Sprintf(":%d", global.Config.Server.Port)
	global.Log.Infof("服务启动成功，监听端口: %s", addr)

	if err := r.Run(addr); err != nil {
		global.Log.Fatalf("服务启动失败: %v", err)
	}
}
