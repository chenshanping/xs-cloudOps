-- 为用户管理“仅本人数据”补充创建人字段
-- 执行方法: mysql -u root -p go-base < server/sql/add_user_created_by_for_user_management_scope.sql

SET NAMES utf8mb4;

SELECT COUNT(*) INTO @has_sys_user_created_by
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_user'
  AND COLUMN_NAME = 'created_by';

SET @sql = IF(
    @has_sys_user_created_by = 0,
    'ALTER TABLE `sys_user` ADD COLUMN `created_by` bigint UNSIGNED NULL DEFAULT 0 COMMENT ''创建人ID'' AFTER `dept_id`',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT COUNT(*) INTO @has_idx_sys_user_created_by
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_user'
  AND INDEX_NAME = 'idx_sys_user_created_by';

SET @sql = IF(
    @has_idx_sys_user_created_by = 0,
    'ALTER TABLE `sys_user` ADD INDEX `idx_sys_user_created_by` (`created_by` ASC)',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT 'sys_user.created_by 已补齐。历史用户如需参与“仅本人=我创建的用户”过滤，请按实际业务回填 created_by。' AS message;
