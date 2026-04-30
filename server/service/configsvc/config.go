package configsvc

import (
	"encoding/json"
	"strings"

	"server/config"
	"server/global"
	"server/model"
	"server/service/oss"
)

type ConfigService struct{}

var Default = &ConfigService{}

const PublicConfigKeysConfigKey = "public_config_keys"

// defaultPublicConfigKeys 定义匿名接口 /api/v1/configs/keys 默认允许返回的配置键。
// 这里应只放“前台展示必需、且不包含敏感信息”的配置。
var defaultPublicConfigKeys = []string{
	// 系统基础
	"sys_name",
	"sys_logo",
	// 注册/登录页展示
	"register_logo",
	"login_bg_image",
	"login_title",
	"login_subtitle",
	"login_bg_color",
	"login_slogan",
	"login_desc",
	"login_features",
	"login_features_max",
	"login_images",
	"login_images_max",
	"enable_register",
	// 前台模式
	"front_mode",
	"user_profile_button_visible",
}

// neverPublicConfigKeys 定义无论任何场景都禁止匿名公开的配置键。
// 即使未来误把这些键加入 defaultPublicConfigKeys，也会在运行时被再次拦截。
var neverPublicConfigKeys = []string{
	// AI / 密钥类
	"ai_config",
	"email_username",
	"email_password",
	// 用户与安全类
	"user_default_password",
	// 存储基础设施类
	"storage_type",
	"storage_local_config",
	"storage_aliyun_config",
	"storage_tencent_config",
	"storage_minio_config",
	// 管理白名单自身也不允许匿名读取
	PublicConfigKeysConfigKey,
}

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

func DefaultPublicConfigKeys() []string {
	return append([]string(nil), defaultPublicConfigKeys...)
}

func DefaultPublicConfigKeysValue() string {
	bytes, err := json.Marshal(DefaultPublicConfigKeys())
	if err != nil {
		return "[]"
	}
	return string(bytes)
}

// GetResolvedPublicConfigKeys 返回最终生效的匿名公开配置键。
// 当前项目策略：只认后端代码定义，不再读取数据库中的自定义白名单。
func (s *ConfigService) GetResolvedPublicConfigKeys() []string {
	return sanitizePublicConfigKeys(DefaultPublicConfigKeys())
}

// GetPublicConfigsByKeys 批量获取允许匿名读取的公开配置
func (s *ConfigService) GetPublicConfigsByKeys(keys []string) (map[string]model.SysConfig, error) {
	filteredKeys := filterPublicConfigKeys(keys, s.GetResolvedPublicConfigKeys())
	if len(filteredKeys) == 0 {
		return map[string]model.SysConfig{}, nil
	}
	return s.GetConfigsByKeys(filteredKeys)
}

func sanitizePublicConfigKeys(keys []string) []string {
	neverPublic := buildConfigKeySet(neverPublicConfigKeys)
	filtered := make([]string, 0, len(keys))
	for _, key := range normalizeConfigKeys(keys) {
		if _, denied := neverPublic[key]; denied {
			continue
		}
		filtered = append(filtered, key)
	}
	return filtered
}

func filterPublicConfigKeys(keys []string, publicKeys []string) []string {
	if len(keys) == 0 {
		return []string{}
	}

	allowed := buildConfigKeySet(publicKeys)
	filtered := make([]string, 0, len(keys))
	for _, key := range normalizeConfigKeys(keys) {
		if _, ok := allowed[key]; ok {
			filtered = append(filtered, key)
		}
	}
	return filtered
}

func normalizeConfigKeys(keys []string) []string {
	normalized := make([]string, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		normalized = append(normalized, key)
		seen[key] = struct{}{}
	}
	return normalized
}

func buildConfigKeySet(keys []string) map[string]struct{} {
	set := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		set[key] = struct{}{}
	}
	return set
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
			config = model.SysConfig{
				Key:   key,
				Value: value,
			}
			if err := tx.Create(&config).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.Model(&config).Update("value", value).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	for key := range configs {
		if key == StorageTypeConfigKey || key == LegacyStorageConfigConfigKey || strings.HasPrefix(key, "storage_") && strings.HasSuffix(key, "_config") {
			oss.ClearClients()
			break
		}
	}

	return nil
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
	return &global.Config.AI
}

// Constants originally from storage.go, needed by config.go for storage config detection.
const (
	StorageTypeConfigKey         = "storage_type"
	LegacyStorageConfigConfigKey = "storage_config"
)
