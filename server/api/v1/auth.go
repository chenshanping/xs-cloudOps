package v1

import (
	"context"

	"server/middleware"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type AuthApi struct{}

var Auth = new(AuthApi)
var authFlow = service.NewAuthFlowService()

// Login 用户登录
func (a *AuthApi) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, flowErr := authFlow.Login(service.AuthLoginInput{
		Username:  req.Username,
		Password:  req.Password,
		CaptchaID: req.CaptchaID,
		Captcha:   req.Captcha,
	}, c.ClientIP(), c.Request.UserAgent())
	if flowErr != nil {
		switch flowErr.Kind {
		case service.AuthFlowErrorKindBadRequest:
			response.BadRequest(c, flowErr.Message)
		default:
			response.Fail(c, flowErr.Message)
		}
		return
	}

	response.OkWithData(c, gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}

// Register 用户注册
func (a *AuthApi) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 检查是否需要验证码
	if service.Captcha.IsRegisterCaptchaEnabled() {
		if req.CaptchaCode == "" {
			response.BadRequest(c, "请输入验证码")
			return
		}
		if !service.Captcha.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
			response.Fail(c, "验证码错误")
			return
		}
	}

	// 检查是否需要邮箱验证
	if service.User.IsEmailVerificationRequired() {
		if req.EmailCode == "" {
			response.BadRequest(c, "请输入邮箱验证码")
			return
		}
		if !service.Captcha.VerifyEmailCode(req.Email, req.EmailCode) {
			response.Fail(c, "邮箱验证码错误或已过期")
			return
		}
	}

	// 注册用户
	if err := service.User.Register(req.Username, req.Password, req.Email); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "注册成功")
}

// SendEmailCode 发送邮箱验证码
func (a *AuthApi) SendEmailCode(c *gin.Context) {
	var req request.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 检查验证码（防止恶意发送）
	if req.CaptchaID != "" && req.Captcha != "" {
		if !service.Captcha.VerifyCaptcha(req.CaptchaID, req.Captcha) {
			response.Fail(c, "验证码错误")
			return
		}
	}

	// 检查邮箱服务是否启用
	if !service.Email.IsEmailEnabled() {
		response.Fail(c, "邮箱服务未配置")
		return
	}

	// 生成并发送验证码
	code, err := service.Captcha.GenerateEmailCode(req.Email)
	if err != nil {
		response.Fail(c, "生成验证码失败")
		return
	}

	if err := service.Email.SendVerificationCode(req.Email, code); err != nil {
		response.Fail(c, "发送邮件失败: "+err.Error())
		return
	}

	response.OkWithMessage(c, "验证码已发送")
}

// ResetPasswordByToken 通过Token重置密码
func (a *AuthApi) ResetPasswordByToken(c *gin.Context) {
	var req request.ResetPasswordByTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if flowErr := authFlow.ResetPasswordByToken(context.Background(), req.Token, req.Password); flowErr != nil {
		response.Fail(c, flowErr.Message)
		return
	}

	response.OkWithMessage(c, "密码重置成功")
}

// ResetPasswordByEmail 通过邮箱验证码重置密码
func (a *AuthApi) ResetPasswordByEmail(c *gin.Context) {
	var req request.ResetPasswordByEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 验证邮箱验证码
	if !service.Captcha.VerifyEmailCode(req.Email, req.EmailCode) {
		response.Fail(c, "邮箱验证码错误或已过期")
		return
	}

	// 根据邮箱获取用户
	user, err := service.User.GetUserByEmail(req.Email)
	if err != nil {
		response.Fail(c, "该邮箱未注册")
		return
	}

	// 重置密码
	if err := service.User.ResetPassword(user.ID, req.NewPassword); err != nil {
		response.Fail(c, "重置密码失败")
		return
	}

	response.OkWithMessage(c, "密码重置成功")
}

// ResetPasswordByUserName 通过用户名重置密码
func (a *AuthApi) ResetPasswordByUserName(c *gin.Context) {
	var req request.ResetPasswordByUserNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	username := req.UserName
	password := req.NewPassword
	captchaId := req.CaptchaId
	captchaCode := req.Captcha
	// 检查验证码
	if captchaId == "" {
		if !service.Captcha.VerifyCaptcha(captchaId, captchaCode) {
			response.Fail(c, "验证码错误")
			return
		}
	}
	user, err := service.User.GetUserByUserName(username)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	// 重置密码
	if user != nil {
		if err := service.User.ResetPassword(user.ID, password); err != nil {
			response.Fail(c, "重置密码失败")
			return
		}
	}

	response.OkWithMessage(c, "密码重置成功")

}

// Logout 用户登出
func (a *AuthApi) Logout(c *gin.Context) {
	// 将当前Token加入黑名单
	if token, exists := c.Get(middleware.ContextTokenKey); exists {
		if tokenStr, ok := token.(string); ok {
			_ = utils.InvalidateToken(tokenStr)
		}
	}
	response.Ok(c)
}

// GetUserInfo 获取当前用户信息（优先从 Redis 缓存获取）
func (a *AuthApi) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	cache, flowErr := authFlow.GetCurrentUserInfo(userID)
	if flowErr != nil {
		response.Fail(c, flowErr.Message)
		return
	}

	response.OkWithData(c, gin.H{
		"user":        cache.User,
		"menus":       cache.Menus,
		"permissions": cache.Permissions,
	})
}

// RefreshToken 刷新Token
func (a *AuthApi) RefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 {
		token = token[7:]
	}

	newToken, err := utils.RefreshToken(token)
	if err != nil {
		response.Unauthorized(c, "Token刷新失败")
		return
	}

	response.OkWithData(c, gin.H{
		"token": newToken,
	})
}
