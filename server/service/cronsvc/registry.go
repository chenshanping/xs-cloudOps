package cronsvc

import (
	"context"
	"sort"
	"sync"
)

type ParamType string

const (
	ParamTypeInt    ParamType = "int"
	ParamTypeString ParamType = "string"
	ParamTypeBool   ParamType = "bool"
)

type ParamDefinition struct {
	Type        ParamType   `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description"`
	Min         *int        `json:"min,omitempty"`
	Max         *int        `json:"max,omitempty"`
}

type TaskFunc func(ctx context.Context, params map[string]interface{}) (string, error)

type TaskHandler struct {
	Code        string                     `json:"code"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	ParamSchema map[string]ParamDefinition `json:"param_schema"`
	Run         TaskFunc                   `json:"-"`
}

type RegisteredTask struct {
	Code        string                     `json:"code"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	ParamSchema map[string]ParamDefinition `json:"param_schema"`
}

var (
	registryMu sync.RWMutex
	registry   = make(map[string]TaskHandler)
)

func Register(handler TaskHandler) {
	if handler.Code == "" || handler.Run == nil {
		return
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[handler.Code] = handler
}

func Get(code string) (TaskHandler, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	handler, ok := registry[code]
	return handler, ok
}

func ListRegisteredTasks() []RegisteredTask {
	registryMu.RLock()
	defer registryMu.RUnlock()
	items := make([]RegisteredTask, 0, len(registry))
	for _, handler := range registry {
		items = append(items, RegisteredTask{
			Code:        handler.Code,
			Name:        handler.Name,
			Description: handler.Description,
			ParamSchema: handler.ParamSchema,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Code < items[j].Code
	})
	return items
}

func init() {
	Register(TaskHandler{
		Code:        "cleanup_login_logs",
		Name:        "清理登录日志",
		Description: "按保留天数分批清理过期登录日志",
		ParamSchema: cleanupLogParamSchema(),
		Run:         cleanupLoginLogs,
	})
	Register(TaskHandler{
		Code:        "cleanup_operation_logs",
		Name:        "清理操作日志",
		Description: "按保留天数分批清理过期操作日志",
		ParamSchema: cleanupLogParamSchema(),
		Run:         cleanupOperationLogs,
	})
}
