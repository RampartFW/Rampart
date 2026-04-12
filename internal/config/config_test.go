package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "rampart-config-test-*")
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "rampart.yaml")
	configData := `
server:
  listen: "127.0.0.1:8080"
backend:
  type: "nftables"
  nftables:
    tableName: "custom_rampart"
`
	os.WriteFile(configPath, []byte(configData), 0644)

	t.Run("Load_From_Path", func(t *testing.T) {
		cfg, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if cfg.Server.Listen != "127.0.0.1:8080" {
			t.Errorf("expected listen 127.0.0.1:8080, got %v", cfg.Server.Listen)
		}
		if cfg.Backend.Type != "nftables" {
			t.Errorf("expected backend nftables, got %v", cfg.Backend.Type)
		}
		if cfg.Backend.Nftables.TableName != "custom_rampart" {
			t.Errorf("expected tableName custom_rampart, got %v", cfg.Backend.Nftables.TableName)
		}
	})

	t.Run("Env_Overrides", func(t *testing.T) {
		os.Setenv("RAMPART_LISTEN", "0.0.0.0:9999")
		defer os.Unsetenv("RAMPART_LISTEN")

		cfg, _ := LoadConfig(configPath)
		if cfg.Server.Listen != "0.0.0.0:9999" {
			t.Errorf("expected env override 0.0.0.0:9999, got %v", cfg.Server.Listen)
		}
	})

	t.Run("Default_Values", func(t *testing.T) {
		cfg := DefaultConfig()
		if cfg.Server.Listen != "0.0.0.0:9443" {
			t.Errorf("wrong default listen: %v", cfg.Server.Listen)
		}
		if cfg.Backend.Type != "auto" {
			t.Errorf("wrong default backend type: %v", cfg.Backend.Type)
		}
	})
}
