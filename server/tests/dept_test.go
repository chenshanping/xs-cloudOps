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
)

func setupDeptTestDB(t *testing.T) *gorm.DB {
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

func TestDeptServiceUpdateDeptRejectsCycle(t *testing.T) {
	db := setupDeptTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	parent := model.SysDept{Name: "研发部", ParentID: 0, Ancestors: "0", Sort: 2, Status: 1}
	child := model.SysDept{Name: "后端组", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}

	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	parent.ParentID = root.ID
	parent.Ancestors = fmt.Sprintf("0,%d", root.ID)
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("create parent: %v", err)
	}
	child.ParentID = parent.ID
	child.Ancestors = fmt.Sprintf("0,%d,%d", root.ID, parent.ID)
	if err := db.Create(&child).Error; err != nil {
		t.Fatalf("create child: %v", err)
	}

	err := Dept.UpdateDept(parent.ID, &request.UpdateDeptRequest{
		ParentID: child.ID,
		Name:     parent.Name,
		Sort:     parent.Sort,
		Status:   parent.Status,
		Remark:   parent.Remark,
	})
	if err == nil {
		t.Fatalf("expected cycle validation error")
	}
}

func TestDeptServiceDeleteDeptRejectsWhenChildrenOrUsersExist(t *testing.T) {
	db := setupDeptTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}

	parent := model.SysDept{Name: "业务部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("create parent: %v", err)
	}

	child := model.SysDept{Name: "销售组", ParentID: parent.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, parent.ID), Sort: 1, Status: 1}
	if err := db.Create(&child).Error; err != nil {
		t.Fatalf("create child: %v", err)
	}

	if err := Dept.DeleteDept(parent.ID); err == nil {
		t.Fatalf("expected delete to fail when child dept exists")
	}

	if err := db.Delete(&child).Error; err != nil {
		t.Fatalf("delete child: %v", err)
	}

	user := model.SysUser{
		Username: "dept-user",
		Password: "pwd",
		Nickname: "部门用户",
		Status:   1,
		DeptID:   parent.ID,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := Dept.DeleteDept(parent.ID); err == nil {
		t.Fatalf("expected delete to fail when assigned user exists")
	}
}

func TestDeptServiceGetManageableDeptTreeBuildsCountsAndBindability(t *testing.T) {
	db := setupDeptTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}

	parent := model.SysDept{Name: "业务部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("create parent: %v", err)
	}

	childA := model.SysDept{Name: "销售一组", ParentID: parent.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, parent.ID), Sort: 1, Status: 1}
	if err := db.Create(&childA).Error; err != nil {
		t.Fatalf("create childA: %v", err)
	}

	childB := model.SysDept{Name: "销售二组", ParentID: parent.ID, Ancestors: fmt.Sprintf("0,%d,%d", root.ID, parent.ID), Sort: 2, Status: 1}
	if err := db.Create(&childB).Error; err != nil {
		t.Fatalf("create childB: %v", err)
	}

	roleAll := model.SysRole{Name: "管理员", Code: "admin", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{Username: "admin-tree", Password: "pwd", Nickname: "管理员", Status: 1, DeptID: root.ID, Roles: []model.SysRole{roleAll}}
	parentUser := model.SysUser{Username: "parent-user", Password: "pwd", Nickname: "父部门用户", Status: 1, DeptID: parent.ID}
	childAUser1 := model.SysUser{Username: "child-a-1", Password: "pwd", Nickname: "子部门用户1", Status: 1, DeptID: childA.ID}
	childAUser2 := model.SysUser{Username: "child-a-2", Password: "pwd", Nickname: "子部门用户2", Status: 1, DeptID: childA.ID}
	childBUser := model.SysUser{Username: "child-b-1", Password: "pwd", Nickname: "子部门用户3", Status: 1, DeptID: childB.ID}
	unassignedUser := model.SysUser{Username: "unassigned", Password: "pwd", Nickname: "未绑定用户", Status: 1, DeptID: 0}

	for _, user := range []*model.SysUser{&operator, &parentUser, &childAUser1, &childAUser2, &childBUser, &unassignedUser} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create user %s: %v", user.Username, err)
		}
	}

	tree, unassignedCount, err := Dept.GetManageableDeptTree(operator.ID)
	if err != nil {
		t.Fatalf("GetManageableDeptTree error: %v", err)
	}

	if unassignedCount != 1 {
		t.Fatalf("unassigned count = %d, want 1", unassignedCount)
	}

	if len(tree) != 1 {
		t.Fatalf("tree root len = %d, want 1", len(tree))
	}

	if tree[0].DirectUserCount != 1 {
		t.Fatalf("root direct count = %d, want 1", tree[0].DirectUserCount)
	}
	if tree[0].TotalUserCount != 5 {
		t.Fatalf("root total count = %d, want 5", tree[0].TotalUserCount)
	}
	if tree[0].Bindable {
		t.Fatalf("root should not be bindable")
	}

	if len(tree[0].Children) != 1 {
		t.Fatalf("root children len = %d, want 1", len(tree[0].Children))
	}

	parentNode := tree[0].Children[0]
	if parentNode.DirectUserCount != 1 {
		t.Fatalf("parent direct count = %d, want 1", parentNode.DirectUserCount)
	}
	if parentNode.TotalUserCount != 4 {
		t.Fatalf("parent total count = %d, want 4", parentNode.TotalUserCount)
	}
	if parentNode.Bindable {
		t.Fatalf("parent should not be bindable")
	}

	if len(parentNode.Children) != 2 {
		t.Fatalf("parent children len = %d, want 2", len(parentNode.Children))
	}

	if !parentNode.Children[0].Bindable {
		t.Fatalf("leaf child should be bindable")
	}
	if parentNode.Children[0].DirectUserCount != 2 || parentNode.Children[0].TotalUserCount != 2 {
		t.Fatalf("childA counts = (%d, %d), want (2, 2)", parentNode.Children[0].DirectUserCount, parentNode.Children[0].TotalUserCount)
	}
}

