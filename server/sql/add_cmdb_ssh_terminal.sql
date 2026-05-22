-- 第二阶段 CMDB：SSH 在线终端
-- 执行方法: mysql -u root -p go-base < server/sql/add_cmdb_ssh_terminal.sql

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `cmdb_terminal_session` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `host_id` bigint unsigned NOT NULL COMMENT '主机ID',
  `user_id` bigint unsigned NOT NULL COMMENT '发起用户ID',
  `username_snapshot` varchar(100) NOT NULL DEFAULT '' COMMENT '发起用户名快照',
  `credential_id_snapshot` bigint unsigned NOT NULL COMMENT '凭据ID快照',
  `client_ip` varchar(45) NOT NULL DEFAULT '' COMMENT '来源IP',
  `status` varchar(20) NOT NULL DEFAULT 'prepared' COMMENT '会话状态:prepared/active/closed/failed',
  `start_time` datetime(3) DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '结束时间',
  `idle_timeout_seconds` int NOT NULL DEFAULT 1800 COMMENT '空闲超时秒数',
  `disconnect_reason` varchar(50) NOT NULL DEFAULT '' COMMENT '断开原因',
  `forced_by_user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '强制断开操作人ID',
  `host_key_fingerprint` varchar(255) NOT NULL DEFAULT '' COMMENT '主机指纹',
  `last_activity_at` datetime(3) DEFAULT NULL COMMENT '最后活动时间',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_cmdb_terminal_session_host_id` (`host_id`),
  KEY `idx_cmdb_terminal_session_user_id` (`user_id`),
  KEY `idx_cmdb_terminal_session_status` (`status`),
  KEY `idx_cmdb_terminal_session_last_activity` (`last_activity_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='终端会话表';

CREATE TABLE IF NOT EXISTS `cmdb_terminal_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `session_id` bigint unsigned NOT NULL COMMENT '终端会话ID',
  `seq` bigint unsigned NOT NULL DEFAULT 0 COMMENT '日志序号',
  `stream_type` varchar(20) NOT NULL COMMENT '流类型:input/output/system',
  `content` longtext COMMENT '日志内容',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_cmdb_terminal_log_session_id` (`session_id`),
  KEY `idx_cmdb_terminal_log_session_seq` (`session_id`, `seq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='终端日志表';

CREATE TABLE IF NOT EXISTS `cmdb_host_ssh_fingerprint` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `host_id` bigint unsigned NOT NULL COMMENT '主机ID',
  `algorithm` varchar(50) NOT NULL DEFAULT '' COMMENT '指纹算法',
  `fingerprint` varchar(255) NOT NULL DEFAULT '' COMMENT '指纹值',
  `first_seen_at` datetime(3) DEFAULT NULL COMMENT '首次记录时间',
  `last_verified_at` datetime(3) DEFAULT NULL COMMENT '最后校验时间',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cmdb_host_ssh_fingerprint_host_id` (`host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机SSH指纹表';

SELECT @cmdb_root_id := `id`
FROM `sys_menu`
WHERE `permission` = 'cmdb' AND `type` = 1 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @cmdb_root_id, '终端工作台', '/cmdb/terminal/:hostId', 'cmdb/terminal/index', 'CodeOutlined', 5, 2, 'cmdb:terminal:list', 1, 1, NOW(3), NOW(3)
WHERE @cmdb_root_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` sm
      WHERE sm.`permission` = 'cmdb:terminal:list' AND sm.`type` = 2 AND sm.`deleted_at` IS NULL
  );

UPDATE `sys_menu`
SET `path` = '/cmdb/terminal/:hostId',
    `component` = 'cmdb/terminal/index',
    `hidden` = 1,
    `updated_at` = NOW(3)
WHERE `permission` = 'cmdb:terminal:list'
  AND `type` = 2
  AND `deleted_at` IS NULL
  AND (
    `path` <> '/cmdb/terminal/:hostId'
    OR `component` <> 'cmdb/terminal/index'
    OR `hidden` <> 1
  );

SELECT @cmdb_terminal_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'cmdb:terminal:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @cmdb_terminal_menu_id, t.`name`, '', '', '', t.`sort`, 3, t.`permission`, 1, 0, NOW(3), NOW(3)
FROM (
    SELECT '连接' AS `name`, 1 AS `sort`, 'cmdb:terminal:connect' AS `permission`
    UNION ALL SELECT '查看', 2, 'cmdb:terminal:view'
    UNION ALL SELECT '断开', 3, 'cmdb:terminal:disconnect'
    UNION ALL SELECT '强制断开', 4, 'cmdb:terminal:force_disconnect'
    UNION ALL SELECT '审计日志', 5, 'cmdb:terminal:audit'
) t
WHERE @cmdb_terminal_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` sm
      WHERE sm.`permission` = t.`permission` AND sm.`type` = 3 AND sm.`deleted_at` IS NULL
  );

