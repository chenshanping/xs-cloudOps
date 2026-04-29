package logsvc

import (
	"server/global"
	"server/model"
	"server/model/request"
)

type LogService struct{}

var Default = &LogService{}

func (s *LogService) GetOperationLogList(req *request.LogListRequest) ([]model.SysOperationLog, int64, error) {
	var logs []model.SysOperationLog
	var total int64

	db := global.DB.Model(&model.SysOperationLog{})

	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
	}
	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.Group != "" {
		db = db.Where("`group` = ?", req.Group)
	}
	if req.Summary != "" {
		db = db.Where("summary LIKE ?", "%"+req.Summary+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.BusinessCode != nil {
		db = db.Where("business_code = ?", *req.BusinessCode)
	}
	if req.StartTime != "" {
		db = db.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		db = db.Where("created_at <= ?", req.EndTime)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	orderStr := "id DESC"
	if req.SortField != "" {
		order := "ASC"
		if req.SortOrder == "descend" {
			order = "DESC"
		}
		orderStr = req.SortField + " " + order
	}
	if err := db.Offset(offset).Limit(req.PageSize).Order(orderStr).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (s *LogService) GetLoginLogList(req *request.LogListRequest) ([]model.SysLoginLog, int64, error) {
	var logs []model.SysLoginLog
	var total int64
	db := global.DB.Model(&model.SysLoginLog{})
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.StartTime != "" {
		db = db.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		db = db.Where("created_at <= ?", req.EndTime)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := req.GetOffset()
	if err := db.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

type RouteGroupCount struct {
	Group string `json:"group"`
	Count int64  `json:"count"`
}

func (s *LogService) GetRouteGroups() ([]RouteGroupCount, error) {
	var groups []RouteGroupCount
	err := global.DB.Model(&model.SysOperationLog{}).
		Select("`group`, COUNT(*) as count").
		Group("`group`").
		Order("count DESC").
		Find(&groups).Error
	return groups, err
}

// CreateLoginLog 创建登录日志
func (s *LogService) CreateLoginLog(log *model.SysLoginLog) error {
	return global.DB.Create(log).Error
}
