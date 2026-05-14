-- 补齐系统管理核心页面按钮权限
-- 执行方法: mysql -u root -p go-base < server/sql/add_system_admin_button_permissions.sql

SET NAMES utf8mb4;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT 0, '系统管理', '/system', 'Layout', 'setting', 1, 1, '', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_menu` WHERE `path` = '/system' AND `type` = 1 AND `deleted_at` IS NULL
);

SELECT @system_menu_id := `id`
FROM `sys_menu`
WHERE `path` = '/system' AND `type` = 1 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @system_menu_id, '角色管理', '/system/role', 'system/role/index', 'team', 2, 2, 'system:role:list', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:role:list' AND `deleted_at` IS NULL);

SELECT @role_menu_id := `id` FROM `sys_menu` WHERE `permission` = 'system:role:list' AND `deleted_at` IS NULL LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @role_menu_id, '新增', '', '', '', 1, 3, 'system:role:add', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:role:add' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @role_menu_id, '编辑', '', '', '', 2, 3, 'system:role:edit', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:role:edit' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @role_menu_id, '删除', '', '', '', 3, 3, 'system:role:delete', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:role:delete' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @role_menu_id, '分配权限', '', '', '', 4, 3, 'system:role:assign', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:role:assign' AND `deleted_at` IS NULL);

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @system_menu_id, '部门管理', '/system/dept', 'system/dept/index', 'apartment', 3, 2, 'system:dept:list', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dept:list' AND `deleted_at` IS NULL);

SELECT @dept_menu_id := `id` FROM `sys_menu` WHERE `permission` = 'system:dept:list' AND `deleted_at` IS NULL LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @dept_menu_id, '新增', '', '', '', 1, 3, 'system:dept:add', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dept:add' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @dept_menu_id, '编辑', '', '', '', 2, 3, 'system:dept:edit', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dept:edit' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @dept_menu_id, '删除', '', '', '', 3, 3, 'system:dept:delete', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dept:delete' AND `deleted_at` IS NULL);

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @system_menu_id, '菜单管理', '/system/menu', 'system/menu/index', 'menu', 4, 2, 'system:menu:list', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:menu:list' AND `deleted_at` IS NULL);

SELECT @menu_menu_id := `id` FROM `sys_menu` WHERE `permission` = 'system:menu:list' AND `deleted_at` IS NULL LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @menu_menu_id, '新增', '', '', '', 1, 3, 'system:menu:add', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:menu:add' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @menu_menu_id, '编辑', '', '', '', 2, 3, 'system:menu:edit', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:menu:edit' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @menu_menu_id, '删除', '', '', '', 3, 3, 'system:menu:delete', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:menu:delete' AND `deleted_at` IS NULL);

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @system_menu_id, '字典管理', '/system/dict', 'system/dict/index', 'AntDesignOutlined', 5, 2, 'system:dict:list', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dict:list' AND `deleted_at` IS NULL);

SELECT @dict_menu_id := `id` FROM `sys_menu` WHERE `permission` = 'system:dict:list' AND `deleted_at` IS NULL LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @dict_menu_id, '新增', '', '', '', 1, 3, 'system:dict:add', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dict:add' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @dict_menu_id, '编辑', '', '', '', 2, 3, 'system:dict:edit', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dict:edit' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @dict_menu_id, '删除', '', '', '', 3, 3, 'system:dict:delete', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:dict:delete' AND `deleted_at` IS NULL);

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @system_menu_id, 'API管理', '/system/api', 'system/api/index', 'api', 6, 2, 'system:api:list', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:api:list' AND `deleted_at` IS NULL);

SELECT @api_menu_id := `id` FROM `sys_menu` WHERE `permission` = 'system:api:list' AND `deleted_at` IS NULL LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @api_menu_id, '新增', '', '', '', 1, 3, 'system:api:add', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:api:add' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @api_menu_id, '编辑', '', '', '', 2, 3, 'system:api:edit', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:api:edit' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @api_menu_id, '删除', '', '', '', 3, 3, 'system:api:delete', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:api:delete' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @api_menu_id, '同步', '', '', '', 4, 3, 'system:api:sync', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:api:sync' AND `deleted_at` IS NULL);

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @system_menu_id, '系统配置', '/system/config', 'system/config/index', 'setting', 7, 2, 'system:config:list', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:config:list' AND `deleted_at` IS NULL);

SELECT @config_menu_id := `id` FROM `sys_menu` WHERE `permission` = 'system:config:list' AND `deleted_at` IS NULL LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @config_menu_id, '编辑', '', '', '', 1, 3, 'system:config:edit', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:config:edit' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @config_menu_id, '测试', '', '', '', 2, 3, 'system:config:test', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:config:test' AND `deleted_at` IS NULL);

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @system_menu_id, '文件管理', '/system/file', 'system/file/index', 'folder', 8, 2, 'system:file:list', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:file:list' AND `deleted_at` IS NULL);

SELECT @file_menu_id := `id` FROM `sys_menu` WHERE `permission` = 'system:file:list' AND `deleted_at` IS NULL LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @file_menu_id, '上传', '', '', '', 1, 3, 'system:file:upload', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:file:upload' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @file_menu_id, '删除', '', '', '', 2, 3, 'system:file:delete', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:file:delete' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @file_menu_id, '批量删除', '', '', '', 3, 3, 'system:file:batchDelete', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:file:batchDelete' AND `deleted_at` IS NULL);
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @file_menu_id, '文件迁移', '', '', '', 4, 3, 'system:file:migrate', 1, 0, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `permission` = 'system:file:migrate' AND `deleted_at` IS NULL);

INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT sr.`id`, sm.`id`
FROM `sys_role` sr
JOIN `sys_menu` sm ON sm.`permission` IN (
    'system:role:list', 'system:role:add', 'system:role:edit', 'system:role:delete', 'system:role:assign',
    'system:dept:list', 'system:dept:add', 'system:dept:edit', 'system:dept:delete',
    'system:menu:list', 'system:menu:add', 'system:menu:edit', 'system:menu:delete',
    'system:dict:list', 'system:dict:add', 'system:dict:edit', 'system:dict:delete',
    'system:api:list', 'system:api:add', 'system:api:edit', 'system:api:delete', 'system:api:sync',
    'system:config:list', 'system:config:edit', 'system:config:test',
    'system:file:list', 'system:file:upload', 'system:file:delete', 'system:file:batchDelete', 'system:file:migrate'
)
WHERE sr.`code` IN ('admin', 'system_admin')
  AND sm.`deleted_at` IS NULL
  AND NOT EXISTS (
      SELECT 1
      FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = sr.`id`
        AND target.`sys_menu_id` = sm.`id`
  );
