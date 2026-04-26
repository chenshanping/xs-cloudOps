package request

// 分页请求
type PageRequest struct {
	Page     int `json:"page" form:"page" comment:"页码"`
	PageSize int `json:"page_size" form:"page_size" comment:"每页数量"`
}

func (p *PageRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return (p.Page - 1) * p.PageSize
}

// 登录请求
type LoginRequest struct {
	Username  string `json:"username" binding:"required" comment:"用户名"`
	Password  string `json:"password" binding:"required" comment:"密码"`
	CaptchaID string `json:"captcha_id" comment:"验证码ID"`
	Captcha   string `json:"captcha" comment:"验证码"`
}

// 注册请求
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=20" comment:"用户名"`
	Password    string `json:"password" binding:"required,min=6" comment:"密码"`
	Email       string `json:"email" binding:"email" comment:"邮箱"`
	EmailCode   string `json:"email_code" comment:"邮箱验证码"`
	CaptchaID   string `json:"captcha_id" comment:"验证码ID"`
	CaptchaCode string `json:"captcha_code" comment:"验证码"`
}

// 发送邮箱验证码请求
type SendEmailCodeRequest struct {
	Email     string `json:"email" binding:"required,email" comment:"邮箱地址"`
	CaptchaID string `json:"captcha_id" comment:"验证码ID"`
	Captcha   string `json:"captcha" comment:"验证码"`
}

type TestEmailRequest struct {
	Email string `json:"email" binding:"required,email" comment:"接收测试邮件的邮箱地址"`
}

// 忘记密码请求
type ForgotPasswordRequest struct {
	Email     string `json:"email" binding:"required,email" comment:"邮箱地址"`
	CaptchaID string `json:"captcha_id" comment:"验证码ID"`
	Captcha   string `json:"captcha" comment:"验证码"`
}

// 重置密码请求
type ResetPasswordByTokenRequest struct {
	Token    string `json:"token" binding:"required" comment:"重置令牌"`
	Password string `json:"password" binding:"required,min=6" comment:"新密码"`
}

// 重置密码参数（通过用户名 + 图形验证码）
type ResetPasswordByUserNameRequest struct {
	UserName    string `json:"username" binding:"required" comment:"用户名"`
	NewPassword string `json:"new_password" binding:"required,min=6" comment:"新密码"`
	CaptchaId   string `json:"captcha_id" binding:"required" comment:"验证码ID"`
	Captcha     string `json:"captcha" binding:"required" comment:"验证码"`
}

// 重置密码参数（通过邮箱验证码）
type ResetPasswordByEmailRequest struct {
	Email       string `json:"email" binding:"required,email" comment:"邮箱地址"`
	EmailCode   string `json:"email_code" binding:"required" comment:"邮箱验证码"`
	NewPassword string `json:"new_password" binding:"required,min=6" comment:"新密码"`
}

// 创建用户请求
type CreateUserRequest struct {
	Username     string `json:"username" binding:"required" comment:"用户名"`
	Password     string `json:"password" binding:"required,min=6" comment:"密码"`
	Nickname     string `json:"nickname" comment:"昵称"`
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

// 重置密码请求
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6" comment:"新密码"`
}

// 用户列表请求
type UserListRequest struct {
	PageRequest
	Username       string `json:"username" form:"username" comment:"用户名"`
	Status         *int   `json:"status" form:"status" comment:"状态"`
	RoleId         *int   `json:"role_id" form:"role_id" comment:"角色ID"`
	DeptId         *int   `json:"dept_id" form:"dept_id" comment:"部门ID"`
	UnassignedDept bool   `json:"unassigned_dept" form:"unassigned_dept" comment:"是否筛选未绑定部门用户"`
}

// 批量修改用户状态请求
type BatchUserStatusRequest struct {
	Ids    []uint `json:"ids" binding:"required" comment:"用户ID列表"`
	Status int    `json:"status" binding:"oneof=0 1" comment:"状态(0:禁用,1:启用)"`
}

// 创建角色请求
type CreateRoleRequest struct {
	Name      string `json:"name" binding:"required" comment:"角色名称"`
	Code      string `json:"code" binding:"required" comment:"角色编码"`
	Sort      int    `json:"sort" comment:"排序"`
	Status    int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	DataScope int    `json:"data_scope" comment:"数据范围"`
	Remark    string `json:"remark" comment:"备注"`
	DeptIds   []uint `json:"dept_ids" comment:"自定义数据范围部门ID列表"`
}

// 更新角色请求
type UpdateRoleRequest struct {
	Name      string `json:"name" comment:"角色名称"`
	Code      string `json:"code" comment:"角色编码"`
	Sort      int    `json:"sort" comment:"排序"`
	Status    int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	DataScope int    `json:"data_scope" comment:"数据范围"`
	Remark    string `json:"remark" comment:"备注"`
	DeptIds   []uint `json:"dept_ids" comment:"自定义数据范围部门ID列表"`
}

