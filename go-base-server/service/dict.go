package service

import (
	"encoding/json"
	"errors"

	"go-base-server/global"
	"go-base-server/model"
	"go-base-server/model/request"
)

type DictService struct{}

var Dict = new(DictService)

// ==================== 字典类型 ====================

// GetDictTypeList 获取字典类型列表（分页）
func (s *DictService) GetDictTypeList(req *request.DictTypeListRequest) ([]model.SysDictType, int64, error) {
	var list []model.SysDictType
	var total int64

	db := global.DB.Model(&model.SysDictType{})

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Type != "" {
		db = db.Where("type LIKE ?", "%"+req.Type+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	if err := db.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetAllDictTypes 获取所有字典类型（不分页）
func (s *DictService) GetAllDictTypes() ([]model.SysDictType, error) {
	var list []model.SysDictType
	err := global.DB.Where("status = ?", 1).Order("created_at DESC").Find(&list).Error
	return list, err
}

// GetDictType 获取字典类型详情
func (s *DictService) GetDictType(id uint) (*model.SysDictType, error) {
	var dictType model.SysDictType
	if err := global.DB.First(&dictType, id).Error; err != nil {
		return nil, errors.New("字典类型不存在")
	}
	return &dictType, nil
}

// CreateDictType 创建字典类型
func (s *DictService) CreateDictType(req *request.CreateDictTypeRequest) error {
	// 检查类型是否已存在
	var count int64
	global.DB.Model(&model.SysDictType{}).Where("type = ?", req.Type).Count(&count)
	if count > 0 {
		return errors.New("字典类型已存在")
	}

	dictType := model.SysDictType{
		Name:   req.Name,
		Type:   req.Type,
		Status: req.Status,
		Remark: req.Remark,
	}

	return global.DB.Create(&dictType).Error
}

// UpdateDictType 更新字典类型
func (s *DictService) UpdateDictType(id uint, req *request.UpdateDictTypeRequest) error {
	var dictType model.SysDictType
	if err := global.DB.First(&dictType, id).Error; err != nil {
		return errors.New("字典类型不存在")
	}

	oldType := dictType.Type // 保存旧类型用于清除缓存

	// 如果修改了类型，检查是否与其他记录冲突
	if req.Type != "" && req.Type != dictType.Type {
		var count int64
		global.DB.Model(&model.SysDictType{}).Where("type = ? AND id != ?", req.Type, id).Count(&count)
		if count > 0 {
			return errors.New("字典类型已存在")
		}
		// 同步更新字典数据的 dict_type
		global.DB.Model(&model.SysDictData{}).Where("dict_type = ?", dictType.Type).Update("dict_type", req.Type)
	}

	dictType.Name = req.Name
	dictType.Type = req.Type
	dictType.Status = req.Status
	dictType.Remark = req.Remark

	err := global.DB.Save(&dictType).Error
	if err == nil {
		// 清除旧类型和新类型的缓存
		_ = Cache.ClearDictCache(oldType)
		if req.Type != oldType {
			_ = Cache.ClearDictCache(req.Type)
		}
	}
	return err
}

// DeleteDictType 删除字典类型
func (s *DictService) DeleteDictType(id uint) error {
	var dictType model.SysDictType
	if err := global.DB.First(&dictType, id).Error; err != nil {
		return errors.New("字典类型不存在")
	}

	typeCode := dictType.Type // 保存字典类型用于清除缓存

	// 删除关联的字典数据
	if err := global.DB.Where("dict_type = ?", dictType.Type).Delete(&model.SysDictData{}).Error; err != nil {
		return err
	}

	err := global.DB.Delete(&dictType).Error
	if err == nil {
		// 清除该字典类型的缓存
		_ = Cache.ClearDictCache(typeCode)
	}
	return err
}

// ==================== 字典数据 ====================

// GetDictDataList 获取字典数据列表（分页）
func (s *DictService) GetDictDataList(req *request.DictDataListRequest) ([]model.SysDictData, int64, error) {
	var list []model.SysDictData
	var total int64

	db := global.DB.Model(&model.SysDictData{}).Where("dict_type = ?", req.DictType)

	if req.Label != "" {
		db = db.Where("label LIKE ?", "%"+req.Label+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	if err := db.Order("sort ASC, created_at DESC").Offset(offset).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetDictDataByType 根据字典类型获取字典数据（不分页，用于下拉框）
// 优先从 Redis 缓存获取
func (s *DictService) GetDictDataByType(dictType string) ([]model.SysDictData, error) {
	// 尝试从缓存获取
	if data, err := Cache.GetDictFromCache(dictType); err == nil {
		var list []model.SysDictData
		if err := json.Unmarshal(data, &list); err == nil {
			return list, nil
		}
	}

	// 缓存未命中，查询数据库
	var list []model.SysDictData
	err := global.DB.Where("dict_type = ? AND status = ?", dictType, 1).
		Order("sort ASC, created_at DESC").Find(&list).Error
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(list); err == nil {
		_ = Cache.SetDictToCache(dictType, data)
	}

	return list, nil
}

// GetDictData 获取字典数据详情
func (s *DictService) GetDictData(id uint) (*model.SysDictData, error) {
	var dictData model.SysDictData
	if err := global.DB.First(&dictData, id).Error; err != nil {
		return nil, errors.New("字典数据不存在")
	}
	return &dictData, nil
}

// CreateDictData 创建字典数据
func (s *DictService) CreateDictData(req *request.CreateDictDataRequest) error {
	// 检查字典类型是否存在
	var typeCount int64
	global.DB.Model(&model.SysDictType{}).Where("type = ?", req.DictType).Count(&typeCount)
	if typeCount == 0 {
		return errors.New("字典类型不存在")
	}

	// 检查同一类型下值是否已存在
	var count int64
	global.DB.Model(&model.SysDictData{}).Where("dict_type = ? AND value = ?", req.DictType, req.Value).Count(&count)
	if count > 0 {
		return errors.New("该字典类型下已存在相同键值")
	}

	dictData := model.SysDictData{
		DictType:  req.DictType,
		Label:     req.Label,
		Value:     req.Value,
		Sort:      req.Sort,
		Status:    req.Status,
		TagType:   req.TagType,
		IsDefault: req.IsDefault,
		Remark:    req.Remark,
	}

	err := global.DB.Create(&dictData).Error
	if err == nil {
		// 清除该字典类型的缓存
		_ = Cache.ClearDictCache(req.DictType)
	}
	return err
}

// UpdateDictData 更新字典数据
func (s *DictService) UpdateDictData(id uint, req *request.UpdateDictDataRequest) error {
	var dictData model.SysDictData
	if err := global.DB.First(&dictData, id).Error; err != nil {
		return errors.New("字典数据不存在")
	}

	dictType := dictData.DictType // 保存字典类型用于清除缓存

	// 如果修改了值，检查是否与同类型下其他记录冲突
	if req.Value != "" && req.Value != dictData.Value {
		var count int64
		global.DB.Model(&model.SysDictData{}).Where("dict_type = ? AND value = ? AND id != ?", dictData.DictType, req.Value, id).Count(&count)
		if count > 0 {
			return errors.New("该字典类型下已存在相同键值")
		}
	}

	dictData.Label = req.Label
	dictData.Value = req.Value
	dictData.Sort = req.Sort
	dictData.Status = req.Status
	dictData.TagType = req.TagType
	dictData.IsDefault = req.IsDefault
	dictData.Remark = req.Remark

	err := global.DB.Save(&dictData).Error
	if err == nil {
		// 清除该字典类型的缓存
		_ = Cache.ClearDictCache(dictType)
	}
	return err
}

// DeleteDictData 删除字典数据
func (s *DictService) DeleteDictData(id uint) error {
	var dictData model.SysDictData
	if err := global.DB.First(&dictData, id).Error; err != nil {
		return errors.New("字典数据不存在")
	}
	dictType := dictData.DictType // 保存字典类型用于清除缓存
	err := global.DB.Delete(&dictData).Error
	if err == nil {
		// 清除该字典类型的缓存
		_ = Cache.ClearDictCache(dictType)
	}
	return err
}

// GetDictLabel 根据字典类型和值获取标签（用于导出）
func (s *DictService) GetDictLabel(dictType string, value string) string {
	list, err := s.GetDictDataByType(dictType)
	if err != nil {
		return value
	}
	for _, item := range list {
		if item.Value == value {
			return item.Label
		}
	}
	return value
}

// GetDictValue 根据字典类型和标签获取值（用于导入）
func (s *DictService) GetDictValue(dictType string, label string) string {
	list, err := s.GetDictDataByType(dictType)
	if err != nil {
		return ""
	}
	for _, item := range list {
		if item.Label == label {
			return item.Value
		}
	}
	return ""
}
