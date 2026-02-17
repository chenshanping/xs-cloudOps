package modules

import (
	"github.com/gin-gonic/gin"

	v1 "go-base-server/api/v1"
	"go-base-server/model/request"
	"go-base-server/router/registry"
)

func init() {
	RegisterModule(&{{.ModelName}}Module{})
}

type {{.ModelName}}Module struct{}

func (m *{{.ModelName}}Module) Name() string {
	return "{{.Description}}"
}

func (m *{{.ModelName}}Module) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *{{.ModelName}}Module) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/{{.RoutePath}}", m.Name(), "{{.Description}}列表", v1.{{.ModelName}}.Get{{.ModelName}}List,
		registry.WithAuth(), registry.WithRequest(request.{{.ModelName}}ListRequest{}))
	R(rg, "GET", "/{{.RoutePath}}/options", m.Name(), "{{.Description}}选项", v1.{{.ModelName}}.Get{{.ModelName}}Options, registry.WithAuth())
{{- if .HasCreatedBy}}
	R(rg, "GET", "/{{.RoutePath}}/creator/options", m.Name(), "{{.Description}}创建人选项", v1.{{.ModelName}}.Get{{.ModelName}}CreatorOptions, registry.WithAuth())
{{- end}}
	R(rg, "GET", "/{{.RoutePath}}/:id", m.Name(), "{{.Description}}详情", v1.{{.ModelName}}.Get{{.ModelName}}, registry.WithAuth())
	R(rg, "POST", "/{{.RoutePath}}", m.Name(), "创建{{.Description}}", v1.{{.ModelName}}.Create{{.ModelName}},
		registry.WithAuth(), registry.WithRequest(request.Create{{.ModelName}}Request{}))
	R(rg, "PUT", "/{{.RoutePath}}/:id", m.Name(), "更新{{.Description}}", v1.{{.ModelName}}.Update{{.ModelName}},
		registry.WithAuth(), registry.WithRequest(request.Update{{.ModelName}}Request{}))
	R(rg, "DELETE", "/{{.RoutePath}}/:id", m.Name(), "删除{{.Description}}", v1.{{.ModelName}}.Delete{{.ModelName}}, registry.WithAuth())
	R(rg, "DELETE", "/{{.RoutePath}}/batch", m.Name(), "批量删除{{.Description}}", v1.{{.ModelName}}.BatchDelete{{.ModelName}},
		registry.WithAuth(), registry.WithRequest(request.BatchDelete{{.ModelName}}Request{}))
{{- if .EnableImportExport}}
	// 导入导出
	R(rg, "GET", "/{{.RoutePath}}/export", m.Name(), "导出{{.Description}}", v1.{{.ModelName}}.Export{{.ModelName}}, registry.WithAuth())
	R(rg, "POST", "/{{.RoutePath}}/import", m.Name(), "导入{{.Description}}", v1.{{.ModelName}}.Import{{.ModelName}}, registry.WithAuth())
	R(rg, "GET", "/{{.RoutePath}}/template", m.Name(), "下载导入模板", v1.{{.ModelName}}.DownloadTemplate{{.ModelName}}, registry.WithAuth())
{{- end}}
{{- if .HasAudit}}
	R(rg, "POST", "/{{.RoutePath}}/:id/audit", m.Name(), "审批{{.Description}}", v1.{{.ModelName}}.Audit{{.ModelName}},
		registry.WithAuth(), registry.WithRequest(request.Audit{{.ModelName}}Request{}))
{{- end}}
{{- if .GenerateFrontendApi}}
	// 前台接口（私有路由，需要登录，但不做 created_by 过滤）
	R(rg, "GET", "/{{.RoutePath}}/frontend", m.Name(), "前台{{.Description}}列表", v1.{{.ModelName}}.GetFrontend{{.ModelName}}List,
		registry.WithAuth(), registry.WithRequest(request.Frontend{{.ModelName}}ListRequest{}))
	R(rg, "GET", "/{{.RoutePath}}/frontend/:id", m.Name(), "前台{{.Description}}详情", v1.{{.ModelName}}.GetFrontend{{.ModelName}}, registry.WithAuth())
{{- end}}
{{- if .LinkToUser}}
	// 当前用户接口（获取/保存自己的扩展信息）
	R(rg, "GET", "/{{.RoutePath}}/my", m.Name(), "获取我的{{.Description}}", v1.{{.ModelName}}.GetMy{{.ModelName}}, registry.WithAuth())
	R(rg, "POST", "/{{.RoutePath}}/my", m.Name(), "保存我的{{.Description}}", v1.{{.ModelName}}.SaveMy{{.ModelName}},
		registry.WithAuth(), registry.WithRequest(request.SaveMy{{.ModelName}}Request{}))
{{- end}}
{{- if .HasStats}}
	// 统计接口
{{- range $chart := .StatsCharts}}
	R(rg, "GET", "/{{$.RoutePath}}/stats/{{$chart.Column}}", m.Name(), "{{$.Description}}{{$chart.Title}}统计", v1.{{$.ModelName}}.Get{{$.ModelName}}Stats{{$chart.Field}}, registry.WithAuth())
{{- end}}
{{- if .HasStatsTrend}}
	R(rg, "GET", "/{{.RoutePath}}/stats/trend", m.Name(), "{{.Description}}趋势统计", v1.{{.ModelName}}.Get{{.ModelName}}TrendStats, registry.WithAuth())
{{- end}}
{{- end}}
}
