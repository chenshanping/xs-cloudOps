package echart

import (
	"time"

	"server/global"
	"server/model"
)

type EchartService struct{}

var Default = &EchartService{}

// GetUserRoleStats 获取用户角色统计数据
func (s *EchartService) GetUserRoleStats() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := global.DB.Model(&model.SysRole{}).
		Select("sys_role.name as name, COUNT(sys_user_role.sys_user_id) as value").
		Joins("LEFT JOIN sys_user_role ON sys_role.id = sys_user_role.sys_role_id").
		Group("sys_role.id").
		Find(&results).Error
	return results, err
}

// GetUserStatusStats 获取用户状态统计
func (s *EchartService) GetUserStatusStats() ([]map[string]interface{}, error) {
	var activeCount int64
	var inactiveCount int64

	global.DB.Model(&model.SysUser{}).Where("status = 1").Count(&activeCount)
	global.DB.Model(&model.SysUser{}).Where("status = 0").Count(&inactiveCount)

	return []map[string]interface{}{
		{"name": "启用", "value": activeCount},
		{"name": "禁用", "value": inactiveCount},
	}, nil
}

// GetRoleStatusStats 获取角色状态统计
func (s *EchartService) GetRoleStatusStats() ([]map[string]interface{}, error) {
	var activeCount int64
	var inactiveCount int64

	global.DB.Model(&model.SysRole{}).Where("status = 1").Count(&activeCount)
	global.DB.Model(&model.SysRole{}).Where("status = 0").Count(&inactiveCount)

	return []map[string]interface{}{
		{"name": "启用", "value": activeCount},
		{"name": "禁用", "value": inactiveCount},
	}, nil
}

// GetUserRegisterTrend 获取用户注册趋势（最近30天）
func (s *EchartService) GetUserRegisterTrend() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	err := global.DB.Model(&model.SysUser{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", thirtyDaysAgo).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&results).Error
	return results, err
}
