package model

{{- /* 计算是否使用BaseModel */ -}}
{{- $useBaseModel := and .HasCreatedAt .HasUpdatedAt .HasDeletedAt -}}
{{- /* 时间包导入判断：不使用BaseModel时需要导入（CreatedAt/UpdatedAt/Audit），或使用BaseModel但有审批或自定义time字段 */ -}}
{{- $needTime := or (and (not $useBaseModel) (or .HasCreatedAt .HasUpdatedAt)) .HasAudit .HasTimeField -}}
{{- $needGorm := and (not $useBaseModel) .HasDeletedAt -}}
{{- if or $needTime $needGorm }}
import (
{{- if $needTime  }}
	"time"
{{- end}}
{{- if $needGorm }}
	"gorm.io/gorm"
{{- end}}
)
{{- end}}

// {{.ModelName}} {{.Description}}
type {{.ModelName}} struct {
{{- if $useBaseModel }}
	BaseModel
{{- else }}
	ID uint `json:"id" gorm:"primarykey"`
{{- if .HasCreatedAt }}
	CreatedAt time.Time `json:"created_at"`
{{- end}}
{{- if .HasUpdatedAt }}
	UpdatedAt time.Time `json:"updated_at"`
{{- end}}
{{- if .HasDeletedAt }}
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
{{- end}}
{{- end}}
{{- range .Columns}}
{{- if not .IsPrimaryKey}}
{{- if and (ne .ColumnName "created_at") (ne .ColumnName "updated_at") (ne .ColumnName "deleted_at")}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	{{.FieldName}}FileID  uint     `json:"{{.JsonName}}_file_id" gorm:"comment:{{.Comment}}文件ID"`
	{{.FieldName}}File    *SysFile `json:"-" gorm:"foreignKey:{{.FieldName}}FileID;references:ID"`
	{{.FieldName}}URL     string   `json:"{{.JsonName}}_url" gorm:"-"`
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
	{{.FieldName}}FileIDs string   `json:"{{.JsonName}}_file_ids" gorm:"type:varchar(500);comment:{{.Comment}}文件ID列表"`
	{{.FieldName}}URLs    []string `json:"{{.JsonName}}_urls" gorm:"-"`
	{{.FieldName}}Names   []string `json:"{{.JsonName}}_names" gorm:"-"`
{{- else}}
	{{.FieldName}} {{.FieldType}} `json:"{{.JsonName}}" gorm:"{{.GormTag}}"`{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- if .LinkToUser}}
	UserID uint     `json:"user_id" gorm:"uniqueIndex;comment:关联用户ID"`
	User   *SysUser `json:"user" gorm:"foreignKey:UserID"`
{{- end}}
{{- if .HasCreatedBy}}
	CreatedBy uint     `json:"created_by" gorm:"comment:创建人 ID"`
	Creator   *SysUser `json:"creator" gorm:"foreignKey:CreatedBy"`
{{- if .HasCreatedByProfile}}
	CreatorProfile *{{.CreatedByProfileModel}} `json:"creator_profile" gorm:"-"` // 创建者身份信息
{{- end}}
{{- end}}
{{- if .HasAudit}}
	AuditStatus int        `json:"audit_status" gorm:"type:tinyint;default:0;comment:审批状态 0-待审批 1-审批通过 2-审批拒绝"`
	AuditRemark string     `json:"audit_remark" gorm:"type:varchar(500);comment:审批备注"`
	AuditTime   *time.Time `json:"audit_time" gorm:"comment:审批时间"`
	AuditBy     uint       `json:"audit_by" gorm:"comment:审批人id"`
	Auditor     *SysUser   `json:"auditor" gorm:"foreignKey:AuditBy"`
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "hasMany"}}
	{{.FieldName}} []{{.RelatedModel}} `json:"{{.JsonName}}" gorm:"foreignKey:{{.ForeignKey}}"`
{{- else if eq .RelationType "many2many"}}
	{{.FieldName}} []{{.RelatedModel}} `json:"{{.JsonName}}" gorm:"many2many:{{.JoinTable}};"`
{{- else if eq .RelationType "belongsTo"}}
	{{.FieldName}} *{{.RelatedModel}} `json:"{{.JsonName}}" gorm:"foreignKey:{{.ForeignKey | ToPascalCase}}"`
{{- end}}
{{- end}}
}

func ({{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}

// FillFileURLs 填充文件URL
func (m *{{.ModelName}}) FillFileURLs() {
{{- range .Columns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	if m.{{.FieldName}}File != nil {
		m.{{.FieldName}}URL = m.{{.FieldName}}File.URL
	}
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
	// {{.FieldName}}URLs 需要在service中根据{{.FieldName}}FileIDs查询填充
{{- end}}
{{- end}}
{{- if .LinkToUser}}
	if m.User != nil {
		m.User.FillAvatarURL()
	}
{{- end}}
{{- if .HasCreatedBy}}
	if m.Creator != nil {
		m.Creator.FillAvatarURL()
	}
{{- if .HasCreatedByProfile}}
	if m.CreatorProfile != nil {
		m.CreatorProfile.FillFileURLs()
	}
{{- end}}
{{- end}}
{{- if .HasAudit}}
	if m.Auditor != nil {
		m.Auditor.FillAvatarURL()
	}
{{- end}}
}
