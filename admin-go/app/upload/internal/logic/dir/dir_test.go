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
