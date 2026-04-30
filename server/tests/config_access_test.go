package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	v1 "server/api/v1"
	"server/config"
	"server/global"
	"server/middleware"
	"server/model"
	"server/model/response"
	"server/service"
	"server/utils"
)

type configAccessResponse struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Data    map[string]model.SysConfig `json:"data"`
}

func setupConfigAccessTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysConfig{}, &model.SysRole{}); err != nil {
		t.Fatalf("auto migrate config access models: %v", err)
	}

	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})

	return db
}

func setupConfigAccessRedis(t *testing.T) *miniredis.Miniredis {
	t.Helper()

	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	previousRedis := global.Redis
	global.Redis = client
	t.Cleanup(func() {
		_ = client.Close()
		server.Close()
		global.Redis = previousRedis
	})

	return server
}

func setupConfigAccessJWTConfig(t *testing.T) {
	t.Helper()

	previousConfig := global.Config
	global.Config = &config.Config{
		JWT: config.JWT{
			Secret:        "config-access-test-secret",
			Expires:       3600,
			RefreshWindow: 3600,
			Issuer:        "config-access-test",
		},
	}
	t.Cleanup(func() {
		global.Config = previousConfig
	})
}

func setupConfigAccessEnforcer(t *testing.T, db *gorm.DB) *casbin.Enforcer {
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

func seedConfigAccessData(t *testing.T, db *gorm.DB) {
	t.Helper()

	configs := []model.SysConfig{
		{Name: "系统名称", Key: "sys_name", Value: "Go Base Test", ValueType: "string"},
		{Name: "登录标题", Key: "login_title", Value: "欢迎回来", ValueType: "string"},
		{Name: "用户身份按钮显示", Key: "user_profile_button_visible", Value: "true", ValueType: "string"},
		{Name: "邮箱密码", Key: "email_password", Value: "smtp-secret", ValueType: "string"},
		{Name: "MinIO配置", Key: "storage_minio_config", Value: `{"secret_access_key":"minio-secret"}`, ValueType: "json"},
		{
			Name:      "公开配置键",
			Key:       service.PublicConfigKeysConfigKey,
			Value:     `["sys_name","user_profile_button_visible","email_password"]`,
			ValueType: "json",
		},
	}
	for _, item := range configs {
		if err := db.Create(&item).Error; err != nil {
			t.Fatalf("create config %s: %v", item.Key, err)
		}
	}
}

func seedConfigAccessRoles(t *testing.T, db *gorm.DB) (model.SysRole, model.SysRole) {
	t.Helper()

	configAdminRole := model.SysRole{Name: "配置管理员", Code: "config-admin", Status: 1}
	normalRole := model.SysRole{Name: "普通用户", Code: "user", Status: 1}

	if err := db.Create(&configAdminRole).Error; err != nil {
		t.Fatalf("create config admin role: %v", err)
	}
	if err := db.Create(&normalRole).Error; err != nil {
		t.Fatalf("create normal role: %v", err)
	}

	return configAdminRole, normalRole
}

func buildConfigAccessTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	public := router.Group("/api/v1")
	public.POST("/configs/keys", v1.Config.GetConfigsByKeys)

	private := router.Group("/api/v1")
	private.Use(middleware.JWTAuth())
	private.Use(middleware.CasbinAuth())
	private.GET("/configs", v1.Config.GetConfigList)

	return router
}

func performPublicConfigAccessRequest(t *testing.T, router *gin.Engine, keys []string, token string) configAccessResponse {
	t.Helper()

	body, err := json.Marshal(gin.H{"keys": keys})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/configs/keys", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	var resp configAccessResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v, body=%s", err, recorder.Body.String())
	}
	return resp
}

func performPrivateConfigListRequest(t *testing.T, router *gin.Engine, token string) response.Response {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/configs", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	var resp response.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v, body=%s", err, recorder.Body.String())
	}
	return resp
}

func issueConfigAccessToken(t *testing.T, userID uint, username string, role model.SysRole) string {
	t.Helper()

	token, err := utils.GenerateToken(userID, username, []uint{role.ID}, []string{role.Code})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return token
}

func decodeConfigListData(t *testing.T, resp response.Response) []model.SysConfig {
	t.Helper()

	bytes, err := json.Marshal(resp.Data)
	if err != nil {
		t.Fatalf("marshal config list data: %v", err)
	}

	var configs []model.SysConfig
	if err := json.Unmarshal(bytes, &configs); err != nil {
		t.Fatalf("unmarshal config list data: %v", err)
	}
	return configs
}

