package dir

import (
	"testing"

	"gbaseadmin/app/upload/internal/model"
)

func TestNormalizeDirInputs(t *testing.T) {
	createIn := &model.DirCreateInput{
		Name: " 图片目录 ",
		Path: " /upload/image ",
	}
	normalizeDirCreateInput(createIn)
	if createIn.Name != "图片目录" || createIn.Path != "/upload/image" {
		t.Fatalf("normalizeDirCreateInput mismatch: %+v", createIn)
	}

	updateIn := &model.DirUpdateInput{
		Name: " 文档目录 ",
		Path: " /upload/docs ",
	}
	normalizeDirUpdateInput(updateIn)
	if updateIn.Name != "文档目录" || updateIn.Path != "/upload/docs" {
		t.Fatalf("normalizeDirUpdateInput mismatch: %+v", updateIn)
	}

	listIn := &model.DirListInput{Keyword: " media "}
	normalizeDirListInput(listIn)
	if listIn.Keyword != "media" {
		t.Fatalf("normalizeDirListInput mismatch: %+v", listIn)
	}

	treeIn := &model.DirTreeInput{Keyword: " assets "}
	normalizeDirTreeInput(treeIn)
	if treeIn.Keyword != "assets" {
		t.Fatalf("normalizeDirTreeInput mismatch: %+v", treeIn)
	}
}

func TestValidateDirFields(t *testing.T) {
	if err := validateDirFields("", "/upload/demo", 0, 1); err == nil || err.Error() != "目录名称不能为空" {
		t.Fatalf("validateDirFields blank name mismatch: %v", err)
	}
	if err := validateDirFields("目录", "", 0, 1); err == nil || err.Error() != "目录路径不能为空" {
		t.Fatalf("validateDirFields blank path mismatch: %v", err)
	}
	if err := validateDirFields("目录", "/upload/demo", 0, 1); err != nil {
		t.Fatalf("validateDirFields should succeed: %v", err)
	}
	if err := validateDirFields("目录", "/upload/demo", -1, 1); err == nil || err.Error() != "排序不能小于0" {
		t.Fatalf("validateDirFields negative sort mismatch: %v", err)
	}
	if err := validateDirFields("目录", "/upload/demo", 0, 2); err == nil || err.Error() != "状态值不合法" {
		t.Fatalf("validateDirFields invalid status mismatch: %v", err)
	}
	if _, err := (&sDir{}).Detail(nil, 0); err == nil || err.Error() != "目录不存在或已删除" {
		t.Fatalf("Detail invalid id mismatch: %v", err)
	}
}

func TestDirInputValidation(t *testing.T) {
	dirSvc := &sDir{}
	if err := dirSvc.Create(nil, nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Create nil input mismatch: %v", err)
	}
	if err := dirSvc.Create(nil, &model.DirCreateInput{Name: " ", Path: "/upload/demo"}); err == nil || err.Error() != "目录名称不能为空" {
		t.Fatalf("Create blank name mismatch: %v", err)
	}
	if err := dirSvc.Update(nil, &model.DirUpdateInput{ID: 1, Name: "目录", Path: " "}); err == nil || err.Error() != "目录路径不能为空" {
		t.Fatalf("Update blank path mismatch: %v", err)
	}
}
