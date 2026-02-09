package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"go-base-server/model/request"
	"go-base-server/model/response"
	"go-base-server/service"
)

type MenuApi struct{}

var Menu = new(MenuApi)

// 获取菜单列表(树形)
func (a *MenuApi) GetMenuTree(c *gin.Context) {
	menus, err := service.Menu.GetMenuTree()
	if err != nil {
		response.Fail(c, "获取菜单列表失败")
		return
	}
	response.OkWithData(c, menus)
}

// 获取菜单详情
func (a *MenuApi) GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	menu, err := service.Menu.GetMenu(uint(id))
	if err != nil {
		response.Fail(c, "获取菜单信息失败")
		return
	}

	response.OkWithData(c, menu)
}

// 创建菜单
func (a *MenuApi) CreateMenu(c *gin.Context) {
	var req request.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Menu.CreateMenu(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "创建成功")
}

// 更新菜单
func (a *MenuApi) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Menu.UpdateMenu(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// 删除菜单
func (a *MenuApi) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Menu.DeleteMenu(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

// 获取用户菜单(根据角色过滤)
func (a *MenuApi) GetUserMenus(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Fail(c, "用户ID为Nil")
		return
	}

	menus, err := service.Menu.GetUserMenus(userID)
	if err != nil {
		response.Fail(c, "获取菜单失败")
		return
	}

	response.OkWithData(c, menus)
}
