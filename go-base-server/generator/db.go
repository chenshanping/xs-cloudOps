package generator

import (
	"strings"

	"go-base-server/global"
)

// GetTables 获取数据库所有表
func GetTables() ([]TableInfo, error) {
	var tables []TableInfo

	sql := `
		SELECT 
			TABLE_NAME as table_name,
			IFNULL(TABLE_COMMENT, '') as table_comment
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = DATABASE()
		ORDER BY TABLE_NAME
	`

	if err := global.DB.Raw(sql).Scan(&tables).Error; err != nil {
		return nil, err
	}

	return tables, nil
}

// GetTableColumns 获取表的字段信息
func GetTableColumns(tableName string) ([]ColumnInfo, error) {
	var columns []ColumnInfo

	sql := `
		SELECT 
			COLUMN_NAME as column_name,
			DATA_TYPE as data_type,
			COLUMN_TYPE as column_type,
			IFNULL(COLUMN_COMMENT, '') as column_comment,
			COLUMN_KEY as column_key,
			IS_NULLABLE as is_nullable,
			COLUMN_DEFAULT as column_default
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`

	if err := global.DB.Raw(sql, tableName).Scan(&columns).Error; err != nil {
		return nil, err
	}

	return columns, nil
}

// ConvertToColumnConfig 将数据库字段信息转换为字段配置
func ConvertToColumnConfig(columns []ColumnInfo) []ColumnConfig {
	configs := make([]ColumnConfig, 0, len(columns))

	for _, col := range columns {
		config := ColumnConfig{
			ColumnName:    col.ColumnName,
			FieldName:     ToPascalCase(col.ColumnName),
			FieldType:     MapDBTypeToGoType(col.DataType, col.IsNullable == "YES"),
			JsonName:      col.ColumnName,
			TsType:        MapDBTypeToTsType(col.DataType),
			Comment:       col.ColumnComment,
			IsPrimaryKey:  col.ColumnKey == "PRI",
			IsRequired:    col.IsNullable == "NO" && col.ColumnKey != "PRI",
			IsSearchable:  isSearchableField(col.ColumnName, col.DataType),
			IsListVisible: isListVisibleField(col.ColumnName),
			IsFormVisible: isFormVisibleField(col.ColumnName),
			FormType:      inferFormType(col.ColumnName, col.DataType),
		}
		configs = append(configs, config)
	}

	return configs
}

// ToPascalCase 转换为大驼峰命名
func ToPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// ToCamelCase 转换为小驼峰命名
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if len(pascal) > 0 {
		return strings.ToLower(pascal[:1]) + pascal[1:]
	}
	return pascal
}

// ToSnakeCase 转换为下划线命名
func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// MapDBTypeToGoType 数据库类型转Go类型
func MapDBTypeToGoType(dbType string, nullable bool) string {
	dbType = strings.ToLower(dbType)

	typeMap := map[string]string{
		"int":        "int",
		"tinyint":    "int",
		"smallint":   "int",
		"mediumint":  "int",
		"bigint":     "int64",
		"float":      "float64",
		"double":     "float64",
		"decimal":    "float64",
		"char":       "string",
		"varchar":    "string",
		"text":       "string",
		"mediumtext": "string",
		"longtext":   "string",
		"datetime":   "time.Time",
		"timestamp":  "time.Time",
		"date":       "time.Time",
		"time":       "string",
		"json":       "string",
		"blob":       "[]byte",
	}

	if goType, ok := typeMap[dbType]; ok {
		return goType
	}
	return "string"
}

// MapDBTypeToTsType 数据库类型转TypeScript类型
func MapDBTypeToTsType(dbType string) string {
	dbType = strings.ToLower(dbType)

	typeMap := map[string]string{
		"int":        "number",
		"tinyint":    "number",
		"smallint":   "number",
		"mediumint":  "number",
		"bigint":     "number",
		"float":      "number",
		"double":     "number",
		"decimal":    "number",
		"char":       "string",
		"varchar":    "string",
		"text":       "string",
		"mediumtext": "string",
		"longtext":   "string",
		"datetime":   "string",
		"timestamp":  "string",
		"date":       "string",
		"time":       "string",
		"json":       "any",
		"blob":       "string",
	}

	if tsType, ok := typeMap[dbType]; ok {
		return tsType
	}
	return "string"
}

// isSearchableField 判断字段是否可搜索
func isSearchableField(name, dataType string) bool {
	// 主键、时间字段、大文本不默认搜索
	if name == "id" || name == "created_at" || name == "updated_at" || name == "deleted_at" {
		return false
	}
	if dataType == "text" || dataType == "mediumtext" || dataType == "longtext" || dataType == "blob" {
		return false
	}
	// 名称、标题等字段默认可搜索
	if strings.Contains(name, "name") || strings.Contains(name, "title") {
		return true
	}
	return false
}

// isListVisibleField 判断字段是否在列表显示
func isListVisibleField(name string) bool {
	// 这些字段默认不在列表显示
	hiddenFields := []string{"deleted_at", "password", "content", "description", "remark"}
	for _, f := range hiddenFields {
		if name == f {
			return false
		}
	}
	return true
}

// isFormVisibleField 判断字段是否在表单显示
func isFormVisibleField(name string) bool {
	// 这些字段默认不在表单显示
	hiddenFields := []string{"id", "created_at", "updated_at", "deleted_at"}
	for _, f := range hiddenFields {
		if name == f {
			return false
		}
	}
	return true
}

// inferFormType 推断表单组件类型
func inferFormType(name, dataType string) string {
	// 根据字段名推断
	if strings.Contains(name, "status") || strings.Contains(name, "type") {
		return "select"
	}
	if strings.Contains(name, "content") || strings.Contains(name, "description") || strings.Contains(name, "remark") {
		return "textarea"
	}
	if strings.Contains(name, "image") || strings.Contains(name, "avatar") || strings.Contains(name, "logo") {
		return "image"
	}
	if strings.Contains(name, "file") || strings.Contains(name, "attachment") {
		return "upload"
	}
	if strings.Contains(name, "enabled") || strings.Contains(name, "hidden") {
		return "switch"
	}

	// 根据数据类型推断
	switch dataType {
	case "text", "mediumtext", "longtext":
		return "textarea"
	case "datetime", "timestamp":
		return "datetime"
	case "date":
		return "date"
	case "int", "tinyint", "smallint", "mediumint", "bigint", "float", "double", "decimal":
		return "number"
	default:
		return "input"
	}
}
