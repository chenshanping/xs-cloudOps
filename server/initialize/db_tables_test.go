package initialize

import (
	"testing"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/service"
)

func setupInitializeTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysMenu{}, &model.SysRole{}, &model.SysDept{}, &model.SysUser{}); err != nil {
		t.Fatalf("auto migrate base role/menu: %v", err)
	}
	if err := db.AutoMigrate(&model.SysApi{}, &model.SysConfig{}, &model.LegacyStorageRecord{}, &model.SysFile{}, &model.SysFileChunk{}, &model.SysDictType{}, &model.SysDictData{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	if err := db.Exec("ALTER TABLE sys_file ADD COLUMN storage_id integer").Error; err != nil {
		t.Fatalf("add sys_file.storage_id: %v", err)
	}
	if err := db.Exec("ALTER TABLE sys_file_chunk ADD COLUMN storage_id integer").Error; err != nil {
		t.Fatalf("add sys_file_chunk.storage_id: %v", err)
	}

	previousDB := global.DB
	previousLog := global.Log
	global.DB = db
	global.Log = zap.NewNop().Sugar()
	t.Cleanup(func() {
		global.DB = previousDB
		global.Log = previousLog
	})

	return db
}

func ensureInitializeRoleSuperAdminColumn(t *testing.T, db *gorm.DB) {
	t.Helper()

	if db.Migrator().HasColumn(&model.SysRole{}, "is_super_admin") {
		return
	}
	if err := db.Exec("ALTER TABLE sys_role ADD COLUMN is_super_admin numeric NOT NULL DEFAULT 0").Error; err != nil {
		t.Fatalf("add is_super_admin column: %v", err)
	}
}

func TestEnsureDeptMenusDoesNotOverwriteCustomizedIcon(t *testing.T) {
	db := setupInitializeTestDB(t)

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
		Icon:       "custom-icon",
		Sort:       3,
		Type:       2,
		Permission: "system:dept:list",
		Status:     1,
		Hidden:     0,
	}
	if err := db.Create(&deptMenu).Error; err != nil {
		t.Fatalf("create dept menu: %v", err)
	}

	ensureDeptMenus()

	var updated model.SysMenu
	if err := db.First(&updated, deptMenu.ID).Error; err != nil {
		t.Fatalf("reload dept menu: %v", err)
	}
	if updated.Icon != "custom-icon" {
		t.Fatalf("dept menu icon = %s, want %s", updated.Icon, "custom-icon")
	}
}

func TestEnsureBuiltInDataDoesNotOverwriteCustomizedRootDept(t *testing.T) {
	db := setupInitializeTestDB(t)

	rootDept := model.SysDept{
		ParentID:  0,
		Ancestors: "custom-ancestors",
		Name:      "平台",
		Sort:      88,
		Status:    0,
		Remark:    "保留已有根部门配置",
	}
	if err := db.Create(&rootDept).Error; err != nil {
		t.Fatalf("create root dept: %v", err)
	}

	ensureBuiltInData()

	var updated model.SysDept
	if err := db.First(&updated, rootDept.ID).Error; err != nil {
		t.Fatalf("reload root dept: %v", err)
	}
	if updated.Ancestors != rootDept.Ancestors {
		t.Fatalf("root dept ancestors overwritten = %s, want %s", updated.Ancestors, rootDept.Ancestors)
	}
	if updated.Sort != rootDept.Sort {
		t.Fatalf("root dept sort overwritten = %d, want %d", updated.Sort, rootDept.Sort)
	}
	if updated.Status != rootDept.Status {
		t.Fatalf("root dept status overwritten = %d, want %d", updated.Status, rootDept.Status)
	}
	if updated.Remark != rootDept.Remark {
		t.Fatalf("root dept remark overwritten = %s, want %s", updated.Remark, rootDept.Remark)
	}
}

