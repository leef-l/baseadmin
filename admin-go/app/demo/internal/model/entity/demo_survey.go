// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoSurvey is the golang structure for table demo_survey.
type DemoSurvey struct {
	Id           uint64      `orm:"id"            description:"问卷ID（Snowflake）"`                         // 问卷ID（Snowflake）
	SurveyNo     string      `orm:"survey_no"     description:"问卷编号|search:eq|priority:100"`             // 问卷编号|search:eq|priority:100
	Title        string      `orm:"title"         description:"问卷标题|search:like|keyword:on|priority:95"` // 问卷标题|search:like|keyword:on|priority:95
	Poster       string      `orm:"poster"        description:"海报"`                                      // 海报
	QuestionJson string      `orm:"question_json" description:"问题JSON"`                                  // 问题JSON
	IntroContent string      `orm:"intro_content" description:"问卷介绍"`                                    // 问卷介绍
	PublishAt    *gtime.Time `orm:"publish_at"    description:"发布时间"`                                    // 发布时间
	ExpireAt     *gtime.Time `orm:"expire_at"     description:"过期时间"`                                    // 过期时间
	IsAnonymous  int         `orm:"is_anonymous"  description:"是否匿名:0=否,1=是"`                            // 是否匿名:0=否,1=是
	Status       int         `orm:"status"        description:"状态:0=草稿,1=已发布,2=已下架"`                     // 状态:0=草稿,1=已发布,2=已下架
	TenantId     uint64      `orm:"tenant_id"     description:"租户"`                                      // 租户
	MerchantId   uint64      `orm:"merchant_id"   description:"商户"`                                      // 商户
	CreatedBy    uint64      `orm:"created_by"    description:"创建人ID"`                                   // 创建人ID
	DeptId       uint64      `orm:"dept_id"       description:"所属部门ID"`                                  // 所属部门ID
	CreatedAt    *gtime.Time `orm:"created_at"    description:"创建时间"`                                    // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"    description:"更新时间"`                                    // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"    description:"软删除时间，非 NULL 表示已删除"`                      // 软删除时间，非 NULL 表示已删除
}
