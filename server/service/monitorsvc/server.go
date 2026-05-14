package monitorsvc

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func (s *MonitorService) CollectServerInfo(ctx context.Context) (*ServerInfo, error) {
	hostInfo, _ := host.InfoWithContext(ctx)
	cpuInfos, _ := cpu.InfoWithContext(ctx)
	physicalCore, _ := cpu.CountsWithContext(ctx, false)
	logicalCore, _ := cpu.CountsWithContext(ctx, true)
	percentList, _ := cpu.PercentWithContext(ctx, 0, false)
	virtualMem, _ := mem.VirtualMemoryWithContext(ctx)
	swapMem, _ := mem.SwapMemoryWithContext(ctx)
	loadAvg, _ := load.AvgWithContext(ctx)
	executable, _ := os.Executable()

	info := &ServerInfo{
		DataSource:  DataSourceHost,
		CollectedAt: nowString(),
		Host: HostInfo{
			Architecture: runtime.GOARCH,
		},
		CPU: CPUInfo{
			PhysicalCore: physicalCore,
			LogicalCore:  logicalCore,
		},
		Process: ProcessInfo{
			PID:           os.Getpid(),
			StartedAt:     timeString(processStartedAt),
			UptimeSeconds: int64(time.Since(processStartedAt).Seconds()),
			GoVersion:     runtime.Version(),
			NumCPU:        runtime.NumCPU(),
			GOMAXPROCS:    runtime.GOMAXPROCS(0),
			BinaryName:    binaryName(executable),
		},
	}

	if hostInfo != nil {
		info.Host.Hostname = hostInfo.Hostname
		info.Host.OS = hostInfo.OS
		info.Host.Platform = hostInfo.Platform
		info.Host.PlatformVersion = hostInfo.PlatformVersion
		info.Host.KernelVersion = hostInfo.KernelVersion
		info.Host.BootTime = timeString(time.Unix(int64(hostInfo.BootTime), 0))
		info.Host.UptimeSeconds = hostInfo.Uptime
	}
	if len(cpuInfos) > 0 {
		info.CPU.ModelName = cpuInfos[0].ModelName
	}
	if len(percentList) > 0 {
		info.CPU.UsagePercent = roundPercent(percentList[0])
	}
	if virtualMem != nil {
		info.Memory = MemoryInfo{
			Total:        virtualMem.Total,
			Used:         virtualMem.Used,
			Free:         virtualMem.Available,
			UsagePercent: roundPercent(virtualMem.UsedPercent),
		}
	}
	if swapMem != nil {
		info.Swap = SwapInfo{
			Total:        swapMem.Total,
			Used:         swapMem.Used,
			Free:         swapMem.Free,
			UsagePercent: roundPercent(swapMem.UsedPercent),
		}
	}
	if loadAvg != nil {
		info.Load = LoadInfo{Load1: loadAvg.Load1, Load5: loadAvg.Load5, Load15: loadAvg.Load15}
	}
	info.Disks = collectDisks(ctx)
	return info, nil
}

func collectDisks(ctx context.Context) []DiskPartition {
	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return []DiskPartition{}
	}
	items := make([]DiskPartition, 0, len(partitions))
	seen := make(map[string]struct{})
	for _, partition := range partitions {
		if partition.Mountpoint == "" {
			continue
		}
		if _, ok := seen[partition.Mountpoint]; ok {
			continue
		}
		seen[partition.Mountpoint] = struct{}{}
		usage, err := disk.UsageWithContext(ctx, partition.Mountpoint)
		if err != nil || usage == nil {
			continue
		}
		items = append(items, DiskPartition{
			Mountpoint:   partition.Mountpoint,
			FsType:       partition.Fstype,
			Total:        usage.Total,
			Used:         usage.Used,
			Free:         usage.Free,
			UsagePercent: roundPercent(usage.UsedPercent),
		})
	}
	return items
}
