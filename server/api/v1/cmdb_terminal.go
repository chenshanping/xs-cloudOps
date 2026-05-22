package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"server/middleware"
	"server/model/request"
	"server/model/response"
	"server/service"
)

type CmdbTerminalApi struct{}

var CmdbTerminal = new(CmdbTerminalApi)

var terminalWSUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (a *CmdbTerminalApi) CreateSession(c *gin.Context) {
	var req request.CreateCmdbTerminalSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	item, err := service.CmdbTerminal.CreateSession(&req, middleware.GetUserID(c), middleware.GetUsername(c), c.ClientIP())
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, item)
}

func (a *CmdbTerminalApi) GetSessions(c *gin.Context) {
	var req request.CmdbTerminalSessionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	list, total, err := service.CmdbTerminal.ListSessions(&req, middleware.GetUserID(c), middleware.GetUserRoleCodes(c))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithPage(c, list, total, req.Page, req.PageSize)
}

func (a *CmdbTerminalApi) GetSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	item, err := service.CmdbTerminal.GetSession(uint(id), middleware.GetUserID(c), middleware.GetUserRoleCodes(c))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, item)
}

func (a *CmdbTerminalApi) GetSessionLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	var req request.CmdbTerminalLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	list, total, err := service.CmdbTerminal.ListLogs(uint(id), middleware.GetUserID(c), middleware.GetUserRoleCodes(c), &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithPage(c, list, total, req.Page, req.PageSize)
}

func (a *CmdbTerminalApi) DisconnectSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.CmdbTerminal.DisconnectSession(uint(id), middleware.GetUserID(c), middleware.GetUserRoleCodes(c), false); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "会话已断开")
}

func (a *CmdbTerminalApi) ForceDisconnectSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.CmdbTerminal.DisconnectSession(uint(id), middleware.GetUserID(c), middleware.GetUserRoleCodes(c), true); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "会话已强制断开")
}

func (a *CmdbTerminalApi) ServeWebSocket(c *gin.Context) {
	conn, err := terminalWSUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	if err := service.CmdbTerminal.ServeWebSocket(conn, c.Query("token")); err != nil {
		_ = conn.WriteJSON(map[string]string{
			"type":    "status",
			"status":  "error",
			"message": err.Error(),
		})
		_ = conn.Close()
	}
}
