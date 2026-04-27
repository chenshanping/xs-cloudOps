package request

// 创建API请求
type CreateApiRequest struct {
	Path        string `json:"path" binding:"required" comment:"API路径"`
	Method      string `json:"method" binding:"required" comment:"请求方法"`
	Group       string `json:"group" comment:"API分组"`
	Description string `json:"description" comment:"描述"`
}

// 更新API请求
type UpdateApiRequest struct {
	Path        string `json:"path" comment:"API路径"`
	Method      string `json:"method" comment:"请求方法"`
	Group       string `json:"group" comment:"API分组"`
	Description string `json:"description" comment:"描述"`
}

// API列表请求
type ApiListRequest struct {
	PageRequest
	Path   string `json:"path" form:"path" comment:"API路径"`
	Method string `json:"method" form:"method" comment:"请求方法"`
	Group  string `json:"group" form:"group" comment:"API分组"`
}
