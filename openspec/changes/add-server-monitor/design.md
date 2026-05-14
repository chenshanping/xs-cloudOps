## Context

go-base 当前已有「登录日志 / 操作日志」两个运维入口，沉淀在 `web/src/views/admin/monitor/` 目录。技术栈：

- 后端 Go 1.21+ + Gin + GORM + Casbin + JWT，Redis、MySQL、可选 OSS（七牛/阿里/腾讯/MinIO/AWS）
- 前端 Vue 3 + TypeScript + Vite + Pinia + Ant Design Vue
- 权限模型：Casbin RBAC + sys_menu/按钮 + sys_menu_api 绑定，启动 bootstrap 在 `server/initialize/db_tables.go`

参考过的相邻实现：

- `server/service/ai/` 服务层组织
- `server/service/file/file.go` 服务包结构 + 错误返回 + DTO
- `server/api/v1/log/`、`server/router/modules/log.go`（Login/Operation log API）—— 路由注册参考
- `server/initialize/db_tables.go` 中 `ensureAIToolMenus()` / `ensureAIMenuApiBindings()` —— 幂等 bootstrap 模板
- `web/src/views/admin/monitor/operation-log/index.vue`、`web/src/views/admin/dashboard/index.vue` —— 后台工具页的紧凑卡片/Tab 风格参考

约束：

- **必须遵守 AGENTS.md**：security-sensitive policy 走代码不走 DB；不引入数据库 FK；不引入代码生成器
- **不允许 `KEYS *`**，仅 `SCAN`；分批处理，单次操作设迭代上限
- **缓存清理白名单是代码级 const**，不在 DB / 配置 / 前端可改
- 启动 bootstrap 必须保留 SysApi 已有 NeedAuth/metadata（参考 `ensureFileUploadApiAccess` 写法）
- 系统级 admin code `admin` 自动授权；`system_admin` 由用户自行配置（按现有 bootstrap 行为）

## Goals / Non-Goals

**Goals:**

- 一个 5 Tab 的服务监控页面（系统 / Runtime / 数据库 / 缓存 / OSS），全部按需拉取
- 后端 5 类只读接口 + 1 类受控写入（缓存按前缀清理）
- 完整的 Casbin 权限链路：菜单 → 按钮 → API → sys_menu_api 绑定，启动幂等 + SQL 升级脚本双轨
- 缓存清理操作必须落 `sys_log_operation`，包含 actor、prefix、影响 key 数量

**Non-Goals:**

- 不做历史趋势图、不做时序存储、不做告警/阈值配置
- 不做远程命令执行、SSH、文件浏览
- 不做 `OPTIMIZE TABLE` / 表碎片整理
- 不在 DB / 系统配置中暴露白名单（前缀白名单是代码常量）
- 不为该模块开后台轮询；用户主动触发才采集

## Decisions

### D1：使用 `gopsutil/v3` 采集主机指标

**选择**：`github.com/shirou/gopsutil/v3`

**理由**：
- 跨平台（Win/Linux/macOS）一致 API，社区主流，MIT
- 不需要 root 权限，不需要 cgo
- 历史稳定，go-base 部署不锁定 OS

**替代方案**：
- 自己读 `/proc` —— 不跨平台，Windows 缺失
- 调用系统命令 —— 不安全、不跨平台、性能差

### D2：采集策略——按需 vs 后台轮询

**选择**：完全按需。前端用户切到某 Tab 才请求；后端每次请求实时采集。

**理由**：
- 监控页非高频访问，轮询无收益
- 避免常驻 goroutine 引入额外故障面
- `runtime.MemStats` 和 `gopsutil` 单次调用开销低（数毫秒）

**风险**：单次采集会触发 Go runtime 的 stop-the-world 一次（`MemStats`）。`gopsutil` CPU 采样默认 1s 间隔——我们用**非阻塞调用**（`cpu.Percent(0, false)`），第一次返回 0 但后续准确，可接受。

