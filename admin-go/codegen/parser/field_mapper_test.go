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

func TestBuildFieldMetaTokenSuffixNarrowed(t *testing.T) {
	inviteToken := buildFieldMeta(columnInfo{
		ColumnName:    "invite_token",
		DataType:      "varchar",
		ColumnType:    "varchar(255)",
		IsNullable:    "NO",
		ColumnComment: "邀请令牌",
	})
	if inviteToken.IsPassword {
		t.Fatalf("invite_token should NOT be treated as password: %+v", inviteToken)
	}

	apiToken := buildFieldMeta(columnInfo{
		ColumnName:    "wechat_api_token",
		DataType:      "varchar",
		ColumnType:    "varchar(255)",
		IsNullable:    "NO",
		ColumnComment: "微信API令牌",
	})
	if !apiToken.IsPassword {
		t.Fatalf("wechat_api_token should be treated as password: %+v", apiToken)
	}

	deviceToken := buildFieldMeta(columnInfo{
		ColumnName:    "device_token",
		DataType:      "varchar",
		ColumnType:    "varchar(255)",
		IsNullable:    "NO",
		ColumnComment: "设备推送令牌",
	})
	if deviceToken.IsPassword {
		t.Fatalf("device_token should NOT be treated as password: %+v", deviceToken)
	}
}

func TestBuildFieldMetaImageAvatarSuffix(t *testing.T) {
	userAvatar := buildFieldMeta(columnInfo{
		ColumnName:    "user_avatar",
		DataType:      "varchar",
		ColumnType:    "varchar(500)",
		IsNullable:    "YES",
		ColumnComment: "用户头像",
	})
	if userAvatar.Component != ComponentImageUpload {
		t.Fatalf("user_avatar should map to ImageUpload: got=%s", userAvatar.Component)
	}
}

func TestBuildFieldMetaMoneySuffixExpanded(t *testing.T) {
	deposit := buildFieldMeta(columnInfo{
		ColumnName:    "security_deposit",
		DataType:      "int",
		ColumnType:    "int(11)",
		IsNullable:    "NO",
		ColumnComment: "押金",
	})
	if !deposit.IsMoney {
		t.Fatalf("security_deposit should be treated as money: %+v", deposit)
	}

	refund := buildFieldMeta(columnInfo{
		ColumnName:    "order_refund",
		DataType:      "int",
		ColumnType:    "int(11)",
		IsNullable:    "NO",
		ColumnComment: "退款金额",
	})
	if !refund.IsMoney {
		t.Fatalf("order_refund should be treated as money: %+v", refund)
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

func TestBuildFieldMetaStatusWithoutEnumFallsBackToInput(t *testing.T) {
	status := buildFieldMeta(columnInfo{
		ColumnName:    "status",
		DataType:      "tinyint",
		ColumnType:    "tinyint(4)",
		IsNullable:    "NO",
		ColumnComment: "状态",
	})
	if status.Component != ComponentInput {
		t.Fatalf("status without enum should fallback to Input: got=%s", status.Component)
	}

	isActive := buildFieldMeta(columnInfo{
		ColumnName:    "is_active",
		DataType:      "tinyint",
		ColumnType:    "tinyint(1)",
		IsNullable:    "NO",
		ColumnComment: "是否激活",
	})
	if isActive.Component != ComponentSwitch {
		t.Fatalf("is_active without enum should default to Switch: got=%s", isActive.Component)
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
		{name: "image avatar suffix", field: FieldMeta{Name: "user_avatar"}, want: ComponentImageUpload},
		{name: "password", field: FieldMeta{Name: "api_secret_key", IsPassword: true}, want: ComponentPassword},
		{name: "url", field: FieldMeta{Name: "link_url"}, want: ComponentInputUrl},
		{name: "richtext exact", field: FieldMeta{Name: "content"}, want: ComponentRichText},
		{name: "richtext suffix", field: FieldMeta{Name: "body_content"}, want: ComponentRichText},
		{name: "datetime", field: FieldMeta{Name: "publish_at"}, want: ComponentDateTimePicker},
		{name: "json editor", field: FieldMeta{Name: "extra_json"}, want: ComponentJsonEditor},
		{name: "file upload", field: FieldMeta{Name: "contract_file"}, want: ComponentFileUpload},
	}

	for _, tc := range tests {
		if got := MapComponent(tc.field); got != tc.want {
			t.Fatalf("%s component mismatch: got=%s want=%s", tc.name, got, tc.want)
		}
	}
}
