package service

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
	"server/utils"
)

type UserService struct{}

var User = new(UserService)

func normalizeUserIDs(ids []uint) []uint {
	uniqueIDs := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}
	return uniqueIDs
}

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

// 获取用户选项列表（轻量级，仅返回 id/username/nickname）
func (s *UserService) GetUserOptions(operatorID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	scope, err := ResolveUserDataScope(operatorID)
	if err != nil {
		return nil, err
	}

	db := global.DB.Model(&model.SysUser{})
	db = ApplyUserDataScope(db, scope, "sys_user")
	err = db.
		Select("id, username, nickname").
		Where("status = ?", 1).
		Order("id ASC").
		Find(&results).Error
	return results, err
}

// 用户登录
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

// 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*model.SysUser, error) {
	var user model.SysUser
	if err := global.DB.Preload("Roles").Preload("Dept").Preload("AvatarFile").First(&user, userID).Error; err != nil {
		return nil, err
	}
	user.FillAvatarURL()
	return &user, nil
}

// 获取当前操作者可访问的用户详情
func (s *UserService) GetManagedUserInfo(operatorID, targetUserID uint) (*model.SysUser, error) {
	return EnsureUserManageable(operatorID, targetUserID)
}

// 获取用户列表
func (s *UserService) GetUserList(operatorID uint, req *request.UserListRequest) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64

	scope, err := ResolveUserDataScope(operatorID)
	if err != nil {
		return nil, 0, err
	}

	db := global.DB.Model(&model.SysUser{})
	db = ApplyUserDataScope(db, scope, "sys_user")

	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.UnassignedDept {
		db = db.Where("sys_user.dept_id IS NULL OR sys_user.dept_id = 0")
	} else if req.DeptId != nil && *req.DeptId > 0 {
		deptIDs, err := getDeptAndDescendantIDs([]uint{uint(*req.DeptId)})
		if err != nil {
			return nil, 0, err
		}
		if len(deptIDs) == 0 {
			return []model.SysUser{}, 0, nil
		}
		db = db.Where("sys_user.dept_id IN ?", deptIDs)
	}
	// 按角色ID过滤
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

	// 填充头像URL
	for i := range users {
		users[i].FillAvatarURL()
	}

	return users, total, nil
}

// 创建用户
func (s *UserService) CreateUser(operatorID uint, req *request.CreateUserRequest) error {
	if err := EnsureDeptManageable(operatorID, req.DeptID); err != nil {
		return err
	}

	// 检查用户名是否存在
	var count int64
	global.DB.Model(&model.SysUser{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := model.SysUser{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
		DeptID:   req.DeptID,
	}

	// 头像：优先使用文件ID，其次兼容直接传URL
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

		// 分配角色
		global.Log.Info("CreateUser: RoleIds=", req.RoleIds, ", UserID=", user.ID)
		if len(req.RoleIds) > 0 {
			var roles []model.SysRole
			if err := tx.Where("id IN ?", req.RoleIds).Find(&roles).Error; err != nil {
				return err
			}
			global.Log.Info("CreateUser: Found roles=", len(roles))
			if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}

		return nil
	})
}

// 更新用户
func (s *UserService) UpdateUser(operatorID, id uint, req *request.UpdateUserRequest) error {
	user, err := EnsureUserManageable(operatorID, id)
	if err != nil {
		return err
	}

	allowLegacyDeptRetain := false
	if req.DeptID == user.DeptID {
		isLeaf, err := IsDeptLeaf(user.DeptID)
		if err != nil {
			return err
		}
		allowLegacyDeptRetain = !isLeaf
	}
	if !allowLegacyDeptRetain {
		if err := EnsureDeptManageable(operatorID, req.DeptID); err != nil {
			return err
		}
	}

	updates := map[string]interface{}{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
		"status":   req.Status,
		"dept_id":  req.DeptID,
	}

	// 获取用户当前角色ID列表（用于判断角色是否变化）
	var oldRoleIds []uint
	global.DB.Table("sys_user_role").Where("sys_user_id = ?", id).Pluck("sys_role_id", &oldRoleIds)

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		targetUser := model.SysUser{BaseModel: model.BaseModel{ID: id}}

		// 头像：优先使用文件ID
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

		// 更新角色
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
		Cache.ClearUserInfoCache(id) // 清除缓存

		// 检查角色是否变化，如果变化则让用户 Token 失效
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
			// 角色变化，让用户 Token 失效，强制重新登录
			_ = utils.RemoveUserToken(id)
		}
	}
	return err
}

// 删除用户
func (s *UserService) DeleteUser(operatorID, id uint) error {
	if _, err := EnsureUserManageable(operatorID, id); err != nil {
		return err
	}
	return s.deleteUserByID(id)
}

