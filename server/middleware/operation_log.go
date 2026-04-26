package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"server/global"
	"server/model"
	"server/router/registry"
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
			// 检查是否为文件上传（multipart/form-data）
			contentType := c.Request.Header.Get("Content-Type")
			if strings.Contains(contentType, "multipart/form-data") {
				// 文件上传请求，记录文件信息而不是完整内容
				if err := c.Request.ParseMultipartForm(32 << 20); err == nil { // 32MB
					if c.Request.MultipartForm != nil && len(c.Request.MultipartForm.File) > 0 {
						var fileInfos []string
						for fieldName, files := range c.Request.MultipartForm.File {
							for _, file := range files {
								fileInfos = append(fileInfos, 
									fmt.Sprintf("%s: %s (%.2fKB)", fieldName, file.Filename, float64(file.Size)/1024))
							}
						}
						requestBody = "[文件上传] " + strings.Join(fileInfos, ", ")
					}
				}
			} else {
				// 普通请求，读取请求体
				bodyBytes, _ := io.ReadAll(c.Request.Body)
				requestBody = string(bodyBytes)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
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

		// 检查是否为文件下载响应（通过Content-Type判断）
		contentType := c.Writer.Header().Get("Content-Type")
		isFileDownload := strings.Contains(contentType, "application/octet-stream") ||
			strings.Contains(contentType, "application/vnd.openxmlformats-officedocument") ||
			strings.Contains(contentType, "application/vnd.ms-excel") ||
			strings.Contains(contentType, "application/pdf") ||
			strings.Contains(contentType, "image/") ||
			strings.Contains(contentType, "video/") ||
			strings.Contains(contentType, "audio/")

		var responseBody string
		var businessCode int

		if isFileDownload {
			// 文件下载响应，不记录响应体
			responseBody = "[文件下载]"
			businessCode = 200 // 文件下载默认认为成功
		} else {
			// 获取完整响应体用于解析业务码
			fullResponseBody := writer.body.String()
			businessCode = parseBusinessCode(fullResponseBody)

			// 限制响应体长度
			responseBody = fullResponseBody
			if len(responseBody) > 1000 {
				responseBody = responseBody[:1000] + "..."
			}
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
