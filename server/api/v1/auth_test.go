package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/response"
	"server/utils"
)

func setupAuthAPITestEnv(t *testing.T) *gorm.DB {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.SysUser{}); err != nil {
		t.Fatalf("auto migrate auth api models: %v", err)
	}

	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}

	previousDB := global.DB
	previousRedis := global.Redis
	global.DB = db
	global.Redis = redis.NewClient(&redis.Options{Addr: redisServer.Addr()})

	t.Cleanup(func() {
		_ = global.Redis.Close()
		redisServer.Close()
		global.DB = previousDB
		global.Redis = previousRedis
	})

	return db
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) response.Response {
	t.Helper()

	var resp response.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v, body=%s", err, recorder.Body.String())
	}
	return resp
}

func TestResetPasswordByUserNameRejectsWhenCaptchaVerificationFails(t *testing.T) {
	setupAuthAPITestEnv(t)

	hashed, err := utils.HashPassword("old-password")
	if err != nil {
		t.Fatalf("hash old password: %v", err)
	}
	user := model.SysUser{
		Username: "alice",
		Password: hashed,
		Nickname: "alice",
		Status:   1,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password-by-username", strings.NewReader(`{"username":"alice","new_password":"123456","captcha_id":"fake-id","captcha":"fake-code"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	Auth.ResetPasswordByUserName(c)

	resp := decodeResponse(t, recorder)
	if resp.Code != response.ERROR || resp.Message != "验证码错误" {
		t.Fatalf("expected captcha error, got code=%d message=%q", resp.Code, resp.Message)
	}

	var saved model.SysUser
	if err := global.DB.Where("username = ?", "alice").First(&saved).Error; err != nil {
		t.Fatalf("query user: %v", err)
	}
	if !utils.CheckPassword("old-password", saved.Password) {
		t.Fatalf("expected password to remain unchanged")
	}
}

func TestResetPasswordByUserNameAllowsWhenCaptchaVerificationSucceeds(t *testing.T) {
	setupAuthAPITestEnv(t)

	hashed, err := utils.HashPassword("old-password")
	if err != nil {
		t.Fatalf("hash old password: %v", err)
	}
	user := model.SysUser{
		Username: "bob",
		Password: hashed,
		Nickname: "bob",
		Status:   1,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := global.Redis.Set(context.Background(), "captcha:captcha-ok", "pass-ok", 0).Err(); err != nil {
		t.Fatalf("seed captcha: %v", err)
	}

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password-by-username", strings.NewReader(`{"username":"bob","new_password":"654321","captcha_id":"captcha-ok","captcha":"pass-ok"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	Auth.ResetPasswordByUserName(c)

	resp := decodeResponse(t, recorder)
	if resp.Code != response.SUCCESS || resp.Message != "密码重置成功" {
		t.Fatalf("expected success, got code=%d message=%q", resp.Code, resp.Message)
	}

	var saved model.SysUser
	if err := global.DB.Where("username = ?", "bob").First(&saved).Error; err != nil {
		t.Fatalf("query user: %v", err)
	}
	if !utils.CheckPassword("654321", saved.Password) {
		t.Fatalf("expected password to be updated")
	}
}

func TestResetPasswordByUserNameReturnsFriendlyValidationMessage(t *testing.T) {
	setupAuthAPITestEnv(t)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password-by-username", strings.NewReader(`{"username":"bob","new_password":"123","captcha_id":"captcha-ok","captcha":"pass-ok"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	Auth.ResetPasswordByUserName(c)

	resp := decodeResponse(t, recorder)
	if resp.Code != response.BAD_REQUEST || resp.Message != "新密码至少 6 位" {
		t.Fatalf("expected friendly validation message, got code=%d message=%q", resp.Code, resp.Message)
	}
}

func TestResetPasswordByEmailDoesNotRevealUnknownEmail(t *testing.T) {
	setupAuthAPITestEnv(t)

	if err := global.Redis.Set(context.Background(), "email_code:missing@example.com", "123456", 0).Err(); err != nil {
		t.Fatalf("seed email code: %v", err)
	}

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password-by-email", strings.NewReader(`{"email":"missing@example.com","email_code":"123456","new_password":"654321"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	Auth.ResetPasswordByEmail(c)

	resp := decodeResponse(t, recorder)
	if resp.Code != response.SUCCESS || resp.Message != "密码重置成功" {
		t.Fatalf("expected generic success message, got code=%d message=%q", resp.Code, resp.Message)
	}
}
