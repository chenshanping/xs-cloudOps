-- 为productType模块补充导出导入权限
-- 执行方法: mysql -u root -p ecobreed < go-base-server/sql/fix_productType_permissions.sql

-- 查询productType主菜单ID
SET @parent_id = (SELECT id FROM sys_menu WHERE permission = 'product_type:list' AND type = 2 LIMIT 1);

-- 检查是否找到主菜单
SELECT CASE 
    WHEN @parent_id IS NULL THEN '错误: 找不到product_type主菜单，请先生成代码'
    ELSE CONCAT('找到主菜单ID: ', @parent_id)
END AS status;

-- 如果主菜单存在，添加导出权限（如果不存在）
INSERT INTO sys_menu (parent_id, name, path, component, icon, sort, type, permission, status, hidden, created_at, updated_at)
SELECT @parent_id, '导出', '', '', '', 5, 3, 'product_type:export', 1, 0, NOW(), NOW()
WHERE @parent_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_menu 
      WHERE parent_id = @parent_id 
      AND permission = 'product_type:export'
  );

-- 添加导入权限（如果不存在）
INSERT INTO sys_menu (parent_id, name, path, component, icon, sort, type, permission, status, hidden, created_at, updated_at)
SELECT @parent_id, '导入', '', '', '', 6, 3, 'product_type:import', 1, 0, NOW(), NOW()
WHERE @parent_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_menu 
      WHERE parent_id = @parent_id 
      AND permission = 'product_type:import'
  );

-- 分配权限给管理员角色（角色ID=1）
INSERT INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu 
WHERE permission IN ('product_type:export', 'product_type:import')
  AND id NOT IN (SELECT menu_id FROM sys_role_menu WHERE role_id = 1)
ON DUPLICATE KEY UPDATE role_id=role_id;

-- 验证结果
SELECT 
    '权限补充完成！' AS message,
    (SELECT COUNT(*) FROM sys_menu WHERE permission LIKE 'product_type:%') AS total_permissions,
    (SELECT COUNT(*) FROM sys_menu WHERE permission IN ('product_type:export', 'product_type:import')) AS export_import_count;

-- 显示所有权限
SELECT id, parent_id, name, permission, type, sort 
FROM sys_menu 
WHERE permission LIKE 'product_type%'
ORDER BY parent_id, sort;

-- 提示
SELECT '请重新登录系统以刷新权限缓存！' AS next_step;
