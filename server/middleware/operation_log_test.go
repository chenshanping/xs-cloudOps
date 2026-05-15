package middleware

import (
	"strings"
	"testing"
)

func TestShouldRecordOperationLog(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		route    string
		summary  string
		expected bool
	}{
		{
			name:     "skip ordinary get list",
			method:   "GET",
			path:     "/api/v1/users",
			route:    "/api/v1/users",
			summary:  "用户列表",
			expected: false,
		},
		{
			name:     "keep detail get",
			method:   "GET",
			path:     "/api/v1/users/1",
			route:    "/api/v1/users/:id",
			summary:  "用户详情",
			expected: true,
		},
		{
			name:     "keep export get",
			method:   "GET",
			path:     "/api/v1/users/export",
			route:    "/api/v1/users/export",
			summary:  "导出用户",
			expected: true,
		},
		{
			name:     "keep config by key get",
			method:   "GET",
			path:     "/api/v1/configs/key/sys_name",
			route:    "/api/v1/configs/key/:key",
			summary:  "根据key获取配置",
			expected: true,
		},
		{
			name:     "keep write request",
			method:   "POST",
			path:     "/api/v1/users",
			route:    "/api/v1/users",
			summary:  "创建用户",
			expected: true,
		},
		{
			name:     "skip log module path",
			method:   "GET",
			path:     "/api/v1/logs/operation",
			route:    "/api/v1/logs/operation",
			summary:  "操作日志列表",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldRecordOperationLog(tt.method, tt.path, tt.route, tt.summary)
			if got != tt.expected {
				t.Fatalf("shouldRecordOperationLog(%q, %q, %q, %q) = %v, want %v", tt.method, tt.path, tt.route, tt.summary, got, tt.expected)
			}
		})
	}
}

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
