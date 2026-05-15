zai## 1. 探查 & 准备

- [x] 1.1 阅读 `ensureLogAuditMenus` / `ensureFileUploadApiAccess` / `ensureAIMenuApiBindings`，确认幂等模式（FirstOrCreate + Attrs，保留已有 NeedAuth/metadata）
- [x] 1.2 阅读 `server/api/v1/log.go`、`server/router/modules/log.go`，确认 API 用 `package v1` flat 文件 + service 单例 + `R(rg, ...)` helper 注册
- [x] 1.3 探查菜单：当前日志入口存在历史结构差异。最终决策：统一到 `运维监控 (/monitor)`，下挂 `服务监控`、`操作日志`、`登录日志`
- [x] 1.4 全局对象在 `server/global/global.go`：`global.DB *gorm.DB`、`global.Redis *redis.Client`、`global.Log *zap.SugaredLogger`、`global.Enforcer`；OSS 实例不在 global，需要从 `service.Storage` / `configsvc` 解析当前激活配置
- [x] 1.5 添加 `github.com/shirou/gopsutil/v3` 依赖（go 1.24，go.mod 已支持新版）

## 2. 后端 service 层

- [x] 2.1 创建 `server/service/monitorsvc/dto.go`：定义 `ServerInfo`、`RuntimeInfo`、`DBStats`、`RedisInfo`、`OssHealth`、`DependencyHealth`、`ClearCacheResult` 各 DTO，字段白名单化
- [x] 2.2 创建 `server/service/monitorsvc/server.go`：`CollectServerInfo()` 用 gopsutil 采 Host/CPU/Mem/Disk/Load + 进程 PID/StartedAt；CPU 用 `cpu.Percent(0, false)` 非阻塞调用
- [x] 2.3 创建 `server/service/monitorsvc/runtime.go`：`CollectRuntime()` 读 `runtime.ReadMemStats` + `runtime.NumGoroutine` + `runtime.NumCPU` + `runtime.Version` + 启动时间
- [x] 2.4 创建 `server/service/monitorsvc/db.go`：`CollectDB(ctx)` 读 `gorm.DB.DB().Stats()`，加一次 `Ping(ctx)` 决定 `reachable`，5s 超时
- [x] 2.5 创建 `server/service/monitorsvc/cache.go`：定义实际缓存前缀白名单 `cache:userinfo:` / `cache:usermenus:` / `cache:userperms:` / `cache:dict:` / `captcha:`；实现 Redis INFO、SCAN 计数、白名单清理与 `cache_prefix_not_allowed`
- [x] 2.6 创建 `server/service/monitorsvc/oss.go`：`CollectOss(ctx)` 读当前激活 OSS 配置；本地存储返回未启用；云存储用 3s timeout + Exists 探测，不返回签名 URL/凭证
- [x] 2.7 创建 `server/service/monitorsvc/dependency.go`：`CollectDependency(ctx)` 并行调用 DB ping / Redis ping / OSS 健康，1s 总超时控制
- [x] 2.8 单元测试 `server/service/monitorsvc/cache_test.go`：白名单拒绝 / 空 prefix 拒绝 / 通配符 prefix 拒绝 / 迭代上限触发 truncated；用 `miniredis` 做隔离测试

## 3. 后端 API + 路由

- [x] 3.1 创建 `server/api/v1/monitor.go`：`Server`、`Runtime`、`DB`、`Redis`、`ClearCache`、`Oss`、`Dependency` 7 个 handler，统一用现有 response helper
- [x] 3.2 `ClearCache` handler 走现有 `middleware/operation_log.go` 自动记录成功/非白名单/失败响应；路由 metadata 为「清理 Redis 缓存」，非白名单返回业务码 400
- [x] 3.3 创建 `server/router/modules/monitor.go`：注册 `/monitor/server`、`/runtime`、`/db`、`/redis`、`/redis/clear`、`/oss`、`/dependency`，全部通过 `registry.WithAuth()` 注册并由统一私有路由中间件鉴权
- [x] 3.4 后端构建：`cd server && go build ./...` —— 已通过

## 4. 启动 bootstrap & SQL 升级脚本

