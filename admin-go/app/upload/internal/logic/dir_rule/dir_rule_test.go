package dir_rule

import (
	"errors"
	"testing"

	"gbaseadmin/app/upload/internal/model"
)

func TestNormalizeDirRuleInputs(t *testing.T) {
	createIn := &model.DirRuleCreateInput{
		Category:     2,
		FileType:     " .TXT,doc；pdf ",
		StorageTypes: " 1,2；3 ",
		SavePath:     " {Y}/{m}/{d} ",
	}
	normalizeDirRuleCreateInput(createIn)
	if createIn.SavePath != "{Y}/{m}/{d}" {
		t.Fatalf("normalizeDirRuleCreateInput mismatch: %+v", createIn)
	}
	if createIn.FileType != "txt,doc,pdf" {
		t.Fatalf("normalizeDirRuleCreateInput fileType mismatch: %+v", createIn)
	}
	if createIn.StorageTypes != "1,2,3" {
		t.Fatalf("normalizeDirRuleCreateInput storageTypes mismatch: %+v", createIn)
	}

	interfaceIn := &model.DirRuleCreateInput{
		Category:     3,
		FileType:     " /upload/file , https://example.com/system/users?id=1 , /upload/file/* ",
		StorageTypes: "1,2,3",
		SavePath:     " route/{Y} ",
	}
	normalizeDirRuleCreateInput(interfaceIn)
	if interfaceIn.FileType != "/upload/file,/system/users,/upload/file/*" {
		t.Fatalf("normalizeDirRuleCreateInput source matcher mismatch: %+v", interfaceIn)
	}

	updateIn := &model.DirRuleUpdateInput{
		Category:     1,
		FileType:     "txt",
		StorageTypes: " 1；1,3 ",
		SavePath:     " docs/{Y} ",
	}
	normalizeDirRuleUpdateInput(updateIn)
	if updateIn.SavePath != "docs/{Y}" {
		t.Fatalf("normalizeDirRuleUpdateInput mismatch: %+v", updateIn)
	}
	if updateIn.FileType != "" {
		t.Fatalf("normalizeDirRuleUpdateInput fileType mismatch: %+v", updateIn)
	}
	if updateIn.StorageTypes != "1,3" {
		t.Fatalf("normalizeDirRuleUpdateInput storageTypes mismatch: %+v", updateIn)
	}

	listIn := &model.DirRuleListInput{Keyword: " upload "}
	normalizeDirRuleListInput(listIn)
	if listIn.Keyword != "upload" {
		t.Fatalf("normalizeDirRuleListInput mismatch: %+v", listIn)
	}
}

func TestDirRuleInputValidation(t *testing.T) {
	dirRuleSvc := &sDirRule{}
	if err := dirRuleSvc.Create(nil, nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Create nil input mismatch: %v", err)
	}
	var typedNil *model.DirRuleCreateInput
	if err := dirRuleSvc.Create(nil, typedNil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Create typed nil input mismatch: %v", err)
	}
	if err := dirRuleSvc.Update(nil, nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Update nil input mismatch: %v", err)
	}
	if err := validateDirRuleFields(0, 1, 1, "", "1", ""); err == nil || err.Error() != "目录ID不能为空" {
		t.Fatalf("validateDirRuleFields missing dirID mismatch: %v", err)
	}
	if err := validateDirRuleFields(1, 9, 1, "", "1", ""); err == nil || err.Error() != "类别值不合法" {
		t.Fatalf("validateDirRuleFields invalid category mismatch: %v", err)
	}
	if err := validateDirRuleFields(1, 1, 1, "", "9", ""); err == nil || err.Error() != "适用存储值不合法" {
		t.Fatalf("validateDirRuleFields invalid storageTypes mismatch: %v", err)
	}
	if err := validateDirRuleFields(1, 2, 1, "", "1", ""); err == nil || err.Error() != "文件类型不能为空" {
		t.Fatalf("validateDirRuleFields missing fileType mismatch: %v", err)
	}
	if err := validateDirRuleFields(1, 2, 1, "txt,doc$", "1", ""); err == nil || err.Error() != "文件类型格式不正确" {
		t.Fatalf("validateDirRuleFields invalid fileType mismatch: %v", err)
	}
	if err := validateDirRuleFields(1, 3, 1, "", "1", ""); err == nil || err.Error() != "来源标识不能为空" {
		t.Fatalf("validateDirRuleFields missing source matcher mismatch: %v", err)
	}
	if err := validateDirRuleFields(1, 1, 1, "", "1,2", "../cert"); err == nil || err.Error() != "父级目录规则仅支持本地存储" {
		t.Fatalf("validateDirRuleFields parent path storage mismatch: %v", err)
	}
	if _, err := dirRuleSvc.Detail(nil, 0); err == nil || err.Error() != "目录规则不存在或已删除" {
		t.Fatalf("Detail invalid id mismatch: %v", err)
	}
}

func TestShouldRetryDirRuleListWithoutStorageTypes(t *testing.T) {
	err := errors.New("SELECT COUNT(1) FROM `upload_dir_rule` WHERE `storage_types` LIKE '%demo%': Error 1054 (42S22): Unknown column 'storage_types' in 'where clause'")
	if !shouldRetryDirRuleListWithoutStorageTypes(err, "demo") {
		t.Fatal("expected missing storage_types column error to trigger retry")
	}
	if shouldRetryDirRuleListWithoutStorageTypes(err, "   ") {
		t.Fatal("blank keyword should not trigger retry")
	}
	if shouldRetryDirRuleListWithoutStorageTypes(errors.New("boom"), "demo") {
		t.Fatal("unrelated errors should not trigger retry")
	}
}
