package cmdbhost

import (
	"bytes"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
)

func setupCmdbHostServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(
		&model.CmdbHostGroup{},
		&model.CmdbHostTag{},
		&model.CmdbHostTagRel{},
		&model.CmdbSshCredential{},
		&model.CmdbHost{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})
	return db
}

func TestHostService_DeleteGroup_AllowsSameNameRecreate(t *testing.T) {
	setupCmdbHostServiceTestDB(t)
	svc := Default
	req := &request.CreateCmdbHostGroupRequest{Name: "生产环境", Status: 1}
	if err := svc.CreateGroup(req); err != nil {
		t.Fatalf("create group: %v", err)
	}
	var group model.CmdbHostGroup
	if err := global.DB.Where("name = ?", "生产环境").First(&group).Error; err != nil {
		t.Fatalf("load group: %v", err)
	}
	if err := svc.DeleteGroup(group.ID); err != nil {
		t.Fatalf("delete group: %v", err)
	}
	if err := svc.CreateGroup(req); err != nil {
		t.Fatalf("recreate group with same name: %v", err)
	}
}

func TestHostService_DeleteTag_AllowsSameNameRecreate(t *testing.T) {
	setupCmdbHostServiceTestDB(t)
	svc := Default
	req := &request.CreateCmdbHostTagRequest{Name: "prod"}
	if err := svc.CreateTag(req); err != nil {
		t.Fatalf("create tag: %v", err)
	}
	var tag model.CmdbHostTag
	if err := global.DB.Where("name = ?", "prod").First(&tag).Error; err != nil {
		t.Fatalf("load tag: %v", err)
	}
	if err := svc.DeleteTag(tag.ID); err != nil {
		t.Fatalf("delete tag: %v", err)
	}
	if err := svc.CreateTag(req); err != nil {
		t.Fatalf("recreate tag with same name: %v", err)
	}
}

func TestHostService_DeleteHost_AllowsSameNameRecreate(t *testing.T) {
	setupCmdbHostServiceTestDB(t)
	svc := Default

	group := model.CmdbHostGroup{Name: "默认分组", Status: 1}
	if err := global.DB.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	credential := model.CmdbSshCredential{Name: "默认凭据", AuthType: model.CmdbCredentialAuthTypePassword, Username: "root", Password: "secret"}
	if err := global.DB.Create(&credential).Error; err != nil {
		t.Fatalf("seed credential: %v", err)
	}

	hostReq := &request.CreateCmdbHostRequest{
		Name:         "host-01",
		GroupID:      group.ID,
		SshHost:      "127.0.0.1",
		SshPort:      22,
		CredentialID: credential.ID,
	}
	if _, err := svc.CreateHost(hostReq); err != nil {
		t.Fatalf("create host: %v", err)
	}
	var host model.CmdbHost
	if err := global.DB.Where("name = ?", "host-01").First(&host).Error; err != nil {
		t.Fatalf("load host: %v", err)
	}
	if err := svc.DeleteHost(host.ID); err != nil {
		t.Fatalf("delete host: %v", err)
	}
	if _, err := svc.CreateHost(hostReq); err != nil {
		t.Fatalf("recreate host with same name: %v", err)
	}
}

func TestHostService_CreateHost_VerifyFailureStillPersists(t *testing.T) {
	setupCmdbHostServiceTestDB(t)
	svc := Default

	group := model.CmdbHostGroup{Name: "默认分组", Status: 1}
	if err := global.DB.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	credential := model.CmdbSshCredential{Name: "默认凭据", AuthType: model.CmdbCredentialAuthTypePassword, Username: "root", Password: "secret"}
	if err := global.DB.Create(&credential).Error; err != nil {
		t.Fatalf("seed credential: %v", err)
	}

	item, err := svc.CreateHost(&request.CreateCmdbHostRequest{
		Name:         "host-verify-fail",
		GroupID:      group.ID,
		SshHost:      "127.0.0.1",
		SshPort:      65022,
		CredentialID: credential.ID,
	})
	if err != nil {
		t.Fatalf("create host: %v", err)
	}
	if item.VerifyStatus != model.CmdbHostVerifyStatusFailed {
		t.Fatalf("verify status = %s, want %s", item.VerifyStatus, model.CmdbHostVerifyStatusFailed)
	}
	var host model.CmdbHost
	if err := global.DB.Where("name = ?", "host-verify-fail").First(&host).Error; err != nil {
		t.Fatalf("host not persisted: %v", err)
	}
	if host.VerifyStatus != model.CmdbHostVerifyStatusFailed {
		t.Fatalf("persisted verify status = %s, want %s", host.VerifyStatus, model.CmdbHostVerifyStatusFailed)
	}
}

