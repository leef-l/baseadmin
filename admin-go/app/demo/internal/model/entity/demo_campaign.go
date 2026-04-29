// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoCampaign is the golang structure for table demo_campaign.
type DemoCampaign struct {
	Id           uint64      `orm:"id"            description:"活动ID（Snowflake）"`                         // 活动ID（Snowflake）
	CampaignNo   string      `orm:"campaign_no"   description:"活动编号|search:eq|priority:100"`             // 活动编号|search:eq|priority:100
	Title        string      `orm:"title"         description:"活动标题|search:like|keyword:on|priority:95"` // 活动标题|search:like|keyword:on|priority:95
	Banner       string      `orm:"banner"        description:"横幅图"`                                     // 横幅图
	Type         int         `orm:"type"          description:"活动类型:1=免费,2=付费,3=公开,4=私密"`                // 活动类型:1=免费,2=付费,3=公开,4=私密
	Channel      int         `orm:"channel"       description:"投放渠道:1=官网,2=小程序,3=短信,4=线下"`               // 投放渠道:1=官网,2=小程序,3=短信,4=线下
	BudgetAmount int         `orm:"budget_amount" description:"预算金额（分）"`                                 // 预算金额（分）
	LandingUrl   string      `orm:"landing_url"   description:"落地页URL"`                                  // 落地页URL
	RuleJson     string      `orm:"rule_json"     description:"规则JSON"`                                  // 规则JSON
	IntroContent string      `orm:"intro_content" description:"活动介绍"`                                    // 活动介绍
	StartAt      *gtime.Time `orm:"start_at"      description:"开始时间"`                                    // 开始时间
	EndAt        *gtime.Time `orm:"end_at"        description:"结束时间"`                                    // 结束时间
	IsPublic     int         `orm:"is_public"     description:"是否公开:0=否,1=是"`                            // 是否公开:0=否,1=是
	Status       int         `orm:"status"        description:"状态:0=草稿,1=已发布,2=已下架"`                     // 状态:0=草稿,1=已发布,2=已下架
	TenantId     uint64      `orm:"tenant_id"     description:"租户"`                                      // 租户
	MerchantId   uint64      `orm:"merchant_id"   description:"商户"`                                      // 商户
	CreatedBy    uint64      `orm:"created_by"    description:"创建人ID"`                                   // 创建人ID
	DeptId       uint64      `orm:"dept_id"       description:"所属部门ID"`                                  // 所属部门ID
	CreatedAt    *gtime.Time `orm:"created_at"    description:"创建时间"`                                    // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"    description:"更新时间"`                                    // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"    description:"软删除时间，非 NULL 表示已删除"`                      // 软删除时间，非 NULL 表示已删除
}
