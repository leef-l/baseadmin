package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gbaseadmin/codegen/generator/menu"
	"gbaseadmin/codegen/generator/util"
)

func TestValidateOnlyFlag(t *testing.T) {
	valid := []string{"", "backend", "frontend", "menu"}
	for _, value := range valid {
		if err := validateOnlyFlag(value); err != nil {
			t.Fatalf("validateOnlyFlag(%q) should succeed: %v", value, err)
		}
	}

	if err := validateOnlyFlag("all"); err == nil {
		t.Fatal("validateOnlyFlag should reject unknown value")
	}
}

func TestFailureCollectorTracksFailures(t *testing.T) {
	var collector failureCollector

	collector.Add("table system_dept", nil)
	if collector.HasFailures() {
		t.Fatal("collector should ignore nil errors")
	}

	collector.Add("table system_dept", errors.New("render failed"))
	if !collector.HasFailures() {
		t.Fatal("collector should report failures")
	}

	if len(collector.items) != 1 {
		t.Fatalf("unexpected failure count: %d", len(collector.items))
	}

	if got, want := collector.items[0], "table system_dept: render failed"; got != want {
		t.Fatalf("unexpected failure item: got=%q want=%q", got, want)
	}
}

func TestParseTableNames(t *testing.T) {
	got, err := parseTableNames(" system_dept, ,system_role,system_dept , upload_dir ")
	if err != nil {
		t.Fatalf("parseTableNames failed: %v", err)
	}
	want := []string{"system_dept", "system_role", "upload_dir"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseTableNames mismatch: got=%v want=%v", got, want)
	}

	if _, err := parseTableNames(" , , "); err == nil {
		t.Fatal("parseTableNames should reject empty input")
	}
}

func TestBuildMenuGeneratorConfigCopiesMaps(t *testing.T) {
	isShow := 0
	cfg := &Config{
		MenuApps: map[string]MenuAppConfig{
			"system": {Title: "系统管理", Icon: "SettingOutlined", Sort: 10},
		},
		MenuModules: map[string]MenuModuleConfig{
			"system/dept": {Sort: 20, IsShow: &isShow, Icon: "ApartmentOutlined"},
		},
	}

	menuCfg := buildMenuGeneratorConfig(cfg, true, true)
	if !menuCfg.Force || !menuCfg.DryRun {
		t.Fatalf("unexpected flags: %+v", menuCfg)
	}
	if got := menuCfg.MenuApps["system"]; got.Title != "系统管理" || got.Icon != "SettingOutlined" || got.Sort != 10 {
		t.Fatalf("unexpected menu app config: %+v", got)
	}
	if got := menuCfg.MenuModules["system/dept"]; got.Sort != 20 || got.Icon != "ApartmentOutlined" || got.IsShow == nil || *got.IsShow != 0 {
		t.Fatalf("unexpected menu module config: %+v", got)
	}

	menuCfg.MenuApps["system"] = menu.MenuAppConfig{Title: "changed"}
	if cfg.MenuApps["system"].Title != "系统管理" {
		t.Fatal("buildMenuGeneratorConfig should not mutate source config maps")
	}
}

func TestNormalizeTableName(t *testing.T) {
	tests := []struct {
		appName string
		input   string
		want    string
	}{
		{appName: "system", input: "dept", want: "system_dept"},
		{appName: "system", input: "system_role", want: "system_role"},
		{appName: "upload", input: " upload_dir ", want: "upload_dir"},
		{appName: "demo", input: `"tag"`, want: "demo_tag"},
		{appName: "", input: "dept", want: "dept"},
	}

	for _, tc := range tests {
		if got := normalizeTableName(tc.appName, tc.input); got != tc.want {
			t.Fatalf("normalizeTableName(%q, %q) = %q, want %q", tc.appName, tc.input, got, tc.want)
		}
	}
}

func TestExtractDAOFileTableName(t *testing.T) {
	dir := t.TempDir()
	daoPath := filepath.Join(dir, "dept.go")
	content := `func NewDeptDao() *DeptDao {
	return &DeptDao{
		table:"system_dept",
	}
}`
	if err := os.WriteFile(daoPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write dao file: %v", err)
	}

	got, err := extractDAOFileTableName(daoPath)
	if err != nil {
		t.Fatalf("extractDAOFileTableName failed: %v", err)
	}
	if got != "system_dept" {
		t.Fatalf("unexpected table name: %q", got)
	}
}

func TestExtractDAOFileTableNameFallsBackToComment(t *testing.T) {
	dir := t.TempDir()
	daoPath := filepath.Join(dir, "users.go")
	content := `// UsersDao is the data access object for the table system_users.
type UsersDao struct{}`
	if err := os.WriteFile(daoPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write dao file: %v", err)
	}

	got, err := extractDAOFileTableName(daoPath)
	if err != nil {
		t.Fatalf("extractDAOFileTableName failed: %v", err)
	}
	if got != "system_users" {
		t.Fatalf("unexpected table name: %q", got)
	}
}

