package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"go-base-server/middleware"
	"go-base-server/model/request"
	"go-base-server/model/response"
	"go-base-server/service"
)

type {{.ModelName}}Api struct{}

var {{.ModelName}} = new({{.ModelName}}Api)

// Get{{.ModelName}}List 获取{{.Description}}列表
func (a *{{.ModelName}}Api) Get{{.ModelName}}List(c *gin.Context) {
	var req request.{{.ModelName}}ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

{{- if .DataIsolation}}
	userID := middleware.GetUserID(c)
	roleIDs := middleware.GetUserRoleIDs(c) // 角色ID从中间件上下文获取
	isAdmin := CheckIsAdmin(roleIDs, "{{.AdminRoleIds}}")
	list, total, err := service.{{.ModelName}}.Get{{.ModelName}}List(&req, userID, isAdmin)
{{- else}}
	list, total, err := service.{{.ModelName}}.Get{{.ModelName}}List(&req)
{{- end}}
	if err != nil {
		response.Fail(c, "获取列表失败")
		return
	}

	response.OkWithPage(c, list, total, req.Page, req.PageSize)
}

// Get{{.ModelName}} 获取{{.Description}}详情
func (a *{{.ModelName}}Api) Get{{.ModelName}}(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	data, err := service.{{.ModelName}}.Get{{.ModelName}}(uint(id))
	if err != nil {
		response.Fail(c, "获取详情失败")
		return
	}

	response.OkWithData(c, data)
}

// Create{{.ModelName}} 创建{{.Description}}
func (a *{{.ModelName}}Api) Create{{.ModelName}}(c *gin.Context) {
	var req request.Create{{.ModelName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := service.{{.ModelName}}.Create{{.ModelName}}(&req, userID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "创建成功")
}

// Update{{.ModelName}} 更新{{.Description}}
func (a *{{.ModelName}}Api) Update{{.ModelName}}(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	var req request.Update{{.ModelName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := service.{{.ModelName}}.Update{{.ModelName}}(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// Delete{{.ModelName}} 删除{{.Description}}
func (a *{{.ModelName}}Api) Delete{{.ModelName}}(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	if err := service.{{.ModelName}}.Delete{{.ModelName}}(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

// BatchDelete{{.ModelName}} 批量删除{{.Description}}
func (a *{{.ModelName}}Api) BatchDelete{{.ModelName}}(c *gin.Context) {
	var req request.BatchDelete{{.ModelName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	if err := service.{{.ModelName}}.BatchDelete{{.ModelName}}(req.Ids); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "批量删除成功")
}

// Get{{.ModelName}}Options 获取{{.Description}}选项列表
func (a *{{.ModelName}}Api) Get{{.ModelName}}Options(c *gin.Context) {
	displayField := c.DefaultQuery("display_field", "name")
	countTable := c.Query("count_table")
	countForeignKey := c.Query("count_field")
	excludeDeleted := c.Query("exclude_deleted") == "true"
	// 数据隔离：统计时按创建人过滤
	var countCreatedBy uint = 0
	if ccb := c.Query("count_created_by"); ccb != "" {
		if id, err := strconv.ParseUint(ccb, 10, 64); err == nil {
			countCreatedBy = uint(id)
		}
	}

{{- if .DataIsolation}}
	userID := middleware.GetUserID(c)
	roleIDs := middleware.GetUserRoleIDs(c)
	isAdmin := CheckIsAdmin(roleIDs, "{{.AdminRoleIds}}")
	list, err := service.{{.ModelName}}.Get{{.ModelName}}Options(displayField, countTable, countForeignKey, excludeDeleted, countCreatedBy, userID, isAdmin)
{{- else}}
	list, err := service.{{.ModelName}}.Get{{.ModelName}}Options(displayField, countTable, countForeignKey, excludeDeleted, countCreatedBy)
{{- end}}
	if err != nil {
		response.Fail(c, "获取选项列表失败")
		return
	}
	response.OkWithData(c, list)
}
{{- if .HasCreatedBy}}

// Get{{.ModelName}}CreatorOptions 获取创建人选项列表
func (a *{{.ModelName}}Api) Get{{.ModelName}}CreatorOptions(c *gin.Context) {
	list, err := service.{{.ModelName}}.Get{{.ModelName}}CreatorOptions()
	if err != nil {
		response.Fail(c, "获取创建人列表失败")
		return
	}
	response.OkWithData(c, list)
}
{{- end}}
{{- if .HasAudit}}

// Audit{{.ModelName}} 审批{{.Description}}
func (a *{{.ModelName}}Api) Audit{{.ModelName}}(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.Audit{{.ModelName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	if err := service.{{.ModelName}}.Audit{{.ModelName}}(uint(id), &req, userID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "审批成功")
}
{{- end}}
{{- if .GenerateFrontendApi}}

// GetFrontend{{.ModelName}}List 获取前台{{.Description}}列表（前台用户使用）
func (a *{{.ModelName}}Api) GetFrontend{{.ModelName}}List(c *gin.Context) {
	var req request.Frontend{{.ModelName}}ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	list, total, err := service.{{.ModelName}}.GetFrontend{{.ModelName}}List(&req)
	if err != nil {
		response.Fail(c, "获取列表失败")
		return
	}

	response.OkWithPage(c, list, total, req.Page, req.PageSize)
}

// GetFrontend{{.ModelName}} 获取前台{{.Description}}详情（前台用户使用）
func (a *{{.ModelName}}Api) GetFrontend{{.ModelName}}(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	data, err := service.{{.ModelName}}.GetFrontend{{.ModelName}}(uint(id))
	if err != nil {
		response.Fail(c, "数据不存在或未发布")
		return
	}
	response.OkWithData(c, data)
}
{{- end}}
{{- if .LinkToUser}}

// GetMy{{.ModelName}} 获取当前用户的{{.Description}}信息
func (a *{{.ModelName}}Api) GetMy{{.ModelName}}(c *gin.Context) {
	userID := middleware.GetUserID(c)
	data, err := service.{{.ModelName}}.GetMy{{.ModelName}}(userID)
	if err != nil {
		response.Fail(c, "获取信息失败")
		return
	}
	response.OkWithData(c, data)
}

// SaveMy{{.ModelName}} 保存当前用户的{{.Description}}信息
func (a *{{.ModelName}}Api) SaveMy{{.ModelName}}(c *gin.Context) {
	var req request.SaveMy{{.ModelName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := service.{{.ModelName}}.SaveMy{{.ModelName}}(userID, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "保存成功")
}
{{- end}}
