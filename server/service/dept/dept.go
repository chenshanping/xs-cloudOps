package dept

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
	"server/service/configsvc"
	"server/service/core"
)

type DeptService struct{}

var Default = &DeptService{}

func (s *DeptService) GetDeptTree() ([]model.SysDept, error) {
	var depts []model.SysDept
	if err := global.DB.Order("sort ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	return s.buildDeptTree(depts, 0), nil
}

func (s *DeptService) GetManageableDeptTree(operatorID uint) ([]model.SysDept, int64, error) {
	return s.GetManageableDeptTreeForResource(operatorID, core.DataScopeResourceDeptManagement)
}

func (s *DeptService) GetManageableDeptTreeForResource(operatorID uint, resourceCode string) ([]model.SysDept, int64, error) {
	scope, err := core.ResolveUserDataScopeForResource(operatorID, resourceCode)
	if err != nil {
		return nil, 0, err
	}

	var allDepts []model.SysDept
	if err := global.DB.Order("sort ASC, id ASC").Find(&allDepts).Error; err != nil {
		return nil, 0, err
	}

	directUserCountMap, unassignedCount, err := s.getVisibleDeptUserCounts(scope)
	if err != nil {
		return nil, 0, err
	}

	allowedManageableSet := make(map[uint]struct{})
	if scope.All {
		for _, dept := range allDepts {
			allowedManageableSet[dept.ID] = struct{}{}
		}
	} else {
		for _, deptID := range scope.DeptIDs {
			allowedManageableSet[deptID] = struct{}{}
		}
		for deptID := range directUserCountMap {
			if deptID == 0 {
				continue
			}
			allowedManageableSet[deptID] = struct{}{}
		}
	}

	includedSet := expandDeptIDsWithAncestors(allDepts, allowedManageableSet)
	filteredDepts := filterDeptsBySet(allDepts, includedSet)

	tree := s.buildDeptTree(filteredDepts, 0)
	s.decorateDeptTree(tree, directUserCountMap, allowedManageableSet)
	return tree, unassignedCount, nil
}

func (s *DeptService) GetManageableDeptTreeWithDefaultsForResource(operatorID uint, resourceCode string) ([]model.SysDept, int64, string, error) {
	tree, unassignedCount, err := s.GetManageableDeptTreeForResource(operatorID, resourceCode)
	if err != nil {
		return nil, 0, "", err
	}

	defaultAvatarURL, err := s.getDefaultUserAvatarURL()
	if err != nil {
		return nil, 0, "", err
	}

	return tree, unassignedCount, defaultAvatarURL, nil
}

func (s *DeptService) GetDept(id uint) (*model.SysDept, error) {
	var dept model.SysDept
	if err := global.DB.First(&dept, id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (s *DeptService) GetManagedDept(operatorID, id uint) (*model.SysDept, error) {
	dept, err := core.EnsureDeptAccessibleForResource(operatorID, id, core.DataScopeResourceDeptManagement)
	if err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *DeptService) CreateDept(req *request.CreateDeptRequest) error {
	ancestors, err := buildDeptAncestors(req.ParentID)
	if err != nil {
		return err
	}

	dept := model.SysDept{
		ParentID:  req.ParentID,
		Ancestors: ancestors,
		Name:      req.Name,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	err = global.DB.Create(&dept).Error
	if err == nil {
		core.Default.ClearAllUserInfoCache()
	}
	return err
}

func (s *DeptService) CreateManagedDept(operatorID uint, req *request.CreateDeptRequest) error {
	if err := core.EnsureDeptParentManageableForResource(operatorID, req.ParentID, core.DataScopeResourceDeptManagement); err != nil {
		return err
	}
	return s.CreateDept(req)
}

func (s *DeptService) UpdateDept(id uint, req *request.UpdateDeptRequest) error {
	var dept model.SysDept
	if err := global.DB.First(&dept, id).Error; err != nil {
		return errors.New("部门不存在")
	}

	if req.ParentID == id {
		return errors.New("上级部门不能选择自己")
	}

	newAncestors, err := buildDeptAncestors(req.ParentID)
	if err != nil {
		return err
	}

	if req.ParentID != 0 {
		var parent model.SysDept
		if err := global.DB.First(&parent, req.ParentID).Error; err != nil {
			return errors.New("上级部门不存在")
		}
		if ancestorsContainID(parent.Ancestors, id) {
			return errors.New("不能将部门移动到自己的下级部门下")
		}
	}

	oldLineage := buildDeptLineage(dept.Ancestors, dept.ID)
	newLineage := buildDeptLineage(newAncestors, dept.ID)

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"parent_id": req.ParentID,
			"ancestors": newAncestors,
			"name":      req.Name,
			"sort":      req.Sort,
			"status":    req.Status,
			"remark":    req.Remark,
		}
		if err := tx.Model(&dept).Updates(updates).Error; err != nil {
			return err
		}

		if oldLineage == newLineage {
			return nil
		}

		var descendants []model.SysDept
		if err := tx.Where("ancestors = ? OR ancestors LIKE ?", oldLineage, oldLineage+",%").Find(&descendants).Error; err != nil {
			return err
		}

		for _, child := range descendants {
			newChildAncestors := strings.Replace(child.Ancestors, oldLineage, newLineage, 1)
			if err := tx.Model(&model.SysDept{}).Where("id = ?", child.ID).Update("ancestors", newChildAncestors).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err == nil {
		core.Default.ClearAllUserInfoCache()
	}
	return err
}

func (s *DeptService) UpdateManagedDept(operatorID, id uint, req *request.UpdateDeptRequest) error {
	if _, err := core.EnsureDeptAccessibleForResource(operatorID, id, core.DataScopeResourceDeptManagement); err != nil {
		return err
	}
	if err := core.EnsureDeptParentManageableForResource(operatorID, req.ParentID, core.DataScopeResourceDeptManagement); err != nil {
		return err
	}
	return s.UpdateDept(id, req)
}

func (s *DeptService) DeleteDept(id uint) error {
	var dept model.SysDept
	if err := global.DB.First(&dept, id).Error; err != nil {
		return errors.New("部门不存在")
	}

	var childCount int64
	if err := global.DB.Model(&model.SysDept{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("存在下级部门，无法删除")
	}

	var userCount int64
	if err := global.DB.Model(&model.SysUser{}).Where("dept_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return errors.New("该部门下存在用户，无法删除")
	}

	var roleDeptCount int64
	if err := global.DB.Table("sys_role_dept").Where("sys_dept_id = ?", id).Count(&roleDeptCount).Error; err != nil {
		return err
	}
	var featureScopeDeptCount int64
	if err := global.DB.Table("sys_role_data_scope_dept").Where("sys_dept_id = ?", id).Count(&featureScopeDeptCount).Error; err != nil {
		return err
	}
	if roleDeptCount > 0 || featureScopeDeptCount > 0 {
		return errors.New("该部门已被角色数据范围引用，无法删除")
	}

	err := global.DB.Delete(&dept).Error
	if err == nil {
		core.Default.ClearAllUserInfoCache()
	}
	return err
}

func (s *DeptService) DeleteManagedDept(operatorID, id uint) error {
	if _, err := core.EnsureDeptAccessibleForResource(operatorID, id, core.DataScopeResourceDeptManagement); err != nil {
		return err
	}
	return s.DeleteDept(id)
}

func (s *DeptService) buildDeptTree(depts []model.SysDept, parentID uint) []model.SysDept {
	tree := make([]model.SysDept, 0)
	for _, dept := range depts {
		if dept.ParentID == parentID {
			dept.Children = s.buildDeptTree(depts, dept.ID)
			tree = append(tree, dept)
		}
	}

	sort.Slice(tree, func(i, j int) bool {
		if tree[i].Sort == tree[j].Sort {
			return tree[i].ID < tree[j].ID
		}
		return tree[i].Sort < tree[j].Sort
	})
	return tree
}

func (s *DeptService) decorateDeptTree(nodes []model.SysDept, directUserCountMap map[uint]int64, allowedManageableSet map[uint]struct{}) int64 {
	var total int64
	for i := range nodes {
		childTotal := s.decorateDeptTree(nodes[i].Children, directUserCountMap, allowedManageableSet)
		nodes[i].DirectUserCount = directUserCountMap[nodes[i].ID]
		nodes[i].HasChildren = len(nodes[i].Children) > 0
		nodes[i].Manageable = true
		if _, ok := allowedManageableSet[nodes[i].ID]; !ok {
			nodes[i].Manageable = false
		}
		nodes[i].Bindable = nodes[i].Manageable && !nodes[i].HasChildren
		nodes[i].TotalUserCount = nodes[i].DirectUserCount + childTotal
		total += nodes[i].TotalUserCount
	}
	return total
}

func (s *DeptService) getVisibleDeptUserCounts(scope *core.UserDataScope) (map[uint]int64, int64, error) {
	countMap := make(map[uint]int64)

	type deptUserCount struct {
		DeptID uint  `gorm:"column:dept_id"`
		Count  int64 `gorm:"column:count"`
	}

	db := global.DB.Model(&model.SysUser{})
	db = core.ApplyUserDataScope(db, scope, "sys_user")

	var rows []deptUserCount
	if err := db.Select("dept_id, COUNT(*) AS count").Group("dept_id").Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	var unassignedCount int64
	for _, row := range rows {
		if row.DeptID == 0 {
			unassignedCount += row.Count
			continue
		}
		countMap[row.DeptID] = row.Count
	}

	return countMap, unassignedCount, nil
}

func (s *DeptService) getDefaultUserAvatarURL() (string, error) {
	config, err := configsvc.Default.GetConfigByKey("register_logo")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(config.Value), nil
}

func expandDeptIDsWithAncestors(allDepts []model.SysDept, deptSet map[uint]struct{}) map[uint]struct{} {
	if len(deptSet) == 0 {
		return map[uint]struct{}{}
	}

	includedSet := make(map[uint]struct{}, len(deptSet))
	for deptID := range deptSet {
		includedSet[deptID] = struct{}{}
	}

	deptMap := make(map[uint]model.SysDept, len(allDepts))
	for _, dept := range allDepts {
		deptMap[dept.ID] = dept
	}

	for deptID := range deptSet {
		dept, ok := deptMap[deptID]
		if !ok {
			continue
		}
		for _, part := range strings.Split(dept.Ancestors, ",") {
			part = strings.TrimSpace(part)
			if part == "" || part == "0" {
				continue
			}
			ancestorID64, err := strconv.ParseUint(part, 10, 64)
			if err != nil {
				continue
			}
			includedSet[uint(ancestorID64)] = struct{}{}
		}
	}

	return includedSet
}

func filterDeptsBySet(allDepts []model.SysDept, includedSet map[uint]struct{}) []model.SysDept {
	if len(includedSet) == 0 {
		return nil
	}

	filtered := make([]model.SysDept, 0, len(includedSet))
	for _, dept := range allDepts {
		if _, ok := includedSet[dept.ID]; ok {
			filtered = append(filtered, dept)
		}
	}
	return filtered
}

func buildDeptAncestors(parentID uint) (string, error) {
	if parentID == 0 {
		return "0", nil
	}

	var parent model.SysDept
	if err := global.DB.First(&parent, parentID).Error; err != nil {
		return "", errors.New("上级部门不存在")
	}

	return buildDeptLineage(parent.Ancestors, parent.ID), nil
}

func buildDeptLineage(ancestors string, id uint) string {
	if ancestors == "" || ancestors == "0" {
		return fmt.Sprintf("0,%d", id)
	}
	return fmt.Sprintf("%s,%d", ancestors, id)
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

func GetManagedDeptIDs(operatorID uint) ([]uint, error) {
	scope, err := core.ResolveUserDataScopeForResource(operatorID, core.DataScopeResourceDeptManagement)
	if err != nil {
		return nil, err
	}

	if scope.All {
		var depts []model.SysDept
		if err := global.DB.Select("id").Find(&depts).Error; err != nil {
			return nil, err
		}
		ids := make([]uint, len(depts))
		for i, d := range depts {
			ids[i] = d.ID
		}
		return ids, nil
	}

	ids := make([]uint, len(scope.DeptIDs))
	copy(ids, scope.DeptIDs)
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, nil
}
