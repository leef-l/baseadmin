// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberShopOrder is the golang structure of table member_shop_order for DAO operations like Where/Data.
type MemberShopOrder struct {
	g.Meta      `orm:"table:member_shop_order, do:true"`
	Id          any         // ID（Snowflake）
	OrderNo     any         // 订单号|search:eq|keyword:on|priority:100
	UserId      any         // 购买会员|ref:member_user.nickname|search:select
	GoodsId     any         // 商品|ref:member_shop_goods.title|search:select
	GoodsTitle  any         // 商品名称（快照）
	GoodsCover  any         // 商品封面（快照）
	Quantity    any         // 购买数量
	TotalPrice  any         // 订单总价（分）
	PayWallet   any         // 支付钱包:1=优惠券余额
	OrderStatus any         // 订单状态:1=已完成,2=已取消|search:select
	Remark      any         // 订单备注|search:off
	Status      any         // 状态:0=关闭,1=开启|search:select
	TenantId    any         // 租户
	MerchantId  any         // 商户
	CreatedBy   any         // 创建人ID
	DeptId      any         // 所属部门ID
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 软删除时间，非 NULL 表示已删除
}
