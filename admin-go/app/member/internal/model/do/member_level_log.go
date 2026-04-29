// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberLevelLog is the golang structure of table member_level_log for DAO operations like Where/Data.
type MemberLevelLog struct {
	g.Meta     `orm:"table:member_level_log, do:true"`
	Id         any         // ID（Snowflake）
	UserId     any         // 会员|ref:member_user.nickname|search:select
	OldLevelId any         // 变更前等级|ref:member_level.name
	NewLevelId any         // 变更后等级|ref:member_level.name
	ChangeType any         // 变更类型:1=自动升级,2=后台调整,3=过期降级|search:select
	ExpireAt   *gtime.Time // 新等级到期时间
	Remark     any         // 变更说明|search:off
	TenantId   any         // 租户
	MerchantId any         // 商户
	CreatedBy  any         // 创建人ID
	DeptId     any         // 所属部门ID
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 软删除时间，非 NULL 表示已删除
}
