-- 第一阶段 CMDB：主机分组、标签、SSH凭据、主机台账基础表
-- 执行方法: mysql -u root -p go-base < server/sql/add_cmdb_host_foundation.sql

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `cmdb_host_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '分组名称',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序值',
  `remark` varchar(500) NOT NULL DEFAULT '' COMMENT '备注',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态:1启用,0禁用',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cmdb_host_group_name` (`name`),
  KEY `idx_cmdb_host_group_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机分组表';

CREATE TABLE IF NOT EXISTS `cmdb_host_tag` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '标签名称',
  `color` varchar(30) NOT NULL DEFAULT '' COMMENT '标签颜色',
  `remark` varchar(500) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cmdb_host_tag_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机标签表';

CREATE TABLE IF NOT EXISTS `cmdb_ssh_credential` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '凭据名称',
  `auth_type` varchar(20) NOT NULL COMMENT '认证方式:password/private_key',
  `username` varchar(100) NOT NULL COMMENT '登录用户名',
  `password` text COMMENT '登录密码',
  `private_key` longtext COMMENT '私钥内容',
  `passphrase` text COMMENT '私钥口令',
  `remark` varchar(500) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cmdb_ssh_credential_name` (`name`),
  KEY `idx_cmdb_ssh_credential_auth_type` (`auth_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='SSH凭据表';

CREATE TABLE IF NOT EXISTS `cmdb_host` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '主机名称',
  `group_id` bigint unsigned NOT NULL COMMENT '主机分组ID',
  `environment` varchar(50) NOT NULL DEFAULT '' COMMENT '环境标识',
  `owner` varchar(100) NOT NULL DEFAULT '' COMMENT '负责人',
  `private_ip` varchar(45) NOT NULL DEFAULT '' COMMENT '内网IP',
  `public_ip` varchar(45) NOT NULL DEFAULT '' COMMENT '公网IP',
  `ssh_host` varchar(255) NOT NULL COMMENT 'SSH连接地址',
  `ssh_port` int NOT NULL DEFAULT 22 COMMENT 'SSH端口',
  `credential_id` bigint unsigned NOT NULL COMMENT 'SSH凭据ID',
  `remark` varchar(500) NOT NULL DEFAULT '' COMMENT '备注',
  `verify_status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '校验状态:pending/success/failed',
  `verify_message` varchar(500) NOT NULL DEFAULT '' COMMENT '校验结果说明',
  `last_verified_at` datetime(3) DEFAULT NULL COMMENT '最后校验时间',
  `hostname` varchar(255) NOT NULL DEFAULT '' COMMENT '主机名',
  `os` varchar(100) NOT NULL DEFAULT '' COMMENT '操作系统',
  `platform` varchar(100) NOT NULL DEFAULT '' COMMENT '发行版标识',
  `platform_version` varchar(255) NOT NULL DEFAULT '' COMMENT '系统版本',
  `kernel_version` varchar(100) NOT NULL DEFAULT '' COMMENT '内核版本',
  `architecture` varchar(50) NOT NULL DEFAULT '' COMMENT '系统架构',
  `cpu_cores` int NOT NULL DEFAULT 0 COMMENT 'CPU核数',
  `memory_mb` bigint NOT NULL DEFAULT 0 COMMENT '内存大小MB',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cmdb_host_name` (`name`),
  KEY `idx_cmdb_host_group_id` (`group_id`),
  KEY `idx_cmdb_host_credential_id` (`credential_id`),
  KEY `idx_cmdb_host_verify_status` (`verify_status`),
  KEY `idx_cmdb_host_private_ip` (`private_ip`),
  KEY `idx_cmdb_host_public_ip` (`public_ip`),
  KEY `idx_cmdb_host_ssh_host` (`ssh_host`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机台账表';

SELECT COUNT(*) INTO @cmdb_host_platform_version_exists
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'cmdb_host'
  AND COLUMN_NAME = 'platform_version';

SET @sql = IF(@cmdb_host_platform_version_exists = 0,
  'ALTER TABLE `cmdb_host` ADD COLUMN `platform_version` varchar(255) NOT NULL DEFAULT '''' COMMENT ''系统版本'' AFTER `platform`',
  'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SELECT COUNT(*) INTO @cmdb_host_kernel_version_exists
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'cmdb_host'
  AND COLUMN_NAME = 'kernel_version';

SET @sql = IF(@cmdb_host_kernel_version_exists = 0,
  'ALTER TABLE `cmdb_host` ADD COLUMN `kernel_version` varchar(100) NOT NULL DEFAULT '''' COMMENT ''内核版本'' AFTER `platform_version`',
  'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `cmdb_host_tag_rel` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `host_id` bigint unsigned NOT NULL COMMENT '主机ID',
  `tag_id` bigint unsigned NOT NULL COMMENT '标签ID',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cmdb_host_tag_rel` (`host_id`,`tag_id`),
  KEY `idx_cmdb_host_tag_rel_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机标签关联表';

SELECT @cmdb_root_id := `id`
FROM `sys_menu`
WHERE `permission` = 'cmdb' AND `type` = 1 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT 0, 'CMDB管理', '/cmdb', 'Layout', 'ClusterOutlined', 31, 1, 'cmdb', 1, 0, NOW(3), NOW(3)
WHERE @cmdb_root_id IS NULL;

SELECT @cmdb_root_id := `id`
FROM `sys_menu`
WHERE `permission` = 'cmdb' AND `type` = 1 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @cmdb_root_id, t.`name`, t.`path`, t.`component`, t.`icon`, t.`sort`, 2, t.`permission`, 1, 0, NOW(3), NOW(3)
FROM (
    SELECT '主机管理' AS `name`, '/cmdb/host' AS `path`, 'cmdb/host/index' AS `component`, 'DesktopOutlined' AS `icon`, 1 AS `sort`, 'cmdb:host:list' AS `permission`
    UNION ALL SELECT '主机分组', '/cmdb/group', 'cmdb/group/index', 'AppstoreOutlined', 2, 'cmdb:group:list'
    UNION ALL SELECT '主机标签', '/cmdb/tag', 'cmdb/tag/index', 'TagsOutlined', 3, 'cmdb:tag:list'
    UNION ALL SELECT 'SSH凭据', '/cmdb/credential', 'cmdb/credential/index', 'KeyOutlined', 4, 'cmdb:credential:list'
) t
WHERE @cmdb_root_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` sm
      WHERE sm.`permission` = t.`permission` AND sm.`type` = 2 AND sm.`deleted_at` IS NULL
  );