- [x] 4.1 在 `server/initialize/db_tables.go` 中实现 `ensureServerMonitorMenuApi()`：幂等创建/更新 `运维监控`（如缺）、`服务监控` 菜单、7 个按钮、7 条 SysApi（`need_auth=true`）、对应 sys_menu_api 绑定；保留已有 description/NeedAuth/metadata
- [x] 4.2 在 `server/initialize/db_tables.go` 主初始化流程中调用 `ensureServerMonitorMenuApi()`，位置在角色权限刷新之前
- [x] 4.3 给内置 `admin` 角色追加 7 个权限码（沿用现有 admin 自动授权模式；**不**触碰 `system_admin`）
- [x] 4.4 创建 `server/sql/add_server_monitor.sql`：与 bootstrap 等价的存量库升级脚本（菜单 + 按钮 + SysApi + sys_menu_api + admin 角色授权 + 基于菜单按钮继承 API 的 Casbin 规则）；**已遵守** `.codex/skills/go-base-sql-upgrade-guardrails/SKILL.md`
- [x] 4.5 集成测试：在 `server/initialize/db_tables_test.go` 加用例验证：（a）首次启动创建全部行；（b）二次启动不覆盖手工修改字段；（c）`admin` 角色获得全部权限且 `system_admin` 不被重授
- [x] 4.6 调整日志菜单结构：`ensureLogAuditMenus()` 将 `操作日志` / `登录日志` 迁移到 `运维监控` 下，只更新父级；空的历史 `系统管理 > 操作审计` 目录自动隐藏；新增 `server/sql/move_log_audit_menus_to_monitor.sql` 覆盖存量库

## 5. 前端 API & 页面

- [x] 5.1 创建 `web/src/api/monitor.ts`：导出 `getServerInfo` / `getRuntimeInfo` / `getDbStats` / `getRedisInfo` / `clearCacheByPrefix` / `getOssHealth` / `getDependencyHealth` 7 个方法 + 对应 TS 类型
- [x] 5.2 创建 `web/src/views/admin/monitor/server/useServerMonitor.ts`：每个 Tab 独立 loading + 数据 ref + 错误态；不开后台轮询，提供 `refresh<Tab>()` 方法
- [x] 5.3 创建 `web/src/views/admin/monitor/server/index.vue`：Ant Design `a-tabs`，5 个 Tab，顶部嵌入 `ServerDependencyBadge`，保持紧凑后台工具页风格
- [x] 5.4 创建 `components/ServerOverviewPanel.vue`：`a-descriptions` + `a-statistic` 展示主机指标；`a-progress` 展示 CPU / Mem / Disk 用量
- [x] 5.5 创建 `components/ServerRuntimePanel.vue`：goroutine / 内存 / GC 统计卡片
- [x] 5.6 创建 `components/ServerDbPanel.vue`：连接池数值卡片 + reachable 徽章
- [x] 5.7 创建 `components/ServerCachePanel.vue`：Redis INFO + 白名单前缀清理卡片；清理按钮 `v-permission="'monitor:cache:clear'"`，点击后 `Modal.confirm` 要求输入完整 prefix 二次确认；清理结果 toast 显示 `deleted` + `truncated`
- [x] 5.8 创建 `components/ServerOssPanel.vue`：根据 `enabled` 显示「未启用」或健康详情；不显示任何凭证字段
- [x] 5.9 创建 `components/ServerDependencyBadge.vue`：DB/Redis/OSS 三色徽章 + 重新检测按钮
- [x] 5.10 路由为 menu-driven：后端 `component=monitor/server/index` 对应 `web/src/views/admin/monitor/server/index.vue`，无需静态路由注册
- [x] 5.11 静态结构自检：已读回新页面/API 文件，确认 light/dark CSS 变量、无营销化 hero 区、`v-permission` 与按钮权限匹配

## 6. 验证 & 收尾

- [x] 6.1 后端单测 + 集成测试：`cd server && go test ./service/monitorsvc ./initialize -count=1` 已通过
- [x] 6.2 后端整体构建：`cd server && go build ./...` 已通过
- [ ] 6.3 路由烟雾测试（手动）：用三种身份（admin / 普通管理员有权限 / 普通管理员无权限）分别访问 7 个接口，验证 200 / 200 / 403
- [ ] 6.4 缓存清理白名单测试（手动）：用合法 prefix `cache:userinfo:` 触发；用非法 prefix `user:` 触发，确认 400；确认操作日志均已落库
- [ ] 6.5 前端手动点击清单（交给用户）：Tab 切换、刷新、清理弹窗二次确认、深色主题、403 提示、断网/Redis 关闭场景下的健康徽章红灯
- [x] 6.6 在 `openspec/changes/add-server-monitor/` 中确认 `proposal.md` / `design.md` / `specs/server-monitor/spec.md` / `tasks.md` 全部存在，`openspec status --change add-server-monitor` 已通过
- [ ] 6.7 实施完成后由用户运行 `/openspec-archive-change` 归档，spec 落到 `openspec/specs/server-monitor/spec.md`
