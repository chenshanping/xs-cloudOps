-- 为角色增加按业务功能配置的数据权限覆盖能力
-- 执行方法: mysql -u root -p go-base < server/sql/add_role_feature_data_scopes.sql

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `sys_role_data_scope` (
    `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NULL DEFAULT NULL,
    `updated_at` datetime(3) NULL DEFAULT NULL,
    `deleted_at` datetime(3) NULL DEFAULT NULL,
    `role_id` bigint UNSIGNED NOT NULL COMMENT '角色ID',
    `resource_code` varchar(100) NULL DEFAULT NULL COMMENT '业务功能资源码',
    `data_scope` bigint NULL DEFAULT 1 COMMENT '数据范围 1全部 2自定义 3本部门 4本部门及下级 5仅本人',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `idx_role_resource` (`role_id` ASC, `resource_code` ASC) USING BTREE,
    INDEX `idx_sys_role_data_scope_deleted_at` (`deleted_at` ASC) USING BTREE,
    CONSTRAINT `fk_sys_role_data_scope_role` FOREIGN KEY (`role_id`) REFERENCES `sys_role` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SELECT COUNT(*) INTO @has_idx_role_resource
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_role_data_scope'
  AND INDEX_NAME = 'idx_role_resource';

SET @sql = IF(
    @has_idx_role_resource = 0,
    'ALTER TABLE `sys_role_data_scope` ADD UNIQUE INDEX `idx_role_resource` (`role_id` ASC, `resource_code` ASC)',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT COUNT(*) INTO @has_idx_role_feature_deleted_at
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_role_data_scope'
  AND INDEX_NAME = 'idx_sys_role_data_scope_deleted_at';

SET @sql = IF(
    @has_idx_role_feature_deleted_at = 0,
    'ALTER TABLE `sys_role_data_scope` ADD INDEX `idx_sys_role_data_scope_deleted_at` (`deleted_at` ASC)',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `sys_role_data_scope_dept` (
    `sys_role_data_scope_id` bigint UNSIGNED NOT NULL,
    `sys_dept_id` bigint UNSIGNED NOT NULL,
    PRIMARY KEY (`sys_role_data_scope_id`, `sys_dept_id`) USING BTREE,
    INDEX `fk_sys_role_data_scope_dept_sys_dept` (`sys_dept_id` ASC) USING BTREE,
    CONSTRAINT `fk_sys_role_data_scope_dept_scope` FOREIGN KEY (`sys_role_data_scope_id`) REFERENCES `sys_role_data_scope` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT `fk_sys_role_data_scope_dept_sys_dept` FOREIGN KEY (`sys_dept_id`) REFERENCES `sys_dept` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT '/api/v1/roles/:id/data-scopes', 'PUT', '角色管理', '分配数据权限', '', '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` WHERE `path` = '/api/v1/roles/:id/data-scopes' AND `method` = 'PUT'
);

SET @role_data_scope_api_id = (
    SELECT `id`
    FROM `sys_api`
    WHERE `path` = '/api/v1/roles/:id/data-scopes' AND `method` = 'PUT'
    ORDER BY `id` ASC
    LIMIT 1
);

INSERT INTO `sys_role_api` (`sys_role_id`, `sys_api_id`)
SELECT sra.`sys_role_id`, @role_data_scope_api_id
FROM `sys_role_api` sra
JOIN `sys_api` sa ON sa.`id` = sra.`sys_api_id`
WHERE sa.`path` = '/api/v1/roles/:id/apis'
  AND sa.`method` = 'PUT'
  AND @role_data_scope_api_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1
      FROM `sys_role_api` target
      WHERE target.`sys_role_id` = sra.`sys_role_id`
        AND target.`sys_api_id` = @role_data_scope_api_id
  )
GROUP BY sra.`sys_role_id`;

INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
SELECT 'p', sr.`code`, '/api/v1/roles/:id/data-scopes', 'PUT', '', '', ''
FROM `sys_role` sr
JOIN `sys_role_api` sra ON sra.`sys_role_id` = sr.`id`
JOIN `sys_api` sa ON sa.`id` = sra.`sys_api_id`
WHERE sa.`path` = '/api/v1/roles/:id/apis'
  AND sa.`method` = 'PUT'
  AND NOT EXISTS (
      SELECT 1
      FROM `casbin_rule` target
      WHERE target.`ptype` = 'p'
        AND target.`v0` = sr.`code`
        AND target.`v1` = '/api/v1/roles/:id/data-scopes'
        AND target.`v2` = 'PUT'
        AND IFNULL(target.`v3`, '') = ''
        AND IFNULL(target.`v4`, '') = ''
        AND IFNULL(target.`v5`, '') = ''
  );

SELECT '角色功能级数据权限升级脚本执行完成，请重新登录刷新权限缓存。' AS message;
