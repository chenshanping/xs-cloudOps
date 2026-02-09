package v1

import (
	"strconv"
	"strings"
)

// CheckIsAdmin 检查用户角色是否为管理员
// roleIDs: 用户的角色ID列表
// adminRoleIDsStr: 管理员角色ID字符串，逗号分隔
func CheckIsAdmin(roleIDs []uint, adminRoleIDsStr string) bool {
	if adminRoleIDsStr == "" {
		return false
	}
	adminRoleIDs := []uint{}
	for _, idStr := range strings.Split(adminRoleIDsStr, ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr != "" {
			if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
				adminRoleIDs = append(adminRoleIDs, uint(id))
			}
		}
	}

	for _, rid := range roleIDs {
		for _, ar := range adminRoleIDs {
			if rid == ar {
				return true
			}
		}
	}
	return false
}
