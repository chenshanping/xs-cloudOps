package initialize

import (
	"testing"

	"github.com/glebarez/sqlite"
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

	if err := db.AutoMigrate(&model.SysMenu{}, &model.SysRole{}); err != nil {
		t.Fatalf("auto migrate base role/menu: %v", err)
	}
	if err := db.AutoMigrate(&model.SysApi{}, &model.SysConfig{}, &model.LegacyStorageRecord{}, &model.SysFile{}, &model.SysFileChunk{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	if err := db.Exec("ALTER TABLE sys_file ADD COLUMN storage_id integer").Error; err != nil {
		t.Fatalf("add sys_file.storage_id: %v", err)
	}
	if err := db.Exec("ALTER TABLE sys_file_chunk ADD COLUMN storage_id integer").Error; err != nil {
		t.Fatalf("add sys_file_chunk.storage_id: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
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
		Name:          "existing.txt",
		Path:          "existing.txt",
		URL:           "/api/v1/upload/existing.txt",
		MD5:           "existing-md5",
		StorageType:   string(model.StorageTypeMinio),
		Status:        1,
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
