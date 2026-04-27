-- 为系统补齐部门管理与数据权限底座
-- 执行方法: mysql -u root -p go-base < server/sql/add_department_permission_foundation.sql

SET NAMES utf8mb4;

-- 1. 结构升级
CREATE TABLE IF NOT EXISTS `sys_dept` (
    `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NULL DEFAULT NULL,
    `updated_at` datetime(3) NULL DEFAULT NULL,
    `deleted_at` datetime(3) NULL DEFAULT NULL,
    `parent_id` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT '父部门ID',
    `ancestors` varchar(500) NULL DEFAULT NULL COMMENT '祖级列表',
    `name` varchar(100) NULL DEFAULT NULL COMMENT '部门名称',
    `sort` bigint NULL DEFAULT 0 COMMENT '排序',
    `status` bigint NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
    `remark` varchar(255) NULL DEFAULT NULL COMMENT '备注',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `idx_sys_dept_deleted_at` (`deleted_at` ASC) USING BTREE,
    INDEX `idx_sys_dept_parent_id` (`parent_id` ASC) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SELECT COUNT(*) INTO @has_sys_user_dept_id
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_user'
  AND COLUMN_NAME = 'dept_id';

SET @sql = IF(
    @has_sys_user_dept_id = 0,
    'ALTER TABLE `sys_user` ADD COLUMN `dept_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT ''部门ID'' AFTER `status`',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT COUNT(*) INTO @has_idx_sys_user_dept_id
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_user'
  AND INDEX_NAME = 'idx_sys_user_dept_id';

SET @sql = IF(
    @has_idx_sys_user_dept_id = 0,
    'ALTER TABLE `sys_user` ADD INDEX `idx_sys_user_dept_id` (`dept_id` ASC)',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT COUNT(*) INTO @has_sys_role_data_scope
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'sys_role'
  AND COLUMN_NAME = 'data_scope';

SET @sql = IF(
    @has_sys_role_data_scope = 0,
    'ALTER TABLE `sys_role` ADD COLUMN `data_scope` bigint NULL DEFAULT 1 COMMENT ''数据范围 1全部 2自定义 3本部门 4本部门及下级 5仅本人'' AFTER `status`',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `sys_role_dept` (
    `sys_role_id` bigint UNSIGNED NOT NULL,
    `sys_dept_id` bigint UNSIGNED NOT NULL,
    PRIMARY KEY (`sys_role_id`, `sys_dept_id`) USING BTREE,
    INDEX `fk_sys_role_dept_sys_dept` (`sys_dept_id` ASC) USING BTREE,
    CONSTRAINT `fk_sys_role_dept_sys_dept` FOREIGN KEY (`sys_dept_id`) REFERENCES `sys_dept` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT `fk_sys_role_dept_sys_role` FOREIGN KEY (`sys_role_id`) REFERENCES `sys_role` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- 2. 根部门与兼容回填
INSERT INTO `sys_dept` (`parent_id`, `ancestors`, `name`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
SELECT 0, '0', '平台', 1, 1, '系统根部门', NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_dept` WHERE `parent_id` = 0 AND `name` = '平台' AND `deleted_at` IS NULL
);

SET @root_dept_id = (
    SELECT `id`
    FROM `sys_dept`
    WHERE `parent_id` = 0 AND `name` = '平台' AND `deleted_at` IS NULL
    ORDER BY `id` ASC
    LIMIT 1
);

UPDATE `sys_user`
SET `dept_id` = @root_dept_id
WHERE (@root_dept_id IS NOT NULL)
  AND (`dept_id` IS NULL OR `dept_id` = 0);

UPDATE `sys_role`
SET `data_scope` = 1
WHERE `data_scope` IS NULL OR `data_scope` = 0;

-- 3. 系统配置
INSERT INTO `sys_config` (`name`, `key`, `value`, `value_type`, `remark`, `created_at`, `updated_at`)
SELECT '部门模块显示', 'dept_module_enabled', 'true', 'string', '后台菜单中是否显示部门管理模块', NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_config` WHERE `key` = 'dept_module_enabled'
);

-- 4. 菜单与按钮权限
SET @system_menu_id = (
    SELECT `id`
    FROM `sys_menu`
    WHERE `path` = '/system' AND `type` = 1
    ORDER BY `id` ASC
    LIMIT 1
);

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @system_menu_id, '部门管理', '/system/dept', 'system/dept/index', 'apartment', 3, 2, 'system:dept:list', 1, 0, NOW(3), NOW(3)
WHERE @system_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dept:list'
  );

SET @dept_menu_id = (
    SELECT `id`
    FROM `sys_menu`
    WHERE `permission` = 'system:dept:list'
    ORDER BY `id` ASC
    LIMIT 1
);

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @dept_menu_id, '新增', '', '', '', 1, 3, 'system:dept:add', 1, 0, NOW(3), NOW(3)
WHERE @dept_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dept:add'
  );

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @dept_menu_id, '编辑', '', '', '', 2, 3, 'system:dept:edit', 1, 0, NOW(3), NOW(3)
WHERE @dept_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dept:edit'
  );

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @dept_menu_id, '删除', '', '', '', 3, 3, 'system:dept:delete', 1, 0, NOW(3), NOW(3)
WHERE @dept_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dept:delete'
  );

