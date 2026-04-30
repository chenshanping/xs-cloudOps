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

	models, err := json.Marshal([]appconfig.AIModel{
		{ID: "gpt-4o", Name: "GPT-4o", Description: "default"},
	})
	if err != nil {
		t.Fatalf("marshal models: %v", err)
	}
	if err := db.Create(&model.AIProviderConfig{
		Name:       "OpenAI",
		APIKey:     "sk-test",
		BaseURL:    "https://api.openai.com/v1",
		ModelsJSON: string(models),
		IsDefault:  true,
		Sort:       0,
	}).Error; err != nil {
		t.Fatalf("create ai provider config: %v", err)
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

func TestAIServiceGetModelsReadsProvidersInsteadOfLegacySysConfig(t *testing.T) {
	db := testutil.SetupStorageServiceTestDB(t)

	newModels, err := json.Marshal([]appconfig.AIModel{
		{ID: "provider-model", Name: "Provider Model", Description: "from providers"},
	})
	if err != nil {
		t.Fatalf("marshal provider models: %v", err)
	}
	if err := db.Create(&model.AIProviderConfig{
		Name:       "OpenAI",
		APIKey:     "sk-provider",
		BaseURL:    "https://api.openai.com/v1",
		ModelsJSON: string(newModels),
		IsDefault:  true,
		Sort:       0,
	}).Error; err != nil {
		t.Fatalf("create ai provider config: %v", err)
	}

	if err := db.Create(&model.SysConfig{
		Name:      "AI配置",
		Key:       "ai_config",
		Value:     `{"default_provider":"Legacy","providers":[{"name":"Legacy","api_key":"sk-legacy","base_url":"https://legacy.example/v1","models":[{"id":"legacy-model","name":"Legacy Model","description":"legacy"}]}]}`,
		ValueType: "json",
		Remark:    "legacy ai config",
	}).Error; err != nil {
		t.Fatalf("create legacy ai_config: %v", err)
	}

	models, err := Default.GetModels()
	if err != nil {
		t.Fatalf("GetModels error: %v", err)
	}
	if len(models) != 1 || models[0].ID != "provider-model" {
		t.Fatalf("models = %#v, want one provider-model", models)
	}
}

func TestAIServiceGetAdminConfigFallsBackToLegacySysConfigAndMigrates(t *testing.T) {
	db := testutil.SetupStorageServiceTestDB(t)

	saved := `{"default_provider":"OpenAI","providers":[{"name":"OpenAI","api_key":"sk-test","base_url":"https://api.openai.com/v1","models":[{"id":"gpt-4o","name":"GPT-4o","description":"default"}]},{"name":"DashScope","api_key":"sk-dash","base_url":"https://dashscope.aliyuncs.com/compatible-mode/v1","models":[{"id":"qwen-max","name":"qwen-max","description":"dash"}]}]}`
	if err := db.Create(&model.SysConfig{
		Name:      "AI配置",
		Key:       "ai_config",
		Value:     saved,
		ValueType: "json",
		Remark:    "AI平台配置",
	}).Error; err != nil {
		t.Fatalf("create legacy ai_config: %v", err)
	}

	cfg, err := Default.GetAdminConfig()
	if err != nil {
		t.Fatalf("GetAdminConfig error: %v", err)
	}
	if cfg.DefaultProvider != "OpenAI" {
		t.Fatalf("default provider = %s, want %s", cfg.DefaultProvider, "OpenAI")
	}

	var migrated []model.AIProviderConfig
	if err := db.Order("sort ASC, id ASC").Find(&migrated).Error; err != nil {
		t.Fatalf("reload migrated ai providers: %v", err)
	}
	if len(migrated) != 2 {
		t.Fatalf("migrated provider count = %d, want 2", len(migrated))
	}
	if migrated[0].Name != "OpenAI" || !migrated[0].IsDefault {
		t.Fatalf("first migrated provider = %#v, want OpenAI default", migrated[0])
	}
	if migrated[1].Name != "DashScope" || migrated[1].IsDefault {
		t.Fatalf("second migrated provider = %#v, want DashScope non-default", migrated[1])
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

	oldModels, err := json.Marshal([]appconfig.AIModel{{ID: "old-model", Name: "Old Model", Description: "old"}})
	if err != nil {
		t.Fatalf("marshal old models: %v", err)
	}
	if err := db.Create(&model.AIProviderConfig{
		Name:       "Old",
		APIKey:     "sk-old",
		BaseURL:    "https://old.example/v1",
		ModelsJSON: string(oldModels),
		IsDefault:  true,
		Sort:       0,
	}).Error; err != nil {
		t.Fatalf("create existing provider: %v", err)
	}

	if err := db.Create(&model.SysConfig{
		Name:      "AI配置",
		Key:       "ai_config",
		Value:     `{"default_provider":"Old","providers":[]}`,
		ValueType: "json",
		Remark:    "旧备注",
	}).Error; err != nil {
		t.Fatalf("create legacy ai_config: %v", err)
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

	var stored []model.AIProviderConfig
	if err := db.Order("sort ASC, id ASC").Find(&stored).Error; err != nil {
		t.Fatalf("reload ai providers: %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("provider count = %d, want 1", len(stored))
	}
	if !stored[0].IsDefault {
		t.Fatalf("stored provider should be default: %#v", stored[0])
	}

	var persistedModels []appconfig.AIModel
	if err := json.Unmarshal([]byte(stored[0].ModelsJSON), &persistedModels); err != nil {
		t.Fatalf("unmarshal stored models_json: %v", err)
	}
	if stored[0].Name != next.DefaultProvider {
		t.Fatalf("stored provider name = %s, want %s", stored[0].Name, next.DefaultProvider)
	}
	if len(persistedModels) != 1 || persistedModels[0].ID != "model-a" {
		t.Fatalf("persisted models = %#v", persistedModels)
	}

	var count int64
	if err := db.Model(&model.AIProviderConfig{}).Count(&count).Error; err != nil {
		t.Fatalf("count ai providers: %v", err)
	}
	if count != 1 {
		t.Fatalf("ai provider row count = %d, want 1", count)
	}

	var legacy model.SysConfig
	if err := db.Where("`key` = ?", "ai_config").First(&legacy).Error; err != nil {
		t.Fatalf("reload legacy ai_config: %v", err)
	}
	if legacy.Value != `{"default_provider":"Old","providers":[]}` {
		t.Fatalf("legacy ai_config value changed = %s", legacy.Value)
	}
}

func TestAIServiceCreateConversationUsesProvidersDefaultModel(t *testing.T) {
	db := testutil.SetupStorageServiceTestDB(t)
	if err := db.AutoMigrate(&model.AIConversation{}); err != nil {
		t.Fatalf("auto migrate ai conversations: %v", err)
	}

	models, err := json.Marshal([]appconfig.AIModel{
		{ID: "provider-default", Name: "Provider Default", Description: ""},
	})
	if err != nil {
		t.Fatalf("marshal provider models: %v", err)
	}
	if err := db.Create(&model.AIProviderConfig{
		Name:       "OpenAI",
		APIKey:     "sk-provider",
		BaseURL:    "https://api.openai.com/v1",
		ModelsJSON: string(models),
		IsDefault:  true,
		Sort:       0,
	}).Error; err != nil {
		t.Fatalf("create ai provider config: %v", err)
	}

	if err := db.Create(&model.SysConfig{
		Name:      "AI配置",
		Key:       "ai_config",
		Value:     `{"default_provider":"Legacy","providers":[{"name":"Legacy","api_key":"sk-legacy","base_url":"https://legacy.example/v1","models":[{"id":"legacy-default","name":"Legacy Default","description":""}]}]}`,
		ValueType: "json",
		Remark:    "legacy ai config",
	}).Error; err != nil {
		t.Fatalf("create legacy ai_config: %v", err)
	}

	conversation, err := Default.CreateConversation(1, CreateConversationInput{})
	if err != nil {
		t.Fatalf("CreateConversation error: %v", err)
	}
	if conversation.Model != "provider-default" {
		t.Fatalf("conversation model = %s, want provider-default", conversation.Model)
	}
}

func TestAIServiceResolveProviderUsesProvidersInsteadOfLegacySysConfig(t *testing.T) {
	db := testutil.SetupStorageServiceTestDB(t)

	models, err := json.Marshal([]appconfig.AIModel{
		{ID: "provider-model", Name: "Provider Model", Description: ""},
	})
	if err != nil {
		t.Fatalf("marshal provider models: %v", err)
	}
	if err := db.Create(&model.AIProviderConfig{
		Name:       "OpenAI",
		APIKey:     "sk-provider",
		BaseURL:    "https://api.openai.com/v1",
		ModelsJSON: string(models),
		IsDefault:  true,
		Sort:       0,
	}).Error; err != nil {
		t.Fatalf("create ai provider config: %v", err)
	}

	if err := db.Create(&model.SysConfig{
		Name:      "AI配置",
		Key:       "ai_config",
		Value:     `{"default_provider":"Legacy","providers":[{"name":"Legacy","api_key":"sk-legacy","base_url":"https://legacy.example/v1","models":[{"id":"provider-model","name":"Legacy Override","description":""}]}]}`,
		ValueType: "json",
		Remark:    "legacy ai config",
	}).Error; err != nil {
		t.Fatalf("create legacy ai_config: %v", err)
	}

	provider, err := Default.resolveProvider("provider-model")
	if err != nil {
		t.Fatalf("resolveProvider error: %v", err)
	}
	if provider.Name != "OpenAI" || provider.APIKey != "sk-provider" {
		t.Fatalf("provider = %#v, want OpenAI from ai_providers", provider)
	}
}
