-- AI 对话与 AI 配置：补齐按钮权限和批量删除对话 API 授权
-- 执行方法: mysql -u root -p go-base < server/sql/add_ai_chat_config_button_permissions.sql

SET NAMES utf8mb4;

-- 1) 定位 AI 页面菜单
SELECT @ai_chat_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'ai:chat:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

SELECT @ai_config_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'ai:config:list' AND `type` = 2 AND `deleted_at` IS NULL
LIMIT 1;

-- 2) 补齐 AI 对话按钮权限
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @ai_chat_menu_id, t.`name`, '', '', '', t.`sort`, 3, t.`permission`, 1, 0, NOW(3), NOW(3)
FROM (
    SELECT '新建对话' AS `name`, 1 AS `sort`, 'ai:chat:create' AS `permission`
    UNION ALL SELECT '发送对话', 2, 'ai:chat:send'
    UNION ALL SELECT '编辑标题', 3, 'ai:chat:update'
    UNION ALL SELECT '删除对话', 4, 'ai:chat:delete'
    UNION ALL SELECT '批量删除', 5, 'ai:chat:batchDelete'
    UNION ALL SELECT '清空上下文', 6, 'ai:chat:clearContext'
    UNION ALL SELECT '上传附件', 7, 'ai:chat:upload'
) t
WHERE @ai_chat_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` sm
      WHERE sm.`permission` = t.`permission` AND sm.`deleted_at` IS NULL
  );

-- 3) 补齐 AI 配置按钮权限
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @ai_config_menu_id, t.`name`, '', '', '', t.`sort`, 3, t.`permission`, 1, 0, NOW(3), NOW(3)
FROM (
    SELECT '新增平台' AS `name`, 1 AS `sort`, 'ai:config:createProvider' AS `permission`
    UNION ALL SELECT '编辑平台', 2, 'ai:config:editProvider'
    UNION ALL SELECT '删除平台', 3, 'ai:config:deleteProvider'
    UNION ALL SELECT '新增模型', 4, 'ai:config:createModel'
    UNION ALL SELECT '编辑模型', 5, 'ai:config:editModel'
    UNION ALL SELECT '删除模型', 6, 'ai:config:deleteModel'
    UNION ALL SELECT '导入模型', 7, 'ai:config:importModel'
    UNION ALL SELECT '测试模型', 8, 'ai:config:test'
    UNION ALL SELECT '保存配置', 9, 'ai:config:save'
) t
WHERE @ai_config_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu` sm
      WHERE sm.`permission` = t.`permission` AND sm.`deleted_at` IS NULL
  );

-- 4) 授权给已拥有对应页面权限的角色，避免既有管理员升级后按钮不可见
INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT srm.`sys_role_id`, button_menu.`id`
FROM `sys_role_menu` srm
JOIN `sys_menu` page_menu ON page_menu.`id` = srm.`sys_menu_id` AND page_menu.`deleted_at` IS NULL
JOIN `sys_menu` button_menu ON button_menu.`deleted_at` IS NULL
WHERE page_menu.`permission` = 'ai:chat:list'
  AND button_menu.`permission` IN (
      'ai:chat:create', 'ai:chat:send', 'ai:chat:update', 'ai:chat:delete',
      'ai:chat:batchDelete', 'ai:chat:clearContext', 'ai:chat:upload'
  )
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = srm.`sys_role_id`
        AND target.`sys_menu_id` = button_menu.`id`
  );

INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT srm.`sys_role_id`, button_menu.`id`
FROM `sys_role_menu` srm
JOIN `sys_menu` page_menu ON page_menu.`id` = srm.`sys_menu_id` AND page_menu.`deleted_at` IS NULL
JOIN `sys_menu` button_menu ON button_menu.`deleted_at` IS NULL
WHERE page_menu.`permission` = 'ai:config:list'
  AND button_menu.`permission` IN (
      'ai:config:createProvider', 'ai:config:editProvider', 'ai:config:deleteProvider',
      'ai:config:createModel', 'ai:config:editModel', 'ai:config:deleteModel',
      'ai:config:importModel', 'ai:config:test', 'ai:config:save'
  )
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` target
      WHERE target.`sys_role_id` = srm.`sys_role_id`
        AND target.`sys_menu_id` = button_menu.`id`
  );

-- 5) 补齐菜单/API 绑定所需 API 元数据
INSERT INTO `sys_api` (`path`, `method`, `group`, `description`, `request_params`, `response_params`, `need_auth`, `created_at`, `updated_at`)
SELECT t.`path`, t.`method`, t.`api_group`, t.`description`, '', '', 1, NOW(3), NOW(3)
FROM (
    SELECT '/api/v1/ai/conversations' AS `path`, 'GET' AS `method`, 'AI对话' AS `api_group`, '获取对话列表' AS `description`
    UNION ALL SELECT '/api/v1/ai/conversations', 'POST', 'AI对话', '创建对话'
    UNION ALL SELECT '/api/v1/ai/conversations/batch', 'DELETE', 'AI对话', '批量删除对话'
    UNION ALL SELECT '/api/v1/ai/conversations/:id', 'GET', 'AI对话', '获取对话详情'
    UNION ALL SELECT '/api/v1/ai/conversations/:id', 'PUT', 'AI对话', '更新对话标题'
    UNION ALL SELECT '/api/v1/ai/conversations/:id', 'DELETE', 'AI对话', '删除对话'
    UNION ALL SELECT '/api/v1/ai/conversations/:id/messages', 'GET', 'AI对话', '获取对话消息'
    UNION ALL SELECT '/api/v1/ai/conversations/:id/clear-context', 'POST', 'AI对话', '清空上下文'
    UNION ALL SELECT '/api/v1/ai/chat', 'POST', 'AI对话', 'AI对话'
    UNION ALL SELECT '/api/v1/ai/chat/stream', 'POST', 'AI对话', 'AI流式对话'
    UNION ALL SELECT '/api/v1/ai/config', 'GET', 'AI配置', '获取AI配置'
    UNION ALL SELECT '/api/v1/ai/config', 'PUT', 'AI配置', '保存AI配置'
    UNION ALL SELECT '/api/v1/ai/test', 'POST', 'AI配置', '测试AI配置'
    UNION ALL SELECT '/api/v1/ai/providers/models/fetch', 'POST', 'AI配置', '拉取平台模型列表'
    UNION ALL SELECT '/api/v1/ai/admin/users', 'GET', 'AI对话历史', 'AI活跃用户列表'
    UNION ALL SELECT '/api/v1/ai/admin/conversations', 'GET', 'AI对话历史', '对话历史列表'
    UNION ALL SELECT '/api/v1/ai/admin/conversations/:id/messages', 'GET', 'AI对话历史', '对话历史消息'
    UNION ALL SELECT '/api/v1/ai/admin/conversations/:id', 'DELETE', 'AI对话历史', '删除历史对话'
    UNION ALL SELECT '/api/v1/files/credential', 'POST', '文件管理', '获取上传凭证'
    UNION ALL SELECT '/api/v1/files/check-md5', 'POST', '文件管理', 'MD5秒传检查'
    UNION ALL SELECT '/api/v1/files/save', 'POST', '文件管理', '保存上传文件'
    UNION ALL SELECT '/api/v1/files/multipart/init', 'POST', '文件管理', '初始化分片上传'
    UNION ALL SELECT '/api/v1/files/multipart/parts', 'GET', '文件管理', '获取已上传分片'
    UNION ALL SELECT '/api/v1/files/multipart/complete', 'POST', '文件管理', '完成分片上传'
    UNION ALL SELECT '/api/v1/files/multipart/abort', 'POST', '文件管理', '取消分片上传'
    UNION ALL SELECT '/api/v1/files/upload/local', 'POST', '文件管理', '本地文件上传'
    UNION ALL SELECT '/api/v1/files/upload/chunk', 'POST', '文件管理', '上传分片'
) t
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_api` sa
    WHERE sa.`path` = t.`path` AND sa.`method` = t.`method` AND sa.`deleted_at` IS NULL
);

