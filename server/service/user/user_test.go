package user

import (
	"database/sql"
	"strconv"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
)

func setupUserServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysConfig{}, &model.SysFile{}, &model.SysRole{}, &model.SysUser{}); err != nil {
		t.Fatalf("auto migrate user service models: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})

	return db
}

func TestDeleteUserByIDClearsAvatarFileReferenceBeforeSoftDelete(t *testing.T) {
	setupUserServiceTestDB(t)

	file := model.SysFile{
		Name:   "avatar.png",
		Path:   "avatar.png",
		URL:    "/api/v1/upload/avatar.png",
		Status: 1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	user := model.SysUser{
		Username:     "avatar-user",
		Password:     "pwd",
		Nickname:     "Avatar User",
		Status:       1,
		AvatarFileID: file.ID,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := Default.deleteUserByID(user.ID); err != nil {
		t.Fatalf("deleteUserByID error: %v", err)
	}

	var deleted model.SysUser
	if err := global.DB.Unscoped().First(&deleted, user.ID).Error; err != nil {
		t.Fatalf("query deleted user: %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatalf("expected user to be soft deleted")
	}
	var avatarFileID sql.NullInt64
	if err := global.DB.Raw("SELECT avatar_file_id FROM sys_user WHERE id = ?", user.ID).Scan(&avatarFileID).Error; err != nil {
		t.Fatalf("query raw avatar_file_id: %v", err)
	}
	if avatarFileID.Valid {
		t.Fatalf("deleted user avatar_file_id should be NULL, got %d", avatarFileID.Int64)
	}
	if !strings.Contains(deleted.Username, "_deleted_") {
		t.Fatalf("deleted username = %q, want suffixed deleted username", deleted.Username)
	}
}

func TestResolveRegisterLogoAvatarFileIDPrefersBoundFileID(t *testing.T) {
	setupUserServiceTestDB(t)

	file := model.SysFile{
		Name:   "register-logo.png",
		Path:   "register-logo.png",
		URL:    "/api/v1/upload/register-logo.png",
		Status: 1,
	}
	if err := global.DB.Create(&file).Error; err != nil {
		t.Fatalf("create file: %v", err)
	}

	configs := []model.SysConfig{
		{Name: "注册默认头像", Key: "register_logo", Value: "/stale/register-logo.png", ValueType: "string"},
		{Name: "注册默认头像文件ID", Key: "register_logo_file_id", Value: strconv.FormatUint(uint64(file.ID), 10), ValueType: "string"},
	}
	if err := global.DB.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}

	fileID, err := resolveRegisterLogoAvatarFileID()
	if err != nil {
		t.Fatalf("resolveRegisterLogoAvatarFileID error: %v", err)
	}
	if fileID != file.ID {
		t.Fatalf("resolveRegisterLogoAvatarFileID = %d, want %d", fileID, file.ID)
	}
}

func TestRegisterWithoutDefaultAvatarStoresNullAvatarFileID(t *testing.T) {
	setupUserServiceTestDB(t)

	if err := Default.Register("new-user", "123456", "new-user@example.com"); err != nil {
		t.Fatalf("register error: %v", err)
	}

	var user model.SysUser
	if err := global.DB.Where("username = ?", "new-user").First(&user).Error; err != nil {
		t.Fatalf("query user: %v", err)
	}

	var avatarFileID sql.NullInt64
	if err := global.DB.Raw("SELECT avatar_file_id FROM sys_user WHERE id = ?", user.ID).Scan(&avatarFileID).Error; err != nil {
		t.Fatalf("query raw avatar_file_id: %v", err)
	}
	if avatarFileID.Valid {
		t.Fatalf("registered user avatar_file_id should be NULL, got %d", avatarFileID.Int64)
	}
}
