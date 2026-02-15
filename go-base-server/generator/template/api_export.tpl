// Export{{.ModelName}} 导出{{.Description}}
func (a *{{.ModelName}}Api) Export{{.ModelName}}(c *gin.Context) {
	var req request.{{.ModelName}}QueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 获取数据
	list, err := service.{{.ModelName}}Service.GetAll{{.ModelName}}(&req)
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
			service.DictService.GetDictLabel("{{.DictType}}", fmt.Sprint(item.{{.FieldName}})),
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

	for i, row := range rows {
		rowNum := i + 2 // Excel行号（从2开始，1是表头）
		
		// 构建数据对象
		var data model.{{.ModelName}}
		colIndex := 0

		{{- range .FormColumns}}
		{{- if eq .FormType "select"}}
		{{- if .DictType}}
		// {{.Comment}} - 数据字典转换
		if colIndex < len(row) && row[colIndex] != "" {
			dictValue := service.DictService.GetDictValue("{{.DictType}}", row[colIndex])
			if dictValue != "" {
				if val, err := utils.ParseCellValue(dictValue, "{{.FieldType}}"); err == nil {
					data.{{.FieldName}} = val.({{.FieldType}})
				} else {
					errors = append(errors, fmt.Sprintf("第%d行{{.Comment}}格式错误: %v", rowNum, err))
					failCount++
					continue
				}
			}
		}
		{{- else}}
		// {{.Comment}}
		if colIndex < len(row) && row[colIndex] != "" {
			if val, err := utils.ParseCellValue(row[colIndex], "{{.FieldType}}"); err == nil {
				data.{{.FieldName}} = val.({{.FieldType}})
			} else {
				errors = append(errors, fmt.Sprintf("第%d行{{.Comment}}格式错误: %v", rowNum, err))
				failCount++
				continue
			}
		}
		{{- end}}
		{{- else if eq .FormType "switch"}}
		// {{.Comment}} - 开关转换
		if colIndex < len(row) && row[colIndex] != "" {
			{{- if .SwitchValues}}
			if row[colIndex] == "{{.SwitchValues.ActiveText}}" {
				data.{{.FieldName}} = {{.SwitchValues.ActiveValue}}
			} else {
				data.{{.FieldName}} = {{.SwitchValues.InactiveValue}}
			}
			{{- else}}
			if val, err := utils.ParseCellValue(row[colIndex], "{{.FieldType}}"); err == nil {
				data.{{.FieldName}} = val.({{.FieldType}})
			}
			{{- end}}
		}
		{{- else}}
		// {{.Comment}}
		if colIndex < len(row) && row[colIndex] != "" {
			if val, err := utils.ParseCellValue(row[colIndex], "{{.FieldType}}"); err == nil {
				data.{{.FieldName}} = val.({{.FieldType}})
			} else {
				errors = append(errors, fmt.Sprintf("第%d行{{.Comment}}格式错误: %v", rowNum, err))
				failCount++
				continue
			}
		}
		{{- end}}
		colIndex++
		{{- end}}

		{{- range .Relations}}
		{{- if eq .RelationType "belongsTo"}}
		// {{.Comment}} - 关联查询
		if colIndex < len(row) && row[colIndex] != "" {
			var related model.{{.RelatedModel}}
			if err := global.DB.Where("{{.DisplayField}} = ?", row[colIndex]).First(&related).Error; err == nil {
				data.{{.ForeignKey}} = related.ID
			} else {
				errors = append(errors, fmt.Sprintf("第%d行{{.Comment}}不存在: %s", rowNum, row[colIndex]))
				failCount++
				continue
			}
		}
		colIndex++
		{{- end}}
		{{- end}}

		{{- if .HasCreatedBy}}
		// 设置创建人
		userID := middleware.GetUserID(c)
		data.CreatedBy = userID
		{{- end}}

		// 保存数据
		if err := service.{{.ModelName}}Service.Create{{.ModelName}}(&data); err != nil {
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

	// 设置表头
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
	if err := exporter.SetHeaders(headers); err != nil {
		response.Fail(c, "生成模板失败: "+err.Error())
		return
	}

	// 添加示例数据行（可选）
	example := []interface{}{
		{{- range .FormColumns}}
		{{- if eq .FormType "select"}}
		{{- if .DictType}}
		"示例值", // {{.Comment}} - 请填写数据字典中的标签
		{{- else if .SelectOptions}}
		"{{(index .SelectOptions 0).Label}}", // {{.Comment}}
		{{- else}}
		"示例值",
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
		"示例值",
		{{- end}}
		{{- end}}
		{{- range .Relations}}
		{{- if eq .RelationType "belongsTo"}}
		"示例{{.Comment}}", // {{.Comment}} - 请填写{{.DisplayField}}
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
