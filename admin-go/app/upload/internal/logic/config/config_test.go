package config

import (
	"testing"

	"gbaseadmin/app/upload/internal/model"
)

func TestPickSensitiveValueTreatsBlankAsUnset(t *testing.T) {
	if got := pickSensitiveValue("", "fallback"); got != "fallback" {
		t.Fatalf("empty input mismatch: %q", got)
	}
	if got := pickSensitiveValue("   ", "fallback"); got != "fallback" {
		t.Fatalf("blank input mismatch: %q", got)
	}
	if got := pickSensitiveValue("new-secret", "fallback"); got != "new-secret" {
		t.Fatalf("new input mismatch: %q", got)
	}
}

func TestSanitizeConfigOutputClearsSecrets(t *testing.T) {
	detail := &model.ConfigDetailOutput{
		OssAccessKey: "ak",
		OssSecretKey: "sk",
		CosSecretID:  "sid",
		CosSecretKey: "skey",
	}
	sanitizeConfigOutput(detail)
	if detail.OssAccessKey != "" || detail.OssSecretKey != "" || detail.CosSecretID != "" || detail.CosSecretKey != "" {
		t.Fatalf("detail secrets should be cleared: %+v", detail)
	}

	list := &model.ConfigListOutput{
		OssAccessKey: "ak",
		OssSecretKey: "sk",
		CosSecretID:  "sid",
		CosSecretKey: "skey",
	}
	sanitizeConfigOutput(list)
	if list.OssAccessKey != "" || list.OssSecretKey != "" || list.CosSecretID != "" || list.CosSecretKey != "" {
		t.Fatalf("list secrets should be cleared: %+v", list)
	}
}

func TestValidateConfigFields(t *testing.T) {
	if err := validateConfigFields(1, "resource/upload", "", "", "", "", "", "", "", ""); err != nil {
		t.Fatalf("local storage validation should succeed: %v", err)
	}
	if err := validateConfigFields(2, "", "oss-cn-shanghai.aliyuncs.com", "demo", "ak", "sk", "", "", "", ""); err != nil {
		t.Fatalf("oss validation should succeed: %v", err)
	}
	if err := validateConfigFields(3, "", "", "", "", "", "ap-guangzhou", "demo-123", "sid", "skey"); err != nil {
		t.Fatalf("cos validation should succeed: %v", err)
	}
	if err := validateConfigFields(3, "", "", "", "", "", "", "demo-123", "sid", "skey"); err == nil {
		t.Fatal("incomplete COS config should fail")
	}
}

func TestNormalizeConfigInputs(t *testing.T) {
	createIn := &model.ConfigCreateInput{
		Name:         " 本地存储 ",
		LocalPath:    " resource/upload ",
		OssEndpoint:  " oss-cn-shanghai.aliyuncs.com ",
		OssBucket:    " bucket ",
		OssAccessKey: " ak ",
		OssSecretKey: " sk ",
		CosRegion:    " ap-guangzhou ",
		CosBucket:    " bucket-123 ",
		CosSecretID:  " sid ",
		CosSecretKey: " skey ",
	}
	normalizeConfigCreateInput(createIn)
	if createIn.Name != "本地存储" || createIn.LocalPath != "resource/upload" || createIn.OssBucket != "bucket" || createIn.CosRegion != "ap-guangzhou" {
		t.Fatalf("normalizeConfigCreateInput mismatch: %+v", createIn)
	}

	listIn := &model.ConfigListInput{Keyword: " demo "}
	normalizeConfigListInput(listIn)
	if listIn.Keyword != "demo" {
		t.Fatalf("normalizeConfigListInput mismatch: %+v", listIn)
	}
}
