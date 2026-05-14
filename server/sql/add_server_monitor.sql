-- 服务监控：新增运维监控目录、服务监控菜单、按钮权限、API 元数据与菜单 API 绑定
-- 执行方法: mysql -u root -p go-base < server/sql/add_server_monitor.sql

SET NAMES utf8mb4;

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

SELECT @server_monitor_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:server:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @monitor_root_id, '服务监控', '/monitor/server', 'monitor/server/index', 'DashboardOutlined', 1, 2, 'monitor:server:list', 1, 0, NOW(3), NOW(3)
WHERE @monitor_root_id IS NOT NULL
  AND @server_monitor_menu_id IS NULL;

SELECT @server_monitor_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:server:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

UPDATE `sys_menu`
SET `parent_id` = @monitor_root_id,
    `updated_at` = NOW(3)
WHERE `id` = @server_monitor_menu_id
  AND @monitor_root_id IS NOT NULL
  AND `parent_id` <> @monitor_root_id;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @server_monitor_menu_id, t.`name`, '', '', '', t.`sort`, 3, t.`permission`, 1, 0, NOW(3), NOW(3)
FROM (
    SELECT '查看系统' AS `name`, 1 AS `sort`, 'monitor:server:view' AS `permission`
    UNION ALL SELECT '查看运行时', 2, 'monitor:runtime:view'
    UNION ALL SELECT '查看数据库', 3, 'monitor:db:view'
    UNION ALL SELECT '查看缓存', 4, 'monitor:cache:view'
    UNION ALL SELECT '清理缓存', 5, 'monitor:cache:clear'
    UNION ALL SELECT '查看 OSS', 6, 'monitor:oss:view'
    UNION ALL SELECT '健康概览', 7, 'monitor:dependency:view'
) t
WHERE @server_monitor_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` sm
      WHERE sm.`permission` = t.`permission` AND sm.`type` = 3 AND sm.`deleted_at` IS NULL
  );

UPDATE `sys_menu` button_menu
JOIN (
    SELECT 'monitor:server:view' AS `permission`
    UNION ALL SELECT 'monitor:runtime:view'
    UNION ALL SELECT 'monitor:db:view'
    UNION ALL SELECT 'monitor:cache:view'
    UNION ALL SELECT 'monitor:cache:clear'
    UNION ALL SELECT 'monitor:oss:view'
    UNION ALL SELECT 'monitor:dependency:view'
) t ON t.`permission` = button_menu.`permission`
SET button_menu.`parent_id` = @server_monitor_menu_id,
    button_menu.`updated_at` = NOW(3)
WHERE @server_monitor_menu_id IS NOT NULL
  AND button_menu.`deleted_at` IS NULL
  AND button_menu.`type` = 3
  AND button_menu.`parent_id` <> @server_monitor_menu_id;

INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT t.`path`, t.`method`, '服务监控', t.`description`, '', '', 1, NOW(3), NOW(3)
FROM (
    SELECT '/api/v1/monitor/server' AS `path`, 'GET' AS `method`, '服务器指标' AS `description`
    UNION ALL SELECT '/api/v1/monitor/runtime', 'GET', '运行时指标'
    UNION ALL SELECT '/api/v1/monitor/db', 'GET', '数据库连接池指标'
    UNION ALL SELECT '/api/v1/monitor/redis', 'GET', 'Redis 缓存指标'
    UNION ALL SELECT '/api/v1/monitor/redis/clear', 'POST', '清理 Redis 缓存'
    UNION ALL SELECT '/api/v1/monitor/oss', 'GET', '对象存储健康'
    UNION ALL SELECT '/api/v1/monitor/dependency', 'GET', '依赖健康概览'
) t
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` sa
    WHERE sa.`path` = t.`path` AND sa.`method` = t.`method` AND sa.`deleted_at` IS NULL
);

INSERT INTO `sys_menu_api` (`sys_menu_id`, `sys_api_id`)
SELECT menu_item.`id`, api_item.`id`
FROM (
    SELECT 'monitor:server:view' AS `menu_permission`, '/api/v1/monitor/server' AS `api_path`, 'GET' AS `api_method`
    UNION ALL SELECT 'monitor:runtime:view', '/api/v1/monitor/runtime', 'GET'
    UNION ALL SELECT 'monitor:db:view', '/api/v1/monitor/db', 'GET'
    UNION ALL SELECT 'monitor:cache:view', '/api/v1/monitor/redis', 'GET'
    UNION ALL SELECT 'monitor:cache:clear', '/api/v1/monitor/redis/clear', 'POST'
    UNION ALL SELECT 'monitor:oss:view', '/api/v1/monitor/oss', 'GET'
    UNION ALL SELECT 'monitor:dependency:view', '/api/v1/monitor/dependency', 'GET'
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
  AND menu_item.`permission` IN (
      'monitor:server:list',
      'monitor:server:view',
      'monitor:runtime:view',
      'monitor:db:view',
      'monitor:cache:view',
      'monitor:cache:clear',
      'monitor:oss:view',
      'monitor:dependency:view'
  )
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = sr.`id`
        AND target.`sys_menu_id` = menu_item.`id`
  );

INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT sr.`id`, menu_item.`id`
FROM `sys_role` sr
JOIN `sys_menu` menu_item ON menu_item.`id` = @monitor_root_id AND menu_item.`deleted_at` IS NULL
WHERE sr.`code` = 'admin'
  AND sr.`deleted_at` IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = sr.`id`
        AND target.`sys_menu_id` = menu_item.`id`
  );

INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
SELECT DISTINCT 'p', sr.`code`, api_item.`path`, api_item.`method`, '', '', ''
FROM `sys_role` sr
JOIN `sys_menu` menu_item ON menu_item.`permission` IN (
    'monitor:server:view',
    'monitor:runtime:view',
    'monitor:db:view',
    'monitor:cache:view',
    'monitor:cache:clear',
    'monitor:oss:view',
    'monitor:dependency:view'
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

SELECT '服务监控菜单、按钮权限、API 元数据与菜单 API 绑定已补齐；请重启服务或重载 Casbin 策略，并重新登录刷新权限缓存。' AS message;
