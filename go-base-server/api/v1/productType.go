package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"go-base-server/middleware"
	"go-base-server/model/request"
	"go-base-server/model/response"
	"go-base-server/service"
)

type ProductTypeApi struct{}

var ProductType = new(ProductTypeApi)

// GetProductTypeList 获取产品类型列表
func (a *ProductTypeApi) GetProductTypeList(c *gin.Context) {
	var req request.ProductTypeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	list, total, err := service.ProductType.GetProductTypeList(&req)
	if err != nil {
		response.Fail(c, "获取列表失败")
		return
	}

	response.OkWithPage(c, list, total, req.Page, req.PageSize)
}

// GetProductType 获取产品类型详情
func (a *ProductTypeApi) GetProductType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	data, err := service.ProductType.GetProductType(uint(id))
	if err != nil {
		response.Fail(c, "获取详情失败")
		return
	}

	response.OkWithData(c, data)
}

// CreateProductType 创建产品类型
func (a *ProductTypeApi) CreateProductType(c *gin.Context) {
	var req request.CreateProductTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := service.ProductType.CreateProductType(&req, userID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "创建成功")
}

// UpdateProductType 更新产品类型
func (a *ProductTypeApi) UpdateProductType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	var req request.UpdateProductTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := service.ProductType.UpdateProductType(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// DeleteProductType 删除产品类型
func (a *ProductTypeApi) DeleteProductType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	if err := service.ProductType.DeleteProductType(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

// BatchDeleteProductType 批量删除产品类型
func (a *ProductTypeApi) BatchDeleteProductType(c *gin.Context) {
	var req request.BatchDeleteProductTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	if err := service.ProductType.BatchDeleteProductType(req.Ids); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "批量删除成功")
}

// GetProductTypeOptions 获取产品类型选项列表
func (a *ProductTypeApi) GetProductTypeOptions(c *gin.Context) {
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
	list, err := service.ProductType.GetProductTypeOptions(displayField, countTable, countForeignKey, excludeDeleted, countCreatedBy)
	if err != nil {
		response.Fail(c, "获取选项列表失败")
		return
	}
	response.OkWithData(c, list)
}
