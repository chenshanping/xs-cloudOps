package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"go-base-server/global"
	"go-base-server/model"
	"go-base-server/router/registry"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 需要跳过记录的路径前缀
var skipLogPaths = []string{
	"/api/v1/logs/",
}

// shouldSkipLog 判断是否跳过记录日志
func shouldSkipLog(path string) bool {
	for _, prefix := range skipLogPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// getRouteInfo 获取路由元信息
func getRouteInfo(method, path string) (group, summary string) {
	routes := registry.GetAllRoutes()
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return route.Group, route.Summary
		}
	}
	return "", ""
}

// parseBusinessCode 从响应体解析业务状态码
func parseBusinessCode(responseBody string) int {
	var resp struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal([]byte(responseBody), &resp); err == nil {
		return resp.Code
	}
	return 0
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 跳过不需要记录的路径
		if shouldSkipLog(path) {
			c.Next()
			return
		}

		start := time.Now()

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 包装响应写入器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = writer

		c.Next()

		// 记录日志
		latency := time.Since(start).Milliseconds()

		// 获取完整响应体用于解析业务码
		fullResponseBody := writer.body.String()
		businessCode := parseBusinessCode(fullResponseBody)

		// 限制响应体长度
		responseBody := fullResponseBody
		if len(responseBody) > 1000 {
			responseBody = responseBody[:1000] + "..."
		}

		// 获取路由元信息（使用 FullPath 获取路由模板，如 /api/v1/roles/:id/menus）
		routePath := c.FullPath()
		group, summary := getRouteInfo(c.Request.Method, routePath)

		log := model.SysOperationLog{
			UserID:       GetUserID(c),
			Username:     GetUsername(c),
			IP:           c.ClientIP(),
			Method:       c.Request.Method,
			Path:         path,
			Group:        group,
			Summary:      summary,
			Request:      requestBody,
			Response:     responseBody,
			Status:       c.Writer.Status(),
			BusinessCode: businessCode,
			Latency:      latency,
			UserAgent:    c.Request.UserAgent(),
		}

		// 异步写入数据库
		go func() {
			if err := global.DB.Create(&log).Error; err != nil {
				global.Log.Errorf("记录操作日志失败: %v", err)
			}
		}()
	}
}
