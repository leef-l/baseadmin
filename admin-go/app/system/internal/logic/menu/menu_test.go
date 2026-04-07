package menu

import (
	"testing"

	"gbaseadmin/app/system/internal/model"
)

func TestNormalizeMenuInputs(t *testing.T) {
	createIn := &model.MenuCreateInput{
		Title:      " 系统管理 ",
		Path:       " /system ",
		Component:  " layouts/index ",
		Permission: " system:list ",
		Icon:       " setting ",
		LinkURL:    " https://example.com ",
	}
	normalizeMenuCreateInput(createIn)
	if createIn.Title != "系统管理" || createIn.Path != "/system" || createIn.Component != "layouts/index" || createIn.Permission != "system:list" || createIn.Icon != "setting" || createIn.LinkURL != "https://example.com" {
		t.Fatalf("normalizeMenuCreateInput mismatch: %+v", createIn)
	}

	updateIn := &model.MenuUpdateInput{
		Title:      " 菜单 ",
		Path:       " /menu ",
		Component:  " menu/index ",
		Permission: " menu:list ",
		Icon:       " menu ",
		LinkURL:    " https://foo.bar ",
	}
	normalizeMenuUpdateInput(updateIn)
	if updateIn.Title != "菜单" || updateIn.Path != "/menu" || updateIn.Component != "menu/index" || updateIn.Permission != "menu:list" || updateIn.Icon != "menu" || updateIn.LinkURL != "https://foo.bar" {
		t.Fatalf("normalizeMenuUpdateInput mismatch: %+v", updateIn)
	}

	listIn := &model.MenuListInput{Keyword: " sys "}
	normalizeMenuListInput(listIn)
	if listIn.Keyword != "sys" {
		t.Fatalf("normalizeMenuListInput mismatch: %+v", listIn)
	}

	treeIn := &model.MenuTreeInput{Keyword: " dashboard "}
	normalizeMenuTreeInput(treeIn)
	if treeIn.Keyword != "dashboard" {
		t.Fatalf("normalizeMenuTreeInput mismatch: %+v", treeIn)
	}
}
