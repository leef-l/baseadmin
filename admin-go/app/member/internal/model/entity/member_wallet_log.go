// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWalletLog is the golang structure for table member_wallet_log.
type MemberWalletLog struct {
	Id             uint64      `orm:"id"               description:"ID（Snowflake）"`                                             // ID（Snowflake）
	UserId         uint64      `orm:"user_id"          description:"会员|ref:member_user.nickname|search:select"`                 // 会员|ref:member_user.nickname|search:select
	WalletType     int         `orm:"wallet_type"      description:"钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select"`                 // 钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select
	ChangeType     int         `orm:"change_type"      description:"变动类型:1=充值,2=消费,3=推广奖,4=仓库卖出收入,5=平台扣除,6=后台调整|search:select"` // 变动类型:1=充值,2=消费,3=推广奖,4=仓库卖出收入,5=平台扣除,6=后台调整|search:select
	ChangeAmount   int64       `orm:"change_amount"    description:"变动金额（分，正增负减）"`                                              // 变动金额（分，正增负减）
	BeforeBalance  int64       `orm:"before_balance"   description:"变动前余额（分）"`                                                  // 变动前余额（分）
	AfterBalance   int64       `orm:"after_balance"    description:"变动后余额（分）"`                                                  // 变动后余额（分）
	RelatedOrderNo string      `orm:"related_order_no" description:"关联单号|search:eq|keyword:on"`                                 // 关联单号|search:eq|keyword:on
	Remark         string      `orm:"remark"           description:"备注说明|search:off"`                                           // 备注说明|search:off
	TenantId       uint64      `orm:"tenant_id"        description:"租户"`                                                        // 租户
	MerchantId     uint64      `orm:"merchant_id"      description:"商户"`                                                        // 商户
	CreatedBy      uint64      `orm:"created_by"       description:"创建人ID"`                                                     // 创建人ID
	DeptId         uint64      `orm:"dept_id"          description:"所属部门ID"`                                                    // 所属部门ID
	CreatedAt      *gtime.Time `orm:"created_at"       description:"创建时间"`                                                      // 创建时间
	UpdatedAt      *gtime.Time `orm:"updated_at"       description:"更新时间"`                                                      // 更新时间
	DeletedAt      *gtime.Time `orm:"deleted_at"       description:"软删除时间，非 NULL 表示已删除"`                                        // 软删除时间，非 NULL 表示已删除
}
