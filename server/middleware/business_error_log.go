package middleware

import (
	"server/global"
	"server/model/response"
	"time"

	"github.com/gin-gonic/gin"
)

func BusinessErrorLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		value, exists := c.Get(response.ContextResponseCodeKey)
		if !exists {
			return
		}

		businessCode, ok := value.(int)
		if !ok || businessCode == response.SUCCESS {
			return
		}

		message, _ := c.Get(response.ContextResponseMessageKey)
		messageText, _ := message.(string)

		errorText := ""
		if len(c.Errors) > 0 {
			errorText = c.Errors.String()
		}
		if errorText == "" {
			errorText = messageText
		}

		global.Log.Warnw("业务请求失败",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"route", c.FullPath(),
			"status", c.Writer.Status(),
			"business_code", businessCode,
			"message", messageText,
			"error", errorText,
			"latency_ms", time.Since(start).Milliseconds(),
			"user_id", GetUserID(c),
			"username", GetUsername(c),
			"ip", c.ClientIP(),
		)
	}
}
