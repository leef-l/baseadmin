// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoSurvey is the golang structure of table demo_survey for DAO operations like Where/Data.
type DemoSurvey struct {
	g.Meta       `orm:"table:demo_survey, do:true"`
	Id           any         // 问卷ID（Snowflake）
	SurveyNo     any         // 问卷编号|search:eq|priority:100
	Title        any         // 问卷标题|search:like|keyword:on|priority:95
	Poster       any         // 海报
	QuestionJson any         // 问题JSON
	IntroContent any         // 问卷介绍
	PublishAt    *gtime.Time // 发布时间
	ExpireAt     *gtime.Time // 过期时间
	IsAnonymous  any         // 是否匿名:0=否,1=是
	Status       any         // 状态:0=草稿,1=已发布,2=已下架
	TenantId     any         // 租户
	MerchantId   any         // 商户
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间，非 NULL 表示已删除
}
