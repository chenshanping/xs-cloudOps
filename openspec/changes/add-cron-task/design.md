## Context

go-base 后端长期常驻进程（Gin），适合进程内任务调度。当前没有清理过期日志的机制，操作日志/登录日志只增不减。

参考实现：
- `server/service/ai/` 服务层组织
- `server/service/file/file.go` 服务包结构
- `server/initialize/db_tables.go` bootstrap 模式
- `server/router/modules/log.go` 路由注册风格

约束：
- AGENTS.md：security-sensitive policy 走代码；soft-delete + unique index 必须显式处理；无数据库 FK
- 单实例部署假设（多实例后续 change 处理）
- `robfig/cron/v3` 是 Go 生态的事实标准

## Goals / Non-Goals

**Goals:**

- 一个内嵌进程的 cron 调度器，启动时从 DB 加载启用任务、运行时支持增删改启停
- 「任务函数代码注册（白名单） + DB 存调度配置」严格分离的安全模型
- 每次执行落日志：开始时间、结束时间、状态（success/failure/skipped）、耗时、stdout/err 摘要（截断至 4KB）
- 内置 2 个低风险任务：清理过期登录日志、清理过期操作日志
- 完整 RBAC：8 个权限码 + 菜单/按钮/API/绑定 + 启动 bootstrap + SQL 升级脚本

**Non-Goals:**

- 不支持任意命令/脚本/反射调用
- 不支持分布式调度去重
- 不做任务编排 / DAG / 依赖触发
- 不持久化任务大返回值

## Decisions

### D1：调度库 — `robfig/cron/v3`

**理由**：MIT、活跃、API 稳定、支持秒级表达式（v3 默认 5 段，可选 6 段）、Job 接口简洁、自带 `Parser`、社区使用面最广。

**替代**：
- `go-co-op/gocron` —— API 偏重，依赖更多
- 自己写 ticker —— cron 表达式校验复杂，不重复造轮子

### D2：任务白名单 = 代码常量 map

```go
// server/service/cronsvc/registry.go
type TaskFunc func(ctx context.Context, params map[string]any) (summary string, err error)

var registry = map[string]TaskHandler{
    "cleanup_login_logs": {
        Run: cleanupLoginLogs,
        ParamSchema: map[string]string{
            "retain_days": "int",
            "batch_limit": "int",
        },
        Description: "清理早于 N 天的登录日志，每批最多删除 batch_limit 行",
    },
    "cleanup_operation_logs": { ... },
}
```

DB 中 `task_code` 不在 `registry` → 创建/启用直接 `400`。

### D3：参数 schema 校验

任务参数走 JSON 字段 `params`。每个任务在 `ParamSchema` 中声明合法字段名 + 类型；保存时只保留 schema 中声明的字段，其他丢弃；类型不符返回 400。

不引入 JSON Schema 库；手写一个 30 行的简单校验器，保持依赖最小。

### D4：调度器生命周期

```
initialize.Setup()
  → cronsvc.Start(ctx)        # 启动 cron.Cron
  → cronsvc.LoadFromDB()      # 加载 status=enabled 的任务
graceful shutdown
  → cronsvc.Stop()            # cron.Stop() + 等待运行中任务完成（30s 超时）
```

### D5：同任务串行执行

每个任务在 registry 中绑定一个 `sync.Mutex`。Cron 触发时 `TryLock`：成功则跑，失败则写一条 `status=skipped` 日志并跳过本次。

理由：避免长任务在下次触发时叠加；管理员可通过执行日志看到 skipped 频率，从而调整 cron 表达式。

### D6：任务执行日志结构

```sql
CREATE TABLE sys_cron_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL,
    task_code VARCHAR(64) NOT NULL,
    started_at DATETIME NOT NULL,
    finished_at DATETIME NULL,
    duration_ms BIGINT NULL,
    status VARCHAR(16) NOT NULL,   -- success / failure / skipped / running
    summary TEXT NULL,             -- 任务返回摘要，最长 4KB
    error_message TEXT NULL,
    triggered_by VARCHAR(32) NOT NULL, -- schedule / manual / startup
    actor_user_id BIGINT NULL,     -- 手动触发时
    created_at DATETIME NOT NULL,
    INDEX idx_task_id (task_id),
    INDEX idx_started_at (started_at)
);
```

无 FK（遵守 AGENTS.md），`task_id` 仅引用，应用层处理。

### D7：表 `sys_cron_task` 字段

```sql
CREATE TABLE sys_cron_task (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(64) NOT NULL,         -- 实例 code（用户起名）
    task_code VARCHAR(64) NOT NULL,    -- 注册中心 key
    name VARCHAR(128) NOT NULL,
    cron_expr VARCHAR(64) NOT NULL,
    params JSON NULL,
    status TINYINT NOT NULL DEFAULT 0, -- 0 disabled / 1 enabled
    last_run_at DATETIME NULL,
    last_status VARCHAR(16) NULL,
    last_duration_ms BIGINT NULL,
    next_run_at DATETIME NULL,
    remark VARCHAR(255) NULL,
    sort INT NOT NULL DEFAULT 0,
    created_by BIGINT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    INDEX idx_status (status),
    INDEX idx_task_code (task_code)
);
```

**Soft delete + unique index**：根据 AGENTS.md「Soft delete + unique index 必须显式解决」，本表 `code` 不设 unique 索引（同名实例允许，因为现实中可能"删了又重建"是常态）。改为应用层在 service 中校验「未删除记录中 code 唯一」，由 service 层主动去重。

