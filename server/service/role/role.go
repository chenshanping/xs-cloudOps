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
	if err := global.DB.
		Preload("Menus").
		Preload("Apis").
		Preload("Depts").
		Preload("FeatureDataScopes").
		Preload("FeatureDataScopes.Depts").
		First(&role, id).Error; err != nil {
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
	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		if len(ids) == 0 {
			if err := tx.Model(&role).Association("Menus").Replace([]model.SysMenu{}); err != nil {
				return err
			}
			return nil
		}

		menus, err := loadRoleMenusWithAncestors(tx, ids)
		if err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Menus").Replace(menus); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return s.syncRolePoliciesByRoleID(roleID)
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
	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		if len(ids) == 0 {
			if err := tx.Model(&role).Association("Apis").Replace([]model.SysApi{}); err != nil {
				return err
			}
			return nil
		}

		var apis []model.SysApi
		// 只关联 need_auth=true 的 API，公开接口不需要权限控制，避免幽灵关联与无用的 Casbin 规则。
		if err := tx.Where("id IN ? AND need_auth = ?", ids, true).Find(&apis).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Apis").Replace(apis); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return s.syncRolePoliciesByRoleID(roleID)
}

func (s *RoleService) AssignDataScopes(roleID uint, assignments []request.RoleFeatureDataScopeAssignment) error {
	var role model.SysRole
	if err := global.DB.First(&role, roleID).Error; err != nil {
		return err
	}
	_ = role

	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		normalizedAssignments, err := normalizePersistedRoleFeatureDataScopeAssignmentsTx(tx, roleID, assignments)
		if err != nil {
			return err
		}

		var existingScopes []model.SysRoleDataScope
		if err := tx.Preload("Depts").Where("role_id = ?", roleID).Find(&existingScopes).Error; err != nil {
			return err
		}
		for i := range existingScopes {
			if err := tx.Model(&existingScopes[i]).Association("Depts").Replace([]model.SysDept{}); err != nil {
				return err
			}
		}
		if err := tx.Unscoped().Where("role_id = ?", roleID).Delete(&model.SysRoleDataScope{}).Error; err != nil {
			return err
		}

		for _, assignment := range normalizedAssignments {
			scope := model.SysRoleDataScope{
				RoleID:       roleID,
				ResourceCode: assignment.ResourceCode,
				DataScope:    assignment.DataScope,
			}
			if err := tx.Create(&scope).Error; err != nil {
				return err
			}
			if err := replaceRoleFeatureScopeCustomDepts(tx, &scope, assignment.DataScope, assignment.DeptIds); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	core.Default.ClearUserCacheByRoleID(roleID)
	return nil
}

func (s *RoleService) SavePermissions(roleID uint, req *request.SaveRolePermissionsRequest) error {
	return s.savePermissions(roleID, req, true)
}

func (s *RoleService) SavePermissionsForOperator(roleID uint, req *request.SaveRolePermissionsRequest, operatorRoleIDs []uint) error {
	allowDirectApis, err := s.HasSuperAdminRoleIDs(operatorRoleIDs)
	if err != nil {
		return err
	}
	return s.savePermissions(roleID, req, allowDirectApis)
}

func (s *RoleService) savePermissions(roleID uint, req *request.SaveRolePermissionsRequest, allowDirectApis bool) error {
	var role model.SysRole
	if err := global.DB.First(&role, roleID).Error; err != nil {
		return err
	}

	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		normalizedAssignments, err := normalizePersistedRoleFeatureDataScopeAssignmentsTx(tx, roleID, req.Scopes)
		if err != nil {
			return err
		}

		menus, err := loadRoleMenusWithAncestors(tx, core.NormalizeIDs(req.MenuIds))
		if err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Menus").Replace(menus); err != nil {
			return err
		}

		if allowDirectApis {
			directApis, err := loadApisByIDs(tx, core.NormalizeIDs(req.DirectApiIds))
			if err != nil {
				return err
			}
			if err := tx.Model(&role).Association("Apis").Replace(directApis); err != nil {
				return err
			}
		}

		if err := replaceRoleFeatureDataScopes(tx, roleID, normalizedAssignments); err != nil {
			return err
		}

		effectiveApis, err := s.loadEffectiveRoleApisTx(tx, role.ID)
		if err != nil {
			return err
		}
		return syncRoleApiPoliciesTx(tx, role.Code, effectiveApis)
	}); err != nil {
		return err
	}

	if err := reloadRuntimePolicies(); err != nil {
		return err
	}
	core.Default.ClearUserCacheByRoleID(roleID)
	return nil
}

