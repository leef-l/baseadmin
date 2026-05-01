package v1

import "github.com/gogf/gf/v2/frame/g"

// ----- 获取 -----

// MemberBizConfigGetReq 获取会员业务配置（单例）。
type MemberBizConfigGetReq struct {
	g.Meta `path:"/member/biz-config" method:"get" tags:"会员业务配置" summary:"获取业务配置"`
}

// MemberBizConfigGetRes 业务配置完整内容。
type MemberBizConfigGetRes struct {
	g.Meta `mime:"application/json"`
	*BizConfigData
}

// ----- 保存 -----

// MemberBizConfigSaveReq 整体覆盖业务配置。
type MemberBizConfigSaveReq struct {
	g.Meta `path:"/member/biz-config" method:"put" tags:"会员业务配置" summary:"保存业务配置"`
	*BizConfigData
}

// MemberBizConfigSaveRes 保存响应。
type MemberBizConfigSaveRes struct {
	g.Meta `mime:"application/json"`
}

// BizConfigData 与后端 logic.bizconfig.Config 同构。
type BizConfigData struct {
	Purchase               BizConfigPurchase   `json:"purchase"`
	Consign                BizConfigConsign    `json:"consign"`
	SelfRebateTiers        []BizConfigRebate   `json:"selfRebateTiers"`
	SelfTurnoverRewardRate float64             `json:"selfTurnoverRewardRate" v:"min:0|max:100" dc:"自购按金额比例返奖励钱包（百分比）"`
	DirectPromoteRate      float64             `json:"directPromoteRate"      v:"min:0|max:100" dc:"直推按金额比例返推广钱包（百分比）"`
}

// BizConfigPurchase 进货时间窗 + 工作日。
type BizConfigPurchase struct {
	StartTime       string `json:"startTime"       v:"required|regex:^\\d{2}:\\d{2}$#进货开始时间不能为空|格式 HH:MM"`
	EndTime         string `json:"endTime"         v:"required|regex:^\\d{2}:\\d{2}$#进货结束时间不能为空|格式 HH:MM"`
	AllowedWeekdays []int  `json:"allowedWeekdays" v:"required" dc:"允许的工作日 1=Mon...7=Sun"`
}

// BizConfigConsign 寄售时间窗（endTime 为空表示无截止）。
type BizConfigConsign struct {
	StartTime string  `json:"startTime" v:"required|regex:^\\d{2}:\\d{2}$#寄售开始时间不能为空|格式 HH:MM"`
	EndTime   *string `json:"endTime"`
}

// BizConfigRebate 自购阶梯返佣档位。
type BizConfigRebate struct {
	NthOrder   int     `json:"nthOrder"   v:"required|min:1" dc:"第 N 单触发"`
	RewardYuan float64 `json:"rewardYuan" v:"required|min:0" dc:"奖励金额（元）"`
}
