package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"go-base-server/middleware"
	"go-base-server/model/request"
	"go-base-server/model/response"
	"go-base-server/service"
)

type ProductApi struct{}

var Product = new(ProductApi)

// GetProductList 获取产品信息列表
func (a *ProductApi) GetProductList(c *gin.Context) {
	var req request.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	list, total, err := service.Product.GetProductList(&req)
	if err != nil {
		response.Fail(c, "获取列表失败")
		return
	}

	response.OkWithPage(c, list, total, req.Page, req.PageSize)
}

// GetProduct 获取产品信息详情
func (a *ProductApi) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	data, err := service.Product.GetProduct(uint(id))
	if err != nil {
		response.Fail(c, "获取详情失败")
		return
	}

	response.OkWithData(c, data)
}

// CreateProduct 创建产品信息
func (a *ProductApi) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := service.Product.CreateProduct(&req, userID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "创建成功")
}

// UpdateProduct 更新产品信息
func (a *ProductApi) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := service.Product.UpdateProduct(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// DeleteProduct 删除产品信息
func (a *ProductApi) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	if err := service.Product.DeleteProduct(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

// BatchDeleteProduct 批量删除产品信息
func (a *ProductApi) BatchDeleteProduct(c *gin.Context) {
	var req request.BatchDeleteProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	if err := service.Product.BatchDeleteProduct(req.Ids); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "批量删除成功")
}

// GetProductOptions 获取产品信息选项列表
func (a *ProductApi) GetProductOptions(c *gin.Context) {
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
	list, err := service.Product.GetProductOptions(displayField, countTable, countForeignKey, excludeDeleted, countCreatedBy)
	if err != nil {
		response.Fail(c, "获取选项列表失败")
		return
	}
	response.OkWithData(c, list)
}

// GetProductStatsTypeId 获取产品信息按产品类型分组统计
func (a *ProductApi) GetProductStatsTypeId(c *gin.Context) {
	data, err := service.Product.GetProductStatsTypeId()
	if err != nil {
		response.Fail(c, "获取统计数据失败")
		return
	}
	response.OkWithData(c, data)
}

// GetProductStatsStatus 获取产品信息按产品状态分组统计
func (a *ProductApi) GetProductStatsStatus(c *gin.Context) {
	data, err := service.Product.GetProductStatsStatus()
	if err != nil {
		response.Fail(c, "获取统计数据失败")
		return
	}
	response.OkWithData(c, data)
}

// GetProductTrendStats 获取产品信息趋势统计
func (a *ProductApi) GetProductTrendStats(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}
	data, err := service.Product.GetProductTrendStats(days)
	if err != nil {
		response.Fail(c, "获取趋势数据失败")
		return
	}
	response.OkWithData(c, data)
}
