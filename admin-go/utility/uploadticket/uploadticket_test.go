package uploadticket

import (
	"testing"
	"time"
)

func TestSignAndVerify(t *testing.T) {
	claims := &Claims{
		Scene:               "avatar",
		Dir:                 "play/avatar/1",
		MaxSize:             1024,
		AllowedExts:         []string{"jpg", "png"},
		AllowedMimePrefixes: []string{"image/"},
		ExpiresAt:           time.Now().Add(time.Minute).Unix(),
		Nonce:               "nonce",
	}

	token, err := Sign(claims, "secret")
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	got, err := Verify(token, "secret")
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if got.Scene != claims.Scene || got.Dir != claims.Dir || got.MaxSize != claims.MaxSize {
		t.Fatalf("claims mismatch: %+v", got)
	}
}

func TestVerifyRejectsInvalidSignature(t *testing.T) {
	token, err := Sign(&Claims{
		Scene:     "avatar",
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
	}, "secret")
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	if _, err = Verify(token, "other-secret"); err == nil {
		t.Fatal("verify should reject invalid signature")
	}
}

func TestVerifyRejectsExpiredTicket(t *testing.T) {
	token, err := Sign(&Claims{
		Scene:     "avatar",
		ExpiresAt: time.Now().Add(-time.Minute).Unix(),
	}, "secret")
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	if _, err = Verify(token, "secret"); err == nil {
		t.Fatal("verify should reject expired ticket")
	}
}
