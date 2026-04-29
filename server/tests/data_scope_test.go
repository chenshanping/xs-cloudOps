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
