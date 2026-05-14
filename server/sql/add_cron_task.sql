-- 定时任务计划：新增任务表、执行日志表、菜单权限、API 绑定与内置禁用任务
-- 执行方法: mysql -u root -p go-base < server/sql/add_cron_task.sql

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `sys_cron_task` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `code` varchar(64) NOT NULL COMMENT '任务实例编码',
  `task_code` varchar(64) NOT NULL COMMENT '注册任务编码',
  `name` varchar(128) NOT NULL COMMENT '任务名称',
  `cron_expr` varchar(64) NOT NULL COMMENT 'Cron表达式',
  `params` json NULL COMMENT '任务参数',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态 0禁用 1启用',
  `last_run_at` datetime(3) NULL DEFAULT NULL COMMENT '上次执行时间',
  `last_status` varchar(16) NULL DEFAULT NULL COMMENT '上次执行状态',
  `last_duration_ms` bigint NULL DEFAULT 0 COMMENT '上次执行耗时毫秒',
  `next_run_at` datetime(3) NULL DEFAULT NULL COMMENT '下次执行时间',
  `remark` varchar(255) NULL DEFAULT NULL COMMENT '备注',
  `sort` bigint NOT NULL DEFAULT 0 COMMENT '排序',
  `created_by` bigint UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_cron_task_deleted_at` (`deleted_at`) USING BTREE,
  INDEX `idx_sys_cron_task_code` (`code`) USING BTREE,
  INDEX `idx_sys_cron_task_task_code` (`task_code`) USING BTREE,
  INDEX `idx_sys_cron_task_status` (`status`) USING BTREE
) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `sys_cron_log` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_id` bigint UNSIGNED NOT NULL COMMENT '任务ID',
  `task_code` varchar(64) NOT NULL COMMENT '注册任务编码',
  `started_at` datetime(3) NOT NULL COMMENT '开始时间',
  `finished_at` datetime(3) NULL DEFAULT NULL COMMENT '结束时间',
  `duration_ms` bigint NULL DEFAULT 0 COMMENT '耗时毫秒',
  `status` varchar(16) NOT NULL COMMENT '执行状态',
  `summary` text NULL COMMENT '执行摘要',
  `error_message` text NULL COMMENT '错误信息',
  `triggered_by` varchar(32) NOT NULL COMMENT '触发方式',
  `actor_user_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '手动触发用户ID',
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_cron_log_task_id` (`task_id`) USING BTREE,
  INDEX `idx_sys_cron_log_started_at` (`started_at`) USING BTREE,
  INDEX `idx_sys_cron_log_status` (`status`) USING BTREE
) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

SELECT @monitor_root_id := `id`
FROM `sys_menu`
WHERE `path` = '/monitor' AND `type` = 1 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT 0, '运维监控', '/monitor', 'Layout', 'MonitorOutlined', 30, 1, '', 1, 0, NOW(3), NOW(3)
WHERE @monitor_root_id IS NULL;

SELECT @monitor_root_id := `id`
FROM `sys_menu`
WHERE `path` = '/monitor' AND `type` = 1 AND `deleted_at` IS NULL
LIMIT 1;

SELECT @cron_task_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:cron:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @monitor_root_id, '定时任务', '/monitor/cron-task', 'monitor/cron-task/index', 'ScheduleOutlined', 4, 2, 'monitor:cron:list', 1, 0, NOW(3), NOW(3)
WHERE @monitor_root_id IS NOT NULL
  AND @cron_task_menu_id IS NULL;

SELECT @cron_task_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:cron:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

UPDATE `sys_menu`
SET `parent_id` = @monitor_root_id,
    `updated_at` = NOW(3)
WHERE `id` = @cron_task_menu_id
  AND @monitor_root_id IS NOT NULL
  AND `parent_id` <> @monitor_root_id;

SELECT @cron_log_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:cron:logs:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @monitor_root_id, '任务执行日志', '/monitor/cron-log', 'monitor/cron-log/index', 'HistoryOutlined', 5, 2, 'monitor:cron:logs:list', 1, 0, NOW(3), NOW(3)
WHERE @monitor_root_id IS NOT NULL
  AND @cron_log_menu_id IS NULL;

SELECT @cron_log_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:cron:logs:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

UPDATE `sys_menu`
SET `parent_id` = @monitor_root_id,
    `updated_at` = NOW(3)
WHERE `id` = @cron_log_menu_id
  AND @monitor_root_id IS NOT NULL
  AND `parent_id` <> @monitor_root_id;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT t.`parent_id`, t.`name`, '', '', '', t.`sort`, 3, t.`permission`, 1, 0, NOW(3), NOW(3)
