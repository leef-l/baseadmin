// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberRebindLog is the golang structure for table member_rebind_log.
type MemberRebindLog struct {
	Id          uint64      `orm:"id"            description:"ID（Snowflake）"`                             // ID（Snowflake）
	UserId      uint64      `orm:"user_id"       description:"会员|ref:member_user.nickname|search:select"` // 会员|ref:member_user.nickname|search:select
	OldParentId uint64      `orm:"old_parent_id" description:"原上级|ref:member_user.nickname"`              // 原上级|ref:member_user.nickname
	NewParentId uint64      `orm:"new_parent_id" description:"新上级|ref:member_user.nickname"`              // 新上级|ref:member_user.nickname
	Reason      string      `orm:"reason"        description:"换绑原因|search:off"`                           // 换绑原因|search:off
	OperatorId  uint64      `orm:"operator_id"   description:"操作人|ref:system_users.username"`             // 操作人|ref:system_users.username
	TenantId    uint64      `orm:"tenant_id"     description:"租户"`                                        // 租户
	MerchantId  uint64      `orm:"merchant_id"   description:"商户"`                                        // 商户
	CreatedBy   uint64      `orm:"created_by"    description:"创建人ID"`                                     // 创建人ID
	DeptId      uint64      `orm:"dept_id"       description:"所属部门ID"`                                    // 所属部门ID
	CreatedAt   *gtime.Time `orm:"created_at"    description:"创建时间"`                                      // 创建时间
	UpdatedAt   *gtime.Time `orm:"updated_at"    description:"更新时间"`                                      // 更新时间
	DeletedAt   *gtime.Time `orm:"deleted_at"    description:"软删除时间，非 NULL 表示已删除"`                        // 软删除时间，非 NULL 表示已删除
}
