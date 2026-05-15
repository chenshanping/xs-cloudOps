package response

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"

	"server/model/request"
)

func validateWithBindingTag(t *testing.T, sample any) error {
	t.Helper()

	v := validator.New()
	v.SetTagName("binding")
	return v.Struct(sample)
}

func TestBindErrorMessageForRequiredField(t *testing.T) {
	err := validateWithBindingTag(t, request.LoginRequest{})
	if err == nil {
		t.Fatal("expected validation error")
	}

	if got := BindErrorMessage(err, request.LoginRequest{}); got != "用户名不能为空" {
		t.Fatalf("expected 用户名不能为空, got %q", got)
	}
}

func TestBindErrorMessageForMinLength(t *testing.T) {
	err := validateWithBindingTag(t, request.ResetPasswordByUserNameRequest{
		UserName:    "alice",
		NewPassword: "123",
		CaptchaId:   "captcha-id",
		Captcha:     "captcha",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	if got := BindErrorMessage(err, request.ResetPasswordByUserNameRequest{}); got != "新密码至少 6 位" {
		t.Fatalf("expected 新密码至少 6 位, got %q", got)
	}
}

func TestBindErrorMessageForEmail(t *testing.T) {
	err := validateWithBindingTag(t, request.SendEmailCodeRequest{
		Email: "bad-email",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	if got := BindErrorMessage(err, request.SendEmailCodeRequest{}); got != "邮箱地址格式不正确" {
		t.Fatalf("expected 邮箱地址格式不正确, got %q", got)
	}
}

func TestBindErrorMessageFallback(t *testing.T) {
	if got := BindErrorMessage(errors.New("random"), nil); got != "参数错误" {
		t.Fatalf("expected 参数错误, got %q", got)
	}
}
