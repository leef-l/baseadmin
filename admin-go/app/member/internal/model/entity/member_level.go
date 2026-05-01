// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberLevel is the golang structure for table member_level.
type MemberLevel struct {
	Id                 uint64      `orm:"id"                   description:"等级ID（Snowflake）"`                          // 等级ID（Snowflake）
	Name               string      `orm:"name"                 description:"等级名称|search:like|keyword:on|priority:100"` // 等级名称|search:like|keyword:on|priority:100
	LevelNo            uint        `orm:"level_no"             description:"等级编号（越大越高）|search:eq"`                     // 等级编号（越大越高）|search:eq
	Icon               string      `orm:"icon"                 description:"等级图标"`                                     // 等级图标
	DurationDays       uint        `orm:"duration_days"        description:"有效天数（0=永久）"`                               // 有效天数（0=永久）
	NeedActiveCount    uint        `orm:"need_active_count"    description:"升级所需有效用户数"`                                // 升级所需有效用户数
	NeedTeamTurnover   uint64      `orm:"need_team_turnover"   description:"升级所需团队营业额（分）"`                             // 升级所需团队营业额（分）
	DailyPurchaseLimit uint        `orm:"daily_purchase_limit" description:"该等级每日限购单数|search:eq"`                      // 该等级每日限购单数|search:eq
	IsTop              int         `orm:"is_top"               description:"是否最高等级:0=否,1=是|search:select"`             // 是否最高等级:0=否,1=是|search:select
	AutoDeploy         int         `orm:"auto_deploy"          description:"到达后自动部署站点:0=否,1=是"`                        // 到达后自动部署站点:0=否,1=是
	Remark             string      `orm:"remark"               description:"等级说明|search:off"`                          // 等级说明|search:off
	Sort               int         `orm:"sort"                 description:"排序（升序）"`                                   // 排序（升序）
	Status             int         `orm:"status"               description:"状态:0=关闭,1=开启|search:select"`               // 状态:0=关闭,1=开启|search:select
	TenantId           uint64      `orm:"tenant_id"            description:"租户"`                                       // 租户
	MerchantId         uint64      `orm:"merchant_id"          description:"商户"`                                       // 商户
	CreatedBy          uint64      `orm:"created_by"           description:"创建人ID"`                                    // 创建人ID
	DeptId             uint64      `orm:"dept_id"              description:"所属部门ID"`                                   // 所属部门ID
	CreatedAt          *gtime.Time `orm:"created_at"           description:"创建时间"`                                     // 创建时间
	UpdatedAt          *gtime.Time `orm:"updated_at"           description:"更新时间"`                                     // 更新时间
	DeletedAt          *gtime.Time `orm:"deleted_at"           description:"软删除时间，非 NULL 表示已删除"`                       // 软删除时间，非 NULL 表示已删除
}
