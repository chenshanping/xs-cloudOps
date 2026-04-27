-- 移除历史 sys_storage / storage_id / 文件级 storage_config 依赖
-- 执行方法: mysql -u root -p go-base < server/sql/remove_legacy_storage_dependencies.sql

SET NAMES utf8mb4;

SET @default_local_storage_config = '{"base_path":"uploads","base_url":"/api/v1/upload"}';
SET @default_empty_storage_config = '{}';

SELECT COUNT(*) INTO @has_sys_storage
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_storage';

SELECT COUNT(*) INTO @has_sys_file_storage_type
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_file'
  AND COLUMN_NAME = 'storage_type';

SET @sql = IF(
    @has_sys_file_storage_type = 0,
    'ALTER TABLE `sys_file` ADD COLUMN `storage_type` varchar(20) NULL DEFAULT NULL COMMENT ''存储类型快照'' AFTER `md5`',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT COUNT(*) INTO @has_sys_file_chunk_storage_type
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_file_chunk'
  AND COLUMN_NAME = 'storage_type';

SET @sql = IF(
    @has_sys_file_chunk_storage_type = 0,
    'ALTER TABLE `sys_file_chunk` ADD COLUMN `storage_type` varchar(20) NULL DEFAULT NULL COMMENT ''存储类型快照'' AFTER `chunk_hash`',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT COUNT(*) INTO @has_sys_file_storage_config
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_file'
  AND COLUMN_NAME = 'storage_config';

SELECT COUNT(*) INTO @has_sys_file_chunk_storage_config
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_file_chunk'
  AND COLUMN_NAME = 'storage_config';

SELECT COUNT(*) INTO @has_sys_file_storage_id
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_file'
  AND COLUMN_NAME = 'storage_id';

SELECT COUNT(*) INTO @has_sys_file_chunk_storage_id
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_file_chunk'
  AND COLUMN_NAME = 'storage_id';

SET @legacy_default_storage_type = NULL;
SET @legacy_system_storage_config = (
    SELECT `value`
    FROM `sys_config`
    WHERE `key` = 'storage_config'
    ORDER BY `id` ASC
    LIMIT 1
);
SET @legacy_local_storage_config = NULL;
SET @legacy_aliyun_storage_config = NULL;
SET @legacy_tencent_storage_config = NULL;
SET @legacy_minio_storage_config = NULL;

SET @sql = IF(
    @has_sys_storage > 0,
    'SET @legacy_default_storage_type = COALESCE(
        (SELECT `type` FROM `sys_storage` WHERE `is_default` = 1 AND `status` = 1 ORDER BY `id` ASC LIMIT 1),
        (SELECT `type` FROM `sys_storage` WHERE `status` = 1 ORDER BY `id` ASC LIMIT 1)
    )',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(
    @has_sys_storage > 0,
    'SET @legacy_local_storage_config = (
        SELECT `config`
        FROM `sys_storage`
        WHERE `type` = ''local'' AND `status` = 1
        ORDER BY `is_default` DESC, `id` ASC
        LIMIT 1
    )',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(
    @has_sys_storage > 0,
    'SET @legacy_aliyun_storage_config = (
        SELECT `config`
        FROM `sys_storage`
        WHERE `type` = ''aliyun'' AND `status` = 1
        ORDER BY `is_default` DESC, `id` ASC
        LIMIT 1
    )',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(
    @has_sys_storage > 0,
    'SET @legacy_tencent_storage_config = (
        SELECT `config`
        FROM `sys_storage`
        WHERE `type` = ''tencent'' AND `status` = 1
        ORDER BY `is_default` DESC, `id` ASC
        LIMIT 1
    )',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(
    @has_sys_storage > 0,
    'SET @legacy_minio_storage_config = (
        SELECT `config`
        FROM `sys_storage`
        WHERE `type` = ''minio'' AND `status` = 1
        ORDER BY `is_default` DESC, `id` ASC
        LIMIT 1
    )',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @default_storage_type = COALESCE(
    (SELECT `value` FROM `sys_config` WHERE `key` = 'storage_type' ORDER BY `id` ASC LIMIT 1),
    @legacy_default_storage_type,
    'local'
);

