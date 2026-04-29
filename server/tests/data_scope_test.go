package tests

import (
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
	. "server/service"
	"server/utils"
)

func setupDataScopeTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(
		&model.SysFile{},
		&model.SysConfig{},
		&model.SysDept{},
		&model.SysRole{},
		&model.SysRoleDataScope{},
		&model.SysUser{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})

	return db
}

func createRoleFeatureScopeTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS sys_role_data_scope (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			role_id INTEGER NOT NULL,
			resource_code TEXT NOT NULL,
			data_scope INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sys_role_data_scope_dept (
			sys_role_data_scope_id INTEGER NOT NULL,
			sys_dept_id INTEGER NOT NULL
		)`,
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatalf("create role feature scope table failed: %v", err)
		}
	}
}

func ensureSysUserCreatedByColumn(t *testing.T, db *gorm.DB) {
	t.Helper()

	if db.Migrator().HasColumn("sys_user", "created_by") {
		return
	}

	if err := db.Exec("ALTER TABLE sys_user ADD COLUMN created_by INTEGER DEFAULT 0").Error; err != nil {
		t.Fatalf("add sys_user.created_by failed: %v", err)
	}
}

func insertRoleFeatureScope(
	t *testing.T,
	db *gorm.DB,
	roleID uint,
	resourceCode string,
	dataScope int,
	deptIDs ...uint,
) {
	t.Helper()

	result := db.Exec(
		"INSERT INTO sys_role_data_scope (role_id, resource_code, data_scope) VALUES (?, ?, ?)",
		roleID,
		resourceCode,
		dataScope,
	)
	if result.Error != nil {
		t.Fatalf("insert role feature scope failed: %v", result.Error)
	}

	if len(deptIDs) == 0 {
		return
	}

	var scopeID uint
	if err := db.Raw(
		"SELECT id FROM sys_role_data_scope WHERE role_id = ? AND resource_code = ? ORDER BY id DESC LIMIT 1",
		roleID,
		resourceCode,
	).Scan(&scopeID).Error; err != nil {
		t.Fatalf("load role feature scope id failed: %v", err)
	}

	for _, deptID := range deptIDs {
		if err := db.Exec(
			"INSERT INTO sys_role_data_scope_dept (sys_role_data_scope_id, sys_dept_id) VALUES (?, ?)",
			scopeID,
			deptID,
		).Error; err != nil {
			t.Fatalf("bind role feature scope dept failed: %v", err)
		}
	}
}

func TestUserServiceGetUserListRespectsDataScope(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create dept %s: %v", root.Name, err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create dept %s: %v", deptA.Name, err)
	}
	deptChild := model.SysDept{Name: "后端组", ParentID: deptA.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, deptA.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptChild).Error; err != nil {
		t.Fatalf("create dept %s: %v", deptChild.Name, err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create dept %s: %v", deptB.Name, err)
	}

	roleAll := model.SysRole{Name: "全部", Code: "all", DataScope: model.DataScopeAll, Status: 1}
	roleDept := model.SysRole{Name: "本部门", Code: "dept", DataScope: model.DataScopeDept, Status: 1}
	roleDeptTree := model.SysRole{Name: "本部门及下级", Code: "dept_tree", DataScope: model.DataScopeDeptAndChildren, Status: 1}
	roleSelf := model.SysRole{Name: "本人", Code: "self", DataScope: model.DataScopeSelf, Status: 1}
	roleCustom := model.SysRole{Name: "自定义", Code: "custom", DataScope: model.DataScopeCustom, Status: 1}

	for _, role := range []*model.SysRole{&roleAll, &roleDept, &roleDeptTree, &roleSelf, &roleCustom} {
		if err := db.Create(role).Error; err != nil {
			t.Fatalf("create role %s: %v", role.Code, err)
		}
	}
	if err := db.Model(&roleCustom).Association("Depts").Append(&deptB); err != nil {
		t.Fatalf("bind custom dept: %v", err)
	}

	admin := model.SysUser{Username: "admin-scope", Password: "pwd", Nickname: "管理员", Status: 1, DeptID: deptA.ID, Roles: []model.SysRole{roleAll}}
	deptUser := model.SysUser{Username: "dept-user", Password: "pwd", Nickname: "部门管理员", Status: 1, DeptID: deptA.ID, Roles: []model.SysRole{roleDept}}
	deptTreeUser := model.SysUser{Username: "tree-user", Password: "pwd", Nickname: "树管理员", Status: 1, DeptID: deptA.ID, Roles: []model.SysRole{roleDeptTree}}
	selfUser := model.SysUser{Username: "self-user", Password: "pwd", Nickname: "本人管理员", Status: 1, DeptID: deptA.ID, Roles: []model.SysRole{roleSelf}}
	customUser := model.SysUser{Username: "custom-user", Password: "pwd", Nickname: "自定义管理员", Status: 1, DeptID: deptA.ID, Roles: []model.SysRole{roleCustom}}
	childMember := model.SysUser{Username: "child-member", Password: "pwd", Nickname: "下级用户", Status: 1, DeptID: deptChild.ID}
	otherMember := model.SysUser{Username: "other-member", Password: "pwd", Nickname: "其他部门用户", Status: 1, DeptID: deptB.ID}

	for _, user := range []*model.SysUser{&admin, &deptUser, &deptTreeUser, &selfUser, &customUser, &childMember, &otherMember} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}
	if err := db.Model(&model.SysUser{}).Where("id = ?", childMember.ID).Update("created_by", selfUser.ID).Error; err != nil {
		t.Fatalf("set child member creator: %v", err)
	}

	tests := []struct {
		name       string
		operatorID uint
		wantTotal  int64
	}{
		{name: "all scope", operatorID: admin.ID, wantTotal: 7},
		{name: "dept scope", operatorID: deptUser.ID, wantTotal: 5},
		{name: "dept and children scope", operatorID: deptTreeUser.ID, wantTotal: 6},
		{name: "self scope", operatorID: selfUser.ID, wantTotal: 1},
		{name: "custom scope", operatorID: customUser.ID, wantTotal: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, total, err := User.GetUserList(tt.operatorID, &request.UserListRequest{
				PageRequest: request.PageRequest{Page: 1, PageSize: 20},
			})
			if err != nil {
				t.Fatalf("GetUserList error: %v", err)
			}
			if total != tt.wantTotal {
				t.Fatalf("GetUserList total = %d, want %d", total, tt.wantTotal)
			}
			if int64(len(list)) != tt.wantTotal {
				t.Fatalf("GetUserList len = %d, want %d", len(list), tt.wantTotal)
			}
		})
	}
}

func TestUserServiceGetUserListPrefersFeatureScopeOverRoleDefault(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleAll := model.SysRole{Name: "全部", Code: "all-user-feature", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "feature-operator",
		Password: "pwd",
		Nickname: "资源级操作人",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	deptMember := model.SysUser{Username: "feature-dept", Password: "pwd", Nickname: "同部门用户", Status: 1, DeptID: deptA.ID}
	otherMember := model.SysUser{Username: "feature-other", Password: "pwd", Nickname: "其他部门用户", Status: 1, DeptID: deptB.ID}
	for _, user := range []*model.SysUser{&operator, &deptMember, &otherMember} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}
	if err := db.Model(&model.SysUser{}).Where("id = ?", deptMember.ID).Update("created_by", operator.ID).Error; err != nil {
		t.Fatalf("set dept member creator: %v", err)
	}

	insertRoleFeatureScope(t, db, roleAll.ID, "system:user-management", model.DataScopeSelf)

	list, total, err := User.GetUserList(operator.ID, &request.UserListRequest{
		PageRequest: request.PageRequest{Page: 1, PageSize: 20},
	})
	if err != nil {
		t.Fatalf("GetUserList error: %v", err)
	}
	if total != 1 {
		t.Fatalf("feature scoped total = %d, want 1", total)
	}
	if len(list) != 1 || list[0].Username != deptMember.Username {
		t.Fatalf("unexpected feature scoped users: %+v", list)
	}
}

func TestUserServiceGetUserListUsesCreatorForSelfScopeInUserManagement(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)
	ensureSysUserCreatedByColumn(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}

	roleAll := model.SysRole{Name: "全部", Code: "all-user-created-by", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "creator-operator",
		Password: "pwd",
		Nickname: "创建人操作员",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	createdMember := model.SysUser{Username: "created-member", Password: "pwd", Nickname: "我创建的用户", Status: 1, DeptID: deptA.ID}
	otherMember := model.SysUser{Username: "other-member", Password: "pwd", Nickname: "他人创建的用户", Status: 1, DeptID: deptA.ID}
	for _, user := range []*model.SysUser{&operator, &createdMember, &otherMember} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	if err := db.Exec("UPDATE sys_user SET created_by = ? WHERE id = ?", operator.ID, createdMember.ID).Error; err != nil {
		t.Fatalf("set created member creator failed: %v", err)
	}
	if err := db.Exec("UPDATE sys_user SET created_by = ? WHERE id = ?", otherMember.ID, otherMember.ID).Error; err != nil {
		t.Fatalf("set other member creator failed: %v", err)
	}

	insertRoleFeatureScope(t, db, roleAll.ID, "system:user-management", model.DataScopeSelf)

	list, total, err := User.GetUserList(operator.ID, &request.UserListRequest{
		PageRequest: request.PageRequest{Page: 1, PageSize: 20},
	})
	if err != nil {
		t.Fatalf("GetUserList error: %v", err)
	}
	if total != 1 {
		t.Fatalf("creator scoped total = %d, want 1", total)
	}
	if len(list) != 1 || list[0].Username != createdMember.Username {
		t.Fatalf("unexpected creator scoped users: %+v", list)
	}
}

func TestUserServiceGetManagedUserInfoUsesCreatorForSelfScopeInUserManagement(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)
	ensureSysUserCreatedByColumn(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}

	roleAll := model.SysRole{Name: "全部", Code: "all-user-detail-created-by", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "detail-creator-operator",
		Password: "pwd",
		Nickname: "详情创建人操作员",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	createdMember := model.SysUser{Username: "detail-created-member", Password: "pwd", Nickname: "我创建的详情用户", Status: 1, DeptID: deptA.ID}
	otherMember := model.SysUser{Username: "detail-other-member", Password: "pwd", Nickname: "他人创建的详情用户", Status: 1, DeptID: deptA.ID}
	for _, user := range []*model.SysUser{&operator, &createdMember, &otherMember} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	if err := db.Exec("UPDATE sys_user SET created_by = ? WHERE id = ?", operator.ID, createdMember.ID).Error; err != nil {
		t.Fatalf("set created member creator failed: %v", err)
	}
	if err := db.Exec("UPDATE sys_user SET created_by = ? WHERE id = ?", otherMember.ID, otherMember.ID).Error; err != nil {
		t.Fatalf("set other member creator failed: %v", err)
	}

	insertRoleFeatureScope(t, db, roleAll.ID, "system:user-management", model.DataScopeSelf)

	if _, err := User.GetManagedUserInfo(operator.ID, createdMember.ID); err != nil {
		t.Fatalf("GetManagedUserInfo should allow created member: %v", err)
	}
	if _, err := User.GetManagedUserInfo(operator.ID, otherMember.ID); err == nil {
		t.Fatalf("GetManagedUserInfo should reject non-created member")
	}
}

func TestUserServiceCreateUserSetsCreatedBy(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}

	roleAll := model.SysRole{Name: "全部", Code: "all-user-created-by-writer", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "creator-writer",
		Password: "pwd",
		Nickname: "创建人写入操作员",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	req := &request.CreateUserRequest{
		Username: "created-by-target",
		Password: "123456",
		Nickname: "创建人写入目标",
		Gender:   1,
		Status:   1,
		DeptID:   deptA.ID,
	}
	if err := User.CreateUser(operator.ID, req); err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}

	var created model.SysUser
	if err := db.Where("username = ?", req.Username).First(&created).Error; err != nil {
		t.Fatalf("load created user: %v", err)
	}
	if created.CreatedBy != operator.ID {
		t.Fatalf("created user created_by = %d, want %d", created.CreatedBy, operator.ID)
	}
}

func TestUserServiceCreateUserUsesUserManagementScopeForDeptBinding(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "开发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptChild := model.SysDept{Name: "后端组", ParentID: deptA.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, deptA.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptChild).Error; err != nil {
		t.Fatalf("create deptChild: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleSelf := model.SysRole{Name: "仅本人默认", Code: "self-user-create-scope", DataScope: model.DataScopeSelf, Status: 1}
	if err := db.Create(&roleSelf).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "user-create-scope-operator",
		Password: "pwd",
		Nickname: "用户创建范围操作员",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleSelf},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	insertRoleFeatureScope(t, db, roleSelf.ID, "system:user-management", model.DataScopeDeptAndChildren)

	if err := User.CreateUser(operator.ID, &request.CreateUserRequest{
		Username: "scoped-created-user",
		Password: "123456",
		Nickname: "范围内用户",
		Gender:   1,
		Status:   1,
		DeptID:   deptChild.ID,
	}); err != nil {
		t.Fatalf("CreateUser should allow in-scope dept: %v", err)
	}

	if err := User.CreateUser(operator.ID, &request.CreateUserRequest{
		Username: "out-of-scope-created-user",
		Password: "123456",
		Nickname: "范围外用户",
		Gender:   1,
		Status:   1,
		DeptID:   deptB.ID,
	}); err == nil {
		t.Fatalf("CreateUser should reject out-of-scope dept")
	}
}

func TestUserServiceUpdateUserUsesUserManagementScopeForDeptBinding(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "开发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptChild := model.SysDept{Name: "后端组", ParentID: deptA.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, deptA.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptChild).Error; err != nil {
		t.Fatalf("create deptChild: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleSelf := model.SysRole{Name: "仅本人默认", Code: "self-user-update-scope", DataScope: model.DataScopeSelf, Status: 1}
	if err := db.Create(&roleSelf).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "user-update-scope-operator",
		Password: "pwd",
		Nickname: "用户更新范围操作员",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleSelf},
	}
	target := model.SysUser{
		Username: "user-update-scope-target",
		Password: "pwd",
		Nickname: "可管理目标用户",
		Status:   1,
		DeptID:   deptA.ID,
	}
	for _, user := range []*model.SysUser{&operator, &target} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	insertRoleFeatureScope(t, db, roleSelf.ID, "system:user-management", model.DataScopeDeptAndChildren)

	if err := User.UpdateUser(operator.ID, target.ID, &request.UpdateUserRequest{
		Nickname: "已更新到子部门",
		Email:    "scoped-update@example.com",
		Phone:    "13800001111",
		Status:   1,
		DeptID:   deptChild.ID,
		RoleIds:  nil,
	}); err != nil {
		t.Fatalf("UpdateUser should allow in-scope dept: %v", err)
	}

	if err := User.UpdateUser(operator.ID, target.ID, &request.UpdateUserRequest{
		Nickname: "尝试移动到范围外",
		Email:    "out-of-scope-update@example.com",
		Phone:    "13800002222",
		Status:   1,
		DeptID:   deptB.ID,
		RoleIds:  nil,
	}); err == nil {
		t.Fatalf("UpdateUser should reject out-of-scope dept")
	}
}

func TestUserServiceGetUserListMergesFeatureScopeAndRoleDefaultAcrossRoles(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleDept := model.SysRole{Name: "本部门", Code: "dept-default-role", DataScope: model.DataScopeDept, Status: 1}
	roleSelf := model.SysRole{Name: "本人", Code: "self-default-role", DataScope: model.DataScopeSelf, Status: 1}
	for _, role := range []*model.SysRole{&roleDept, &roleSelf} {
		if err := db.Create(role).Error; err != nil {
			t.Fatalf("create role %s: %v", role.Code, err)
		}
	}

	operator := model.SysUser{
		Username: "multi-operator",
		Password: "pwd",
		Nickname: "多角色操作人",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleDept, roleSelf},
	}
	deptMember := model.SysUser{Username: "multi-dept", Password: "pwd", Nickname: "同部门用户", Status: 1, DeptID: deptA.ID}
	otherMember := model.SysUser{Username: "multi-other", Password: "pwd", Nickname: "其他部门用户", Status: 1, DeptID: deptB.ID}
	for _, user := range []*model.SysUser{&operator, &deptMember, &otherMember} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	insertRoleFeatureScope(t, db, roleSelf.ID, "system:user-management", model.DataScopeCustom, deptB.ID)

	list, total, err := User.GetUserList(operator.ID, &request.UserListRequest{
		PageRequest: request.PageRequest{Page: 1, PageSize: 20},
	})
	if err != nil {
		t.Fatalf("GetUserList error: %v", err)
	}
	if total != 3 {
		t.Fatalf("merged feature scope total = %d, want 3", total)
	}
	if len(list) != 3 {
		t.Fatalf("merged feature scope len = %d, want 3", len(list))
	}
}

func TestDeptServiceGetManageableDeptTreePrefersFeatureScopeOverRoleDefault(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleAll := model.SysRole{Name: "全部", Code: "all-dept-feature", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "dept-feature-operator",
		Password: "pwd",
		Nickname: "部门资源级操作人",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	insertRoleFeatureScope(t, db, roleAll.ID, "system:dept-management", model.DataScopeCustom, deptB.ID)

	tree, _, err := Dept.GetManageableDeptTree(operator.ID)
	if err != nil {
		t.Fatalf("GetManageableDeptTree error: %v", err)
	}
	if len(tree) != 1 {
		t.Fatalf("feature dept tree root len = %d, want 1", len(tree))
	}
	if len(tree[0].Children) != 1 || tree[0].Children[0].ID != deptB.ID {
		t.Fatalf("unexpected feature scoped dept tree: %+v", tree)
	}
}

func TestDeptServiceGetManageableDeptTreeForResourceUsesUserManagementScope(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleAll := model.SysRole{Name: "全部", Code: "all-user-tree-feature", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "user-tree-feature-operator",
		Password: "pwd",
		Nickname: "用户树资源级操作人",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	deptAMember := model.SysUser{Username: "user-tree-dept-a", Password: "pwd", Nickname: "研发部用户", Status: 1, DeptID: deptA.ID}
	deptBMember := model.SysUser{Username: "user-tree-dept-b", Password: "pwd", Nickname: "市场部用户", Status: 1, DeptID: deptB.ID}
	for _, user := range []*model.SysUser{&operator, &deptAMember, &deptBMember} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	insertRoleFeatureScope(t, db, roleAll.ID, "system:user-management", model.DataScopeCustom, deptB.ID)

	tree, _, err := Dept.GetManageableDeptTreeForResource(operator.ID, "system:user-management")
	if err != nil {
		t.Fatalf("GetManageableDeptTreeForResource error: %v", err)
	}
	if len(tree) != 1 {
		t.Fatalf("feature dept tree root len = %d, want 1", len(tree))
	}
	if tree[0].TotalUserCount != 1 {
		t.Fatalf("root total count = %d, want 1", tree[0].TotalUserCount)
	}
	if len(tree[0].Children) != 1 || tree[0].Children[0].ID != deptB.ID {
		t.Fatalf("unexpected feature scoped dept tree: %+v", tree)
	}
	if tree[0].Children[0].DirectUserCount != 1 || tree[0].Children[0].TotalUserCount != 1 {
		t.Fatalf("unexpected deptB counts: %+v", tree[0].Children[0])
	}
}

func TestDeptServiceGetManageableDeptTreeForResourceUsesCreatedUsersForSelfScope(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)
	ensureSysUserCreatedByColumn(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleAll := model.SysRole{Name: "全部", Code: "all-user-tree-self-scope", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "user-tree-self-operator",
		Password: "pwd",
		Nickname: "用户树本人范围操作人",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	createdMember := model.SysUser{Username: "user-tree-created", Password: "pwd", Nickname: "我创建的市场部用户", Status: 1, DeptID: deptB.ID}
	otherMember := model.SysUser{Username: "user-tree-other", Password: "pwd", Nickname: "他人创建的市场部用户", Status: 1, DeptID: deptB.ID}
	for _, user := range []*model.SysUser{&operator, &createdMember, &otherMember} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}
	if err := db.Exec("UPDATE sys_user SET created_by = ? WHERE id = ?", operator.ID, createdMember.ID).Error; err != nil {
		t.Fatalf("set created member creator failed: %v", err)
	}
	if err := db.Exec("UPDATE sys_user SET created_by = ? WHERE id = ?", otherMember.ID, otherMember.ID).Error; err != nil {
		t.Fatalf("set other member creator failed: %v", err)
	}

	insertRoleFeatureScope(t, db, roleAll.ID, "system:user-management", model.DataScopeSelf)

	tree, _, err := Dept.GetManageableDeptTreeForResource(operator.ID, "system:user-management")
	if err != nil {
		t.Fatalf("GetManageableDeptTreeForResource error: %v", err)
	}
	if len(tree) != 1 {
		t.Fatalf("feature dept tree root len = %d, want 1", len(tree))
	}
	if tree[0].TotalUserCount != 1 {
		t.Fatalf("root total count = %d, want 1", tree[0].TotalUserCount)
	}
	if len(tree[0].Children) != 1 || tree[0].Children[0].ID != deptB.ID {
		t.Fatalf("unexpected self-scoped dept tree: %+v", tree)
	}
	if tree[0].Children[0].DirectUserCount != 1 || tree[0].Children[0].TotalUserCount != 1 {
		t.Fatalf("unexpected self-scoped dept counts: %+v", tree[0].Children[0])
	}
}

func TestDeptServiceGetManageableDeptTreeForResourceUsesDeptAndChildrenScopeWithoutCreatedUsers(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "开发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptChild := model.SysDept{Name: "后端组", ParentID: deptA.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, deptA.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptChild).Error; err != nil {
		t.Fatalf("create deptChild: %v", err)
	}

	roleDeptTree := model.SysRole{Name: "本部门及子级", Code: "dept-tree-user-tree-scope", DataScope: model.DataScopeDeptAndChildren, Status: 1}
	if err := db.Create(&roleDeptTree).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "dept-tree-scope-operator",
		Password: "pwd",
		Nickname: "开发部操作人",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleDeptTree},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	tree, _, err := Dept.GetManageableDeptTreeForResource(operator.ID, "system:user-management")
	if err != nil {
		t.Fatalf("GetManageableDeptTreeForResource error: %v", err)
	}
	if len(tree) != 1 {
		t.Fatalf("tree root len = %d, want 1", len(tree))
	}
	if len(tree[0].Children) != 1 || tree[0].Children[0].ID != deptA.ID {
		t.Fatalf("unexpected dept tree root children: %+v", tree)
	}
	if len(tree[0].Children[0].Children) != 1 || tree[0].Children[0].Children[0].ID != deptChild.ID {
		t.Fatalf("unexpected dept tree descendants: %+v", tree[0].Children[0])
	}
}

func TestRoleServiceGetRoleIncludesFeatureDataScopes(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}

	role := model.SysRole{Name: "角色A", Code: "role-feature-detail", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	insertRoleFeatureScope(t, db, role.ID, "system:user-management", model.DataScopeCustom, deptA.ID)

	detail, err := Role.GetRole(role.ID)
	if err != nil {
		t.Fatalf("GetRole error: %v", err)
	}
	if len(detail.FeatureDataScopes) != 1 {
		t.Fatalf("feature data scope len = %d, want 1", len(detail.FeatureDataScopes))
	}
	if detail.FeatureDataScopes[0].ResourceCode != "system:user-management" {
		t.Fatalf("feature resource code = %s, want %s", detail.FeatureDataScopes[0].ResourceCode, "system:user-management")
	}
	if len(detail.FeatureDataScopes[0].Depts) != 1 || detail.FeatureDataScopes[0].Depts[0].ID != deptA.ID {
		t.Fatalf("unexpected feature data scope depts: %+v", detail.FeatureDataScopes[0].Depts)
	}
}

func TestRoleServiceAssignDataScopesReplacesExistingOverrides(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	role := model.SysRole{Name: "角色B", Code: "role-feature-assign", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	if err := Role.AssignDataScopes(role.ID, []request.RoleFeatureDataScopeAssignment{
		{
			ResourceCode: "system:user-management",
			DataScope:    model.DataScopeCustom,
			DeptIds:      []uint{deptA.ID},
		},
		{
			ResourceCode: "system:dept-management",
			DataScope:    model.DataScopeDept,
		},
	}); err != nil {
		t.Fatalf("AssignDataScopes first error: %v", err)
	}

	if err := Role.AssignDataScopes(role.ID, []request.RoleFeatureDataScopeAssignment{
		{
			ResourceCode: "system:user-management",
			DataScope:    model.DataScopeCustom,
			DeptIds:      []uint{deptB.ID},
		},
	}); err != nil {
		t.Fatalf("AssignDataScopes second error: %v", err)
	}

	detail, err := Role.GetRole(role.ID)
	if err != nil {
		t.Fatalf("GetRole error: %v", err)
	}
	if len(detail.FeatureDataScopes) != 1 {
		t.Fatalf("feature data scope len = %d, want 1", len(detail.FeatureDataScopes))
	}
	scope := detail.FeatureDataScopes[0]
	if scope.ResourceCode != "system:user-management" {
		t.Fatalf("feature resource code = %s, want %s", scope.ResourceCode, "system:user-management")
	}
	if scope.DataScope != model.DataScopeCustom {
		t.Fatalf("feature data scope = %d, want %d", scope.DataScope, model.DataScopeCustom)
	}
	if len(scope.Depts) != 1 || scope.Depts[0].ID != deptB.ID {
		t.Fatalf("feature scope depts = %+v, want deptB only", scope.Depts)
	}
}

func TestUserServiceUpdateUserStatusRejectsOutOfScopeTarget(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create dept %s: %v", root.Name, err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create dept %s: %v", deptA.Name, err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create dept %s: %v", deptB.Name, err)
	}

	roleDept := model.SysRole{Name: "本部门", Code: "dept", DataScope: model.DataScopeDept, Status: 1}
	if err := db.Create(&roleDept).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{Username: "operator", Password: "pwd", Nickname: "操作人", Status: 1, DeptID: deptA.ID, Roles: []model.SysRole{roleDept}}
	target := model.SysUser{Username: "target", Password: "pwd", Nickname: "目标人", Status: 1, DeptID: deptB.ID}

	for _, user := range []*model.SysUser{&operator, &target} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	err := User.UpdateUserStatus(operator.ID, target.ID, 1)
	if err == nil {
		t.Fatalf("expected out-of-scope update to fail")
	}
}

func TestEnsureDeptManageableRejectsNonLeafDept(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	parent := model.SysDept{Name: "业务部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("create parent: %v", err)
	}
	child := model.SysDept{Name: "业务一组", ParentID: parent.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, parent.ID), Sort: 1, Status: 1}
	if err := db.Create(&child).Error; err != nil {
		t.Fatalf("create child: %v", err)
	}

	roleAll := model.SysRole{Name: "管理员", Code: "admin", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	operator := model.SysUser{Username: "operator", Password: "pwd", Nickname: "操作人", Status: 1, DeptID: root.ID, Roles: []model.SysRole{roleAll}}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	if err := EnsureDeptManageable(operator.ID, parent.ID); err == nil {
		t.Fatalf("expected parent dept to be rejected for user binding")
	}

	if err := EnsureDeptManageable(operator.ID, child.ID); err != nil {
		t.Fatalf("expected leaf dept to be allowed, got %v", err)
	}
}

func TestUserServiceGetUserListSupportsDeptTreeAndUnassignedFilter(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	parent := model.SysDept{Name: "业务部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("create parent: %v", err)
	}
	child := model.SysDept{Name: "业务一组", ParentID: parent.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, parent.ID), Sort: 1, Status: 1}
	if err := db.Create(&child).Error; err != nil {
		t.Fatalf("create child: %v", err)
	}
	other := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&other).Error; err != nil {
		t.Fatalf("create other: %v", err)
	}

	roleAll := model.SysRole{Name: "管理员", Code: "admin", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{Username: "admin-filter", Password: "pwd", Nickname: "管理员", Status: 1, DeptID: root.ID, Roles: []model.SysRole{roleAll}}
	parentUser := model.SysUser{Username: "parent-user", Password: "pwd", Nickname: "父部门用户", Status: 1, DeptID: parent.ID}
	childUser := model.SysUser{Username: "child-user", Password: "pwd", Nickname: "子部门用户", Status: 1, DeptID: child.ID}
	otherUser := model.SysUser{Username: "other-user", Password: "pwd", Nickname: "其他部门用户", Status: 1, DeptID: other.ID}
	unassignedUser := model.SysUser{Username: "unassigned-user", Password: "pwd", Nickname: "未绑定用户", Status: 1, DeptID: 0}

	for _, user := range []*model.SysUser{&operator, &parentUser, &childUser, &otherUser, &unassignedUser} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	parentID := int(parent.ID)
	list, total, err := User.GetUserList(operator.ID, &request.UserListRequest{
		PageRequest: request.PageRequest{Page: 1, PageSize: 20},
		DeptId:      &parentID,
	})
	if err != nil {
		t.Fatalf("GetUserList dept tree error: %v", err)
	}
	if total != 2 {
		t.Fatalf("dept tree total = %d, want 2", total)
	}
	if len(list) != 2 {
		t.Fatalf("dept tree len = %d, want 2", len(list))
	}

	list, total, err = User.GetUserList(operator.ID, &request.UserListRequest{
		PageRequest:    request.PageRequest{Page: 1, PageSize: 20},
		UnassignedDept: true,
	})
	if err != nil {
		t.Fatalf("GetUserList unassigned error: %v", err)
	}
	if total != 1 {
		t.Fatalf("unassigned total = %d, want 1", total)
	}
	if len(list) != 1 || list[0].Username != "unassigned-user" {
		t.Fatalf("unexpected unassigned users result: %+v", list)
	}
}

func TestUserServiceUpdateUserPersistsDeptID(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "111", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "222", ParentID: 0, Ancestors: "0", Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleAdmin := model.SysRole{Name: "管理员", Code: "admin", DataScope: model.DataScopeAll, Status: 1}
	roleUser := model.SysRole{Name: "普通用户", Code: "user", DataScope: model.DataScopeSelf, Status: 1}
	for _, role := range []*model.SysRole{&roleAdmin, &roleUser} {
		if err := db.Create(role).Error; err != nil {
			t.Fatalf("create role %s: %v", role.Code, err)
		}
	}

	operator := model.SysUser{Username: "admin", Password: "pwd", Nickname: "管理员", Status: 1, DeptID: root.ID, Roles: []model.SysRole{roleAdmin}}
	target := model.SysUser{Username: "target", Password: "pwd", Nickname: "目标用户", Status: 1, DeptID: root.ID, Roles: []model.SysRole{roleUser}}
	for _, user := range []*model.SysUser{&operator, &target} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	err := User.UpdateUser(operator.ID, target.ID, &request.UpdateUserRequest{
		Nickname: "已更新",
		Email:    "target@example.com",
		Phone:    "13800000000",
		Status:   1,
		DeptID:   deptB.ID,
		RoleIds:  []uint{roleUser.ID},
	})
	if err != nil {
		t.Fatalf("UpdateUser error: %v", err)
	}

	var updated model.SysUser
	if err := db.Preload("Dept").First(&updated, target.ID).Error; err != nil {
		t.Fatalf("reload target: %v", err)
	}
	if updated.DeptID != deptB.ID {
		t.Fatalf("updated dept_id = %d, want %d", updated.DeptID, deptB.ID)
	}
}

func TestUserServiceUpdateUserAllowsUnchangedLegacyParentDept(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	parent := model.SysDept{Name: "111", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("create parent: %v", err)
	}
	child := model.SysDept{Name: "333", ParentID: parent.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, parent.ID), Sort: 1, Status: 1}
	if err := db.Create(&child).Error; err != nil {
		t.Fatalf("create child: %v", err)
	}

	roleAdmin := model.SysRole{Name: "管理员", Code: "admin", DataScope: model.DataScopeAll, Status: 1}
	roleUser := model.SysRole{Name: "普通用户", Code: "user", DataScope: model.DataScopeSelf, Status: 1}
	for _, role := range []*model.SysRole{&roleAdmin, &roleUser} {
		if err := db.Create(role).Error; err != nil {
			t.Fatalf("create role %s: %v", role.Code, err)
		}
	}

	operator := model.SysUser{Username: "admin-legacy", Password: "pwd", Nickname: "管理员", Status: 1, DeptID: root.ID, Roles: []model.SysRole{roleAdmin}}
	target := model.SysUser{Username: "legacy-user", Password: "pwd", Nickname: "历史用户", Status: 1, DeptID: parent.ID, Roles: []model.SysRole{roleUser}}
	for _, user := range []*model.SysUser{&operator, &target} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	err := User.UpdateUser(operator.ID, target.ID, &request.UpdateUserRequest{
		Nickname: "已更新昵称",
		Email:    "legacy@example.com",
		Phone:    "13800000001",
		Status:   1,
		DeptID:   parent.ID,
		RoleIds:  []uint{roleUser.ID},
	})
	if err != nil {
		t.Fatalf("UpdateUser should allow unchanged legacy parent dept: %v", err)
	}

	var updated model.SysUser
	if err := db.First(&updated, target.ID).Error; err != nil {
		t.Fatalf("reload target: %v", err)
	}
	if updated.DeptID != parent.ID {
		t.Fatalf("updated dept_id = %d, want %d", updated.DeptID, parent.ID)
	}
	if updated.Nickname != "已更新昵称" {
		t.Fatalf("updated nickname = %s, want %s", updated.Nickname, "已更新昵称")
	}
}

func TestUserServiceResetManagedUserPasswordUsesConfiguredDefault(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}

	roleAdmin := model.SysRole{Name: "管理员", Code: "admin", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAdmin).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{Username: "admin-reset", Password: "old-hash", Nickname: "管理员", Status: 1, DeptID: root.ID, Roles: []model.SysRole{roleAdmin}}
	target := model.SysUser{Username: "target-reset", Password: "before-hash", Nickname: "目标用户", Status: 1, DeptID: root.ID}
	for _, user := range []*model.SysUser{&operator, &target} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	if err := db.Create(&model.SysConfig{
		Name:      "用户默认密码",
		Key:       "user_default_password",
		Value:     "ResetAbc1234",
		ValueType: "string",
		Remark:    "后台用户管理重置密码默认值",
	}).Error; err != nil {
		t.Fatalf("create config: %v", err)
	}

	if err := User.ResetManagedUserPassword(operator.ID, target.ID); err != nil {
		t.Fatalf("ResetManagedUserPassword error: %v", err)
	}

	var updated model.SysUser
	if err := db.First(&updated, target.ID).Error; err != nil {
		t.Fatalf("reload target: %v", err)
	}
	if !utils.CheckPassword("ResetAbc1234", updated.Password) {
		t.Fatalf("expected configured default password to be applied")
	}
}

func TestUserServiceResetManagedUserPasswordFallsBackToBuiltInDefault(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}

	roleAdmin := model.SysRole{Name: "管理员", Code: "admin", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAdmin).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{Username: "admin-fallback", Password: "old-hash", Nickname: "管理员", Status: 1, DeptID: root.ID, Roles: []model.SysRole{roleAdmin}}
	target := model.SysUser{Username: "target-fallback", Password: "before-hash", Nickname: "目标用户", Status: 1, DeptID: root.ID}
	for _, user := range []*model.SysUser{&operator, &target} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	if err := User.ResetManagedUserPassword(operator.ID, target.ID); err != nil {
		t.Fatalf("ResetManagedUserPassword error: %v", err)
	}

	var updated model.SysUser
	if err := db.First(&updated, target.ID).Error; err != nil {
		t.Fatalf("reload target: %v", err)
	}
	if !utils.CheckPassword("123456", updated.Password) {
		t.Fatalf("expected built-in fallback password to be applied")
	}
}

func TestUserServiceBatchResetManagedUserPasswords(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleAdmin := model.SysRole{Name: "管理员", Code: "admin", DataScope: model.DataScopeAll, Status: 1}
	roleDept := model.SysRole{Name: "本部门", Code: "dept", DataScope: model.DataScopeDept, Status: 1}
	for _, role := range []*model.SysRole{&roleAdmin, &roleDept} {
		if err := db.Create(role).Error; err != nil {
			t.Fatalf("create role %s: %v", role.Code, err)
		}
	}

	admin := model.SysUser{Username: "admin-batch-reset", Password: "pwd", Nickname: "管理员", Status: 1, DeptID: root.ID, Roles: []model.SysRole{roleAdmin}}
	deptOperator := model.SysUser{Username: "dept-operator-reset", Password: "pwd", Nickname: "部门管理员", Status: 1, DeptID: deptA.ID, Roles: []model.SysRole{roleDept}}
	inScopeA := model.SysUser{Username: "in-scope-a", Password: "old-a", Nickname: "范围内A", Status: 1, DeptID: deptA.ID}
	inScopeB := model.SysUser{Username: "in-scope-b", Password: "old-b", Nickname: "范围内B", Status: 1, DeptID: deptA.ID}
	outOfScope := model.SysUser{Username: "out-of-scope", Password: "old-c", Nickname: "范围外", Status: 1, DeptID: deptB.ID}
	for _, user := range []*model.SysUser{&admin, &deptOperator, &inScopeA, &inScopeB, &outOfScope} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	if err := db.Create(&model.SysConfig{
		Name:      "用户默认密码",
		Key:       "user_default_password",
		Value:     "BatchPwd9988",
		ValueType: "string",
		Remark:    "后台用户管理重置密码默认值",
	}).Error; err != nil {
		t.Fatalf("create config: %v", err)
	}

	if err := User.BatchResetManagedUserPasswords(deptOperator.ID, []uint{inScopeA.ID, outOfScope.ID}); err == nil {
		t.Fatalf("expected out-of-scope batch reset to fail")
	}

	if err := User.BatchResetManagedUserPasswords(admin.ID, []uint{inScopeA.ID, inScopeB.ID}); err != nil {
		t.Fatalf("BatchResetManagedUserPasswords error: %v", err)
	}

	var users []model.SysUser
	if err := db.Where("id IN ?", []uint{inScopeA.ID, inScopeB.ID}).Find(&users).Error; err != nil {
		t.Fatalf("reload users: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("reloaded users len = %d, want 2", len(users))
	}
	for _, user := range users {
		if !utils.CheckPassword("BatchPwd9988", user.Password) {
			t.Fatalf("expected configured batch password for user %s", user.Username)
		}
	}
}

func TestRoleServicePersistsCustomDeptDataScope(t *testing.T) {
	db := setupDataScopeTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	if err := Role.CreateRole(&request.CreateRoleRequest{
		Name:      "自定义范围角色",
		Code:      "custom-scope-role",
		Sort:      1,
		Status:    1,
		DataScope: model.DataScopeCustom,
		DeptIds:   []uint{deptA.ID, deptB.ID},
		Remark:    "测试自定义部门范围",
	}); err != nil {
		t.Fatalf("CreateRole error: %v", err)
	}

	var created model.SysRole
	if err := db.Preload("Depts").Where("code = ?", "custom-scope-role").First(&created).Error; err != nil {
		t.Fatalf("load created role: %v", err)
	}
	if created.DataScope != model.DataScopeCustom {
		t.Fatalf("created role data_scope = %d, want %d", created.DataScope, model.DataScopeCustom)
	}
	if len(created.Depts) != 2 {
		t.Fatalf("created role custom dept len = %d, want 2", len(created.Depts))
	}

	if err := Role.UpdateRole(created.ID, &request.UpdateRoleRequest{
		Name:      created.Name,
		Code:      created.Code,
		Sort:      created.Sort,
		Status:    created.Status,
		DataScope: model.DataScopeDeptAndChildren,
		DeptIds:   nil,
		Remark:    created.Remark,
	}); err != nil {
		t.Fatalf("UpdateRole error: %v", err)
	}

	var updated model.SysRole
	if err := db.Preload("Depts").First(&updated, created.ID).Error; err != nil {
		t.Fatalf("reload role: %v", err)
	}
	if updated.DataScope != model.DataScopeDeptAndChildren {
		t.Fatalf("updated role data_scope = %d, want %d", updated.DataScope, model.DataScopeDeptAndChildren)
	}
	if len(updated.Depts) != 0 {
		t.Fatalf("updated role custom dept len = %d, want 0", len(updated.Depts))
	}
}
