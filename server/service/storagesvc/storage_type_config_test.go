package storagesvc

import (
	"testing"

	"server/model"
	"server/service/configsvc"
	"server/testutil"
)

func TestStorageServiceGetStorageByTypeUsesTypeSpecificConfig(t *testing.T) {
	db := testutil.SetupStorageServiceTestDB(t)

	configs := []model.SysConfig{
		{Key: configsvc.StorageTypeConfigKey, Value: string(model.StorageTypeMinio)},
		{Key: StorageConfigKey(model.StorageTypeLocal), Value: `{"base_path":"local-dir","base_url":"/api/v1/upload"}`},
		{Key: StorageConfigKey(model.StorageTypeMinio), Value: `{"endpoint":"127.0.0.1:9000","access_key_id":"minio","secret_access_key":"secret","bucket_name":"files","use_ssl":false}`},
	}
	if err := db.Create(&configs).Error; err != nil {
		t.Fatalf("create configs: %v", err)
	}

	localStorage, err := Default.GetStorageByType(model.StorageTypeLocal)
	if err != nil {
		t.Fatalf("GetStorageByType(local) error: %v", err)
	}
	if localStorage.Type != model.StorageTypeLocal {
		t.Fatalf("local storage type = %s, want %s", localStorage.Type, model.StorageTypeLocal)
	}
	if localStorage.Config != `{"base_path":"local-dir","base_url":"/api/v1/upload"}` {
		t.Fatalf("local storage config = %s", localStorage.Config)
	}

	minioStorage, err := Default.GetStorageByType(model.StorageTypeMinio)
	if err != nil {
		t.Fatalf("GetStorageByType(minio) error: %v", err)
	}
	if minioStorage.Type != model.StorageTypeMinio {
		t.Fatalf("minio storage type = %s, want %s", minioStorage.Type, model.StorageTypeMinio)
	}
	if minioStorage.Config != `{"endpoint":"127.0.0.1:9000","access_key_id":"minio","secret_access_key":"secret","bucket_name":"files","use_ssl":false}` {
		t.Fatalf("minio storage config = %s", minioStorage.Config)
	}
}
