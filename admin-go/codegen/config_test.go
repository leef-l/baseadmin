package main

import (
	"os"
	"path/filepath"
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
}

func TestResolveConfigPathKeepsAbsolutePath(t *testing.T) {
	configDir := filepath.Join(t.TempDir(), "admin-go", "codegen")
	absolute := filepath.Join(t.TempDir(), "already-absolute")

	if got := resolveConfigPath(configDir, absolute); got != absolute {
		t.Fatalf("resolveConfigPath should keep absolute path, got %q", got)
	}
}