// 创建部门请求
type CreateDeptRequest struct {
	ParentID uint   `json:"parent_id" comment:"父部门ID"`
	Name     string `json:"name" binding:"required" comment:"部门名称"`
	Sort     int    `json:"sort" comment:"排序"`
	Status   int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	Remark   string `json:"remark" comment:"备注"`
}

// 更新部门请求
type UpdateDeptRequest struct {
	ParentID uint   `json:"parent_id" comment:"父部门ID"`
	Name     string `json:"name" binding:"required" comment:"部门名称"`
	Sort     int    `json:"sort" comment:"排序"`
	Status   int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	Remark   string `json:"remark" comment:"备注"`
}

// 分配菜单请求
type AssignMenusRequest struct {
	MenuIds []uint `json:"menu_ids" comment:"菜单ID列表"`
}

// 分配API请求
type AssignApisRequest struct {
	ApiIds []uint `json:"api_ids" comment:"API ID列表"`
}

// 创建菜单请求
type CreateMenuRequest struct {
	ParentID   uint   `json:"parent_id" comment:"父菜单ID"`
	Name       string `json:"name" binding:"required" comment:"菜单名称"`
	Path       string `json:"path" comment:"路由路径"`
	Component  string `json:"component" comment:"组件路径"`
	Icon       string `json:"icon" comment:"图标"`
	Sort       int    `json:"sort" comment:"排序"`
	Type       int    `json:"type" binding:"required,oneof=1 2 3" comment:"类型(1:目录,2:菜单,3:按钮)"`
	Permission string `json:"permission" comment:"权限标识"`
	Status     int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	Hidden     int    `json:"hidden" comment:"是否隐藏(0:显示,1:隐藏)"`
}

// 更新菜单请求
type UpdateMenuRequest struct {
	ParentID   uint   `json:"parent_id" comment:"父菜单ID"`
	Name       string `json:"name" comment:"菜单名称"`
	Path       string `json:"path" comment:"路由路径"`
	Component  string `json:"component" comment:"组件路径"`
	Icon       string `json:"icon" comment:"图标"`
	Sort       int    `json:"sort" comment:"排序"`
	Type       int    `json:"type" comment:"类型(1:目录,2:菜单,3:按钮)"`
	Permission string `json:"permission" comment:"权限标识"`
	Status     int    `json:"status" comment:"状态(0:禁用,1:启用)"`
	Hidden     int    `json:"hidden" comment:"是否隐藏(0:显示,1:隐藏)"`
}

// 创建API请求
type CreateApiRequest struct {
	Path        string `json:"path" binding:"required" comment:"API路径"`
	Method      string `json:"method" binding:"required" comment:"请求方法"`
	Group       string `json:"group" comment:"API分组"`
	Description string `json:"description" comment:"描述"`
}

// 更新API请求
type UpdateApiRequest struct {
	Path        string `json:"path" comment:"API路径"`
	Method      string `json:"method" comment:"请求方法"`
	Group       string `json:"group" comment:"API分组"`
	Description string `json:"description" comment:"描述"`
}

// API列表请求
type ApiListRequest struct {
	PageRequest
	Path   string `json:"path" form:"path" comment:"API路径"`
	Method string `json:"method" form:"method" comment:"请求方法"`
	Group  string `json:"group" form:"group" comment:"API分组"`
}

// 更新个人资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" comment:"昵称"`
	Email    string `json:"email" comment:"邮箱"`
	Phone    string `json:"phone" comment:"手机号"`
}

// 日志列表请求
type LogListRequest struct {
	PageRequest
	Username     string `json:"username" form:"username" comment:"用户名"`
	Method       string `json:"method" form:"method" comment:"请求方法"`
	Path         string `json:"path" form:"path" comment:"请求路径"`
	Group        string `json:"group" form:"group" comment:"路由分组"`
	Summary      string `json:"summary" form:"summary" comment:"路由描述"`
	Status       *int   `json:"status" form:"status" comment:"HTTP状态码"`
	BusinessCode *int   `json:"business_code" form:"business_code" comment:"业务状态码"`
	StartTime    string `json:"start_time" form:"start_time" comment:"开始时间"`
	EndTime      string `json:"end_time" form:"end_time" comment:"结束时间"`
	SortField    string `json:"sort_field" form:"sort_field" comment:"排序字段"`
	SortOrder    string `json:"sort_order" form:"sort_order" comment:"排序方式(ascend/descend)"`
}

// 慢查询日志列表请求
type SlowLogListRequest struct {
	PageRequest
	SQL        string  `json:"sql" form:"sql" comment:"SQL语句"`
	MinLatency float64 `json:"min_latency" form:"min_latency" comment:"最小耗时(毫秒)"`
}
