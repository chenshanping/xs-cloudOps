package initialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigFromPath(t *testing.T) {
	tempDir := t.TempDir()
	configContent := []byte(`server:
  host: custom-host:9100
  port: 9100
  mode: release
`)

	testCases := []struct {
		name       string
		configPath string
	}{
		{
			name:       "yaml extension",
			configPath: filepath.Join(tempDir, "custom-config.yaml"),
		},
		{
			name:       "no extension yaml content",
			configPath: filepath.Join(tempDir, "custom-config"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(tc.configPath, configContent, 0o600); err != nil {
				t.Fatalf("write temp config: %v", err)
			}

			cfg, vp, err := LoadConfig(tc.configPath)
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
			if vp.ConfigFileUsed() != tc.configPath {
				t.Fatalf("expected config file %q, got %q", tc.configPath, vp.ConfigFileUsed())
			}
		})
	}
}

func TestInitConfigReturnsErrorForMissingConfig(t *testing.T) {
	missingConfigPath := filepath.Join(t.TempDir(), "missing-config")

	err := InitConfig(missingConfigPath)
	if err == nil {
		t.Fatal("expected missing config to return an error")
	}
	if !strings.Contains(err.Error(), "读取配置文件失败") {
		t.Fatalf("expected read failure message, got %q", err.Error())
	}
}
