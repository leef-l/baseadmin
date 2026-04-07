package alipay

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/url"
	"testing"
)

func TestParsePublicKeySupportsPKCS1(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	})
	parsed, err := parsePublicKey(string(pemBytes))
	if err != nil {
		t.Fatalf("parsePublicKey failed: %v", err)
	}
	if parsed.N.Cmp(key.PublicKey.N) != 0 || parsed.E != key.PublicKey.E {
		t.Fatal("parsed public key mismatch")
	}
}

func TestNormalizePEMWrapsRawContent(t *testing.T) {
	raw := base64.StdEncoding.EncodeToString([]byte("demo-key"))
	normalized := normalizePEM(raw, "PUBLIC KEY")
	if normalized == raw {
		t.Fatal("normalizePEM should wrap raw content")
	}
}

func TestNormalizePEMBlankInput(t *testing.T) {
	if got := normalizePEM("   ", "PUBLIC KEY"); got != "" {
		t.Fatalf("normalizePEM blank input mismatch: %q", got)
	}
}

func TestCreateOrderIncludesReturnURL(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	client := &Client{
		cfg: Config{
			AppID:      "app-id",
			NotifyURL:  "https://notify.example/callback",
			ReturnURL:  " https://return.example/done ",
			GatewayURL: "https://openapi.alipay.com/gateway.do",
		},
		privateKey: key,
	}

	payURL, err := client.CreateOrder(context.Background(), "ORD-1", 1234, "demo")
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	parsed, err := url.Parse(payURL)
	if err != nil {
		t.Fatalf("url.Parse failed: %v", err)
	}
	if got := parsed.Query().Get("return_url"); got != "https://return.example/done" {
		t.Fatalf("return_url mismatch: %q", got)
	}
}

func TestVerifyNotifyRequiresOutTradeNo(t *testing.T) {
	client := &Client{}
	params := url.Values{
		"sign":         {"demo-sign"},
		"trade_status": {"TRADE_SUCCESS"},
	}
	if _, err := client.VerifyNotify(context.Background(), params); err == nil {
		t.Fatal("VerifyNotify should reject missing out_trade_no")
	}
}
