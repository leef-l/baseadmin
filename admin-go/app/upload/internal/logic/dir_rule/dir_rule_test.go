package dir_rule

import (
	"testing"

	"gbaseadmin/app/upload/internal/model"
)

func TestNormalizeDirRuleInputs(t *testing.T) {
	createIn := &model.DirRuleCreateInput{SavePath: " {Y}/{m}/{d} "}
	normalizeDirRuleCreateInput(createIn)
	if createIn.SavePath != "{Y}/{m}/{d}" {
		t.Fatalf("normalizeDirRuleCreateInput mismatch: %+v", createIn)
	}

	updateIn := &model.DirRuleUpdateInput{SavePath: " docs/{Y} "}
	normalizeDirRuleUpdateInput(updateIn)
	if updateIn.SavePath != "docs/{Y}" {
		t.Fatalf("normalizeDirRuleUpdateInput mismatch: %+v", updateIn)
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
}
