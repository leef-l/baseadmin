package password

import (
	"testing"

	"github.com/gogf/gf/v2/crypto/gsha256"
	"golang.org/x/crypto/bcrypt"
)

func TestHashAndVerifyRoundTrip(t *testing.T) {
	hashed, err := Hash("admin123")
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	if !Verify(hashed, "admin123") {
		t.Fatal("Verify should accept freshly hashed password")
	}
	if Verify(hashed, "wrong") {
		t.Fatal("Verify should reject wrong password")
	}
}

func TestVerifySupportsLegacySHA256(t *testing.T) {
	legacy := gsha256.Encrypt("admin123")
	if !Verify(legacy, "admin123") {
		t.Fatal("Verify should support legacy sha256 digest")
	}
	if !NeedsRehash(legacy) {
		t.Fatal("legacy sha256 digest should require rehash")
	}
}

func TestNeedsRehashForLowCostBCrypt(t *testing.T) {
	lowCostHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword failed: %v", err)
	}
	if !NeedsRehash(string(lowCostHash)) {
		t.Fatal("low-cost bcrypt hash should require rehash")
	}

	currentHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword failed: %v", err)
	}
	if NeedsRehash(string(currentHash)) {
		t.Fatal("default-cost bcrypt hash should not require rehash")
	}
}

func TestVerifyAndNeedsRehashTrimStoredHash(t *testing.T) {
	legacy := "  " + gsha256.Encrypt("admin123") + "  "
	if !Verify(legacy, "admin123") {
		t.Fatal("Verify should trim legacy digests")
	}

	currentHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword failed: %v", err)
	}
	if !Verify(" \n"+string(currentHash)+"\t ", "admin123") {
		t.Fatal("Verify should trim bcrypt hashes")
	}
	if NeedsRehash(" " + string(currentHash) + " ") {
		t.Fatal("NeedsRehash should ignore surrounding whitespace")
	}
}
