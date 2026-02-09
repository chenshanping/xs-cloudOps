package service

import (
	"errors"

	"go-base-server/global"
	"go-base-server/model"
	"go-base-server/model/request"
)

type ProductTypeService struct{}

var ProductType = new(ProductTypeService)

// GetProductTypeList 获取产品类型列表
func (s *ProductTypeService) GetProductTypeList(req *request.ProductTypeListRequest) ([]model.ProductType, int64, error) {
	var list []model.ProductType
	var total int64

	db := global.DB.Model(&model.ProductType{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	// 排序处理
	orderBy := "id DESC"
	if req.SortField != "" && req.SortOrder != "" {
		// 前端传入排序
		allowedFields := map[string]bool{
			"id": true,
			"created_at": true,
		}
		if allowedFields[req.SortField] {
			order := "ASC"
			if req.SortOrder == "descend" || req.SortOrder == "desc" {
				order = "DESC"
			}
			orderBy = req.SortField + " " + order
		}
	}
	query := db.Offset(offset).Limit(req.PageSize).Order(orderBy)
	if err := query.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	// 填充文件URL
	for i := range list {
		list[i].FillFileURLs()
	}

	return list, total, nil
}

// GetProductType 获取产品类型详情
func (s *ProductTypeService) GetProductType(id uint) (*model.ProductType, error) {
	var data model.ProductType
	query := global.DB
	if err := query.First(&data, id).Error; err != nil {
		return nil, err
	}
	data.FillFileURLs()
	return &data, nil
}

// CreateProductType 创建产品类型
func (s *ProductTypeService) CreateProductType(req *request.CreateProductTypeRequest, userID uint) error {

	data := model.ProductType{
		Name: req.Name,
		Icon: req.Icon,
	}
	return global.DB.Create(&data).Error
}

// UpdateProductType 更新产品类型
func (s *ProductTypeService) UpdateProductType(id uint, req *request.UpdateProductTypeRequest) error {
	var data model.ProductType
	if err := global.DB.First(&data, id).Error; err != nil {
		return errors.New("数据不存在")
	}

	updates := map[string]interface{}{
		"name": req.Name,
		"icon": req.Icon,
	}
	return global.DB.Model(&data).Updates(updates).Error
}

// DeleteProductType 删除产品类型
func (s *ProductTypeService) DeleteProductType(id uint) error {
	var data model.ProductType
	if err := global.DB.First(&data, id).Error; err != nil {
		return errors.New("数据不存在")
	}
	return global.DB.Delete(&data).Error
}

// BatchDeleteProductType 批量删除产品类型
func (s *ProductTypeService) BatchDeleteProductType(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return global.DB.Where("id IN ?", ids).Delete(&model.ProductType{}).Error
}

// GetProductTypeOptions 获取产品类型选项列表（带可选关联统计）
// excludeDeleted: 是否排除软删除数据（统计表有deleted_at字段时传true）
// countCreatedBy: 统计时按创建人过滤（数据隔离用，传当前用户ID，0表示不过滤）
func (s *ProductTypeService) GetProductTypeOptions(displayField, countTable, countForeignKey string, excludeDeleted bool, countCreatedBy uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	if displayField == "" {
		displayField = "name"
	}

	// 无统计关联时，简单查询
	if countTable == "" || countForeignKey == "" {
		err := global.DB.Model(&model.ProductType{}).
			Select("id, `" + displayField + "` as name").
			Order("id ASC").
			Find(&results).Error
		// 转换[]byte为string（GORM返回map时字符串字段会是[]byte类型）
		for i := range results {
			if nameBytes, ok := results[i]["name"].([]byte); ok {
				results[i]["name"] = string(nameBytes)
			}
		}
		return results, err
	}

	// 有统计关联时，使用子查询
	subQuery := global.DB.Table(countTable).
		Select(countForeignKey + " as fk, COUNT(*) as cnt")

	// 排除软删除数据
	if excludeDeleted {
		subQuery = subQuery.Where("deleted_at IS NULL")
	}
	// 数据隔离：统计时按创建人过滤
	if countCreatedBy > 0 {
		subQuery = subQuery.Where("created_by = ?", countCreatedBy)
	}
	subQuery = subQuery.Group(countForeignKey)
	err := global.DB.Table("product_type").
		Select("product_type.id, product_type.`" + displayField + "` as name, COALESCE(sub.cnt, 0) as count").
		Joins("LEFT JOIN (?) as sub ON product_type.id = sub.fk", subQuery).
		Order("product_type.id ASC").
		Find(&results).Error
	// 转换[]byte为string
	for i := range results {
		if nameBytes, ok := results[i]["name"].([]byte); ok {
			results[i]["name"] = string(nameBytes)
		}
	}

	return results, err
}