func TestEnsureBuiltInDataCreatesAIMenus(t *testing.T) {
	db := setupInitializeTestDB(t)

	adminRole := model.SysRole{
		Name:   "超级管理员",
		Code:   "admin",
		Status: 1,
	}
	if err := db.Create(&adminRole).Error; err != nil {
		t.Fatalf("create admin role: %v", err)
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

	ensureBuiltInData()

	var aiToolsMenu model.SysMenu
	if err := db.Where("permission = ? AND type = ?", "ai:tools", 1).First(&aiToolsMenu).Error; err != nil {
		t.Fatalf("load ai tools menu: %v", err)
	}
	if aiToolsMenu.Path != "/ai-tools" {
		t.Fatalf("ai tools path = %s, want %s", aiToolsMenu.Path, "/ai-tools")
	}

	var aiChatMenu model.SysMenu
	if err := db.Where("permission = ? AND type = ?", "ai:chat:list", 2).First(&aiChatMenu).Error; err != nil {
		t.Fatalf("load ai chat menu: %v", err)
	}
	if aiChatMenu.ParentID != aiToolsMenu.ID {
		t.Fatalf("ai chat parent_id = %d, want %d", aiChatMenu.ParentID, aiToolsMenu.ID)
	}

	var aiConfigMenu model.SysMenu
	if err := db.Where("permission = ? AND type = ?", "ai:config:list", 2).First(&aiConfigMenu).Error; err != nil {
		t.Fatalf("load ai config menu: %v", err)
	}
	if aiConfigMenu.ParentID != aiToolsMenu.ID {
		t.Fatalf("ai config parent_id = %d, want %d", aiConfigMenu.ParentID, aiToolsMenu.ID)
	}
	if aiConfigMenu.Component != "ai/config/index" {
		t.Fatalf("ai config component = %s, want %s", aiConfigMenu.Component, "ai/config/index")
	}
}

func TestBackfillDepartmentFoundationAppliesCompatibleDefaultsOnlyToLegacyRows(t *testing.T) {
	db := setupInitializeTestDB(t)
	ensureInitializeRoleSuperAdminColumn(t, db)

	rootDept := model.SysDept{
		ParentID:  0,
		Ancestors: "0",
		Name:      "平台",
		Sort:      1,
		Status:    1,
		Remark:    "系统根部门",
	}
	if err := db.Create(&rootDept).Error; err != nil {
		t.Fatalf("create root dept: %v", err)
	}

	legacyUser := model.SysUser{Username: "legacy-user", Password: "pwd", Nickname: "旧用户", Status: 1, DeptID: 0}
	existingUser := model.SysUser{Username: "existing-user", Password: "pwd", Nickname: "已有部门用户", Status: 1, DeptID: rootDept.ID}
	if err := db.Create(&legacyUser).Error; err != nil {
		t.Fatalf("create legacy user: %v", err)
	}
	if err := db.Create(&existingUser).Error; err != nil {
		t.Fatalf("create existing user: %v", err)
	}

	legacyRole := model.SysRole{Name: "旧角色", Code: "legacy-role", DataScope: 0, Status: 1}
	existingRole := model.SysRole{Name: "已有角色", Code: "existing-role", DataScope: model.DataScopeSelf, Status: 1}
	if err := db.Create(&legacyRole).Error; err != nil {
		t.Fatalf("create legacy role: %v", err)
	}
	if err := db.Create(&existingRole).Error; err != nil {
		t.Fatalf("create existing role: %v", err)
	}

	backfillDepartmentFoundation(rootDept.ID)

	var updatedLegacyUser model.SysUser
	if err := db.First(&updatedLegacyUser, legacyUser.ID).Error; err != nil {
		t.Fatalf("reload legacy user: %v", err)
	}
	if updatedLegacyUser.DeptID != rootDept.ID {
		t.Fatalf("legacy user dept_id = %d, want %d", updatedLegacyUser.DeptID, rootDept.ID)
	}

	var updatedExistingUser model.SysUser
	if err := db.First(&updatedExistingUser, existingUser.ID).Error; err != nil {
		t.Fatalf("reload existing user: %v", err)
	}
	if updatedExistingUser.DeptID != existingUser.DeptID {
		t.Fatalf("existing user dept_id overwritten = %d, want %d", updatedExistingUser.DeptID, existingUser.DeptID)
	}

	var updatedLegacyRole model.SysRole
	if err := db.First(&updatedLegacyRole, legacyRole.ID).Error; err != nil {
		t.Fatalf("reload legacy role: %v", err)
	}
	if updatedLegacyRole.DataScope != model.DataScopeAll {
		t.Fatalf("legacy role data_scope = %d, want %d", updatedLegacyRole.DataScope, model.DataScopeAll)
	}

	var updatedExistingRole model.SysRole
	if err := db.First(&updatedExistingRole, existingRole.ID).Error; err != nil {
		t.Fatalf("reload existing role: %v", err)
	}
	if updatedExistingRole.DataScope != existingRole.DataScope {
		t.Fatalf("existing role data_scope overwritten = %d, want %d", updatedExistingRole.DataScope, existingRole.DataScope)
	}
}

func TestEnsureBuiltInDataMarksHistoricalAdminRolesAsExplicitSuperAdmin(t *testing.T) {
	db := setupInitializeTestDB(t)
	ensureInitializeRoleSuperAdminColumn(t, db)

	roles := []model.SysRole{
		{Name: "超级管理员", Code: "admin", Status: 1},
		{Name: "系统管理员", Code: "system_admin", Status: 1},
		{Name: "普通角色", Code: "auditor", Status: 1},
	}
	if err := db.Create(&roles).Error; err != nil {
		t.Fatalf("create roles: %v", err)
	}

	ensureBuiltInData()

	type roleFlag struct {
		ID           uint
		Code         string
		IsSuperAdmin int
	}

	var results []roleFlag
	if err := db.Raw("SELECT id, code, is_super_admin FROM sys_role ORDER BY id ASC").Scan(&results).Error; err != nil {
		t.Fatalf("load role flags: %v", err)
	}

	flags := make(map[string]int, len(results))
	for _, item := range results {
		flags[item.Code] = item.IsSuperAdmin
	}
	if flags["admin"] != 1 {
		t.Fatalf("admin is_super_admin = %d, want 1", flags["admin"])
	}
	if flags["system_admin"] != 1 {
		t.Fatalf("system_admin is_super_admin = %d, want 1", flags["system_admin"])
	}
	if flags["auditor"] != 0 {
		t.Fatalf("auditor is_super_admin = %d, want 0", flags["auditor"])
	}
}

func TestEnsureBuiltInDataGrantsAdminAIApiAccess(t *testing.T) {
	db := setupInitializeTestDB(t)

	adminRole := model.SysRole{
		Name:   "超级管理员",
		Code:   "admin",
		Status: 1,
	}
	if err := db.Create(&adminRole).Error; err != nil {
		t.Fatalf("create admin role: %v", err)
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

	ensureBuiltInData()

	var aiTestAPI model.SysApi
	if err := db.Where("path = ? AND method = ?", "/api/v1/ai/test", "POST").First(&aiTestAPI).Error; err != nil {
		t.Fatalf("load ai test api: %v", err)
	}

	var relationCount int64
	if err := db.Table("sys_role_api").
		Where("sys_role_id = ? AND sys_api_id = ?", adminRole.ID, aiTestAPI.ID).
		Count(&relationCount).Error; err != nil {
		t.Fatalf("count ai test api relation: %v", err)
	}
	if relationCount != 1 {
		t.Fatalf("admin ai test api relation count = %d, want 1", relationCount)
	}
}

func TestEnsureDeptApiAccessDoesNotOverwriteCustomizedApiMetadata(t *testing.T) {
	db := setupInitializeTestDB(t)

	sourceAPIs := []model.SysApi{
		{Path: "/api/v1/menus", Method: "GET", Group: "菜单管理", Description: "菜单列表", NeedAuth: true},
		{Path: "/api/v1/menus/:id", Method: "GET", Group: "菜单管理", Description: "菜单详情", NeedAuth: true},
		{Path: "/api/v1/menus", Method: "POST", Group: "菜单管理", Description: "创建菜单", NeedAuth: true},
		{Path: "/api/v1/menus/:id", Method: "PUT", Group: "菜单管理", Description: "更新菜单", NeedAuth: true},
		{Path: "/api/v1/menus/:id", Method: "DELETE", Group: "菜单管理", Description: "删除菜单", NeedAuth: true},
	}
	if err := db.Create(&sourceAPIs).Error; err != nil {
		t.Fatalf("create source apis: %v", err)
	}

	customTreeAPI := model.SysApi{
		Path:        "/api/v1/depts/tree",
		Method:      "GET",
		Group:       "自定义部门分组",
		Description: "保留已有部门树描述",
		NeedAuth:    false,
	}
	if err := db.Create(&customTreeAPI).Error; err != nil {
		t.Fatalf("create custom dept api: %v", err)
	}

	ensureDeptApiAccess()

	var updated model.SysApi
	if err := db.First(&updated, customTreeAPI.ID).Error; err != nil {
		t.Fatalf("reload custom dept api: %v", err)
	}
	if updated.Group != customTreeAPI.Group {
		t.Fatalf("dept api group overwritten = %s, want %s", updated.Group, customTreeAPI.Group)
	}
	if updated.Description != customTreeAPI.Description {
		t.Fatalf("dept api description overwritten = %s, want %s", updated.Description, customTreeAPI.Description)
	}
	if updated.NeedAuth != customTreeAPI.NeedAuth {
		t.Fatalf("dept api need_auth overwritten = %t, want %t", updated.NeedAuth, customTreeAPI.NeedAuth)
	}
}

func TestEnsureBuiltInDataGrantsProviderModelFetchApiToAdminRoles(t *testing.T) {
	db := setupInitializeTestDB(t)

	roles := []model.SysRole{
		{Name: "超级管理员", Code: "admin", Status: 1},
		{Name: "系统管理员", Code: "system_admin", Status: 1},
	}
	if err := db.Create(&roles).Error; err != nil {
		t.Fatalf("create roles: %v", err)
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

	ensureBuiltInData()

	var fetchAPI model.SysApi
	if err := db.Where("path = ? AND method = ?", "/api/v1/ai/providers/models/fetch", "POST").First(&fetchAPI).Error; err != nil {
		t.Fatalf("load provider model fetch api: %v", err)
	}

	for _, role := range roles {
		var relationCount int64
		if err := db.Table("sys_role_api").
			Where("sys_role_id = ? AND sys_api_id = ?", role.ID, fetchAPI.ID).
			Count(&relationCount).Error; err != nil {
			t.Fatalf("count role api relation (%s): %v", role.Code, err)
		}
		if relationCount != 1 {
			t.Fatalf("%s fetch api relation count = %d, want 1", role.Code, relationCount)
		}
	}
}

func TestEnsureBuiltInDataGrantsAIConfigReadWriteApisToAdminRoles(t *testing.T) {
	db := setupInitializeTestDB(t)

	roles := []model.SysRole{
		{Name: "超级管理员", Code: "admin", Status: 1},
		{Name: "系统管理员", Code: "system_admin", Status: 1},
	}
	if err := db.Create(&roles).Error; err != nil {
		t.Fatalf("create roles: %v", err)
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

	ensureBuiltInData()

	expected := []struct {
		path   string
		method string
	}{
		{path: "/api/v1/ai/config", method: "GET"},
		{path: "/api/v1/ai/config", method: "PUT"},
	}

	for _, item := range expected {
		var api model.SysApi
		if err := db.Where("path = ? AND method = ?", item.path, item.method).First(&api).Error; err != nil {
			t.Fatalf("load ai config api (%s %s): %v", item.method, item.path, err)
		}

		for _, role := range roles {
			var relationCount int64
			if err := db.Table("sys_role_api").
				Where("sys_role_id = ? AND sys_api_id = ?", role.ID, api.ID).
				Count(&relationCount).Error; err != nil {
				t.Fatalf("count role api relation (%s %s %s): %v", role.Code, item.method, item.path, err)
			}
			if relationCount != 1 {
				t.Fatalf("%s ai config api relation count (%s %s) = %d, want 1", role.Code, item.method, item.path, relationCount)
			}
		}
	}
}

func TestEnsureAIApiAccessDoesNotOverwriteCustomizedAIConfigApiMetadata(t *testing.T) {
	db := setupInitializeTestDB(t)

	customAPI := model.SysApi{
		Path:        "/api/v1/ai/config",
		Method:      "GET",
		Group:       "自定义AI分组",
		Description: "保留已有 AI 配置读取描述",
		NeedAuth:    false,
	}
	if err := db.Create(&customAPI).Error; err != nil {
		t.Fatalf("create custom ai config api: %v", err)
	}

	ensureAIApiAccess()

	var updated model.SysApi
	if err := db.First(&updated, customAPI.ID).Error; err != nil {
		t.Fatalf("reload ai config api: %v", err)
	}
	if updated.Group != customAPI.Group {
		t.Fatalf("ai config api group overwritten = %s, want %s", updated.Group, customAPI.Group)
	}
	if updated.Description != customAPI.Description {
		t.Fatalf("ai config api description overwritten = %s, want %s", updated.Description, customAPI.Description)
	}
	if updated.NeedAuth != customAPI.NeedAuth {
		t.Fatalf("ai config api need_auth overwritten = %t, want %t", updated.NeedAuth, customAPI.NeedAuth)
	}
}

func TestEnsureAIToolMenusDoesNotOverwriteCustomizedMetadata(t *testing.T) {
	db := setupInitializeTestDB(t)

	aiToolsMenu := model.SysMenu{
		ParentID:   0,
		Name:       "自定义AI工具",
		Path:       "/custom-ai-tools",
		Component:  "Layout",
		Icon:       "custom-root-icon",
		Sort:       99,
		Type:       1,
		Permission: "ai:tools",
		Status:     1,
		Hidden:     1,
	}
	if err := db.Create(&aiToolsMenu).Error; err != nil {
		t.Fatalf("create ai tools menu: %v", err)
	}

	aiConfigMenu := model.SysMenu{
		ParentID:   aiToolsMenu.ID,
		Name:       "自定义AI配置",
		Path:       "/custom-ai/config",
		Component:  "custom/ai/config/index",
		Icon:       "custom-child-icon",
		Sort:       88,
		Type:       2,
		Permission: "ai:config:list",
		Status:     1,
		Hidden:     1,
	}
	if err := db.Create(&aiConfigMenu).Error; err != nil {
		t.Fatalf("create ai config menu: %v", err)
	}

	ensureAIToolMenus()

	var updatedToolsMenu model.SysMenu
	if err := db.First(&updatedToolsMenu, aiToolsMenu.ID).Error; err != nil {
		t.Fatalf("reload ai tools menu: %v", err)
	}
	if updatedToolsMenu.Path != aiToolsMenu.Path {
		t.Fatalf("ai tools path overwritten = %s, want %s", updatedToolsMenu.Path, aiToolsMenu.Path)
	}
	if updatedToolsMenu.Icon != aiToolsMenu.Icon {
		t.Fatalf("ai tools icon overwritten = %s, want %s", updatedToolsMenu.Icon, aiToolsMenu.Icon)
	}

	var updatedConfigMenu model.SysMenu
	if err := db.First(&updatedConfigMenu, aiConfigMenu.ID).Error; err != nil {
		t.Fatalf("reload ai config menu: %v", err)
	}
	if updatedConfigMenu.Path != aiConfigMenu.Path {
		t.Fatalf("ai config path overwritten = %s, want %s", updatedConfigMenu.Path, aiConfigMenu.Path)
	}
	if updatedConfigMenu.Component != aiConfigMenu.Component {
		t.Fatalf("ai config component overwritten = %s, want %s", updatedConfigMenu.Component, aiConfigMenu.Component)
	}
}

func TestEnsureSystemStorageConfigsMigratesLegacyDefaultStorage(t *testing.T) {
	db := setupInitializeTestDB(t)

	legacy := model.LegacyStorageRecord{
		Name:      "旧默认存储",
		Type:      model.StorageTypeLocal,
		Config:    `{"base_path":"legacy-uploads","base_url":"/api/v1/upload"}`,
		IsDefault: 1,
		Status:    1,
	}
	if err := db.Create(&legacy).Error; err != nil {
		t.Fatalf("create legacy storage: %v", err)
	}

	ensureSystemStorageConfigs()

	var storageType model.SysConfig
	if err := db.Where("`key` = ?", service.StorageTypeConfigKey).First(&storageType).Error; err != nil {
		t.Fatalf("load storage type config: %v", err)
	}
	if storageType.Value != string(model.StorageTypeLocal) {
		t.Fatalf("storage_type = %s, want %s", storageType.Value, model.StorageTypeLocal)
	}

	var storageConfig model.SysConfig
	if err := db.Where("`key` = ?", service.StorageConfigKey(model.StorageTypeLocal)).First(&storageConfig).Error; err != nil {
		t.Fatalf("load storage local config: %v", err)
	}
	if storageConfig.Value != legacy.Config {
		t.Fatalf("storage_local_config = %s, want %s", storageConfig.Value, legacy.Config)
	}
}

func TestCleanupStorageBuiltInDataRemovesLegacyMenuAndApis(t *testing.T) {
	db := setupInitializeTestDB(t)

	role := model.SysRole{Name: "管理员", Code: "admin", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	menu := model.SysMenu{
		Name:       "存储管理",
		Path:       "/system/storage",
		Component:  "system/storage/index",
		Type:       2,
		Permission: "system:storage:list",
		Status:     1,
	}
	if err := db.Create(&menu).Error; err != nil {
		t.Fatalf("create legacy menu: %v", err)
	}
	if err := db.Exec("INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES (?, ?)", role.ID, menu.ID).Error; err != nil {
		t.Fatalf("create legacy role menu relation: %v", err)
	}

	api := model.SysApi{
		Path:        "/api/v1/storages",
		Method:      "GET",
		Group:       "存储管理",
		Description: "存储配置列表",
	}
	if err := db.Create(&api).Error; err != nil {
		t.Fatalf("create legacy api: %v", err)
	}
	if err := db.Exec("INSERT INTO sys_role_api (sys_role_id, sys_api_id) VALUES (?, ?)", role.ID, api.ID).Error; err != nil {
		t.Fatalf("create legacy role api relation: %v", err)
	}

	cleanupStorageBuiltInData()

	var menuCount int64
	if err := db.Model(&model.SysMenu{}).Where("path = ?", "/system/storage").Count(&menuCount).Error; err != nil {
		t.Fatalf("count legacy menus: %v", err)
	}
	if menuCount != 0 {
		t.Fatalf("legacy storage menu count = %d, want 0", menuCount)
	}

	var apiCount int64
	if err := db.Model(&model.SysApi{}).Where("path = ?", "/api/v1/storages").Count(&apiCount).Error; err != nil {
		t.Fatalf("count legacy apis: %v", err)
	}
	if apiCount != 0 {
		t.Fatalf("legacy storage api count = %d, want 0", apiCount)
	}
}

func TestCleanupSlowLogBuiltInDataRemovesLegacyMenuAndApis(t *testing.T) {
	db := setupInitializeTestDB(t)

	role := model.SysRole{Name: "管理员", Code: "admin", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	menu := model.SysMenu{
		ParentID:   2,
		Name:       "慢查询日志",
		Path:       "/monitor/show-log",
		Component:  "monitor/show-log/index",
		Type:       2,
		Permission: "monitor:show-log:list",
		Status:     1,
	}
	if err := db.Create(&menu).Error; err != nil {
		t.Fatalf("create slow log menu: %v", err)
	}
	if err := db.Exec("INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES (?, ?)", role.ID, menu.ID).Error; err != nil {
		t.Fatalf("create slow log role menu relation: %v", err)
	}

	apis := []model.SysApi{
		{Path: "/api/v1/logs/slow", Method: "GET", Group: "日志管理", Description: "慢查询日志列表"},
		{Path: "/api/v1/logs/slow/:id", Method: "DELETE", Group: "日志管理", Description: "删除慢查询日志"},
	}
	if err := db.Create(&apis).Error; err != nil {
		t.Fatalf("create slow log apis: %v", err)
	}
	for _, api := range apis {
		if err := db.Exec("INSERT INTO sys_role_api (sys_role_id, sys_api_id) VALUES (?, ?)", role.ID, api.ID).Error; err != nil {
			t.Fatalf("create slow log role api relation: %v", err)
		}
	}

	cleanupSlowLogBuiltInData()

	var menuCount int64
	if err := db.Model(&model.SysMenu{}).Where("path = ? OR permission = ?", "/monitor/show-log", "monitor:show-log:list").Count(&menuCount).Error; err != nil {
		t.Fatalf("count slow log menus: %v", err)
	}
	if menuCount != 0 {
		t.Fatalf("slow log menu count = %d, want 0", menuCount)
	}

	var apiCount int64
	if err := db.Model(&model.SysApi{}).Where("path IN ?", []string{"/api/v1/logs/slow", "/api/v1/logs/slow/:id"}).Count(&apiCount).Error; err != nil {
		t.Fatalf("count slow log apis: %v", err)
	}
	if apiCount != 0 {
		t.Fatalf("slow log api count = %d, want 0", apiCount)
	}

	var roleMenuCount int64
	if err := db.Table("sys_role_menu").Where("sys_role_id = ?", role.ID).Count(&roleMenuCount).Error; err != nil {
		t.Fatalf("count slow log role menu relations: %v", err)
	}
	if roleMenuCount != 0 {
		t.Fatalf("slow log role menu relation count = %d, want 0", roleMenuCount)
	}

	var roleAPICount int64
	if err := db.Table("sys_role_api").Where("sys_role_id = ?", role.ID).Count(&roleAPICount).Error; err != nil {
		t.Fatalf("count slow log role api relations: %v", err)
	}
	if roleAPICount != 0 {
		t.Fatalf("slow log role api relation count = %d, want 0", roleAPICount)
	}
}

func TestEnsureFileStorageSnapshotsBackfillsMissingRowsOnly(t *testing.T) {
	db := setupInitializeTestDB(t)

	legacy := model.LegacyStorageRecord{
		Name:      "旧存储",
		Type:      model.StorageTypeLocal,
		Config:    `{"base_path":"legacy-dir","base_url":"/api/v1/upload"}`,
		IsDefault: 1,
		Status:    1,
	}
	if err := db.Create(&legacy).Error; err != nil {
		t.Fatalf("create legacy storage: %v", err)
	}

	missingSnapshot := model.SysFile{
		Name:   "missing.txt",
		Path:   "missing.txt",
		URL:    "/api/v1/upload/missing.txt",
		MD5:    "missing-md5",
		Status: 1,
	}
	if err := db.Create(&missingSnapshot).Error; err != nil {
		t.Fatalf("create missing snapshot file: %v", err)
	}
	if err := db.Exec("UPDATE sys_file SET storage_id = ? WHERE id = ?", legacy.ID, missingSnapshot.ID).Error; err != nil {
		t.Fatalf("set missing snapshot storage_id: %v", err)
	}

	existingSnapshot := model.SysFile{
		Name:        "existing.txt",
		Path:        "existing.txt",
		URL:         "/api/v1/upload/existing.txt",
		MD5:         "existing-md5",
		StorageType: string(model.StorageTypeMinio),
		Status:      1,
	}
	if err := db.Create(&existingSnapshot).Error; err != nil {
		t.Fatalf("create existing snapshot file: %v", err)
	}
	if err := db.Exec("UPDATE sys_file SET storage_id = ? WHERE id = ?", legacy.ID, existingSnapshot.ID).Error; err != nil {
		t.Fatalf("set existing snapshot storage_id: %v", err)
	}

	ensureFileStorageSnapshots()

	var updatedMissing model.SysFile
	if err := db.First(&updatedMissing, missingSnapshot.ID).Error; err != nil {
		t.Fatalf("reload missing snapshot file: %v", err)
	}
	if updatedMissing.StorageType != string(model.StorageTypeLocal) {
		t.Fatalf("missing snapshot storage_type = %s, want %s", updatedMissing.StorageType, model.StorageTypeLocal)
	}

	var updatedExisting model.SysFile
	if err := db.First(&updatedExisting, existingSnapshot.ID).Error; err != nil {
		t.Fatalf("reload existing snapshot file: %v", err)
	}
	if updatedExisting.StorageType != existingSnapshot.StorageType {
		t.Fatalf("existing snapshot storage_type overwritten = %s, want %s", updatedExisting.StorageType, existingSnapshot.StorageType)
	}
}

func TestEnsureBuiltInDataCreatesGenderDictWithoutOverwritingCustomizedData(t *testing.T) {
	db := setupInitializeTestDB(t)

	customType := model.SysDictType{
		Name:   "自定义性别名称",
		Type:   "sys_gender",
		Status: 1,
		Remark: "保留已有字典类型配置",
	}
	if err := db.Create(&customType).Error; err != nil {
		t.Fatalf("create custom gender dict type: %v", err)
	}

	customData := model.SysDictData{
		DictType:  "sys_gender",
		Label:     "男士",
		Value:     "1",
		Sort:      88,
		Status:    1,
		TagType:   "cyan",
		IsDefault: 0,
		Remark:    "保留已有字典项配置",
	}
	if err := db.Create(&customData).Error; err != nil {
		t.Fatalf("create custom gender dict data: %v", err)
	}

	ensureBuiltInData()

	var dictType model.SysDictType
	if err := db.Where("type = ?", "sys_gender").First(&dictType).Error; err != nil {
		t.Fatalf("load gender dict type: %v", err)
	}
	if dictType.Name != customType.Name {
		t.Fatalf("gender dict type name overwritten = %s, want %s", dictType.Name, customType.Name)
	}

	var existingMale model.SysDictData
	if err := db.Where("dict_type = ? AND value = ?", "sys_gender", "1").First(&existingMale).Error; err != nil {
		t.Fatalf("load existing male dict data: %v", err)
	}
	if existingMale.Label != customData.Label {
		t.Fatalf("gender dict data label overwritten = %s, want %s", existingMale.Label, customData.Label)
	}
	if existingMale.TagType != customData.TagType {
		t.Fatalf("gender dict data tag_type overwritten = %s, want %s", existingMale.TagType, customData.TagType)
	}

	var count int64
	if err := db.Model(&model.SysDictData{}).Where("dict_type = ?", "sys_gender").Count(&count).Error; err != nil {
		t.Fatalf("count gender dict data: %v", err)
	}
	if count != 3 {
		t.Fatalf("gender dict data count = %d, want 3", count)
	}
}
