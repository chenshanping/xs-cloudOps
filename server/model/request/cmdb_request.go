package request

import "mime/multipart"

type CmdbHostGroupListRequest struct {
	Name   string `json:"name" form:"name" comment:"分组名称"`
	Status *int   `json:"status" form:"status" comment:"状态"`
}

type CreateCmdbHostGroupRequest struct {
	Name   string `json:"name" binding:"required,max=100" comment:"分组名称"`
	Sort   int    `json:"sort" comment:"排序值"`
	Remark string `json:"remark" binding:"max=500" comment:"备注"`
	Status int    `json:"status" binding:"oneof=0 1" comment:"状态"`
}

type UpdateCmdbHostGroupRequest struct {
	Name   string `json:"name" binding:"required,max=100" comment:"分组名称"`
	Sort   int    `json:"sort" comment:"排序值"`
	Remark string `json:"remark" binding:"max=500" comment:"备注"`
	Status int    `json:"status" binding:"oneof=0 1" comment:"状态"`
}

type CmdbHostTagListRequest struct {
	Name string `json:"name" form:"name" comment:"标签名称"`
}

type CreateCmdbHostTagRequest struct {
	Name   string `json:"name" binding:"required,max=100" comment:"标签名称"`
	Color  string `json:"color" binding:"max=30" comment:"标签颜色"`
	Remark string `json:"remark" binding:"max=500" comment:"备注"`
}

type UpdateCmdbHostTagRequest struct {
	Name   string `json:"name" binding:"required,max=100" comment:"标签名称"`
	Color  string `json:"color" binding:"max=30" comment:"标签颜色"`
	Remark string `json:"remark" binding:"max=500" comment:"备注"`
}

type CmdbCredentialListRequest struct {
	PageRequest
	Name     string `json:"name" form:"name" comment:"凭据名称"`
	AuthType string `json:"auth_type" form:"auth_type" comment:"认证方式"`
}

type CreateCmdbCredentialRequest struct {
	Name       string `json:"name" binding:"required,max=100" comment:"凭据名称"`
	AuthType   string `json:"auth_type" binding:"required,oneof=password private_key" comment:"认证方式"`
	Username   string `json:"username" binding:"required,max=100" comment:"登录用户名"`
	Password   string `json:"password" binding:"omitempty,max=5000" comment:"登录密码"`
	PrivateKey string `json:"private_key" binding:"omitempty" comment:"私钥内容"`
	Passphrase string `json:"passphrase" binding:"omitempty,max=5000" comment:"私钥口令"`
	Remark     string `json:"remark" binding:"max=500" comment:"备注"`
}

type UpdateCmdbCredentialRequest struct {
	Name       string `json:"name" binding:"required,max=100" comment:"凭据名称"`
	AuthType   string `json:"auth_type" binding:"required,oneof=password private_key" comment:"认证方式"`
	Username   string `json:"username" binding:"required,max=100" comment:"登录用户名"`
	Password   string `json:"password" binding:"omitempty,max=5000" comment:"登录密码"`
	PrivateKey string `json:"private_key" binding:"omitempty" comment:"私钥内容"`
	Passphrase string `json:"passphrase" binding:"omitempty,max=5000" comment:"私钥口令"`
	Remark     string `json:"remark" binding:"max=500" comment:"备注"`
}

type CmdbHostListRequest struct {
	PageRequest
	Keyword      string `json:"keyword" form:"keyword" comment:"关键字"`
	GroupID      *uint  `json:"group_id" form:"group_id" comment:"分组ID"`
	TagID        *uint  `json:"tag_id" form:"tag_id" comment:"标签ID"`
	VerifyStatus string `json:"verify_status" form:"verify_status" comment:"校验状态"`
	Environment  string `json:"environment" form:"environment" comment:"环境标识"`
}

type CreateCmdbHostRequest struct {
	Name         string `json:"name" binding:"required,max=100" comment:"主机名称"`
	GroupID      uint   `json:"group_id" binding:"required" comment:"主机分组ID"`
	TagIDs       []uint `json:"tag_ids" comment:"标签ID列表"`
	Environment  string `json:"environment" binding:"max=50" comment:"环境标识"`
	Owner        string `json:"owner" binding:"max=100" comment:"负责人"`
	PrivateIP    string `json:"private_ip" binding:"max=45" comment:"内网IP"`
	PublicIP     string `json:"public_ip" binding:"max=45" comment:"公网IP"`
	SshHost      string `json:"ssh_host" binding:"required,max=255" comment:"SSH连接地址"`
	SshPort      int    `json:"ssh_port" binding:"omitempty,min=1,max=65535" comment:"SSH端口"`
	CredentialID uint   `json:"credential_id" binding:"required" comment:"SSH凭据ID"`
	Remark       string `json:"remark" binding:"max=500" comment:"备注"`
}

type UpdateCmdbHostRequest struct {
	Name         string `json:"name" binding:"required,max=100" comment:"主机名称"`
	GroupID      uint   `json:"group_id" binding:"required" comment:"主机分组ID"`
	TagIDs       []uint `json:"tag_ids" comment:"标签ID列表"`
	Environment  string `json:"environment" binding:"max=50" comment:"环境标识"`
	Owner        string `json:"owner" binding:"max=100" comment:"负责人"`
	PrivateIP    string `json:"private_ip" binding:"max=45" comment:"内网IP"`
	PublicIP     string `json:"public_ip" binding:"max=45" comment:"公网IP"`
	SshHost      string `json:"ssh_host" binding:"required,max=255" comment:"SSH连接地址"`
	SshPort      int    `json:"ssh_port" binding:"omitempty,min=1,max=65535" comment:"SSH端口"`
	CredentialID uint   `json:"credential_id" binding:"required" comment:"SSH凭据ID"`
	Remark       string `json:"remark" binding:"max=500" comment:"备注"`
}

type CmdbHostImportItem struct {
	Name           string   `json:"name" comment:"主机名称"`
	GroupName      string   `json:"group_name" comment:"分组名称"`
	TagNames       []string `json:"tag_names" comment:"标签名称列表"`
	Environment    string   `json:"environment" comment:"环境标识"`
	Owner          string   `json:"owner" comment:"负责人"`
	PrivateIP      string   `json:"private_ip" comment:"内网IP"`
	PublicIP       string   `json:"public_ip" comment:"公网IP"`
	SshHost        string   `json:"ssh_host" comment:"SSH连接地址"`
	SshPort        int      `json:"ssh_port" comment:"SSH端口"`
	CredentialName string   `json:"credential_name" comment:"SSH凭据名称"`
	Remark         string   `json:"remark" comment:"备注"`
}

type CmdbHostImportRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required" comment:"导入文件"`
}
