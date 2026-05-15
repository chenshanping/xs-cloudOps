package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"server/global"
)

type PublicAuthGuard struct{}

var PublicGuard = &PublicAuthGuard{}

const (
	registerRateLimitMessage       = "注册过于频繁，请稍后再试"
	emailCodeRateLimitMessage      = "验证码发送过于频繁，请稍后再试"
	passwordResetRateLimitMessage  = "密码重置操作过于频繁，请稍后再试"
	refreshRateLimitMessage        = "Token刷新过于频繁，请稍后再试"
	rateLimitStorageErrorMessage   = "服务繁忙，请稍后再试"
)

func (g *PublicAuthGuard) CheckRegister(ip, username, email string) error {
	return g.checkAll(
		registerRateLimitMessage,
		limitRule{key: g.composeKey("register", "ip", ip), limit: 10, window: 15 * time.Minute},
		limitRule{key: g.composeKey("register", "username", username), limit: 5, window: 15 * time.Minute},
		limitRule{key: g.composeKey("register", "email", email), limit: 5, window: 15 * time.Minute},
	)
}

func (g *PublicAuthGuard) CheckSendEmailCode(ip, email string) error {
	return g.checkAll(
		emailCodeRateLimitMessage,
		limitRule{key: g.composeKey("send-email-code", "ip", ip), limit: 10, window: 15 * time.Minute},
		limitRule{key: g.composeKey("send-email-code", "email", email), limit: 3, window: 10 * time.Minute},
	)
}

func (g *PublicAuthGuard) CheckResetPasswordByToken(ip, token string) error {
	return g.checkAll(
		passwordResetRateLimitMessage,
		limitRule{key: g.composeKey("reset-by-token", "ip", ip), limit: 10, window: 15 * time.Minute},
		limitRule{key: g.composeKey("reset-by-token", "token", g.tokenFingerprint(token)), limit: 5, window: 15 * time.Minute},
	)
}

func (g *PublicAuthGuard) CheckResetPasswordByEmail(ip, email string) error {
	return g.checkAll(
		passwordResetRateLimitMessage,
		limitRule{key: g.composeKey("reset-by-email", "ip", ip), limit: 10, window: 15 * time.Minute},
		limitRule{key: g.composeKey("reset-by-email", "email", email), limit: 5, window: 15 * time.Minute},
	)
}

func (g *PublicAuthGuard) CheckResetPasswordByUsername(ip, username string) error {
	return g.checkAll(
		passwordResetRateLimitMessage,
		limitRule{key: g.composeKey("reset-by-username", "ip", ip), limit: 10, window: 15 * time.Minute},
		limitRule{key: g.composeKey("reset-by-username", "username", username), limit: 5, window: 15 * time.Minute},
	)
}

func (g *PublicAuthGuard) CheckRefreshToken(ip, token string) error {
	return g.checkAll(
		refreshRateLimitMessage,
		limitRule{key: g.composeKey("refresh", "ip", ip), limit: 30, window: 15 * time.Minute},
		limitRule{key: g.composeKey("refresh", "token", g.tokenFingerprint(token)), limit: 10, window: 15 * time.Minute},
	)
}

type limitRule struct {
	key    string
	limit  int64
	window time.Duration
}

func (g *PublicAuthGuard) checkAll(limitMessage string, rules ...limitRule) error {
	for _, rule := range rules {
		if rule.key == "" {
			continue
		}
		if err := g.consume(rule); err != nil {
			if errors.Is(err, errLimitStorageUnavailable) {
				return errors.New(rateLimitStorageErrorMessage)
			}
			return errors.New(limitMessage)
		}
	}
	return nil
}

var errLimitStorageUnavailable = errors.New("rate_limit_storage_unavailable")

func (g *PublicAuthGuard) consume(rule limitRule) error {
	if global.Redis == nil {
		return errLimitStorageUnavailable
	}

	ctx := context.Background()
	count, err := global.Redis.Incr(ctx, rule.key).Result()
	if err != nil {
		return errLimitStorageUnavailable
	}
	if count == 1 {
		if err := global.Redis.Expire(ctx, rule.key, rule.window).Err(); err != nil {
			return errLimitStorageUnavailable
		}
	}
	if count > rule.limit {
		return errors.New("rate_limited")
	}
	return nil
}

func (g *PublicAuthGuard) composeKey(action, dimension, raw string) string {
	value := strings.TrimSpace(strings.ToLower(raw))
	if value == "" {
		return ""
	}
	return fmt.Sprintf("public_auth_rl:%s:%s:%s", action, dimension, value)
}

func (g *PublicAuthGuard) tokenFingerprint(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:8])
}
