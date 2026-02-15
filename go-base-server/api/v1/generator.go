package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"go-base-server/generator"
	"go-base-server/model/response"
	"go-base-server/service"
)

type GeneratorApi struct{}

var GeneratorAPI = new(GeneratorApi)

// GetTables 获取数据库表列表
func (a *GeneratorApi) GetTables(c *gin.Context) {
	tables, err := service.Generator.GetTables()
	if err != nil {
		response.Fail(c, "获取表列表失败: "+err.Error())
		return
	}
	response.OkWithData(c, tables)
}

// GetTableColumns 获取表字段信息
func (a *GeneratorApi) GetTableColumns(c *gin.Context) {
	tableName := c.Param("name")
	if tableName == "" {
		response.BadRequest(c, "表名不能为空")
		return
	}

	columns, err := service.Generator.GetTableColumns(tableName)
	if err != nil {
		response.Fail(c, "获取字段信息失败: "+err.Error())
		return
	}
	response.OkWithData(c, columns)
}

// Preview 预览生成的代码
func (a *GeneratorApi) Preview(c *gin.Context) {
	var config generator.GeneratorConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := service.Generator.Preview(&config)
	if err != nil {
		response.Fail(c, "预览失败: "+err.Error())
		return
	}
	response.OkWithData(c, result)
}

// Generate 生成代码
func (a *GeneratorApi) Generate(c *gin.Context) {
	var config generator.GeneratorConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := service.Generator.Generate(&config)
	if err != nil {
		response.Fail(c, "生成失败: "+err.Error())
		return
	}
	response.OkWithData(c, result)
}

// DeleteModule 删除模块
func (a *GeneratorApi) DeleteModule(c *gin.Context) {
	moduleName := c.Param("name")
	if moduleName == "" {
		response.BadRequest(c, "模块名不能为空")
		return
	}

	if err := service.Generator.DeleteModule(moduleName); err != nil {
		response.Fail(c, "删除失败: "+err.Error())
		return
	}
	response.OkWithMessage(c, "删除成功")
}

// GetGeneratedModules 获取已生成的模块列表
func (a *GeneratorApi) GetGeneratedModules(c *gin.Context) {
	modules, err := service.Generator.GetGeneratedModules()
	if err != nil {
		response.Fail(c, "获取模块列表失败: "+err.Error())
		return
	}
	response.OkWithData(c, modules)
}

// SaveConfig 新增配置（不生成代码）
func (a *GeneratorApi) SaveConfig(c *gin.Context) {
	var config generator.GeneratorConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// 新增时确保ID为0
	config.ID = 0

	if _, err := service.Generator.SaveConfig(&config); err != nil {
		response.Fail(c, "保存失败: "+err.Error())
		return
	}
	response.OkWithMessage(c, "保存成功")
}

// UpdateConfig 更新配置（不生成代码）
func (a *GeneratorApi) UpdateConfig(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "ID不能为空")
		return
	}

	var config generator.GeneratorConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 从 URL 获取 ID
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil || id == 0 {
		response.BadRequest(c, "ID无效")
		return
	}
	config.ID = id

	if _, err := service.Generator.SaveConfig(&config); err != nil {
		response.Fail(c, "更新失败: "+err.Error())
		return
	}
	response.OkWithMessage(c, "更新成功")
}

// GetConfigs 获取已保存的配置列表
func (a *GeneratorApi) GetConfigs(c *gin.Context) {
	list, err := service.Generator.GetConfigs()
	if err != nil {
		response.Fail(c, "获取配置列表失败: "+err.Error())
		return
	}
	response.OkWithData(c, list)
}

// GetConfig 获取单个配置详情
func (a *GeneratorApi) GetConfig(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "ID不能为空")
		return
	}

	config, err := service.Generator.GetConfig(id)
	if err != nil {
		response.Fail(c, "获取配置失败: "+err.Error())
		return
	}
	response.OkWithData(c, config)
}

// DeleteConfig 删除配置
func (a *GeneratorApi) DeleteConfig(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "ID不能为空")
		return
	}

	if err := service.Generator.DeleteConfig(id); err != nil {
		response.Fail(c, "删除失败: "+err.Error())
		return
	}
	response.OkWithMessage(c, "删除成功")
}

// ExecuteSQL 执行建表SQL
func (a *GeneratorApi) ExecuteSQL(c *gin.Context) {
	var req struct {
		SQL string `json:"sql"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Generator.ExecuteSQL(req.SQL); err != nil {
		response.Fail(c, "执行SQL失败: "+err.Error())
		return
	}
	response.OkWithMessage(c, "SQL执行成功")
}

// ExportConfig 导出配置为JSON
func (a *GeneratorApi) ExportConfig(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "ID不能为空")
		return
	}

	config, err := service.Generator.GetConfig(id)
	if err != nil {
		response.Fail(c, "获取配置失败: "+err.Error())
		return
	}

	// 导出为JSON字符串
	jsonStr, err := generator.ExportConfigToString(config)
	if err != nil {
		response.Fail(c, "导出失败: "+err.Error())
		return
	}

	// 设置响应头，触发下载
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s_config.json", config.ModuleName))
	c.String(200, jsonStr)
}

// ImportConfig 导入配置
func (a *GeneratorApi) ImportConfig(c *gin.Context) {
	var req struct {
		ConfigJSON string `json:"config_json" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 从JSON字符串导入配置
	config, err := generator.ImportConfigFromString(req.ConfigJSON)
	if err != nil {
		response.Fail(c, "导入失败: "+err.Error())
		return
	}

	// 保存导入的配置（ID设为0，作为新配置保存）
	config.ID = 0
	savedConfig, err := service.Generator.SaveConfig(config)
	if err != nil {
		response.Fail(c, "保存导入的配置失败: "+err.Error())
		return
	}

	response.OkWithData(c, savedConfig)
}

// ImportConfigPreview 预览导入的配置（不保存）
func (a *GeneratorApi) ImportConfigPreview(c *gin.Context) {
	var req struct {
		ConfigJSON string `json:"config_json" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 从JSON字符串导入配置
	config, err := generator.ImportConfigFromString(req.ConfigJSON)
	if err != nil {
		response.Fail(c, "解析失败: "+err.Error())
		return
	}

	// 返回解析后的配置供前端预览
	response.OkWithData(c, config)
}
