package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")

	configContent := `
server:
  port: 8080
  mode: debug
  trusted_proxies:
    - "127.0.0.1"

database:
  driver: sqlite
  dsn: test.db

jwt:
  secret: test-secret-key
  expire_hour: 24

upload:
  path: ./uploads
  max_size: 10485760
  allow_ext:
    - ".jpg"
    - ".png"

site:
  name: "Test Forum"
  brief: "A test forum"
  page_size: 20

email:
  enabled: false
  host: smtp.example.com
  port: 587
  username: test@example.com
  password: password
  from: test@example.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want 8080", cfg.Server.Port)
	}
	if cfg.Server.Mode != "debug" {
		t.Errorf("Server.Mode = %q, want %q", cfg.Server.Mode, "debug")
	}
	if len(cfg.Server.TrustedProxies) != 1 || cfg.Server.TrustedProxies[0] != "127.0.0.1" {
		t.Errorf("Server.TrustedProxies = %v, want [127.0.0.1]", cfg.Server.TrustedProxies)
	}
	if cfg.Database.Driver != "sqlite" {
		t.Errorf("Database.Driver = %q, want %q", cfg.Database.Driver, "sqlite")
	}
	if cfg.Database.DSN != "test.db" {
		t.Errorf("Database.DSN = %q, want %q", cfg.Database.DSN, "test.db")
	}
	if cfg.JWT.Secret != "test-secret-key" {
		t.Errorf("JWT.Secret = %q, want %q", cfg.JWT.Secret, "test-secret-key")
	}
	if cfg.JWT.ExpireHour != 24 {
		t.Errorf("JWT.ExpireHour = %d, want 24", cfg.JWT.ExpireHour)
	}
	if cfg.Site.Name != "Test Forum" {
		t.Errorf("Site.Name = %q, want %q", cfg.Site.Name, "Test Forum")
	}
	if cfg.Site.PageSize != 20 {
		t.Errorf("Site.PageSize = %d, want 20", cfg.Site.PageSize)
	}
	if Global == nil {
		t.Error("Global config is nil after Load()")
	}
}

func TestLoadInvalidPath(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Load() expected error for nonexistent file, got nil")
	}
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")

	initialContent := `
server:
  port: 8080
  mode: debug

database:
  driver: sqlite
  dsn: test.db

jwt:
  secret: test-secret
  expire_hour: 24

upload:
  path: ./uploads
  max_size: 10485760
  allow_ext: []

site:
  name: "Original Name"
  brief: "Original Brief"
  page_size: 20

email:
  enabled: false
  host: ""
  port: 0
  username: ""
  password: ""
  from: ""
`
	if err := os.WriteFile(configPath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	cfg.Site.Name = "Updated Name"
	cfg.Site.Brief = "Updated Brief"
	cfg.Site.PageSize = 30

	if err := Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	newCfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() after save error = %v", err)
	}

	if newCfg.Site.Name != "Updated Name" {
		t.Errorf("Site.Name after save = %q, want %q", newCfg.Site.Name, "Updated Name")
	}
	if newCfg.Site.Brief != "Updated Brief" {
		t.Errorf("Site.Brief after save = %q, want %q", newCfg.Site.Brief, "Updated Brief")
	}
	if newCfg.Site.PageSize != 30 {
		t.Errorf("Site.PageSize after save = %d, want 30", newCfg.Site.PageSize)
	}
}
