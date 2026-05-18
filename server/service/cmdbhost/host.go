package cmdbhost

import (
	"errors"

	"server/global"
	"server/model"
	"server/model/request"
)

type Service struct{}

var Default = &Service{}

func (s *Service) ListHosts(req *request.CmdbHostListRequest) ([]HostListItem, int64, error) {
	var list []model.CmdbHost
	var total int64
	db := global.DB.Model(&model.CmdbHost{})
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("name LIKE ? OR ssh_host LIKE ? OR private_ip LIKE ? OR public_ip LIKE ? OR hostname LIKE ?", like, like, like, like, like)
	}
	if req.GroupID != nil {
		db = db.Where("group_id = ?", *req.GroupID)
	}
	if req.VerifyStatus != "" {
		db = db.Where("verify_status = ?", req.VerifyStatus)
	}
	if req.Environment != "" {
		db = db.Where("environment = ?", req.Environment)
	}
	if req.TagID != nil {
		db = db.Joins("JOIN cmdb_host_tag_rel rel ON rel.host_id = cmdb_host.id").Where("rel.tag_id = ?", *req.TagID)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("updated_at DESC").Offset(req.GetOffset()).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return s.buildHostListItems(list)
}

func (s *Service) GetHost(id uint) (*HostListItem, error) {
	var item model.CmdbHost
	if err := global.DB.First(&item, id).Error; err != nil {
		return nil, errors.New("主机不存在")
	}
	list, _, err := s.buildHostListItems([]model.CmdbHost{item})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("主机不存在")
	}
	return &list[0], nil
}

func (s *Service) buildHostListItems(items []model.CmdbHost) ([]HostListItem, int64, error) {
	result := make([]HostListItem, 0, len(items))
	for _, item := range items {
		var group model.CmdbHostGroup
		if err := global.DB.First(&group, item.GroupID).Error; err != nil {
			return nil, 0, err
		}
		var credential model.CmdbSshCredential
		if err := global.DB.First(&credential, item.CredentialID).Error; err != nil {
			return nil, 0, err
		}
		var rels []model.CmdbHostTagRel
		if err := global.DB.Where("host_id = ?", item.ID).Find(&rels).Error; err != nil {
			return nil, 0, err
		}
		tagIDs := make([]uint, 0, len(rels))
		for _, rel := range rels {
			tagIDs = append(tagIDs, rel.TagID)
		}
		var tags []model.CmdbHostTag
		if len(tagIDs) > 0 {
			if err := global.DB.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
				return nil, 0, err
			}
		}
		tagViews := make([]HostTagView, 0, len(tags))
		for _, tag := range tags {
			tagViews = append(tagViews, HostTagView{ID: tag.ID, Name: tag.Name, Color: tag.Color})
		}
		result = append(result, HostListItem{
			ID:          item.ID,
			Name:        item.Name,
			Group:       HostGroupOption{ID: group.ID, Name: group.Name, Sort: group.Sort, Status: group.Status, Remark: group.Remark},
			Tags:        tagViews,
			Environment: item.Environment,
			Owner:       item.Owner,
			Remark:      item.Remark,
			PrivateIP:   item.PrivateIP,
			PublicIP:    item.PublicIP,
			SshHost:     item.SshHost,
			SshPort:     item.SshPort,
			CredentialSummary: HostCredentialSummary{
				ID:       credential.ID,
				Name:     credential.Name,
				AuthType: credential.AuthType,
				Username: credential.Username,
			},
			VerifyStatus:    item.VerifyStatus,
			VerifyMessage:   item.VerifyMessage,
			LastVerifiedAt:  item.LastVerifiedAt,
			Hostname:        item.Hostname,
			OS:              item.OS,
			Platform:        item.Platform,
			PlatformVersion: item.PlatformVersion,
			KernelVersion:   item.KernelVersion,
			Architecture:    item.Architecture,
			CpuCores:        item.CpuCores,
			MemoryMB:        item.MemoryMB,
			UpdatedAt:       item.UpdatedAt,
		})
	}
	return result, int64(len(result)), nil
}

