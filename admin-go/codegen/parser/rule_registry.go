package parser

type fieldRuleRegistry struct {
	displayFieldPriority     []string
	hiddenFieldNames         map[string]struct{}
	imageFieldNames          map[string]struct{}
	imageFieldSuffixes       []string
	searchableTextFieldNames map[string]struct{}
	searchableTextSuffixes   []string
	exactSearchFieldNames    map[string]struct{}
	exactSearchFieldSuffixes []string
	moneyFieldNames          map[string]struct{}
	moneyFieldSuffixes       []string
}

var rules = fieldRuleRegistry{
	displayFieldPriority: []string{
		"title", "sku_no", "name", "username", "nickname",
		"real_name", "label", "phone", "mobile",
		"order_no", "code", "no", "email",
	},
	hiddenFieldNames: toStructSet(
		"id", "created_at", "updated_at", "deleted_at", "created_by", "dept_id",
	),
	imageFieldNames: toStructSet(
		"avatar", "cover", "logo", "banner", "thumbnail", "poster",
	),
	imageFieldSuffixes: []string{
		"_image", "_img", "_photo", "_pic", "_cover", "_banner", "_logo", "_thumbnail", "_poster", "_avatar",
	},
	searchableTextFieldNames: toStructSet(
		"title", "name", "username", "nickname",
		"phone", "mobile", "email", "real_name",
		"order_no", "remark", "description", "summary",
		"intro", "address", "contact", "contact_name",
		"link_url", "url", "keyword",
	),
	searchableTextSuffixes: []string{
		"_name", "_title", "_remark", "_desc", "_description", "_summary",
		"_intro", "_phone", "_mobile", "_email", "_address", "_keyword", "_no",
	},
	exactSearchFieldNames: toStructSet("no", "code", "sn"),
	exactSearchFieldSuffixes: []string{
		"_no", "_code", "_sn",
	},
	moneyFieldNames: toStructSet(
		"price", "amount", "balance", "income_total", "income_balance",
	),
	moneyFieldSuffixes: []string{
		"_price", "_amount", "_balance", "_income", "_fee", "_cost", "_deposit", "_refund",
	},
}

func toStructSet(items ...string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

func hasRuleName(set map[string]struct{}, name string) bool {
	_, ok := set[name]
	return ok
}

func hasRuleSuffix(name string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if len(name) >= len(suffix) && name[len(name)-len(suffix):] == suffix {
			return true
		}
	}
	return false
}

func isHiddenFieldName(name string) bool {
	return hasRuleName(rules.hiddenFieldNames, name)
}

func isImageFieldName(name string) bool {
	return hasRuleName(rules.imageFieldNames, name) || hasRuleSuffix(name, rules.imageFieldSuffixes)
}

func isSearchableTextFieldName(name string) bool {
	return hasRuleName(rules.searchableTextFieldNames, name) || hasRuleSuffix(name, rules.searchableTextSuffixes)
}

func isExactSearchFieldName(name string) bool {
	return hasRuleName(rules.exactSearchFieldNames, name) || hasRuleSuffix(name, rules.exactSearchFieldSuffixes)
}

func isMoneyFieldName(name string) bool {
	return hasRuleName(rules.moneyFieldNames, name) || hasRuleSuffix(name, rules.moneyFieldSuffixes)
}

func displayFieldPriorityOrder() []string {
	return append([]string(nil), rules.displayFieldPriority...)
}
