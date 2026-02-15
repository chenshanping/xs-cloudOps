package utils

import (
	"bytes"
	"fmt"
	"strconv"
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
func (e *ExcelExporter) AddDataValidation(colIndex int, options []string, startRow, endRow int) error {
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
