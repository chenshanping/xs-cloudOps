package response

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func BindErrorMessage(err error, sample any) string {
	if err == nil {
		return "参数错误"
	}

	if message := bindValidationErrorMessage(err, sample); message != "" {
		return message
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return "请求体格式错误"
	}

	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		return "请求体格式错误"
	}

	message := strings.ToLower(strings.TrimSpace(err.Error()))
	switch {
	case message == "eof":
		return "请求体不能为空"
	case strings.Contains(message, "invalid character"),
		strings.Contains(message, "cannot unmarshal"),
		strings.Contains(message, "json"):
		return "请求体格式错误"
	default:
		return "参数错误"
	}
}

func bindValidationErrorMessage(err error, sample any) string {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) || len(validationErrs) == 0 {
		return ""
	}

	fieldMap := buildFieldLabelMap(sample)
	firstErr := validationErrs[0]
	label := fieldMap[firstErr.StructField()]
	if label == "" {
		label = firstErr.Field()
	}

	switch firstErr.Tag() {
	case "required":
		return label + "不能为空"
	case "min":
		return label + "至少 " + firstErr.Param() + " 位"
	case "max":
		return label + "最多 " + firstErr.Param() + " 位"
	case "email":
		return label + "格式不正确"
	case "oneof":
		return label + "取值无效"
	default:
		return label + "格式不正确"
	}
}

func buildFieldLabelMap(sample any) map[string]string {
	labels := map[string]string{}
	if sample == nil {
		return labels
	}

	typ := reflect.TypeOf(sample)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return labels
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		label := strings.TrimSpace(field.Tag.Get("comment"))
		if label == "" {
			label = strings.TrimSpace(strings.Split(field.Tag.Get("json"), ",")[0])
		}
		if label == "" {
			label = field.Name
		}
		labels[field.Name] = label
	}

	return labels
}
