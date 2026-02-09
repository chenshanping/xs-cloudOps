package v1

import (
	"go-base-server/model/request"
	"strconv"

	"go-base-server/model"
	"go-base-server/model/response"
	"go-base-server/service"

	"github.com/gin-gonic/gin"
)

type ConfigApi struct{}

var Config = new(ConfigApi)

// GetConfigList 获取配置列表
func (a *ConfigApi) GetConfigList(c *gin.Context) {
	key := c.Query("key")
	configs, err := service.Config.GetConfigList(key)
	if err != nil {
		response.Fail(c, "获取配置列表失败")
		return
	}
	response.OkWithData(c, configs)
}

// GetConfigByKey 根据key获取配置
func (a *ConfigApi) GetConfigByKey(c *gin.Context) {
	key := c.Param("key")
	config, err := service.Config.GetConfigByKey(key)
	if err != nil {
		response.NotFound(c, "配置不存在")
		return
	}
	response.OkWithData(c, config)
}

// GetConfigsByKeys 批量获取配置
func (a *ConfigApi) GetConfigsByKeys(c *gin.Context) {
	var req struct {
		Keys []string `json:"keys"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	configs, err := service.Config.GetConfigsByKeys(req.Keys)
	if err != nil {
		response.Fail(c, "获取配置失败")
		return
	}
	response.OkWithData(c, configs)
}

// CreateConfig 创建配置
func (a *ConfigApi) CreateConfig(c *gin.Context) {
	var config model.SysConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.Config.CreateConfig(&config); err != nil {
		response.Fail(c, "创建配置失败")
		return
	}
	response.OkWithData(c, config)
}

// UpdateConfig 更新配置
func (a *ConfigApi) UpdateConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.Config.UpdateConfig(uint(id), data); err != nil {
		response.Fail(c, "更新配置失败")
		return
	}
	response.OkWithMessage(c, "更新成功")
}

// BatchUpdateConfigs 批量更新配置
func (a *ConfigApi) BatchUpdateConfigs(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.Config.BatchUpdateConfigs(req); err != nil {
		response.Fail(c, "更新配置失败")
		return
	}
	response.OkWithMessage(c, "更新成功")
}

// DeleteConfig 删除配置
func (a *ConfigApi) DeleteConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := service.Config.DeleteConfig(uint(id)); err != nil {
		response.Fail(c, "删除配置失败")
		return
	}
	response.OkWithMessage(c, "删除成功")
}

// SendTestEmail 发送测试邮件
func (a *ConfigApi) SendTestEmail(c *gin.Context) {
	//var req struct {
	//	Email string `json:"email" binding:"required,email"`
	//}
	var req request.TestEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入有效的邮箱地址")
		return
	}

	// 获取系统名称
	sysName := "Go RBAC Admin"
	if config, err := service.Config.GetConfigByKey("sys_name"); err == nil && config.Value != "" {
		sysName = config.Value
	}

	subject := "【" + sysName + "】测试邮件"
	body := `
		<div style="padding: 20px; font-family: Arial, sans-serif;">
			<h2 style="color: #1890ff;">` + sysName + `</h2>
			<p>您好！</p>
			<p>这是一封测试邮件，用于验证您的SMTP邮箱配置是否正确。</p>
			<p style="color: #52c41a; font-size: 18px; font-weight: bold; padding: 20px 0;">✅ 邮箱配置正确，邮件发送成功！</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="color: #999; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
		</div>
	`

	if err := service.Email.SendEmail(req.Email, subject, body); err != nil {
		response.Fail(c, "发送失败: "+err.Error())
		return
	}

	response.OkWithMessage(c, "测试邮件发送成功")
}