func TestDeptServiceGetManageableDeptTreeReturnsEmptySliceWhenNoDeptVisible(t *testing.T) {
	db := setupDeptTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}

	role := model.SysRole{Name: "空自定义范围角色", Code: "dept-empty-custom", DataScope: model.DataScopeCustom, Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "dept-empty-user",
		Password: "pwd",
		Nickname: "空范围用户",
		Status:   1,
		DeptID:   root.ID,
		Roles:    []model.SysRole{role},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	tree, unassignedCount, err := Dept.GetManageableDeptTree(operator.ID)
	if err != nil {
		t.Fatalf("GetManageableDeptTree error: %v", err)
	}
	if tree == nil {
		t.Fatalf("expected empty dept tree slice, got nil")
	}
	if len(tree) != 0 {
		t.Fatalf("tree len = %d, want 0", len(tree))
	}
	if unassignedCount != 0 {
		t.Fatalf("unassigned count = %d, want 0", unassignedCount)
	}
}

func TestDeptServiceGetManageableDeptTreeWithDefaultsReturnsRegisterLogo(t *testing.T) {
	db := setupDeptTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}

	roleAll := model.SysRole{Name: "管理员", Code: "dept-default-avatar-role", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "dept-default-avatar-operator",
		Password: "pwd",
		Nickname: "部门树默认头像操作员",
		Status:   1,
		DeptID:   root.ID,
		Roles:    []model.SysRole{roleAll},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	if err := db.Create(&model.SysConfig{
		Name:      "注册默认头像",
		Key:       "register_logo",
		Value:     "https://cdn.example.com/register-default-avatar.png",
		ValueType: "string",
		Remark:    "用户管理新增默认头像",
	}).Error; err != nil {
		t.Fatalf("create config: %v", err)
	}

	tree, unassignedCount, defaultAvatarURL, err := Dept.GetManageableDeptTreeWithDefaultsForResource(operator.ID, "system:user-management")
	if err != nil {
		t.Fatalf("GetManageableDeptTreeWithDefaultsForResource error: %v", err)
	}
	if tree == nil {
		t.Fatalf("expected dept tree slice, got nil")
	}
	if unassignedCount != 0 {
		t.Fatalf("unassigned count = %d, want 0", unassignedCount)
	}
	if defaultAvatarURL != "https://cdn.example.com/register-default-avatar.png" {
		t.Fatalf("default avatar url = %q", defaultAvatarURL)
	}
}

func TestDeptServiceDeleteDeptRejectsWhenReferencedByRoleDataScope(t *testing.T) {
	db := setupDeptTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}

	target := model.SysDept{Name: "研发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&target).Error; err != nil {
		t.Fatalf("create target dept: %v", err)
	}

	role := model.SysRole{Name: "自定义部门角色", Code: "dept-custom-role", DataScope: model.DataScopeCustom, Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	if err := db.Model(&role).Association("Depts").Append(&target); err != nil {
		t.Fatalf("bind role dept: %v", err)
	}

	if err := Dept.DeleteDept(target.ID); err == nil {
		t.Fatalf("expected delete to fail when dept is referenced by role data scope")
	}
}

func TestDeptServiceDeleteDeptRejectsWhenReferencedByFeatureDataScope(t *testing.T) {
	db := setupDeptTestDB(t)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}

	target := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&target).Error; err != nil {
		t.Fatalf("create target dept: %v", err)
	}

	role := model.SysRole{Name: "功能范围角色", Code: "dept-feature-scope-role", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	scope := model.SysRoleDataScope{
		RoleID:       role.ID,
		ResourceCode: "system:user-management",
		DataScope:    model.DataScopeCustom,
	}
	if err := db.Create(&scope).Error; err != nil {
		t.Fatalf("create feature scope: %v", err)
	}
	if err := db.Model(&scope).Association("Depts").Append(&target); err != nil {
		t.Fatalf("bind feature scope dept: %v", err)
	}

	if err := Dept.DeleteDept(target.ID); err == nil {
		t.Fatalf("expected delete to fail when dept is referenced by feature data scope")
	}
}

