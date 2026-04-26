-- 为用户管理补充批量启用/批量禁用按钮权限和批量状态接口授权
-- 执行方法: mysql -u root -p go-base < server/sql/fix_user_batch_status_permissions.sql

SET @user_menu_id = (
    SELECT id
    FROM sys_menu
    WHERE permission = 'system:user:list' AND type = 2
    ORDER BY id ASC
    LIMIT 1
);

-- 新增用户管理按钮: 批量启用 / 批量禁用
INSERT INTO sys_menu (parent_id, name, path, component, icon, sort, type, permission, status, hidden, created_at, updated_at)
SELECT @user_menu_id, '批量启用', '', '', '', 5, 3, 'system:user:batchEnable', 1, 0, NOW(), NOW()
WHERE @user_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_menu WHERE permission = 'system:user:batchEnable'
  );

INSERT INTO sys_menu (parent_id, name, path, component, icon, sort, type, permission, status, hidden, created_at, updated_at)
SELECT @user_menu_id, '批量禁用', '', '', '', 6, 3, 'system:user:batchDisable', 1, 0, NOW(), NOW()
WHERE @user_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_menu WHERE permission = 'system:user:batchDisable'
  );

SET @batch_status_api_id = (
    SELECT id
    FROM sys_api
    WHERE path = '/api/v1/users/batch-status' AND method = 'PUT'
    ORDER BY id ASC
    LIMIT 1
);

-- 新增批量状态 API
INSERT INTO sys_api (path, method, `group`, description, request_params, response_params, need_auth, created_at, updated_at)
SELECT '/api/v1/users/batch-status', 'PUT', '用户管理', '批量修改用户状态', '', '', 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_api WHERE path = '/api/v1/users/batch-status' AND method = 'PUT'
);

SET @batch_status_api_id = (
    SELECT id
    FROM sys_api
    WHERE path = '/api/v1/users/batch-status' AND method = 'PUT'
    ORDER BY id ASC
    LIMIT 1
);

SET @batch_enable_menu_id = (
    SELECT id
    FROM sys_menu
    WHERE permission = 'system:user:batchEnable'
    ORDER BY id ASC
    LIMIT 1
);

SET @batch_disable_menu_id = (
    SELECT id
    FROM sys_menu
    WHERE permission = 'system:user:batchDisable'
    ORDER BY id ASC
    LIMIT 1
);

-- 将新按钮授权给所有原本就有 system:user:edit 的角色
INSERT INTO sys_role_menu (sys_role_id, sys_menu_id)
SELECT DISTINCT srm.sys_role_id, @batch_enable_menu_id
FROM sys_role_menu srm
JOIN sys_menu sm ON sm.id = srm.sys_menu_id
WHERE sm.permission = 'system:user:edit'
  AND @batch_enable_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_role_menu target
      WHERE target.sys_role_id = srm.sys_role_id
        AND target.sys_menu_id = @batch_enable_menu_id
  );

INSERT INTO sys_role_menu (sys_role_id, sys_menu_id)
SELECT DISTINCT srm.sys_role_id, @batch_disable_menu_id
FROM sys_role_menu srm
JOIN sys_menu sm ON sm.id = srm.sys_menu_id
WHERE sm.permission = 'system:user:edit'
  AND @batch_disable_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_role_menu target
      WHERE target.sys_role_id = srm.sys_role_id
        AND target.sys_menu_id = @batch_disable_menu_id
  );

-- 将批量状态 API 授权给所有原本就有单用户状态修改 API 的角色
INSERT INTO sys_role_api (sys_role_id, sys_api_id)
SELECT DISTINCT sra.sys_role_id, @batch_status_api_id
FROM sys_role_api sra
JOIN sys_api sa ON sa.id = sra.sys_api_id
WHERE sa.path = '/api/v1/users/:id/status' AND sa.method = 'PUT'
  AND @batch_status_api_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_role_api target
      WHERE target.sys_role_id = sra.sys_role_id
        AND target.sys_api_id = @batch_status_api_id
  );

-- casbin_rule 基于已有单用户状态修改接口授权角色补齐
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT DISTINCT 'p', cr.v0, '/api/v1/users/batch-status', 'PUT', '', '', ''
FROM casbin_rule cr
WHERE cr.ptype = 'p'
  AND cr.v1 = '/api/v1/users/:id/status'
  AND cr.v2 = 'PUT'
  AND NOT EXISTS (
      SELECT 1 FROM casbin_rule target
      WHERE target.ptype = 'p'
        AND target.v0 = cr.v0
        AND target.v1 = '/api/v1/users/batch-status'
        AND target.v2 = 'PUT'
        AND IFNULL(target.v3, '') = ''
        AND IFNULL(target.v4, '') = ''
        AND IFNULL(target.v5, '') = ''
  );

SELECT '用户批量状态权限补充完成，请重新登录刷新权限缓存。' AS message;
