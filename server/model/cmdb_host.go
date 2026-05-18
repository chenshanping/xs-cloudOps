package model

import "time"

const (
	CmdbHostVerifyStatusPending = "pending"
	CmdbHostVerifyStatusSuccess = "success"
	CmdbHostVerifyStatusFailed  = "failed"
)

type CmdbHost struct {
	ID              uint       `json:"id" gorm:"primaryKey;comment:主键ID"`
	Name            string     `json:"name" gorm:"size:100;uniqueIndex;not null;comment:主机名称"`
	GroupID         uint       `json:"group_id" gorm:"index;not null;comment:主机分组ID"`
	Environment     string     `json:"environment" gorm:"size:50;default:'';comment:环境标识"`
	Owner           string     `json:"owner" gorm:"size:100;default:'';comment:负责人"`
	PrivateIP       string     `json:"private_ip" gorm:"size:45;index;default:'';comment:内网IP"`
	PublicIP        string     `json:"public_ip" gorm:"size:45;index;default:'';comment:公网IP"`
	SshHost         string     `json:"ssh_host" gorm:"size:255;index;not null;comment:SSH连接地址"`
	SshPort         int        `json:"ssh_port" gorm:"default:22;not null;comment:SSH端口"`
	CredentialID    uint       `json:"credential_id" gorm:"index;not null;comment:SSH凭据ID"`
	Remark          string     `json:"remark" gorm:"size:500;default:'';comment:备注"`
	VerifyStatus    string     `json:"verify_status" gorm:"size:20;index;default:'pending';not null;comment:校验状态:pending/success/failed"`
	VerifyMessage   string     `json:"verify_message" gorm:"size:500;default:'';comment:校验结果说明"`
	LastVerifiedAt  *time.Time `json:"last_verified_at" gorm:"comment:最后校验时间"`
	Hostname        string     `json:"hostname" gorm:"size:255;default:'';comment:主机名"`
	OS              string     `json:"os" gorm:"size:100;default:'';comment:操作系统"`
	Platform        string     `json:"platform" gorm:"size:100;default:'';comment:发行版标识"`
	PlatformVersion string     `json:"platform_version" gorm:"size:255;default:'';comment:系统版本"`
	KernelVersion   string     `json:"kernel_version" gorm:"size:100;default:'';comment:内核版本"`
	Architecture    string     `json:"architecture" gorm:"size:50;default:'';comment:系统架构"`
	CpuCores        int        `json:"cpu_cores" gorm:"default:0;not null;comment:CPU核数"`
	MemoryMB        int64      `json:"memory_mb" gorm:"default:0;not null;comment:内存大小MB"`
	CreatedAt       time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"comment:更新时间"`
}

func (CmdbHost) TableName() string {
	return "cmdb_host"
}
