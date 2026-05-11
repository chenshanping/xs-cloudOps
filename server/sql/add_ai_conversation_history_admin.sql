-- 为 AI 对话历史（管理员视角）补充菜单、按钮权限、API 接口与角色授权
-- 执行方法: mysql -u root -p go-base < server/sql/add_ai_conversation_history_admin.sql

-- 1) 找到 AI 工具一级目录（type=1, permission='ai:tools'）
SET @ai_tools_menu_id = (
    SELECT id
    FROM sys_menu
    WHERE permission = 'ai:tools' AND type = 1
    ORDER BY id ASC
    LIMIT 1
);

-- 2) 新增 Type=2 页面菜单：AI对话历史
INSERT INTO sys_menu (parent_id, name, path, component, icon, sort, type, permission, status, hidden, created_at, updated_at)
SELECT @ai_tools_menu_id, '对话历史', '/ai/history', 'ai/history/index', 'history', 3, 2, 'ai:history:list', 1, 0, NOW(), NOW()
WHERE @ai_tools_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_menu WHERE permission = 'ai:history:list'
  );

SET @ai_history_menu_id = (
    SELECT id
    FROM sys_menu
    WHERE permission = 'ai:history:list' AND type = 2
    ORDER BY id ASC
    LIMIT 1
);

-- 3) 新增 Type=3 按钮：查看消息
INSERT INTO sys_menu (parent_id, name, path, component, icon, sort, type, permission, status, hidden, created_at, updated_at)
SELECT @ai_history_menu_id, '查看消息', '', '', '', 1, 3, 'ai:history:view', 1, 0, NOW(), NOW()
WHERE @ai_history_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_menu WHERE permission = 'ai:history:view'
  );

-- 4) 新增 Type=3 按钮：删除对话
INSERT INTO sys_menu (parent_id, name, path, component, icon, sort, type, permission, status, hidden, created_at, updated_at)
SELECT @ai_history_menu_id, '删除对话', '', '', '', 2, 3, 'ai:history:delete', 1, 0, NOW(), NOW()
WHERE @ai_history_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_menu WHERE permission = 'ai:history:delete'
  );

-- 5) 将三个菜单授权给 admin / system_admin 角色
INSERT INTO sys_role_menu (sys_role_id, sys_menu_id)
SELECT sr.id, sm.id
FROM sys_role sr
JOIN sys_menu sm ON sm.permission IN ('ai:history:list', 'ai:history:view', 'ai:history:delete')
WHERE sr.code IN ('admin', 'system_admin')
  AND NOT EXISTS (
      SELECT 1 FROM sys_role_menu target
      WHERE target.sys_role_id = sr.id
        AND target.sys_menu_id = sm.id
  );

-- 6) 新增四个 API 接口
INSERT INTO sys_api (path, method, `group`, description, request_params, response_params, need_auth, created_at, updated_at)
SELECT '/api/v1/ai/admin/users', 'GET', 'AI对话历史', 'AI活跃用户列表', '', '', 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_api WHERE path = '/api/v1/ai/admin/users' AND method = 'GET'
);

INSERT INTO sys_api (path, method, `group`, description, request_params, response_params, need_auth, created_at, updated_at)
SELECT '/api/v1/ai/admin/conversations', 'GET', 'AI对话历史', '对话历史列表', '', '', 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_api WHERE path = '/api/v1/ai/admin/conversations' AND method = 'GET'
);

INSERT INTO sys_api (path, method, `group`, description, request_params, response_params, need_auth, created_at, updated_at)
SELECT '/api/v1/ai/admin/conversations/:id/messages', 'GET', 'AI对话历史', '对话历史消息', '', '', 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_api WHERE path = '/api/v1/ai/admin/conversations/:id/messages' AND method = 'GET'
);

INSERT INTO sys_api (path, method, `group`, description, request_params, response_params, need_auth, created_at, updated_at)
SELECT '/api/v1/ai/admin/conversations/:id', 'DELETE', 'AI对话历史', '删除历史对话', '', '', 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_api WHERE path = '/api/v1/ai/admin/conversations/:id' AND method = 'DELETE'
);

-- 7) 将四个 API 授权给 admin / system_admin 角色
INSERT INTO sys_role_api (sys_role_id, sys_api_id)
SELECT sr.id, sa.id
FROM sys_role sr
JOIN sys_api sa ON (sa.path = '/api/v1/ai/admin/users' AND sa.method = 'GET')
                OR (sa.path = '/api/v1/ai/admin/conversations' AND sa.method = 'GET')
                OR (sa.path = '/api/v1/ai/admin/conversations/:id/messages' AND sa.method = 'GET')
                OR (sa.path = '/api/v1/ai/admin/conversations/:id' AND sa.method = 'DELETE')
WHERE sr.code IN ('admin', 'system_admin')
  AND NOT EXISTS (
      SELECT 1 FROM sys_role_api target
      WHERE target.sys_role_id = sr.id
        AND target.sys_api_id = sa.id
  );

-- 8) 同步写入 casbin 策略
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT DISTINCT 'p', sr.code, sa.path, sa.method, '', '', ''
FROM sys_role sr
JOIN sys_api sa ON (sa.path = '/api/v1/ai/admin/users' AND sa.method = 'GET')
                OR (sa.path = '/api/v1/ai/admin/conversations' AND sa.method = 'GET')
                OR (sa.path = '/api/v1/ai/admin/conversations/:id/messages' AND sa.method = 'GET')
                OR (sa.path = '/api/v1/ai/admin/conversations/:id' AND sa.method = 'DELETE')
WHERE sr.code IN ('admin', 'system_admin')
  AND NOT EXISTS (
      SELECT 1 FROM casbin_rule target
      WHERE target.ptype = 'p'
        AND target.v0 = sr.code
        AND target.v1 = sa.path
        AND target.v2 = sa.method
        AND IFNULL(target.v3, '') = ''
        AND IFNULL(target.v4, '') = ''
        AND IFNULL(target.v5, '') = ''
  );

SELECT 'AI 对话历史管理菜单/按钮/接口/授权补充完成，请重新登录刷新权限缓存。' AS message;
