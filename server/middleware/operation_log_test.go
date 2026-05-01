package middleware

import (
	"strings"
	"testing"
)

func TestSanitizeLogPayloadMasksSensitiveJSONFields(t *testing.T) {
	payload := `{"username":"demo","password":"123456","token":"secret-token","profile":{"secret_key":"sk-test","access_key":"ak-test","email_code":"8899"},"normal":"keep"}`

	sanitized := sanitizeLogPayload(payload)

	for _, value := range []string{"123456", "secret-token", "sk-test", "ak-test", "8899"} {
		if strings.Contains(sanitized, value) {
			t.Fatalf("sanitized payload still contains sensitive value %q: %s", value, sanitized)
		}
	}
	if !strings.Contains(sanitized, `"normal":"keep"`) {
		t.Fatalf("sanitized payload should preserve non-sensitive fields: %s", sanitized)
	}
}

func TestSanitizeLogPayloadLeavesNonJSONMarkersUntouched(t *testing.T) {
	payload := "[文件上传] file: demo.txt (12.00KB)"

	sanitized := sanitizeLogPayload(payload)

	if sanitized != payload {
		t.Fatalf("non-json payload changed = %s, want %s", sanitized, payload)
	}
}
