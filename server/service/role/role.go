package role

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
	"server/service/core"
)

type RoleService struct{}

var Default = &RoleService{}

// 获取角色列表
func (s *RoleService) GetRoleList() ([]model.SysRole, error) {
	var roles []model.SysRole
	if err := global.DB.Preload("Depts").Order("sort ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// 获取角色详情
func (s *RoleService) GetRole(id uint) (*model.SysRole, error) {
	var role model.SysRole
	if err := global.DB.Preload("Depts").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// 创建角色
func (s *RoleService) CreateRole(req *request.CreateRoleRequest) error {
	dataScope := req.DataScope
	if dataScope == 0 {
		dataScope = model.DataScopeAll
	}
	if err := validateRoleDataScope(dataScope, req.DeptIds); err != nil {
		return err
	}

	role := model.SysRole{
		Name:      req.Name,
		Code:      req.Code,
		DataScope: dataScope,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
	}
	if err := global.DB.Create(&role).Error; err != nil {
		return err
	}

	if dataScope == model.DataScopeCustom && len(req.DeptIds) > 0 {
		var depts []model.SysDept
		global.DB.Where("id IN ?", req.DeptIds).Find(&depts)
		global.DB.Model(&role).Association("Depts").Append(depts)
	}

	return nil
}

// 更新角色
func (s *RoleService) UpdateRole(id uint, req *request.UpdateRoleRequest) error {
	var role model.SysRole
	if err := global.DB.First(&role, id).Error; err != nil {
		return err
	}

	dataScope := req.DataScope
	if dataScope == 0 {
		role.DataScope = model.DataScopeAll
	} else {
		role.DataScope = dataScope
	}
	if err := validateRoleDataScope(role.DataScope, req.DeptIds); err != nil {
		return err
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"code":       req.Code,
		"data_scope": role.DataScope,
		"sort":       req.Sort,
		"status":     req.Status,
		"remark":     req.Remark,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&role).Updates(updates).Error; err != nil {
			return err
		}
		if err := replaceRoleCustomDepts(tx, &role, role.DataScope, req.DeptIds); err != nil {
			return err
		}
		core.Default.ClearUserCacheByRoleID(id)
		return nil
	})
}

// 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	var role model.SysRole
	if err := global.DB.First(&role, id).Error; err != nil {
		return err
	}
	if err := global.DB.Delete(&role).Error; err != nil {
		return err
	}
	core.Default.ClearUserCacheByRoleID(id)
	return nil
}

// 分配菜单
func (s *RoleService) AssignMenus(roleID uint, menuIDs []uint) error {
	var role model.SysRole
	if err := global.DB.First(&role, roleID).Error; err != nil {
		return err
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&role).Association("Menus").Replace(menuIDs); err != nil {
			return err
		}
		core.Default.ClearUserCacheByRoleID(roleID)
		return nil
	})
}

// 分配接口
func (s *RoleService) AssignApis(roleID uint, apiIDs []uint) error {
	var role model.SysRole
	if err := global.DB.First(&role, roleID).Error; err != nil {
		return err
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&role).Association("Apis").Replace(apiIDs); err != nil {
			return err
		}
		core.Default.ClearUserCacheByRoleID(roleID)
		return nil
	})
}

func validateRoleDataScope(dataScope int, deptIDs []uint) error {
	switch dataScope {
	case model.DataScopeAll, model.DataScopeCustom, model.DataScopeDept, model.DataScopeDeptAndChildren, model.DataScopeSelf:
		// valid
	default:
		return fmt.Errorf("无效的数据权限范围: %d", dataScope)
	}
	if dataScope == model.DataScopeCustom && len(core.NormalizeIDs(deptIDs)) == 0 {
		return errors.New("自定义数据权限必须选择部门")
	}
	return nil
}

func replaceRoleCustomDepts(tx *gorm.DB, role *model.SysRole, dataScope int, deptIDs []uint) error {
	if dataScope != model.DataScopeCustom {
		return tx.Model(role).Association("Depts").Replace([]model.SysDept{})
	}
	ids := core.NormalizeIDs(deptIDs)
	if len(ids) == 0 {
		return tx.Model(role).Association("Depts").Replace([]model.SysDept{})
	}
	var depts []model.SysDept
	if err := tx.Where("id IN ?", ids).Find(&depts).Error; err != nil {
		return err
	}
	return tx.Model(role).Association("Depts").Replace(depts)
}
