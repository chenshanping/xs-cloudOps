package cmdbcredential

import (
	"errors"

	"server/global"
	"server/model"
	"server/model/request"
)

type Service struct{}

var Default = &Service{}

func (s *Service) List(req *request.CmdbCredentialListRequest) ([]CredentialSummary, int64, error) {
	var list []model.CmdbSshCredential
	var total int64

	db := global.DB.Model(&model.CmdbSshCredential{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.AuthType != "" {
		db = db.Where("auth_type = ?", req.AuthType)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("updated_at DESC").Offset(req.GetOffset()).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	result := make([]CredentialSummary, 0, len(list))
	for _, item := range list {
		var bindCount int64
		if err := global.DB.Model(&model.CmdbHost{}).Where("credential_id = ?", item.ID).Count(&bindCount).Error; err != nil {
			return nil, 0, err
		}
		result = append(result, buildSummary(item, bindCount))
	}

	return result, total, nil
}

func (s *Service) ListOptions() ([]CredentialSummary, error) {
	var list []model.CmdbSshCredential
	if err := global.DB.Order("updated_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	result := make([]CredentialSummary, 0, len(list))
	for _, item := range list {
		var bindCount int64
		if err := global.DB.Model(&model.CmdbHost{}).Where("credential_id = ?", item.ID).Count(&bindCount).Error; err != nil {
			return nil, err
		}
		result = append(result, buildSummary(item, bindCount))
	}
	return result, nil
}

func (s *Service) Get(id uint) (*CredentialSummary, error) {
	var item model.CmdbSshCredential
	if err := global.DB.First(&item, id).Error; err != nil {
		return nil, errors.New("SSH凭据不存在")
	}
	var bindCount int64
	if err := global.DB.Model(&model.CmdbHost{}).Where("credential_id = ?", item.ID).Count(&bindCount).Error; err != nil {
		return nil, err
	}
	summary := buildSummary(item, bindCount)
	return &summary, nil
}

func (s *Service) Create(req *request.CreateCmdbCredentialRequest) error {
	if err := validateCredentialPayload(req.AuthType, req.Password, req.PrivateKey); err != nil {
		return err
	}
	var count int64
	if err := global.DB.Model(&model.CmdbSshCredential{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("凭据名称已存在")
	}
	item := model.CmdbSshCredential{
		Name:       req.Name,
		AuthType:   req.AuthType,
		Username:   req.Username,
		Password:   req.Password,
		PrivateKey: req.PrivateKey,
		Passphrase: req.Passphrase,
		Remark:     req.Remark,
	}
	return global.DB.Create(&item).Error
}

func (s *Service) Update(id uint, req *request.UpdateCmdbCredentialRequest) error {
	if err := validateCredentialPayload(req.AuthType, req.Password, req.PrivateKey); err != nil {
		return err
	}
	var item model.CmdbSshCredential
	if err := global.DB.First(&item, id).Error; err != nil {
		return errors.New("SSH凭据不存在")
	}
	var count int64
	if err := global.DB.Model(&model.CmdbSshCredential{}).Where("name = ? AND id <> ?", req.Name, id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("凭据名称已存在")
	}

	item.Name = req.Name
	item.AuthType = req.AuthType
	item.Username = req.Username
	item.Remark = req.Remark
	if req.Password != "" || req.AuthType == model.CmdbCredentialAuthTypePassword {
		item.Password = req.Password
		if req.AuthType == model.CmdbCredentialAuthTypePassword {
			item.PrivateKey = ""
			item.Passphrase = ""
		}
	}
	if req.PrivateKey != "" || req.AuthType == model.CmdbCredentialAuthTypePrivateKey {
		item.PrivateKey = req.PrivateKey
		item.Passphrase = req.Passphrase
		if req.AuthType == model.CmdbCredentialAuthTypePrivateKey {
			item.Password = ""
		}
	}

	return global.DB.Save(&item).Error
}

func (s *Service) Delete(id uint) error {
	var item model.CmdbSshCredential
	if err := global.DB.First(&item, id).Error; err != nil {
		return errors.New("SSH凭据不存在")
	}
	var bindCount int64
	if err := global.DB.Model(&model.CmdbHost{}).Where("credential_id = ?", id).Count(&bindCount).Error; err != nil {
		return err
	}
	if bindCount > 0 {
		return errors.New("SSH凭据仍被主机使用中")
	}
	return global.DB.Delete(&item).Error
}

func validateCredentialPayload(authType, password, privateKey string) error {
	switch authType {
	case model.CmdbCredentialAuthTypePassword:
		if password == "" {
			return errors.New("密码认证方式必须填写登录密码")
		}
	case model.CmdbCredentialAuthTypePrivateKey:
		if privateKey == "" {
			return errors.New("私钥认证方式必须填写私钥内容")
		}
	default:
		return errors.New("不支持的认证方式")
	}
	return nil
}
