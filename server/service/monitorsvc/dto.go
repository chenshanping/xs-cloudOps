package monitorsvc

type DataSource string

const (
	DataSourceHost      DataSource = "host"
	DataSourceContainer DataSource = "container"
)

type HostInfo struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	KernelVersion   string `json:"kernel_version"`
	Architecture    string `json:"architecture"`
	BootTime        string `json:"boot_time"`
	UptimeSeconds   uint64 `json:"uptime_seconds"`
}

type CPUInfo struct {
	ModelName    string  `json:"model_name"`
	PhysicalCore int     `json:"physical_core"`
	LogicalCore  int     `json:"logical_core"`
	UsagePercent float64 `json:"usage_percent"`
}

type MemoryInfo struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePercent float64 `json:"usage_percent"`
}

type SwapInfo struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePercent float64 `json:"usage_percent"`
}

type LoadInfo struct {
	Load1  float64 `json:"load_1"`
	Load5  float64 `json:"load_5"`
	Load15 float64 `json:"load_15"`
}

type DiskPartition struct {
	Mountpoint   string  `json:"mountpoint"`
	FsType       string  `json:"fs_type"`
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePercent float64 `json:"usage_percent"`
}

type ProcessInfo struct {
	PID           int    `json:"pid"`
	StartedAt     string `json:"started_at"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	GoVersion     string `json:"go_version"`
	NumCPU        int    `json:"num_cpu"`
	GOMAXPROCS    int    `json:"go_max_procs"`
	BinaryName    string `json:"binary_name"`
}

type ServerInfo struct {
	DataSource  DataSource      `json:"data_source"`
	Host        HostInfo        `json:"host"`
	CPU         CPUInfo         `json:"cpu"`
	Memory      MemoryInfo      `json:"memory"`
	Swap        SwapInfo        `json:"swap"`
	Load        LoadInfo        `json:"load"`
	Disks       []DiskPartition `json:"disks"`
	Process     ProcessInfo     `json:"process"`
	CollectedAt string          `json:"collected_at"`
}

type RuntimeInfo struct {
	Goroutines           int     `json:"goroutines"`
	NumCPU               int     `json:"num_cpu"`
	GOMAXPROCS           int     `json:"go_max_procs"`
	GoVersion            string  `json:"go_version"`
	HeapAlloc            uint64  `json:"heap_alloc"`
	HeapInuse            uint64  `json:"heap_inuse"`
	HeapSys              uint64  `json:"heap_sys"`
	HeapObjects          uint64  `json:"heap_objects"`
	StackInuse           uint64  `json:"stack_inuse"`
	StackSys             uint64  `json:"stack_sys"`
	NextGC               uint64  `json:"next_gc"`
	LastGC               string  `json:"last_gc"`
	NumGC                uint32  `json:"num_gc"`
	NumForcedGC          uint32  `json:"num_forced_gc"`
	PauseTotalNs         uint64  `json:"pause_total_ns"`
	GCCPUFraction        float64 `json:"gc_cpu_fraction"`
	ProcessUptimeSeconds int64   `json:"process_uptime_seconds"`
	CollectedAt          string  `json:"collected_at"`
}

type DBStats struct {
	Reachable          bool   `json:"reachable"`
	PingLatencyMs      int64  `json:"ping_latency_ms"`
	Error              string `json:"error,omitempty"`
	MaxOpenConnections int    `json:"max_open_connections"`
	OpenConnections    int    `json:"open_connections"`
	InUse              int    `json:"in_use"`
	Idle               int    `json:"idle"`
	WaitCount          int64  `json:"wait_count"`
	WaitDurationMs     int64  `json:"wait_duration_ms"`
	MaxIdleClosed      int64  `json:"max_idle_closed"`
	MaxIdleTimeClosed  int64  `json:"max_idle_time_closed"`
	MaxLifetimeClosed  int64  `json:"max_lifetime_closed"`
	CollectedAt        string `json:"collected_at"`
}

type RedisPrefixCount struct {
	Prefix    string `json:"prefix"`
	Count     int64  `json:"count"`
	Truncated bool   `json:"truncated"`
	Error     string `json:"error,omitempty"`
}

type RedisInfo struct {
	Reachable              bool               `json:"reachable"`
	PingLatencyMs          int64              `json:"ping_latency_ms"`
	Error                  string             `json:"error,omitempty"`
	Version                string             `json:"version"`
	UptimeSeconds          int64              `json:"uptime_seconds"`
	ConnectedClients       int64              `json:"connected_clients"`
	UsedMemory             int64              `json:"used_memory"`
	UsedMemoryPeak         int64              `json:"used_memory_peak"`
	UsedMemoryHuman        string             `json:"used_memory_human"`
	TotalCommandsProcessed int64              `json:"total_commands_processed"`
	KeyspaceHits           int64              `json:"keyspace_hits"`
	KeyspaceMisses         int64              `json:"keyspace_misses"`
	HitRate                float64            `json:"hit_rate"`
	DBSize                 int64              `json:"db_size"`
	AllowedPrefixes        []string           `json:"allowed_prefixes"`
	PrefixCounts           []RedisPrefixCount `json:"prefix_counts"`
	CollectedAt            string             `json:"collected_at"`
}

type ClearCacheResult struct {
	Prefix    string `json:"prefix"`
	Deleted   int64  `json:"deleted"`
	Truncated bool   `json:"truncated"`
}

type OssHealth struct {
	Enabled     bool   `json:"enabled"`
	Provider    string `json:"provider,omitempty"`
	Reachable   bool   `json:"reachable"`
	LatencyMs   int64  `json:"latency_ms"`
	Error       string `json:"error,omitempty"`
	CollectedAt string `json:"collected_at"`
}

type DependencyHealth struct {
	DB    DBHealth    `json:"db"`
	Redis RedisHealth `json:"redis"`
	OSS   OssHealth   `json:"oss"`
}

type DBHealth struct {
	Reachable     bool   `json:"reachable"`
	PingLatencyMs int64  `json:"ping_latency_ms"`
	Error         string `json:"error,omitempty"`
}

type RedisHealth struct {
	Reachable     bool   `json:"reachable"`
	PingLatencyMs int64  `json:"ping_latency_ms"`
	Error         string `json:"error,omitempty"`
}