INSERT INTO `sys_config` (`name`, `key`, `value`, `value_type`, `remark`, `created_at`, `updated_at`)
SELECT '存储类型', 'storage_type', @default_storage_type, 'string', '当前文件上传使用的存储类型', NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_config` WHERE `key` = 'storage_type'
);

UPDATE `sys_config`
SET `value` = @default_storage_type,
    `updated_at` = NOW(3)
WHERE `key` = 'storage_type'
  AND (`value` IS NULL OR `value` = '');

INSERT INTO `sys_config` (`name`, `key`, `value`, `value_type`, `remark`, `created_at`, `updated_at`)
SELECT '本地存储配置', 'storage_local_config',
       COALESCE(
           CASE WHEN @default_storage_type = 'local' THEN NULLIF(@legacy_system_storage_config, '') END,
           NULLIF(@legacy_local_storage_config, ''),
           @default_local_storage_config
       ),
       'json', '本地存储的已保存配置(JSON)', NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_config` WHERE `key` = 'storage_local_config'
);

INSERT INTO `sys_config` (`name`, `key`, `value`, `value_type`, `remark`, `created_at`, `updated_at`)
SELECT '阿里云 OSS 配置', 'storage_aliyun_config',
       COALESCE(
           CASE WHEN @default_storage_type = 'aliyun' THEN NULLIF(@legacy_system_storage_config, '') END,
           NULLIF(@legacy_aliyun_storage_config, ''),
           @default_empty_storage_config
       ),
       'json', '阿里云 OSS 的已保存配置(JSON)', NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_config` WHERE `key` = 'storage_aliyun_config'
);

INSERT INTO `sys_config` (`name`, `key`, `value`, `value_type`, `remark`, `created_at`, `updated_at`)
SELECT '腾讯云 COS 配置', 'storage_tencent_config',
       COALESCE(
           CASE WHEN @default_storage_type = 'tencent' THEN NULLIF(@legacy_system_storage_config, '') END,
           NULLIF(@legacy_tencent_storage_config, ''),
           @default_empty_storage_config
       ),
       'json', '腾讯云 COS 的已保存配置(JSON)', NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_config` WHERE `key` = 'storage_tencent_config'
);

