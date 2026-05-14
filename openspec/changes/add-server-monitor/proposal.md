## Why

go-base 已经有「登录日志 / 操作日志」两个运维入口，但缺少**实时**的服务运行健康可见性。当线上出现卡顿、内存吃紧、Redis 连接异常、DB 连接池打满、缓存堆积时，管理员只能登服务器或通过外部监控查看，**无法在系统内自助诊断**，也无法在受控边界内做"按前缀清理缓存"这类轻运维动作。

参考 SaiAdmin 的「服务监控」模块对应能力，结合 Go 原生 `runtime` 和 GORM `sql.DBStats` 的优势，新增一个**只读为主、写操作严格受限**的服务监控模块。

## What Changes

- **新增/调整** `运维监控` 一级菜单分组，下挂 `服务监控`、`操作日志`、`登录日志`；存量 `系统管理 > 操作审计` 下的日志菜单迁移到 `运维监控`
- **新增** 后端 `service/monitorsvc/` 服务层，聚合主机、Go 运行时、DB 连接池、Redis、OSS 的运行指标
- **新增** 受控的「按前缀清理 Redis 缓存」功能；**前缀白名单写在代码里**（初版按当前代码实际前缀：`cache:userinfo:`、`cache:usermenus:`、`cache:userperms:`、`cache:dict:`、`captcha:`），不可在数据库或前端配置
- **新增** 7 条按钮型权限码：`monitor:server:view`、`monitor:runtime:view`、`monitor:db:view`、`monitor:cache:view`、`monitor:cache:clear`、`monitor:oss:view`、`monitor:dependency:view`
- **新增** 启动 bootstrap：`server/initialize/db_tables.go` 中追加 `ensureServerMonitorMenuApi()`，幂等创建菜单/按钮/API/sys_menu_api 绑定
- **新增** 历史升级 SQL：`server/sql/add_server_monitor.sql`、`server/sql/move_log_audit_menus_to_monitor.sql`，给已存在安装补齐菜单/权限/绑定并迁移日志菜单父级
- **修改** 路由：`server/router/modules/` 下新增 `monitor.go`，注册 `/api/v1/monitor/server`、`/runtime`、`/db`、`/redis`、`/redis/clear`、`/oss`、`/dependency`
- **修改** 前端：`web/src/views/admin/monitor/server/` 新页面，保持现有 Ant Design Vue 后台工具页风格，所有写按钮走 `v-permission`
- **不修改** 现有 `登录日志`、`操作日志` 页面功能，仅调整菜单父级

## Capabilities

### New Capabilities

- `server-monitor`: 服务运行时健康可见性。聚合主机系统指标、Go 进程运行时（goroutine/heap/GC）、数据库连接池、Redis 信息、OSS 健康，并在受控前缀下清理 Redis 缓存。

### Modified Capabilities

（无 — 不改动既有能力契约）

## Impact

**代码与目录**:
- 后端：`server/service/monitorsvc/`（新）、`server/api/v1/monitor.go`（新）、`server/router/modules/monitor.go`（新）、`server/initialize/db_tables.go`（追加幂等 bootstrap 函数）
- 前端：`web/src/views/admin/monitor/server/`（新，含 `index.vue` + 各 Tab 子组件）、`web/src/api/monitor.ts`（新）
- SQL：`server/sql/add_server_monitor.sql`、`server/sql/move_log_audit_menus_to_monitor.sql`（新，仅供存量库升级）

**依赖**:
- 引入 `github.com/shirou/gopsutil/v3`（CPU / Mem / Disk / Host / Load 跨平台采集；MIT，活跃维护）。**这是本次唯一新增的第三方依赖**。
- 复用现有 `redis.Client`、`gorm.DB`、`oss` 实例，不新增连接

**性能**:
- 监控接口默认**按需调用**，不开后台轮询；前端用户主动点 Tab 才拉取
- `gopsutil` 单次采集毫秒级；`runtime.MemStats` 几乎零开销
- "按 key 前缀计数"使用 `SCAN MATCH count=200`，分批扫描，限定最长 5000 次迭代防 DoS

**安全**:
- 全部接口需 JWT + Casbin 鉴权
- 缓存清理强制白名单前缀；越权前缀直接 400 拒绝（policy 写在代码里，不在配置）
- 主机指标不暴露绝对路径、不暴露环境变量、不暴露用户名等敏感字段
- 操作日志：缓存清理动作必须落操作日志（actor + prefix + 影响 key 数量）

**回滚**:
- 删除菜单 + 按钮 + API + sys_menu_api 绑定的若干行（可通过菜单管理页面手动删除即可）
- 删除前端路由文件夹、删除后端模块目录
- `gopsutil` 依赖回滚通过 `go mod tidy` 移除
- 不涉及任何业务表 schema 变更，无数据迁移

## Out of Scope

- **不做**任意命令执行 / SSH 远程命令 / 服务器文件浏览。这类能力即使加权限码也会成为攻击面收口，留给专业运维平台。
- **不做** SaiAdmin 的「数据表维护」（`OPTIMIZE TABLE` / 碎片整理 / DROP）—— 已在评估阶段排除。
- **不做** 前端实时图表（如 ECharts 历史趋势）；初版只展示**当前快照**。历史趋势需要时序存储，另起 change。
- **不做**告警 / 推送 / 阈值配置；仅展示，不响应。
- **不在** 数据库 / 系统配置中暴露任何缓存清理白名单的开关；白名单是**代码定义的安全边界**，符合 AGENTS.md「security-sensitive system config 应该 backend-owned policy」。
- **不**新增 cron 任务清单 / 调度行为；定时任务由独立 OpenSpec change `add-cron-task` 承载。

## Spec Lifecycle

实施完成后：
- `server-monitor` 作为新 capability 落到 `openspec/specs/server-monitor/spec.md`
- 本 change 走 `/openspec-archive-change` 归档
- 不需要更新 / 归档既有 spec
