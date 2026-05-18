package cmdbhost

import "time"

type HostGroupOption struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

type HostTagOption struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	Remark string `json:"remark"`
}

type HostTagView struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type HostCredentialSummary struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	AuthType string `json:"auth_type"`
	Username string `json:"username"`
}

type HostListItem struct {
	ID                uint                  `json:"id"`
	Name              string                `json:"name"`
	Group             HostGroupOption       `json:"group"`
	Tags              []HostTagView         `json:"tags"`
	Environment       string                `json:"environment"`
	Owner             string                `json:"owner"`
	Remark            string                `json:"remark"`
	PrivateIP         string                `json:"private_ip"`
	PublicIP          string                `json:"public_ip"`
	SshHost           string                `json:"ssh_host"`
	SshPort           int                   `json:"ssh_port"`
	CredentialSummary HostCredentialSummary `json:"credential_summary"`
	VerifyStatus      string                `json:"verify_status"`
	VerifyMessage     string                `json:"verify_message"`
	LastVerifiedAt    *time.Time            `json:"last_verified_at"`
	Hostname          string                `json:"hostname"`
	OS                string                `json:"os"`
	Platform          string                `json:"platform"`
	PlatformVersion   string                `json:"platform_version"`
	KernelVersion     string                `json:"kernel_version"`
	Architecture      string                `json:"architecture"`
	CpuCores          int                   `json:"cpu_cores"`
	MemoryMB          int64                 `json:"memory_mb"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

type HostImportRowResult struct {
	Row           int    `json:"row"`
	Name          string `json:"name"`
	Created       bool   `json:"created"`
	VerifySuccess bool   `json:"verify_success"`
	ErrorMessage  string `json:"error_message"`
	VerifyMessage string `json:"verify_message"`
}

type HostImportResult struct {
	Total        int                   `json:"total"`
	SuccessCount int                   `json:"success_count"`
	FailureCount int                   `json:"failure_count"`
	Rows         []HostImportRowResult `json:"rows"`
}
