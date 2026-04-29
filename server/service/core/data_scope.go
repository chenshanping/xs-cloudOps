package core

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"server/global"
	"server/model"
)

// UserDataScope holds resolved data scope information for an operator.
type UserDataScope struct {
	OperatorID   uint
	OperatorDept uint
	All          bool
	AllowSelf    bool
	CreatorIDs   []uint
	DeptIDs      []uint
}

// NormalizeIDs deduplicates and removes zero IDs.
func NormalizeIDs(ids []uint) []uint {
	result := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

// ResolveUserDataScope resolves the default data scope for the given operator.
func ResolveUserDataScope(operatorID uint) (*UserDataScope, error) {
	return ResolveUserDataScopeForResource(operatorID, "")
}

// ResolveUserDataScopeForResource resolves the data scope for the given operator and resource.
func ResolveUserDataScopeForResource(operatorID uint, resourceCode string) (*UserDataScope, error) {
	var user model.SysUser
	query := global.DB.Preload("Roles.Depts")
	if resourceCode != "" {
		query = query.Preload("Roles.FeatureDataScopes", "resource_code = ?", resourceCode).
			Preload("Roles.FeatureDataScopes.Depts")
	}
	if err := query.First(&user, operatorID).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	return buildUserDataScope(&user, resourceCode)
}

func buildUserDataScope(user *model.SysUser, resourceCode string) (*UserDataScope, error) {
	scope := &UserDataScope{
		OperatorID:   user.ID,
		OperatorDept: user.DeptID,
	}

	creatorSet := make(map[uint]struct{})
	deptSet := make(map[uint]struct{})
	hasConfiguredScope := false
	allowSelf := false

	for _, role := range user.Roles {
		dataScope := role.DataScope
		scopeDepts := role.Depts
		if resourceCode != "" {
			if featureScope := findRoleFeatureDataScope(role.FeatureDataScopes, resourceCode); featureScope != nil {
				dataScope = featureScope.DataScope
				scopeDepts = featureScope.Depts
			}
		}
		if dataScope == 0 {
			dataScope = model.DataScopeAll
		}

		if role.Code == "admin" || role.Code == "super_admin" || dataScope == model.DataScopeAll {
			scope.All = true
			scope.AllowSelf = true
			scope.DeptIDs = nil
			return scope, nil
		}

		hasConfiguredScope = true

		switch dataScope {
		case model.DataScopeCustom:
			for _, dept := range scopeDepts {
				deptSet[dept.ID] = struct{}{}
			}
		case model.DataScopeDept:
			if user.DeptID > 0 {
				deptSet[user.DeptID] = struct{}{}
			}
		case model.DataScopeDeptAndChildren:
			if user.DeptID > 0 {
				deptIDs, err := GetDeptAndDescendantIDs([]uint{user.DeptID})
				if err != nil {
					return nil, err
				}
				for _, deptID := range deptIDs {
					deptSet[deptID] = struct{}{}
				}
			}
		case model.DataScopeSelf:
			if resourceCode == DataScopeResourceUserManagement {
				creatorSet[user.ID] = struct{}{}
				continue
			}
			allowSelf = true
		}
	}

	if !hasConfiguredScope {
		allowSelf = true
	}

	scope.AllowSelf = allowSelf
	scope.CreatorIDs = mapKeysToSortedSlice(creatorSet)
	scope.DeptIDs = mapKeysToSortedSlice(deptSet)
	return scope, nil
}

func findRoleFeatureDataScope(
	featureScopes []model.SysRoleDataScope,
	resourceCode string,
) *model.SysRoleDataScope {
	for i := range featureScopes {
		if featureScopes[i].ResourceCode == resourceCode {
			return &featureScopes[i]
		}
	}
	return nil
}

// ApplyUserDataScope applies the data scope filter to a query.
func ApplyUserDataScope(db *gorm.DB, scope *UserDataScope, userTable string) *gorm.DB {
	if scope == nil {
		return db.Where("1 = 0")
	}
	if scope.All {
		return db
	}

	conditions := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)

	if len(scope.DeptIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("%s.dept_id IN ?", userTable))
		args = append(args, scope.DeptIDs)
	}
	if len(scope.CreatorIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("%s.created_by IN ?", userTable))
		args = append(args, scope.CreatorIDs)
	}
	if scope.AllowSelf {
		conditions = append(conditions, fmt.Sprintf("%s.id = ?", userTable))
		args = append(args, scope.OperatorID)
	}

	if len(conditions) == 0 {
		return db.Where("1 = 0")
	}

	return db.Where(strings.Join(conditions, " OR "), args...)
}

func deptIDInScope(scope *UserDataScope, deptID uint) bool {
	if scope == nil || deptID == 0 {
		return false
	}
	if scope.All {
		return true
	}
	for _, allowedDeptID := range scope.DeptIDs {
		if allowedDeptID == deptID {
			return true
		}
	}
	return false
}

func ensureDeptInScopeForResource(
	operatorID,
	deptID uint,
	resourceCode string,
) (*model.SysDept, *UserDataScope, error) {
	if deptID == 0 {
		return nil, nil, errors.New("请选择所属部门")
	}

	var dept model.SysDept
	if err := global.DB.First(&dept, deptID).Error; err != nil {
		return nil, nil, errors.New("所属部门不存在")
	}

	scope, err := ResolveUserDataScopeForResource(operatorID, resourceCode)
	if err != nil {
		return nil, nil, err
	}
	if !deptIDInScope(scope, deptID) {
		return nil, nil, errors.New("无权操作该部门")
	}
	return &dept, scope, nil
}

