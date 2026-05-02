package initialize

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromPath(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "custom-config.yaml")

	configContent := []byte(`server:
  host: custom-host:9100
  port: 9100
  mode: release
`)

	if err := os.WriteFile(configPath, configContent, 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	cfg, vp, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("load config from explicit path: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config to be loaded")
	}
	if vp == nil {
		t.Fatal("expected viper instance to be returned")
	}
	if cfg.Server.Host != "custom-host:9100" {
		t.Fatalf("expected host to be loaded from temp config, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 9100 {
		t.Fatalf("expected port 9100, got %d", cfg.Server.Port)
	}
	if vp.ConfigFileUsed() != configPath {
		t.Fatalf("expected config file %q, got %q", configPath, vp.ConfigFileUsed())
	}
}
