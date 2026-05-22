package request

type CreateCmdbTerminalSessionRequest struct {
	HostID uint `json:"host_id" binding:"required" comment:"主机ID"`
}

type CmdbTerminalSessionListRequest struct {
	PageRequest
	HostID *uint  `json:"host_id" form:"host_id" comment:"主机ID"`
	UserID *uint  `json:"user_id" form:"user_id" comment:"用户ID"`
	Status string `json:"status" form:"status" binding:"omitempty,oneof=prepared active closed failed" comment:"会话状态"`
}

type CmdbTerminalLogListRequest struct {
	PageRequest
	StreamType string `json:"stream_type" form:"stream_type" binding:"omitempty,oneof=input output system" comment:"流类型"`
}
