package service

import (
	"errors"
	"sort"

	"server/global"
	"server/model"
	"server/model/request"
)

type MenuService struct{}

var Menu = new(MenuService)

// 获取菜单列表(树形)
func (s *MenuService) GetMenuTree() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	if err := global.DB.Order("sort ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return s.buildMenuTree(menus, 0), nil
}

// 获取菜单详情
func (s *MenuService) GetMenu(id uint) (*model.SysMenu, error) {
	var menu model.SysMenu
	if err := global.DB.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// 创建菜单
func (s *MenuService) CreateMenu(req *request.CreateMenuRequest) error {
	menu := model.SysMenu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Sort:       req.Sort,
		Type:       req.Type,
		Permission: req.Permission,
		Status:     req.Status,
		Hidden:     req.Hidden,
	}
	err := global.DB.Create(&menu).Error
	if err == nil {
		Cache.ClearAllUserInfoCache() // 菜单变更清除所有缓存
	}
	return err
}

// 更新菜单
func (s *MenuService) UpdateMenu(id uint, req *request.UpdateMenuRequest) error {
	var menu model.SysMenu
	if err := global.DB.First(&menu, id).Error; err != nil {
		return errors.New("菜单不存在")
	}

	updates := map[string]interface{}{
		"parent_id":  req.ParentID,
		"name":       req.Name,
		"path":       req.Path,
		"component":  req.Component,
		"icon":       req.Icon,
		"sort":       req.Sort,
		"type":       req.Type,
		"permission": req.Permission,
		"status":     req.Status,
		"hidden":     req.Hidden,
	}

	err := global.DB.Model(&menu).Updates(updates).Error
	if err == nil {
		Cache.ClearAllUserInfoCache() // 菜单变更清除所有缓存
	}
	return err
}

// 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	// 检查是否有子菜单
	var count int64
	global.DB.Model(&model.SysMenu{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("存在子菜单，无法删除")
	}

	err := global.DB.Delete(&model.SysMenu{}, id).Error
	if err == nil {
		Cache.ClearAllUserInfoCache() // 菜单变更清除所有缓存
	}
	return err
}

// 构建菜单树
func (s *MenuService) buildMenuTree(menus []model.SysMenu, parentID uint) []model.SysMenu {
	var tree []model.SysMenu
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = s.buildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	// 按Sort排序
	sort.Slice(tree, func(i, j int) bool {
		return tree[i].Sort < tree[j].Sort
	})
	return tree
}

// 获取用户权限列表（包含按钮权限）
func (s *MenuService) GetUserPermissions(userID uint) ([]string, error) {
	var user model.SysUser
	if err := global.DB.Preload("Roles.Menus").First(&user, userID).Error; err != nil {
		return nil, err
	}

	permissionSet := make(map[string]bool)
	for _, role := range user.Roles {
		// 超级管理员拥有所有权限（ID=1 或 Code=admin）
		if role.ID == 1 || role.Code == "super_admin" {
			return []string{"*"}, nil
		}
		for _, menu := range role.Menus {
			if menu.Status == 1 && menu.Permission != "" {
				permissionSet[menu.Permission] = true
			}
		}
	}

	permissions := make([]string, 0, len(permissionSet))
	for perm := range permissionSet {
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

// 获取用户菜单
func (s *MenuService) GetUserMenus(userID uint) ([]model.SysMenu, error) {
	var user model.SysUser
	if err := global.DB.Preload("Roles.Menus").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 收集所有有权限的菜单ID（排除按钮）
	allowedMenuIDs := make(map[uint]bool)
	for _, role := range user.Roles {
		for _, menu := range role.Menus {
			if menu.Status == 1 && menu.Type != 3 { // 排除按钮
				allowedMenuIDs[menu.ID] = true
			}
		}
	}

	// 获取所有菜单
	var allMenus []model.SysMenu
	if err := global.DB.Order("sort ASC").Find(&allMenus).Error; err != nil {
		return nil, err
	}

	// 筛选有权限的菜单（包含父菜单）
	menuMap := make(map[uint]model.SysMenu)
	var toProcess []model.SysMenu

	for _, menu := range allMenus {
		if allowedMenuIDs[menu.ID] {
			toProcess = append(toProcess, menu)
		}
	}

	// 收集所有需要的菜单（包含父菜单）
	for _, menu := range toProcess {
		menuMap[menu.ID] = menu
		// 收集父菜单
		parentID := menu.ParentID
		for parentID != 0 {
			var parent model.SysMenu
			if err := global.DB.First(&parent, parentID).Error; err != nil {
				break
			}
			if _, exists := menuMap[parent.ID]; !exists {
				menuMap[parent.ID] = parent
			}
			parentID = parent.ParentID
		}
	}

	var menus []model.SysMenu
	for _, menu := range menuMap {
		if menu.Status == 1 {
			menus = append(menus, menu)
		}
	}

	return s.buildMenuTree(menus, 0), nil
}
