package configsvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"

	"server/config"
	"server/global"
	"server/model"
	"server/service/filesvc"
	"server/service/oss"
)

type ConfigService struct{}

var Default = &ConfigService{}

const PublicConfigKeysConfigKey = "public_config_keys"

const (
	SysLogoFileIDConfigKey         = "sys_logo_file_id"
	RegisterLogoFileIDConfigKey    = "register_logo_file_id"
	LoginBGImageFileIDConfigKey    = "login_bg_image_file_id"
	SliderCaptchaBgFileIDConfigKey = "slider_captcha_bg_file_id"
)

var imageFileReferenceConfigLabels = map[string]string{
	SysLogoFileIDConfigKey:         "系统 Logo",
	RegisterLogoFileIDConfigKey:    "注册默认头像",
	LoginBGImageFileIDConfigKey:    "登录页背景图",
	SliderCaptchaBgFileIDConfigKey: "滑动验证码背景",
}

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
	if err := db.Order("id asc").Find(&configs).Error; err != nil {
		return nil, err
	}
	return s.resolveConfigList(configs)
}

// GetConfigByKey 根据key获取配置
func (s *ConfigService) GetConfigByKey(key string) (*model.SysConfig, error) {
	requestKeys := []string{key}
	if fileIDKey, ok := ImageURLToFileIDKeyMap()[key]; ok {
		requestKeys = append(requestKeys, fileIDKey)
	}

	configs, err := s.GetConfigsByKeys(requestKeys)
	if err != nil {
		return nil, err
	}
	config, ok := configs[key]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &config, nil
}

