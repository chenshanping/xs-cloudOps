package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExcelExporter Excel导出器
type ExcelExporter struct {
	file      *excelize.File
	sheetName string
	rowIndex  int
}

// NewExcelExporter 创建Excel导出器
func NewExcelExporter(sheetName string) *ExcelExporter {
	f := excelize.NewFile()
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1") // 删除默认sheet

	return &ExcelExporter{
		file:      f,
		sheetName: sheetName,
		rowIndex:  1,
	}
}

// SetHeaders 设置表头
func (e *ExcelExporter) SetHeaders(headers []string) error {
	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", columnName(i), e.rowIndex)
		if err := e.file.SetCellValue(e.sheetName, cell, header); err != nil {
			return err
		}
	}
	e.rowIndex++
	return nil
}

// SetHeadersWithComments 设置表头并添加批注
func (e *ExcelExporter) SetHeadersWithComments(headers []string, comments []string) error {
	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", columnName(i), e.rowIndex)
		if err := e.file.SetCellValue(e.sheetName, cell, header); err != nil {
			return err
		}
		// 添加批注

	}
	e.rowIndex++
	return nil
}

// AddDataValidation 为指定列添加下拉数据验证
func (e *ExcelExporter) AddDataValidation(colIndex int, options []string, startRow int, endRow int) error {
	if len(options) == 0 {
		return nil
	}

	// 构建下拉选项字符串
	formula := `"` + options[0]
	for i := 1; i < len(options); i++ {
		formula += "," + options[i]
	}
	formula += `"`

	// 设置数据验证范围
	colName := columnName(colIndex)
	sqref := fmt.Sprintf("%s%d:%s%d", colName, startRow, colName, endRow)

	dv := excelize.NewDataValidation(true)
	dv.Sqref = sqref
	dv.SetDropList(options)

	return e.file.AddDataValidation(e.sheetName, dv)
}

// AddRow 添加一行数据
func (e *ExcelExporter) AddRow(values []interface{}) error {
	for i, value := range values {
		cell := fmt.Sprintf("%s%d", columnName(i), e.rowIndex)
		if err := e.file.SetCellValue(e.sheetName, cell, value); err != nil {
			return err
		}
	}
	e.rowIndex++
	return nil
}

// SaveToBuffer 保存到缓冲区
func (e *ExcelExporter) SaveToBuffer() ([]byte, error) {
	buffer, err := e.file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// ExcelImporter Excel导入器
type ExcelImporter struct {
	file      *excelize.File
	sheetName string
}

// NewExcelImporter 创建Excel导入器
func NewExcelImporter(fileData []byte) (*ExcelImporter, error) {
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return nil, err
	}

	sheetName := f.GetSheetName(0)
	return &ExcelImporter{
		file:      f,
		sheetName: sheetName,
	}, nil
}

// GetHeaders 获取表头
func (e *ExcelImporter) GetHeaders() ([]string, error) {
	rows, err := e.file.GetRows(e.sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("Excel文件为空")
	}
	return rows[0], nil
}

// GetRows 获取所有数据行（不包含表头）
func (e *ExcelImporter) GetRows() ([][]string, error) {
	rows, err := e.file.GetRows(e.sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) <= 1 {
		return [][]string{}, nil
	}
	return rows[1:], nil
}

// columnName 将列索引转换为Excel列名 (0->A, 1->B, ..., 26->AA)
func columnName(index int) string {
	name := ""
	index++
	for index > 0 {
		index--
		name = string(rune('A'+index%26)) + name
		index /= 26
	}
	return name
}

// ParseCellValue 解析单元格值为指定类型
func ParseCellValue(value string, fieldType string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}

	switch fieldType {
	case "int":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return int(v), nil
	case "int32":
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(v), nil
	case "int64":
		return strconv.ParseInt(value, 10, 64)
	case "uint":
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return uint(v), nil
	case "uint32":
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return uint32(v), nil
	case "uint64":
		return strconv.ParseUint(value, 10, 64)
	case "float32":
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		return float32(v), nil
	case "float64":
		return strconv.ParseFloat(value, 64)
	case "bool":
		return strconv.ParseBool(value)
	case "time.Time":
		// 尝试多种时间格式
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02",
			"2006/01/02 15:04:05",
			"2006/01/02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, value); err == nil {
				return t, nil
			}
		}
		return nil, fmt.Errorf("无法解析时间: %s", value)
	default:
		return value, nil
	}
}

// ==================== 通用导入校验框架 ====================

// ImportField 导入字段定义
type ImportField struct {
	Header   string                            // 期望的表头名称
	Key      string                            // 字段键名
	Required bool                              // 是否必填
	Type     string                            // 数据类型: string/int/uint/float64/time.Time
	MaxLen   int                               // 最大长度（仅string有效，0=不限）
	Enum     []string                          // 枚举值（可选，为空则不校验枚举）
	Validate func(value string, row int) error // 自定义校验函数（可选）
}