// EnsureDeptManageable validates that the operator can manage the given department.
func EnsureDeptManageable(operatorID, deptID uint) error {
	return EnsureDeptManageableForResource(operatorID, deptID, "")
}

func EnsureDeptManageableForResource(operatorID, deptID uint, resourceCode string) error {
	if _, _, err := ensureDeptInScopeForResource(operatorID, deptID, resourceCode); err != nil {
		return err
	}

	var childCount int64
	if err := global.DB.Model(&model.SysDept{}).Where("parent_id = ?", deptID).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("存在下级部门的部门不能直接绑定用户")
	}
	return nil
}

func EnsureDeptAccessibleForResource(
	operatorID,
	deptID uint,
	resourceCode string,
) (*model.SysDept, error) {
	dept, _, err := ensureDeptInScopeForResource(operatorID, deptID, resourceCode)
	if err != nil {
		return nil, err
	}
	return dept, nil
}

func EnsureDeptParentManageableForResource(
	operatorID,
	parentID uint,
	resourceCode string,
) error {
	if parentID == 0 {
		scope, err := ResolveUserDataScopeForResource(operatorID, resourceCode)
		if err != nil {
			return err
		}
		if !scope.All {
			return errors.New("无权在顶级部门下创建或移动部门")
		}
		return nil
	}
	_, err := EnsureDeptAccessibleForResource(operatorID, parentID, resourceCode)
	return err
}

// IsDeptLeaf checks whether the given department has no children.
func IsDeptLeaf(deptID uint) (bool, error) {
	if deptID == 0 {
		return false, nil
	}

	var childCount int64
	if err := global.DB.Model(&model.SysDept{}).Where("parent_id = ?", deptID).Count(&childCount).Error; err != nil {
		return false, err
	}
	return childCount == 0, nil
}

// EnsureUserManageable validates and returns the user if it's within the operator's scope.
func EnsureUserManageable(operatorID, targetUserID uint) (*model.SysUser, error) {
	return EnsureUserManageableForResource(operatorID, targetUserID, "")
}

func EnsureUserManageableForResource(
	operatorID,
	targetUserID uint,
	resourceCode string,
) (*model.SysUser, error) {
	scope, err := ResolveUserDataScopeForResource(operatorID, resourceCode)
	if err != nil {
		return nil, err
	}

	var user model.SysUser
	query := global.DB.Preload("Roles").Preload("Dept").Preload("AvatarFile").Model(&model.SysUser{}).Where("sys_user.id = ?", targetUserID)
	query = ApplyUserDataScope(query, scope, "sys_user")
	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在或无权访问")
		}
		return nil, err
	}

	user.FillAvatarURL()
	return &user, nil
}

// EnsureUsersManageable validates that all the given user IDs are within the operator's scope.
func EnsureUsersManageable(operatorID uint, ids []uint) ([]model.SysUser, error) {
	return EnsureUsersManageableForResource(operatorID, ids, "")
}

func EnsureUsersManageableForResource(
	operatorID uint,
	ids []uint,
	resourceCode string,
) ([]model.SysUser, error) {
	normalized := NormalizeIDs(ids)
	if len(normalized) == 0 {
		return nil, errors.New("请选择要操作的用户")
	}

	scope, err := ResolveUserDataScopeForResource(operatorID, resourceCode)
	if err != nil {
		return nil, err
	}

	var users []model.SysUser
	query := global.DB.Preload("Roles").Model(&model.SysUser{}).Where("sys_user.id IN ?", normalized)
	query = ApplyUserDataScope(query, scope, "sys_user")
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	if len(users) != len(normalized) {
		return nil, errors.New("存在超出数据权限范围的用户")
	}

	return users, nil
}

// GetDeptAndDescendantIDs returns all department IDs that are the given root IDs
// or descendants of them, by scanning the ancestors column.
func GetDeptAndDescendantIDs(rootIDs []uint) ([]uint, error) {
	normalized := NormalizeIDs(rootIDs)
	if len(normalized) == 0 {
		return nil, nil
	}

	var depts []model.SysDept
	if err := global.DB.Model(&model.SysDept{}).Select("id", "ancestors").Find(&depts).Error; err != nil {
		return nil, err
	}

	rootSet := make(map[uint]struct{}, len(normalized))
	for _, id := range normalized {
		rootSet[id] = struct{}{}
	}

	resultSet := make(map[uint]struct{})
	for _, dept := range depts {
		if _, ok := rootSet[dept.ID]; ok {
			resultSet[dept.ID] = struct{}{}
			continue
		}
		for ancestorID := range rootSet {
			if ancestorsContainID(dept.Ancestors, ancestorID) {
				resultSet[dept.ID] = struct{}{}
				break
			}
		}
	}

	return mapKeysToSortedSlice(resultSet), nil
}

func ancestorsContainID(ancestors string, targetID uint) bool {
	if ancestors == "" {
		return false
	}

	for _, part := range strings.Split(ancestors, ",") {
		if strings.TrimSpace(part) == "" {
			continue
		}
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil {
			continue
		}
		if uint(id) == targetID {
			return true
		}
	}
	return false
}

func mapKeysToSortedSlice(items map[uint]struct{}) []uint {
	result := make([]uint, 0, len(items))
	for id := range items {
		result = append(result, id)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}