### D8：API 路径与权限

| 接口 | 方法 | 权限码 |
|---|---|---|
| `/admin/monitor/cron-task` | GET | `monitor:cron:view` |
| `/admin/monitor/cron-task` | POST | `monitor:cron:create` |
| `/admin/monitor/cron-task/:id` | PUT | `monitor:cron:update` |
| `/admin/monitor/cron-task/:id` | DELETE | `monitor:cron:delete` |
| `/admin/monitor/cron-task/:id/enable` | POST | `monitor:cron:enable` |
| `/admin/monitor/cron-task/:id/disable` | POST | `monitor:cron:disable` |
| `/admin/monitor/cron-task/:id/run` | POST | `monitor:cron:runNow` |
| `/admin/monitor/cron-log` | GET | `monitor:cron:logs:view` |

`runNow` 同步触发但不阻塞 HTTP 响应：handler 起独立 goroutine 跑任务，立即返回 `log_id`，前端轮询 `/cron-log?id=` 获取结果。

### D9：菜单结构

```
运维管理
├── 登录日志              （已有）
├── 操作日志              （已有）
├── 服务监控              （add-server-monitor change）
└── 定时任务              （本 change）
    ├── 任务列表
    ├── 任务执行日志
    └── 7 个按钮权限
```

### D10：内置任务实现

```go
// cleanup_login_logs
func cleanupLoginLogs(ctx context.Context, params map[string]any) (string, error) {
    retain := getInt(params, "retain_days", 30)
    batch := getInt(params, "batch_limit", 1000)
    cutoff := time.Now().AddDate(0, 0, -retain)
    var total int64
    for {
        res := db.WithContext(ctx).
            Where("created_at < ?", cutoff).
            Limit(batch).
            Delete(&model.SysLogLogin{})
        if res.Error != nil { return "", res.Error }
        total += res.RowsAffected
        if res.RowsAffected < int64(batch) { break }
        if ctx.Err() != nil { break }
    }
    return fmt.Sprintf("deleted=%d cutoff=%s", total, cutoff.Format(time.RFC3339)), nil
}
```

`cleanup_operation_logs` 同构。

### D11：前端目录

```
web/src/views/admin/monitor/cron-task/
├── index.vue              # 列表 + ProTable
├── components/
│   ├── CronTaskFormDrawer.vue  # 创建/编辑 Drawer
│   ├── CronTaskParamsForm.vue  # 根据 task_code 渲染参数表单
│   └── CronRunNowConfirm.vue   # runNow 二次确认
└── useCronTaskPage.ts

web/src/views/admin/monitor/cron-log/
└── index.vue              # 执行日志列表
```

任务参数表单根据后端返回的 `param_schema` 动态生成（数字字段 = `a-input-number`，字符串 = `a-input`）。

## Risks / Trade-offs

- [代码注册的任务有限，运维提想加任务必须改代码并发版]
  → Trade-off：换来安全。这是刻意决定。文档里要写明「新增任务请提 PR」。

- [单实例假设，多实例部署会重复执行]
  → Mitigation：部署文档明确单实例 / 多实例只在一个节点启用调度器（环境变量 `CRON_ENABLED=false` 关闭非主节点的调度，但 API 仍可读）。多实例分布式锁留独立 change。

- [长任务被 SIGTERM 强杀]
  → Mitigation：graceful shutdown 等待 30s；任务实现内部尊重 ctx，定期 `ctx.Err()` 检查；超时未完成的任务在执行日志中标 `running`，下次启动时自动改 `failure: 异常退出`。

- [`runNow` 被高频点击导致积压]
  → Mitigation：D5 的 mutex 已自然限流；前端按钮在请求中 disable。

- [执行日志膨胀]
  → Mitigation：内置 `cleanup_cron_logs` 任务，但**不在 v1 启用**（避免互锁），v1 提供后台手动清理按钮 + 接受运维定期归档。

- [`params` JSON 字段在不同 DB 兼容性]
  → Mitigation：MySQL 5.7+ 支持 JSON；如需 PostgreSQL/SQLite，model 层用 `datatypes.JSON` 或者 string 存。先 MySQL，文档注明。

## Migration Plan

1. 拉 `robfig/cron/v3` 依赖
2. AutoMigrate 创建 2 张新表
3. bootstrap 创建菜单/按钮/API/绑定 + 给 admin 角色授权 8 个权限码
4. 存量库可选执行 `add_cron_task.sql`
5. 启动调度器，加载 `status=enabled` 任务（首次部署没有，无任何调度发生）
6. 管理员手工创建任务实例并启用

**回滚**：
- 重启服务（无调度器即可）
- 反向 SQL 删菜单/按钮/API/绑定
- 可选 drop 2 张表（注意会丢执行历史）
- `go mod tidy` 清依赖

风险等级：中。涉及新表 + 后台调度器，但**默认任务全部 disabled**，安装后不自动产生任何副作用。

## Open Questions

- 是否需要 cron 表达式可视化（如 `cronstrue` 中文翻译）？建议 v1 只校验 + 返回下次执行时间，文字翻译留 v2。
- `triggered_by=startup` 是否启用？建议**不启用**——重启即触发任务在生产是反模式（容易级联故障）。仅 `schedule` + `manual`。
