-- 为 sys_file 表添加 storage_bucket 列，记录文件上传时的桶名/路径快照
-- 执行方法: mysql -u root -p go-base < server/sql/add_file_storage_bucket_column.sql

SET NAMES utf8mb4;

-- 添加 storage_bucket 列（如果不存在）
SELECT COUNT(*) INTO @col_exists
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_file'
  AND COLUMN_NAME = 'storage_bucket';

SET @sql = IF(@col_exists = 0,
    'ALTER TABLE `sys_file` ADD COLUMN `storage_bucket` varchar(100) DEFAULT NULL COMMENT ''存储桶/路径快照'' AFTER `storage_type`',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
