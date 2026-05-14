## ADDED Requirements

### Requirement: 系统主机指标查询

The system SHALL provide an authenticated endpoint that returns a snapshot of host-level runtime metrics (OS, CPU, memory, disk, load, host info, process uptime) for administrators with the `monitor:server:view` permission.

#### Scenario: Authorized administrator queries host metrics

- **WHEN** an administrator with `monitor:server:view` calls `GET /api/v1/monitor/server`
- **THEN** the system returns a 200 response containing OS family/version, hostname, CPU model and per-core usage, memory total/used/free, disk usage per mounted partition, system load averages, Go process PID and start time, and a `data_source` field indicating `host` or `container`

#### Scenario: Unauthorized user is rejected

- **WHEN** a user without `monitor:server:view` calls `GET /api/v1/monitor/server`
- **THEN** the system returns a 403 response and does not invoke any host-metric collection code

### Requirement: Go 运行时指标查询

The system SHALL expose Go runtime metrics including goroutine count, memory statistics, GC pause and count, number of CPUs, and Go version, gated by the `monitor:runtime:view` permission.

#### Scenario: Runtime snapshot is returned

- **WHEN** an administrator with `monitor:runtime:view` calls `GET /api/v1/monitor/runtime`
- **THEN** the system returns a 200 response with `goroutines`, `heap_alloc`, `heap_inuse`, `heap_sys`, `next_gc`, `last_gc`, `num_gc`, `pause_total_ns`, `num_cpu`, `go_version`, and `process_uptime_seconds`

#### Scenario: Repeated calls do not block

- **WHEN** an administrator calls `GET /api/v1/monitor/runtime` 10 times in succession
- **THEN** every call completes within 200 ms and the system does not start any background goroutine that persists between calls

### Requirement: 数据库连接池指标查询

The system SHALL expose the GORM-managed database connection pool statistics, gated by `monitor:db:view`.

#### Scenario: Pool stats are reported

- **WHEN** an administrator with `monitor:db:view` calls `GET /api/v1/monitor/db`
- **THEN** the system returns `max_open_connections`, `open_connections`, `in_use`, `idle`, `wait_count`, `wait_duration_ms`, `max_idle_closed`, `max_lifetime_closed`

#### Scenario: Database is unreachable

- **WHEN** the database is down at the moment the endpoint is called
- **THEN** the system returns a 200 response with `reachable: false` and a non-fatal error message, instead of returning 500

### Requirement: Redis 缓存指标查询

The system SHALL expose Redis runtime information and per-prefix key counts, gated by `monitor:cache:view`. Prefix scanning MUST use `SCAN` with a per-call iteration cap and MUST NOT use `KEYS`.

#### Scenario: Redis info is returned

- **WHEN** an administrator with `monitor:cache:view` calls `GET /api/v1/monitor/redis`
- **THEN** the system returns selected Redis `INFO` fields (connected_clients, used_memory, used_memory_peak, total_commands_processed, keyspace_hits, keyspace_misses, uptime_in_seconds), per-prefix key counts for the cache prefix whitelist, and a `truncated` flag per prefix indicating whether scan reached the iteration cap

#### Scenario: Scan iteration cap is enforced

- **WHEN** scanning a prefix would exceed the iteration cap
- **THEN** the system stops scanning, returns the partial count, and sets `truncated: true` for that prefix

### Requirement: 受控的 Redis 缓存清理

The system SHALL allow privileged administrators with `monitor:cache:clear` to delete Redis keys matching a prefix, BUT ONLY if the prefix is in the code-level whitelist. The whitelist MUST NOT be configurable from the database, system config UI, or any frontend input.

#### Scenario: Whitelisted prefix is cleared

- **WHEN** an administrator with `monitor:cache:clear` calls `POST /api/v1/monitor/redis/clear` with body `{"prefix": "cache:userinfo:"}`
- **THEN** the system scans and deletes all keys matching `cache:userinfo:*`, returns `{deleted: <n>, truncated: <bool>}`, and writes an operation log entry containing actor, prefix, and deleted count

