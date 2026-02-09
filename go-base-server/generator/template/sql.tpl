-- {{.Description}} 建表SQL
-- 生成时间: {{.GenerateTime}}
-- 模块: {{.ModuleName}}

CREATE TABLE IF NOT EXISTS `{{.TableName}}` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
{{- range .Columns}}
{{- if and (ne .ColumnName "id") (ne .ColumnName "created_at") (ne .ColumnName "updated_at") (ne .ColumnName "deleted_at")}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
  `{{.ColumnName}}_file_id` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '{{.Comment}}文件ID',
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
  `{{.ColumnName}}_file_ids` VARCHAR(500) NULL DEFAULT '' COMMENT '{{.Comment}}文件ID列表',
{{- else}}
  `{{.ColumnName}}` {{.SqlType}} {{if .IsRequired}}NOT NULL{{else}}NULL{{end}}{{if .DefaultValue}} DEFAULT {{.DefaultValue}}{{end}} COMMENT '{{.Comment}}',
{{- end}}
{{- end}}
{{- end}}
{{- if .LinkToUser}}
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '关联用户ID',
{{- end}}
{{- if .HasCreatedBy}}
  `created_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '创建人 ID',
{{- end}}
{{- if .HasAudit}}
  `audit_status` TINYINT NULL DEFAULT 0 COMMENT '审批状态 0-待审批 1-审批通过 2-审批拒绝',
  `audit_remark` VARCHAR(500) NULL DEFAULT '' COMMENT '审批备注',
  `audit_time` DATETIME NULL COMMENT '审批时间',
  `audit_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '审批人ID',
{{- end}}
{{- if .HasCreatedAt}}
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
{{- end}}
{{- if .HasUpdatedAt}}
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
{{- end}}
{{- if .HasDeletedAt}}
  `deleted_at` DATETIME NULL COMMENT '删除时间',
{{- end}}
  PRIMARY KEY (`id`)
{{- if .LinkToUser}},
  UNIQUE INDEX `uk_user_id` (`user_id`)
{{- end}}
{{- range .UniqueColumns}},
  UNIQUE INDEX `uk_{{.ColumnName}}` (`{{.ColumnName}}`)
{{- end}}
{{- if .HasDeletedAt}},
  INDEX `idx_deleted_at` (`deleted_at`)
{{- end}}
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='{{.Description}}';
