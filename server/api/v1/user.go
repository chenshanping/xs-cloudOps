package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"server/global"
	"server/middleware"
	"server/model/request"
	"server/model/response"
	"server/service"
)

type UserApi struct{}

var User = new(UserApi)

// 获取用户选项列表（轻量级，用于下拉选择）
func (a *UserApi) GetUserOptions(c *gin.Context) {
	operatorID := middleware.GetUserID(c)
	list, err := service.User.GetUserOptions(operatorID)
	if err != nil {
		response.Fail(c, "获取用户选项失败")
		return
	}
	response.OkWithData(c, list)
}

// 获取用户列表
func (a *UserApi) GetUserList(c *gin.Context) {
	var req request.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	users, total, err := service.User.GetUserList(operatorID, &req)
	if err != nil {
		response.Fail(c, "获取用户列表失败")
		return
	}

	response.OkWithPage(c, users, total, req.Page, req.PageSize)
}

// 获取用户详情
func (a *UserApi) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	user, err := service.User.GetManagedUserInfo(operatorID, uint(id))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, user)
}

// 创建用户
func (a *UserApi) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	if err := service.User.CreateUser(operatorID, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "创建成功")
}

// 更新用户
func (a *UserApi) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	if err := service.User.UpdateUser(operatorID, uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// 删除用户
func (a *UserApi) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	if err := service.User.DeleteUser(operatorID, uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

// 批量删除用户
func (a *UserApi) BatchDeleteUsers(c *gin.Context) {
	var req struct {
		Ids []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if len(req.Ids) == 0 {
		response.BadRequest(c, "请选择要删除的用户")
		return
	}

	operatorID := middleware.GetUserID(c)
	successCount, failedMsgs := service.User.BatchDeleteUsers(operatorID, req.Ids)

	if len(failedMsgs) == 0 {
		response.OkWithMessage(c, "batch_delete_success")
	} else if successCount > 0 {
		response.OkWithData(c, gin.H{
			"success_count": successCount,
			"failed_count":  len(failedMsgs),
			"failed_msgs":   failedMsgs,
		})
	} else {
		response.Fail(c, "删除失败")
	}
}

// 修改用户状态
func (a *UserApi) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	if err := service.User.UpdateUserStatus(operatorID, uint(id), req.Status); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "修改成功")
}

// 批量修改用户状态
func (a *UserApi) BatchUpdateUserStatus(c *gin.Context) {
	var req request.BatchUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	if err := service.User.BatchUpdateUserStatus(operatorID, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "修改成功")
}

// 重置密码
func (a *UserApi) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	if err := service.User.ResetManagedUserPassword(operatorID, uint(id), req.Password); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "重置成功")
}

// 修改密码
func (a *UserApi) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.User.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "修改成功")
}

// 获取个人资料
func (a *UserApi) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := service.User.GetUserInfo(userID)
	if err != nil {
		response.Fail(c, "获取个人资料失败")
		return
	}
	response.OkWithData(c, user)
}

// 更新个人资料
func (a *UserApi) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.User.UpdateProfile(userID, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// 更新头像（绑定文件表ID）
func (a *UserApi) UpdateAvatar(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		FileID uint `json:"file_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.User.UpdateAvatar(userID, req.FileID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// 强制用户下线
func (a *UserApi) ForceOffline(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 不能强制下线自己
	currentUserID := middleware.GetUserID(c)
	if err := service.User.ForceOffline(currentUserID, uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "已强制该用户下线")
}

// 获取当前用户所有身份信息
func (a *UserApi) GetUserProfiles(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取用户角色编码列表
	var userRoles []string
	user, err := service.User.GetUserInfo(userID)
	if err == nil && user != nil {
		for _, role := range user.Roles {
			userRoles = append(userRoles, role.Code)
		}
	}

	profiles := global.Profiles.GetUserProfiles(userID, userRoles)
	response.OkWithData(c, profiles)
}

// 获取系统已注册的所有身份类型
func (a *UserApi) GetRegisteredProfiles(c *gin.Context) {
	profiles := global.Profiles.GetRegisteredProfiles()
	response.OkWithData(c, profiles)
}

// 获取指定用户的身份信息（管理员使用）
func (a *UserApi) GetUserProfilesById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 获取用户角色编码列表
	var userRoles []string
	operatorID := middleware.GetUserID(c)
	user, err := service.User.GetManagedUserInfo(operatorID, uint(id))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	for _, role := range user.Roles {
		userRoles = append(userRoles, role.Code)
	}

	profiles := global.Profiles.GetUserProfiles(uint(id), userRoles)
	response.OkWithData(c, profiles)
}
