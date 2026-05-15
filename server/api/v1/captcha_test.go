package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"server/model/response"
)

func TestGetSliderCaptchaRejected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/captcha/slider", nil)

	CaptchaAPI.GetSliderCaptcha(c)

	var resp response.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Code != response.BAD_REQUEST || resp.Message != "当前版本暂不支持滑动验证码" {
		t.Fatalf("expected slider rejection, got code=%d message=%q", resp.Code, resp.Message)
	}
}
