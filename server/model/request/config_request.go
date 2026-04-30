package request

// 批量读取配置请求
type ConfigKeysRequest struct {
	Keys []string `json:"keys" binding:"required" comment:"配置键列表"`
}

// 测试邮件请求
type TestEmailRequest struct {
	Email string `json:"email" binding:"required,email" comment:"接收测试邮件的邮箱地址"`
}
