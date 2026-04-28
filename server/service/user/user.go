package user

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
	"server/service/configsvc"
	"server/service/core"
	"server/utils"
)

type UserService struct{}

var Default = &UserService{}

const (
	userDefaultPasswordConfigKey = "user_default_password"
	userDefaultPasswordFallback  = "123456"
)

func isProtectedBatchStatusUser(user model.SysUser) bool {
	if user.ID == 1 || user.Username == "admin" {
		return true
	}
	for _, role := range user.Roles {
		if role.ID == 1 || role.Code == "admin" || role.Code == "super_admin" {
			return true
		}
	}
	return false
}

func validateBatchUserStatusTargets(ids []uint, status int, operatorID uint, users []model.SysUser) error {
	if len(ids) == 0 {
		return errors.New("请选择要修改状态的用户")
	}
	if status != 0 && status != 1 {
		return errors.New("状态值无效")
	}

	if status == 0 {
		for _, id := range ids {
			if id == operatorID {
				return errors.New("不能批量禁用自己")
			}
		}
	}

	for _, user := range users {
		if isProtectedBatchStatusUser(user) {
			return fmt.Errorf("用户「%s」为受保护管理员，不能批量修改状态", user.Username)
		}
	}

	return nil
}

func validateUserGender(gender int) error {
	if gender < 0 || gender > 2 {
		return errors.New("性别值无效")
	}
	return nil
}

func (s *UserService) GetUserOptions(operatorID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	scope, err := core.ResolveUserDataScope(operatorID)
	if err != nil {
		return nil, err
	}

	db := global.DB.Model(&model.SysUser{})
	db = core.ApplyUserDataScope(db, scope, "sys_user")
	err = db.
		Select("id, username, nickname").
		Where("status = ?", 1).
		Order("id ASC").
		Find(&results).Error
	return results, err
}

