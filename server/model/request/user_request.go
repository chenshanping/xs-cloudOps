package request

// 创建用户请求
type CreateUserRequest struct {
	Username     string `json:"username" binding:"required" comment:"用户名"`
	Password     string `json:"password" binding:"required,min=6" comment:"密码"`
	Nickname     string `json:"nickname" comment:"昵称"`
	Gender       int    `json:"gender" comment:"性别(0:未知,1:男,2:女)"`
	Email        string `json:"email" comment:"邮箱"`
	Phone        string `json:"phone" comment:"手机号"`
	Avatar       string `json:"avatar" comment:"头像地址"`
	AvatarFileID uint   `json:"avatar_file_id" comment:"头像文件ID"`
	Status       int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	DeptID       uint   `json:"dept_id" comment:"部门ID"`
	RoleIds      []uint `json:"role_ids" comment:"角色ID列表"`
}

// 更新用户请求
type UpdateUserRequest struct {
	Nickname     string `json:"nickname" comment:"昵称"`
	Gender       int    `json:"gender" comment:"性别(0:未知,1:男,2:女)"`
	Email        string `json:"email" comment:"邮箱"`
	Phone        string `json:"phone" comment:"手机号"`
	Avatar       string `json:"avatar" comment:"头像地址"`
	AvatarFileID uint   `json:"avatar_file_id" comment:"头像文件ID"`
	Status       int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	DeptID       uint   `json:"dept_id" comment:"部门ID"`
	RoleIds      []uint `json:"role_ids" comment:"角色ID列表"`
}

// 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" comment:"旧密码"`
	NewPassword string `json:"new_password" binding:"required,min=6" comment:"新密码"`
}

// 用户列表请求
type UserListRequest struct {
	PageRequest
	Username       string `json:"username" form:"username" comment:"用户名"`
	Status         *int   `json:"status" form:"status" comment:"状态"`
	Gender         *int   `json:"gender" form:"gender" comment:"性别"`
	RoleId         *int   `json:"role_id" form:"role_id" comment:"角色ID"`
	DeptId         *int   `json:"dept_id" form:"dept_id" comment:"部门ID"`
	UnassignedDept bool   `json:"unassigned_dept" form:"unassigned_dept" comment:"是否筛选未绑定部门用户"`
}

// 批量修改用户状态请求
type BatchUserStatusRequest struct {
	Ids    []uint `json:"ids" binding:"required" comment:"用户ID列表"`
	Status int    `json:"status" binding:"oneof=0 1" comment:"状态(0:禁用,1:启用)"`
}

// 批量重置密码请求
type BatchResetPasswordRequest struct {
	Ids []uint `json:"ids" binding:"required" comment:"用户ID列表"`
}

// 更新个人资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" comment:"昵称"`
	Email    string `json:"email" comment:"邮箱"`
	Phone    string `json:"phone" comment:"手机号"`
}
