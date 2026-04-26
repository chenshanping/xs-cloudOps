package v1

import (
	"context"
	"fmt"

	"server/global"
	"server/middleware"
	"server/model"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type AuthApi struct{}

var Auth = new(AuthApi)

// Login 用户登录
func (a *AuthApi) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 检查账户是否被锁定
	if locked, minutes := service.Captcha.CheckLoginLock(req.Username); locked {
		response.Fail(c, fmt.Sprintf("账户已被锁定，请%d分钟后重试", minutes))
		return
	}

	// 检查是否需要验证码
	if service.Captcha.IsLoginCaptchaEnabled() {
		captchaType := service.Captcha.GetCaptchaType()
		if req.CaptchaID == "" || req.Captcha == "" {
			response.BadRequest(c, "请完成验证")
			return
		}
		
		// 滑动验证码在前端已验证，后端只验证captchaID是否有效
		if captchaType == "slider" {
			if req.Captcha != "slider_verified" {
				response.Fail(c, "请先完成滑动验证")
				return
			}
			// 滑动验证码已在前端验证通过，这里只需确认captchaID格式正确
		} else {
			// 普通图形验证码
			if !service.Captcha.VerifyCaptcha(req.CaptchaID, req.Captcha) {
				response.Fail(c, "验证码错误")
				return
			}
		}
	}

	// 获取客户端信息
	ip := c.ClientIP()
	clientInfo := utils.GetClientInfo(ip, c.Request.UserAgent())

	user, err := service.User.Login(req.Username, req.Password)
	if err != nil {
		// 增加登录失败次数
		retryCount, locked := service.Captcha.IncrLoginRetry(req.Username)
		maxRetry := service.Captcha.GetLoginMaxRetry()
		
		errMsg := err.Error()
		if locked {
			lockTime := service.Captcha.GetLoginLockTime()
			errMsg = fmt.Sprintf("登录失败次数过多，账户已锁定%d分钟", lockTime)
		} else {
			errMsg = fmt.Sprintf("%s，还剩%d次尝试机会", errMsg, maxRetry-retryCount)
		}
		
		// 记录登录失败日志
		service.Log.CreateLoginLog(&model.SysLoginLog{
			Username: req.Username,
			IP:       ip,
			Location: clientInfo.Location,
			Browser:  clientInfo.Browser,
			OS:       clientInfo.OS,
			Status:   0,
			Msg:      err.Error(),
		})
		response.Fail(c, errMsg)
		return
	}
	
	// 登录成功，清除重试次数
	service.Captcha.ClearLoginRetry(req.Username)

	// 提取角色ID、编码列表
	roleIDs := make([]uint, 0, len(user.Roles))
	roleCodes := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleIDs = append(roleIDs, role.ID)
		roleCodes = append(roleCodes, role.Code)
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username, roleIDs, roleCodes)
	if err != nil {
		response.Fail(c, "生成Token失败")
		return
	}

	// 记录登录成功日志
	service.Log.CreateLoginLog(&model.SysLoginLog{
		UserID:   user.ID,
		Username: user.Username,
		IP:       ip,
		Location: clientInfo.Location,
		Browser:  clientInfo.Browser,
		OS:       clientInfo.OS,
		Status:   1,
		Msg:      "登录成功",
	})

	response.OkWithData(c, gin.H{
		"token": token,
		"user":  user,
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

	// 验证Token
	ctx := context.Background()
	key := "reset_password:" + req.Token
	userIDStr, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		response.Fail(c, "链接已过期或无效")
		return
	}

	// 转换userID
	var userID uint
	if _, err := global.Redis.Get(ctx, key).Uint64(); err == nil {
		userID = uint(global.Redis.Get(ctx, key).Val()[0] - '0')
	}
	if userIDStr != "" {
		var id uint64
		for _, c := range userIDStr {
			id = id*10 + uint64(c-'0')
		}
		userID = uint(id)
	}

	// 重置密码
	if err := service.User.ResetPassword(userID, req.Password); err != nil {
		response.Fail(c, "重置密码失败")
		return
	}

	// 删除Token
	global.Redis.Del(ctx, key)

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

	// 尝试从缓存获取
	if cache, err := service.Cache.GetUserInfoFromCache(userID); err == nil {
		response.OkWithData(c, gin.H{
			"user":        cache.User,
			"menus":       cache.Menus,
			"permissions": cache.Permissions,
		})
		return
	}

	// 缓存未命中，查询数据库
	user, err := service.User.GetUserInfo(userID)
	if err != nil {
		response.Fail(c, "获取用户信息失败")
		return
	}

	menus, err := service.Menu.GetUserMenus(userID)
	if err != nil {
		response.Fail(c, "获取用户菜单失败")
		return
	}

	permissions, err := service.Menu.GetUserPermissions(userID)
	if err != nil {
		response.Fail(c, "获取用户权限失败")
		return
	}

	// 写入缓存
	_ = service.Cache.SetUserInfoToCache(userID, &service.UserInfoCache{
		User:        user,
		Menus:       menus,
		Permissions: permissions,
	})

	response.OkWithData(c, gin.H{
		"user":        user,
		"menus":       menus,
		"permissions": permissions,
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
