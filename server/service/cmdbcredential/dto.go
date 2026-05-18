package cmdbcredential

import "server/model"

type CredentialSummary struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	AuthType  string `json:"auth_type"`
	Username  string `json:"username"`
	Remark    string `json:"remark"`
	BindCount int64  `json:"bind_count"`
}

func buildSummary(item model.CmdbSshCredential, bindCount int64) CredentialSummary {
	return CredentialSummary{
		ID:        item.ID,
		Name:      item.Name,
		AuthType:  item.AuthType,
		Username:  item.Username,
		Remark:    item.Remark,
		BindCount: bindCount,
	}
}
