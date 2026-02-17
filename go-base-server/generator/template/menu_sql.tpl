{{- if .HasMenu}}
-- {{.Description}} 菜单SQL
-- 生成时间: {{.GenerateTime}}
-- 模块: {{.ModuleName}}
-- 说明: 请根据实际情况修改父菜单ID和排序号

-- 目录/菜单 (type: 1=目录, 2=菜单, 3=按钮)
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`) VALUES
({{.MenuConfig.ParentID}}, '{{.MenuConfig.MenuName}}', '/{{.ModuleName}}', '{{.ComponentPath}}', '{{.MenuConfig.MenuIcon}}', {{.MenuConfig.MenuSort}}, 2, '{{.MenuConfig.Permission}}:list', 1, 0, NOW(), NOW());

-- 获取刚插入的菜单ID (用于插入按钮权限)
SET @menu_id = LAST_INSERT_ID();

-- 按钮权限
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `sort`, `type`, `permission`, `status`, `hidden`, `created_at`, `updated_at`) VALUES
(@menu_id, '查看', '', '', '', 1, 3, '{{.MenuConfig.Permission}}:list', 1, 0, NOW(), NOW()),
(@menu_id, '新增', '', '', '', 2, 3, '{{.MenuConfig.Permission}}:add', 1, 0, NOW(), NOW()),
(@menu_id, '编辑', '', '', '', 3, 3, '{{.MenuConfig.Permission}}:edit', 1, 0, NOW(), NOW()),
(@menu_id, '删除', '', '', '', 4, 3, '{{.MenuConfig.Permission}}:delete', 1, 0, NOW(), NOW()){{if .EnableImportExport}},
(@menu_id, '导出', '', '', '', 5, 3, '{{.MenuConfig.Permission}}:export', 1, 0, NOW(), NOW()),
(@menu_id, '导入', '', '', '', 6, 3, '{{.MenuConfig.Permission}}:import', 1, 0, NOW(), NOW()){{end}}{{if .HasAudit}},
(@menu_id, '审批', '', '', '', 7, 3, '{{.MenuConfig.Permission}}:audit', 1, 0, NOW(), NOW()){{end}};
{{- end}}
