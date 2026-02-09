{{- if .LinkToUser}}
-- {{.Description}} 角色SQL
-- 生成时间: {{.GenerateTime}}
-- 模块: {{.ModuleName}}
-- 说明: 该角色用于限定哪些用户可以填写{{.ProfileName}}信息

-- 检查角色是否存在，不存在则创建
INSERT INTO `sys_role` (`name`, `code`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
SELECT '{{.ProfileName}}', '{{.ProfileRoleCode}}', 0, 1, '{{.Description}}角色', NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (
    SELECT 1 FROM `sys_role` WHERE `code` = '{{.ProfileRoleCode}}' AND `deleted_at` IS NULL
);
{{- end}}
