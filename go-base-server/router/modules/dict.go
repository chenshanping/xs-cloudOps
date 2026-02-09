package modules

import (
	"github.com/gin-gonic/gin"

	v1 "go-base-server/api/v1"
	"go-base-server/model/request"
	"go-base-server/router/registry"
)

func init() {
	RegisterModule(&DictModule{})
}

type DictModule struct{}

func (m *DictModule) Name() string {
	return "字典管理"
}

func (m *DictModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 公开接口：根据类型获取字典数据（供下拉框等使用）
	R(rg, "GET", "/dict/type/:type", m.Name(), "获取字典数据", v1.Dict.GetDictDataByType)
}

func (m *DictModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	// 字典类型
	R(rg, "GET", "/dict/types", m.Name(), "字典类型列表", v1.Dict.GetDictTypeList, registry.WithAuth())
	R(rg, "GET", "/dict/types/all", m.Name(), "所有字典类型", v1.Dict.GetAllDictTypes, registry.WithAuth())
	R(rg, "GET", "/dict/types/:id", m.Name(), "字典类型详情", v1.Dict.GetDictType, registry.WithAuth())
	R(rg, "POST", "/dict/types", m.Name(), "创建字典类型", v1.Dict.CreateDictType,
		registry.WithAuth(), registry.WithRequest(request.CreateDictTypeRequest{}))
	R(rg, "PUT", "/dict/types/:id", m.Name(), "更新字典类型", v1.Dict.UpdateDictType,
		registry.WithAuth(), registry.WithRequest(request.UpdateDictTypeRequest{}))
	R(rg, "DELETE", "/dict/types/:id", m.Name(), "删除字典类型", v1.Dict.DeleteDictType, registry.WithAuth())

	// 字典数据
	R(rg, "GET", "/dict/data", m.Name(), "字典数据列表", v1.Dict.GetDictDataList, registry.WithAuth())
	R(rg, "GET", "/dict/data/:id", m.Name(), "字典数据详情", v1.Dict.GetDictData, registry.WithAuth())
	R(rg, "POST", "/dict/data", m.Name(), "创建字典数据", v1.Dict.CreateDictData,
		registry.WithAuth(), registry.WithRequest(request.CreateDictDataRequest{}))
	R(rg, "PUT", "/dict/data/:id", m.Name(), "更新字典数据", v1.Dict.UpdateDictData,
		registry.WithAuth(), registry.WithRequest(request.UpdateDictDataRequest{}))
	R(rg, "DELETE", "/dict/data/:id", m.Name(), "删除字典数据", v1.Dict.DeleteDictData, registry.WithAuth())
}
