package auth

import (
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"

	"server/global"
)

func setupPublicAuthGuardRedis(t *testing.T) {
	t.Helper()

	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	previousRedis := global.Redis
	global.Redis = client
	t.Cleanup(func() {
		_ = client.Close()
		server.Close()
		global.Redis = previousRedis
	})
}

func TestPublicAuthGuardCheckSendEmailCodeRateLimitedByEmail(t *testing.T) {
	setupPublicAuthGuardRedis(t)

	for i := 0; i < 3; i++ {
		if err := PublicGuard.CheckSendEmailCode("127.0.0.1", "rate@example.com"); err != nil {
			t.Fatalf("unexpected rate limit before threshold: %v", err)
		}
	}

	err := PublicGuard.CheckSendEmailCode("127.0.0.1", "rate@example.com")
	if err == nil || err.Error() != emailCodeRateLimitMessage {
		t.Fatalf("expected email code rate limit message, got %v", err)
	}
}

func TestPublicAuthGuardCheckRefreshTokenRateLimitedByToken(t *testing.T) {
	setupPublicAuthGuardRedis(t)

	for i := 0; i < 10; i++ {
		if err := PublicGuard.CheckRefreshToken("127.0.0.1", "token-A"); err != nil {
			t.Fatalf("unexpected refresh rate limit before threshold: %v", err)
		}
	}

	err := PublicGuard.CheckRefreshToken("127.0.0.1", "token-A")
	if err == nil || err.Error() != refreshRateLimitMessage {
		t.Fatalf("expected refresh rate limit message, got %v", err)
	}
}
