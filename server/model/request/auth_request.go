package request

// 登录请求
type LoginRequest struct {
	Username  string `json:"username" binding:"required" comment:"用户名"`
	Password  string `json:"password" binding:"required" comment:"密码"`
	CaptchaID string `json:"captcha_id" comment:"验证码ID"`
	Captcha   string `json:"captcha" comment:"验证码"`
}

// 注册请求
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=20" comment:"用户名"`
	Password    string `json:"password" binding:"required,min=6" comment:"密码"`
	Email       string `json:"email" binding:"email" comment:"邮箱"`
	EmailCode   string `json:"email_code" comment:"邮箱验证码"`
	CaptchaID   string `json:"captcha_id" comment:"验证码ID"`
	CaptchaCode string `json:"captcha_code" comment:"验证码"`
}

// 发送邮箱验证码请求
type SendEmailCodeRequest struct {
	Email     string `json:"email" binding:"required,email" comment:"邮箱地址"`
	CaptchaID string `json:"captcha_id" comment:"验证码ID"`
	Captcha   string `json:"captcha" comment:"验证码"`
}

// 忘记密码请求
type ForgotPasswordRequest struct {
	Email     string `json:"email" binding:"required,email" comment:"邮箱地址"`
	CaptchaID string `json:"captcha_id" comment:"验证码ID"`
	Captcha   string `json:"captcha" comment:"验证码"`
}

// 重置密码请求
type ResetPasswordByTokenRequest struct {
	Token    string `json:"token" binding:"required" comment:"重置令牌"`
	Password string `json:"password" binding:"required,min=6" comment:"新密码"`
}

// 重置密码参数（通过用户名 + 图形验证码）
type ResetPasswordByUserNameRequest struct {
	UserName    string `json:"username" binding:"required" comment:"用户名"`
	NewPassword string `json:"new_password" binding:"required,min=6" comment:"新密码"`
	CaptchaId   string `json:"captcha_id" binding:"required" comment:"验证码ID"`
	Captcha     string `json:"captcha" binding:"required" comment:"验证码"`
}

// 重置密码参数（通过邮箱验证码）
type ResetPasswordByEmailRequest struct {
	Email       string `json:"email" binding:"required,email" comment:"邮箱地址"`
	EmailCode   string `json:"email_code" binding:"required" comment:"邮箱验证码"`
	NewPassword string `json:"new_password" binding:"required,min=6" comment:"新密码"`
}
