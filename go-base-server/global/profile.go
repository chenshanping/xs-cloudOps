package global

import (
	"reflect"
	"strings"
	"sync"
)

// FieldConfig 字段配置
type FieldConfig struct {
	Key      string `json:"key"`      // 字段名
	Label    string `json:"label"`    // 中文标签
	Required bool   `json:"required"` // 是否必填
	Type     string `json:"type"`     // 字段类型: text, image, file, images, files, select
}

// UserProfile 用户身份信息
type UserProfile struct {
	Key         string        `json:"key"`         // 身份标识，如 "doctor", "merchant"
	Name        string        `json:"name"`        // 显示名称，如 "医生", "商家"
	Description string        `json:"description"` // 描述
	Icon        string        `json:"icon"`        // 图标（ant-design图标名）
	Data        interface{}   `json:"data"`        // 身份数据
	HasProfile  bool          `json:"has_profile"` // 当前用户是否有此身份
	IsComplete  bool          `json:"is_complete"` // 是否已完善（所有必填字段都有值）
	Fields      []FieldConfig `json:"fields"`      // 字段配置（标签、是否必填等）
}

// ProfileHandler 身份处理器接口
type ProfileHandler struct {
	Key         string        // 身份标识
	Name        string        // 显示名称
	Description string        // 描述
	Icon        string        // 图标
	Fields      []FieldConfig // 字段配置
	RoleCode    string        // 限定角色编码（为空表示不限制，任何有数据的用户都能看到）
	GetProfile  func(userID uint) (interface{}, error) // 获取用户身份数据
}

// ProfileRegistry 用户身份注册表
type ProfileRegistry struct {
	handlers map[string]*ProfileHandler
	mu       sync.RWMutex
}

// 全局身份注册表
var Profiles = &ProfileRegistry{
	handlers: make(map[string]*ProfileHandler),
}

// Register 注册身份处理器
func (r *ProfileRegistry) Register(handler *ProfileHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[handler.Key] = handler
}

// Unregister 注销身份处理器
func (r *ProfileRegistry) Unregister(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handlers, key)
}

// GetUserProfiles 获取用户所有身份信息
// userRoles: 用户拥有的角色编码列表
// 逻辑：所有已注册的身份都显示给登录用户，让用户可以申请任何身份
func (r *ProfileRegistry) GetUserProfiles(userID uint, userRoles []string) []UserProfile {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var profiles []UserProfile
	for _, handler := range r.handlers {
		if handler.GetProfile == nil {
			continue
		}

		// 获取用户身份数据
		data, err := handler.GetProfile(userID)
		if err != nil {
			// 查询出错，跳过
			continue
		}
		hasProfile := !isNilInterface(data)

		// 所有已注册的身份都显示，让用户可以申请
		profile := UserProfile{
			Key:         handler.Key,
			Name:        handler.Name,
			Description: handler.Description,
			Icon:        handler.Icon,
			Fields:      handler.Fields,
			Data:        data,
			HasProfile:  hasProfile,
			IsComplete:  hasProfile && r.checkComplete(data, handler.Fields),
		}

		profiles = append(profiles, profile)
	}

	return profiles
}

// checkComplete 检查所有必填字段是否都有值
func (r *ProfileRegistry) checkComplete(data interface{}, fields []FieldConfig) bool {
	if data == nil || len(fields) == 0 {
		return false
	}

	// 将 data 转换为 map
	dataMap := structToMap(data)
	if dataMap == nil {
		return false
	}

	// 检查所有必填字段
	for _, field := range fields {
		if !field.Required {
			continue
		}
		value, exists := dataMap[field.Key]
		if !exists || isEmptyValue(value) {
			return false
		}
	}
	return true
}

// structToMap 将结构体转换为 map
func structToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		// 取 json tag 的第一部分（可能有 omitempty 等）
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			jsonTag = jsonTag[:idx]
		}
		result[jsonTag] = v.Field(i).Interface()
	}
	return result
}

// isNilInterface 检查接口值是否为 nil（处理 (*T)(nil) 的情况）
func isNilInterface(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Ptr && rv.IsNil()
}

// isEmptyValue 检查值是否为空
func isEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(val).Int() == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(val).Uint() == 0
	case float32, float64:
		return reflect.ValueOf(val).Float() == 0
	case bool:
		return !val
	default:
		rv := reflect.ValueOf(v)
		return rv.IsZero()
	}
}

// GetRegisteredProfiles 获取所有已注册的身份类型（不含用户数据）
func (r *ProfileRegistry) GetRegisteredProfiles() []UserProfile {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var profiles []UserProfile
	for _, handler := range r.handlers {
		profiles = append(profiles, UserProfile{
			Key:         handler.Key,
			Name:        handler.Name,
			Description: handler.Description,
			Icon:        handler.Icon,
			Fields:      handler.Fields,
		})
	}

	return profiles
}

// HasHandler 检查是否存在指定身份处理器
func (r *ProfileRegistry) HasHandler(key string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.handlers[key]
	return exists
}

// GetUserBoundProfiles 检查用户绑定了哪些身份，返回身份名称列表
func (r *ProfileRegistry) GetUserBoundProfiles(userID uint) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var boundProfiles []string
	for _, handler := range r.handlers {
		if handler.GetProfile == nil {
			continue
		}
		data, err := handler.GetProfile(userID)
		if err != nil {
			continue
		}
		if !isNilInterface(data) {
			boundProfiles = append(boundProfiles, handler.Name)
		}
	}
	return boundProfiles
}
