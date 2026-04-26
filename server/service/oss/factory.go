package oss

import (
	"fmt"
	"sync"

	"server/model"
)

var (
	clientCache = make(map[uint]Client)
	cacheMutex  sync.RWMutex
)

// NewClient 根据存储配置创建客户端
func NewClient(storage *model.SysStorage) (Client, error) {
	switch model.StorageType(storage.Type) {
	case model.StorageTypeLocal:
		return NewLocalClient(storage.Config)
	case model.StorageTypeAliyun:
		return NewAliyunClient(storage.Config)
	case model.StorageTypeTencent:
		return NewTencentClient(storage.Config)
	case model.StorageTypeMinio:
		return NewMinioClient(storage.Config)
	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", storage.Type)
	}
}

// GetClient 获取缓存的客户端，如果不存在则创建
func GetClient(storage *model.SysStorage) (Client, error) {
	cacheMutex.RLock()
	if client, ok := clientCache[storage.ID]; ok {
		cacheMutex.RUnlock()
		return client, nil
	}
	cacheMutex.RUnlock()

	// 创建新客户端
	client, err := NewClient(storage)
	if err != nil {
		return nil, err
	}

	// 缓存客户端
	cacheMutex.Lock()
	clientCache[storage.ID] = client
	cacheMutex.Unlock()

	return client, nil
}

// RemoveClient 从缓存中移除客户端
func RemoveClient(storageID uint) {
	cacheMutex.Lock()
	delete(clientCache, storageID)
	cacheMutex.Unlock()
}

// ClearClients 清除所有缓存的客户端
func ClearClients() {
	cacheMutex.Lock()
	clientCache = make(map[uint]Client)
	cacheMutex.Unlock()
}
