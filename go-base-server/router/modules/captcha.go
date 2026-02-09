package modules

import (
	"github.com/gin-gonic/gin"

	v1 "go-base-server/api/v1"
)

func init() {
	RegisterModule(&CaptchaModule{})
}

type CaptchaModule struct{}

func (m *CaptchaModule) Name() string {
	return "验证码管理"
}

func (m *CaptchaModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 公开路由，无需认证
	R(rg, "GET", "/captcha", m.Name(), "获取图形验证码", v1.CaptchaAPI.GetCaptcha)
	R(rg, "GET", "/captcha/config", m.Name(), "获取验证码配置", v1.CaptchaAPI.GetCaptchaConfig)
	R(rg, "GET", "/captcha/slider", m.Name(), "获取滑动验证码", v1.CaptchaAPI.GetSliderCaptcha)
	R(rg, "POST", "/captcha/slider/verify", m.Name(), "验证滑动验证码", v1.CaptchaAPI.VerifySliderCaptcha)
}

func (m *CaptchaModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	// 验证码相关接口都是公开的
}
