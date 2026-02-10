package service

import (
	"errors"

	"go-base-server/global"
	"go-base-server/model"
	"go-base-server/model/request"
)

type ProductService struct{}

var Product = new(ProductService)

// GetProductList 获取产品信息列表
func (s *ProductService) GetProductList(req *request.ProductListRequest) ([]model.Product, int64, error) {
	var list []model.Product
	var total int64

	db := global.DB.Model(&model.Product{})
	if req.Name != nil {
		db = db.Where("name = ?", *req.Name)
	}
	if req.Num != nil {
		db = db.Where("num >= ?", *req.Num)
	}
	if req.TypeId != nil {
		db = db.Where("type_id = ?", *req.TypeId)
	}

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
			"num": true,
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
	query = query.Preload("ProductType")
	if err := query.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	// 填充文件URL
	for i := range list {
		list[i].FillFileURLs()
	}

	return list, total, nil
}

// GetProduct 获取产品信息详情
func (s *ProductService) GetProduct(id uint) (*model.Product, error) {
	var data model.Product
	query := global.DB
	query = query.Preload("ProductType")
	if err := query.First(&data, id).Error; err != nil {
		return nil, err
	}
	data.FillFileURLs()
	return &data, nil
}

// CreateProduct 创建产品信息
func (s *ProductService) CreateProduct(req *request.CreateProductRequest, userID uint) error {

	data := model.Product{
		Name: req.Name,
		Num: req.Num,
		Price: req.Price,
		Status: req.Status,
		TypeId: req.TypeId,
	}
	return global.DB.Create(&data).Error
}

// UpdateProduct 更新产品信息
func (s *ProductService) UpdateProduct(id uint, req *request.UpdateProductRequest) error {
	var data model.Product
	if err := global.DB.First(&data, id).Error; err != nil {
		return errors.New("数据不存在")
	}

	updates := map[string]interface{}{
		"name": req.Name,
		"num": req.Num,
		"price": req.Price,
		"status": req.Status,
		"type_id": req.TypeId,
	}
	return global.DB.Model(&data).Updates(updates).Error
}

// DeleteProduct 删除产品信息
func (s *ProductService) DeleteProduct(id uint) error {
	var data model.Product
	if err := global.DB.First(&data, id).Error; err != nil {
		return errors.New("数据不存在")
	}
	return global.DB.Delete(&data).Error
}

// BatchDeleteProduct 批量删除产品信息
func (s *ProductService) BatchDeleteProduct(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return global.DB.Where("id IN ?", ids).Delete(&model.Product{}).Error
}

// GetProductOptions 获取产品信息选项列表（带可选关联统计）
// excludeDeleted: 是否排除软删除数据（统计表有deleted_at字段时传true）
// countCreatedBy: 统计时按创建人过滤（数据隔离用，传当前用户ID，0表示不过滤）
func (s *ProductService) GetProductOptions(displayField, countTable, countForeignKey string, excludeDeleted bool, countCreatedBy uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	if displayField == "" {
		displayField = "name"
	}

	// 无统计关联时，简单查询
	if countTable == "" || countForeignKey == "" {
		err := global.DB.Model(&model.Product{}).
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
	err := global.DB.Table("product").
		Select("product.id, product.`" + displayField + "` as name, COALESCE(sub.cnt, 0) as count").
		Joins("LEFT JOIN (?) as sub ON product.id = sub.fk", subQuery).
		Order("product.id ASC").
		Find(&results).Error
	// 转换[]byte为string
	for i := range results {
		if nameBytes, ok := results[i]["name"].([]byte); ok {
			results[i]["name"] = string(nameBytes)
		}
	}

	return results, err
}

// GetProductStatsTypeId 获取产品信息按产品类型分组统计
func (s *ProductService) GetProductStatsTypeId() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := global.DB.Table("product").
		Select("type_id as group_key, COUNT(*) as value").
		Where("deleted_at IS NULL").
		Group("type_id").
		Order("value DESC").
		Find(&results).Error
	return results, err
}

// GetProductStatsStatus 获取产品信息按产品状态分组统计
func (s *ProductService) GetProductStatsStatus() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := global.DB.Table("product").
		Select("status as group_key, COUNT(*) as value").
		Where("deleted_at IS NULL").
		Group("status").
		Order("value DESC").
		Find(&results).Error
	return results, err
}

// GetProductTrendStats 获取产品信息趋势统计
func (s *ProductService) GetProductTrendStats(days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	if days <= 0 {
		days = 30
	}
	err := global.DB.Table("product").
		Select("DATE(created_at) as date, COUNT(*) as value").
		Where("created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", days).
		Where("deleted_at IS NULL").
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&results).Error
	return results, err
}
