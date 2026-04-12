package shared

import "testing"

func TestNormalizeUploadRuleSource(t *testing.T) {
	cases := map[string]string{
		" /upload/file ":                        "/upload/file",
		"#/upload/file?tab=all":                 "/upload/file",
		"https://example.com/system/users?id=1": "/system/users",
		"system/users/edit":                     "/system/users/edit",
		`\\system\\users\\edit\\`:               "/system/users/edit",
		"":                                      "",
		"   ":                                   "",
		"/":                                     "/",
	}

	for input, want := range cases {
		if got := NormalizeUploadRuleSource(input); got != want {
			t.Fatalf("NormalizeUploadRuleSource(%q) mismatch: got=%q want=%q", input, got, want)
		}
	}
}
