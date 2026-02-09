package model
import (
	"time"
)

// ProductType 产品类型
type ProductType struct {
	ID uint `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name string `json:"name" gorm:"size:255;comment:产品类型名称"` // 产品类型名称
	Icon string `json:"icon" gorm:"size:255;comment:类型图标"` // 类型图标
}

func (ProductType) TableName() string {
	return "product_type"
}

// FillFileURLs 填充文件URL
func (m *ProductType) FillFileURLs() {
}
