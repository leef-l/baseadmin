// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberWallet is the golang structure of table member_wallet for DAO operations like Where/Data.
type MemberWallet struct {
	g.Meta       `orm:"table:member_wallet, do:true"`
	Id           any         // ID（Snowflake）
	UserId       any         // 会员|ref:member_user.nickname|search:select
	WalletType   any         // 钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select
	Balance      any         // 当前余额（分）
	TotalIncome  any         // 累计收入（分）
	TotalExpense any         // 累计支出（分）
	FrozenAmount any         // 冻结金额（分）
	Status       any         // 状态:0=冻结,1=正常|search:select
	TenantId     any         // 租户
	MerchantId   any         // 商户
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间，非 NULL 表示已删除
}
