package modules

import (
	"github.com/gin-gonic/gin"

	v1 "go-base-server/api/v1"
	"go-base-server/model/request"
	"go-base-server/router/registry"
)

func init() {
	RegisterModule(&ProductModule{})
}

type ProductModule struct{}

func (m *ProductModule) Name() string {
	return "产品信息"
}

func (m *ProductModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *ProductModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/product", m.Name(), "产品信息列表", v1.Product.GetProductList,
		registry.WithAuth(), registry.WithRequest(request.ProductListRequest{}))
	R(rg, "GET", "/product/options", m.Name(), "产品信息选项", v1.Product.GetProductOptions, registry.WithAuth())
	R(rg, "GET", "/product/:id", m.Name(), "产品信息详情", v1.Product.GetProduct, registry.WithAuth())
	R(rg, "POST", "/product", m.Name(), "创建产品信息", v1.Product.CreateProduct,
		registry.WithAuth(), registry.WithRequest(request.CreateProductRequest{}))
	R(rg, "PUT", "/product/:id", m.Name(), "更新产品信息", v1.Product.UpdateProduct,
		registry.WithAuth(), registry.WithRequest(request.UpdateProductRequest{}))
	R(rg, "DELETE", "/product/:id", m.Name(), "删除产品信息", v1.Product.DeleteProduct, registry.WithAuth())
	R(rg, "DELETE", "/product/batch", m.Name(), "批量删除产品信息", v1.Product.BatchDeleteProduct,
		registry.WithAuth(), registry.WithRequest(request.BatchDeleteProductRequest{}))
	// 导入导出
	R(rg, "GET", "/product/export", m.Name(), "导出产品信息", v1.Product.ExportProduct, registry.WithAuth())
	R(rg, "POST", "/product/import", m.Name(), "导入产品信息", v1.Product.ImportProduct, registry.WithAuth())
	R(rg, "GET", "/product/template", m.Name(), "下载导入模板", v1.Product.DownloadTemplateProduct, registry.WithAuth())
	// 统计接口
	R(rg, "GET", "/product/stats/type_id", m.Name(), "产品信息产品类型统计", v1.Product.GetProductStatsTypeId, registry.WithAuth())
	R(rg, "GET", "/product/stats/status", m.Name(), "产品信息产品状态统计", v1.Product.GetProductStatsStatus, registry.WithAuth())
	R(rg, "GET", "/product/stats/trend", m.Name(), "产品信息趋势统计", v1.Product.GetProductTrendStats, registry.WithAuth())
}
