package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadVerifyConfigRequiresFields(t *testing.T) {
	t.Setenv("MYSQL_HOST", "127.0.0.1")
	t.Setenv("MYSQL_PORT", "")
	t.Setenv("MYSQL_USER", "root")
	t.Setenv("MYSQL_PASSWORD", "secret")
	t.Setenv("MYSQL_DATABASE", "demo")

	if _, err := loadVerifyConfig(); err == nil {
		t.Fatal("expected missing port to fail loadVerifyConfig")
	}
}

func TestLoadVerifyConfigSuccess(t *testing.T) {
	t.Setenv("MYSQL_HOST", "127.0.0.1")
	t.Setenv("MYSQL_PORT", "3306")
	t.Setenv("MYSQL_USER", "root")
	t.Setenv("MYSQL_PASSWORD", "secret")
	t.Setenv("MYSQL_DATABASE", "demo")

	cfg, err := loadVerifyConfig()
	if err != nil {
		t.Fatalf("loadVerifyConfig failed: %v", err)
	}
	if cfg.Host != "127.0.0.1" || cfg.Port != "3306" || cfg.User != "root" || cfg.DBName != "demo" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestParseVerifyOptionsDefaultsToAll(t *testing.T) {
	opts, err := parseVerifyOptions(nil)
	if err != nil {
		t.Fatalf("parseVerifyOptions failed: %v", err)
	}
	if opts.Stage != verifyStageAll || opts.TempRoot != "" || opts.KeepTemp {
		t.Fatalf("unexpected default options: %+v", opts)
	}
}

func TestValidateVerifyOptionsRequiresTempRootForHeavyStages(t *testing.T) {
	if err := validateVerifyOptions(verifyOptions{Stage: verifyStageDAO}); err == nil {
		t.Fatal("dao stage should require temp root")
	}
	if err := validateVerifyOptions(verifyOptions{Stage: verifyStageBuild}); err == nil {
		t.Fatal("build stage should require temp root")
	}
	if err := validateVerifyOptions(verifyOptions{Stage: verifyStageRender}); err != nil {
		t.Fatalf("render stage should succeed: %v", err)
	}
}

func TestGeneratedAppNames(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module demo\n"), 0o644); err != nil {
		t.Fatalf("write go.mod failed: %v", err)
	}
	for _, appName := range []string{"verifydemo", "sample"} {
		if err := os.MkdirAll(filepath.Join(root, "app", appName), 0o755); err != nil {
			t.Fatalf("mkdir app failed: %v", err)
		}
	}

	got, err := generatedAppNames(root)
	if err != nil {
		t.Fatalf("generatedAppNames failed: %v", err)
	}
	want := []string{"sample", "verifydemo"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("generatedAppNames mismatch: got=%v want=%v", got, want)
	}
}