UPDATE `sys_menu` page_menu
JOIN (
    SELECT 'cmdb:host:list' AS `permission`
    UNION ALL SELECT 'cmdb:group:list'
    UNION ALL SELECT 'cmdb:tag:list'
    UNION ALL SELECT 'cmdb:credential:list'
) t ON t.`permission` = page_menu.`permission`
SET page_menu.`parent_id` = @cmdb_root_id,
    page_menu.`updated_at` = NOW(3)
WHERE @cmdb_root_id IS NOT NULL
  AND page_menu.`deleted_at` IS NULL
  AND page_menu.`type` = 2
  AND page_menu.`parent_id` <> @cmdb_root_id;

SELECT @cmdb_host_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'cmdb:host:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

SELECT @cmdb_group_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'cmdb:group:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

SELECT @cmdb_tag_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'cmdb:tag:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

SELECT @cmdb_credential_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'cmdb:credential:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT CASE
         WHEN t.`parent_permission` = 'cmdb:host:list' THEN @cmdb_host_menu_id
         WHEN t.`parent_permission` = 'cmdb:group:list' THEN @cmdb_group_menu_id
         WHEN t.`parent_permission` = 'cmdb:tag:list' THEN @cmdb_tag_menu_id
         WHEN t.`parent_permission` = 'cmdb:credential:list' THEN @cmdb_credential_menu_id
       END,
       t.`name`, '', '', '', t.`sort`, 3, t.`permission`, 1, 0, NOW(3), NOW(3)
