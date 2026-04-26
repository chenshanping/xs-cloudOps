package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"server/global"
	"server/model"
	"server/service/oss"
)

type StorageService struct{}

var Storage = new(StorageService)

// GetStorageList 获取存储配置列表
func (s *StorageService) GetStorageList() ([]model.SysStorage, error) {
	var storages []model.SysStorage
	err := global.DB.Order("id asc").Find(&storages).Error
	return storages, err
}

// GetStorageByID 根据ID获取存储配置
func (s *StorageService) GetStorageByID(id uint) (*model.SysStorage, error) {
	var storage model.SysStorage
	err := global.DB.First(&storage, id).Error
	if err != nil {
		return nil, err
	}
	return &storage, nil
}

// GetDefaultStorage 获取默认存储配置
func (s *StorageService) GetDefaultStorage() (*model.SysStorage, error) {
	var storage model.SysStorage
	err := global.DB.Where("is_default = ? AND status = ?", 1, 1).First(&storage).Error
	if err != nil {
		// 如果没有默认存储，尝试获取第一个启用的存储
		err = global.DB.Where("status = ?", 1).First(&storage).Error
		if err != nil {
			return nil, err
		}
	}
	return &storage, nil
}

// CreateStorage 创建存储配置
func (s *StorageService) CreateStorage(storage *model.SysStorage) error {
	// 如果设为默认，先取消其他默认
	if storage.IsDefault == 1 {
		global.DB.Model(&model.SysStorage{}).Where("is_default = ?", 1).Update("is_default", 0)
	}
	return global.DB.Create(storage).Error
}

// UpdateStorage 更新存储配置
func (s *StorageService) UpdateStorage(id uint, data map[string]interface{}) error {
	// 清除缓存的客户端
	oss.RemoveClient(id)

	// 如果设为默认，先取消其他默认
	if isDefault, ok := data["is_default"]; ok && isDefault == 1 {
		global.DB.Model(&model.SysStorage{}).Where("is_default = ? AND id != ?", 1, id).Update("is_default", 0)
	}
	return global.DB.Model(&model.SysStorage{}).Where("id = ?", id).Updates(data).Error
}

// DeleteStorage 删除存储配置
func (s *StorageService) DeleteStorage(id uint) error {
	// 检查是否有文件使用此存储
	var count int64
	global.DB.Model(&model.SysFile{}).Where("storage_id = ?", id).Count(&count)
	if count > 0 {
		return global.DB.Model(&model.SysStorage{}).Where("id = ?", id).Update("status", 0).Error
	}

	// 清除缓存的客户端
	oss.RemoveClient(id)

	return global.DB.Delete(&model.SysStorage{}, id).Error
}

// SetDefaultStorage 设置默认存储
func (s *StorageService) SetDefaultStorage(id uint) error {
	// 取消所有默认
	global.DB.Model(&model.SysStorage{}).Where("is_default = ?", 1).Update("is_default", 0)
	// 设置新默认
	return global.DB.Model(&model.SysStorage{}).Where("id = ?", id).Update("is_default", 1).Error
}

// TestStorage 测试存储配置是否可用
// 注意：仅创建客户端无法验证 AccessKey/Secret 是否正确，这里增加一次真实网络请求校验。
func (s *StorageService) TestStorage(storage *model.SysStorage) error {
	// 先尝试创建客户端，校验基础配置是否正确（JSON 结构、必填字段等）
	client, err := oss.NewClient(storage)
	if err != nil {
		return err
	}

	// 本地存储不需要额外校验，能创建客户端就算成功
	if storage.Type == model.StorageTypeLocal {
		return nil
	}

	// 远程对象存储：通过带签名的 URL 发起一次 GET 请求，验证签名/权限是否正确
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用一个不存在的 key 即可，只要能正常鉴权，返回 404 也是成功
	key := "__storage_test_object__"
	url, err := client.GetSignedURL(ctx, key, 30*time.Second)
	if err != nil {
		return err
	}

	// 某些实现返回的是后端相对路径（如本地/MinIO 代理），这类走后端逻辑的在这里不再做额外 HTTP 校验
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

	// 403/401 基本可以判定为签名/鉴权失败
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		fmt.Println(resp.Header)
		return fmt.Errorf("认证失败，请检查 AccessKey/Secret、Bucket、Region 等配置 (status=%d)", resp.StatusCode)
	}

	// 404 代表 bucket 存在且鉴权通过，只是对象不存在，这里视为成功
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	// 其他 4xx/5xx 状态也当作失败返回
	if resp.StatusCode >= 400 {
		return fmt.Errorf("测试请求失败，状态码: %d", resp.StatusCode)
	}

	return nil
}
