## 1. 探查 & 准备

- [x] 1.1 阅读 `server/initialize/db_tables.go`、`server/initialize/init.go` 现有启动流程，确认调度器启动注入位置
- [x] 1.2 阅读 `server/model/sys_log_login.go` + `sys_log_operation.go`，确认表名 / 字段（用于内置清理任务）
- [x] 1.3 阅读 `server/service/file/file.go` 服务层风格 + `server/router/modules/log.go` 路由风格
- [x] 1.4 探查 `运维管理` 一级菜单是否存在（与 `add-server-monitor` change 的探查共用结论；如本 change 先实施则需创建该一级目录）
- [x] 1.5 添加依赖 `github.com/robfig/cron/v3`，`go mod tidy`

## 2. 数据模型 & 迁移

- [x] 2.1 创建 `server/model/sys_cron_task.go`：字段按 design.md D7；`Params` 用 `datatypes.JSON`（先 MySQL）；BaseModel 嵌入；不引入 FK；不在 `code` 上加 unique index
- [x] 2.2 创建 `server/model/sys_cron_log.go`：字段按 design.md D6；索引 `(task_id)`、`(started_at)`
- [x] 2.3 在 `server/initialize/gorm.go`（或现有 AutoMigrate 聚合点）注册两张新表
- [x] 2.4 baseline schema 文件（`server/sql/init.sql` 或对等）追加这两张表的 DDL
- [x] 2.5 检查现有 `server/sql/` 升级脚本命名规则，遵循 `.windsurf/workflows/sql-upgrade-guardrails.md`

## 3. 后端 service 层

- [x] 3.1 创建 `server/service/cronsvc/registry.go`：`TaskHandler{Run, ParamSchema, Description}`；`registry` map；`Register/Get/List` 公开方法（仅供同包初始化使用）
- [x] 3.2 创建 `server/service/cronsvc/tasks_cleanup.go`：实现 `cleanupLoginLogs` 和 `cleanupOperationLogs`，遵守 ctx 取消、按批删除、返回 summary 字符串
- [x] 3.3 在 `cronsvc/init.go`（或 `registry.go` init）中调用 `Register("cleanup_login_logs", ...)`、`Register("cleanup_operation_logs", ...)`
- [x] 3.4 创建 `server/service/cronsvc/scheduler.go`：包装 `cron.Cron`；`Start(ctx)`、`Stop(timeout)`、`AddTask(model)`、`RemoveTask(id)`、`UpdateTask(model)`、`RunNow(id, actorID)`；任务实例携带 `sync.Mutex` 保证同任务串行
- [x] 3.5 创建 `server/service/cronsvc/loader.go`：`LoadFromDB()` 启动时加载 enabled 任务到调度器
- [x] 3.6 创建 `server/service/cronsvc/logsink.go`：`StartLog(taskID, triggeredBy, actorID) → logID`、`FinishLog(logID, status, summary, err)`；panic 恢复
- [x] 3.7 创建 `server/service/cronsvc/params.go`：根据 `ParamSchema` 校验/裁剪/类型转换 params
- [x] 3.8 创建 `server/service/cronsvc/expr.go`：用 `cron.Parser` 校验表达式，返回 `next_run_at`
- [x] 3.9 创建 `server/service/cronsvc/service.go`：CRUD + Enable/Disable/RunNow 业务逻辑（包含「同名 code 在未删除中唯一」的应用层校验）
- [x] 3.10 单测 `server/service/cronsvc/scheduler_test.go`：白名单拒绝、参数 schema 校验、表达式校验、同任务并发被 skipped、graceful shutdown
- [x] 3.11 单测 `server/service/cronsvc/tasks_cleanup_test.go`：用内存 sqlite 验证清理只删过期行、batch 分批、ctx 取消

## 4. 后端 API + 路由

- [x] 4.1 创建 `server/api/v1/cron/task.go`：8 个 handler 对应 D8 表格
- [x] 4.2 创建 `server/api/v1/cron/log.go`：分页查询 + 单条详情
- [x] 4.3 `runNow` handler：起 goroutine 执行，立即返回 `log_id`；并发执行同任务时返回 409
- [x] 4.4 所有写操作落操作日志（包括失败原因）
- [x] 4.5 创建 `server/router/modules/cron.go`：注册路由，挂 JWT + Casbin 中间件；在 `server/router/router.go` 聚合点引入
- [x] 4.6 后端构建：`cd server && go build ./...`

