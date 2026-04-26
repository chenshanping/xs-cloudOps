package initialize

import (
	"fmt"
	"server/config"
	"server/global"

	"github.com/spf13/viper"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("解析配置文件失败: %w", err))
	}

	global.Config = &cfg
	global.Viper = v
}
