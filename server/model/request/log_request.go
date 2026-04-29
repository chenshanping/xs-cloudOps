package request

// 日志列表请求
type LogListRequest struct {
	PageRequest
	Username     string `json:"username" form:"username" comment:"用户名"`
	Method       string `json:"method" form:"method" comment:"请求方法"`
	Path         string `json:"path" form:"path" comment:"请求路径"`
	Group        string `json:"group" form:"group" comment:"路由分组"`
	Summary      string `json:"summary" form:"summary" comment:"路由描述"`
	Status       *int   `json:"status" form:"status" comment:"HTTP状态码"`
	BusinessCode *int   `json:"business_code" form:"business_code" comment:"业务状态码"`
	StartTime    string `json:"start_time" form:"start_time" comment:"开始时间"`
	EndTime      string `json:"end_time" form:"end_time" comment:"结束时间"`
	SortField    string `json:"sort_field" form:"sort_field" comment:"排序字段"`
	SortOrder    string `json:"sort_order" form:"sort_order" comment:"排序方式(ascend/descend)"`
}
