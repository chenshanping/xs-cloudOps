package cmdbterminal

import "time"

type SessionConnectPayload struct {
	SessionID          uint   `json:"session_id"`
	WSToken            string `json:"ws_token"`
	WSURL              string `json:"ws_url"`
	IdleTimeoutSeconds int    `json:"idle_timeout_seconds"`
}

type SessionListItem struct {
	ID                   uint       `json:"id"`
	HostID               uint       `json:"host_id"`
	HostName             string     `json:"host_name"`
	UserID               uint       `json:"user_id"`
	UsernameSnapshot     string     `json:"username_snapshot"`
	CredentialIDSnapshot uint       `json:"credential_id_snapshot"`
	ClientIP             string     `json:"client_ip"`
	Status               string     `json:"status"`
	StartTime            *time.Time `json:"start_time"`
	EndTime              *time.Time `json:"end_time"`
	IdleTimeoutSeconds   int        `json:"idle_timeout_seconds"`
	DisconnectReason     string     `json:"disconnect_reason"`
	ForcedByUserID       uint       `json:"forced_by_user_id"`
	HostKeyFingerprint   string     `json:"host_key_fingerprint"`
	LastActivityAt       *time.Time `json:"last_activity_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type SessionLogItem struct {
	ID         uint      `json:"id"`
	SessionID  uint      `json:"session_id"`
	Seq        uint64    `json:"seq"`
	StreamType string    `json:"stream_type"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
