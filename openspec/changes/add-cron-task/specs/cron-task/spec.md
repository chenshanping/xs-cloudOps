## ADDED Requirements

### Requirement: 任务白名单与代码注册

The system SHALL execute ONLY tasks whose `task_code` is registered in a code-level whitelist (`server/service/cronsvc/registry.go`). The database, frontend, and any configuration source MUST NOT be able to introduce a new executable task at runtime.

#### Scenario: Creating a task with whitelisted code succeeds

- **WHEN** an administrator with `monitor:cron:create` calls `POST /admin/monitor/cron-task` with `task_code = "cleanup_login_logs"` and a valid cron expression
- **THEN** the system creates a record in `sys_cron_task` with `status = disabled` by default

#### Scenario: Creating a task with non-whitelisted code is rejected

- **WHEN** an administrator submits `task_code = "rm_rf_root"` (not in registry)
- **THEN** the system returns 400 with code `cron_task_code_not_registered` and does not insert any row

#### Scenario: Enabling a task with non-whitelisted code is rejected

- **WHEN** a `sys_cron_task` row somehow exists with an invalid `task_code` and an administrator attempts to enable it
- **THEN** the system returns 400, does not register the task with the scheduler, and does not run it

### Requirement: 任务参数 schema 校验

Each registered task SHALL declare a parameter schema. The system MUST persist only the parameter fields declared in the schema and MUST reject save attempts where a declared field has the wrong type.

#### Scenario: Valid params are persisted

- **WHEN** a task is saved with `params = {"retain_days": 30, "batch_limit": 1000}` matching the schema
- **THEN** the system persists exactly those two fields

#### Scenario: Unknown params are stripped

- **WHEN** a task is saved with `params = {"retain_days": 30, "evil": "<script>"}`
- **THEN** the system persists only `{"retain_days": 30}` and ignores `evil`

#### Scenario: Wrong type is rejected

- **WHEN** a task is saved with `params = {"retain_days": "not-a-number"}`
- **THEN** the system returns 400 and does not persist the change

### Requirement: cron 表达式校验

The system SHALL validate cron expressions using the `robfig/cron` parser at task save time and reject invalid expressions before persisting.

#### Scenario: Valid expression is accepted

- **WHEN** an administrator saves a task with `cron_expr = "0 2 * * *"`
- **THEN** the system parses the expression, computes `next_run_at`, and persists the task

#### Scenario: Invalid expression is rejected

- **WHEN** an administrator saves a task with `cron_expr = "@@@@"`
- **THEN** the system returns 400 with code `cron_expr_invalid`

### Requirement: 调度器生命周期

The system SHALL start a single in-process scheduler at server startup, load all `status = enabled` tasks, and gracefully stop the scheduler on shutdown waiting at most 30 seconds for in-flight tasks.

#### Scenario: Enabled tasks are scheduled at startup

- **WHEN** the server starts with two enabled rows in `sys_cron_task`
- **THEN** the scheduler is registered with both tasks and `next_run_at` is computed for each

#### Scenario: Disabled tasks are not scheduled

- **WHEN** the server starts with disabled rows
- **THEN** those tasks are NOT added to the scheduler

#### Scenario: Graceful shutdown waits for in-flight tasks

- **WHEN** the server receives shutdown signal while a task is running
- **THEN** the scheduler stops accepting new triggers and waits up to 30 seconds for the running task to finish before exiting

### Requirement: 任务 CRUD 与启停

The system SHALL provide REST endpoints for listing, creating, updating, deleting, enabling, and disabling cron tasks, each gated by the matching permission code.

#### Scenario: Enable adds task to running scheduler

- **WHEN** an administrator with `monitor:cron:enable` enables a task at runtime
- **THEN** the scheduler picks up the new task without server restart, and the task starts running on its cron expression

#### Scenario: Disable removes task from scheduler

- **WHEN** an administrator with `monitor:cron:disable` disables a running task
- **THEN** the scheduler removes the task entry; subsequent cron ticks do not invoke the task; if a task is currently executing it is allowed to finish

