// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWallet is the golang structure for table member_wallet.
type MemberWallet struct {
	Id           uint64      `orm:"id"            description:"ID（Snowflake）"`                             // ID（Snowflake）
	UserId       uint64      `orm:"user_id"       description:"会员|ref:member_user.nickname|search:select"` // 会员|ref:member_user.nickname|search:select
	WalletType   int         `orm:"wallet_type"   description:"钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select"` // 钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select
	Balance      int64       `orm:"balance"       description:"当前余额（分）"`                                   // 当前余额（分）
	TotalIncome  uint64      `orm:"total_income"  description:"累计收入（分）"`                                   // 累计收入（分）
	TotalExpense uint64      `orm:"total_expense" description:"累计支出（分）"`                                   // 累计支出（分）
	FrozenAmount uint64      `orm:"frozen_amount" description:"冻结金额（分）"`                                   // 冻结金额（分）
	Status       int         `orm:"status"        description:"状态:0=冻结,1=正常|search:select"`                // 状态:0=冻结,1=正常|search:select
	TenantId     uint64      `orm:"tenant_id"     description:"租户"`                                        // 租户
	MerchantId   uint64      `orm:"merchant_id"   description:"商户"`                                        // 商户
	CreatedBy    uint64      `orm:"created_by"    description:"创建人ID"`                                     // 创建人ID
	DeptId       uint64      `orm:"dept_id"       description:"所属部门ID"`                                    // 所属部门ID
	CreatedAt    *gtime.Time `orm:"created_at"    description:"创建时间"`                                      // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"    description:"更新时间"`                                      // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"    description:"软删除时间，非 NULL 表示已删除"`                        // 软删除时间，非 NULL 表示已删除
}
