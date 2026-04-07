package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestE2EVerifySQLCoversDemoSchemaCoreFields(t *testing.T) {
	demoSQL := mustReadCodegenSQL(t, "demo.sql")
	e2eSQL := mustReadCodegenSQL(t, "e2e_verify.sql")

	tablePairs := []struct {
		demoTable string
		e2eTable  string
		fields    []string
		snippets  []string
	}{
		{
			demoTable: "demo_category",
			e2eTable:  "verifydemo_category",
			fields:    []string{"id", "parent_id", "name", "icon", "sort", "status", "created_by", "dept_id", "created_at", "updated_at", "deleted_at"},
			snippets:  []string{"COMMENT='分类'"},
		},
		{
			demoTable: "demo_article",
			e2eTable:  "verifydemo_article",
			fields: []string{
				"id", "category_id", "user_id", "title", "order_no", "cover", "attachment_file",
				"body_content", "extra_json", "link_url", "status", "type", "is_top", "price",
				"pay_password", "sort", "icon", "email", "phone", "remark", "level", "extra_field",
				"publish_at", "expire_at", "created_by", "dept_id", "created_at", "updated_at", "deleted_at",
			},
			snippets: []string{
				"|ref:system_users.username",
				"等级:dict:article_level",
				"价格（分）",
				"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章'",
			},
		},
		{
			demoTable: "demo_tag",
			e2eTable:  "verifydemo_tag",
			fields:    []string{"id", "name", "color", "sort", "status", "created_by", "dept_id", "created_at", "updated_at", "deleted_at"},
			snippets:  []string{"COMMENT='标签'"},
		},
	}

	for _, tc := range tablePairs {
		demoStmt := extractCreateTable(t, demoSQL, tc.demoTable)
		e2eStmt := extractCreateTable(t, e2eSQL, tc.e2eTable)
		for _, field := range tc.fields {
			if !strings.Contains(demoStmt, "`"+field+"`") {
				t.Fatalf("%s should contain field %s", tc.demoTable, field)
			}
			if !strings.Contains(e2eStmt, "`"+field+"`") {
				t.Fatalf("%s should contain field %s", tc.e2eTable, field)
			}
		}
		for _, snippet := range tc.snippets {
			if !strings.Contains(demoStmt, snippet) {
				t.Fatalf("%s should contain snippet %q", tc.demoTable, snippet)
			}
			if !strings.Contains(e2eStmt, snippet) {
				t.Fatalf("%s should contain snippet %q", tc.e2eTable, snippet)
			}
		}
	}
}

func mustReadCodegenSQL(t *testing.T, name string) string {
	t.Helper()
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "sql", name))
	if err != nil {
		t.Fatalf("read %s failed: %v", name, err)
	}
	return string(content)
}

func extractCreateTable(t *testing.T, sqlText, tableName string) string {
	t.Helper()
	pattern := "(?is)CREATE TABLE IF NOT EXISTS `" + regexp.QuoteMeta(tableName) + "`\\s*\\(.*?\\)\\s*ENGINE=InnoDB.*?;"
	re := regexp.MustCompile(pattern)
	stmt := re.FindString(sqlText)
	if stmt == "" {
		t.Fatalf("create table for %s not found", tableName)
	}
	return stmt
}
