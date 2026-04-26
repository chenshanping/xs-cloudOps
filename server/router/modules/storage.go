package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/router/registry"
)

func init() {
	RegisterModule(&StorageModule{})
}

type StorageModule struct{}

func (m *StorageModule) Name() string {
	return "存储管理"
}

func (m *StorageModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *StorageModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/storages", m.Name(), "存储配置列表", v1.Storage.GetStorageList, registry.WithAuth())
	R(rg, "GET", "/storages/:id", m.Name(), "存储配置详情", v1.Storage.GetStorage, registry.WithAuth())
	R(rg, "POST", "/storages", m.Name(), "创建存储配置", v1.Storage.CreateStorage, registry.WithAuth())
	R(rg, "PUT", "/storages/:id", m.Name(), "更新存储配置", v1.Storage.UpdateStorage, registry.WithAuth())
	R(rg, "DELETE", "/storages/:id", m.Name(), "删除存储配置", v1.Storage.DeleteStorage, registry.WithAuth())
	R(rg, "PUT", "/storages/:id/default", m.Name(), "设置默认存储", v1.Storage.SetDefaultStorage, registry.WithAuth())
	R(rg, "POST", "/storages/test", m.Name(), "测试存储配置", v1.Storage.TestStorage, registry.WithAuth())
}
