// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberRebindLog is the golang structure of table member_rebind_log for DAO operations like Where/Data.
type MemberRebindLog struct {
	g.Meta      `orm:"table:member_rebind_log, do:true"`
	Id          any         // ID（Snowflake）
	UserId      any         // 会员|ref:member_user.nickname|search:select
	OldParentId any         // 原上级|ref:member_user.nickname
	NewParentId any         // 新上级|ref:member_user.nickname
	Reason      any         // 换绑原因|search:off
	OperatorId  any         // 操作人|ref:system_users.username
	TenantId    any         // 租户
	MerchantId  any         // 商户
	CreatedBy   any         // 创建人ID
	DeptId      any         // 所属部门ID
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 软删除时间，非 NULL 表示已删除
}
