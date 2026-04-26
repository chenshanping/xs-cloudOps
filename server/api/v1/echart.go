package v1

import (
	"github.com/gin-gonic/gin"

	"server/model/response"
	"server/service"
)

type EchartApi struct{}

var Echart = new(EchartApi)

// 获取用户角色占比统计
func (a *EchartApi) GetUserRoleStats(c *gin.Context) {
	stats, err := service.Echart.GetUserRoleStats()
	if err != nil {
		response.Fail(c, "获取用户角色统计失败")
		return
	}

	response.OkWithData(c, stats)
}

// 获取用户状态统计
func (a *EchartApi) GetUserStatusStats(c *gin.Context) {
	stats, err := service.Echart.GetUserStatusStats()
	if err != nil {
		response.Fail(c, "获取用户状态统计失败")
		return
	}

	response.OkWithData(c, stats)
}

// 获取角色状态统计
func (a *EchartApi) GetRoleStatusStats(c *gin.Context) {
	stats, err := service.Echart.GetRoleStatusStats()
	if err != nil {
		response.Fail(c, "获取角色状态统计失败")
		return
	}

	response.OkWithData(c, stats)
}

// 获取用户注册趋势（近30天）
func (a *EchartApi) GetUserRegisterTrend(c *gin.Context) {
	stats, err := service.Echart.GetUserRegisterTrend()
	if err != nil {
		response.Fail(c, "获取用户注册趋势失败")
		return
	}

	response.OkWithData(c, stats)
}
