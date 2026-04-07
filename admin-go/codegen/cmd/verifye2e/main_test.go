package main

import (
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
