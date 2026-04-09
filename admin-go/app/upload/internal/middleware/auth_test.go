package middleware

import "testing"

func TestResolveUploadPermission(t *testing.T) {
	cases := []struct {
		method string
		path   string
		want   string
	}{
		{method: "GET", path: "/api/upload/dir/tree", want: "upload:dir:list"},
		{method: "POST", path: "/api/upload/uploader/upload", want: "upload:file:create"},
		{method: "DELETE", path: "/api/upload/config/delete", want: "upload:config:delete"},
		{method: "GET", path: "/api/upload/dir_rule/detail", want: "upload:dir_rule:list"},
		{method: "DELETE", path: "/api/upload/file/batch-delete", want: "upload:file:batch-delete"},
		{method: "POST", path: "/api/upload/file/export", want: denyPermission},
		{method: "GET", path: "/api/upload/report/list", want: denyPermission},
	}
	for _, tc := range cases {
		if got := resolveUploadPermission(tc.method, tc.path); got != tc.want {
			t.Fatalf("resolveUploadPermission(%q, %q) = %q, want %q", tc.method, tc.path, got, tc.want)
		}
	}
}
