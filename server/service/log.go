package service

import (
	"server/global"
	"server/model"
	"server/model/request"
)

type LogService struct{}

var Log = new(LogService)

// 获取操作日志列表
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

// 获取登录日志列表
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
	if err := db.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// 记录登录日志
func (s *LogService) CreateLoginLog(log *model.SysLoginLog) error {
	return global.DB.Create(log).Error
}

// 路由分组统计
type RouteGroupCount struct {
	Group string `json:"group"`
	Count int64  `json:"count"`
}

// 获取路由分组列表（从数据库统计）
func (s *LogService) GetRouteGroups() ([]RouteGroupCount, error) {
	var groups []RouteGroupCount
	err := global.DB.Model(&model.SysOperationLog{}).
		Select("`group`, COUNT(*) as count").
		Group("`group`").
		Order("count DESC").
		Find(&groups).Error
	return groups, err
}

// 获取慢查询日志列表
func (s *LogService) GetSlowLogList(req *request.SlowLogListRequest) ([]model.SysSlowLog, int64, error) {
	var logs []model.SysSlowLog
	var total int64

	db := global.DB.Model(&model.SysSlowLog{})

	if req.SQL != "" {
		db = db.Where("sql LIKE ?", "%"+req.SQL+"%")
	}
	if req.MinLatency > 0 {
		db = db.Where("latency >= ?", req.MinLatency)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	if err := db.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