func TestHostService_ImportHosts_ReturnsMixedResults(t *testing.T) {
	setupCmdbHostServiceTestDB(t)
	svc := Default

	group := model.CmdbHostGroup{Name: "默认分组", Status: 1}
	if err := global.DB.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	credential := model.CmdbSshCredential{Name: "默认凭据", AuthType: model.CmdbCredentialAuthTypePassword, Username: "root", Password: "secret"}
	if err := global.DB.Create(&credential).Error; err != nil {
		t.Fatalf("seed credential: %v", err)
	}
	tag := model.CmdbHostTag{Name: "prod"}
	if err := global.DB.Create(&tag).Error; err != nil {
		t.Fatalf("seed tag: %v", err)
	}

	file := excelize.NewFile()
	sheet := file.GetSheetName(0)
	rows := [][]string{
		{"主机名称", "分组名称", "标签名称列表", "环境标识", "SSH连接地址", "SSH端口", "SSH凭据名称"},
		{"host-import-1", "默认分组", "prod", "prod", "127.0.0.1", "65022", "默认凭据"},
		{"host-import-2", "不存在分组", "prod", "prod", "127.0.0.1", "22", "默认凭据"},
	}
	for rowIndex, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, rowIndex+1)
		if err := file.SetSheetRow(sheet, cell, &row); err != nil {
			t.Fatalf("set row %d: %v", rowIndex, err)
		}
	}
	buffer, err := file.WriteToBuffer()
	if err != nil {
		t.Fatalf("write buffer: %v", err)
	}

	result, err := svc.ImportHosts(bytes.Clone(buffer.Bytes()))
	if err != nil {
		t.Fatalf("import hosts: %v", err)
	}
	if result.Total != 2 {
		t.Fatalf("total = %d, want 2", result.Total)
	}
	if result.SuccessCount != 1 {
		t.Fatalf("success_count = %d, want 1", result.SuccessCount)
	}
	if result.FailureCount != 1 {
		t.Fatalf("failure_count = %d, want 1", result.FailureCount)
	}
}

func TestHostService_GetImportTemplate_UsesCompactHeaders(t *testing.T) {
	setupCmdbHostServiceTestDB(t)
	svc := Default

	buf, filename, err := svc.GetImportTemplate()
	if err != nil {
		t.Fatalf("get import template: %v", err)
	}
	if filename != "主机导入模板.xlsx" {
		t.Fatalf("filename = %s", filename)
	}

	file, err := excelize.OpenReader(bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("open template: %v", err)
	}
	defer file.Close()

	rows, err := file.GetRows(file.GetSheetName(0))
	if err != nil {
		t.Fatalf("get rows: %v", err)
	}
	if len(rows) < 3 {
		t.Fatalf("template rows = %d, want at least 3", len(rows))
	}
	gotHeaders := rows[0]
	if len(gotHeaders) != len(cmdbHostImportHeaders) {
		t.Fatalf("header count = %d, want %d", len(gotHeaders), len(cmdbHostImportHeaders))
	}
	for i, header := range cmdbHostImportHeaders {
		if gotHeaders[i] != header {
			t.Fatalf("header[%d] = %s, want %s", i, gotHeaders[i], header)
		}
	}
}

func TestResolvePlatformInfo_UsesDistributionNameInsteadOfKernelFamily(t *testing.T) {
	osName, platform, version := resolvePlatformInfo(
		"Linux",
		"CentOS Linux",
		"centos",
		"CentOS Linux 7 (Core)",
		"",
	)
	if osName != "CentOS Linux" {
		t.Fatalf("osName = %s, want CentOS Linux", osName)
	}
	if platform != "centos" {
		t.Fatalf("platform = %s, want centos", platform)
	}
	if version != "CentOS Linux 7 (Core)" {
		t.Fatalf("version = %s, want CentOS Linux 7 (Core)", version)
	}
}

func TestResolvePlatformInfo_FallsBackToLegacyRelease(t *testing.T) {
	osName, platform, version := resolvePlatformInfo(
		"Linux",
		"",
		"",
		"",
		"CentOS release 6.10 (Final)",
	)
	if osName != "CentOS Linux" {
		t.Fatalf("osName = %s, want CentOS Linux", osName)
	}
	if platform != "centos" {
		t.Fatalf("platform = %s, want centos", platform)
	}
	if version != "6.10" {
		t.Fatalf("version = %s, want 6.10", version)
	}
}