UPDATE `sys_api` sa
JOIN (
    SELECT '/api/v1/ai/conversations' AS `path`, 'GET' AS `method`, 'AI对话' AS `api_group`, '获取对话列表' AS `description`
    UNION ALL SELECT '/api/v1/ai/conversations', 'POST', 'AI对话', '创建对话'
    UNION ALL SELECT '/api/v1/ai/conversations/batch', 'DELETE', 'AI对话', '批量删除对话'
    UNION ALL SELECT '/api/v1/ai/conversations/:id', 'GET', 'AI对话', '获取对话详情'
    UNION ALL SELECT '/api/v1/ai/conversations/:id', 'PUT', 'AI对话', '更新对话标题'
    UNION ALL SELECT '/api/v1/ai/conversations/:id', 'DELETE', 'AI对话', '删除对话'
    UNION ALL SELECT '/api/v1/ai/conversations/:id/messages', 'GET', 'AI对话', '获取对话消息'
    UNION ALL SELECT '/api/v1/ai/conversations/:id/clear-context', 'POST', 'AI对话', '清空上下文'
    UNION ALL SELECT '/api/v1/ai/chat', 'POST', 'AI对话', 'AI对话'
    UNION ALL SELECT '/api/v1/ai/chat/stream', 'POST', 'AI对话', 'AI流式对话'
    UNION ALL SELECT '/api/v1/ai/config', 'GET', 'AI配置', '获取AI配置'
    UNION ALL SELECT '/api/v1/ai/config', 'PUT', 'AI配置', '保存AI配置'
    UNION ALL SELECT '/api/v1/ai/test', 'POST', 'AI配置', '测试AI配置'
    UNION ALL SELECT '/api/v1/ai/providers/models/fetch', 'POST', 'AI配置', '拉取平台模型列表'
    UNION ALL SELECT '/api/v1/ai/admin/users', 'GET', 'AI对话历史', 'AI活跃用户列表'
    UNION ALL SELECT '/api/v1/ai/admin/conversations', 'GET', 'AI对话历史', '对话历史列表'
    UNION ALL SELECT '/api/v1/ai/admin/conversations/:id/messages', 'GET', 'AI对话历史', '对话历史消息'
    UNION ALL SELECT '/api/v1/ai/admin/conversations/:id', 'DELETE', 'AI对话历史', '删除历史对话'
    UNION ALL SELECT '/api/v1/files/credential', 'POST', '文件管理', '获取上传凭证'
    UNION ALL SELECT '/api/v1/files/check-md5', 'POST', '文件管理', 'MD5秒传检查'
    UNION ALL SELECT '/api/v1/files/save', 'POST', '文件管理', '保存上传文件'
    UNION ALL SELECT '/api/v1/files/multipart/init', 'POST', '文件管理', '初始化分片上传'
    UNION ALL SELECT '/api/v1/files/multipart/parts', 'GET', '文件管理', '获取已上传分片'
    UNION ALL SELECT '/api/v1/files/multipart/complete', 'POST', '文件管理', '完成分片上传'
    UNION ALL SELECT '/api/v1/files/multipart/abort', 'POST', '文件管理', '取消分片上传'
    UNION ALL SELECT '/api/v1/files/upload/local', 'POST', '文件管理', '本地文件上传'
    UNION ALL SELECT '/api/v1/files/upload/chunk', 'POST', '文件管理', '上传分片'
) t ON t.`path` = sa.`path` AND t.`method` = sa.`method`
SET sa.`group` = t.`api_group`,
    sa.`description` = t.`description`,
    sa.`need_auth` = 1,
    sa.`updated_at` = NOW(3)
WHERE sa.`deleted_at` IS NULL;

