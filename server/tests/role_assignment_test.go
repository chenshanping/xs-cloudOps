package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/middleware"
	"server/model"
	"server/model/response"
	. "server/service"
)

func setupRoleAssignmentTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysRole{}, &model.SysMenu{}, &model.SysApi{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	if err := db.AutoMigrate(&model.SysUser{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})

	return db
}

func ensureRoleAssignmentSuperAdminColumn(t *testing.T, db *gorm.DB) {
	t.Helper()

	if db.Migrator().HasColumn(&model.SysRole{}, "is_super_admin") {
		return
	}
	if err := db.Exec("ALTER TABLE sys_role ADD COLUMN is_super_admin numeric NOT NULL DEFAULT 0").Error; err != nil {
		t.Fatalf("add is_super_admin column: %v", err)
	}
}

func setupRoleAssignmentTestEnforcer(t *testing.T, db *gorm.DB) *casbin.Enforcer {
	t.Helper()

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		t.Fatalf("new casbin adapter: %v", err)
	}
	modelText := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
`
	m, err := casbinModel.NewModelFromString(modelText)
	if err != nil {
		t.Fatalf("new casbin model: %v", err)
	}
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		t.Fatalf("new casbin enforcer: %v", err)
	}
	enforcer.AddFunction("keyMatch2", util.KeyMatch2Func)

	previousEnforcer := global.Enforcer
	global.Enforcer = enforcer
	t.Cleanup(func() {
		global.Enforcer = previousEnforcer
	})

	return enforcer
}

func TestRoleServiceAssignMenusReplacesRoleMenus(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)

	role := model.SysRole{Name: "测试角色", Code: "test-role", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	menuA := model.SysMenu{Name: "菜单A", Path: "/a", Component: "a", Type: 2, Permission: "perm:a", Status: 1}
	menuB := model.SysMenu{Name: "菜单B", Path: "/b", Component: "b", Type: 2, Permission: "perm:b", Status: 1}
	if err := db.Create(&menuA).Error; err != nil {
		t.Fatalf("create menuA: %v", err)
	}
	if err := db.Create(&menuB).Error; err != nil {
		t.Fatalf("create menuB: %v", err)
	}

	if err := Role.AssignMenus(role.ID, []uint{menuA.ID, menuB.ID}); err != nil {
		t.Fatalf("AssignMenus error: %v", err)
	}

	var roleMenuCount int64
	if err := db.Table("sys_role_menu").Where("sys_role_id = ?", role.ID).Count(&roleMenuCount).Error; err != nil {
		t.Fatalf("count sys_role_menu: %v", err)
	}
	if roleMenuCount != 2 {
		t.Fatalf("sys_role_menu count = %d, want 2", roleMenuCount)
	}

	var updated model.SysRole
	if err := db.Preload("Menus").First(&updated, role.ID).Error; err != nil {
		t.Fatalf("reload role: %v", err)
	}
	if len(updated.Menus) != 2 {
		t.Fatalf("assigned menus len = %d, want 2", len(updated.Menus))
	}
}

func TestRoleServiceAssignApisReplacesRoleApis(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)
	enforcer := setupRoleAssignmentTestEnforcer(t, db)

	role := model.SysRole{Name: "测试角色", Code: "test-role-api", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	apiA := model.SysApi{Path: "/api/a", Method: "GET", Group: "测试", Description: "A", NeedAuth: true}
	apiB := model.SysApi{Path: "/api/b", Method: "POST", Group: "测试", Description: "B", NeedAuth: true}
	if err := db.Create(&apiA).Error; err != nil {
		t.Fatalf("create apiA: %v", err)
	}
	if err := db.Create(&apiB).Error; err != nil {
		t.Fatalf("create apiB: %v", err)
	}

	if err := Role.AssignApis(role.ID, []uint{apiA.ID, apiB.ID}); err != nil {
		t.Fatalf("AssignApis error: %v", err)
	}

	var roleApiCount int64
	if err := db.Table("sys_role_api").Where("sys_role_id = ?", role.ID).Count(&roleApiCount).Error; err != nil {
		t.Fatalf("count sys_role_api: %v", err)
	}
	if roleApiCount != 2 {
		t.Fatalf("sys_role_api count = %d, want 2", roleApiCount)
	}

	var updated model.SysRole
	if err := db.Preload("Apis").First(&updated, role.ID).Error; err != nil {
		t.Fatalf("reload role: %v", err)
	}
	if len(updated.Apis) != 2 {
		t.Fatalf("assigned apis len = %d, want 2", len(updated.Apis))
	}

	ok, err := enforcer.Enforce(role.Code, apiA.Path, apiA.Method)
	if err != nil {
		t.Fatalf("enforce apiA: %v", err)
	}
	if !ok {
		t.Fatalf("expected casbin policy for apiA")
	}

	ok, err = enforcer.Enforce(role.Code, apiB.Path, apiB.Method)
	if err != nil {
		t.Fatalf("enforce apiB: %v", err)
	}
	if !ok {
		t.Fatalf("expected casbin policy for apiB")
	}

	var casbinPolicyCount int64
	if err := db.Table("casbin_rule").Where("v0 = ?", role.Code).Count(&casbinPolicyCount).Error; err != nil {
		t.Fatalf("count casbin_rule: %v", err)
	}
	if casbinPolicyCount != 2 {
		t.Fatalf("casbin policy count = %d, want 2", casbinPolicyCount)
	}
}

func TestRoleServiceAssignApisDoesNotBypassAdminRole(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)
	ensureRoleAssignmentSuperAdminColumn(t, db)
	enforcer := setupRoleAssignmentTestEnforcer(t, db)

	role := model.SysRole{Name: "管理员", Code: "admin", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	allowedAPI := model.SysApi{Path: "/api/admin/allowed", Method: "GET", Group: "测试", Description: "允许", NeedAuth: true}
	deniedAPI := model.SysApi{Path: "/api/admin/denied", Method: "GET", Group: "测试", Description: "拒绝", NeedAuth: true}
	if err := db.Create(&allowedAPI).Error; err != nil {
		t.Fatalf("create allowed api: %v", err)
	}
	if err := db.Create(&deniedAPI).Error; err != nil {
		t.Fatalf("create denied api: %v", err)
	}

	if err := Role.AssignApis(role.ID, []uint{allowedAPI.ID}); err != nil {
		t.Fatalf("AssignApis error: %v", err)
	}

	ok, err := enforcer.Enforce(role.Code, allowedAPI.Path, allowedAPI.Method)
	if err != nil {
		t.Fatalf("enforce allowed api: %v", err)
	}
	if !ok {
		t.Fatalf("expected admin to access assigned api")
	}

	ok, err = enforcer.Enforce(role.Code, deniedAPI.Path, deniedAPI.Method)
	if err != nil {
		t.Fatalf("enforce denied api: %v", err)
	}
	if ok {
		t.Fatalf("expected admin to be denied for unassigned api")
	}
}

func TestCasbinAuthAllowsExplicitSuperAdminRoleWithoutAssignedAPI(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)
	ensureRoleAssignmentSuperAdminColumn(t, db)
	setupRoleAssignmentTestEnforcer(t, db)

	role := model.SysRole{Name: "显式超管", Code: "ops-admin", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	if err := db.Exec("UPDATE sys_role SET is_super_admin = 1 WHERE id = ?", role.ID).Error; err != nil {
		t.Fatalf("mark role as explicit super admin: %v", err)
	}

	user := model.SysUser{Username: "super-admin-user", Password: "pwd", Nickname: "超管用户", Status: 1}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	router := buildCasbinAuthTestEngine(user.ID, []uint{role.ID}, []string{role.Code})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/protected/resource", nil))

	resp := decodeTestResponse(t, recorder)
	if resp.Code != response.SUCCESS {
		t.Fatalf("response code = %d, want %d, body=%s", resp.Code, response.SUCCESS, recorder.Body.String())
	}
}

func TestCasbinAuthDisablingExplicitSuperAdminRestoresOrdinaryRoleChecks(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)
	ensureRoleAssignmentSuperAdminColumn(t, db)
	setupRoleAssignmentTestEnforcer(t, db)

	role := model.SysRole{Name: "可切换超管", Code: "toggle-admin", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	if err := db.Exec("UPDATE sys_role SET is_super_admin = 1 WHERE id = ?", role.ID).Error; err != nil {
		t.Fatalf("mark role as explicit super admin: %v", err)
	}

	user := model.SysUser{Username: "toggle-user", Password: "pwd", Nickname: "切换用户", Status: 1}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	router := buildCasbinAuthTestEngine(user.ID, []uint{role.ID}, []string{role.Code})

	firstRecorder := httptest.NewRecorder()
	router.ServeHTTP(firstRecorder, httptest.NewRequest(http.MethodGet, "/api/v1/protected/resource", nil))
	firstResp := decodeTestResponse(t, firstRecorder)
	if firstResp.Code != response.SUCCESS {
		t.Fatalf("enabled super admin response code = %d, want %d, body=%s", firstResp.Code, response.SUCCESS, firstRecorder.Body.String())
	}

	if err := db.Exec("UPDATE sys_role SET is_super_admin = 0 WHERE id = ?", role.ID).Error; err != nil {
		t.Fatalf("disable explicit super admin: %v", err)
	}

	secondRecorder := httptest.NewRecorder()
	router.ServeHTTP(secondRecorder, httptest.NewRequest(http.MethodGet, "/api/v1/protected/resource", nil))
	secondResp := decodeTestResponse(t, secondRecorder)
	if secondResp.Code != response.FORBIDDEN {
		t.Fatalf("disabled super admin response code = %d, want %d, body=%s", secondResp.Code, response.FORBIDDEN, secondRecorder.Body.String())
	}
}

func TestRoleServiceGetRoleLoadsPermissionAssociations(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)

	role := model.SysRole{Name: "查看权限角色", Code: "role-with-permissions", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	menu := model.SysMenu{Name: "菜单A", Path: "/perm-a", Component: "a", Type: 2, Permission: "perm:a", Status: 1}
	api := model.SysApi{Path: "/api/perm-a", Method: "GET", Group: "测试", Description: "权限接口", NeedAuth: true}
	if err := db.Create(&menu).Error; err != nil {
		t.Fatalf("create menu: %v", err)
	}
	if err := db.Create(&api).Error; err != nil {
		t.Fatalf("create api: %v", err)
	}

	if err := db.Exec("INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES (?, ?)", role.ID, menu.ID).Error; err != nil {
		t.Fatalf("insert sys_role_menu: %v", err)
	}
	if err := db.Exec("INSERT INTO sys_role_api (sys_role_id, sys_api_id) VALUES (?, ?)", role.ID, api.ID).Error; err != nil {
		t.Fatalf("insert sys_role_api: %v", err)
	}

	detail, err := Role.GetRole(role.ID)
	if err != nil {
		t.Fatalf("GetRole error: %v", err)
	}
	if len(detail.Menus) != 1 {
		t.Fatalf("detail menus len = %d, want 1", len(detail.Menus))
	}
	if len(detail.Apis) != 1 {
		t.Fatalf("detail apis len = %d, want 1", len(detail.Apis))
	}
}

func TestRoleServiceAssignMenusIncludesAncestorMenusForButtonPermission(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)

	rootMenu := model.SysMenu{Name: "系统管理", Path: "/system", Component: "Layout", Type: 1, Status: 1}
	if err := db.Create(&rootMenu).Error; err != nil {
		t.Fatalf("create root menu: %v", err)
	}

	pageMenu := model.SysMenu{
		ParentID:   rootMenu.ID,
		Name:       "用户管理",
		Path:       "/system/user",
		Component:  "system/user/index",
		Type:       2,
		Permission: "system:user:list",
		Status:     1,
	}
	if err := db.Create(&pageMenu).Error; err != nil {
		t.Fatalf("create page menu: %v", err)
	}

	buttonMenu := model.SysMenu{
		ParentID:   pageMenu.ID,
		Name:       "新增",
		Type:       3,
		Permission: "system:user:create",
		Status:     1,
	}
	if err := db.Create(&buttonMenu).Error; err != nil {
		t.Fatalf("create button menu: %v", err)
	}

	role := model.SysRole{Name: "按钮权限角色", Code: "role-button-only", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	if err := Role.AssignMenus(role.ID, []uint{buttonMenu.ID}); err != nil {
		t.Fatalf("AssignMenus error: %v", err)
	}

	var updated model.SysRole
	if err := db.Preload("Menus").First(&updated, role.ID).Error; err != nil {
		t.Fatalf("reload role: %v", err)
	}
	if len(updated.Menus) != 3 {
		t.Fatalf("assigned menus len = %d, want 3", len(updated.Menus))
	}

	user := model.SysUser{
		Username: "button-role-user",
		Password: "pwd",
		Nickname: "按钮角色用户",
		Status:   1,
		Roles:    []model.SysRole{{BaseModel: model.BaseModel{ID: role.ID}}},
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	menus, err := Menu.GetUserMenus(user.ID)
	if err != nil {
		t.Fatalf("GetUserMenus error: %v", err)
	}
	if len(menus) != 1 {
		t.Fatalf("visible menu roots len = %d, want 1", len(menus))
	}
	if menus[0].Path != "/system" {
		t.Fatalf("root menu path = %s, want /system", menus[0].Path)
	}
	if len(menus[0].Children) != 1 || menus[0].Children[0].Path != "/system/user" {
		t.Fatalf("expected user page menu to be preserved, got %#v", menus[0].Children)
	}
}

func TestRoleServiceGetRoleListIncludesUserStatsWithoutUsersSummary(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)

	role := model.SysRole{Name: "带用户角色", Code: "role-with-users", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	users := []model.SysUser{
		{Username: "alice", Password: "pwd", Nickname: "爱丽丝", Status: 1},
		{Username: "bob", Password: "pwd", Nickname: "鲍勃", Status: 1},
	}
	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			t.Fatalf("create user %d: %v", i, err)
		}
		if err := db.Exec("INSERT INTO sys_user_role (sys_user_id, sys_role_id) VALUES (?, ?)", users[i].ID, role.ID).Error; err != nil {
			t.Fatalf("bind role for user %d: %v", i, err)
		}
	}

	roles, err := Role.GetRoleList()
	if err != nil {
		t.Fatalf("GetRoleList error: %v", err)
	}

	var target map[string]any
	for _, item := range roles {
		raw, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshal role: %v", err)
		}

		var decoded map[string]any
		if err := json.Unmarshal(raw, &decoded); err != nil {
			t.Fatalf("unmarshal role: %v", err)
		}
		if int(decoded["id"].(float64)) == int(role.ID) {
			target = decoded
			break
		}
	}
	if target == nil {
		t.Fatalf("target role %d not found in role list", role.ID)
	}

	if got := int(target["user_count"].(float64)); got != 2 {
		t.Fatalf("user_count = %d, want 2", got)
	}

	if _, exists := target["users"]; exists {
		t.Fatalf("users field should be omitted from role list response: %#v", target["users"])
	}
}

func TestRoleServiceDeleteRoleBlocksLinkedUsers(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)

	role := model.SysRole{Name: "不可删除角色", Code: "role-delete-guard", Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	user := model.SysUser{Username: "guard-user", Password: "pwd", Nickname: "删除保护用户", Status: 1}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Exec("INSERT INTO sys_user_role (sys_user_id, sys_role_id) VALUES (?, ?)", user.ID, role.ID).Error; err != nil {
		t.Fatalf("bind role: %v", err)
	}

	err := Role.DeleteRole(role.ID)
	if err == nil {
		t.Fatalf("expected DeleteRole to fail when linked users exist")
	}

	var count int64
	if err := db.Model(&model.SysRole{}).Where("id = ?", role.ID).Count(&count).Error; err != nil {
		t.Fatalf("count role: %v", err)
	}
	if count != 1 {
		t.Fatalf("role count after blocked delete = %d, want 1", count)
	}
}

func buildCasbinAuthTestEngine(userID uint, roleIDs []uint, roleCodes []string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.ContextUserIDKey, userID)
		c.Set(middleware.ContextRoleIDsKey, roleIDs)
		c.Set(middleware.ContextRoleCodesKey, roleCodes)
		c.Next()
	})
	router.Use(middleware.CasbinAuth())
	router.GET("/api/v1/protected/resource", func(c *gin.Context) {
		response.Ok(c)
	})
	return router
}

func decodeTestResponse(t *testing.T, recorder *httptest.ResponseRecorder) response.Response {
	t.Helper()

	var resp response.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response body: %v, body=%s", err, recorder.Body.String())
	}
	return resp
}
