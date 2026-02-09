
-- 养殖户认证 角色SQL
-- 生成时间: 2026-02-07 17:35:20
-- 模块: farmer_certification
-- 说明: 该角色用于限定哪些用户可以填写养殖户信息

-- 检查角色是否存在，不存在则创建
INSERT INTO `sys_role` (`name`, `code`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
SELECT '养殖户', 'farmer_certification', 0, 1, '养殖户认证角色', NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_role` WHERE `code` = 'farmer_certification' AND `deleted_at` IS NULL
);
