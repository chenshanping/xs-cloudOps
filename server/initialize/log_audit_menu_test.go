package initialize

import (
	"testing"

	"server/model"
)

func TestEnsureLogAuditMenusMigratesExistingMenusToMonitorWithoutOverwritingMetadata(t *testing.T) {
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

	auditMenu := model.SysMenu{
		ParentID:  systemMenu.ID,
		Name:      "操作审计",
		Path:      "/system/operation-audit",
		Component: "Layout",
		Icon:      "audit",
		Sort:      8,
		Type:      1,
		Status:    1,
		Hidden:    0,
	}
	if err := db.Create(&auditMenu).Error; err != nil {
		t.Fatalf("create operation audit menu: %v", err)
	}

	operationMenu := model.SysMenu{
		ParentID:   auditMenu.ID,
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
		ParentID:   auditMenu.ID,
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

	var updatedMonitorMenu model.SysMenu
	if err := db.First(&updatedMonitorMenu, monitorMenu.ID).Error; err != nil {
		t.Fatalf("reload monitor menu: %v", err)
	}
	if updatedMonitorMenu.Name != monitorMenu.Name {
		t.Fatalf("monitor menu name overwritten = %s, want %s", updatedMonitorMenu.Name, monitorMenu.Name)
	}

	var updatedAuditMenu model.SysMenu
	if err := db.First(&updatedAuditMenu, auditMenu.ID).Error; err != nil {
		t.Fatalf("reload operation audit menu: %v", err)
	}
	if updatedAuditMenu.Hidden != 1 {
		t.Fatalf("empty legacy audit menu hidden = %d, want 1", updatedAuditMenu.Hidden)
	}

	var updatedOperationMenu model.SysMenu
	if err := db.First(&updatedOperationMenu, operationMenu.ID).Error; err != nil {
		t.Fatalf("reload operation log menu: %v", err)
	}
	if updatedOperationMenu.ParentID != monitorMenu.ID {
		t.Fatalf("operation log parent_id = %d, want %d", updatedOperationMenu.ParentID, monitorMenu.ID)
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
	if updatedOperationMenu.Sort != operationMenu.Sort {
		t.Fatalf("operation log sort overwritten = %d, want %d", updatedOperationMenu.Sort, operationMenu.Sort)
	}

	var updatedLoginMenu model.SysMenu
	if err := db.First(&updatedLoginMenu, loginMenu.ID).Error; err != nil {
		t.Fatalf("reload login log menu: %v", err)
	}
	if updatedLoginMenu.ParentID != monitorMenu.ID {
		t.Fatalf("login log parent_id = %d, want %d", updatedLoginMenu.ParentID, monitorMenu.ID)
	}
	if updatedLoginMenu.Path != loginMenu.Path {
		t.Fatalf("login log path overwritten = %s, want %s", updatedLoginMenu.Path, loginMenu.Path)
	}
	if updatedLoginMenu.Component != loginMenu.Component {
		t.Fatalf("login log component overwritten = %s, want %s", updatedLoginMenu.Component, loginMenu.Component)
	}
}

func TestEnsureLogAuditMenusKeepsLegacyAuditMenuVisibleWhenItHasOtherChildren(t *testing.T) {
	db := setupInitializeTestDB(t)

	systemMenu := model.SysMenu{ParentID: 0, Name: "系统管理", Path: "/system", Component: "Layout", Type: 1, Status: 1}
	if err := db.Create(&systemMenu).Error; err != nil {
		t.Fatalf("create system menu: %v", err)
	}
	monitorMenu := model.SysMenu{ParentID: 0, Name: "运维监控", Path: "/monitor", Component: "Layout", Type: 1, Status: 1}
	if err := db.Create(&monitorMenu).Error; err != nil {
		t.Fatalf("create monitor menu: %v", err)
	}
	auditMenu := model.SysMenu{
		ParentID:  systemMenu.ID,
		Name:      "操作审计",
		Path:      "/system/operation-audit",
		Component: "Layout",
		Type:      1,
		Status:    1,
		Hidden:    0,
	}
	if err := db.Create(&auditMenu).Error; err != nil {
		t.Fatalf("create operation audit menu: %v", err)
	}
	customChild := model.SysMenu{
		ParentID:   auditMenu.ID,
		Name:       "自定义审计",
		Path:       "/system/operation-audit/custom",
		Component:  "system/operation-audit/custom",
		Type:       2,
		Permission: "system:operation-audit:custom",
		Status:     1,
	}
	if err := db.Create(&customChild).Error; err != nil {
		t.Fatalf("create custom audit child: %v", err)
	}

	ensureLogAuditMenus()

	var updatedAuditMenu model.SysMenu
	if err := db.First(&updatedAuditMenu, auditMenu.ID).Error; err != nil {
		t.Fatalf("reload operation audit menu: %v", err)
	}
	if updatedAuditMenu.Hidden != 0 {
		t.Fatalf("non-empty legacy audit menu hidden = %d, want 0", updatedAuditMenu.Hidden)
	}
}
