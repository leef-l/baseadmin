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