-- 5. API 元数据
INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT '/api/v1/depts/tree', 'GET', '部门管理', '部门树', '', '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` WHERE `path` = '/api/v1/depts/tree' AND `method` = 'GET'
);

INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT '/api/v1/depts/:id', 'GET', '部门管理', '部门详情', '', '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` WHERE `path` = '/api/v1/depts/:id' AND `method` = 'GET'
);

INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT '/api/v1/depts', 'POST', '部门管理', '创建部门', '', '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` WHERE `path` = '/api/v1/depts' AND `method` = 'POST'
);

INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT '/api/v1/depts/:id', 'PUT', '部门管理', '更新部门', '', '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` WHERE `path` = '/api/v1/depts/:id' AND `method` = 'PUT'
);

INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT '/api/v1/depts/:id', 'DELETE', '部门管理', '删除部门', '', '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` WHERE `path` = '/api/v1/depts/:id' AND `method` = 'DELETE'
);

SET @dept_tree_api_id = (
    SELECT `id` FROM `sys_api` WHERE `path` = '/api/v1/depts/tree' AND `method` = 'GET' ORDER BY `id` ASC LIMIT 1
);
SET @dept_detail_api_id = (
    SELECT `id` FROM `sys_api` WHERE `path` = '/api/v1/depts/:id' AND `method` = 'GET' ORDER BY `id` ASC LIMIT 1
);
SET @dept_create_api_id = (
    SELECT `id` FROM `sys_api` WHERE `path` = '/api/v1/depts' AND `method` = 'POST' ORDER BY `id` ASC LIMIT 1
);
SET @dept_update_api_id = (
    SELECT `id` FROM `sys_api` WHERE `path` = '/api/v1/depts/:id' AND `method` = 'PUT' ORDER BY `id` ASC LIMIT 1
);
SET @dept_delete_api_id = (
    SELECT `id` FROM `sys_api` WHERE `path` = '/api/v1/depts/:id' AND `method` = 'DELETE' ORDER BY `id` ASC LIMIT 1
);

-- 6. 给内置管理员角色补齐菜单与 API 权限
INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT sr.`id`, sm.`id`
FROM `sys_role` sr
JOIN `sys_menu` sm ON sm.`permission` IN ('system:dept:list', 'system:dept:add', 'system:dept:edit', 'system:dept:delete')
WHERE sr.`code` IN ('admin', 'system_admin')
  AND NOT EXISTS (
      SELECT 1
      FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = sr.`id`
        AND target.`sys_menu_id` = sm.`id`
  );

INSERT INTO `sys_role_api` (`sys_role_id`, `sys_api_id`)
SELECT sr.`id`, api_id_map.`api_id`
FROM `sys_role` sr
JOIN (
    SELECT @dept_tree_api_id AS `api_id`
    UNION ALL SELECT @dept_detail_api_id
    UNION ALL SELECT @dept_create_api_id
    UNION ALL SELECT @dept_update_api_id
    UNION ALL SELECT @dept_delete_api_id
) api_id_map
WHERE sr.`code` IN ('admin', 'system_admin')
  AND api_id_map.`api_id` IS NOT NULL
  AND NOT EXISTS (
      SELECT 1
      FROM `sys_role_api` target
      WHERE target.`sys_role_id` = sr.`id`
        AND target.`sys_api_id` = api_id_map.`api_id`
  );

INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
SELECT 'p', sr.`code`, api_map.`path`, api_map.`method`, '', '', ''
FROM `sys_role` sr
JOIN (
    SELECT '/api/v1/depts/tree' AS `path`, 'GET' AS `method`
    UNION ALL SELECT '/api/v1/depts/:id', 'GET'
    UNION ALL SELECT '/api/v1/depts', 'POST'
    UNION ALL SELECT '/api/v1/depts/:id', 'PUT'
    UNION ALL SELECT '/api/v1/depts/:id', 'DELETE'
) api_map
WHERE sr.`code` IN ('admin', 'system_admin')
  AND NOT EXISTS (
      SELECT 1
      FROM `casbin_rule` target
      WHERE target.`ptype` = 'p'
        AND target.`v0` = sr.`code`
        AND target.`v1` = api_map.`path`
        AND target.`v2` = api_map.`method`
        AND IFNULL(target.`v3`, '') = ''
        AND IFNULL(target.`v4`, '') = ''
        AND IFNULL(target.`v5`, '') = ''
  );

SELECT '部门管理与数据权限底座升级脚本执行完成，请重新登录刷新权限缓存。' AS message;
