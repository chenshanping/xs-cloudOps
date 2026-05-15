package captcha

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
)

func setupCaptchaServiceTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysConfig{}); err != nil {
		t.Fatalf("auto migrate captcha config: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})
}

func TestGetCaptchaTypeFallsBackFromSliderToDigit(t *testing.T) {
	setupCaptchaServiceTestDB(t)

	if err := global.DB.Create(&model.SysConfig{
		Key:   "login_captcha_type",
		Value: "slider",
	}).Error; err != nil {
		t.Fatalf("create captcha config: %v", err)
	}

	if got := Default.GetCaptchaType(); got != "digit" {
		t.Fatalf("expected digit fallback, got %q", got)
	}
}
