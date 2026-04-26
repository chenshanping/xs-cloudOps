package modules

import (
	v1 "server/api/v1"
	"server/model/request"
	"server/router/registry"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterModule(&AuthModule{})
}

type AuthModule struct{}

func (m *AuthModule) Name() string {
	return "认证管理"
}

func (m *AuthModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	R(rg, "POST", "/auth/login", m.Name(), "登录", v1.Auth.Login,
		registry.WithRequest(request.LoginRequest{}))
	R(rg, "POST", "/auth/register", m.Name(), "注册", v1.Auth.Register,
		registry.WithRequest(request.RegisterRequest{}))
	R(rg, "POST", "/auth/send-email-code", m.Name(), "发送邮箱验证码", v1.Auth.SendEmailCode,
		registry.WithRequest(request.SendEmailCodeRequest{}))
	R(rg, "POST", "/auth/reset-password", m.Name(), "重置密码(Token)", v1.Auth.ResetPasswordByToken,
		registry.WithRequest(request.ResetPasswordByTokenRequest{}))
	R(rg, "POST", "/auth/reset-password-by-username", m.Name(), "重置密码(用户名)", v1.Auth.ResetPasswordByUserName,
		registry.WithRequest(request.ResetPasswordByUserNameRequest{}))
	R(rg, "POST", "/auth/reset-password-by-email", m.Name(), "重置密码(邮箱)", v1.Auth.ResetPasswordByEmail,
		registry.WithRequest(request.ResetPasswordByEmailRequest{}))
	// 刷新Token放在公开路由，允许过期 Token 请求
	R(rg, "POST", "/auth/refresh", m.Name(), "刷新Token", v1.Auth.RefreshToken)
}

func (m *AuthModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "POST", "/auth/logout", m.Name(), "登出", v1.Auth.Logout, registry.WithAuth())
	R(rg, "GET", "/auth/userinfo", m.Name(), "获取用户信息", v1.Auth.GetUserInfo, registry.WithAuth())
}

// R 注册路由并记录元信息
func R(rg *gin.RouterGroup, method, path, group, summary string, handler gin.HandlerFunc, opts ...registry.RouteOption) {
	fullPath := rg.BasePath() + path
	registry.Register(method, fullPath, group, summary, handler, opts...)
	switch method {
	case "GET":
		rg.GET(path, handler)
	case "POST":
		rg.POST(path, handler)
	case "PUT":
		rg.PUT(path, handler)
	case "DELETE":
		rg.DELETE(path, handler)
	}
}
