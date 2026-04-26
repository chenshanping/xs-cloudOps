package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	mathRand "math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"

	"server/global"
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

// GetCaptchaType 获取验证码类型
func (s *CaptchaService) GetCaptchaType() string {
	config, err := Config.GetConfigByKey("login_captcha_type")
	if err != nil {
		return "digit" // 默认数字验证码
	}
	return config.Value
}

// GetSliderCaptchaBg 获取滑动验证码背景图
func (s *CaptchaService) GetSliderCaptchaBg() string {
	config, err := Config.GetConfigByKey("slider_captcha_bg")
	if err != nil {
		return "" // 默认为空，使用渐变背景
	}
	return config.Value
}

// SliderCaptchaData 滑动验证码数据
type SliderCaptchaData struct {
	CaptchaID string `json:"captcha_id"`
	BgWidth   int    `json:"bg_width"`
	BgHeight  int    `json:"bg_height"`
	SliderY   int    `json:"slider_y"`
}

// GenerateSliderCaptcha 生成滑动验证码
func (s *CaptchaService) GenerateSliderCaptcha() (*SliderCaptchaData, error) {
	// 生成唯一ID
	captchaID := uuid.New().String()
	
	// 设置参数
	bgWidth := 280
	bgHeight := 160
	sliderSize := 40
	
	// 生成随机X位置（滑块的正确位置）- 确保x在有效范围内
	mathRand.Seed(time.Now().UnixNano())
	xPos := mathRand.Intn(bgWidth-sliderSize*2-20) + sliderSize + 10 // 确保滑块不会太边缘
	yPos := mathRand.Intn(bgHeight-sliderSize-20) + 10
	
	// 将正确位置存储到Redis
	ctx := context.Background()
	key := "slider_captcha:" + captchaID
	data := map[string]int{"x": xPos, "y": yPos}
	jsonData, _ := json.Marshal(data)
	global.Redis.Set(ctx, key, string(jsonData), 5*time.Minute)
	
	return &SliderCaptchaData{
		CaptchaID: captchaID,
		BgWidth:   bgWidth,
		BgHeight:  bgHeight,
		SliderY:   yPos,
	}, nil
}

// VerifySliderCaptcha 验证滑动验证码
func (s *CaptchaService) VerifySliderCaptcha(captchaID string, x int) bool {
	ctx := context.Background()
	key := "slider_captcha:" + captchaID
	
	jsonData, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	
	// 解析存储的位置
	var data map[string]int
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return false
	}
	
	// 删除验证码
	global.Redis.Del(ctx, key)
	
	// 允许一定的误差（±5像素）
	targetX := data["x"]
	return x >= targetX-5 && x <= targetX+5
}

// GetSliderCaptchaTarget 获取滑动验证码目标位置（仅用于生成图片）
func (s *CaptchaService) GetSliderCaptchaTarget(captchaID string) (int, int, bool) {
	ctx := context.Background()
	key := "slider_captcha:" + captchaID
	
	jsonData, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return 0, 0, false
	}
	
	var data map[string]int
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return 0, 0, false
	}
	
	return data["x"], data["y"], true
}

// GenerateCaptcha 生成图形验证码
func (s *CaptchaService) GenerateCaptcha() (id, b64s string, err error) {
	captchaType := s.GetCaptchaType()
	var driver base64Captcha.Driver

	switch captchaType {
	case "math":
		// 算术验证码
		driver = base64Captcha.NewDriverMath(
			80,  // 高度
			240, // 宽度
			5,   // 干扰线数量
			base64Captcha.OptionShowHollowLine,
			nil, // 背景色
			nil, // 字体库
			nil, // 字体
		)
	case "string":
		// 字符串验证码
		driver = base64Captcha.NewDriverString(
			80,                             // 高度
			240,                            // 宽度
			5,                              // 干扰线数量
			base64Captcha.OptionShowHollowLine,
			5,                              // 验证码长度
			"ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789", // 可用字符
			nil,
			nil,
			nil,
		)
	default:
		// 数字验证码
		driver = base64Captcha.NewDriverDigit(
			80,  // 高度
			240, // 宽度
			5,   // 验证码长度
			0.7, // 最大倾斜角度
			80,  // 干扰点数量
		)
	}

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
	// 注册不再需要图形验证码
	return false
}

// GetLoginMaxRetry 获取最大重试次数
func (s *CaptchaService) GetLoginMaxRetry() int {
	config, err := Config.GetConfigByKey("login_max_retry")
	if err != nil {
		return 5
	}
	var val int
	for _, c := range config.Value {
		val = val*10 + int(c-'0')
	}
	if val <= 0 {
		return 5
	}
	return val
}

// GetLoginLockTime 获取锁定时间(分钟)
func (s *CaptchaService) GetLoginLockTime() int {
	config, err := Config.GetConfigByKey("login_lock_time")
	if err != nil {
		return 15
	}
	var val int
	for _, c := range config.Value {
		val = val*10 + int(c-'0')
	}
	if val <= 0 {
		return 15
	}
	return val
}

// CheckLoginLock 检查用户是否被锁定
func (s *CaptchaService) CheckLoginLock(username string) (bool, int) {
	ctx := context.Background()
	key := "login_lock:" + username
	ttl, err := global.Redis.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		return false, 0
	}
	return true, int(ttl.Minutes())
}

// IncrLoginRetry 增加登录失败次数
func (s *CaptchaService) IncrLoginRetry(username string) (int, bool) {
	ctx := context.Background()
	retryKey := "login_retry:" + username
	lockKey := "login_lock:" + username

	maxRetry := s.GetLoginMaxRetry()
	lockTime := s.GetLoginLockTime()

	// 增加重试次数
	count, _ := global.Redis.Incr(ctx, retryKey).Result()
	global.Redis.Expire(ctx, retryKey, time.Duration(lockTime)*time.Minute)

	// 检查是否达到最大重试次数
	if int(count) >= maxRetry {
		// 锁定账户
		global.Redis.Set(ctx, lockKey, "1", time.Duration(lockTime)*time.Minute)
		global.Redis.Del(ctx, retryKey)
		return int(count), true
	}

	return int(count), false
}

// ClearLoginRetry 清除登录重试次数
func (s *CaptchaService) ClearLoginRetry(username string) {
	ctx := context.Background()
	global.Redis.Del(ctx, "login_retry:"+username)
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
