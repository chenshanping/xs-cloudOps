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
