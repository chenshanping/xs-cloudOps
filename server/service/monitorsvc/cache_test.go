package monitorsvc

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"server/global"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func TestClearByPrefixRejectsUnsafePrefix(t *testing.T) {
	svc := &MonitorService{}
	cases := []string{"", "*", "cache:userinfo:*", "user:"}
	for _, prefix := range cases {
		_, err := svc.ClearByPrefix(context.Background(), prefix)
		if !errors.Is(err, ErrCachePrefixNotAllowed) {
			t.Fatalf("prefix %q expected ErrCachePrefixNotAllowed, got %v", prefix, err)
		}
	}
}

func TestClearByPrefixDeletesAllowedPrefixOnly(t *testing.T) {
	restore := setupRedisForTest(t)
	defer restore()

	ctx := context.Background()
	if err := global.Redis.Set(ctx, "cache:userinfo:1", "a", 0).Err(); err != nil {
		t.Fatal(err)
	}
	if err := global.Redis.Set(ctx, "cache:userinfo:2", "b", 0).Err(); err != nil {
		t.Fatal(err)
	}
	if err := global.Redis.Set(ctx, "cache:dict:sys", "c", 0).Err(); err != nil {
		t.Fatal(err)
	}

	result, err := Default.ClearByPrefix(ctx, "cache:userinfo:")
	if err != nil {
		t.Fatal(err)
	}
	if result.Deleted != 2 || result.Truncated {
		t.Fatalf("unexpected result: %+v", result)
	}
	if n, err := global.Redis.Exists(ctx, "cache:userinfo:1", "cache:userinfo:2").Result(); err != nil || n != 0 {
		t.Fatalf("userinfo keys still exist: n=%d err=%v", n, err)
	}
	if n, err := global.Redis.Exists(ctx, "cache:dict:sys").Result(); err != nil || n != 1 {
		t.Fatalf("dict key should remain: n=%d err=%v", n, err)
	}
}

func TestCollectRedisCountsAllowedPrefixes(t *testing.T) {
	restore := setupRedisForTest(t)
	defer restore()

	ctx := context.Background()
	_ = global.Redis.Set(ctx, "cache:userperms:1", "a", 0).Err()
	_ = global.Redis.Set(ctx, "captcha:abc", "b", 0).Err()

	info, err := Default.CollectRedis(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Reachable {
		t.Fatalf("redis should be reachable: %+v", info)
	}
	counts := make(map[string]int64)
	for _, item := range info.PrefixCounts {
		counts[item.Prefix] = item.Count
	}
	if counts["cache:userperms:"] != 1 || counts["captcha:"] != 1 {
		t.Fatalf("unexpected counts: %+v", counts)
	}
}

func TestScanRedisPrefixWithLimitCanTruncate(t *testing.T) {
	restore := setupRedisForTest(t)
	defer restore()

	ctx := context.Background()
	for i := 0; i < 300; i++ {
		if err := global.Redis.Set(ctx, "cache:dict:item:"+strconv.Itoa(i), "v", 0).Err(); err != nil {
			t.Fatal(err)
		}
	}

	count, truncated, err := scanRedisPrefixWithLimit(ctx, "cache:dict:", false, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 || !truncated {
		t.Fatalf("expected truncated partial scan, count=%d truncated=%v", count, truncated)
	}
}

func setupRedisForTest(t *testing.T) func() {
	t.Helper()
	server, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	previous := global.Redis
	global.Redis = redis.NewClient(&redis.Options{Addr: server.Addr()})
	return func() {
		_ = global.Redis.Close()
		global.Redis = previous
		server.Close()
	}
}
