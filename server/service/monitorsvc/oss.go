package monitorsvc

import (
	"context"
	"time"

	"server/model"
	"server/service/oss"
	"server/service/storagesvc"
)

func (s *MonitorService) CollectOss(ctx context.Context) (*OssHealth, error) {
	result := &OssHealth{CollectedAt: nowString()}
	profile, err := storagesvc.Default.GetSystemStorage()
	if err != nil {
		result.Error = safeError(err)
		return result, nil
	}
	if profile == nil || profile.Type == model.StorageTypeLocal {
		result.Enabled = false
		result.Provider = string(model.StorageTypeLocal)
		result.Reachable = true
		return result, nil
	}

	result.Enabled = true
	result.Provider = string(profile.Type)
	probeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	start := time.Now()
	client, err := oss.GetClient(profile)
	if err != nil {
		result.LatencyMs = millisSince(start)
		result.Error = safeError(err)
		return result, nil
	}
	_, err = client.Exists(probeCtx, "__go_base_monitor_probe__")
	result.LatencyMs = millisSince(start)
	if err != nil {
		result.Error = safeError(err)
		return result, nil
	}
	result.Reachable = true
	return result, nil
}
