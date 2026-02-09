package v1

import (
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

// SaveConfig 保存配置（不生成代码）
func (a *GeneratorApi) SaveConfig(c *gin.Context) {
	var config generator.GeneratorConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if _, err := service.Generator.SaveConfig(&config); err != nil {
		response.Fail(c, "保存失败: "+err.Error())
		return
	}
	response.OkWithMessage(c, "保存成功")
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