### D3：缓存清理白名单——代码常量

**选择**：

```go
// server/service/monitorsvc/cache.go
var allowedCachePrefixes = []string{
    "cache:userinfo:",  // 用户信息缓存
    "cache:usermenus:", // 用户菜单缓存
    "cache:userperms:", // 用户权限缓存
    "cache:dict:",      // 字典缓存
    "captcha:",         // 验证码
}
```

任何不在白名单的 prefix 直接 `400 Bad Request`。

**理由**：AGENTS.md 明确 security-sensitive 策略走代码：「Default to backend-coded policy. Do not add admin UI or DB-configurable toggles for these controls」。

**替代方案**（拒绝）：
- DB 表存白名单 + UI 编辑 → 等于内置后门，超管误操作可清空 session/token 等关键 key
- 前端硬编码 → 绕过简单（直接打 API）

### D4：扫描 Redis key——SCAN 而非 KEYS

**选择**：`SCAN MATCH <prefix>* COUNT 200`，最大迭代 5000 次（即扫到 ~100 万 key 就停）。

```go
const (
    scanBatch    = 200
    scanMaxIter  = 5000
)
```

**理由**：
- `KEYS` 在生产环境是大忌（O(N) 阻塞）
- 5000 次迭代上限防止恶意/异常前缀匹配大量 key 导致接口长时间挂起
- 超出上限返回部分结果 + `truncated: true` 标记，前端友好提示

### D5：DB 指标——`sql.DBStats`

**选择**：直接读 GORM 底层 `*sql.DB` 的 `Stats()`。

返回字段：MaxOpenConnections / OpenConnections / InUse / Idle / WaitCount / WaitDuration / MaxIdleClosed / MaxLifetimeClosed。**不**返回当前慢查询、不返回 schema 信息。

**理由**：标准库自带，零成本，足够诊断连接池打满 / 等待激增。

### D6：OSS 健康——已配置才检测

**选择**：读 `configsvc` 当前激活的 OSS 配置，若启用则做一次 `Bucket()` / 列出根目录 1 项的轻探测；未启用直接返回 `{enabled: false}`。

**理由**：
- 五种 OSS 后端实现已存在 `server/service/oss/`
- 健康探测走最便宜路径，避免误产生流量费

### D7：API 路径与权限码

| 接口 | 方法 | 权限码 |
|---|---|---|
| `/api/v1/monitor/server` | GET | `monitor:server:view` |
| `/api/v1/monitor/runtime` | GET | `monitor:runtime:view` |
| `/api/v1/monitor/db` | GET | `monitor:db:view` |
| `/api/v1/monitor/redis` | GET | `monitor:cache:view` |
| `/api/v1/monitor/redis/clear` | POST | `monitor:cache:clear` |
| `/api/v1/monitor/oss` | GET | `monitor:oss:view` |
| `/api/v1/monitor/dependency` | GET | `monitor:dependency:view` |

`/dependency` 是聚合健康概览（DB ping + Redis ping + OSS 健康），用于首屏 Tab 顶部的健康徽章。

### D8：菜单结构

```
运维监控 (一级目录)
├── 服务监控        （新）
│   ├── 查看系统    button → monitor:server:view
│   ├── 查看运行时  button → monitor:runtime:view
│   ├── 查看数据库  button → monitor:db:view
│   ├── 查看缓存    button → monitor:cache:view
│   ├── 清理缓存    button → monitor:cache:clear
│   ├── 查看 OSS   button → monitor:oss:view
│   └── 健康概览    button → monitor:dependency:view
├── 操作日志        （迁移既有菜单）
└── 登录日志        （迁移既有菜单）
```

`登录日志 / 操作日志` 属于运维可观测性入口，随本 change 从 `系统管理 > 操作审计` 迁移到 `运维监控` 下。迁移只更新父级，不覆盖既有日志菜单的名称、路径、组件、图标、排序等可维护元数据；若历史 `操作审计` 目录迁移后无子菜单，则隐藏该空目录。

