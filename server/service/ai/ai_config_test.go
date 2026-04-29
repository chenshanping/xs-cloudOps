package ai

import (
	"encoding/json"
	"testing"

	appconfig "server/config"
	"server/global"
	"server/model"
	"server/testutil"
)

func TestAIServiceGetAdminConfigReadsPersistedAIConfig(t *testing.T) {
	db := testutil.SetupStorageServiceTestDB(t)

	saved := `{"default_provider":"OpenAI","providers":[{"name":"OpenAI","api_key":"sk-test","base_url":"https://api.openai.com/v1","models":[{"id":"gpt-4o","name":"GPT-4o","description":"default"}]}]}`
	if err := db.Create(&model.SysConfig{
		Name:      "AI配置",
		Key:       "ai_config",
		Value:     saved,
		ValueType: "json",
		Remark:    "AI平台配置",
	}).Error; err != nil {
		t.Fatalf("create ai_config: %v", err)
	}

	cfg, err := Default.GetAdminConfig()
	if err != nil {
		t.Fatalf("GetAdminConfig error: %v", err)
	}
	if cfg.DefaultProvider != "OpenAI" {
		t.Fatalf("default provider = %s, want %s", cfg.DefaultProvider, "OpenAI")
	}
	if len(cfg.Providers) != 1 {
		t.Fatalf("providers len = %d, want 1", len(cfg.Providers))
	}
	if len(cfg.Providers[0].Models) != 1 || cfg.Providers[0].Models[0].ID != "gpt-4o" {
		t.Fatalf("models = %#v, want one gpt-4o", cfg.Providers[0].Models)
	}
}

func TestAIServiceGetAdminConfigReturnsEmptyConfigWhenMissing(t *testing.T) {
	testutil.SetupStorageServiceTestDB(t)

	previousAIConfig := global.Config.AI
	global.Config.AI = appconfig.AI{
		DefaultProvider: "fallback",
		Providers: []appconfig.AIProvider{
			{Name: "fallback"},
		},
	}
	t.Cleanup(func() {
		global.Config.AI = previousAIConfig
	})

	cfg, err := Default.GetAdminConfig()
	if err != nil {
		t.Fatalf("GetAdminConfig error: %v", err)
	}
	if cfg.DefaultProvider != "" {
		t.Fatalf("default provider = %s, want empty", cfg.DefaultProvider)
	}
	if len(cfg.Providers) != 0 {
		t.Fatalf("providers len = %d, want 0", len(cfg.Providers))
	}
}

func TestAIServiceSaveAdminConfigUpsertsCompatibleAIConfig(t *testing.T) {
	db := testutil.SetupStorageServiceTestDB(t)

	if err := db.Create(&model.SysConfig{
		Name:      "AI配置",
		Key:       "ai_config",
		Value:     `{"default_provider":"Old","providers":[]}`,
		ValueType: "json",
		Remark:    "旧备注",
	}).Error; err != nil {
		t.Fatalf("create existing ai_config: %v", err)
	}

	next := &appconfig.AI{
		DefaultProvider: "NewProvider",
		Providers: []appconfig.AIProvider{
			{
				Name:    "NewProvider",
				APIKey:  "sk-new",
				BaseURL: "https://provider.example/v1",
				Models: []appconfig.AIModel{
					{ID: "model-a", Name: "Model A", Description: "alpha"},
				},
			},
		},
	}

	if err := Default.SaveAdminConfig(next); err != nil {
		t.Fatalf("SaveAdminConfig error: %v", err)
	}

	var stored model.SysConfig
	if err := db.Where("`key` = ?", "ai_config").First(&stored).Error; err != nil {
		t.Fatalf("reload ai_config: %v", err)
	}
	if stored.ValueType != "json" {
		t.Fatalf("value_type = %s, want json", stored.ValueType)
	}

	var persisted appconfig.AI
	if err := json.Unmarshal([]byte(stored.Value), &persisted); err != nil {
		t.Fatalf("unmarshal stored value: %v", err)
	}
	if persisted.DefaultProvider != next.DefaultProvider {
		t.Fatalf("default provider = %s, want %s", persisted.DefaultProvider, next.DefaultProvider)
	}
	if len(persisted.Providers) != 1 || persisted.Providers[0].Models[0].ID != "model-a" {
		t.Fatalf("persisted providers = %#v", persisted.Providers)
	}

	var count int64
	if err := db.Model(&model.SysConfig{}).Where("`key` = ?", "ai_config").Count(&count).Error; err != nil {
		t.Fatalf("count ai_config rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("ai_config row count = %d, want 1", count)
	}
}
