package service

import (
	"time"

	"go-base-server/global"
	"go-base-server/model"
)

type EchartService struct{}

var Echart = new(EchartService)

// ChartItem 图表数据项
type ChartItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// TrendItem 趋势数据项
type TrendItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// dbTrendItem 数据库查询结果临时结构
type dbTrendItem struct {
	Date  time.Time
	Count int64
}

// 获取用户角色占比统计
func (s *EchartService) GetUserRoleStats() ([]ChartItem, error) {
	var result []ChartItem

	// 查询每个角色的用户数量
	err := global.DB.Table("sys_role").
		Select("sys_role.name as name, COUNT(DISTINCT sys_user_role.sys_user_id) as value").
		Joins("LEFT JOIN sys_user_role ON sys_role.id = sys_user_role.sys_role_id").
		Where("sys_role.status = ?", 1).
		Group("sys_role.id, sys_role.name").
		Order("value DESC").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// 如果没有数据，返回空数组
	if result == nil {
		result = []ChartItem{}
	}

	return result, nil
}

// 获取用户状态统计
func (s *EchartService) GetUserStatusStats() ([]ChartItem, error) {
	var result []ChartItem

	// 统计正常用户
	var activeCount int64
	if err := global.DB.Model(&model.SysUser{}).Where("status = ?", 1).Count(&activeCount).Error; err != nil {
		return nil, err
	}
	result = append(result, ChartItem{Name: "正常", Value: activeCount})

	// 统计禁用用户
	var disabledCount int64
	if err := global.DB.Model(&model.SysUser{}).Where("status = ?", 0).Count(&disabledCount).Error; err != nil {
		return nil, err
	}
	result = append(result, ChartItem{Name: "禁用", Value: disabledCount})

	return result, nil
}

// 获取角色状态统计
func (s *EchartService) GetRoleStatusStats() ([]ChartItem, error) {
	var result []ChartItem

	// 统计启用角色
	var activeCount int64
	if err := global.DB.Model(&model.SysRole{}).Where("status = ?", 1).Count(&activeCount).Error; err != nil {
		return nil, err
	}
	result = append(result, ChartItem{Name: "启用", Value: activeCount})

	// 统计禁用角色
	var disabledCount int64
	if err := global.DB.Model(&model.SysRole{}).Where("status = ?", 0).Count(&disabledCount).Error; err != nil {
		return nil, err
	}
	result = append(result, ChartItem{Name: "禁用", Value: disabledCount})

	return result, nil
}

// 获取用户注册趋势（近30天）
func (s *EchartService) GetUserRegisterTrend() ([]TrendItem, error) {
	var dbResult []dbTrendItem

	// 获取近30天的日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -29)

	// 查询近30天每天的注册用户数
	// 使用 < nextDay 代替 <= endDate 23:59:59
	nextDay := endDate.AddDate(0, 0, 1)
	err := global.DB.Table("sys_user").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ? AND created_at < ?", startDate.Format("2006-01-02"), nextDay.Format("2006-01-02")).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dbResult).Error

	if err != nil {
		return nil, err
	}

	// 补全没有数据的日期
	dateMap := make(map[string]int64)
	for _, item := range dbResult {
		// 将 time.Time 格式化为统一的字符串格式
		dateStr := item.Date.Format("2006-01-02")
		dateMap[dateStr] = item.Count
	}

	var fullResult []TrendItem
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		count := dateMap[dateStr]
		fullResult = append(fullResult, TrendItem{
			Date:  dateStr,
			Count: count,
		})
	}

	return fullResult, nil
}
