package uploader

import (
	"testing"
	"time"
)

func TestNormalizeUploadRuleFileTypes(t *testing.T) {
	got := normalizeUploadRuleFileTypes(" .TXT,doc；pdf txt ")
	if got != "txt,doc,pdf" {
		t.Fatalf("normalizeUploadRuleFileTypes mismatch: %q", got)
	}
}

func TestDirRuleFileTypeMatches(t *testing.T) {
	if !dirRuleFileTypeMatches("txt,doc,pdf", ".DOC") {
		t.Fatal("expected dirRuleFileTypeMatches to match normalized extension")
	}
	if dirRuleFileTypeMatches("txt,doc,pdf", "xls") {
		t.Fatal("expected dirRuleFileTypeMatches to reject unmatched extension")
	}
}

func TestRenderUploadRulePath(t *testing.T) {
	now := time.Date(2026, 4, 12, 18, 30, 15, 0, time.Local)
	got := renderUploadRulePath("docs/{Y}/{m}/{ext}", now, ".PDF")
	if got != "docs/2026/04/pdf" {
		t.Fatalf("renderUploadRulePath mismatch: %q", got)
	}
}

func TestBuildObjectKey(t *testing.T) {
	got := buildObjectKey("docs/2026/04", "file.txt")
	if got != "docs/2026/04/file.txt" {
		t.Fatalf("buildObjectKey mismatch: %q", got)
	}
}