#### Scenario: Update with invalid changes does not affect schedule

- **WHEN** an administrator submits an update with an invalid cron expression
- **THEN** the system returns 400, the persisted task remains unchanged, and the scheduler entry remains unchanged

### Requirement: 同任务串行执行

The system SHALL ensure that two invocations of the same task do not execute concurrently. If the previous invocation is still running when a cron tick fires, the new invocation MUST be skipped and a log row recorded with `status = skipped`.

#### Scenario: Overlapping trigger is skipped

- **GIVEN** task A is currently running and its cron expression fires again
- **WHEN** the scheduler attempts to invoke task A
- **THEN** the new invocation is skipped, a `sys_cron_log` row is written with `status = skipped`, and the running invocation is not interrupted

### Requirement: 手动触发（runNow）

The system SHALL allow administrators with `monitor:cron:runNow` to manually trigger a task immediately. The HTTP request MUST return promptly with the new log id; task execution runs asynchronously.

#### Scenario: Manual run returns log id immediately

- **WHEN** an administrator calls `POST /admin/monitor/cron-task/:id/run`
- **THEN** the system creates a `sys_cron_log` row with `status = running`, `triggered_by = manual`, returns the log id within 200 ms, and runs the task in a separate goroutine

#### Scenario: Manual run respects same-task serial constraint

- **WHEN** an administrator triggers `runNow` while the same task is already running
- **THEN** the system returns 409 with code `cron_task_already_running` and does not start a duplicate execution

### Requirement: 执行日志记录

Every task execution attempt SHALL produce a `sys_cron_log` row containing `task_id`, `task_code`, `started_at`, `finished_at`, `duration_ms`, `status` (`success`/`failure`/`skipped`/`running`), `summary` (truncated to 4KB), `error_message` (on failure), `triggered_by` (`schedule` or `manual`), and `actor_user_id` (manual only).

#### Scenario: Successful run is logged

- **WHEN** a task completes successfully
- **THEN** a log row exists with `status = success`, `duration_ms` populated, `summary` set, `error_message` empty

#### Scenario: Failing task is logged with error

- **WHEN** a task panics or returns a non-nil error
- **THEN** the panic is recovered, a log row is written with `status = failure`, `error_message` containing the truncated stack trace or error string, and the scheduler continues to operate

### Requirement: 内置任务

The system SHALL ship two built-in registered tasks at v1:

1. `cleanup_login_logs` — deletes `sys_log_login` rows older than `retain_days` in batches of `batch_limit`
2. `cleanup_operation_logs` — deletes `sys_log_operation` rows older than `retain_days` in batches of `batch_limit`

Both tasks MUST be installed in `disabled` state by the SQL upgrade script.

#### Scenario: cleanup_login_logs removes only old rows

- **GIVEN** `sys_log_login` contains rows from yesterday and rows from 60 days ago
- **WHEN** the task runs with `retain_days = 30, batch_limit = 1000`
- **THEN** only rows older than 30 days are deleted; recent rows are untouched; the summary reports the deletion count

#### Scenario: cleanup respects ctx cancellation

- **WHEN** the task is in the middle of batched deletion and the context is cancelled
- **THEN** the task stops at the next batch boundary and returns a partial-success summary

### Requirement: 菜单与权限自动注册

On server startup the system SHALL idempotently ensure the `定时任务` and `任务执行日志` menus, the 8 button permissions, the corresponding SysApi rows with `need_auth=true`, and the `sys_menu_api` bindings. Existing custom edits MUST be preserved.

#### Scenario: First-time bootstrap

- **WHEN** the server starts on a fresh database
- **THEN** the bootstrap creates the menus, buttons, SysApi rows, bindings, and grants the 8 codes to the built-in `admin` role

#### Scenario: Subsequent startup preserves admin edits

- **WHEN** the server restarts after manual edits to a button label
- **THEN** the bootstrap does not overwrite the edited label and does not regrant permissions removed from non-`admin` roles
