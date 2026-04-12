package dir_rule

import (
	"testing"

	"gbaseadmin/app/upload/internal/model"
)

func TestNormalizeDirRuleInputs(t *testing.T) {
	createIn := &model.DirRuleCreateInput{
		Category: 2,
		FileType: " .TXT,doc；pdf ",
		SavePath: " {Y}/{m}/{d} ",
	}
	normalizeDirRuleCreateInput(createIn)
	if createIn.SavePath != "{Y}/{m}/{d}" {
		t.Fatalf("normalizeDirRuleCreateInput mismatch: %+v", createIn)
	}
	if createIn.FileType != "txt,doc,pdf" {
		t.Fatalf("normalizeDirRuleCreateInput fileType mismatch: %+v", createIn)
	}

	updateIn := &model.DirRuleUpdateInput{
		Category: 1,
		FileType: "txt",
		SavePath: " docs/{Y} ",
	}
	normalizeDirRuleUpdateInput(updateIn)
	if updateIn.SavePath != "docs/{Y}" {
		t.Fatalf("normalizeDirRuleUpdateInput mismatch: %+v", updateIn)
	}
	if updateIn.FileType != "" {
		t.Fatalf("normalizeDirRuleUpdateInput fileType mismatch: %+v", updateIn)
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
	if err := validateDirRuleFields(0, 1, 1, ""); err == nil || err.Error() != "目录ID不能为空" {
		t.Fatalf("validateDirRuleFields missing dirID mismatch: %v", err)
	}
	if err := validateDirRuleFields(1, 9, 1, ""); err == nil || err.Error() != "类别值不合法" {
		t.Fatalf("validateDirRuleFields invalid category mismatch: %v", err)
	}
	if err := validateDirRuleFields(1, 2, 1, ""); err == nil || err.Error() != "文件类型不能为空" {
		t.Fatalf("validateDirRuleFields missing fileType mismatch: %v", err)
	}
	if err := validateDirRuleFields(1, 2, 1, "txt,doc$"); err == nil || err.Error() != "文件类型格式不正确" {
		t.Fatalf("validateDirRuleFields invalid fileType mismatch: %v", err)
	}
	if _, err := dirRuleSvc.Detail(nil, 0); err == nil || err.Error() != "目录规则不存在或已删除" {
		t.Fatalf("Detail invalid id mismatch: %v", err)
	}
}
