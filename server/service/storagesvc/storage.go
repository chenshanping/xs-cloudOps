package storagesvc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"server/model"
	"server/service/configsvc"
	"server/service/oss"
)

type StorageService struct{}

var Default = &StorageService{}

func StorageConfigKey(storageType model.StorageType) string {
	if storageType == "" {
		storageType = Default.DefaultStorageType()
	}
	return fmt.Sprintf("storage_%s_config", storageType)
}

func (s *StorageService) SupportedStorageTypes() []model.StorageType {
	return []model.StorageType{
		model.StorageTypeLocal,
		model.StorageTypeAliyun,
		model.StorageTypeTencent,
		model.StorageTypeMinio,
	}
}

func (s *StorageService) DefaultStorageType() model.StorageType {
	return model.StorageTypeLocal
}

const defaultLocalStorageConfig = `{"base_path":"uploads","base_url":"/api/v1/upload"}`

func (s *StorageService) DefaultStorageConfig(storageType model.StorageType) string {
	switch storageType {
	case model.StorageTypeLocal:
		return defaultLocalStorageConfig
	default:
		return "{}"
	}
}

func (s *StorageService) NormalizeStorageConfig(storageType model.StorageType, configJSON string) (string, error) {
	if strings.TrimSpace(configJSON) == "" {
		configJSON = s.DefaultStorageConfig(storageType)
	}

	switch storageType {
	case model.StorageTypeLocal:
		var config model.LocalConfig
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			return "", fmt.Errorf("解析本地存储配置失败: %v", err)
		}
		if strings.TrimSpace(config.BasePath) == "" {
			config.BasePath = "uploads"
		}
		if strings.TrimSpace(config.BaseURL) == "" {
			config.BaseURL = "/api/v1/upload"
		}
		normalized, err := json.Marshal(config)
		if err != nil {
			return "", err
		}
		return string(normalized), nil
	case model.StorageTypeAliyun:
		var config model.AliyunOSSConfig
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			return "", fmt.Errorf("解析阿里云存储配置失败: %v", err)
		}
		normalized, err := json.Marshal(config)
		if err != nil {
			return "", err
		}
		return string(normalized), nil
	case model.StorageTypeTencent:
		var config model.TencentCOSConfig
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			return "", fmt.Errorf("解析腾讯云存储配置失败: %v", err)
		}
		normalized, err := json.Marshal(config)
		if err != nil {
			return "", err
		}
		return string(normalized), nil
	case model.StorageTypeMinio:
		var config model.MinioConfig
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			return "", fmt.Errorf("解析 MinIO 存储配置失败: %v", err)
		}
		normalized, err := json.Marshal(config)
		if err != nil {
			return "", err
		}
		return string(normalized), nil
	default:
		return "", fmt.Errorf("不支持的存储类型: %s", storageType)
	}
}

func (s *StorageService) BuildStorageProfile(storageType model.StorageType, configJSON string) (*model.StorageProfile, error) {
	normalized, err := s.NormalizeStorageConfig(storageType, configJSON)
	if err != nil {
		return nil, err
	}

	return &model.StorageProfile{
		Name:   "系统配置",
		Type:   storageType,
		Config: normalized,
	}, nil
}

func (s *StorageService) GetStorageByType(storageType model.StorageType) (*model.StorageProfile, error) {
	if storageType == "" {
		storageType = s.DefaultStorageType()
	}

	configs, err := configsvc.Default.GetConfigsByKeys([]string{
		StorageConfigKey(storageType),
		configsvc.LegacyStorageConfigConfigKey,
	})
	if err != nil {
		return nil, err
	}

	configValue := ""
	if config, ok := configs[StorageConfigKey(storageType)]; ok && strings.TrimSpace(config.Value) != "" {
		configValue = config.Value
	} else if legacyConfig, ok := configs[configsvc.LegacyStorageConfigConfigKey]; ok {
		configValue = legacyConfig.Value
	}

	return s.BuildStorageProfile(storageType, configValue)
}

func (s *StorageService) GetSystemStorage() (*model.StorageProfile, error) {
	configs, err := configsvc.Default.GetConfigsByKeys([]string{configsvc.StorageTypeConfigKey})
	if err != nil {
		return nil, err
	}

	storageType := s.DefaultStorageType()
	if config, ok := configs[configsvc.StorageTypeConfigKey]; ok && strings.TrimSpace(config.Value) != "" {
		storageType = model.StorageType(config.Value)
	}

	return s.GetStorageByType(storageType)
}

func (s *StorageService) GetDefaultStorage() (*model.StorageProfile, error) {
	return s.GetSystemStorage()
}

func (s *StorageService) UpdateSystemStorage(storageType model.StorageType, configJSON string) error {
	storage, err := s.BuildStorageProfile(storageType, configJSON)
	if err != nil {
		return err
	}
	if err := s.TestStorage(storage); err != nil {
		return err
	}
	if err := configsvc.Default.BatchUpdateConfigs(map[string]string{
		configsvc.StorageTypeConfigKey:       string(storage.Type),
		StorageConfigKey(storage.Type): storage.Config,
	}); err != nil {
		return err
	}

	oss.ClearClients()
	return nil
}

func (s *StorageService) TestStorage(storage *model.StorageProfile) error {
	client, err := oss.NewClient(storage)
	if err != nil {
		return err
	}

	if storage.Type == model.StorageTypeLocal {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := "__storage_test_object__"
	url, err := client.GetSignedURL(ctx, key, 30*time.Second)
	if err != nil {
		return err
	}

	if strings.HasPrefix(url, "/") {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("认证失败，请检查 AccessKey/Secret、Bucket、Region 等配置 (status=%d)", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("测试请求失败，状态码: %d", resp.StatusCode)
	}

	return nil
}
