package ai

import (
	"encoding/json"
	"errors"
	"strings"

	appconfig "server/config"
	"server/global"
	"server/model"

	"gorm.io/gorm"
)

const legacyAdminAIConfigKey = "ai_config"

var emptyAdminAIConfig = appconfig.AI{}

func (s *AIService) GetAdminConfig() (*appconfig.AI, error) {
	cfg, found, err := s.getAdminConfigFromProviders()
	if err != nil {
		return nil, err
	}
	if found {
		return cfg, nil
	}

	cfg, found, err = s.getAdminConfigFromLegacyConfig()
	if err != nil {
		return nil, err
	}
	if found {
		if err := s.replaceAIProvidersFromConfig(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	empty := emptyAdminAIConfig
	return &empty, nil
}

func (s *AIService) SaveAdminConfig(cfg *appconfig.AI) error {
	if cfg == nil {
		cfg = &appconfig.AI{}
	}
	return s.replaceAIProvidersFromConfig(cfg)
}

func (s *AIService) getAdminConfigFromProviders() (*appconfig.AI, bool, error) {
	var records []model.AIProviderConfig
	if err := global.DB.Order("sort ASC, id ASC").Find(&records).Error; err != nil {
		return nil, false, err
	}
	if len(records) == 0 {
		return nil, false, nil
	}

	cfg, err := buildAdminConfigFromProviderRows(records)
	if err != nil {
		return nil, false, err
	}
	return cfg, true, nil
}

func (s *AIService) getAdminConfigFromLegacyConfig() (*appconfig.AI, bool, error) {
	var record model.SysConfig
	err := global.DB.Where("`key` = ?", legacyAdminAIConfigKey).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if strings.TrimSpace(record.Value) == "" {
		return nil, false, nil
	}

	cfg, err := unmarshalAdminConfig(record.Value)
	if err != nil {
		return nil, false, err
	}
	return cfg, true, nil
}

func (s *AIService) replaceAIProvidersFromConfig(cfg *appconfig.AI) error {
	rows, err := buildProviderRowsFromAdminConfig(cfg)
	if err != nil {
		return err
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.AIProviderConfig{}).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		return tx.Create(&rows).Error
	})
}

func buildAdminConfigFromProviderRows(records []model.AIProviderConfig) (*appconfig.AI, error) {
	cfg := &appconfig.AI{
		Providers: make([]appconfig.AIProvider, 0, len(records)),
	}

	for _, record := range records {
		models := make([]appconfig.AIModel, 0)
		if strings.TrimSpace(record.ModelsJSON) != "" {
			if err := json.Unmarshal([]byte(record.ModelsJSON), &models); err != nil {
				return nil, err
			}
		}

		cfg.Providers = append(cfg.Providers, appconfig.AIProvider{
			Name:    record.Name,
			APIKey:  record.APIKey,
			BaseURL: record.BaseURL,
			Models:  models,
		})
		if record.IsDefault && cfg.DefaultProvider == "" {
			cfg.DefaultProvider = record.Name
		}
	}

	if cfg.DefaultProvider == "" && len(cfg.Providers) > 0 {
		cfg.DefaultProvider = cfg.Providers[0].Name
	}

	return cfg, nil
}

func buildProviderRowsFromAdminConfig(cfg *appconfig.AI) ([]model.AIProviderConfig, error) {
	if cfg == nil || len(cfg.Providers) == 0 {
		return nil, nil
	}

	defaultName := strings.TrimSpace(cfg.DefaultProvider)
	if defaultName == "" {
		defaultName = cfg.Providers[0].Name
	}

	hasDefault := false
	for _, provider := range cfg.Providers {
		if provider.Name == defaultName {
			hasDefault = true
			break
		}
	}
	if !hasDefault {
		defaultName = cfg.Providers[0].Name
	}

	rows := make([]model.AIProviderConfig, 0, len(cfg.Providers))
	for index, provider := range cfg.Providers {
		modelsJSON, err := json.Marshal(provider.Models)
		if err != nil {
			return nil, err
		}
		rows = append(rows, model.AIProviderConfig{
			Name:       provider.Name,
			APIKey:     provider.APIKey,
			BaseURL:    provider.BaseURL,
			ModelsJSON: string(modelsJSON),
			IsDefault:  provider.Name == defaultName,
			Sort:       index,
		})
	}

	return rows, nil
}

func unmarshalAdminConfig(raw string) (*appconfig.AI, error) {
	var cfg appconfig.AI
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
