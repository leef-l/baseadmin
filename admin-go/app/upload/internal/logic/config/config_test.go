package config

import (
	"testing"

	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/model/entity"
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
	if err := validateConfigMeta(0, 0, 10, 1); err == nil || err.Error() != "存储类型值不合法" {
		t.Fatalf("validateConfigMeta invalid storage mismatch: %v", err)
	}
	if err := validateConfigMeta(1, 3, 10, 1); err == nil || err.Error() != "是否默认值不合法" {
		t.Fatalf("validateConfigMeta invalid isDefault mismatch: %v", err)
	}
	if err := validateConfigMeta(1, 0, 0, 1); err == nil || err.Error() != "最大文件大小必须大于0" {
		t.Fatalf("validateConfigMeta invalid maxSize mismatch: %v", err)
	}
}

func TestNormalizeConfigInputs(t *testing.T) {
	createIn := &model.ConfigCreateInput{
		Name:         " 本地存储 ",
		LocalPath:    " ./resource/upload//nested/ ",
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
	if createIn.Name != "本地存储" || createIn.LocalPath != "resource/upload/nested" || createIn.OssBucket != "bucket" || createIn.CosRegion != "ap-guangzhou" {
		t.Fatalf("normalizeConfigCreateInput mismatch: %+v", createIn)
	}

	updateIn := &model.ConfigUpdateInput{
		Name:      " 云存储 ",
		LocalPath: "   ",
	}
	normalizeConfigUpdateInput(updateIn)
	if updateIn.Name != "云存储" || updateIn.LocalPath != "" {
		t.Fatalf("normalizeConfigUpdateInput mismatch: %+v", updateIn)
	}

	listIn := &model.ConfigListInput{Keyword: " demo "}
	normalizeConfigListInput(listIn)
	if listIn.Keyword != "demo" {
		t.Fatalf("normalizeConfigListInput mismatch: %+v", listIn)
	}
}

func TestConfigInputValidation(t *testing.T) {
	configSvc := &sConfig{}
	if err := configSvc.Create(nil, nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Create nil input mismatch: %v", err)
	}
	if err := configSvc.Create(nil, &model.ConfigCreateInput{Name: " "}); err == nil || err.Error() != "配置名称不能为空" {
		t.Fatalf("Create blank name mismatch: %v", err)
	}
	if err := configSvc.Update(nil, &model.ConfigUpdateInput{ID: 1, Name: " "}); err == nil || err.Error() != "配置名称不能为空" {
		t.Fatalf("Update blank name mismatch: %v", err)
	}
	if _, err := configSvc.Detail(nil, 0); err == nil || err.Error() != "上传配置不存在或已删除" {
		t.Fatalf("Detail invalid id mismatch: %v", err)
	}
}

func TestConfigMatchesFileURL(t *testing.T) {
	ossConfig := &entity.UploadConfig{
		Storage:     2,
		OssBucket:   "demo-bucket",
		OssEndpoint: "oss-cn-shanghai.aliyuncs.com",
	}
	if !configMatchesFileURL(ossConfig, "https://demo-bucket.oss-cn-shanghai.aliyuncs.com/archive/demo.png") {
		t.Fatal("expected oss file url to match config")
	}

	cosConfig := &entity.UploadConfig{
		Storage:   3,
		CosBucket: "demo-1250000000",
		CosRegion: "ap-guangzhou",
	}
	if !configMatchesFileURL(cosConfig, "https://demo-1250000000.cos.ap-guangzhou.myqcloud.com/archive/demo.png") {
		t.Fatal("expected cos file url to match config")
	}

	localConfig := &entity.UploadConfig{Storage: 1}
	if configMatchesFileURL(localConfig, "/upload/demo.png") {
		t.Fatal("local config should not rely on file url matching")
	}
}