### D9：前端目录与组件拆分

```
web/src/views/admin/monitor/server/
├── index.vue                   # Tab 容器，保持紧凑后台工具页风格
├── components/
│   ├── ServerOverviewPanel.vue # 系统 Tab
│   ├── ServerRuntimePanel.vue  # Go runtime Tab
│   ├── ServerDbPanel.vue       # DB Tab
│   ├── ServerCachePanel.vue    # Redis Tab + 清理操作
│   ├── ServerOssPanel.vue      # OSS Tab
│   └── ServerDependencyBadge.vue # 顶部健康徽章
└── useServerMonitor.ts         # 各 Tab 拉取 + loading 状态
```

每个 Panel 内部用 `a-descriptions` + `a-statistic` + `a-progress`，遵循现有 admin 设计语言；不引入图表库。

写按钮统一 `v-permission="'monitor:cache:clear'"`，再次确认后才发请求（`Modal.confirm` + 输入 prefix 二次确认）。

## Risks / Trade-offs

- [指标暴露过多敏感信息]
  → Mitigation：返回字段白名单化；不返回 `os.Hostname()` 之外的网络/路径；不返回环境变量。所有 DTO 字段在 `dto.go` 中明确声明。

- [缓存清理误删导致系统不可用]
  → Mitigation：前端二次确认必须输入完整 prefix；后端再次校验白名单；操作日志强制落库。允许的三个前缀都是**可重新生成**的（菜单/权限缓存有自动 invalidation，captcha 短期过期）。

- [`SCAN` 在大库上慢]
  → Mitigation：迭代上限 + `truncated` 标记 + 前端 5s loading；超出预期会显示「扫描已截断，请联系运维」。

- [`gopsutil` 在某些容器化环境（cgroup v1/v2 差异）读到宿主机指标]
  → Mitigation：在 DTO 里标记 `data_source: "host" | "container"`；初版按 host 实现，container-aware 留给后续 change。

- [`runtime.MemStats` 触发 STW]
  → Mitigation：每次调用 STW 数十微秒，可忽略。但**不要**在循环中重复调用，单请求只取一次。

- [`gopsutil` 引入间接依赖膨胀]
  → Mitigation：v3 主包依赖 ~5 个间接依赖，体积可控；锁定到一个稳定版本，不引入 v4 alpha。

## Migration Plan

**部署步骤**：
1. 后端代码 + `go mod tidy` 拉 `gopsutil/v3`
2. 启动后 `ensureServerMonitorMenuApi()` 自动建菜单/按钮/API/绑定（首次启动幂等）
3. 存量库可选执行 `server/sql/add_server_monitor.sql`（与 bootstrap 等价，仅在 bootstrap 未运行的环境使用）
4. 默认 `admin` 角色自动获得 7 个权限码（按现有 bootstrap 行为）
5. 其他角色由超管自行勾选

**回滚**：
1. 数据库：手工删除菜单 + 按钮 + API + sys_menu_api 绑定相关行（或反向 SQL 脚本）；不涉及业务数据
2. 代码：删除 `server/service/monitorsvc/`、`server/api/v1/monitor.go`、`server/router/modules/monitor.go`、bootstrap 函数；从前端删 `web/src/views/admin/monitor/server/`
3. `go mod tidy` 清理 `gopsutil` 依赖

**风险等级**：低。无业务表 schema 变更，回滚干净。

## Open Questions

- **运维监控一级目录**是否已存在？→ 已确认需统一到 `运维监控 (/monitor)`，并将现有审计日志菜单迁移到该目录。
- **OSS 健康探测的频次保护**：用户连点 Tab 是否需要节流？初版不节流，由前端 loading 阻断；如成本暴露，后续加 5s 缓存。

