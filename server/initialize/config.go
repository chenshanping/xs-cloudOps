package initialize

import (
	"fmt"
	"server/config"
	"server/global"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string) (*config.Config, *viper.Viper, error) {
	v := viper.New()
	if configPath == "" {
		configPath = "config.yaml"
	}
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, v, nil
}

func InitConfig(configPath string) {
	cfg, v, err := LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	global.Config = cfg
	global.Viper = v
}
