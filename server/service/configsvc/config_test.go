package configsvc

import (
	"strconv"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
)

func setupConfigServiceTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysConfig{}, &model.SysFile{}, &model.SysFileReference{}); err != nil {
		t.Fatalf("auto migrate config service models: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})
}

func TestBatchUpdateConfigsSyncsImageFileRefsAndResolvesDerivedURL(t *testing.T) {
	setupConfigServiceTestDB(t)

	file := model.SysFile{
		Name:   "logo.png",
		Path:   "logo.png",
		URL:    "https://cdn.example.com/logo.png",
		Status: 1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	configs := []model.SysConfig{
		{Name: "系统Logo", Key: "sys_logo", Value: "https://stale.example.com/logo.png", ValueType: "string"},
		{Name: "系统Logo文件ID", Key: SysLogoFileIDConfigKey, Value: "", ValueType: "string"},
	}
	if err := global.DB.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}

	if err := Default.BatchUpdateConfigs(map[string]string{
		SysLogoFileIDConfigKey: strconv.FormatUint(uint64(file.ID), 10),
	}); err != nil {
		t.Fatalf("BatchUpdateConfigs bind file id: %v", err)
	}

	var ref model.SysFileReference
	if err := global.DB.Where("ref_table = ? AND ref_id = ? AND ref_field = ?", "sys_config", configs[1].ID, SysLogoFileIDConfigKey).
		First(&ref).Error; err != nil {
		t.Fatalf("query config file ref: %v", err)
	}
	if ref.FileID != file.ID {
		t.Fatalf("config file ref file_id = %d, want %d", ref.FileID, file.ID)
	}

	resolved, err := Default.GetConfigByKey("sys_logo")
	if err != nil {
		t.Fatalf("GetConfigByKey sys_logo: %v", err)
	}
	if resolved.Value != file.URL {
		t.Fatalf("resolved sys_logo = %q, want %q", resolved.Value, file.URL)
	}

	if err := Default.BatchUpdateConfigs(map[string]string{
		SysLogoFileIDConfigKey: "",
	}); err != nil {
		t.Fatalf("BatchUpdateConfigs clear file id: %v", err)
	}

	var refCount int64
	if err := global.DB.Model(&model.SysFileReference{}).
		Where("ref_table = ? AND ref_id = ?", "sys_config", configs[1].ID).
		Count(&refCount).Error; err != nil {
		t.Fatalf("count config file refs: %v", err)
	}
	if refCount != 0 {
		t.Fatalf("expected config file refs to be cleared, got %d", refCount)
	}

	resolved, err = Default.GetConfigByKey("sys_logo")
	if err != nil {
		t.Fatalf("GetConfigByKey sys_logo after clear: %v", err)
	}
	if resolved.Value != "" {
		t.Fatalf("resolved sys_logo after clear = %q, want empty", resolved.Value)
	}
}
