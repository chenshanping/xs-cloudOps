package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"go-base-server/model/response"
	"go-base-server/utils"
)

const (
	ContextUserIDKey    = "user_id"
	ContextUsernameKey  = "username"
	ContextRoleIDsKey   = "role_ids"
	ContextRoleCodesKey = "role_codes"
	ContextTokenKey     = "token"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c, "未登录或Token已过期")
			c.Abort()
			return
		}

		// 去掉Bearer前缀
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, "Token无效或已过期")
			c.Abort()
			return
		}

		// 验证Token在Redis中的状态（黑名单检查）
		if err := utils.ValidateTokenInRedis(token, claims.UserID); err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextRoleIDsKey, claims.RoleIDs)
		c.Set(ContextRoleCodesKey, claims.RoleCodes)
		c.Set(ContextTokenKey, token)
		c.Next()
	}
}

// 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get(ContextUserIDKey); exists {
		return userID.(uint)
	}
	return 0
}

// 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get(ContextUsernameKey); exists {
		return username.(string)
	}
	return ""
}

// 从上下文获取用户角色ID列表（直接从 JWT token 中解析，无需查询数据库）
func GetUserRoleIDs(c *gin.Context) []uint {
	if roleIDs, exists := c.Get(ContextRoleIDsKey); exists {
		if ids, ok := roleIDs.([]uint); ok {
			return ids
		}
	}
	return []uint{}
}

// 从上下文获取用户角色编码列表（用于Casbin权限校验）
func GetUserRoleCodes(c *gin.Context) []string {
	if roleCodes, exists := c.Get(ContextRoleCodesKey); exists {
		if codes, ok := roleCodes.([]string); ok {
			return codes
		}
	}
	return []string{}
}
