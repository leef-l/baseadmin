package dbmigrate

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBaselineDashboardMenusGrantedToAdminRole(t *testing.T) {
	content, err := readBaselineMigration()
	if err != nil {
		t.Fatalf("readBaselineMigration failed: %v", err)
	}
	requiredMenuRows := []string{
		"(1000000000000000060, 0, '仪表盘'",
		"(1000000000000000061, 1000000000000000060, '分析页'",
		"(1000000000000000062, 1000000000000000060, '工作台'",
	}
	for _, row := range requiredMenuRows {
		if !strings.Contains(content, row) {
			t.Fatalf("baseline migration missing menu seed: %s", row)
		}
	}
	requiredRoleMenuRows := []string{
		"(1000000000000000002, 1000000000000000060)",
		"(1000000000000000002, 1000000000000000061)",
		"(1000000000000000002, 1000000000000000062)",
	}
	for _, row := range requiredRoleMenuRows {
		if !strings.Contains(content, row) {
			t.Fatalf("baseline migration missing admin role grant: %s", row)
		}
	}
}

func readBaselineMigration() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if ok {
		path := filepath.Join(filepath.Dir(currentFile), "..", "..", "database", "migrations", "000001_baseline_system_upload.up.sql")
		if data, err := os.ReadFile(path); err == nil {
			return string(data), nil
		}
	}
	path := filepath.Join("database", "migrations", "000001_baseline_system_upload.up.sql")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
