package jwt

import (
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseTokenRoundTrip(t *testing.T) {
	token, err := GenerateToken(123, "admin", 9, 7, 5)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserID != 123 || claims.Username != "admin" || claims.DeptID != 9 || claims.TenantID != 7 || claims.MerchantID != 5 {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}

func TestParseTokenRejectsUnexpectedSigningMethod(t *testing.T) {
	now := time.Now()
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS512, Claims{
		UserID:   1,
		Username: "admin",
		DeptID:   2,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(now.Add(time.Hour)),
			IssuedAt:  gojwt.NewNumericDate(now),
			Issuer:    "gbaseadmin",
		},
	})
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("SignedString failed: %v", err)
	}

	if _, err := ParseToken(tokenStr); err == nil {
		t.Fatal("ParseToken should reject non-HS256 tokens")
	}
}

func TestVerifyAnyTokenAcceptsMemberToken(t *testing.T) {
	token, err := GenerateMemberToken(1001, "13800000000", 1, 2002, "coach")
	if err != nil {
		t.Fatalf("GenerateMemberToken failed: %v", err)
	}

	if !VerifyAnyToken(token) {
		t.Fatal("VerifyAnyToken should accept valid member token")
	}
}

func TestParseTokenTrimsWhitespace(t *testing.T) {
	token, err := GenerateToken(123, "admin", 9, 0, 0)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if _, err := ParseToken(" \n" + token + "\t "); err != nil {
		t.Fatalf("ParseToken should accept surrounding whitespace: %v", err)
	}
}

func TestNormalizeSecretFallsBackOnBlank(t *testing.T) {
	if got := normalizeSecret("  ", "fallback"); got != "fallback" {
		t.Fatalf("normalizeSecret fallback mismatch: %q", got)
	}
	if got := normalizeSecret("  real-secret  ", "fallback"); got != "real-secret" {
		t.Fatalf("normalizeSecret trim mismatch: %q", got)
	}
}