func (s *RoleService) SyncRolePoliciesForMenus(menuIDs []uint) error {
	ids := core.NormalizeIDs(menuIDs)
	if len(ids) == 0 {
		return nil
	}

	var roleIDs []uint
	if err := global.DB.Table("sys_role_menu").Distinct("sys_role_id").Where("sys_menu_id IN ?", ids).Pluck("sys_role_id", &roleIDs).Error; err != nil {
		return err
	}
	roleIDs = core.NormalizeIDs(roleIDs)
	for _, roleID := range roleIDs {
		if err := s.syncRolePoliciesByRoleID(roleID); err != nil {
			return err
		}
	}
	return nil
}

func (s *RoleService) syncRolePoliciesByRoleID(roleID uint) error {
	var role model.SysRole
	if err := global.DB.First(&role, roleID).Error; err != nil {
		return err
	}

	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		effectiveApis, err := s.loadEffectiveRoleApisTx(tx, role.ID)
		if err != nil {
			return err
		}
		return syncRoleApiPoliciesTx(tx, role.Code, effectiveApis)
	}); err != nil {
		return err
	}

	if err := reloadRuntimePolicies(); err != nil {
		return err
	}
	core.Default.ClearUserCacheByRoleID(roleID)
	return nil
}

func (s *RoleService) loadEffectiveRoleApisTx(tx *gorm.DB, roleID uint) ([]model.SysApi, error) {
	var role model.SysRole
	if err := tx.Preload("Apis").Preload("Menus.Apis").First(&role, roleID).Error; err != nil {
		return nil, err
	}

	apiMap := make(map[uint]model.SysApi)
	for _, api := range role.Apis {
		if !api.NeedAuth {
			continue
		}
		apiMap[api.ID] = api
	}
	for _, menu := range role.Menus {
		for _, api := range menu.Apis {
			if !api.NeedAuth {
				continue
			}
			apiMap[api.ID] = api
		}
	}

	apiIDs := make([]uint, 0, len(apiMap))
	for id := range apiMap {
		apiIDs = append(apiIDs, id)
	}
	sort.Slice(apiIDs, func(i, j int) bool { return apiIDs[i] < apiIDs[j] })

	apis := make([]model.SysApi, 0, len(apiIDs))
	for _, id := range apiIDs {
		apis = append(apis, apiMap[id])
	}
	return apis, nil
}

