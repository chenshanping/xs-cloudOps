package menu

import (
	"errors"
	"sort"
	"strings"

	"server/global"
	"server/model"
	"server/model/request"
	"server/service/configsvc"
	"server/service/core"
)

type MenuService struct{}

var Default = &MenuService{}

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
		core.Default.ClearAllUserInfoCache() // 菜单变更清除所有缓存
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
		core.Default.ClearAllUserInfoCache() // 菜单变更清除所有缓存
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
		core.Default.ClearAllUserInfoCache() // 菜单变更清除所有缓存
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

func (s *MenuService) pruneEmptyDirectories(menus []model.SysMenu) []model.SysMenu {
	pruned := make([]model.SysMenu, 0, len(menus))
	for _, menu := range menus {
		menu.Children = s.pruneEmptyDirectories(menu.Children)
		if menu.Type == 1 && len(menu.Children) == 0 {
			continue
		}
		pruned = append(pruned, menu)
	}
	return pruned
}

// 获取用户权限列表（包含按钮权限）
func (s *MenuService) GetUserPermissions(userID uint) ([]string, error) {
	var user model.SysUser
	if err := global.DB.Preload("Roles.Menus").First(&user, userID).Error; err != nil {
		return nil, err
	}

	permissionSet := make(map[string]bool)
	for _, role := range user.Roles {
		if role.IsSuperAdmin {
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

	// 获取所有菜单
	var allMenus []model.SysMenu
	if err := global.DB.Order("sort ASC").Find(&allMenus).Error; err != nil {
		return nil, err
	}
	deptModuleEnabled := s.isDeptModuleEnabled()
	allMenuMap := make(map[uint]model.SysMenu, len(allMenus))
	for _, menu := range allMenus {
		allMenuMap[menu.ID] = menu
	}

	if hasSuperAdminRole(user.Roles) {
		visibleMenus := make([]model.SysMenu, 0, len(allMenus))
		for _, menu := range allMenus {
			if !deptModuleEnabled && (menu.Path == "/system/dept" || menu.Permission == "system:dept:list") {
				continue
			}
			if menu.Status == 1 && menu.Type != 3 {
				visibleMenus = append(visibleMenus, menu)
			}
		}
		tree := s.buildMenuTree(visibleMenus, 0)
		tree = s.pruneEmptyDirectories(tree)
		if tree == nil {
			return []model.SysMenu{}, nil
		}
		return tree, nil
	}

	menuMap := make(map[uint]model.SysMenu)
	for _, role := range user.Roles {
		for _, roleMenu := range role.Menus {
			if roleMenu.Status != 1 {
				continue
			}

			currentID := roleMenu.ID
			for currentID != 0 {
				menu, exists := allMenuMap[currentID]
				if !exists {
					break
				}

				if !deptModuleEnabled && (menu.Path == "/system/dept" || menu.Permission == "system:dept:list") {
					break
				}
				if menu.Status == 1 && menu.Type != 3 {
					menuMap[menu.ID] = menu
				}
				currentID = menu.ParentID
			}
		}
	}

	var menus []model.SysMenu
	for _, menu := range menuMap {
		if !deptModuleEnabled && (menu.Path == "/system/dept" || menu.Permission == "system:dept:list") {
			continue
		}
		if menu.Status == 1 {
			menus = append(menus, menu)
		}
	}

	tree := s.buildMenuTree(menus, 0)
	tree = s.pruneEmptyDirectories(tree)
	if tree == nil {
		return []model.SysMenu{}, nil
	}
	return tree, nil
}

func hasSuperAdminRole(roles []model.SysRole) bool {
	for _, role := range roles {
		if role.IsSuperAdmin {
			return true
		}
	}
	return false
}

func (s *MenuService) isDeptModuleEnabled() bool {
	config, err := configsvc.Default.GetConfigByKey("dept_module_enabled")
	if err != nil {
		return true
	}
	return config.Value != "0" && strings.ToLower(config.Value) != "false"
}
