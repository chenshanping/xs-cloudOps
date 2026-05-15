package utils

import (
	"context"
	"strings"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"

	"server/config"
	"server/global"
)

func setupJWTTestEnv(t *testing.T) {
	t.Helper()

	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	previousRedis := global.Redis
	previousConfig := global.Config
	global.Redis = client
	global.Config = &config.Config{
		JWT: config.JWT{
			Secret:        "jwt-test-secret",
			Expires:       1,
			RefreshWindow: 60,
			Issuer:        "jwt-test",
		},
	}

	t.Cleanup(func() {
		_ = client.Close()
		server.Close()
		global.Redis = previousRedis
		global.Config = previousConfig
	})
}

func TestGenerateTokenStoresWhitelistBeyondAccessExpiry(t *testing.T) {
	setupJWTTestEnv(t)

	token, err := GenerateToken(7, "alice", []uint{1}, []string{"admin"})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	ttl, err := global.Redis.TTL(context.Background(), TokenKey+"7").Result()
	if err != nil {
		t.Fatalf("query ttl: %v", err)
	}
	if ttl <= time.Second {
		t.Fatalf("expected whitelist ttl to cover refresh window, got %v", ttl)
	}
}

func TestRefreshTokenRejectsWhenNotCurrentWhitelistedToken(t *testing.T) {
	setupJWTTestEnv(t)

	token, err := GenerateToken(7, "alice", []uint{1}, []string{"admin"})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	if err := global.Redis.Set(context.Background(), TokenKey+"7", "another-token", tokenWhitelistTTL()).Err(); err != nil {
		t.Fatalf("overwrite whitelist: %v", err)
	}

	_, err = RefreshToken(token)
	if err == nil || !strings.Contains(err.Error(), "token已失效") {
		t.Fatalf("expected token invalid error, got %v", err)
	}
}
