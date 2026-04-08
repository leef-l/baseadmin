package shared

import (
	"path/filepath"
	"testing"
)

func TestNormalizeLocalStoragePath(t *testing.T) {
	if got := NormalizeLocalStoragePath("  resource/upload/  "); got != "resource/upload" {
		t.Fatalf("NormalizeLocalStoragePath mismatch: %q", got)
	}
	if got := NormalizeLocalStoragePath(" "); got != DefaultLocalStoragePath {
		t.Fatalf("NormalizeLocalStoragePath blank mismatch: %q", got)
	}
	if got := NormalizeLocalStoragePath(" . "); got != DefaultLocalStoragePath {
		t.Fatalf("NormalizeLocalStoragePath dot mismatch: %q", got)
	}
}

func TestBuildLocalFileURL(t *testing.T) {
	if got := BuildLocalFileURL("2026-04-08", "demo.png"); got != "/upload/2026-04-08/demo.png" {
		t.Fatalf("BuildLocalFileURL mismatch: %q", got)
	}
	if got := BuildLocalFileURL(" 2026-04-08/../ ", `\nested\demo.png `); got != "/upload/2026-04-08/nested/demo.png" {
		t.Fatalf("BuildLocalFileURL sanitized mismatch: %q", got)
	}
}

func TestLocalStoragePhysicalPath(t *testing.T) {
	want := filepath.Join(DefaultLocalStoragePath, "2026-04-08", "demo.png")
	if got := LocalStoragePhysicalPath("/upload/2026-04-08/demo.png"); got != want {
		t.Fatalf("LocalStoragePhysicalPath mismatch: got=%q want=%q", got, want)
	}
	safeWant := filepath.Join(DefaultLocalStoragePath, "nested", "demo.png")
	if got := LocalStoragePhysicalPath("/upload/../../nested/../demo.png"); got != safeWant {
		t.Fatalf("LocalStoragePhysicalPath should stay inside root: got=%q want=%q", got, safeWant)
	}
}
