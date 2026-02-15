package v1

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go-base-server/global"
	"go-base-server/middleware"
	"go-base-server/model"
	"go-base-server/model/request"
	"go-base-server/model/response"
	"go-base-server/service"
	"go-base-server/utils"
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

// ExportProduct 导出产品信息
func (a *ProductApi) ExportProduct(c *gin.Context) {
	var req request.ProductQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	list, err := service.Product.GetAllProduct(&req)
	if err != nil {
		response.Fail(c, "获取数据失败: "+err.Error())
		return
	}

	// 创建Excel导出器
	exporter := utils.NewExcelExporter("产品信息")

	// 设置表头
	headers := []string{
		"产品名称",
		"产品数量",
		"产品单价",
		"状态",
		"产品类型",
	}
	if err := exporter.SetHeaders(headers); err != nil {
		response.Fail(c, "设置表头失败: "+err.Error())
		return
	}

	// 添加数据行
	for _, item := range list {
		row := []interface{}{
			item.Name,
			item.Num,
			item.Price,
			service.Dict.GetDictLabel("common_status", fmt.Sprint(item.Status)),
			func() string {
				if item.ProductType != nil {
					return item.ProductType.Name
				}
				return ""
			}(),
		}
		if err := exporter.AddRow(row); err != nil {
			response.Fail(c, "添加数据行失败: "+err.Error())
			return
		}
	}

	// 生成文件
	buffer, err := exporter.SaveToBuffer()
	if err != nil {
		response.Fail(c, "生成Excel失败: "+err.Error())
		return
	}

	// 返回文件
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=product_%s.xlsx", time.Now().Format("20060102150405")))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer)
}

// ImportProduct 导入产品信息
func (a *ProductApi) ImportProduct(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请上传文件")
		return
	}

	// 读取文件内容
	f, err := file.Open()
	if err != nil {
		response.Fail(c, "读取文件失败: "+err.Error())
		return
	}
	defer f.Close()

	fileData := make([]byte, file.Size)
	if _, err := f.Read(fileData); err != nil {
		response.Fail(c, "读取文件内容失败: "+err.Error())
		return
	}

	// 创建Excel导入器
	importer, err := utils.NewExcelImporter(fileData)
	if err != nil {
		response.Fail(c, "解析Excel失败: "+err.Error())
		return
	}

	// 获取数据行
	rows, err := importer.GetRows()
	if err != nil {
		response.Fail(c, "读取数据失败: "+err.Error())
		return
	}

	if len(rows) == 0 {
		response.Fail(c, "Excel文件无数据")
		return
	}

	// 解析并导入数据
	successCount := 0
	failCount := 0
	var errors []string

	for i, row := range rows {
		rowNum := i + 2 // Excel行号（从2开始，1是表头）

		// 跳过空行
		isEmpty := true
		for _, cell := range row {
			if cell != "" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		// 构建请求对象
		var createReq request.CreateProductRequest
		colIndex := 0
		hasError := false
		var rowErrors []string
		// 产品名称
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				if val, err := utils.ParseCellValue(cellValue, "string"); err == nil {
					createReq.Name = val.(string)
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("产品名称格式错误: %v", err))
					hasError = true
				}
			} else {
				rowErrors = append(rowErrors, "产品名称为必填项")
				hasError = true
			}
		} else {
			rowErrors = append(rowErrors, "产品名称为必填项")
			hasError = true
		}
		colIndex++
		// 产品数量
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				if val, err := utils.ParseCellValue(cellValue, "int"); err == nil {
					createReq.Num = val.(int)
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("产品数量格式错误: %v", err))
					hasError = true
				}
			}
		}
		colIndex++
		// 产品单价
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				if val, err := utils.ParseCellValue(cellValue, "float64"); err == nil {
					createReq.Price = val.(float64)
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("产品单价格式错误: %v", err))
					hasError = true
				}
			}
		}
		colIndex++
		// 状态 - 数据字典转换和验证
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				dictValue := service.Dict.GetDictValue("common_status", cellValue)
				if dictValue == "" {
					// 获取所有有效值
					dictList, _ := service.Dict.GetDictDataByType("common_status")
					validValues := make([]string, 0, len(dictList))
					for _, item := range dictList {
						validValues = append(validValues, item.Label)
					}
					rowErrors = append(rowErrors, fmt.Sprintf("状态值\"%s\"不存在，有效值: %s", cellValue, strings.Join(validValues, ", ")))
					hasError = true
				} else {
					if val, err := utils.ParseCellValue(dictValue, "string"); err == nil {
						createReq.Status = val.(string)
					} else {
						rowErrors = append(rowErrors, fmt.Sprintf("状态格式错误: %v", err))
						hasError = true
					}
				}
			} else {
				rowErrors = append(rowErrors, "状态为必填项")
				hasError = true
			}
		}
		colIndex++
		// 产品类型 - 关联查询和验证
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				var related model.ProductType
				if err := global.DB.Where("name = ?", cellValue).First(&related).Error; err == nil {
					createReq.TypeId = related.ID
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("产品类型\"%s\"不存在，请先创建该产品类型", cellValue))
					hasError = true
				}
			} else {
				rowErrors = append(rowErrors, "产品类型为必填项")
				hasError = true
			}
		} else {
			rowErrors = append(rowErrors, "产品类型为必填项")
			hasError = true
		}
		colIndex++

		// 如果有错误，记录并跳过
		if hasError {
			errors = append(errors, fmt.Sprintf("第%d行: %s", rowNum, strings.Join(rowErrors, "; ")))
			failCount++
			continue
		}

		// 保存数据
		if err := service.Product.CreateProduct(&createReq, 0); err != nil {
			errors = append(errors, fmt.Sprintf("第%d行保存失败: %v", rowNum, err))
			failCount++
		} else {
			successCount++
		}
	}

	// 返回导入结果
	result := map[string]interface{}{
		"success_count": successCount,
		"fail_count":    failCount,
		"total":         len(rows),
		"errors":        errors,
	}
	response.OkWithData(c, result)
}

