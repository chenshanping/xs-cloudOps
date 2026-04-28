package captcha

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
	"server/service/configsvc"
)

type CaptchaService struct{}

var Default = &CaptchaService{}

type redisStore struct {
	expiration time.Duration
}

func newRedisStore() *redisStore {
	return &redisStore{
		expiration: 5 * time.Minute,
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

type SliderCaptchaData struct {
	CaptchaID string `json:"captcha_id"`
	BgWidth   int    `json:"bg_width"`
	BgHeight  int    `json:"bg_height"`
	SliderY   int    `json:"slider_y"`
}

func (s *CaptchaService) GetCaptchaType() string {
	config, err := configsvc.Default.GetConfigByKey("login_captcha_type")
	if err != nil {
		return "digit"
	}
	return config.Value
}

func (s *CaptchaService) GetSliderCaptchaBg() string {
	config, err := configsvc.Default.GetConfigByKey("slider_captcha_bg")
	if err != nil {
		return ""
	}
	return config.Value
}

func (s *CaptchaService) GenerateSliderCaptcha() (*SliderCaptchaData, error) {
	captchaID := uuid.New().String()

	bgWidth := 280
	bgHeight := 160
	sliderSize := 40

	mathRand.Seed(time.Now().UnixNano())
	xPos := mathRand.Intn(bgWidth-sliderSize*2-20) + sliderSize + 10
	yPos := mathRand.Intn(bgHeight-sliderSize-20) + 10

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

func (s *CaptchaService) VerifySliderCaptcha(captchaID string, x int) bool {
	ctx := context.Background()
	key := "slider_captcha:" + captchaID

	jsonData, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	var data map[string]int
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return false
	}

	global.Redis.Del(ctx, key)

	targetX := data["x"]
	return x >= targetX-5 && x <= targetX+5
}

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

func (s *CaptchaService) GenerateCaptcha() (id, b64s string, err error) {
	captchaType := s.GetCaptchaType()
	var driver base64Captcha.Driver

	switch captchaType {
	case "math":
		driver = base64Captcha.NewDriverMath(
			80,
			240,
			5,
			base64Captcha.OptionShowHollowLine,
			nil,
			nil,
			nil,
		)
	case "string":
		driver = base64Captcha.NewDriverString(
			80,
			240,
			5,
			base64Captcha.OptionShowHollowLine,
			5,
			"ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789",
			nil,
			nil,
			nil,
		)
	default:
		driver = base64Captcha.NewDriverDigit(
			80,
			240,
			5,
			0.7,
			80,
		)
	}

	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err = captcha.Generate()
	return
}

func (s *CaptchaService) VerifyCaptcha(id, code string) bool {
	if id == "" || code == "" {
		return false
	}
	return store.Verify(id, code, true)
}

func (s *CaptchaService) IsLoginCaptchaEnabled() bool {
	config, err := configsvc.Default.GetConfigByKey("login_captcha_enabled")
	if err != nil {
		return false
	}
	return config.Value == "1" || config.Value == "true"
}

func (s *CaptchaService) IsRegisterCaptchaEnabled() bool {
	return false
}

func (s *CaptchaService) GetLoginMaxRetry() int {
	config, err := configsvc.Default.GetConfigByKey("login_max_retry")
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

func (s *CaptchaService) GetLoginLockTime() int {
	config, err := configsvc.Default.GetConfigByKey("login_lock_time")
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

func (s *CaptchaService) CheckLoginLock(username string) (bool, int) {
	ctx := context.Background()
	key := "login_lock:" + username
	ttl, err := global.Redis.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		return false, 0
	}
	return true, int(ttl.Minutes())
}

func (s *CaptchaService) IncrLoginRetry(username string) (int, bool) {
	ctx := context.Background()
	retryKey := "login_retry:" + username
	lockKey := "login_lock:" + username

	maxRetry := s.GetLoginMaxRetry()
	lockTime := s.GetLoginLockTime()

	count, _ := global.Redis.Incr(ctx, retryKey).Result()
	global.Redis.Expire(ctx, retryKey, time.Duration(lockTime)*time.Minute)

	if int(count) >= maxRetry {
		global.Redis.Set(ctx, lockKey, "1", time.Duration(lockTime)*time.Minute)
		global.Redis.Del(ctx, retryKey)
		return int(count), true
	}

	return int(count), false
}

func (s *CaptchaService) ClearLoginRetry(username string) {
	ctx := context.Background()
	global.Redis.Del(ctx, "login_retry:"+username)
}

func (s *CaptchaService) IsRegisterEmailVerifyEnabled() bool {
	config, err := configsvc.Default.GetConfigByKey("register_email_verify")
	if err != nil {
		return false
	}
	return config.Value == "1" || config.Value == "true"
}

func (s *CaptchaService) GenerateEmailCode(email string) (string, error) {
	code := generateRandomCode(6)
	ctx := context.Background()
	key := "email_code:" + email
	err := global.Redis.Set(ctx, key, code, 10*time.Minute).Err()
	return code, err
}

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

func generateRandomCode(length int) string {
	const digits = "0123456789"
	code := make([]byte, length)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code[i] = digits[n.Int64()]
	}
	return string(code)
}
