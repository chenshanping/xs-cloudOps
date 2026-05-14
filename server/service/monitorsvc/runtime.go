package monitorsvc

import (
	"context"
	"runtime"
	"time"
)

func (s *MonitorService) CollectRuntime(ctx context.Context) (*RuntimeInfo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	lastGC := ""
	if memStats.LastGC > 0 {
		lastGC = timeString(time.Unix(0, int64(memStats.LastGC)))
	}

	return &RuntimeInfo{
		Goroutines:           runtime.NumGoroutine(),
		NumCPU:               runtime.NumCPU(),
		GOMAXPROCS:           runtime.GOMAXPROCS(0),
		GoVersion:            runtime.Version(),
		HeapAlloc:            memStats.HeapAlloc,
		HeapInuse:            memStats.HeapInuse,
		HeapSys:              memStats.HeapSys,
		HeapObjects:          memStats.HeapObjects,
		StackInuse:           memStats.StackInuse,
		StackSys:             memStats.StackSys,
		NextGC:               memStats.NextGC,
		LastGC:               lastGC,
		NumGC:                memStats.NumGC,
		NumForcedGC:          memStats.NumForcedGC,
		PauseTotalNs:         memStats.PauseTotalNs,
		GCCPUFraction:        memStats.GCCPUFraction,
		ProcessUptimeSeconds: int64(time.Since(processStartedAt).Seconds()),
		CollectedAt:          nowString(),
	}, nil
}
