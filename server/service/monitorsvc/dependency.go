package monitorsvc

import (
	"context"
	"time"
)

func (s *MonitorService) CollectDependency(ctx context.Context) (*DependencyHealth, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	result := &DependencyHealth{}
	dbCh := make(chan DBHealth, 1)
	redisCh := make(chan RedisHealth, 1)
	ossCh := make(chan OssHealth, 1)

	go func() {
		dbCh <- s.CollectDBHealth(ctx)
	}()
	go func() {
		redisCh <- s.CollectRedisHealth(ctx)
	}()
	go func() {
		ossHealth, _ := s.CollectOss(ctx)
		if ossHealth != nil {
			ossCh <- *ossHealth
			return
		}
		ossCh <- OssHealth{CollectedAt: nowString(), Error: "oss health unavailable"}
	}()

	result.DB = receiveDBHealth(ctx, dbCh)
	result.Redis = receiveRedisHealth(ctx, redisCh)
	result.OSS = receiveOssHealth(ctx, ossCh)
	return result, nil
}

func receiveDBHealth(ctx context.Context, ch <-chan DBHealth) DBHealth {
	select {
	case value := <-ch:
		return value
	case <-ctx.Done():
		return DBHealth{Reachable: false, Error: safeError(ctx.Err())}
	}
}

func receiveRedisHealth(ctx context.Context, ch <-chan RedisHealth) RedisHealth {
	select {
	case value := <-ch:
		return value
	case <-ctx.Done():
		return RedisHealth{Reachable: false, Error: safeError(ctx.Err())}
	}
}

func receiveOssHealth(ctx context.Context, ch <-chan OssHealth) OssHealth {
	select {
	case value := <-ch:
		return value
	case <-ctx.Done():
		return OssHealth{CollectedAt: nowString(), Error: safeError(ctx.Err())}
	}
}
