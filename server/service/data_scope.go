package service

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

type UserDataScope struct {
	OperatorID   uint
	OperatorDept uint
	All          bool
	AllowSelf    bool
	DeptIDs      []uint
}

func ResolveUserDataScope(operatorID uint) (*UserDataScope, error) {
	var user model.SysUser
	if err := global.DB.Preload("Roles.Depts").First(&user, operatorID).Error; err != nil {
		return nil, err
	}

	scope := &UserDataScope{
		OperatorID:   operatorID,
		OperatorDept: user.DeptID,
	}

	deptSet := make(map[uint]struct{})
	hasConfiguredScope := false
	allowSelf := false

	for _, role := range user.Roles {
		dataScope := role.DataScope
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
			for _, dept := range role.Depts {
				deptSet[dept.ID] = struct{}{}
			}
		case model.DataScopeDept:
			if user.DeptID > 0 {
				deptSet[user.DeptID] = struct{}{}
			}
		case model.DataScopeDeptAndChildren:
			if user.DeptID > 0 {
				deptIDs, err := getDeptAndDescendantIDs([]uint{user.DeptID})
				if err != nil {
					return nil, err
				}
				for _, deptID := range deptIDs {
					deptSet[deptID] = struct{}{}
				}
			}
		case model.DataScopeSelf:
			allowSelf = true
		}
	}

	if !hasConfiguredScope {
		allowSelf = true
	}

	scope.AllowSelf = allowSelf
	scope.DeptIDs = mapKeysToSortedSlice(deptSet)
	return scope, nil
}

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
	if scope.AllowSelf {
		conditions = append(conditions, fmt.Sprintf("%s.id = ?", userTable))
		args = append(args, scope.OperatorID)
	}

	if len(conditions) == 0 {
		return db.Where("1 = 0")
	}

	return db.Where(strings.Join(conditions, " OR "), args...)
}

func EnsureDeptManageable(operatorID, deptID uint) error {
	if deptID == 0 {
		return errors.New("请选择所属部门")
	}

	var dept model.SysDept
	if err := global.DB.First(&dept, deptID).Error; err != nil {
		return errors.New("所属部门不存在")
	}

	scope, err := ResolveUserDataScope(operatorID)
	if err != nil {
		return err
	}

	allowed := scope.All
	if !allowed {
		for _, allowedDeptID := range scope.DeptIDs {
			if allowedDeptID == deptID {
				allowed = true
				break
			}
		}
	}
	if !allowed {
		return errors.New("无权操作该部门")
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

func EnsureUserManageable(operatorID, targetUserID uint) (*model.SysUser, error) {
	scope, err := ResolveUserDataScope(operatorID)
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

func EnsureUsersManageable(operatorID uint, ids []uint) ([]model.SysUser, error) {
	normalized := normalizeUserIDs(ids)
	if len(normalized) == 0 {
		return nil, errors.New("请选择要操作的用户")
	}

	scope, err := ResolveUserDataScope(operatorID)
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

func getDeptAndDescendantIDs(rootIDs []uint) ([]uint, error) {
	normalized := normalizeUserIDs(rootIDs)
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
