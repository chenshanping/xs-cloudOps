package v1

import (
	"fmt"
	"strconv"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
{{- if or .HasRelations .HasCreatedBy}}

	"go-base-server/global"
{{- end}}
	"go-base-server/middleware"
{{- if .HasRelations}}
	"go-base-server/model"
{{- end}}
	"go-base-server/model/request"
	"go-base-server/model/response"
	"go-base-server/service"
	"go-base-server/utils"
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
{{- if .HasStats}}
{{- range $chart := .StatsCharts}}

// Get{{$.ModelName}}Stats{{$chart.Field}} 获取{{$.Description}}按{{$chart.Title}}分组统计
func (a *{{$.ModelName}}Api) Get{{$.ModelName}}Stats{{$chart.Field}}(c *gin.Context) {
	data, err := service.{{$.ModelName}}.Get{{$.ModelName}}Stats{{$chart.Field}}()
	if err != nil {
		response.Fail(c, "获取统计数据失败")
		return
	}
	response.OkWithData(c, data)
}
{{- end}}
{{- if .HasStatsTrend}}

// Get{{.ModelName}}TrendStats 获取{{.Description}}趋势统计
func (a *{{.ModelName}}Api) Get{{.ModelName}}TrendStats(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}
	data, err := service.{{.ModelName}}.Get{{.ModelName}}TrendStats(days)
	if err != nil {
		response.Fail(c, "获取趋势数据失败")
		return
	}
	response.OkWithData(c, data)
}
{{- end}}
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

