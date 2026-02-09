-- {{.Description}} 关联中间表SQL
-- 生成时间: {{.GenerateTime}}
-- 模块: {{.ModuleName}}

{{range .Relations}}
{{if eq .RelationType "many2many"}}
-- {{.Comment}}与{{$.Description}}关联中间表
CREATE TABLE IF NOT EXISTS `{{.JoinTable}}` (
  `{{$.TableName}}_id` BIGINT UNSIGNED NOT NULL COMMENT '{{$.Description}}ID',
  `{{.RelatedTable}}_id` BIGINT UNSIGNED NOT NULL COMMENT '{{.Comment}}ID',
  PRIMARY KEY (`{{$.TableName}}_id`, `{{.RelatedTable}}_id`),
  INDEX `idx_{{$.TableName}}_id` (`{{$.TableName}}_id`),
  INDEX `idx_{{.RelatedTable}}_id` (`{{.RelatedTable}}_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='{{$.Description}}-{{.Comment}}关联表';

{{end}}
{{end}}
