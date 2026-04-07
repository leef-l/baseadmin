package wxpay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRequestURI(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "full url with query", raw: "https://api.mch.weixin.qq.com/v3/pay/transactions/h5?foo=bar", want: "/v3/pay/transactions/h5?foo=bar"},
		{name: "full url without path", raw: "https://api.mch.weixin.qq.com", want: "/"},
		{name: "path only", raw: "/v3/pay/transactions/h5", want: "/v3/pay/transactions/h5"},
	}

	for _, tc := range tests {
		if got := requestURI(tc.raw); got != tc.want {
			t.Fatalf("%s mismatch: got=%q want=%q", tc.name, got, tc.want)
		}
	}
}

func TestRandomString(t *testing.T) {
	value := randomString(32)
	if len(value) != 32 {
		t.Fatalf("randomString length mismatch: %d", len(value))
	}
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, ch := range value {
		if !strings.ContainsRune(alphabet, ch) {
			t.Fatalf("randomString returned unexpected rune: %q in %q", ch, value)
		}
	}
	if got := randomString(0); got != "" {
		t.Fatalf("randomString(0) should return empty string, got %q", got)
	}
}

func TestVerifySignatureFallbackWithoutPlatformKey(t *testing.T) {
	client := &Client{}
	if err := client.verifySignature([]byte("hello"), base64.StdEncoding.EncodeToString([]byte("sig"))); err != nil {
		t.Fatalf("verifySignature fallback should accept valid base64: %v", err)
	}
}

func TestVerifySignatureWithPlatformKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	message := []byte("timestamp\nnonce\nbody\n")
	sum := sha256.Sum256(message)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, sum[:])
	if err != nil {
		t.Fatalf("SignPKCS1v15 failed: %v", err)
	}

	client := &Client{platformPublicKey: &privateKey.PublicKey}
	if err := client.verifySignature(message, base64.StdEncoding.EncodeToString(sig)); err != nil {
		t.Fatalf("verifySignature with platform key failed: %v", err)
	}
}

func TestMerchantSummary(t *testing.T) {
	client := &Client{cfg: Config{MchID: " mch123 "}}
	if got := client.merchantSummary(); got != "mch123:false" {
		t.Fatalf("merchantSummary mismatch without verifier: %q", got)
	}
	client.platformPublicKey = &rsa.PublicKey{}
	if got := client.merchantSummary(); got != "mch123:true" {
		t.Fatalf("merchantSummary mismatch with verifier: %q", got)
	}
}

func TestRequestClientDefaultsToConfiguredTimeoutClient(t *testing.T) {
	client := &Client{}
	if got := client.requestClient(); got != defaultHTTPClient {
		t.Fatal("requestClient should fall back to package default client")
	}
	if got := defaultHTTPClient.Timeout; got != 15*time.Second {
		t.Fatalf("defaultHTTPClient timeout mismatch: %v", got)
	}
}

func TestRequestClientPrefersInjectedClient(t *testing.T) {
	custom := &http.Client{Timeout: time.Second}
	client := &Client{httpClient: custom}
	if got := client.requestClient(); got != custom {
		t.Fatal("requestClient should prefer injected client")
	}
}
