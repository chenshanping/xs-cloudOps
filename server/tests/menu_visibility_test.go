package tests

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
	. "server/service"
)

func setupMenuVisibilityTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(
		&model.SysConfig{},
		&model.SysMenu{},
		&model.SysRole{},
		&model.SysUser{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})

	return db
}

func ensureRoleSuperAdminColumn(t *testing.T, db *gorm.DB) {
	t.Helper()

	if db.Migrator().HasColumn(&model.SysRole{}, "is_super_admin") {
		return
	}
	if err := db.Exec("ALTER TABLE sys_role ADD COLUMN is_super_admin numeric NOT NULL DEFAULT 0").Error; err != nil {
		t.Fatalf("add is_super_admin column: %v", err)
	}
}

func TestMenuServiceGetUserMenusHidesDeptMenuWhenConfigDisabled(t *testing.T) {
	db := setupMenuVisibilityTestDB(t)

	if err := db.Create(&model.SysConfig{
		Name:      "部门模块显示",
		Key:       "dept_module_enabled",
		Value:     "false",
		ValueType: "string",
		Remark:    "隐藏部门管理模块",
	}).Error; err != nil {
		t.Fatalf("create dept config: %v", err)
	}

	systemMenu := model.SysMenu{
		ParentID:  0,
		Name:      "系统管理",
		Path:      "/system",
		Component: "Layout",
		Icon:      "setting",
		Sort:      1,
		Type:      1,
		Status:    1,
	}
	if err := db.Create(&systemMenu).Error; err != nil {
		t.Fatalf("create system menu: %v", err)
	}

	deptMenu := model.SysMenu{
		ParentID:   systemMenu.ID,
		Name:       "部门管理",
		Path:       "/system/dept",
		Component:  "system/dept/index",
		Icon:       "apartment",
		Sort:       1,
		Type:       2,
		Permission: "system:dept:list",
		Status:     1,
	}
	userMenu := model.SysMenu{
		ParentID:   systemMenu.ID,
		Name:       "用户管理",
		Path:       "/system/user",
		Component:  "system/user/index",
		Icon:       "user",
		Sort:       2,
		Type:       2,
		Permission: "system:user:list",
		Status:     1,
	}
	if err := db.Create(&deptMenu).Error; err != nil {
		t.Fatalf("create dept menu: %v", err)
	}
	if err := db.Create(&userMenu).Error; err != nil {
		t.Fatalf("create user menu: %v", err)
	}

	role := model.SysRole{
		Name:   "管理员",
		Code:   "admin-lite",
		Status: 1,
		Menus:  []model.SysMenu{deptMenu, userMenu},
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "menu-operator",
		Password: "pwd",
		Nickname: "菜单操作人",
		Status:   1,
		Roles:    []model.SysRole{role},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	menus, err := Menu.GetUserMenus(operator.ID)
	if err != nil {
		t.Fatalf("GetUserMenus error: %v", err)
	}

	if containsMenuPath(menus, "/system/dept") {
		t.Fatalf("department menu should be hidden when dept_module_enabled=false")
	}
	if !containsMenuPath(menus, "/system/user") {
		t.Fatalf("user menu should remain visible when department module is hidden")
	}
}

func TestMenuServiceGetUserMenusIncludesPageForButtonOnlyPermissions(t *testing.T) {
	db := setupMenuVisibilityTestDB(t)

	rootMenu := model.SysMenu{
		ParentID:  0,
		Name:      "系统管理",
		Path:      "/system",
		Component: "Layout",
		Sort:      1,
		Type:      1,
		Status:    1,
	}
	pageMenu := model.SysMenu{
		ParentID:   0,
		Name:       "用户管理",
		Path:       "/system/user",
		Component:  "system/user/index",
		Sort:       2,
		Type:       2,
		Permission: "system:user:list",
		Status:     1,
	}
	buttonMenu := model.SysMenu{
		ParentID:   pageMenu.ID,
		Name:       "重置密码",
		Path:       "",
		Component:  "",
		Sort:       1,
		Type:       3,
		Permission: "system:user:reset-password",
		Status:     1,
	}
	if err := db.Create(&rootMenu).Error; err != nil {
		t.Fatalf("create root menu: %v", err)
	}
	pageMenu.ParentID = rootMenu.ID
	if err := db.Create(&pageMenu).Error; err != nil {
		t.Fatalf("create page menu: %v", err)
	}
	buttonMenu.ParentID = pageMenu.ID
	if err := db.Create(&buttonMenu).Error; err != nil {
		t.Fatalf("create button menu: %v", err)
	}

	role := model.SysRole{
		Name:   "按钮权限角色",
		Code:   "button-only-role",
		Status: 1,
		Menus:  []model.SysMenu{buttonMenu},
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "button-only-user",
		Password: "pwd",
		Nickname: "按钮权限用户",
		Status:   1,
		Roles:    []model.SysRole{role},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	menus, err := Menu.GetUserMenus(operator.ID)
	if err != nil {
		t.Fatalf("GetUserMenus error: %v", err)
	}
	if len(menus) != 1 {
		t.Fatalf("visible menu roots len = %d, want 1", len(menus))
	}
	if menus[0].Path != "/system" {
		t.Fatalf("root menu path = %s, want /system", menus[0].Path)
	}
	if len(menus[0].Children) != 1 || menus[0].Children[0].Path != "/system/user" {
		t.Fatalf("expected user page menu to be preserved, got %#v", menus[0].Children)
	}
}

func TestMenuServiceGetUserMenusPrunesEmptyDirectoryMenus(t *testing.T) {
	db := setupMenuVisibilityTestDB(t)

	rootMenu := model.SysMenu{
		ParentID:  0,
		Name:      "系统管理",
		Path:      "/system",
		Component: "Layout",
		Sort:      1,
		Type:      1,
		Status:    1,
	}
	if err := db.Create(&rootMenu).Error; err != nil {
		t.Fatalf("create root menu: %v", err)
	}

	role := model.SysRole{
		Name:   "空目录角色",
		Code:   "empty-directory-role",
		Status: 1,
		Menus:  []model.SysMenu{rootMenu},
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "empty-directory-user",
		Password: "pwd",
		Nickname: "空目录用户",
		Status:   1,
		Roles:    []model.SysRole{role},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	menus, err := Menu.GetUserMenus(operator.ID)
	if err != nil {
		t.Fatalf("GetUserMenus error: %v", err)
	}
	if len(menus) != 0 {
		t.Fatalf("visible menu roots len = %d, want 0, menus=%#v", len(menus), menus)
	}
}

func TestMenuServiceGetUserPermissionsDoesNotBypassAdminRoleID(t *testing.T) {
	db := setupMenuVisibilityTestDB(t)
	ensureRoleSuperAdminColumn(t, db)

	buttonMenu := model.SysMenu{
		Name:       "新增",
		Path:       "",
		Component:  "",
		Sort:       1,
		Type:       3,
		Permission: "system:user:create",
		Status:     1,
	}
	if err := db.Create(&buttonMenu).Error; err != nil {
		t.Fatalf("create button menu: %v", err)
	}

	role := model.SysRole{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "管理员",
		Code:      "admin",
		Status:    1,
		Menus:     []model.SysMenu{buttonMenu},
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "admin-perm-user",
		Password: "pwd",
		Nickname: "管理员权限用户",
		Status:   1,
		Roles:    []model.SysRole{{BaseModel: model.BaseModel{ID: role.ID}}},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	perms, err := Menu.GetUserPermissions(operator.ID)
	if err != nil {
		t.Fatalf("GetUserPermissions error: %v", err)
	}
	if len(perms) != 1 || perms[0] != "system:user:create" {
		t.Fatalf("permissions = %#v, want only assigned button permission", perms)
	}
}

func TestMenuServiceGetUserPermissionsReturnsWildcardForExplicitSuperAdminRole(t *testing.T) {
	db := setupMenuVisibilityTestDB(t)
	ensureRoleSuperAdminColumn(t, db)

	buttonMenu := model.SysMenu{
		Name:       "新增",
		Sort:       1,
		Type:       3,
		Permission: "system:user:create",
		Status:     1,
	}
	if err := db.Create(&buttonMenu).Error; err != nil {
		t.Fatalf("create button menu: %v", err)
	}

	role := model.SysRole{
		Name:   "显式超管角色",
		Code:   "ops-admin",
		Status: 1,
		Menus:  []model.SysMenu{buttonMenu},
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	if err := db.Exec("UPDATE sys_role SET is_super_admin = 1 WHERE id = ?", role.ID).Error; err != nil {
		t.Fatalf("mark role as explicit super admin: %v", err)
	}

	operator := model.SysUser{
		Username: "explicit-super-admin-user",
		Password: "pwd",
		Nickname: "显式超管用户",
		Status:   1,
		Roles:    []model.SysRole{{BaseModel: model.BaseModel{ID: role.ID}}},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	perms, err := Menu.GetUserPermissions(operator.ID)
	if err != nil {
		t.Fatalf("GetUserPermissions error: %v", err)
	}
	if len(perms) != 1 || perms[0] != "*" {
		t.Fatalf("permissions = %#v, want wildcard for explicit super admin", perms)
	}
}

func TestMenuServiceGetUserMenusReturnsFullTreeForExplicitSuperAdminRole(t *testing.T) {
	db := setupMenuVisibilityTestDB(t)
	ensureRoleSuperAdminColumn(t, db)

	systemMenu := model.SysMenu{
		ParentID:  0,
		Name:      "系统管理",
		Path:      "/system",
		Component: "Layout",
		Sort:      1,
		Type:      1,
		Status:    1,
	}
	userMenu := model.SysMenu{
		ParentID:   0,
		Name:       "用户管理",
		Path:       "/system/user",
		Component:  "system/user/index",
		Sort:       1,
		Type:       2,
		Permission: "system:user:list",
		Status:     1,
	}
	roleMenu := model.SysMenu{
		ParentID:   0,
		Name:       "角色管理",
		Path:       "/system/role",
		Component:  "system/role/index",
		Sort:       2,
		Type:       2,
		Permission: "system:role:list",
		Status:     1,
	}
	if err := db.Create(&systemMenu).Error; err != nil {
		t.Fatalf("create system menu: %v", err)
	}
	userMenu.ParentID = systemMenu.ID
	roleMenu.ParentID = systemMenu.ID
	if err := db.Create(&userMenu).Error; err != nil {
		t.Fatalf("create user menu: %v", err)
	}
	if err := db.Create(&roleMenu).Error; err != nil {
		t.Fatalf("create role menu: %v", err)
	}

	role := model.SysRole{
		Name:   "菜单全开超管",
		Code:   "menu-super-admin",
		Status: 1,
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	if err := db.Exec("UPDATE sys_role SET is_super_admin = 1 WHERE id = ?", role.ID).Error; err != nil {
		t.Fatalf("mark role as explicit super admin: %v", err)
	}

	operator := model.SysUser{
		Username: "menu-super-admin-user",
		Password: "pwd",
		Nickname: "菜单超管用户",
		Status:   1,
		Roles:    []model.SysRole{{BaseModel: model.BaseModel{ID: role.ID}}},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	menus, err := Menu.GetUserMenus(operator.ID)
	if err != nil {
		t.Fatalf("GetUserMenus error: %v", err)
	}
	if len(menus) != 1 {
		t.Fatalf("visible menu roots len = %d, want 1", len(menus))
	}
	if menus[0].Path != "/system" {
		t.Fatalf("root menu path = %s, want /system", menus[0].Path)
	}
	if len(menus[0].Children) != 2 {
		t.Fatalf("visible child menu len = %d, want 2", len(menus[0].Children))
	}
}

func containsMenuPath(menus []model.SysMenu, target string) bool {
	for _, menu := range menus {
		if menu.Path == target {
			return true
		}
		if containsMenuPath(menu.Children, target) {
			return true
		}
	}
	return false
}
