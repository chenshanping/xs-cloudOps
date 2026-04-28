package auth

import (
	"context"
	"fmt"
	"strconv"

	"server/global"
	"server/model"
	"server/service/captcha"
	"server/service/core"
	"server/service/logsvc"
	"server/service/menu"
	"server/service/user"
	"server/utils"
)

type AuthFlowErrorKind string

const (
	AuthFlowErrorKindBadRequest   AuthFlowErrorKind = "bad_request"
	AuthFlowErrorKindFail         AuthFlowErrorKind = "fail"
	AuthFlowErrorKindUnauthorized AuthFlowErrorKind = "unauthorized"
)

type AuthFlowError struct {
	Kind    AuthFlowErrorKind
	Message string
	Err     error
}

func (e *AuthFlowError) Error() string {
	return e.Message
}

type AuthLoginResult struct {
	Token string
	User  *model.SysUser
}

type AuthLoginInput struct {
	Username  string
	Password  string
	CaptchaID string
	Captcha   string
}

type AuthFlowService struct {
	CheckLoginLock        func(username string) (bool, int)
	IsLoginCaptchaEnabled func() bool
	GetCaptchaType        func() string
	VerifyCaptcha         func(id, code string) bool
	GetClientInfo         func(ip, userAgent string) utils.ClientInfo
	LoginUser             func(username, password string) (*model.SysUser, error)
	IncrLoginRetry        func(username string) (int, bool)
	GetLoginMaxRetry      func() int
	GetLoginLockTime      func() int
	ClearLoginRetry       func(username string)
	CreateLoginLog        func(log *model.SysLoginLog) error
	GenerateToken         func(userID uint, username string, roleIDs []uint, roleCodes []string) (string, error)
	GetUserInfoFromCache  func(userID uint) (*core.UserInfoCache, error)
	GetUserInfo           func(userID uint) (*model.SysUser, error)
	GetUserMenus          func(userID uint) ([]model.SysMenu, error)
	GetUserPermissions    func(userID uint) ([]string, error)
	SetUserInfoToCache    func(userID uint, cache *core.UserInfoCache) error
	RedisGet              func(ctx context.Context, key string) (string, error)
	RedisDel              func(ctx context.Context, key string) error
	ResetPassword         func(userID uint, password string) error
}

func NewAuthFlowService() *AuthFlowService {
	return &AuthFlowService{
		CheckLoginLock:        captcha.Default.CheckLoginLock,
		IsLoginCaptchaEnabled: captcha.Default.IsLoginCaptchaEnabled,
		GetCaptchaType:        captcha.Default.GetCaptchaType,
		VerifyCaptcha:         captcha.Default.VerifyCaptcha,
		GetClientInfo:         utils.GetClientInfo,
		LoginUser:             user.Default.Login,
		IncrLoginRetry:        captcha.Default.IncrLoginRetry,
		GetLoginMaxRetry:      captcha.Default.GetLoginMaxRetry,
		GetLoginLockTime:      captcha.Default.GetLoginLockTime,
		ClearLoginRetry:       captcha.Default.ClearLoginRetry,
		CreateLoginLog:        logsvc.Default.CreateLoginLog,
		GenerateToken:         utils.GenerateToken,
		GetUserInfoFromCache:  core.Default.GetUserInfoFromCache,
		GetUserInfo:           user.Default.GetUserInfo,
		GetUserMenus:          menu.Default.GetUserMenus,
		GetUserPermissions:    menu.Default.GetUserPermissions,
		SetUserInfoToCache:    core.Default.SetUserInfoToCache,
		RedisGet: func(ctx context.Context, key string) (string, error) {
			return global.Redis.Get(ctx, key).Result()
		},
		RedisDel: func(ctx context.Context, key string) error {
			return global.Redis.Del(ctx, key).Err()
		},
		ResetPassword: user.Default.ResetPassword,
	}
}

