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
	if got := BuildLocalFileURLWithBase("resource/upload", "../cert", "demo.pem"); got != "/resource/cert/demo.pem" {
		t.Fatalf("BuildLocalFileURLWithBase parent mismatch: %q", got)
	}
	if got := BuildLocalFileURLWithBase("resource/upload", "2026-04-08", "demo.png"); got != "/upload/2026-04-08/demo.png" {
		t.Fatalf("BuildLocalFileURLWithBase upload mismatch: %q", got)
	}
}

func TestLocalStoragePhysicalPath(t *testing.T) {
	want := filepath.Join(DefaultLocalStoragePath, "2026-04-08", "demo.png")
	if got := LocalStoragePhysicalPath("/upload/2026-04-08/demo.png"); got != want {
		t.Fatalf("LocalStoragePhysicalPath mismatch: got=%q want=%q", got, want)
	}
	parentWant := filepath.Join("resource", "cert", "demo.pem")
	if got := LocalStoragePhysicalPath("/resource/cert/demo.pem"); got != parentWant {
		t.Fatalf("LocalStoragePhysicalPath parent mismatch: got=%q want=%q", got, parentWant)
	}
}

func TestResolveLocalStorageDir(t *testing.T) {
	if got, ok := ResolveLocalStorageDir("resource/upload", "../cert"); !ok || got != filepath.Join("resource", "cert") {
		t.Fatalf("ResolveLocalStorageDir parent mismatch: got=%q ok=%v", got, ok)
	}
	if _, ok := ResolveLocalStorageDir("resource/upload", "../../outside"); ok {
		t.Fatal("ResolveLocalStorageDir should reject escaping public root")
	}
}
