package users

import (
	"reflect"
	"testing"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

func TestCompactRoleIDs(t *testing.T) {
	input := []snowflake.JsonInt64{0, 3, 3, -1, 2, 2, 5}
	want := []snowflake.JsonInt64{3, 2, 5}
	if got := compactRoleIDs(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("compactRoleIDs mismatch: got=%v want=%v", got, want)
	}
}

func TestCompactRoleIDsReturnsNilWhenEmpty(t *testing.T) {
	if got := compactRoleIDs(nil); got != nil {
		t.Fatalf("compactRoleIDs(nil) should return nil, got %v", got)
	}
	if got := compactRoleIDs([]snowflake.JsonInt64{0, -1}); got != nil {
		t.Fatalf("compactRoleIDs should drop non-positive ids and return nil, got %v", got)
	}
}

func TestNormalizeUsersInputs(t *testing.T) {
	createIn := &model.UsersCreateInput{
		Username: " admin ",
		Nickname: " 管理员 ",
		Email:    " admin@example.com ",
		Avatar:   " /avatar.png ",
	}
	normalizeUsersWriteInput(createIn)
	if createIn.Username != "admin" || createIn.Nickname != "管理员" || createIn.Email != "admin@example.com" || createIn.Avatar != "/avatar.png" {
		t.Fatalf("normalizeUsersWriteInput mismatch: %+v", createIn)
	}

	listIn := &model.UsersListInput{
		Keyword:  " demo ",
		Username: " admin ",
		Nickname: " nick ",
		Email:    " user@example.com ",
	}
	normalizeUsersListInput(listIn)
	if listIn.Keyword != "demo" || listIn.Username != "admin" || listIn.Nickname != "nick" || listIn.Email != "user@example.com" {
		t.Fatalf("normalizeUsersListInput mismatch: %+v", listIn)
	}
}
