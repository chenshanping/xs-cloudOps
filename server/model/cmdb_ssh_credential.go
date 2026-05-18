package model

import "time"

const (
	CmdbCredentialAuthTypePassword   = "password"
	CmdbCredentialAuthTypePrivateKey = "private_key"
)

type CmdbSshCredential struct {
	ID         uint      `json:"id" gorm:"primaryKey;comment:主键ID"`
	Name       string    `json:"name" gorm:"size:100;uniqueIndex;not null;comment:凭据名称"`
	AuthType   string    `json:"auth_type" gorm:"size:20;index;not null;comment:认证方式:password/private_key"`
	Username   string    `json:"username" gorm:"size:100;not null;comment:登录用户名"`
	Password   string    `json:"-" gorm:"type:text;comment:登录密码"`
	PrivateKey string    `json:"-" gorm:"type:longtext;comment:私钥内容"`
	Passphrase string    `json:"-" gorm:"type:text;comment:私钥口令"`
	Remark     string    `json:"remark" gorm:"size:500;default:'';comment:备注"`
	CreatedAt  time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (CmdbSshCredential) TableName() string {
	return "cmdb_ssh_credential"
}
