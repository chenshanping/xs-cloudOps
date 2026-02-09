package model
import (
	"time"
)

// Product 产品信息
type Product struct {
	ID uint `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TypeId uint `json:"type_id" gorm:"comment:类型id"` // 类型id
	Name string `json:"name" gorm:"size:255;comment:产品名称"` // 产品名称
	ProductType *ProductType `json:"product_type" gorm:"foreignKey:TypeId"`
}

func (Product) TableName() string {
	return "product"
}

// FillFileURLs 填充文件URL
func (m *Product) FillFileURLs() {
}
