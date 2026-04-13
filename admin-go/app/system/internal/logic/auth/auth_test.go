package auth

import (
	"reflect"
	"testing"
	"time"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

func TestUserCacheKeys(t *testing.T) {
	if got := userCacheKeys(0); got != nil {
		t.Fatalf("userCacheKeys should ignore non-positive ids, got %v", got)
	}

	want := []string{
		"system:auth:info:42",
		"system:auth:menus:42",
		"system:authz:perms:42",
	}
	if got := userCacheKeys(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("userCacheKeys mismatch: got=%v want=%v", got, want)
	}
}

func TestLoginFailCacheKeyNormalizesInputs(t *testing.T) {
	if got := loginFailCacheKey(" Admin ", " 127.0.0.1 "); got != "system:auth:login_fail:admin:127.0.0.1" {
		t.Fatalf("loginFailCacheKey mismatch: %q", got)
	}

	if got := loginFailCacheKey(" ", " "); got != "system:auth:login_fail:unknown:unknown" {
		t.Fatalf("loginFailCacheKey should fallback to unknown parts, got %q", got)
	}
}

func TestCompactInt64s(t *testing.T) {
	input := []int64{0, 3, 3, -1, 2, 2, 5}
	want := []int64{3, 2, 5}
	if got := compactInt64s(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("compactInt64s mismatch: got=%v want=%v", got, want)
	}
}

func TestCompactPermissions(t *testing.T) {
	input := []string{" view:user ", "", "view:user", "edit:user", "  "}
	want := []string{"view:user", "edit:user"}
	if got := compactPermissions(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("compactPermissions mismatch: got=%v want=%v", got, want)
	}
}

func TestBuildMenuTree(t *testing.T) {
	root := &model.AuthMenuOutput{ID: snowflake.JsonInt64(1), ParentID: 0, Title: "root"}
	child := &model.AuthMenuOutput{ID: snowflake.JsonInt64(2), ParentID: snowflake.JsonInt64(1), Title: "child"}
	orphan := &model.AuthMenuOutput{ID: snowflake.JsonInt64(3), ParentID: snowflake.JsonInt64(99), Title: "orphan"}

	tree := buildMenuTree([]*model.AuthMenuOutput{root, child, orphan})
	if len(tree) != 2 {
		t.Fatalf("buildMenuTree top-level count mismatch: %d", len(tree))
	}
	if len(root.Children) != 1 || root.Children[0] != child {
		t.Fatalf("child node should attach to root: %+v", root.Children)
	}
	if tree[1] != orphan {
		t.Fatalf("orphan node should remain top-level: %+v", tree[1])
	}
}

func TestNormalizeAuthLoginInput(t *testing.T) {
	in := &model.AuthLoginInput{
		Username: " admin ",
		Password: " 123456 ",
	}
	normalizeAuthLoginInput(in)
	if in.Username != "admin" {
		t.Fatalf("normalizeAuthLoginInput username mismatch: %+v", in)
	}
	if in.Password != " 123456 " {
		t.Fatalf("normalizeAuthLoginInput should not trim password: %+v", in)
	}
}

func TestIssueTicketInputValidation(t *testing.T) {
	authSvc := &sAuth{}
	if _, err := authSvc.IssueTicket(nil, nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("IssueTicket nil input mismatch: %v", err)
	}
	if _, err := authSvc.IssueTicket(nil, &model.AuthIssueTicketInput{
		UserID:    0,
		TargetApp: "crm",
	}); err == nil || err.Error() != "用户ID不能为空" {
		t.Fatalf("IssueTicket user id mismatch: %v", err)
	}
	if _, err := authSvc.IssueTicket(nil, &model.AuthIssueTicketInput{
		UserID:    1,
		TargetApp: " ",
	}); err == nil || err.Error() != "目标应用不能为空" {
		t.Fatalf("IssueTicket target app mismatch: %v", err)
	}
}

func TestNormalizeAuthIssueTicketInput(t *testing.T) {
	in := &model.AuthIssueTicketInput{TargetApp: " crm "}
	normalizeAuthIssueTicketInput(in)
	if in.TargetApp != "crm" {
		t.Fatalf("normalizeAuthIssueTicketInput mismatch: %+v", in)
	}
}

func TestMarkTicketUsedLocal(t *testing.T) {
	key := "system:auth:ticket:used:test-local"
	if replayed := markTicketUsedLocal(key, 50*time.Millisecond); replayed {
		t.Fatal("first local ticket use should not be treated as replay")
	}
	if replayed := markTicketUsedLocal(key, 50*time.Millisecond); !replayed {
		t.Fatal("second local ticket use should be treated as replay")
	}
	time.Sleep(110 * time.Millisecond)
	if replayed := markTicketUsedLocal(key, 50*time.Millisecond); replayed {
		t.Fatal("expired local ticket marker should be evicted")
	}
}

func TestChangePasswordInputValidation(t *testing.T) {
	authSvc := &sAuth{}
	if err := authSvc.ChangePassword(nil, nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("ChangePassword nil input mismatch: %v", err)
	}
	var typedNil *model.AuthChangePasswordInput
	if err := authSvc.ChangePassword(nil, typedNil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("ChangePassword typed nil input mismatch: %v", err)
	}
	if err := authSvc.ChangePassword(nil, &model.AuthChangePasswordInput{
		UserID:      1,
		OldPassword: " ",
		NewPassword: "abc12345",
	}); err == nil || err.Error() != "旧密码不能为空" {
		t.Fatalf("ChangePassword blank old password mismatch: %v", err)
	}
	if err := authSvc.ChangePassword(nil, &model.AuthChangePasswordInput{
		UserID:      1,
		OldPassword: "abc123",
		NewPassword: "",
	}); err == nil || err.Error() != "新密码不能为空" {
		t.Fatalf("ChangePassword blank new password mismatch: %v", err)
	}
	if err := authSvc.ChangePassword(nil, &model.AuthChangePasswordInput{
		UserID:      1,
		OldPassword: "abc12345",
		NewPassword: "short1",
	}); err == nil || err.Error() != "密码长度需为8-64位" {
		t.Fatalf("ChangePassword weak password mismatch: %v", err)
	}
	if err := authSvc.ChangePassword(nil, &model.AuthChangePasswordInput{
		UserID:      1,
		OldPassword: "abc12345",
		NewPassword: "abc12345",
	}); err == nil || err.Error() != "新密码不能与旧密码相同" {
		t.Fatalf("ChangePassword same password mismatch: %v", err)
	}
}

func TestNormalizeAuthChangePasswordInput(t *testing.T) {
	in := &model.AuthChangePasswordInput{
		OldPassword: " old-pass ",
		NewPassword: " new-pass ",
	}
	normalizeAuthChangePasswordInput(in)
	if in.OldPassword != "old-pass" || in.NewPassword != "new-pass" {
		t.Fatalf("normalizeAuthChangePasswordInput mismatch: %+v", in)
	}
}

func TestCollectRoleHelpers(t *testing.T) {
	roles := []roleSnapshot{
		{ID: 2, Title: " 编辑 ", IsAdmin: 0},
		{ID: 1, Title: "管理员", IsAdmin: 1},
		{ID: 2, Title: "编辑", IsAdmin: 0},
	}
	if got := collectRoleIDs(roles); !reflect.DeepEqual(got, []int64{2, 1}) {
		t.Fatalf("collectRoleIDs mismatch: %v", got)
	}
	if got := collectRoleTitles(roles); !reflect.DeepEqual(got, []string{"编辑", "管理员"}) {
		t.Fatalf("collectRoleTitles mismatch: %v", got)
	}
	if !hasAdminRole(roles) {
		t.Fatal("hasAdminRole should detect admin role")
	}
}
