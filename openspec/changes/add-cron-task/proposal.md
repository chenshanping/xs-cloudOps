## Why

go-base 当前没有内置任务调度能力。**系统每次重启后**，过期登录日志、过期操作日志会一直堆积，运维只能手工 SQL 清理或写外部 cron 调外部脚本，**没有可观测**（执行成功/失败、耗时、上次执行时间均不可见）。

参考 SaiAdmin 的「定时任务 + 执行日志」设计，结合 Go 生态成熟的 `robfig/cron/v3`，新增**任务函数代码注册（白名单）+ DB 存调度配置**的定时任务模块。**绝不**支持从 DB 读取可执行命令字符串——那等于内置 RCE。

## What Changes

- **新增** `运维管理 > 定时任务` 菜单 + `运维管理 > 任务执行日志` 子菜单
- **新增** 后端 `service/cronsvc/`：调度器（基于 `robfig/cron/v3`）+ 任务注册中心（`map[code]TaskFunc` 代码硬编码白名单）+ DB 表 `sys_cron_task` / `sys_cron_log`
- **新增** **代码注册的内置任务白名单**（初版两个）：
  - `cleanup_login_logs`：按保留天数清理 `sys_log_login` 过期记录
  - `cleanup_operation_logs`：按保留天数清理 `sys_log_operation` 过期记录
  - 任务参数（保留天数、单次最大删除条数）走 `params` JSON 字段，但**不允许动态变更任务函数**
- **新增** 任务 CRUD（增删改查 + 启停 + 立即执行）+ 执行日志查询/清理
- **新增** 8 条按钮型权限码：`monitor:cron:create`、`update`、`delete`、`enable`、`disable`、`runNow`、`view`、`logs:view`
- **新增** 启动 bootstrap：`ensureCronTaskMenuApi()` 幂等创建菜单/按钮/API/sys_menu_api 绑定 + 启动时把 DB 中 `status=enabled` 的任务全部加载到调度器
- **新增** 历史升级 SQL：`server/sql/add_cron_task.sql`（菜单 + 权限 + API 绑定 + 两个内置任务的种子记录，默认 `status=disabled` 不自动跑）
- **新增** 表：`sys_cron_task`、`sys_cron_log`，纳入 baseline schema 与 GORM AutoMigrate
- **修改** `server/initialize/`：在 GORM 初始化与 router 注册之间插入 `cronsvc.Start(ctx)`；`Shutdown` 时优雅停止（等待运行中任务最多 30s）
- **不修改** 任何业务表 schema

## Capabilities

### New Capabilities

- `cron-task`: 系统调度能力。后台可视化管理代码白名单中的任务调度（cron 表达式、启停、单次手动触发），并保留每次执行的成功/失败/耗时日志。

### Modified Capabilities

（无）

## Impact

**代码**:
- 后端：`server/service/cronsvc/`（新）、`server/api/v1/cron/`（新）、`server/router/modules/cron.go`（新）、`server/model/sys_cron_task.go` + `sys_cron_log.go`（新）、`server/initialize/db_tables.go` 追加 bootstrap、`server/initialize/cron.go`（新，启动调度器入口）
- 前端：`web/src/views/admin/monitor/cron-task/`、`web/src/views/admin/monitor/cron-log/`（新）、`web/src/api/cron.ts`（新）
- SQL：`server/sql/add_cron_task.sql`（新）

**依赖**:
- `github.com/robfig/cron/v3`（MIT，最广泛使用，活跃维护）。**本次唯一新增依赖**

**性能**:
- 调度器进程内单例，goroutine 数量 = 启用的任务数；空载内存几十 KB 量级
- 任务执行采用「同任务串行 + 不同任务并行」：每个任务实例携带 `mu sync.Mutex`，前一次未结束则跳过本次（写日志记录 `skipped: true`）
- 单实例部署假设；如需多实例部署，独立 change 增加分布式锁

**安全**:
- **任务函数完全代码注册**，DB 仅存 `task_code`、cron 表达式、参数 JSON、状态；`task_code` 不在白名单则拒绝创建/启用
- 任务参数限定为已注册任务声明的参数 schema（如 `cleanup_login_logs` 只接受 `retain_days: int`、`batch_limit: int`），其他字段忽略
- 所有写操作进操作日志
- cron 表达式必须通过 `cron.Parser` 校验
- 单次 `runNow` 必须额外鉴权 `monitor:cron:runNow`，不复用 `update`

**回滚**:
- 停止调度器（重启服务即停）
- 删除菜单 + 按钮 + API + sys_menu_api 绑定
- 可选 drop `sys_cron_task` + `sys_cron_log` 表（包含执行历史）
- `go mod tidy` 清理 `robfig/cron`

## Out of Scope

- **不支持**从 DB 读 shell 命令 / Go 表达式 / 任意脚本字符串执行——这是核心安全边界，**永久不做**
- **不支持** 任务依赖编排（前置任务、链式触发）—— 复杂度太高，单独提案
- **不支持** 分布式部署去重锁（单实例假设）；多实例上线前需另起 change
- **不支持** 任务结果持久化大对象（执行日志只存状态 + 摘要文本，不存任务返回的 payload）
- **不支持**前端动态注册新任务函数 / 上传 .so / 反射执行
- **不在** 系统配置中暴露任务白名单的开关；白名单 = 代码常量

## Spec Lifecycle

- 新 capability `cron-task` → 写入 `openspec/changes/add-cron-task/specs/cron-task/spec.md`
- 实施完成后归档到 `openspec/specs/cron-task/spec.md`
- 不需要修改/归档既有 spec
