package dept

import (
	"testing"

	"gbaseadmin/app/system/internal/model"
)

func TestNormalizeDeptInputs(t *testing.T) {
	createIn := &model.DeptCreateInput{
		Title:    " 技术部 ",
		Username: " admin ",
		Email:    " admin@example.com ",
	}
	normalizeDeptCreateInput(createIn)
	if createIn.Title != "技术部" || createIn.Username != "admin" || createIn.Email != "admin@example.com" {
		t.Fatalf("normalizeDeptCreateInput mismatch: %+v", createIn)
	}

	updateIn := &model.DeptUpdateInput{
		Title:    " 运营部 ",
		Username: " ops ",
		Email:    " ops@example.com ",
	}
	normalizeDeptUpdateInput(updateIn)
	if updateIn.Title != "运营部" || updateIn.Username != "ops" || updateIn.Email != "ops@example.com" {
		t.Fatalf("normalizeDeptUpdateInput mismatch: %+v", updateIn)
	}

	listIn := &model.DeptListInput{Keyword: " root "}
	normalizeDeptListInput(listIn)
	if listIn.Keyword != "root" {
		t.Fatalf("normalizeDeptListInput mismatch: %+v", listIn)
	}

	treeIn := &model.DeptTreeInput{Keyword: " center "}
	normalizeDeptTreeInput(treeIn)
	if treeIn.Keyword != "center" {
		t.Fatalf("normalizeDeptTreeInput mismatch: %+v", treeIn)
	}
}

func TestDeptInputValidation(t *testing.T) {
	deptSvc := &sDept{}
	if err := deptSvc.Create(nil, nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Create nil input mismatch: %v", err)
	}
	if err := deptSvc.Create(nil, &model.DeptCreateInput{Title: " "}); err == nil || err.Error() != "部门名称不能为空" {
		t.Fatalf("Create blank title mismatch: %v", err)
	}
	if err := deptSvc.Update(nil, &model.DeptUpdateInput{ID: 1, Title: " "}); err == nil || err.Error() != "部门名称不能为空" {
		t.Fatalf("Update blank title mismatch: %v", err)
	}
}
