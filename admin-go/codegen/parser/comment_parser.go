package parser

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

var enumLabelIdents = map[string]string{
	"启用": "Enabled", "禁用": "Disabled",
	"正常": "Normal", "异常": "Abnormal",
	"有效": "Valid", "无效": "Invalid",
	"是": "Yes", "否": "No",
	"男": "Male", "女": "Female",
	"开启": "On", "关闭": "Off",
	"显示": "Show", "隐藏": "Hide",
	"已完成": "Done", "进行中": "InProgress",
	"待处理": "Pending", "已取消": "Cancelled",
	"待审核": "PendingReview", "已通过": "Approved", "已拒绝": "Rejected",
	"待支付": "Unpaid", "已支付": "Paid", "已退款": "Refunded",
	"草稿": "Draft", "已发布": "Published", "已下架": "Offline",
	"目录": "Dir", "菜单": "Menu", "按钮": "Button",
	"普通": "Regular", "VIP": "VIP", "管理员": "Admin",
	"成功": "Success", "失败": "Failed",
	"充值": "Recharge", "消费": "Consume", "提现": "Withdraw",
	"置顶": "Pinned", "推荐": "Recommended", "热门": "Hot", "精华": "Featured",
	"外部链接": "ExternalLink", "内部链接": "InternalLink",
	"待确认": "Unconfirmed", "已确认": "Confirmed",
	"待发货": "Unshipped", "已发货": "Shipped", "已签收": "Received",
	"冻结": "Frozen", "解冻": "Unfrozen",
	"上架": "Online", "下架": "OffShelf",
	"免费": "Free", "付费": "Paid",
	"公开": "Public", "私密": "Private",
	"全部": "All", "本部门及以下": "DeptAndBelow", "本部门": "DeptOnly", "仅本人": "SelfOnly", "自定义": "Custom",
	"内容": "Content", "行为": "Behavior", "申诉": "Appeal",
}

type CommentMeta struct {
	Label          string
	ShortLabel     string
	TooltipText    string
	EnumValues     []EnumValue
	DictType       string
	RefTableHint   string
	RefDisplayHint string
	SearchMode     string
	KeywordMode    string
	SearchPriority int
}

// extractParentheses 从 label 中提取中文括号（）或英文括号()内的内容
// "排序（升序）" → shortLabel="排序", tooltip="升序"
// "部门名称"     → shortLabel="部门名称", tooltip=""
func extractParentheses(label string) (shortLabel string, tooltip string) {
	// 优先匹配中文括号（使用 LastIndex 匹配最外层右括号）
	if idx := strings.Index(label, "（"); idx >= 0 {
		if end := strings.LastIndex(label, "）"); end > idx {
			shortLabel = strings.TrimSpace(label[:idx])
			tooltip = strings.TrimSpace(label[idx+len("（") : end])
			return
		}
	}
	// 再匹配英文括号
	if idx := strings.Index(label, "("); idx >= 0 {
		if end := strings.LastIndex(label, ")"); end > idx {
			shortLabel = strings.TrimSpace(label[:idx])
			tooltip = strings.TrimSpace(label[idx+1 : end])
			return
		}
	}
	return label, ""
}

// ParseComment 解析字段备注
// 输入：状态:0=关闭,1=开启
// 输出：label="状态", shortLabel="状态", tooltipText="", enums=[{0,关闭},{1,开启}]
// 输入：排序（升序）
// 输出：label="排序（升序）", shortLabel="排序", tooltipText="升序", enums=[]
// 输入：部门名称
// 输出：label="部门名称", shortLabel="部门名称", tooltipText="", enums=[]
func ParseCommentMeta(comment string) CommentMeta {
	comment = strings.TrimSpace(comment)
	if comment == "" {
		return CommentMeta{}
	}

	parts := strings.Split(comment, "|")
	mainPart := strings.TrimSpace(parts[0])
	meta := CommentMeta{}
	for _, directive := range parts[1:] {
		directive = strings.TrimSpace(directive)
		if strings.HasPrefix(directive, "ref:") {
			refValue := strings.TrimSpace(strings.TrimPrefix(directive, "ref:"))
			if refValue == "" {
				continue
			}
			if dotIdx := strings.LastIndex(refValue, "."); dotIdx > 0 && dotIdx < len(refValue)-1 {
				meta.RefTableHint = strings.TrimSpace(refValue[:dotIdx])
				meta.RefDisplayHint = strings.TrimSpace(refValue[dotIdx+1:])
			} else {
				meta.RefTableHint = refValue
			}
			continue
		}
		if strings.HasPrefix(directive, "search:") {
			meta.SearchMode = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(directive, "search:")))
			continue
		}
		if strings.HasPrefix(directive, "keyword:") {
			meta.KeywordMode = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(directive, "keyword:")))
			continue
		}
		if strings.HasPrefix(directive, "priority:") {
			if v, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(directive, "priority:"))); err == nil {
				meta.SearchPriority = v
			}
			continue
		}
		if strings.HasPrefix(directive, "search-priority:") {
			if v, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(directive, "search-priority:"))); err == nil {
				meta.SearchPriority = v
			}
			continue
		}
	}

	// 查找冒号分隔符（支持中文冒号和英文冒号）
	sepIdx := -1
	sepLen := 0
	for i, ch := range mainPart {
		if ch == ':' || ch == '：' {
			sepIdx = i
			sepLen = utf8.RuneLen(ch)
			break
		}
	}

	// 没有冒号，整个备注就是 label
	if sepIdx < 0 {
		meta.Label = mainPart
		meta.ShortLabel, meta.TooltipText = extractParentheses(meta.Label)
		return meta
	}

	meta.Label = strings.TrimSpace(mainPart[:sepIdx])
	enumPart := strings.TrimSpace(mainPart[sepIdx+sepLen:])
	if len(meta.Label) == 0 {
		meta.Label = mainPart
		meta.ShortLabel, meta.TooltipText = extractParentheses(meta.Label)
		return meta
	}

	meta.ShortLabel, meta.TooltipText = extractParentheses(meta.Label)

	// 解析枚举部分：0=关闭,1=开启
	if enumPart == "" {
		return meta
	}

	// 如果是字典引用（如 dict:gender），不解析枚举
	if strings.HasPrefix(enumPart, "dict:") {
		meta.DictType = strings.TrimSpace(strings.TrimPrefix(enumPart, "dict:"))
		return meta
	}

	enumPart = strings.ReplaceAll(enumPart, "，", ",")
	pairs := strings.Split(enumPart, ",")
	seenValues := make(map[string]struct{})
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		eqIdx := strings.Index(pair, "=")
		if eqIdx < 0 {
			continue
		}
		val := strings.TrimSpace(pair[:eqIdx])
		lab := strings.TrimSpace(pair[eqIdx+1:])
		if val == "" || lab == "" {
			continue
		}
		if _, exists := seenValues[val]; exists {
			continue
		}
		seenValues[val] = struct{}{}
		meta.EnumValues = append(meta.EnumValues, EnumValue{Value: val, Label: lab, NameIdent: labelToIdent(lab)})
	}

	return meta
}

func ParseComment(comment string) (label string, shortLabel string, tooltipText string, enums []EnumValue) {
	meta := ParseCommentMeta(comment)
	return meta.Label, meta.ShortLabel, meta.TooltipText, meta.EnumValues
}

// labelToIdent 将中文枚举标签转为语义化 Go 标识符
func labelToIdent(label string) string {
	if ident, ok := enumLabelIdents[label]; ok {
		return ident
	}
	return ""
}
