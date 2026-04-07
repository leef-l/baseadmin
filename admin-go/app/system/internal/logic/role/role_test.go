package role

import (
	"testing"

	"gbaseadmin/app/system/internal/model"
)

func TestNormalizeRoleInputs(t *testing.T) {
	createIn := &model.RoleCreateInput{Title: " 管理员 "}
	normalizeRoleCreateInput(createIn)
	if createIn.Title != "管理员" {
		t.Fatalf("normalizeRoleCreateInput mismatch: %+v", createIn)
	}

	updateIn := &model.RoleUpdateInput{Title: " 编辑 "}
	normalizeRoleUpdateInput(updateIn)
	if updateIn.Title != "编辑" {
		t.Fatalf("normalizeRoleUpdateInput mismatch: %+v", updateIn)
	}

	listIn := &model.RoleListInput{Keyword: " admin "}
	normalizeRoleListInput(listIn)
	if listIn.Keyword != "admin" {
		t.Fatalf("normalizeRoleListInput mismatch: %+v", listIn)
	}

	treeIn := &model.RoleTreeInput{Keyword: " operator "}
	normalizeRoleTreeInput(treeIn)
	if treeIn.Keyword != "operator" {
		t.Fatalf("normalizeRoleTreeInput mismatch: %+v", treeIn)
	}
}
