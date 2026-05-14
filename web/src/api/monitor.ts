import request from '@/utils/request'
import type { ApiResponse } from '@/types'

export type DataSource = 'host' | 'container'

export interface HostInfo {
  hostname: string
  os: string
  platform: string
  platform_version: string
  kernel_version: string
  architecture: string
  boot_time: string
  uptime_seconds: number
}

export interface CPUInfo {
  model_name: string
  physical_core: number
  logical_core: number
  usage_percent: number
}

export interface MemoryInfo {
  total: number
  used: number
  free: number
  usage_percent: number
}

export interface SwapInfo {
  total: number
  used: number
  free: number
  usage_percent: number
}

export interface LoadInfo {
  load_1: number
  load_5: number
  load_15: number
}

export interface DiskPartition {
  mountpoint: string
  fs_type: string
  total: number
  used: number
  free: number
  usage_percent: number
}

export interface ProcessInfo {
  pid: number
  started_at: string
  uptime_seconds: number
  go_version: string
  num_cpu: number
  go_max_procs: number
  binary_name: string
}

export interface ServerInfo {
  data_source: DataSource
  host: HostInfo
  cpu: CPUInfo
  memory: MemoryInfo
  swap: SwapInfo
  load: LoadInfo
  disks: DiskPartition[]
  process: ProcessInfo
  collected_at: string
}

export interface RuntimeInfo {
  goroutines: number
  num_cpu: number
  go_max_procs: number
  go_version: string
  heap_alloc: number
  heap_inuse: number
  heap_sys: number
  heap_objects: number
  stack_inuse: number
  stack_sys: number
  next_gc: number
  last_gc: string
  num_gc: number
  num_forced_gc: number
  pause_total_ns: number
  gc_cpu_fraction: number
  process_uptime_seconds: number
  collected_at: string
}

export interface DBStats {
  reachable: boolean
  ping_latency_ms: number
  error?: string
  max_open_connections: number
  open_connections: number
  in_use: number
  idle: number
  wait_count: number
  wait_duration_ms: number
  max_idle_closed: number
  max_idle_time_closed: number
  max_lifetime_closed: number
  collected_at: string
}

export interface RedisPrefixCount {
  prefix: string
  count: number
  truncated: boolean
  error?: string
}

export interface RedisInfo {
  reachable: boolean
  ping_latency_ms: number
  error?: string
  version: string
  uptime_seconds: number
  connected_clients: number
  used_memory: number
  used_memory_peak: number
  used_memory_human: string
  total_commands_processed: number
  keyspace_hits: number
  keyspace_misses: number
  hit_rate: number
  db_size: number
  allowed_prefixes: string[]
  prefix_counts: RedisPrefixCount[]
  collected_at: string
}

export interface ClearCacheResult {
  prefix: string
  deleted: number
  truncated: boolean
}

export interface OssHealth {
  enabled: boolean
  provider?: string
  reachable: boolean
  latency_ms: number
  error?: string
  collected_at: string
}

export interface DBHealth {
  reachable: boolean
  ping_latency_ms: number
  error?: string
}

export interface RedisHealth {
  reachable: boolean
  ping_latency_ms: number
  error?: string
}

export interface DependencyHealth {
  db: DBHealth
  redis: RedisHealth
  oss: OssHealth
}

export function getServerInfo() {
  return request.get<any, ApiResponse<ServerInfo>>('/monitor/server', { silent: true })
}

export function getRuntimeInfo() {
  return request.get<any, ApiResponse<RuntimeInfo>>('/monitor/runtime', { silent: true })
}

export function getDbStats() {
  return request.get<any, ApiResponse<DBStats>>('/monitor/db', { silent: true })
}

export function getRedisInfo() {
  return request.get<any, ApiResponse<RedisInfo>>('/monitor/redis', { silent: true })
}

export function clearCacheByPrefix(prefix: string) {
  return request.post<any, ApiResponse<ClearCacheResult>>('/monitor/redis/clear', { prefix })
}

export function getOssHealth() {
  return request.get<any, ApiResponse<OssHealth>>('/monitor/oss', { silent: true })
}

export function getDependencyHealth() {
  return request.get<any, ApiResponse<DependencyHealth>>('/monitor/dependency', { silent: true })
}
