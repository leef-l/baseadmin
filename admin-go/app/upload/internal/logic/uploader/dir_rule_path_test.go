package uploader

import (
	"testing"
	"time"

	"gbaseadmin/app/upload/internal/model/entity"
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

func TestDirRuleSourceMatches(t *testing.T) {
	if !dirRuleSourceMatches("/upload/file,/system/users/*", "https://example.com/system/users/edit?id=1") {
		t.Fatal("expected dirRuleSourceMatches to match normalized route")
	}
	if dirRuleSourceMatches("/upload/image", "/system/users") {
		t.Fatal("expected dirRuleSourceMatches to reject unmatched source")
	}
}

func TestRenderUploadRulePath(t *testing.T) {
	now := time.Date(2026, 4, 12, 18, 30, 15, 0, time.Local)
	got := renderUploadRulePath("docs/{Y}/{m}/{ext}", now, ".PDF")
	if got != "docs/2026/04/pdf" {
		t.Fatalf("renderUploadRulePath mismatch: %q", got)
	}
}

func TestRenderUploadRulePathKeepsParentRelative(t *testing.T) {
	now := time.Date(2026, 4, 12, 18, 30, 15, 0, time.Local)
	got := renderUploadRulePath("../cert/{ext}", now, ".PEM")
	if got != "../cert/pem" {
		t.Fatalf("renderUploadRulePath parent mismatch: %q", got)
	}
}

func TestBuildObjectKey(t *testing.T) {
	got := buildObjectKey("docs/2026/04", "file.txt")
	if got != "docs/2026/04/file.txt" {
		t.Fatalf("buildObjectKey mismatch: %q", got)
	}
}

func TestSelectUploadRulePathPrefersTypeRule(t *testing.T) {
	rules := []*entity.UploadDirRule{
		{Category: 1, StorageTypes: "2,3", SavePath: "cloud/default/{Y-m-d}"},
		{Category: 1, StorageTypes: "1", SavePath: "local/default/{Y-m-d}"},
		{Category: 2, FileType: "txt,doc", StorageTypes: "2", SavePath: "cloud/docs/{ext}"},
		{Category: 2, FileType: "txt,doc", StorageTypes: "1", SavePath: "local/docs/{ext}"},
		{Category: 2, FileType: "jpg,png", StorageTypes: "1", SavePath: "images/{ext}"},
	}

	got := selectUploadRulePath(rules, 1, ".DOC", "")
	if got != "local/docs/{ext}" {
		t.Fatalf("selectUploadRulePath type-rule mismatch: %q", got)
	}
}

func TestSelectUploadRulePathPrefersSourceRule(t *testing.T) {
	rules := []*entity.UploadDirRule{
		{Category: 1, StorageTypes: "1", SavePath: "local/default/{Y-m-d}"},
		{Category: 2, FileType: "txt,doc", StorageTypes: "1", SavePath: "local/docs/{ext}"},
		{Category: 3, FileType: "/system/users/*", StorageTypes: "1", SavePath: "routes/users/{Y-m-d}"},
	}

	got := selectUploadRulePath(rules, 1, ".DOC", "/system/users/edit")
	if got != "routes/users/{Y-m-d}" {
		t.Fatalf("selectUploadRulePath source-rule mismatch: %q", got)
	}
}

func TestSelectUploadRulePathFallsBackToDefaultRule(t *testing.T) {
	rules := []*entity.UploadDirRule{
		{Category: 1, StorageTypes: "2,3", SavePath: "cloud/default/{Y-m-d}"},
		{Category: 1, StorageTypes: "1", SavePath: "local/default/{Y-m-d}"},
		{Category: 2, FileType: "jpg,png", StorageTypes: "1", SavePath: "images/{ext}"},
	}

	got := selectUploadRulePath(rules, 2, "pdf", "")
	if got != "cloud/default/{Y-m-d}" {
		t.Fatalf("selectUploadRulePath default fallback mismatch: %q", got)
	}
}

func TestDirRuleSupportsStorageType(t *testing.T) {
	if !dirRuleSupportsStorageType("", 1) {
		t.Fatal("blank storage types should match legacy rules")
	}
	if !dirRuleSupportsStorageType("1, 3", 3) {
		t.Fatal("expected storage type matcher to find normalized value")
	}
	if dirRuleSupportsStorageType("1,3", 2) {
		t.Fatal("expected storage type matcher to reject unmatched storage")
	}
}

func TestHasParentRelativeDir(t *testing.T) {
	if !hasParentRelativeDir("../cert") {
		t.Fatal("expected hasParentRelativeDir to detect parent path")
	}
	if hasParentRelativeDir("cert/demo") {
		t.Fatal("expected hasParentRelativeDir to ignore normal path")
	}
}
