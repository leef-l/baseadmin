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
	if treeIn.AssignableOnly {
		t.Fatalf("normalizeRoleTreeInput should preserve assignableOnly: %+v", treeIn)
	}
}

func TestRoleInputValidation(t *testing.T) {
	roleSvc := &sRole{}
	if err := roleSvc.Create(nil, nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Create nil input mismatch: %v", err)
	}
	if err := roleSvc.Create(nil, &model.RoleCreateInput{Title: " "}); err == nil || err.Error() != "角色名称不能为空" {
		t.Fatalf("Create blank title mismatch: %v", err)
	}
	if err := roleSvc.Update(nil, &model.RoleUpdateInput{ID: 1, Title: " "}); err == nil || err.Error() != "角色名称不能为空" {
		t.Fatalf("Update blank title mismatch: %v", err)
	}
	if err := validateRoleFields("管理员", 0, 1, 0, 0); err == nil || err.Error() != "数据范围值不合法" {
		t.Fatalf("validateRoleFields invalid dataScope mismatch: %v", err)
	}
	if err := validateRoleFields("管理员", 1, 1, 2, 0); err == nil || err.Error() != "超级管理员值不合法" {
		t.Fatalf("validateRoleFields invalid isAdmin mismatch: %v", err)
	}
	if _, err := roleSvc.Detail(nil, 0); err == nil || err.Error() != "角色不存在或已删除" {
		t.Fatalf("Detail invalid id mismatch: %v", err)
	}
}

func TestIsBuiltinAdminRoleIgnoresInvalidID(t *testing.T) {
	roleSvc := &sRole{}
	got, err := roleSvc.isBuiltinAdminRole(nil, 0)
	if err != nil {
		t.Fatalf("isBuiltinAdminRole returned err: %v", err)
	}
	if got {
		t.Fatal("isBuiltinAdminRole should ignore invalid id")
	}
}

func TestGrantAndGetRoleRelationsRejectInvalidRole(t *testing.T) {
	roleSvc := &sRole{}

	if err := roleSvc.GrantMenu(nil, &model.RoleGrantMenuInput{ID: 0}); err == nil || err.Error() != "角色不存在或已删除" {
		t.Fatalf("GrantMenu invalid role mismatch: %v", err)
	}

	if err := roleSvc.GrantDept(nil, &model.RoleGrantDeptInput{ID: 0, DataScope: 1}); err == nil || err.Error() != "角色不存在或已删除" {
		t.Fatalf("GrantDept invalid role mismatch: %v", err)
	}

	if _, err := roleSvc.GetMenuIDs(nil, 0); err == nil || err.Error() != "角色不存在或已删除" {
		t.Fatalf("GetMenuIDs invalid role mismatch: %v", err)
	}

	if _, err := roleSvc.GetDeptIDs(nil, 0); err == nil || err.Error() != "角色不存在或已删除" {
		t.Fatalf("GetDeptIDs invalid role mismatch: %v", err)
	}
}

func TestAdminRoleMutationValidation(t *testing.T) {
	if err := validateAdminRoleMutationAllowed(false, false); err != nil {
		t.Fatalf("validateAdminRoleMutationAllowed should ignore non-admin role: %v", err)
	}
	if err := validateAdminRoleMutationAllowed(true, true); err != nil {
		t.Fatalf("validateAdminRoleMutationAllowed should allow admin actor: %v", err)
	}
	if err := validateAdminRoleMutationAllowed(true, false); err == nil || err.Error() != "只有超级管理员可以操作超级管理员角色" {
		t.Fatalf("validateAdminRoleMutationAllowed mismatch: %v", err)
	}
}