// ImportError 导入错误详情
type ImportError struct {
	Row     int    `json:"row"`     // Excel行号（从2开始，1是表头）
	Column  string `json:"column"`  // 列名（表头名称）
	Value   string `json:"value"`   // 原始值
	Message string `json:"message"` // 错误描述
}

// ImportResult 导入校验结果
type ImportResult struct {
	TotalCount   int                      `json:"total_count"`   // 总行数
	SuccessCount int                      `json:"success_count"` // 成功条数
	FailedCount  int                      `json:"failed_count"`  // 失败条数
	Errors       []ImportError            `json:"errors"`        // 错误详情列表
	Data         []map[string]interface{} `json:"-"`             // 解析后的有效数据
}

// ValidateHeaders 校验表头是否与字段定义匹配
func ValidateHeaders(actual []string, fields []ImportField) error {
	expected := make([]string, len(fields))
	for i, f := range fields {
		expected[i] = f.Header
	}

	if len(actual) < len(expected) {
		return fmt.Errorf("表头列数不匹配，期望 %d 列，实际 %d 列。请下载最新导入模板", len(expected), len(actual))
	}

	var mismatched []string
	for i, exp := range expected {
		got := strings.TrimSpace(actual[i])
		if got != exp {
			mismatched = append(mismatched, fmt.Sprintf("第%d列期望【%s】实际【%s】", i+1, exp, got))
		}
	}

	if len(mismatched) > 0 {
		return fmt.Errorf("表头不匹配: %s。请下载最新导入模板", strings.Join(mismatched, "；"))
	}
	return nil
}

// ValidateImport 通用导入校验：校验表头 + 逐行校验数据
func ValidateImport(importer *ExcelImporter, fields []ImportField) (*ImportResult, error) {
	result := &ImportResult{}

	// 1. 获取并校验表头
	headers, err := importer.GetHeaders()
	if err != nil {
		return nil, fmt.Errorf("读取表头失败: %w", err)
	}
	if err := ValidateHeaders(headers, fields); err != nil {
		return nil, err
	}

	// 2. 获取数据行
	rows, err := importer.GetRows()
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("导入文件中没有数据行")
	}

	result.TotalCount = len(rows)

	// 3. 逐行校验
	for i, row := range rows {
		rowNum := i + 2 // Excel行号（第1行是表头）
		rowData := make(map[string]interface{})
		rowValid := true

		for j, field := range fields {
			value := ""
			if j < len(row) {
				value = strings.TrimSpace(row[j])
			}

			// 必填检查
			if field.Required && value == "" {
				result.Errors = append(result.Errors, ImportError{
					Row:     rowNum,
					Column:  field.Header,
					Value:   value,
					Message: fmt.Sprintf("【%s】不能为空", field.Header),
				})
				rowValid = false
				continue
			}

			// 空值跳过后续校验
			if value == "" {
				rowData[field.Key] = nil
				continue
			}

			// 最大长度检查
			if field.MaxLen > 0 && len([]rune(value)) > field.MaxLen {
				result.Errors = append(result.Errors, ImportError{
					Row:     rowNum,
					Column:  field.Header,
					Value:   value,
					Message: fmt.Sprintf("【%s】长度不能超过%d个字符", field.Header, field.MaxLen),
				})
				rowValid = false
				continue
			}

			// 枚举检查
			if len(field.Enum) > 0 {
				found := false
				for _, e := range field.Enum {
					if value == e {
						found = true
						break
					}
				}
				if !found {
					result.Errors = append(result.Errors, ImportError{
						Row:     rowNum,
						Column:  field.Header,
						Value:   value,
						Message: fmt.Sprintf("【%s】值无效，可选值: %s", field.Header, strings.Join(field.Enum, "/")),
					})
					rowValid = false
					continue
				}
			}

			// 类型检查
			parsed, parseErr := ParseCellValue(value, field.Type)
			if parseErr != nil {
				result.Errors = append(result.Errors, ImportError{
					Row:     rowNum,
					Column:  field.Header,
					Value:   value,
					Message: fmt.Sprintf("【%s】格式错误: %s", field.Header, parseErr.Error()),
				})
				rowValid = false
				continue
			}

			// 自定义校验
			if field.Validate != nil {
				if vErr := field.Validate(value, rowNum); vErr != nil {
					result.Errors = append(result.Errors, ImportError{
						Row:     rowNum,
						Column:  field.Header,
						Value:   value,
						Message: vErr.Error(),
					})
					rowValid = false
					continue
				}
			}

			rowData[field.Key] = parsed
		}

		if rowValid {
			result.Data = append(result.Data, rowData)
			result.SuccessCount++
		} else {
			result.FailedCount++
		}
	}

	return result, nil
}
