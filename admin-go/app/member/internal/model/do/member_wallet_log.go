// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWalletLog is the golang structure of table member_wallet_log for DAO operations like Where/Data.
type MemberWalletLog struct {
	g.Meta         `orm:"table:member_wallet_log, do:true"`
	Id             any         // ID（Snowflake）
	UserId         any         // 会员|ref:member_user.nickname|search:select
	WalletType     any         // 钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select
	ChangeType     any         // 变动类型:1=充值,2=消费,3=推广奖,4=仓库卖出收入,5=平台扣除,6=后台调整|search:select
	ChangeAmount   any         // 变动金额（分，正增负减）
	BeforeBalance  any         // 变动前余额（分）
	AfterBalance   any         // 变动后余额（分）
	RelatedOrderNo any         // 关联单号|search:eq|keyword:on
	Remark         any         // 备注说明|search:off
	TenantId       any         // 租户
	MerchantId     any         // 商户
	CreatedBy      any         // 创建人ID
	DeptId         any         // 所属部门ID
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
	DeletedAt      *gtime.Time // 软删除时间，非 NULL 表示已删除
}
