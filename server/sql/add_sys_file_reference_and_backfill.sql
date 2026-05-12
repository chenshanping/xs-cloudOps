-- 创建 sys_file_reference 表，并回填用户头像、系统配置图片、AI 对话附件的文件引用关系
-- 执行方法: mysql -u root -p go-base < server/sql/add_sys_file_reference_and_backfill.sql

SET NAMES utf8mb4;

SELECT COUNT(*) INTO @tbl_exists
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_file_reference';

SET @sql = IF(@tbl_exists = 0,
    'CREATE TABLE `sys_file_reference` (
        `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) NULL DEFAULT NULL,
        `updated_at` datetime(3) NULL DEFAULT NULL,
        `file_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT ''文件ID'',
        `ref_table` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT ''引用表名'',
        `ref_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT ''引用记录ID'',
        `ref_field` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT ''引用字段'',
        PRIMARY KEY (`id`) USING BTREE,
        UNIQUE INDEX `uk_sys_file_reference`(`ref_field` ASC, `ref_table` ASC, `ref_id` ASC, `file_id` ASC) USING BTREE,
        INDEX `idx_sys_file_reference_file_id`(`file_id` ASC) USING BTREE,
        INDEX `idx_ref_table_id_field`(`ref_table` ASC, `ref_id` ASC, `ref_field` ASC) USING BTREE
    ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

INSERT INTO `sys_config` (`created_at`, `updated_at`, `name`, `key`, `value`, `value_type`, `remark`)
SELECT NOW(3), NOW(3), '系统Logo文件ID', 'sys_logo_file_id', '', 'string', '系统Logo关联文件ID'
WHERE NOT EXISTS (SELECT 1 FROM `sys_config` WHERE `key` = 'sys_logo_file_id');

INSERT INTO `sys_config` (`created_at`, `updated_at`, `name`, `key`, `value`, `value_type`, `remark`)
SELECT NOW(3), NOW(3), '注册默认头像文件ID', 'register_logo_file_id', '', 'string', '注册默认头像关联文件ID'
WHERE NOT EXISTS (SELECT 1 FROM `sys_config` WHERE `key` = 'register_logo_file_id');

INSERT INTO `sys_config` (`created_at`, `updated_at`, `name`, `key`, `value`, `value_type`, `remark`)
SELECT NOW(3), NOW(3), '登录页背景图文件ID', 'login_bg_image_file_id', '', 'string', '登录页背景图关联文件ID'
WHERE NOT EXISTS (SELECT 1 FROM `sys_config` WHERE `key` = 'login_bg_image_file_id');

INSERT INTO `sys_config` (`created_at`, `updated_at`, `name`, `key`, `value`, `value_type`, `remark`)
SELECT NOW(3), NOW(3), '滑动验证码背景文件ID', 'slider_captcha_bg_file_id', '', 'string', '滑动验证码背景关联文件ID'
WHERE NOT EXISTS (SELECT 1 FROM `sys_config` WHERE `key` = 'slider_captcha_bg_file_id');

UPDATE `sys_config` AS target
JOIN `sys_config` AS source ON source.`key` = 'sys_logo'
JOIN `sys_file` AS file ON file.`url` = source.`value` AND file.`status` = 1
SET target.`value` = CAST(file.`id` AS CHAR), target.`updated_at` = NOW(3)
WHERE target.`key` = 'sys_logo_file_id'
  AND (target.`value` IS NULL OR target.`value` = '' OR target.`value` = '0');

UPDATE `sys_config` AS target
JOIN `sys_config` AS source ON source.`key` = 'register_logo'
JOIN `sys_file` AS file ON file.`url` = source.`value` AND file.`status` = 1
SET target.`value` = CAST(file.`id` AS CHAR), target.`updated_at` = NOW(3)
WHERE target.`key` = 'register_logo_file_id'
  AND (target.`value` IS NULL OR target.`value` = '' OR target.`value` = '0');

UPDATE `sys_config` AS target
JOIN `sys_config` AS source ON source.`key` = 'login_bg_image'
JOIN `sys_file` AS file ON file.`url` = source.`value` AND file.`status` = 1
SET target.`value` = CAST(file.`id` AS CHAR), target.`updated_at` = NOW(3)
WHERE target.`key` = 'login_bg_image_file_id'
  AND (target.`value` IS NULL OR target.`value` = '' OR target.`value` = '0');

UPDATE `sys_config` AS target
JOIN `sys_config` AS source ON source.`key` = 'slider_captcha_bg'
JOIN `sys_file` AS file ON file.`url` = source.`value` AND file.`status` = 1
SET target.`value` = CAST(file.`id` AS CHAR), target.`updated_at` = NOW(3)
WHERE target.`key` = 'slider_captcha_bg_file_id'
  AND (target.`value` IS NULL OR target.`value` = '' OR target.`value` = '0');

INSERT IGNORE INTO `sys_file_reference` (`created_at`, `updated_at`, `file_id`, `ref_table`, `ref_id`, `ref_field`)
SELECT NOW(3), NOW(3), user.`avatar_file_id`, 'sys_user', user.`id`, 'avatar'
FROM `sys_user` AS user
WHERE user.`deleted_at` IS NULL
  AND user.`avatar_file_id` IS NOT NULL
  AND user.`avatar_file_id` > 0;

INSERT IGNORE INTO `sys_file_reference` (`created_at`, `updated_at`, `file_id`, `ref_table`, `ref_id`, `ref_field`)
SELECT NOW(3), NOW(3), CAST(config.`value` AS UNSIGNED), 'sys_config', config.`id`, config.`key`
FROM `sys_config` AS config
WHERE config.`deleted_at` IS NULL
  AND config.`key` IN ('sys_logo_file_id', 'register_logo_file_id', 'login_bg_image_file_id', 'slider_captcha_bg_file_id')
  AND config.`value` REGEXP '^[0-9]+$'
  AND CAST(config.`value` AS UNSIGNED) > 0;

INSERT IGNORE INTO `sys_file_reference` (`created_at`, `updated_at`, `file_id`, `ref_table`, `ref_id`, `ref_field`)
SELECT NOW(3), NOW(3), jt.`file_id`, 'ai_message', refs.`id`, 'attachment'
FROM (
    SELECT msg.`id`, msg.`file_ids`
    FROM `ai_messages` AS msg
    WHERE msg.`deleted_at` IS NULL
      AND msg.`file_ids` IS NOT NULL
      AND msg.`file_ids` <> ''
      AND JSON_VALID(msg.`file_ids`)
) AS refs
JOIN JSON_TABLE(refs.`file_ids`, '$[*]' COLUMNS (`file_id` BIGINT PATH '$')) AS jt
WHERE jt.`file_id` > 0;