// DownloadTemplateProduct 下载导入模板
func (a *ProductApi) DownloadTemplateProduct(c *gin.Context) {
	// 创建Excel导出器
	exporter := utils.NewExcelExporter("产品信息模板")

	// 准备表头和批注
	headers := []string{
		"产品名称(必填)",
		"产品数量(选填)",
		"产品单价(选填)",
		"状态(必填)",
		"产品类型(必填)",
	}

	comments := []string{
		"请填写产品名称",
		"请填写整数",
		"请填写数字",
		"数据字典: common_status\\n请从下拉列表中选择",
		"关联字段: name\\n请填写已存在的产品类型名称",
	}

	// 设置表头和批注
	if err := exporter.SetHeadersWithComments(headers, comments); err != nil {
		response.Fail(c, "设置表头失败: "+err.Error())
		return
	}

	// 为数据字典和开关字段添加下拉验证
	colIndex := 0
	colIndex++
	colIndex++
	colIndex++
	// 状态 - 数据字典下拉
	{
		dictList, _ := service.Dict.GetDictDataByType("common_status")
		options := make([]string, 0, len(dictList))
		for _, item := range dictList {
			options = append(options, item.Label)
		}
		if len(options) > 0 {
			exporter.AddDataValidation(colIndex, options, 2, 1000)
		}
	}
	colIndex++
	// 产品类型 - 关联字段下拉
	{
		// 使用options接口获取选项（已过滤status等条件）
		optionsList, _ := service.ProductType.GetProductTypeOptions("name", "", "", false, 0)
		options := make([]string, 0, len(optionsList))
		for _, item := range optionsList {
			if name, ok := item["name"].(string); ok {
				options = append(options, name)
			}
		}
		if len(options) > 0 && len(options) <= 100 {
			// 只有选项数量不超过100时才添加下拉验证（Excel限制）
			exporter.AddDataValidation(colIndex, options, 2, 1000)
		}
	}
	colIndex++

	// 添加示例数据行
	example := []interface{}{
		"",
		0,
		0.0,
		"", // 状态 - 从下拉列表选择
		"", // 产品类型 - 填写name
	}
	exporter.AddRow(example)

	// 生成文件
	buffer, err := exporter.SaveToBuffer()
	if err != nil {
		response.Fail(c, "生成模板失败: "+err.Error())
		return
	}

	// 返回文件
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=product_template.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer)
}