FROM (
    SELECT 'cmdb:host:list' AS `parent_permission`, '查看' AS `name`, 1 AS `sort`, 'cmdb:host:view' AS `permission`
    UNION ALL SELECT 'cmdb:host:list', '新增', 2, 'cmdb:host:create'
    UNION ALL SELECT 'cmdb:host:list', '编辑', 3, 'cmdb:host:update'
    UNION ALL SELECT 'cmdb:host:list', '删除', 4, 'cmdb:host:delete'
    UNION ALL SELECT 'cmdb:host:list', '导入', 5, 'cmdb:host:import'
    UNION ALL SELECT 'cmdb:host:list', '校验', 6, 'cmdb:host:verify'
    UNION ALL SELECT 'cmdb:group:list', '新增', 1, 'cmdb:group:create'
    UNION ALL SELECT 'cmdb:group:list', '编辑', 2, 'cmdb:group:update'
    UNION ALL SELECT 'cmdb:group:list', '删除', 3, 'cmdb:group:delete'
    UNION ALL SELECT 'cmdb:tag:list', '新增', 1, 'cmdb:tag:create'
    UNION ALL SELECT 'cmdb:tag:list', '编辑', 2, 'cmdb:tag:update'
    UNION ALL SELECT 'cmdb:tag:list', '删除', 3, 'cmdb:tag:delete'
    UNION ALL SELECT 'cmdb:credential:list', '查看', 1, 'cmdb:credential:view'
    UNION ALL SELECT 'cmdb:credential:list', '新增', 2, 'cmdb:credential:create'
    UNION ALL SELECT 'cmdb:credential:list', '编辑', 3, 'cmdb:credential:update'
    UNION ALL SELECT 'cmdb:credential:list', '删除', 4, 'cmdb:credential:delete'
) t
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_menu` sm
    WHERE sm.`permission` = t.`permission` AND sm.`type` = 3 AND sm.`deleted_at` IS NULL
);

UPDATE `sys_menu` button_menu
JOIN (
    SELECT 'cmdb:host:list' AS `parent_permission`, 'cmdb:host:view' AS `permission`
    UNION ALL SELECT 'cmdb:host:list', 'cmdb:host:create'
    UNION ALL SELECT 'cmdb:host:list', 'cmdb:host:update'
    UNION ALL SELECT 'cmdb:host:list', 'cmdb:host:delete'
    UNION ALL SELECT 'cmdb:host:list', 'cmdb:host:import'
    UNION ALL SELECT 'cmdb:host:list', 'cmdb:host:verify'
    UNION ALL SELECT 'cmdb:group:list', 'cmdb:group:create'
    UNION ALL SELECT 'cmdb:group:list', 'cmdb:group:update'
    UNION ALL SELECT 'cmdb:group:list', 'cmdb:group:delete'
    UNION ALL SELECT 'cmdb:tag:list', 'cmdb:tag:create'
    UNION ALL SELECT 'cmdb:tag:list', 'cmdb:tag:update'
    UNION ALL SELECT 'cmdb:tag:list', 'cmdb:tag:delete'
    UNION ALL SELECT 'cmdb:credential:list', 'cmdb:credential:view'
    UNION ALL SELECT 'cmdb:credential:list', 'cmdb:credential:create'
    UNION ALL SELECT 'cmdb:credential:list', 'cmdb:credential:update'
    UNION ALL SELECT 'cmdb:credential:list', 'cmdb:credential:delete'
) t ON t.`permission` = button_menu.`permission`
SET button_menu.`parent_id` = CASE
    WHEN t.`parent_permission` = 'cmdb:host:list' THEN @cmdb_host_menu_id
    WHEN t.`parent_permission` = 'cmdb:group:list' THEN @cmdb_group_menu_id
    WHEN t.`parent_permission` = 'cmdb:tag:list' THEN @cmdb_tag_menu_id
    WHEN t.`parent_permission` = 'cmdb:credential:list' THEN @cmdb_credential_menu_id
END,
button_menu.`updated_at` = NOW(3)
WHERE button_menu.`deleted_at` IS NULL
  AND button_menu.`type` = 3;

INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT t.`path`, t.`method`, 'CMDB管理', t.`description`, '', '', 1, NOW(3), NOW(3)
FROM (
    SELECT '/api/v1/cmdb/host-groups' AS `path`, 'GET' AS `method`, '主机分组列表' AS `description`
    UNION ALL SELECT '/api/v1/cmdb/host-groups', 'POST', '创建主机分组'
    UNION ALL SELECT '/api/v1/cmdb/host-groups/:id', 'PUT', '更新主机分组'
    UNION ALL SELECT '/api/v1/cmdb/host-groups/:id', 'DELETE', '删除主机分组'
    UNION ALL SELECT '/api/v1/cmdb/host-tags', 'GET', '主机标签列表'
    UNION ALL SELECT '/api/v1/cmdb/host-tags', 'POST', '创建主机标签'
    UNION ALL SELECT '/api/v1/cmdb/host-tags/:id', 'PUT', '更新主机标签'
    UNION ALL SELECT '/api/v1/cmdb/host-tags/:id', 'DELETE', '删除主机标签'
    UNION ALL SELECT '/api/v1/cmdb/ssh-credentials', 'GET', 'SSH凭据列表'
    UNION ALL SELECT '/api/v1/cmdb/ssh-credentials/options', 'GET', 'SSH凭据选项'
    UNION ALL SELECT '/api/v1/cmdb/ssh-credentials/:id', 'GET', 'SSH凭据详情'
    UNION ALL SELECT '/api/v1/cmdb/ssh-credentials', 'POST', '创建SSH凭据'
    UNION ALL SELECT '/api/v1/cmdb/ssh-credentials/:id', 'PUT', '更新SSH凭据'
    UNION ALL SELECT '/api/v1/cmdb/ssh-credentials/:id', 'DELETE', '删除SSH凭据'
    UNION ALL SELECT '/api/v1/cmdb/hosts', 'GET', '主机列表'
    UNION ALL SELECT '/api/v1/cmdb/hosts/:id', 'GET', '主机详情'
    UNION ALL SELECT '/api/v1/cmdb/hosts', 'POST', '创建主机'
    UNION ALL SELECT '/api/v1/cmdb/hosts/:id', 'PUT', '更新主机'
    UNION ALL SELECT '/api/v1/cmdb/hosts/:id', 'DELETE', '删除主机'
    UNION ALL SELECT '/api/v1/cmdb/hosts/:id/verify', 'POST', '校验主机'
    UNION ALL SELECT '/api/v1/cmdb/hosts/import-template', 'GET', '下载主机导入模板'
    UNION ALL SELECT '/api/v1/cmdb/hosts/import', 'POST', '导入主机'
) t
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` sa
    WHERE sa.`path` = t.`path` AND sa.`method` = t.`method` AND sa.`deleted_at` IS NULL
);

