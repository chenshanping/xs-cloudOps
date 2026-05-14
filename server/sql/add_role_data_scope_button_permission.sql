-- 角色管理：新增独立的"数据权限"按钮（system:role:dataScope）
-- 背景：原"分配权限"抽屉里的"数据权限" tab 抽离为独立按钮，按钮需单独可控
-- 执行方法: mysql -u root -p go-base < server/sql/add_role_data_scope_button_permission.sql

SET NAMES utf8mb4;

-- 1) 定位角色管理菜单（系统管理 → 角色管理）
SELECT @role_menu_id := `id`
FROM `sys_menu`
WHERE `permission` = 'system:role:list' AND `deleted_at` IS NULL
LIMIT 1;

-- 2) 插入按钮权限（幂等：按 permission 去重）
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`)
SELECT @role_menu_id, '数据权限', '', '', '', 5, 3, 'system:role:dataScope', 1, 0, NOW(3), NOW(3)
WHERE @role_menu_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM `sys_menu`
      WHERE `permission` = 'system:role:dataScope' AND `deleted_at` IS NULL
  );

-- 3) 将新按钮授权给所有原本就拥有"分配权限"的角色，保证既有管理员开箱可见
INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT DISTINCT srm.`sys_role_id`, target_menu.`id`
FROM `sys_role_menu` srm
JOIN `sys_menu` sm ON sm.`id` = srm.`sys_menu_id` AND sm.`deleted_at` IS NULL
JOIN `sys_menu` target_menu ON target_menu.`permission` = 'system:role:dataScope' AND target_menu.`deleted_at` IS NULL
WHERE sm.`permission` = 'system:role:assign'
  AND NOT EXISTS (
      SELECT 1 FROM `sys_role_menu` t
      WHERE t.`sys_role_id` = srm.`sys_role_id`
        AND t.`sys_menu_id` = target_menu.`id`
  );
