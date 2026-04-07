package response

import "testing"

func TestResolveMessageUsesFallbackForEmptyOverrides(t *testing.T) {
	if got := resolveMessage("fallback"); got != "fallback" {
		t.Fatalf("resolveMessage fallback mismatch: %q", got)
	}
	if got := resolveMessage("fallback", "   "); got != "fallback" {
		t.Fatalf("resolveMessage should keep fallback for blank override, got %q", got)
	}
	if got := resolveMessage("fallback", "   ", "custom"); got != "custom" {
		t.Fatalf("resolveMessage should pick first non-empty override, got %q", got)
	}
	if got := resolveMessage("fallback", "custom"); got != "custom" {
		t.Fatalf("resolveMessage custom mismatch: %q", got)
	}
}
