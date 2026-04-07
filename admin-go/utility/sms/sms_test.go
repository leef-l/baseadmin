package sms

import (
	"strings"
	"testing"
)

func TestGenerateCodeReturnsSixDigits(t *testing.T) {
	code := generateCode()
	if len(code) != 6 {
		t.Fatalf("generateCode length mismatch: %q", code)
	}
	if strings.Trim(code, "0123456789") != "" {
		t.Fatalf("generateCode should only contain digits: %q", code)
	}
}
