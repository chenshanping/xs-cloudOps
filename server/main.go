package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"server/global"
	"server/initialize"
	"server/router"
	"syscall"
	"time"

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

	serverCtx, stopSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignal()

	cronCtx, cancelCron := context.WithCancel(serverCtx)
	defer cancelCron()
	initialize.InitCronScheduler(cronCtx)
	defer initialize.ShutdownCronScheduler()

	// 设置Gin模式
	gin.SetMode(global.Config.Server.Mode)

	// 初始化路由
	r := router.InitRouter()

	// 启动服务
	addr := fmt.Sprintf(":%d", global.Config.Server.Port)
	global.Log.Infof("服务启动成功，监听端口: %s", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		<-serverCtx.Done()
		cancelCron()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			global.Log.Errorf("服务优雅停止失败: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		global.Log.Fatalf("服务启动失败: %v", err)
	}
}
