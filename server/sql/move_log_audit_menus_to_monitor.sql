-- 日志审计菜单迁移至运维监控目录
-- 执行方法: mysql -u root -p go-base < server/sql/move_log_audit_menus_to_monitor.sql

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

SELECT @operation_log_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:operation-log:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @monitor_root_id, '操作日志', '/monitor/operation-log', 'monitor/operation-log/index', 'file-text', 2, 2, 'monitor:operation-log:list', 1, 0, NOW(3), NOW(3)
WHERE @monitor_root_id IS NOT NULL
  AND @operation_log_menu_id IS NULL;

SELECT @login_log_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:login-log:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @monitor_root_id, '登录日志', '/monitor/login-log', 'monitor/login-log/index', 'login', 3, 2, 'monitor:login-log:list', 1, 0, NOW(3), NOW(3)
WHERE @monitor_root_id IS NOT NULL
  AND @login_log_menu_id IS NULL;

SELECT @operation_log_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:operation-log:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

SELECT @login_log_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:login-log:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

UPDATE `sys_menu`
SET `parent_id` = @monitor_root_id,
    `updated_at` = NOW(3)
WHERE @monitor_root_id IS NOT NULL
  AND `deleted_at` IS NULL
  AND `id` IN (@operation_log_menu_id, @login_log_menu_id)
  AND `parent_id` <> @monitor_root_id;

INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT sr.`id`, menu_item.`id`
FROM `sys_role` sr
JOIN `sys_menu` menu_item ON menu_item.`deleted_at` IS NULL
WHERE sr.`code` = 'admin'
  AND sr.`deleted_at` IS NULL
  AND menu_item.`id` IN (@monitor_root_id, @operation_log_menu_id, @login_log_menu_id)
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` existing
      WHERE existing.`sys_role_id` = sr.`id`
        AND existing.`sys_menu_id` = menu_item.`id`
  );

INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT existing_menu.`sys_role_id`, @monitor_root_id
FROM `sys_role_menu` existing_menu
JOIN `sys_menu` log_menu ON log_menu.`id` = existing_menu.`sys_menu_id` AND log_menu.`deleted_at` IS NULL
WHERE @monitor_root_id IS NOT NULL
  AND log_menu.`permission` IN ('monitor:operation-log:list', 'monitor:login-log:list')
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = existing_menu.`sys_role_id`
        AND target.`sys_menu_id` = @monitor_root_id
  );

SELECT @legacy_audit_menu_id := `id`
FROM `sys_menu`
WHERE `path` = '/system/operation-audit' AND `type` = 1 AND `deleted_at` IS NULL
LIMIT 1;

UPDATE `sys_menu` audit_menu
LEFT JOIN `sys_menu` child_menu ON child_menu.`parent_id` = audit_menu.`id` AND child_menu.`deleted_at` IS NULL
SET audit_menu.`hidden` = 1,
    audit_menu.`updated_at` = NOW(3)
WHERE audit_menu.`id` = @legacy_audit_menu_id
  AND audit_menu.`hidden` <> 1
  AND child_menu.`id` IS NULL;