FROM (
    SELECT @cron_task_menu_id AS `parent_id`, '查看' AS `name`, 1 AS `sort`, 'monitor:cron:view' AS `permission`
    UNION ALL SELECT @cron_task_menu_id, '新增', 2, 'monitor:cron:create'
    UNION ALL SELECT @cron_task_menu_id, '编辑', 3, 'monitor:cron:update'
    UNION ALL SELECT @cron_task_menu_id, '删除', 4, 'monitor:cron:delete'
    UNION ALL SELECT @cron_task_menu_id, '启用', 5, 'monitor:cron:enable'
    UNION ALL SELECT @cron_task_menu_id, '停用', 6, 'monitor:cron:disable'
    UNION ALL SELECT @cron_task_menu_id, '立即执行', 7, 'monitor:cron:runNow'
    UNION ALL SELECT @cron_log_menu_id, '查看日志', 1, 'monitor:cron:logs:view'
) t
WHERE t.`parent_id` IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` sm
      WHERE sm.`permission` = t.`permission` AND sm.`type` = 3 AND sm.`deleted_at` IS NULL
  );

UPDATE `sys_menu` button_menu
JOIN (
    SELECT @cron_task_menu_id AS `parent_id`, 'monitor:cron:view' AS `permission`
    UNION ALL SELECT @cron_task_menu_id, 'monitor:cron:create'
    UNION ALL SELECT @cron_task_menu_id, 'monitor:cron:update'
    UNION ALL SELECT @cron_task_menu_id, 'monitor:cron:delete'
    UNION ALL SELECT @cron_task_menu_id, 'monitor:cron:enable'
    UNION ALL SELECT @cron_task_menu_id, 'monitor:cron:disable'
    UNION ALL SELECT @cron_task_menu_id, 'monitor:cron:runNow'
    UNION ALL SELECT @cron_log_menu_id, 'monitor:cron:logs:view'
) t ON t.`permission` = button_menu.`permission`
SET button_menu.`parent_id` = t.`parent_id`,
    button_menu.`updated_at` = NOW(3)
WHERE t.`parent_id` IS NOT NULL
  AND button_menu.`deleted_at` IS NULL
  AND button_menu.`type` = 3
  AND button_menu.`parent_id` <> t.`parent_id`;

INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT t.`path`, t.`method`, '定时任务', t.`description`, '', '', 1, NOW(3), NOW(3)
FROM (
    SELECT '/api/v1/monitor/cron-task' AS `path`, 'GET' AS `method`, '定时任务列表' AS `description`
    UNION ALL SELECT '/api/v1/monitor/cron-task', 'POST', '创建定时任务'
    UNION ALL SELECT '/api/v1/monitor/cron-task/:id', 'PUT', '更新定时任务'
    UNION ALL SELECT '/api/v1/monitor/cron-task/:id', 'DELETE', '删除定时任务'
    UNION ALL SELECT '/api/v1/monitor/cron-task/:id/enable', 'POST', '启用定时任务'
    UNION ALL SELECT '/api/v1/monitor/cron-task/:id/disable', 'POST', '停用定时任务'
    UNION ALL SELECT '/api/v1/monitor/cron-task/:id/run', 'POST', '立即执行定时任务'
    UNION ALL SELECT '/api/v1/monitor/cron-task/registry', 'GET', '定时任务注册列表'
    UNION ALL SELECT '/api/v1/monitor/cron-log', 'GET', '定时任务执行日志'
    UNION ALL SELECT '/api/v1/monitor/cron-log/:id', 'GET', '定时任务执行日志详情'
) t
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` sa
    WHERE sa.`path` = t.`path` AND sa.`method` = t.`method` AND sa.`deleted_at` IS NULL
);

