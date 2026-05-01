package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"server/model/request"
	"server/model/response"
	"server/service"
	"server/service/core"
)

type RoleApi struct{}

var Role = new(RoleApi)

// 获取角色列表
func (a *RoleApi) GetRoleList(c *gin.Context) {
	roles, err := service.Role.GetRoleList()
	if err != nil {
		response.Fail(c, "获取角色列表失败")
		return
	}
	response.OkWithData(c, roles)
}

// 获取角色详情
func (a *RoleApi) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	role, err := service.Role.GetRole(uint(id))
	if err != nil {
		response.Fail(c, "获取角色信息失败")
		return
	}

	response.OkWithData(c, role)
}

// 获取数据权限资源列表
func (a *RoleApi) GetDataScopeResources(c *gin.Context) {
	response.OkWithData(c, core.SupportedDataScopeResources())
}

// 创建角色
func (a *RoleApi) CreateRole(c *gin.Context) {
	var req request.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Role.CreateRole(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "创建成功")
}

// 更新角色
func (a *RoleApi) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Role.UpdateRole(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// 删除角色
func (a *RoleApi) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Role.DeleteRole(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

// 分配菜单权限
func (a *RoleApi) AssignMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.AssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Role.AssignMenus(uint(id), req.MenuIds); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "分配成功")
}

// 分配API权限
func (a *RoleApi) AssignApis(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.AssignApisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Role.AssignApis(uint(id), req.ApiIds); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "分配成功")
}

// 分配业务功能数据权限
func (a *RoleApi) AssignDataScopes(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.AssignRoleDataScopesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Role.AssignDataScopes(uint(id), req.Scopes); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "分配成功")
}

// 统一保存角色权限
func (a *RoleApi) SavePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.SaveRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Role.SavePermissions(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "保存成功")
}
