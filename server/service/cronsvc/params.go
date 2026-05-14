package cronsvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"

	"gorm.io/datatypes"
)

func ValidateParams(raw []byte, schema map[string]ParamDefinition) (map[string]interface{}, datatypes.JSON, error) {
	input := map[string]interface{}{}
	if len(raw) > 0 && string(raw) != "null" {
		decoderInput := map[string]interface{}{}
		if err := json.Unmarshal(raw, &decoderInput); err != nil {
			return nil, nil, fmt.Errorf("任务参数不是合法JSON")
		}
		input = decoderInput
	}

	cleaned := make(map[string]interface{}, len(schema))
	for key, def := range schema {
		value, exists := input[key]
		if !exists || value == nil {
			if def.Default != nil {
				cleaned[key] = def.Default
				continue
			}
			if def.Required {
				return nil, nil, fmt.Errorf("缺少任务参数: %s", key)
			}
			continue
		}

		converted, err := convertParamValue(key, value, def)
		if err != nil {
			return nil, nil, err
		}
		cleaned[key] = converted
	}

	encoded, err := json.Marshal(cleaned)
	if err != nil {
		return nil, nil, err
	}
	return cleaned, datatypes.JSON(encoded), nil
}

func ParamsToMap(params datatypes.JSON) (map[string]interface{}, error) {
	if len(params) == 0 || string(params) == "null" {
		return map[string]interface{}{}, nil
	}
	var output map[string]interface{}
	if err := json.Unmarshal(params, &output); err != nil {
		return nil, err
	}
	return output, nil
}

func convertParamValue(key string, value interface{}, def ParamDefinition) (interface{}, error) {
	switch def.Type {
	case ParamTypeInt:
		intValue, err := toInt(value)
		if err != nil {
			return nil, fmt.Errorf("任务参数%s必须是整数", key)
		}
		if def.Min != nil && intValue < *def.Min {
			return nil, fmt.Errorf("任务参数%s不能小于%d", key, *def.Min)
		}
		if def.Max != nil && intValue > *def.Max {
			return nil, fmt.Errorf("任务参数%s不能大于%d", key, *def.Max)
		}
		return intValue, nil
	case ParamTypeString:
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("任务参数%s必须是字符串", key)
		}
		return strValue, nil
	case ParamTypeBool:
		boolValue, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("任务参数%s必须是布尔值", key)
		}
		return boolValue, nil
	default:
		return nil, fmt.Errorf("任务参数%s类型不支持", key)
	}
}

func toInt(value interface{}) (int, error) {
	switch typed := value.(type) {
	case int:
		return typed, nil
	case int64:
		return int(typed), nil
	case float64:
		if math.Trunc(typed) != typed {
			return 0, errors.New("not integer")
		}
		return int(typed), nil
	case json.Number:
		intValue, err := strconv.Atoi(typed.String())
		if err != nil {
			return 0, err
		}
		return intValue, nil
	default:
		return 0, errors.New("not integer")
	}
}

func getIntParam(params map[string]interface{}, key string, fallback int) int {
	if value, ok := params[key]; ok {
		if intValue, err := toInt(value); err == nil {
			return intValue
		}
	}
	return fallback
}
