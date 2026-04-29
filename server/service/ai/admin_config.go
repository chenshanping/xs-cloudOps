package ai

import (
	"encoding/json"
	"errors"

	appconfig "server/config"
	"server/global"
	"server/model"

	"gorm.io/gorm"
)

const adminAIConfigKey = "ai_config"

var emptyAdminAIConfig = appconfig.AI{}

func (s *AIService) GetAdminConfig() (*appconfig.AI, error) {
	var record model.SysConfig
	err := global.DB.Where("`key` = ?", adminAIConfigKey).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cfg := emptyAdminAIConfig
		return &cfg, nil
	}
	if err != nil {
		return nil, err
	}
	if record.Value == "" {
		cfg := emptyAdminAIConfig
		return &cfg, nil
	}

	var cfg appconfig.AI
	if err := json.Unmarshal([]byte(record.Value), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *AIService) SaveAdminConfig(cfg *appconfig.AI) error {
	if cfg == nil {
		cfg = &appconfig.AI{}
	}

	value, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	var record model.SysConfig
	err = global.DB.Where("`key` = ?", adminAIConfigKey).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		record = model.SysConfig{
			Name:      "AI配置",
			Key:       adminAIConfigKey,
			Value:     string(value),
			ValueType: "json",
			Remark:    "AI平台配置，包含平台名称、API Key、基础URL和模型列表",
		}
		return global.DB.Create(&record).Error
	}
	if err != nil {
		return err
	}

	return global.DB.Model(&record).Updates(map[string]interface{}{
		"value":      string(value),
		"value_type": "json",
	}).Error
}
