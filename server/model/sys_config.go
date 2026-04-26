package model

type SysConfig struct {
	BaseModel
	Name      string `json:"name" gorm:"size:100;comment:配置名称"`
	Key       string `json:"key" gorm:"size:100;uniqueIndex;comment:配置键"`
	Value     string `json:"value" gorm:"type:text;comment:配置值"`
	ValueType string `json:"value_type" gorm:"size:20;default:string;comment:值类型 string/json"`
	Remark    string `json:"remark" gorm:"size:255;comment:备注"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}
