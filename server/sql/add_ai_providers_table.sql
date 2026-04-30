-- 新增 AI 平台配置表，并兼容迁移历史 sys_config.ai_config 数据
-- 执行方法: mysql -u root -p go-base < server/sql/add_ai_providers_table.sql

SET NAMES utf8mb4;

SELECT COUNT(*) INTO @has_ai_providers
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'ai_providers';

SET @sql = IF(
    @has_ai_providers = 0,
    'CREATE TABLE `ai_providers` (
        `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) NULL DEFAULT NULL,
        `updated_at` datetime(3) NULL DEFAULT NULL,
        `deleted_at` datetime(3) NULL DEFAULT NULL,
        `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT ''平台名称'',
        `api_key` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT ''平台API Key'',
        `base_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT ''平台Base URL'',
        `models_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT ''模型配置JSON'',
        `is_default` tinyint(1) NULL DEFAULT 0 COMMENT ''是否默认平台'',
        `sort` bigint NULL DEFAULT 0 COMMENT ''排序'',
        PRIMARY KEY (`id`) USING BTREE,
        UNIQUE INDEX `idx_ai_providers_name`(`name` ASC) USING BTREE,
        INDEX `idx_ai_providers_deleted_at`(`deleted_at` ASC) USING BTREE
    ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @legacy_sys_config_json = (
    SELECT `value`
    FROM `sys_config`
    WHERE `key` = 'ai_config'
      AND `deleted_at` IS NULL
    ORDER BY `id` ASC
    LIMIT 1
);

SET @source_ai_config_json = COALESCE(
    NULLIF(@legacy_sys_config_json, '')
);

SET @default_provider_name = JSON_UNQUOTE(JSON_EXTRACT(@source_ai_config_json, '$.default_provider'));

SET @default_provider_match_count = 0;
SET @sql = IF(
    @source_ai_config_json IS NOT NULL AND @source_ai_config_json <> '',
    'SELECT COUNT(*) INTO @default_provider_match_count
     FROM JSON_TABLE(
         @source_ai_config_json,
         ''$.providers[*]''
         COLUMNS (
             provider_name varchar(100) PATH ''$.name''
         )
     ) AS jt
     WHERE jt.provider_name = @default_provider_name',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

INSERT INTO `ai_providers` (`name`, `api_key`, `base_url`, `models_json`, `is_default`, `sort`, `created_at`, `updated_at`)
SELECT
    jt.provider_name,
    jt.api_key,
    jt.base_url,
    JSON_EXTRACT(jt.provider_json, '$.models'),
    CASE
        WHEN @default_provider_match_count > 0 AND jt.provider_name = @default_provider_name THEN 1
        WHEN @default_provider_match_count = 0 AND jt.sort_order = 1 THEN 1
        ELSE 0
    END,
    jt.sort_order - 1,
    NOW(3),
    NOW(3)
FROM JSON_TABLE(
    @source_ai_config_json,
    '$.providers[*]'
    COLUMNS (
        sort_order FOR ORDINALITY,
        provider_json JSON PATH '$',
        provider_name varchar(100) PATH '$.name',
        api_key text PATH '$.api_key',
        base_url varchar(255) PATH '$.base_url'
    )
) AS jt
WHERE @source_ai_config_json IS NOT NULL
  AND @source_ai_config_json <> ''
  AND NOT EXISTS (
      SELECT 1
      FROM `ai_providers`
      WHERE `deleted_at` IS NULL
  );

SELECT 'AI 平台配置表补齐完成：ai_providers 已存在，历史 sys_config.ai_config 已在缺失时迁移。' AS message;
