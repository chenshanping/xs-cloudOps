package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/middleware"
	"server/model"
	"server/model/response"
	. "server/service"
)

func TestRolePermissionSmokeSavePermissionsPreservesUnknownScope(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)
	setupRoleAssignmentTestEnforcer(t, db)

	operatorRole := model.SysRole{Name: "权限管理员", Code: "role-permission-smoke-operator", Status: 1}
	if err := db.Create(&operatorRole).Error; err != nil {
		t.Fatalf("create operator role: %v", err)
	}

	saveAPI := model.SysApi{
		Path:        "/api/v1/roles/:id/permissions",
		Method:      "PUT",
		Group:       "角色管理",
		Description: "统一保存角色权限",
		NeedAuth:    true,
	}
	if err := db.Create(&saveAPI).Error; err != nil {
		t.Fatalf("create save api: %v", err)
	}
	if err := Role.AssignApis(operatorRole.ID, []uint{saveAPI.ID}); err != nil {
		t.Fatalf("assign operator api: %v", err)
	}

	targetRole := model.SysRole{Name: "目标角色", Code: "role-permission-smoke-target", Status: 1}
	if err := db.Create(&targetRole).Error; err != nil {
		t.Fatalf("create target role: %v", err)
	}

	existingScopes := []model.SysRoleDataScope{
		{RoleID: targetRole.ID, ResourceCode: "system:user-management", DataScope: model.DataScopeDept},
		{RoleID: targetRole.ID, ResourceCode: "legacy:archived-resource", DataScope: model.DataScopeSelf},
	}
	for _, scope := range existingScopes {
		currentScope := scope
		if err := db.Create(&currentScope).Error; err != nil {
			t.Fatalf("create existing scope %s: %v", currentScope.ResourceCode, err)
		}
	}

	router := buildRolePermissionSmokeRouter(1001, []uint{operatorRole.ID}, []string{operatorRole.Code})

	body := `{
		"menu_ids": [],
		"direct_api_ids": [],
		"scopes": [
			{"resource_code":"system:user-management","data_scope":1,"dept_ids":[]},
			{"resource_code":"legacy:archived-resource","data_scope":5,"dept_ids":[]}
		]
	}`
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/roles/"+toUintString(targetRole.ID)+"/permissions",
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, req)

	resp := decodeTestResponse(t, recorder)
	if resp.Code != response.SUCCESS {
		t.Fatalf("save permissions smoke response code = %d, want %d, body=%s", resp.Code, response.SUCCESS, recorder.Body.String())
	}

	detail, err := Role.GetRole(targetRole.ID)
	if err != nil {
		t.Fatalf("GetRole error: %v", err)
	}
	if len(detail.FeatureDataScopes) != 2 {
		t.Fatalf("feature data scope len = %d, want 2", len(detail.FeatureDataScopes))
	}

	scopeMap := make(map[string]model.SysRoleDataScope, len(detail.FeatureDataScopes))
	for _, scope := range detail.FeatureDataScopes {
		scopeMap[scope.ResourceCode] = scope
	}

	if updated := scopeMap["system:user-management"]; updated.DataScope != model.DataScopeAll {
		t.Fatalf("known scope data_scope = %d, want %d", updated.DataScope, model.DataScopeAll)
	}

	preserved, ok := scopeMap["legacy:archived-resource"]
	if !ok {
		t.Fatalf("expected unknown scope to remain after smoke save")
	}
	if preserved.DataScope != model.DataScopeSelf {
		t.Fatalf("unknown scope data_scope = %d, want %d", preserved.DataScope, model.DataScopeSelf)
	}
}

