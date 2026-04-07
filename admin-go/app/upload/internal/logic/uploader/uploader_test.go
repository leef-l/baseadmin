package uploader

import (
	"encoding/json"
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
