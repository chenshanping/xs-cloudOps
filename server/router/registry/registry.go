package registry

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// FieldInfo 字段信息
type FieldInfo struct {
	Name        string `json:"name"`        // 字段名(json tag)
	Type        string `json:"type"`        // 字段类型
	Description string `json:"description"` // 字段描述(来自 comment tag 或结构体注释)
	Required    bool   `json:"required"`    // 是否必填
	In          string `json:"in"`          // 参数位置: query, body, path
}

// RouteInfo 路由元信息
type RouteInfo struct {
	Method      string          // HTTP 方法
	Path        string          // 路由路径
	Group       string          // 分组名称
	Summary     string          // 描述
	Description string          // 详细描述
	Request     interface{}     // 请求参数结构体
	Response    interface{}     // 响应结构体
	NeedAuth    bool            // 是否需要认证
	Handler     gin.HandlerFunc // 处理函数
}

// RouteOption 路由配置选项
type RouteOption func(*RouteInfo)

// WithRequest 设置请求参数
func WithRequest(req interface{}) RouteOption {
	return func(r *RouteInfo) {
		r.Request = req
	}
}

// WithResponse 设置响应参数
func WithResponse(resp interface{}) RouteOption {
	return func(r *RouteInfo) {
		r.Response = resp
	}
}

// WithDescription 设置详细描述
func WithDescription(desc string) RouteOption {
	return func(r *RouteInfo) {
		r.Description = desc
	}
}

// WithAuth 设置需要认证
func WithAuth() RouteOption {
	return func(r *RouteInfo) {
		r.NeedAuth = true
	}
}

// 存储所有注册的路由信息
var registeredRoutes []RouteInfo

// Register 注册路由信息（兼容旧版本）
func Register(method, path, group, summary string, handler gin.HandlerFunc, opts ...RouteOption) {
	route := RouteInfo{
		Method:  method,
		Path:    path,
		Group:   group,
		Summary: summary,
		Handler: handler,
	}
	// 应用选项
	for _, opt := range opts {
		opt(&route)
	}
	registeredRoutes = append(registeredRoutes, route)
}

// GetAllRoutes 获取所有已注册的路由信息
func GetAllRoutes() []RouteInfo {
	return registeredRoutes
}

// GetTypeName 获取类型名称
func GetTypeName(i interface{}) string {
	if i == nil {
		return ""
	}
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
