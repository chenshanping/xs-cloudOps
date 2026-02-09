package initialize

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"go-base-server/global"
)

func InitRedis() {
	cfg := global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(fmt.Errorf("连接Redis失败: %w", err))
	}

	global.Redis = client
	global.Log.Info("Redis连接成功")
}
