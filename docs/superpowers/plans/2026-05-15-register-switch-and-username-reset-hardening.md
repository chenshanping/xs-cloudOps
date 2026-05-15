# Register Switch And Username Reset Hardening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enforce the `enable_register` system config on the backend registration path and fix the captcha bypass in the public username-based password reset API.

**Architecture:** Keep the registration policy in `server/service/user` so all callers of `User.Register(...)` share the same backend guard. Fix the password reset bypass in the auth API handler by requiring captcha presence and successful verification before password reset continues.

**Tech Stack:** Go, Gin, GORM, sqlite test DB, existing service-layer and handler-layer tests

---

### Task 1: Enforce Register Switch In User Service

**Files:**
- Modify: `server/service/user/user.go`
- Modify: `server/service/user/user_test.go`

- [ ] **Step 1: Write the failing tests**

Add service tests that seed `sys_config.enable_register` and assert:

```go
func TestRegisterRejectsWhenRegisterDisabled(t *testing.T) {
	setupUserServiceTestDB(t)
	if err := global.DB.Create(&model.SysConfig{
		Key:   "enable_register",
		Value: "false",
	}).Error; err != nil {
		t.Fatalf("create enable_register config: %v", err)
	}

	err := Default.Register("closed-user", "123456", "closed-user@example.com")
	if err == nil || err.Error() != "系统已关闭注册" {
		t.Fatalf("expected register closed error, got %v", err)
	}

	var count int64
	if err := global.DB.Model(&model.SysUser{}).Where("username = ?", "closed-user").Count(&count).Error; err != nil {
		t.Fatalf("count user: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no user created, got %d", count)
	}
}

func TestRegisterAllowsWhenRegisterEnabled(t *testing.T) {
	setupUserServiceTestDB(t)
	if err := global.DB.Create(&model.SysConfig{
		Key:   "enable_register",
		Value: "true",
	}).Error; err != nil {
		t.Fatalf("create enable_register config: %v", err)
	}

	if err := Default.Register("open-user", "123456", "open-user@example.com"); err != nil {
		t.Fatalf("register error: %v", err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./service/user -run "TestRegisterRejectsWhenRegisterDisabled|TestRegisterAllowsWhenRegisterEnabled" -count=1`
Expected: FAIL because `Register(...)` currently ignores `enable_register`.

- [ ] **Step 3: Write minimal implementation**

Add a focused helper in `server/service/user/user.go`:

```go
func (s *UserService) IsRegisterEnabled() bool {
	config, err := configsvc.Default.GetConfigByKey("enable_register")
	if err != nil {
		return false
	}
	value := strings.TrimSpace(strings.ToLower(config.Value))
	return value == "1" || value == "true"
}
```

Call it at the top of `Register(...)`:

```go
if !s.IsRegisterEnabled() {
	return errors.New("系统已关闭注册")
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./service/user -run "TestRegisterRejectsWhenRegisterDisabled|TestRegisterAllowsWhenRegisterEnabled" -count=1`
Expected: PASS

- [ ] **Step 5: Verify neighboring service tests**

Run: `go test ./service/user -count=1`
Expected: PASS

### Task 2: Fix Username Reset Captcha Bypass

**Files:**
- Modify: `server/api/v1/auth.go`
- Create: `server/api/v1/auth_test.go`

- [ ] **Step 1: Write the failing tests**

Add handler tests that prove:

```go
func TestResetPasswordByUserNameRejectsWhenCaptchaVerificationFails(t *testing.T) {
	// request contains a non-empty captcha_id and captcha,
	// but VerifyCaptcha returns false.
	// Expected: handler must not call GetUserByUserName or ResetPassword
	// and must return a failure response.
}

func TestResetPasswordByUserNameAllowsWhenCaptchaVerificationSucceeds(t *testing.T) {
	// request contains valid captcha inputs,
	// VerifyCaptcha returns true,
	// user lookup and reset proceed successfully.
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./api/v1 -run "TestResetPasswordByUserNameRejectsWhenCaptchaVerificationFails|TestResetPasswordByUserNameAllowsWhenCaptchaVerificationSucceeds" -count=1`
Expected: FAIL because the current handler only verifies captcha when `captcha_id == ""`.

- [ ] **Step 3: Write minimal implementation**

Replace the current captcha branch in `ResetPasswordByUserName` with:

```go
if captchaId == "" || captchaCode == "" {
	response.BadRequest(c, "参数错误")
	return
}
if !service.Captcha.VerifyCaptcha(captchaId, captchaCode) {
	response.Fail(c, "验证码错误")
	return
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./api/v1 -run "TestResetPasswordByUserNameRejectsWhenCaptchaVerificationFails|TestResetPasswordByUserNameAllowsWhenCaptchaVerificationSucceeds" -count=1`
Expected: PASS

- [ ] **Step 5: Run broader backend verification**

Run: `go test ./service/user ./api/v1 -count=1`
Expected: PASS

- [ ] **Step 6: Run final compile verification**

Run: `go build ./...`
Expected: PASS
