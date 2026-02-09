package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"go-base-server/global"
	"go-base-server/model"
	"go-base-server/model/request"
)

type RoleService struct{}

var Role = new(RoleService)

// 获取角色列表
func (s *RoleService) GetRoleList() ([]model.SysRole, error) {
	var roles []model.SysRole
	if err := global.DB.Order("sort ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// 获取角色详情
func (s *RoleService) GetRole(id uint) (*model.SysRole, error) {
	var role model.SysRole
	if err := global.DB.Preload("Menus").Preload("Apis").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// 创建角色
func (s *RoleService) CreateRole(req *request.CreateRoleRequest) error {
	// 检查角色编码是否存在
	var count int64
	global.DB.Model(&model.SysRole{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return errors.New("角色编码已存在")
	}

	role := model.SysRole{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}

	return global.DB.Create(&role).Error
}

// 更新角色
func (s *RoleService) UpdateRole(id uint, req *request.UpdateRoleRequest) error {
	var role model.SysRole
	if err := global.DB.First(&role, id).Error; err != nil {
		return errors.New("角色不存在")
	}

	// 检查角色编码是否被其他角色使用
	if req.Code != "" && req.Code != role.Code {
		var count int64
		global.DB.Model(&model.SysRole{}).Where("code = ? AND id != ?", req.Code, id).Count(&count)
		if count > 0 {
			return errors.New("角色编码已存在")
		}
	}

	//updates := map[string]interface{}{
	//	"name":   req.Name,
	//	"code":   req.Code,
	//	"sort":   req.Sort,
	//	"status": req.Status,
	//	"remark": req.Remark,
	//}
	role.Name = req.Name
	role.Code = req.Code
	role.Sort = req.Sort
	role.Status = req.Status
	role.Remark = req.Remark
	err := global.DB.Save(&role).Error
	if err == nil {
		Cache.ClearUserCacheByRoleID(id) // 清除该角色用户的缓存
	}
	return err
}

// 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	var role model.SysRole
	if err := global.DB.First(&role, id).Error; err != nil {
		return errors.New("角色不存在")
	}

	// 检查是否有用户绑定该角色
	var userCount int64
	global.DB.Table("sys_user_role").Where("sys_role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		return errors.New("该角色已绑定用户，无法删除")
	}

	// 先清除缓存（删除前查询用户）
	Cache.ClearUserCacheByRoleID(id)

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 清除角色关联
		if err := tx.Model(&role).Association("Menus").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Apis").Clear(); err != nil {
			return err
		}
		// 软删除前修改 code，避免唯一索引冲突
		deletedCode := fmt.Sprintf("%s_deleted_%d_%d", role.Code, role.ID, time.Now().Unix())
		if err := tx.Model(&role).Update("code", deletedCode).Error; err != nil {
			return err
		}
		return tx.Delete(&role).Error
	})
}

// 分配菜单
func (s *RoleService) AssignMenus(roleID uint, menuIds []uint) error {
	var role model.SysRole
	if err := global.DB.First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}

	// 先清除旧关联
	if err := global.DB.Model(&role).Association("Menus").Clear(); err != nil {
		return err
	}

	// 如果有新的菜单 ID，添加关联
	if len(menuIds) > 0 {
		var menus []model.SysMenu
		if err := global.DB.Where("id IN ?", menuIds).Find(&menus).Error; err != nil {
			return err
		}
		if err := global.DB.Model(&role).Association("Menus").Append(&menus); err != nil {
			return err
		}
	}

	// 清除该角色用户的缓存
	Cache.ClearUserCacheByRoleID(roleID)
	return nil
}

// 分配API
func (s *RoleService) AssignApis(roleID uint, apiIds []uint) error {
	var role model.SysRole
	if err := global.DB.First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}

	// 先清除旧关联
	if err := global.DB.Model(&role).Association("Apis").Clear(); err != nil {
		return err
	}

	// 如果有新的API ID，添加关联
	var apis []model.SysApi
	if len(apiIds) > 0 {
		if err := global.DB.Where("id IN ?", apiIds).Find(&apis).Error; err != nil {
			return err
		}
		if err := global.DB.Model(&role).Association("Apis").Append(&apis); err != nil {
			return err
		}
	}

	// 更新Casbin策略
	return s.UpdateCasbinPolicy(role.Code, apis)
}

// 更新Casbin策略
func (s *RoleService) UpdateCasbinPolicy(roleCode string, apis []model.SysApi) error {
	// 删除旧策略
	_, err := global.Enforcer.RemoveFilteredPolicy(0, roleCode)
	if err != nil {
		return err
	}

	// 添加新策略
	var policies [][]string
	for _, api := range apis {
		policies = append(policies, []string{roleCode, api.Path, api.Method})
	}

	if len(policies) > 0 {
		_, err = global.Enforcer.AddPolicies(policies)
		if err != nil {
			return err
		}
	}

	return global.Enforcer.SavePolicy()
}
