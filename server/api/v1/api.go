package v1

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	"server/model"
	"server/model/request"
	"server/model/response"
	"server/router/registry"
	"server/service"
)

type ApiApi struct{}

var Api = new(ApiApi)

// 获取API列表
func (a *ApiApi) GetApiList(c *gin.Context) {
	var req request.ApiListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	apis, total, err := service.Api.GetApiList(&req)
	if err != nil {
		response.Fail(c, "获取API列表失败")
		return
	}

	response.OkWithPage(c, apis, total, req.Page, req.PageSize)
}

// 获取全部API
func (a *ApiApi) GetAllApis(c *gin.Context) {
	apis, err := service.Api.GetAllApis()
	if err != nil {
		response.Fail(c, "获取API列表失败")
		return
	}
	response.OkWithData(c, apis)
}

// 获取API详情
func (a *ApiApi) GetApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	api, err := service.Api.GetApi(uint(id))
	if err != nil {
		response.Fail(c, "获取API信息失败")
		return
	}

	response.OkWithData(c, api)
}

// 创建API
func (a *ApiApi) CreateApi(c *gin.Context) {
	var req request.CreateApiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Api.CreateApi(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "创建成功")
}

// 更新API
func (a *ApiApi) UpdateApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateApiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Api.UpdateApi(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// 删除API
func (a *ApiApi) DeleteApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Api.DeleteApi(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

// 获取API分组
func (a *ApiApi) GetApiGroups(c *gin.Context) {
	groups, err := service.Api.GetApiGroups()
	if err != nil {
		response.Fail(c, "获取API分组失败")
		return
	}
	response.OkWithData(c, groups)
}

// 同步API路由到数据库
func (a *ApiApi) SyncApis(c *gin.Context) {
	// 获取所有注册的路由
	routes := registry.GetAllRoutes()

	var apis []model.SysApi
	for _, route := range routes {
		// 解析请求参数
		var requestParams string
		if route.Request != nil {
			// 根据 HTTP 方法决定参数位置
			paramIn := "body"
			if route.Method == "GET" || route.Method == "DELETE" {
				paramIn = "query"
			}
			fields := registry.ParseStructFields(route.Request, paramIn)
			if len(fields) > 0 {
				fieldsJSON, _ := json.Marshal(fields)
				requestParams = string(fieldsJSON)
			}
		}

		// 解析响应参数
		var responseParams string
		if route.Response != nil {
			fields := registry.ParseStructFields(route.Response, "body")
			if len(fields) > 0 {
				fieldsJSON, _ := json.Marshal(fields)
				responseParams = string(fieldsJSON)
			}
		}

		apis = append(apis, model.SysApi{
			Path:           route.Path,
			Method:         route.Method,
			Group:          route.Group,
			Description:    route.Summary,
			RequestParams:  requestParams,
			ResponseParams: responseParams,
			NeedAuth:       route.NeedAuth,
		})
	}

	added, updated, deleted, err := service.Api.SyncApis(apis)
	if err != nil {
		response.Fail(c, "同步API失败: "+err.Error())
		return
	}

	response.OkWithData(c, map[string]interface{}{
		"added":   added,
		"updated": updated,
		"deleted": deleted,
		"total":   len(apis),
	})
}
