package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gbaseadmin/codegen/parser"
)

type scopeContract struct {
	AllowedApps         []string `json:"allowedApps"`
	SupportedComponents []string `json:"supportedComponents"`
}

func loadScopeContract(t *testing.T) scopeContract {
	t.Helper()
	path := filepath.Join("..", "..", "contracts", "baseadmin-scope.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read scope contract: %v", err)
	}
	var contract scopeContract
	if err := json.Unmarshal(data, &contract); err != nil {
		t.Fatalf("unmarshal scope contract: %v", err)
	}
	return contract
}

func TestDefaultAllowedAppsMatchScopeContract(t *testing.T) {
	contract := loadScopeContract(t)
	if got := defaultAllowedApps(); !reflect.DeepEqual(got, contract.AllowedApps) {
		t.Fatalf("defaultAllowedApps mismatch: got=%v want=%v", got, contract.AllowedApps)
	}
}

func TestSupportedComponentsMatchScopeContract(t *testing.T) {
	contract := loadScopeContract(t)
	if got := parser.SupportedComponentNames(); !reflect.DeepEqual(got, contract.SupportedComponents) {
		t.Fatalf("supported components mismatch: got=%v want=%v", got, contract.SupportedComponents)
	}
}

func TestValidateMetaScopeRejectsUnexpectedApp(t *testing.T) {
	cfg := &Config{AllowedApps: defaultAllowedApps()}
	meta := &parser.TableMeta{AppName: "demo"}
	err := validateMetaScope(cfg, meta)
	if err == nil || !strings.Contains(err.Error(), "只允许生成应用") {
		t.Fatalf("validateMetaScope unexpected error: %v", err)
	}
}

func TestValidateMetaScopeRejectsDictFields(t *testing.T) {
	cfg := &Config{AllowedApps: defaultAllowedApps()}
	meta := &parser.TableMeta{
		AppName: "system",
		Fields: []parser.FieldMeta{
			{Name: "status", DictType: "article_status"},
		},
	}
	err := validateMetaScope(cfg, meta)
	if err == nil || !strings.Contains(err.Error(), "未保留 system/dict 模块") {
		t.Fatalf("validateMetaScope dict error mismatch: %v", err)
	}
}

func TestValidateMetaScopeAllowsDictFieldsWhenConfigured(t *testing.T) {
	cfg := &Config{
		AllowedApps:            []string{"demo"},
		AllowMissingDictModule: true,
	}
	meta := &parser.TableMeta{
		AppName: "demo",
		Fields: []parser.FieldMeta{
			{Name: "level", DictType: "article_level"},
		},
	}
	if err := validateMetaScope(cfg, meta); err != nil {
		t.Fatalf("validateMetaScope should allow dict fields when explicitly enabled: %v", err)
	}
}

func TestValidateMetaScopeAllowsDictFieldsWhenFrontendDictModuleExists(t *testing.T) {
	root := t.TempDir()
	dictAPIPath := filepath.Join(root, "api", "system", "dict", "index.ts")
	if err := os.MkdirAll(filepath.Dir(dictAPIPath), 0o755); err != nil {
		t.Fatalf("mkdir dict api dir: %v", err)
	}
	if err := os.WriteFile(dictAPIPath, []byte("export async function getDictByType() { return []; }\n"), 0o644); err != nil {
		t.Fatalf("write dict api file: %v", err)
	}

	cfg := &Config{
		AllowedApps: []string{"demo"},
		Frontend: FrontendConfig{
			Output: root,
		},
	}
	meta := &parser.TableMeta{
		AppName: "demo",
		Fields: []parser.FieldMeta{
			{Name: "level", DictType: "article_level"},
		},
	}
	if err := validateMetaScope(cfg, meta); err != nil {
		t.Fatalf("validateMetaScope should allow dict fields when frontend dict module exists: %v", err)
	}
}
