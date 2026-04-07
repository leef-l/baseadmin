package uploader

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBuildUniqueNameWithoutExtension(t *testing.T) {
	now := time.Unix(1712457600, 123456789)
	got := buildUniqueName(now, 7, "")
	if strings.Contains(got, ".") {
		t.Fatalf("unique name without extension should not contain dot: %q", got)
	}
}

func TestBuildUniqueNameWithExtension(t *testing.T) {
	now := time.Unix(1712457600, 123456789)
	got := buildUniqueName(now, 7, "png")
	if !strings.HasSuffix(got, ".png") {
		t.Fatalf("unique name should keep extension suffix, got %q", got)
	}
}

func TestBuildUniqueNameNormalizesExtension(t *testing.T) {
	now := time.Unix(1712457600, 123456789)
	got := buildUniqueName(now, 7, " .PNG ")
	if !strings.HasSuffix(got, ".png") {
		t.Fatalf("unique name should normalize extension, got %q", got)
	}
}

func TestGetInt64AndGetStringSupportCommonScanTypes(t *testing.T) {
	values := map[string]interface{}{
		"max_size":   []byte("20"),
		"storage":    "3",
		"local_path": []byte("resource/custom"),
		"json_num":   json.Number("8"),
	}

	if got := getInt64(values, "max_size"); got != 20 {
		t.Fatalf("getInt64 max_size mismatch: %d", got)
	}
	if got := getInt64(values, "storage"); got != 3 {
		t.Fatalf("getInt64 storage mismatch: %d", got)
	}
	if got := getInt64(values, "json_num"); got != 8 {
		t.Fatalf("getInt64 json_num mismatch: %d", got)
	}
	if got := getString(values, "local_path"); got != "resource/custom" {
		t.Fatalf("getString local_path mismatch: %q", got)
	}
}

func TestGetInt64DropsOverflowUint64(t *testing.T) {
	values := map[string]interface{}{
		"too_big": ^uint64(0),
	}
	if got := getInt64(values, "too_big"); got != 0 {
		t.Fatalf("getInt64 should ignore overflow uint64 values, got %d", got)
	}
}

func TestNormalizeLocalStoragePath(t *testing.T) {
	if got := normalizeLocalStoragePath("  resource/upload/  "); got != "resource/upload" {
		t.Fatalf("normalizeLocalStoragePath mismatch: %q", got)
	}
	if got := normalizeLocalStoragePath(" "); got != defaultLocalStoragePath {
		t.Fatalf("normalizeLocalStoragePath blank mismatch: %q", got)
	}
}

func TestBuildLocalFileURL(t *testing.T) {
	if got := buildLocalFileURL("2026-04-08", "demo.png"); got != "/upload/2026-04-08/demo.png" {
		t.Fatalf("buildLocalFileURL mismatch: %q", got)
	}
}

func TestLocalStoragePhysicalPath(t *testing.T) {
	want := filepath.Join(defaultLocalStoragePath, "2026-04-08", "demo.png")
	if got := localStoragePhysicalPath("/upload/2026-04-08/demo.png"); got != want {
		t.Fatalf("localStoragePhysicalPath mismatch: got=%q want=%q", got, want)
	}
}