func TestScanExistingModulesIncludesLogicAndController(t *testing.T) {
	appDir := t.TempDir()
	logicDir := filepath.Join(appDir, "internal", "logic", "dept")
	controllerDir := filepath.Join(appDir, "internal", "controller", "role")
	emptyControllerDir := filepath.Join(appDir, "internal", "controller", "empty")
	if err := os.MkdirAll(logicDir, 0o755); err != nil {
		t.Fatalf("mkdir logic dir: %v", err)
	}
	if err := os.MkdirAll(controllerDir, 0o755); err != nil {
		t.Fatalf("mkdir controller dir: %v", err)
	}
	if err := os.MkdirAll(emptyControllerDir, 0o755); err != nil {
		t.Fatalf("mkdir empty controller dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(logicDir, "dept.go"), []byte("package dept\n"), 0o644); err != nil {
		t.Fatalf("write logic file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(controllerDir, "role.go"), []byte("package role\n"), 0o644); err != nil {
		t.Fatalf("write controller file: %v", err)
	}

	got := scanExistingModules(appDir, []string{"users"})
	want := []string{"dept", "role", "users"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("scanExistingModules mismatch: got=%v want=%v", got, want)
	}
}

func TestScanExistingTablesNormalizesLegacyHackConfigAndDAOFiles(t *testing.T) {
	appDir := filepath.Join(t.TempDir(), "system")
	if err := os.MkdirAll(filepath.Join(appDir, "hack"), 0o755); err != nil {
		t.Fatalf("mkdir hack: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(appDir, "internal", "dao", "internal"), 0o755); err != nil {
		t.Fatalf("mkdir dao internal: %v", err)
	}

	hackContent := "tables: \"dept,role_menu,system_users\"\n"
	if err := os.WriteFile(filepath.Join(appDir, "hack", "config.yaml"), []byte(hackContent), 0o644); err != nil {
		t.Fatalf("write hack config: %v", err)
	}
	daoContent := `func NewUsersDao() *UsersDao {
	return &UsersDao{
		table:    "system_users",
	}
}`
	if err := os.WriteFile(filepath.Join(appDir, "internal", "dao", "internal", "users.go"), []byte(daoContent), 0o644); err != nil {
		t.Fatalf("write dao file: %v", err)
	}

	got := scanExistingTables(appDir, []string{"system_dept", "system_role"})
	want := []string{"system_dept", "system_role", "system_role_menu", "system_users"}
	if len(got) != len(want) {
		t.Fatalf("scanExistingTables len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("scanExistingTables mismatch: got=%v want=%v", got, want)
		}
	}
}

func TestReadHackConfigTablesSupportsYaml(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := `
gfcli:
  gen:
    dao:
      - tables: "dept, role_menu"
      - tables: "system_users"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write hack config: %v", err)
	}

	got := readHackConfigTables(path)
	want := []string{"dept", "role_menu", "system_users"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("readHackConfigTables mismatch: got=%v want=%v", got, want)
	}
}

func TestReadHackConfigTablesFallbackCollectsMultipleLegacyLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.txt")
	content := `
tables: "dept, role_menu"
other: keep
tables: "system_users,role_menu"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write hack config: %v", err)
	}

	got := readHackConfigTables(path)
	want := []string{"dept", "role_menu", "system_users"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("legacy fallback mismatch: got=%v want=%v", got, want)
	}
}

func TestSplitCSVValuesTrimsInlineHashComments(t *testing.T) {
	got := splitCSVValues(`"dept, role_menu" # legacy comment`)
	want := []string{"dept", "role_menu"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitCSVValues mismatch: got=%v want=%v", got, want)
	}

	got = splitCSVValues(`dept,system_users # keep comment out`)
	want = []string{"dept", "system_users"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitCSVValues unquoted mismatch: got=%v want=%v", got, want)
	}
}

func TestRenderTemplateSkipReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	tplPath := filepath.Join(dir, "simple.tpl")
	outPath := filepath.Join(dir, "out.txt")
	if err := os.WriteFile(tplPath, []byte("hello {{.Name}}"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	if err := os.WriteFile(outPath, []byte("existing"), 0o644); err != nil {
		t.Fatalf("write output: %v", err)
	}

	written, err := renderTemplate(tplPath, outPath, map[string]string{"Name": "world"}, false, util.NewTemplateCache())
	if err != nil {
		t.Fatalf("renderTemplate failed: %v", err)
	}
	if written {
		t.Fatal("renderTemplate should report skipped write when overwrite is false")
	}

	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(content) != "existing" {
		t.Fatalf("output should remain unchanged, got %q", string(content))
	}
}

func TestWriteFileIfAbsentRejectsExistingDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "target")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir target dir: %v", err)
	}

	written, err := writeFileIfAbsent(dir, []byte("x"))
	if err == nil {
		t.Fatal("writeFileIfAbsent should reject existing directory")
	}
	if written {
		t.Fatal("writeFileIfAbsent should not report written for directory target")
	}
}

func TestRenderTemplateRejectsExistingDirectory(t *testing.T) {
	dir := t.TempDir()
	tplPath := filepath.Join(dir, "simple.tpl")
	outDir := filepath.Join(dir, "existing-dir")
	if err := os.WriteFile(tplPath, []byte("hello {{.Name}}"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}

	written, err := renderTemplate(tplPath, outDir, map[string]string{"Name": "world"}, false, util.NewTemplateCache())
	if err == nil {
		t.Fatal("renderTemplate should reject existing directory output")
	}
	if written {
		t.Fatal("renderTemplate should not report write when output path is a directory")
	}
}

func TestRenderTemplateSkipsUnchangedOverwrite(t *testing.T) {
	dir := t.TempDir()
	tplPath := filepath.Join(dir, "simple.tpl")
	outPath := filepath.Join(dir, "out.txt")
	if err := os.WriteFile(tplPath, []byte("hello {{.Name}}"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	if err := os.WriteFile(outPath, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("write output: %v", err)
	}

	written, err := renderTemplate(tplPath, outPath, map[string]string{"Name": "world"}, true, util.NewTemplateCache())
	if err != nil {
		t.Fatalf("renderTemplate failed: %v", err)
	}
	if written {
		t.Fatal("renderTemplate should skip identical content even when overwrite is true")
	}
}
