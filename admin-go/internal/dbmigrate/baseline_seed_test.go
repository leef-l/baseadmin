package dbmigrate

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBaselineCoreMenusGrantedToAdminRole(t *testing.T) {
	content, err := readBaselineMigration()
	if err != nil {
		t.Fatalf("readBaselineMigration failed: %v", err)
	}
	requiredMenuRows := []string{
		"(1000000000000000010, 0, '系统管理'",
		"(314253730209861632, 0, '上传管理'",
		"'system:dept:batch-delete'",
		"'system:role:batch-delete'",
		"'system:menu:batch-delete'",
		"'system:user:batch-delete'",
		"'upload:config:batch-delete'",
		"'upload:dir:batch-delete'",
		"'upload:dir_rule:batch-delete'",
		"'upload:file:batch-delete'",
	}
	for _, row := range requiredMenuRows {
		if !strings.Contains(content, row) {
			t.Fatalf("baseline migration missing menu seed: %s", row)
		}
	}
	requiredRoleMenuRows := []string{
		"(1000000000000000002, 1000000000000000010)",
		"(1000000000000000002, 314253730209861632)",
	}
	for _, row := range requiredRoleMenuRows {
		if !strings.Contains(content, row) {
			t.Fatalf("baseline migration missing admin role grant: %s", row)
		}
	}
	forbiddenRows := []string{
		"'仪表盘'",
		"'分析页'",
		"'工作台'",
	}
	for _, row := range forbiddenRows {
		if strings.Contains(content, row) {
			t.Fatalf("baseline migration should not contain removed seed: %s", row)
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