-- 6) 按方案 C 建立菜单/按钮到 API 的绑定：角色只勾菜单按钮即可继承后端 API 权限
INSERT INTO `sys_menu_api` (`sys_menu_id`, `sys_api_id`)
SELECT menu_item.`id`, api_item.`id`
FROM (
    SELECT 'ai:chat:list' AS `menu_permission`, '/api/v1/ai/conversations' AS `api_path`, 'GET' AS `api_method`
    UNION ALL SELECT 'ai:chat:list', '/api/v1/ai/conversations/:id', 'GET'
    UNION ALL SELECT 'ai:chat:list', '/api/v1/ai/conversations/:id/messages', 'GET'
    UNION ALL SELECT 'ai:chat:create', '/api/v1/ai/conversations', 'POST'
    UNION ALL SELECT 'ai:chat:send', '/api/v1/ai/chat', 'POST'
    UNION ALL SELECT 'ai:chat:send', '/api/v1/ai/chat/stream', 'POST'
    UNION ALL SELECT 'ai:chat:update', '/api/v1/ai/conversations/:id', 'PUT'
    UNION ALL SELECT 'ai:chat:delete', '/api/v1/ai/conversations/:id', 'DELETE'
    UNION ALL SELECT 'ai:chat:batchDelete', '/api/v1/ai/conversations/batch', 'DELETE'
    UNION ALL SELECT 'ai:chat:clearContext', '/api/v1/ai/conversations/:id/clear-context', 'POST'
    UNION ALL SELECT 'ai:chat:upload', '/api/v1/files/credential', 'POST'
    UNION ALL SELECT 'ai:chat:upload', '/api/v1/files/check-md5', 'POST'
    UNION ALL SELECT 'ai:chat:upload', '/api/v1/files/save', 'POST'
    UNION ALL SELECT 'ai:chat:upload', '/api/v1/files/multipart/init', 'POST'
    UNION ALL SELECT 'ai:chat:upload', '/api/v1/files/multipart/parts', 'GET'
    UNION ALL SELECT 'ai:chat:upload', '/api/v1/files/multipart/complete', 'POST'
    UNION ALL SELECT 'ai:chat:upload', '/api/v1/files/multipart/abort', 'POST'
    UNION ALL SELECT 'ai:chat:upload', '/api/v1/files/upload/local', 'POST'
    UNION ALL SELECT 'ai:chat:upload', '/api/v1/files/upload/chunk', 'POST'
    UNION ALL SELECT 'ai:config:list', '/api/v1/ai/config', 'GET'
    UNION ALL SELECT 'ai:config:save', '/api/v1/ai/config', 'PUT'
    UNION ALL SELECT 'ai:config:test', '/api/v1/ai/test', 'POST'
    UNION ALL SELECT 'ai:config:importModel', '/api/v1/ai/providers/models/fetch', 'POST'
    UNION ALL SELECT 'ai:history:list', '/api/v1/ai/admin/users', 'GET'
    UNION ALL SELECT 'ai:history:list', '/api/v1/ai/admin/conversations', 'GET'
    UNION ALL SELECT 'ai:history:view', '/api/v1/ai/admin/conversations/:id/messages', 'GET'
    UNION ALL SELECT 'ai:history:delete', '/api/v1/ai/admin/conversations/:id', 'DELETE'
) t
JOIN `sys_menu` menu_item ON menu_item.`permission` = t.`menu_permission` AND menu_item.`deleted_at` IS NULL
JOIN `sys_api` api_item ON api_item.`path` = t.`api_path` AND api_item.`method` = t.`api_method` AND api_item.`deleted_at` IS NULL
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_menu_api` existing
    WHERE existing.`sys_menu_id` = menu_item.`id`
      AND existing.`sys_api_id` = api_item.`id`
);

-- 7) 同步已拥有菜单按钮的角色到 Casbin 策略表；应用运行中需重载策略或重启后生效
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
SELECT DISTINCT 'p', sr.`code`, api_item.`path`, api_item.`method`, '', '', ''
FROM `sys_role_menu` srm
JOIN `sys_role` sr ON sr.`id` = srm.`sys_role_id` AND sr.`deleted_at` IS NULL
JOIN `sys_menu_api` sma ON sma.`sys_menu_id` = srm.`sys_menu_id`
JOIN `sys_api` api_item ON api_item.`id` = sma.`sys_api_id` AND api_item.`deleted_at` IS NULL
WHERE NOT EXISTS (
    SELECT 1 FROM `casbin_rule` target
    WHERE target.`ptype` = 'p'
      AND target.`v0` = sr.`code`
      AND target.`v1` = api_item.`path`
      AND target.`v2` = api_item.`method`
      AND IFNULL(target.`v3`, '') = ''
      AND IFNULL(target.`v4`, '') = ''
      AND IFNULL(target.`v5`, '') = ''
);

SELECT 'AI 对话与 AI 配置按钮权限补齐完成；请重启服务或重载 Casbin 策略，并重新登录刷新权限缓存。' AS message;
