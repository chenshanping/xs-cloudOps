package config

import "fmt"

type Config struct {
	Server Server `mapstructure:"server" yaml:"server"`
	MySQL  MySQL  `mapstructure:"mysql" yaml:"mysql"`
	Redis  Redis  `mapstructure:"redis" yaml:"redis"`
	JWT    JWT    `mapstructure:"jwt" yaml:"jwt"`
	Casbin Casbin `mapstructure:"casbin" yaml:"casbin"`
	Log    Log    `mapstructure:"log" yaml:"log"`
	AI     AI     `mapstructure:"ai" yaml:"ai"`
}

type Server struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port int    `mapstructure:"port" yaml:"port"`
	Mode string `mapstructure:"mode" yaml:"mode"`
}

type MySQL struct {
	Host         string `mapstructure:"host" yaml:"host"`
	Port         int    `mapstructure:"port" yaml:"port"`
	Username     string `mapstructure:"username" yaml:"username"`
	Password     string `mapstructure:"password" yaml:"password"`
	DBName       string `mapstructure:"dbname" yaml:"dbname"`
	Charset      string `mapstructure:"charset" yaml:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
}

func (m *MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", m.Username, m.Password, m.Host, m.Port, m.DBName, m.Charset)
}

type Redis struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
}

type JWT struct {
	Secret        string `mapstructure:"secret" yaml:"secret"`
	Expires       int64  `mapstructure:"expires" yaml:"expires"`
	RefreshWindow int64  `mapstructure:"refresh_window" yaml:"refresh_window"` // Token过期后允许刷新的时间窗口(秒)
	Issuer        string `mapstructure:"issuer" yaml:"issuer"`
}

type Casbin struct {
	ModelPath string `mapstructure:"model_path" yaml:"model_path"`
}

type Log struct {
	Level      string `mapstructure:"level" yaml:"level"`
	Format     string `mapstructure:"format" yaml:"format"`
	Directory  string `mapstructure:"directory" yaml:"directory"`
	Filename   string `mapstructure:"filename" yaml:"filename"`
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`
	Compress   bool   `mapstructure:"compress" yaml:"compress"`
	Stdout     bool   `mapstructure:"stdout" yaml:"stdout"`
}

type AI struct {
	DefaultProvider string       `mapstructure:"default_provider" yaml:"default_provider" json:"default_provider"`
	Providers       []AIProvider `mapstructure:"providers" yaml:"providers" json:"providers"`
}

type AIProvider struct {
	Name    string    `mapstructure:"name" yaml:"name" json:"name"`
	APIKey  string    `mapstructure:"api_key" yaml:"api_key" json:"api_key"`
	BaseURL string    `mapstructure:"base_url" yaml:"base_url" json:"base_url"`
	Models  []AIModel `mapstructure:"models" yaml:"models" json:"models"`
}

type AIModel struct {
	ID               string   `mapstructure:"id" yaml:"id" json:"id"`
	Name             string   `mapstructure:"name" yaml:"name" json:"name"`
	Group            string   `mapstructure:"group" yaml:"group" json:"group,omitempty"`
	Description      string   `mapstructure:"description" yaml:"description" json:"description"`
	IsThinking       bool     `mapstructure:"is_thinking" yaml:"is_thinking" json:"is_thinking"`
	SupportVision    bool     `mapstructure:"support_vision" yaml:"support_vision" json:"support_vision"`
	SupportTools     bool     `mapstructure:"support_tools" yaml:"support_tools" json:"support_tools"`
	SearchStrategy   string   `mapstructure:"search_strategy" yaml:"search_strategy" json:"search_strategy"`
	SupportEmbedding bool     `mapstructure:"support_embedding" yaml:"support_embedding" json:"support_embedding"`
	SupportRerank    bool     `mapstructure:"support_rerank" yaml:"support_rerank" json:"support_rerank"`
	IsFree           bool     `mapstructure:"is_free" yaml:"is_free" json:"is_free"`
	Temperature      *float64 `mapstructure:"temperature" yaml:"temperature" json:"temperature,omitempty"`
	ContextWindow    *int     `mapstructure:"context_window" yaml:"context_window" json:"context_window,omitempty"`
	Tags             []string `mapstructure:"tags" yaml:"tags" json:"tags,omitempty"`
}

// GetProvider 根据名称获取平台配置
func (ai *AI) GetProvider(name string) *AIProvider {
	for i := range ai.Providers {
		if ai.Providers[i].Name == name {
			return &ai.Providers[i]
		}
	}
	return nil
}

// GetDefaultProvider 获取默认平台配置
func (ai *AI) GetDefaultProvider() *AIProvider {
	if ai.DefaultProvider != "" {
		if p := ai.GetProvider(ai.DefaultProvider); p != nil {
			return p
		}
	}
	// 如果没有默认平台或找不到，返回第一个
	if len(ai.Providers) > 0 {
		return &ai.Providers[0]
	}
	return nil
}

// GetProviderByModel 根据模型ID找到对应的平台
func (ai *AI) GetProviderByModel(modelID string) *AIProvider {
	for i := range ai.Providers {
		for _, m := range ai.Providers[i].Models {
			if m.ID == modelID {
				return &ai.Providers[i]
			}
		}
	}
	return nil
}
