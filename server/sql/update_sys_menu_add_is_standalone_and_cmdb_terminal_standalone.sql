-- 为系统菜单补齐独立页字段，并将 CMDB 终端页切换为独立页面
-- 执行方法: mysql -u root -p go-base < server/sql/update_sys_menu_add_is_standalone_and_cmdb_terminal_standalone.sql

SET NAMES utf8mb4;

SELECT COUNT(*) INTO @sys_menu_is_standalone_exists
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_menu'
  AND COLUMN_NAME = 'is_standalone';

SET @sql = IF(
  @sys_menu_is_standalone_exists = 0,
  'ALTER TABLE `sys_menu` ADD COLUMN `is_standalone` tinyint NOT NULL DEFAULT 0 COMMENT ''是否独立页 0否 1是'' AFTER `hidden`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

UPDATE `sys_menu`
SET `is_standalone` = 1,
    `hidden` = 1,
    `path` = '/cmdb/terminal/:hostId',
    `component` = 'cmdb/terminal/index',
    `updated_at` = NOW(3)
WHERE `permission` = 'cmdb:terminal:list'
  AND `type` = 2
  AND `deleted_at` IS NULL
  AND (
    `is_standalone` <> 1
    OR `hidden` <> 1
    OR `path` <> '/cmdb/terminal/:hostId'
    OR `component` <> 'cmdb/terminal/index'
  );

SELECT 'sys_menu 独立页字段与 CMDB 终端菜单独立页配置已补齐。' AS message;
