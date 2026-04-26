package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/model/request"
	"server/router/registry"
)

func init() {
	RegisterModule(&UserModule{})
}

type UserModule struct{}

func (m *UserModule) Name() string {
	return "用户管理"
}

func (m *UserModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *UserModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	// 个人中心
	R(rg, "GET", "/user/profile", "个人中心", "获取个人资料", v1.User.GetProfile, registry.WithAuth())
	R(rg, "PUT", "/user/profile", "个人中心", "更新个人资料", v1.User.UpdateProfile,
		registry.WithAuth(), registry.WithRequest(request.UpdateProfileRequest{}))
	R(rg, "PUT", "/user/avatar", "个人中心", "更新头像", v1.User.UpdateAvatar, registry.WithAuth())
	R(rg, "PUT", "/user/password", "个人中心", "修改密码", v1.User.ChangePassword,
		registry.WithAuth(), registry.WithRequest(request.ChangePasswordRequest{}))
	R(rg, "GET", "/user/profiles", "个人中心", "获取用户所有身份", v1.User.GetUserProfiles, registry.WithAuth())
	R(rg, "GET", "/user/profiles/types", "个人中心", "获取所有身份类型", v1.User.GetRegisteredProfiles, registry.WithAuth())

	// 用户管理
	R(rg, "GET", "/users", m.Name(), "用户列表", v1.User.GetUserList,
		registry.WithAuth(), registry.WithRequest(request.UserListRequest{}))
	R(rg, "GET", "/users/options", m.Name(), "用户选项", v1.User.GetUserOptions, registry.WithAuth())
	R(rg, "GET", "/users/:id", m.Name(), "用户详情", v1.User.GetUser, registry.WithAuth())
	R(rg, "POST", "/users", m.Name(), "创建用户", v1.User.CreateUser,
		registry.WithAuth(), registry.WithRequest(request.CreateUserRequest{}))
	R(rg, "PUT", "/users/:id", m.Name(), "更新用户", v1.User.UpdateUser,
		registry.WithAuth(), registry.WithRequest(request.UpdateUserRequest{}))
	R(rg, "DELETE", "/users/:id", m.Name(), "删除用户", v1.User.DeleteUser, registry.WithAuth())
	R(rg, "DELETE", "/users/batch", m.Name(), "批量删除用户", v1.User.BatchDeleteUsers, registry.WithAuth())
	R(rg, "PUT", "/users/:id/status", m.Name(), "修改用户状态", v1.User.UpdateUserStatus, registry.WithAuth())
	R(rg, "PUT", "/users/:id/password", m.Name(), "重置密码", v1.User.ResetPassword,
		registry.WithAuth(), registry.WithRequest(request.ResetPasswordRequest{}))
	R(rg, "POST", "/users/:id/offline", m.Name(), "强制下线", v1.User.ForceOffline, registry.WithAuth())
	R(rg, "GET", "/users/:id/profiles", m.Name(), "用户身份", v1.User.GetUserProfilesById, registry.WithAuth())
}