// GetConfigsByKeys 批量获取配置
func (s *ConfigService) GetConfigsByKeys(keys []string) (map[string]model.SysConfig, error) {
	normalizedKeys := normalizeConfigKeys(keys)
	queryKeys := expandConfigQueryKeys(normalizedKeys)
	if len(queryKeys) == 0 {
		return map[string]model.SysConfig{}, nil
	}

	var configs []model.SysConfig
	if err := global.DB.Where("`key` IN ?", queryKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	result := make(map[string]model.SysConfig, len(configs))
	for _, c := range configs {
		result[c.Key] = c
	}
	result, err := s.ResolveImageConfigURLs(result)
	if err != nil {
		return nil, err
	}

	filtered := make(map[string]model.SysConfig, len(normalizedKeys))
	for _, key := range normalizedKeys {
		if config, ok := result[key]; ok {
			filtered[key] = config
		}
	}
	return filtered, nil
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

func ImageFileReferenceConfigKeys() []string {
	return []string{
		SysLogoFileIDConfigKey,
		RegisterLogoFileIDConfigKey,
		LoginBGImageFileIDConfigKey,
		SliderCaptchaBgFileIDConfigKey,
	}
}

// ImageFileIDToURLKeyMap 返回 file_id 配置键 → URL 配置键的映射。
// 迁移文件后需要同步更新 URL 配置键的值。
func ImageFileIDToURLKeyMap() map[string]string {
	return map[string]string{
		SysLogoFileIDConfigKey:         "sys_logo",
		RegisterLogoFileIDConfigKey:    "register_logo",
		LoginBGImageFileIDConfigKey:    "login_bg_image",
		SliderCaptchaBgFileIDConfigKey: "slider_captcha_bg",
	}
}

func ImageURLToFileIDKeyMap() map[string]string {
	result := make(map[string]string, len(ImageFileIDToURLKeyMap()))
	for fileIDKey, urlKey := range ImageFileIDToURLKeyMap() {
		result[urlKey] = fileIDKey
	}
	return result
}

func ImageFileReferenceLabel(key string) string {
	if label, ok := imageFileReferenceConfigLabels[key]; ok {
		return label
	}
	return key
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
	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		updatedConfigs := make(map[string]model.SysConfig, len(configs))
		for key, value := range configs {
			var config model.SysConfig
			err := tx.Where("`key` = ?", key).First(&config).Error
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				config = model.SysConfig{
					Key:   key,
					Value: value,
				}
				if err := tx.Create(&config).Error; err != nil {
					return err
				}
			case err != nil:
				return err
			default:
				if err := tx.Model(&config).Update("value", value).Error; err != nil {
					return err
				}
				config.Value = value
			}
			updatedConfigs[key] = config
		}

		return s.syncConfigImageFileRefs(tx, updatedConfigs)
	}); err != nil {
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

func (s *ConfigService) ResolveImageConfigURLs(configs map[string]model.SysConfig) (map[string]model.SysConfig, error) {
	resolved := make(map[string]model.SysConfig, len(configs)+len(ImageFileIDToURLKeyMap()))
	for key, config := range configs {
		resolved[key] = config
	}

	fileIDByKey := make(map[string]uint, len(ImageFileIDToURLKeyMap()))
	fileIDs := make([]uint, 0, len(ImageFileIDToURLKeyMap()))
	for fileIDKey := range ImageFileIDToURLKeyMap() {
		config, ok := resolved[fileIDKey]
		if !ok {
			continue
		}
		fileID, ok := parseConfigFileID(config.Value)
		if !ok {
			continue
		}
		fileIDByKey[fileIDKey] = fileID
		fileIDs = append(fileIDs, fileID)
	}

	fileURLMap := make(map[uint]string, len(fileIDs))
	if len(fileIDs) > 0 {
		var files []model.SysFile
		if err := global.DB.Select("id", "url").Where("id IN ? AND status = ?", deduplicateUint(fileIDs), 1).Find(&files).Error; err != nil {
			return nil, err
		}
		for _, file := range files {
			fileURLMap[file.ID] = file.URL
		}
	}

	for fileIDKey, urlKey := range ImageFileIDToURLKeyMap() {
		derived := resolved[urlKey]
		derived.Key = urlKey
		derived.Value = ""
		if fileID, ok := fileIDByKey[fileIDKey]; ok {
			derived.Value = fileURLMap[fileID]
		}
		resolved[urlKey] = derived
	}

	return resolved, nil
}

func (s *ConfigService) resolveConfigList(configs []model.SysConfig) ([]model.SysConfig, error) {
	configMap := make(map[string]model.SysConfig, len(configs))
	for _, config := range configs {
		configMap[config.Key] = config
	}

	resolved, err := s.ResolveImageConfigURLs(configMap)
	if err != nil {
		return nil, err
	}

	existingKeys := make(map[string]struct{}, len(configs))
	for i := range configs {
		configs[i] = resolved[configs[i].Key]
		existingKeys[configs[i].Key] = struct{}{}
	}

	urlKeys := make([]string, 0, len(ImageURLToFileIDKeyMap()))
	for urlKey := range ImageURLToFileIDKeyMap() {
		urlKeys = append(urlKeys, urlKey)
	}
	sort.Strings(urlKeys)
	for _, urlKey := range urlKeys {
		if _, exists := existingKeys[urlKey]; exists {
			continue
		}
		if config, ok := resolved[urlKey]; ok {
			configs = append(configs, config)
		}
	}

	return configs, nil
}

func (s *ConfigService) syncConfigImageFileRefs(tx *gorm.DB, updatedConfigs map[string]model.SysConfig) error {
	for _, fileIDKey := range ImageFileReferenceConfigKeys() {
		config, ok := updatedConfigs[fileIDKey]
		if !ok {
			continue
		}

		fileID, ok := parseConfigFileID(config.Value)
		if !ok {
			if err := filesvc.Reference.ClearRefs(tx, "sys_config", config.ID); err != nil {
				return err
			}
			continue
		}

		var file model.SysFile
		if err := tx.Select("id").Where("id = ? AND status = ?", fileID, 1).First(&file).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%s文件不存在", ImageFileReferenceLabel(fileIDKey))
			}
			return err
		}

		if err := filesvc.Reference.ReplaceRefs(tx, "sys_config", config.ID, []filesvc.FileRef{{
			FileID: file.ID,
			Field:  config.Key,
		}}); err != nil {
			return err
		}
	}

	return nil
}

func expandConfigQueryKeys(keys []string) []string {
	result := append([]string(nil), keys...)
	urlToFileID := ImageURLToFileIDKeyMap()
	seen := buildConfigKeySet(result)
	for _, key := range keys {
		fileIDKey, ok := urlToFileID[key]
		if !ok {
			continue
		}
		if _, exists := seen[fileIDKey]; exists {
			continue
		}
		result = append(result, fileIDKey)
		seen[fileIDKey] = struct{}{}
	}
	return result
}

func parseConfigFileID(rawValue string) (uint, bool) {
	rawValue = strings.TrimSpace(rawValue)
	if rawValue == "" || rawValue == "0" {
		return 0, false
	}

	var fileID uint
	for _, ch := range rawValue {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		fileID = fileID*10 + uint(ch-'0')
	}
	if fileID == 0 {
		return 0, false
	}
	return fileID, true
}

func deduplicateUint(values []uint) []uint {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[uint]struct{}, len(values))
	result := make([]uint, 0, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
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
