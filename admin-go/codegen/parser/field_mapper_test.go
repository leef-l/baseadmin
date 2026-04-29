package parser

import "testing"

func TestBuildFieldMetaDerivesStableFlags(t *testing.T) {
	orderNo := buildFieldMeta(columnInfo{
		ColumnName:    "order_no",
		DataType:      "varchar",
		ColumnType:    "varchar(64)",
		IsNullable:    "NO",
		ColumnComment: "订单号",
	})
	if !orderNo.IsSearchable || !orderNo.IsExactSearch {
		t.Fatalf("order_no flags mismatch: %+v", orderNo)
	}

	balance := buildFieldMeta(columnInfo{
		ColumnName:    "income_balance",
		DataType:      "int",
		ColumnType:    "int(11)",
		IsNullable:    "NO",
		ColumnComment: "余额",
	})
	if !balance.IsMoney {
		t.Fatalf("income_balance should be treated as money: %+v", balance)
	}
	if balance.Component != ComponentInputNumber {
		t.Fatalf("income_balance component mismatch: %+v", balance)
	}
}

func TestBuildFieldMetaAddsSystemScopeRefHints(t *testing.T) {
	tenant := buildFieldMeta(columnInfo{
		ColumnName:    "tenant_id",
		DataType:      "bigint",
		ColumnType:    "bigint unsigned",
		IsNullable:    "NO",
		ColumnComment: "租户ID，0 表示平台",
	})
	if tenant.RefTableHint != "system_tenant" || tenant.RefDisplayHint != "name" {
		t.Fatalf("tenant_id ref hint mismatch: %+v", tenant)
	}
	if tenant.Label != "租户" || tenant.ShortLabel != "租户" || tenant.TooltipText != "" {
		t.Fatalf("tenant_id presentation mismatch: %+v", tenant)
	}

	merchant := buildFieldMeta(columnInfo{
		ColumnName:    "merchant_id",
		DataType:      "bigint",
		ColumnType:    "bigint unsigned",
		IsNullable:    "NO",
		ColumnComment: "商户ID，0 表示租户级/平台级",
	})
	if merchant.RefTableHint != "system_merchant" || merchant.RefDisplayHint != "name" {
		t.Fatalf("merchant_id ref hint mismatch: %+v", merchant)
	}
	if merchant.Label != "商户" || merchant.ShortLabel != "商户" || merchant.TooltipText != "" {
		t.Fatalf("merchant_id presentation mismatch: %+v", merchant)
	}
}

func TestMapComponentRecognizesCoreComponentTypes(t *testing.T) {
	tests := []struct {
		name  string
		field FieldMeta
		want  string
	}{
		{name: "parent", field: FieldMeta{Name: "parent_id"}, want: ComponentTreeSelectSingle},
		{name: "multi ids", field: FieldMeta{Name: "role_ids"}, want: ComponentSelectMulti},
		{name: "image exact", field: FieldMeta{Name: "cover"}, want: ComponentImageUpload},
		{name: "image suffix", field: FieldMeta{Name: "product_logo"}, want: ComponentImageUpload},
		{name: "password", field: FieldMeta{Name: "api_secret_key"}, want: ComponentPassword},
		{name: "url", field: FieldMeta{Name: "link_url"}, want: ComponentInputUrl},
	}

	for _, tc := range tests {
		if got := MapComponent(tc.field); got != tc.want {
			t.Fatalf("%s component mismatch: got=%s want=%s", tc.name, got, tc.want)
		}
	}
}