#### Scenario: Non-whitelisted prefix is rejected

- **WHEN** an administrator calls `POST /api/v1/monitor/redis/clear` with `{"prefix": "user:"}` (not in whitelist)
- **THEN** the system returns a 400 response with code `cache_prefix_not_allowed` and does not touch Redis

#### Scenario: Empty prefix is rejected

- **WHEN** an administrator calls `POST /api/v1/monitor/redis/clear` with `{"prefix": ""}` or `{"prefix": "*"}`
- **THEN** the system returns a 400 response and does not scan or delete any key

#### Scenario: User without permission is rejected

- **WHEN** a user without `monitor:cache:clear` calls `POST /api/v1/monitor/redis/clear` with any body
- **THEN** the system returns 403 and does not access Redis

### Requirement: OSS 健康检测

The system SHALL report whether the active OSS provider is configured and reachable, gated by `monitor:oss:view`. When no OSS is enabled the endpoint SHALL still return 200 with `enabled: false`.

#### Scenario: OSS is enabled and healthy

- **WHEN** an administrator calls `GET /api/v1/monitor/oss` and an OSS provider is active
- **THEN** the system returns `{enabled: true, provider: <name>, reachable: true, latency_ms: <n>}` after a low-cost probe

#### Scenario: OSS is not configured

- **WHEN** no OSS provider is configured
- **THEN** the system returns `{enabled: false}` without attempting any network call

#### Scenario: OSS probe fails

- **WHEN** the OSS provider is configured but the probe times out or errors
- **THEN** the system returns `{enabled: true, provider: <name>, reachable: false, error: <safe message>}` and does not leak credentials or signed URLs in the response

### Requirement: 依赖健康概览

The system SHALL provide a single aggregated dependency health endpoint, gated by `monitor:dependency:view`, that summarizes DB / Redis / OSS reachability for use as a top-of-page badge.

#### Scenario: All dependencies healthy

- **WHEN** an administrator calls `GET /api/v1/monitor/dependency`
- **THEN** the system returns `{db: {reachable: true}, redis: {reachable: true}, oss: {enabled: <bool>, reachable: <bool>}}` within 1 second

#### Scenario: One dependency unhealthy

- **WHEN** Redis is unreachable but DB is healthy
- **THEN** the endpoint returns 200 with `redis.reachable: false` and continues to report DB health truthfully, without short-circuiting

### Requirement: 菜单与权限自动注册

On server startup the system SHALL idempotently ensure the existence of the 运维监控 root menu, 服务监控 menu, 操作日志 menu, 登录日志 menu, its 7 button permissions, the 7 SysApi rows with `need_auth=true`, and the corresponding sys_menu_api bindings. Existing custom metadata on these rows MUST be preserved except that 操作日志 and 登录日志 SHALL be parented under 运维监控.

#### Scenario: First-time bootstrap

- **WHEN** the server starts on a fresh database
- **THEN** the bootstrap function creates the 服务监控 menu, the 7 buttons, the 7 SysApi rows, and the menu↔api bindings, and grants all 7 permission codes to the built-in `admin` role

#### Scenario: Subsequent startup preserves admin edits

- **WHEN** the server restarts after an administrator manually edited a button label or SysApi description
- **THEN** the bootstrap function does not overwrite the edited fields and does not regrant permissions that were intentionally removed from non-`admin` roles

#### Scenario: Existing log menus are moved under operation monitor

- **WHEN** existing 操作日志 or 登录日志 menu rows are still parented under 系统管理 > 操作审计
- **THEN** startup moves those menu rows under 运维监控 without overwriting their name, path, component, icon, sort, status, or hidden fields
- **AND** if the legacy 操作审计 directory has no remaining child menus after the move, the system hides that empty directory
