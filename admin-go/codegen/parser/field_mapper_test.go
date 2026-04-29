package parser

import (
	"database/sql"
	"testing"
)

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

func TestBuildFieldMetaEnumComponentAutoDerivation(t *testing.T) {
	// 2值 status → Switch
	status2 := buildFieldMeta(columnInfo{
		ColumnName:    "status",
		DataType:      "tinyint",
		ColumnType:    "tinyint(1)",
		IsNullable:    "NO",
		ColumnComment: "状态:0=关闭,1=开启",
	})
	if status2.Component != ComponentSwitch {
		t.Fatalf("status with 2 enum values should be Switch: got=%s", status2.Component)
	}

	// 3值 status → Radio
	status3 := buildFieldMeta(columnInfo{
		ColumnName:    "status",
		DataType:      "tinyint",
		ColumnType:    "tinyint(1)",
		IsNullable:    "NO",
		ColumnComment: "状态:0=草稿,1=已发布,2=已下架",
	})
	if status3.Component != ComponentRadio {
		t.Fatalf("status with 3 enum values should be Radio: got=%s", status3.Component)
	}

	// 4值 type → Select (非 status/is_ 前缀)
	type4 := buildFieldMeta(columnInfo{
		ColumnName:    "type",
		DataType:      "tinyint",
		ColumnType:    "tinyint(1)",
		IsNullable:    "NO",
		ColumnComment: "类型:1=普通,2=置顶,3=推荐,4=热门",
	})
	if type4.Component != ComponentSelect {
		t.Fatalf("type with 4 enum values should be Select: got=%s", type4.Component)
	}
}

func TestBuildFieldMetaEmptyCommentFallback(t *testing.T) {
	field := buildFieldMeta(columnInfo{
		ColumnName:    "extra_field",
		DataType:      "varchar",
		ColumnType:    "varchar(100)",
		IsNullable:    "YES",
		ColumnComment: "",
	})
	if field.Label != "ExtraField" || field.ShortLabel != "ExtraField" {
		t.Fatalf("empty comment should fallback to CamelCase: Label=%s ShortLabel=%s", field.Label, field.ShortLabel)
	}
}

func TestBuildFieldMetaValidationRules(t *testing.T) {
	// email
	email := buildFieldMeta(columnInfo{
		ColumnName:    "email",
		DataType:      "varchar",
		ColumnType:    "varchar(100)",
		IsNullable:    "NO",
		ColumnComment: "邮箱",
	})
	hasEmail := false
	for _, r := range email.ValidationRules {
		if r == "email" {
			hasEmail = true
		}
	}
	if !hasEmail {
		t.Fatalf("email field should have 'email' validation rule: %v", email.ValidationRules)
	}

	// phone
	phone := buildFieldMeta(columnInfo{
		ColumnName:    "phone",
		DataType:      "varchar",
		ColumnType:    "varchar(30)",
		IsNullable:    "NO",
		ColumnComment: "手机号",
	})
	hasPhone := false
	for _, r := range phone.ValidationRules {
		if r == "phone-loose" {
			hasPhone = true
		}
	}
	if !hasPhone {
		t.Fatalf("phone field should have 'phone-loose' validation rule: %v", phone.ValidationRules)
	}

	// url
	linkUrl := buildFieldMeta(columnInfo{
		ColumnName:    "link_url",
		DataType:      "varchar",
		ColumnType:    "varchar(500)",
		IsNullable:    "YES",
		ColumnComment: "链接",
	})
	hasUrl := false
	for _, r := range linkUrl.ValidationRules {
		if r == "url" {
			hasUrl = true
		}
	}
	if !hasUrl {
		t.Fatalf("link_url field should have 'url' validation rule: %v", linkUrl.ValidationRules)
	}

	// max-length
	title := buildFieldMeta(columnInfo{
		ColumnName:    "title",
		DataType:      "varchar",
		ColumnType:    "varchar(200)",
		IsNullable:    "NO",
		ColumnComment: "标题",
		CharMaxLength: func() sql.NullInt64 { return sql.NullInt64{Int64: 200, Valid: true} }(),
	})
	hasMaxLen := false
	for _, r := range title.ValidationRules {
		if r == "max-length:200" {
			hasMaxLen = true
		}
	}
	if !hasMaxLen {
		t.Fatalf("title field should have 'max-length:200' validation rule: %v", title.ValidationRules)
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
