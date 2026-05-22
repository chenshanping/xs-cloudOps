package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/model/request"
	"server/router/registry"
)

func init() {
	RegisterModule(&CmdbModule{})
}

type CmdbModule struct{}

func (m *CmdbModule) Name() string {
	return "CMDB管理"
}

func (m *CmdbModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/cmdb/terminal/ws", m.Name(), "SSH终端WebSocket接入", v1.CmdbTerminal.ServeWebSocket)
}

func (m *CmdbModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	R(rg, "GET", "/cmdb/host-groups", m.Name(), "主机分组列表", v1.Cmdb.GetHostGroups, registry.WithAuth())
	R(rg, "POST", "/cmdb/host-groups", m.Name(), "创建主机分组", v1.Cmdb.CreateHostGroup, registry.WithAuth(), registry.WithRequest(request.CreateCmdbHostGroupRequest{}))
	R(rg, "PUT", "/cmdb/host-groups/:id", m.Name(), "更新主机分组", v1.Cmdb.UpdateHostGroup, registry.WithAuth(), registry.WithRequest(request.UpdateCmdbHostGroupRequest{}))
	R(rg, "DELETE", "/cmdb/host-groups/:id", m.Name(), "删除主机分组", v1.Cmdb.DeleteHostGroup, registry.WithAuth())

	R(rg, "GET", "/cmdb/host-tags", m.Name(), "主机标签列表", v1.Cmdb.GetHostTags, registry.WithAuth())
	R(rg, "POST", "/cmdb/host-tags", m.Name(), "创建主机标签", v1.Cmdb.CreateHostTag, registry.WithAuth(), registry.WithRequest(request.CreateCmdbHostTagRequest{}))
	R(rg, "PUT", "/cmdb/host-tags/:id", m.Name(), "更新主机标签", v1.Cmdb.UpdateHostTag, registry.WithAuth(), registry.WithRequest(request.UpdateCmdbHostTagRequest{}))
	R(rg, "DELETE", "/cmdb/host-tags/:id", m.Name(), "删除主机标签", v1.Cmdb.DeleteHostTag, registry.WithAuth())

	R(rg, "GET", "/cmdb/ssh-credentials", m.Name(), "SSH凭据列表", v1.Cmdb.GetCredentials, registry.WithAuth())
	R(rg, "GET", "/cmdb/ssh-credentials/options", m.Name(), "SSH凭据选项", v1.Cmdb.GetCredentialOptions, registry.WithAuth())
	R(rg, "GET", "/cmdb/ssh-credentials/:id", m.Name(), "SSH凭据详情", v1.Cmdb.GetCredential, registry.WithAuth())
	R(rg, "POST", "/cmdb/ssh-credentials", m.Name(), "创建SSH凭据", v1.Cmdb.CreateCredential, registry.WithAuth(), registry.WithRequest(request.CreateCmdbCredentialRequest{}))
	R(rg, "PUT", "/cmdb/ssh-credentials/:id", m.Name(), "更新SSH凭据", v1.Cmdb.UpdateCredential, registry.WithAuth(), registry.WithRequest(request.UpdateCmdbCredentialRequest{}))
	R(rg, "DELETE", "/cmdb/ssh-credentials/:id", m.Name(), "删除SSH凭据", v1.Cmdb.DeleteCredential, registry.WithAuth())

	R(rg, "GET", "/cmdb/hosts", m.Name(), "主机列表", v1.Cmdb.GetHosts, registry.WithAuth())
	R(rg, "GET", "/cmdb/hosts/:id", m.Name(), "主机详情", v1.Cmdb.GetHost, registry.WithAuth())
	R(rg, "POST", "/cmdb/hosts", m.Name(), "创建主机", v1.Cmdb.CreateHost, registry.WithAuth(), registry.WithRequest(request.CreateCmdbHostRequest{}))
	R(rg, "PUT", "/cmdb/hosts/:id", m.Name(), "更新主机", v1.Cmdb.UpdateHost, registry.WithAuth(), registry.WithRequest(request.UpdateCmdbHostRequest{}))
	R(rg, "DELETE", "/cmdb/hosts/:id", m.Name(), "删除主机", v1.Cmdb.DeleteHost, registry.WithAuth())
	R(rg, "POST", "/cmdb/hosts/:id/verify", m.Name(), "校验主机", v1.Cmdb.VerifyHost, registry.WithAuth())
	R(rg, "GET", "/cmdb/hosts/import-template", m.Name(), "下载主机导入模板", v1.Cmdb.GetHostImportTemplate, registry.WithAuth())
	R(rg, "POST", "/cmdb/hosts/import", m.Name(), "导入主机", v1.Cmdb.ImportHosts, registry.WithAuth())

	R(rg, "POST", "/cmdb/terminal/sessions", m.Name(), "创建SSH终端会话", v1.CmdbTerminal.CreateSession, registry.WithAuth(), registry.WithRequest(request.CreateCmdbTerminalSessionRequest{}))
	R(rg, "GET", "/cmdb/terminal/sessions", m.Name(), "终端会话列表", v1.CmdbTerminal.GetSessions, registry.WithAuth())
	R(rg, "GET", "/cmdb/terminal/sessions/:id", m.Name(), "终端会话详情", v1.CmdbTerminal.GetSession, registry.WithAuth())
	R(rg, "GET", "/cmdb/terminal/sessions/:id/logs", m.Name(), "终端会话日志", v1.CmdbTerminal.GetSessionLogs, registry.WithAuth())
	R(rg, "POST", "/cmdb/terminal/sessions/:id/disconnect", m.Name(), "断开终端会话", v1.CmdbTerminal.DisconnectSession, registry.WithAuth())
	R(rg, "POST", "/cmdb/terminal/sessions/:id/force-disconnect", m.Name(), "强制断开终端会话", v1.CmdbTerminal.ForceDisconnectSession, registry.WithAuth())
}
