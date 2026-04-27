package request

// 测试邮件请求
type TestEmailRequest struct {
	Email string `json:"email" binding:"required,email" comment:"接收测试邮件的邮箱地址"`
}
