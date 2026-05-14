package router

import (
	"runtime/debug"
	"testing"

	"github.com/gin-gonic/gin"

	"server/config"
	"server/global"
)

func TestInitRouterDoesNotPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousConfig := global.Config
	global.Config = &config.Config{Server: config.Server{Host: "localhost"}}
	t.Cleanup(func() {
		global.Config = previousConfig
	})
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("InitRouter panicked: %v\n%s", recovered, string(debug.Stack()))
		}
	}()
	_ = InitRouter()
}