// Export{{.ModelName}} 导出{{.Description}}
func (a *{{.ModelName}}Api) Export{{.ModelName}}(c *gin.Context) {
	var req request.{{.ModelName}}QueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

{{- if .DataIsolation}}
	userID := middleware.GetUserID(c)
	roleIDs := middleware.GetUserRoleIDs(c)
	isAdmin := CheckIsAdmin(roleIDs, "{{.AdminRoleIds}}")
	list, err := service.{{.ModelName}}.GetAll{{.ModelName}}(&req, userID, isAdmin)
{{- else}}
	list, err := service.{{.ModelName}}.GetAll{{.ModelName}}(&req)
{{- end}}
	if err != nil {
		response.Fail(c, "获取数据失败: "+err.Error())
		return
	}

	// 创建Excel导出器
	exporter := utils.NewExcelExporter("{{.Description}}")

	// 设置表头
	headers := []string{
		{{- range .ListColumns}}
		"{{.Comment}}",
		{{- end}}
		{{- range .Relations}}
		{{- if eq .RelationType "belongsTo"}}
		"{{.Comment}}",
		{{- end}}
		{{- end}}
	}
	if err := exporter.SetHeaders(headers); err != nil {
		response.Fail(c, "设置表头失败: "+err.Error())
		return
	}

	// 添加数据行
	for _, item := range list {
		row := []interface{}{
			{{- range .ListColumns}}
			{{- if eq .FormType "select"}}
			{{- if .DictType}}
			service.Dict.GetDictLabel("{{.DictType}}", fmt.Sprint(item.{{.FieldName}})),
			{{- else}}
			item.{{.FieldName}},
			{{- end}}
			{{- else if eq .FormType "switch"}}
			{{- if .SwitchValues}}
			func() string { if item.{{.FieldName}} == {{.SwitchValues.ActiveValue}} { return "{{.SwitchValues.ActiveText}}" }; return "{{.SwitchValues.InactiveText}}" }(),
			{{- else}}
			item.{{.FieldName}},
			{{- end}}
			{{- else if eq .FieldType "time.Time"}}
			item.{{.FieldName}}.Format("2006-01-02 15:04:05"),
			{{- else}}
			item.{{.FieldName}},
			{{- end}}
			{{- end}}
			{{- range .Relations}}
			{{- if eq .RelationType "belongsTo"}}
			func() string { if item.{{.FieldName}} != nil { return item.{{.FieldName}}.{{ToPascalCase .DisplayField}} }; return "" }(),
			{{- end}}
			{{- end}}
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
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename={{.ModuleName}}_%s.xlsx", time.Now().Format("20060102150405")))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer)
}

// Import{{.ModelName}} 导入{{.Description}}
func (a *{{.ModelName}}Api) Import{{.ModelName}}(c *gin.Context) {
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
{{- if .HasCreatedBy}}
	userID := middleware.GetUserID(c)
{{- end}}

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
		var createReq request.Create{{.ModelName}}Request
		colIndex := 0
		hasError := false
		var rowErrors []string

		{{- range .FormColumns}}
		{{- if eq .FormType "select"}}
		{{- if .DictType}}
		// {{.Comment}} - 数据字典转换和验证
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				dictValue := service.Dict.GetDictValue("{{.DictType}}", cellValue)
				if dictValue == "" {
					// 获取所有有效值
					dictList, _ := service.Dict.GetDictDataByType("{{.DictType}}")
					validValues := make([]string, 0, len(dictList))
					for _, item := range dictList {
						validValues = append(validValues, item.Label)
					}
					rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}值\"%s\"不存在，有效值: %s", cellValue, strings.Join(validValues, ", ")))
					hasError = true
				} else {
					if val, err := utils.ParseCellValue(dictValue, "{{.FieldType}}"); err == nil {
						{{- if and .IsRequired (or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64"))}}
						v := val.({{.FieldType}})
						createReq.{{.FieldName}} = &v
						{{- else}}
						createReq.{{.FieldName}} = val.({{.FieldType}})
						{{- end}}
					} else {
						rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}格式错误: %v", err))
						hasError = true
					}
				}
			}{{if .IsRequired}} else {
				rowErrors = append(rowErrors, "{{.Comment}}为必填项")
				hasError = true
			}{{end}}
		}
		{{- else if .SelectOptions}}
		// {{.Comment}} - 固定选项验证
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				validOptions := map[string]{{.FieldType}}{
					{{- range .SelectOptions}}
					"{{.Label}}": {{.Value}},
					{{- end}}
				}
				if val, ok := validOptions[cellValue]; ok {
					{{- if and .IsRequired (or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64"))}}
					createReq.{{.FieldName}} = &val
					{{- else}}
					createReq.{{.FieldName}} = val
					{{- end}}
				} else {
					validValues := []string{
						{{- range .SelectOptions}}
						"{{.Label}}",
						{{- end}}
					}
					rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}值\"%s\"无效，有效值: %s", cellValue, strings.Join(validValues, ", ")))
					hasError = true
				}
			}{{if .IsRequired}} else {
				rowErrors = append(rowErrors, "{{.Comment}}为必填项")
				hasError = true
			}{{end}}
		}
		{{- else}}
		// {{.Comment}}
		if colIndex < len(row) && row[colIndex] != "" {
			if val, err := utils.ParseCellValue(row[colIndex], "{{.FieldType}}"); err == nil {
				{{- if and .IsRequired (or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64"))}}
				v := val.({{.FieldType}})
				createReq.{{.FieldName}} = &v
				{{- else}}
				createReq.{{.FieldName}} = val.({{.FieldType}})
				{{- end}}
			} else {
				rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}格式错误: %v", err))
				hasError = true
			}
		}{{if .IsRequired}} else {
			rowErrors = append(rowErrors, "{{.Comment}}为必填项")
			hasError = true
		}{{end}}
		{{- end}}
		{{- else if eq .FormType "switch"}}
		// {{.Comment}} - 开关转换和验证
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				{{- if .SwitchValues}}
				if cellValue == "{{.SwitchValues.ActiveText}}" {
					createReq.{{.FieldName}} = {{.SwitchValues.ActiveValue}}
				} else if cellValue == "{{.SwitchValues.InactiveText}}" {
					createReq.{{.FieldName}} = {{.SwitchValues.InactiveValue}}
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}值\"%s\"无效，有效值: {{.SwitchValues.ActiveText}}, {{.SwitchValues.InactiveText}}", cellValue))
					hasError = true
				}
				{{- else}}
				if val, err := utils.ParseCellValue(cellValue, "{{.FieldType}}"); err == nil {
					createReq.{{.FieldName}} = val.({{.FieldType}})
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}格式错误: %v", err))
					hasError = true
				}
				{{- end}}
			}{{if .IsRequired}} else {
				rowErrors = append(rowErrors, "{{.Comment}}为必填项")
				hasError = true
			}{{end}}
		}
		{{- else}}
		// {{.Comment}}
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				if val, err := utils.ParseCellValue(cellValue, "{{.FieldType}}"); err == nil {
					{{- if and .IsRequired (or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64"))}}
					v := val.({{.FieldType}})
					createReq.{{.FieldName}} = &v
					{{- else}}
					createReq.{{.FieldName}} = val.({{.FieldType}})
					{{- end}}
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}格式错误: %v", err))
					hasError = true
				}
			}{{if .IsRequired}} else {
				rowErrors = append(rowErrors, "{{.Comment}}为必填项")
				hasError = true
			}{{end}}
		}{{if .IsRequired}} else {
			rowErrors = append(rowErrors, "{{.Comment}}为必填项")
			hasError = true
		}{{end}}
		{{- end}}
		colIndex++
		{{- end}}

		{{- range .Relations}}
		{{- if eq .RelationType "belongsTo"}}
		// {{.Comment}} - 关联查询和验证
		if colIndex < len(row) {
			cellValue := row[colIndex]
			if cellValue != "" {
				{{- if .UseOptionsApi}}
				// 使用options接口验证（已过滤status等条件）
				optionsList, _ := service.{{.RelatedModel}}.Get{{.RelatedModel}}Options("{{.DisplayField}}", "", "", false, 0)
				found := false
				var relatedID uint
				for _, item := range optionsList {
					if name, ok := item["name"].(string); ok && name == cellValue {
						if id, ok := item["id"].(uint); ok {
							relatedID = id
							found = true
							break
						}
						// 处理id可能是float64的情况（JSON解析）
						if id, ok := item["id"].(float64); ok {
							relatedID = uint(id)
							found = true
							break
						}
					}
				}
				if found {
					createReq.{{.ForeignKey | ToPascalCase}} = relatedID
				} else {
					// 获取所有有效值
					validValues := make([]string, 0, len(optionsList))
					for _, item := range optionsList {
						if name, ok := item["name"].(string); ok {
							validValues = append(validValues, name)
						}
					}
					if len(validValues) > 0 {
						rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}\"%s\"不存在，有效值: %s", cellValue, strings.Join(validValues, ", ")))
					} else {
						rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}\"%s\"不存在，请先创建该{{.Comment}}", cellValue))
					}
					hasError = true
				}
				{{- else}}
				// 直接查询关联表
				var related model.{{.RelatedModel}}
				if err := global.DB.Where("{{.DisplayField}} = ?", cellValue).First(&related).Error; err == nil {
					createReq.{{.ForeignKey | ToPascalCase}} = related.ID
				} else {
					rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}\"%s\"不存在，请先创建该{{.Comment}}", cellValue))
					hasError = true
				}
				{{- end}}
			}{{if .IsRequired}} else {
				rowErrors = append(rowErrors, "{{.Comment}}为必填项")
				hasError = true
			}{{end}}
		}{{if .IsRequired}} else {
			rowErrors = append(rowErrors, "{{.Comment}}为必填项")
			hasError = true
		}{{end}}
		colIndex++
		{{- end}}
		{{- end}}

		// 如果有错误，记录并跳过
		if hasError {
			errors = append(errors, fmt.Sprintf("第%d行: %s", rowNum, strings.Join(rowErrors, "; ")))
			failCount++
			continue
		}

		// 保存数据
{{- if .HasCreatedBy}}
		if err := service.{{.ModelName}}.Create{{.ModelName}}(&createReq, userID); err != nil {
{{- else}}
		if err := service.{{.ModelName}}.Create{{.ModelName}}(&createReq, 0); err != nil {
{{- end}}
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

// DownloadTemplate{{.ModelName}} 下载导入模板
func (a *{{.ModelName}}Api) DownloadTemplate{{.ModelName}}(c *gin.Context) {
	// 创建Excel导出器
	exporter := utils.NewExcelExporter("{{.Description}}模板")

	// 准备表头和批注
	headers := []string{
		{{- range .FormColumns}}
		"{{.Comment}}{{if .IsRequired}}(必填){{else}}(选填){{end}}",
		{{- end}}
		{{- range .Relations}}
		{{- if eq .RelationType "belongsTo"}}
		"{{.Comment}}{{if .IsRequired}}(必填){{else}}(选填){{end}}",
		{{- end}}
		{{- end}}
	}
	
	comments := []string{
		{{- range .FormColumns}}
		{{- if eq .FormType "select"}}
		{{- if .DictType}}
		"数据字典: {{.DictType}}\\n请从下拉列表中选择",
		{{- else if .SelectOptions}}
		"可选值: {{range $i, $opt := .SelectOptions}}{{if $i}}, {{end}}{{$opt.Label}}{{end}}\\n请从下拉列表中选择",
		{{- else}}
		"请填写有效值",
		{{- end}}
		{{- else if eq .FormType "switch"}}
		{{- if .SwitchValues}}
		"可选值: {{.SwitchValues.ActiveText}}, {{.SwitchValues.InactiveText}}\\n请从下拉列表中选择",
		{{- else}}
		"可选值: 是, 否\\n请从下拉列表中选择",
		{{- end}}
		{{- else if eq .FieldType "int"}}
		"请填写整数",
		{{- else if eq .FieldType "float64"}}
		"请填写数字",
		{{- else if eq .FieldType "time.Time"}}
		"时间格式: 2006-01-02 15:04:05 或 2006-01-02",
		{{- else}}
		"请填写{{.Comment}}",
		{{- end}}
		{{- end}}
		{{- range .Relations}}
		{{- if eq .RelationType "belongsTo"}}
		"关联字段: {{.DisplayField}}\\n请填写已存在的{{.Comment}}名称",
		{{- end}}
		{{- end}}
	}
	
	// 设置表头和批注
	if err := exporter.SetHeadersWithComments(headers, comments); err != nil {
		response.Fail(c, "设置表头失败: "+err.Error())
		return
	}

	// 为数据字典和开关字段添加下拉验证
	colIndex := 0
	{{- range .FormColumns}}
	{{- if eq .FormType "select"}}
	{{- if .DictType}}
	// {{.Comment}} - 数据字典下拉
	{
		dictList, _ := service.Dict.GetDictDataByType("{{.DictType}}")
		options := make([]string, 0, len(dictList))
		for _, item := range dictList {
			options = append(options, item.Label)
		}
		if len(options) > 0 {
			exporter.AddDataValidation(colIndex, options, 2, 1000)
		}
	}
	{{- else if .SelectOptions}}
	// {{.Comment}} - 固定选项下拉
	{
		options := []string{
			{{- range .SelectOptions}}
			"{{.Label}}",
			{{- end}}
		}
		exporter.AddDataValidation(colIndex, options, 2, 1000)
	}
	{{- end}}
	{{- else if eq .FormType "switch"}}
	// {{.Comment}} - 开关下拉
	{
		{{- if .SwitchValues}}
		options := []string{"{{.SwitchValues.ActiveText}}", "{{.SwitchValues.InactiveText}}"}
		{{- else}}
		options := []string{"是", "否"}
		{{- end}}
		exporter.AddDataValidation(colIndex, options, 2, 1000)
	}
	{{- end}}
	colIndex++
	{{- end}}
	{{- range .Relations}}
	{{- if eq .RelationType "belongsTo"}}
	// {{.Comment}} - 关联字段下拉
	{{- if .UseOptionsApi}}
	{
		// 使用options接口获取选项（已过滤status等条件）
		optionsList, _ := service.{{.RelatedModel}}.Get{{.RelatedModel}}Options("{{.DisplayField}}", "", "", false, 0)
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
	{{- else}}
	{
		// 直接查询关联表
		var relatedList []model.{{.RelatedModel}}
		if err := global.DB.Find(&relatedList).Error; err == nil {
			options := make([]string, 0, len(relatedList))
			for _, item := range relatedList {
				options = append(options, item.{{ToPascalCase .DisplayField}})
			}
			if len(options) > 0 && len(options) <= 100 {
				// 只有选项数量不超过100时才添加下拉验证（Excel限制）
				exporter.AddDataValidation(colIndex, options, 2, 1000)
			}
		}
	}
	{{- end}}
	colIndex++
	{{- end}}
	{{- end}}

	// 添加示例数据行
	example := []interface{}{
		{{- range .FormColumns}}
		{{- if eq .FormType "select"}}
		{{- if .DictType}}
		"", // {{.Comment}} - 从下拉列表选择
		{{- else if .SelectOptions}}
		"{{(index .SelectOptions 0).Label}}", // {{.Comment}}
		{{- else}}
		"",
		{{- end}}
		{{- else if eq .FormType "switch"}}
		{{- if .SwitchValues}}
		"{{.SwitchValues.ActiveText}}", // {{.Comment}}
		{{- else}}
		"是",
		{{- end}}
		{{- else if eq .FieldType "int"}}
		0,
		{{- else if eq .FieldType "float64"}}
		0.0,
		{{- else if eq .FieldType "time.Time"}}
		"2024-01-01 00:00:00",
		{{- else}}
		"",
		{{- end}}
		{{- end}}
		{{- range .Relations}}
		{{- if eq .RelationType "belongsTo"}}
		"", // {{.Comment}} - 填写{{.DisplayField}}
		{{- end}}
		{{- end}}
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
	c.Header("Content-Disposition", "attachment; filename={{.ModuleName}}_template.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer)
}
