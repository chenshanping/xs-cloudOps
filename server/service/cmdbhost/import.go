package cmdbhost

import (
	"errors"
	"fmt"
	"strings"

	"server/global"
	"server/model"
	"server/model/request"
	"server/utils"
)

var cmdbHostImportHeaders = []string{
	"主机名称",
	"分组名称",
	"标签名称列表",
	"环境标识",
	"SSH连接地址",
	"SSH端口",
	"SSH凭据名称",
}

func (s *Service) GetImportTemplate() ([]byte, string, error) {
	exporter := utils.NewExcelExporter("主机导入模板")
	if err := exporter.SetHeaders(cmdbHostImportHeaders); err != nil {
		return nil, "", err
	}
	if err := exporter.AddRow([]interface{}{
		"prod-web-01", "生产环境", "web,prod", "prod", "10.10.10.11", 22, "生产Root凭据",
	}); err != nil {
		return nil, "", err
	}
	if err := exporter.AddRow([]interface{}{
		"test-db-01", "测试环境", "db,test", "test", "172.16.1.20", 22, "测试DBA凭据",
	}); err != nil {
		return nil, "", err
	}
	buf, err := exporter.SaveToBuffer()
	if err != nil {
		return nil, "", err
	}
	return buf, "主机导入模板.xlsx", nil
}

func (s *Service) ImportHosts(fileData []byte) (*HostImportResult, error) {
	importer, err := utils.NewExcelImporter(fileData)
	if err != nil {
		return nil, errors.New("打开导入文件失败")
	}
	headers, err := importer.GetHeaders()
	if err != nil {
		return nil, errors.New("读取导入文件失败")
	}
	if err := utils.ValidateHeaders(headers, buildCmdbHostImportFields()); err != nil {
		return nil, err
	}
	rows, err := importer.GetRows()
	if err != nil {
		return nil, errors.New("读取导入文件失败")
	}
	result := &HostImportResult{Rows: make([]HostImportRowResult, 0)}
	if len(rows) == 0 {
		return result, nil
	}
	for index, row := range rows {
		rowNo := index + 2
		result.Total++
		item := parseImportRow(row)
		rowResult := HostImportRowResult{Row: rowNo, Name: item.Name}
		if item.Name == "" || item.GroupName == "" || item.SshHost == "" || item.CredentialName == "" {
			rowResult.ErrorMessage = "必填列缺失"
			result.FailureCount++
			result.Rows = append(result.Rows, rowResult)
			continue
		}
		groupID, err := lookupGroupIDByName(item.GroupName)
		if err != nil {
			rowResult.ErrorMessage = err.Error()
			result.FailureCount++
			result.Rows = append(result.Rows, rowResult)
			continue
		}
		credentialID, err := lookupCredentialIDByName(item.CredentialName)
		if err != nil {
			rowResult.ErrorMessage = err.Error()
			result.FailureCount++
			result.Rows = append(result.Rows, rowResult)
			continue
		}
		tagIDs, err := lookupTagIDsByNames(item.TagNames)
		if err != nil {
			rowResult.ErrorMessage = err.Error()
			result.FailureCount++
			result.Rows = append(result.Rows, rowResult)
			continue
		}
		host, createErr := s.CreateHost(&request.CreateCmdbHostRequest{
			Name:         item.Name,
			GroupID:      groupID,
			TagIDs:       tagIDs,
			Environment:  item.Environment,
			Owner:        item.Owner,
			PrivateIP:    item.PrivateIP,
			PublicIP:     item.PublicIP,
			SshHost:      item.SshHost,
			SshPort:      item.SshPort,
			CredentialID: credentialID,
			Remark:       item.Remark,
		})
		if createErr != nil {
			rowResult.ErrorMessage = createErr.Error()
			result.FailureCount++
			result.Rows = append(result.Rows, rowResult)
			continue
		}
		rowResult.Created = true
		rowResult.VerifySuccess = host.VerifyStatus == model.CmdbHostVerifyStatusSuccess
		rowResult.VerifyMessage = host.VerifyMessage
		result.SuccessCount++
		result.Rows = append(result.Rows, rowResult)
	}
	return result, nil
}

func buildCmdbHostImportFields() []utils.ImportField {
	return []utils.ImportField{
		{Header: "主机名称", Key: "name", Required: true, Type: "string", MaxLen: 100},
		{Header: "分组名称", Key: "group_name", Required: true, Type: "string", MaxLen: 100},
		{Header: "标签名称列表", Key: "tag_names", Type: "string", MaxLen: 500},
		{Header: "环境标识", Key: "environment", Type: "string", MaxLen: 50},
		{Header: "SSH连接地址", Key: "ssh_host", Required: true, Type: "string", MaxLen: 255},
		{Header: "SSH端口", Key: "ssh_port", Type: "int"},
		{Header: "SSH凭据名称", Key: "credential_name", Required: true, Type: "string", MaxLen: 100},
	}
}

func parseImportRow(row []string) request.CmdbHostImportItem {
	item := request.CmdbHostImportItem{}
	if len(row) > 0 {
		item.Name = strings.TrimSpace(row[0])
	}
	if len(row) > 1 {
		item.GroupName = strings.TrimSpace(row[1])
	}
	if len(row) > 2 {
		rawTags := strings.TrimSpace(row[2])
		if rawTags != "" {
			for _, tag := range strings.Split(rawTags, ",") {
				name := strings.TrimSpace(tag)
				if name != "" {
					item.TagNames = append(item.TagNames, name)
				}
			}
		}
	}
	if len(row) > 3 {
		item.Environment = strings.TrimSpace(row[3])
	}
	if len(row) > 4 {
		item.SshHost = strings.TrimSpace(row[4])
	}
	if len(row) > 5 {
		fmt.Sscanf(strings.TrimSpace(row[5]), "%d", &item.SshPort)
	}
	if len(row) > 6 {
		item.CredentialName = strings.TrimSpace(row[6])
	}
	return item
}

func lookupGroupIDByName(name string) (uint, error) {
	var item model.CmdbHostGroup
	if err := global.DB.Where("name = ?", name).First(&item).Error; err != nil {
		return 0, errors.New("分组不存在: " + name)
	}
	return item.ID, nil
}

func lookupCredentialIDByName(name string) (uint, error) {
	var item model.CmdbSshCredential
	if err := global.DB.Where("name = ?", name).First(&item).Error; err != nil {
		return 0, errors.New("SSH凭据不存在: " + name)
	}
	return item.ID, nil
}

func lookupTagIDsByNames(names []string) ([]uint, error) {
	if len(names) == 0 {
		return nil, nil
	}
	var tags []model.CmdbHostTag
	if err := global.DB.Where("name IN ?", names).Find(&tags).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(tags))
	found := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		ids = append(ids, tag.ID)
		found[tag.Name] = struct{}{}
	}
	for _, name := range names {
		if _, ok := found[name]; !ok {
			return nil, errors.New("标签不存在: " + name)
		}
	}
	return ids, nil
}
