package appticket

import (
	"context"
	"testing"
	"time"
)

func TestGenerateAndParseRoundTrip(t *testing.T) {
	ticket, err := Generate(context.Background(), "admin", "crm", defaultAppID)
	if err != nil {
		t.Fatalf("Generate should succeed: %v", err)
	}

	claims, err := Parse(context.Background(), ticket)
	if err != nil {
		t.Fatalf("Parse should succeed: %v", err)
	}
	if claims.Username != "admin" {
		t.Fatalf("username mismatch: %q", claims.Username)
	}
	if claims.SourceApp != "crm" {
		t.Fatalf("source app mismatch: %q", claims.SourceApp)
	}
	if claims.TargetApp != defaultAppID {
		t.Fatalf("target app mismatch: %q", claims.TargetApp)
	}
	if err := ValidateTarget(context.Background(), claims); err != nil {
		t.Fatalf("ValidateTarget should accept current app id: %v", err)
	}
}

func TestReplayCacheKeyFallsBackToTicketHash(t *testing.T) {
	key := ReplayCacheKey(nil, " demo-ticket ")
	if key == "system:auth:ticket:used:" {
		t.Fatal("ReplayCacheKey should not return an empty suffix")
	}
}

func TestReplayTTLUsesExpiry(t *testing.T) {
	claims := &Claims{}
	claims.ExpiresAt = nil
	if got := ReplayTTL(claims); got < time.Second {
		t.Fatalf("ReplayTTL should keep a positive ttl, got %v", got)
	}
}
