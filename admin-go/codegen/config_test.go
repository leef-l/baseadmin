package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigResolvesRelativeOutputPathsFromConfigDir(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, "admin-go", "codegen")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}

	configPath := filepath.Join(configDir, "codegen.yaml")
	content := `
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: secret
  dbname: demo
backend:
  output: ../app
frontend:
  output: ../../vue-vben-admin/apps/web-antd/src
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if got, want := cfg.Backend.Output, filepath.Join(root, "admin-go", "app"); got != want {
		t.Fatalf("backend output mismatch: got=%q want=%q", got, want)
	}
	if got, want := cfg.Frontend.Output, filepath.Join(root, "vue-vben-admin", "apps", "web-antd", "src"); got != want {
		t.Fatalf("frontend output mismatch: got=%q want=%q", got, want)
	}
	if got, want := cfg.AllowedApps, defaultAllowedApps(); len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("allowed apps mismatch: got=%v want=%v", got, want)
	}
}

func TestResolveConfigPathKeepsAbsolutePath(t *testing.T) {
	configDir := filepath.Join(t.TempDir(), "admin-go", "codegen")
	absolute := filepath.Join(t.TempDir(), "already-absolute")

	if got := resolveConfigPath(configDir, absolute); got != absolute {
		t.Fatalf("resolveConfigPath should keep absolute path, got %q", got)
	}
}

func TestLoadConfigRejectsMenuAppsOutsideAllowedScope(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, "admin-go", "codegen")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}

	configPath := filepath.Join(configDir, "codegen.yaml")
	content := `
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: secret
  dbname: demo
backend:
  output: ../app
frontend:
  output: ../../vue-vben-admin/apps/web-antd/src
allowed_apps:
  - system
  - upload
menu_apps:
  demo:
    title: 演示
    icon: AppstoreOutlined
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil || !strings.Contains(err.Error(), "menu_apps 包含未授权应用") {
		t.Fatalf("unexpected error: %v", err)
	}
}