INSERT INTO `sys_menu_api` (`sys_menu_id`, `sys_api_id`)
SELECT menu_item.`id`, api_item.`id`
FROM (
    SELECT 'cmdb:group:list' AS `menu_permission`, '/api/v1/cmdb/host-groups' AS `api_path`, 'GET' AS `api_method`
    UNION ALL SELECT 'cmdb:group:create', '/api/v1/cmdb/host-groups', 'POST'
    UNION ALL SELECT 'cmdb:group:update', '/api/v1/cmdb/host-groups/:id', 'PUT'
    UNION ALL SELECT 'cmdb:group:delete', '/api/v1/cmdb/host-groups/:id', 'DELETE'
    UNION ALL SELECT 'cmdb:tag:list', '/api/v1/cmdb/host-tags', 'GET'
    UNION ALL SELECT 'cmdb:tag:create', '/api/v1/cmdb/host-tags', 'POST'
    UNION ALL SELECT 'cmdb:tag:update', '/api/v1/cmdb/host-tags/:id', 'PUT'
    UNION ALL SELECT 'cmdb:tag:delete', '/api/v1/cmdb/host-tags/:id', 'DELETE'
    UNION ALL SELECT 'cmdb:credential:list', '/api/v1/cmdb/ssh-credentials', 'GET'
    UNION ALL SELECT 'cmdb:credential:list', '/api/v1/cmdb/ssh-credentials/options', 'GET'
    UNION ALL SELECT 'cmdb:credential:view', '/api/v1/cmdb/ssh-credentials/:id', 'GET'
    UNION ALL SELECT 'cmdb:credential:create', '/api/v1/cmdb/ssh-credentials', 'POST'
    UNION ALL SELECT 'cmdb:credential:update', '/api/v1/cmdb/ssh-credentials/:id', 'PUT'
    UNION ALL SELECT 'cmdb:credential:delete', '/api/v1/cmdb/ssh-credentials/:id', 'DELETE'
    UNION ALL SELECT 'cmdb:host:list', '/api/v1/cmdb/hosts', 'GET'
    UNION ALL SELECT 'cmdb:host:view', '/api/v1/cmdb/hosts/:id', 'GET'
    UNION ALL SELECT 'cmdb:host:create', '/api/v1/cmdb/hosts', 'POST'
    UNION ALL SELECT 'cmdb:host:update', '/api/v1/cmdb/hosts/:id', 'PUT'
    UNION ALL SELECT 'cmdb:host:delete', '/api/v1/cmdb/hosts/:id', 'DELETE'
    UNION ALL SELECT 'cmdb:host:verify', '/api/v1/cmdb/hosts/:id/verify', 'POST'
    UNION ALL SELECT 'cmdb:host:import', '/api/v1/cmdb/hosts/import-template', 'GET'
    UNION ALL SELECT 'cmdb:host:import', '/api/v1/cmdb/hosts/import', 'POST'
) t
JOIN `sys_menu` menu_item ON menu_item.`permission` = t.`menu_permission` AND menu_item.`deleted_at` IS NULL
JOIN `sys_api` api_item ON api_item.`path` = t.`api_path` AND api_item.`method` = t.`api_method` AND api_item.`deleted_at` IS NULL
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_menu_api` existing
    WHERE existing.`sys_menu_id` = menu_item.`id`
      AND existing.`sys_api_id` = api_item.`id`
);

INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT sr.`id`, menu_item.`id`
FROM `sys_role` sr
JOIN `sys_menu` menu_item ON menu_item.`deleted_at` IS NULL
WHERE sr.`code` IN ('admin', 'system_admin')
  AND sr.`deleted_at` IS NULL
  AND menu_item.`permission` IN (
      'cmdb',
      'cmdb:host:list',
      'cmdb:host:view',
      'cmdb:host:create',
      'cmdb:host:update',
      'cmdb:host:delete',
      'cmdb:host:import',
      'cmdb:host:verify',
      'cmdb:group:list',
      'cmdb:group:create',
      'cmdb:group:update',
      'cmdb:group:delete',
      'cmdb:tag:list',
      'cmdb:tag:create',
      'cmdb:tag:update',
      'cmdb:tag:delete',
      'cmdb:credential:list',
      'cmdb:credential:view',
      'cmdb:credential:create',
      'cmdb:credential:update',
      'cmdb:credential:delete'
  )
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = sr.`id`
        AND target.`sys_menu_id` = menu_item.`id`
  );

INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
SELECT DISTINCT 'p', sr.`code`, api_item.`path`, api_item.`method`, '', '', ''
FROM `sys_role` sr
JOIN `sys_menu` menu_item ON menu_item.`deleted_at` IS NULL
JOIN `sys_role_menu` srm ON srm.`sys_role_id` = sr.`id` AND srm.`sys_menu_id` = menu_item.`id`
JOIN `sys_menu_api` sma ON sma.`sys_menu_id` = menu_item.`id`
JOIN `sys_api` api_item ON api_item.`id` = sma.`sys_api_id` AND api_item.`deleted_at` IS NULL
WHERE sr.`code` IN ('admin', 'system_admin')
  AND sr.`deleted_at` IS NULL
  AND menu_item.`permission` IN (
      'cmdb:host:list',
      'cmdb:host:view',
      'cmdb:host:create',
      'cmdb:host:update',
      'cmdb:host:delete',
      'cmdb:host:import',
      'cmdb:host:verify',
      'cmdb:group:list',
      'cmdb:group:create',
      'cmdb:group:update',
      'cmdb:group:delete',
      'cmdb:tag:list',
      'cmdb:tag:create',
      'cmdb:tag:update',
      'cmdb:tag:delete',
      'cmdb:credential:list',
      'cmdb:credential:view',
      'cmdb:credential:create',
      'cmdb:credential:update',
      'cmdb:credential:delete'
  )
  AND NOT EXISTS (
      SELECT 1 FROM `casbin_rule` target
      WHERE target.`ptype` = 'p'
        AND target.`v0` = sr.`code`
        AND target.`v1` = api_item.`path`
        AND target.`v2` = api_item.`method`
        AND IFNULL(target.`v3`, '') = ''
        AND IFNULL(target.`v4`, '') = ''
        AND IFNULL(target.`v5`, '') = ''
  );

SELECT '第一阶段 CMDB 基础表、菜单、按钮权限、API 元数据与 Casbin 策略已补齐，请重启服务或重新登录刷新权限缓存。' AS message;
