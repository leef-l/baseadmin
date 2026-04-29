// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoCampaign is the golang structure of table demo_campaign for DAO operations like Where/Data.
type DemoCampaign struct {
	g.Meta       `orm:"table:demo_campaign, do:true"`
	Id           any         // 活动ID（Snowflake）
	CampaignNo   any         // 活动编号|search:eq|priority:100
	Title        any         // 活动标题|search:like|keyword:on|priority:95
	Banner       any         // 横幅图
	Type         any         // 活动类型:1=免费,2=付费,3=公开,4=私密
	Channel      any         // 投放渠道:1=官网,2=小程序,3=短信,4=线下
	BudgetAmount any         // 预算金额（分）
	LandingUrl   any         // 落地页URL
	RuleJson     any         // 规则JSON
	IntroContent any         // 活动介绍
	StartAt      *gtime.Time // 开始时间
	EndAt        *gtime.Time // 结束时间
	IsPublic     any         // 是否公开:0=否,1=是
	Status       any         // 状态:0=草稿,1=已发布,2=已下架
	TenantId     any         // 租户
	MerchantId   any         // 商户
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间，非 NULL 表示已删除
}
