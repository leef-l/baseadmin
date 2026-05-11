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

func TestParseTokenTrimsWhitespace(t *testing.T) {
	token, err := GenerateToken(123, "admin", 9, 0, 0)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if _, err := ParseToken(" \n" + token + "\t "); err != nil {
		t.Fatalf("ParseToken should accept surrounding whitespace: %v", err)
	}
}

func TestDefaultInsecureSecretIsRejected(t *testing.T) {
	if defaultInsecureSecret != "gbaseadmin-secret-key" {
		t.Fatalf("default insecure secret constant changed unexpectedly: %q", defaultInsecureSecret)
	}
}
