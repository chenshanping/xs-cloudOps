package registry

import (
	"reflect"
	"strings"
)

type DefinitionField struct {
	Name        string
	Description string
	Required    bool
	ParamType   string
	SwaggerType string
}

// ParseStructFields 解析结构体字段信息，供 API 同步使用。
func ParseStructFields(obj interface{}, paramIn string) []FieldInfo {
	fields := ParseDefinitionFields(obj)
	if len(fields) == 0 {
		return nil
	}

	result := make([]FieldInfo, 0, len(fields))
	for _, field := range fields {
		result = append(result, FieldInfo{
			Name:        field.Name,
			Type:        field.ParamType,
			Description: field.Description,
			Required:    field.Required,
			In:          paramIn,
		})
	}
	return result
}

// ParseDefinitionFields 解析结构体字段元数据，供 Swagger 和 API 同步共用。
func ParseDefinitionFields(obj interface{}) []DefinitionField {
	if obj == nil {
		return nil
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	var fields []DefinitionField
	parseDefinitionFieldsRecursive(t, &fields)
	return fields
}

func parseDefinitionFieldsRecursive(t reflect.Type, fields *[]DefinitionField) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			embeddedType := field.Type
			if embeddedType.Kind() == reflect.Ptr {
				embeddedType = embeddedType.Elem()
			}
			if embeddedType.Kind() == reflect.Struct {
				parseDefinitionFieldsRecursive(embeddedType, fields)
			}
			continue
		}

		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = field.Tag.Get("form")
			if name == "" || name == "-" {
				continue
			}
		}
		name = strings.Split(name, ",")[0]

		description := field.Tag.Get("comment")
		if description == "" {
			description = field.Tag.Get("label")
		}
		if description == "" {
			description = field.Name
		}

		required := strings.Contains(field.Tag.Get("binding"), "required")
		*fields = append(*fields, DefinitionField{
			Name:        name,
			Description: description,
			Required:    required,
			ParamType:   goTypeToParamString(field.Type),
			SwaggerType: goTypeToSwaggerString(field.Type),
		})
	}
}

func goTypeToParamString(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		elemType := goTypeToParamString(t.Elem())
		return "array[" + elemType + "]"
	case reflect.Ptr:
		return goTypeToParamString(t.Elem())
	default:
		return "object"
	}
}

func goTypeToSwaggerString(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Ptr:
		return goTypeToSwaggerString(t.Elem())
	default:
		return "object"
	}
}