func (s *UserService) Login(username, password string) (*model.SysUser, error) {
	var user model.SysUser
	if err := global.DB.Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

func (s *UserService) GetUserInfo(userID uint) (*model.SysUser, error) {
	var user model.SysUser
	if err := global.DB.Preload("Roles").Preload("Dept").Preload("AvatarFile").First(&user, userID).Error; err != nil {
		return nil, err
	}
	user.FillAvatarURL()
	return &user, nil
}

func (s *UserService) GetManagedUserInfo(operatorID, targetUserID uint) (*model.SysUser, error) {
	return core.EnsureUserManageable(operatorID, targetUserID)
}

func (s *UserService) GetUserList(operatorID uint, req *request.UserListRequest) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64

	scope, err := core.ResolveUserDataScope(operatorID)
	if err != nil {
		return nil, 0, err
	}

	db := global.DB.Model(&model.SysUser{})
	db = core.ApplyUserDataScope(db, scope, "sys_user")

	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.Gender != nil {
		db = db.Where("gender = ?", *req.Gender)
	}
	if req.UnassignedDept {
		db = db.Where("sys_user.dept_id IS NULL OR sys_user.dept_id = 0")
	} else if req.DeptId != nil && *req.DeptId > 0 {
		deptIDs, err := core.GetDeptAndDescendantIDs([]uint{uint(*req.DeptId)})
		if err != nil {
			return nil, 0, err
		}
		if len(deptIDs) == 0 {
			return []model.SysUser{}, 0, nil
		}
		db = db.Where("sys_user.dept_id IN ?", deptIDs)
	}
	if req.RoleId != nil && *req.RoleId > 0 {
		db = db.Where("id IN (SELECT sys_user_id FROM sys_user_role WHERE sys_role_id = ?)", *req.RoleId)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	if err := db.Preload("Roles").Preload("Dept").Preload("AvatarFile").Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	for i := range users {
		users[i].FillAvatarURL()
	}

	return users, total, nil
}

func (s *UserService) CreateUser(operatorID uint, req *request.CreateUserRequest) error {
	if err := validateUserGender(req.Gender); err != nil {
		return err
	}
	if err := core.EnsureDeptManageable(operatorID, req.DeptID); err != nil {
		return err
	}

	var count int64
	global.DB.Model(&model.SysUser{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := model.SysUser{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
		DeptID:   req.DeptID,
	}

	if req.AvatarFileID > 0 {
		var file model.SysFile
		if err := global.DB.First(&file, req.AvatarFileID).Error; err != nil {
			return errors.New("头像文件不存在")
		}
		user.AvatarFileID = file.ID
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if len(req.RoleIds) > 0 {
			var roles []model.SysRole
			if err := tx.Where("id IN ?", req.RoleIds).Find(&roles).Error; err != nil {
				return err
			}
			if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *UserService) UpdateUser(operatorID, id uint, req *request.UpdateUserRequest) error {
	if err := validateUserGender(req.Gender); err != nil {
		return err
	}
	user, err := core.EnsureUserManageable(operatorID, id)
	if err != nil {
		return err
	}

	allowLegacyDeptRetain := false
	if req.DeptID == user.DeptID {
		isLeaf, err := core.IsDeptLeaf(user.DeptID)
		if err != nil {
			return err
		}
		allowLegacyDeptRetain = !isLeaf
	}
	if !allowLegacyDeptRetain {
		if err := core.EnsureDeptManageable(operatorID, req.DeptID); err != nil {
			return err
		}
	}

	updates := map[string]interface{}{
		"nickname": req.Nickname,
		"gender":   req.Gender,
		"email":    req.Email,
		"phone":    req.Phone,
		"status":   req.Status,
		"dept_id":  req.DeptID,
	}

	var oldRoleIds []uint
	global.DB.Table("sys_user_role").Where("sys_user_id = ?", id).Pluck("sys_role_id", &oldRoleIds)

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		targetUser := model.SysUser{BaseModel: model.BaseModel{ID: id}}

		if req.AvatarFileID > 0 {
			var file model.SysFile
			if err := tx.First(&file, req.AvatarFileID).Error; err != nil {
				return errors.New("头像文件不存在")
			}
			updates["avatar_file_id"] = file.ID
		} else if req.Avatar != "" {
			updates["avatar"] = req.Avatar
			updates["avatar_file_id"] = 0
		}

		if err := tx.Model(&model.SysUser{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		var roles []model.SysRole
		if len(req.RoleIds) > 0 {
			if err := tx.Where("id IN ?", req.RoleIds).Find(&roles).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&targetUser).Association("Roles").Replace(roles); err != nil {
			return err
		}

		return nil
	})
	if err == nil {
		core.Default.ClearUserInfoCache(id)

		newRoleIds := make([]uint, len(req.RoleIds))
		copy(newRoleIds, req.RoleIds)
		sort.Slice(oldRoleIds, func(i, j int) bool { return oldRoleIds[i] < oldRoleIds[j] })
		sort.Slice(newRoleIds, func(i, j int) bool { return newRoleIds[i] < newRoleIds[j] })

		rolesChanged := len(oldRoleIds) != len(newRoleIds)
		if !rolesChanged {
			for i := range oldRoleIds {
				if oldRoleIds[i] != newRoleIds[i] {
					rolesChanged = true
					break
				}
			}
		}

		if rolesChanged {
			_ = utils.RemoveUserToken(id)
		}
	}
	return err
}

func (s *UserService) DeleteUser(operatorID, id uint) error {
	if _, err := core.EnsureUserManageable(operatorID, id); err != nil {
		return err
	}
	return s.deleteUserByID(id)
}

func (s *UserService) deleteUserByID(id uint) error {
	var user model.SysUser
	if err := global.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	boundProfiles := global.Profiles.GetUserBoundProfiles(id)
	if len(boundProfiles) > 0 {
		return fmt.Errorf("该用户已绑定身份：%s，无法删除", strings.Join(boundProfiles, "、"))
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}
		deletedUsername := fmt.Sprintf("%s_deleted_%d_%d", user.Username, user.ID, time.Now().Unix())
		if err := tx.Model(&user).Update("username", deletedUsername).Error; err != nil {
			return err
		}
		return tx.Delete(&user).Error
	})
	if err == nil {
		core.Default.ClearUserInfoCache(id)
	}
	return err
}

func (s *UserService) BatchDeleteUsers(operatorID uint, ids []uint) (int, []string) {
	var successCount int
	var failedMsgs []string

	normalized := core.NormalizeIDs(ids)
	if _, err := core.EnsureUsersManageable(operatorID, normalized); err != nil {
		return 0, []string{err.Error()}
	}

	for _, id := range normalized {
		if err := s.deleteUserByID(id); err != nil {
			failedMsgs = append(failedMsgs, fmt.Sprintf("ID %d: %s", id, err.Error()))
		} else {
			successCount++
		}
	}

	return successCount, failedMsgs
}

func (s *UserService) UpdateUserStatus(operatorID, id uint, status int) error {
	if status != 0 && status != 1 {
		return errors.New("状态值无效")
	}
	if _, err := core.EnsureUserManageable(operatorID, id); err != nil {
		return err
	}

	err := global.DB.Model(&model.SysUser{}).Where("id = ?", id).Update("status", status).Error
	if err == nil {
		core.Default.ClearUserInfoCache(id)
	}
	return err
}

func (s *UserService) BatchUpdateUserStatus(operatorID uint, req *request.BatchUserStatusRequest) error {
	ids := core.NormalizeIDs(req.Ids)
	if len(ids) == 0 {
		return errors.New("请选择要修改状态的用户")
	}

	users, err := core.EnsureUsersManageable(operatorID, ids)
	if err != nil {
		return err
	}

	if err := validateBatchUserStatusTargets(ids, req.Status, operatorID, users); err != nil {
		return err
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.SysUser{}).
			Where("id IN ?", ids).
			Updates(map[string]interface{}{"status": req.Status}).Error; err != nil {
			return err
		}

		if req.Status == 0 {
			for _, id := range ids {
				if err := utils.RemoveUserToken(id); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, id := range ids {
		core.Default.ClearUserInfoCache(id)
	}
	return nil
}

func (s *UserService) managedUserDefaultPassword() string {
	config, err := configsvc.Default.GetConfigByKey(userDefaultPasswordConfigKey)
	if err != nil {
		return userDefaultPasswordFallback
	}

	password := strings.TrimSpace(config.Value)
	if password == "" {
		return userDefaultPasswordFallback
	}

	return password
}

func (s *UserService) ResetManagedUserPassword(operatorID, id uint) error {
	if _, err := core.EnsureUserManageable(operatorID, id); err != nil {
		return err
	}
	return s.ResetPassword(id, s.managedUserDefaultPassword())
}

func (s *UserService) BatchResetManagedUserPasswords(operatorID uint, ids []uint) error {
	normalized := core.NormalizeIDs(ids)
	if len(normalized) == 0 {
		return errors.New("请选择要重置密码的用户")
	}

	if _, err := core.EnsureUsersManageable(operatorID, normalized); err != nil {
		return err
	}

	return s.ResetUsersPassword(normalized, s.managedUserDefaultPassword())
}

func (s *UserService) ForceOffline(operatorID, id uint) error {
	if operatorID == id {
		return errors.New("不能强制下线自己")
	}
	if _, err := core.EnsureUserManageable(operatorID, id); err != nil {
		return err
	}
	return utils.RemoveUserToken(id)
}

func (s *UserService) ResetPassword(id uint, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	err = global.DB.Model(&model.SysUser{}).Where("id = ?", id).Update("password", hashedPassword).Error
	if err == nil {
		core.Default.ClearUserInfoCache(id)
	}
	return err
}

func (s *UserService) ResetUsersPassword(ids []uint, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	if err := global.DB.Model(&model.SysUser{}).Where("id IN ?", ids).Update("password", hashedPassword).Error; err != nil {
		return err
	}

	for _, id := range ids {
		core.Default.ClearUserInfoCache(id)
	}

	return nil
}

func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	var user model.SysUser
	if err := global.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return global.DB.Model(&user).Update("password", hashedPassword).Error
}

func (s *UserService) UpdateProfile(id uint, req *request.UpdateProfileRequest) error {
	updates := map[string]interface{}{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
	}
	err := global.DB.Model(&model.SysUser{}).Where("id = ?", id).Updates(updates).Error
	if err == nil {
		core.Default.ClearUserInfoCache(id)
	}
	return err
}

func (s *UserService) UpdateAvatar(id uint, fileID uint) error {
	var file model.SysFile
	if err := global.DB.First(&file, fileID).Error; err != nil {
		return errors.New("文件不存在")
	}

	err := global.DB.Model(&model.SysUser{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"avatar_file_id": file.ID,
		}).Error
	if err == nil {
		core.Default.ClearUserInfoCache(id)
	}
	return err
}

func (s *UserService) Register(username, password, email string) error {
	var count int64
	global.DB.Model(&model.SysUser{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	global.DB.Model(&model.SysUser{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return errors.New("邮箱已被使用")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := model.SysUser{
		Username: username,
		Password: hashedPassword,
		Nickname: username,
		Email:    email,
		Status:   1,
	}

	if config, err := configsvc.Default.GetConfigByKey("register_logo"); err == nil && config.Value != "" {
		var file model.SysFile
		if err := global.DB.Where("url = ? AND status = ?", config.Value, 1).First(&file).Error; err == nil {
			user.AvatarFileID = file.ID
		}
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		var role model.SysRole
		if err := tx.Where("code = ?", "user").First(&role).Error; err == nil {
			if err := tx.Model(&user).Association("Roles").Append(&role); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *UserService) GetUserByEmail(email string) (*model.SysUser, error) {
	var user model.SysUser
	if err := global.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByUserName(username string) (*model.SysUser, error) {
	var user model.SysUser
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) IsEmailVerificationRequired() bool {
	config, err := configsvc.Default.GetConfigByKey("register_email_verify")
	if err != nil {
		return false
	}
	return config.Value == "1" || config.Value == "true"
}
