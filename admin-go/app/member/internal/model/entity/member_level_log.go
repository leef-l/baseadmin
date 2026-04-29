// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberLevelLog is the golang structure for table member_level_log.
type MemberLevelLog struct {
	Id         uint64      `orm:"id"           description:"ID（Snowflake）"`                             // ID（Snowflake）
	UserId     uint64      `orm:"user_id"      description:"会员|ref:member_user.nickname|search:select"` // 会员|ref:member_user.nickname|search:select
	OldLevelId uint64      `orm:"old_level_id" description:"变更前等级|ref:member_level.name"`               // 变更前等级|ref:member_level.name
	NewLevelId uint64      `orm:"new_level_id" description:"变更后等级|ref:member_level.name"`               // 变更后等级|ref:member_level.name
	ChangeType int         `orm:"change_type"  description:"变更类型:1=自动升级,2=后台调整,3=过期降级|search:select"`   // 变更类型:1=自动升级,2=后台调整,3=过期降级|search:select
	ExpireAt   *gtime.Time `orm:"expire_at"    description:"新等级到期时间"`                                   // 新等级到期时间
	Remark     string      `orm:"remark"       description:"变更说明|search:off"`                           // 变更说明|search:off
	TenantId   uint64      `orm:"tenant_id"    description:"租户"`                                        // 租户
	MerchantId uint64      `orm:"merchant_id"  description:"商户"`                                        // 商户
	CreatedBy  uint64      `orm:"created_by"   description:"创建人ID"`                                     // 创建人ID
	DeptId     uint64      `orm:"dept_id"      description:"所属部门ID"`                                    // 所属部门ID
	CreatedAt  *gtime.Time `orm:"created_at"   description:"创建时间"`                                      // 创建时间
	UpdatedAt  *gtime.Time `orm:"updated_at"   description:"更新时间"`                                      // 更新时间
	DeletedAt  *gtime.Time `orm:"deleted_at"   description:"软删除时间，非 NULL 表示已删除"`                        // 软删除时间，非 NULL 表示已删除
}
