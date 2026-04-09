package middleware

import "testing"

func TestResolveSystemPermission(t *testing.T) {
	cases := []struct {
		method string
		path   string
		want   string
	}{
		{method: "GET", path: "/api/system/auth/info", want: ""},
		{method: "GET", path: "/api/system/dept/tree", want: "system:dept:list"},
		{method: "DELETE", path: "/api/system/menu/delete", want: "system:menu:delete"},
		{method: "POST", path: "/api/system/role/grant-menu", want: "system:role:grant:menu"},
		{method: "GET", path: "/api/system/role/dept-ids", want: "system:role:grant:dept"},
		{method: "PUT", path: "/api/system/users/reset-password", want: "system:user:update"},
		{method: "POST", path: "/api/system/dept/export-now", want: denyPermission},
		{method: "GET", path: "/api/system/auth/refresh", want: denyPermission},
		{method: "GET", path: "/api/system/report/list", want: denyPermission},
	}
	for _, tc := range cases {
		if got := resolveSystemPermission(tc.method, tc.path); got != tc.want {
			t.Fatalf("resolveSystemPermission(%q, %q) = %q, want %q", tc.method, tc.path, got, tc.want)
		}
	}
}

func TestSplitRouteAction(t *testing.T) {
	module, action := splitRouteAction("/api/system/", "/api/system/users/list")
	if module != "users" || action != "list" {
		t.Fatalf("splitRouteAction mismatch: module=%q action=%q", module, action)
	}
}
