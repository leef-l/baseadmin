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
