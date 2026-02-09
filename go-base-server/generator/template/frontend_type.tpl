// {{.Description}}
{{- if .Relations}}
{{- range .Relations}}
import type { {{.RelatedModel}} } from './{{.RelatedTable}}'
{{- end}}
{{- end}}

// {{.Description}}响应体
export interface {{.ModelName}} {
  id: number
{{- range .Columns}}
{{- if or (eq .FormType "images") (eq .FormType "files")}}
  {{.JsonName}}_file_ids?: string
  {{.JsonName}}_urls?: string[]
  {{.JsonName}}_names?: string[]
{{- else if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
  {{.JsonName}}_file_id?: number
  {{.JsonName}}_url?: string
{{- else}}
  {{.JsonName}}{{if not .IsRequired}}?{{end}}: {{.TsType}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
  {{.JsonName}}?: {{.RelatedModel}}
{{- else if eq .RelationType "many2many"}}
  {{.JsonName}}_ids?: number[]
  {{.JsonName}}?: {{.RelatedModel}}[]
{{- else if eq .RelationType "hasMany"}}
  {{.JsonName}}?: {{.RelatedModel}}[]
{{- end}}
{{- end}}
{{- if .HasCreatedAt}}
  created_at?: string
{{- end}}
{{- if .HasUpdatedAt}}
  updated_at?: string
{{- end}}
{{- if .LinkToUser}}
  user_id?: number
  user?: { id: number; username: string; nickname?: string; avatar_file_url?: string }
{{- end}}
{{- if .HasCreatedBy}}
  created_by?: number
  creator?: { id: number; username: string; nickname?: string }
{{- if .HasCreatedByProfile}}
  creator_profile?: { id: number; user_id: number; {{.CreatedByProfileField}}: string }
{{- end}}
{{- end}}
{{- if .HasAudit}}
  audit_status?: number  // 0-待审批 1-审批通过 2-审批拒绝
  audit_remark?: string
  audit_time?: string
  audit_by?: number
  auditor?: { id: number; username: string; nickname?: string }
{{- end}}
}

// 创建{{.Description}}请求体
export interface Create{{.ModelName}}Request {
{{- if .LinkToUser}}
  user_id: number
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "images") (eq .FormType "files")}}
  {{.JsonName}}_file_ids{{if not .IsRequired}}?{{end}}: string
{{- else if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
  {{.JsonName}}_file_id{{if not .IsRequired}}?{{end}}: number
{{- else}}
  {{.JsonName}}{{if not .IsRequired}}?{{end}}: {{.TsType}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
  {{.ForeignKeyJson}}?: number
{{- else if eq .RelationType "many2many"}}
  {{.JsonName}}_ids?: number[]
{{- end}}
{{- end}}
}

// 更新{{.Description}}请求体
export type Update{{.ModelName}}Request = Partial<Create{{.ModelName}}Request>

// {{.Description}}查询参数
export interface {{.ModelName}}Query {
  page?: number
  page_size?: number
{{- range .SearchColumns}}
  {{.JsonName}}?: {{.TsType}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
  {{.ForeignKeyJson}}?: number
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
  created_by?: number
{{- end}}
  sort_field?: string
  sort_order?: 'ascend' | 'descend'
}

// 选项项（用于下拉选择）
export interface OptionItem {
  id: number
  name: string
  count?: number
}
{{- if .LinkToUser}}

// 保存我的{{.Description}}请求体
export interface SaveMy{{.ModelName}}Request {
{{- range .FormColumns}}
{{- if or (eq .FormType "images") (eq .FormType "files")}}
  {{.JsonName}}_file_ids?: string
{{- else if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
  {{.JsonName}}_file_id?: number
{{- else}}
  {{.JsonName}}?: {{.TsType}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
  {{.ForeignKeyJson}}?: number
{{- else if eq .RelationType "many2many"}}
  {{.JsonName}}_ids?: number[]
{{- end}}
{{- end}}
}
{{- end}}
