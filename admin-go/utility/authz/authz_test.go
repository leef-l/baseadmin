package authz

import (
	"reflect"
	"testing"
)

func TestPermissionCacheKey(t *testing.T) {
	if got := PermissionCacheKey(0); got != "" {
		t.Fatalf("PermissionCacheKey(0) = %q, want empty", got)
	}
	if got := PermissionCacheKey(42); got != "system:authz:perms:42" {
		t.Fatalf("PermissionCacheKey(42) = %q", got)
	}
}

func TestCollectRoleIDsDeduplicatesAndSkipsInvalid(t *testing.T) {
	got := collectRoleIDs([]roleRow{
		{RoleID: 2},
		{RoleID: 2},
		{RoleID: 0},
		{RoleID: 3},
	})
	want := []int64{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collectRoleIDs mismatch: got=%v want=%v", got, want)
	}
}

func TestCompactPermissionsDeduplicatesAndTrims(t *testing.T) {
	got := compactPermissions([]permissionRow{
		{Permission: " system:user:list "},
		{Permission: "system:user:list"},
		{Permission: ""},
		{Permission: "upload:file:create"},
	})
	want := []string{"system:user:list", "upload:file:create"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("compactPermissions mismatch: got=%v want=%v", got, want)
	}
}

func TestHasAdminRole(t *testing.T) {
	if !hasAdminRole([]roleRow{{RoleID: 1, IsAdmin: 1}}) {
		t.Fatal("hasAdminRole should detect admin role")
	}
	if hasAdminRole([]roleRow{{RoleID: 1, IsAdmin: 0}}) {
		t.Fatal("hasAdminRole should ignore non-admin role")
	}
}
