package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"server/middleware"
	"server/model/request"
	"server/model/response"
	"server/service"
)

type DeptApi struct{}

var Dept = new(DeptApi)

func (a *DeptApi) GetDeptTree(c *gin.Context) {
	tree, err := service.Dept.GetDeptTree()
	if err != nil {
		response.Fail(c, "获取部门树失败")
		return
	}
	response.OkWithData(c, tree)
}

func (a *DeptApi) GetManageableDeptTree(c *gin.Context) {
	operatorID := middleware.GetUserID(c)
	tree, unassignedCount, err := service.Dept.GetManageableDeptTree(operatorID)
	if err != nil {
		response.Fail(c, "获取可管理部门树失败")
		return
	}
	response.OkWithData(c, gin.H{
		"tree":                  tree,
		"unassigned_user_count": unassignedCount,
	})
}

func (a *DeptApi) GetDept(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	dept, err := service.Dept.GetDept(uint(id))
	if err != nil {
		response.Fail(c, "获取部门详情失败")
		return
	}
	response.OkWithData(c, dept)
}

func (a *DeptApi) CreateDept(c *gin.Context) {
	var req request.CreateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Dept.CreateDept(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "创建成功")
}

func (a *DeptApi) UpdateDept(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Dept.UpdateDept(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "更新成功")
}

func (a *DeptApi) DeleteDept(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Dept.DeleteDept(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "删除成功")
}
