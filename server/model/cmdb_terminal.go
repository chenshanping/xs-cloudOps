package model

import "time"

const (
	CmdbTerminalSessionStatusPrepared = "prepared"
	CmdbTerminalSessionStatusActive   = "active"
	CmdbTerminalSessionStatusClosed   = "closed"
	CmdbTerminalSessionStatusFailed   = "failed"
)

const (
	CmdbTerminalStreamTypeInput  = "input"
	CmdbTerminalStreamTypeOutput = "output"
	CmdbTerminalStreamTypeSystem = "system"
)

type CmdbTerminalSession struct {
	ID                   uint       `json:"id" gorm:"primaryKey;comment:主键ID"`
	HostID               uint       `json:"host_id" gorm:"index;not null;comment:主机ID"`
	UserID               uint       `json:"user_id" gorm:"index;not null;comment:发起用户ID"`
	UsernameSnapshot     string     `json:"username_snapshot" gorm:"size:100;default:'';not null;comment:发起用户名快照"`
	CredentialIDSnapshot uint       `json:"credential_id_snapshot" gorm:"index;not null;comment:凭据ID快照"`
	ClientIP             string     `json:"client_ip" gorm:"size:45;default:'';not null;comment:来源IP"`
	Status               string     `json:"status" gorm:"size:20;index;default:'prepared';not null;comment:会话状态:prepared/active/closed/failed"`
	StartTime            *time.Time `json:"start_time" gorm:"comment:开始时间"`
	EndTime              *time.Time `json:"end_time" gorm:"comment:结束时间"`
	IdleTimeoutSeconds   int        `json:"idle_timeout_seconds" gorm:"default:1800;not null;comment:空闲超时秒数"`
	DisconnectReason     string     `json:"disconnect_reason" gorm:"size:50;default:'';not null;comment:断开原因"`
	ForcedByUserID       uint       `json:"forced_by_user_id" gorm:"default:0;not null;comment:强制断开操作人ID"`
	HostKeyFingerprint   string     `json:"host_key_fingerprint" gorm:"size:255;default:'';not null;comment:主机指纹"`
	LastActivityAt       *time.Time `json:"last_activity_at" gorm:"index;comment:最后活动时间"`
	CreatedAt            time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"comment:更新时间"`
}

func (CmdbTerminalSession) TableName() string {
	return "cmdb_terminal_session"
}

type CmdbTerminalLog struct {
	ID         uint      `json:"id" gorm:"primaryKey;comment:主键ID"`
	SessionID  uint      `json:"session_id" gorm:"index:idx_cmdb_terminal_log_session_seq,priority:1;not null;comment:终端会话ID"`
	Seq        uint64    `json:"seq" gorm:"index:idx_cmdb_terminal_log_session_seq,priority:2;not null;comment:日志序号"`
	StreamType string    `json:"stream_type" gorm:"size:20;index;not null;comment:流类型:input/output/system"`
	Content    string    `json:"content" gorm:"type:longtext;comment:日志内容"`
	CreatedAt  time.Time `json:"created_at" gorm:"comment:创建时间"`
}

func (CmdbTerminalLog) TableName() string {
	return "cmdb_terminal_log"
}

type CmdbHostSshFingerprint struct {
	ID             uint       `json:"id" gorm:"primaryKey;comment:主键ID"`
	HostID         uint       `json:"host_id" gorm:"uniqueIndex;not null;comment:主机ID"`
	Algorithm      string     `json:"algorithm" gorm:"size:50;default:'';not null;comment:指纹算法"`
	Fingerprint    string     `json:"fingerprint" gorm:"size:255;default:'';not null;comment:指纹值"`
	FirstSeenAt    *time.Time `json:"first_seen_at" gorm:"comment:首次记录时间"`
	LastVerifiedAt *time.Time `json:"last_verified_at" gorm:"comment:最后校验时间"`
	CreatedAt      time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"comment:更新时间"`
}

func (CmdbHostSshFingerprint) TableName() string {
	return "cmdb_host_ssh_fingerprint"
}
