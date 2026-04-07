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

func TestNormalizeConfigTrimsFieldsAndDefaultsProvider(t *testing.T) {
	cfg := normalizeConfig(smsConfig{
		Provider:        "  ",
		AccessKeyId:     "  id  ",
		AccessKeySecret: "  secret  ",
		SignName:        "  sign  ",
		TemplateCode:    "  tpl  ",
	})
	if cfg.Provider != "aliyun" {
		t.Fatalf("provider mismatch: %q", cfg.Provider)
	}
	if cfg.AccessKeyId != "id" || cfg.AccessKeySecret != "secret" || cfg.SignName != "sign" || cfg.TemplateCode != "tpl" {
		t.Fatalf("normalizeConfig should trim fields: %+v", cfg)
	}
}
