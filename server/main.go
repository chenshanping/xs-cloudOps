package main

import (
	"flag"
	"fmt"
	"log"
	"server/global"
	"server/initialize"
	"server/router"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("c", "config.yaml", "配置文件路径")
	flag.Parse()

	// 初始化配置
	if err := initialize.InitConfig(*configPath); err != nil {
		log.Fatalf("startup failed: %v", err)
	}

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