INSERT INTO `sys_config` (`name`, `key`, `value`, `value_type`, `remark`, `created_at`, `updated_at`)
SELECT 'MinIO 配置', 'storage_minio_config',
       COALESCE(
           CASE WHEN @default_storage_type = 'minio' THEN NULLIF(@legacy_system_storage_config, '') END,
           NULLIF(@legacy_minio_storage_config, ''),
           @default_empty_storage_config
       ),
       'json', 'MinIO 的已保存配置(JSON)', NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_config` WHERE `key` = 'storage_minio_config'
);

UPDATE `sys_config`
SET `value` = COALESCE(
        CASE WHEN @default_storage_type = 'local' THEN NULLIF(@legacy_system_storage_config, '') END,
        NULLIF(@legacy_local_storage_config, ''),
        @default_local_storage_config
    ),
    `updated_at` = NOW(3)
WHERE `key` = 'storage_local_config'
  AND (`value` IS NULL OR `value` = '');

UPDATE `sys_config`
SET `value` = COALESCE(
        CASE WHEN @default_storage_type = 'aliyun' THEN NULLIF(@legacy_system_storage_config, '') END,
        NULLIF(@legacy_aliyun_storage_config, ''),
        @default_empty_storage_config
    ),
    `updated_at` = NOW(3)
WHERE `key` = 'storage_aliyun_config'
  AND (`value` IS NULL OR `value` = '');

UPDATE `sys_config`
SET `value` = COALESCE(
        CASE WHEN @default_storage_type = 'tencent' THEN NULLIF(@legacy_system_storage_config, '') END,
        NULLIF(@legacy_tencent_storage_config, ''),
        @default_empty_storage_config
    ),
    `updated_at` = NOW(3)
WHERE `key` = 'storage_tencent_config'
  AND (`value` IS NULL OR `value` = '');

UPDATE `sys_config`
SET `value` = COALESCE(
        CASE WHEN @default_storage_type = 'minio' THEN NULLIF(@legacy_system_storage_config, '') END,
        NULLIF(@legacy_minio_storage_config, ''),
        @default_empty_storage_config
    ),
    `updated_at` = NOW(3)
WHERE `key` = 'storage_minio_config'
  AND (`value` IS NULL OR `value` = '');

SET @sql = IF(
    @has_sys_storage > 0 AND @has_sys_file_storage_id > 0,
    'UPDATE `sys_file` AS f
      JOIN `sys_storage` AS s ON s.`id` = f.`storage_id`
      SET f.`storage_type` = s.`type`
      WHERE f.`storage_type` IS NULL OR f.`storage_type` = ''''',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(
    @has_sys_storage > 0 AND @has_sys_file_chunk_storage_id > 0,
    'UPDATE `sys_file_chunk` AS c
      JOIN `sys_storage` AS s ON s.`id` = c.`storage_id`
      SET c.`storage_type` = s.`type`
      WHERE c.`storage_type` IS NULL OR c.`storage_type` = ''''',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

UPDATE `sys_file`
SET `storage_type` = @default_storage_type
WHERE `storage_type` IS NULL
   OR `storage_type` = '';

UPDATE `sys_file_chunk`
SET `storage_type` = @default_storage_type
WHERE `storage_type` IS NULL
   OR `storage_type` = '';

SET @sys_file_storage_fk = (
    SELECT `CONSTRAINT_NAME`
    FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'sys_file'
      AND COLUMN_NAME = 'storage_id'
      AND REFERENCED_TABLE_NAME IS NOT NULL
    ORDER BY `CONSTRAINT_NAME` ASC
    LIMIT 1
);

SET @sys_file_chunk_storage_fk = (
    SELECT `CONSTRAINT_NAME`
    FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'sys_file_chunk'
      AND COLUMN_NAME = 'storage_id'
      AND REFERENCED_TABLE_NAME IS NOT NULL
    ORDER BY `CONSTRAINT_NAME` ASC
    LIMIT 1
);

SET @sql = CASE
    WHEN @has_sys_file_storage_id > 0 AND @has_sys_file_storage_config > 0 AND @sys_file_storage_fk IS NOT NULL
        THEN CONCAT('ALTER TABLE `sys_file` DROP FOREIGN KEY `', @sys_file_storage_fk, '`, DROP COLUMN `storage_id`, DROP COLUMN `storage_config`')
    WHEN @has_sys_file_storage_id > 0 AND @has_sys_file_storage_config > 0
        THEN 'ALTER TABLE `sys_file` DROP COLUMN `storage_id`, DROP COLUMN `storage_config`'
    WHEN @has_sys_file_storage_id > 0 AND @sys_file_storage_fk IS NOT NULL
        THEN CONCAT('ALTER TABLE `sys_file` DROP FOREIGN KEY `', @sys_file_storage_fk, '`, DROP COLUMN `storage_id`')
    WHEN @has_sys_file_storage_id > 0
        THEN 'ALTER TABLE `sys_file` DROP COLUMN `storage_id`'
    WHEN @has_sys_file_storage_config > 0
        THEN 'ALTER TABLE `sys_file` DROP COLUMN `storage_config`'
    ELSE 'SELECT 1'
END;
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = CASE
    WHEN @has_sys_file_chunk_storage_id > 0 AND @has_sys_file_chunk_storage_config > 0 AND @sys_file_chunk_storage_fk IS NOT NULL
        THEN CONCAT('ALTER TABLE `sys_file_chunk` DROP FOREIGN KEY `', @sys_file_chunk_storage_fk, '`, DROP COLUMN `storage_id`, DROP COLUMN `storage_config`')
    WHEN @has_sys_file_chunk_storage_id > 0 AND @has_sys_file_chunk_storage_config > 0
        THEN 'ALTER TABLE `sys_file_chunk` DROP COLUMN `storage_id`, DROP COLUMN `storage_config`'
    WHEN @has_sys_file_chunk_storage_id > 0 AND @sys_file_chunk_storage_fk IS NOT NULL
        THEN CONCAT('ALTER TABLE `sys_file_chunk` DROP FOREIGN KEY `', @sys_file_chunk_storage_fk, '`, DROP COLUMN `storage_id`')
    WHEN @has_sys_file_chunk_storage_id > 0
        THEN 'ALTER TABLE `sys_file_chunk` DROP COLUMN `storage_id`'
    WHEN @has_sys_file_chunk_storage_config > 0
        THEN 'ALTER TABLE `sys_file_chunk` DROP COLUMN `storage_config`'
    ELSE 'SELECT 1'
END;
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(
    @has_sys_storage > 0,
    'DROP TABLE `sys_storage`',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT '历史存储依赖清理完成：已切换为 storage_type + 分类型系统配置，sys_storage / storage_id / 文件级 storage_config 已移除。' AS message;
