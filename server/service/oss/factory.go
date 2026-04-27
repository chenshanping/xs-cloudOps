package oss

import (
	"fmt"
	"sync"

	"server/model"
)

var (
	clientCache         = make(map[string]Client)
	customClientBuilder = make(map[model.StorageType]func(*model.StorageProfile) (Client, error))
	cacheMutex          sync.RWMutex
)

// NewClient 根据存储配置创建客户端
func NewClient(storage *model.StorageProfile) (Client, error) {
	cacheMutex.RLock()
	builder := customClientBuilder[model.StorageType(storage.Type)]
	cacheMutex.RUnlock()
	if builder != nil {
		return builder(storage)
	}

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
func GetClient(storage *model.StorageProfile) (Client, error) {
	cacheKey := storage.CacheKey()

	cacheMutex.RLock()
	if client, ok := clientCache[cacheKey]; ok {
		cacheMutex.RUnlock()
		return client, nil
	}
	cacheMutex.RUnlock()

	client, err := NewClient(storage)
	if err != nil {
		return nil, err
	}

	cacheMutex.Lock()
	clientCache[cacheKey] = client
	cacheMutex.Unlock()

	return client, nil
}

// ClearClients 清除所有缓存的客户端
func ClearClients() {
	cacheMutex.Lock()
	clientCache = make(map[string]Client)
	cacheMutex.Unlock()
}

// RegisterClientBuilderForTest 注册测试客户端构造器
func RegisterClientBuilderForTest(storageType model.StorageType, builder func(*model.StorageProfile) (Client, error)) func() {
	cacheMutex.Lock()
	customClientBuilder[storageType] = builder
	clientCache = make(map[string]Client)
	cacheMutex.Unlock()

	return func() {
		cacheMutex.Lock()
		delete(customClientBuilder, storageType)
		clientCache = make(map[string]Client)
		cacheMutex.Unlock()
	}
}
