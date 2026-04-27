package service

import (
	"context"
	"errors"
	"testing"

	"server/model"
	"server/utils"
)

func TestAuthFlowServiceLoginReturnsLockedError(t *testing.T) {
	svc := NewAuthFlowService()
	svc.CheckLoginLock = func(username string) (bool, int) {
		if username != "alice" {
			t.Fatalf("unexpected username: %s", username)
		}
		return true, 5
	}

	_, err := svc.Login(AuthLoginInput{Username: "alice"}, "127.0.0.1", "ua")
	if err == nil {
		t.Fatalf("expected locked error")
	}
	if err.Kind != AuthFlowErrorKindFail {
		t.Fatalf("expected fail kind, got %s", err.Kind)
	}
	if err.Message != "账户已被锁定，请5分钟后重试" {
		t.Fatalf("unexpected message: %s", err.Message)
	}
}

func TestAuthFlowServiceLoginSuccessReturnsTokenAndUser(t *testing.T) {
	svc := NewAuthFlowService()
	svc.CheckLoginLock = func(string) (bool, int) { return false, 0 }
	svc.IsLoginCaptchaEnabled = func() bool { return false }
	svc.GetClientInfo = func(ip, ua string) utils.ClientInfo {
		return utils.ClientInfo{Location: "本地", Browser: "Chrome", OS: "Windows"}
	}
	svc.LoginUser = func(username, password string) (*model.SysUser, error) {
		if username != "alice" || password != "secret" {
			t.Fatalf("unexpected credentials: %s/%s", username, password)
		}
		return &model.SysUser{
			BaseModel: model.BaseModel{ID: 7},
			Username:  "alice",
			Roles: []model.SysRole{
				{BaseModel: model.BaseModel{ID: 1}, Code: "admin"},
			},
		}, nil
	}
	svc.ClearLoginRetry = func(username string) {
		if username != "alice" {
			t.Fatalf("unexpected clear retry user: %s", username)
		}
	}
	logs := 0
	svc.CreateLoginLog = func(log *model.SysLoginLog) error {
		logs++
		if log.Status != 1 {
			t.Fatalf("expected success log")
		}
		return nil
	}
	svc.GenerateToken = func(userID uint, username string, roleIDs []uint, roleCodes []string) (string, error) {
		if userID != 7 || username != "alice" {
			t.Fatalf("unexpected token args: %d/%s", userID, username)
		}
		if len(roleIDs) != 1 || roleIDs[0] != 1 {
			t.Fatalf("unexpected role ids: %#v", roleIDs)
		}
		if len(roleCodes) != 1 || roleCodes[0] != "admin" {
			t.Fatalf("unexpected role codes: %#v", roleCodes)
		}
		return "signed-token", nil
	}

	result, err := svc.Login(AuthLoginInput{
		Username: "alice",
		Password: "secret",
	}, "127.0.0.1", "ua")
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if result.Token != "signed-token" {
		t.Fatalf("unexpected token: %s", result.Token)
	}
	if result.User == nil || result.User.Username != "alice" {
		t.Fatalf("unexpected user: %#v", result.User)
	}
	if logs != 1 {
		t.Fatalf("expected 1 success log, got %d", logs)
	}
}

func TestAuthFlowServiceGetCurrentUserInfoReturnsCacheFirst(t *testing.T) {
	svc := NewAuthFlowService()
	cached := &UserInfoCache{
		User:        &model.SysUser{BaseModel: model.BaseModel{ID: 8}, Username: "bob"},
		Menus:       []model.SysMenu{{BaseModel: model.BaseModel{ID: 2}, Name: "Dashboard"}},
		Permissions: []string{"system:user:list"},
	}
	svc.GetUserInfoFromCache = func(userID uint) (*UserInfoCache, error) {
		if userID != 8 {
			t.Fatalf("unexpected user id: %d", userID)
		}
		return cached, nil
	}
	svc.GetUserInfo = func(uint) (*model.SysUser, error) {
		t.Fatalf("should not hit database when cache exists")
		return nil, nil
	}

	result, err := svc.GetCurrentUserInfo(8)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if result != cached {
		t.Fatalf("expected cached result")
	}
}

func TestAuthFlowServiceResetPasswordByTokenParsesUserID(t *testing.T) {
	svc := NewAuthFlowService()
	svc.RedisGet = func(ctx context.Context, key string) (string, error) {
		if key != "reset_password:token-1" {
			t.Fatalf("unexpected redis key: %s", key)
		}
		return "42", nil
	}
	svc.ResetPassword = func(userID uint, password string) error {
		if userID != 42 || password != "new-password" {
			t.Fatalf("unexpected reset args: %d/%s", userID, password)
		}
		return nil
	}
	deleted := false
	svc.RedisDel = func(ctx context.Context, key string) error {
		deleted = true
		if key != "reset_password:token-1" {
			t.Fatalf("unexpected delete key: %s", key)
		}
		return nil
	}

	if err := svc.ResetPasswordByToken(context.Background(), "token-1", "new-password"); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if !deleted {
		t.Fatalf("expected token to be deleted")
	}
}

func TestAuthFlowServiceResetPasswordByTokenRejectsInvalidUserID(t *testing.T) {
	svc := NewAuthFlowService()
	svc.RedisGet = func(context.Context, string) (string, error) {
		return "not-a-number", nil
	}

	err := svc.ResetPasswordByToken(context.Background(), "bad-token", "new-password")
	if err == nil {
		t.Fatalf("expected invalid user id error")
	}
	if err.Kind != AuthFlowErrorKindFail {
		t.Fatalf("expected fail kind, got %s", err.Kind)
	}
}

func TestAuthFlowServiceLoginTracksRetryOnFailure(t *testing.T) {
	svc := NewAuthFlowService()
	svc.CheckLoginLock = func(string) (bool, int) { return false, 0 }
	svc.IsLoginCaptchaEnabled = func() bool { return false }
	svc.GetClientInfo = func(string, string) utils.ClientInfo { return utils.ClientInfo{} }
	svc.LoginUser = func(string, string) (*model.SysUser, error) { return nil, errors.New("密码错误") }
	svc.IncrLoginRetry = func(string) (int, bool) { return 1, false }
	svc.GetLoginMaxRetry = func() int { return 5 }
	logs := 0
	svc.CreateLoginLog = func(log *model.SysLoginLog) error {
		logs++
		if log.Status != 0 {
			t.Fatalf("expected failure log")
		}
		return nil
	}

	_, err := svc.Login(AuthLoginInput{Username: "alice", Password: "bad"}, "127.0.0.1", "ua")
	if err == nil {
		t.Fatalf("expected login error")
	}
	if err.Message != "密码错误，还剩4次尝试机会" {
		t.Fatalf("unexpected message: %s", err.Message)
	}
	if logs != 1 {
		t.Fatalf("expected 1 failure log, got %d", logs)
	}
}
