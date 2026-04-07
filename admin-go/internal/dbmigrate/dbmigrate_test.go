package dbmigrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEnsureMySQLDSNParams(t *testing.T) {
	got := ensureMySQLDSNParams("user:pass@tcp(127.0.0.1:3306)/gbaseadmin?charset=utf8mb4")
	if !strings.Contains(got, "charset=utf8mb4") {
		t.Fatalf("dsn should preserve charset: %q", got)
	}
	if !strings.Contains(got, "parseTime=true") {
		t.Fatalf("dsn should add parseTime: %q", got)
	}
	if !strings.Contains(got, "multiStatements=true") {
		t.Fatalf("dsn should add multiStatements: %q", got)
	}
}

func TestEnsureMySQLDSNParamsFillsBlankValues(t *testing.T) {
	got := ensureMySQLDSNParams("user:pass@tcp(127.0.0.1:3306)/gbaseadmin?charset=&parseTime= &multiStatements=")
	if !strings.Contains(got, "charset=utf8mb4") {
		t.Fatalf("dsn should backfill blank charset: %q", got)
	}
	if !strings.Contains(got, "parseTime=true") {
		t.Fatalf("dsn should backfill blank parseTime: %q", got)
	}
	if !strings.Contains(got, "multiStatements=true") {
		t.Fatalf("dsn should backfill blank multiStatements: %q", got)
	}
}

func TestGoFrameLinkToMySQLDSN(t *testing.T) {
	got, err := goFrameLinkToMySQLDSN("mysql:user:pass@tcp(127.0.0.1:3306)/gbaseadmin")
	if err != nil {
		t.Fatalf("goFrameLinkToMySQLDSN failed: %v", err)
	}
	if !strings.HasPrefix(got, "user:pass@tcp(127.0.0.1:3306)/gbaseadmin?") {
		t.Fatalf("unexpected dsn: %q", got)
	}
}

func TestNormalizeMigrationName(t *testing.T) {
	if got := normalizeMigrationName(" Init Users Table "); got != "init_users_table" {
		t.Fatalf("normalizeMigrationName mismatch: %q", got)
	}
	if got := normalizeMigrationName("___"); got != "migration" {
		t.Fatalf("normalizeMigrationName fallback mismatch: %q", got)
	}
}

func TestLoadDatabaseLink(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := "" +
		"database:\n" +
		"  default:\n" +
		"    link: \"mysql:user:pass@tcp(127.0.0.1:3306)/demo\"\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	got, err := loadDatabaseLink(path, "default")
	if err != nil {
		t.Fatalf("loadDatabaseLink failed: %v", err)
	}
	if got != "mysql:user:pass@tcp(127.0.0.1:3306)/demo" {
		t.Fatalf("loadDatabaseLink mismatch: %q", got)
	}
}

func TestCreateMigrationFiles(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 4, 7, 12, 30, 45, 0, time.UTC)
	if err := createMigrationFiles(dir, "add user index", now); err != nil {
		t.Fatalf("createMigrationFiles failed: %v", err)
	}
	upPath := filepath.Join(dir, "20260407123045_add_user_index.up.sql")
	downPath := filepath.Join(dir, "20260407123045_add_user_index.down.sql")
	if _, err := os.Stat(upPath); err != nil {
		t.Fatalf("up migration not created: %v", err)
	}
	if _, err := os.Stat(downPath); err != nil {
		t.Fatalf("down migration not created: %v", err)
	}
}

func TestCreateMigrationFilesRejectsDuplicateVersionAndName(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 4, 7, 12, 30, 45, 0, time.UTC)
	if err := createMigrationFiles(dir, "add user index", now); err != nil {
		t.Fatalf("first createMigrationFiles failed: %v", err)
	}
	if err := createMigrationFiles(dir, "add user index", now); err == nil {
		t.Fatal("createMigrationFiles should reject duplicate files")
	}
}

func TestEnvMySQLDSN(t *testing.T) {
	t.Setenv("DB_HOST", "mysql")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_USER", "demo")
	t.Setenv("DB_PASS", "secret")
	t.Setenv("DB_NAME", "gbaseadmin")

	got := envMySQLDSN()
	if !strings.HasPrefix(got, "demo:secret@tcp(mysql:3306)/gbaseadmin?") {
		t.Fatalf("envMySQLDSN mismatch: %q", got)
	}
}

func TestResolveConfigPathSupportsAbsoluteGFConfigFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("database:\n  default:\n    link: \"mysql:user:pass@tcp(localhost:3306)/demo\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	t.Setenv("GF_GCFG_PATH", filepath.Join(dir, "ignored"))
	t.Setenv("GF_GCFG_FILE", configPath)
	got, err := resolveConfigPath("")
	if err != nil {
		t.Fatalf("resolveConfigPath failed: %v", err)
	}
	if got != configPath {
		t.Fatalf("resolveConfigPath mismatch: got=%q want=%q", got, configPath)
	}
}