func (s *Service) CreateHost(req *request.CreateCmdbHostRequest) (*HostListItem, error) {
	if err := s.validateHostRefs(req.GroupID, req.CredentialID, req.TagIDs); err != nil {
		return nil, err
	}
	var count int64
	if err := global.DB.Model(&model.CmdbHost{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("主机名称已存在")
	}
	port := req.SshPort
	if port <= 0 {
		port = 22
	}
	item := model.CmdbHost{
		Name:         req.Name,
		GroupID:      req.GroupID,
		Environment:  req.Environment,
		Owner:        req.Owner,
		PrivateIP:    req.PrivateIP,
		PublicIP:     req.PublicIP,
		SshHost:      req.SshHost,
		SshPort:      port,
		CredentialID: req.CredentialID,
		Remark:       req.Remark,
		VerifyStatus: model.CmdbHostVerifyStatusPending,
	}
	if err := global.DB.Create(&item).Error; err != nil {
		return nil, err
	}
	if err := s.replaceHostTags(item.ID, req.TagIDs); err != nil {
		return nil, err
	}
	_ = s.VerifyHost(item.ID)
	return s.GetHost(item.ID)
}

func (s *Service) UpdateHost(id uint, req *request.UpdateCmdbHostRequest) (*HostListItem, error) {
	if err := s.validateHostRefs(req.GroupID, req.CredentialID, req.TagIDs); err != nil {
		return nil, err
	}
	var item model.CmdbHost
	if err := global.DB.First(&item, id).Error; err != nil {
		return nil, errors.New("主机不存在")
	}
	var count int64
	if err := global.DB.Model(&model.CmdbHost{}).Where("name = ? AND id <> ?", req.Name, id).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("主机名称已存在")
	}
	port := req.SshPort
	if port <= 0 {
		port = 22
	}
	item.Name = req.Name
	item.GroupID = req.GroupID
	item.Environment = req.Environment
	item.Owner = req.Owner
	item.PrivateIP = req.PrivateIP
	item.PublicIP = req.PublicIP
	item.SshHost = req.SshHost
	item.SshPort = port
	item.CredentialID = req.CredentialID
	item.Remark = req.Remark
	item.VerifyStatus = model.CmdbHostVerifyStatusPending
	item.VerifyMessage = ""
	item.LastVerifiedAt = nil
	if err := global.DB.Save(&item).Error; err != nil {
		return nil, err
	}
	if err := s.replaceHostTags(item.ID, req.TagIDs); err != nil {
		return nil, err
	}
	_ = s.VerifyHost(item.ID)
	return s.GetHost(item.ID)
}

func (s *Service) DeleteHost(id uint) error {
	var item model.CmdbHost
	if err := global.DB.First(&item, id).Error; err != nil {
		return errors.New("主机不存在")
	}
	if err := global.DB.Where("host_id = ?", id).Delete(&model.CmdbHostTagRel{}).Error; err != nil {
		return err
	}
	return global.DB.Delete(&item).Error
}

func (s *Service) validateHostRefs(groupID, credentialID uint, tagIDs []uint) error {
	var groupCount int64
	if err := global.DB.Model(&model.CmdbHostGroup{}).Where("id = ?", groupID).Count(&groupCount).Error; err != nil {
		return err
	}
	if groupCount == 0 {
		return errors.New("主机分组不存在")
	}
	var credentialCount int64
	if err := global.DB.Model(&model.CmdbSshCredential{}).Where("id = ?", credentialID).Count(&credentialCount).Error; err != nil {
		return err
	}
	if credentialCount == 0 {
		return errors.New("SSH凭据不存在")
	}
	if len(tagIDs) > 0 {
		var tagCount int64
		if err := global.DB.Model(&model.CmdbHostTag{}).Where("id IN ?", tagIDs).Count(&tagCount).Error; err != nil {
			return err
		}
		if tagCount != int64(len(uniqueUintSlice(tagIDs))) {
			return errors.New("存在无效的主机标签")
		}
	}
	return nil
}

func (s *Service) replaceHostTags(hostID uint, tagIDs []uint) error {
	if err := global.DB.Where("host_id = ?", hostID).Delete(&model.CmdbHostTagRel{}).Error; err != nil {
		return err
	}
	for _, tagID := range uniqueUintSlice(tagIDs) {
		if err := global.DB.Create(&model.CmdbHostTagRel{HostID: hostID, TagID: tagID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func uniqueUintSlice(items []uint) []uint {
	seen := make(map[uint]struct{}, len(items))
	result := make([]uint, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}
