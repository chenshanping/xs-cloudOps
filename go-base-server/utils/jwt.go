package utils

import (
	"context"
	"go-base-server/global"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Redis key 前缀
	TokenKey     = "jwt:token:"     // 白名单：用户当前有效 Token
	BlacklistKey = "jwt:blacklist:" // 黑名单：已失效的 Token
)

type Claims struct {
	UserID    uint     `json:"user_id"`
	Username  string   `json:"username"`
	RoleIDs   []uint   `json:"role_ids"`
	RoleCodes []string `json:"role_codes"`
	jwt.RegisteredClaims
}

// 生成Token并存入Redis
func GenerateToken(userID uint, username string, roleIDs []uint, roleCodes []string) (string, error) {
	cfg := global.Config.JWT
	claims := Claims{
		UserID:    userID,
		Username:  username,
		RoleIDs:   roleIDs,
		RoleCodes: roleCodes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.Expires) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}

	// 存入 Redis（白名单模式）
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", TokenKey, userID)
	err = global.Redis.Set(ctx, key, tokenString, time.Duration(cfg.Expires)*time.Second).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 解析Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 解析Token（忽略过期时间，用于刷新场景）
func ParseTokenIgnoreExpired(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.Secret), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 验证Token是否在Redis中有效（未被拉黑）
func ValidateTokenInRedis(tokenString string, userID uint) error {
	ctx := context.Background()

	// 检查是否在黑名单
	blackKey := BlacklistKey + tokenString
	exists, err := global.Redis.Exists(ctx, blackKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("token已失效")
	}

	// 检查白名单
	key := fmt.Sprintf("%s%d", TokenKey, userID)
	storedToken, err := global.Redis.Get(ctx, key).Result()
	if err != nil || storedToken != tokenString {
		return errors.New("token已失效或已在其他设备登录")
	}

	return nil
}

// 将Token加入黑名单（登出时调用）
func InvalidateToken(tokenString string) error {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return err
	}

	ctx := context.Background()
	// 计算剩余过期时间
	expTime := claims.ExpiresAt.Time
	ttl := time.Until(expTime)
	if ttl <= 0 {
		return nil // 已过期无需处理
	}

	// 加入黑名单
	key := BlacklistKey + tokenString
	return global.Redis.Set(ctx, key, "1", ttl).Err()
}

// 删除用户的Token（强制下线）
// 将当前 Token 加入黑名单，并删除白名单
func RemoveUserToken(userID uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", TokenKey, userID)

	// 获取当前 Token 并加入黑名单
	if tokenString, err := global.Redis.Get(ctx, key).Result(); err == nil && tokenString != "" {
		// 将 Token 加入黑名单，设置较长的过期时间（覆盖刷新窗口）
		refreshWindow := global.Config.JWT.RefreshWindow
		if refreshWindow <= 0 {
			refreshWindow = 7 * 24 * 3600
		}
		blackKey := BlacklistKey + tokenString
		_ = global.Redis.Set(ctx, blackKey, "1", time.Duration(refreshWindow)*time.Second).Err()
	}

	// 删除白名单
	return global.Redis.Del(ctx, key).Err()
}

// 刷新Token（允许在Token过期后一定时间内刷新）
// 注意：刷新时会重新查询数据库获取最新角色信息
func RefreshToken(tokenString string) (string, error) {
	// 使用忽略过期的解析，允许刷新过期的Token
	claims, err := ParseTokenIgnoreExpired(tokenString)
	if err != nil {
		return "", err
	}

	// 限制刷新窗口：Token过期超过配置时间则不允许刷新
	refreshWindow := global.Config.JWT.RefreshWindow
	if refreshWindow <= 0 {
		refreshWindow = 7 * 24 * 3600 // 默认7天
	}
	if claims.ExpiresAt != nil {
		expiredDuration := time.Since(claims.ExpiresAt.Time)
		if expiredDuration > time.Duration(refreshWindow)*time.Second {
			return "", errors.New("token已过期太久，请重新登录")
		}
	}

	ctx := context.Background()

	// 检查 Token 是否在黑名单中（被强制下线）
	blackKey := BlacklistKey + tokenString
	exists, _ := global.Redis.Exists(ctx, blackKey).Result()
	if exists > 0 {

		// 检查用户 Token 是否在白名单中（如果被删除则不允许刷新）
		ctx := context.Background()
		key := fmt.Sprintf("%s%d", TokenKey, claims.UserID)
		storedToken, err := global.Redis.Get(ctx, key).Result()
		if err != nil || storedToken != tokenString {
			return "", errors.New("token已失效，请重新登录")
		}
	}

	// 旧Token加入黑名单（如果还未过期）
	_ = InvalidateToken(tokenString)

	// 注意：这里仍然使用旧的角色信息，因为角色变化时已经删除了 Token
	// 如果能走到这里，说明角色没有变化

	return GenerateToken(claims.UserID, claims.Username, claims.RoleIDs, claims.RoleCodes)
}
