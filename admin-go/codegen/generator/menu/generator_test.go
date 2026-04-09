package menu

import (
	"database/sql"
	"reflect"
	"strings"
	"testing"

	"gbaseadmin/codegen/parser"
)

func TestBuildButtonSpecsMatchesFeatureFlags(t *testing.T) {
	meta := &parser.TableMeta{
		AppName:      "demo",
		ModuleName:   "article",
		HasImport:    true,
		HasBatchEdit: true,
	}

	got := buildButtonSpecs(meta)
	want := []buttonSpec{
		{suffix: "新增", permission: "demo:article:create", sort: 1},
		{suffix: "修改", permission: "demo:article:update", sort: 2},
		{suffix: "删除", permission: "demo:article:delete", sort: 3},
		{suffix: "批量删除", permission: "demo:article:batch-delete", sort: 4},
		{suffix: "查看", permission: "demo:article:detail", sort: 5},
		{suffix: "导出", permission: "demo:article:export", sort: 6},
		{suffix: "导入", permission: "demo:article:import", sort: 7},
		{suffix: "批量编辑", permission: "demo:article:batch-update", sort: 8},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildButtonSpecs mismatch:\nwant=%+v\ngot=%+v", want, got)
	}
}

func TestBuildButtonSpecsSkipsUnsupportedActions(t *testing.T) {
	meta := &parser.TableMeta{
		AppName:      "demo",
		ModuleName:   "category",
		HasParentID:  true,
		HasImport:    false,
		HasBatchEdit: false,
	}

	got := buildButtonSpecs(meta)
	for _, btn := range got {
		if btn.permission == "demo:category:batch-delete" || btn.permission == "demo:category:import" || btn.permission == "demo:category:batch-update" {
			t.Fatalf("unexpected button permission for tree/non-import module: %+v", got)
		}
	}
}

func TestBuildMenuTitleFallsBackWhenCommentEmpty(t *testing.T) {
	meta := &parser.TableMeta{
		ModelName:  "UserRole",
		ModuleName: "user_role",
	}
	if got := buildMenuTitle(meta); got != "UserRole" {
		t.Fatalf("buildMenuTitle fallback mismatch: %q", got)
	}
}

func TestGuessMenuIconUsesModuleDefaults(t *testing.T) {
	cases := []struct {
		name string
		meta *parser.TableMeta
		want string
	}{
		{
			name: "dir_rule",
			meta: &parser.TableMeta{ModuleName: "dir_rule", Comment: "文件目录规则表"},
			want: "PartitionOutlined",
		},
		{
			name: "users",
			meta: &parser.TableMeta{ModuleName: "users", Comment: "用户表"},
			want: "UserOutlined",
		},
		{
			name: "menu",
			meta: &parser.TableMeta{ModuleName: "menu", Comment: "菜单表"},
			want: "MenuOutlined",
		},
	}

	for _, tt := range cases {
		if got := guessMenuIcon(tt.meta); got != tt.want {
			t.Fatalf("%s icon mismatch: want=%s got=%s", tt.name, tt.want, got)
		}
	}
}

func TestGuessMenuIconFallsBackToKeywords(t *testing.T) {
	meta := &parser.TableMeta{
		ModuleName: "profile_entry",
		Comment:    "审计日志表",
	}
	if got := guessMenuIcon(meta); got != "ProfileOutlined" {
		t.Fatalf("guessMenuIcon keyword mismatch: %q", got)
	}
}

func TestGuessMenuIconFallsBackToDefault(t *testing.T) {
	meta := &parser.TableMeta{
		ModuleName: "widget",
		Comment:    "控件中心",
	}
	if got := guessMenuIcon(meta); got != "AppstoreOutlined" {
		t.Fatalf("guessMenuIcon default mismatch: %q", got)
	}
}

func TestCleanTitleTrimsRepeatedSuffixes(t *testing.T) {
	if got := cleanTitle("标签管理表 "); got != "标签" {
		t.Fatalf("cleanTitle mismatch: %q", got)
	}
}

func TestDashCaseReplacesUnderscores(t *testing.T) {
	if got := dashCase("user_role_detail"); got != "user-role-detail" {
		t.Fatalf("dashCase mismatch: %q", got)
	}
}

func TestGenerateIDRemainsMonotonicWhenClockMovesBackwards(t *testing.T) {
	timestamps := []int64{sfEpoch + 1000, sfEpoch + 999, sfEpoch + 1001}
	idx := 0
	gen := &sfGenerator{
		workerID: 1,
		nowMillis: func() int64 {
			value := timestamps[idx]
			if idx < len(timestamps)-1 {
				idx++
			}
			return value
		},
	}

	first := gen.nextID()
	second := gen.nextID()
	third := gen.nextID()

	if !(first < second && second < third) {
		t.Fatalf("IDs should remain monotonic, got %d, %d, %d", first, second, third)
	}
}

func TestFindMenuIDRejectsNilDB(t *testing.T) {
	if _, err := findMenuID(nil, "path", "/demo", menuTypeDirectory); err != sql.ErrConnDone {
		t.Fatalf("findMenuID nil db mismatch: %v", err)
	}
}

func TestGenerateBatchWithEmptyInputReturnsEarly(t *testing.T) {
	gen := New(Config{DSN: "invalid-dsn"})
	count, err := gen.GenerateBatch(nil)
	if err != nil {
		t.Fatalf("GenerateBatch should ignore empty input: %v", err)
	}
	if count != 0 {
		t.Fatalf("GenerateBatch empty count mismatch: %d", count)
	}
}

func TestValidateBatchMetasRejectsDuplicateMenuPath(t *testing.T) {
	metas := []*parser.TableMeta{
		{TableName: "system_user", AppName: "system", ModuleName: "user", ModelName: "User", Comment: "用户"},
		{TableName: "system_user", AppName: "system", ModuleName: "user", ModelName: "User", Comment: "用户"},
	}

	err := validateBatchMetas(metas)
	if err == nil || !strings.Contains(err.Error(), "菜单 path 冲突") {
		t.Fatalf("expected duplicate path error, got %v", err)
	}
}

func TestValidateBatchMetasRejectsNilMeta(t *testing.T) {
	metas := []*parser.TableMeta{
		nil,
	}

	err := validateBatchMetas(metas)
	if err == nil || !strings.Contains(err.Error(), "元数据为空") {
		t.Fatalf("expected nil meta error, got %v", err)
	}
}

func TestValidateBatchMetasRejectsMissingIdentity(t *testing.T) {
	metas := []*parser.TableMeta{
		{TableName: "broken_table", AppName: "", ModuleName: "", Comment: "损坏"},
	}

	err := validateBatchMetas(metas)
	if err == nil || !strings.Contains(err.Error(), "缺少应用名") {
		t.Fatalf("expected missing app error, got %v", err)
	}
}
