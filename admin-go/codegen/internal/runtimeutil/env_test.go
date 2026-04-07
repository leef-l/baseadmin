package runtimeutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadEnvFileIfExistsDoesNotOverrideExisting(t *testing.T) {
	t.Setenv("MYSQL_HOST", "existing-host")
	t.Setenv("MYSQL_PORT", "")

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("MYSQL_HOST=file-host\nMYSQL_PORT=40001\n"), 0o644); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	if err := LoadEnvFileIfExists(envPath); err != nil {
		t.Fatalf("LoadEnvFileIfExists failed: %v", err)
	}

	if got := os.Getenv("MYSQL_HOST"); got != "existing-host" {
		t.Fatalf("MYSQL_HOST should keep existing value, got %q", got)
	}
	if got := os.Getenv("MYSQL_PORT"); got != "" {
		t.Fatalf("MYSQL_PORT should keep existing empty env, got %q", got)
	}

	os.Unsetenv("MYSQL_PORT")
	if err := LoadEnvFileIfExists(envPath); err != nil {
		t.Fatalf("LoadEnvFileIfExists failed on second pass: %v", err)
	}
	if got := os.Getenv("MYSQL_PORT"); got != "40001" {
		t.Fatalf("MYSQL_PORT should be loaded from env file, got %q", got)
	}
}

func TestLoadEnvFileIfExistsSupportsExportAndQuotes(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "# comment\nexport MYSQL_DATABASE=\"demo\"\nMYSQL_PASSWORD='secret'\nINVALID_LINE\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	if err := LoadEnvFileIfExists(envPath); err != nil {
		t.Fatalf("LoadEnvFileIfExists failed: %v", err)
	}

	if got := os.Getenv("MYSQL_DATABASE"); got != "demo" {
		t.Fatalf("MYSQL_DATABASE mismatch, got %q", got)
	}
	if got := os.Getenv("MYSQL_PASSWORD"); got != "secret" {
		t.Fatalf("MYSQL_PASSWORD mismatch, got %q", got)
	}
}

func TestLoadEnvFileIfExistsTrimsInlineCommentsForUnquotedValues(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "MYSQL_HOST=127.0.0.1 # local\nMYSQL_TAG=v1#stable\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	if err := LoadEnvFileIfExists(envPath); err != nil {
		t.Fatalf("LoadEnvFileIfExists failed: %v", err)
	}

	if got := os.Getenv("MYSQL_HOST"); got != "127.0.0.1" {
		t.Fatalf("MYSQL_HOST mismatch, got %q", got)
	}
	if got := os.Getenv("MYSQL_TAG"); got != "v1#stable" {
		t.Fatalf("MYSQL_TAG should keep literal # when not a comment, got %q", got)
	}
}

func TestLoadEnvFileIfExistsSupportsLongValues(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	longValue := strings.Repeat("x", 70*1024)
	content := "LONG_TOKEN=" + longValue + "\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	if err := LoadEnvFileIfExists(envPath); err != nil {
		t.Fatalf("LoadEnvFileIfExists failed: %v", err)
	}

	if got := os.Getenv("LONG_TOKEN"); got != longValue {
		t.Fatalf("LONG_TOKEN mismatch, len=%d want=%d", len(got), len(longValue))
	}
}
