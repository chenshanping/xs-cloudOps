-- 定时任务：移除独立任务执行日志页面，并将日志查看入口合并到定时任务抽屉
-- 执行方法: mysql -u root -p go-base < server/sql/update_cron_task_merge_log_menu_into_drawer.sql

SET NAMES utf8mb4;

SELECT @monitor_root_id := `id`
FROM `sys_menu`
WHERE `path` = '/monitor' AND `type` = 1 AND `deleted_at` IS NULL
LIMIT 1;

SELECT @cron_task_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:cron:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

SELECT @cron_log_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'monitor:cron:logs:list' AND `type` = 2
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @cron_task_menu_id, '查看日志', '', '', '', 8, 3, 'monitor:cron:logs:view', 1, 0, NOW(3), NOW(3)
WHERE @cron_task_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu`
      WHERE `permission` = 'monitor:cron:logs:view' AND `type` = 3 AND `deleted_at` IS NULL
  );

UPDATE `sys_menu`
SET `parent_id` = @cron_task_menu_id,
    `sort` = 8,
    `updated_at` = NOW(3)
WHERE `permission` = 'monitor:cron:logs:view'
  AND `type` = 3
  AND `deleted_at` IS NULL
  AND @cron_task_menu_id IS NOT NULL
  AND (`parent_id` <> @cron_task_menu_id OR `sort` <> 8);

INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT source_rel.`sys_role_id`, @monitor_root_id
FROM `sys_role_menu` source_rel
JOIN `sys_menu` source_menu ON source_menu.`id` = source_rel.`sys_menu_id`
WHERE source_menu.`permission` = 'monitor:cron:logs:list'
  AND source_menu.`type` = 2
  AND @monitor_root_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = source_rel.`sys_role_id`
        AND target.`sys_menu_id` = @monitor_root_id
  );

INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT source_rel.`sys_role_id`, @cron_task_menu_id
FROM `sys_role_menu` source_rel
JOIN `sys_menu` source_menu ON source_menu.`id` = source_rel.`sys_menu_id`
WHERE source_menu.`permission` = 'monitor:cron:logs:list'
  AND source_menu.`type` = 2
  AND @cron_task_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = source_rel.`sys_role_id`
        AND target.`sys_menu_id` = @cron_task_menu_id
  );

DELETE target
FROM `sys_role_menu` target
JOIN `sys_menu` source_menu ON source_menu.`id` = target.`sys_menu_id`
WHERE source_menu.`permission` = 'monitor:cron:logs:list'
  AND source_menu.`type` = 2;

UPDATE `sys_menu`
SET `deleted_at` = NOW(3),
    `updated_at` = NOW(3)
WHERE `permission` = 'monitor:cron:logs:list'
  AND `type` = 2
  AND `deleted_at` IS NULL;

SELECT '定时任务执行日志已从独立菜单迁移到定时任务抽屉；请重载权限并重新登录刷新菜单缓存。' AS message;
