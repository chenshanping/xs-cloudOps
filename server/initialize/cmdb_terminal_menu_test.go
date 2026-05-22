package initialize

import (
	"testing"

	"server/model"
)

func TestEnsureCmdbMenuApiCreatesTerminalMenusBindingsAndGrants(t *testing.T) {
	db := setupInitializeTestDB(t)

	adminRole := model.SysRole{Name: "超级管理员", Code: "admin", Status: 1}
	systemAdminRole := model.SysRole{Name: "系统管理员", Code: "system_admin", Status: 1}
	if err := db.Create(&[]model.SysRole{adminRole, systemAdminRole}).Error; err != nil {
		t.Fatalf("create roles: %v", err)
	}
	if err := db.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		t.Fatalf("reload admin role: %v", err)
	}
	if err := db.Where("code = ?", "system_admin").First(&systemAdminRole).Error; err != nil {
		t.Fatalf("reload system admin role: %v", err)
	}

	ensureCmdbMenuApi()

	expectedPermissions := []string{
		"cmdb:terminal:list",
		"cmdb:terminal:connect",
		"cmdb:terminal:view",
		"cmdb:terminal:disconnect",
		"cmdb:terminal:force_disconnect",
		"cmdb:terminal:audit",
	}
	for _, permission := range expectedPermissions {
		var count int64
		if err := db.Model(&model.SysMenu{}).Where("permission = ?", permission).Count(&count).Error; err != nil {
			t.Fatalf("count menu %s: %v", permission, err)
		}
		if count != 1 {
			t.Fatalf("menu %s count = %d, want 1", permission, count)
		}
	}

	var terminalMenu model.SysMenu
	if err := db.Where("permission = ?", "cmdb:terminal:list").First(&terminalMenu).Error; err != nil {
		t.Fatalf("load terminal menu: %v", err)
	}
	if terminalMenu.Path != "/cmdb/terminal/:hostId" {
		t.Fatalf("terminal menu path = %s, want /cmdb/terminal/:hostId", terminalMenu.Path)
	}
	if terminalMenu.Hidden != 1 {
		t.Fatalf("terminal menu hidden = %d, want 1", terminalMenu.Hidden)
	}
	if terminalMenu.IsStandalone != 1 {
		t.Fatalf("terminal menu is_standalone = %d, want 1", terminalMenu.IsStandalone)
	}

	var apiCount int64
	if err := db.Model(&model.SysApi{}).Where("path LIKE ?", "/api/v1/cmdb/terminal%").Count(&apiCount).Error; err != nil {
		t.Fatalf("count terminal apis: %v", err)
	}
	if apiCount != 7 {
		t.Fatalf("terminal api count = %d, want 7", apiCount)
	}

	for _, permission := range expectedPermissions {
		var count int64
		if err := db.Raw("SELECT COUNT(*) FROM sys_role_menu srm JOIN sys_menu sm ON sm.id = srm.sys_menu_id WHERE srm.sys_role_id = ? AND sm.permission = ?", adminRole.ID, permission).Scan(&count).Error; err != nil {
			t.Fatalf("count admin role menu %s: %v", permission, err)
		}
		if count != 1 {
			t.Fatalf("admin role menu %s count = %d, want 1", permission, count)
		}
		if err := db.Raw("SELECT COUNT(*) FROM sys_role_menu srm JOIN sys_menu sm ON sm.id = srm.sys_menu_id WHERE srm.sys_role_id = ? AND sm.permission = ?", systemAdminRole.ID, permission).Scan(&count).Error; err != nil {
			t.Fatalf("count system_admin role menu %s: %v", permission, err)
		}
		if count != 1 {
			t.Fatalf("system_admin role menu %s count = %d, want 1", permission, count)
		}
	}
}
