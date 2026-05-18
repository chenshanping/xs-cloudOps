package cmdbhost

import (
	"errors"

	"server/global"
	"server/model"
	"server/model/request"
)

func (s *Service) ListTags(req *request.CmdbHostTagListRequest) ([]model.CmdbHostTag, error) {
	db := global.DB.Model(&model.CmdbHostTag{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	var list []model.CmdbHostTag
	if err := db.Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Service) CreateTag(req *request.CreateCmdbHostTagRequest) error {
	var count int64
	if err := global.DB.Model(&model.CmdbHostTag{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("主机标签名称已存在")
	}
	item := model.CmdbHostTag{
		Name:   req.Name,
		Color:  req.Color,
		Remark: req.Remark,
	}
	return global.DB.Create(&item).Error
}

func (s *Service) UpdateTag(id uint, req *request.UpdateCmdbHostTagRequest) error {
	var item model.CmdbHostTag
	if err := global.DB.First(&item, id).Error; err != nil {
		return errors.New("主机标签不存在")
	}
	var count int64
	if err := global.DB.Model(&model.CmdbHostTag{}).Where("name = ? AND id <> ?", req.Name, id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("主机标签名称已存在")
	}
	item.Name = req.Name
	item.Color = req.Color
	item.Remark = req.Remark
	return global.DB.Save(&item).Error
}

func (s *Service) DeleteTag(id uint) error {
	var item model.CmdbHostTag
	if err := global.DB.First(&item, id).Error; err != nil {
		return errors.New("主机标签不存在")
	}
	if err := global.DB.Where("tag_id = ?", id).Delete(&model.CmdbHostTagRel{}).Error; err != nil {
		return err
	}
	return global.DB.Delete(&item).Error
}
