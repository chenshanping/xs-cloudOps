package modules

import (
	"github.com/gin-gonic/gin"

	v1 "go-base-server/api/v1"
	"go-base-server/model/request"
	"go-base-server/router/registry"
)

func init() {
	RegisterModule(&ProductTypeModule{})
}

type ProductTypeModule struct{}

func (m *ProductTypeModule) Name() string {
	return "产品类型"
}

func (m *ProductTypeModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *ProductTypeModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/productType", m.Name(), "产品类型列表", v1.ProductType.GetProductTypeList,
		registry.WithAuth(), registry.WithRequest(request.ProductTypeListRequest{}))
	R(rg, "GET", "/productType/options", m.Name(), "产品类型选项", v1.ProductType.GetProductTypeOptions, registry.WithAuth())
	R(rg, "GET", "/productType/:id", m.Name(), "产品类型详情", v1.ProductType.GetProductType, registry.WithAuth())
	R(rg, "POST", "/productType", m.Name(), "创建产品类型", v1.ProductType.CreateProductType,
		registry.WithAuth(), registry.WithRequest(request.CreateProductTypeRequest{}))
	R(rg, "PUT", "/productType/:id", m.Name(), "更新产品类型", v1.ProductType.UpdateProductType,
		registry.WithAuth(), registry.WithRequest(request.UpdateProductTypeRequest{}))
	R(rg, "DELETE", "/productType/:id", m.Name(), "删除产品类型", v1.ProductType.DeleteProductType, registry.WithAuth())
	R(rg, "DELETE", "/productType/batch", m.Name(), "批量删除产品类型", v1.ProductType.BatchDeleteProductType,
		registry.WithAuth(), registry.WithRequest(request.BatchDeleteProductTypeRequest{}))
	// 导入导出
	R(rg, "GET", "/productType/export", m.Name(), "导出产品类型", v1.ProductType.ExportProductType, registry.WithAuth())
	R(rg, "POST", "/productType/import", m.Name(), "导入产品类型", v1.ProductType.ImportProductType, registry.WithAuth())
	R(rg, "GET", "/productType/template", m.Name(), "下载导入模板", v1.ProductType.DownloadTemplateProductType, registry.WithAuth())
}
