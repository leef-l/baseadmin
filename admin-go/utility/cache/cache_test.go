package cache

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestNormalizeKeysDropsEmptyAndDuplicates(t *testing.T) {
	got := normalizeKeys([]string{" user:1 ", "", "user:2", "user:1", "  "})
	want := []string{"user:1", "user:2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeKeys mismatch: got=%v want=%v", got, want)
	}
}

func TestNormalizeKey(t *testing.T) {
	if got := normalizeKey(" user:1 "); got != "user:1" {
		t.Fatalf("normalizeKey should trim spaces, got %q", got)
	}
	if got := normalizeKey("   "); got != "" {
		t.Fatalf("normalizeKey should collapse blank input to empty, got %q", got)
	}
}

func TestTTLSecondsRoundsUp(t *testing.T) {
	if got := ttlSeconds(1500 * time.Millisecond); got != 2 {
		t.Fatalf("ttlSeconds should round up fractional seconds, got %d", got)
	}
	if got := ttlSeconds(0); got != 1 {
		t.Fatalf("ttlSeconds should keep minimum ttl of 1 second, got %d", got)
	}
}

func TestShouldDeleteInvalidJSON(t *testing.T) {
	if shouldDeleteInvalidJSON(&json.SyntaxError{}) != true {
		t.Fatal("syntax errors should trigger cache eviction")
	}
	if shouldDeleteInvalidJSON(&json.InvalidUnmarshalError{}) {
		t.Fatal("invalid unmarshal errors should not trigger cache eviction")
	}
}

func TestParseCachedInt64(t *testing.T) {
	if got, err := parseCachedInt64(" 42 "); err != nil || got != 42 {
		t.Fatalf("parseCachedInt64 mismatch: got=%d err=%v", got, err)
	}
	if _, err := parseCachedInt64("not-a-number"); err == nil {
		t.Fatal("parseCachedInt64 should reject invalid values")
	}
}