func (s *UserService) deleteUserByID(id uint) error {
	var user model.SysUser
	if err := global.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 检查用户是否绑定了身份
	boundProfiles := global.Profiles.GetUserBoundProfiles(id)
	if len(boundProfiles) > 0 {
		return fmt.Errorf("该用户已绑定身份：%s，无法删除", strings.Join(boundProfiles, "、"))
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 删除用户角色关联
		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}
		// 软删除前修改 username，避免唯一索引冲突
		deletedUsername := fmt.Sprintf("%s_deleted_%d_%d", user.Username, user.ID, time.Now().Unix())
		if err := tx.Model(&user).Update("username", deletedUsername).Error; err != nil {
			return err
		}
		// 删除用户
		return tx.Delete(&user).Error
	})
	if err == nil {
		Cache.ClearUserInfoCache(id) // 清除缓存
	}
	return err
}

// 批量删除用户
func (s *UserService) BatchDeleteUsers(operatorID uint, ids []uint) (int, []string) {
	var successCount int
	var failedMsgs []string

	normalized := normalizeUserIDs(ids)
	if _, err := EnsureUsersManageable(operatorID, normalized); err != nil {
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

// 修改用户状态
func (s *UserService) UpdateUserStatus(operatorID, id uint, status int) error {
	if status != 0 && status != 1 {
		return errors.New("状态值无效")
	}
	if _, err := EnsureUserManageable(operatorID, id); err != nil {
		return err
	}

	err := global.DB.Model(&model.SysUser{}).Where("id = ?", id).Update("status", status).Error
	if err == nil {
		Cache.ClearUserInfoCache(id)
	}
	return err
}

// 批量修改用户状态
func (s *UserService) BatchUpdateUserStatus(operatorID uint, req *request.BatchUserStatusRequest) error {
	ids := normalizeUserIDs(req.Ids)
	if len(ids) == 0 {
		return errors.New("请选择要修改状态的用户")
	}

	users, err := EnsureUsersManageable(operatorID, ids)
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
		Cache.ClearUserInfoCache(id)
	}
	return nil
}

// 管理端重置密码
func (s *UserService) ResetManagedUserPassword(operatorID, id uint, password string) error {
	if _, err := EnsureUserManageable(operatorID, id); err != nil {
		return err
	}
	return s.ResetPassword(id, password)
}

// 管理端强制用户下线
func (s *UserService) ForceOffline(operatorID, id uint) error {
	if operatorID == id {
		return errors.New("不能强制下线自己")
	}
	if _, err := EnsureUserManageable(operatorID, id); err != nil {
		return err
	}
	return utils.RemoveUserToken(id)
}

// 重置密码
func (s *UserService) ResetPassword(id uint, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	return global.DB.Model(&model.SysUser{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// 修改密码
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

// 更新个人资料
func (s *UserService) UpdateProfile(id uint, req *request.UpdateProfileRequest) error {
	updates := map[string]interface{}{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
	}
	err := global.DB.Model(&model.SysUser{}).Where("id = ?", id).Updates(updates).Error
	if err == nil {
		Cache.ClearUserInfoCache(id)
	}
	return err
}

// 更新头像（使用文件ID）
func (s *UserService) UpdateAvatar(id uint, fileID uint) error {
	// 检查文件是否存在
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
		Cache.ClearUserInfoCache(id)
	}
	return err
}

// 注册用户
func (s *UserService) Register(username, password, email string) error {
	// 检查用户名是否存在
	var count int64
	global.DB.Model(&model.SysUser{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否存在
	global.DB.Model(&model.SysUser{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return errors.New("邮箱已被使用")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := model.SysUser{
		Username: username,
		Password: hashedPassword,
		Nickname: username,
		Email:    email,
		Status:   1, // 默认启用
	}

	// 设置默认头像
	if config, err := Config.GetConfigByKey("register_logo"); err == nil && config.Value != "" {
		var file model.SysFile
		if err := global.DB.Where("url = ? AND status = ?", config.Value, 1).First(&file).Error; err == nil {
			user.AvatarFileID = file.ID
		}
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 分配默认角色（user角色）
		var role model.SysRole
		if err := tx.Where("code = ?", "user").First(&role).Error; err == nil {
			if err := tx.Model(&user).Association("Roles").Append(&role); err != nil {
				return err
			}
		}

		return nil
	})
}

// 根据邮箱获取用户
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

// 根据邮箱获取用户
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

// 检查注册是否需要邮箱验证
func (s *UserService) IsEmailVerificationRequired() bool {
	config, err := Config.GetConfigByKey("register_email_verify")
	if err != nil {
		return false // 默认不需要
	}
	return config.Value == "1" || config.Value == "true"
}
