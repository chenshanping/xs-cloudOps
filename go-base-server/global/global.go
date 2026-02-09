package global

import (
	"go-base-server/config"

	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config   *config.Config
	Viper    *viper.Viper
	Log      *zap.SugaredLogger
	DB       *gorm.DB
	Redis    *redis.Client
	Enforcer *casbin.Enforcer
)
