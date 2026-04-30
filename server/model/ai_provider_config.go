package model

type AIProviderConfig struct {
	BaseModel
	Name       string `json:"name" gorm:"size:100;uniqueIndex;comment:平台名称"`
	APIKey     string `json:"api_key" gorm:"type:text;comment:平台API Key"`
	BaseURL    string `json:"base_url" gorm:"size:255;comment:平台Base URL"`
	ModelsJSON string `json:"models_json" gorm:"type:longtext;comment:模型配置JSON"`
	IsDefault  bool   `json:"is_default" gorm:"default:false;comment:是否默认平台"`
	Sort       int    `json:"sort" gorm:"default:0;comment:排序"`
}

func (AIProviderConfig) TableName() string {
	return "ai_providers"
}
