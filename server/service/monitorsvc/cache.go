package monitorsvc

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"server/global"
)

var ErrCachePrefixNotAllowed = errors.New("cache_prefix_not_allowed")

const (
	redisScanBatch   int64 = 200
	redisScanMaxIter       = 5000
)

var allowedCachePrefixes = []string{
	"cache:userinfo:",
	"cache:usermenus:",
	"cache:userperms:",
	"cache:dict:",
	"captcha:",
}

func (s *MonitorService) AllowedCachePrefixes() []string {
	result := make([]string, len(allowedCachePrefixes))
	copy(result, allowedCachePrefixes)
	return result
}

func (s *MonitorService) CollectRedis(ctx context.Context) (*RedisInfo, error) {
	result := &RedisInfo{
		AllowedPrefixes: s.AllowedCachePrefixes(),
		CollectedAt:     nowString(),
	}
	if global.Redis == nil {
		result.Error = "redis not initialized"
		return result, nil
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	start := time.Now()
	if err := global.Redis.Ping(pingCtx).Err(); err != nil {
		result.Error = safeError(err)
		result.PingLatencyMs = millisSince(start)
		return result, nil
	}
	result.Reachable = true
	result.PingLatencyMs = millisSince(start)

	if info, err := global.Redis.Info(ctx).Result(); err == nil {
		fillRedisInfo(result, parseRedisInfo(info))
	} else {
		result.Error = safeError(err)
	}
	if size, err := global.Redis.DBSize(ctx).Result(); err == nil {
		result.DBSize = size
	}

	result.PrefixCounts = make([]RedisPrefixCount, 0, len(allowedCachePrefixes))
	for _, prefix := range allowedCachePrefixes {
		count, truncated, err := scanRedisPrefix(ctx, prefix, false)
		item := RedisPrefixCount{Prefix: prefix, Count: count, Truncated: truncated}
		if err != nil {
			item.Error = safeError(err)
		}
		result.PrefixCounts = append(result.PrefixCounts, item)
	}
	return result, nil
}

func (s *MonitorService) CollectRedisHealth(ctx context.Context) RedisHealth {
	if global.Redis == nil {
		return RedisHealth{Reachable: false, Error: "redis not initialized"}
	}
	pingCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	start := time.Now()
	if err := global.Redis.Ping(pingCtx).Err(); err != nil {
		return RedisHealth{Reachable: false, PingLatencyMs: millisSince(start), Error: safeError(err)}
	}
	return RedisHealth{Reachable: true, PingLatencyMs: millisSince(start)}
}

func (s *MonitorService) ClearByPrefix(ctx context.Context, prefix string) (*ClearCacheResult, error) {
	prefix = strings.TrimSpace(prefix)
	if !isAllowedCachePrefix(prefix) {
		return nil, ErrCachePrefixNotAllowed
	}
	deleted, truncated, err := scanRedisPrefix(ctx, prefix, true)
	if err != nil {
		return nil, err
	}
	return &ClearCacheResult{Prefix: prefix, Deleted: deleted, Truncated: truncated}, nil
}

func isAllowedCachePrefix(prefix string) bool {
	if prefix == "" || strings.ContainsAny(prefix, "*?[]") {
		return false
	}
	for _, allowed := range allowedCachePrefixes {
		if prefix == allowed {
			return true
		}
	}
	return false
}

func scanRedisPrefix(ctx context.Context, prefix string, deleteKeys bool) (int64, bool, error) {
	return scanRedisPrefixWithLimit(ctx, prefix, deleteKeys, redisScanBatch, redisScanMaxIter)
}

func scanRedisPrefixWithLimit(ctx context.Context, prefix string, deleteKeys bool, batch int64, maxIter int) (int64, bool, error) {
	if global.Redis == nil {
		return 0, false, errors.New("redis not initialized")
	}
	if batch <= 0 {
		batch = redisScanBatch
	}
	if maxIter <= 0 {
		return 0, true, nil
	}
	var cursor uint64
	var total int64
	truncated := false
	for iter := 0; ; iter++ {
		if iter >= maxIter {
			truncated = cursor != 0
			break
		}
		keys, nextCursor, err := global.Redis.Scan(ctx, cursor, prefix+"*", batch).Result()
		if err != nil {
			return total, truncated, err
		}
		if deleteKeys && len(keys) > 0 {
			deleted, err := global.Redis.Del(ctx, keys...).Result()
			if err != nil {
				return total, truncated, err
			}
			total += deleted
		} else {
			total += int64(len(keys))
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return total, truncated, nil
}

func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(info, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		result[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return result
}

func fillRedisInfo(result *RedisInfo, fields map[string]string) {
	result.Version = fields["redis_version"]
	result.UsedMemoryHuman = fields["used_memory_human"]
	result.UptimeSeconds = redisInt(fields, "uptime_in_seconds")
	result.ConnectedClients = redisInt(fields, "connected_clients")
	result.UsedMemory = redisInt(fields, "used_memory")
	result.UsedMemoryPeak = redisInt(fields, "used_memory_peak")
	result.TotalCommandsProcessed = redisInt(fields, "total_commands_processed")
	result.KeyspaceHits = redisInt(fields, "keyspace_hits")
	result.KeyspaceMisses = redisInt(fields, "keyspace_misses")
	total := result.KeyspaceHits + result.KeyspaceMisses
	if total > 0 {
		result.HitRate = roundPercent(float64(result.KeyspaceHits) / float64(total) * 100)
	}
}

func redisInt(fields map[string]string, key string) int64 {
	value, _ := strconv.ParseInt(fields[key], 10, 64)
	return value
}