func loadApisByIDs(tx *gorm.DB, ids []uint) ([]model.SysApi, error) {
	if len(ids) == 0 {
		return []model.SysApi{}, nil
	}

	var apis []model.SysApi
	if err := tx.Where("id IN ?", ids).Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

func replaceRoleFeatureDataScopes(
	tx *gorm.DB,
	roleID uint,
	assignments []request.RoleFeatureDataScopeAssignment,
) error {
	var existingScopes []model.SysRoleDataScope
	if err := tx.Preload("Depts").Where("role_id = ?", roleID).Find(&existingScopes).Error; err != nil {
		return err
	}
	for i := range existingScopes {
		if err := tx.Model(&existingScopes[i]).Association("Depts").Replace([]model.SysDept{}); err != nil {
			return err
		}
	}
	if err := tx.Unscoped().Where("role_id = ?", roleID).Delete(&model.SysRoleDataScope{}).Error; err != nil {
		return err
	}

	for _, assignment := range assignments {
		scope := model.SysRoleDataScope{
			RoleID:       roleID,
			ResourceCode: assignment.ResourceCode,
			DataScope:    assignment.DataScope,
		}
		if err := tx.Create(&scope).Error; err != nil {
			return err
		}
		if err := replaceRoleFeatureScopeCustomDepts(tx, &scope, assignment.DataScope, assignment.DeptIds); err != nil {
			return err
		}
	}
	return nil
}

type casbinRuleRow struct {
	ID    uint   `gorm:"column:id;primaryKey"`
	PType string `gorm:"column:ptype"`
	V0    string `gorm:"column:v0"`
	V1    string `gorm:"column:v1"`
	V2    string `gorm:"column:v2"`
	V3    string `gorm:"column:v3"`
	V4    string `gorm:"column:v4"`
	V5    string `gorm:"column:v5"`
}

func (casbinRuleRow) TableName() string {
	return "casbin_rule"
}

func syncRoleApiPoliciesTx(tx *gorm.DB, roleCode string, apis []model.SysApi) error {
	if err := tx.Where("ptype = ? AND v0 = ?", "p", roleCode).Delete(&casbinRuleRow{}).Error; err != nil {
		return err
	}

	if len(apis) == 0 {
		return nil
	}

	rules := make([]casbinRuleRow, 0, len(apis))
	for _, api := range apis {
		rules = append(rules, casbinRuleRow{
			PType: "p",
			V0:    roleCode,
			V1:    api.Path,
			V2:    api.Method,
		})
	}
	return tx.Create(&rules).Error
}

func reloadRuntimePolicies() error {
	if global.Enforcer == nil {
		return nil
	}
	if err := global.Enforcer.LoadPolicy(); err != nil {
		return err
	}
	return nil
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

func validateRoleFeatureDataScope(resourceCode string, dataScope int, deptIDs []uint) error {
	if !core.IsSupportedDataScopeResource(resourceCode) {
		return fmt.Errorf("不支持的数据权限资源: %s", resourceCode)
	}
	if err := validateRoleDataScope(dataScope, deptIDs); err != nil {
		return err
	}
	if core.SupportsDataScopeForResource(resourceCode, dataScope) {
		return nil
	}

	switch dataScope {
	case model.DataScopeCustom:
		return fmt.Errorf("数据权限资源 %s 不支持自定义部门范围", resourceCode)
	case model.DataScopeDept:
		return fmt.Errorf("数据权限资源 %s 不支持本部门范围", resourceCode)
	case model.DataScopeDeptAndChildren:
		return fmt.Errorf("数据权限资源 %s 不支持本部门及下级范围", resourceCode)
	case model.DataScopeSelf:
		return fmt.Errorf("数据权限资源 %s 不支持仅本人范围", resourceCode)
	default:
		return fmt.Errorf("数据权限资源 %s 不支持该数据范围", resourceCode)
	}
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

func normalizeRoleFeatureDataScopeAssignments(
	assignments []request.RoleFeatureDataScopeAssignment,
) ([]request.RoleFeatureDataScopeAssignment, error) {
	normalized := make([]request.RoleFeatureDataScopeAssignment, 0, len(assignments))
	indexByCode := make(map[string]int, len(assignments))
	for _, assignment := range assignments {
		if err := validateRoleFeatureDataScope(assignment.ResourceCode, assignment.DataScope, assignment.DeptIds); err != nil {
			return nil, err
		}

		normalizedAssignment := request.RoleFeatureDataScopeAssignment{
			ResourceCode: assignment.ResourceCode,
			DataScope:    assignment.DataScope,
			DeptIds:      core.NormalizeIDs(assignment.DeptIds),
		}

		if idx, exists := indexByCode[assignment.ResourceCode]; exists {
			normalized[idx] = normalizedAssignment
			continue
		}

		indexByCode[assignment.ResourceCode] = len(normalized)
		normalized = append(normalized, normalizedAssignment)
	}

	return normalized, nil
}

func normalizePersistedRoleFeatureDataScopeAssignmentsTx(
	tx *gorm.DB,
	roleID uint,
	assignments []request.RoleFeatureDataScopeAssignment,
) ([]request.RoleFeatureDataScopeAssignment, error) {
	filteredAssignments := make([]request.RoleFeatureDataScopeAssignment, 0, len(assignments))
	for _, assignment := range assignments {
		if !core.IsSupportedDataScopeResource(assignment.ResourceCode) {
			continue
		}
		filteredAssignments = append(filteredAssignments, assignment)
	}

	normalizedAssignments, err := normalizeRoleFeatureDataScopeAssignments(filteredAssignments)
	if err != nil {
		return nil, err
	}

	var existingScopes []model.SysRoleDataScope
	if err := tx.Preload("Depts").Where("role_id = ?", roleID).Find(&existingScopes).Error; err != nil {
		return nil, err
	}

	preservedUnknownScopes := make([]request.RoleFeatureDataScopeAssignment, 0)
	for _, scope := range existingScopes {
		if core.IsSupportedDataScopeResource(scope.ResourceCode) {
			continue
		}
		preservedUnknownScopes = append(preservedUnknownScopes, request.RoleFeatureDataScopeAssignment{
			ResourceCode: scope.ResourceCode,
			DataScope:    scope.DataScope,
			DeptIds:      extractDeptIDs(scope.Depts),
		})
	}

	return append(normalizedAssignments, preservedUnknownScopes...), nil
}

func extractDeptIDs(depts []model.SysDept) []uint {
	ids := make([]uint, 0, len(depts))
	for _, dept := range depts {
		ids = append(ids, dept.ID)
	}
	return ids
}

func replaceRoleFeatureScopeCustomDepts(
	tx *gorm.DB,
	scope *model.SysRoleDataScope,
	dataScope int,
	deptIDs []uint,
) error {
	if dataScope != model.DataScopeCustom {
		return tx.Model(scope).Association("Depts").Replace([]model.SysDept{})
	}
	ids := core.NormalizeIDs(deptIDs)
	if len(ids) == 0 {
		return tx.Model(scope).Association("Depts").Replace([]model.SysDept{})
	}
	var depts []model.SysDept
	if err := tx.Where("id IN ?", ids).Find(&depts).Error; err != nil {
		return err
	}
	return tx.Model(scope).Association("Depts").Replace(depts)
}
