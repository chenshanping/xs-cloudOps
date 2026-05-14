package monitorsvc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"server/global"
)

func (s *MonitorService) CollectDB(ctx context.Context) (*DBStats, error) {
	result := &DBStats{CollectedAt: nowString()}
	if global.DB == nil {
		result.Error = "database not initialized"
		return result, nil
	}

	sqlDB, err := global.DB.DB()
	if err != nil {
		result.Error = safeError(err)
		return result, nil
	}

	stats := sqlDB.Stats()
	fillDBStats(result, stats)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	start := time.Now()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		result.Error = safeError(err)
		result.PingLatencyMs = millisSince(start)
		return result, nil
	}
	result.Reachable = true
	result.PingLatencyMs = millisSince(start)
	return result, nil
}

func (s *MonitorService) CollectDBHealth(ctx context.Context) DBHealth {
	if global.DB == nil {
		return DBHealth{Reachable: false, Error: "database not initialized"}
	}
	sqlDB, err := global.DB.DB()
	if err != nil {
		return DBHealth{Reachable: false, Error: safeError(err)}
	}
	pingCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	start := time.Now()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return DBHealth{Reachable: false, PingLatencyMs: millisSince(start), Error: safeError(fmt.Errorf("ping database: %w", err))}
	}
	return DBHealth{Reachable: true, PingLatencyMs: millisSince(start)}
}

func fillDBStats(result *DBStats, stats sql.DBStats) {
	result.MaxOpenConnections = stats.MaxOpenConnections
	result.OpenConnections = stats.OpenConnections
	result.InUse = stats.InUse
	result.Idle = stats.Idle
	result.WaitCount = stats.WaitCount
	result.WaitDurationMs = stats.WaitDuration.Milliseconds()
	result.MaxIdleClosed = stats.MaxIdleClosed
	result.MaxIdleTimeClosed = stats.MaxIdleTimeClosed
	result.MaxLifetimeClosed = stats.MaxLifetimeClosed
}