UPDATE `sys_menu` button_menu
JOIN (
    SELECT 'cmdb:terminal:connect' AS `permission`
    UNION ALL SELECT 'cmdb:terminal:view'
    UNION ALL SELECT 'cmdb:terminal:disconnect'
    UNION ALL SELECT 'cmdb:terminal:force_disconnect'
    UNION ALL SELECT 'cmdb:terminal:audit'
) t ON t.`permission` = button_menu.`permission`
SET button_menu.`parent_id` = @cmdb_terminal_menu_id,
    button_menu.`updated_at` = NOW(3)
WHERE @cmdb_terminal_menu_id IS NOT NULL
  AND button_menu.`deleted_at` IS NULL
  AND button_menu.`type` = 3;

INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT t.`path`, t.`method`, 'CMDB管理', t.`description`, '', '', t.`need_auth`, NOW(3), NOW(3)
FROM (
    SELECT '/api/v1/cmdb/terminal/sessions' AS `path`, 'POST' AS `method`, '创建SSH终端会话' AS `description`, 1 AS `need_auth`
    UNION ALL SELECT '/api/v1/cmdb/terminal/sessions', 'GET', '终端会话列表', 1
    UNION ALL SELECT '/api/v1/cmdb/terminal/sessions/:id', 'GET', '终端会话详情', 1
    UNION ALL SELECT '/api/v1/cmdb/terminal/sessions/:id/logs', 'GET', '终端会话日志', 1
    UNION ALL SELECT '/api/v1/cmdb/terminal/sessions/:id/disconnect', 'POST', '断开终端会话', 1
    UNION ALL SELECT '/api/v1/cmdb/terminal/sessions/:id/force-disconnect', 'POST', '强制断开终端会话', 1
    UNION ALL SELECT '/api/v1/cmdb/terminal/ws', 'GET', 'SSH终端WebSocket接入', 0
) t
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` sa
    WHERE sa.`path` = t.`path` AND sa.`method` = t.`method` AND sa.`deleted_at` IS NULL
);

INSERT INTO `sys_menu_api` (`sys_menu_id`, `sys_api_id`)
SELECT menu_item.`id`, api_item.`id`
FROM (
    SELECT 'cmdb:terminal:connect' AS `menu_permission`, '/api/v1/cmdb/terminal/sessions' AS `api_path`, 'POST' AS `api_method`
    UNION ALL SELECT 'cmdb:terminal:list', '/api/v1/cmdb/terminal/sessions', 'GET'
    UNION ALL SELECT 'cmdb:terminal:view', '/api/v1/cmdb/terminal/sessions/:id', 'GET'
    UNION ALL SELECT 'cmdb:terminal:audit', '/api/v1/cmdb/terminal/sessions/:id/logs', 'GET'
    UNION ALL SELECT 'cmdb:terminal:disconnect', '/api/v1/cmdb/terminal/sessions/:id/disconnect', 'POST'
    UNION ALL SELECT 'cmdb:terminal:force_disconnect', '/api/v1/cmdb/terminal/sessions/:id/force-disconnect', 'POST'
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
      'cmdb:terminal:list',
      'cmdb:terminal:connect',
      'cmdb:terminal:view',
      'cmdb:terminal:disconnect',
      'cmdb:terminal:force_disconnect',
      'cmdb:terminal:audit'
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
      'cmdb:terminal:list',
      'cmdb:terminal:connect',
      'cmdb:terminal:view',
      'cmdb:terminal:disconnect',
      'cmdb:terminal:force_disconnect',
      'cmdb:terminal:audit'
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

SELECT '第二阶段 CMDB SSH 在线终端表、菜单、按钮权限、API 元数据与 Casbin 策略已补齐，请重启服务或重新登录刷新权限缓存。' AS message;
