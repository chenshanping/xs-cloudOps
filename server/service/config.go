package service

import (
	"encoding/json"

	"server/config"
	"server/global"
	"server/model"
)

type ConfigService struct{}

var Config = new(ConfigService)

// GetConfigList 获取配置列表
func (s *ConfigService) GetConfigList(key string) ([]model.SysConfig, error) {
	var configs []model.SysConfig
	db := global.DB.Model(&model.SysConfig{})
	if key != "" {
		db = db.Where("`key` LIKE ?", "%"+key+"%")
	}
	err := db.Order("id asc").Find(&configs).Error
	return configs, err
}

// GetConfigByKey 根据key获取配置
func (s *ConfigService) GetConfigByKey(key string) (*model.SysConfig, error) {
	var config model.SysConfig
	err := global.DB.Where("`key` = ?", key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetConfigsByKeys 批量获取配置
func (s *ConfigService) GetConfigsByKeys(keys []string) (map[string]model.SysConfig, error) {
	var configs []model.SysConfig
	err := global.DB.Where("`key` IN ?", keys).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]model.SysConfig)
	for _, c := range configs {
		result[c.Key] = c
	}
	return result, nil
}

// CreateConfig 创建配置
func (s *ConfigService) CreateConfig(config *model.SysConfig) error {
	return global.DB.Create(config).Error
}

// UpdateConfig 更新配置
func (s *ConfigService) UpdateConfig(id uint, data map[string]interface{}) error {
	return global.DB.Model(&model.SysConfig{}).Where("id = ?", id).Updates(data).Error
}

// UpdateConfigByKey 根据key更新配置值
func (s *ConfigService) UpdateConfigByKey(key string, value string) error {
	return global.DB.Model(&model.SysConfig{}).Where("`key` = ?", key).Update("value", value).Error
}

// BatchUpdateConfigs 批量更新配置（不存在则创建）
func (s *ConfigService) BatchUpdateConfigs(configs map[string]string) error {
	tx := global.DB.Begin()
	for key, value := range configs {
		var config model.SysConfig
		err := tx.Where("`key` = ?", key).First(&config).Error
		if err != nil {
			// key 不存在，创建新配置
			config = model.SysConfig{
				Key:   key,
				Value: value,
			}
			if err := tx.Create(&config).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// key 存在，更新值
			if err := tx.Model(&config).Update("value", value).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

// DeleteConfig 删除配置
func (s *ConfigService) DeleteConfig(id uint) error {
	return global.DB.Delete(&model.SysConfig{}, id).Error
}

// GetAIConfig 获取AI配置（优先从数据库读取，兜底用配置文件）
func (s *ConfigService) GetAIConfig() *config.AI {
	cfg, err := s.GetConfigByKey("ai_config")
	if err == nil && cfg.Value != "" {
		var aiConfig config.AI
		if err := json.Unmarshal([]byte(cfg.Value), &aiConfig); err == nil {
			return &aiConfig
		}
	}
	// 兜底使用配置文件
	return &global.Config.AI
}
