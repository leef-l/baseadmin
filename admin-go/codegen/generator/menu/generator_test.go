package menu

import (
	"database/sql"
	"reflect"
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
