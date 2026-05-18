package cmdbcredential

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/model/request"
)

func setupCredentialServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.CmdbSshCredential{}, &model.CmdbHost{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	previousDB := global.DB
	global.DB = db
	t.Cleanup(func() {
		global.DB = previousDB
	})
	return db
}

func TestCredentialService_Delete_AllowsSameNameRecreate(t *testing.T) {
	setupCredentialServiceTestDB(t)
	svc := Default

	createReq := &request.CreateCmdbCredentialRequest{
		Name:     "root-key",
		AuthType: model.CmdbCredentialAuthTypePassword,
		Username: "root",
		Password: "secret",
	}
	if err := svc.Create(createReq); err != nil {
		t.Fatalf("create credential: %v", err)
	}

	var credential model.CmdbSshCredential
	if err := global.DB.Where("name = ?", "root-key").First(&credential).Error; err != nil {
		t.Fatalf("load credential: %v", err)
	}
	if err := svc.Delete(credential.ID); err != nil {
		t.Fatalf("delete credential: %v", err)
	}
	if err := svc.Create(createReq); err != nil {
		t.Fatalf("recreate credential with same name: %v", err)
	}
}