func TestDeptServiceGetManagedDeptUsesDeptManagementScope(t *testing.T) {
	db := setupDeptTestDB(t)
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

	roleAll := model.SysRole{Name: "默认全部", Code: "dept-managed-detail-all", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "dept-managed-detail-operator",
		Password: "pwd",
		Nickname: "部门详情范围操作员",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	insertRoleFeatureScope(t, db, roleAll.ID, "system:dept-management", model.DataScopeDeptAndChildren)

	if _, err := Dept.GetManagedDept(operator.ID, deptChild.ID); err != nil {
		t.Fatalf("GetManagedDept should allow in-scope dept: %v", err)
	}
	if _, err := Dept.GetManagedDept(operator.ID, deptB.ID); err == nil {
		t.Fatalf("GetManagedDept should reject out-of-scope dept")
	}
}

func TestDeptServiceCreateDeptUsesDeptManagementScope(t *testing.T) {
	db := setupDeptTestDB(t)
	createRoleFeatureScopeTables(t, db)

	root := model.SysDept{Name: "平台", ParentID: 0, Ancestors: "0", Sort: 1, Status: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	deptA := model.SysDept{Name: "开发部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 1, Status: 1}
	if err := db.Create(&deptA).Error; err != nil {
		t.Fatalf("create deptA: %v", err)
	}
	deptB := model.SysDept{Name: "市场部", ParentID: root.ID, Ancestors: fmt.Sprintf("0,%d", root.ID), Sort: 2, Status: 1}
	if err := db.Create(&deptB).Error; err != nil {
		t.Fatalf("create deptB: %v", err)
	}

	roleAll := model.SysRole{Name: "默认全部", Code: "dept-create-scope-all", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "dept-create-scope-operator",
		Password: "pwd",
		Nickname: "部门创建范围操作员",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	insertRoleFeatureScope(t, db, roleAll.ID, "system:dept-management", model.DataScopeDeptAndChildren)

	if err := Dept.CreateManagedDept(operator.ID, &request.CreateDeptRequest{
		ParentID: deptA.ID,
		Name:     "测试子部门",
		Sort:     1,
		Status:   1,
		Remark:   "",
	}); err != nil {
		t.Fatalf("CreateManagedDept should allow in-scope parent: %v", err)
	}

	if err := Dept.CreateManagedDept(operator.ID, &request.CreateDeptRequest{
		ParentID: deptB.ID,
		Name:     "越权子部门",
		Sort:     1,
		Status:   1,
		Remark:   "",
	}); err == nil {
		t.Fatalf("CreateManagedDept should reject out-of-scope parent")
	}
}

func TestDeptServiceUpdateAndDeleteDeptUseDeptManagementScope(t *testing.T) {
	db := setupDeptTestDB(t)
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

	roleAll := model.SysRole{Name: "默认全部", Code: "dept-update-scope-all", DataScope: model.DataScopeAll, Status: 1}
	if err := db.Create(&roleAll).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	operator := model.SysUser{
		Username: "dept-update-scope-operator",
		Password: "pwd",
		Nickname: "部门编辑范围操作员",
		Status:   1,
		DeptID:   deptA.ID,
		Roles:    []model.SysRole{roleAll},
	}
	if err := db.Create(&operator).Error; err != nil {
		t.Fatalf("create operator: %v", err)
	}

	insertRoleFeatureScope(t, db, roleAll.ID, "system:dept-management", model.DataScopeDeptAndChildren)

	if err := Dept.UpdateManagedDept(operator.ID, deptChild.ID, &request.UpdateDeptRequest{
		ParentID: deptA.ID,
		Name:     "后端组-更新",
		Sort:     2,
		Status:   1,
		Remark:   "范围内更新",
	}); err != nil {
		t.Fatalf("UpdateManagedDept should allow in-scope dept: %v", err)
	}

	if err := Dept.UpdateManagedDept(operator.ID, deptB.ID, &request.UpdateDeptRequest{
		ParentID: root.ID,
		Name:     "市场部-更新",
		Sort:     3,
		Status:   1,
		Remark:   "范围外更新",
	}); err == nil {
		t.Fatalf("UpdateManagedDept should reject out-of-scope dept")
	}

	if err := Dept.DeleteManagedDept(operator.ID, deptB.ID); err == nil {
		t.Fatalf("DeleteManagedDept should reject out-of-scope dept")
	}
}
