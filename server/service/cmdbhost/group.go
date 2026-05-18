package cmdbhost

import (
	"errors"

	"server/global"
	"server/model"
	"server/model/request"
)

func (s *Service) ListGroups(req *request.CmdbHostGroupListRequest) ([]model.CmdbHostGroup, error) {
	db := global.DB.Model(&model.CmdbHostGroup{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	var list []model.CmdbHostGroup
	if err := db.Order("sort ASC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Service) CreateGroup(req *request.CreateCmdbHostGroupRequest) error {
	var count int64
	if err := global.DB.Model(&model.CmdbHostGroup{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("主机分组名称已存在")
	}
	item := model.CmdbHostGroup{
		Name:   req.Name,
		Sort:   req.Sort,
		Remark: req.Remark,
		Status: req.Status,
	}
	return global.DB.Create(&item).Error
}

func (s *Service) UpdateGroup(id uint, req *request.UpdateCmdbHostGroupRequest) error {
	var item model.CmdbHostGroup
	if err := global.DB.First(&item, id).Error; err != nil {
		return errors.New("主机分组不存在")
	}
	var count int64
	if err := global.DB.Model(&model.CmdbHostGroup{}).Where("name = ? AND id <> ?", req.Name, id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("主机分组名称已存在")
	}
	item.Name = req.Name
	item.Sort = req.Sort
	item.Remark = req.Remark
	item.Status = req.Status
	return global.DB.Save(&item).Error
}

func (s *Service) DeleteGroup(id uint) error {
	var item model.CmdbHostGroup
	if err := global.DB.First(&item, id).Error; err != nil {
		return errors.New("主机分组不存在")
	}
	var bindCount int64
	if err := global.DB.Model(&model.CmdbHost{}).Where("group_id = ?", id).Count(&bindCount).Error; err != nil {
		return err
	}
	if bindCount > 0 {
		return errors.New("主机分组仍被主机使用中")
	}
	return global.DB.Delete(&item).Error
}