func (s *AuthFlowService) Login(input AuthLoginInput, ip, userAgent string) (*AuthLoginResult, *AuthFlowError) {
	if locked, minutes := s.CheckLoginLock(input.Username); locked {
		return nil, &AuthFlowError{
			Kind:    AuthFlowErrorKindFail,
			Message: fmt.Sprintf("账户已被锁定，请%d分钟后重试", minutes),
		}
	}

	if s.IsLoginCaptchaEnabled() {
		captchaType := s.GetCaptchaType()
		if input.CaptchaID == "" || input.Captcha == "" {
			return nil, &AuthFlowError{
				Kind:    AuthFlowErrorKindBadRequest,
				Message: "请完成验证",
			}
		}
		if captchaType == "slider" {
			if input.Captcha != "slider_verified" {
				return nil, &AuthFlowError{
					Kind:    AuthFlowErrorKindFail,
					Message: "请先完成滑动验证",
				}
			}
		} else if !s.VerifyCaptcha(input.CaptchaID, input.Captcha) {
			return nil, &AuthFlowError{
				Kind:    AuthFlowErrorKindFail,
				Message: "验证码错误",
			}
		}
	}

	clientInfo := s.GetClientInfo(ip, userAgent)
	user, err := s.LoginUser(input.Username, input.Password)
	if err != nil {
		retryCount, locked := s.IncrLoginRetry(input.Username)
		maxRetry := s.GetLoginMaxRetry()
		errMsg := err.Error()
		if locked {
			errMsg = fmt.Sprintf("登录失败次数过多，账户已锁定%d分钟", s.GetLoginLockTime())
		} else {
			errMsg = fmt.Sprintf("%s，还剩%d次尝试机会", errMsg, maxRetry-retryCount)
		}
		_ = s.CreateLoginLog(&model.SysLoginLog{
			Username: input.Username,
			IP:       ip,
			Location: clientInfo.Location,
			Browser:  clientInfo.Browser,
			OS:       clientInfo.OS,
			Status:   0,
			Msg:      err.Error(),
		})
		return nil, &AuthFlowError{
			Kind:    AuthFlowErrorKindFail,
			Message: errMsg,
			Err:     err,
		}
	}

	s.ClearLoginRetry(input.Username)

	roleIDs := make([]uint, 0, len(user.Roles))
	roleCodes := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleIDs = append(roleIDs, role.ID)
		roleCodes = append(roleCodes, role.Code)
	}

	token, err := s.GenerateToken(user.ID, user.Username, roleIDs, roleCodes)
	if err != nil {
		return nil, &AuthFlowError{
			Kind:    AuthFlowErrorKindFail,
			Message: "生成Token失败",
			Err:     err,
		}
	}

	_ = s.CreateLoginLog(&model.SysLoginLog{
		UserID:   user.ID,
		Username: user.Username,
		IP:       ip,
		Location: clientInfo.Location,
		Browser:  clientInfo.Browser,
		OS:       clientInfo.OS,
		Status:   1,
		Msg:      "登录成功",
	})

	return &AuthLoginResult{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthFlowService) GetCurrentUserInfo(userID uint) (*core.UserInfoCache, *AuthFlowError) {
	if cache, err := s.GetUserInfoFromCache(userID); err == nil {
		return cache, nil
	}

	user, err := s.GetUserInfo(userID)
	if err != nil {
		return nil, &AuthFlowError{
			Kind:    AuthFlowErrorKindFail,
			Message: "获取用户信息失败",
			Err:     err,
		}
	}

	menus, err := s.GetUserMenus(userID)
	if err != nil {
		return nil, &AuthFlowError{
			Kind:    AuthFlowErrorKindFail,
			Message: "获取用户菜单失败",
			Err:     err,
		}
	}

	permissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return nil, &AuthFlowError{
			Kind:    AuthFlowErrorKindFail,
			Message: "获取用户权限失败",
			Err:     err,
		}
	}

	cache := &core.UserInfoCache{
		User:        user,
		Menus:       menus,
		Permissions: permissions,
	}
	_ = s.SetUserInfoToCache(userID, cache)
	return cache, nil
}

func (s *AuthFlowService) ResetPasswordByToken(ctx context.Context, token, password string) *AuthFlowError {
	key := "reset_password:" + token
	userIDStr, err := s.RedisGet(ctx, key)
	if err != nil {
		return &AuthFlowError{
			Kind:    AuthFlowErrorKindFail,
			Message: "链接已过期或无效",
			Err:     err,
		}
	}

	userIDValue, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return &AuthFlowError{
			Kind:    AuthFlowErrorKindFail,
			Message: "链接已过期或无效",
			Err:     err,
		}
	}

	if err := s.ResetPassword(uint(userIDValue), password); err != nil {
		return &AuthFlowError{
			Kind:    AuthFlowErrorKindFail,
			Message: "重置密码失败",
			Err:     err,
		}
	}

	_ = s.RedisDel(ctx, key)
	return nil
}
