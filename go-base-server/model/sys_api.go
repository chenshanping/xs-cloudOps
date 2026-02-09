package model

// FieldInfo 字段信息
type FieldInfo struct {
	Name        string `json:"name"`        // 字段名
	Type        string `json:"type"`        // 字段类型
	Description string `json:"description"` // 字段描述
	Required    bool   `json:"required"`    // 是否必填
	In          string `json:"in"`          // 参数位置: query, body, path
}

type SysApi struct {
	BaseModel
	Path           string `json:"path" gorm:"size:200;comment:API路径"`
	Method         string `json:"method" gorm:"size:10;comment:请求方法"`
	Group          string `json:"group" gorm:"size:50;comment:API分组"`
	Description    string `json:"description" gorm:"size:255;comment:描述"`
	RequestParams  string `json:"request_params" gorm:"type:text;comment:请求参数JSON"`
	ResponseParams string `json:"response_params" gorm:"type:text;comment:响应参数JSON"`
	NeedAuth       bool   `json:"need_auth" gorm:"comment:是否需要认证"`
}

func (SysApi) TableName() string {
	return "sys_api"
}
