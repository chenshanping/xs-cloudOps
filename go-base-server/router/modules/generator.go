package modules

import (
	"github.com/gin-gonic/gin"

	v1 "go-base-server/api/v1"
	"go-base-server/router/registry"
)

func init() {
	RegisterModule(&GeneratorModule{})
}

type GeneratorModule struct{}

func (m *GeneratorModule) Name() string {
	return "代码生成"
}

func (m *GeneratorModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *GeneratorModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/generator/tables", m.Name(), "获取数据库表列表", v1.GeneratorAPI.GetTables, registry.WithAuth())
	R(rg, "GET", "/generator/tables/:name/columns", m.Name(), "获取表字段信息", v1.GeneratorAPI.GetTableColumns, registry.WithAuth())
	R(rg, "POST", "/generator/preview", m.Name(), "预览生成代码", v1.GeneratorAPI.Preview, registry.WithAuth())
	R(rg, "POST", "/generator/generate", m.Name(), "生成代码", v1.GeneratorAPI.Generate, registry.WithAuth())
	R(rg, "GET", "/generator/modules", m.Name(), "获取已生成模块", v1.GeneratorAPI.GetGeneratedModules, registry.WithAuth())
	R(rg, "DELETE", "/generator/modules/:name", m.Name(), "删除已生成模块", v1.GeneratorAPI.DeleteModule, registry.WithAuth())
	// 配置保存/获取/删除
	R(rg, "POST", "/generator/configs", m.Name(), "新增配置", v1.GeneratorAPI.SaveConfig, registry.WithAuth())
	R(rg, "PUT", "/generator/configs/:id", m.Name(), "更新配置", v1.GeneratorAPI.UpdateConfig, registry.WithAuth())
	R(rg, "GET", "/generator/configs", m.Name(), "获取配置列表", v1.GeneratorAPI.GetConfigs, registry.WithAuth())
	R(rg, "GET", "/generator/configs/:id", m.Name(), "获取配置详情", v1.GeneratorAPI.GetConfig, registry.WithAuth())
	R(rg, "DELETE", "/generator/configs/:id", m.Name(), "删除配置", v1.GeneratorAPI.DeleteConfig, registry.WithAuth())
	// 执行SQL
	R(rg, "POST", "/generator/execute-sql", m.Name(), "执行建表SQL", v1.GeneratorAPI.ExecuteSQL, registry.WithAuth())
	// 导入导出
	R(rg, "GET", "/generator/configs/:id/export", m.Name(), "导出配置", v1.GeneratorAPI.ExportConfig, registry.WithAuth())
	R(rg, "POST", "/generator/configs/import", m.Name(), "导入配置", v1.GeneratorAPI.ImportConfig, registry.WithAuth())
	R(rg, "POST", "/generator/configs/import-preview", m.Name(), "预览导入配置", v1.GeneratorAPI.ImportConfigPreview, registry.WithAuth())
}
