package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/model/request"
	"server/router/registry"
)

func init() {
	RegisterModule(&DeptModule{})
}

type DeptModule struct{}

func (m *DeptModule) Name() string {
	return "部门管理"
}

func (m *DeptModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *DeptModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/depts/tree", m.Name(), "部门树", v1.Dept.GetDeptTree, registry.WithAuth())
	R(rg, "GET", "/depts/manageable-tree", m.Name(), "可管理部门树", v1.Dept.GetManageableDeptTree, registry.WithAuth())
	R(rg, "GET", "/depts/:id", m.Name(), "部门详情", v1.Dept.GetDept, registry.WithAuth())
	R(rg, "POST", "/depts", m.Name(), "创建部门", v1.Dept.CreateDept,
		registry.WithAuth(), registry.WithRequest(request.CreateDeptRequest{}))
	R(rg, "PUT", "/depts/:id", m.Name(), "更新部门", v1.Dept.UpdateDept,
		registry.WithAuth(), registry.WithRequest(request.UpdateDeptRequest{}))
	R(rg, "DELETE", "/depts/:id", m.Name(), "删除部门", v1.Dept.DeleteDept, registry.WithAuth())
}