func TestRolePermissionSmokeRejectsUnsupportedScopeWithoutMutation(t *testing.T) {
	db := setupRoleAssignmentTestDB(t)
	setupRoleAssignmentTestEnforcer(t, db)

	operatorRole := model.SysRole{Name: "权限管理员", Code: "role-permission-smoke-guard", Status: 1}
	if err := db.Create(&operatorRole).Error; err != nil {
		t.Fatalf("create operator role: %v", err)
	}

	saveAPI := model.SysApi{
		Path:        "/api/v1/roles/:id/permissions",
		Method:      "PUT",
		Group:       "角色管理",
		Description: "统一保存角色权限",
		NeedAuth:    true,
	}
	if err := db.Create(&saveAPI).Error; err != nil {
		t.Fatalf("create save api: %v", err)
	}
	if err := Role.AssignApis(operatorRole.ID, []uint{saveAPI.ID}); err != nil {
		t.Fatalf("assign operator api: %v", err)
	}

	rootMenu := model.SysMenu{Name: "系统管理", Path: "/system", Component: "Layout", Type: 1, Status: 1}
	if err := db.Create(&rootMenu).Error; err != nil {
		t.Fatalf("create root menu: %v", err)
	}

	pageMenu := model.SysMenu{
		ParentID:   rootMenu.ID,
		Name:       "菜单A",
		Path:       "/system/a",
		Component:  "system/a/index",
		Type:       2,
		Permission: "system:a:list",
		Status:     1,
	}
	if err := db.Create(&pageMenu).Error; err != nil {
		t.Fatalf("create page menu: %v", err)
	}

	roleAPI := model.SysApi{Path: "/api/v1/a", Method: "GET", Group: "测试", Description: "A", NeedAuth: true}
	if err := db.Create(&roleAPI).Error; err != nil {
		t.Fatalf("create role api: %v", err)
	}

	targetRole := model.SysRole{Name: "回滚目标角色", Code: "role-permission-smoke-rollover", Status: 1}
	if err := db.Create(&targetRole).Error; err != nil {
		t.Fatalf("create target role: %v", err)
	}
	if err := Role.AssignMenus(targetRole.ID, []uint{pageMenu.ID}); err != nil {
		t.Fatalf("AssignMenus error: %v", err)
	}
	if err := Role.AssignApis(targetRole.ID, []uint{roleAPI.ID}); err != nil {
		t.Fatalf("AssignApis error: %v", err)
	}

	router := buildRolePermissionSmokeRouter(1002, []uint{operatorRole.ID}, []string{operatorRole.Code})

	body := `{
		"menu_ids": [],
		"direct_api_ids": [],
		"scopes": [
			{"resource_code":"system:dept-management","data_scope":5,"dept_ids":[]}
		]
	}`
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/roles/"+toUintString(targetRole.ID)+"/permissions",
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, req)

	resp := decodeTestResponse(t, recorder)
	if resp.Code == response.SUCCESS {
		t.Fatalf("expected smoke save to reject unsupported scope, body=%s", recorder.Body.String())
	}
	if !strings.Contains(resp.Message, "仅本人") {
		t.Fatalf("unexpected error message: %s", resp.Message)
	}

	var roleMenuCount int64
	if err := db.Table("sys_role_menu").Where("sys_role_id = ?", targetRole.ID).Count(&roleMenuCount).Error; err != nil {
		t.Fatalf("count sys_role_menu: %v", err)
	}
	if roleMenuCount == 0 {
		t.Fatalf("expected role menus to remain after rejected smoke save")
	}

	var roleAPICount int64
	if err := db.Table("sys_role_api").Where("sys_role_id = ?", targetRole.ID).Count(&roleAPICount).Error; err != nil {
		t.Fatalf("count sys_role_api: %v", err)
	}
	if roleAPICount == 0 {
		t.Fatalf("expected role apis to remain after rejected smoke save")
	}
}

func buildRolePermissionSmokeRouter(userID uint, roleIDs []uint, roleCodes []string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.ContextUserIDKey, userID)
		c.Set(middleware.ContextRoleIDsKey, roleIDs)
		c.Set(middleware.ContextRoleCodesKey, roleCodes)
		c.Next()
	})
	router.Use(middleware.CasbinAuth())
	router.PUT("/api/v1/roles/:id/permissions", v1.Role.SavePermissions)
	return router
}
