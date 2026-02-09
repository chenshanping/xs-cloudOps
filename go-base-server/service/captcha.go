package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/mojocn/base64Captcha"

	"go-base-server/global"
)

type CaptchaService struct{}

var Captcha = new(CaptchaService)

// Redis 存储验证码
type redisStore struct {
	expiration time.Duration
}

func newRedisStore() *redisStore {
	return &redisStore{
		expiration: 5 * time.Minute, // 验证码5分钟过期
	}
}

func (s *redisStore) Set(id string, value string) error {
	ctx := context.Background()
	return global.Redis.Set(ctx, "captcha:"+id, value, s.expiration).Err()
}

func (s *redisStore) Get(id string, clear bool) string {
	ctx := context.Background()
	key := "captcha:" + id
	val, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	if clear {
		global.Redis.Del(ctx, key)
	}
	return val
}

func (s *redisStore) Verify(id, answer string, clear bool) bool {
	val := s.Get(id, clear)
	return val == answer
}

var store = newRedisStore()

// GenerateCaptcha 生成图形验证码
func (s *CaptchaService) GenerateCaptcha() (id, b64s string, err error) {
	// 配置验证码参数
	driver := base64Captcha.NewDriverDigit(
		80,  // 高度
		240, // 宽度
		5,   // 验证码长度
		0.7, // 最大倾斜角度
		80,  // 干扰点数量
	)

	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err = captcha.Generate()
	return
}

// VerifyCaptcha 验证验证码
func (s *CaptchaService) VerifyCaptcha(id, code string) bool {
	if id == "" || code == "" {
		return false
	}
	return store.Verify(id, code, true)
}

// IsCaptchaEnabled 检查是否启用登录验证码
func (s *CaptchaService) IsLoginCaptchaEnabled() bool {
	config, err := Config.GetConfigByKey("login_captcha_enabled")
	if err != nil {
		return false // 默认不启用
	}
	return config.Value == "1" || config.Value == "true"
}

// IsRegisterCaptchaEnabled 检查是否启用注册验证码
func (s *CaptchaService) IsRegisterCaptchaEnabled() bool {
	config, err := Config.GetConfigByKey("register_captcha_enabled")
	if err != nil {
		return false
	}
	return config.Value == "1" || config.Value == "true"
}

func (s *CaptchaService) IsRegisterEmailVerifyEnabled() bool {
	config, err := Config.GetConfigByKey("register_email_verify")
	if err != nil {
		return false
	}
	return config.Value == "1" || config.Value == "true"

}

// GenerateEmailCode 生成邮箱验证码（6位数字）
func (s *CaptchaService) GenerateEmailCode(email string) (string, error) {
	code := generateRandomCode(6)
	ctx := context.Background()
	key := "email_code:" + email
	err := global.Redis.Set(ctx, key, code, 10*time.Minute).Err() // 10分钟过期
	return code, err
}

// VerifyEmailCode 验证邮箱验证码
func (s *CaptchaService) VerifyEmailCode(email, code string) bool {
	ctx := context.Background()
	key := "email_code:" + email
	val, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	if val == code {
		global.Redis.Del(ctx, key)
		return true
	}
	return false
}

// generateRandomCode 生成指定长度的随机数字验证码
func generateRandomCode(length int) string {
	const digits = "0123456789"
	code := make([]byte, length)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code[i] = digits[n.Int64()]
	}
	return string(code)
}