func TestConfigAccessUsesBackendPublicWhitelistAndPrivateConfigPermissions(t *testing.T) {
	db := setupConfigAccessTestDB(t)
	setupConfigAccessRedis(t)
	setupConfigAccessJWTConfig(t)
	enforcer := setupConfigAccessEnforcer(t, db)
	seedConfigAccessData(t, db)
	configAdminRole, normalRole := seedConfigAccessRoles(t, db)

	if _, err := enforcer.AddPolicy(configAdminRole.Code, "/api/v1/configs", http.MethodGet); err != nil {
		t.Fatalf("grant config list policy: %v", err)
	}
	if err := enforcer.SavePolicy(); err != nil {
		t.Fatalf("save casbin policies: %v", err)
	}

	router := buildConfigAccessTestEngine()

	t.Run("anonymous request reads backend-defined public keys", func(t *testing.T) {
		resp := performPublicConfigAccessRequest(t, router, []string{"sys_name", "user_profile_button_visible", "login_title"}, "")
		if resp.Code != response.SUCCESS {
			t.Fatalf("response code = %d, want %d", resp.Code, response.SUCCESS)
		}
		if got := resp.Data["sys_name"].Value; got != "Go Base Test" {
			t.Fatalf("sys_name = %q, want %q", got, "Go Base Test")
		}
		if got := resp.Data["user_profile_button_visible"].Value; got != "true" {
			t.Fatalf("user_profile_button_visible = %q, want %q", got, "true")
		}
		if got := resp.Data["login_title"].Value; got != "欢迎回来" {
			t.Fatalf("login_title = %q, want %q", got, "欢迎回来")
		}
	})

	t.Run("database public whitelist config no longer changes exposed keys", func(t *testing.T) {
		resp := performPublicConfigAccessRequest(t, router, []string{"login_title"}, "")
		if resp.Code != response.SUCCESS {
			t.Fatalf("response code = %d, want %d", resp.Code, response.SUCCESS)
		}
		if got := resp.Data["login_title"].Value; got != "欢迎回来" {
			t.Fatalf("login_title = %q, want %q", got, "欢迎回来")
		}
	})

	t.Run("configured whitelist cannot force sensitive keys public", func(t *testing.T) {
		resp := performPublicConfigAccessRequest(t, router, []string{"email_password", "storage_minio_config"}, "")
		if resp.Code != response.SUCCESS {
			t.Fatalf("response code = %d, want %d", resp.Code, response.SUCCESS)
		}
		if _, exists := resp.Data["email_password"]; exists {
			t.Fatalf("email_password should never be returned publicly")
		}
		if _, exists := resp.Data["storage_minio_config"]; exists {
			t.Fatalf("storage_minio_config should never be returned publicly")
		}
	})

	t.Run("valid token does not elevate public config access", func(t *testing.T) {
		token := issueConfigAccessToken(t, 1001, "normal-user", normalRole)
		resp := performPublicConfigAccessRequest(t, router, []string{"sys_name", "email_password"}, token)
		if resp.Code != response.SUCCESS {
			t.Fatalf("response code = %d, want %d", resp.Code, response.SUCCESS)
		}
		if got := resp.Data["sys_name"].Value; got != "Go Base Test" {
			t.Fatalf("sys_name = %q, want %q", got, "Go Base Test")
		}
		if _, exists := resp.Data["email_password"]; exists {
			t.Fatalf("valid token unexpectedly elevated public config access")
		}
	})

	t.Run("private config list requires config permission", func(t *testing.T) {
		token := issueConfigAccessToken(t, 1002, "normal-user", normalRole)
		resp := performPrivateConfigListRequest(t, router, token)
		if resp.Code != response.FORBIDDEN {
			t.Fatalf("response code = %d, want %d", resp.Code, response.FORBIDDEN)
		}
	})

	t.Run("config admin can read protected config through private list api", func(t *testing.T) {
		token := issueConfigAccessToken(t, 1003, "config-admin", configAdminRole)
		resp := performPrivateConfigListRequest(t, router, token)
		if resp.Code != response.SUCCESS {
			t.Fatalf("response code = %d, want %d, message=%s", resp.Code, response.SUCCESS, resp.Message)
		}

		configs := decodeConfigListData(t, resp)
		byKey := make(map[string]model.SysConfig, len(configs))
		for _, item := range configs {
			byKey[item.Key] = item
		}

		if got := byKey["email_password"].Value; got != "smtp-secret" {
			t.Fatalf("email_password = %q, want %q", got, "smtp-secret")
		}
		if got := byKey["storage_minio_config"].Value; got != `{"secret_access_key":"minio-secret"}` {
			t.Fatalf("storage_minio_config = %q, want expected secret json", got)
		}
	})
}

func TestConfigAccessFallsBackToDefaultPublicKeysWhenWhitelistConfigInvalid(t *testing.T) {
	db := setupConfigAccessTestDB(t)
	setupConfigAccessRedis(t)
	setupConfigAccessJWTConfig(t)

	configs := []model.SysConfig{
		{Name: "登录标题", Key: "login_title", Value: "欢迎回来", ValueType: "string"},
		{Name: "公开配置键", Key: service.PublicConfigKeysConfigKey, Value: `{"invalid":true}`, ValueType: "json"},
	}
	for _, item := range configs {
		if err := db.Create(&item).Error; err != nil {
			t.Fatalf("create config %s: %v", item.Key, err)
		}
	}

	router := buildConfigAccessTestEngine()
	resp := performPublicConfigAccessRequest(t, router, []string{"login_title"}, "")
	if resp.Code != response.SUCCESS {
		t.Fatalf("response code = %d, want %d", resp.Code, response.SUCCESS)
	}
	if got := resp.Data["login_title"].Value; got != "欢迎回来" {
		t.Fatalf("login_title = %q, want %q", got, "欢迎回来")
	}
}
