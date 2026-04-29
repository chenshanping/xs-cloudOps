package role

import (
	"errors"
	"fmt"
	"sort"

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

	if err := global.DB.
		Model(&model.SysRole{}).
		Select("sys_role.*, COALESCE((SELECT COUNT(*) FROM sys_user_role WHERE sys_user_role.sys_role_id = sys_role.id), 0) AS user_count").
		Preload("Depts").
		Order("sort ASC").
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// 获取角色详情
func (s *RoleService) GetRole(id uint) (*model.SysRole, error) {
	var role model.SysRole
	if err := global.DB.Preload("Menus").Preload("Apis").Preload("Depts").First(&role, id).Error; err != nil {
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
		Name:         req.Name,
		Code:         req.Code,
		DataScope:    dataScope,
		Sort:         req.Sort,
		Status:       req.Status,
		IsSuperAdmin: req.IsSuperAdmin,
		Remark:       req.Remark,
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
		"name":           req.Name,
		"code":           req.Code,
		"data_scope":     role.DataScope,
		"sort":           req.Sort,
		"status":         req.Status,
		"is_super_admin": req.IsSuperAdmin,
		"remark":         req.Remark,
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

	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		var userCount int64
		if err := tx.Table("sys_user_role").Where("sys_role_id = ?", id).Count(&userCount).Error; err != nil {
			return err
		}
		if userCount > 0 {
			return errors.New("当前角色下存在用户，无法删除")
		}
		return tx.Delete(&role).Error
	}); err != nil {
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
	ids := core.NormalizeIDs(menuIDs)
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if len(ids) == 0 {
			if err := tx.Model(&role).Association("Menus").Replace([]model.SysMenu{}); err != nil {
				return err
			}
			core.Default.ClearUserCacheByRoleID(roleID)
			return nil
		}

		menus, err := loadRoleMenusWithAncestors(tx, ids)
		if err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Menus").Replace(menus); err != nil {
			return err
		}
		core.Default.ClearUserCacheByRoleID(roleID)
		return nil
	})
}

func loadRoleMenusWithAncestors(tx *gorm.DB, ids []uint) ([]model.SysMenu, error) {
	var selectedMenus []model.SysMenu
	if err := tx.Where("id IN ?", ids).Find(&selectedMenus).Error; err != nil {
		return nil, err
	}
	if len(selectedMenus) == 0 {
		return []model.SysMenu{}, nil
	}

	menuMap := make(map[uint]model.SysMenu, len(selectedMenus))
	parentIDs := make([]uint, 0, len(selectedMenus))
	for _, menu := range selectedMenus {
		menuMap[menu.ID] = menu
		if menu.ParentID > 0 {
			parentIDs = append(parentIDs, menu.ParentID)
		}
	}

	for len(parentIDs) > 0 {
		parentIDs = core.NormalizeIDs(parentIDs)
		var parentMenus []model.SysMenu
		if err := tx.Where("id IN ?", parentIDs).Find(&parentMenus).Error; err != nil {
			return nil, err
		}

		nextParentIDs := make([]uint, 0, len(parentMenus))
		for _, parent := range parentMenus {
			if _, exists := menuMap[parent.ID]; !exists {
				menuMap[parent.ID] = parent
			}
			if parent.ParentID > 0 {
				nextParentIDs = append(nextParentIDs, parent.ParentID)
			}
		}
		parentIDs = nextParentIDs
	}

	menuIDs := make([]uint, 0, len(menuMap))
	for id := range menuMap {
		menuIDs = append(menuIDs, id)
	}
	sort.Slice(menuIDs, func(i, j int) bool {
		return menuIDs[i] < menuIDs[j]
	})

	menus := make([]model.SysMenu, 0, len(menuIDs))
	for _, id := range menuIDs {
		menus = append(menus, menuMap[id])
	}
	return menus, nil
}

// 分配接口
func (s *RoleService) AssignApis(roleID uint, apiIDs []uint) error {
	var role model.SysRole
	if err := global.DB.First(&role, roleID).Error; err != nil {
		return err
	}
	ids := core.NormalizeIDs(apiIDs)
	var assignedApis []model.SysApi
	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		if len(ids) == 0 {
			if err := tx.Model(&role).Association("Apis").Replace([]model.SysApi{}); err != nil {
				return err
			}
			return nil
		}

		var apis []model.SysApi
		if err := tx.Where("id IN ?", ids).Find(&apis).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Apis").Replace(apis); err != nil {
			return err
		}
		assignedApis = apis
		return nil
	}); err != nil {
		return err
	}

	if err := syncRoleApiPolicies(&role, assignedApis); err != nil {
		return err
	}
	core.Default.ClearUserCacheByRoleID(roleID)
	return nil
}

func syncRoleApiPolicies(role *model.SysRole, apis []model.SysApi) error {
	if global.Enforcer == nil {
		return nil
	}

	if _, err := global.Enforcer.RemoveFilteredPolicy(0, role.Code); err != nil {
		return err
	}

	if len(apis) > 0 {
		policies := make([][]string, 0, len(apis))
		for _, api := range apis {
			policies = append(policies, []string{role.Code, api.Path, api.Method})
		}
		if _, err := global.Enforcer.AddPolicies(policies); err != nil {
			return err
		}
	}

	return global.Enforcer.SavePolicy()
}

func (s *RoleService) HasSuperAdminRoleIDs(roleIDs []uint) (bool, error) {
	ids := core.NormalizeIDs(roleIDs)
	if len(ids) == 0 {
		return false, nil
	}

	var count int64
	if err := global.DB.Model(&model.SysRole{}).
		Where("id IN ?", ids).
		Where("is_super_admin = ?", true).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
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
