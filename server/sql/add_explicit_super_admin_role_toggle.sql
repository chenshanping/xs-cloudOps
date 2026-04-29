-- 为角色增加显式超管开关
-- 执行方法: mysql -u root -p go-base < server/sql/add_explicit_super_admin_role_toggle.sql

SET NAMES utf8mb4;

SELECT COUNT(*) INTO @has_sys_role_is_super_admin
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_role'
  AND COLUMN_NAME = 'is_super_admin';

SET @sql = IF(
    @has_sys_role_is_super_admin = 0,
    'ALTER TABLE `sys_role` ADD COLUMN `is_super_admin` TINYINT(1) NOT NULL DEFAULT 0 COMMENT ''是否显式超管 1是 0否'' AFTER `status`',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

UPDATE `sys_role`
SET `is_super_admin` = 1
WHERE `code` IN ('admin', 'system_admin')
  AND (`is_super_admin` IS NULL OR `is_super_admin` = 0);

SELECT '显式超管角色开关升级完成，请重新登录刷新权限缓存。' AS message;
