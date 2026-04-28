-- 为用户管理补充默认密码配置、批量重置密码按钮权限和批量重置密码接口授权
-- 执行方法: mysql -u root -p go-base < server/sql/add_user_password_reset_controls.sql

INSERT INTO sys_config (name, `key`, `value`, value_type, remark, created_at, updated_at)
SELECT '用户默认密码', 'user_default_password', '123456', 'string', '后台用户管理单条/批量重置密码默认值', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_config WHERE `key` = 'user_default_password'
);

SET @user_menu_id = (
    SELECT id
    FROM sys_menu
    WHERE permission = 'system:user:list' AND type = 2
    ORDER BY id ASC
    LIMIT 1
);

INSERT INTO sys_menu (parent_id, name, path, component, icon, sort, type, permission, status, hidden, created_at, updated_at)
SELECT @user_menu_id, '用户重置密码', '', '', '', 4, 3, 'system:user:resetPwd', 1, 0, NOW(), NOW()
WHERE @user_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_menu WHERE permission = 'system:user:resetPwd'
  );

INSERT INTO sys_menu (parent_id, name, path, component, icon, sort, type, permission, status, hidden, created_at, updated_at)
SELECT @user_menu_id, '批量重置密码', '', '', '', 7, 3, 'system:user:batchResetPwd', 1, 0, NOW(), NOW()
WHERE @user_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_menu WHERE permission = 'system:user:batchResetPwd'
  );

SET @single_reset_menu_id = (
    SELECT id
    FROM sys_menu
    WHERE permission = 'system:user:resetPwd'
    ORDER BY id ASC
    LIMIT 1
);

SET @batch_reset_menu_id = (
    SELECT id
    FROM sys_menu
    WHERE permission = 'system:user:batchResetPwd'
    ORDER BY id ASC
    LIMIT 1
);

INSERT INTO sys_role_menu (sys_role_id, sys_menu_id)
SELECT DISTINCT srm.sys_role_id, @single_reset_menu_id
FROM sys_role_menu srm
JOIN sys_menu sm ON sm.id = srm.sys_menu_id
WHERE sm.permission = 'system:user:edit'
  AND @single_reset_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_role_menu target
      WHERE target.sys_role_id = srm.sys_role_id
        AND target.sys_menu_id = @single_reset_menu_id
  );

INSERT INTO sys_role_menu (sys_role_id, sys_menu_id)
SELECT DISTINCT srm.sys_role_id, @batch_reset_menu_id
FROM sys_role_menu srm
JOIN sys_menu sm ON sm.id = srm.sys_menu_id
WHERE sm.permission = 'system:user:resetPwd'
  AND @batch_reset_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_role_menu target
      WHERE target.sys_role_id = srm.sys_role_id
        AND target.sys_menu_id = @batch_reset_menu_id
  );

INSERT INTO sys_api (path, method, `group`, description, request_params, response_params, need_auth, created_at, updated_at)
SELECT '/api/v1/users/batch-password', 'PUT', '用户管理', '批量重置密码', '', '', 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_api WHERE path = '/api/v1/users/batch-password' AND method = 'PUT'
);

SET @batch_reset_api_id = (
    SELECT id
    FROM sys_api
    WHERE path = '/api/v1/users/batch-password' AND method = 'PUT'
    ORDER BY id ASC
    LIMIT 1
);

INSERT INTO sys_role_api (sys_role_id, sys_api_id)
SELECT DISTINCT sra.sys_role_id, @batch_reset_api_id
FROM sys_role_api sra
JOIN sys_api sa ON sa.id = sra.sys_api_id
WHERE sa.path = '/api/v1/users/:id/password' AND sa.method = 'PUT'
  AND @batch_reset_api_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_role_api target
      WHERE target.sys_role_id = sra.sys_role_id
        AND target.sys_api_id = @batch_reset_api_id
  );

INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT DISTINCT 'p', cr.v0, '/api/v1/users/batch-password', 'PUT', '', '', ''
FROM casbin_rule cr
WHERE cr.ptype = 'p'
  AND cr.v1 = '/api/v1/users/:id/password'
  AND cr.v2 = 'PUT'
  AND NOT EXISTS (
      SELECT 1 FROM casbin_rule target
      WHERE target.ptype = 'p'
        AND target.v0 = cr.v0
        AND target.v1 = '/api/v1/users/batch-password'
        AND target.v2 = 'PUT'
        AND IFNULL(target.v3, '') = ''
        AND IFNULL(target.v4, '') = ''
        AND IFNULL(target.v5, '') = ''
  );

SELECT '用户默认密码与批量重置密码权限补充完成，请重新登录刷新权限缓存。' AS message;
