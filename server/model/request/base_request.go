package request

// 分页请求
type PageRequest struct {
	Page     int `json:"page" form:"page" comment:"页码"`
	PageSize int `json:"page_size" form:"page_size" comment:"每页数量"`
}

func (p *PageRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return (p.Page - 1) * p.PageSize
}
