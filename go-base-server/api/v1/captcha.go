package v1

import (
	"github.com/gin-gonic/gin"

	"go-base-server/model/response"
	"go-base-server/service"
)

type CaptchaApi struct{}

var CaptchaAPI = new(CaptchaApi)

// GetCaptcha 获取图形验证码
func (a *CaptchaApi) GetCaptcha(c *gin.Context) {
	id, b64s, err := service.Captcha.GenerateCaptcha()
	if err != nil {
		response.Fail(c, "生成验证码失败")
		return
	}

	response.OkWithData(c, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}

// GetCaptchaConfig 获取验证码配置
func (a *CaptchaApi) GetCaptchaConfig(c *gin.Context) {
	response.OkWithData(c, gin.H{
		"login_captcha_enabled":    service.Captcha.IsLoginCaptchaEnabled(),
		"register_captcha_enabled": service.Captcha.IsRegisterCaptchaEnabled(),
		"register_email_verify":    service.Captcha.IsRegisterEmailVerifyEnabled(),
	})
}
