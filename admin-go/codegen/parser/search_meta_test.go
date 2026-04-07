package parser

import "testing"

func TestApplySearchMetaDefaultHeuristics(t *testing.T) {
	username := buildFieldMeta(columnInfo{
		ColumnName:    "username",
		DataType:      "varchar",
		ColumnType:    "varchar(64)",
		IsNullable:    "NO",
		ColumnComment: "登录用户名",
	})
	if !username.SearchEnabled || username.SearchOperator != "like" || username.SearchComponent != "Input" {
		t.Fatalf("username search meta mismatch: %+v", username)
	}
	if !username.KeywordEnabled {
		t.Fatalf("username should participate in keyword search")
	}

	orderNo := buildFieldMeta(columnInfo{
		ColumnName:    "order_no",
		DataType:      "varchar",
		ColumnType:    "varchar(64)",
		IsNullable:    "NO",
		ColumnComment: "订单号",
	})
	if !orderNo.SearchEnabled || orderNo.SearchOperator != "eq" {
		t.Fatalf("order_no should use exact search: %+v", orderNo)
	}
	if orderNo.KeywordEnabled {
		t.Fatalf("order_no should not participate in keyword search by default")
	}

	status := buildFieldMeta(columnInfo{
		ColumnName:    "status",
		DataType:      "tinyint",
		ColumnType:    "tinyint(1)",
		IsNullable:    "NO",
		ColumnComment: "状态:0=关闭,1=开启",
	})
	if !status.SearchEnabled || status.SearchComponent != "Select" || status.SearchOperator != "eq" {
		t.Fatalf("status should use select exact search: %+v", status)
	}

	startAt := buildFieldMeta(columnInfo{
		ColumnName:    "start_at",
		DataType:      "datetime",
		ColumnType:    "datetime",
		IsNullable:    "YES",
		ColumnComment: "开始时间",
	})
	if !startAt.SearchEnabled || !startAt.SearchRange || startAt.SearchFormField != "startAtRange" {
		t.Fatalf("start_at should use range search: %+v", startAt)
	}
}

func TestApplySearchMetaCommentOverrides(t *testing.T) {
	field := buildFieldMeta(columnInfo{
		ColumnName:    "summary",
		DataType:      "varchar",
		ColumnType:    "varchar(255)",
		IsNullable:    "YES",
		ColumnComment: "摘要|search:off|keyword:only|priority:99",
	})
	if field.SearchEnabled {
		t.Fatalf("summary should not render individual search control: %+v", field)
	}
	if !field.KeywordEnabled || !field.KeywordOnly {
		t.Fatalf("summary should be keyword-only: %+v", field)
	}
	if field.SearchPriority != 99 {
		t.Fatalf("summary priority mismatch: %+v", field)
	}

	explicit := buildFieldMeta(columnInfo{
		ColumnName:    "remark",
		DataType:      "varchar",
		ColumnType:    "varchar(255)",
		IsNullable:    "YES",
		ColumnComment: "备注|search:eq|keyword:off|priority:88",
	})
	if !explicit.SearchEnabled || explicit.SearchOperator != "eq" || explicit.SearchComponent != "Input" {
		t.Fatalf("remark explicit eq search mismatch: %+v", explicit)
	}
	if explicit.KeywordEnabled {
		t.Fatalf("remark should not participate in keyword search: %+v", explicit)
	}
	if explicit.SearchPriority != 88 {
		t.Fatalf("remark priority mismatch: %+v", explicit)
	}
}

func TestFinalizeSearchMetaKeywordPromotion(t *testing.T) {
	meta := &TableMeta{
		Fields: []FieldMeta{
			{Name: "username", SearchEnabled: true, SearchOperator: "like", SearchComponent: "Input", KeywordEnabled: true, SearchPriority: 95},
			{Name: "nickname", SearchEnabled: true, SearchOperator: "like", SearchComponent: "Input", KeywordEnabled: true, SearchPriority: 94},
			{Name: "email", SearchEnabled: true, SearchOperator: "like", SearchComponent: "Input", KeywordEnabled: true, SearchPriority: 90},
			{Name: "dept_id", SearchEnabled: true, SearchOperator: "eq", SearchComponent: "Select", SearchPriority: 82},
		},
	}

	finalizeSearchMeta(meta)

	if !meta.HasKeywordSearch {
		t.Fatalf("keyword search should be enabled")
	}
	if len(meta.SearchFields) != 3 {
		t.Fatalf("expected 3 visible search fields, got %d", len(meta.SearchFields))
	}
	if meta.SearchFields[0].Name != "username" || meta.SearchFields[1].Name != "nickname" || meta.SearchFields[2].Name != "dept_id" {
		t.Fatalf("unexpected visible search fields: %+v", meta.SearchFields)
	}
	if len(meta.KeywordSearchFields) != 3 {
		t.Fatalf("expected 3 keyword search fields, got %d", len(meta.KeywordSearchFields))
	}
}
