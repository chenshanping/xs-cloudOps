package model

// Product 产品信息
type Product struct {
	BaseModel
	TypeId uint `json:"type_id" gorm:"comment:产品类型"` // 产品类型
	Name string `json:"name" gorm:"size:255;comment:产品名称"` // 产品名称
	Num int `json:"num" gorm:"comment:产品数量"` // 产品数量
	Price float64 `json:"price" gorm:"comment:产品单价"` // 产品单价
	Status string `json:"status" gorm:"size:255;comment:状态"` // 状态
	ProductType *ProductType `json:"product_type" gorm:"foreignKey:TypeId"`
}

func (Product) TableName() string {
	return "product"
}

// FillFileURLs 填充文件URL
func (m *Product) FillFileURLs() {
}
