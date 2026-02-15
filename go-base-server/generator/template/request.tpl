package request

// {{.ModelName}}QueryRequest {{.Description}}查询请求（用于导出）
type {{.ModelName}}QueryRequest struct {
{{- range .SearchColumns}}
	{{.FieldName}} {{if or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64")}}*{{.FieldType}}{{else if eq .SearchType "eq"}}*{{.FieldType}}{{else}}string{{end}} `json:"{{.JsonName}}" form:"{{.JsonName}}"{{if .Comment}} comment:"{{.Comment}}"{{end}}`
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
	{{.ForeignKey | ToPascalCase}} *uint `json:"{{.ForeignKeyJson}}" form:"{{.ForeignKeyJson}}" comment:"{{if .Comment}}{{.Comment}}ID{{else}}{{.FieldName}}ID{{end}}"`
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
	CreatedBy *uint `json:"created_by" form:"created_by" comment:"创建人 ID"`
{{- end}}
}

// {{.ModelName}}ListRequest {{.Description}}列表请求
type {{.ModelName}}ListRequest struct {
	PageRequest
{{- range .SearchColumns}}
	{{.FieldName}} {{if or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64")}}*{{.FieldType}}{{else if eq .SearchType "eq"}}*{{.FieldType}}{{else}}string{{end}} `json:"{{.JsonName}}" form:"{{.JsonName}}"{{if .Comment}} comment:"{{.Comment}}"{{end}}`
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
	{{.ForeignKey | ToPascalCase}} *uint `json:"{{.ForeignKeyJson}}" form:"{{.ForeignKeyJson}}" comment:"{{if .Comment}}{{.Comment}}ID{{else}}{{.FieldName}}ID{{end}}"`
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
	CreatedBy *uint `json:"created_by" form:"created_by" comment:"创建人 ID"`
{{- end}}
	SortField string `json:"sort_field" form:"sort_field" comment:"排序字段"`
	SortOrder string `json:"sort_order" form:"sort_order" comment:"排序方式 asc/desc"`
}

// Create{{.ModelName}}Request 创建{{.Description}}请求
type Create{{.ModelName}}Request struct {
{{- if .LinkToUser}}
	UserID uint `json:"user_id" binding:"required" comment:"关联用户ID"`
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	{{.FieldName}}FileID uint `json:"{{.JsonName}}_file_id" comment:"{{if .Comment}}{{.Comment}}文件ID{{else}}{{.FieldName}}文件ID{{end}}"`
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
	{{.FieldName}}FileIDs string `json:"{{.JsonName}}_file_ids" comment:"{{if .Comment}}{{.Comment}}文件ID列表{{else}}{{.FieldName}}文件ID列表{{end}}"`
{{- else}}
	{{.FieldName}} {{if and .IsRequired (or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64"))}}*{{end}}{{.FieldType}} `json:"{{.JsonName}}"{{if .IsRequired}} binding:"required"{{end}}{{if .Comment}} comment:"{{.Comment}}"{{end}}`
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
	{{.ForeignKey | ToPascalCase}} uint `json:"{{.ForeignKeyJson}}" comment:"{{if .Comment}}{{.Comment}}ID{{else}}{{.FieldName}}ID{{end}}"`
{{- else if eq .RelationType "many2many"}}
	{{.FieldName}}Ids []uint `json:"{{.JsonName}}_ids" comment:"{{.FieldName}}ID列表"`
{{- end}}
{{- end}}
}

// Update{{.ModelName}}Request 更新{{.Description}}请求
type Update{{.ModelName}}Request struct {
{{- if .LinkToUser}}
	UserID uint `json:"user_id" comment:"关联用户ID"`
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	{{.FieldName}}FileID uint `json:"{{.JsonName}}_file_id" comment:"{{if .Comment}}{{.Comment}}文件ID{{else}}{{.FieldName}}文件ID{{end}}"`
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
	{{.FieldName}}FileIDs string `json:"{{.JsonName}}_file_ids" comment:"{{if .Comment}}{{.Comment}}文件ID列表{{else}}{{.FieldName}}文件ID列表{{end}}"`
{{- else}}
	{{.FieldName}} {{.FieldType}} `json:"{{.JsonName}}"{{if .Comment}} comment:"{{.Comment}}"{{end}}`
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
	{{.ForeignKey | ToPascalCase}} uint `json:"{{.ForeignKeyJson}}" comment:"{{if .Comment}}{{.Comment}}ID{{else}}{{.FieldName}}ID{{end}}"`
{{- else if eq .RelationType "many2many"}}
	{{.FieldName}}Ids []uint `json:"{{.JsonName}}_ids" comment:"{{.FieldName}}ID列表"`
{{- end}}
{{- end}}
}

// BatchDelete{{.ModelName}}Request 批量删除{{.Description}}请求
type BatchDelete{{.ModelName}}Request struct {
	Ids []uint `json:"ids" binding:"required" comment:"ID列表"`
}
{{- if .HasAudit}}

// Audit{{.ModelName}}Request 审批{{.Description}}请求
type Audit{{.ModelName}}Request struct {
	AuditStatus int    `json:"audit_status" binding:"required,oneof=1 2" comment:"审批状态 1-通过 2-拒绝"`
	AuditRemark string `json:"audit_remark" binding:"required" comment:"审批备注"`
}
{{- end}}
{{- if .GenerateFrontendApi}}

// Frontend{{.ModelName}}ListRequest 前台{{.Description}}列表请求
type Frontend{{.ModelName}}ListRequest struct {
	PageRequest
{{- range .SearchColumns}}
{{- if ne .ColumnName "status"}}
	{{.FieldName}} {{if or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64")}}*{{.FieldType}}{{else if eq .SearchType "eq"}}*{{.FieldType}}{{else}}string{{end}} `json:"{{.JsonName}}" form:"{{.JsonName}}"{{if .Comment}} comment:"{{.Comment}}"{{end}}`
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
	{{.ForeignKey | ToPascalCase}} *uint `json:"{{.ForeignKeyJson}}" form:"{{.ForeignKeyJson}}" comment:"{{if .Comment}}{{.Comment}}ID{{else}}{{.FieldName}}ID{{end}}"`
{{- end}}
{{- end}}
}
{{- end}}
{{- if .LinkToUser}}

// SaveMy{{.ModelName}}Request 保存我的{{.Description}}请求
type SaveMy{{.ModelName}}Request struct {
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	{{.FieldName}}FileID uint `json:"{{.JsonName}}_file_id" comment:"{{if .Comment}}{{.Comment}}文件ID{{else}}{{.FieldName}}文件ID{{end}}"`
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
	{{.FieldName}}FileIDs string `json:"{{.JsonName}}_file_ids" comment:"{{if .Comment}}{{.Comment}}文件ID列表{{else}}{{.FieldName}}文件ID列表{{end}}"`
{{- else}}
	{{.FieldName}} {{.FieldType}} `json:"{{.JsonName}}"{{if .Comment}} comment:"{{.Comment}}"{{end}}`
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
	{{.ForeignKey | ToPascalCase}} uint `json:"{{.ForeignKeyJson}}" comment:"{{if .Comment}}{{.Comment}}ID{{else}}{{.FieldName}}ID{{end}}"`
{{- else if eq .RelationType "many2many"}}
	{{.FieldName}}Ids []uint `json:"{{.JsonName}}_ids" comment:"{{.FieldName}}ID列表"`
{{- end}}
{{- end}}
}
{{- end}}
