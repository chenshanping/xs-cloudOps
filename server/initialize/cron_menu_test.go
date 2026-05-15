package initialize

import (
	"testing"

	"server/model"
)

func TestEnsureCronTaskMenuApiCreatesMenusButtonsApisAndSeeds(t *testing.T) {
	db := setupInitializeTestDB(t)
	adminRole := model.SysRole{Name: "Admin", Code: "admin", Status: 1}
	if err := db.Create(&adminRole).Error; err != nil {
		t.Fatalf("create admin role: %v", err)
	}
	monitorMenu := model.SysMenu{ParentID: 0, Name: "自定义监控", Path: "/monitor", Component: "Layout", Icon: "custom-monitor", Sort: 99, Type: 1, Status: 1}
	if err := db.Create(&monitorMenu).Error; err != nil {
		t.Fatalf("create monitor menu: %v", err)
	}

	ensureCronTaskMenuApi()

	var updatedMonitor model.SysMenu
	if err := db.First(&updatedMonitor, monitorMenu.ID).Error; err != nil {
		t.Fatalf("reload monitor menu: %v", err)
	}
	if updatedMonitor.Name != monitorMenu.Name || updatedMonitor.Icon != monitorMenu.Icon || updatedMonitor.Sort != monitorMenu.Sort {
		t.Fatalf("monitor root overwritten: name=%s icon=%s sort=%d", updatedMonitor.Name, updatedMonitor.Icon, updatedMonitor.Sort)
	}

	expectedPermissions := []string{
		"monitor:cron:list",
		"monitor:cron:view",
		"monitor:cron:create",
		"monitor:cron:update",
		"monitor:cron:delete",
		"monitor:cron:enable",
		"monitor:cron:disable",
		"monitor:cron:runNow",
		"monitor:cron:logs:view",
	}
	var menuCount int64
	if err := db.Model(&model.SysMenu{}).Where("permission IN ?", expectedPermissions).Count(&menuCount).Error; err != nil {
		t.Fatalf("count cron menus: %v", err)
	}
	if menuCount != int64(len(expectedPermissions)) {
		t.Fatalf("cron menu count = %d, want %d", menuCount, len(expectedPermissions))
	}
	var taskMenu model.SysMenu
	if err := db.Where("permission = ?", "monitor:cron:list").First(&taskMenu).Error; err != nil {
		t.Fatalf("load cron task menu: %v", err)
	}
	var logButton model.SysMenu
	if err := db.Where("permission = ?", "monitor:cron:logs:view").First(&logButton).Error; err != nil {
		t.Fatalf("load cron log button: %v", err)
	}
	if logButton.ParentID != taskMenu.ID {
		t.Fatalf("cron log button parent = %d, want %d", logButton.ParentID, taskMenu.ID)
	}

	var apiCount int64
	if err := db.Model(&model.SysApi{}).Where("`group` = ?", "定时任务").Count(&apiCount).Error; err != nil {
		t.Fatalf("count cron apis: %v", err)
	}
	if apiCount != 10 {
		t.Fatalf("cron api count = %d, want 10", apiCount)
	}

	var bindingCount int64
	if err := db.Table("sys_menu_api").Count(&bindingCount).Error; err != nil {
		t.Fatalf("count menu api bindings: %v", err)
	}
	if bindingCount != 10 {
		t.Fatalf("cron menu api binding count = %d, want 10", bindingCount)
	}

	var adminRelationCount int64
	if err := db.Raw("SELECT COUNT(*) FROM sys_role_menu srm JOIN sys_menu sm ON sm.id = srm.sys_menu_id WHERE srm.sys_role_id = ? AND sm.permission IN ?", adminRole.ID, expectedPermissions).Scan(&adminRelationCount).Error; err != nil {
		t.Fatalf("count admin role menu relation: %v", err)
	}
	if adminRelationCount != int64(len(expectedPermissions)) {
		t.Fatalf("admin role menu relation count = %d, want %d", adminRelationCount, len(expectedPermissions))
	}

	var taskCount int64
	if err := db.Model(&model.SysCronTask{}).Where("code IN ? AND status = ?", []string{"cleanup_login_logs_default", "cleanup_operation_logs_default"}, model.CronTaskStatusDisabled).Count(&taskCount).Error; err != nil {
		t.Fatalf("count default cron tasks: %v", err)
	}
	if taskCount != 2 {
		t.Fatalf("default cron task count = %d, want 2", taskCount)
	}

	ensureCronTaskMenuApi()
	if err := db.Model(&model.SysMenu{}).Where("permission IN ?", expectedPermissions).Count(&menuCount).Error; err != nil {
		t.Fatalf("recount cron menus: %v", err)
	}
	if menuCount != int64(len(expectedPermissions)) {
		t.Fatalf("cron menu count after rerun = %d, want %d", menuCount, len(expectedPermissions))
	}
	if err := db.Model(&model.SysApi{}).Where("`group` = ?", "定时任务").Count(&apiCount).Error; err != nil {
		t.Fatalf("recount cron apis: %v", err)
	}
	if apiCount != 10 {
		t.Fatalf("cron api count after rerun = %d, want 10", apiCount)
	}
	if err := db.Table("sys_menu_api").Count(&bindingCount).Error; err != nil {
		t.Fatalf("recount menu api bindings: %v", err)
	}
	if bindingCount != 10 {
		t.Fatalf("cron menu api binding count after rerun = %d, want 10", bindingCount)
	}
}

func TestEnsureCronTaskMenuApiDoesNotOverwriteExistingPageMenu(t *testing.T) {
	db := setupInitializeTestDB(t)
	monitorMenu := model.SysMenu{ParentID: 0, Name: "运维监控", Path: "/monitor", Component: "Layout", Type: 1, Status: 1}
	if err := db.Create(&monitorMenu).Error; err != nil {
		t.Fatalf("create monitor menu: %v", err)
	}
	taskMenu := model.SysMenu{
		ParentID:   monitorMenu.ID,
		Name:       "自定义定时任务",
		Path:       "/custom/cron-task",
		Component:  "custom/cron-task/index",
		Icon:       "custom-cron-icon",
		Sort:       88,
		Type:       2,
		Permission: "monitor:cron:list",
		Status:     1,
		Hidden:     1,
	}
	if err := db.Create(&taskMenu).Error; err != nil {
		t.Fatalf("create cron task menu: %v", err)
	}

	ensureCronTaskMenuApi()

	var updated model.SysMenu
	if err := db.First(&updated, taskMenu.ID).Error; err != nil {
		t.Fatalf("reload cron task menu: %v", err)
	}
	if updated.Name != taskMenu.Name || updated.Path != taskMenu.Path || updated.Component != taskMenu.Component || updated.Icon != taskMenu.Icon || updated.Sort != taskMenu.Sort || updated.Hidden != taskMenu.Hidden {
		t.Fatalf("cron task page overwritten: got name=%s path=%s component=%s icon=%s sort=%d hidden=%d", updated.Name, updated.Path, updated.Component, updated.Icon, updated.Sort, updated.Hidden)
	}
}
