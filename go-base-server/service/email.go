package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"

	"go-base-server/global"
)

type EmailService struct{}

var Email = new(EmailService)

// EmailConfig 邮箱配置
type EmailConfig struct {
	Host     string // SMTP服务器地址
	Port     int    // SMTP端口
	Username string // 发件人邮箱
	Password string // 邮箱密码/授权码
	FromName string // 发件人名称
}

// GetEmailConfig 从数据库获取邮箱配置
func (s *EmailService) GetEmailConfig() (*EmailConfig, error) {
	keys := []string{
		"email_smtp_host",
		"email_smtp_port",
		"email_username",
		"email_password",
		"email_from_name",
	}

	configs, err := Config.GetConfigsByKeys(keys)
	if err != nil {
		return nil, err
	}

	port := 587 // 默认端口
	if p, ok := configs["email_smtp_port"]; ok && p.Value != "" {
		if portInt, err := strconv.Atoi(p.Value); err == nil {
			port = portInt
		}
	}

	config := &EmailConfig{
		Port: port,
	}

	if v, ok := configs["email_smtp_host"]; ok {
		config.Host = v.Value
	}
	if v, ok := configs["email_username"]; ok {
		config.Username = v.Value
	}
	if v, ok := configs["email_password"]; ok {
		config.Password = v.Value
	}
	if v, ok := configs["email_from_name"]; ok {
		config.FromName = v.Value
	}

	return config, nil
}

// SendEmail 发送邮件
func (s *EmailService) SendEmail(to, subject, body string) error {
	config, err := s.GetEmailConfig()
	if err != nil {
		return fmt.Errorf("获取邮箱配置失败: %v", err)
	}

	if config.Host == "" || config.Username == "" || config.Password == "" {
		return fmt.Errorf("邮箱配置不完整")
	}

	// 构建发件人显示名称
	from := config.Username
	if config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", config.FromName, config.Username)
	}

	// 构建邮件内容
	msg := s.buildMessage(from, to, subject, body)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// 根据端口决定是否使用TLS
	if config.Port == 465 {
		return s.sendMailSSL(addr, auth, config.Username, to, msg)
	}

	return smtp.SendMail(addr, auth, config.Username, []string{to}, msg)
}

// buildMessage 构建邮件消息
func (s *EmailService) buildMessage(from, to, subject, body string) []byte {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(body)

	return []byte(message.String())
}

// sendMailSSL 使用SSL发送邮件(465端口)
func (s *EmailService) sendMailSSL(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	host := strings.Split(addr, ":")[0]

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

// SendVerificationCode 发送验证码邮件
func (s *EmailService) SendVerificationCode(to, code string) error {
	// 获取系统名称
	sysName := "Go RBAC Admin"
	if config, err := Config.GetConfigByKey("sys_name"); err == nil && config.Value != "" {
		sysName = config.Value
	}

	subject := fmt.Sprintf("【%s】邮箱验证码", sysName)
	body := fmt.Sprintf(`
		<div style="padding: 20px; font-family: Arial, sans-serif;">
			<h2 style="color: #1890ff;">%s</h2>
			<p>您好！</p>
			<p>您的验证码是：</p>
			<div style="font-size: 32px; font-weight: bold; color: #1890ff; padding: 20px 0;">%s</div>
			<p>验证码有效期为 <strong>10分钟</strong>，请尽快使用。</p>
			<p>如果您没有请求此验证码，请忽略此邮件。</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="color: #999; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
		</div>
	`, sysName, code)

	return s.SendEmail(to, subject, body)
}

// SendResetPasswordEmail 发送重置密码邮件
func (s *EmailService) SendResetPasswordEmail(to, token string) error {
	// 获取系统名称和前端地址
	sysName := "Go RBAC Admin"
	frontendURL := "http://localhost:5173"

	if config, err := Config.GetConfigByKey("sys_name"); err == nil && config.Value != "" {
		sysName = config.Value
	}
	if config, err := Config.GetConfigByKey("frontend_url"); err == nil && config.Value != "" {
		frontendURL = config.Value
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	subject := fmt.Sprintf("【%s】重置密码", sysName)
	body := fmt.Sprintf(`
		<div style="padding: 20px; font-family: Arial, sans-serif;">
			<h2 style="color: #1890ff;">%s</h2>
			<p>您好！</p>
			<p>您正在申请重置密码，请点击下面的链接进行重置：</p>
			<p style="padding: 20px 0;">
				<a href="%s" style="background: #1890ff; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 4px;">重置密码</a>
			</p>
			<p>如果按钮无法点击，请复制以下链接到浏览器打开：</p>
			<p style="color: #666; word-break: break-all;">%s</p>
			<p>链接有效期为 <strong>30分钟</strong>，请尽快使用。</p>
			<p>如果您没有请求重置密码，请忽略此邮件。</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="color: #999; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
		</div>
	`, sysName, resetURL, resetURL)

	return s.SendEmail(to, subject, body)
}

// IsEmailEnabled 检查邮箱功能是否启用
func (s *EmailService) IsEmailEnabled() bool {
	config, err := s.GetEmailConfig()
	if err != nil {
		global.Log.Warnf("获取邮箱配置失败: %v", err)
		return false
	}
	return config.Host != "" && config.Username != "" && config.Password != ""
}
