package v1

import "github.com/gogf/gf/v2/frame/g"

// PortalBizConfigReq C 端拉取业务时间配置（用于商城/寄售页倒计时）。
type PortalBizConfigReq struct {
	g.Meta `path:"/biz-config" method:"get" tags:"会员-公共" summary:"业务时间配置"`
}

// PortalBizConfigRes 业务配置（脱敏，仅返回时间窗与开放工作日）。
type PortalBizConfigRes struct {
	g.Meta          `mime:"application/json"`
	PurchaseStart   string  `json:"purchaseStart" dc:"进货开始 HH:MM"`
	PurchaseEnd     string  `json:"purchaseEnd"   dc:"进货结束 HH:MM"`
	PurchaseDays    []int   `json:"purchaseDays"  dc:"允许进货的工作日 1=Mon...7=Sun"`
	ConsignStart    string  `json:"consignStart"  dc:"寄售开始 HH:MM"`
	ConsignEnd      *string `json:"consignEnd"    dc:"寄售结束 HH:MM（null 表示无截止）"`
	ServerTimestamp int64   `json:"serverTimestamp" dc:"服务器当前时间戳（秒），前端用于倒计时校时"`
}
