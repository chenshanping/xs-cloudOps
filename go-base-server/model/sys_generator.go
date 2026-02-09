package model

// SysGenerator 代码生成记录
type SysGenerator struct {
	BaseModel
	GenTableName string `json:"table_name" gorm:"column:table_name;size:128;index"`
	ModuleName   string `json:"module_name" gorm:"size:128;index"`
	Description  string `json:"description" gorm:"size:255"`
	ConfigJSON   string `json:"config_json" gorm:"type:longtext"`
}

func (SysGenerator) TableName() string {
	return "sys_generator"
}
