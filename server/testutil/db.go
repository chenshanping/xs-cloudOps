package testutil

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/config"
	"server/global"
	"server/model"
)

// SetupFileServiceTestDB creates an in-memory SQLite database with the models
// needed for file service integration tests, replaces global.DB, and restores
// it on test cleanup.
func SetupFileServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysConfig{}, &model.SysFile{}, &model.SysUser{}, &model.SysFileChunk{}, &model.AIMessage{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	previousDB := global.DB
	previousConfig := global.Config
	global.DB = db
	if global.Config == nil {
		global.Config = &config.Config{}
	}
	t.Cleanup(func() {
		global.DB = previousDB
		global.Config = previousConfig
	})

	return db
}

// SetupStorageServiceTestDB creates an in-memory SQLite database with the models
// needed for storage service tests, replaces global.DB, and restores it on test
// cleanup.
func SetupStorageServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysConfig{}, &model.SysFile{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	previousDB := global.DB
	previousConfig := global.Config
	global.DB = db
	if global.Config == nil {
		global.Config = &config.Config{}
	}
	t.Cleanup(func() {
		global.DB = previousDB
		global.Config = previousConfig
	})

	return db
}
