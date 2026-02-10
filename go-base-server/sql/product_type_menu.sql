
-- 产品类型 菜单SQL
-- 生成时间: 2026-02-10 22:34:58
-- 模块: product_type
-- 说明: 请根据实际情况修改父菜单ID和排序号

-- 目录/菜单 (type: 1=目录, 2=菜单, 3=按钮)
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`) VALUES
(0, '产品类型', '/product_type', 'product_type/index', '', 0, 2, 'product_type:list', 1, 0, NOW(), NOW());

-- 获取刚插入的菜单ID (用于插入按钮权限)
SET @menu_id = LAST_INSERT_ID();

-- 按钮权限
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`) VALUES
(@menu_id, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0, NOW(), NOW()),
(@menu_id, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0, NOW(), NOW()),
(@menu_id, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0, NOW(), NOW()),
(@menu_id, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0, NOW(), NOW());
