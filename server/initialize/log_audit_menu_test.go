package initialize

import (
	"testing"

	"server/model"
)

func TestEnsureLogAuditMenusMigratesExistingMenusWithoutOverwritingMetadata(t *testing.T) {
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

	monitorMenu := model.SysMenu{
		ParentID:  0,
		Name:      "系统监控",
		Path:      "/monitor",
		Component: "Layout",
		Icon:      "monitor",
		Sort:      2,
		Type:      1,
		Status:    1,
	}
	if err := db.Create(&monitorMenu).Error; err != nil {
		t.Fatalf("create monitor menu: %v", err)
	}

	operationMenu := model.SysMenu{
		ParentID:   monitorMenu.ID,
		Name:       "自定义操作日志",
		Path:       "/custom/operation-log",
		Component:  "custom/operation-log/index",
		Icon:       "custom-operation-icon",
		Sort:       91,
		Type:       2,
		Permission: "monitor:operation-log:list",
		Status:     1,
		Hidden:     1,
	}
	if err := db.Create(&operationMenu).Error; err != nil {
		t.Fatalf("create operation log menu: %v", err)
	}

	loginMenu := model.SysMenu{
		ParentID:   monitorMenu.ID,
		Name:       "自定义登录日志",
		Path:       "/custom/login-log",
		Component:  "custom/login-log/index",
		Icon:       "custom-login-icon",
		Sort:       92,
		Type:       2,
		Permission: "monitor:login-log:list",
		Status:     1,
		Hidden:     1,
	}
	if err := db.Create(&loginMenu).Error; err != nil {
		t.Fatalf("create login log menu: %v", err)
	}

	ensureLogAuditMenus()

	var auditMenu model.SysMenu
	if err := db.Where("path = ? AND type = ?", "/system/operation-audit", 1).First(&auditMenu).Error; err != nil {
		t.Fatalf("load operation audit menu: %v", err)
	}
	if auditMenu.ParentID != systemMenu.ID {
		t.Fatalf("audit menu parent_id = %d, want %d", auditMenu.ParentID, systemMenu.ID)
	}

	var updatedOperationMenu model.SysMenu
	if err := db.First(&updatedOperationMenu, operationMenu.ID).Error; err != nil {
		t.Fatalf("reload operation log menu: %v", err)
	}
	if updatedOperationMenu.ParentID != auditMenu.ID {
		t.Fatalf("operation log parent_id = %d, want %d", updatedOperationMenu.ParentID, auditMenu.ID)
	}
	if updatedOperationMenu.Path != operationMenu.Path {
		t.Fatalf("operation log path overwritten = %s, want %s", updatedOperationMenu.Path, operationMenu.Path)
	}
	if updatedOperationMenu.Icon != operationMenu.Icon {
		t.Fatalf("operation log icon overwritten = %s, want %s", updatedOperationMenu.Icon, operationMenu.Icon)
	}
	if updatedOperationMenu.Name != operationMenu.Name {
		t.Fatalf("operation log name overwritten = %s, want %s", updatedOperationMenu.Name, operationMenu.Name)
	}

	var updatedLoginMenu model.SysMenu
	if err := db.First(&updatedLoginMenu, loginMenu.ID).Error; err != nil {
		t.Fatalf("reload login log menu: %v", err)
	}
	if updatedLoginMenu.ParentID != auditMenu.ID {
		t.Fatalf("login log parent_id = %d, want %d", updatedLoginMenu.ParentID, auditMenu.ID)
	}
	if updatedLoginMenu.Path != loginMenu.Path {
		t.Fatalf("login log path overwritten = %s, want %s", updatedLoginMenu.Path, loginMenu.Path)
	}
	if updatedLoginMenu.Component != loginMenu.Component {
		t.Fatalf("login log component overwritten = %s, want %s", updatedLoginMenu.Component, loginMenu.Component)
	}
}
