package auth

import (
	"reflect"
	"testing"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

func TestUserCacheKeys(t *testing.T) {
	if got := userCacheKeys(0); got != nil {
		t.Fatalf("userCacheKeys should ignore non-positive ids, got %v", got)
	}

	want := []string{
		"system:auth:info:42",
		"system:auth:menus:42",
	}
	if got := userCacheKeys(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("userCacheKeys mismatch: got=%v want=%v", got, want)
	}
}

func TestLoginFailCacheKeyNormalizesInputs(t *testing.T) {
	if got := loginFailCacheKey(" Admin ", " 127.0.0.1 "); got != "system:auth:login_fail:admin:127.0.0.1" {
		t.Fatalf("loginFailCacheKey mismatch: %q", got)
	}

	if got := loginFailCacheKey(" ", " "); got != "system:auth:login_fail:unknown:unknown" {
		t.Fatalf("loginFailCacheKey should fallback to unknown parts, got %q", got)
	}
}

func TestCompactInt64s(t *testing.T) {
	input := []int64{0, 3, 3, -1, 2, 2, 5}
	want := []int64{3, 2, 5}
	if got := compactInt64s(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("compactInt64s mismatch: got=%v want=%v", got, want)
	}
}

func TestCompactPermissions(t *testing.T) {
	input := []string{" view:user ", "", "view:user", "edit:user", "  "}
	want := []string{"view:user", "edit:user"}
	if got := compactPermissions(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("compactPermissions mismatch: got=%v want=%v", got, want)
	}
}

func TestBuildMenuTree(t *testing.T) {
	root := &model.AuthMenuOutput{ID: snowflake.JsonInt64(1), ParentID: 0, Title: "root"}
	child := &model.AuthMenuOutput{ID: snowflake.JsonInt64(2), ParentID: snowflake.JsonInt64(1), Title: "child"}
	orphan := &model.AuthMenuOutput{ID: snowflake.JsonInt64(3), ParentID: snowflake.JsonInt64(99), Title: "orphan"}

	tree := buildMenuTree([]*model.AuthMenuOutput{root, child, orphan})
	if len(tree) != 2 {
		t.Fatalf("buildMenuTree top-level count mismatch: %d", len(tree))
	}
	if len(root.Children) != 1 || root.Children[0] != child {
		t.Fatalf("child node should attach to root: %+v", root.Children)
	}
	if tree[1] != orphan {
		t.Fatalf("orphan node should remain top-level: %+v", tree[1])
	}
}
