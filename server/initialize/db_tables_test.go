package initialize

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
)

func setupInitializeTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysMenu{}, &model.SysRole{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
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