INSERT INTO `sys_menu_api` (`sys_menu_id`, `sys_api_id`)
SELECT menu_item.`id`, api_item.`id`
FROM (
    SELECT 'monitor:cron:view' AS `menu_permission`, '/api/v1/monitor/cron-task' AS `api_path`, 'GET' AS `api_method`
    UNION ALL SELECT 'monitor:cron:view', '/api/v1/monitor/cron-task/registry', 'GET'
    UNION ALL SELECT 'monitor:cron:create', '/api/v1/monitor/cron-task', 'POST'
    UNION ALL SELECT 'monitor:cron:update', '/api/v1/monitor/cron-task/:id', 'PUT'
    UNION ALL SELECT 'monitor:cron:delete', '/api/v1/monitor/cron-task/:id', 'DELETE'
    UNION ALL SELECT 'monitor:cron:enable', '/api/v1/monitor/cron-task/:id/enable', 'POST'
    UNION ALL SELECT 'monitor:cron:disable', '/api/v1/monitor/cron-task/:id/disable', 'POST'
    UNION ALL SELECT 'monitor:cron:runNow', '/api/v1/monitor/cron-task/:id/run', 'POST'
    UNION ALL SELECT 'monitor:cron:logs:view', '/api/v1/monitor/cron-log', 'GET'
    UNION ALL SELECT 'monitor:cron:logs:view', '/api/v1/monitor/cron-log/:id', 'GET'
) t
JOIN `sys_menu` menu_item ON menu_item.`permission` = t.`menu_permission` AND menu_item.`deleted_at` IS NULL
JOIN `sys_api` api_item ON api_item.`path` = t.`api_path` AND api_item.`method` = t.`api_method` AND api_item.`deleted_at` IS NULL
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_menu_api` existing
    WHERE existing.`sys_menu_id` = menu_item.`id`
      AND existing.`sys_api_id` = api_item.`id`
);

INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT sr.`id`, menu_item.`id`
FROM `sys_role` sr
JOIN `sys_menu` menu_item ON menu_item.`deleted_at` IS NULL
WHERE sr.`code` = 'admin'
  AND sr.`deleted_at` IS NULL
  AND (
      menu_item.`id` = @monitor_root_id
      OR menu_item.`permission` IN (
          'monitor:cron:list',
          'monitor:cron:logs:list',
          'monitor:cron:view',
          'monitor:cron:create',
          'monitor:cron:update',
          'monitor:cron:delete',
          'monitor:cron:enable',
          'monitor:cron:disable',
          'monitor:cron:runNow',
          'monitor:cron:logs:view'
      )
  )
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = sr.`id`
        AND target.`sys_menu_id` = menu_item.`id`
  );

INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
SELECT DISTINCT 'p', sr.`code`, api_item.`path`, api_item.`method`, '', '', ''
FROM `sys_role` sr
JOIN `sys_menu` menu_item ON menu_item.`permission` IN (
    'monitor:cron:view',
    'monitor:cron:create',
    'monitor:cron:update',
    'monitor:cron:delete',
    'monitor:cron:enable',
    'monitor:cron:disable',
    'monitor:cron:runNow',
    'monitor:cron:logs:view'
) AND menu_item.`deleted_at` IS NULL
JOIN `sys_role_menu` srm ON srm.`sys_role_id` = sr.`id` AND srm.`sys_menu_id` = menu_item.`id`
JOIN `sys_menu_api` sma ON sma.`sys_menu_id` = menu_item.`id`
JOIN `sys_api` api_item ON api_item.`id` = sma.`sys_api_id` AND api_item.`deleted_at` IS NULL
WHERE sr.`deleted_at` IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM `casbin_rule` target
      WHERE target.`ptype` = 'p'
        AND target.`v0` = sr.`code`
        AND target.`v1` = api_item.`path`
        AND target.`v2` = api_item.`method`
        AND IFNULL(target.`v3`, '') = ''
        AND IFNULL(target.`v4`, '') = ''
        AND IFNULL(target.`v5`, '') = ''
  );

INSERT INTO `sys_cron_task` (`code`, `task_code`, `name`, `cron_expr`, `params`, `status`, `remark`, `sort`, `created_at`, `updated_at`)
SELECT 'cleanup_login_logs_default', 'cleanup_login_logs', '清理登录日志', '0 2 * * *', '{"batch_limit":1000,"retain_days":30}', 0, '内置任务，默认停用', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_cron_task` task
    WHERE task.`code` = 'cleanup_login_logs_default' AND task.`deleted_at` IS NULL
);

INSERT INTO `sys_cron_task` (`code`, `task_code`, `name`, `cron_expr`, `params`, `status`, `remark`, `sort`, `created_at`, `updated_at`)
SELECT 'cleanup_operation_logs_default', 'cleanup_operation_logs', '清理操作日志', '0 2 * * *', '{"batch_limit":1000,"retain_days":30}', 0, '内置任务，默认停用', 2, NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_cron_task` task
    WHERE task.`code` = 'cleanup_operation_logs_default' AND task.`deleted_at` IS NULL
);

SELECT '定时任务计划表、菜单、按钮权限、API 元数据、菜单 API 绑定与内置停用任务已补齐；请重启服务或重载 Casbin 策略，并重新登录刷新权限缓存。' AS message;
