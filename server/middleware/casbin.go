package middleware

import (
	"strings"

	"server/global"
	"server/model/response"
	rolesvc "server/service/role"

	"github.com/gin-gonic/gin"
)

// 白名单路径，所有登录用户都可访问
var casbinWhiteList = map[string]bool{
	"/api/v1/auth/userinfo": true,
	"/api/v1/auth/logout":   true,
	"/api/v1/auth/refresh":  true,
	"/api/v1/user/password": true,
	"/api/v1/user/profile":  true,
}

// 白名单路径后缀，匹配以这些后缀结尾的路径
var casbinWhiteSuffixes = []string{
	"/my", // 用户身份相关的 /my 接口，允许登录用户填写自己的信息
}

// 基于用户角色进行Casbin校验：遍历用户的所有角色，用角色编码作为sub进行Enforce
func CasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		// 白名单直接放行
		if casbinWhiteList[path] {
			c.Next()
			return
		}

		// 后缀白名单检查（支持 /xxx/my 类型的接口）
		for _, suffix := range casbinWhiteSuffixes {
			if strings.HasSuffix(path, suffix) {
				c.Next()
				return
			}
		}

		roleIDs := GetUserRoleIDs(c)
		hasSuperAdmin, err := rolesvc.Default.HasSuperAdminRoleIDs(roleIDs)
		if err != nil {
			global.Log.Errorf("查询显式超管角色失败: %v", err)
			response.Fail(c, "权限校验失败")
			c.Abort()
			return
		}
		if hasSuperAdmin {
			c.Next()
			return
		}

		// 从JWT解析角色编码，避免每次请求查询数据库
		roleCodes := GetUserRoleCodes(c)
		if len(roleCodes) == 0 {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		// 逐个角色校验
		for _, code := range roleCodes {
			ok, err := global.Enforcer.Enforce(code, path, method)
			if err != nil {
				global.Log.Errorf("Casbin权限校验失败: %v", err)
				response.Fail(c, "权限校验失败")
				c.Abort()
				return
			}
			if ok {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "无权限访问")
		c.Abort()
	}
}
