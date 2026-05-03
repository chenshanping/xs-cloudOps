## Why

现有 CMDB 只能看到主机的静态信息（分组、IP、SSH 凭据）和 SSH 测试是否能连通，**看不到主机实时运行状态**：CPU、内存、磁盘、负载、是否还活着。

继续只靠“SSH 测试通过 = 在线”会有三个问题：

1. 测试通过不等于**现在还活着**，离线感知滞后。
2. 无法回答最基础的 Ops 问题：这台机现在累不累、磁盘是不是满了。
3. 后续要做主机驱动的自动化（发告警、阻止把负载高的机器当部署目标等）没有数据基础。

借鉴同类项目（`deviops`）的监控架构，它用 **自研 agent + Prometheus + Pushgateway**，功能完整但对当前阶段太重：

- 需要额外部署 Prometheus、Pushgateway 两个服务
- agent 由后端**动态编译**（要求服务器有 Go 工具链）
- 心跳 token 硬编码共享，鉴权偏弱

本轮选择**最小闭环 MVP**：

- 交付预编译 agent 二进制，通过 CMDB 已有的 SSH 凭据自动部署
- agent 直接把 JSON 指标 + 心跳 POST 到后端新接口
- 后端只存**最新快照**到数据库，不引入 Prometheus
- 每台主机独立 token，支持吊销

把“Prometheus + 历史曲线 + 告警”显式留到下一轮，避免一次引入两个外部服务 + 两张新表 + 一套新协议 + 一套曲线 UI。

## What Changes

- 新增一个独立的 `caelor-agent` Go 子项目，负责 CPU / 内存 / 磁盘 / 负载 / 运行时长的采集，按固定周期向后端上报 JSON 指标与心跳。
- 新增后端 `CMDB 主机监控` 模块：
  - 主机 agent 记录表（token、状态、最后心跳、版本）。
  - 主机指标最新快照表。
  - Agent 注册/心跳/指标上报接口（公开路由，token 校验）。
  - 主机监控查询接口（管理员 JWT 授权）。
  - 定时任务：超时无心跳 → 主机标为离线。
- 扩展 SSH 凭据现有部署能力：
  - 新增“部署 agent”操作，后端通过 SSH 推送 agent 二进制 + 启动脚本 + per-host token。
  - 新增“卸载 agent”操作。
- 前端 CMDB 主机模块：
  - 列表增列：CPU、内存、最后心跳。
  - 详情抽屉新增系统监控面板：当前快照。
  - 操作列增加 `部署监控 / 卸载监控` 按钮。
- 权限与菜单：
  - 新增 `cmdb:host:monitor:deploy` / `cmdb:host:monitor:view` 权限码。
  - 后台监控相关路由统一挂在 CMDB 模块下。

### Out of Scope

- **不引入 Prometheus 或 Pushgateway**。
- **不实现**监控指标的历史查询、曲线、告警规则、通知通道。
- **不替换** `cmdb_host.status` 的现有语义；在线判定仍沿用现有字段，本次只在超时时把 `status` 从 1 改为 3。
- **不动态编译 agent**；交付的是预编译多平台二进制。
- 不支持 agent 自升级；升级靠“卸载 + 重新部署”。
- 不改造现有 SSH 终端、凭据、主机 CRUD 的业务语义。
- 不新增邮件/飞书/企微通知通道。

## Capabilities

### New Capabilities

- `cmdb-host-monitoring-mvp`: 已绑定 SSH 凭据的主机可以被部署监控 agent，后端必须按 per-host token 接收心跳和基础资源指标，并持久化最新快照；在 agent 心跳超时时必须把主机标为离线。

### Modified Capabilities

- None. 当前 `cmdb_host_management` 能力尚未作为 source-of-truth spec 归档，故本轮以“新能力”呈现。

## Impact

- Backend:
  - 新增 `server/model/cmdb_host_agent.go`, `server/model/cmdb_host_metric.go`
  - 新增 `server/service/cmdb/agent.go`, `server/service/cmdb/metric.go`
  - 新增 `server/api/v1/cmdb_monitor.go`
  - 扩展 `server/router/modules/cmdb.go`（公开心跳/上报路由 + 私有监控路由）
  - 扩展 `server/initialize/db_tables.go` AutoMigrate 列表 + 菜单/API 元数据
  - 新增 `server/sql/add_cmdb_host_monitoring.sql`（无外键、MySQL 8、幂等）
  - 新增定时任务：`server/cron/cmdb_offline_checker.go` 或在现有调度器注册
  - 新增静态资源目录 `server/assets/agent/`（存放 agent 多平台二进制）
- Frontend:
  - 扩展 `web/src/types/cmdb.ts`
  - 扩展 `web/src/api/cmdb.ts`
  - 扩展 `web/src/views/admin/cmdb/host/index.vue`（列 / 操作 / 抽屉）
  - 扩展 `web/src/views/admin/cmdb/host/components/HostDetailDrawer.vue`
- New repo:
  - `caelor-agent/`（独立 Go 模块，可以放 `apps/caelor-agent` 或单独仓库，本轮用 monorepo 路径）
- Permissions / Menu:
  - 新菜单按钮权限码 `cmdb:host:monitor:deploy`, `cmdb:host:monitor:view`
  - 默认授权给 admin / system_admin
- DB:
  - 两张新表：`cmdb_host_agent`, `cmdb_host_metric`
  - 已有 `cmdb_host.status` 行为补充：超时 = 3
- Rollback:
  - 关闭定时 offline checker
  - 下线监控相关路由
  - 前端回滚操作列与抽屉面板
  - 保留两张新表数据不回滚
  - SQL 升级脚本为幂等 add-column style，不回退结构
