package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"server/global"
	"server/model"
)

type CacheService struct{}

var Default = &CacheService{}

const (
	UserInfoCacheKey    = "cache:userinfo:"  // 用户信息缓存
	UserMenusCacheKey   = "cache:usermenus:" // 用户菜单缓存
	UserPermsCacheKey   = "cache:userperms:" // 用户权限缓存
	DictCacheKey        = "cache:dict:"      // 数据字典缓存
	CacheExpireTime     = 30 * time.Minute   // 缓存过期时间
	DictCacheExpireTime = 24 * time.Hour     // 字典缓存过期时间（24小时）
)

// UserInfoCache 用户信息缓存结构
type UserInfoCache struct {
	User        *model.SysUser  `json:"user"`
	Menus       []model.SysMenu `json:"menus"`
	Permissions []string        `json:"permissions"`
}

// GetUserInfoFromCache 从缓存获取用户信息
func (s *CacheService) GetUserInfoFromCache(userID uint) (*UserInfoCache, error) {
	if global.Redis == nil {
		return nil, fmt.Errorf("redis not initialized")
	}
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", UserInfoCacheKey, userID)

	data, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err // 缓存未命中
	}

	var cache UserInfoCache
	if err := json.Unmarshal([]byte(data), &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// SetUserInfoToCache 设置用户信息缓存
func (s *CacheService) SetUserInfoToCache(userID uint, cache *UserInfoCache) error {
	if global.Redis == nil {
		return nil
	}
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", UserInfoCacheKey, userID)

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return global.Redis.Set(ctx, key, data, CacheExpireTime).Err()
}

// ClearUserInfoCache 清除用户信息缓存
func (s *CacheService) ClearUserInfoCache(userID uint) error {
	if global.Redis == nil {
		return nil
	}
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", UserInfoCacheKey, userID)
	return global.Redis.Del(ctx, key).Err()
}

// ClearAllUserInfoCache 清除所有用户信息缓存（角色/菜单变更时使用）
func (s *CacheService) ClearAllUserInfoCache() error {
	if global.Redis == nil {
		return nil
	}
	ctx := context.Background()

	// 使用 SCAN 查找所有用户缓存 key
	var cursor uint64
	var keys []string
	for {
		var batch []string
		var err error
		batch, cursor, err = global.Redis.Scan(ctx, cursor, UserInfoCacheKey+"*", 100).Result()
		if err != nil {
			return err
		}
		keys = append(keys, batch...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		return global.Redis.Del(ctx, keys...).Err()
	}
	return nil
}

// ClearUserCacheByRoleID 清除指定角色的用户缓存
func (s *CacheService) ClearUserCacheByRoleID(roleID uint) error {
	if global.Redis == nil {
		return nil
	}
	// 查找拥有该角色的用户
	var userRoles []struct {
		SysUserID uint `gorm:"column:sys_user_id"`
	}
	if err := global.DB.Table("sys_user_role").Where("sys_role_id = ?", roleID).Find(&userRoles).Error; err != nil {
		return err
	}

	ctx := context.Background()
	for _, ur := range userRoles {
		key := fmt.Sprintf("%s%d", UserInfoCacheKey, ur.SysUserID)
		global.Redis.Del(ctx, key)
	}
	return nil
}

// ==================== 数据字典缓存 ====================

// GetDictFromCache 从缓存获取字典数据
func (s *CacheService) GetDictFromCache(dictType string) ([]byte, error) {
	if global.Redis == nil {
		return nil, fmt.Errorf("redis not initialized")
	}
	ctx := context.Background()
	key := DictCacheKey + dictType
	return global.Redis.Get(ctx, key).Bytes()
}

// SetDictToCache 设置字典数据缓存
func (s *CacheService) SetDictToCache(dictType string, data []byte) error {
	if global.Redis == nil {
		return nil
	}
	ctx := context.Background()
	key := DictCacheKey + dictType
	return global.Redis.Set(ctx, key, data, DictCacheExpireTime).Err()
}

// ClearDictCache 清除指定字典类型的缓存
func (s *CacheService) ClearDictCache(dictType string) error {
	if global.Redis == nil {
		return nil
	}
	ctx := context.Background()
	key := DictCacheKey + dictType
	return global.Redis.Del(ctx, key).Err()
}

// ClearAllDictCache 清除所有字典缓存
func (s *CacheService) ClearAllDictCache() error {
	if global.Redis == nil {
		return nil
	}
	ctx := context.Background()

	var cursor uint64
	var keys []string
	for {
		var batch []string
		var err error
		batch, cursor, err = global.Redis.Scan(ctx, cursor, DictCacheKey+"*", 100).Result()
		if err != nil {
			return err
		}
		keys = append(keys, batch...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		return global.Redis.Del(ctx, keys...).Err()
	}
	return nil
}
