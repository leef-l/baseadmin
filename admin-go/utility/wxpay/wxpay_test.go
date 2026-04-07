package wxpay

import (
	"strings"
	"testing"
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
