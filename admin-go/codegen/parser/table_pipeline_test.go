package parser

import (
	"strings"
	"testing"
)

func TestSplitTableIdentity(t *testing.T) {
	cases := []struct {
		tableName  string
		wantApp    string
		wantModule string
	}{
		{tableName: "system_dept", wantApp: "system", wantModule: "dept"},
		{tableName: "upload_dir_rule", wantApp: "upload", wantModule: "dir_rule"},
		{tableName: "users", wantApp: "", wantModule: "users"},
	}

	for _, tc := range cases {
		got := splitTableIdentity(tc.tableName)
		if got.appName != tc.wantApp || got.moduleName != tc.wantModule {
			t.Fatalf("splitTableIdentity(%q) mismatch: got=%+v", tc.tableName, got)
		}
	}
}

func TestValidateRequiredScopeFields(t *testing.T) {
	meta := &TableMeta{
		TableName:     "demo_article",
		HasTenantID:   true,
		HasMerchantID: true,
		HasCreatedBy:  true,
		HasDeptID:     true,
	}
	if err := validateRequiredScopeFields(meta); err != nil {
		t.Fatalf("validateRequiredScopeFields returned unexpected error: %v", err)
	}
}

func TestValidateRequiredScopeFieldsRejectsMissingField(t *testing.T) {
	meta := &TableMeta{
		TableName:    "demo_article",
		HasTenantID:  true,
		HasCreatedBy: true,
	}
	err := validateRequiredScopeFields(meta)
	if err == nil {
		t.Fatal("validateRequiredScopeFields should reject missing merchant_id and dept_id")
	}
	if !strings.Contains(err.Error(), "merchant_id") || !strings.Contains(err.Error(), "dept_id") {
		t.Fatalf("error should mention missing merchant_id and dept_id, got: %v", err)
	}
}

func TestResolveReferenceFieldsPrefersAppPrefixedTable(t *testing.T) {
	p := &Parser{
		tableColumnsCache: map[string]map[string]struct{}{
			"system_role": {
				"id":         {},
				"name":       {},
				"parent_id":  {},
				"deleted_at": {},
			},
			"role": {
				"id":    {},
				"title": {},
			},
		},
	}

	identity := splitTableIdentity("system_user")
	meta := buildTableMetaSkeleton(identity, "用户表")
	appendColumnFields(meta, []columnInfo{{
		ColumnName:    "role_id",
		DataType:      "bigint",
		ColumnType:    "bigint(20)",
		IsNullable:    "NO",
		ColumnComment: "角色ID",
	}}, nil)

	if err := p.resolveReferenceFields(meta, identity); err != nil {
		t.Fatalf("resolveReferenceFields failed: %v", err)
	}

	field := meta.Fields[0]
	if field.RefTable != "role" || field.RefTableDB != "system_role" {
		t.Fatalf("unexpected reference target: %+v", field)
	}
	if field.RefDisplayField != "name" || field.RefFieldName != "RoleName" || field.RefFieldJSON != "roleName" {
		t.Fatalf("unexpected reference display mapping: %+v", field)
	}
	if field.RefTableApp != "system" || !field.RefIsTree || !field.RefHasDeletedAt {
		t.Fatalf("unexpected reference flags: %+v", field)
	}
}

func TestResolveReferenceFieldsUsesHints(t *testing.T) {
	p := &Parser{
		tableColumnsCache: map[string]map[string]struct{}{
			"system_users": {
				"id":         {},
				"username":   {},
				"deleted_at": {},
			},
		},
	}

	identity := splitTableIdentity("audit_log")
	meta := buildTableMetaSkeleton(identity, "日志表")
	appendColumnFields(meta, []columnInfo{{
		ColumnName:    "user_id",
		DataType:      "bigint",
		ColumnType:    "bigint(20)",
		IsNullable:    "NO",
		ColumnComment: "创建人|ref:system_users.username",
	}}, nil)

	if err := p.resolveReferenceFields(meta, identity); err != nil {
		t.Fatalf("resolveReferenceFields failed: %v", err)
	}

	field := meta.Fields[0]
	if field.RefTable != "users" || field.RefTableDB != "system_users" || field.RefTableApp != "system" {
		t.Fatalf("unexpected hinted reference target: %+v", field)
	}
	if field.RefDisplayField != "username" || field.RefFieldName != "UsersUsername" {
		t.Fatalf("unexpected hinted display mapping: %+v", field)
	}
	if field.RefIsTree || !field.RefHasDeletedAt {
		t.Fatalf("unexpected hinted flags: %+v", field)
	}
}

func TestResolveReferenceFieldsRejectsMissingDisplayField(t *testing.T) {
	p := &Parser{
		tableColumnsCache: map[string]map[string]struct{}{
			"system_role": {"id": {}},
			"role":        {"id": {}},
		},
	}

	identity := splitTableIdentity("system_user")
	meta := buildTableMetaSkeleton(identity, "用户表")
	appendColumnFields(meta, []columnInfo{{
		ColumnName:    "role_id",
		DataType:      "bigint",
		ColumnType:    "bigint(20)",
		IsNullable:    "NO",
		ColumnComment: "角色ID",
	}}, nil)

	err := p.resolveReferenceFields(meta, identity)
	if err == nil {
		t.Fatal("resolveReferenceFields should fail when display field is missing")
	}
	if !strings.Contains(err.Error(), "字段 role_id 是外键") {
		t.Fatalf("unexpected error: %v", err)
	}
}
