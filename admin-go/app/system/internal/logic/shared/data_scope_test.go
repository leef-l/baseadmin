package shared

import (
	"reflect"
	"testing"
)

func TestSelectDataScopeModes(t *testing.T) {
	rows := []dataScopeRoleRow{
		{RoleID: 1, DataScope: 5},
		{RoleID: 2, DataScope: 3},
		{RoleID: 3, DataScope: 4},
	}
	got := selectDataScopeModes(rows)
	if !got.IncludeCurrentDept || !got.IncludeSelf || !got.IncludeCustomDept {
		t.Fatalf("selectDataScopeModes mismatch: %+v", got)
	}
	if got.AllowAll || got.IncludeDeptAndChildren {
		t.Fatalf("selectDataScopeModes unexpected flags: %+v", got)
	}
	if got := selectDataScopeModes(nil); !got.IncludeSelf {
		t.Fatalf("selectDataScopeModes nil mismatch: %+v", got)
	}
}

func TestSelectDataScopeModesAllowAll(t *testing.T) {
	got := selectDataScopeModes([]dataScopeRoleRow{
		{RoleID: 1, DataScope: 1},
		{RoleID: 2, DataScope: 5},
	})
	if !got.AllowAll {
		t.Fatalf("selectDataScopeModes should allow all: %+v", got)
	}
}

func TestCustomScopeRoleIDs(t *testing.T) {
	rows := []dataScopeRoleRow{
		{RoleID: 2, DataScope: 5},
		{RoleID: 2, DataScope: 5},
		{RoleID: 3, DataScope: 4},
		{RoleID: 4, DataScope: 9},
	}
	want := []int64{2}
	if got := customScopeRoleIDs(rows); !reflect.DeepEqual(got, want) {
		t.Fatalf("customScopeRoleIDs mismatch: got=%v want=%v", got, want)
	}
}

func TestExpandDeptTree(t *testing.T) {
	tree := map[int64][]int64{
		0: {1},
		1: {2, 3},
		2: {4},
	}
	want := []int64{1, 2, 3, 4}
	if got := expandDeptTree([]int64{1}, tree); !reflect.DeepEqual(got, want) {
		t.Fatalf("expandDeptTree mismatch: got=%v want=%v", got, want)
	}
}

func TestContainsInt64(t *testing.T) {
	if !containsInt64([]int64{2, 5, 9}, 5) {
		t.Fatal("containsInt64 should find existing value")
	}
	if containsInt64([]int64{2, 5, 9}, 7) {
		t.Fatal("containsInt64 should reject missing value")
	}
}
