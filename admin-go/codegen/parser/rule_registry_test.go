package parser

import (
	"reflect"
	"testing"
)

func TestRuleRegistryFieldNameHelpers(t *testing.T) {
	if !isHiddenFieldName("created_at") || isHiddenFieldName("title") {
		t.Fatal("hidden field registry mismatch")
	}
	if !isImageFieldName("cover") || !isImageFieldName("product_logo") || isImageFieldName("title") {
		t.Fatal("image field registry mismatch")
	}
	if !isSearchableTextFieldName("title") || !isSearchableTextFieldName("user_name") || isSearchableTextFieldName("amount") {
		t.Fatal("searchable text registry mismatch")
	}
	if !isExactSearchFieldName("order_no") || !isExactSearchFieldName("code") || isExactSearchFieldName("title") {
		t.Fatal("exact search registry mismatch")
	}
	if !isMoneyFieldName("price") || !isMoneyFieldName("service_fee") || isMoneyFieldName("title") {
		t.Fatal("money field registry mismatch")
	}
}

func TestDisplayFieldPriorityOrderStable(t *testing.T) {
	got := displayFieldPriorityOrder()
	want := []string{
		"title", "sku_no", "name", "username", "nickname",
		"real_name", "label", "phone", "mobile",
		"order_no", "code", "no", "email",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("display priority mismatch: got=%v want=%v", got, want)
	}

	got[0] = "mutated"
	if again := displayFieldPriorityOrder(); again[0] != "title" {
		t.Fatalf("display priority should return a copy, got=%v", again)
	}
}
