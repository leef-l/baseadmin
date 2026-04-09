package uploader

import (
	"context"
	"encoding/json"
	"errors"
	"os"
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

func TestRandomSuffixRange(t *testing.T) {
	for i := 0; i < 32; i++ {
		got := randomSuffix(10000)
		if got < 0 || got >= 10000 {
			t.Fatalf("randomSuffix out of range: %d", got)
		}
	}
	if got := randomSuffix(1); got != 0 {
		t.Fatalf("randomSuffix(1) mismatch: %d", got)
	}
}

func TestNewStorageProvider(t *testing.T) {
	if _, ok := newStorageProvider(uploadStorageConfig{StorageType: 1}).(localStorageProvider); !ok {
		t.Fatal("storage type 1 should use localStorageProvider")
	}
	if _, ok := newStorageProvider(uploadStorageConfig{StorageType: 2}).(ossStorageProvider); !ok {
		t.Fatal("storage type 2 should use ossStorageProvider")
	}
	if _, ok := newStorageProvider(uploadStorageConfig{StorageType: 3}).(cosStorageProvider); !ok {
		t.Fatal("storage type 3 should use cosStorageProvider")
	}
}

func TestLocalStorageProviderRollbackRemovesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "demo.txt")
	if err := os.WriteFile(path, []byte("demo"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	result, err := localStorageProvider{}.Store(context.Background(), storeRequest{
		DateDir:       "2026-04-10",
		LocalFilePath: path,
		UniqueName:    "demo.txt",
	})
	if err != nil {
		t.Fatalf("local store failed: %v", err)
	}
	if result.FileURL != "/upload/2026-04-10/demo.txt" {
		t.Fatalf("local store url mismatch: %q", result.FileURL)
	}

	runCleanupHook(context.Background(), "rollback", result.OnRollback)
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("rollback should remove local file, got err=%v", err)
	}
}

func TestCheckStorageReadyValidatesConfigShape(t *testing.T) {
	if err := checkStorageReady(context.Background(), uploadStorageConfig{StorageType: 1, LocalPath: t.TempDir()}); err != nil {
		t.Fatalf("local storage should be ready: %v", err)
	}
	if err := checkStorageReady(context.Background(), uploadStorageConfig{StorageType: 2, OSS: ossConfig{Endpoint: "oss-cn-shanghai.aliyuncs.com"}}); err == nil {
		t.Fatal("incomplete oss config should fail")
	}
	if err := checkStorageReady(context.Background(), uploadStorageConfig{StorageType: 3, COS: cosConfig{Region: "ap-guangzhou"}}); err == nil {
		t.Fatal("incomplete cos config should fail")
	}
}
