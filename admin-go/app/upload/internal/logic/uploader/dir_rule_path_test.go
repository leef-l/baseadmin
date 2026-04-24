package uploader

import (
	"testing"
	"time"
)

func TestNormalizeUploadRuleFileTypes(t *testing.T) {
	got := normalizeUploadRuleFileTypes(" .TXT,doc；pdf txt image/* application/pdf ")
	if got != "txt,doc,pdf,image/*,application/pdf" {
		t.Fatalf("normalizeUploadRuleFileTypes mismatch: %q", got)
	}
}

func TestDirRuleFileTypeMatches(t *testing.T) {
	if !dirRuleFileTypeMatches("txt,doc,pdf", ".DOC", "application/msword") {
		t.Fatal("expected dirRuleFileTypeMatches to match normalized extension")
	}
	if !dirRuleFileTypeMatches("image,image/*,application/pdf", ".webp", "image/webp") {
		t.Fatal("expected dirRuleFileTypeMatches to match mime alias and prefix")
	}
	if !dirRuleFileTypeMatches("application/pdf", ".bin", "application/pdf; charset=binary") {
		t.Fatal("expected dirRuleFileTypeMatches to match exact mime")
	}
	if dirRuleFileTypeMatches("txt,doc,pdf", "xls", "application/vnd.ms-excel") {
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
	got := renderUploadRulePath("docs/{Y}/{m}/{ext}", now, ".PDF", 0)
	if got != "docs/2026/04/pdf" {
		t.Fatalf("renderUploadRulePath mismatch: %q", got)
	}
}

func TestRenderUploadRulePathSupportsSystemUserID(t *testing.T) {
	now := time.Date(2026, 4, 12, 18, 30, 15, 0, time.Local)
	got := renderUploadRulePath("users/{systemUserId}/{Y-m-d}/{ext}", now, ".PNG", 12345)
	if got != "users/12345/2026-04-12/png" {
		t.Fatalf("renderUploadRulePath system user mismatch: %q", got)
	}
}

func TestRenderUploadRulePathKeepsParentRelative(t *testing.T) {
	now := time.Date(2026, 4, 12, 18, 30, 15, 0, time.Local)
	got := renderUploadRulePath("../cert/{ext}", now, ".PEM", 0)
	if got != "../cert/pem" {
		t.Fatalf("renderUploadRulePath parent mismatch: %q", got)
	}
}

func TestRenderUploadRulePathSupportsUpAlias(t *testing.T) {
	now := time.Date(2026, 4, 12, 18, 30, 15, 0, time.Local)
	got := renderUploadRulePath("@up/cert/{ext}", now, ".PEM", 0)
	if got != "../cert/pem" {
		t.Fatalf("renderUploadRulePath @up mismatch: %q", got)
	}
}

func TestBuildObjectKey(t *testing.T) {
	got := buildObjectKey("docs/2026/04", "file.txt")
	if got != "docs/2026/04/file.txt" {
		t.Fatalf("buildObjectKey mismatch: %q", got)
	}
}

func TestSelectUploadRulePathPrefersTypeRule(t *testing.T) {
	rules := []*uploadDirRuleRecord{
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
	rules := []*uploadDirRuleRecord{
		{Category: 1, StorageTypes: "1", SavePath: "local/default/{Y-m-d}"},
		{Category: 2, FileType: "txt,doc", StorageTypes: "1", SavePath: "local/docs/{ext}"},
		{Category: 3, FileType: "/system/users/*", StorageTypes: "1", SavePath: "routes/users/{Y-m-d}"},
	}

	got := selectUploadRulePath(rules, 1, ".DOC", "/system/users/edit")
	if got != "routes/users/{Y-m-d}" {
		t.Fatalf("selectUploadRulePath source-rule mismatch: %q", got)
	}
}

func TestSelectUploadRuleReturnsMatchedDirID(t *testing.T) {
	rules := []*uploadDirRuleRecord{
		{Id: 1, DirId: 11, Category: 1, StorageTypes: "1", SavePath: "local/default/{Y-m-d}"},
		{Id: 2, DirId: 22, Category: 3, FileType: "/system/users/*", StorageTypes: "1", SavePath: "routes/users/{Y-m-d}"},
	}

	got := selectUploadRule(rules, 1, ".DOC", "application/msword", "/system/users/edit")
	if got == nil {
		t.Fatal("expected selectUploadRule to return matched rule")
	}
	if int64(got.DirId) != 22 {
		t.Fatalf("selectUploadRule dir mismatch: %d", got.DirId)
	}
	if got.SavePath != "routes/users/{Y-m-d}" {
		t.Fatalf("selectUploadRule savePath mismatch: %q", got.SavePath)
	}
}

func TestSelectUploadRulePathFallsBackToDefaultRule(t *testing.T) {
	rules := []*uploadDirRuleRecord{
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
	if !hasParentRelativeDir("@up/cert") {
		t.Fatal("expected hasParentRelativeDir to detect @up path")
	}
	if hasParentRelativeDir("cert/demo") {
		t.Fatal("expected hasParentRelativeDir to ignore normal path")
	}
}
