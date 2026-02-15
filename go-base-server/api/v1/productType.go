package v1

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-base-server/middleware"
	"go-base-server/model/request"
	"go-base-server/model/response"
	"go-base-server/service"
	"go-base-server/utils"

	"github.com/gin-gonic/gin"
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

// ExportProductType 导出产品类型
func (a *ProductTypeApi) ExportProductType(c *gin.Context) {
	var req request.ProductTypeQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	list, err := service.ProductType.GetAllProductType(&req)
	if err != nil {
		response.Fail(c, "获取数据失败: "+err.Error())
		return
	}

	// 创建Excel导出器
	exporter := utils.NewExcelExporter("产品类型")

	// 设置表头
	headers := []string{
		"产品类型名称",
		"类型图标",
		"",
	}
	if err := exporter.SetHeaders(headers); err != nil {
		response.Fail(c, "设置表头失败: "+err.Error())
		return
	}

	// 添加数据行
	for _, item := range list {
		row := []interface{}{
			item.Name,
			item.Icon,
			service.Dict.GetDictLabel("common_status", fmt.Sprint(item.Status)),
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
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=productType_%s.xlsx", time.Now().Format("20060102150405")))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer)
}

// ImportProductType 导入产品类型
func (a *ProductTypeApi) ImportProductType(c *gin.Context) {
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
		var createReq request.CreateProductTypeRequest
		colIndex := 0
		hasError := false
		var rowErrors []string
		// 产品类型名称
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				if val, err := utils.ParseCellValue(cellValue, "string"); err == nil {
					createReq.Name = val.(string)
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("产品类型名称格式错误: %v", err))
					hasError = true
				}
			}
		}
		colIndex++
		// 类型图标
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				if val, err := utils.ParseCellValue(cellValue, "string"); err == nil {
					createReq.Icon = val.(string)
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("类型图标格式错误: %v", err))
					hasError = true
				}
			}
		}
		colIndex++
		//  - 数据字典转换和验证
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
					rowErrors = append(rowErrors, fmt.Sprintf("值\"%s\"不存在，有效值: %s", cellValue, strings.Join(validValues, ", ")))
					hasError = true
				} else {
					if val, err := utils.ParseCellValue(dictValue, "string"); err == nil {
						createReq.Status = val.(string)
					} else {
						rowErrors = append(rowErrors, fmt.Sprintf("格式错误: %v", err))
						hasError = true
					}
				}
			}
		}
		colIndex++

		// 如果有错误，记录并跳过
		if hasError {
			errors = append(errors, fmt.Sprintf("第%d行: %s", rowNum, strings.Join(rowErrors, "; ")))
			failCount++
			continue
		}

		// 保存数据
		if err := service.ProductType.CreateProductType(&createReq, 0); err != nil {
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

// DownloadTemplateProductType 下载导入模板
func (a *ProductTypeApi) DownloadTemplateProductType(c *gin.Context) {
	// 创建Excel导出器
	exporter := utils.NewExcelExporter("产品类型模板")

	// 准备表头和批注
	headers := []string{
		"产品类型名称(选填)",
		"类型图标(选填)",
		"(选填)",
	}

	comments := []string{
		"请填写产品类型名称",
		"请填写类型图标",
		"数据字典: common_status\\n请从下拉列表中选择",
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
	//  - 数据字典下拉
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

	// 添加示例数据行
	example := []interface{}{
		"",
		"",
		"", //  - 从下拉列表选择
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
	c.Header("Content-Disposition", "attachment; filename=productType_template.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer)
}
