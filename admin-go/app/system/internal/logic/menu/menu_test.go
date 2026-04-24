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

func TestMenuInputValidation(t *testing.T) {
	menuSvc := &sMenu{}
	if err := menuSvc.Create(nil, nil); err == nil || err.Error() != "请求参数不能为空" {
		t.Fatalf("Create nil input mismatch: %v", err)
	}
	if err := menuSvc.Create(nil, &model.MenuCreateInput{Title: " "}); err == nil || err.Error() != "菜单名称不能为空" {
		t.Fatalf("Create blank title mismatch: %v", err)
	}
	if err := menuSvc.Update(nil, &model.MenuUpdateInput{ID: 1, Title: " "}); err == nil || err.Error() != "菜单名称不能为空" {
		t.Fatalf("Update blank title mismatch: %v", err)
	}
	if err := validateMenuFields("系统管理", 0, "", "", "", "", 1, 0, 1, 1); err == nil || err.Error() != "菜单类型值不合法" {
		t.Fatalf("validateMenuFields invalid type mismatch: %v", err)
	}
	if err := validateMenuFields("系统管理", 1, "/system", "", "", "", 1, 0, 1, -1); err == nil || err.Error() != "排序不能小于0" {
		t.Fatalf("validateMenuFields negative sort mismatch: %v", err)
	}
	if err := validateMenuFields("系统管理", 1, "/system", "", "", "", 3, 0, 1, 0); err == nil || err.Error() != "是否显示值不合法" {
		t.Fatalf("validateMenuFields invalid isShow mismatch: %v", err)
	}
	if err := validateMenuFields("菜单", 2, "/system/menu", "", "", "", 1, 0, 1, 0); err == nil || err.Error() != "菜单类型必须填写前端组件路径" {
		t.Fatalf("validateMenuFields missing component mismatch: %v", err)
	}
	if err := validateMenuFields("按钮", 3, "", "", "", "", 1, 0, 1, 0); err == nil || err.Error() != "按钮类型必须填写权限标识" {
		t.Fatalf("validateMenuFields missing permission mismatch: %v", err)
	}
	if err := validateMenuFields("外链", 4, "/docs", "", "", "", 1, 0, 1, 0); err == nil || err.Error() != "外链/内链类型必须填写地址" {
		t.Fatalf("validateMenuFields missing linkURL mismatch: %v", err)
	}
	if err := validateMenuFields("目录", 1, "system", "", "", "", 1, 0, 1, 0); err == nil || err.Error() != "前端路由路径必须以 / 开头" {
		t.Fatalf("validateMenuFields invalid path prefix mismatch: %v", err)
	}
	if _, err := menuSvc.Detail(nil, 0); err == nil || err.Error() != "菜单不存在或已删除" {
		t.Fatalf("Detail invalid id mismatch: %v", err)
	}
}

func TestNormalizeMenuTypeFields(t *testing.T) {
	path := "/system/button"
	component := "system/button/index"
	linkURL := "https://example.com"
	normalizeMenuTypeFields(3, &path, &component, &linkURL)
	if path != "" || component != "" || linkURL != "" {
		t.Fatalf("normalizeMenuTypeFields button mismatch: path=%q component=%q linkURL=%q", path, component, linkURL)
	}

	isCache := 1
	normalizeMenuCache(&isCache, 1)
	if isCache != 0 {
		t.Fatalf("normalizeMenuCache directory mismatch: got=%d", isCache)
	}

	isCache = 1
	normalizeMenuCache(&isCache, 4)
	if isCache != 0 {
		t.Fatalf("normalizeMenuCache link mismatch: got=%d", isCache)
	}

	isCache = 2
	normalizeMenuCache(&isCache, 2)
	if isCache != 0 {
		t.Fatalf("normalizeMenuCache invalid menu cache mismatch: got=%d", isCache)
	}

	isCache = 1
	normalizeMenuCache(&isCache, 2)
	if isCache != 1 {
		t.Fatalf("normalizeMenuCache menu mismatch: got=%d", isCache)
	}
}
