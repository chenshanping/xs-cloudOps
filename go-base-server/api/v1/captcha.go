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
		"login_captcha_type":       service.Captcha.GetCaptchaType(),
		"register_captcha_enabled": service.Captcha.IsRegisterCaptchaEnabled(),
		"register_email_verify":    service.Captcha.IsRegisterEmailVerifyEnabled(),
		"slider_captcha_bg":        service.Captcha.GetSliderCaptchaBg(),
	})
}

// GetSliderCaptcha 获取滑动验证码
func (a *CaptchaApi) GetSliderCaptcha(c *gin.Context) {
	data, err := service.Captcha.GenerateSliderCaptcha()
	if err != nil {
		response.Fail(c, "生成验证码失败")
		return
	}
	
	// 获取目标位置用于前端生成图片
	x, y, ok := service.Captcha.GetSliderCaptchaTarget(data.CaptchaID)
	if !ok {
		response.Fail(c, "获取验证码失败")
		return
	}
	
	response.OkWithData(c, gin.H{
		"captcha_id": data.CaptchaID,
		"bg_width":   data.BgWidth,
		"bg_height":  data.BgHeight,
		"slider_y":   y,
		"target_x":   x, // 前端需要知道缺口位置来绘制图片
	})
}

// VerifySliderCaptcha 验证滑动验证码
func (a *CaptchaApi) VerifySliderCaptcha(c *gin.Context) {
	var req struct {
		CaptchaID string `json:"captcha_id" binding:"required"`
		X         int    `json:"x"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	
	if service.Captcha.VerifySliderCaptcha(req.CaptchaID, req.X) {
		response.OkWithData(c, gin.H{"success": true})
	} else {
		response.Fail(c, "验证失败")
	}
}
