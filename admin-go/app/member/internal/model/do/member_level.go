// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberLevel is the golang structure of table member_level for DAO operations like Where/Data.
type MemberLevel struct {
	g.Meta             `orm:"table:member_level, do:true"`
	Id                 any         // 等级ID（Snowflake）
	Name               any         // 等级名称|search:like|keyword:on|priority:100
	LevelNo            any         // 等级编号（越大越高）|search:eq
	Icon               any         // 等级图标
	DurationDays       any         // 有效天数（0=永久）
	NeedActiveCount    any         // 升级所需有效用户数
	NeedTeamTurnover   any         // 升级所需团队营业额（分）
	DailyPurchaseLimit any         // 该等级每日限购单数|search:eq
	IsTop              any         // 是否最高等级:0=否,1=是|search:select
	AutoDeploy         any         // 到达后自动部署站点:0=否,1=是
	Remark             any         // 等级说明|search:off
	Sort               any         // 排序（升序）
	Status             any         // 状态:0=关闭,1=开启|search:select
	TenantId           any         // 租户
	MerchantId         any         // 商户
	CreatedBy          any         // 创建人ID
	DeptId             any         // 所属部门ID
	CreatedAt          *gtime.Time // 创建时间
	UpdatedAt          *gtime.Time // 更新时间
	DeletedAt          *gtime.Time // 软删除时间，非 NULL 表示已删除
}
