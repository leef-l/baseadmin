package parser

import "testing"

func TestParseCommentMetaSupportsTooltipEnumAndDirectives(t *testing.T) {
	meta := ParseCommentMeta("状态（用于展示）:0=关闭,1=开启|ref:system_users.username|search:eq|keyword:off|priority:88")

	if meta.Label != "状态（用于展示）" || meta.ShortLabel != "状态" || meta.TooltipText != "用于展示" {
		t.Fatalf("label parsing mismatch: %+v", meta)
	}
	if len(meta.EnumValues) != 2 || meta.EnumValues[1].NameIdent != "On" {
		t.Fatalf("enum parsing mismatch: %+v", meta.EnumValues)
	}
	if meta.RefTableHint != "system_users" || meta.RefDisplayHint != "username" {
		t.Fatalf("ref parsing mismatch: %+v", meta)
	}
	if meta.SearchMode != "eq" || meta.KeywordMode != "off" || meta.SearchPriority != 88 {
		t.Fatalf("directive parsing mismatch: %+v", meta)
	}
}

func TestParseCommentSupportsDictSyntax(t *testing.T) {
	meta := ParseCommentMeta("性别:dict:gender")
	if meta.Label != "性别" || meta.DictType != "gender" {
		t.Fatalf("dict parsing mismatch: %+v", meta)
	}
	if len(meta.EnumValues) != 0 {
		t.Fatalf("dict syntax should not parse enum values: %+v", meta.EnumValues)
	}
}

func TestParseCommentReturnsSimpleFields(t *testing.T) {
	label, shortLabel, tooltip, enums := ParseComment("排序（升序）")
	if label != "排序（升序）" || shortLabel != "排序" || tooltip != "升序" {
		t.Fatalf("ParseComment mismatch: label=%q short=%q tooltip=%q", label, shortLabel, tooltip)
	}
	if len(enums) != 0 {
		t.Fatalf("simple ParseComment should not include enums: %+v", enums)
	}
}