## 5. 启动 bootstrap & SQL 升级

- [x] 5.1 创建 `server/initialize/cron.go`：在 GORM 初始化、router 注册之前调用 `cronsvc.Start(ctx)` + `cronsvc.LoadFromDB()`；shutdown hook 调用 `Stop(30s)`
- [x] 5.2 在 `server/initialize/db_tables.go` 实现 `ensureCronTaskMenuApi()`：幂等创建「定时任务 / 任务执行日志」菜单 + 8 个按钮 + SysApi（`need_auth=true`）+ sys_menu_api 绑定；保留已有 metadata；位置在角色权限刷新之前
- [x] 5.3 给 `admin` 角色追加 8 个权限码（沿用现有 admin 自动授权模式；不触碰 `system_admin`）
- [x] 5.4 创建 `server/sql/add_cron_task.sql`：等价升级脚本（菜单/按钮/SysApi/绑定/Casbin 规则）+ **种子两条 disabled 状态的内置任务**（cleanup_login_logs 默认 retain_days=30, batch_limit=1000；cleanup_operation_logs 同；cron 表达式 `0 2 * * *`）
- [x] 5.5 集成测试 `server/initialize/db_tables_test.go`：首次启动创建全部行；二次启动不覆盖手工修改；admin 获得 8 个权限码

## 6. 前端 API & 页面

- [x] 6.1 创建 `web/src/api/cron.ts`：list / create / update / delete / enable / disable / runNow / logList / logDetail / registry（拉取已注册任务白名单和参数 schema）共 10 个方法 + TS 类型
- [x] 6.2 创建 `web/src/views/admin/monitor/cron-task/index.vue`：ProTable 列表，列含 名称 / task_code / cron 表达式 / 状态开关 / 上次执行时间 / 上次状态 / 下次执行时间 / 操作；操作含编辑、删除、立即执行、查看日志
- [x] 6.3 创建 `cron-task/components/CronTaskFormDrawer.vue`：Drawer 创建/编辑；task_code 选择拉自 `/registry`，禁止手输；选中后联动展示 `param_schema` 描述
- [x] 6.4 创建 `cron-task/components/CronTaskParamsForm.vue`：根据 `param_schema` 动态渲染表单字段（int → `a-input-number`，string → `a-input`，bool → `a-switch`）
- [x] 6.5 创建 `cron-task/components/CronRunNowConfirm.vue`：`Modal.confirm` 二次确认；触发后跳转到执行日志页并自动选中刚生成的 log_id
- [x] 6.6 创建 `cron-task/useCronTaskPage.ts`：列表分页、状态切换、刷新逻辑
- [x] 6.7 创建 `web/src/views/admin/monitor/cron-log/index.vue`：执行日志 ProTable，过滤 task_id / status / 时间范围；点击查看 summary/error 详情
- [x] 6.8 视觉自检：与 role 抽屉的 permission-shell 风格统一；深色模式排版正常

## 7. 验证 & 收尾

- [x] 7.1 后端单测：`cd server && go test ./service/cronsvc/... ./initialize/... -count=1`
- [x] 7.2 后端整体构建：`cd server && go build ./...`
- [ ] 7.3 手动烟雾：创建 `cleanup_login_logs` 任务（disabled）→ runNow → 查看日志 success；改 retain_days=999 → 再 runNow → 不删任何行
- [x] 7.4 安全测试：通过 API 直接 POST `task_code = "evil"` 应被 400 拒绝；params 注入字段被丢弃
- [x] 7.5 同任务串行测试：把 cleanup 改成长任务（暂时 mock sleep）→ 高频触发 runNow → 应返回 409
- [ ] 7.6 graceful shutdown 测试：runNow 长任务 → SIGTERM → 等待 ≤30s 后退出；执行日志最终为 success 或 failure（不留 running）
- [ ] 7.7 前端手动点击清单（交给用户）：CRUD、启停、立即执行、深色模式、403 提示、参数表单根据 task_code 切换
- [x] 7.8 跑 `openspec status --change add-cron-task` 确认 4 个 artifact 完成
- [ ] 7.9 实施完成后由用户运行 `/openspec-archive-change` 归档
